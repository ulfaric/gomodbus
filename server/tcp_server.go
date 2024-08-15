package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
	"github.com/ulfaric/gomodbus/pdu"
	"gopkg.in/yaml.v3"
)

type TCPServer struct {
	Host      string
	Port      int
	ByteOrder string
	WordOrder string
	Slaves    map[byte]*Slave
	mu        sync.Mutex

	UseTLS  bool
	CertFile string
	KeyFile  string
	CAFile   string
}



type TCPConfig struct {
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	ByteOrder string `yaml:"byteOrder"`
	WordOrder string `yaml:"wordOrder"`
	UseTLS    bool   `yaml:"useTLS"`
	CertFile  string `yaml:"certFile"`
	KeyFile   string `yaml:"keyFile"`
	CAFile    string `yaml:"caFile"` // New field for custom CA

	Slaves []struct {
		UnitID byte `yaml:"unitID"`
		Coils  []struct {
			Address uint16 `yaml:"address"`
			Value   bool   `yaml:"value"`
		} `yaml:"coils"`
		DiscreteInputs []struct {
			Address uint16 `yaml:"address"`
			Value   bool   `yaml:"value"`
		} `yaml:"discreteInputs"`
		HoldingRegisters []struct {
			Address uint16 `yaml:"address"`
			Value   uint16 `yaml:"value"`
		} `yaml:"holdingRegisters"`
		InputRegisters []struct {
			Address uint16 `yaml:"address"`
			Value   uint16 `yaml:"value"`
		} `yaml:"inputRegisters"`
	} `yaml:"slaves"`
}



func NewTCPServer(host, byteOrder, wordOrder string, port int) (*TCPServer, error) {
	// read config file
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Printf("Failed to read config file: %v", err)
	}

	// parse config file
	var config TCPConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Printf("Failed to parse config file: %v", err)
	}

	if config.ByteOrder != gomodbus.BigEndian && config.ByteOrder != gomodbus.LittleEndian {
		return nil, fmt.Errorf("invalid byte order: %s", config.ByteOrder)
	}
	if config.WordOrder != gomodbus.BigEndian && config.WordOrder != gomodbus.LittleEndian {
		return nil, fmt.Errorf("invalid word order: %s", config.WordOrder)
	}

	slaves := make(map[byte]*Slave)
	for _, slaveConfig := range config.Slaves {
		slave := &Slave{}

		// Set coils and their legal addresses
		for _, coil := range slaveConfig.Coils {
			slave.Coils[coil.Address] = coil.Value
			slave.LegalCoilsAddress[coil.Address] = true
		}

		// Set discrete inputs and their legal addresses
		for _, input := range slaveConfig.DiscreteInputs {
			slave.DiscreteInputs[input.Address] = input.Value
			slave.LegalDiscreteInputsAddress[input.Address] = true
		}

		// Set holding registers and their legal addresses
		for _, register := range slaveConfig.HoldingRegisters {
			slave.HoldingRegisters[register.Address] = register.Value
			slave.LegalHoldingRegistersAddress[register.Address] = true
		}

		// Set input registers and their legal addresses
		for _, register := range slaveConfig.InputRegisters {
			slave.InputRegisters[register.Address] = register.Value
			slave.LegalInputRegistersAddress[register.Address] = true
		}

		slaves[slaveConfig.UnitID] = slave
	}

	return &TCPServer{
		Host:      host,
		Port:      port,
		ByteOrder: byteOrder,
		WordOrder: wordOrder,
		Slaves:    slaves,
		UseTLS:    config.UseTLS,
		CertFile:  config.CertFile,
		KeyFile:   config.KeyFile,
		CAFile:    config.CAFile, // Load the CA file from config
	}, nil
}



func (s *TCPServer) AddSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Slaves[unitID] = &Slave{}
}

