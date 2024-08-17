package tests

import (
	"testing"
	"time"

	"github.com/tarm/serial"
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/server"
)

// Start a simple Modbus server with legal addresses set
func startModbusSerialServer(t *testing.T) *server.SerialServer {
	s := server.SerialServer{
		SerialConfig: &serial.Config{
			Name:        "/dev/ttyUSB0",
			Baud:        9600,
			ReadTimeout: 1 * time.Second,
		},
		ByteOrder: gomodbus.BigEndian,
		WordOrder: gomodbus.BigEndian,
		Slaves:    make(map[byte]*server.Slave),
	}

	// Add a slave
	s.AddSlave(1)

	// Set legal addresses and initialize discrete inputs for the slave
	slave := s.Slaves[1]
	for i := 0; i < 65535; i++ {
		// Initialize some discrete inputs for testing
		if i%2 == 0 {
			slave.DiscreteInputs[uint16(i)] = true
		}
	}

	go func() {
		if err := s.Start(); err != nil {
			t.Errorf("failed to start Modbus server: %v", err)
		}
	}()
	time.Sleep(2 * time.Second) // Give the server some time to start

	return &s
}

func TestStartSerialServer(t *testing.T) {
	startModbusSerialServer(t)
}