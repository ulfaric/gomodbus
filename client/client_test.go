package client

import (
	"testing"
	"gomodbus/pdu"
	"gomodbus/adu"
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
	pdu := pdu.ReadCoilsPDU(startingAddress, quantity)
	t.Logf("PDU: %v", pdu)
	adu := adu.ReadCoilsTCPADU(transactionID, startingAddress, quantity, unitID)
	t.Logf("ADU: %v", adu)
	adu_bytes := adu.ToBytes()
	t.Logf("ADU bytes: %v", adu_bytes)
	client.conn.Write(adu_bytes)
	response := make([]byte, 10)
	n, err := client.conn.Read(response)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}
	t.Logf("Response: %v", response)
    if n < 9 {
        t.Fatalf("response too short")
    }
    byteCount := response[8]
    if int(byteCount) > n-9 {
        t.Fatalf("response too short")
    }

    coils := make([]bool, quantity)
    for i := uint16(0); i < quantity; i++ {
        byteIndex := 9 + i/8
        bitIndex := i % 8
        coils[i] = response[byteIndex]&(1<<bitIndex) != 0
    }
	t.Logf("Coils: %v", coils)
}