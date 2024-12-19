package client

import (
	"fmt"
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

	if client.conn == nil {
		gomodbus.Logger.Sugar().Errorf("serial client not connected")
		return nil, fmt.Errorf("serial client not connected")
	}

	time.Sleep(time.Duration(client.bufferTime) * time.Millisecond)

	buffer := make([]byte, 512)

	n, err := client.conn.Read(buffer)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to read response: %v", err)
		return nil, err
	}

	responseBytes := buffer[:n]
	gomodbus.Logger.Sugar().Debugf("serial client received response: %v", responseBytes)

	responseADU := &adu.SerialADU{}
	err = responseADU.FromBytes(responseBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("serial client failed to parse response ADU: %v", err)
		return nil, err
	}
	return responseADU.PDU, nil

}
