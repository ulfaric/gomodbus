package socket

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
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
	path     string
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
		path:   "./modbus.sock",
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start begins listening for incoming connections on the Unix socket.
func (s *Socket) Start() {

	// Remove the socket file if it already exists
	if _, err := os.Stat(s.path); err == nil {
		gomodbus.Logger.Sugar().Infof("removing existing socket file: %s", s.path)
		os.RemoveAll(s.path)
	}

	// Create a new unix socket listener
	listener, err := net.Listen("unix", s.path)
	if err != nil {
		panic(err)
	}
	s.listener = listener
	gomodbus.Logger.Sugar().Infof("ModBus TCP socket started on %s", s.path)
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

	if bodyLength > 0 {
		bodyBuffer := make([]byte, bodyLength)
		_, err = io.ReadFull(conn, bodyBuffer)
		if err != nil {
			return nil, nil, err
		}
		return &header, bodyBuffer, nil
	}
	
	return &header, nil, nil
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
			gomodbus.Logger.Sugar().Debugf("close modbus socket connection...")
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second)) // Set a read deadline
			header, bodyBuffer, err := receiveRequest(conn)
			if err != nil {
				continue
			}

			switch header.Type {

			// Handle sever functions
			case RequestType_SetTCPServer:
				s.handleSetTCPServer(bodyBuffer, conn)
			case RequestType_SetSerialServer:
				s.handleSetSerialServer(bodyBuffer, conn)
			case RequestType_StartServer:
				s.handleStartServer(conn)
			case RequestType_StopServer:
				s.handleStopServer(conn)

			// Handle slave functions
			case RequestType_AddSlave:
				s.handleAddSlave(bodyBuffer, conn)
			case RequestType_RemoveSlave:
				s.handleRemoveSlave(bodyBuffer, conn)

			// Handle Client functions
			case RequestType_SetTCPClient:
				s.handleSetTCPClient(bodyBuffer, conn)
			case RequestType_SetSerialClient:
				s.handleSetSerialClient(bodyBuffer, conn)
			case RequestType_Connect:
				s.handleConnect(conn)
			case RequestType_Disconnect:
				s.handleDisconnect(conn)

			// Handle Relay functions
			case RequestType_AddRelay:
				s.handleAddRelay(bodyBuffer, conn)
			case RequestType_DeleteRelay:
				s.handleDeleteRelay(bodyBuffer, conn)
			case RequestType_SetRelay:
				s.handleSetRelay(bodyBuffer, conn)
			case RequestType_ReadRelay:
				s.handleReadRelay(bodyBuffer, conn)
			case RequestType_WriteRelay:
				s.handleWriteRelay(bodyBuffer, conn)

			// Handle Register functions
			case RequestType_AddRegister:
				s.handleAddRegister(bodyBuffer, conn)
			case RequestType_DeleteRegister:
				s.handleDeleteRegister(bodyBuffer, conn)
			case RequestType_SetRegister:
				s.handleSetRegister(bodyBuffer, conn)
			case RequestType_ReadRegister:
				s.handleReadRegister(bodyBuffer, conn)
			case RequestType_WriteRegister:
				s.handleWriteRegister(bodyBuffer, conn)
			}
		}
	}
}

// handleSetTCPServer set the modbus server as tcp server
func (s *Socket) handleSetTCPServer(bodyBuffer []byte, conn net.Conn) {

	// unmarshal the request
	var request TCPServerRequest
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal NewTCPServerRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	// create the tcp server
	s.server = server.NewTCPServer(
		request.Host,
		request.Port,
		request.UseTls,
		request.CertFile,
		request.KeyFile,
		request.CaFile,
	)

	gomodbus.Logger.Sugar().Debugf("created new TCP server on %s:%d", request.Host, request.Port)

	// send ack
	s.sendACK(conn)
}

// handleSetSerialServer set the modbus server as serial server
func (s *Socket) handleSetSerialServer(bodyBuffer []byte, conn net.Conn) {
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
	)

	gomodbus.Logger.Sugar().Debugf("created new Serial server on %s", request.Port)

	s.sendACK(conn)
}

// handleStartServer start the modbus server, if successful, it will send an ACK. Otherwise, it will send a NACK.
func (s *Socket) handleStartServer(conn net.Conn) {
	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	err := s.server.Start()
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to start server: %v", err)
		s.sendNACK(conn, err)
		return
	}
	gomodbus.Logger.Sugar().Infof("started ModBus server")
	s.sendACK(conn)
}

