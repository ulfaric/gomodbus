package client

import (
	"bytes"
	"fmt"

	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/pdu"
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
		exceptionMessage, exists := gomodbus.ModbusException[exceptionCode]
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
	request := pdu.NewPDUReadCoils(address, quantity)
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
	var response pdu.PDUReadResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return nil, err
	}

	// Check if the function code is correct
	if response.FunctionCode != gomodbus.ReadCoil {
		gomodbus.Logger.Sugar().Errorf("invalid function code for ReadCoils response: %v", err)
		return nil, fmt.Errorf("invalid function code, expected %d, got %d", gomodbus.ReadCoil, response.FunctionCode)
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
	request := pdu.NewPDUReadDiscreteInputs(address, quantity)
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
	var response pdu.PDUReadResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return nil, err
	}

	// Check if the function code is correct
	if response.FunctionCode != gomodbus.ReadDiscreteInput {
		gomodbus.Logger.Sugar().Errorf("invalid function code for ReadDiscreteInputs response: %v", err)
		return nil, fmt.Errorf("invalid function code, expected %d, got %d", gomodbus.ReadDiscreteInput, response.FunctionCode)
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
func ReadHoldingRegisters(c Client, unitID byte, address uint16, quantity uint16) ([][]byte, error) {
	// Create a new PDU for reading holding registers
	request := pdu.NewPDUReadHoldingRegisters(address, quantity)
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
	var response pdu.PDUReadResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return nil, err
	}

	// Check if the function code is correct
	if response.FunctionCode != gomodbus.ReadHoldingRegister {
		gomodbus.Logger.Sugar().Errorf("invalid function code for ReadHoldingRegisters response: %v", err)
		return nil, fmt.Errorf("invalid function code, expected %d, got %d", gomodbus.ReadHoldingRegister, response.FunctionCode)
	}

	// Check if the data length is correct
	if len(response.Data) != int(quantity)*2 {
		gomodbus.Logger.Sugar().Errorf("invalid data length for ReadHoldingRegisters response: %v", err)
		return nil, fmt.Errorf("invalid data length, expected %d, got %d", quantity*2, len(response.Data))
	}

    // Parse the register values from the response PDU
    registers := make([][]byte, quantity)
    for i := uint16(0); i < quantity; i++ {
        registers[i] = response.Data[i*2 : (i+1)*2]
    }

    return registers, nil
}

// ReadInputRegisters reads the values of input registers from a Modbus server
func ReadInputRegisters(c Client, unitID byte, address uint16, quantity uint16) ([][]byte, error) {
	// Create a new PDU for reading input registers
	request := pdu.NewPDUReadInputRegisters(address, quantity)
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
	var response pdu.PDUReadResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return nil, err
	}

	// Check if the function code is correct
	if response.FunctionCode != gomodbus.ReadInputRegister {
		gomodbus.Logger.Sugar().Errorf("invalid function code for ReadInputRegisters response: %v", err)
		return nil, fmt.Errorf("invalid function code, expected %d, got %d", gomodbus.ReadInputRegister, response.FunctionCode)
	}

	// Check if the data length is correct
	if len(response.Data) != int(quantity)*2 {
		gomodbus.Logger.Sugar().Errorf("invalid data length for ReadInputRegisters response: %v", err)
		return nil, fmt.Errorf("invalid data length, expected %d, got %d", quantity*2, len(response.Data))
	}

	// Parse the register values from the response PDU
	registers := make([][]byte, quantity)
	for i := uint16(0); i < quantity; i++ {
		registers[i] = response.Data[i*2 : (i+1)*2]
	}

	return registers, nil
}

