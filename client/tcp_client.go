package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
)

type TCPClient struct {
	Host          string
	Port          int
	conn          net.Conn
	UseTLS        bool
	CertFile      string
	KeyFile       string
	CAFile        string // New field for custom CA
	transactionID uint16
}

func NewTCPClient(host string, port int, useTLS bool, certFile, keyFile, caFile string) Client {
	return &TCPClient{
		Host:     host,
		Port:     port,
		UseTLS:   useTLS,
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
	}
}

func (client *TCPClient) Connect() error {
	// Check if the host is a valid IP address
	ip := net.ParseIP(client.Host)
	if ip == nil {
		gomodbus.Logger.Error("invalid host IP address")
		return fmt.Errorf("invalid host IP address")
	}

	address := fmt.Sprintf("%s:%d", client.Host, client.Port)

	var conn net.Conn
	var err error

	if client.UseTLS {
		// Load client TLS certificate and key if provided
		var cert tls.Certificate
		if client.CertFile != "" && client.KeyFile != "" {
			cert, err = tls.LoadX509KeyPair(client.CertFile, client.KeyFile)
			if err != nil {
				return fmt.Errorf("failed to load TLS certificate and key: %v", err)
			}
		}

		// Configure TLS
		tlsConfig := &tls.Config{}
		if client.CertFile != "" && client.KeyFile != "" {
			tlsConfig.Certificates = []tls.Certificate{cert}
		}

		if client.CAFile != "" {
			// Load custom CA
			caCert, err := os.ReadFile(client.CAFile)
			if err != nil {
				return fmt.Errorf("failed to read CA file: %v", err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
			tlsConfig.RootCAs = caCertPool
		} else {
			// Use the system's root CAs for verification
			rootCAs, err := x509.SystemCertPool()
			if err != nil {
				return fmt.Errorf("failed to load system root CAs: %v", err)
			}
			tlsConfig.RootCAs = rootCAs
		}

		// Establish a TLS connection to the server
		conn, err = tls.Dial("tcp", address, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to connect to %s with TLS: %v", address, err)
		}
	} else {
		// Establish a non-TLS connection to the server
		conn, err = net.Dial("tcp", address)
		if err != nil {
			return fmt.Errorf("failed to connect to %s: %v", address, err)
		}
	}

	client.conn = conn
	gomodbus.Logger.Sugar().Debugf("client connected to server: %s, %d", address, client.Port)
	return nil
}

func (client *TCPClient) Disconnect() error {
	if client.conn != nil {
		err := client.conn.Close()
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to close connection: %v", err)
			return err
		}
		gomodbus.Logger.Sugar().Debugf("client disconnected from server: %s, %d", client.Host, client.Port)
		client.conn = nil
	}
	return nil
}

func (client *TCPClient) SendRequest(unitID byte, pduBytes []byte) error {
	adu := adu.NewTCPADU(client.transactionID, unitID, pduBytes)
	aduBytes := adu.ToBytes()
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("client failed to send request: %v", err)
		return err
	}
	client.transactionID++
	gomodbus.Logger.Sugar().Debugf("client sent request: %v", aduBytes)
	return nil
}

func (client *TCPClient) ReceiveResponse() ([]byte, error) {
	// Define a buffer with a size that can accommodate the expected response
	buffer := make([]byte, 512) // Adjust the size as needed

	// Read exactly the number of bytes into the buffer
	n, err := client.conn.Read(buffer)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("client failed to receive response: %v", err)
		return nil, err
	}

	// Slice the buffer to the actual number of bytes read
	responseBytes := buffer[:n]

	gomodbus.Logger.Sugar().Debugf("client received response: %v", responseBytes)

	responseADU := &adu.TCPADU{}
	err = responseADU.FromBytes(responseBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("client failed to parse response ADU: %v", err)
		return nil, err
	}
	return responseADU.PDU, nil
}

