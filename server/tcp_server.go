package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

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
	mu        sync.Mutex

	UseTLS   bool
	CertFile string
	KeyFile  string
	CAFile   string

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewTCPServer(host string, port int, useTLS bool, byteOrder, wordOrder, certFile, keyFile, caFile string) Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &TCPServer{
		Host:     host,
		Port:     port,
		Slaves:   make(map[byte]*Slave),
		UseTLS:   useTLS,
		CertFile: certFile,
		KeyFile:  keyFile,
		CAFile:   caFile,
		ctx:      ctx,
		cancel:   cancel,
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

func (s *TCPServer) GetSlave(unitID byte) (*Slave, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slave, ok := s.Slaves[unitID]
	if !ok {
		return nil, fmt.Errorf("slave not found")
	}
	return slave, nil
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
		cert, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
		if err != nil {
			gomodbus.Logger.Error("failed to load TLS certificate and key", zap.Error(err))
			return fmt.Errorf("failed to load TLS certificate and key: %v", err)
		}

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

	gomodbus.Logger.Info("Waiting for connections...")
	for {
		select {
		case <-s.ctx.Done():
			gomodbus.Logger.Info("Shutting down server...")
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				gomodbus.Logger.Error("failed to accept connection", zap.Error(err))
				continue
			}
			gomodbus.Logger.Info("Connection accepted from", zap.String("address", conn.RemoteAddr().String()))
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

func (s *TCPServer) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer s.wg.Done()

	buffer := make([]byte, 512)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, err := conn.Read(buffer)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
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

			responsePDU, _ := processRequest(requestADU.PDU, slave)

			responseADU := adu.NewTCPADU(requestADU.TransactionID, requestADU.UnitID, responsePDU)
			_, err = conn.Write(responseADU.ToBytes())
			if err != nil {
				gomodbus.Logger.Error("failed to write response", zap.Error(err))
				return
			}
		}
	}
}
