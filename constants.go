package gomodbus

const (
	ReadCoil               byte = 0x01
	ReadDiscreteInput      byte = 0x02
	ReadHoldingRegister    byte = 0x03
	ReadInputRegister      byte = 0x04
	WriteSingleCoil        byte = 0x05
	WriteSingleRegister    byte = 0x06
	WriteMultipleCoils     byte = 0x0F
	WriteMultipleRegisters byte = 0x10
)
