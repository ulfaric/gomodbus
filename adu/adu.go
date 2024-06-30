package adu

import (
	"bytes"
	"encoding/binary"
	"gomodbus"
	"gomodbus/pdu"
)

type TCP_ADU struct {
	TransactionID uint16
	ProtocolID    uint16
	Length        uint16
	UnitID        byte
	PDU           pdu.PDU
}

type RTU_ADU struct {
	DeviceAddress byte
	PDU           pdu.PDU
	CRC           uint16
}

func ReadCoilsTCPADU(transactionID, startingAddress, quantity uint16, unitID byte) TCP_ADU {
	pdu := pdu.ReadCoilsPDU(startingAddress, quantity)
	adu := TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        6,
		UnitID:        unitID,
		PDU:           pdu,
	}
	return adu
}

func (adu *TCP_ADU) ToBytes() []byte {
	buffer := new(bytes.Buffer)

	// Write TransactionID, ProtocolID, and Length as uint16 to the buffer
	binary.Write(buffer, binary.BigEndian, adu.TransactionID)
	binary.Write(buffer, binary.BigEndian, adu.ProtocolID)
	binary.Write(buffer, binary.BigEndian, adu.Length)

	// Write UnitID and FunctionCode as bytes
	buffer.WriteByte(adu.UnitID)
	buffer.WriteByte(byte(adu.PDU.FunctionCode))

	// Write each data element as uint16
	for _, data := range adu.PDU.Data {
		binary.Write(buffer, binary.BigEndian, data)
	}

	return buffer.Bytes()
}

func (adu *TCP_ADU) FromBytes(adu_bytes []byte) {
	adu.TransactionID = binary.BigEndian.Uint16(adu_bytes[0:2])
	adu.ProtocolID = binary.BigEndian.Uint16(adu_bytes[2:4])
	adu.Length = binary.BigEndian.Uint16(adu_bytes[4:6])
	adu.UnitID = adu_bytes[6]
	adu.PDU.FunctionCode = gomodbus.FunctionCode(adu_bytes[7])

	// Calculate the correct number of uint16 elements in the Data slice
	dataLength := (len(adu_bytes) - 8) / 2
	adu.PDU.Data = make([]uint16, dataLength)

	// Correctly iterate over the byte slice to fill the Data slice
	for i := 0; i < dataLength; i++ {
		adu.PDU.Data[i] = binary.BigEndian.Uint16(adu_bytes[i*2+8 : i*2+10])
	}
}
