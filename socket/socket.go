package socket

import (
	"fmt"
	"net"

	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/client"
	"github.com/ulfaric/gomodbus/protobuf"
	"github.com/ulfaric/gomodbus/server"
	"google.golang.org/protobuf/proto"
)

type Socket struct {
	Port   int
	server server.Server
	client client.Client
}

func NewSocket(port int) *Socket {
	return &Socket{Port: port}
}

func (s *Socket) Start() error {
	address := fmt.Sprintf("%s:%d", "127.0.0.1", s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	gomodbus.Logger.Sugar().Infof("goModBus wrapper started listening on port %d", s.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()
		gomodbus.Logger.Sugar().Infof("goModBus wrapper accepted connection from address %s", conn.LocalAddr().String())
		go s.handleConnection(conn)
	}
}

func (s *Socket) handleConnection(conn net.Conn) {

	for {
		// Read the header
		headerBuffer := make([]byte, 2)
		_, err := conn.Read(headerBuffer)
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("wrapper failed reading header: %v", err)
			return
		}
		header := &protobuf.Header{}
		err = proto.Unmarshal(headerBuffer, header)
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling header: %v", err)
			return
		}

		bodyBuffer := make([]byte, header.Length)
		_, err = conn.Read(bodyBuffer)
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("wrapper failed reading body: %v", err)
			return
		}

		switch header.Function {
		case protobuf.Function_NEW_TCP_SERVER:
			gomodbus.Logger.Debug("wrapper received NEW_TCP_SERVER")
			s.handleNewTCPServer(bodyBuffer)
		case protobuf.Function_NEW_RTU_SERVER:
			gomodbus.Logger.Debug("wrapper received NEW_RTU_SERVER")
			s.handleNewRTUServer(bodyBuffer)
		case protobuf.Function_START:
			gomodbus.Logger.Debug("wrapper received START")
			s.handleStart()
		case protobuf.Function_STOP:
			gomodbus.Logger.Debug("wrapper received STOP")
			s.handleStop()
		case protobuf.Function_ADD_COILS:
			gomodbus.Logger.Debug("wrapper received ADD_COILS")
			s.addCoils(bodyBuffer)
		case protobuf.Function_DELETE_COILS:
			gomodbus.Logger.Debug("wrapper received DELETE_COILS")
			s.deleteCoils(bodyBuffer)
		case protobuf.Function_ADD_DISCRETE_INPUTS:
			gomodbus.Logger.Debug("wrapper received ADD_DISCRETE_INPUTS")
			s.addDiscreteInputs(bodyBuffer)
		case protobuf.Function_DELETE_DISCRETE_INPUTS:
			gomodbus.Logger.Debug("wrapper received DELETE_DISCRETE_INPUTS")
			s.deleteDiscreteInputs(bodyBuffer)
		case protobuf.Function_ADD_HOLDING_REGISTERS:
			gomodbus.Logger.Debug("wrapper received ADD_HOLDING_REGISTERS")
			s.addHoldingRegisters(bodyBuffer)
		case protobuf.Function_DELETE_HOLDING_REGISTERS:
			gomodbus.Logger.Debug("wrapper received DELETE_HOLDING_REGISTERS")
			s.deleteHoldingRegisters(bodyBuffer)
		case protobuf.Function_ADD_INPUT_REGISTERS:
			gomodbus.Logger.Debug("wrapper received ADD_INPUT_REGISTERS")
			s.addInputRegisters(bodyBuffer)
		case protobuf.Function_DELETE_INPUT_REGISTERS:
			gomodbus.Logger.Debug("wrapper received DELETE_INPUT_REGISTERS")
			s.deleteInputRegisters(bodyBuffer)
		default:
			gomodbus.Logger.Sugar().Errorf("wrapper received unknown function: %d", header.Function)
		}
	}
}

func (s *Socket) handleNewTCPServer(bodyBuffer []byte) {
	request := &protobuf.NewTCPServerRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling NEW_TCP_SERVER: %v", err)
		return
	}

	s.server = server.NewTCPServer(
		request.Host,
		int(request.Port),
		request.UseTLS,
		request.ByteOrder,
		request.WordOrder,
		request.CertFile,
		request.KeyFile,
		request.CaFile,
	)
	gomodbus.Logger.Sugar().Debugf("wrapper created new TCP server %s:%d", request.Host, request.Port)
}

func (s *Socket) handleNewRTUServer(bodyBuffer []byte) {

}

func (s *Socket) handleStart() {
	go s.server.Start()
}

func (s *Socket) handleStop() {
	s.server.Stop()
}

func (s *Socket) addCoils(bodyBuffer []byte) {
	request := &protobuf.AddCoilsRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling ADD_COILS: %v", err)
		return
	}
	server.AddCoils(s.server, byte(request.UnitID), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Debugf("wrapper added coils at address %d with values_count %d", request.Address, len(request.Values))
}

func (s *Socket) deleteCoils(bodyBuffer []byte) {
	request := &protobuf.DeleteCoilsRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling DELETE_COILS: %v", err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteCoils(s.server, byte(request.UnitID), addresses)
}

func (s *Socket) addDiscreteInputs(bodyBuffer []byte) {
	request := &protobuf.AddDiscreteInputsRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling ADD_DISCRETE_INPUTS: %v", err)
		return
	}
	server.AddDiscreteInputs(s.server, byte(request.UnitID), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Debugf("wrapper added discrete inputs at address %d with values_count %d", request.Address, len(request.Values))
}

func (s *Socket) deleteDiscreteInputs(bodyBuffer []byte) {
	request := &protobuf.DeleteDiscreteInputsRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling DELETE_DISCRETE_INPUTS: %v", err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteDiscreteInputs(s.server, byte(request.UnitID), addresses)
}

func (s *Socket) addHoldingRegisters(bodyBuffer []byte) {
	request := &protobuf.AddHoldingRegistersRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling ADD_HOLDING_REGISTERS: %v", err)
		return
	}
	server.AddHoldingRegisters(s.server, byte(request.UnitID), uint16(request.Address), request.Values)
}

func (s *Socket) deleteHoldingRegisters(bodyBuffer []byte) {
	request := &protobuf.DeleteHoldingRegistersRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling DELETE_HOLDING_REGISTERS: %v", err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteHoldingRegisters(s.server, byte(request.UnitID), addresses)
}

func (s *Socket) addInputRegisters(bodyBuffer []byte) {
	request := &protobuf.AddInputRegistersRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling ADD_INPUT_REGISTERS: %v", err)
		return
	}
	server.AddInputRegisters(s.server, byte(request.UnitID), uint16(request.Address), request.Values)
}

func (s *Socket) deleteInputRegisters(bodyBuffer []byte) {
	request := &protobuf.DeleteInputRegistersRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling DELETE_INPUT_REGISTERS: %v", err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteInputRegisters(s.server, byte(request.UnitID), addresses)
}
