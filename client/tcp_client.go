package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
	"github.com/ulfaric/gomodbus/pdu"
)

type TCPClient struct {
	Host     string
	Port     int
	conn     net.Conn
	UseTLS   bool
	CertFile string
	KeyFile  string
	CAFile   string // New field for custom CA
}


func (client *TCPClient) Connect() error {
	// Check if the host is a valid IP address
	ip := net.ParseIP(client.Host)
	if ip == nil {
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
	return nil
}


func (client *TCPClient) Close() error {
	if client.conn != nil {
		return client.conn.Close()
	}
	return nil
}

func (client *TCPClient) ReadCoils(transactionID, startingAddress, quantity, unitID int) ([]bool, error) {
	// Create a PDU for the Read Coils request
	readPDU := pdu.NewPDUReadCoils(uint16(startingAddress), uint16(quantity))
	tcpADU := adu.NewTCPADU(uint16(transactionID), byte(unitID), readPDU.ToBytes())
	aduBytes := tcpADU.ToBytes()

	// Send the request
	log.Printf("Sending Read Coils Request: %x", aduBytes)
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		return nil, err
	}

	// Calculate the expected response length
	responseLength, err := gomodbus.CalculateADULength(gomodbus.ReadCoil, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	log.Printf("Received Read Coils Response: %x", response)	

	// Check for ADU errors
	err = gomodbus.CheckModbusError(response)
	if err != nil {
		return nil, err
	}

	// Parse the response ADU
	responseADU := &adu.TCPADU{}
	err = responseADU.FromBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response ADU: %v", err)
	}

	// Verify the function code in the response PDU
	if responseADU.PDU[0] != gomodbus.ReadCoil {
		return nil, fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.ReadCoil, responseADU.PDU[0])
	}

	// Parse the response PDU
	readResponsePDU := &pdu.PDU_ReadResponse{}
	readResponsePDU.FromBytes(responseADU.PDU)

	// Parse the coil status from the response PDU
	coils := make([]bool, quantity)
	for i := 0; i < quantity; i++ {
		byteIndex := i / 8
		bitIndex := i % 8
		coils[i] = (readResponsePDU.Data[byteIndex] & (1 << bitIndex)) != 0
	}

	return coils, nil
}


func (client *TCPClient) ReadDiscreteInputs(transactionID, startingAddress, quantity, unitID int) ([]bool, error) {
	// Create a PDU for the Read Discrete Inputs request
	readPDU := pdu.NewPDUReadDiscreteInputs(uint16(startingAddress), uint16(quantity))
	tcpADU := adu.NewTCPADU(uint16(transactionID), byte(unitID), readPDU.ToBytes())
	aduBytes := tcpADU.ToBytes()

	// Send the request
	log.Printf("Sending Read Discrete Inputs Request: %x", aduBytes)
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		return nil, err
	}

	// Calculate the expected response length
	responseLength, err := gomodbus.CalculateADULength(gomodbus.ReadDiscreteInput, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	log.Printf("Received Read Discrete Inputs Response: %x", response)

	// Check for ADU errors
	err = gomodbus.CheckModbusError(response)
	if err != nil {
		return nil, err
	}

	// Parse the response ADU
	responseADU := &adu.TCPADU{}
	err = responseADU.FromBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response ADU: %v", err)
	}

	// Verify the function code in the response PDU
	if responseADU.PDU[0] != gomodbus.ReadDiscreteInput {
		return nil, fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.ReadDiscreteInput, responseADU.PDU[0])
	}

	// Parse the response PDU
	readResponsePDU := &pdu.PDU_ReadResponse{}
	readResponsePDU.FromBytes(responseADU.PDU)

	// Parse the input status from the response PDU
	inputs := make([]bool, quantity)
	for i := 0; i < quantity; i++ {
		byteIndex := i / 8
		bitIndex := i % 8
		inputs[i] = (readResponsePDU.Data[byteIndex] & (1 << bitIndex)) != 0
	}

	return inputs, nil
}


