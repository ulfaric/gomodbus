package pdu

import (
	"bytes"
	"encoding/binary"
	"github.com/ulfaric/gomodbus"
)

type PDU_Read struct {
	FunctionCode    byte
	startingAddress uint16
	quantity        uint16
}

func New_PDU_ReadCoils(startingAddress, quantity uint16) *PDU_Read {
	return &PDU_Read{
		FunctionCode:    gomodbus.ReadCoil,
		startingAddress: startingAddress,
		quantity:        quantity,
	}
}

func New_PDU_ReadDiscreteInputs(startingAddress, quantity uint16) *PDU_Read {
	return &PDU_Read{
		FunctionCode:    gomodbus.ReadDiscreteInput,
		startingAddress: startingAddress,
		quantity:        quantity,
	}
}

func New_PDU_ReadHoldingRegisters(startingAddress, quantity uint16) *PDU_Read {
	return &PDU_Read{
		FunctionCode:    gomodbus.ReadHoldingRegister,
		startingAddress: startingAddress,
		quantity:        quantity,
	}
}

func New_PDU_ReadInputRegisters(startingAddress, quantity uint16) *PDU_Read {
	return &PDU_Read{
		FunctionCode:    gomodbus.ReadInputRegister,
		startingAddress: startingAddress,
		quantity:        quantity,
	}
}

func (pdu *PDU_Read) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(pdu.FunctionCode))
	binary.Write(buffer, binary.BigEndian, pdu.startingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.quantity)
	return buffer.Bytes()
}

func (pdu *PDU_Read) FromBytes(data []byte) {
	buffer := bytes.NewBuffer(data)
	functionCode, err := buffer.ReadByte()
	if err != nil {
		panic(err)
	}
	if functionCode != gomodbus.ReadCoil {
		panic("Invalid function code")
	}
	pdu.FunctionCode = functionCode
	binary.Read(buffer, binary.BigEndian, &pdu.startingAddress)
	binary.Read(buffer, binary.BigEndian, &pdu.quantity)
}

func (pdu *PDU_Read) Validate(data []byte) bool {
	return len(data) == 5
}

type PDU_WriteSingleCoil struct {
	FunctionCode    byte
	OutputAddress   uint16
	OutputValue     uint16
}

// New_PDU_WriteSingleCoil creates a new PDU_WriteSingleCoil
func New_PDU_WriteSingleCoil(outputAddress uint16, value bool) *PDU_WriteSingleCoil {
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

// ToBytes converts the PDU to a byte slice
func (pdu *PDU_WriteSingleCoil) ToBytes() []byte {
	buffer := new(bytes.Buffer)

	// Write the PDU fields to the buffer
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.OutputAddress)
	binary.Write(buffer, binary.BigEndian, pdu.OutputValue)

	return buffer.Bytes()
}

type PDU_WriteMultipleCoils struct {
	FunctionCode      byte
	StartingAddress   uint16
	QuantityOfOutputs uint16
	ByteCount         byte
	OutputValues      []byte
}

func New_PDU_WriteMultipleCoils(startingAddress uint16, values []bool) *PDU_WriteMultipleCoils {
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

	// Write the PDU fields to the buffer
	buffer.WriteByte(byte(pdu.FunctionCode))
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.QuantityOfOutputs)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.OutputValues)

	return buffer.Bytes()
}

// PDU_WriteSingleRegister represents the PDU for the Write Single Register function
type PDU_WriteSingleRegister struct {
	FunctionCode    byte
	RegisterAddress uint16
	RegisterValue   uint16
}

// New_PDU_WriteSingleRegister creates a new PDU_WriteSingleRegister
func New_PDU_WriteSingleRegister(registerAddress uint16, registerValue uint16) *PDU_WriteSingleRegister {
	return &PDU_WriteSingleRegister{
		FunctionCode:    gomodbus.WriteSingleRegister,
		RegisterAddress: registerAddress,
		RegisterValue:   registerValue,
	}
}

// ToBytes converts the PDU to a byte slice
func (pdu *PDU_WriteSingleRegister) ToBytes() []byte {
	buffer := new(bytes.Buffer)

	// Write the PDU fields to the buffer
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.RegisterAddress)
	binary.Write(buffer, binary.BigEndian, pdu.RegisterValue)

	return buffer.Bytes()
}

// PDU_WriteMultipleRegisters represents the PDU for the Write Multiple Registers function
type PDU_WriteMultipleRegisters struct {
	FunctionCode      byte
	StartingAddress   uint16
	QuantityOfOutputs uint16
	ByteCount         byte
	OutputValues      []byte
}

// New_PDU_WriteMultipleRegisters creates a new PDU_WriteMultipleRegisters
func New_PDU_WriteMultipleRegisters(startingAddress uint16, values []uint16) *PDU_WriteMultipleRegisters {
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

// ToBytes converts the PDU to a byte slice
func (pdu *PDU_WriteMultipleRegisters) ToBytes() []byte {
	buffer := new(bytes.Buffer)

	// Write the PDU fields to the buffer
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.QuantityOfOutputs)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.OutputValues)

	return buffer.Bytes()
}