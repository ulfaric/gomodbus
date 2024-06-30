package adu

import (
	"bytes"
	"encoding/binary"
	"gomodbus/pdu"
)

type TCP_ADU struct {
	TransactionID uint16
	ProtocolID    uint16
	Length        uint16
	UnitID        byte
	PDU           []byte
}

func ReadCoilsTCPADU(transactionID, startingAddress, quantity uint16, unitID byte) TCP_ADU {
	pdu := pdu.ReadCoilsPDU(startingAddress, quantity)
	pdu_bytes := pdu.ToBytes()
	lengh := len(pdu_bytes) + 1
	adu := TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(lengh),
		UnitID:        unitID,
		PDU:           pdu_bytes,
	}
	return adu
}

func WriteSingleCoilTCPADU(transactionID, address, value uint16, unitID byte) TCP_ADU {
	pdu := pdu.WriteSingleCoilPDU(address, value)
	pdu_bytes := pdu.ToBytes()
	lengh := len(pdu_bytes) + 1
	adu := TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(lengh),
		UnitID:        unitID,
		PDU:           pdu_bytes,
	}
	return adu
}

func WriteMultipleCoilsTCPADU(transactionID, startingAddress, quantity uint16, byteCount byte, values []uint16, unitID byte) TCP_ADU {
	pdu := pdu.WriteMultipleCoilsPDU(startingAddress, quantity, byteCount, values)
	pdu_bytes := pdu.ToBytes()
	lengh := len(pdu_bytes) + 1
	adu := TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(lengh),
		UnitID:        unitID,
		PDU:           pdu_bytes,
	}
	return adu
}

func ReadDiscreteInputsTCPADU(transactionID, startingAddress, quantity uint16, unitID byte) TCP_ADU {
	pdu := pdu.ReadDiscreteInputsPDU(startingAddress, quantity)
	pdu_bytes := pdu.ToBytes()
	lengh := len(pdu_bytes) + 1
	adu := TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(lengh),
		UnitID:        unitID,
		PDU:           pdu_bytes,
	}
	return adu
}

func ReadHoldingRegistersTCPADU(transactionID, startingAddress, quantity uint16, unitID byte) TCP_ADU {
	pdu := pdu.ReadHoldingRegistersPDU(startingAddress, quantity)
	pdu_bytes := pdu.ToBytes()
	lengh := len(pdu_bytes) + 1
	adu := TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(lengh),
		UnitID:        unitID,
		PDU:           pdu_bytes,
	}
	return adu
}

func ReadInputRegistersTCPADU(transactionID, startingAddress, quantity uint16, unitID byte) TCP_ADU {
	pdu := pdu.ReadInputRegistersPDU(startingAddress, quantity)
	pdu_bytes := pdu.ToBytes()
	lengh := len(pdu_bytes) + 1
	adu := TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(lengh),
		UnitID:        unitID,
		PDU:           pdu_bytes,
	}
	return adu
}

func WriteSingleRegisterTCPADU(transactionID, address, value uint16, unitID byte) TCP_ADU {
	pdu := pdu.WriteSingleRegisterPDU(address, value)
	pdu_bytes := pdu.ToBytes()
	lengh := len(pdu_bytes) + 1
	adu := TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(lengh),
		UnitID:        unitID,
		PDU:           pdu_bytes,
	}
	return adu
}

func WriteMultipleRegistersTCPADU(transactionID, startingAddress, quantity uint16, byteCount byte, values []uint16, unitID byte) TCP_ADU {
	pdu := pdu.WriteMultipleRegistersPDU(startingAddress, quantity, byteCount, values)
	pdu_bytes := pdu.ToBytes()
	lengh := len(pdu_bytes) + 1
	adu := TCP_ADU{
		TransactionID: transactionID,
		ProtocolID:    0,
		Length:        uint16(lengh),
		UnitID:        unitID,
		PDU:           pdu_bytes,
	}
	return adu
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