// WriteSingleCoil writes a single coil value to a Modbus server
func WriteSingleCoil(c Client, unitID byte, address uint16, value bool) error {
	// Create a new PDU for writing a single coil
	request := pdu.NewPDUWriteSingleCoil(address, value)
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
	var response pdu.PDUWriteSingleCoilResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return err
	}

	// Check if the function code is correct
	if response.FunctionCode != gomodbus.WriteSingleCoil {
		gomodbus.Logger.Sugar().Errorf("invalid function code for WriteSingleCoil response: %v", err)
		return fmt.Errorf("invalid function code, expected %d, got %d", gomodbus.WriteSingleCoil, response.FunctionCode)
	}

	// Check if the response matches the request
	if response.Address != address || response.Value != value {
		gomodbus.Logger.Sugar().Errorf("response does not match request: %v", err)
		return fmt.Errorf("response does not match request, expected address %d, got %d, expected value %t, got %t", address, response.Address, value, response.Value)
	}

	return nil
}

// WriteMultipleCoils writes multiple coil values to a Modbus server
func WriteMultipleCoils(c Client, unitID byte, address uint16, values []bool) error {
	// Create a new PDU for writing multiple coils
	request := pdu.NewPDUWriteMultipleCoils(address, values)
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
	var response pdu.PDUWriteMultipleCoilsResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return err
	}

	// Check if the function code is correct
	if response.FunctionCode != gomodbus.WriteMultipleCoils {
		gomodbus.Logger.Sugar().Errorf("invalid function code for WriteMultipleCoils response: %v", err)
		return fmt.Errorf("invalid function code, expected %d, got %d", gomodbus.WriteMultipleCoils, response.FunctionCode)
	}

	// Check if the response matches the request
	if response.Address != address || response.Quantity != uint16(len(values)) {
		gomodbus.Logger.Sugar().Errorf("response does not match request: %v", err)
		return fmt.Errorf("response does not match request, expected address %d, got %d, expected quantity %d, got %d", address, response.Address, len(values), response.Quantity)
	}

	return nil
}

// WriteSingleRegister writes a single register value to a Modbus server
func WriteSingleRegister(c Client, unitID byte, address uint16, value []byte) error {
	// Create a new PDU for writing a single register
	request := pdu.NewPDUWriteSingleRegister(address, value)
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
	var response pdu.PDUWriteSingleRegisterResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return err
	}

	// Check if the function code is correct
	if response.FunctionCode != gomodbus.WriteSingleRegister {
		gomodbus.Logger.Sugar().Errorf("invalid function code for WriteSingleRegister response: %v", err)
		return fmt.Errorf("invalid function code, expected %d, got %d", gomodbus.WriteSingleRegister, response.FunctionCode)
	}

	// Check if the response matches the request
	if response.Address != address || !bytes.Equal(response.Value, value) {
		gomodbus.Logger.Sugar().Errorf("response does not match request: %v", err)
		return fmt.Errorf("response does not match request, expected address %d, got %d, expected value %v, got %v", address, response.Address, value, response.Value)
	}

	return nil
}

// WriteMultipleRegisters writes multiple register values to a Modbus server
func WriteMultipleRegisters(c Client, unitID byte, address uint16, quantity uint16, values []byte) error {
	// Create a new PDU for writing multiple registers
	request := pdu.NewPDUWriteMultipleRegisters(address, quantity, values)
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
	var response pdu.PDUWriteMultipleRegistersResponse
	err = response.FromBytes(responseBytes)
	if err != nil {
		return err
	}

	// Check if the function code is correct
	if response.FunctionCode != gomodbus.WriteMultipleRegisters {
		gomodbus.Logger.Sugar().Errorf("invalid function code for WriteMultipleRegisters response: %v", err)
		return fmt.Errorf("invalid function code, expected %d, got %d", gomodbus.WriteMultipleRegisters, response.FunctionCode)
	}

	// Check if the response matches the request
	if response.Address != address || response.Quantity != uint16(len(values)/2) {
		gomodbus.Logger.Sugar().Errorf("response does not match request: %v", err)
		return fmt.Errorf("response does not match request, expected address %d, got %d, expected quantity %d, got %d", address, response.Address, len(values)/2, response.Quantity)
	}

	return nil
}
