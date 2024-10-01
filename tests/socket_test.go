package tests

import (
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"

	"github.com/ulfaric/gomodbus/socket"
	"google.golang.org/protobuf/proto"
)

func TestSocket(t *testing.T) {
	// Create and start the socket
	sock := socket.NewSocket()
	go sock.Start()
	defer sock.Stop()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)

	// Connect to the socket
	conn, err := net.Dial("tcp", "127.0.0.1:1993")
	if err != nil {
		t.Fatalf("Failed to connect to socket: %v", err)
	}
	defer conn.Close()

	// Test NewTCPServerRequest
	testNewTCPServerRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test NewSlaveRequest
	testNewSlaveRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test AddCoilsRequest
	testAddCoilsRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test AddDiscreteInputsRequest
	testAddDiscreteInputsRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test AddHoldingRegistersRequest
	testAddHoldingRegistersRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test AddInputRegistersRequest
	testAddInputRegistersRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test StartServerRequest
	testStartServerRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test NewTCPClientRequest
	testNewTCPClientRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test ConnectClientRequest
	testConnectClientRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test ReadCoilsRequest
	testReadCoilsRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test ReadDiscreteInputsRequest
	testReadDiscreteInputsRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test ReadHoldingRegistersRequest
	testReadHoldingRegistersRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test ReadInputRegistersRequest
	testReadInputRegistersRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test WriteCoilRequest
	testWriteCoilRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test WriteCoilsRequest
	testWriteCoilsRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test WriteRegisterRequest
	testWriteRegisterRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test WriteRegistersRequest
	testWriteRegistersRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test DeleteCoilsRequest
	testDeleteCoilsRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test DeleteDiscreteInputsRequest
	testDeleteDiscreteInputsRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test DeleteHoldingRegistersRequest
	testDeleteHoldingRegistersRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test DeleteInputRegistersRequest
	testDeleteInputRegistersRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test RemoveSlaveRequest
	testRemoveSlaveRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test DisconnectClientRequest
	testDisconnectClientRequest(t, conn)
	time.Sleep(1 * time.Second)

	// Test StopServerRequest
	testStopServerRequest(t, conn)
	time.Sleep(2 * time.Second)
}


func testNewTCPServerRequest(t *testing.T, conn net.Conn) {
	request := &socket.TCPServerRequest{
		Host:      "127.0.0.1",
		Port:      1502,
		UseTls:    false,
		ByteOrder: "big",
		WordOrder: "big",
	}
	sendRequest(t, conn, socket.RequestType_NewTCPServerRequest, request)
	checkResponse(t, conn)
}

func testNewSerialServerRequest(t *testing.T, conn net.Conn) {
	request := &socket.SerialServerRequest{
		Port:      "/dev/ttyS0",
		Baudrate:  9600,
		Databits:  8,
		Parity:    0,
		Stopbits:  1,
		ByteOrder: "big",
		WordOrder: "big",
	}
	sendRequest(t, conn, socket.RequestType_NewSerialServerRequest, request)
	checkResponse(t, conn)
}

func testNewSlaveRequest(t *testing.T, conn net.Conn) {
	request := &socket.SlaveRequest{
		UnitId: 1,
	}
	sendRequest(t, conn, socket.RequestType_NewSlaveRequest, request)
	checkResponse(t, conn)
}

func testRemoveSlaveRequest(t *testing.T, conn net.Conn) {
	request := &socket.SlaveRequest{
		UnitId: 1,
	}
	sendRequest(t, conn, socket.RequestType_RemoveSlaveRequest, request)
	checkResponse(t, conn)
}

func testAddCoilsRequest(t *testing.T, conn net.Conn) {
	request := &socket.CoilsRequest{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Values:       []bool{true, false, true},
	}
	sendRequest(t, conn, socket.RequestType_AddCoilsRequest, request)
	checkResponse(t, conn)
}

func testDeleteCoilsRequest(t *testing.T, conn net.Conn) {
	request := &socket.DeleteRequest{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Addresses:    []uint32{0, 1, 2},
	}
	sendRequest(t, conn, socket.RequestType_DeleteCoilsRequest, request)
	checkResponse(t, conn)
}

