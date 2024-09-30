package socket

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	"github.com/tarm/serial"
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/client"
	"github.com/ulfaric/gomodbus/server"
	proto "google.golang.org/protobuf/proto"
)

// Socket struct represents a Unix socket server for ModBus communication.
type Socket struct {
	addr     string
	listener net.Listener

	server server.Server
	client client.Client

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
			gomodbus.Logger.Info("Shutting down ModBus TCP socket...")
			return
		default:
			conn, err := s.listener.Accept()
			if err != nil {
				gomodbus.Logger.Sugar().Errorf("failed to accept connection: %v", err)
				continue
			}
			gomodbus.Logger.Sugar().Debugf("connection accepted: %v", conn.LocalAddr())
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
	buffer := make([]byte, 16)
	_, err = io.ReadFull(conn, buffer)
	if err != nil {
		return nil, nil, err
	}

	headerLength := binary.BigEndian.Uint64(buffer[0:8])
	bodyLength := binary.BigEndian.Uint64(buffer[8:16])

	gomodbus.Logger.Sugar().Debugf("headerLength: %d, bodyLength: %d", headerLength, bodyLength)
	headerBuffer := make([]byte, headerLength)
	_, err = io.ReadFull(conn, headerBuffer)
	if err != nil {
		return nil, nil, err
	}

	header := Header{}
	err = proto.Unmarshal(headerBuffer, &header)
	if err != nil {
		return nil, nil, err
	}

	bodyBuffer := make([]byte, bodyLength)
	_, err = io.ReadFull(conn, bodyBuffer)
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

func (*Socket) sendACK(conn net.Conn) {
	response := &Header{
		Type:   RequestType_ACK,
		Length: 0,
	}
	gomodbus.Logger.Sugar().Debugf("sending ACK")
	sendResponse(conn, response, nil)
}

func (*Socket) sendNACK(conn net.Conn, err error) {
	response := &Header{
		Type:   RequestType_NACK,
		Length: 0,
		Error:  err.Error(),
	}
	gomodbus.Logger.Sugar().Debugf("sending NACK")
	sendResponse(conn, response, nil)
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
			conn.SetReadDeadline(time.Now().Add(5 * time.Second)) // Set a read deadline
			header, bodyBuffer, err := receiveRequest(conn)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue // Timeout, check context and retry
				}
				gomodbus.Logger.Sugar().Errorf("failed to read from connection: %v", err)
				return
			}

			switch header.Type {

			// Handle sever functions
			case RequestType_NewTCPServerRequest:
				gomodbus.Logger.Sugar().Debugf("received NewTCPServerRequest")
				s.handleNewTCPServerRequest(bodyBuffer, conn)

			case RequestType_NewSerialServerRequest:
				gomodbus.Logger.Sugar().Debugf("received NewSerialServerRequest")
				s.handleNewSerialServerRequest(bodyBuffer, conn)

			case RequestType_NewSlaveRequest:
				gomodbus.Logger.Sugar().Debugf("received NewSlaveRequest")
				s.handleNewSlaveRequest(bodyBuffer, conn)

			case RequestType_RemoveSlaveRequest:
				gomodbus.Logger.Sugar().Debugf("received RemoveSlaveRequest")
				s.handleRemoveSlaveRequest(bodyBuffer, conn)

			case RequestType_AddCoilsRequest:
				gomodbus.Logger.Sugar().Debugf("received AddCoilsRequest")
				s.handleAddCoilsRequest(bodyBuffer, conn)

			case RequestType_DeleteCoilsRequest:
				gomodbus.Logger.Sugar().Debugf("received DeleteCoilsRequest")
				s.handleDeleteCoilsRequest(bodyBuffer, conn)

			case RequestType_AddDiscreteInputsRequest:
				gomodbus.Logger.Sugar().Debugf("received AddDiscreteInputsRequest")
				s.handleAddDiscreteInputsRequest(bodyBuffer, conn)

			case RequestType_DeleteDiscreteInputsRequest:
				gomodbus.Logger.Sugar().Debugf("received DeleteDiscreteInputsRequest")
				s.handleDeleteDiscreteInputsRequest(bodyBuffer, conn)

			case RequestType_AddHoldingRegistersRequest:
				gomodbus.Logger.Sugar().Debugf("received AddHoldingRegistersRequest")
				s.handleAddHoldingRegistersRequest(bodyBuffer, conn)

			case RequestType_DeleteHoldingRegistersRequest:
				gomodbus.Logger.Sugar().Debugf("received DeleteHoldingRegistersRequest")
				s.handleDeleteHoldingRegistersRequest(bodyBuffer, conn)

			case RequestType_AddInputRegistersRequest:
				gomodbus.Logger.Sugar().Debugf("received AddInputRegistersRequest")
				s.handleAddInputRegistersRequest(bodyBuffer, conn)

			case RequestType_DeleteInputRegistersRequest:
				gomodbus.Logger.Sugar().Debugf("received DeleteInputRegistersRequest")
				s.handleDeleteInputRegistersRequest(bodyBuffer, conn)

			case RequestType_StartServerRequest:
				gomodbus.Logger.Sugar().Debugf("received StartServerRequest")
				s.handleStartServerRequest(conn)

			case RequestType_StopServerRequest:
				gomodbus.Logger.Sugar().Debugf("received StopServerRequest")
				s.handleStopServerRequest(conn)

			// Handle Client functions
			case RequestType_NewTCPClientRequest:
				gomodbus.Logger.Sugar().Debugf("received NewTCPClientRequest")
				s.handleNewTCPClientRequest(bodyBuffer, conn)

			case RequestType_NewSerialClientRequest:
				gomodbus.Logger.Sugar().Debugf("received NewSerialClientRequest")
				s.handleNewSerialClientRequest(bodyBuffer, conn)

			case RequestType_ReadCoilsRequest:
				gomodbus.Logger.Sugar().Debugf("received ReadCoilsRequest")
				s.handleReadCoilsRequest(bodyBuffer, conn)

			case RequestType_ReadDiscreteInputsRequest:
				gomodbus.Logger.Sugar().Debugf("received ReadDiscreteInputsRequest")
				s.handleReadDiscreteInputsRequest(bodyBuffer, conn)

			case RequestType_ReadHoldingRegistersRequest:
				gomodbus.Logger.Sugar().Debugf("received ReadHoldingRegistersRequest")
				s.handleReadHoldingRegistersRequest(bodyBuffer, conn)

			case RequestType_ReadInputRegistersRequest:
				gomodbus.Logger.Sugar().Debugf("received ReadInputRegistersRequest")
				s.handleReadInputRegistersRequest(bodyBuffer, conn)

			case RequestType_WriteCoilRequest:
				gomodbus.Logger.Sugar().Debugf("received WriteCoilRequest")
				s.handleWriteCoilRequest(bodyBuffer, conn)

			case RequestType_WriteCoilsRequest:
				gomodbus.Logger.Sugar().Debugf("received WriteCoilsRequest")
				s.handleWriteCoilsRequest(bodyBuffer, conn)

			case RequestType_WriteRegisterRequest:
				gomodbus.Logger.Sugar().Debugf("received WriteRegisterRequest")
				s.handleWriteRegisterRequest(bodyBuffer, conn)

			case RequestType_WriteRegistersRequest:
				gomodbus.Logger.Sugar().Debugf("received WriteRegistersRequest")
				s.handleWriteRegistersRequest(bodyBuffer, conn)

			case RequestType_ConnectRequest:
				gomodbus.Logger.Sugar().Debugf("received ConnectRequest")
				s.handleConnectRequest(conn)

			case RequestType_DisconnectRequest:
				gomodbus.Logger.Sugar().Debugf("received DisconnectRequest")
				s.handleDisconnectRequest(conn)
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
		s.sendNACK(conn, err)
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

	gomodbus.Logger.Sugar().Debugf("created new TCP server on %s:%d", request.Host, request.Port)

	s.sendACK(conn)
}

