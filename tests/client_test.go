package tests

import (
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/client"
	"reflect"
	"testing"
	"math"
)

func TestClientReadCoils(t *testing.T) {
	c := client.TCPClient{
		Host: "192.168.5.3",
		Port: 1502,
	}

	err := c.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Test case 1: Read 1 coil starting from address 0
	transactionID := 10
	quantity := 1
	startingAddress := 0
	unitID := 0
	value := true
	// t.Logf("Writing single coil at address %v : %v", startingAddress, value)
	// err = c.WriteSingleCoil(transactionID, startingAddress, unitID, value)
	// if err != nil {
	// 	t.Fatalf("Failed to write single coil: %v", err)
	// }
	coil, err := c.ReadCoils(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}
	t.Logf("Coil at address %v : %v", startingAddress, coil)
	if !reflect.DeepEqual(coil[0], value) {
		t.Errorf("Expected %v, but got %v", value, coil)
	}

	// // Test case 2: Read 3 coils starting from address 0
	// quantity = 3
	// values := []bool{true, false, true}
	// t.Logf("Writing multiple coils at address %v - %v : %v", startingAddress, startingAddress+quantity, values)
	// err = c.WriteMultipleCoils(transactionID, startingAddress, unitID, values)
	// if err != nil {
	// 	t.Fatalf("Failed to write multiple coils: %v", err)
	// }
	// coils, err := c.ReadCoils(transactionID, startingAddress, quantity, unitID)
	// if err != nil {
	// 	t.Fatalf("Failed to read coils: %v", err)
	// }
	// if !reflect.DeepEqual(coils, values) {
	// 	t.Errorf("Expected %v, but got %v", values, coils)
	// }
	// t.Logf("Coils at address %v - %v : %v", startingAddress, startingAddress+quantity, coils)
}

func TestClientReadDiscreteInput(t *testing.T) {
	client := client.TCPClient{
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

func TestClientHoldingRegisters(t *testing.T) {
	c := client.TCPClient{
		Host: "192.168.5.3",
		Port: 1502,
	}

	err := c.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Test case 1: Write and read 1 holding registers starting from address 0
	transactionID := 10
	quantity := 1
	startingAddress := 0
	unitID := 0
	value := uint16(10000)
	t.Logf("Writing single register at address %v : %v", startingAddress, value)
	err = c.WriteSingleRegister(transactionID, startingAddress, unitID, value)
	if err != nil {
		t.Fatalf("Failed to write single register: %v", err)
	}
	registers, err := c.ReadHoldingRegisters(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}
	if !reflect.DeepEqual(registers[0], value) {
		t.Errorf("Expected %v, but got %v", value, registers)
	}
	t.Logf("Holding register at address %v : %v", startingAddress, registers[0])

	// Test case 2: Write and read 4 holding registers starting from address 0
	quantity = 4
	value_64 := int64(-10000)
	t.Logf("Wrting value: %v as 64 bits into 4 registers", value_64)
	values := gomodbus.EncodeModbusRegisters(value_64, gomodbus.BigEndian, gomodbus.BigEndian)
	t.Logf("Writing multiple registers at address %v - %v : %v", startingAddress, startingAddress+quantity, values)
	err = c.WriteMultipleRegisters(transactionID, startingAddress, unitID, values)
	if err != nil {
		t.Fatalf("Failed to write multiple registers: %v", err)
	}
	registers, err = c.ReadHoldingRegisters(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}
	if !reflect.DeepEqual(registers, values) {
		t.Errorf("Expected %v, but got %v", values, registers)
	}
	t.Logf("Holding registers at address %v - %v : %v", startingAddress, startingAddress+quantity, registers)
	decoded_value := int64(gomodbus.DecodeModbusRegisters(registers, gomodbus.BigEndian, gomodbus.BigEndian).(uint64))
	if !reflect.DeepEqual(decoded_value, value_64) {
		t.Errorf("Expected %v, but got %v", value_64, decoded_value)
	}
	t.Logf("Value: %v", decoded_value)

	// Test case 3: Write and read 2 holding registers starting from address 0
	quantity = 2
	value_32 := float32(22.33)
	t.Logf("Wrting value: %v as 32 bits into 2 registers", value_32)
	values = gomodbus.EncodeModbusRegisters(value_32, gomodbus.BigEndian, gomodbus.BigEndian)
	t.Logf("Writing multiple registers at address %v - %v : %v", startingAddress, startingAddress+quantity, values)
	err = c.WriteMultipleRegisters(transactionID, startingAddress, unitID, values)
	if err != nil {
		t.Fatalf("Failed to write multiple registers: %v", err)
	}
	registers, err = c.ReadHoldingRegisters(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}
	if !reflect.DeepEqual(registers, values) {

		t.Errorf("Expected %v, but got %v", values, registers)
	}
	t.Logf("Holding registers at address %v - %v : %v", startingAddress, startingAddress+quantity, registers)
	decoded_value_32 := math.Float32frombits(gomodbus.DecodeModbusRegisters(registers, gomodbus.BigEndian, gomodbus.BigEndian).(uint32))
	if !reflect.DeepEqual(decoded_value_32,value_32) {
		t.Errorf("Expected %v, but got %v", value_32, decoded_value_32)
	}
	t.Logf("Value: %v", decoded_value_32)
}

func EncodeModbusRegisters(value_64 int64, s1, s2 string) {
	panic("unimplemented")
}

func TestClientReadInputRegisters(t *testing.T) {
	c := client.TCPClient{
		Host: "192.168.5.3",
		Port: 1502,
	}

	err := c.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	// Test case 1: Read 3 input registers starting from address 0
	transactionID := 10
	quantity := 4
	startingAddress := 0
	unitID := 0
	registers, err := c.ReadInputRegisters(transactionID, startingAddress, quantity, unitID)
	if err != nil {
		t.Fatalf("Failed to read input registers: %v", err)
	}
	t.Logf("Response: %v", registers)
}

