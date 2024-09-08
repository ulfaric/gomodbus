package pdu

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ulfaric/gomodbus"
	"go.uber.org/zap"
)

// PDUReadResponse represents a response for a read request
type PDUReadResponse struct {
	FunctionCode byte
	ByteCount    byte
	Data         []byte
}

// NewPDUReadCoilsResponse creates a new PDUReadResponse for reading coils
func NewPDUReadCoilsResponse(coils []bool) *PDUReadResponse {
	coilsBytes := make([]byte, (len(coils)+7)/8)
	for i, coil := range coils {
		if coil {
			coilsBytes[i/8] |= 1 << (i % 8)
		}
	}
	return &PDUReadResponse{
		FunctionCode: gomodbus.ReadCoil,
		ByteCount:    byte(len(coilsBytes)),
		Data:         coilsBytes,
	}
}

// NewPDUReadDiscreteInputsResponse creates a new PDUReadResponse for reading discrete inputs
func NewPDUReadDiscreteInputsResponse(inputs []bool) *PDUReadResponse {
	inputsBytes := make([]byte, (len(inputs)+7)/8)
	for i, input := range inputs {
		if input {
			inputsBytes[i/8] |= 1 << (i % 8)
		}
	}
	return &PDUReadResponse{
		FunctionCode: gomodbus.ReadDiscreteInput,
		ByteCount:    byte(len(inputsBytes)),
		Data:         inputsBytes,
	}
}

// NewPDUReadHoldingRegistersResponse creates a new PDUReadResponse for reading holding registers
func NewPDUReadHoldingRegistersResponse(registers []byte) *PDUReadResponse {
	return &PDUReadResponse{
		FunctionCode: gomodbus.ReadHoldingRegister,
		ByteCount:    byte(len(registers)),
		Data:         registers,
	}
}

// NewPDUReadInputRegistersResponse creates a new PDUReadResponse for reading input registers
func NewPDUReadInputRegistersResponse(registers []byte) *PDUReadResponse {
	return &PDUReadResponse{
		FunctionCode: gomodbus.ReadInputRegister,
		ByteCount:    byte(len(registers)),
		Data:         registers,
	}
}


// ToBytes converts PDUReadResponse to bytes
func (pdu *PDUReadResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.Data)
	return buffer.Bytes()
}

// FromBytes parses bytes into PDUReadResponse
func (pdu *PDUReadResponse) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var err error

	// Read FunctionCode
	pdu.FunctionCode, err = buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing FunctionCode for PDUReadResponse", zap.Error(err))
		return err
	}

	// Read ByteCount
	pdu.ByteCount, err = buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing ByteCount for PDUReadResponse", zap.Error(err))
		return err
	}

	// Read Data
	pdu.Data = buffer.Bytes()
	if len(pdu.Data) != int(pdu.ByteCount) {
		gomodbus.Logger.Sugar().Errorf("data length mismatch, expected %d, got %d", pdu.ByteCount, len(pdu.Data))
		return fmt.Errorf("data length mismatch, expected %d, got %d", pdu.ByteCount, len(pdu.Data))
	}

	return nil
}

// PDUWriteSingleCoilResponse represents a response for a write single coil request
type PDUWriteSingleCoilResponse struct {
	FunctionCode byte
	Address      uint16
	Value        bool
}

// NewWriteSingleCoilResponse creates a new PDUWriteSingleCoilResponse
func NewWriteSingleCoilResponse(address uint16, value bool) *PDUWriteSingleCoilResponse {
	return &PDUWriteSingleCoilResponse{
		FunctionCode: gomodbus.WriteSingleCoil,
		Address:      address,
		Value:        value,
	}
}

// ToBytes converts PDUWriteSingleCoilResponse to bytes
func (pdu *PDUWriteSingleCoilResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.Address)
	if pdu.Value {
		binary.Write(buffer, binary.BigEndian, uint16(0xFF00))
	} else {
		binary.Write(buffer, binary.BigEndian, uint16(0x0000))
	}
	return buffer.Bytes()
}

// FromBytes parses bytes into PDUWriteSingleCoilResponse
func (pdu *PDUWriteSingleCoilResponse) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var err error

	// Read FunctionCode
	pdu.FunctionCode, err = buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing FunctionCode for PDUWriteSingleCoilResponse", zap.Error(err))
		return err
	}

	// Read Address
	err = binary.Read(buffer, binary.BigEndian, &pdu.Address)
	if err != nil {
		gomodbus.Logger.Error("error parsing Address for PDUWriteSingleCoilResponse", zap.Error(err))
		return err
	}

	// Read Value
	var coilValue uint16
	err = binary.Read(buffer, binary.BigEndian, &coilValue)
	if err != nil {
		gomodbus.Logger.Error("error parsing coil value for PDUWriteSingleCoilResponse", zap.Error(err))
		return err
	}
	pdu.Value = coilValue == 0xFF00

	return nil
}

// PDUWriteMultipleCoilsResponse represents a response for a write multiple coils request
type PDUWriteMultipleCoilsResponse struct {
	FunctionCode byte
	Address      uint16
	Quantity     uint16
}

// NewWriteMultipleCoilsResponse creates a new PDUWriteMultipleCoilsResponse
func NewWriteMultipleCoilsResponse(address, quantity uint16) *PDUWriteMultipleCoilsResponse {
	return &PDUWriteMultipleCoilsResponse{
		FunctionCode: gomodbus.WriteMultipleCoils,
		Address:      address,
		Quantity:     quantity,
	}
}

