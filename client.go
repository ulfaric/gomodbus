package gomodbus

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Client interface defines the methods for a Modbus client
type Client interface {
	Connect() error
	Disconnect() error
	SendRequest(unitID byte, pduBytes []byte) error
	ReceiveResponse() ([]byte, error)
}

func CheckForException(responsePDU []byte) (bool, string) {
	if responsePDU[0]&0x80 != 0 { // Check if the MSB is set, indicating an exception
		exceptionCode := responsePDU[1]
		exceptionMessage, exists := ModbusException[exceptionCode]
		if !exists {
			return true, "Unknown exception code"
		}
		return true, exceptionMessage
	}
	return false, ""
}

// ReadCoils reads the status of coils from a Modbus server
func ReadCoils(c Client, unitID byte, address uint16, quantity uint16) ([]bool, error) {
	// Create a new PDU for reading coils
	request := NewPDUReadCoils(address, quantity)
	requestBytes := request.ToBytes()

	// Send the request
	err := c.SendRequest(unitID, requestBytes)
	if err != nil {
		return nil, err
	}

	// Receive the response
	responseBytes, err := c.ReceiveResponse()
	if err != nil {
		return nil, err
	}

	exception, message := CheckForException(responseBytes)
	if exception {
		return nil, fmt.Errorf("receivedModbus exception: %s", message)
	}

	// Parse the response
	var response PDUReadResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return nil, err
	}

	// Check if the function code is correct
	if response.FunctionCode != ReadCoil {
		Logger.Sugar().Errorf("invalid function code for ReadCoils response: %v", err)
		return nil, fmt.Errorf("invalid function code, expected %d, got %d", ReadCoil, response.FunctionCode)
	}

	// Parse the coil status from the response PDU
	coils := make([]bool, quantity)
	for i := uint16(0); i < quantity; i++ {
		byteIndex := i / 8
		bitIndex := uint8(i % 8)
		coils[i] = (response.Data[byteIndex] & (1 << bitIndex)) != 0
	}

	return coils, nil
}

// ReadDiscreteInputs reads the status of discrete inputs from a Modbus server
func ReadDiscreteInputs(c Client, unitID byte, address uint16, quantity uint16) ([]bool, error) {
	// Create a new PDU for reading discrete inputs
	request := NewPDUReadDiscreteInputs(address, quantity)
	requestBytes := request.ToBytes()

	// Send the request
	err := c.SendRequest(unitID, requestBytes)
	if err != nil {
		return nil, err
	}

	// Receive the response
	responseBytes, err := c.ReceiveResponse()
	if err != nil {
		return nil, err
	}

	exception, message := CheckForException(responseBytes)
	if exception {
		return nil, fmt.Errorf("received Modbus exception: %s", message)
	}

	// Parse the response
	var response PDUReadResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return nil, err
	}

	// Check if the function code is correct
	if response.FunctionCode != ReadDiscreteInput {
		Logger.Sugar().Errorf("invalid function code for ReadDiscreteInputs response: %v", err)
		return nil, fmt.Errorf("invalid function code, expected %d, got %d", ReadDiscreteInput, response.FunctionCode)
	}

	// Parse the discrete input status from the response PDU
	inputs := make([]bool, quantity)
	for i := uint16(0); i < quantity; i++ {
		byteIndex := i / 8
		bitIndex := uint8(i % 8)
		inputs[i] = (response.Data[byteIndex] & (1 << bitIndex)) != 0
	}

	return inputs, nil
}

