package server

import (
	"encoding/binary"
	"fmt"

	"log"
	"os"
	"sync"

	"github.com/tarm/serial"
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
	"gopkg.in/yaml.v3"
)

type SerialServer struct {
	SerialConfig *serial.Config
	ByteOrder    string
	WordOrder    string
	Slaves       map[byte]*Slave
	mu           sync.Mutex
}

type Config struct {
	SerialPort string `yaml:"serialPort"`
	BaudRate   int    `yaml:"baudRate"`
	ByteOrder  string `yaml:"byteOrder"`
	WordOrder  string `yaml:"wordOrder"`
	Slaves     []struct {
		UnitID byte `yaml:"unitID"`
		Coils  []struct {
			Address uint16 `yaml:"address"`
			Value   bool   `yaml:"value"`
		} `yaml:"coils"`
		DiscreteInputs []struct {
			Address uint16 `yaml:"address"`
			Value   bool   `yaml:"value"`
		} `yaml:"discreteInputs"`
		HoldingRegisters []struct {
			Address uint16 `yaml:"address"`
			Value   uint16 `yaml:"value"`
		} `yaml:"holdingRegisters"`
		InputRegisters []struct {
			Address uint16 `yaml:"address"`
			Value   uint16 `yaml:"value"`
		} `yaml:"inputRegisters"`
	} `yaml:"slaves"`
}

func NewServer(configFile string) (*SerialServer, error) {
	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse config file
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	if config.ByteOrder != gomodbus.BigEndian && config.ByteOrder != gomodbus.LittleEndian {
		return nil, fmt.Errorf("invalid byte order: %s", config.ByteOrder)
	}
	if config.WordOrder != gomodbus.BigEndian && config.WordOrder != gomodbus.LittleEndian {
		return nil, fmt.Errorf("invalid word order: %s", config.WordOrder)
	}

	slaves := make(map[byte]*Slave)
	for _, slaveConfig := range config.Slaves {
		slave := &Slave{}

		// Set coils and their legal addresses
		for _, coil := range slaveConfig.Coils {
			slave.Coils[coil.Address] = coil.Value
			slave.LegalCoilsAddress[coil.Address] = true
		}

		// Set discrete inputs and their legal addresses
		for _, input := range slaveConfig.DiscreteInputs {
			slave.DiscreteInputs[input.Address] = input.Value
			slave.LegalDiscreteInputsAddress[input.Address] = true
		}

		// Set holding registers and their legal addresses
		for _, register := range slaveConfig.HoldingRegisters {
			slave.HoldingRegisters[register.Address] = register.Value
			slave.LegalHoldingRegistersAddress[register.Address] = true
		}

		// Set input registers and their legal addresses
		for _, register := range slaveConfig.InputRegisters {
			slave.InputRegisters[register.Address] = register.Value
			slave.LegalInputRegistersAddress[register.Address] = true
		}

		slaves[slaveConfig.UnitID] = slave
	}

	serialConfig := &serial.Config{
		Name: config.SerialPort,
		Baud: config.BaudRate,
	}

	return &SerialServer{
		SerialConfig: serialConfig,
		ByteOrder:    config.ByteOrder,
		WordOrder:    config.WordOrder,
		Slaves:       slaves,
	}, nil
}

func (s *SerialServer) AddSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Slaves[unitID] = &Slave{}
}

func (s *SerialServer) Start() error {
	port, err := serial.OpenPort(s.SerialConfig)
	if err != nil {
		return fmt.Errorf("failed to open serial port: %v", err)
	}
	defer port.Close()

	log.Printf("Modbus RTU server started on %s", s.SerialConfig.Name)
	buffer := make([]byte, 512)

	for {
		n, err := port.Read(buffer)
		if err != nil {
			log.Printf("Failed to read from serial port: %v", err)
			continue
		}

		request := buffer[:n]
		response, err := s.processRequest(request)
		if err != nil {
			log.Printf("Failed to process request: %v", err)
			response = s.exceptionResponse(request, 0x04) // Server Device Failure
		}

		_, err = port.Write(response)
		if err != nil {
			log.Printf("Failed to write response: %v", err)
			continue
		}
	}
}

