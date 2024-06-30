package pdu

import (
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