// ReadHoldingRegisters reads the values of holding registers from a Modbus server
func ReadHoldingRegisters(c Client, unitID byte, address uint16, quantity uint16) ([]uint16, error) {
	// Create a new PDU for reading holding registers
	request := NewPDUReadHoldingRegisters(address, quantity)
	requestBytes := request.ToBytes()

	// Send the request
	err := c.SendRequest(unitID, requestBytes)
	if err != nil {
		return nil, err
	}

	// Receive the response
	responseBytes, err := c.ReceiveResponse()
	if err != nil {
		return nil, err
	}

	exception, message := CheckForException(responseBytes)
	if exception {
		return nil, fmt.Errorf("received Modbus exception: %s", message)
	}

	// Parse the response
	var response PDUReadResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return nil, err
	}

	// Check if the function code is correct
	if response.FunctionCode != ReadHoldingRegister {
		Logger.Sugar().Errorf("invalid function code for ReadHoldingRegisters response: %v", err)
		return nil, fmt.Errorf("invalid function code, expected %d, got %d", ReadHoldingRegister, response.FunctionCode)
	}

	// Check if the data length is correct
	if len(response.Data) != int(quantity)*2 {
		Logger.Sugar().Errorf("invalid data length for ReadHoldingRegisters response: %v", err)
		return nil, fmt.Errorf("invalid data length, expected %d, got %d", quantity*2, len(response.Data))
	}

	// Parse the register values from the response PDU
	registers := make([]uint16, quantity)
	for i := uint16(0); i < quantity; i++ {
		registers[i] = binary.BigEndian.Uint16(response.Data[i*2 : (i+1)*2])
	}

	return registers, nil
}

// ReadInputRegisters reads the values of input registers from a Modbus server
func ReadInputRegisters(c Client, unitID byte, address uint16, quantity uint16) ([]uint16, error) {
	// Create a new PDU for reading input registers
	request := NewPDUReadInputRegisters(address, quantity)
	requestBytes := request.ToBytes()

	// Send the request
	err := c.SendRequest(unitID, requestBytes)
	if err != nil {
		return nil, err
	}

	// Receive the response
	responseBytes, err := c.ReceiveResponse()
	if err != nil {
		return nil, err
	}

	exception, message := CheckForException(responseBytes)
	if exception {
		return nil, fmt.Errorf("received Modbus exception: %s", message)
	}

	// Parse the response
	var response PDUReadResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return nil, err
	}

	// Check if the function code is correct
	if response.FunctionCode != ReadInputRegister {
		Logger.Sugar().Errorf("invalid function code for ReadInputRegisters response: %v", err)
		return nil, fmt.Errorf("invalid function code, expected %d, got %d", ReadInputRegister, response.FunctionCode)
	}

	// Check if the data length is correct
	if len(response.Data) != int(quantity)*2 {
		Logger.Sugar().Errorf("invalid data length for ReadInputRegisters response: %v", err)
		return nil, fmt.Errorf("invalid data length, expected %d, got %d", quantity*2, len(response.Data))
	}

	// Parse the register values from the response PDU
	registers := make([]uint16, quantity)
	for i := uint16(0); i < quantity; i++ {
		registers[i] = binary.BigEndian.Uint16(response.Data[i*2 : (i+1)*2])
	}

	return registers, nil
}

// WriteCoil writes a single coil value to a Modbus server
func WriteCoil(c Client, unitID byte, address uint16, value bool) error {
	// Create a new PDU for writing a single coil
	request := NewPDUWriteSingleCoil(address, value)
	requestBytes := request.ToBytes()

	// Send the request
	err := c.SendRequest(unitID, requestBytes)
	if err != nil {
		return err
	}

	// Receive the response
	responseBytes, err := c.ReceiveResponse()
	if err != nil {
		return err
	}

	exception, message := CheckForException(responseBytes)
	if exception {
		return fmt.Errorf("received Modbus exception: %s", message)
	}

	// Parse the response
	var response PDUWriteSingleCoilResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return err
	}

	// Check if the function code is correct
	if response.FunctionCode != WriteSingleCoil {
		Logger.Sugar().Errorf("invalid function code for WriteSingleCoil response: %v", err)
		return fmt.Errorf("invalid function code, expected %d, got %d", WriteSingleCoil, response.FunctionCode)
	}

	// Check if the response matches the request
	if response.Address != address || response.Value != value {
		Logger.Sugar().Errorf("response does not match request: %v", err)
		return fmt.Errorf("response does not match request, expected address %d, got %d, expected value %t, got %t", address, response.Address, value, response.Value)
	}

	return nil
}

