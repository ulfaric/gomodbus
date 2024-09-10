package socket

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"net"
	"sync"

	"github.com/ulfaric/gomodbus"
)

type Socket struct {
	addr     *net.UnixAddr
	listener *net.UnixListener

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewSocket(conn net.Conn) *Socket {
	return &Socket{
		addr: &net.UnixAddr{
			Name: "/tmp/modbus.sock",
			Net:  "unix",
		},
	}
}

func (s *Socket) Start() {
	listener, err := net.ListenUnix("unix", s.addr)
	if err != nil {
		panic(err)
	}
	s.listener = listener
	gomodbus.Logger.Sugar().Infof("ModBus unix socket started on %s", s.addr.Name)
	defer s.listener.Close()

	for {
		select {
		case <-s.ctx.Done():
			gomodbus.Logger.Info("Shutting down server...")
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to accept connection: %v", err)
				continue
			}
			gomodbus.Logger.Sugar().Infof("connection accepted: %v", conn.LocalAddr())
			s.wg.Add(1)
			go s.handleConnection(conn)
		}
	}
}

func (s *Socket) Stop() {
	s.cancel()
	s.wg.Wait()
}

func (s *Socket) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			buffer := make([]byte, 16)
			_, err := conn.Read(buffer)
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to read pre-header from connection: %v", err)
				continue
			}

			headerLength := binary.BigEndian.Uint64(buffer[0:7])
			bodyLength := binary.BigEndian.Uint64(buffer[8:15])

			headerBuffer := make([]byte, headerLength)
			_, err = conn.Read(headerBuffer)
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to read header from connection: %v", err)
				continue
			}

			header := Header{}
			err = json.Unmarshal(headerBuffer, &header)
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to unmarshal header: %v", err)
				continue
			}

			bodyBuffer := make([]byte, bodyLength)
			_, err = conn.Read(bodyBuffer)
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to read body from connection: %v", err)
			}

			switch header.Type {
			case NewTCPServerRequest:
				gomodbus.Logger.Sugar().Infof("NewTCPServerRequest: %v", header.Type)
			case NewSerialServerRequest:
				gomodbus.Logger.Sugar().Infof("NewSerialServerRequest: %v", header.Type)
			case NewSlaveRequest:
				gomodbus.Logger.Sugar().Infof("NewSlaveRequest: %v", header.Type)
			case RemoveSlaveRequest:
				gomodbus.Logger.Sugar().Infof("RemoveSlaveRequest: %v", header.Type)
			case AddCoilsRequest:
				gomodbus.Logger.Sugar().Infof("AddCoilsRequest: %v", header.Type)
			case DeleteCoilsRequest:
				gomodbus.Logger.Sugar().Infof("DeleteCoilsRequest: %v", header.Type)
			case AddDiscreteInputsRequest:
				gomodbus.Logger.Sugar().Infof("AddDiscreteInputsRequest: %v", header.Type)
			case DeleteDiscreteInputsRequest:
				gomodbus.Logger.Sugar().Infof("DeleteDiscreteInputsRequest: %v", header.Type)
			case AddHoldingRegistersRequest:
				gomodbus.Logger.Sugar().Infof("AddHoldingRegistersRequest: %v", header.Type)
			case DeleteHoldingRegistersRequest:
				gomodbus.Logger.Sugar().Infof("DeleteHoldingRegistersRequest: %v", header.Type)
			case AddInputRegistersRequest:
				gomodbus.Logger.Sugar().Infof("AddInputRegistersRequest: %v", header.Type)
			case DeleteInputRegistersRequest:
				gomodbus.Logger.Sugar().Infof("DeleteInputRegistersRequest: %v", header.Type)
			}

		}
	}
}
