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
	"github.com/ulfaric/gomodbus/pdu"
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
	// Parse the request as a Serial ADU
	serialADU := &adu.SerialADU{}
	err := serialADU.FromBytes(request)
	if err != nil {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}

	// Fetch the slave based on Unit ID
	slave, err := s.getSlave(serialADU.Address)
	if err != nil {
		return s.exceptionResponse(request, 0x04), nil // Server Device Failure
	}

	// Process the request based on the function code
	var responsePDU []byte
	switch serialADU.PDU[0] { // Function code is the first byte in PDU
	case 0x01:
		responsePDU, err = s.readCoils(slave, serialADU.PDU)
	case 0x02:
		responsePDU, err = s.readDiscreteInputs(slave, serialADU.PDU)
	case 0x03:
		responsePDU, err = s.readHoldingRegisters(slave, serialADU.PDU)
	case 0x04:
		responsePDU, err = s.readInputRegisters(slave, serialADU.PDU)
	case 0x05:
		responsePDU, err = s.writeSingleCoil(slave, serialADU.PDU)
	case 0x06:
		responsePDU, err = s.writeSingleRegister(slave, serialADU.PDU)
	case 0x0F:
		responsePDU, err = s.writeMultipleCoils(slave, serialADU.PDU)
	case 0x10:
		responsePDU, err = s.writeMultipleRegisters(slave, serialADU.PDU)
	default:
		return s.exceptionResponse(request, 0x01), nil // Illegal Function
	}

	if err != nil {
		return s.exceptionResponse(request, 0x04), err // Server Device Failure
	}

	// Create the response ADU using the response PDU
	responseADU := adu.NewSerialADU(serialADU.Address, responsePDU)
	return responseADU.ToBytes(), nil
}

func (s *SerialServer) readCoils(slave *Slave, pduBytes []byte) ([]byte, error) {
	pduRead := &pdu.PDU_Read{}
	pduRead.FromBytes(pduBytes)

	startAddress := pduRead.StartingAddress
	quantity := pduRead.Quantity

	if startAddress+quantity > 65535 || !slave.LegalCoilsAddress[startAddress] {
		return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
	}

	byteCount := (quantity + 7) / 8
	responseData := make([]byte, byteCount)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalCoilsAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
		}
		if slave.Coils[startAddress+uint16(i)] {
			responseData[i/8] |= 1 << (i % 8)
		}
	}

	responsePDU := pdu.NewPDUReadResponse(0x01, responseData)
	return responsePDU.ToBytes(), nil
}

func (s *SerialServer) readDiscreteInputs(slave *Slave, pduBytes []byte) ([]byte, error) {
	pduRead := &pdu.PDU_Read{}
	pduRead.FromBytes(pduBytes)

	startAddress := pduRead.StartingAddress
	quantity := pduRead.Quantity

	if startAddress+quantity > 65535 || !slave.LegalDiscreteInputsAddress[startAddress] {
		return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
	}

	byteCount := (quantity + 7) / 8
	responseData := make([]byte, byteCount)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalDiscreteInputsAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
		}
		if slave.DiscreteInputs[startAddress+uint16(i)] {
			responseData[i/8] |= 1 << (i % 8)
		}
	}

	responsePDU := pdu.NewPDUReadResponse(0x02, responseData)
	return responsePDU.ToBytes(), nil
}

func (s *SerialServer) readHoldingRegisters(slave *Slave, pduBytes []byte) ([]byte, error) {
	pduRead := &pdu.PDU_Read{}
	pduRead.FromBytes(pduBytes)

	startAddress := pduRead.StartingAddress
	quantity := pduRead.Quantity

	if startAddress+quantity > 65535 || !slave.LegalHoldingRegistersAddress[startAddress] {
		return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
	}

	responseData := make([]byte, 2*quantity)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
		}
		value := slave.HoldingRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(responseData[2*i:], value)
	}

	responsePDU := pdu.NewPDUReadResponse(0x03, responseData)
	return responsePDU.ToBytes(), nil
}

