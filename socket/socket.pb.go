// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.28.0
// source: socket.proto

package socket

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RequestType int32

const (
	RequestType_ACK                           RequestType = 0
	RequestType_NACK                          RequestType = 1
	RequestType_NewTCPServerRequest           RequestType = 2
	RequestType_NewSerialServerRequest        RequestType = 3
	RequestType_NewSlaveRequest               RequestType = 4
	RequestType_RemoveSlaveRequest            RequestType = 5
	RequestType_AddCoilsRequest               RequestType = 6
	RequestType_DeleteCoilsRequest            RequestType = 7
	RequestType_AddDiscreteInputsRequest      RequestType = 8
	RequestType_DeleteDiscreteInputsRequest   RequestType = 9
	RequestType_AddHoldingRegistersRequest    RequestType = 10
	RequestType_DeleteHoldingRegistersRequest RequestType = 11
	RequestType_AddInputRegistersRequest      RequestType = 12
	RequestType_DeleteInputRegistersRequest   RequestType = 13
)

// Enum value maps for RequestType.
var (
	RequestType_name = map[int32]string{
		0:  "ACK",
		1:  "NACK",
		2:  "NewTCPServerRequest",
		3:  "NewSerialServerRequest",
		4:  "NewSlaveRequest",
		5:  "RemoveSlaveRequest",
		6:  "AddCoilsRequest",
		7:  "DeleteCoilsRequest",
		8:  "AddDiscreteInputsRequest",
		9:  "DeleteDiscreteInputsRequest",
		10: "AddHoldingRegistersRequest",
		11: "DeleteHoldingRegistersRequest",
		12: "AddInputRegistersRequest",
		13: "DeleteInputRegistersRequest",
	}
	RequestType_value = map[string]int32{
		"ACK":                           0,
		"NACK":                          1,
		"NewTCPServerRequest":           2,
		"NewSerialServerRequest":        3,
		"NewSlaveRequest":               4,
		"RemoveSlaveRequest":            5,
		"AddCoilsRequest":               6,
		"DeleteCoilsRequest":            7,
		"AddDiscreteInputsRequest":      8,
		"DeleteDiscreteInputsRequest":   9,
		"AddHoldingRegistersRequest":    10,
		"DeleteHoldingRegistersRequest": 11,
		"AddInputRegistersRequest":      12,
		"DeleteInputRegistersRequest":   13,
	}
)

func (x RequestType) Enum() *RequestType {
	p := new(RequestType)
	*p = x
	return p
}

func (x RequestType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RequestType) Descriptor() protoreflect.EnumDescriptor {
	return file_socket_proto_enumTypes[0].Descriptor()
}

func (RequestType) Type() protoreflect.EnumType {
	return &file_socket_proto_enumTypes[0]
}

func (x RequestType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RequestType.Descriptor instead.
func (RequestType) EnumDescriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{0}
}

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   RequestType `protobuf:"varint,1,opt,name=type,proto3,enum=socket.RequestType" json:"type,omitempty"`
	Length uint64      `protobuf:"varint,2,opt,name=length,proto3" json:"length,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{0}
}

func (x *Header) GetType() RequestType {
	if x != nil {
		return x.Type
	}
	return RequestType_ACK
}

func (x *Header) GetLength() uint64 {
	if x != nil {
		return x.Length
	}
	return 0
}

type TCPServerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Host      string `protobuf:"bytes,1,opt,name=host,proto3" json:"host,omitempty"`
	Port      int32  `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	UseTls    bool   `protobuf:"varint,3,opt,name=use_tls,json=useTls,proto3" json:"use_tls,omitempty"`
	ByteOrder string `protobuf:"bytes,4,opt,name=byte_order,json=byteOrder,proto3" json:"byte_order,omitempty"`
	WordOrder string `protobuf:"bytes,5,opt,name=word_order,json=wordOrder,proto3" json:"word_order,omitempty"`
	CertFile  string `protobuf:"bytes,6,opt,name=cert_file,json=certFile,proto3" json:"cert_file,omitempty"`
	KeyFile   string `protobuf:"bytes,7,opt,name=key_file,json=keyFile,proto3" json:"key_file,omitempty"`
	CaFile    string `protobuf:"bytes,8,opt,name=ca_file,json=caFile,proto3" json:"ca_file,omitempty"`
}

func (x *TCPServerRequest) Reset() {
	*x = TCPServerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TCPServerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TCPServerRequest) ProtoMessage() {}

func (x *TCPServerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TCPServerRequest.ProtoReflect.Descriptor instead.
func (*TCPServerRequest) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{1}
}

