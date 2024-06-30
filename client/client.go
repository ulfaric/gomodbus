package client

import (
	"fmt"
	"net"
	"gomodbus/adu"
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

func (client *TCPClient) ReadCoils(transactionID, startingAddress, quantity uint16, unitID byte) ([]byte, error) {
	adu := adu.ReadCoilsTCPADU(transactionID, startingAddress, quantity, unitID)
	adu_bytes := adu.ToBytes()
	_, err := client.conn.Write(adu_bytes)
	if err != nil {
		return nil, err
	}
	// read the response
	var response []byte
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}
	return response, nil
}