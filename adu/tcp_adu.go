package adu

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
	buffer := bytes.NewBuffer(data)

	if err := binary.Read(buffer, binary.BigEndian, &adu.TransactionID); err != nil {
		return err
	}

	if err := binary.Read(buffer, binary.BigEndian, &adu.ProtocolID); err != nil {
		return err
	}

	if err := binary.Read(buffer, binary.BigEndian, &adu.Length); err != nil {
		return err
	}

	unitID, err := buffer.ReadByte()
	if err != nil {
		return err
	}
	adu.UnitID = unitID

	// Ensure the remaining bytes match the expected PDU length
	if int(adu.Length-1) != buffer.Len() {
		return fmt.Errorf("invalid PDU length")
	}

	adu.PDU = buffer.Bytes()
	return nil
}
