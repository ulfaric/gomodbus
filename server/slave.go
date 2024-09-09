package server

import "sync"

type Slave struct {
	Coils            map[uint16]bool
	DiscreteInputs   map[uint16]bool
	HoldingRegisters map[uint16][]byte
	InputRegisters   map[uint16][]byte
	mu               sync.Mutex
}

func (s *Slave) AddCoils(address uint16, values []bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		s.Coils[address+uint16(i)] = value
	}
}

func (s *Slave) DeleteCoils(addresses []uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, address := range addresses {
		delete(s.Coils, address)
	}
}

func (s *Slave) AddDiscreteInputs(address uint16, values []bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		s.DiscreteInputs[address+uint16(i)] = value
	}
}

func (s *Slave) DeleteDiscreteInputs(addresses []uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, address := range addresses {
		delete(s.DiscreteInputs, address)
	}
}

func (s *Slave) AddHoldingRegisters(address uint16, values [][]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		s.HoldingRegisters[address+uint16(i)] = value
	}
}

func (s *Slave) DeleteHoldingRegisters(addresses []uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, address := range addresses {
		delete(s.HoldingRegisters, address)
	}
}

func (s *Slave) AddInputRegisters(address uint16, values [][]byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, value := range values {
		s.InputRegisters[address+uint16(i)] = value
	}
}

func (s *Slave) DeleteInputRegisters(addresses []uint16) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, address := range addresses {
		delete(s.InputRegisters, address)
	}
}