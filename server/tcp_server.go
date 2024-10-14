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
)

type TCPServer struct {
	Host      string
	Port      int32
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

// NewTCPServer creates a new TCP server.
func NewTCPServer(host string, port int32, useTLS bool, certFile, keyFile, caFile string) Server {
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

// AddSlave adds a new slave to the server.
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

// GetSlave gets a slave by unit ID.
func (s *TCPServer) GetSlave(unitID byte) (*Slave, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	slave, ok := s.Slaves[unitID]
	if !ok {
		return nil, fmt.Errorf("slave not found")
	}
	return slave, nil
}

// RemoveSlave removes a slave by unit ID.
func (s *TCPServer) RemoveSlave(unitID byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Slaves, unitID)
}

// Start starts the TCP server.
func (s *TCPServer) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	var listener net.Listener
	var err error

	if s.UseTLS {
		cert, err := tls.LoadX509KeyPair(s.CertFile, s.KeyFile)
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to load TLS certificate and key: %v", err)
			return fmt.Errorf("failed to load TLS certificate and key: %v", err)
		}

		var tlsConfig *tls.Config
		if s.CAFile != "" {
			caCert, err := os.ReadFile(s.CAFile)
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to read CA file: %v", err)
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
			gomodbus.Logger.Sugar().Errorf("failed to listen on %s with TLS: %v", addr, err)
			return fmt.Errorf("failed to listen on %s with TLS: %v", addr, err)
		}
		gomodbus.Logger.Sugar().Infof("Modbus server started with TLS on %s", addr)
	} else {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to listen on %s: %v", addr, err)
			return fmt.Errorf("failed to listen on %s: %v", addr, err)
		}
		gomodbus.Logger.Sugar().Infof("Modbus server started on %s", addr)
	}

	s.wg.Add(1)
	go s.acceptConnection(listener)
	return nil
}

// Stop stops the TCP server.
func (s *TCPServer) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

// acceptConnection accepts connections from the TCP server.
func (s *TCPServer) acceptConnection(listener net.Listener) {
	defer s.wg.Done()
	defer listener.Close()
	gomodbus.Logger.Info("Waiting for connections...")
	for {
		select {
		case <-s.ctx.Done():
			gomodbus.Logger.Info("Shutting down server...")
			return
		default:
			tcpListener, ok := listener.(*net.TCPListener)
			if !ok {
				gomodbus.Logger.Sugar().Errorf("listener is not a TCP listener")
				continue
			}
			tcpListener.SetDeadline(time.Now().Add(time.Second * 1))
			conn, err := tcpListener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && (opErr.Timeout() || opErr.Err.Error() == "use of closed network connection") {
					continue
				}
				gomodbus.Logger.Sugar().Errorf("failed to accept connection: %v", err)
				continue
			}
			gomodbus.Logger.Sugar().Infof("Connection accepted from %s", conn.RemoteAddr().String())
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

// handleConnection handles the connection from the TCP server.
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer s.wg.Done()

	buffer := make([]byte, 512)
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			conn.SetReadDeadline(time.Now().Add(time.Second * 1))
			n, err := conn.Read(buffer)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				return
			}
			gomodbus.Logger.Sugar().Debugf("server received request: %v", buffer[:n])

			requestADU := &adu.TCPADU{}
			err = requestADU.FromBytes(buffer[:n])
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
				continue
			}

			slave, ok := s.Slaves[requestADU.UnitID]
			if !ok {
				gomodbus.Logger.Sugar().Errorf("slave not found: %v", requestADU.UnitID)
				reponsePDU := pdu.NewPDUErrorResponse(requestADU.PDU[0], 0x04)
				response := adu.NewTCPADU(requestADU.TransactionID, requestADU.UnitID, reponsePDU.ToBytes())
				_, err = conn.Write(response.ToBytes())
				if err != nil {
					gomodbus.Logger.Sugar().Errorf("failed to write response: %v", err)
					return
				}
				continue
			}

			responsePDU, _ := processRequest(requestADU.PDU, slave)

			responseADU := adu.NewTCPADU(requestADU.TransactionID, requestADU.UnitID, responsePDU)
			_, err = conn.Write(responseADU.ToBytes())
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to write response: %v", err)
				return
			}
		}
	}
}
