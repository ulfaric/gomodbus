package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ulfaric/gomodbus"
	"gopkg.in/yaml.v3"
)

type TCPServer struct {
	Host      string
	Port      int
	ByteOrder string
	WordOrder string
	Slaves    map[byte]*Slave
	mu        sync.Mutex
}

type TCPConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	ByteOrder string `yaml:"byteOrder"`
	WordOrder string `yaml:"wordOrder"`
	Slaves    []struct {
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

func NewTCPServer(host, byteOrder, wordOrder string, port int) (*TCPServer, error) {
	// read config file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Failed to read config file: %v", err)
	}

	// parse config file
	var config TCPConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Failed to parse config file: %v", err)
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

	return &TCPServer{
		Host:      host,
		Port:      port,
		ByteOrder: byteOrder,
		WordOrder: wordOrder,
		Slaves:    slaves,
	}, nil
}

func (s *TCPServer) AddSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Slaves[unitID] = &Slave{}
}

func (s *TCPServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", addr, err)
	}
	defer listener.Close()

	log.Printf("Modbus server started on %s", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 512)
	for {
		conn.SetDeadline(time.Now().Add(5 * time.Minute))
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		request := buffer[:n]
		response, err := s.processRequest(request)
		if err != nil {
			log.Printf("Failed to process request: %v", err)
			response = s.exceptionResponse(request, 0x04) // Server Device Failure
		}

		_, err = conn.Write(response)
		if err != nil {
			log.Printf("Failed to write response: %v", err)
			return
		}
	}
}

func (s *TCPServer) getSlave(unitID byte) (*Slave, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slave, exists := s.Slaves[unitID]
	if !exists {
		return nil, fmt.Errorf("slave %d not found", unitID)
	}
	return slave, nil
}

func (s *TCPServer) processRequest(request []byte) ([]byte, error) {
	if len(request) < 8 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}

	transactionID := request[0:2]
	protocolID := request[2:4]
	length := binary.BigEndian.Uint16(request[4:6])
	unitID := request[6]
	functionCode := request[7]
	data := request[8:]

	// Verify if the actual length matches the length specified in the ADU
	if len(request) != int(length)+6 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}

	slave, err := s.getSlave(unitID)
	if err != nil {
		return s.exceptionResponse(request, 0x04), nil // Server Device Failure
	}

	var response []byte
	switch functionCode {
	case 0x01:
		response, err = s.readCoils(slave, unitID, transactionID, protocolID, data, request)
	case 0x02:
		response, err = s.readDiscreteInputs(slave, unitID, transactionID, protocolID, data, request)
	case 0x03:
		response, err = s.readHoldingRegisters(slave, unitID, transactionID, protocolID, data, request)
	case 0x04:
		response, err = s.readInputRegisters(slave, unitID, transactionID, protocolID, data, request)
	case 0x05:
		response, err = s.writeSingleCoil(slave, unitID, transactionID, protocolID, data, request)
	case 0x06:
		response, err = s.writeSingleRegister(slave, unitID, transactionID, protocolID, data, request)
	case 0x0F:
		response, err = s.writeMultipleCoils(slave, unitID, transactionID, protocolID, data, request)
	case 0x10:
		response, err = s.writeMultipleRegisters(slave, unitID, transactionID, protocolID, data, request)
	default:
		response = s.exceptionResponse(request, 0x01) // Illegal Function
	}

	return response, err
}

func (s *TCPServer) readCoils(slave *Slave, unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])

	if startAddress+quantity > 65535 || !slave.LegalCoilsAddress[startAddress] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	byteCount := (quantity + 7) / 8
	responseLength := 3 + byteCount
	response := make([]byte, 7+responseLength)

	copy(response[0:2], transactionID)
	copy(response[2:4], protocolID)
	binary.BigEndian.PutUint16(response[4:6], uint16(responseLength))
	response[6] = unitID
	response[7] = 0x01
	response[8] = byte(byteCount)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalCoilsAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		if slave.Coils[startAddress+uint16(i)] {
			response[9+i/8] |= 1 << (i % 8)
		}
	}

	return response, nil
}

func (s *TCPServer) readDiscreteInputs(slave *Slave, unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])

	if startAddress+quantity > 65535 || !slave.LegalDiscreteInputsAddress[startAddress] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	byteCount := (quantity + 7) / 8
	responseLength := 3 + byteCount
	response := make([]byte, 7+responseLength)

	copy(response[0:2], transactionID)
	copy(response[2:4], protocolID)
	binary.BigEndian.PutUint16(response[4:6], uint16(responseLength))
	response[6] = unitID
	response[7] = 0x02
	response[8] = byte(byteCount)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalDiscreteInputsAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		if slave.DiscreteInputs[startAddress+uint16(i)] {
			response[9+i/8] |= 1 << (i % 8)
		}
	}

	return response, nil
}

