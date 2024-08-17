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
		}

		// Set discrete inputs and their legal addresses
		for _, input := range slaveConfig.DiscreteInputs {
			slave.DiscreteInputs[input.Address] = input.Value
		}

		// Set holding registers and their legal addresses
		for _, register := range slaveConfig.HoldingRegisters {
			slave.HoldingRegisters[register.Address] = register.Value
		}

		// Set input registers and their legal addresses
		for _, register := range slaveConfig.InputRegisters {
			slave.InputRegisters[register.Address] = register.Value
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
	s.Slaves[unitID] = &Slave{
		Coils:            make(map[uint16]bool),
		DiscreteInputs:   make(map[uint16]bool),
		HoldingRegisters: make(map[uint16]uint16),
		InputRegisters:   make(map[uint16]uint16),
	}
}

// Start starts the server
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

		log.Printf("Sending Response: %x", response)
		_, err = port.Write(response)
		if err != nil {
			log.Printf("Failed to write response: %v", err)
			continue
		}
	}
}

// getSlave gets the slave based on the unit ID
func (s *SerialServer) getSlave(unitID byte) (*Slave, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slave, exists := s.Slaves[unitID]
	if !exists {
		return nil, fmt.Errorf("slave %d not found", unitID)
	}
	return slave, nil
}

// processRequest processes the request
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
	var responseADU []byte
	switch serialADU.PDU[0] { // Function code is the first byte in PDU
	case 0x01:
		responseADU, err = s.handleReadCoils(slave, serialADU)
		if err != nil {
			log.Printf("Failed to handle Read Coils: %v", err)
		}
	case 0x02:
		responseADU, err = s.handleReadDiscreteInputs(slave, serialADU)
		if err != nil {
			log.Printf("Failed to handle Read Discrete Inputs: %v", err)
		}
	case 0x03:
		responseADU, err = s.hanldeReadHoldingRegisters(slave, serialADU)
		if err != nil {
			log.Printf("Failed to handle Read Holding Registers: %v", err)
		}
	case 0x04:
		responseADU, err = s.handleReadInputRegisters(slave, serialADU)
		if err != nil {
			log.Printf("Failed to handle Read Input Registers: %v", err)
		}
	case 0x05:
		responseADU, err = s.handleWriteSingleCoil(slave, serialADU)
		if err != nil {
			log.Printf("Failed to handle Write Single Coil: %v", err)
		}
	case 0x06:
		responseADU, err = s.handleWriteSingleRegister(slave, serialADU)
		if err != nil {
			log.Printf("Failed to handle Write Single Register: %v", err)
		}
	case 0x0F:
		responseADU, err = s.handleWriteMultipleCoils(slave, serialADU)
		if err != nil {
			log.Printf("Failed to handle Write Multiple Coils: %v", err)
		}
	case 0x10:
		responseADU, err = s.handleWriteMultipleRegisters(slave, serialADU)
		if err != nil {
			log.Printf("Failed to handle Write Multiple Registers: %v", err)
		}
	default:
		return s.exceptionResponse(request, 0x01), nil // Illegal Function
	}
	// Create the response ADU using the response PDU
	return responseADU, nil
}

// handleReadCoils handles the Read Coils function
func (s *SerialServer) handleReadCoils(slave *Slave, requestADU *adu.SerialADU) ([]byte, error) {
	pduRead := &pdu.PDU_Read{}
	pduRead.FromBytes(requestADU.PDU)

	startAddress := pduRead.StartingAddress
	quantity := pduRead.Quantity

	if startAddress+quantity > 65535 {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if _, exists := slave.Coils[startAddress+uint16(i)]; !exists {
			return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
		}
	}

	byteCount := (quantity + 7) / 8
	responseData := make([]byte, byteCount)
	for i := 0; i < int(quantity); i++ {
		if slave.Coils[startAddress+uint16(i)] {
			responseData[i/8] |= 1 << (i % 8)
		}
	}

	responsePDU := pdu.NewPDUReadResponse(0x01, responseData)
	return responsePDU.ToBytes(), nil
}

// handleReadDiscreteInputs handles the Read Discrete Inputs function
func (s *SerialServer) handleReadDiscreteInputs(slave *Slave, requestADU *adu.SerialADU) ([]byte, error) {
	pduRead := &pdu.PDU_Read{}
	pduRead.FromBytes(requestADU.PDU)

	startAddress := pduRead.StartingAddress
	quantity := pduRead.Quantity

	if startAddress+quantity > 65535 {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if _, exists := slave.DiscreteInputs[startAddress+uint16(i)]; !exists {
			return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
		}
	}

	byteCount := (quantity + 7) / 8
	responseData := make([]byte, byteCount)

	for i := 0; i < int(quantity); i++ {
		if slave.DiscreteInputs[startAddress+uint16(i)] {
			responseData[i/8] |= 1 << (i % 8)
		}
	}

	responsePDU := pdu.NewPDUReadResponse(0x02, responseData)
	return responsePDU.ToBytes(), nil
}