func (x *TCPServerRequest) GetHost() string {
	if x != nil {
		return x.Host
	}
	return ""
}

func (x *TCPServerRequest) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *TCPServerRequest) GetUseTls() bool {
	if x != nil {
		return x.UseTls
	}
	return false
}

func (x *TCPServerRequest) GetByteOrder() string {
	if x != nil {
		return x.ByteOrder
	}
	return ""
}

func (x *TCPServerRequest) GetWordOrder() string {
	if x != nil {
		return x.WordOrder
	}
	return ""
}

func (x *TCPServerRequest) GetCertFile() string {
	if x != nil {
		return x.CertFile
	}
	return ""
}

func (x *TCPServerRequest) GetKeyFile() string {
	if x != nil {
		return x.KeyFile
	}
	return ""
}

func (x *TCPServerRequest) GetCaFile() string {
	if x != nil {
		return x.CaFile
	}
	return ""
}

type SerialServerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Port      string `protobuf:"bytes,1,opt,name=port,proto3" json:"port,omitempty"`
	Baudrate  int32  `protobuf:"varint,2,opt,name=baudrate,proto3" json:"baudrate,omitempty"`
	Databits  uint32 `protobuf:"varint,3,opt,name=databits,proto3" json:"databits,omitempty"`
	Parity    uint32 `protobuf:"varint,4,opt,name=parity,proto3" json:"parity,omitempty"`
	Stopbits  uint32 `protobuf:"varint,5,opt,name=stopbits,proto3" json:"stopbits,omitempty"`
	ByteOrder string `protobuf:"bytes,6,opt,name=byte_order,json=byteOrder,proto3" json:"byte_order,omitempty"`
	WordOrder string `protobuf:"bytes,7,opt,name=word_order,json=wordOrder,proto3" json:"word_order,omitempty"`
}

func (x *SerialServerRequest) Reset() {
	*x = SerialServerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SerialServerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SerialServerRequest) ProtoMessage() {}

func (x *SerialServerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SerialServerRequest.ProtoReflect.Descriptor instead.
func (*SerialServerRequest) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{2}
}

func (x *SerialServerRequest) GetPort() string {
	if x != nil {
		return x.Port
	}
	return ""
}

func (x *SerialServerRequest) GetBaudrate() int32 {
	if x != nil {
		return x.Baudrate
	}
	return 0
}

func (x *SerialServerRequest) GetDatabits() uint32 {
	if x != nil {
		return x.Databits
	}
	return 0
}

func (x *SerialServerRequest) GetParity() uint32 {
	if x != nil {
		return x.Parity
	}
	return 0
}

func (x *SerialServerRequest) GetStopbits() uint32 {
	if x != nil {
		return x.Stopbits
	}
	return 0
}

func (x *SerialServerRequest) GetByteOrder() string {
	if x != nil {
		return x.ByteOrder
	}
	return ""
}

func (x *SerialServerRequest) GetWordOrder() string {
	if x != nil {
		return x.WordOrder
	}
	return ""
}

type SlaveRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UnitId uint32 `protobuf:"varint,1,opt,name=unit_id,json=unitId,proto3" json:"unit_id,omitempty"`
}

func (x *SlaveRequest) Reset() {
	*x = SlaveRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SlaveRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SlaveRequest) ProtoMessage() {}

func (x *SlaveRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SlaveRequest.ProtoReflect.Descriptor instead.
func (*SlaveRequest) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{3}
}

