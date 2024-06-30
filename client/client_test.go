package client

import (
	"testing"
)


func TestReadCoils(t *testing.T) {
	client := TCPClient{
		Host: "192.168.5.3",
		Port: 1502,
	}

	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Test case 1: Read 1 coil starting from address 0
	transactionID := uint16(10)
	quantity := uint16(3)
	startingAddress := uint16(2)
	unitID := byte(0)
	response, err := client.ReadCoils(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}
	t.Logf("Response: %v", response)
}