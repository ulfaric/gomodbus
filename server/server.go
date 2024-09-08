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
	RemoveSlave(unitID byte)
}

func HandleReadCoils(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

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

func HandleReadDiscreteInputs(slave *Slave, requestPDUBytes []byte) ([]byte, error) {
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

func HandleReadHoldingRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

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

func HandleReadInputRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

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
