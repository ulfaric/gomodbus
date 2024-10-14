package gomodbus

import (
	"fmt"
	"sync"

	"context"
	"github.com/tarm/serial"
)

type SerialServer struct {
	Port      string
	BaudRate  int
	DataBits  byte
	Parity    byte
	StopBits  byte
	ByteOrder string
	WordOrder string
	Slaves    map[byte]*Slave

	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewSerialServer(port string, baudRate int, dataBits byte, parity byte, stopBits byte, byteOrder, wordOrder string) Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &SerialServer{
		Port:      port,
		BaudRate:  baudRate,
		DataBits:  dataBits,
		Parity:    parity,
		StopBits:  stopBits,
		ByteOrder: byteOrder,
		WordOrder: wordOrder,
		Slaves:    make(map[byte]*Slave),

		ctx:    ctx,
		cancel: cancel,
	}
}

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

func (s *SerialServer) GetSlave(unitID byte) (*Slave, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slave, ok := s.Slaves[unitID]
	if !ok {
		return nil, fmt.Errorf("slave not found")
	}
	return slave, nil
}

func (s *SerialServer) RemoveSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Slaves, unitID)
}

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
		Logger.Sugar().Errorf("failed to open serial port: %v", err)
		return fmt.Errorf("failed to open serial port: %v", err)
	}
	defer port.Close()

	Logger.Sugar().Infof("Modbus serial server started on %s", s.Port)
	s.wg.Add(1)
	go s.handleRequest(port)
	return nil
}

func (s *SerialServer) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

func (s *SerialServer) handleRequest(port *serial.Port) error {
	defer s.wg.Done()
	Logger.Sugar().Info("Waiting for requests...")
	buffer := make([]byte, 256)
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			n, err := port.Read(buffer)
			if err != nil {
				Logger.Sugar().Errorf("failed to read from serial port: %v", err)
				continue
			}
			Logger.Sugar().Debugf("Server received request: %v", buffer[:n])

			requestADU := &SerialADU{}
			err = requestADU.FromBytes(buffer[:n])
			if err != nil {
				Logger.Sugar().Errorf("failed to parse request: %v", err)
				continue
			}

			slave, ok := s.Slaves[requestADU.UnitID]
			if !ok {
				Logger.Sugar().Errorf("slave not found: %v", requestADU.UnitID)
				responsePDU := NewPDUErrorResponse(requestADU.PDU[0], 0x04)
				response := NewSerialADU(requestADU.UnitID, responsePDU.ToBytes())
				_, err = port.Write(response.ToBytes())
				if err != nil {
					Logger.Sugar().Errorf("failed to write response: %v", err)
				}
				continue
			}

			responsePDU, _ := processRequest(requestADU.PDU, slave)

			responseADU := NewSerialADU(requestADU.UnitID, responsePDU)
			_, err = port.Write(responseADU.ToBytes())
			if err != nil {
				Logger.Sugar().Errorf("failed to write response: %v", err)
			}
		}
	}
}