func testAddDiscreteInputsRequest(t *testing.T, conn net.Conn) {
	request := &socket.DiscreteInputsRequest{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Values:       []bool{true, false, true},
	}
	sendRequest(t, conn, socket.RequestType_AddDiscreteInputsRequest, request)
	checkResponse(t, conn)
}

func testDeleteDiscreteInputsRequest(t *testing.T, conn net.Conn) {
	request := &socket.DeleteRequest{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Addresses:    []uint32{0, 1, 2},
	}
	sendRequest(t, conn, socket.RequestType_DeleteDiscreteInputsRequest, request)
	checkResponse(t, conn)
}

func testAddHoldingRegistersRequest(t *testing.T, conn net.Conn) {
	request := &socket.HoldingRegistersRequest{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Values:       [][]byte{{0x12, 0x34}, {0x56, 0x78}},
	}
	sendRequest(t, conn, socket.RequestType_AddHoldingRegistersRequest, request)
	checkResponse(t, conn)
}

func testDeleteHoldingRegistersRequest(t *testing.T, conn net.Conn) {
	request := &socket.DeleteRequest{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Addresses:    []uint32{0, 1},
	}
	sendRequest(t, conn, socket.RequestType_DeleteHoldingRegistersRequest, request)
	checkResponse(t, conn)
}

func testAddInputRegistersRequest(t *testing.T, conn net.Conn) {
	request := &socket.InputRegistersRequest{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Values:       [][]byte{{0x12, 0x34}, {0x56, 0x78}},
	}
	sendRequest(t, conn, socket.RequestType_AddInputRegistersRequest, request)
	checkResponse(t, conn)
}

func testDeleteInputRegistersRequest(t *testing.T, conn net.Conn) {
	request := &socket.DeleteRequest{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Addresses:    []uint32{0, 1},
	}
	sendRequest(t, conn, socket.RequestType_DeleteInputRegistersRequest, request)
	checkResponse(t, conn)
}

func testStartServerRequest(t *testing.T, conn net.Conn) {
	sendRequest(t, conn, socket.RequestType_StartServerRequest, nil)
	checkResponse(t, conn)
}

func testStopServerRequest(t *testing.T, conn net.Conn) {
	sendRequest(t, conn, socket.RequestType_StopServerRequest, nil)
	checkResponse(t, conn)
}

func testNewTCPClientRequest(t *testing.T, conn net.Conn) {
	request := &socket.TCPClientRequest{
		Host:      "127.0.0.1",
		Port:      1502,
		UseTls:    false,
		ByteOrder: "big",
		WordOrder: "big",
	}
	sendRequest(t, conn, socket.RequestType_NewTCPClientRequest, request)
	checkResponse(t, conn)
}

func testNewSerialClientRequest(t *testing.T, conn net.Conn) {
	request := &socket.SerialClientRequest{
		Port:      "/dev/ttyS0",
		Baudrate:  9600,
		Databits:  8,
		Parity:    0,
		Stopbits:  1,
		ByteOrder: "big",
		WordOrder: "big",
	}
	sendRequest(t, conn, socket.RequestType_NewSerialClientRequest, request)
	checkResponse(t, conn)
}

func testReadCoilsRequest(t *testing.T, conn net.Conn) {
	request := &socket.ReadCoils{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Count:        3,
	}
	sendRequest(t, conn, socket.RequestType_ReadCoilsRequest, request)
	checkResponse(t, conn)
}

func testReadDiscreteInputsRequest(t *testing.T, conn net.Conn) {
	request := &socket.ReadDiscreteInputs{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Count:        3,
	}
	sendRequest(t, conn, socket.RequestType_ReadDiscreteInputsRequest, request)
	checkResponse(t, conn)
}

func testReadHoldingRegistersRequest(t *testing.T, conn net.Conn) {
	request := &socket.ReadHoldingRegisters{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Count:        2,
	}
	sendRequest(t, conn, socket.RequestType_ReadHoldingRegistersRequest, request)
	checkResponse(t, conn)
}

func testReadInputRegistersRequest(t *testing.T, conn net.Conn) {
	request := &socket.ReadInputRegisters{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Count:        2,
	}
	sendRequest(t, conn, socket.RequestType_ReadInputRegistersRequest, request)
	checkResponse(t, conn)
}

