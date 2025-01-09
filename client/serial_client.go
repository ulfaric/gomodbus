package client

import (
	"fmt"
	"sync"
	"time"

	"github.com/tarm/serial"
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
)

type SerialClient struct {
	Port       string
	BaudRate   int
	DataBits   byte
	Parity     serial.Parity
	StopBits   serial.StopBits
	conn       *serial.Port
	bufferTime int
	timeout    int

	mu sync.Mutex
}

func NewSerialClient(port string, baudRate int, dataBits byte, parity serial.Parity, stopBits serial.StopBits, bufferTime int, timeout int) Client {
	return &SerialClient{
		Port:       port,
		BaudRate:   baudRate,
		DataBits:   dataBits,
		Parity:     parity,
		StopBits:   stopBits,
		bufferTime: bufferTime,
		timeout:    timeout,
	}
}

func (client *SerialClient) Connect() error {
	c := &serial.Config{Name: client.Port, Baud: client.BaudRate, Parity: client.Parity, StopBits: client.StopBits, Size: client.DataBits, ReadTimeout: time.Duration(client.bufferTime) * time.Millisecond}
	conn, err := serial.OpenPort(c)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to open serial port: %v", err)
		return err
	}
	client.conn = conn
	gomodbus.Logger.Sugar().Debugf("serial client connected: %s, %d", client.Port, client.BaudRate)
	return nil
}

func (client *SerialClient) Disconnect() error {
	if client.conn != nil {
		err := client.conn.Close()
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to close serial port: %v", err)
			return err
		}
		gomodbus.Logger.Sugar().Debugf("serial client disconnected: %s", client.Port)
	}
	return nil
}

func (client *SerialClient) SendRequest(unitID byte, pduBytes []byte) error {
	client.mu.Lock()
	defer client.mu.Unlock()

	if client.conn == nil {
		gomodbus.Logger.Sugar().Errorf("serial client not connected")
		return fmt.Errorf("serial client not connected")
	}

	serialADU := adu.NewSerialADU(unitID, pduBytes)
	aduBytes := serialADU.ToBytes()

	_, err := client.conn.Write(aduBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to send request: %v", err)
		return err
	}

	gomodbus.Logger.Sugar().Debugf("serial client sent request: %v", aduBytes)
	return nil
}

func (client *SerialClient) ReceiveResponse() ([]byte, error) {
	client.mu.Lock()
	defer client.mu.Unlock()

	if client.conn == nil {
		gomodbus.Logger.Sugar().Errorf("serial client not connected")
		return nil, fmt.Errorf("serial client not connected")
	}

	// Initial delay for buffer
	time.Sleep(time.Duration(client.bufferTime) * time.Millisecond)

	responseChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		buffer := make([]byte, 512)
		n, err := client.conn.Read(buffer)
		if err != nil {
			errChan <- err
			return
		}

		responseBytes := buffer[:n]
		gomodbus.Logger.Sugar().Debugf("serial client received response: %v", responseBytes)

		responseADU := &adu.SerialADU{}
		err = responseADU.FromBytes(responseBytes)
		if err != nil {
			errChan <- err
			return
		}
		responseChan <- responseADU.PDU
	}()

	// Wait for either response or timeout
	select {
	case response := <-responseChan:
		return response, nil
	case err := <-errChan:
		return nil, fmt.Errorf("read error: %v", err)
	case <-time.After(time.Duration(client.timeout) * time.Millisecond):
		client.Disconnect()
		client.Connect()
		return nil, fmt.Errorf("timeout waiting for response after %d ms", client.timeout)
	}
}
