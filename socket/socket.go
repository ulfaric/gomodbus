package socket

import (
	"context"
	"encoding/binary"
	"net"
	"sync"
	"time"

	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/server"
	proto "google.golang.org/protobuf/proto"
)

// Socket struct represents a Unix socket server for ModBus communication.
type Socket struct {
	addr     string
	listener net.Listener

	server server.Server

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewSocket initializes a new Socket instance.
func NewSocket() *Socket {
	ctx, cancel := context.WithCancel(context.Background())
	return &Socket{
		addr:   "127.0.0.1:1993",
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins listening for incoming connections on the Unix socket.
func (s *Socket) Start() {

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}
	s.listener = listener
	gomodbus.Logger.Sugar().Infof("ModBus TCP socket started on %s", s.addr)
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

// Stop gracefully shuts down the server.
func (s *Socket) Stop() {
	s.cancel()
	s.wg.Wait()
}

// receiveRequest reads and parses a request from the connection.
func receiveRequest(conn net.Conn) (h *Header, body []byte, err error) {
	conn.SetReadDeadline(time.Now().Add(1 * time.Second)) // Set a read deadline
	buffer := make([]byte, 16)
	_, err = conn.Read(buffer)
	if err != nil {
		return nil, nil, err
	}

	headerLength := binary.BigEndian.Uint64(buffer[0:8])
	bodyLength := binary.BigEndian.Uint64(buffer[8:16])

	gomodbus.Logger.Sugar().Infof("headerLength: %d, bodyLength: %d", headerLength, bodyLength)
	headerBuffer := make([]byte, headerLength)
	_, err = conn.Read(headerBuffer)
	if err != nil {
		return nil, nil, err
	}

	header := Header{}
	err = proto.Unmarshal(headerBuffer, &header)
	if err != nil {
		return nil, nil, err
	}

	bodyBuffer := make([]byte, bodyLength)
	_, err = conn.Read(bodyBuffer)
	if err != nil {
		return nil, nil, err
	}

	return &header, bodyBuffer, nil
}

// sendResponse sends a response back to the client.
func sendResponse(conn net.Conn, header *Header, body []byte) error {
	headerBuffer, err := proto.Marshal(header)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to marshal response header: %v", err)
		return err
	}

	buffer := make([]byte, 16)
	binary.BigEndian.PutUint64(buffer[0:8], uint64(len(headerBuffer)))
	if body != nil {
		binary.BigEndian.PutUint64(buffer[8:16], uint64(len(body)))
	}

	buffer = append(buffer, headerBuffer...)
	if body != nil {
		buffer = append(buffer, body...)
	}

	_, err = conn.Write(buffer)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to write response: %v", err)
		return err
	}
	return nil
}

// handleConnection processes incoming requests on the connection.
func (s *Socket) handleConnection(conn net.Conn) {
	defer conn.Close()
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			header, bodyBuffer, err := receiveRequest(conn)
            if err != nil {
                if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
                    continue // Timeout, check context and retry
                }
                gomodbus.Logger.Sugar().Errorf("failed to read from connection: %v", err)
                return
            }

			switch header.Type {
			case RequestType_NewTCPServerRequest:
				gomodbus.Logger.Sugar().Infof("received NewTCPServerRequest")
				s.handleNewTCPServerRequest(bodyBuffer, conn)

			case RequestType_NewSerialServerRequest:
				gomodbus.Logger.Sugar().Infof("received NewSerialServerRequest")
				s.handleNewSerialServerRequest(bodyBuffer, conn)

			case RequestType_NewSlaveRequest:
				gomodbus.Logger.Sugar().Infof("received NewSlaveRequest")
				s.handleNewSlaveRequest(bodyBuffer, conn)

			case RequestType_RemoveSlaveRequest:
				gomodbus.Logger.Sugar().Infof("received RemoveSlaveRequest")
				s.handleRemoveSlaveRequest(bodyBuffer, conn)

			case RequestType_AddCoilsRequest:
				gomodbus.Logger.Sugar().Infof("received AddCoilsRequest")
				s.handleAddCoilsRequest(bodyBuffer, conn)

			case RequestType_DeleteCoilsRequest:
				gomodbus.Logger.Sugar().Infof("received DeleteCoilsRequest")
				s.handleDeleteCoilsRequest(bodyBuffer, conn)

			case RequestType_AddDiscreteInputsRequest:
				gomodbus.Logger.Sugar().Infof("received AddDiscreteInputsRequest")
				s.handleAddDiscreteInputsRequest(bodyBuffer, conn)

			case RequestType_DeleteDiscreteInputsRequest:
				gomodbus.Logger.Sugar().Infof("received DeleteDiscreteInputsRequest")
				s.handleDeleteDiscreteInputsRequest(bodyBuffer, conn)

			case RequestType_AddHoldingRegistersRequest:
				gomodbus.Logger.Sugar().Infof("received AddHoldingRegistersRequest")
				s.handleAddHoldingRegistersRequest(bodyBuffer, conn)

			case RequestType_DeleteHoldingRegistersRequest:
				gomodbus.Logger.Sugar().Infof("received DeleteHoldingRegistersRequest")
				s.handleDeleteHoldingRegistersRequest(bodyBuffer, conn)

			case RequestType_AddInputRegistersRequest:
				gomodbus.Logger.Sugar().Infof("received AddInputRegistersRequest")
				s.handleAddInputRegistersRequest(bodyBuffer, conn)

			case RequestType_DeleteInputRegistersRequest:
				gomodbus.Logger.Sugar().Infof("received DeleteInputRegistersRequest")
				s.handleDeleteInputRegistersRequest(bodyBuffer, conn)

			case RequestType_StartServerRequest:
				gomodbus.Logger.Sugar().Infof("received StartServerRequest")
				s.handleStartServerRequest(conn)

			case RequestType_StopServerRequest:
				gomodbus.Logger.Sugar().Infof("received StopServerRequest")
				s.handleStopServerRequest(conn)
			}
		}
	}
}

