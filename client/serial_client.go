package client

import (
	"context"
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
}

func NewSerialClient(port string, baudRate int, dataBits byte, parity serial.Parity, stopBits serial.StopBits, bufferTime int) Client {
	return &SerialClient{
		Port:       port,
		BaudRate:   baudRate,
		DataBits:   dataBits,
		Parity:     parity,
		StopBits:   stopBits,
		bufferTime: bufferTime,
	}
}

func (client *SerialClient) Connect() error {
	c := &serial.Config{Name: client.Port, Baud: client.BaudRate, Parity: client.Parity, StopBits: client.StopBits, Size: client.DataBits}
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serialADU := adu.NewSerialADU(unitID, pduBytes)
	aduBytes := serialADU.ToBytes()

	done := make(chan error, 1)

	go func() {
		_, err := client.conn.Write(aduBytes)
		done <- err
	}()

	select {
	case <-ctx.Done():
		gomodbus.Logger.Sugar().Errorf("serial client write operation timed out or canceled: %v", ctx.Err())
		return ctx.Err()
	case err := <-done:
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("serial client failed to send request: %v", err)
			return err
		}
		gomodbus.Logger.Sugar().Debugf("serial client sent request: %v", aduBytes)
		return nil
	}
}

func (client *SerialClient) ReceiveResponse() ([]byte, error) {

	time.Sleep(time.Duration(client.bufferTime) * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	buffer := make([]byte, 512)
	done := make(chan struct {
		n   int
		err error
	}, 1)

	go func() {
		n, err := client.conn.Read(buffer)
		done <- struct {
			n   int
			err error
		}{n, err}
	}()

	select {
	case <-ctx.Done():
		gomodbus.Logger.Sugar().Errorf("serial client read operation timed out or canceled: %v", ctx.Err())
		return nil, ctx.Err()
	case result := <-done:
		if result.err != nil {
			gomodbus.Logger.Sugar().Errorf("serial client failed to receive response: %v", result.err)
			return nil, result.err
		}
		responseBytes := buffer[:result.n]
		gomodbus.Logger.Sugar().Debugf("serial client received response: %v", responseBytes)

		responseADU := &adu.SerialADU{}
		err := responseADU.FromBytes(responseBytes)
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("serial client failed to parse response ADU: %v", err)
			return nil, err
		}
		return responseADU.PDU, nil
	}
}