func (s *TCPServer) readHoldingRegisters(slave *Slave, unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])

	if startAddress+quantity > 65535 || !slave.LegalHoldingRegistersAddress[startAddress] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	responseLength := 3 + 2*quantity
	response := make([]byte, 7+responseLength)

	copy(response[0:2], transactionID)
	copy(response[2:4], protocolID)
	binary.BigEndian.PutUint16(response[4:6], uint16(responseLength))
	response[6] = unitID
	response[7] = 0x03
	response[8] = byte(2 * quantity)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		value := slave.HoldingRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(response[9+2*i:], value)
	}

	return response, nil
}

func (s *TCPServer) readInputRegisters(slave *Slave, unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])

	if startAddress+quantity > 65535 || !slave.LegalInputRegistersAddress[startAddress] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	responseLength := 3 + 2*quantity
	response := make([]byte, 7+responseLength)

	copy(response[0:2], transactionID)
	copy(response[2:4], protocolID)
	binary.BigEndian.PutUint16(response[4:6], uint16(responseLength))
	response[6] = unitID
	response[7] = 0x04
	response[8] = byte(2 * quantity)

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalInputRegistersAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		value := slave.InputRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(response[9+2*i:], value)
	}

	return response, nil
}

func (s *TCPServer) writeSingleCoil(slave *Slave, unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	address := binary.BigEndian.Uint16(data[:2])
	value := binary.BigEndian.Uint16(data[2:4])

	if address >= 65535 || !slave.LegalCoilsAddress[address] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	if value == 0xFF00 {
		slave.Coils[address] = true
	} else if value == 0x0000 {
		slave.Coils[address] = false
	} else {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}

	response := make([]byte, 12)
	copy(response[0:2], transactionID)
	copy(response[2:4], protocolID)
	binary.BigEndian.PutUint16(response[4:6], 6)
	response[6] = unitID
	response[7] = 0x05
	copy(response[8:10], data[:2])
	copy(response[10:12], data[2:4])

	return response, nil
}

func (s *TCPServer) writeMultipleCoils(slave *Slave, unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 5 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])
	byteCount := data[4]

	if len(data) < int(5+byteCount) || startAddress+quantity > 65535 || !slave.LegalCoilsAddress[startAddress] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalCoilsAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		if data[5+i/8]&(1<<(i%8)) != 0 {
			slave.Coils[startAddress+uint16(i)] = true
		} else {
			slave.Coils[startAddress+uint16(i)] = false
		}
	}

	response := make([]byte, 12)
	copy(response[0:2], transactionID)
	copy(response[2:4], protocolID)
	binary.BigEndian.PutUint16(response[4:6], 6)
	response[6] = unitID
	response[7] = 0x0F
	copy(response[8:10], data[:2])
	copy(response[10:12], data[2:4])

	return response, nil
}

func (s *TCPServer) writeSingleRegister(slave *Slave, unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	address := binary.BigEndian.Uint16(data[:2])
	value := binary.BigEndian.Uint16(data[2:4])

	if address >= 65535 || !slave.LegalHoldingRegistersAddress[address] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	slave.HoldingRegisters[address] = value

	response := make([]byte, 12)
	copy(response[0:2], transactionID)
	copy(response[2:4], protocolID)
	binary.BigEndian.PutUint16(response[4:6], 6)
	response[6] = unitID
	response[7] = 0x06
	copy(response[8:10], data[:2])
	copy(response[10:12], data[2:4])

	return response, nil
}

func (s *TCPServer) writeMultipleRegisters(slave *Slave, unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 5 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])
	byteCount := data[4]

	if len(data) < int(5+byteCount) || startAddress+quantity > 65535 || !slave.LegalHoldingRegistersAddress[startAddress] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if !slave.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		value := binary.BigEndian.Uint16(data[5+2*i:])
		slave.HoldingRegisters[startAddress+uint16(i)] = value
	}

	response := make([]byte, 12)
	copy(response[0:2], transactionID)
	copy(response[2:4], protocolID)
	binary.BigEndian.PutUint16(response[4:6], 6)
	response[6] = unitID
	response[7] = 0x10
	copy(response[8:10], data[:2])
	copy(response[10:12], data[2:4])

	return response, nil
}

func (s *TCPServer) exceptionResponse(request []byte, exceptionCode byte) []byte {
	transactionID := request[0:2]
	protocolID := request[2:4]
	unitID := request[6]
	functionCode := request[7] | 0x80 // Set the highest bit to indicate an error

	response := make([]byte, 9)
	copy(response[0:2], transactionID)
	copy(response[2:4], protocolID)
	binary.BigEndian.PutUint16(response[4:6], 0x03) // Length of remaining bytes
	response[6] = unitID
	response[7] = functionCode
	response[8] = exceptionCode

	return response
}