func (s *TCPServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	var listener net.Listener
	var err error

	if s.UseTLS {
		// Load server TLS certificate and key
		cert, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate and key: %v", err)
		}

		// Load custom CA if provided
		var tlsConfig *tls.Config
		if s.CAFile != "" {
			caCert, err := os.ReadFile(s.CAFile)
			if err != nil {
				return fmt.Errorf("failed to read CA file: %v", err)
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
				ClientCAs:    caCertPool,
				ClientAuth:   tls.RequireAndVerifyClientCert,
			}
		} else {
			tlsConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
		}

		// Start listener with TLS
		listener, err = tls.Listen("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to listen on %s with TLS: %v", addr, err)
		}
		log.Printf("Modbus server started with TLS on %s", addr)
	} else {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to listen on %s: %v", addr, err)
		}
		log.Printf("Modbus server started on %s", addr)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn)
	}
}


func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 512)
	for {
		conn.SetDeadline(time.Now().Add(5 * time.Minute))
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		request := buffer[:n]
		response, err := s.processRequest(request)
		if err != nil {
			log.Printf("Failed to process request: %v", err)
			response = s.exceptionResponse(request, 0x04) // Server Device Failure
		}

		_, err = conn.Write(response)
		if err != nil {
			log.Printf("Failed to write response: %v", err)
			return
		}
	}
}

func (s *TCPServer) getSlave(unitID byte) (*Slave, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slave, exists := s.Slaves[unitID]
	if !exists {
		return nil, fmt.Errorf("slave %d not found", unitID)
	}
	return slave, nil
}

func (s *TCPServer) processRequest(request []byte) ([]byte, error) {
    // Parse the request as a TCP ADU
    tcpADU := &adu.TCPADU{}
    err := tcpADU.FromBytes(request)
    if err != nil {
        return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
    }

    // Verify if the actual length matches the length specified in the ADU
    if len(request) != int(tcpADU.Length)+6 {
        return s.exceptionResponse(request, 0x03), nil // Illegal Data Value
    }

    // Fetch the slave based on Unit ID
    slave, err := s.getSlave(tcpADU.UnitID)
    if err != nil {
        return s.exceptionResponse(request, 0x04), nil // Server Device Failure
    }

	// Process the request based on the function code
	var responsePDU []byte
    switch tcpADU.PDU[0] { // Function code is the first byte in PDU
    case 0x01:
        responsePDU, err = s.readCoils(slave, tcpADU.PDU)
    case 0x02:
        responsePDU, err = s.readDiscreteInputs(slave, tcpADU.PDU)
    case 0x03:
        responsePDU, err = s.readHoldingRegisters(slave, tcpADU.PDU)
    case 0x04:
        responsePDU, err = s.readInputRegisters(slave, tcpADU.PDU)
    case 0x05:
        responsePDU, err = s.writeSingleCoil(slave, tcpADU.PDU)
    case 0x06:
        responsePDU, err = s.writeSingleRegister(slave, tcpADU.PDU)
    case 0x0F:
        responsePDU, err = s.writeMultipleCoils(slave, tcpADU.PDU)
    case 0x10:
        responsePDU, err = s.writeMultipleRegisters(slave, tcpADU.PDU)
    default:
        return s.exceptionResponse(request, 0x01), nil // Illegal Function
    }

    if err != nil {
        return s.exceptionResponse(request, 0x04), err // Server Device Failure
    }

    // Create the response ADU using the response PDU
    responseADU := adu.NewTCPADU(tcpADU.TransactionID, tcpADU.UnitID, responsePDU)

    return responseADU.ToBytes(), nil
}

func (s *TCPServer) readCoils(slave *Slave, pduBytes []byte) ([]byte, error) {
    // Parse the PDU
    pduRead := &pdu.PDU_Read{}
    pduRead.FromBytes(pduBytes)

    startAddress := pduRead.StartingAddress
    quantity := pduRead.Quantity

    // Validate request
    if startAddress+quantity > 65535 || !slave.LegalCoilsAddress[startAddress] {
        return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
    }

    // Prepare response data
    byteCount := (quantity + 7) / 8
    responseData := make([]byte, byteCount)

    for i := 0; i < int(quantity); i++ {
        if !slave.LegalCoilsAddress[startAddress+uint16(i)] {
            return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
        }
        if slave.Coils[startAddress+uint16(i)] {
            responseData[i/8] |= 1 << (i % 8)
        }
    }

    // Create the response PDU
    responsePDU := pdu.NewPDUReadResponse(0x01, responseData)

    return responsePDU.ToBytes(), nil
}

