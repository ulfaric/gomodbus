package server

import (
	"testing"
	"time"

	"gomodbus/client"
	"gomodbus"
)

// Start a simple Modbus server with legal addresses set
func startModbusServer(t *testing.T) *Server {
	s := Server{
		Host:      "127.0.0.1",
		Port:      502,
		ByteOrder: gomodbus.BigEndian,
		WordOrder: gomodbus.BigEndian,
		Slaves:    make(map[byte]*Slave),
	}

	// Add a slave
	s.AddSlave(1)

	// Set legal addresses and initialize discrete inputs for the slave
	slave := s.Slaves[1]
	for i := 0; i < 65535; i++ {
		slave.LegalCoilsAddress[i] = true
		slave.LegalDiscreteInputsAddress[i] = true
		slave.LegalHoldingRegistersAddress[i] = true
		slave.LegalInputRegistersAddress[i] = true

		// Initialize some discrete inputs for testing
		if i%2 == 0 {
			slave.DiscreteInputs[i] = true
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

func createModbusClient(t *testing.T) *client.TCPClient {
	modbusClient := client.TCPClient{
		Host: "127.0.0.1",
		Port: 502,
	}

	// Connect to the server
	err := modbusClient.Connect()
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	t.Cleanup(func() {
		modbusClient.Close()
	})

	return &modbusClient
}

func TestReadCoils(t *testing.T) {
	startModbusServer(t)
	modbusClient := createModbusClient(t)

	transactionID := 1
	startingAddress := 0
	quantity := 10
	unitID := 1

	coils, err := modbusClient.ReadCoils(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}
	t.Logf("Coils: %v", coils)
}

func TestReadDiscreteInputs(t *testing.T) {
	startModbusServer(t)
	modbusClient := createModbusClient(t)

	transactionID := 1
	startingAddress := 0
	quantity := 10
	unitID := 1

	discreteInputs, err := modbusClient.ReadDiscreteInputs(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read discrete inputs: %v", err)
	}
	t.Logf("Discrete Inputs: %v", discreteInputs)
}

func TestReadHoldingRegisters(t *testing.T) {
	startModbusServer(t)
	modbusClient := createModbusClient(t)

	transactionID := 1
	startingAddress := 0
	quantity := 10
	unitID := 1

	holdingRegisters, err := modbusClient.ReadHoldingRegisters(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}
	t.Logf("Holding Registers: %v", holdingRegisters)
}

func TestReadInputRegisters(t *testing.T) {
	startModbusServer(t)
	modbusClient := createModbusClient(t)

	transactionID := 1
	startingAddress := 0
	quantity := 10
	unitID := 1

	inputRegisters, err := modbusClient.ReadInputRegisters(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read input registers: %v", err)
	}
	t.Logf("Input Registers: %v", inputRegisters)
}

func TestWriteSingleCoil(t *testing.T) {
	startModbusServer(t)
	modbusClient := createModbusClient(t)

	transactionID := 1
	startingAddress := 0
	unitID := 1

	// Write single coil
	err := modbusClient.WriteSingleCoil(transactionID, startingAddress, unitID, true)
	if err != nil {
		t.Fatalf("Failed to write single coil: %v", err)
	}
	t.Log("Single Coil Written")

	// Read back the coil to verify the write
	coils, err := modbusClient.ReadCoils(transactionID, startingAddress, 1, unitID)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}
	if len(coils) != 1 || !coils[0] {
		t.Fatalf("Verification failed: expected coil at address %d to be true, got %v", startingAddress, coils)
	}
	t.Log("Single Coil Write Verified")
}

func TestWriteSingleRegister(t *testing.T) {
	startModbusServer(t)
	modbusClient := createModbusClient(t)

	transactionID := 1
	startingAddress := 0
	unitID := 1
	value := uint16(12345)

	// Write single register
	err := modbusClient.WriteSingleRegister(transactionID, startingAddress, unitID, value)
	if err != nil {
		t.Fatalf("Failed to write single register: %v", err)
	}
	t.Log("Single Register Written")

	// Read back the register to verify the write
	registers, err := modbusClient.ReadHoldingRegisters(transactionID, startingAddress, 1, unitID)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}
	if len(registers) != 1 || registers[0] != value {
		t.Fatalf("Verification failed: expected register at address %d to be %d, got %v", startingAddress, value, registers)
	}
	t.Log("Single Register Write Verified")
}

func TestWriteMultipleCoils(t *testing.T) {
	startModbusServer(t)
	modbusClient := createModbusClient(t)

	transactionID := 1
	startingAddress := 0
	unitID := 1
	values := []bool{true, false, true, false, true, false, true, false, true, false}

	// Write multiple coils
	err := modbusClient.WriteMultipleCoils(transactionID, startingAddress, unitID, values)
	if err != nil {
		t.Fatalf("Failed to write multiple coils: %v", err)
	}
	t.Log("Multiple Coils Written")

	// Read back the coils to verify the write
	coils, err := modbusClient.ReadCoils(transactionID, startingAddress, len(values), unitID)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}
	if len(coils) != len(values) {
		t.Fatalf("Verification failed: expected %d coils, got %d", len(values), len(coils))
	}
	for i, v := range values {
		if coils[i] != v {
			t.Fatalf("Verification failed: expected coil at address %d to be %v, got %v", startingAddress+i, v, coils[i])
		}
	}
	t.Log("Multiple Coils Write Verified")
}

func TestWriteMultipleRegisters(t *testing.T) {
	startModbusServer(t)
	modbusClient := createModbusClient(t)

	transactionID := 1
	startingAddress := 0
	unitID := 1
	registerValues := []uint16{1111, 2222, 3333, 4444, 5555, 6666, 7777, 8888, 9999, 1010}

	// Write multiple registers
	err := modbusClient.WriteMultipleRegisters(transactionID, startingAddress, unitID, registerValues)
	if err != nil {
		t.Fatalf("Failed to write multiple registers: %v", err)
	}
	t.Log("Multiple Registers Written")

	// Read back the registers to verify the write
	registers, err := modbusClient.ReadHoldingRegisters(transactionID, startingAddress, len(registerValues), unitID)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}
	if len(registers) != len(registerValues) {
		t.Fatalf("Verification failed: expected %d registers, got %d", len(registerValues), len(registers))
	}
	for i, v := range registerValues {
		if registers[i] != v {
			t.Fatalf("Verification failed: expected register at address %d to be %d, got %d", startingAddress+i, v, registers[i])
		}
	}
	t.Log("Multiple Registers Write Verified")
}
