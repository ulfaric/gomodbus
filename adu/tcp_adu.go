package adu

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ulfaric/gomodbus"
)

type TCPADU struct {
	TransactionID uint16
	ProtocolID    uint16
	Length        uint16
	UnitID        byte
	PDU           []byte
}

func NewTCPADU(transactionID uint16, unitID byte, pdu []byte) *TCPADU {
	return &TCPADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(1 + len(pdu)),
		UnitID:        unitID,
		PDU:           pdu,
	}
}

func (adu *TCPADU) ToBytes() []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 7+len(adu.PDU)))

	// Write TransactionID, ProtocolID, and Length as uint16 to the buffer
	binary.Write(buffer, binary.BigEndian, adu.TransactionID)
	binary.Write(buffer, binary.BigEndian, adu.ProtocolID)
	binary.Write(buffer, binary.BigEndian, adu.Length)

	// Write UnitID
	buffer.WriteByte(adu.UnitID)

	// Write PDU
	buffer.Write(adu.PDU)

	return buffer.Bytes()
}

func (adu *TCPADU) FromBytes(data []byte) error {
	// Check minimum length requirement (2 bytes TransactionID + 2 bytes ProtocolID + 2 bytes Length + 1 byte UnitID + at least 2 bytes PDU)
	if len(data) < 9 {
		gomodbus.Logger.Sugar().Errorf("insufficient data length for TCPADU: got %d bytes, minimum required is 9", len(data))
		return fmt.Errorf("insufficient data length: got %d bytes, minimum required is 9", len(data))
	}

	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, &adu.TransactionID); err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse TransactionID for TCPADU: %v", err)
		return err
	}

	if err := binary.Read(buffer, binary.BigEndian, &adu.ProtocolID); err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse ProtocolID for TCPADU: %v", err)
		return err
	}

	if err := binary.Read(buffer, binary.BigEndian, &adu.Length); err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse Length for TCPADU: %v", err)
		return err
	}

	unitID, err := buffer.ReadByte()
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse UnitID for TCPADU: %v", err)
		return err
	}
	adu.UnitID = unitID

	// Ensure the remaining bytes match the expected PDU length
	if int(adu.Length-1) != buffer.Len() {
		gomodbus.Logger.Sugar().Errorf("invalid PDU length, expected %d, got %d", adu.Length-1, buffer.Len())
		return fmt.Errorf("invalid PDU length, expected %d, got %d", adu.Length-1, buffer.Len())
	}

	adu.PDU = buffer.Bytes()
	return nil
}