// handleStopServer stop the modbus server, if successful, it will send an ACK. Otherwise, it will send a NACK.
func (s *Socket) handleStopServer(conn net.Conn) {
	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	err := s.server.Stop()
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to stop server: %v", err)
		s.sendNACK(conn, err)
		return
	}
	gomodbus.Logger.Sugar().Infof("stopped ModBus server")
	s.sendACK(conn)
}

// handleAddSlaveRequest processes a request to add a new slave.
func (s *Socket) handleAddSlave(bodyBuffer []byte, conn net.Conn) {
	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	var request Slave
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
func (s *Socket) handleRemoveSlave(bodyBuffer []byte, conn net.Conn) {
	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	var request Slave
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

// handleSetTCPClient set the modbus client in Modbus TCP mode
func (s *Socket) handleSetTCPClient(bodyBuffer []byte, conn net.Conn) {
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

// handleSetSerialClient set the modbus client in Modbus Serial mode
func (s *Socket) handleSetSerialClient(bodyBuffer []byte, conn net.Conn) {
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
		int(request.BufferTime),
		int(request.Timeout),
	)

	gomodbus.Logger.Sugar().Infof("created new Serial client to %s", request.Port)

	s.sendACK(conn)
}

// handleConnectRequest processes a request to connect to a slave.
func (s *Socket) handleConnect(conn net.Conn) {

	if s.client == nil {
		gomodbus.Logger.Sugar().Errorf("client not set")
		s.sendNACK(conn, errors.New("client not set"))
		return
	}

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
func (s *Socket) handleDisconnect(conn net.Conn) {

	if s.client == nil {
		gomodbus.Logger.Sugar().Errorf("client not set")
		s.sendNACK(conn, errors.New("client not set"))
		return
	}

	err := s.client.Disconnect()
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to disconnect: %v", err)
		s.sendNACK(conn, err)
		return
	}

	gomodbus.Logger.Sugar().Infof("ModBus client disconnected.")

	s.sendACK(conn)
}

// handleAddRelayRequest processes a request to add a new relay.
func (s *Socket) handleAddRelay(bodyBuffer []byte, conn net.Conn) {

	var request Relay
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddRelayRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	err = server.AddRelays(s.server, byte(request.Slave.UnitId), uint16(request.Address), request.Values, request.Writable)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to add relays: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)

}

// handleDeleteRelayRequest processes a request to delete a relay.
func (s *Socket) handleDeleteRelay(bodyBuffer []byte, conn net.Conn) {

	var request Relay
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal RemoveRelayRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	var addresses []uint16
	for i := uint32(0); i < request.Count; i++ {
		addresses = append(addresses, uint16(request.Address+i))
	}

	err = server.DeleteRelays(s.server, byte(request.Slave.UnitId), addresses, request.Writable)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to remove relays: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)

}

// handleSetRelayRequest processes a request to set a relay.
func (s *Socket) handleSetRelay(bodyBuffer []byte, conn net.Conn) {

	var request Relay
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal SetRelayRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	err = server.SetRelays(s.server, byte(request.Slave.UnitId), uint16(request.Address), request.Values, request.Writable)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to set relays: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)

}

// handleReadRelayRequest processes a request to read a relay.
func (s *Socket) handleReadRelay(bodyBuffer []byte, conn net.Conn) {

	// unmarshal the request
	var request Relay
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal ReadRelayRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	// check if the client is set
	if s.client == nil {
		gomodbus.Logger.Sugar().Errorf("client not set")
		s.sendNACK(conn, errors.New("client not set"))
		return
	}

	// read the relays
	var values = make([]bool, request.Count)
	if request.Writable {
		values, err = client.ReadCoils(s.client, byte(request.Slave.UnitId), uint16(request.Address), uint16(request.Count))
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to read coils: %v", err)
			s.sendNACK(conn, err)
			return
		}
	} else {
		values, err = client.ReadDiscreteInputs(s.client, byte(request.Slave.UnitId), uint16(request.Address), uint16(request.Count))
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to read discrete inputs: %v", err)
			s.sendNACK(conn, err)
			return
		}
	}

	// send the response
	response := &Relay{
		Slave:    request.Slave,
		Address:  request.Address,
		Count:    uint32(len(values)),
		Values:   values,
		Writable: request.Writable,
	}

	buffer, err := proto.Marshal(response)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to marshal Relay: %v", err)
		s.sendNACK(conn, err)
		return
	}

	header := &Header{
		Type:   RequestType_ACK,
		Length: uint64(len(buffer)),
	}

	sendResponse(conn, header, buffer)
}

