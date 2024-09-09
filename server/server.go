package server

import (
	"bytes"

	"github.com/ulfaric/gomodbus"
	"github.com/ulfaric/gomodbus/pdu"
	"go.uber.org/zap"
)

type Server interface {
	Start() error
	Stop() error
	AddSlave(unitID byte)
	GetSlave(unitID byte) (*Slave, error)
	RemoveSlave(unitID byte)
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
		gomodbus.Logger.Error("failed to parse request", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.ReadCoil {
		gomodbus.Logger.Error("invalid function code for read coils", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	coils := make([]bool, requestPDU.Quantity)
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		coil, ok := slave.Coils[i]
		if !ok {
			gomodbus.Logger.Error("coil not found", zap.Error(err))
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
		gomodbus.Logger.Error("failed to parse request", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadDiscreteInput, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.ReadDiscreteInput {
		gomodbus.Logger.Error("invalid function code for read discrete inputs", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadDiscreteInput, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	inputs := make([]bool, requestPDU.Quantity)
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		input, ok := slave.DiscreteInputs[i]
		if !ok {
			gomodbus.Logger.Error("discrete input not found", zap.Error(err))
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
		gomodbus.Logger.Error("failed to parse request", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadHoldingRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.ReadHoldingRegister {
		gomodbus.Logger.Error("invalid function code for read holding registers", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadHoldingRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	var registers bytes.Buffer
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		register, ok := slave.HoldingRegisters[i]
		if !ok {
			gomodbus.Logger.Error("holding register not found", zap.Error(err))
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
		gomodbus.Logger.Error("failed to parse request", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadInputRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.ReadInputRegister {
		gomodbus.Logger.Error("invalid function code for read input registers", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.ReadInputRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	var registers bytes.Buffer
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		register, ok := slave.InputRegisters[i]
		if !ok {
			gomodbus.Logger.Error("input register not found", zap.Error(err))
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
		gomodbus.Logger.Error("failed to parse request", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.WriteSingleCoil {
		gomodbus.Logger.Error("invalid function code for write single coil", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleCoil, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	_, ok := slave.Coils[requestPDU.Address]
	if !ok {
		gomodbus.Logger.Error("coil not found", zap.Error(err))
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
		gomodbus.Logger.Error("failed to parse request", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleCoils, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.WriteMultipleCoils {
		gomodbus.Logger.Error("invalid function code for write multiple coils", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleCoils, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	for i := 0; i < int(requestPDU.Quantity); i++ {
		if _, ok := slave.Coils[requestPDU.StartingAddress+uint16(i)]; !ok {
			gomodbus.Logger.Error("coil not found", zap.Error(err))
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
		gomodbus.Logger.Error("failed to parse request", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.WriteSingleRegister {
		gomodbus.Logger.Error("invalid function code for write single register", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteSingleRegister, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	_, ok := slave.HoldingRegisters[requestPDU.Address]
	if !ok {
		gomodbus.Logger.Error("holding register not found", zap.Error(err))
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
		gomodbus.Logger.Error("failed to parse request", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != gomodbus.WriteMultipleRegisters {
		gomodbus.Logger.Error("invalid function code for write multiple registers", zap.Error(err))
		responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	for i := 0; i < int(requestPDU.Quantity); i++ {
		if _, ok := slave.HoldingRegisters[requestPDU.StartingAddress+uint16(i)]; !ok {
			gomodbus.Logger.Error("holding register not found", zap.Error(err))
			responsePDU := pdu.NewPDUErrorResponse(gomodbus.WriteMultipleRegisters, gomodbus.IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		slave.HoldingRegisters[requestPDU.StartingAddress+uint16(i)] = requestPDU.Values[2*i:]
	}

	responsePDU := pdu.NewWriteMultipleRegistersResponse(requestPDU.StartingAddress, requestPDU.Quantity)
	return responsePDU.ToBytes(), nil
}

