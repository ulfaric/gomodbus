package client

import (
	"reflect"
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
	transactionID := 10
	quantity := 3
	startingAddress := 0
	unitID := 0
	coils, err := client.ReadCoils(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}
	t.Logf("Response: %v", coils)
	expected := []bool{true, true, true}
	if !reflect.DeepEqual(coils, expected) {
		t.Errorf("Expected %v, but got %v", expected, coils)
	}
}

func TestReadDiscreteInput(t *testing.T) {
	client := TCPClient{
		Host: "192.168.5.3",
		Port: 1502,
	}

	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Test case 1: Read 1 coil starting from address 0
	transactionID := 10
	quantity := 3
	startingAddress := 0
	unitID := 0
	coils, err := client.ReadDiscreteInputs(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read discrete inputs: %v", err)
	}
	t.Logf("Response: %v", coils)
	expected := []bool{true, true, true}
	if !reflect.DeepEqual(coils, expected) {
		t.Errorf("Expected %v, but got %v", expected, coils)
	}
}