func (client *TCPClient) ReadHoldingRegisters(transactionID, startingAddress, quantity, unitID int) ([]uint16, error) {
	// Create a PDU for the Read Holding Registers request
	readPDU := pdu.NewPDUReadHoldingRegisters(uint16(startingAddress), uint16(quantity))
	tcpADU := adu.NewTCPADU(uint16(transactionID), byte(unitID), readPDU.ToBytes())
	aduBytes := tcpADU.ToBytes()

	// Send the request
	log.Printf("Sending Read Holding Registers Request: %x", aduBytes)
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		return nil, err
	}

	// Calculate the expected response length
	responseLength, err := gomodbus.CalculateADULength(gomodbus.ReadHoldingRegister, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	log.Printf("Received Read Holding Registers Response: %x", response)

	// Check for ADU errors
	err = gomodbus.CheckModbusError(response)
	if err != nil {
		return nil, err
	}

	// Parse the response ADU
	responseADU := &adu.TCPADU{}
	err = responseADU.FromBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response ADU: %v", err)
	}

	// Verify the function code in the response PDU
	if responseADU.PDU[0] != gomodbus.ReadHoldingRegister {
		return nil, fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.ReadHoldingRegister, responseADU.PDU[0])
	}

	// Parse the response PDU
	readResponsePDU := &pdu.PDU_ReadResponse{}
	readResponsePDU.FromBytes(responseADU.PDU)

	// Parse the register values from the response PDU
	registers := make([]uint16, quantity)
	for i := 0; i < quantity; i++ {
		register := binary.BigEndian.Uint16(readResponsePDU.Data[i*2 : i*2+2])
		registers[i] = register
	}

	return registers, nil
}


func (client *TCPClient) ReadInputRegisters(transactionID, startingAddress, quantity, unitID int) ([]uint16, error) {
	// Create a PDU for the Read Input Registers request
	readPDU := pdu.NewPDUReadInputRegisters(uint16(startingAddress), uint16(quantity))
	tcpADU := adu.NewTCPADU(uint16(transactionID), byte(unitID), readPDU.ToBytes())
	aduBytes := tcpADU.ToBytes()

	// Send the request
	log.Printf("Sending Read Input Registers Request: %x", aduBytes)
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		return nil, err
	}

	// Calculate the expected response length
	responseLength, err := gomodbus.CalculateADULength(gomodbus.ReadInputRegister, quantity)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate response length: %v", err)
	}

	// Read the response
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return nil, err
	}

	log.Printf("Received Read Input Registers Response: %x", response)

	// Check for ADU errors
	err = gomodbus.CheckModbusError(response)
	if err != nil {
		return nil, err
	}

	// Parse the response ADU
	responseADU := &adu.TCPADU{}
	err = responseADU.FromBytes(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response ADU: %v", err)
	}

	// Verify the function code in the response PDU
	if responseADU.PDU[0] != gomodbus.ReadInputRegister {
		return nil, fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.ReadInputRegister, responseADU.PDU[0])
	}

	// Parse the response PDU
	readResponsePDU := &pdu.PDU_ReadResponse{}
	readResponsePDU.FromBytes(responseADU.PDU)

	// Parse the register values from the response PDU
	registers := make([]uint16, quantity)
	for i := 0; i < quantity; i++ {
		register := binary.BigEndian.Uint16(readResponsePDU.Data[i*2 : i*2+2])
		registers[i] = register
	}

	return registers, nil
}


