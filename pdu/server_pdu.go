package pdu

import (
	"bytes"
	"encoding/binary"
)

type PDU_ReadResponse struct {
	FunctionCode byte
	ByteCount    byte
	Data         []byte
}

func New_PDU_ReadResponse(functionCode byte, data []byte) *PDU_ReadResponse {
	return &PDU_ReadResponse{
		FunctionCode: functionCode,
		ByteCount:    byte(len(data)),
		Data:         data,
	}
}

func (pdu *PDU_ReadResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)

	// Write the PDU fields to the buffer
	buffer.WriteByte(pdu.FunctionCode)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.Data)

	return buffer.Bytes()
}

type PDU_ReadRegistersResponse struct {
	FunctionCode byte
	ByteCount    byte
	Data         []byte
}

func New_PDU_ReadRegistersResponse(functionCode byte, data []byte) *PDU_ReadRegistersResponse {
	return &PDU_ReadRegistersResponse{
		FunctionCode: functionCode,
		ByteCount:    byte(len(data)),
		Data:         data,
	}
}

func (pdu *PDU_ReadRegistersResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)

	// Write the PDU fields to the buffer
	buffer.WriteByte(pdu.FunctionCode)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.Data)

	return buffer.Bytes()
}

type PDU_WriteSingleResponse struct {
	FunctionCode    byte
	OutputAddress   uint16
	OutputValue     uint16
}

func New_PDU_WriteSingleResponse(functionCode byte, outputAddress, outputValue uint16) *PDU_WriteSingleResponse {
	return &PDU_WriteSingleResponse{
		FunctionCode:  functionCode,
		OutputAddress: outputAddress,
		OutputValue:   outputValue,
	}
}

func (pdu *PDU_WriteSingleResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)

	// Write the PDU fields to the buffer
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.OutputAddress)
	binary.Write(buffer, binary.BigEndian, pdu.OutputValue)

	return buffer.Bytes()
}

type PDU_WriteMultipleResponse struct {
	FunctionCode    byte
	StartingAddress uint16
	Quantity        uint16
}

func New_PDU_WriteMultipleResponse(functionCode byte, startingAddress, quantity uint16) *PDU_WriteMultipleResponse {
	return &PDU_WriteMultipleResponse{
		FunctionCode:    functionCode,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

func (pdu *PDU_WriteMultipleResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)

	// Write the PDU fields to the buffer
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.Quantity)

	return buffer.Bytes()
}
