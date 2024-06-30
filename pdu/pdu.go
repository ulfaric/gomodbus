package pdu

import (
	"bytes"
	"encoding/binary"
	"gomodbus"
)

type PDU struct {
	FunctionCode gomodbus.FunctionCode
	Data         []uint16
}

func ReadCoilsPDU(startingAddress, quantity uint16) PDU {
	return PDU{
		FunctionCode: gomodbus.ReadCoil,
		Data:         []uint16{startingAddress, quantity},
	}
}

func ReadDiscreteInputsPDU(startingAddress, quantity uint16) PDU {
	return PDU{
		FunctionCode: gomodbus.ReadDiscreteInput,
		Data:         []uint16{startingAddress, quantity},
	}
}

func ReadHoldingRegistersPDU(startingAddress, quantity uint16) PDU {

	return PDU{
		FunctionCode: gomodbus.ReadHoldingRegister,
		Data:         []uint16{startingAddress, quantity},
	}
}

func ReadInputRegistersPDU(startingAddress, quantity uint16) PDU {
	return PDU{
		FunctionCode: gomodbus.ReadInputRegister,
		Data:         []uint16{startingAddress, quantity},
	}
}

func WriteSingleCoilPDU(address, value uint16) PDU {
	return PDU{
		FunctionCode: gomodbus.WriteSingleCoil,
		Data:         []uint16{address, value},
	}
}

func WriteSingleRegisterPDU(address, value uint16) PDU {
	return PDU{
		FunctionCode: gomodbus.WriteSingleRegister,
		Data:         []uint16{address, value},
	}
}

func WriteMultipleCoilsPDU(startingAddress, quantity uint16, byteCount byte, values []uint16) PDU {
	data := []uint16{startingAddress, quantity, uint16(byteCount)}
	data = append(data, values...)
	return PDU{
		FunctionCode: gomodbus.WriteMultipleCoils,
		Data:         data,
	}
}

func WriteMultipleRegistersPDU(startingAddress, quantity uint16, byteCount byte, values []uint16) PDU {
	data := []uint16{startingAddress, quantity, uint16(byteCount)}
	data = append(data, values...)
	return PDU{
		FunctionCode: gomodbus.WriteMultipleRegisters,
		Data:         data,
	}
}

func (pdu *PDU) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(pdu.FunctionCode))

	if pdu.FunctionCode == gomodbus.WriteMultipleCoils || pdu.FunctionCode == gomodbus.WriteMultipleRegisters {
		buffer.WriteByte(byte(pdu.Data[0]))
		for i := 1; i < len(pdu.Data); i++ {
			binary.Write(buffer, binary.BigEndian, pdu.Data[i])
		}
	} else {
		for _, data := range pdu.Data {
			binary.Write(buffer, binary.BigEndian, data)
		}
	}
	return buffer.Bytes()

}