// handleWriteRelayRequest processes a request to write a relay.
func (s *Socket) handleWriteRelay(bodyBuffer []byte, conn net.Conn) {

	var request Relay
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal WriteRelayRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	if s.client == nil {
		gomodbus.Logger.Sugar().Errorf("client not set")
		s.sendNACK(conn, errors.New("client not set"))
		return
	}

	if request.Writable {
		err = client.WriteMultipleCoils(s.client, byte(request.Slave.UnitId), uint16(request.Address), request.Values)
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to write coils: %v", err)
			s.sendNACK(conn, err)
			return
		}
		s.sendACK(conn)
		return
	} else {
		gomodbus.Logger.Sugar().Errorf("relay is not writable")
		s.sendNACK(conn, errors.New("relay is not writable"))
		return
	}

}

// handleAddRegisterRequest processes a request to add a new register.
func (s *Socket) handleAddRegister(bodyBuffer []byte, conn net.Conn) {

	var request Register
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal AddRegisterRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	err = server.AddRegisters(s.server, byte(request.Slave.UnitId), uint16(request.Address), request.Values, request.Writable)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to add registers: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)

}

// handleDeleteRegisterRequest processes a request to delete a register.
func (s *Socket) handleDeleteRegister(bodyBuffer []byte, conn net.Conn) {

	var request Register
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal DeleteRegisterRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	var addresses []uint16
	for i := uint32(0); i < request.Count; i++ {
		addresses = append(addresses, uint16(request.Address+i))
	}

	err = server.DeleteRegisters(s.server, byte(request.Slave.UnitId), addresses, request.Writable)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to delete registers: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)

}

// handleSetRegisterRequest processes a request to set a register.
func (s *Socket) handleSetRegister(bodyBuffer []byte, conn net.Conn) {

	var request Register
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal SetRegisterRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	if s.server == nil {
		gomodbus.Logger.Sugar().Errorf("server not set")
		s.sendNACK(conn, errors.New("server not set"))
		return
	}

	err = server.SetRegisters(s.server, byte(request.Slave.UnitId), uint16(request.Address), request.Values, request.Writable)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to set registers: %v", err)
		s.sendNACK(conn, err)
		return
	}

	s.sendACK(conn)

}

// handleReadRegisterRequest processes a request to read a register.
func (s *Socket) handleReadRegister(bodyBuffer []byte, conn net.Conn) {

	var request Register
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal ReadRegisterRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	if s.client == nil {
		gomodbus.Logger.Sugar().Errorf("client not set")
		s.sendNACK(conn, errors.New("client not set"))
		return
	}

	var values = make([][]byte, request.Count)
	if request.Writable {
		values, err = client.ReadHoldingRegisters(s.client, byte(request.Slave.UnitId), uint16(request.Address), uint16(request.Count))
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to read holding registers: %v", err)
			s.sendNACK(conn, err)
			return
		}
	} else {
		values, err = client.ReadInputRegisters(s.client, byte(request.Slave.UnitId), uint16(request.Address), uint16(request.Count))
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to read input registers: %v", err)
			s.sendNACK(conn, err)
			return
		}
	}

	// send the response
	response := &Register{
		Slave:    request.Slave,
		Address:  request.Address,
		Count:    uint32(len(values)),
		Values:   values,
		Writable: request.Writable,
	}

	buffer, err := proto.Marshal(response)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to marshal Register: %v", err)
		s.sendNACK(conn, err)
		return
	}

	header := &Header{
		Type:   RequestType_ACK,
		Length: uint64(len(buffer)),
	}

	sendResponse(conn, header, buffer)

}

// handleWriteRegisterRequest processes a request to write a register.
func (s *Socket) handleWriteRegister(bodyBuffer []byte, conn net.Conn) {

	var request Register
	err := proto.Unmarshal(bodyBuffer, &request)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to unmarshal WriteRegisterRequest: %v", err)
		s.sendNACK(conn, err)
		return
	}

	if s.client == nil {
		gomodbus.Logger.Sugar().Errorf("client not set")
		s.sendNACK(conn, errors.New("client not set"))
		return
	}

	var values []byte
	for _, value := range request.Values {
		values = append(values, value...)
	}

	if request.Writable {
		err = client.WriteMultipleRegisters(s.client, byte(request.Slave.UnitId), uint16(request.Address), uint16(request.Count), values)
		if err != nil {
			gomodbus.Logger.Sugar().Errorf("failed to write multiple registers: %v", err)
			s.sendNACK(conn, err)
			return
		}
		s.sendACK(conn)
		return
	} else {
		gomodbus.Logger.Sugar().Errorf("register is not writable")
		s.sendNACK(conn, errors.New("register is not writable"))
		return
	}

}
