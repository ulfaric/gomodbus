package client

import (
	"bytes"
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
	case 0x05: // Write Single Coil
		// The response PDU for Write Single Coil includes the function code, output address, and output value
		pduLength = 1 + 2 + 2 // Function Code + Output Address + Output Value
	case 0x06: // Write Single Register
		// The response PDU for Write Single Register includes the function code, register address, and register value
		pduLength = 1 + 2 + 2 // Function Code + Register Address + Register Value
	case 0x0F: // Write Multiple Coils
		// The response PDU for Write Multiple Coils includes the function code, starting address, and quantity of coils
		pduLength = 1 + 2 + 2 // Function Code + Starting Address + Quantity of Coils
	case 0x10: // Write Multiple Registers
		// The response PDU for Write Multiple Registers includes the function code, starting address, and quantity of registers
		pduLength = 1 + 2 + 2 // Function Code + Starting Address + Quantity of Registers
	default:
		// Other function codes (handle accordingly)
		return 0, fmt.Errorf("function code not implemented")
	}

	return MBAP_HEADER_LENGTH + pduLength, nil
}

func DecodeModbusRegisters(registers []uint16, byte_order, word_order string) interface{} {
	swapBytes := func(value uint16) uint16 {
		return (value>>8)&0xFF | (value&0xFF)<<8
	}

	switch len(registers) {
	case 1:
		// Swap bytes for single uint16
		if byte_order == gomodbus.LittleEndian {
			return swapBytes(registers[0])
		} else {
			return registers[0]
		}
	case 2:
		// Combine two uint16 into a single 32-bit integer
		var combined uint32
		if word_order == gomodbus.LittleEndian {
			for i := 1; i >= 0; i-- {
				if byte_order == gomodbus.LittleEndian {
					swapped := swapBytes(registers[i])
					combined = (combined << 16) | uint32(swapped)
				} else {
					combined = (combined << 16) | uint32(registers[i])
				}
			}
		} else {
			for i := 0; i < 2; i++ {
				if byte_order == gomodbus.LittleEndian {
					swapped := swapBytes(registers[i])
					combined = (combined << 16) | uint32(swapped)
				} else {
					combined = (combined << 16) | uint32(registers[i])
				}
			}
		}
		return combined
	case 4:
		// Combine four uint16 into a single 64-bit integer
		var combined uint64
		if word_order == gomodbus.LittleEndian {
			for i := 3; i >= 0; i-- {
				if byte_order == gomodbus.LittleEndian {
					swapped := swapBytes(registers[i])
					combined = (combined << 16) | uint64(swapped)
				} else {
					combined = (combined << 16) | uint64(registers[i])
				}
			}
		} else {
			for i := 0; i < 4; i++ {
				if byte_order == gomodbus.LittleEndian {
					swapped := swapBytes(registers[i])
					combined = (combined << 16) | uint64(swapped)
				} else {
					combined = (combined << 16) | uint64(registers[i])
				}
			}
		}
		return combined
	default:
		panic("The input must contain exactly 1, 2, or 4 uint16 values.")
	}
}

func EncodeModbusRegisters(value interface{}, byte_order, word_order string) []uint16 {
	splitUint32 := func(value uint32) []uint16 {
		return []uint16{uint16(value >> 16), uint16(value & 0xFFFF)}
	}

	splitUint64 := func(value uint64) []uint16 {
		return []uint16{uint16(value >> 48), uint16(value >> 32), uint16(value >> 16), uint16(value)}
	}

	swapBytes := func(value uint16) uint16 {
		return (value>>8)&0xFF | (value&0xFF)<<8
	}

	reverseSlice := func(slice []uint16) {
		for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
			slice[i], slice[j] = slice[j], slice[i]
		}
	}

	switch value := value.(type) {
	case uint16:
		// Swap bytes for single uint16
		if byte_order == gomodbus.LittleEndian {
			return []uint16{swapBytes(value)}
		} else {
			return []uint16{value}
		}
	case int16:
		// convert int16 to uint16 and swap bytes
		if byte_order == gomodbus.LittleEndian {
			return []uint16{swapBytes(uint16(value))}
		} else {
			return []uint16{uint16(value)}
		}
	case uint32:
		// Split 32-bit integer into two uint16
		registers := splitUint32(value)
		if word_order == gomodbus.LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == gomodbus.LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case int32:
		// convert int32 to uint32 and split into two uint16
		registers := splitUint32(uint32(value))
		if word_order == gomodbus.LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == gomodbus.LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case uint64:
		// Split 64-bit integer into four uint16
		registers := splitUint64(value)
		if word_order == gomodbus.LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == gomodbus.LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case int64:
		// Convert int64 to uint64 and split into four uint16
		registers := splitUint64(uint64(value))
		if word_order == gomodbus.LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == gomodbus.LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case float32:
		// Convert float32 to uint32 and split into two uint16
		bits := math.Float32bits(value)
		registers := splitUint32(uint32(bits))
		if word_order == gomodbus.LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == gomodbus.LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case float64:
		// Convert float64 to uint64 and split into four uint16
		bits := math.Float64bits(value)
		registers := splitUint64(bits)
		if word_order == gomodbus.LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == gomodbus.LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	default:
		panic("The input must be a uint16, uint32, or uint64 value.")
	}
}

// CheckModbusError checks if the response ADU indicates an error.
func CheckModbusError(response []byte) error {
	// Minimum length for a valid Modbus response ADU
	if len(response) < 9 {
		return fmt.Errorf("invalid response length: %d", len(response))
	}

	// Extract the function code from the response
	functionCode := response[7]

	// Check if the most significant bit of the function code is set
	if functionCode&0x80 != 0 {
		exceptionCode := response[8]
		return fmt.Errorf("modbus exception: %s", getExceptionMessage(exceptionCode))
	}

	return nil
}