func (s *TCPServer) readDiscreteInputs(slave *Slave, pduBytes []byte) ([]byte, error) {
    // Parse the PDU
    pduRead := &pdu.PDU_Read{}
    pduRead.FromBytes(pduBytes)

    startAddress := pduRead.StartingAddress
    quantity := pduRead.Quantity

    // Validate request
    if startAddress+quantity > 65535 || !slave.LegalDiscreteInputsAddress[startAddress] {
        return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
    }

    // Prepare response data
    byteCount := (quantity + 7) / 8
    responseData := make([]byte, byteCount)

    for i := 0; i < int(quantity); i++ {
        if !slave.LegalDiscreteInputsAddress[startAddress+uint16(i)] {
            return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
        }
        if slave.DiscreteInputs[startAddress+uint16(i)] {
            responseData[i/8] |= 1 << (i % 8)
        }
    }

    // Create the response PDU
    responsePDU := pdu.NewPDUReadResponse(0x02, responseData)

    return responsePDU.ToBytes(), nil
}

func (s *TCPServer) readHoldingRegisters(slave *Slave, pduBytes []byte) ([]byte, error) {
    // Parse the PDU
    pduRead := &pdu.PDU_Read{}
    pduRead.FromBytes(pduBytes)

    startAddress := pduRead.StartingAddress
    quantity := pduRead.Quantity

    // Validate request
    if startAddress+quantity > 65535 || !slave.LegalHoldingRegistersAddress[startAddress] {
        return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
    }

    // Prepare response data
    responseData := make([]byte, 2*quantity)

    for i := 0; i < int(quantity); i++ {
        if !slave.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
            return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
        }
        value := slave.HoldingRegisters[startAddress+uint16(i)]
        binary.BigEndian.PutUint16(responseData[2*i:], value)
    }

    // Create the response PDU
    responsePDU := pdu.NewPDUReadResponse(0x03, responseData)

    return responsePDU.ToBytes(), nil
}

func (s *TCPServer) readInputRegisters(slave *Slave, pduBytes []byte) ([]byte, error) {
    // Parse the PDU
    pduRead := &pdu.PDU_Read{}
    pduRead.FromBytes(pduBytes)

    startAddress := pduRead.StartingAddress
    quantity := pduRead.Quantity

    // Validate request
    if startAddress+quantity > 65535 || !slave.LegalInputRegistersAddress[startAddress] {
        return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
    }

    // Prepare response data
    responseData := make([]byte, 2*quantity)

    for i := 0; i < int(quantity); i++ {
        if !slave.LegalInputRegistersAddress[startAddress+uint16(i)] {
            return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
        }
        value := slave.InputRegisters[startAddress+uint16(i)]
        binary.BigEndian.PutUint16(responseData[2*i:], value)
    }

    // Create the response PDU
    responsePDU := pdu.NewPDUReadResponse(0x04, responseData)

    return responsePDU.ToBytes(), nil
}

func (s *TCPServer) writeSingleCoil(slave *Slave, pduBytes []byte) ([]byte, error) {
    // Parse the PDU
    pduWrite := &pdu.PDU_WriteSingleCoil{}
    pduWrite.FromBytes(pduBytes)

    address := pduWrite.OutputAddress
    value := pduWrite.OutputValue

    // Validate address
    if address >= 65535 || !slave.LegalCoilsAddress[address] {
        return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
    }

    // Update coil value based on the PDU
    if value == 0xFF00 {
        slave.Coils[address] = true
    } else if value == 0x0000 {
        slave.Coils[address] = false
    } else {
        return s.exceptionResponse(pduBytes, 0x03), nil // Illegal Data Value
    }

    // The response is the same as the request for this function
    return pduWrite.ToBytes(), nil
}

