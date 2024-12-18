package server

import (
	"fmt"
	"sync"

	"context"
	"github.com/tarm/serial"
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
	"github.com/ulfaric/gomodbus/pdu"
)

type SerialServer struct {
	Port      string
	BaudRate  int
	DataBits  byte
	Parity    byte
	StopBits  byte
	Slaves    map[byte]*Slave

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewSerialServer creates a new serial server.
func NewSerialServer(port string, baudRate int, dataBits byte, parity byte, stopBits byte) Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &SerialServer{
		Port:      port,
		BaudRate:  baudRate,
		DataBits:  dataBits,
		Parity:    parity,
		StopBits:  stopBits,
		Slaves:    make(map[byte]*Slave),

		ctx:    ctx,
		cancel: cancel,
	}
}

// AddSlave adds a new slave to the server.
func (s *SerialServer) AddSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Slaves[unitID] = &Slave{
		Coils:            make(map[uint16]bool),
		DiscreteInputs:   make(map[uint16]bool),
		HoldingRegisters: make(map[uint16][]byte),
		InputRegisters:   make(map[uint16][]byte),
	}
}

// GetSlave gets a slave by unit ID.
func (s *SerialServer) GetSlave(unitID byte) (*Slave, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slave, ok := s.Slaves[unitID]
	if !ok {
		return nil, fmt.Errorf("slave not found")
	}
	return slave, nil
}

// RemoveSlave removes a slave by unit ID.
func (s *SerialServer) RemoveSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Slaves, unitID)
}

// Start starts the serial server.
func (s *SerialServer) Start() error {
	config := &serial.Config{
		Name:     s.Port,
		Baud:     s.BaudRate,
		Size:     s.DataBits,
		Parity:   serial.Parity(s.Parity),
		StopBits: serial.StopBits(s.StopBits),
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to open serial port: %v", err)
		return fmt.Errorf("failed to open serial port: %v", err)
	}
	defer port.Close()

	gomodbus.Logger.Sugar().Infof("Modbus serial server started on %s", s.Port)
	s.wg.Add(1)
	go s.handleRequest(port)
	return nil
}

// Stop stops the serial server.
func (s *SerialServer) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

// handleRequest handles the request from the serial port.
func (s *SerialServer) handleRequest(port *serial.Port) error {
	defer s.wg.Done()
	gomodbus.Logger.Sugar().Info("Waiting for requests...")
	buffer := make([]byte, 256)
	for {
		select {
		case <-s.ctx.Done():
			gomodbus.Logger.Sugar().Info("Shutting down server...")
			return nil
		default:
			n, err := port.Read(buffer)
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to read from serial port: %v", err)
				continue
			}
			gomodbus.Logger.Sugar().Debugf("Server received request: %v", buffer[:n])

			requestADU := &adu.SerialADU{}
			err = requestADU.FromBytes(buffer[:n])
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
				continue
			}

			slave, ok := s.Slaves[requestADU.UnitID]
			if !ok {
				gomodbus.Logger.Sugar().Errorf("slave not found: %v", requestADU.UnitID)
				responsePDU := pdu.NewPDUErrorResponse(requestADU.PDU[0], 0x04)
				response := adu.NewSerialADU(requestADU.UnitID, responsePDU.ToBytes())
				_, err = port.Write(response.ToBytes())
				if err != nil {
					gomodbus.Logger.Sugar().Errorf("failed to write response: %v", err)
				}
				continue
			}

			responsePDU, _ := processRequest(requestADU.PDU, slave)

			responseADU := adu.NewSerialADU(requestADU.UnitID, responsePDU)
			_, err = port.Write(responseADU.ToBytes())
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to write response: %v", err)
			}
		}
	}
}
