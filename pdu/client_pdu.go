package pdu

import (
	"bytes"
	"encoding/binary"

	"github.com/ulfaric/gomodbus"
	"go.uber.org/zap"
)

// PDURead represents a Modbus PDU for reading operations
type PDURead struct {
	FunctionCode    byte
	StartingAddress uint16
	Quantity        uint16
}

// NewPDUReadCoils creates a new PDU for reading coils
func NewPDUReadCoils(startingAddress, quantity uint16) *PDURead {
	return &PDURead{
		FunctionCode:    gomodbus.ReadCoil,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

// NewPDUReadDiscreteInputs creates a new PDU for reading discrete inputs
func NewPDUReadDiscreteInputs(startingAddress, quantity uint16) *PDURead {
	return &PDURead{
		FunctionCode:    gomodbus.ReadDiscreteInput,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

// NewPDUReadHoldingRegisters creates a new PDU for reading holding registers
func NewPDUReadHoldingRegisters(startingAddress, quantity uint16) *PDURead {
	return &PDURead{
		FunctionCode:    gomodbus.ReadHoldingRegister,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

// NewPDUReadInputRegisters creates a new PDU for reading input registers
func NewPDUReadInputRegisters(startingAddress, quantity uint16) *PDURead {
	return &PDURead{
		FunctionCode:    gomodbus.ReadInputRegister,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

// ToBytes converts the PDURead to a byte slice
func (pdu *PDURead) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(pdu.FunctionCode))
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.Quantity)
	return buffer.Bytes()
}

// FromBytes populates the PDURead fields from a byte slice
func (pdu *PDURead) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing function code for PDURead", zap.Error(err))
		return err
	}
	pdu.FunctionCode = functionCode
	err = binary.Read(buffer, binary.BigEndian, &pdu.StartingAddress)
	if err != nil {
		gomodbus.Logger.Error("error parsing starting address for PDURead", zap.Error(err))
		return err
	}
	err = binary.Read(buffer, binary.BigEndian, &pdu.Quantity)
	if err != nil {
		gomodbus.Logger.Error("error parsing quantity for PDURead", zap.Error(err))
		return err
	}
	return nil
}

// PDUWriteSingleCoil represents a Modbus PDU for writing a single coil
type PDUWriteSingleCoil struct {
	FunctionCode byte
	Address      uint16
	Value        uint16
}

// NewPDUWriteSingleCoil creates a new PDU for writing a single coil
func NewPDUWriteSingleCoil(address uint16, value bool) *PDUWriteSingleCoil {
	var coilValue uint16
	if value {
		coilValue = 0xFF00 // ON
	} else {
		coilValue = 0x0000 // OFF
	}

	return &PDUWriteSingleCoil{
		FunctionCode: gomodbus.WriteSingleCoil,
		Address:      address,
		Value:        coilValue,
	}
}

// ToBytes converts the PDUWriteSingleCoil to a byte slice
func (pdu *PDUWriteSingleCoil) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.Address)
	binary.Write(buffer, binary.BigEndian, pdu.Value)
	return buffer.Bytes()
}

// FromBytes populates the PDUWriteSingleCoil fields from a byte slice
func (pdu *PDUWriteSingleCoil) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing function code for PDUWriteSingleCoil", zap.Error(err))
		return err
	}
	pdu.FunctionCode = functionCode
	err = binary.Read(buffer, binary.BigEndian, &pdu.Address)
	if err != nil {
		gomodbus.Logger.Error("error parsing address for PDUWriteSingleCoil", zap.Error(err))
		return err
	}
	err = binary.Read(buffer, binary.BigEndian, &pdu.Value)
	if err != nil {
		gomodbus.Logger.Error("error parsing value for PDUWriteSingleCoil", zap.Error(err))
		return err
	}
	return nil
}

// PDUWriteMultipleCoils represents a Modbus PDU for writing multiple coils
type PDUWriteMultipleCoils struct {
	FunctionCode    byte
	StartingAddress uint16
	Quantity        uint16
	ByteCount       byte
	Values          []byte
}

