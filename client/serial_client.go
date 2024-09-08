package client

import (
	"github.com/tarm/serial"
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
	"go.uber.org/zap"
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
		gomodbus.Logger.Error("failed to open serial port", zap.Error(err))
		return err
	}
	client.conn = conn
	gomodbus.Logger.Debug("serial client connected", zap.String("port", client.Port), zap.Int("baudRate", client.BaudRate))
	return nil
}

func (client *SerialClient) Disconnect() error {
	if client.conn != nil {
		err := client.conn.Close()
		if err != nil {
			gomodbus.Logger.Error("failed to close serial port", zap.Error(err))
			return err
		}
		gomodbus.Logger.Debug("serial client disconnected", zap.String("port", client.Port))
	}
	return nil
}

func (client *SerialClient) SendRequest(unitID byte, pduBytes []byte) error {
	serialADU := adu.NewSerialADU(unitID, pduBytes)
	aduBytes := serialADU.ToBytes()
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		gomodbus.Logger.Error("serial client failed to send request", zap.Error(err))
		return err
	}
	gomodbus.Logger.Debug("serial client sent request", zap.ByteString("request", aduBytes))
	return nil
}

func (client *SerialClient) ReceiveResponse() ([]byte, error) {
	buffer := make([]byte, 512)
	n, err := client.conn.Read(buffer)
	if err != nil {
		gomodbus.Logger.Error("serial client failed to receive response", zap.Error(err))
		return nil, err
	}
	responseBytes := buffer[:n]
	gomodbus.Logger.Debug("serial client received response", zap.ByteString("response", responseBytes))

	responseADU := &adu.SerialADU{}
	err = responseADU.FromBytes(responseBytes)
	if err != nil {
		gomodbus.Logger.Error("serial client failed to parse response ADU", zap.Error(err))
		return nil, err
	}
	return responseADU.PDU, nil
}
