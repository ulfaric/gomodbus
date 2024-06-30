package gomodbus

type FunctionCode byte

const (
	ReadCoil               FunctionCode = 0x01
	ReadDiscreteInput      FunctionCode = 0x02
	ReadHoldingRegister    FunctionCode = 0x03
	ReadInputRegister      FunctionCode = 0x04
	WriteSingleCoil        FunctionCode = 0x05
	WriteSingleRegister    FunctionCode = 0x06
	WriteMultipleCoils     FunctionCode = 0x0F
	WriteMultipleRegisters FunctionCode = 0x10
)
