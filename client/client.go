package client

import (
	"encoding/binary"
	"fmt"
	"gomodbus"
	"gomodbus/adu"
	"gomodbus/pdu"
	"math"
	"net"
)

type TCPClient struct {
	Host string
	Port int
	conn net.Conn
}

func (client *TCPClient) Connect() error {
	// Check if the host is a valid IP address
	ip := net.ParseIP(client.Host)
	if ip == nil {
		return fmt.Errorf("invalid host IP address")
	}
	// establish a connection to the server
	address := fmt.Sprintf("%s:%d", client.Host, client.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	client.conn = conn
	return nil
}

func (client *TCPClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}

func calculateADULength(functionCode byte, quantity int) (int, error) {
	var pduLength int
	const MBAP_HEADER_LENGTH = 7
	switch functionCode {
	case 0x01, 0x02: // Read Coils, Read Discrete Inputs
		// dataLength is the number of coils/inputs
		byteCount := int(math.Ceil(float64(quantity) / 8.0)) // Number of bytes needed
		pduLength = 1 + 1 + byteCount                        // Function Code + Byte Count + Data
	case 0x03, 0x04: // Read Holding Registers, Read Input Registers
		// dataLength is the number of registers
		byteCount := quantity * 2     // 2 bytes per register
		pduLength = 1 + 1 + byteCount // Function Code + Byte Count + Data
	default:
		// Other function codes (handle accordingly)
		return 0, fmt.Errorf("function code not implemented in this example")
	}

	return MBAP_HEADER_LENGTH + pduLength, nil
}

func (client *TCPClient) ReadCoils(transactionID, startingAddress, quantity, unitID int) ([]bool, error) {
	pdu := pdu.New_PDU_ReadCoils(uint16(startingAddress), uint16(quantity))
	adu := adu.New_TCP_ADU(uint16(transactionID), byte(unitID), pdu.ToBytes())
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return nil, err
	}

	// Calculate the response length
	responseLength, err := calculateADULength(gomodbus.ReadCoil, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	// Verify the function code in the response PDU
	if response[7] != gomodbus.ReadCoil {
		return nil, fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.ReadCoil, response[7])
	}

	// Get the byte count from the response
	byteCount := response[8]
	if int(byteCount) != len(response)-9 {
		return nil, fmt.Errorf("invalid byte count in response, expect %d but received %d", len(response)-9, byteCount)
	}

	// Parse the coil status
	coils := make([]bool, quantity)
	for i := 0; i < quantity; i++ {
		byteIndex := 9 + i/8
		bitIndex := i % 8
		coils[i] = (response[byteIndex] & (1 << bitIndex)) != 0
	}

	return coils, nil
}

func (client *TCPClient) ReadDiscreteInputs(transactionID, startingAddress, quantity, unitID int) ([]bool, error) {
	pdu := pdu.New_PDU_ReadDiscreteInputs(uint16(startingAddress), uint16(quantity))
	adu := adu.New_TCP_ADU(uint16(transactionID), byte(unitID), pdu.ToBytes())
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return nil, err
	}

	// Calculate the response length
	// Calculate the response length
	responseLength, err := calculateADULength(gomodbus.ReadCoil, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	// Verify the function code in the response PDU
	if response[7] != gomodbus.ReadDiscreteInput {
		return nil, fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.ReadDiscreteInput, response[7])
	}

	// Get the byte count from the response
	byteCount := response[8]
	if int(byteCount) != len(response)-9 {
		return nil, fmt.Errorf("invalid byte count in response, expect %d but received %d", len(response)-9, byteCount)
	}

	// Parse the input status
	inputs := make([]bool, quantity)
	for i := 0; i < quantity; i++ {
		byteIndex := 9 + i/8
		bitIndex := i % 8
		inputs[i] = (response[byteIndex] & (1 << bitIndex)) != 0
	}

	return inputs, nil
}

func (client *TCPClient) ReadHoldingRegisters(transactionID, startingAddress, quantity, unitID int) ([]uint16, error) {
	pdu := pdu.New_PDU_ReadHoldingRegisters(uint16(startingAddress), uint16(quantity))
	adu := adu.New_TCP_ADU(uint16(transactionID), byte(unitID), pdu.ToBytes())
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return nil, err
	}

	// Calculate the response length
	responseLength, err := calculateADULength(gomodbus.ReadCoil, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	// Verify the function code in the response PDU
	if response[7] != gomodbus.ReadHoldingRegister {
		return nil, fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.ReadHoldingRegister, response[7])
	}

	// Get the byte count from the response
	byteCount := response[8]
	if int(byteCount) != len(response)-9 {
		return nil, fmt.Errorf("invalid byte count in response, expect %d but received %d", len(response)-9, byteCount)
	}

	// Parse the register values
	registers := make([]uint16, quantity)
	for i := 0; i < quantity; i++ {
		register := binary.BigEndian.Uint16(response[9+i*2 : 9+i*2+2])
		registers[i] = register
	}

	return registers, nil
}

func (client *TCPClient) ReadInputRegisters(transactionID, startingAddress, quantity, unitID int) ([]uint16, error) {
	pdu := pdu.New_PDU_ReadInputRegisters(uint16(startingAddress), uint16(quantity))
	adu := adu.New_TCP_ADU(uint16(transactionID), byte(unitID), pdu.ToBytes())
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return nil, err
	}

	// Calculate the response length
	// Calculate the response length
	responseLength, err := calculateADULength(gomodbus.ReadCoil, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	// Verify the function code in the response PDU
	if response[7] != gomodbus.ReadInputRegister {
		return nil, fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.ReadInputRegister, response[7])
	}

	// Get the byte count from the response
	byteCount := response[8]
	if int(byteCount) != len(response)-9 {
		return nil, fmt.Errorf("invalid byte count in response, expect %d but received %d", len(response)-9, byteCount)
	}

	// Parse the register values
	registers := make([]uint16, quantity)
	for i := 0; i < quantity; i++ {
		register := binary.BigEndian.Uint16(response[9+i*2 : 9+i*2+2])
		registers[i] = register
	}

	return registers, nil
}

func (client *TCPClient) WriteSingleCoil(transactionID, address, unitID int, value bool) error {
	pdu := pdu.New_PDU_WriteSingleCoil(uint16(address), value)
	adu := adu.New_TCP_ADU(uint16(transactionID), byte(unitID), pdu.ToBytes())
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return err
	}

	// Read the response
	response := make([]byte, 12)
	_, err = client.conn.Read(response)
	if err != nil {
		return err
	}

	// Verify the function code in the response PDU
	if response[7] != gomodbus.WriteSingleCoil {
		return fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.WriteSingleCoil, response[7])
	}

	return nil
}