func (client *TCPClient) WriteSingleCoil(transactionID, address, unitID int, value bool) error {
	// Create a PDU for the Write Single Coil request
	writePDU := pdu.NewPDUWriteSingleCoil(uint16(address), value)
	tcpAdu := adu.NewTCPADU(uint16(transactionID), byte(unitID), writePDU.ToBytes())
	aduBytes := tcpAdu.ToBytes()

	// Send the request
	log.Printf("Sending Write Single Coil Request: %x", aduBytes)
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		return err
	}

	// Calculate the expected response length
	responseLength, err := gomodbus.CalculateADULength(gomodbus.WriteSingleCoil, 1)
	if err != nil {
		return fmt.Errorf("failed to calculate response length: %v", err)
	}
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return err
	}

	log.Printf("Received Write Single Coil Response: %x", response)

	// Check for ADU errors
	err = gomodbus.CheckModbusError(response)
	if err != nil {
		return err
	}

	// Parse the response ADU
	responseADU := &adu.TCPADU{}
	err = responseADU.FromBytes(response)
	if err != nil {
		return fmt.Errorf("failed to parse response ADU: %v", err)
	}

	// Verify the function code in the response PDU
	if responseADU.PDU[0] != gomodbus.WriteSingleCoil {
		return fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.WriteSingleCoil, responseADU.PDU[0])
	}

	// Parse the response PDU
	writeResponsePDU := &pdu.PDU_WriteSingleResponse{}
	writeResponsePDU.FromBytes(responseADU.PDU)

	// Verify the address and value in the response PDU
	if writeResponsePDU.OutputAddress != uint16(address) {
		return fmt.Errorf("invalid address in response, expect %v but received %v", address, writeResponsePDU.OutputAddress)
	}

	if (value && writeResponsePDU.OutputValue != 0xFF00) || (!value && writeResponsePDU.OutputValue != 0x0000) {
		return fmt.Errorf("invalid value in response, expect %v but received %v", value, writeResponsePDU.OutputValue)
	}

	return nil
}


func (client *TCPClient) WriteSingleRegister(transactionID, address, unitID int, value uint16) error {
	// Create a PDU for the Write Single Register request
	writePDU := pdu.NewPDUWriteSingleRegister(uint16(address), value)
	tcpADU := adu.NewTCPADU(uint16(transactionID), byte(unitID), writePDU.ToBytes())
	aduBytes := tcpADU.ToBytes()

	// Send the request
	log.Printf("Sending Write Single Register Request: %x", aduBytes)
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		return err
	}

	// Calculate the expected response length
	responseLength, err := gomodbus.CalculateADULength(gomodbus.WriteSingleRegister, 1)
	if err != nil {
		return fmt.Errorf("failed to calculate response length: %v", err)
	}
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return err
	}

	log.Printf("Received Write Single Register Response: %x", response)

	// Check for ADU errors
	err = gomodbus.CheckModbusError(response)
	if err != nil {
		return err
	}

	// Parse the response ADU
	responseADU := &adu.TCPADU{}
	err = responseADU.FromBytes(response)
	if err != nil {
		return fmt.Errorf("failed to parse response ADU: %v", err)
	}

	// Verify the function code in the response PDU
	if responseADU.PDU[0] != gomodbus.WriteSingleRegister {
		return fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.WriteSingleRegister, responseADU.PDU[0])
	}

	// Parse the response PDU
	writeResponsePDU := &pdu.PDU_WriteSingleResponse{}
	writeResponsePDU.FromBytes(responseADU.PDU)

	// Verify the address in the response PDU
	if writeResponsePDU.OutputAddress != uint16(address) {
		return fmt.Errorf("invalid address in response, expect %v but received %v", address, writeResponsePDU.OutputAddress)
	}

	// Verify the value in the response PDU
	if writeResponsePDU.OutputValue != value {
		return fmt.Errorf("invalid value in response, expect %v but received %v", value, writeResponsePDU.OutputValue)
	}

	return nil
}