// handleNewTCPServerRequest processes a request to create a new TCP server.
func (s *Socket) handleNewTCPServerRequest(bodyBuffer []byte, conn net.Conn) {
	var request TCPServerRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal NewTCPServerRequest: %v", err)
		return
	}

	s.server = server.NewTCPServer(
		request.Host,
		request.Port,
		request.UseTls,
		request.ByteOrder,
		request.WordOrder,
		request.CertFile,
		request.KeyFile,
		request.CaFile,
	)

	gomodbus.Logger.Sugar().Infof("created new TCP server on %s:%d", request.Host, request.Port)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleNewSerialServerRequest processes a request to create a new Serial server.
func (s *Socket) handleNewSerialServerRequest(bodyBuffer []byte, conn net.Conn) {
	var request SerialServerRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal NewSerialServerRequest: %v", err)
		return
	}

	s.server = server.NewSerialServer(
		request.Port,
		int(request.Baudrate),
		byte(request.Databits),
		byte(request.Parity),
		byte(request.Stopbits),
		request.ByteOrder,
		request.WordOrder,
	)

	gomodbus.Logger.Sugar().Infof("created new Serial server on %s", request.Port)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleNewSlaveRequest processes a request to add a new slave.
func (s *Socket) handleNewSlaveRequest(bodyBuffer []byte, conn net.Conn) {
	var request SlaveRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal NewSlaveRequest: %v", err)
		return
	}

	s.server.AddSlave(byte(request.UnitId))
	gomodbus.Logger.Sugar().Infof("added slave with unit ID %d", request.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleRemoveSlaveRequest processes a request to remove a slave.
func (s *Socket) handleRemoveSlaveRequest(bodyBuffer []byte, conn net.Conn) {
	var request SlaveRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal RemoveSlaveRequest: %v", err)
		return
	}

	s.server.RemoveSlave(byte(request.UnitId))
	gomodbus.Logger.Sugar().Infof("removed slave with unit ID %d", request.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleAddCoilsRequest processes a request to add coils to a slave.
func (s *Socket) handleAddCoilsRequest(bodyBuffer []byte, conn net.Conn) {
	var request CoilsRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddCoilsRequest: %v", err)
		return
	}

	server.AddCoils(s.server, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Infof("added %d coils to slave with unit ID %d", len(request.Values), request.SlaveRequest.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleDeleteCoilsRequest processes a request to delete coils from a slave.
func (s *Socket) handleDeleteCoilsRequest(bodyBuffer []byte, conn net.Conn) {
	var request DeleteRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal DeleteCoilsRequest: %v", err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteCoils(s.server, byte(request.SlaveRequest.UnitId), addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d coils from slave with unit ID %d", len(request.Addresses), request.SlaveRequest.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleAddDiscreteInputsRequest processes a request to add discrete inputs to a slave.
func (s *Socket) handleAddDiscreteInputsRequest(bodyBuffer []byte, conn net.Conn) {
	var request DiscreteInputsRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddDiscreteInputsRequest: %v", err)
		return
	}
	server.AddDiscreteInputs(s.server, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Infof("added %d discrete inputs to slave with unit ID %d", len(request.Values), request.SlaveRequest.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleDeleteDiscreteInputsRequest processes a request to delete discrete inputs from a slave.
func (s *Socket) handleDeleteDiscreteInputsRequest(bodyBuffer []byte, conn net.Conn) {
	var request DeleteRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal DeleteDiscreteInputsRequest: %v", err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteDiscreteInputs(s.server, byte(request.SlaveRequest.UnitId), addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d discrete inputs from slave with unit ID %d", len(request.Addresses), request.SlaveRequest.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleAddHoldingRegistersRequest processes a request to add holding registers to a slave.
func (s *Socket) handleAddHoldingRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request HoldingRegistersRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddHoldingRegistersRequest: %v", err)
		return
	}
	server.AddHoldingRegisters(s.server, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Infof("added %d holding registers to slave with unit ID %d", len(request.Values), request.SlaveRequest.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleDeleteHoldingRegistersRequest processes a request to delete holding registers from a slave.
func (s *Socket) handleDeleteHoldingRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request DeleteRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal DeleteHoldingRegistersRequest: %v", err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteHoldingRegisters(s.server, byte(request.SlaveRequest.UnitId), addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d holding registers from slave with unit ID %d", len(request.Addresses), request.SlaveRequest.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleAddInputRegistersRequest processes a request to add input registers to a slave.
func (s *Socket) handleAddInputRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request InputRegistersRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddInputRegistersRequest: %v", err)
		return
	}
	server.AddInputRegisters(s.server, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Infof("added %d input registers to slave with unit ID %d", len(request.Values), request.SlaveRequest.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleDeleteInputRegistersRequest processes a request to delete input registers from a slave.
func (s *Socket) handleDeleteInputRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request DeleteRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal DeleteInputRegistersRequest: %v", err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteInputRegisters(s.server, byte(request.SlaveRequest.UnitId), addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d input registers from slave with unit ID %d", len(request.Addresses), request.SlaveRequest.UnitId)

	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}


// handleStartServerRequest processes a request to start the server.
func (s *Socket) handleStartServerRequest(conn net.Conn) {
	go s.server.Start()
	gomodbus.Logger.Sugar().Infof("started ModBus server")
	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}

// handleStopServerRequest processes a request to stop the server.
func (s *Socket) handleStopServerRequest(conn net.Conn) {
	s.server.Stop()
	gomodbus.Logger.Sugar().Infof("stopped ModBus server")
	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	sendResponse(conn, response, nil)
}