// handleNewSerialServerRequest processes a request to create a new Serial server.
func (s *Socket) handleNewSerialServerRequest(bodyBuffer []byte, conn net.Conn) {
	var request SerialServerRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal NewSerialServerRequest: %v", err)
		s.sendNACK(conn, err)
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

	gomodbus.Logger.Sugar().Debugf("created new Serial server on %s", request.Port)

	s.sendACK(conn)
}

// handleNewSlaveRequest processes a request to add a new slave.
func (s *Socket) handleNewSlaveRequest(bodyBuffer []byte, conn net.Conn) {
	var request SlaveRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal NewSlaveRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.server.AddSlave(byte(request.UnitId))
	gomodbus.Logger.Sugar().Infof("added slave with unit ID %d", request.UnitId)

	s.sendACK(conn)
}

// handleRemoveSlaveRequest processes a request to remove a slave.
func (s *Socket) handleRemoveSlaveRequest(bodyBuffer []byte, conn net.Conn) {
	var request SlaveRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal RemoveSlaveRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.server.RemoveSlave(byte(request.UnitId))
	gomodbus.Logger.Sugar().Infof("removed slave with unit ID %d", request.UnitId)

	s.sendACK(conn)
}

// handleAddCoilsRequest processes a request to add coils to a slave.
func (s *Socket) handleAddCoilsRequest(bodyBuffer []byte, conn net.Conn) {
	var request CoilsRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddCoilsRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	server.AddCoils(s.server, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Infof("added %d coils to slave with unit ID %d", len(request.Values), request.SlaveRequest.UnitId)

	s.sendACK(conn)
}

// handleDeleteCoilsRequest processes a request to delete coils from a slave.
func (s *Socket) handleDeleteCoilsRequest(bodyBuffer []byte, conn net.Conn) {
	var request DeleteRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal DeleteCoilsRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteCoils(s.server, byte(request.SlaveRequest.UnitId), addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d coils from slave with unit ID %d", len(request.Addresses), request.SlaveRequest.UnitId)

	s.sendACK(conn)
}

// handleAddDiscreteInputsRequest processes a request to add discrete inputs to a slave.
func (s *Socket) handleAddDiscreteInputsRequest(bodyBuffer []byte, conn net.Conn) {
	var request DiscreteInputsRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddDiscreteInputsRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}
	server.AddDiscreteInputs(s.server, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Infof("added %d discrete inputs to slave with unit ID %d", len(request.Values), request.SlaveRequest.UnitId)

	s.sendACK(conn)
}

// handleDeleteDiscreteInputsRequest processes a request to delete discrete inputs from a slave.
func (s *Socket) handleDeleteDiscreteInputsRequest(bodyBuffer []byte, conn net.Conn) {
	var request DeleteRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal DeleteDiscreteInputsRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteDiscreteInputs(s.server, byte(request.SlaveRequest.UnitId), addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d discrete inputs from slave with unit ID %d", len(request.Addresses), request.SlaveRequest.UnitId)

	s.sendACK(conn)
}

// handleAddHoldingRegistersRequest processes a request to add holding registers to a slave.
func (s *Socket) handleAddHoldingRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request HoldingRegistersRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddHoldingRegistersRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}
	server.AddHoldingRegisters(s.server, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Infof("added %d holding registers to slave with unit ID %d", len(request.Values), request.SlaveRequest.UnitId)

	s.sendACK(conn)
}

// handleDeleteHoldingRegistersRequest processes a request to delete holding registers from a slave.
func (s *Socket) handleDeleteHoldingRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request DeleteRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal DeleteHoldingRegistersRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteHoldingRegisters(s.server, byte(request.SlaveRequest.UnitId), addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d holding registers from slave with unit ID %d", len(request.Addresses), request.SlaveRequest.UnitId)

	s.sendACK(conn)
}

// handleAddInputRegistersRequest processes a request to add input registers to a slave.
func (s *Socket) handleAddInputRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request InputRegistersRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddInputRegistersRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}
	server.AddInputRegisters(s.server, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Values)
	gomodbus.Logger.Sugar().Infof("added %d input registers to slave with unit ID %d", len(request.Values), request.SlaveRequest.UnitId)

	s.sendACK(conn)
}

