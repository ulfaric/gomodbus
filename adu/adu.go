package adu

import (
	"bytes"
	"encoding/binary"
)

type TCP_ADU struct {
	TransactionID uint16
	ProtocolID    uint16
	Length        uint16
	UnitID        byte
	PDU           []byte
}

func New_TCP_ADU(transactionID uint16, unitID byte, pdu []byte) *TCP_ADU {
	return &TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(1 + len(pdu)),
		UnitID:        unitID,
		PDU:           pdu,
	}
}

func (adu *TCP_ADU) ToBytes() []byte {
	buffer := new(bytes.Buffer)

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

func (adu *TCP_ADU) FromBytes(data []byte) {
	buffer := bytes.NewBuffer(data)

	binary.Read(buffer, binary.BigEndian, &adu.TransactionID)
	binary.Read(buffer, binary.BigEndian, &adu.ProtocolID)
	binary.Read(buffer, binary.BigEndian, &adu.Length)
	adu.UnitID, _ = buffer.ReadByte()
	adu.PDU = buffer.Bytes()
}
