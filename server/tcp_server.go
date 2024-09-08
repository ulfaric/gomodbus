package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/adu"
	"github.com/ulfaric/gomodbus/pdu"
	"go.uber.org/zap"
)

type TCPServer struct {
	Host      string
	Port      int
	ByteOrder string
	WordOrder string
	Slaves    map[byte]*Slave
	Status    chan string
	mu        sync.Mutex

	UseTLS   bool
	CertFile string
	KeyFile  string
	CAFile   string
}

func NewTCPServer(host string, port int, useTLS bool, byteOrder, wordOrder, certFile, keyFile, caFile string) *TCPServer {
	return &TCPServer{
		Host:     host,
		Port:     port,
		Slaves:   make(map[byte]*Slave),
		Status:   make(chan string),
		UseTLS:   useTLS,
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
	}
}

func (s *TCPServer) AddSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Slaves[unitID] = &Slave{
		Coils:            make(map[uint16]bool),
		DiscreteInputs:   make(map[uint16]bool),
		HoldingRegisters: make(map[uint16][]byte),
		InputRegisters:   make(map[uint16][]byte),
	}
}

func (s *TCPServer) RemoveSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Slaves, unitID)
}

func (s *TCPServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	var listener net.Listener
	var err error

	if s.UseTLS {
		// Load server TLS certificate and key
		cert, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
		if err != nil {
			gomodbus.Logger.Error("failed to load TLS certificate and key", zap.Error(err))
			return fmt.Errorf("failed to load TLS certificate and key: %v", err)
		}

		// Load custom CA if provided
		var tlsConfig *tls.Config
		if s.CAFile != "" {
			caCert, err := os.ReadFile(s.CAFile)
			if err != nil {
				gomodbus.Logger.Error("failed to read CA file", zap.Error(err))
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
			gomodbus.Logger.Error("failed to listen on %s with TLS", zap.Error(err))
			return fmt.Errorf("failed to listen on %s with TLS: %v", addr, err)
		}
		gomodbus.Logger.Sugar().Infof("Modbus server started with TLS on %s", addr)
	} else {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			gomodbus.Logger.Error("failed to listen on %s", zap.Error(err))
			return fmt.Errorf("failed to listen on %s: %v", addr, err)
		}
		gomodbus.Logger.Sugar().Infof("Modbus server started on %s", addr)
	}
	defer listener.Close()

	s.Status <- "Running"
	gomodbus.Logger.Info("Waiting for connections...")
	for {
		status := <-s.Status
		if status == "Stopping" {
			break
		}
		conn, err := listener.Accept()
		if err != nil {
			gomodbus.Logger.Error("failed to accept connection", zap.Error(err))
			continue
		}
		go s.handleConnection(conn)
	}
	return nil
}

func (s *TCPServer) Stop() error {
	s.Status <- "Stopping"
	return nil
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 256)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			gomodbus.Logger.Error("failed to read from connection", zap.Error(err))
			return
		}

		requestADU := &adu.TCPADU{}
		err = requestADU.FromBytes(buffer[:n])
		if err != nil {
			gomodbus.Logger.Error("failed to parse request", zap.Error(err))
			continue
		}

		slave, ok := s.Slaves[requestADU.UnitID]
		if !ok {
			gomodbus.Logger.Error("slave not found", zap.Uint8("unitID", requestADU.UnitID))
			reponsePDU := pdu.NewPDUErrorResponse(requestADU.PDU[0], 0x04)
			response := adu.NewTCPADU(requestADU.TransactionID, requestADU.UnitID, reponsePDU.ToBytes())
			_, err = conn.Write(response.ToBytes())
			if err != nil {
				gomodbus.Logger.Error("failed to write response", zap.Error(err))
				return
			}
			continue
		}

		responsePDU, err := s.processRequest(requestADU.PDU, slave)
		if err != nil {
			gomodbus.Logger.Error("failed to process request", zap.Error(err))
			responsePDU := pdu.NewPDUErrorResponse(requestADU.PDU[0], 0x07)
			response := adu.NewTCPADU(requestADU.TransactionID, requestADU.UnitID, responsePDU.ToBytes())
			_, err = conn.Write(response.ToBytes())
			if err != nil {
				gomodbus.Logger.Error("failed to write response", zap.Error(err))
				return
			}
			continue
		}

		responseADU := adu.NewTCPADU(requestADU.TransactionID, requestADU.UnitID, responsePDU)
		_, err = conn.Write(responseADU.ToBytes())
		if err != nil {
			gomodbus.Logger.Error("failed to write response", zap.Error(err))
			return
		}
	}
}

func (s *TCPServer) processRequest(request []byte, slave *Slave) ([]byte, error) {
	switch request[0] {
	case gomodbus.ReadCoil:
		responsePDUBytes, err := HandleReadCoils(slave, request)
		if err != nil {
			return nil, err
		}
		return responsePDUBytes, nil
	case gomodbus.ReadDiscreteInput:
		responsePDUBytes, err := HandleReadDiscreteInputs(slave, request)
		if err != nil {
			return nil, err
		}
		return responsePDUBytes, nil
	case gomodbus.ReadHoldingRegister:
		responsePDUBytes, err := HandleReadHoldingRegisters(slave, request)
		if err != nil {
			return nil, err
		}
		return responsePDUBytes, nil
	case gomodbus.ReadInputRegister:
		responsePDUBytes, err := HandleReadInputRegisters(slave, request)
		if err != nil {
			return nil, err
		}
		return responsePDUBytes, nil
	default:
		responsePDU := pdu.NewPDUErrorResponse(request[0], 0x07)
		return responsePDU.ToBytes(), nil
	}
}