// ToBytes converts PDUWriteMultipleCoilsResponse to bytes
func (pdu *PDUWriteMultipleCoilsResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.Address)
	binary.Write(buffer, binary.BigEndian, pdu.Quantity)
	return buffer.Bytes()
}

// FromBytes parses bytes into PDUWriteMultipleCoilsResponse
func (pdu *PDUWriteMultipleCoilsResponse) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var err error

	// Read FunctionCode
	pdu.FunctionCode, err = buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing FunctionCode for PDUWriteMultipleCoilsResponse", zap.Error(err))
		return err
	}

	// Read Address
	err = binary.Read(buffer, binary.BigEndian, &pdu.Address)
	if err != nil {
		gomodbus.Logger.Error("error parsing Address for PDUWriteMultipleCoilsResponse", zap.Error(err))
		return err
	}

	// Read Quantity
	err = binary.Read(buffer, binary.BigEndian, &pdu.Quantity)
	if err != nil {
		gomodbus.Logger.Error("error parsing Quantity for PDUWriteMultipleCoilsResponse", zap.Error(err))
		return err
	}

	return nil
}

// PDUWriteSingleRegisterResponse represents a response for a write single register request
type PDUWriteSingleRegisterResponse struct {
	FunctionCode byte
	Address      uint16
	Value        []byte
}

// NewWriteSingleRegisterResponse creates a new PDUWriteSingleRegisterResponse
func NewWriteSingleRegisterResponse(address uint16, value []byte) *PDUWriteSingleRegisterResponse {
	return &PDUWriteSingleRegisterResponse{
		FunctionCode: gomodbus.WriteSingleRegister,
		Address:      address,
		Value:        value,
	}
}

// ToBytes converts PDUWriteSingleRegisterResponse to bytes
func (pdu *PDUWriteSingleRegisterResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.Address)
	buffer.Write(pdu.Value)
	return buffer.Bytes()
}

// FromBytes parses bytes into PDUWriteSingleRegisterResponse
func (pdu *PDUWriteSingleRegisterResponse) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var err error

	// Read FunctionCode
	pdu.FunctionCode, err = buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing FunctionCode for PDUWriteSingleRegisterResponse", zap.Error(err))
		return err
	}

	// Read Address
	err = binary.Read(buffer, binary.BigEndian, &pdu.Address)
	if err != nil {
		gomodbus.Logger.Error("error parsing Address for PDUWriteSingleRegisterResponse", zap.Error(err))
		return err
	}

	// Read Value
	pdu.Value = buffer.Bytes()
	return nil
}

// PDUWriteMultipleRegistersResponse represents a response for a write multiple registers request
type PDUWriteMultipleRegistersResponse struct {
	FunctionCode byte
	Address      uint16
	Quantity     uint16
}

// NewWriteMultipleRegistersResponse creates a new PDUWriteMultipleRegistersResponse
func NewWriteMultipleRegistersResponse(address, quantity uint16) *PDUWriteMultipleRegistersResponse {
	return &PDUWriteMultipleRegistersResponse{
		FunctionCode: gomodbus.WriteMultipleRegisters,
		Address:      address,
		Quantity:     quantity,
	}
}

// ToBytes converts PDUWriteMultipleRegistersResponse to bytes
func (pdu *PDUWriteMultipleRegistersResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.Address)
	binary.Write(buffer, binary.BigEndian, pdu.Quantity)
	return buffer.Bytes()
}

// FromBytes parses bytes into PDUWriteMultipleRegistersResponse
func (pdu *PDUWriteMultipleRegistersResponse) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var err error

	// Read FunctionCode
	pdu.FunctionCode, err = buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing FunctionCode for PDUWriteMultipleRegistersResponse", zap.Error(err))
		return err
	}

	// Read Address
	err = binary.Read(buffer, binary.BigEndian, &pdu.Address)
	if err != nil {
		gomodbus.Logger.Error("error parsing Address for PDUWriteMultipleRegistersResponse", zap.Error(err))
		return err
	}

	// Read Quantity
	err = binary.Read(buffer, binary.BigEndian, &pdu.Quantity)
	if err != nil {
		gomodbus.Logger.Error("error parsing Quantity for PDUWriteMultipleRegistersResponse", zap.Error(err))
		return err
	}

	return nil
}

// PDUErrorResponse represents a Modbus exception response
type PDUErrorResponse struct {
	FunctionCode  byte
	ExceptionCode byte
}

// NewPDUErrorResponse creates a new PDUErrorResponse
func NewPDUErrorResponse(functionCode, exceptionCode byte) *PDUErrorResponse {
	return &PDUErrorResponse{
		FunctionCode:  functionCode | 0x80, // Set the MSB to indicate an exception
		ExceptionCode: exceptionCode,
	}
}

// ToBytes converts PDUErrorResponse to bytes
func (pdu *PDUErrorResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	buffer.WriteByte(pdu.ExceptionCode)
	return buffer.Bytes()
}

// FromBytes parses bytes into PDUErrorResponse
func (pdu *PDUErrorResponse) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var err error

	// Read FunctionCode
	pdu.FunctionCode, err = buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing FunctionCode for PDUErrorResponse", zap.Error(err))
		return err
	}

	// Read ExceptionCode
	pdu.ExceptionCode, err = buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing ExceptionCode for PDUErrorResponse", zap.Error(err))
		return err
	}

	return nil
}
