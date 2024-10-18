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
	conn.SetReadDeadline(time.Now().Add(30 * time.Second)) // Set a read deadline
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

func TestSetTCPServer(t *testing.T) {
	// Create and start the socket
	sock := socket.NewSocket()
	go sock.Start()
	defer sock.Stop()

	// Give the server a moment to start
	time.Sleep(1 * time.Second)

	// Connect to the socket
	conn, err := net.Dial("unix", "./modbus.sock")
	if err != nil {
		t.Fatalf("Failed to connect to socket: %v", err)
	}
	defer conn.Close()

	request := &socket.TCPServerRequest{
		Host: "127.0.0.1",
		Port:     1502,
		UseTls:   false,
	}
	sendRequest(t, conn, socket.RequestType_SetTCPServer, request)
	checkResponse(t, conn)

	sendRequest(t, conn, socket.RequestType_StartServer, nil)
	checkResponse(t, conn)

	newRequest := &socket.Slave{
		UnitId: 1,
	}
	sendRequest(t, conn, socket.RequestType_AddSlave, newRequest)
	checkResponse(t, conn)

	relayRequest := &socket.Relay{
		Slave: &socket.Slave{
			UnitId: 1,
		},
		Address: 1,
		Count: 10,
		Values: []bool{true, false, true, false, true, false, true, false, true, false},
		Writable: true,
	}
	sendRequest(t, conn, socket.RequestType_AddRelay, relayRequest)
	checkResponse(t, conn)

	registerRequest := &socket.Register{
		Slave: &socket.Slave{
			UnitId: 1,
		},
		Address: 1,
		Count: 10,
		Values: [][]byte{[]byte{1, 2}, []byte{3, 4}, []byte{5, 6}, []byte{7, 8}, []byte{9, 10}, []byte{11, 12}, []byte{13, 14}, []byte{15, 16}, []byte{17, 18}, []byte{19, 20}},
		Writable: true,
	}
	sendRequest(t, conn, socket.RequestType_AddRegister, registerRequest)
	checkResponse(t, conn)

	clientRequest := &socket.TCPClientRequest{
		Host: "127.0.0.1",
		Port: 1502,
		UseTls: false,
	}
	sendRequest(t, conn, socket.RequestType_SetTCPClient, clientRequest)
	checkResponse(t, conn)

	sendRequest(t, conn, socket.RequestType_Connect, nil)
	checkResponse(t, conn)

	registerReadRequest := &socket.Register{
		Slave: &socket.Slave{
			UnitId: 1,
		},
		Address: 1,
		Count: 10,
		Writable: true,
	}
	sendRequest(t, conn, socket.RequestType_ReadRegister, registerReadRequest)
	checkResponse(t, conn)

	deleteRelaysRequest := &socket.Relay{
		Slave: &socket.Slave{
			UnitId: 1,
		},
		Address: 1,
		Count: 10,
		Values: []bool{true, false, true, false, true, false, true, false, true, false},
		Writable: true,
	}
	sendRequest(t, conn, socket.RequestType_DeleteRelay, deleteRelaysRequest)
	checkResponse(t, conn)

	sendRequest(t, conn, socket.RequestType_StopServer, nil)
	checkResponse(t, conn)
}