func (s *SerialServer) readInputRegisters(slave *Slave, pduBytes []byte) ([]byte, error) {
	pduRead := &pdu.PDU_Read{}
	pduRead.FromBytes(pduBytes)

	startAddress := pduRead.StartingAddress
	quantity := pduRead.Quantity

	if startAddress+quantity > 65535 || !slave.LegalInputRegistersAddress[startAddress] {
		return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
	}

	responseData := make([]byte, 2*quantity)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalInputRegistersAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
		}
		value := slave.InputRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(responseData[2*i:], value)
	}

	responsePDU := pdu.NewPDUReadResponse(0x04, responseData)
	return responsePDU.ToBytes(), nil
}

func (s *SerialServer) writeSingleCoil(slave *Slave, pduBytes []byte) ([]byte, error) {
	pduWrite := &pdu.PDU_WriteSingleCoil{}
	pduWrite.FromBytes(pduBytes)

	address := pduWrite.OutputAddress
	value := pduWrite.OutputValue

	if address >= 65535 || !slave.LegalCoilsAddress[address] {
		return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
	}

	if value == 0xFF00 {
		slave.Coils[address] = true
	} else if value == 0x0000 {
		slave.Coils[address] = false
	} else {
		return s.exceptionResponse(pduBytes, 0x03), nil // Illegal Data Value
	}

	return pduWrite.ToBytes(), nil
}

func (s *SerialServer) writeMultipleCoils(slave *Slave, pduBytes []byte) ([]byte, error) {
	pduWrite := &pdu.PDU_WriteMultipleCoils{}
	pduWrite.FromBytes(pduBytes)

	startAddress := pduWrite.StartingAddress
	quantity := pduWrite.QuantityOfOutputs
	byteCount := pduWrite.ByteCount

	if int(byteCount) != len(pduWrite.OutputValues) || startAddress+quantity > 65535 || !slave.LegalCoilsAddress[startAddress] {
		return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalCoilsAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
		}
		if pduWrite.OutputValues[i/8]&(1<<(i%8)) != 0 {
			slave.Coils[startAddress+uint16(i)] = true
		} else {
			slave.Coils[startAddress+uint16(i)] = false
		}
	}

	responsePDU := pdu.NewPDUWriteMultipleResponse(0x0F, startAddress, quantity)
	return responsePDU.ToBytes(), nil
}

func (s *SerialServer) writeSingleRegister(slave *Slave, pduBytes []byte) ([]byte, error) {
	pduWrite := &pdu.PDU_WriteSingleRegister{}
	pduWrite.FromBytes(pduBytes)

	address := pduWrite.RegisterAddress
	value := pduWrite.RegisterValue

	if address >= 65535 || !slave.LegalHoldingRegistersAddress[address] {
		return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
	}

	slave.HoldingRegisters[address] = value

	return pduWrite.ToBytes(), nil
}

func (s *SerialServer) writeMultipleRegisters(slave *Slave, pduBytes []byte) ([]byte, error) {
	pduWrite := &pdu.PDU_WriteMultipleRegisters{}
	pduWrite.FromBytes(pduBytes)

	startAddress := pduWrite.StartingAddress
	quantity := pduWrite.QuantityOfOutputs
	byteCount := pduWrite.ByteCount

	if int(byteCount) != len(pduWrite.OutputValues) || startAddress+quantity > 65535 || !slave.LegalHoldingRegistersAddress[startAddress] {
		return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
		}
		value := binary.BigEndian.Uint16(pduWrite.OutputValues[2*i:])
		slave.HoldingRegisters[startAddress+uint16(i)] = value
	}

	responsePDU := pdu.NewPDUWriteMultipleResponse(0x10, startAddress, quantity)
	return responsePDU.ToBytes(), nil
}

func (s *SerialServer) exceptionResponse(request []byte, exceptionCode byte) []byte {
	serialADU := &adu.SerialADU{}
	err := serialADU.FromBytes(request)
	if err != nil {
		log.Printf("Failed to parse ADU for exception response: %v", err)
		return nil
	}

	functionCode := serialADU.PDU[0] | 0x80 // Set the highest bit to indicate an error
	exceptionPDU := []byte{functionCode, exceptionCode}

	responseADU := adu.NewSerialADU(serialADU.Address, exceptionPDU)
	return responseADU.ToBytes()
}