// handleDeleteInputRegistersRequest processes a request to delete input registers from a slave.
func (s *Socket) handleDeleteInputRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request DeleteRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal DeleteInputRegistersRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}
	addresses := make([]uint16, len(request.Addresses))
	for i, address := range request.Addresses {
		addresses[i] = uint16(address)
	}
	server.DeleteInputRegisters(s.server, byte(request.SlaveRequest.UnitId), addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d input registers from slave with unit ID %d", len(request.Addresses), request.SlaveRequest.UnitId)

	s.sendACK(conn)
}

// handleStartServerRequest processes a request to start the server.
func (s *Socket) handleStartServerRequest(conn net.Conn) {
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.server.Start()
	}()

	err := <-errChan
	if err != nil {
		s.sendNACK(conn, err)
		return
	}
	gomodbus.Logger.Sugar().Infof("started ModBus server")
	s.sendACK(conn)
}

// handleStopServerRequest processes a request to stop the server.
func (s *Socket) handleStopServerRequest(conn net.Conn) {
	s.server.Stop()
	gomodbus.Logger.Sugar().Infof("stopped ModBus server")
	s.sendACK(conn)
}

// handleNewTCPClientRequest processes a request to create a new TCP client.
func (s *Socket) handleNewTCPClientRequest(bodyBuffer []byte, conn net.Conn) {
	var request TCPClientRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal NewTCPClientRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.client = client.NewTCPClient(
		request.Host,
		int(request.Port),
		request.UseTls,
		request.CertFile,
		request.KeyFile,
		request.CaFile,
	)

	gomodbus.Logger.Sugar().Infof("created new TCP client to %s:%d", request.Host, request.Port)

	s.sendACK(conn)
}

