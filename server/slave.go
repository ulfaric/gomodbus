package server

import "sync"

type Slave struct {
	Coils            map[uint16]bool
	DiscreteInputs   map[uint16]bool
	HoldingRegisters map[uint16][]byte
	InputRegisters   map[uint16][]byte
	mu               sync.Mutex
}

// AddCoils adds a new coil to the slave.
func (s *Slave) AddCoils(address uint16, values []bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		s.Coils[address+uint16(i)] = value
	}
}

// SetCoils sets the coils to the slave.
func (s *Slave) SetCoils(address uint16, values []bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		if _, ok := s.Coils[address+uint16(i)]; ok {
			s.Coils[address+uint16(i)] = value
		}
	}
}

// DeleteCoils deletes the coils from the slave.
func (s *Slave) DeleteCoils(addresses []uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, address := range addresses {
		delete(s.Coils, address)
	}
}

// AddDiscreteInputs adds a new discrete input to the slave.
func (s *Slave) AddDiscreteInputs(address uint16, values []bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		s.DiscreteInputs[address+uint16(i)] = value
	}
}

// SetDiscreteInputs sets the discrete inputs to the slave.
func (s *Slave) SetDiscreteInputs(address uint16, values []bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		if _, ok := s.DiscreteInputs[address+uint16(i)]; ok {
			s.DiscreteInputs[address+uint16(i)] = value
		}
	}
}

// DeleteDiscreteInputs deletes the discrete inputs from the slave.
func (s *Slave) DeleteDiscreteInputs(addresses []uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, address := range addresses {
		delete(s.DiscreteInputs, address)
	}
}

// AddHoldingRegisters adds a new holding register to the slave.
func (s *Slave) AddHoldingRegisters(address uint16, values [][]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		s.HoldingRegisters[address+uint16(i)] = value
	}
}

// SetHoldingRegisters sets the holding registers to the slave.
func (s *Slave) SetHoldingRegisters(address uint16, values [][]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		if _, ok := s.HoldingRegisters[address+uint16(i)]; ok {
			s.HoldingRegisters[address+uint16(i)] = value
		}
	}
}

// DeleteHoldingRegisters deletes the holding registers from the slave.
func (s *Slave) DeleteHoldingRegisters(addresses []uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, address := range addresses {
		delete(s.HoldingRegisters, address)
	}
}

// AddInputRegisters adds a new input register to the slave.
func (s *Slave) AddInputRegisters(address uint16, values [][]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		s.InputRegisters[address+uint16(i)] = value
	}
}

// SetInputRegisters sets the input registers to the slave.
func (s *Slave) SetInputRegisters(address uint16, values [][]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		if _, ok := s.InputRegisters[address+uint16(i)]; ok {
			s.InputRegisters[address+uint16(i)] = value
		}
	}
}

// DeleteInputRegisters deletes the input registers from the slave.
func (s *Slave) DeleteInputRegisters(addresses []uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, address := range addresses {
		delete(s.InputRegisters, address)
	}
}
