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
	LittleEndian           string = "little"
	BigEndian              string = "big"

	IllegalFunction                byte = 0x01
	IllegalDataAddress             byte = 0x02
	IllegalDataValue               byte = 0x03
	SlaveDeviceFailure             byte = 0x04
	Acknowledge                    byte = 0x05
	SlaveDeviceBusy                byte = 0x06
	MemoryParityError              byte = 0x08
	NegativeAcknowledge            byte = 0x07 // Added missing error code
	GatewayPathUnavailable         byte = 0x0A
	GatewayTargetDeviceFailedToRespond byte = 0x0B
)

// ModbusExceptionMap maps Modbus exception codes to their corresponding error messages.
var ModbusException = map[byte]string{
	0x01: "Illegal Function",
	0x02: "Illegal Data Address",
	0x03: "Illegal Data Value",
	0x04: "Slave Device Failure",
	0x05: "Acknowledge",
	0x06: "Slave Device Busy",
	0x07: "Negative Acknowledge", // Added missing error message
	0x08: "Memory Parity Error",
	0x0A: "Gateway Path Unavailable",
	0x0B: "Gateway Target Device Failed to Respond",
}
