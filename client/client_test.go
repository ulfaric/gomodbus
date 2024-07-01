package client

import (
	// "gomodbus/pdu"
	// "gomodbus/adu"
	// "gomodbus"
	"reflect"
	"testing"
	"math"
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
	expected := []bool{false, false, false}
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
	expected := []bool{false, false, false}
	if !reflect.DeepEqual(coils, expected) {
		t.Errorf("Expected %v, but got %v", expected, coils)
	}
}

func TestReadHoldingRegisters(t *testing.T) {
	client := TCPClient{
		Host: "192.168.5.3",
		Port: 1502,
	}

	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Test case 1: Read 3 holding registers starting from address 0
	transactionID := 10
	quantity := 2
	startingAddress := 0
	unitID := 0
	registers, err := client.ReadHoldingRegisters(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}
	t.Logf("Response: %v", registers)

	value := DecodeModbusRegisters(registers).(uint32)
	result := math.Float32frombits(value)
	expected := float32(2.2)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	t.Logf("Value: %v", result)
}

func TestReadInputRegisters(t *testing.T) {
	client := TCPClient{
		Host: "192.168.5.3",
		Port: 1502,
	}

	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Test case 1: Read 3 input registers starting from address 0
	transactionID := 10
	quantity := 4
	startingAddress := 0
	unitID := 0
	registers, err := client.ReadInputRegisters(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read input registers: %v", err)
	}
	t.Logf("Response: %v", registers)
	value := DecodeModbusRegisters(registers).(uint64)
	result := int64(value)
	expected := int64(-10000)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
	t.Logf("Value: %v", result)
}
