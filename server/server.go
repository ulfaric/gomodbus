package server

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/tarm/serial"
)

type Server interface {
	Start() error
	Stop() error
}

func NewTCPServerDev()Server{
	s := &TCPServerDev{
		server: server{
			Slaves: make(map[byte]*Slave),
		},
	}
	return s
}

type server struct {
	ByteOrder string
	WordOrder string
	Slaves    map[byte]*Slave
	conn      interface{}
	mu        sync.Mutex
}

func (s *server) AddSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Slaves[unitID] = &Slave{
		Coils:            make(map[uint16]bool),
		DiscreteInputs:   make(map[uint16]bool),
		HoldingRegisters: make(map[uint16]uint16),
		InputRegisters:   make(map[uint16]uint16),
	}
}

func (s *server) RemoveSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Slaves, unitID)
}

func (s *server) ReceiveRequest() ([]byte, error) {
	switch c := s.conn.(type) {
	case *serial.Port:
		data, err := io.ReadAll(c)
		return data, err
	case net.Conn:
		data, err := io.ReadAll(c)
		return data, err
	default:
		return nil, fmt.Errorf("unsupported connection type: %T", c)
	}
}

func (s *server) SendResponse(data []byte) error {
	switch c := s.conn.(type) {
	case *serial.Port:
		_, err := c.Write(data)
		return err
	case net.Conn:
		_, err := c.Write(data)
		return err
	default:
		return fmt.Errorf("unsupported connection type: %T", c)
	}
}


