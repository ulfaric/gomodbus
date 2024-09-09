package server

import (
	"fmt"
	"sync"

	"github.com/tarm/serial"
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
	"github.com/ulfaric/gomodbus/pdu"
	"go.uber.org/zap"
	"context"
)

type SerialServer struct {
	Port      string
	BaudRate  int
	DataBits  byte
	Parity    serial.Parity
	StopBits  serial.StopBits
	ByteOrder string
	WordOrder string
	Slaves    map[byte]*Slave

	mu        sync.Mutex
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
}

func NewSerialServer(port string, baudRate int, dataBits byte, parity serial.Parity, stopBits serial.StopBits, byteOrder, wordOrder string) Server {
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

		ctx:       ctx,
		cancel:    cancel,
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
		Parity:   s.Parity,
		StopBits: s.StopBits,
	}

	port, err := serial.OpenPort(config)
	if err != nil {
		gomodbus.Logger.Error("failed to open serial port", zap.Error(err))
		return fmt.Errorf("failed to open serial port: %v", err)
	}
	defer port.Close()

	gomodbus.Logger.Sugar().Infof("Modbus serial server started on %s", s.Port)

	gomodbus.Logger.Info("Waiting for requests...")

	buffer := make([]byte, 256)
	for {
		select {
		case <-s.ctx.Done():
			return nil
		default:
			n, err := port.Read(buffer)
			if err != nil {
				gomodbus.Logger.Error("failed to read from serial port", zap.Error(err))
				continue
			}

			requestADU := &adu.SerialADU{}
			err = requestADU.FromBytes(buffer[:n])
			if err != nil {
				gomodbus.Logger.Error("failed to parse request", zap.Error(err))
				continue
			}

			slave, ok := s.Slaves[requestADU.UnitID]
			if !ok {
				gomodbus.Logger.Error("slave not found", zap.Uint8("unitID", requestADU.UnitID))
				responsePDU := pdu.NewPDUErrorResponse(requestADU.PDU[0], 0x04)
				response := adu.NewSerialADU(requestADU.UnitID, responsePDU.ToBytes())
				_, err = port.Write(response.ToBytes())
				if err != nil {
					gomodbus.Logger.Error("failed to write response", zap.Error(err))
				}
				continue
			}

			responsePDU, _ := processRequest(requestADU.PDU, slave)

			responseADU := adu.NewSerialADU(requestADU.UnitID, responsePDU)
			_, err = port.Write(responseADU.ToBytes())
			if err != nil {
				gomodbus.Logger.Error("failed to write response", zap.Error(err))
			}
		}
	}
}

func (s *SerialServer) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}
