package client

import (
	"fmt"
	"gomodbus"
	"gomodbus/adu"
	"net"
)

type TCPClient struct {
	Host string
	Port int
	conn net.Conn
}

func (client *TCPClient) Connect() error {
	// Check if the host is a valid IP address
	ip := net.ParseIP(client.Host)
	if ip == nil {
		return fmt.Errorf("invalid host IP address")
	}
	// establish a connection to the server
	address := fmt.Sprintf("%s:%d", client.Host, client.Port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	client.conn = conn
	return nil
}

func (client *TCPClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}

func (client *TCPClient) ReadCoils(transactionID, startingAddress, quantity uint16, unitID byte) ([]bool, error) {
	adu := adu.ReadCoilsTCPADU(transactionID, startingAddress, quantity, unitID)
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return nil, err
	}
	// read the response
	response := make([]byte, 10)
	n, err := client.conn.Read(response)
	if err != nil {
		return nil, err
	}
    if n < 9 {
        return nil, fmt.Errorf("response too short")
    }

    if response[7] != byte(gomodbus.ReadCoil) {
        return nil, fmt.Errorf("invalid function code in response")
    }

    byteCount := response[8]
    if int(byteCount) > n-9 {
        return nil, fmt.Errorf("byte count mismatch")
    }

    coils := make([]bool, quantity)
    for i := uint16(0); i < quantity; i++ {
        byteIndex := 9 + i/8
        bitIndex := i % 8
        coils[i] = response[byteIndex]&(1<<bitIndex) != 0
    }

    return coils, nil
}