func (client *TCPClient) WriteMultipleCoils(transactionID, startingAddress, unitID int, values []bool) error {
	// Create a PDU for the Write Multiple Coils request
	writePDU := pdu.NewPDUWriteMultipleCoils(uint16(startingAddress), values)
	tcpADU := adu.NewTCPADU(uint16(transactionID), byte(unitID), writePDU.ToBytes())
	aduBytes := tcpADU.ToBytes()

	// Send the request
	log.Printf("Sending Write Multiple Coils Request: %x", aduBytes)
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		return err
	}

	// Calculate the expected response length
	responseLength, err := gomodbus.CalculateADULength(gomodbus.WriteMultipleCoils, 1)
	if err != nil {
		return fmt.Errorf("failed to calculate response length: %v", err)
	}
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return err
	}

	// Check for ADU errors
	err = gomodbus.CheckModbusError(response)
	if err != nil {
		return err
	}

	log.Printf("Received Write Multiple Coils Response: %x", response)

	// Parse the response ADU
	responseADU := &adu.TCPADU{}
	err = responseADU.FromBytes(response)
	if err != nil {
		return fmt.Errorf("failed to parse response ADU: %v", err)
	}

	// Verify the function code in the response PDU
	if responseADU.PDU[0] != gomodbus.WriteMultipleCoils {
		return fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.WriteMultipleCoils, responseADU.PDU[0])
	}

	// Parse the response PDU
	writeResponsePDU := &pdu.PDU_WriteMultipleResponse{}
	writeResponsePDU.FromBytes(responseADU.PDU)

	// Verify the starting address in the response PDU
	if writeResponsePDU.StartingAddress != uint16(startingAddress) {
		return fmt.Errorf("invalid starting address in response, expect %v but received %v", startingAddress, writeResponsePDU.StartingAddress)
	}

	// Verify the quantity in the response PDU
	if writeResponsePDU.Quantity != uint16(len(values)) {
		return fmt.Errorf("invalid quantity in response, expect %v but received %v", len(values), writeResponsePDU.Quantity)
	}

	return nil
}


func (client *TCPClient) WriteMultipleRegisters(transactionID, startingAddress, unitID int, values []uint16) error {
	// Create a PDU for the Write Multiple Registers request
	writePDU := pdu.NewPDUWriteMultipleRegisters(uint16(startingAddress), values)
	tcpADU := adu.NewTCPADU(uint16(transactionID), byte(unitID), writePDU.ToBytes())
	aduBytes := tcpADU.ToBytes()

	// Send the request
	log.Printf("Sending Write Multiple Registers Request: %x", aduBytes)
	_, err := client.conn.Write(aduBytes)
	if err != nil {
		return err
	}

	// Calculate the expected response length
	responseLength, err := gomodbus.CalculateADULength(gomodbus.WriteMultipleRegisters, 1)
	if err != nil {
		return fmt.Errorf("failed to calculate response length: %v", err)
	}
	response := make([]byte, responseLength)
	_, err = client.conn.Read(response)
	if err != nil {
		return err
	}

	// Check for ADU errors
	err = gomodbus.CheckModbusError(response)
	if err != nil {
		return err
	}

	log.Printf("Received Write Multiple Registers Response: %x", response)

	// Parse the response ADU
	responseADU := &adu.TCPADU{}
	err = responseADU.FromBytes(response)
	if err != nil {
		return fmt.Errorf("failed to parse response ADU: %v", err)
	}

	// Verify the function code in the response PDU
	if responseADU.PDU[0] != gomodbus.WriteMultipleRegisters {
		return fmt.Errorf("invalid function code in response, expect %x but received %x", gomodbus.WriteMultipleRegisters, responseADU.PDU[0])
	}

	// Parse the response PDU
	writeResponsePDU := &pdu.PDU_WriteMultipleResponse{}
	writeResponsePDU.FromBytes(responseADU.PDU)

	// Verify the starting address in the response PDU
	if writeResponsePDU.StartingAddress != uint16(startingAddress) {
		return fmt.Errorf("invalid starting address in response, expect %v but received %v", startingAddress, writeResponsePDU.StartingAddress)
	}

	// Verify the quantity in the response PDU
	if writeResponsePDU.Quantity != uint16(len(values)) {
		return fmt.Errorf("invalid quantity in response, expect %v but received %v", len(values), writeResponsePDU.Quantity)
	}

	return nil
}

