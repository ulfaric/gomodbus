package server

type Slave struct {
	Coils                        [65535]bool
	LegalCoilsAddress            [65535]bool
	DiscreteInputs               [65535]bool
	LegalDiscreteInputsAddress   [65535]bool
	HoldingRegisters             [65535]uint16
	LegalHoldingRegistersAddress [65535]bool
	InputRegisters               [65535]uint16
	LegalInputRegistersAddress   [65535]bool
}