func (s *SerialServer) getSlave(unitID byte) (*Slave, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slave, exists := s.Slaves[unitID]
	if !exists {
		return nil, fmt.Errorf("slave %d not found", unitID)
	}
	return slave, nil
}

func (s *SerialServer) processRequest(request []byte) ([]byte, error) {
	_adu := &adu.Serial_ADU{}
	err := _adu.FromBytes(request)
	if err != nil {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}

	slave, err := s.getSlave(_adu.Address)
	if err != nil {
		return s.exceptionResponse(request, 0x04), nil // Server Device Failure
	}

	var response []byte
	switch _adu.PDU[0] { // Function code is the first byte of the PDU
	case 0x01:
		response, err = s.readCoils(slave, _adu.PDU)
	case 0x02:
		response, err = s.readDiscreteInputs(slave, _adu.PDU)
	case 0x03:
		response, err = s.readHoldingRegisters(slave, _adu.PDU)
	case 0x04:
		response, err = s.readInputRegisters(slave, _adu.PDU)
	case 0x05:
		response, err = s.writeSingleCoil(slave, _adu.PDU)
	case 0x06:
		response, err = s.writeSingleRegister(slave, _adu.PDU)
	case 0x0F:
		response, err = s.writeMultipleCoils(slave, _adu.PDU)
	case 0x10:
		response, err = s.writeMultipleRegisters(slave, _adu.PDU)
	default:
		response = s.exceptionResponse(request, 0x01) // Illegal Function
	}

	if err != nil {
		return s.exceptionResponse(request, 0x04), nil // Server Device Failure
	}

	// Create the ADU for the response
	responseAdu := adu.New_Serial_ADU(_adu.Address, response)
	return responseAdu.ToBytes(), nil
}

func (s *SerialServer) readCoils(slave *Slave, pdu []byte) ([]byte, error) {
	if len(pdu) < 5 {
		return nil, fmt.Errorf("PDU too short")
	}
	startAddress := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])

	if startAddress+quantity > 65535 || !slave.LegalCoilsAddress[startAddress] {
		return nil, fmt.Errorf("illegal Data Address")
	}

	byteCount := (quantity + 7) / 8
	response := make([]byte, 2+byteCount)

	response[0] = pdu[0] // Function code
	response[1] = byte(byteCount)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalCoilsAddress[startAddress+uint16(i)] {
			return nil, fmt.Errorf("illegal Data Address")
		}
		if slave.Coils[startAddress+uint16(i)] {
			response[2+i/8] |= 1 << (i % 8)
		}
	}

	return response, nil
}

func (s *SerialServer) readDiscreteInputs(slave *Slave, pdu []byte) ([]byte, error) {
	if len(pdu) < 5 {
		return nil, fmt.Errorf("PDU too short")
	}
	startAddress := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])

	if startAddress+quantity > 65535 || !slave.LegalDiscreteInputsAddress[startAddress] {
		return nil, fmt.Errorf("illegal Data Address")
	}

	byteCount := (quantity + 7) / 8
	response := make([]byte, 2+byteCount)

	response[0] = pdu[0] // Function code
	response[1] = byte(byteCount)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalDiscreteInputsAddress[startAddress+uint16(i)] {
			return nil, fmt.Errorf("illegal Data Address")
		}
		if slave.DiscreteInputs[startAddress+uint16(i)] {
			response[2+i/8] |= 1 << (i % 8)
		}
	}

	return response, nil
}

func (s *SerialServer) readHoldingRegisters(slave *Slave, pdu []byte) ([]byte, error) {
	if len(pdu) < 5 {
		return nil, fmt.Errorf("PDU too short")
	}
	startAddress := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])

	if startAddress+quantity > 65535 || !slave.LegalHoldingRegistersAddress[startAddress] {
		return nil, fmt.Errorf("illegal Data Address")
	}

	response := make([]byte, 2+2*quantity)

	response[0] = pdu[0] // Function code
	response[1] = byte(2 * quantity)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
			return nil, fmt.Errorf("illegal Data Address")
		}
		value := slave.HoldingRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(response[2+2*i:], value)
	}

	return response, nil
}