// getExceptionMessage returns a human-readable message for a given exception code.
func getExceptionMessage(exceptionCode byte) string {
	switch exceptionCode {
	case 0x01:
		return "Illegal Function"
	case 0x02:
		return "Illegal Data Address"
	case 0x03:
		return "Illegal Data Value"
	case 0x04:
		return "Server Device Failure"
	default:
		return fmt.Sprintf("Unknown exception code: %02X", exceptionCode)
	}
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

	// Check for ADU errors
	err = CheckModbusError(response)
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

	// Check for ADU errors
	err = CheckModbusError(response)
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
	responseLength, err := calculateADULength(gomodbus.ReadHoldingRegister, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	// Check for ADU errors
	err = CheckModbusError(response)
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
	responseLength, err := calculateADULength(gomodbus.ReadInputRegister, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	// Check for ADU errors
	err = CheckModbusError(response)
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
	responseLength, err := calculateADULength(gomodbus.WriteSingleCoil, 1)
	if err != nil {
		return fmt.Errorf("failed to calculate response length: %v", err)
	}
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return err
	}

	// Check for ADU errors
	err = CheckModbusError(response)
	if err != nil {
		return err
	}

	// Verify the function code in the response PDU
	if response[7] != gomodbus.WriteSingleCoil {
		return fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.WriteSingleCoil, response[7])
	}

	if !bytes.Equal(response[8:10], adu_bytes[8:10]) {
		return fmt.Errorf("invalid address in response, expect %v but received %v", adu_bytes[8:10], response[8:10])
	}

	if !bytes.Equal(response[10:12], adu_bytes[10:12]) {
		return fmt.Errorf("invalid value in response, expect %v but received %v", adu_bytes[10:12], response[10:12])
	}

	return nil
}

func (client *TCPClient) WriteSingleRegister(transactionID, address, unitID int, value uint16) error {
	pdu := pdu.New_PDU_WriteSingleRegister(uint16(address), value)
	adu := adu.New_TCP_ADU(uint16(transactionID), byte(unitID), pdu.ToBytes())
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return err
	}

	// Read the response
	responseLength, err := calculateADULength(gomodbus.WriteSingleRegister, 1)
	if err != nil {
		return fmt.Errorf("failed to calculate response length: %v", err)
	}
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return err
	}

	// Check for ADU errors
	err = CheckModbusError(response)
	if err != nil {
		return err
	}

	// Verify the function code in the response PDU
	if response[7] != gomodbus.WriteSingleRegister {
		return fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.WriteSingleRegister, response[7])
	}

	// Verify the address and value in the response PDU
	if !bytes.Equal(response[8:10], adu_bytes[8:10]) {
		return fmt.Errorf("invalid address in response, expect %v but received %v", adu_bytes[8:10], response[8:10])
	}

	if !bytes.Equal(response[10:12], adu_bytes[10:12]) {
		return fmt.Errorf("invalid value in response, expect %v but received %v", adu_bytes[10:12], response[10:12])
	}
	
	return nil
}

func (client *TCPClient) WriteMultipleCoils(transactionID, startingAddress, unitID int, values []bool) error {
	pdu := pdu.New_PDU_WriteMultipleCoils(uint16(startingAddress), values)
	adu := adu.New_TCP_ADU(uint16(transactionID), byte(unitID), pdu.ToBytes())
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return err
	}

	// Read the response
	responseLength, err := calculateADULength(gomodbus.WriteMultipleCoils, 1)
	if err != nil {
		return fmt.Errorf("failed to calculate response length: %v", err)
	}
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return err
	}

	// Check for ADU errors
	err = CheckModbusError(response)
	if err != nil {
		return err
	}

	// Verify the function code in the response PDU
	if response[7] != gomodbus.WriteMultipleCoils {
		return fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.WriteMultipleCoils, response[7])
	}

	if !bytes.Equal(response[8:10], adu_bytes[8:10]) {
		return fmt.Errorf("invalid starting address in response, expect %v but received %v", adu_bytes[8:10], response[8:10])
	}

	if !bytes.Equal(response[10:12], adu_bytes[10:12]) {
		return fmt.Errorf("invalid quantity in response, expect %v but received %v", adu_bytes[10:12], response[10:12])
	}

	return nil
}

func (client *TCPClient) WriteMultipleRegisters(transactionID, startingAddress, unitID int, values []uint16) error {
	pdu := pdu.New_PDU_WriteMultipleRegisters(uint16(startingAddress), values)
	adu := adu.New_TCP_ADU(uint16(transactionID), byte(unitID), pdu.ToBytes())
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return err
	}

	// Read the response
	responseLength, err := calculateADULength(gomodbus.WriteMultipleRegisters, 1)
	if err != nil {
		return fmt.Errorf("failed to calculate response length: %v", err)
	}
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return err
	}

	// Check for ADU errors
	err = CheckModbusError(response)
	if err != nil {
		return err
	}

	// Verify the function code in the response PDU
	if response[7] != gomodbus.WriteMultipleRegisters {
		return fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.WriteMultipleRegisters, response[7])
	}

	
	if !bytes.Equal(response[8:10], adu_bytes[8:10]) {
		return fmt.Errorf("invalid starting address in response, expect %v but received %v", adu_bytes[8:10], response[8:10])
	}

	if !bytes.Equal(response[10:12], adu_bytes[10:12]) {
		return fmt.Errorf("invalid quantity in response, expect %v but received %v", adu_bytes[10:12], response[10:12])
	}

	return nil
}
