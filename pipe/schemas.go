package pipe

const (
	NewTCPServerRequest = 1
	NewSerialServerRequest = 2

	NewSlaveRequest = 3
	RemoveSlaveRequest = 4

	AddCoilsRequest = 5
	DeleteCoilsRequest = 6
	AddDiscreteInputsRequest = 7
	DeleteDiscreteInputsRequest = 8
	AddHoldingRegistersRequest = 9
	DeleteHoldingRegistersRequest = 10
	AddInputRegistersRequest = 11
	DeleteInputRegistersRequest = 12

	
)

type Header struct {
	Type   uint8  `json:"type"`
	Length uint64 `json:"length"`
}

type TCPServerRequest struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	UseTLS    bool   `json:"use_tls"`
	ByteOrder string `json:"byte_order"`
	WordOrder string `json:"word_order"`
	CertFile  string `json:"cert_file"`
	KeyFile   string `json:"key_file"`
	CAFile    string `json:"ca_file"`
}

type SerialServerRequest struct {
	Port     string `json:"port"`
	BaudRate byte   `json:"baudrate"`
	DataBits byte   `json:"databits"`
	Parity   string `json:"parity"`
	StopBits  byte   `json:"stopbits"`
	ByteOrder string `json:"byte_order"`
	WordOrder string `json:"word_order"`
}

type SlaveRequest struct {
	UnitID byte `json:"unit_id"`
}

type CoilsRequest struct {
	SlaveRequest
	Address uint16 `json:"address"`
	Values  []bool `json:"values"`
}

type DiscreteInputsRequest struct {
	SlaveRequest
	Address uint16 `json:"address"`
	Values  []bool `json:"values"`
}

type HoldingRegistersRequest struct {
	SlaveRequest
	Address uint16 `json:"address"`
	Values  [][]byte `json:"values"`
}

type InputRegistersRequest struct {
	SlaveRequest
	Address uint16 `json:"address"`
	Values  [][]byte `json:"values"`
}
