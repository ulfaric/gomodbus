package pdu

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type PDU_ReadResponse struct {
	FunctionCode byte
	ByteCount    byte
	Data         []byte
}

func NewPDUReadResponse(functionCode byte, data []byte) *PDU_ReadResponse {
	return &PDU_ReadResponse{
		FunctionCode: functionCode,
		ByteCount:    byte(len(data)),
		Data:         data,
	}
}

func (pdu *PDU_ReadResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	buffer.WriteByte(pdu.ByteCount)
	buffer.Write(pdu.Data)
	return buffer.Bytes()
}

func (pdu *PDU_ReadResponse) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var err error

	pdu.FunctionCode, err = buffer.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read FunctionCode: %v", err)
	}

	pdu.ByteCount, err = buffer.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read ByteCount: %v", err)
	}

	pdu.Data = buffer.Bytes()
	if len(pdu.Data) != int(pdu.ByteCount) {
		return fmt.Errorf("data length mismatch, expected %d, got %d", pdu.ByteCount, len(pdu.Data))
	}

	return nil
}

type PDU_WriteSingleResponse struct {
	FunctionCode  byte
	OutputAddress uint16
	OutputValue   uint16
}

func NewPDUWriteSingleResponse(functionCode byte, outputAddress, outputValue uint16) *PDU_WriteSingleResponse {
	return &PDU_WriteSingleResponse{
		FunctionCode:  functionCode,
		OutputAddress: outputAddress,
		OutputValue:   outputValue,
	}
}

func (pdu *PDU_WriteSingleResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.OutputAddress)
	binary.Write(buffer, binary.BigEndian, pdu.OutputValue)
	return buffer.Bytes()
}

func (pdu *PDU_WriteSingleResponse) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var err error

	pdu.FunctionCode, err = buffer.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read FunctionCode: %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &pdu.OutputAddress)
	if err != nil {
		return fmt.Errorf("failed to read OutputAddress: %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &pdu.OutputValue)
	if err != nil {
		return fmt.Errorf("failed to read OutputValue: %v", err)
	}

	return nil
}

type PDU_WriteMultipleResponse struct {
	FunctionCode    byte
	StartingAddress uint16
	Quantity        uint16
}

func NewPDUWriteMultipleResponse(functionCode byte, startingAddress, quantity uint16) *PDU_WriteMultipleResponse {
	return &PDU_WriteMultipleResponse{
		FunctionCode:    functionCode,
		StartingAddress: startingAddress,
		Quantity:        quantity,
	}
}

func (pdu *PDU_WriteMultipleResponse) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(pdu.FunctionCode)
	binary.Write(buffer, binary.BigEndian, pdu.StartingAddress)
	binary.Write(buffer, binary.BigEndian, pdu.Quantity)
	return buffer.Bytes()
}

func (pdu *PDU_WriteMultipleResponse) FromBytes(data []byte) error {
	buffer := bytes.NewBuffer(data)
	var err error

	pdu.FunctionCode, err = buffer.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read FunctionCode: %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &pdu.StartingAddress)
	if err != nil {
		return fmt.Errorf("failed to read StartingAddress: %v", err)
	}

	err = binary.Read(buffer, binary.BigEndian, &pdu.Quantity)
	if err != nil {
		return fmt.Errorf("failed to read Quantity: %v", err)
	}

	return nil
}
