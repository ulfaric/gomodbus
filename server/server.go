package server

import (
	"bytes"
	"errors"
	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/pdu"
)

type Server interface {
	Start() error
	Stop() error
	AddSlave(unitID byte)
	GetSlave(unitID byte) (*Slave, error)
	RemoveSlave(unitID byte)
}

// AddCoils adds new coils to the slave.
func AddCoils(server Server, unitID byte, address uint16, values []bool) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.AddCoils(address, values)
	gomodbus.Logger.Sugar().Infof("added %d coils to slave with unit ID %d", len(values), unitID)
	return nil
}

// SetCoils sets existing coils' values to the slave.
func SetCoils(server Server, unitID byte, address uint16, values []bool) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.SetCoils(address, values)
	gomodbus.Logger.Sugar().Infof("set %d coils on slave with unit ID %d", len(values), unitID)
	return nil
}

// DeleteCoils deletes coils from the slave.
func DeleteCoils(server Server, unitID byte, addresses []uint16) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.DeleteCoils(addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d coils from slave with unit ID %d", len(addresses), unitID)
	return nil
}

// AddDiscreteInputs adds new discrete inputs to the slave.
func AddDiscreteInputs(server Server, unitID byte, address uint16, values []bool) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.AddDiscreteInputs(address, values)
	gomodbus.Logger.Sugar().Infof("added %d discrete inputs to slave with unit ID %d", len(values), unitID)
	return nil
}

// SetDiscreteInputs sets existing discrete inputs' values to the slave.
func SetDiscreteInputs(server Server, unitID byte, address uint16, values []bool) error {

	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.SetDiscreteInputs(address, values)
	gomodbus.Logger.Sugar().Infof("set %d discrete inputs on slave with unit ID %d", len(values), unitID)
	return nil
}

// DeleteDiscreteInputs deletes discrete inputs from the slave.
func DeleteDiscreteInputs(server Server, unitID byte, addresses []uint16) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.DeleteDiscreteInputs(addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d discrete inputs from slave with unit ID %d", len(addresses), unitID)
	return nil
}

// AddRelays adds new relays to the slave. If writeable is true, the relays are coils, otherwise they are discrete inputs.
func AddRelays(server Server, unitID byte, address uint16, values []bool, writeable bool) error {
	if writeable {
		return AddCoils(server, unitID, address, values)
	}
	return AddDiscreteInputs(server, unitID, address, values)
}

// SetRelays sets existing relays' values to the slave. If writeable is true, the relays are coils, otherwise they are discrete inputs.
func SetRelays(server Server, unitID byte, address uint16, values []bool, writeable bool) error {
	if writeable {
		return SetCoils(server, unitID, address, values)
	}
	return SetDiscreteInputs(server, unitID, address, values)
}

// DeleteRelays deletes relays from the slave. If writeable is true, the relays are coils, otherwise they are discrete inputs.
func DeleteRelays(server Server, unitID byte, addresses []uint16, writeable bool) error {
	if writeable {
		return DeleteCoils(server, unitID, addresses)
	}
	return DeleteDiscreteInputs(server, unitID, addresses)
}

// AddHoldingRegisters adds new holding registers to the slave.
func AddHoldingRegisters(server Server, unitID byte, address uint16, values [][]byte) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.AddHoldingRegisters(address, values)
	gomodbus.Logger.Sugar().Infof("added %d holding registers to slave with unit ID %d", len(values), unitID)
	return nil
}

// SetHoldingRegisters sets existing holding registers' values to the slave.
func SetHoldingRegisters(server Server, unitID byte, address uint16, values [][]byte) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.SetHoldingRegisters(address, values)
	gomodbus.Logger.Sugar().Infof("set %d holding registers on slave with unit ID %d", len(values), unitID)
	return nil
}

// DeleteHoldingRegisters deletes holding registers from the slave.
func DeleteHoldingRegisters(server Server, unitID byte, addresses []uint16) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.DeleteHoldingRegisters(addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d holding registers from slave with unit ID %d", len(addresses), unitID)
	return nil
}

