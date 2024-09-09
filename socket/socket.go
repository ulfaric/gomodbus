package wrapper

import (
	"fmt"
	"net"

	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/client"
	"github.com/ulfaric/gomodbus/protobuf"
	"github.com/ulfaric/gomodbus/server"
	"google.golang.org/protobuf/proto"
)

type Wrapper struct {
	Port   int
	server server.Server
	client client.Client
}

func NewWrapper(port int) *Wrapper {
	return &Wrapper{Port: port}
}

func (w *Wrapper) Start() error {
	address := fmt.Sprintf("%s:%d", "127.0.0.1", w.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()
	gomodbus.Logger.Sugar().Infof("goModBus wrapper started listening on port %d", w.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		defer conn.Close()
		gomodbus.Logger.Sugar().Infof("goModBus wrapper accepted connection from address %s", conn.LocalAddr().String())
		go w.handleConnection(conn)
	}
}

func (w *Wrapper) handleConnection(conn net.Conn) {

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
		w.handleNewTCPServer(bodyBuffer)
	case protobuf.Function_NEW_RTU_SERVER:
		gomodbus.Logger.Debug("wrapper received NEW_RTU_SERVER")
		w.handleNewRTUServer(bodyBuffer)
	case protobuf.Function_START:
		gomodbus.Logger.Debug("wrapper received START")
		w.handleStart()
	case protobuf.Function_STOP:
		gomodbus.Logger.Debug("wrapper received STOP")
		w.handleStop()
	case protobuf.Function_ADD_COIL:
		gomodbus.Logger.Debug("wrapper received ADD_COIL")
		w.addCoil(bodyBuffer)
	case protobuf.Function_ADD_COILS:
		gomodbus.Logger.Debug("wrapper received ADD_COILS")
		w.addCoils(bodyBuffer)
	}
}

func (w *Wrapper) handleNewTCPServer(bodyBuffer []byte) {
	request := &protobuf.NewTCPServerRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling NEW_TCP_SERVER: %v", err)
		return
	}

	w.server = server.NewTCPServer(
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

func (w *Wrapper) handleNewRTUServer(bodyBuffer []byte) {

}

func (w *Wrapper) handleStart() {
	go w.server.Start()
}

func (w *Wrapper) handleStop() {
	w.server.Stop()
}

func (w *Wrapper) addCoil(bodyBuffer []byte) {
	request := &protobuf.AddCoilRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling ADD_COIL: %v", err)
		return
	}

	server.AddCoil(w.server, byte(request.UnitID), uint16(request.Address), request.Value)
	gomodbus.Logger.Sugar().Debugf("wrapper added coil at address %d with value %v", request.Address, request.Value)
}

func (w *Wrapper) addCoils(bodyBuffer []byte) {
	request := &protobuf.AddCoilsRequest{}
	err := proto.Unmarshal(bodyBuffer, request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("wrapper failed unmarshalling ADD_COILS: %v", err)
		return
	}
	server.AddCoils(w.server, byte(request.UnitID), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Debugf("wrapper added coils at address %d with values_count %d", request.Address, len(request.Values))
}
