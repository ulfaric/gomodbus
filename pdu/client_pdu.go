package pdu

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ulfaric/gomodbus"
)

type PDU_Read struct {
	FunctionCode    byte
	StartingAddress uint16
	Quantity        uint16
}

func NewPDUReadCoils(startingAddress, quantity uint16) *PDU_Read {
	return &PDU_Read{
		FunctionCode:    gomodbus.ReadCoil,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

func NewPDUReadDiscreteInputs(startingAddress, quantity uint16) *PDU_Read {
	return &PDU_Read{
		FunctionCode:    gomodbus.ReadDiscreteInput,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

func NewPDUReadHoldingRegisters(startingAddress, quantity uint16) *PDU_Read {
	return &PDU_Read{
		FunctionCode:    gomodbus.ReadHoldingRegister,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

func NewPDUReadInputRegisters(startingAddress, quantity uint16) *PDU_Read {
	return &PDU_Read{
		FunctionCode:    gomodbus.ReadInputRegister,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

func (pdu *PDU_Read) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(pdu.FunctionCode))
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.Quantity)
	return buffer.Bytes()
}

func (pdu *PDU_Read) FromBytes(data []byte) error {
	if len(data) < 5 {
		return fmt.Errorf("insufficient data to parse PDU_Read")
	}
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		return err
	}
	pdu.FunctionCode = functionCode
	binary.Read(buffer, binary.BigEndian, &pdu.StartingAddress)
	binary.Read(buffer, binary.BigEndian, &pdu.Quantity)
	return nil
}

type PDU_WriteSingleCoil struct {
	FunctionCode  byte
	OutputAddress uint16
	OutputValue   uint16
}

func NewPDUWriteSingleCoil(outputAddress uint16, value bool) *PDU_WriteSingleCoil {
	var outputValue uint16
	if value {
		outputValue = 0xFF00 // ON
	} else {
		outputValue = 0x0000 // OFF
	}

	return &PDU_WriteSingleCoil{
		FunctionCode:  gomodbus.WriteSingleCoil,
		OutputAddress: outputAddress,
		OutputValue:   outputValue,
	}
}

func (pdu *PDU_WriteSingleCoil) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.OutputAddress)
	binary.Write(buffer, binary.BigEndian, pdu.OutputValue)
	return buffer.Bytes()
}

func (pdu *PDU_WriteSingleCoil) FromBytes(data []byte) error {
	if len(data) < 5 {
		return fmt.Errorf("insufficient data to parse PDU_WriteSingleCoil")
	}
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		return err
	}
	pdu.FunctionCode = functionCode
	binary.Read(buffer, binary.BigEndian, &pdu.OutputAddress)
	binary.Read(buffer, binary.BigEndian, &pdu.OutputValue)
	return nil
}

type PDU_WriteMultipleCoils struct {
	FunctionCode      byte
	StartingAddress   uint16
	QuantityOfOutputs uint16
	ByteCount         byte
	OutputValues      []byte
}

func NewPDUWriteMultipleCoils(startingAddress uint16, values []bool) *PDU_WriteMultipleCoils {
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

	return &PDU_WriteMultipleCoils{
		FunctionCode:      gomodbus.WriteMultipleCoils,
		StartingAddress:   startingAddress,
		QuantityOfOutputs: quantityOfOutputs,
		ByteCount:         byte(byteCount),
		OutputValues:      outputValues,
	}
}

func (pdu *PDU_WriteMultipleCoils) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(pdu.FunctionCode))
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.QuantityOfOutputs)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.OutputValues)
	return buffer.Bytes()
}

func (pdu *PDU_WriteMultipleCoils) FromBytes(data []byte) error {
	if len(data) < 5 {
		return fmt.Errorf("insufficient data to parse PDU_WriteMultipleCoils")
	}
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		return err
	}
	pdu.FunctionCode = functionCode
	binary.Read(buffer, binary.BigEndian, &pdu.StartingAddress)
	binary.Read(buffer, binary.BigEndian, &pdu.QuantityOfOutputs)
	byteCount, err := buffer.ReadByte()
	if err != nil {
		return err
	}
	pdu.ByteCount = byteCount
	pdu.OutputValues = buffer.Bytes()
	return nil
}

type PDU_WriteSingleRegister struct {
	FunctionCode    byte
	RegisterAddress uint16
	RegisterValue   uint16
}

func NewPDUWriteSingleRegister(registerAddress uint16, registerValue uint16) *PDU_WriteSingleRegister {
	return &PDU_WriteSingleRegister{
		FunctionCode:    gomodbus.WriteSingleRegister,
		RegisterAddress: registerAddress,
		RegisterValue:   registerValue,
	}
}

func (pdu *PDU_WriteSingleRegister) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.RegisterAddress)
	binary.Write(buffer, binary.BigEndian, pdu.RegisterValue)
	return buffer.Bytes()
}

func (pdu *PDU_WriteSingleRegister) FromBytes(data []byte) error {
	if len(data) < 5 {
		return fmt.Errorf("insufficient data to parse PDU_WriteSingleRegister")
	}
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		return err
	}
	pdu.FunctionCode = functionCode
	binary.Read(buffer, binary.BigEndian, &pdu.RegisterAddress)
	binary.Read(buffer, binary.BigEndian, &pdu.RegisterValue)
	return nil
}

type PDU_WriteMultipleRegisters struct {
	FunctionCode      byte
	StartingAddress   uint16
	QuantityOfOutputs uint16
	ByteCount         byte
	OutputValues      []byte
}

func NewPDUWriteMultipleRegisters(startingAddress uint16, values []uint16) *PDU_WriteMultipleRegisters {
	quantityOfOutputs := uint16(len(values))
	byteCount := byte(quantityOfOutputs * 2) // Each register is 2 bytes

	// Initialize output values with the required byte count
	outputValues := make([]byte, byteCount)

	// Pack the register values into the output values byte slice
	for i, v := range values {
		binary.BigEndian.PutUint16(outputValues[i*2:], v)
	}

	return &PDU_WriteMultipleRegisters{
		FunctionCode:      gomodbus.WriteMultipleRegisters,
		StartingAddress:   startingAddress,
		QuantityOfOutputs: quantityOfOutputs,
		ByteCount:         byteCount,
		OutputValues:      outputValues,
	}
}

func (pdu *PDU_WriteMultipleRegisters) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.QuantityOfOutputs)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.OutputValues)
	return buffer.Bytes()
}

func (pdu *PDU_WriteMultipleRegisters) FromBytes(data []byte) error {
	if len(data) < 5 {
		return fmt.Errorf("insufficient data to parse PDU_WriteMultipleRegisters")
	}
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		return err
	}
	pdu.FunctionCode = functionCode
	binary.Read(buffer, binary.BigEndian, &pdu.StartingAddress)
	binary.Read(buffer, binary.BigEndian, &pdu.QuantityOfOutputs)
	byteCount, err := buffer.ReadByte()
	if err != nil {
		return err
	}
	pdu.ByteCount = byteCount
	pdu.OutputValues = buffer.Bytes()
	return nil
}
