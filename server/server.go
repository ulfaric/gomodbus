package server

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

const (
	BigEndian    = "big"
LittleEndian = "little"
)

type Server struct {
	Host                         string
	Port                         int
	Coils                        [65535]bool
	LegalCoilsAddress            [65535]bool
	DiscreteInputs               [65535]bool
	LegalDiscreteInputsAddress   [65535]bool
	HoldingRegisters             [65535]uint16
	LegalHoldingRegistersAddress [65535]bool
	InputRegisters               [65535]uint16
	LegalInputRegistersAddress   [65535]bool
	ByteOrder                    string
	WordOrder                    string
}

func NewServer(host, byteOrder, wordOrder string, port int) (*Server, error) {
	if byteOrder != BigEndian && byteOrder != LittleEndian {
		return nil, fmt.Errorf("invalid byte order: %s", byteOrder)
	}
	if wordOrder != BigEndian && wordOrder != LittleEndian {
		return nil, fmt.Errorf("invalid word order: %s", wordOrder)
	}
	return &Server{
		Host:      host,
		Port:      port,
		ByteOrder: byteOrder,
		WordOrder: wordOrder,
	}, nil
}

func (s *Server) Start() error {
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

func (s *Server) handleConnection(conn net.Conn) {
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

func (s *Server) processRequest(request []byte) ([]byte, error) {
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

	var response []byte
	var err error
	switch functionCode {
	case 0x01:
		response, err = s.readCoils(unitID, transactionID, protocolID, data, request)
	case 0x02:
		response, err = s.readDiscreteInputs(unitID, transactionID, protocolID, data, request)
	case 0x03:
		response, err = s.readHoldingRegisters(unitID, transactionID, protocolID, data, request)
	case 0x04:
		response, err = s.readInputRegisters(unitID, transactionID, protocolID, data, request)
	case 0x05:
		response, err = s.writeSingleCoil(unitID, transactionID, protocolID, data, request)
	case 0x06:
		response, err = s.writeSingleRegister(unitID, transactionID, protocolID, data, request)
	case 0x0F:
		response, err = s.writeMultipleCoils(unitID, transactionID, protocolID, data, request)
	case 0x10:
		response, err = s.writeMultipleRegisters(unitID, transactionID, protocolID, data, request)
	default:
		response = s.exceptionResponse(request, 0x01) // Illegal Function
	}

	return response, err
}

func (s *Server) readDiscreteInputs(unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])

	if startAddress+quantity > 65535 || !s.LegalDiscreteInputsAddress[startAddress] {
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
		if !s.LegalDiscreteInputsAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		if s.DiscreteInputs[startAddress+uint16(i)] {
			response[9+i/8] |= 1 << (i % 8)
		}
	}

	return response, nil
}

func (s *Server) readCoils(unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])

	if startAddress+quantity > 65535 || !s.LegalCoilsAddress[startAddress] {
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
		if !s.LegalCoilsAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		if s.Coils[startAddress+uint16(i)] {
			response[9+i/8] |= 1 << (i % 8)
		}
	}

	return response, nil
}

func (s *Server) writeSingleCoil(unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	address := binary.BigEndian.Uint16(data[:2])
	value := binary.BigEndian.Uint16(data[2:4])

	if address >= 65535 || !s.LegalCoilsAddress[address] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	if value == 0xFF00 {
		s.Coils[address] = true
	} else if value == 0x0000 {
		s.Coils[address] = false
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

func (s *Server) writeMultipleCoils(unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 5 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])
	byteCount := data[4]

	if len(data) < int(5+byteCount) || startAddress+quantity > 65535 || !s.LegalCoilsAddress[startAddress] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if !s.LegalCoilsAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		if data[5+i/8]&(1<<(i%8)) != 0 {
			s.Coils[startAddress+uint16(i)] = true
		} else {
			s.Coils[startAddress+uint16(i)] = false
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

func (s *Server) readInputRegisters(unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])

	if startAddress+quantity > 65535 || !s.LegalInputRegistersAddress[startAddress] {
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
		if !s.LegalInputRegistersAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		value := s.InputRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(response[9+2*i:], value)
	}

	return response, nil
}

func (s *Server) readHoldingRegisters(unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])

	if startAddress+quantity > 65535 || !s.LegalHoldingRegistersAddress[startAddress] {
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
		if !s.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		value := s.HoldingRegisters[startAddress+uint16(i)]
		binary.BigEndian.PutUint16(response[9+2*i:], value)
	}

	return response, nil
}

func (s *Server) writeSingleRegister(unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 4 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	address := binary.BigEndian.Uint16(data[:2])
	value := binary.BigEndian.Uint16(data[2:4])

	if address >= 65535 || !s.LegalHoldingRegistersAddress[address] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	s.HoldingRegisters[address] = value

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

func (s *Server) writeMultipleRegisters(unitID byte, transactionID, protocolID, data, request []byte) ([]byte, error) {
	if len(data) < 5 {
		return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
	}
	startAddress := binary.BigEndian.Uint16(data[:2])
	quantity := binary.BigEndian.Uint16(data[2:4])
	byteCount := data[4]

	if len(data) < int(5+byteCount) || startAddress+quantity > 65535 || !s.LegalHoldingRegistersAddress[startAddress] {
		return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
	}

	for i := 0; i < int(quantity); i++ {
		if !s.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
			return s.exceptionResponse(request, 0x02), nil // Illegal Data Address
		}
		value := binary.BigEndian.Uint16(data[5+2*i:])
		s.HoldingRegisters[startAddress+uint16(i)] = value
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


func (s *Server) exceptionResponse(request []byte, exceptionCode byte) []byte {
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
