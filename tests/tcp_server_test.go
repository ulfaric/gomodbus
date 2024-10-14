package tests

import (
	"testing"
	"time"

	c "github.com/ulfaric/gomodbus/client"
	s "github.com/ulfaric/gomodbus/server"
)

func createTestTCPServer() s.Server {
	server := s.NewTCPServer("127.0.0.1", 1502, false, "big", "big", "", "", "")
	server.AddSlave(1)
	// Initialize coils
	slave, _ := server.GetSlave(1)
	for i := 0; i < 10; i++ {
		slave.Coils[uint16(i)] = false
		slave.DiscreteInputs[uint16(i)] = false
		slave.HoldingRegisters[uint16(i)] = []byte{0x00, 0x00}
		slave.InputRegisters[uint16(i)] = []byte{0x00, 0x00}
	}
	return server
}

func TestTCPServer_Coils(t *testing.T) {
	server := createTestTCPServer()
	go server.Start()
	defer server.Stop()

	time.Sleep(1 * time.Second)

	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	defer client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	// Write some coils
	err = c.WriteMultipleCoils(client, 1, 0, []bool{true, false, true})
	if err != nil {
		t.Fatalf("Failed to write coils: %v", err)
	}

	// Read the coils
	coils, err := c.ReadCoils(client, 1, 0, 3)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}

	expected := []bool{true, false, true}
	for i, coil := range coils {
		if coil != expected[i] {
			t.Errorf("Coil %d: expected %v, got %v", i, expected[i], coil)
		}
	}
}

func TestTCPServer_ReadHoldingRegisters(t *testing.T) {
	server := createTestTCPServer()
	go server.Start()
	defer server.Stop()

	time.Sleep(1 * time.Second)

	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Disconnect()

	// Write some registers
	err = c.WriteMultipleRegisters(client, 1, 0, 2, []byte{0x12, 0x34, 0x56, 0x78})
	if err != nil {
		t.Fatalf("Failed to write registers: %v", err)
	}

	// Read the registers
	registers, err := c.ReadHoldingRegisters(client, 1, 0, 2)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}

	expected := []uint16{0x1234, 0x5678}
	for i, reg := range registers {
		if reg != expected[i] {
			t.Errorf("Register %d: expected 0x%04X, got 0x%04X", i, expected[i], reg)
		}
	}
}