// handleNewSerialClientRequest processes a request to create a new Serial client.
func (s *Socket) handleNewSerialClientRequest(bodyBuffer []byte, conn net.Conn) {
	var request SerialClientRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal NewSerialClientRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.client = client.NewSerialClient(
		request.Port,
		int(request.Baudrate),
		byte(request.Databits),
		serial.Parity(request.Parity),
		serial.StopBits(request.Stopbits),
	)

	gomodbus.Logger.Sugar().Infof("created new Serial client to %s", request.Port)

	s.sendACK(conn)
}

// handleReadCoilsRequest processes a request to read coils from a slave.
func (s *Socket) handleReadCoilsRequest(bodyBuffer []byte, conn net.Conn) {
	var request ReadCoils
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal ReadCoilsRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	response, err := client.ReadCoils(s.client, byte(request.SlaveRequest.UnitId), uint16(request.Address), uint16(request.Count))
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to read coils: %v", err)
		s.sendNACK(conn, err)
		return
	}

	readCoilsResponse := &ReadCoilsResponse{
		Values: response,
	}
	readCoilsResponseBuffer, err := proto.Marshal(readCoilsResponse)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to marshal ReadCoilsResponse: %v", err)
		s.sendNACK(conn, err)
		return
	}

	header := &Header{
		Type:   RequestType_ACK,
		Length: uint64(len(readCoilsResponseBuffer)),
	}
	sendResponse(conn, header, readCoilsResponseBuffer)
}

// handleReadDiscreteInputsRequest processes a request to read discrete inputs from a slave.
func (s *Socket) handleReadDiscreteInputsRequest(bodyBuffer []byte, conn net.Conn) {
	var request ReadDiscreteInputs
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal ReadDiscreteInputsRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	response, err := client.ReadDiscreteInputs(s.client, byte(request.SlaveRequest.UnitId), uint16(request.Address), uint16(request.Count))
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to read discrete inputs: %v", err)
		s.sendNACK(conn, err)
		return
	}

	readDiscreteInputsResponse := &ReadDiscreteInputsResponse{
		Values: response,
	}
	readDiscreteInputsResponseBuffer, err := proto.Marshal(readDiscreteInputsResponse)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to marshal ReadDiscreteInputsResponse: %v", err)
		s.sendNACK(conn, err)
		return
	}

	header := &Header{
		Type:   RequestType_ACK,
		Length: uint64(len(readDiscreteInputsResponseBuffer)),
	}
	sendResponse(conn, header, readDiscreteInputsResponseBuffer)
}

