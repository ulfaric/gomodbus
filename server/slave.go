package server

type Slave struct {
	Coils            map[uint16]bool
	DiscreteInputs   map[uint16]bool
	HoldingRegisters map[uint16][]byte
	InputRegisters   map[uint16][]byte
}

func (s *Slave) AddCoil(address uint16, value bool) {
	s.Coils[address] = value
}

func (s *Slave) AddCoils(address uint16, values []bool) {
	for i, value := range values {
		s.Coils[address+uint16(i)] = value
	}
}

func (s *Slave) AddDiscreteInput(address uint16, value bool) {
	s.DiscreteInputs[address] = value
}

func (s *Slave) AddDiscreteInputs(address uint16, values []bool) {
	for i, value := range values {
		s.DiscreteInputs[address+uint16(i)] = value
	}
}

func (s *Slave) AddHoldingRegister(address uint16, value []byte) {
	s.HoldingRegisters[address] = value
}

func (s *Slave) AddHoldingRegisters(address uint16, values [][]byte) {
	for i, value := range values {
		s.HoldingRegisters[address+uint16(i)] = value
	}
}

func (s *Slave) AddInputRegister(address uint16, value []byte) {
	s.InputRegisters[address] = value
}

func (s *Slave) AddInputRegisters(address uint16, values [][]byte) {
	for i, value := range values {
		s.InputRegisters[address+uint16(i)] = value
	}
}
