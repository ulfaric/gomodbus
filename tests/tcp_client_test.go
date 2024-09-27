package tests

import (
	"testing"
	
	"github.com/ulfaric/gomodbus"
	c "github.com/ulfaric/gomodbus/client"
)



func TestTCPClient_Connect(t *testing.T) {
	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}
}

func TestTCPClient_ReadCoils(t *testing.T) {
	gomodbus.EnableDebug()
	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	coils, err :=c.ReadCoils(client, 1, 1, 1)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}
	t.Logf("Coils: %v", coils)

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}

}

func TestTCPClient_ReadDiscreteInputs(t *testing.T) {
	gomodbus.EnableDebug()
	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	discreteInputs, err := c.ReadDiscreteInputs(client, 1, 1, 10)
	if err != nil {
		t.Fatalf("Failed to read discrete inputs: %v", err)
	}
	t.Logf("Discrete inputs: %v", discreteInputs)

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}
}

func TestTCPClient_ReadHoldingRegisters(t *testing.T) {
	gomodbus.EnableDebug()
	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	registers, err := c.ReadHoldingRegisters(client, 1, 1, 10)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}
	t.Logf("Holding registers: %v", registers)

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}
}

func TestTCPClient_ReadInputRegisters(t *testing.T) {
	gomodbus.EnableDebug()
	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	registers, err := c.ReadInputRegisters(client, 1, 1, 10)
	if err != nil {
		t.Fatalf("Failed to read input registers: %v", err)
	}
	t.Logf("Input registers: %v", registers)

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}
}

func TestTCPClient_WriteSingleCoil(t *testing.T) {
	gomodbus.EnableDebug()
	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	err = c.WriteSingleCoil(client, 1, 1, false)
	if err != nil {
		t.Fatalf("Failed to write single coil: %v", err)
	}

	coils, err := c.ReadCoils(client, 1, 1, 1)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}
	t.Logf("Coils: %v", coils)

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}
}

func TestTCPClient_WriteMultipleCoils(t *testing.T) {
	gomodbus.EnableDebug()
	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	values := []bool{true, false, true}
	
	err = c.WriteMultipleCoils(client, 1, 1, values)
	if err != nil {
		t.Fatalf("Failed to write multiple coils: %v", err)
	}

	coils, err := c.ReadCoils(client, 1, 1, 3)
	if err != nil {
		t.Fatalf("Failed to read coils: %v", err)
	}
	t.Logf("Coils: %v", coils)

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}
}

func TestTCPClient_WriteSingleRegister(t *testing.T) {
	gomodbus.EnableDebug()
	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	err = c.WriteSingleRegister(client, 1, 1, []byte{0x12, 0xFF})
	if err != nil {
		t.Fatalf("Failed to write single register: %v", err)
	}

	registers, err := c.ReadHoldingRegisters(client, 1, 1, 1)
	if err != nil {
		t.Fatalf("Failed to read input registers: %v", err)
	}
	t.Logf("Holding registers: %x", registers)

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}
}

func TestTCPClient_WriteMultipleRegisters(t *testing.T) {
	gomodbus.EnableDebug()
	client := c.NewTCPClient("127.0.0.1", 1502, false, "", "", "")
	err := client.Connect()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}

	values := []byte{0x12, 0xFF, 0x34, 0x56}
	err = c.WriteMultipleRegisters(client, 1, 1, 2, values)
	if err != nil {
		t.Fatalf("Failed to write multiple registers: %v", err)
	}

	registers, err := c.ReadHoldingRegisters(client, 1, 1, 2)
	if err != nil {
		t.Fatalf("Failed to read holding registers: %v", err)
	}
	t.Logf("Holding registers: %x", registers)

	err = client.Disconnect()
	if err != nil {
		t.Fatalf("Failed to disconnect: %v", err)
	}
}
	
	