// handleReadHoldingRegistersRequest processes a request to read holding registers from a slave.
func (s *Socket) handleReadHoldingRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request ReadHoldingRegisters
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal ReadHoldingRegistersRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	response, err := client.ReadHoldingRegisters(s.client, byte(request.SlaveRequest.UnitId), uint16(request.Address), uint16(request.Count))
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to read holding registers: %v", err)
		s.sendNACK(conn, err)
		return
	}

	responseBytes := make([]byte, 0)
	for _, value := range response {
		responseBytes = binary.BigEndian.AppendUint16(responseBytes, value)
	}

	readHoldingRegistersResponse := &ReadHoldingRegistersResponse{
		Values: responseBytes,
	}
	readHoldingRegistersResponseBuffer, err := proto.Marshal(readHoldingRegistersResponse)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to marshal ReadHoldingRegistersResponse: %v", err)
		s.sendNACK(conn, err)
		return
	}

	header := &Header{
		Type:   RequestType_ACK,
		Length: uint64(len(readHoldingRegistersResponseBuffer)),
	}
	sendResponse(conn, header, readHoldingRegistersResponseBuffer)
}

// handleReadInputRegistersRequest processes a request to read input registers from a slave.
func (s *Socket) handleReadInputRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request ReadInputRegisters
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal ReadInputRegistersRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	response, err := client.ReadInputRegisters(s.client, byte(request.SlaveRequest.UnitId), uint16(request.Address), uint16(request.Count))
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to read input registers: %v", err)
		s.sendNACK(conn, err)
		return
	}

	responseBytes := make([]byte, 0)
	for _, value := range response {
		responseBytes = binary.BigEndian.AppendUint16(responseBytes, value)
	}

	readInputRegistersResponse := &ReadInputRegistersResponse{
		Values: responseBytes,
	}
	readInputRegistersResponseBuffer, err := proto.Marshal(readInputRegistersResponse)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to marshal ReadInputRegistersResponse: %v", err)
		s.sendNACK(conn, err)
		return
	}

	header := &Header{
		Type:   RequestType_ACK,
		Length: uint64(len(readInputRegistersResponseBuffer)),
	}
	sendResponse(conn, header, readInputRegistersResponseBuffer)
}

// handleWriteCoilRequest processes a request to write a single coil to a slave.
func (s *Socket) handleWriteCoilRequest(bodyBuffer []byte, conn net.Conn) {
	var request WriteCoil
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal WriteCoilRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	err = client.WriteSingleCoil(s.client, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Value)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to write coil: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)
}

// handleWriteCoilsRequest processes a request to write multiple coils to a slave.
func (s *Socket) handleWriteCoilsRequest(bodyBuffer []byte, conn net.Conn) {
	var request WriteCoils
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal WriteCoilsRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}
	err = client.WriteMultipleCoils(s.client, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Values)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to write coils: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)
}

// handleWriteRegisterRequest processes a request to write a single register to a slave.
func (s *Socket) handleWriteRegisterRequest(bodyBuffer []byte, conn net.Conn) {
	var request WriteRegister
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal WriteRegisterRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	err = client.WriteSingleRegister(s.client, byte(request.SlaveRequest.UnitId), uint16(request.Address), request.Value)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to write register: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)
}

// handleWriteRegistersRequest processes a request to write multiple registers to a slave.
func (s *Socket) handleWriteRegistersRequest(bodyBuffer []byte, conn net.Conn) {
	var request WriteRegisters
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal WriteRegistersRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	err = client.WriteMultipleRegisters(s.client, byte(request.SlaveRequest.UnitId), uint16(request.Address), uint16(len(request.Values)/2), request.Values)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to write registers: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)
}

// handleConnectRequest processes a request to connect to a slave.
func (s *Socket) handleConnectRequest(conn net.Conn) {

	err := s.client.Connect()
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to connect: %v", err)
		s.sendNACK(conn, err)
		return
	}

	gomodbus.Logger.Sugar().Infof("ModBus client connected.")

	s.sendACK(conn)
}

// handleDisconnectRequest processes a request to disconnect from a slave.
func (s *Socket) handleDisconnectRequest(conn net.Conn) {

	err := s.client.Disconnect()
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to disconnect: %v", err)
		s.sendNACK(conn, err)
		return
	}

	gomodbus.Logger.Sugar().Infof("ModBus client disconnected.")

	s.sendACK(conn)
}