// AddInputRegisters adds new input registers to the slave.
func AddInputRegisters(server Server, unitID byte, address uint16, values [][]byte) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.AddInputRegisters(address, values)
	gomodbus.Logger.Sugar().Infof("added %d input registers to slave with unit ID %d", len(values), unitID)
	return nil
}

// SetInputRegisters sets existing input registers' values to the slave.
func SetInputRegisters(server Server, unitID byte, address uint16, values [][]byte) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.SetInputRegisters(address, values)
	gomodbus.Logger.Sugar().Infof("set %d input registers on slave with unit ID %d", len(values), unitID)
	return nil
}

// DeleteInputRegisters deletes input registers from the slave.
func DeleteInputRegisters(server Server, unitID byte, addresses []uint16) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.DeleteInputRegisters(addresses)
	gomodbus.Logger.Sugar().Infof("deleted %d input registers from slave with unit ID %d", len(addresses), unitID)
	return nil
}

// AddRegisters adds new registers to the slave. If writeable is true, the registers are holding registers, otherwise they are input registers.
func AddRegisters(server Server, unitID byte, address uint16, values [][]byte, writeable bool) error {
	if writeable {
		return AddHoldingRegisters(server, unitID, address, values)
	}
	return AddInputRegisters(server, unitID, address, values)
}

// SetRegisters sets existing registers' values to the slave. If writeable is true, the registers are holding registers, otherwise they are input registers.
func SetRegisters(server Server, unitID byte, address uint16, values [][]byte, writeable bool) error {
	if writeable {
		return SetHoldingRegisters(server, unitID, address, values)
	}
	return SetInputRegisters(server, unitID, address, values)
}

// DeleteRegisters deletes registers from the slave. If writeable is true, the registers are holding registers, otherwise they are input registers.
func DeleteRegisters(server Server, unitID byte, addresses []uint16, writeable bool) error {
	if writeable {
		return DeleteHoldingRegisters(server, unitID, addresses)
	}
	return DeleteInputRegisters(server, unitID, addresses)
}

// processRequest processes the request and returns the response.
func processRequest(requestPDUBytes []byte, slave *Slave) ([]byte, error) {
	switch requestPDUBytes[0] {
	case gomodbus.ReadCoil:
		responsePDUBytes, err := handleReadCoils(slave, requestPDUBytes)
		return responsePDUBytes, err
	case gomodbus.ReadDiscreteInput:
		responsePDUBytes, err := handleReadDiscreteInputs(slave, requestPDUBytes)
		return responsePDUBytes, err
	case gomodbus.ReadHoldingRegister:
		responsePDUBytes, err := handleReadHoldingRegisters(slave, requestPDUBytes)
		return responsePDUBytes, err
	case gomodbus.ReadInputRegister:
		responsePDUBytes, err := handleReadInputRegisters(slave, requestPDUBytes)
		return responsePDUBytes, err
	case gomodbus.WriteSingleCoil:
		responsePDUBytes, err := handleWriteSingleCoil(slave, requestPDUBytes)
		return responsePDUBytes, err
	case gomodbus.WriteMultipleCoils:
		responsePDUBytes, err := handleWriteMultipleCoils(slave, requestPDUBytes)
		return responsePDUBytes, err
	case gomodbus.WriteSingleRegister:
		responsePDUBytes, err := handleWriteSingleRegister(slave, requestPDUBytes)
		return responsePDUBytes, err
	case gomodbus.WriteMultipleRegisters:
		responsePDUBytes, err := handleWriteMultipleRegisters(slave, requestPDUBytes)
		return responsePDUBytes, err
	default:
		responsePDU := pdu.NewPDUErrorResponse(requestPDUBytes[0], 0x07)
		return responsePDU.ToBytes(), nil
	}
}