func testWriteCoilRequest(t *testing.T, conn net.Conn) {
	request := &socket.WriteCoil{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Value:        true,
	}
	sendRequest(t, conn, socket.RequestType_WriteCoilRequest, request)
	checkResponse(t, conn)
}

func testWriteCoilsRequest(t *testing.T, conn net.Conn) {
	request := &socket.WriteCoils{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Values:       []bool{true, false, true},
	}
	sendRequest(t, conn, socket.RequestType_WriteCoilsRequest, request)
	checkResponse(t, conn)
}

func testWriteRegisterRequest(t *testing.T, conn net.Conn) {
	request := &socket.WriteRegister{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Value:        []byte{0x12, 0x34},
	}
	sendRequest(t, conn, socket.RequestType_WriteRegisterRequest, request)
	checkResponse(t, conn)
}

func testWriteRegistersRequest(t *testing.T, conn net.Conn) {
	request := &socket.WriteRegisters{
		SlaveRequest: &socket.SlaveRequest{UnitId: 1},
		Address:      0,
		Values:       []byte{0x12, 0x34, 0x56, 0x78},
	}
	sendRequest(t, conn, socket.RequestType_WriteRegistersRequest, request)
	checkResponse(t, conn)
}

func checkResponse(t *testing.T, conn net.Conn) {
	header, _, err := receiveResponse(t, conn)
	if err != nil {
		t.Fatalf("Failed to receive response: %v", err)
	}
	if header.Type != socket.RequestType_ACK {
		t.Fatalf("Expected ACK, got %v", header.Type)
	}
}

func sendRequest(t *testing.T, conn net.Conn, requestType socket.RequestType, request proto.Message) {
	header := &socket.Header{
		Type:   requestType,
		Length: 0,
	}
	headerBuffer, err := proto.Marshal(header)
	if err != nil {
		t.Fatalf("Failed to marshal header: %v", err)
	}

	var bodyBuffer []byte
	if request != nil {
		bodyBuffer, err = proto.Marshal(request)
		if err != nil {
			t.Fatalf("Failed to marshal request: %v", err)
		}
		header.Length = uint64(len(bodyBuffer))
	}

	buffer := make([]byte, 16)
	binary.BigEndian.PutUint64(buffer[0:8], uint64(len(headerBuffer)))
	binary.BigEndian.PutUint64(buffer[8:16], header.Length)

	buffer = append(buffer, headerBuffer...)
	if bodyBuffer != nil {
		buffer = append(buffer, bodyBuffer...)
	}

	_, err = conn.Write(buffer)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
}

func receiveResponse(t *testing.T, conn net.Conn) (*socket.Header, []byte, error) {
	conn.SetReadDeadline(time.Now().Add(1 * time.Second)) // Set a read deadline
	buffer := make([]byte, 16)
	_, err := conn.Read(buffer)
	if err != nil {
		return nil, nil, err
	}

	headerLength := binary.BigEndian.Uint64(buffer[0:8])
	conn.SetReadDeadline(time.Now().Add(1 * time.Second)) // Set a read deadline
	headerBuffer := make([]byte, headerLength)
	_, err = io.ReadFull(conn, headerBuffer)
	if err != nil {
		return nil, nil, err
	}

	header := &socket.Header{}
	err = proto.Unmarshal(headerBuffer, header)
	if err != nil {
		return nil, nil, err
	}

	bodyLength := binary.BigEndian.Uint64(buffer[8:16])
	bodyBuffer := make([]byte, bodyLength)
	_, err = io.ReadFull(conn, bodyBuffer)
	if err != nil {
		return nil, nil, err
	}

	return header, bodyBuffer, nil
}

func testConnectClientRequest(t *testing.T, conn net.Conn) {
	sendRequest(t, conn, socket.RequestType_ConnectRequest, nil)
	checkResponse(t, conn)
}

func testDisconnectClientRequest(t *testing.T, conn net.Conn) {
	sendRequest(t, conn, socket.RequestType_DisconnectRequest, nil)
	checkResponse(t, conn)
}

