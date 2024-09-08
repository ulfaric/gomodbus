package server

type Slave struct {
	Coils                        map[uint16]bool
	DiscreteInputs               map[uint16]bool
	HoldingRegisters             map[uint16]uint16
	InputRegisters               map[uint16]uint16
}