// hanldeReadHoldingRegisters handles the Read Holding Registers function
func (s *SerialServer) hanldeReadHoldingRegisters(slave *Slave, requestADU *adu.SerialADU) ([]byte, error) {
	pduRead := &pdu.PDU_Read{}
	pduRead.FromBytes(requestADU.PDU)

	startAddress := pduRead.StartingAddress
	quantity := pduRead.Quantity

	if startAddress+quantity > 65535 {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if _, exists := slave.HoldingRegisters[startAddress+uint16(i)]; !exists {
			return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
		}
	}

	responseData := make([]byte, 2*quantity)
	for i := 0; i < int(quantity); i++ {
		value := slave.HoldingRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(responseData[2*i:], value)
	}

	responsePDU := pdu.NewPDUReadResponse(0x03, responseData)
	return responsePDU.ToBytes(), nil
}

// handleReadInputRegisters handles the Read Input Registers function
func (s *SerialServer) handleReadInputRegisters(slave *Slave, requestADU *adu.SerialADU) ([]byte, error) {
	pduRead := &pdu.PDU_Read{}
	pduRead.FromBytes(requestADU.PDU)

	startAddress := pduRead.StartingAddress
	quantity := pduRead.Quantity

	if startAddress+quantity > 65535 {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if _, exists := slave.InputRegisters[startAddress+uint16(i)]; !exists {
			return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
		}
	}

	responseData := make([]byte, 2*quantity)
	for i := 0; i < int(quantity); i++ {
		value := slave.InputRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(responseData[2*i:], value)
	}

	responsePDU := pdu.NewPDUReadResponse(0x04, responseData)
	return responsePDU.ToBytes(), nil
}

// handleWriteSingleCoil handles the Write Single Coil function
func (s *SerialServer) handleWriteSingleCoil(slave *Slave, requestADU *adu.SerialADU) ([]byte, error) {
	pduWrite := &pdu.PDU_WriteSingleCoil{}
	pduWrite.FromBytes(requestADU.PDU)

	address := pduWrite.OutputAddress
	value := pduWrite.OutputValue

	if address >= 65535 {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	if _, exists := slave.Coils[address]; !exists {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	if value == 0xFF00 {
		slave.Coils[address] = true
	} else if value == 0x0000 {
		slave.Coils[address] = false
	} else {
		return s.exceptionResponse(requestADU.ToBytes(), 0x03), nil // Illegal Data Value
	}

	return pduWrite.ToBytes(), nil
}

// handleWriteMultipleCoils handles the Write Multiple Coils function
func (s *SerialServer) handleWriteMultipleCoils(slave *Slave, requestADU *adu.SerialADU) ([]byte, error) {
	pduWrite := &pdu.PDU_WriteMultipleCoils{}
	pduWrite.FromBytes(requestADU.PDU)

	startAddress := pduWrite.StartingAddress
	quantity := pduWrite.QuantityOfOutputs
	byteCount := pduWrite.ByteCount

	if int(byteCount) != len(pduWrite.OutputValues) || startAddress+quantity > 65535 {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if _, exists := slave.Coils[startAddress+uint16(i)]; !exists {
			return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
		}
	}

	for i := 0; i < int(quantity); i++ {
		if pduWrite.OutputValues[i/8]&(1<<(i%8)) != 0 {
			slave.Coils[startAddress+uint16(i)] = true
		} else {
			slave.Coils[startAddress+uint16(i)] = false
		}
	}

	responsePDU := pdu.NewPDUWriteMultipleResponse(0x0F, startAddress, quantity)
	return responsePDU.ToBytes(), nil
}

// handleWriteSingleRegister handles the Write Single Register function
func (s *SerialServer) handleWriteSingleRegister(slave *Slave, requestADU *adu.SerialADU) ([]byte, error) {
	pduWrite := &pdu.PDU_WriteSingleRegister{}
	pduWrite.FromBytes(requestADU.PDU)

	address := pduWrite.RegisterAddress
	value := pduWrite.RegisterValue

	if address >= 65535 {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	if _, exists := slave.HoldingRegisters[address]; !exists {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	slave.HoldingRegisters[address] = value

	return pduWrite.ToBytes(), nil
}

// handleWriteMultipleRegisters handles the Write Multiple Registers function
func (s *SerialServer) handleWriteMultipleRegisters(slave *Slave, requestADU *adu.SerialADU) ([]byte, error) {
	pduWrite := &pdu.PDU_WriteMultipleRegisters{}
	pduWrite.FromBytes(requestADU.PDU)

	startAddress := pduWrite.StartingAddress
	quantity := pduWrite.QuantityOfOutputs
	byteCount := pduWrite.ByteCount

	if int(byteCount) != len(pduWrite.OutputValues) || startAddress+quantity > 65535 {
		return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if _, exists := slave.HoldingRegisters[startAddress+uint16(i)]; !exists {
			return s.exceptionResponse(requestADU.ToBytes(), 0x02), nil // Illegal Data Address
		}
	}

	for i := 0; i < int(quantity); i++ {
		value := binary.BigEndian.Uint16(pduWrite.OutputValues[2*i:])
		slave.HoldingRegisters[startAddress+uint16(i)] = value
	}

	responsePDU := pdu.NewPDUWriteMultipleResponse(0x10, startAddress, quantity)
	return responsePDU.ToBytes(), nil
}

// exceptionResponse creates an exception response
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
