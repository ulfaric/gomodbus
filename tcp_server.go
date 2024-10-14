package gomodbus

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type TCPServer struct {
	Host      string
	Port      int32
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

func NewTCPServer(host string, port int32, useTLS bool, byteOrder, wordOrder, certFile, keyFile, caFile string) Server {
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
			Logger.Sugar().Errorf("failed to load TLS certificate and key: %v", err)
			return fmt.Errorf("failed to load TLS certificate and key: %v", err)
		}

		var tlsConfig *tls.Config
		if s.CAFile != "" {
			caCert, err := os.ReadFile(s.CAFile)
			if err != nil {
				Logger.Sugar().Errorf("failed to read CA file: %v", err)
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
			Logger.Sugar().Errorf("failed to listen on %s with TLS: %v", addr, err)
			return fmt.Errorf("failed to listen on %s with TLS: %v", addr, err)
		}
		Logger.Sugar().Infof("Modbus server started with TLS on %s", addr)
	} else {
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			Logger.Sugar().Errorf("failed to listen on %s: %v", addr, err)
			return fmt.Errorf("failed to listen on %s: %v", addr, err)
		}
		Logger.Sugar().Infof("Modbus server started on %s", addr)
	}

	s.wg.Add(1)
	go s.acceptConnection(listener)
	return nil
}

func (s *TCPServer) Stop() error {
	s.cancel()
	s.wg.Wait()
	return nil
}

func (s *TCPServer) acceptConnection(listener net.Listener) {
	defer s.wg.Done()
	defer listener.Close()
	Logger.Info("Waiting for connections...")
	for {
		select {
		case <-s.ctx.Done():
			Logger.Info("Shutting down server...")
			return
		default:
			tcpListener, ok := listener.(*net.TCPListener)
			if !ok {
				Logger.Sugar().Errorf("listener is not a TCP listener")
				continue
			}
			tcpListener.SetDeadline(time.Now().Add(time.Second * 1))
			conn, err := tcpListener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && (opErr.Timeout() || opErr.Err.Error() == "use of closed network connection") {
					continue
				}
				Logger.Sugar().Errorf("failed to accept connection: %v", err)
				continue
			}
			Logger.Sugar().Infof("Connection accepted from %s", conn.RemoteAddr().String())
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
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
			conn.SetReadDeadline(time.Now().Add(time.Second * 1))
			n, err := conn.Read(buffer)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				return
			}
			Logger.Sugar().Debugf("server received request: %v", buffer[:n])

			requestADU := &TCPADU{}
			err = requestADU.FromBytes(buffer[:n])
			if err != nil {
				Logger.Sugar().Errorf("failed to parse request: %v", err)
				continue
			}

			slave, ok := s.Slaves[requestADU.UnitID]
			if !ok {
				Logger.Sugar().Errorf("slave not found: %v", requestADU.UnitID)
				reponsePDU := NewPDUErrorResponse(requestADU.PDU[0], 0x04)
				response := NewTCPADU(requestADU.TransactionID, requestADU.UnitID, reponsePDU.ToBytes())
				_, err = conn.Write(response.ToBytes())
				if err != nil {
					Logger.Sugar().Errorf("failed to write response: %v", err)
					return
				}
				continue
			}

			responsePDU, _ := processRequest(requestADU.PDU, slave)

			responseADU := NewTCPADU(requestADU.TransactionID, requestADU.UnitID, responsePDU)
			_, err = conn.Write(responseADU.ToBytes())
			if err != nil {
				Logger.Sugar().Errorf("failed to write response: %v", err)
				return
			}
		}
	}
}