// handleReadCoils handles the read coils request.
func handleReadCoils(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	// Parse the request PDU
	var requestPDU pdu.PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the function code is valid
	if requestPDU.FunctionCode != gomodbus.ReadCoil {
		gomodbus.Logger.Sugar().Errorf("invalid function code for read coils: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Read the coils
	coils := make([]bool, requestPDU.Quantity)
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		coil, ok := slave.Coils[i]
		if !ok {
			gomodbus.Logger.Sugar().Errorf("coil not found: %v", err)
			responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadCoil, gomodbus.IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		coils[i-requestPDU.StartingAddress] = coil
	}

	// Create the response PDU
	responsePDU := pdu.NewPDUReadCoilsResponse(coils)
	return responsePDU.ToBytes(), nil
}

// handleReadDiscreteInputs handles the read discrete inputs request.
func handleReadDiscreteInputs(slave *Slave, requestPDUBytes []byte) ([]byte, error) {
	// Parse the request PDU
	var requestPDU pdu.PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadDiscreteInput, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the function code is valid
	if requestPDU.FunctionCode != gomodbus.ReadDiscreteInput {
		gomodbus.Logger.Sugar().Errorf("invalid function code for read discrete inputs: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadDiscreteInput, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Read the discrete inputs
	inputs := make([]bool, requestPDU.Quantity)
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		input, ok := slave.DiscreteInputs[i]
		if !ok {
			gomodbus.Logger.Sugar().Errorf("discrete input not found: %v", err)
			responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadDiscreteInput, gomodbus.IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		inputs[i-requestPDU.StartingAddress] = input
	}

	// Create the response PDU
	responsePDU := pdu.NewPDUReadDiscreteInputsResponse(inputs)
	return responsePDU.ToBytes(), nil
}

// handleReadHoldingRegisters handles the read holding registers request.
func handleReadHoldingRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	// Parse the request PDU
	var requestPDU pdu.PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadHoldingRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the function code is valid
	if requestPDU.FunctionCode != gomodbus.ReadHoldingRegister {
		gomodbus.Logger.Sugar().Errorf("invalid function code for read holding registers: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadHoldingRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Read the holding registers
	var registers bytes.Buffer
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		register, ok := slave.HoldingRegisters[i]
		if !ok {
			gomodbus.Logger.Sugar().Errorf("holding register not found: %v", err)
			responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadHoldingRegister, gomodbus.IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		registers.Write(register)

	}

	// Create the response PDU
	responsePDU := pdu.NewPDUReadHoldingRegistersResponse(registers.Bytes())
	return responsePDU.ToBytes(), nil
}

// handleReadInputRegisters handles the read input registers request.
func handleReadInputRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	// Parse the request PDU
	var requestPDU pdu.PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadInputRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the function code is valid
	if requestPDU.FunctionCode != gomodbus.ReadInputRegister {
		gomodbus.Logger.Sugar().Errorf("invalid function code for read input registers: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadInputRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Read the input registers
	var registers bytes.Buffer
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		register, ok := slave.InputRegisters[i]
		if !ok {
			gomodbus.Logger.Sugar().Errorf("input register not found: %v", err)
			responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadInputRegister, gomodbus.IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		registers.Write(register)
	}

	// Create the response PDU
	responsePDU := pdu.NewPDUReadInputRegistersResponse(registers.Bytes())
	return responsePDU.ToBytes(), nil
}

// handleWriteSingleCoil handles the write single coil request.
func handleWriteSingleCoil(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	// Parse the request PDU
	var requestPDU pdu.PDUWriteSingleCoil
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the function code is valid
	if requestPDU.FunctionCode != gomodbus.WriteSingleCoil {
		gomodbus.Logger.Sugar().Errorf("invalid function code for write single coil: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the coil exists
	_, ok := slave.Coils[requestPDU.Address]
	if !ok {
		gomodbus.Logger.Sugar().Errorf("coil not found: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleCoil, gomodbus.IllegalDataAddress)
		return responsePDU.ToBytes(), nil
	}

	// Set the coil value
	slave.Coils[requestPDU.Address] = requestPDU.Value == 0xFF00

	// Create the response PDU
	responsePDU := pdu.NewWriteSingleCoilResponse(requestPDU.Address, requestPDU.Value == 0xFF00)
	return responsePDU.ToBytes(), nil
}

// handleWriteMultipleCoils handles the write multiple coils request.
func handleWriteMultipleCoils(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	// Parse the request PDU
	var requestPDU pdu.PDUWriteMultipleCoils
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleCoils, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the function code is valid
	if requestPDU.FunctionCode != gomodbus.WriteMultipleCoils {
		gomodbus.Logger.Sugar().Errorf("invalid function code for write multiple coils: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleCoils, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Set the coil values
	for i := 0; i < int(requestPDU.Quantity); i++ {
		if _, ok := slave.Coils[requestPDU.StartingAddress+uint16(i)]; !ok {
			gomodbus.Logger.Sugar().Errorf("coil not found: %v", err)
			responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleCoils, gomodbus.IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		if requestPDU.Values[i/8]&(1<<(i%8)) != 0 {
			slave.Coils[requestPDU.StartingAddress+uint16(i)] = true
		} else {
			slave.Coils[requestPDU.StartingAddress+uint16(i)] = false
		}
	}

	// Create the response PDU
	responsePDU := pdu.NewWriteMultipleCoilsResponse(requestPDU.StartingAddress, requestPDU.Quantity)
	return responsePDU.ToBytes(), nil
}

// handleWriteSingleRegister handles the write single register request.
func handleWriteSingleRegister(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	// Parse the request PDU
	var requestPDU pdu.PDUWriteSingleRegister
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the function code is valid
	if requestPDU.FunctionCode != gomodbus.WriteSingleRegister {
		gomodbus.Logger.Sugar().Errorf("invalid function code for write single register: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the register exists
	_, ok := slave.HoldingRegisters[requestPDU.Address]
	if !ok {
		gomodbus.Logger.Sugar().Errorf("holding register not found: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleRegister, gomodbus.IllegalDataAddress)
		return responsePDU.ToBytes(), nil
	}

	// Set the register value
	slave.HoldingRegisters[requestPDU.Address] = requestPDU.Value

	// Create the response PDU
	responsePDU := pdu.NewWriteSingleRegisterResponse(requestPDU.Address, requestPDU.Value)
	return responsePDU.ToBytes(), nil
}

// handleWriteMultipleRegisters handles the write multiple registers request.
func handleWriteMultipleRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	// Parse the request PDU
	var requestPDU pdu.PDUWriteMultipleRegisters
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the function code is valid
	if requestPDU.FunctionCode != gomodbus.WriteMultipleRegisters {
		gomodbus.Logger.Sugar().Errorf("invalid function code for write multiple registers: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	// Check if the values length is valid
	if len(requestPDU.Values) != int(requestPDU.Quantity)*2 {
		gomodbus.Logger.Sugar().Errorf("invalid values length for write multiple registers, expected %d, got %d", requestPDU.Quantity*2, len(requestPDU.Values))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalDataValue)
		return responsePDU.ToBytes(), nil
	}

	// Set the register values
	for i := 0; i < int(requestPDU.Quantity); i++ {
		if _, ok := slave.HoldingRegisters[requestPDU.StartingAddress+uint16(i)]; !ok {
			gomodbus.Logger.Sugar().Errorf("holding register not found: %v", err)
			responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		slave.HoldingRegisters[requestPDU.StartingAddress+uint16(i)] = requestPDU.Values[2*i : 2*i+2]
	}

	// Create the response PDU
	responsePDU := pdu.NewWriteMultipleRegistersResponse(requestPDU.StartingAddress, requestPDU.Quantity)
	return responsePDU.ToBytes(), nil
}