func (s *TCPServer) writeMultipleCoils(slave *Slave, pduBytes []byte) ([]byte, error) {
    // Parse the PDU
    pduWrite := &pdu.PDU_WriteMultipleCoils{}
    pduWrite.FromBytes(pduBytes)

    startAddress := pduWrite.StartingAddress
    quantity := pduWrite.QuantityOfOutputs
    byteCount := pduWrite.ByteCount

    // Validate the request
    if int(byteCount) != len(pduWrite.OutputValues) || startAddress+quantity > 65535 || !slave.LegalCoilsAddress[startAddress] {
        return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
    }

    // Update coils based on the PDU
    for i := 0; i < int(quantity); i++ {
        if !slave.LegalCoilsAddress[startAddress+uint16(i)] {
            return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
        }
        if pduWrite.OutputValues[i/8]&(1<<(i%8)) != 0 {
            slave.Coils[startAddress+uint16(i)] = true
        } else {
            slave.Coils[startAddress+uint16(i)] = false
        }
    }

    // Create the response PDU
    responsePDU := pdu.NewPDUWriteMultipleResponse(0x0F, startAddress, quantity)

    return responsePDU.ToBytes(), nil
}

func (s *TCPServer) writeSingleRegister(slave *Slave, pduBytes []byte) ([]byte, error) {
    // Parse the PDU
    pduWrite := &pdu.PDU_WriteSingleRegister{}
    pduWrite.FromBytes(pduBytes)

    address := pduWrite.RegisterAddress
    value := pduWrite.RegisterValue

    // Validate the address
    if address >= 65535 || !slave.LegalHoldingRegistersAddress[address] {
        return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
    }

    // Update the register with the new value
    slave.HoldingRegisters[address] = value

    // The response is the same as the request for this function
    return pduWrite.ToBytes(), nil
}

func (s *TCPServer) writeMultipleRegisters(slave *Slave, pduBytes []byte) ([]byte, error) {
    // Parse the PDU
    pduWrite := &pdu.PDU_WriteMultipleRegisters{}
    pduWrite.FromBytes(pduBytes)

    startAddress := pduWrite.StartingAddress
    quantity := pduWrite.QuantityOfOutputs
    byteCount := pduWrite.ByteCount

    // Validate the request
    if int(byteCount) != len(pduWrite.OutputValues) || startAddress+quantity > 65535 || !slave.LegalHoldingRegistersAddress[startAddress] {
        return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
    }

    // Update registers based on the PDU
    for i := 0; i < int(quantity); i++ {
        if !slave.LegalHoldingRegistersAddress[startAddress+uint16(i)] {
            return s.exceptionResponse(pduBytes, 0x02), nil // Illegal Data Address
        }
        value := binary.BigEndian.Uint16(pduWrite.OutputValues[2*i:])
        slave.HoldingRegisters[startAddress+uint16(i)] = value
    }

    // Create the response PDU
    responsePDU := pdu.NewPDUWriteMultipleResponse(0x10, startAddress, quantity)

    return responsePDU.ToBytes(), nil
}

func (s *TCPServer) exceptionResponse(request []byte, exceptionCode byte) []byte {
    // Extract the necessary fields from the request
    transactionID := binary.BigEndian.Uint16(request[0:2])
    protocolID := binary.BigEndian.Uint16(request[2:4])
    unitID := request[6]
    functionCode := request[7] | 0x80 // Set the highest bit to indicate an error

    // Create the exception PDU
    exceptionPDU := []byte{functionCode, exceptionCode}

    // Create the exception ADU
    exceptionADU := adu.NewTCPADU(transactionID, unitID, exceptionPDU)
    exceptionADU.ProtocolID = protocolID
    exceptionADU.Length = uint16(3) // Length of remaining bytes (Unit ID + Function Code + Exception Code)

    return exceptionADU.ToBytes()
}