func (s *SerialServer) readInputRegisters(slave *Slave, pdu []byte) ([]byte, error) {
	if len(pdu) < 5 {
		return nil, fmt.Errorf("PDU too short")
	}
	startAddress := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])

	if startAddress+quantity > 65535 || !slave.LegalInputRegistersAddress[startAddress] {
		return nil, fmt.Errorf("illegal Data Address")
	}

	response := make([]byte, 2+2*quantity)

	response[0] = pdu[0] // Function code
	response[1] = byte(2 * quantity)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalInputRegistersAddress[startAddress+uint16(i)] {
			return nil, fmt.Errorf("illegal Data Address")
		}
		value := slave.InputRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(response[2+2*i:], value)
	}

	return response, nil
}

func (s *SerialServer) writeSingleCoil(slave *Slave, pdu []byte) ([]byte, error) {
	if len(pdu) < 5 {
		return nil, fmt.Errorf("PDU too short")
	}
	address := binary.BigEndian.Uint16(pdu[1:3])
	value := binary.BigEndian.Uint16(pdu[3:5])

	if address >= 65535 || !slave.LegalCoilsAddress[address] {
		return nil, fmt.Errorf("illegal Data Address")
	}

	if value == 0xFF00 {
		slave.Coils[address] = true
	} else if value == 0x0000 {
		slave.Coils[address] = false
	} else {
		return nil, fmt.Errorf("illegal Data Value")
	}

	// Response is the echo of the request
	return pdu, nil
}

func (s *SerialServer) writeMultipleCoils(slave *Slave, pdu []byte) ([]byte, error) {
	if len(pdu) < 6 {
		return nil, fmt.Errorf("PDU too short")
	}
	startAddress := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])
	byteCount := pdu[5]

	if len(pdu) < int(6+byteCount) || startAddress+quantity > 65535 || !slave.LegalCoilsAddress[startAddress] {
		return nil, fmt.Errorf("illegal Data Address")
	}

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalCoilsAddress[startAddress+uint16(i)] {
			return nil, fmt.Errorf("illegal Data Address")
		}
		if pdu[6+i/8]&(1<<(i%8)) != 0 {
			slave.Coils[startAddress+uint16(i)] = true
		} else {
			slave.Coils[startAddress+uint16(i)] = false
		}
	}

	// Response is the starting address and quantity of coils written
	return pdu[:6], nil
}

func (s *SerialServer) writeSingleRegister(slave *Slave, pdu []byte) ([]byte, error) {
	if len(pdu) < 5 {
		return nil, fmt.Errorf("PDU too short")
	}
	address := binary.BigEndian.Uint16(pdu[1:3])
	value := binary.BigEndian.Uint16(pdu[3:5])

	if address >= 65535 || !slave.LegalHoldingRegistersAddress[address] {
		return nil, fmt.Errorf("illegal Data Address")
	}

	slave.HoldingRegisters[address] = value

	// Response is the echo of the request
	return pdu, nil
}

func (s *SerialServer) writeMultipleRegisters(slave *Slave, pdu []byte) ([]byte, error) {
	if len(pdu) < 6 {
		return nil, fmt.Errorf("PDU too short")
	}
	startAddress := binary.BigEndian.Uint16(pdu[1:3])
	quantity := binary.BigEndian.Uint16(pdu[3:5])
	byteCount := pdu[5]

	if len(pdu) < int(6+byteCount) || startAddress+quantity > 65535 || !slave.LegalHoldingRegistersAddress[startAddress] {
		return nil, fmt.Errorf("illegal Data Address")
	}

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
			return nil, fmt.Errorf("illegal Data Address")
		}
		value := binary.BigEndian.Uint16(pdu[6+2*i:])
		slave.HoldingRegisters[startAddress+uint16(i)] = value
	}

	// Response is the starting address and quantity of registers written
	return pdu[:6], nil
}

func (s *SerialServer) exceptionResponse(request []byte, exceptionCode byte) []byte {
	unitID := request[0]
	functionCode := request[1] | 0x80 // Set the highest bit to indicate an error

	response := []byte{unitID, functionCode, exceptionCode}

	// Create the ADU for the response
	responseAdu := adu.New_Serial_ADU(unitID, response)
	return responseAdu.ToBytes()
}
