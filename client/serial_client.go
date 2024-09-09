package client

import (
	"github.com/tarm/serial"
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
)

type SerialClient struct {
	Port     string
	BaudRate int
	conn     *serial.Port
}

func NewSerialClient(port string, baudRate int) Client {
	return &SerialClient{
		Port:     port,
		BaudRate: baudRate,
	}
}

func (client *SerialClient) Connect() error {
	c := &serial.Config{Name: client.Port, Baud: client.BaudRate}
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
	serialADU := adu.NewSerialADU(unitID, pduBytes)
	aduBytes := serialADU.ToBytes()
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("serial client failed to send request: %v", err)
		return err
	}
	gomodbus.Logger.Sugar().Debugf("serial client sent request: %v", aduBytes)
	return nil
}

func (client *SerialClient) ReceiveResponse() ([]byte, error) {
	buffer := make([]byte, 512)
	n, err := client.conn.Read(buffer)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("serial client failed to receive response: %v", err)
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
