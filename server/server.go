package server

import (
	"bytes"

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

func AddCoil(server Server, unitID byte, address uint16, value bool) {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return
	}
	slave.AddCoil(address, value)
}

func AddCoils(server Server, unitID byte, address uint16, values []bool) {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return
	}
	slave.AddCoils(address, values)
}

func AddDiscreteInput(server Server, unitID byte, address uint16, value bool) {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return
	}
	slave.AddDiscreteInput(address, value)
}

func AddDiscreteInputs(server Server, unitID byte, address uint16, values []bool) {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return
	}
	slave.AddDiscreteInputs(address, values)
}

func AddHoldingRegister(server Server, unitID byte, address uint16, value []byte) {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return
	}
	slave.AddHoldingRegister(address, value)
}

func AddHoldingRegisters(server Server, unitID byte, address uint16, values [][]byte) {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return
	}
	slave.AddHoldingRegisters(address, values)
}

func AddInputRegister(server Server, unitID byte, address uint16, value []byte) {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return
	}
	slave.AddInputRegister(address, value)
}

func AddInputRegisters(server Server, unitID byte, address uint16, values [][]byte) {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("slave not found: %v", err)
		return
	}
	slave.AddInputRegisters(address, values)
}


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

func handleReadCoils(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU pdu.PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.ReadCoil {
		gomodbus.Logger.Sugar().Errorf("invalid function code for read coils: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

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

	responsePDU := pdu.NewPDUReadCoilsResponse(coils)
	return responsePDU.ToBytes(), nil
}

func handleReadDiscreteInputs(slave *Slave, requestPDUBytes []byte) ([]byte, error) {
	var requestPDU pdu.PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadDiscreteInput, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.ReadDiscreteInput {
		gomodbus.Logger.Sugar().Errorf("invalid function code for read discrete inputs: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadDiscreteInput, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

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

	responsePDU := pdu.NewPDUReadDiscreteInputsResponse(inputs)
	return responsePDU.ToBytes(), nil
}

func handleReadHoldingRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU pdu.PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadHoldingRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.ReadHoldingRegister {
		gomodbus.Logger.Sugar().Errorf("invalid function code for read holding registers: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadHoldingRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

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

	responsePDU := pdu.NewPDUReadHoldingRegistersResponse(registers.Bytes())
	return responsePDU.ToBytes(), nil
}

func handleReadInputRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU pdu.PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadInputRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.ReadInputRegister {
		gomodbus.Logger.Sugar().Errorf("invalid function code for read input registers: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadInputRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

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

	responsePDU := pdu.NewPDUReadInputRegistersResponse(registers.Bytes())
	return responsePDU.ToBytes(), nil
}

func handleWriteSingleCoil(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU pdu.PDUWriteSingleCoil
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.WriteSingleCoil {
		gomodbus.Logger.Sugar().Errorf("invalid function code for write single coil: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	_, ok := slave.Coils[requestPDU.Address]
	if !ok {
		gomodbus.Logger.Sugar().Errorf("coil not found: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleCoil, gomodbus.IllegalDataAddress)
		return responsePDU.ToBytes(), nil
	}

	slave.Coils[requestPDU.Address] = requestPDU.Value == 0xFF00

	responsePDU := pdu.NewWriteSingleCoilResponse(requestPDU.Address, requestPDU.Value == 0xFF00)
	return responsePDU.ToBytes(), nil
}

func handleWriteMultipleCoils(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU pdu.PDUWriteMultipleCoils
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleCoils, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.WriteMultipleCoils {
		gomodbus.Logger.Sugar().Errorf("invalid function code for write multiple coils: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleCoils, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

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

	responsePDU := pdu.NewWriteMultipleCoilsResponse(requestPDU.StartingAddress, requestPDU.Quantity)
	return responsePDU.ToBytes(), nil
}

func handleWriteSingleRegister(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU pdu.PDUWriteSingleRegister
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.WriteSingleRegister {
		gomodbus.Logger.Sugar().Errorf("invalid function code for write single register: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	_, ok := slave.HoldingRegisters[requestPDU.Address]
	if !ok {
		gomodbus.Logger.Sugar().Errorf("holding register not found: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleRegister, gomodbus.IllegalDataAddress)
		return responsePDU.ToBytes(), nil
	}

	slave.HoldingRegisters[requestPDU.Address] = requestPDU.Value

	responsePDU := pdu.NewWriteSingleRegisterResponse(requestPDU.Address, requestPDU.Value)
	return responsePDU.ToBytes(), nil
}

func handleWriteMultipleRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU pdu.PDUWriteMultipleRegisters
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		gomodbus.Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.WriteMultipleRegisters {
		gomodbus.Logger.Sugar().Errorf("invalid function code for write multiple registers: %v", err)
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	for i := 0; i < int(requestPDU.Quantity); i++ {
		if _, ok := slave.HoldingRegisters[requestPDU.StartingAddress+uint16(i)]; !ok {
			gomodbus.Logger.Sugar().Errorf("holding register not found: %v", err)
			responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		slave.HoldingRegisters[requestPDU.StartingAddress+uint16(i)] = requestPDU.Values[2*i:]
	}

	responsePDU := pdu.NewWriteMultipleRegistersResponse(requestPDU.StartingAddress, requestPDU.Quantity)
	return responsePDU.ToBytes(), nil
}