// WriteCoils writes multiple coil values to a Modbus server
func WriteCoils(c Client, unitID byte, address uint16, values []bool) error {
	// Create a new PDU for writing multiple coils
	request := NewPDUWriteMultipleCoils(address, values)
	requestBytes := request.ToBytes()

	// Send the request
	err := c.SendRequest(unitID, requestBytes)
	if err != nil {
		return err
	}

	// Receive the response
	responseBytes, err := c.ReceiveResponse()
	if err != nil {
		return err
	}

	exception, message := CheckForException(responseBytes)
	if exception {
		return fmt.Errorf("received Modbus exception: %s", message)
	}

	// Parse the response
	var response PDUWriteMultipleCoilsResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return err
	}

	// Check if the function code is correct
	if response.FunctionCode != WriteMultipleCoils {
		Logger.Sugar().Errorf("invalid function code for WriteMultipleCoils response: %v", err)
		return fmt.Errorf("invalid function code, expected %d, got %d", WriteMultipleCoils, response.FunctionCode)
	}

	// Check if the response matches the request
	if response.Address != address || response.Quantity != uint16(len(values)) {
		Logger.Sugar().Errorf("response does not match request: %v", err)
		return fmt.Errorf("response does not match request, expected address %d, got %d, expected quantity %d, got %d", address, response.Address, len(values), response.Quantity)
	}

	return nil
}

// WriteRegister writes a single register value to a Modbus server
func WriteRegister(c Client, unitID byte, address uint16, value []byte) error {
	// Create a new PDU for writing a single register
	request := NewPDUWriteSingleRegister(address, value)
	requestBytes := request.ToBytes()

	// Send the request
	err := c.SendRequest(unitID, requestBytes)
	if err != nil {
		return err
	}

	// Receive the response
	responseBytes, err := c.ReceiveResponse()
	if err != nil {
		return err
	}

	exception, message := CheckForException(responseBytes)
	if exception {
		return fmt.Errorf("received Modbus exception: %s", message)
	}

	// Parse the response
	var response PDUWriteSingleRegisterResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return err
	}

	// Check if the function code is correct
	if response.FunctionCode != WriteSingleRegister {
		Logger.Sugar().Errorf("invalid function code for WriteSingleRegister response: %v", err)
		return fmt.Errorf("invalid function code, expected %d, got %d", WriteSingleRegister, response.FunctionCode)
	}

	// Check if the response matches the request
	if response.Address != address || !bytes.Equal(response.Value, value) {
		Logger.Sugar().Errorf("response does not match request: %v", err)
		return fmt.Errorf("response does not match request, expected address %d, got %d, expected value %v, got %v", address, response.Address, value, response.Value)
	}

	return nil
}

// WriteRegisters writes multiple register values to a Modbus server
func WriteRegisters(c Client, unitID byte, address uint16, quantity uint16, values []byte) error {
	// Create a new PDU for writing multiple registers
	request := NewPDUWriteMultipleRegisters(address, quantity, values)
	requestBytes := request.ToBytes()

	// Send the request
	err := c.SendRequest(unitID, requestBytes)
	if err != nil {
		return err
	}

	// Receive the response
	responseBytes, err := c.ReceiveResponse()
	if err != nil {
		return err
	}

	exception, message := CheckForException(responseBytes)
	if exception {
		return fmt.Errorf("received Modbus exception: %s", message)
	}

	// Parse the response
	var response PDUWriteMultipleRegistersResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return err
	}

	// Check if the function code is correct
	if response.FunctionCode != WriteMultipleRegisters {
		Logger.Sugar().Errorf("invalid function code for WriteMultipleRegisters response: %v", err)
		return fmt.Errorf("invalid function code, expected %d, got %d", WriteMultipleRegisters, response.FunctionCode)
	}

	// Check if the response matches the request
	if response.Address != address || response.Quantity != uint16(len(values)/2) {
		Logger.Sugar().Errorf("response does not match request: %v", err)
		return fmt.Errorf("response does not match request, expected address %d, got %d, expected quantity %d, got %d", address, response.Address, len(values)/2, response.Quantity)
	}

	return nil
}