// NewPDUWriteMultipleCoils creates a new PDU for writing multiple coils
func NewPDUWriteMultipleCoils(startingAddress uint16, values []bool) *PDUWriteMultipleCoils {
	quantityOfOutputs := uint16(len(values))
	byteCount := (quantityOfOutputs + 7) / 8 // Calculate the number of bytes needed to hold the coil values

	// Initialize output values with the required byte count
	outputValues := make([]byte, byteCount)

	// Pack the boolean values into the output values byte slice
	for i, v := range values {
		if v {
			outputValues[i/8] |= 1 << (i % 8)
		}
	}

	return &PDUWriteMultipleCoils{
		FunctionCode:    gomodbus.WriteMultipleCoils,
		StartingAddress: startingAddress,
		Quantity:        quantityOfOutputs,
		ByteCount:       byte(byteCount),
		Values:          outputValues,
	}
}

// ToBytes converts the PDUWriteMultipleCoils to a byte slice
func (pdu *PDUWriteMultipleCoils) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(pdu.FunctionCode))
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.Quantity)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.Values)
	return buffer.Bytes()
}

// FromBytes populates the PDUWriteMultipleCoils fields from a byte slice
func (pdu *PDUWriteMultipleCoils) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing function code for PDUWriteMultipleCoils", zap.Error(err))
		return err
	}
	pdu.FunctionCode = functionCode
	err = binary.Read(buffer, binary.BigEndian, &pdu.StartingAddress)
	if err != nil {
		gomodbus.Logger.Error("error parsing starting address for PDUWriteMultipleCoils", zap.Error(err))
		return err
	}
	err = binary.Read(buffer, binary.BigEndian, &pdu.Quantity)
	if err != nil {
		gomodbus.Logger.Error("error parsing quantity of outputs for PDUWriteMultipleCoils", zap.Error(err))
		return err
	}
	byteCount, err := buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing byte count for PDUWriteMultipleCoils", zap.Error(err))
		return err
	}
	pdu.ByteCount = byteCount
	pdu.Values = buffer.Bytes()
	return nil
}

// PDUWriteSingleRegister represents a Modbus PDU for writing a single register
type PDUWriteSingleRegister struct {
	FunctionCode byte
	Address      uint16
	Value        []byte
}

// NewPDUWriteSingleRegister creates a new PDU for writing a single register
func NewPDUWriteSingleRegister(registerAddress uint16, registerValue []byte) *PDUWriteSingleRegister {
	return &PDUWriteSingleRegister{
		FunctionCode: gomodbus.WriteSingleRegister,
		Address:      registerAddress,
		Value:        registerValue,
	}
}

// ToBytes converts the PDUWriteSingleRegister to a byte slice
func (pdu *PDUWriteSingleRegister) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.Address)
	buffer.Write(pdu.Value)
	return buffer.Bytes()
}

// FromBytes populates the PDUWriteSingleRegister fields from a byte slice
func (pdu *PDUWriteSingleRegister) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing function code for PDUWriteSingleRegister", zap.Error(err))
		return err
	}
	pdu.FunctionCode = functionCode
	err = binary.Read(buffer, binary.BigEndian, &pdu.Address)
	if err != nil {
		gomodbus.Logger.Error("error parsing address for PDUWriteSingleRegister", zap.Error(err))
		return err
	}
	pdu.Value = buffer.Bytes()
	return nil
}

// PDUWriteMultipleRegisters represents a Modbus PDU for writing multiple registers
type PDUWriteMultipleRegisters struct {
	FunctionCode    byte
	StartingAddress uint16
	Quantity        uint16
	ByteCount       byte
	Values          []byte
}

// NewPDUWriteMultipleRegisters creates a new PDU for writing multiple registers
func NewPDUWriteMultipleRegisters(startingAddress uint16, quantity uint16, values []byte) *PDUWriteMultipleRegisters {
	return &PDUWriteMultipleRegisters{
		FunctionCode:    gomodbus.WriteMultipleRegisters,
		StartingAddress: startingAddress,
		Quantity:        quantity,
		ByteCount:       byte(len(values)),
		Values:          values,
	}
}

// ToBytes converts the PDUWriteMultipleRegisters to a byte slice
func (pdu *PDUWriteMultipleRegisters) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.Quantity)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.Values)
	return buffer.Bytes()
}

// FromBytes populates the PDUWriteMultipleRegisters fields from a byte slice
func (pdu *PDUWriteMultipleRegisters) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing function code", zap.Error(err))
		return err
	}
	pdu.FunctionCode = functionCode
	binary.Read(buffer, binary.BigEndian, &pdu.StartingAddress)
	binary.Read(buffer, binary.BigEndian, &pdu.Quantity)
	byteCount, err := buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Error("error parsing byte count", zap.Error(err))
		return err
	}
	pdu.ByteCount = byteCount
	pdu.Values = buffer.Bytes()
	return nil
}
