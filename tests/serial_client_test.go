package tests

import (
    "testing"

    c "github.com/ulfaric/gomodbus/client"
    "github.com/ulfaric/gomodbus"
)

func TestSerialClient_ReadInputRegister(t *testing.T) {
    gomodbus.EnableDebug()
    client := c.NewSerialClient("/dev/ttyUSB0", 9600, 8, 'N', 1, 500)
    err := client.Connect()
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer client.Disconnect()

    registers, err := c.ReadHoldingRegisters(client, 170, 2004, 1)
    if err != nil {
        t.Fatalf("Failed to read input registers: %v", err)
    }
    t.Logf("Input registers: %x", registers)
}


func TestSerialClient_WriteRegister(t *testing.T) {
    gomodbus.EnableDebug()
    client := c.NewSerialClient("/dev/ttyUSB0", 9600, 8, 'N', 1, 500)
    err := client.Connect()
    if err != nil {
        t.Fatalf("Failed to connect: %v", err)
    }
    defer client.Disconnect()

	value := uint16(50)
	valueBytes, err := gomodbus.Serializer(value, "big", "big")
	if err != nil {
		t.Fatalf("Failed to serialize value: %v", err)
	}
	t.Logf("Value bytes: %x", valueBytes)

    // Write a value to a register
    err = c.WriteMultipleRegisters(client, 170, 3001, 1, valueBytes[0])
    if err != nil {
        t.Fatalf("Failed to write to register: %v", err)
    }
    t.Logf("Successfully wrote to register")
}