func (x *SlaveRequest) GetUnitId() uint32 {
	if x != nil {
		return x.UnitId
	}
	return 0
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlaveRequest *SlaveRequest `protobuf:"bytes,1,opt,name=slave_request,json=slaveRequest,proto3" json:"slave_request,omitempty"`
	Addresses    []uint32      `protobuf:"varint,2,rep,packed,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{4}
}

func (x *DeleteRequest) GetSlaveRequest() *SlaveRequest {
	if x != nil {
		return x.SlaveRequest
	}
	return nil
}

func (x *DeleteRequest) GetAddresses() []uint32 {
	if x != nil {
		return x.Addresses
	}
	return nil
}

type CoilsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlaveRequest *SlaveRequest `protobuf:"bytes,1,opt,name=slave_request,json=slaveRequest,proto3" json:"slave_request,omitempty"`
	Address      uint32        `protobuf:"varint,2,opt,name=address,proto3" json:"address,omitempty"`
	Values       []bool        `protobuf:"varint,3,rep,packed,name=values,proto3" json:"values,omitempty"`
}

func (x *CoilsRequest) Reset() {
	*x = CoilsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CoilsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CoilsRequest) ProtoMessage() {}

func (x *CoilsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CoilsRequest.ProtoReflect.Descriptor instead.
func (*CoilsRequest) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{5}
}

func (x *CoilsRequest) GetSlaveRequest() *SlaveRequest {
	if x != nil {
		return x.SlaveRequest
	}
	return nil
}

func (x *CoilsRequest) GetAddress() uint32 {
	if x != nil {
		return x.Address
	}
	return 0
}

func (x *CoilsRequest) GetValues() []bool {
	if x != nil {
		return x.Values
	}
	return nil
}

type DiscreteInputsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlaveRequest *SlaveRequest `protobuf:"bytes,1,opt,name=slave_request,json=slaveRequest,proto3" json:"slave_request,omitempty"`
	Address      uint32        `protobuf:"varint,2,opt,name=address,proto3" json:"address,omitempty"`
	Values       []bool        `protobuf:"varint,3,rep,packed,name=values,proto3" json:"values,omitempty"`
}

func (x *DiscreteInputsRequest) Reset() {
	*x = DiscreteInputsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscreteInputsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscreteInputsRequest) ProtoMessage() {}

func (x *DiscreteInputsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscreteInputsRequest.ProtoReflect.Descriptor instead.
func (*DiscreteInputsRequest) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{6}
}

func (x *DiscreteInputsRequest) GetSlaveRequest() *SlaveRequest {
	if x != nil {
		return x.SlaveRequest
	}
	return nil
}

func (x *DiscreteInputsRequest) GetAddress() uint32 {
	if x != nil {
		return x.Address
	}
	return 0
}

func (x *DiscreteInputsRequest) GetValues() []bool {
	if x != nil {
		return x.Values
	}
	return nil
}

type DiscreteInputsRequestDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlaveRequest *SlaveRequest `protobuf:"bytes,1,opt,name=slave_request,json=slaveRequest,proto3" json:"slave_request,omitempty"`
	Addresses    []uint32      `protobuf:"varint,2,rep,packed,name=addresses,proto3" json:"addresses,omitempty"`
}

func (x *DiscreteInputsRequestDelete) Reset() {
	*x = DiscreteInputsRequestDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DiscreteInputsRequestDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DiscreteInputsRequestDelete) ProtoMessage() {}

func (x *DiscreteInputsRequestDelete) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DiscreteInputsRequestDelete.ProtoReflect.Descriptor instead.
func (*DiscreteInputsRequestDelete) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{7}
}

func (x *DiscreteInputsRequestDelete) GetSlaveRequest() *SlaveRequest {
	if x != nil {
		return x.SlaveRequest
	}
	return nil
}

func (x *DiscreteInputsRequestDelete) GetAddresses() []uint32 {
	if x != nil {
		return x.Addresses
	}
	return nil
}

type HoldingRegistersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlaveRequest *SlaveRequest `protobuf:"bytes,1,opt,name=slave_request,json=slaveRequest,proto3" json:"slave_request,omitempty"`
	Address      uint32        `protobuf:"varint,2,opt,name=address,proto3" json:"address,omitempty"`
	Values       [][]byte      `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *HoldingRegistersRequest) Reset() {
	*x = HoldingRegistersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HoldingRegistersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HoldingRegistersRequest) ProtoMessage() {}

func (x *HoldingRegistersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HoldingRegistersRequest.ProtoReflect.Descriptor instead.
func (*HoldingRegistersRequest) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{8}
}

func (x *HoldingRegistersRequest) GetSlaveRequest() *SlaveRequest {
	if x != nil {
		return x.SlaveRequest
	}
	return nil
}

func (x *HoldingRegistersRequest) GetAddress() uint32 {
	if x != nil {
		return x.Address
	}
	return 0
}

func (x *HoldingRegistersRequest) GetValues() [][]byte {
	if x != nil {
		return x.Values
	}
	return nil
}

type InputRegistersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SlaveRequest *SlaveRequest `protobuf:"bytes,1,opt,name=slave_request,json=slaveRequest,proto3" json:"slave_request,omitempty"`
	Address      uint32        `protobuf:"varint,2,opt,name=address,proto3" json:"address,omitempty"`
	Values       [][]byte      `protobuf:"bytes,3,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *InputRegistersRequest) Reset() {
	*x = InputRegistersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_socket_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputRegistersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputRegistersRequest) ProtoMessage() {}

func (x *InputRegistersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_socket_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputRegistersRequest.ProtoReflect.Descriptor instead.
func (*InputRegistersRequest) Descriptor() ([]byte, []int) {
	return file_socket_proto_rawDescGZIP(), []int{9}
}

func (x *InputRegistersRequest) GetSlaveRequest() *SlaveRequest {
	if x != nil {
		return x.SlaveRequest
	}
	return nil
}

func (x *InputRegistersRequest) GetAddress() uint32 {
	if x != nil {
		return x.Address
	}
	return 0
}

func (x *InputRegistersRequest) GetValues() [][]byte {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_socket_proto protoreflect.FileDescriptor

var file_socket_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x22, 0x49, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13,
	0x2e, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x6e,
	0x67, 0x74, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x6c, 0x65, 0x6e, 0x67, 0x74,
	0x68, 0x22, 0xe2, 0x01, 0x0a, 0x10, 0x54, 0x43, 0x50, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x6f, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x75, 0x73, 0x65, 0x5f, 0x74, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x06, 0x75, 0x73, 0x65, 0x54, 0x6c, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x79, 0x74, 0x65, 0x5f,
	0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x62, 0x79, 0x74,
	0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x0a, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x77, 0x6f, 0x72, 0x64,
	0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x63, 0x65, 0x72, 0x74, 0x5f, 0x66, 0x69,
	0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x65, 0x72, 0x74, 0x46, 0x69,
	0x6c, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x6b, 0x65, 0x79, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x17, 0x0a,
	0x07, 0x63, 0x61, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x63, 0x61, 0x46, 0x69, 0x6c, 0x65, 0x22, 0xd3, 0x01, 0x0a, 0x13, 0x53, 0x65, 0x72, 0x69, 0x61,
	0x6c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x61, 0x75, 0x64, 0x72, 0x61, 0x74, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x62, 0x61, 0x75, 0x64, 0x72, 0x61, 0x74, 0x65, 0x12, 0x1a,
	0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x69, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x62, 0x69, 0x74, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x61,
	0x72, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x70, 0x61, 0x72, 0x69,
	0x74, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x6f, 0x70, 0x62, 0x69, 0x74, 0x73, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x73, 0x74, 0x6f, 0x70, 0x62, 0x69, 0x74, 0x73, 0x12, 0x1d,
	0x0a, 0x0a, 0x62, 0x79, 0x74, 0x65, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x62, 0x79, 0x74, 0x65, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x1d, 0x0a,
	0x0a, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x77, 0x6f, 0x72, 0x64, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x22, 0x27, 0x0a, 0x0c,
	0x53, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x6e, 0x69, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x75,
	0x6e, 0x69, 0x74, 0x49, 0x64, 0x22, 0x68, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x0d, 0x73, 0x6c, 0x61, 0x76, 0x65, 0x5f,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x53, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x0c, 0x73, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x0d, 0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x22,
	0x7b, 0x0a, 0x0c, 0x43, 0x6f, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x39, 0x0a, 0x0d, 0x73, 0x6c, 0x61, 0x76, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e,
	0x53, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x0c, 0x73, 0x6c,
	0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x08, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x22, 0x84, 0x01, 0x0a,
	0x15, 0x44, 0x69, 0x73, 0x63, 0x72, 0x65, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x73, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x0d, 0x73, 0x6c, 0x61, 0x76, 0x65, 0x5f,
	0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x53, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x52, 0x0c, 0x73, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x08, 0x52, 0x06, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x73, 0x22, 0x76, 0x0a, 0x1b, 0x44, 0x69, 0x73, 0x63, 0x72, 0x65, 0x74, 0x65, 0x49,
	0x6e, 0x70, 0x75, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x39, 0x0a, 0x0d, 0x73, 0x6c, 0x61, 0x76, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x6f, 0x63, 0x6b,
	0x65, 0x74, 0x2e, 0x53, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52,
	0x0c, 0x73, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a,
	0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0d,
	0x52, 0x09, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x65, 0x73, 0x22, 0x86, 0x01, 0x0a, 0x17,
	0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x0d, 0x73, 0x6c, 0x61, 0x76, 0x65,
	0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14,
	0x2e, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x53, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x52, 0x0c, 0x73, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0c, 0x52, 0x06, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x22, 0x84, 0x01, 0x0a, 0x15, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x39,
	0x0a, 0x0d, 0x73, 0x6c, 0x61, 0x76, 0x65, 0x5f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x53,
	0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x0c, 0x73, 0x6c, 0x61,
	0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0c, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x2a, 0xf0, 0x02, 0x0a, 0x0b,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x41,
	0x43, 0x4b, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x41, 0x43, 0x4b, 0x10, 0x01, 0x12, 0x17,
	0x0a, 0x13, 0x4e, 0x65, 0x77, 0x54, 0x43, 0x50, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x4e, 0x65, 0x77, 0x53, 0x65,
	0x72, 0x69, 0x61, 0x6c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x10, 0x03, 0x12, 0x13, 0x0a, 0x0f, 0x4e, 0x65, 0x77, 0x53, 0x6c, 0x61, 0x76, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x04, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x53, 0x6c, 0x61, 0x76, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x05,
	0x12, 0x13, 0x0a, 0x0f, 0x41, 0x64, 0x64, 0x43, 0x6f, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x10, 0x06, 0x12, 0x16, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43,
	0x6f, 0x69, 0x6c, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x07, 0x12, 0x1c, 0x0a,
	0x18, 0x41, 0x64, 0x64, 0x44, 0x69, 0x73, 0x63, 0x72, 0x65, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x08, 0x12, 0x1f, 0x0a, 0x1b, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x69, 0x73, 0x63, 0x72, 0x65, 0x74, 0x65, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x09, 0x12, 0x1e, 0x0a, 0x1a,
	0x41, 0x64, 0x64, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x0a, 0x12, 0x21, 0x0a, 0x1d,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x48, 0x6f, 0x6c, 0x64, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x67,
	0x69, 0x73, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x0b, 0x12,
	0x1c, 0x0a, 0x18, 0x41, 0x64, 0x64, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65, 0x67, 0x69, 0x73,
	0x74, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x0c, 0x12, 0x1f, 0x0a,
	0x1b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x10, 0x0d, 0x42, 0x0a,
	0x5a, 0x08, 0x2e, 0x2f, 0x73, 0x6f, 0x63, 0x6b, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_socket_proto_rawDescOnce sync.Once
	file_socket_proto_rawDescData = file_socket_proto_rawDesc
)

func file_socket_proto_rawDescGZIP() []byte {
	file_socket_proto_rawDescOnce.Do(func() {
		file_socket_proto_rawDescData = protoimpl.X.CompressGZIP(file_socket_proto_rawDescData)
	})
	return file_socket_proto_rawDescData
}

var file_socket_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_socket_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_socket_proto_goTypes = []any{
	(RequestType)(0),                    // 0: socket.RequestType
	(*Header)(nil),                      // 1: socket.Header
	(*TCPServerRequest)(nil),            // 2: socket.TCPServerRequest
	(*SerialServerRequest)(nil),         // 3: socket.SerialServerRequest
	(*SlaveRequest)(nil),                // 4: socket.SlaveRequest
	(*DeleteRequest)(nil),               // 5: socket.DeleteRequest
	(*CoilsRequest)(nil),                // 6: socket.CoilsRequest
	(*DiscreteInputsRequest)(nil),       // 7: socket.DiscreteInputsRequest
	(*DiscreteInputsRequestDelete)(nil), // 8: socket.DiscreteInputsRequestDelete
	(*HoldingRegistersRequest)(nil),     // 9: socket.HoldingRegistersRequest
	(*InputRegistersRequest)(nil),       // 10: socket.InputRegistersRequest
}
var file_socket_proto_depIdxs = []int32{
	0, // 0: socket.Header.type:type_name -> socket.RequestType
	4, // 1: socket.DeleteRequest.slave_request:type_name -> socket.SlaveRequest
	4, // 2: socket.CoilsRequest.slave_request:type_name -> socket.SlaveRequest
	4, // 3: socket.DiscreteInputsRequest.slave_request:type_name -> socket.SlaveRequest
	4, // 4: socket.DiscreteInputsRequestDelete.slave_request:type_name -> socket.SlaveRequest
	4, // 5: socket.HoldingRegistersRequest.slave_request:type_name -> socket.SlaveRequest
	4, // 6: socket.InputRegistersRequest.slave_request:type_name -> socket.SlaveRequest
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_socket_proto_init() }
func file_socket_proto_init() {
	if File_socket_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_socket_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Header); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_socket_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*TCPServerRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_socket_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*SerialServerRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_socket_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*SlaveRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_socket_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*DeleteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_socket_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*CoilsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_socket_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*DiscreteInputsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_socket_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*DiscreteInputsRequestDelete); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_socket_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*HoldingRegistersRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_socket_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*InputRegistersRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_socket_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_socket_proto_goTypes,
		DependencyIndexes: file_socket_proto_depIdxs,
		EnumInfos:         file_socket_proto_enumTypes,
		MessageInfos:      file_socket_proto_msgTypes,
	}.Build()
	File_socket_proto = out.File
	file_socket_proto_rawDesc = nil
	file_socket_proto_goTypes = nil
	file_socket_proto_depIdxs = nil
}
