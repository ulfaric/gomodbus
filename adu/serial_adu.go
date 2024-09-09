package adu

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/ulfaric/gomodbus"
)

type SerialADU struct {
	UnitID byte
	PDU    []byte
	CRC    uint16
}

func NewSerialADU(address byte, pdu []byte) *SerialADU {
	adu := &SerialADU{
		UnitID: address,
		PDU:    pdu,
	}

	// Calculate CRC16 checksum for the ADU (Address + PDU)
	adu.CRC = calculateCRC16(adu.UnitID, adu.PDU)

	return adu
}

func (adu *SerialADU) ToBytes() []byte {
	buffer := new(bytes.Buffer)

	// Write Address
	buffer.WriteByte(adu.UnitID)

	// Write PDU
	buffer.Write(adu.PDU)

	// Write CRC16 as two bytes in Little Endian format
	binary.Write(buffer, binary.LittleEndian, adu.CRC)

	return buffer.Bytes()
}

func (adu *SerialADU) FromBytes(data []byte) error {

	buffer := bytes.NewBuffer(data)

	// Read Address
	err := binary.Read(buffer, binary.BigEndian, &adu.UnitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse UnitID for SerialADU: %v", err)
		return err
	}

	// Read PDU (all bytes except the last two, which are CRC)
	adu.PDU = buffer.Next(len(data) - 3) // Exclude UnitID and CRC

	// Read and validate CRC
	var crc uint16
	binary.Read(buffer, binary.LittleEndian, &crc)

	expectedCRC := calculateCRC16(adu.UnitID, adu.PDU)
	if crc != expectedCRC {
		return fmt.Errorf("CRC mismatch: expected 0x%X, got 0x%X", expectedCRC, crc)
	}

	adu.CRC = crc
	return nil
}

// calculateCRC16 calculates the CRC16 for the given data (Modbus standard)
func calculateCRC16(address byte, pdu []byte) uint16 {
	crc := uint16(0xFFFF) // Initialize the CRC with 0xFFFF

	// Process the Address byte
	crc = updateCRC16(crc, address)

	// Process each byte in the PDU
	for _, b := range pdu {
		crc = updateCRC16(crc, b)
	}

	return crc
}

// updateCRC16 updates the CRC16 value for each byte of the data
func updateCRC16(crc uint16, b byte) uint16 {
	crc ^= uint16(b) // XOR byte with the CRC

	for i := 0; i < 8; i++ { // Process each bit
		if crc&0x0001 != 0 {
			crc = (crc >> 1) ^ 0xA001 // Shift right and XOR with the polynomial
		} else {
			crc >>= 1 // Just shift right if the least significant bit is not set
		}
	}

	return crc
}
