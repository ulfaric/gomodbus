package gomodbus

import (
	"bytes"
	"errors"
)

type Server interface {
	Start() error
	Stop() error
	AddSlave(unitID byte)
	GetSlave(unitID byte) (*Slave, error)
	RemoveSlave(unitID byte)
}

func AddCoils(server Server, unitID byte, address uint16, values []bool) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.AddCoils(address, values)
	return nil
}

func DeleteCoils(server Server, unitID byte, addresses []uint16) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.DeleteCoils(addresses)
	return nil
}


func AddDiscreteInputs(server Server, unitID byte, address uint16, values []bool) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.AddDiscreteInputs(address, values)
	return nil
}

func DeleteDiscreteInputs(server Server, unitID byte, addresses []uint16) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.DeleteDiscreteInputs(addresses)
	return nil
}


func AddHoldingRegisters(server Server, unitID byte, address uint16, values [][]byte) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.AddHoldingRegisters(address, values)
	return nil
}

func DeleteHoldingRegisters(server Server, unitID byte, addresses []uint16) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.DeleteHoldingRegisters(addresses)
	return nil
}


func AddInputRegisters(server Server, unitID byte, address uint16, values [][]byte) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.AddInputRegisters(address, values)
	return nil
}

func DeleteInputRegisters(server Server, unitID byte, addresses []uint16) error {
	slave, err := server.GetSlave(unitID)
	if err != nil {
		Logger.Sugar().Errorf("slave not found: %v", err)
		return errors.New("slave not found")
	}
	slave.DeleteInputRegisters(addresses)
	return nil
}


func processRequest(requestPDUBytes []byte, slave *Slave) ([]byte, error) {
	switch requestPDUBytes[0] {
	case ReadCoil:
		responsePDUBytes, err := handleReadCoils(slave, requestPDUBytes)
		return responsePDUBytes, err
	case ReadDiscreteInput:
		responsePDUBytes, err := handleReadDiscreteInputs(slave, requestPDUBytes)
		return responsePDUBytes, err
	case ReadHoldingRegister:
		responsePDUBytes, err := handleReadHoldingRegisters(slave, requestPDUBytes)
		return responsePDUBytes, err
	case ReadInputRegister:
		responsePDUBytes, err := handleReadInputRegisters(slave, requestPDUBytes)
		return responsePDUBytes, err
	case WriteSingleCoil:
		responsePDUBytes, err := handleWriteSingleCoil(slave, requestPDUBytes)
		return responsePDUBytes, err
	case WriteMultipleCoils:
		responsePDUBytes, err := handleWriteMultipleCoils(slave, requestPDUBytes)
		return responsePDUBytes, err
	case WriteSingleRegister:
		responsePDUBytes, err := handleWriteSingleRegister(slave, requestPDUBytes)
		return responsePDUBytes, err
	case WriteMultipleRegisters:
		responsePDUBytes, err := handleWriteMultipleRegisters(slave, requestPDUBytes)
		return responsePDUBytes, err
	default:
		responsePDU := NewPDUErrorResponse(requestPDUBytes[0], 0x07)
		return responsePDU.ToBytes(), nil
	}
}

func handleReadCoils(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := NewPDUErrorResponse(ReadCoil, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != ReadCoil {
		Logger.Sugar().Errorf("invalid function code for read coils: %v", err)
		responsePDU := NewPDUErrorResponse(ReadCoil, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	coils := make([]bool, requestPDU.Quantity)
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		coil, ok := slave.Coils[i]
		if !ok {
			Logger.Sugar().Errorf("coil not found: %v", err)
			responsePDU := NewPDUErrorResponse(ReadCoil, IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		coils[i-requestPDU.StartingAddress] = coil
	}

	responsePDU := NewPDUReadCoilsResponse(coils)
	return responsePDU.ToBytes(), nil
}

func handleReadDiscreteInputs(slave *Slave, requestPDUBytes []byte) ([]byte, error) {
	var requestPDU PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := NewPDUErrorResponse(ReadDiscreteInput, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != ReadDiscreteInput {
		Logger.Sugar().Errorf("invalid function code for read discrete inputs: %v", err)
		responsePDU := NewPDUErrorResponse(ReadDiscreteInput, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	inputs := make([]bool, requestPDU.Quantity)
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		input, ok := slave.DiscreteInputs[i]
		if !ok {
			Logger.Sugar().Errorf("discrete input not found: %v", err)
			responsePDU := NewPDUErrorResponse(ReadDiscreteInput, IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		inputs[i-requestPDU.StartingAddress] = input
	}

	responsePDU := NewPDUReadDiscreteInputsResponse(inputs)
	return responsePDU.ToBytes(), nil
}

func handleReadHoldingRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := NewPDUErrorResponse(ReadHoldingRegister, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != ReadHoldingRegister {
		Logger.Sugar().Errorf("invalid function code for read holding registers: %v", err)
		responsePDU := NewPDUErrorResponse(ReadHoldingRegister, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	var registers bytes.Buffer
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		register, ok := slave.HoldingRegisters[i]
		if !ok {
			Logger.Sugar().Errorf("holding register not found: %v", err)
			responsePDU := NewPDUErrorResponse(ReadHoldingRegister, IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		registers.Write(register)

	}

	responsePDU := NewPDUReadHoldingRegistersResponse(registers.Bytes())
	return responsePDU.ToBytes(), nil
}

func handleReadInputRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU PDURead
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := NewPDUErrorResponse(ReadInputRegister, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != ReadInputRegister {
		Logger.Sugar().Errorf("invalid function code for read input registers: %v", err)
		responsePDU := NewPDUErrorResponse(ReadInputRegister, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	var registers bytes.Buffer
	for i := requestPDU.StartingAddress; i < requestPDU.StartingAddress+requestPDU.Quantity; i++ {
		register, ok := slave.InputRegisters[i]
		if !ok {
			Logger.Sugar().Errorf("input register not found: %v", err)
			responsePDU := NewPDUErrorResponse(ReadInputRegister, IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		registers.Write(register)
	}

	responsePDU := NewPDUReadInputRegistersResponse(registers.Bytes())
	return responsePDU.ToBytes(), nil
}

func handleWriteSingleCoil(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU PDUWriteSingleCoil
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := NewPDUErrorResponse(WriteSingleCoil, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != WriteSingleCoil {
		Logger.Sugar().Errorf("invalid function code for write single coil: %v", err)
		responsePDU := NewPDUErrorResponse(WriteSingleCoil, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	_, ok := slave.Coils[requestPDU.Address]
	if !ok {
		Logger.Sugar().Errorf("coil not found: %v", err)
		responsePDU := NewPDUErrorResponse(WriteSingleCoil, IllegalDataAddress)
		return responsePDU.ToBytes(), nil
	}

	slave.Coils[requestPDU.Address] = requestPDU.Value == 0xFF00

	responsePDU := NewWriteSingleCoilResponse(requestPDU.Address, requestPDU.Value == 0xFF00)
	return responsePDU.ToBytes(), nil
}

func handleWriteMultipleCoils(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU PDUWriteMultipleCoils
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := NewPDUErrorResponse(WriteMultipleCoils, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != WriteMultipleCoils {
		Logger.Sugar().Errorf("invalid function code for write multiple coils: %v", err)
		responsePDU := NewPDUErrorResponse(WriteMultipleCoils, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	for i := 0; i < int(requestPDU.Quantity); i++ {
		if _, ok := slave.Coils[requestPDU.StartingAddress+uint16(i)]; !ok {
			Logger.Sugar().Errorf("coil not found: %v", err)
			responsePDU := NewPDUErrorResponse(WriteMultipleCoils, IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		if requestPDU.Values[i/8]&(1<<(i%8)) != 0 {
			slave.Coils[requestPDU.StartingAddress+uint16(i)] = true
		} else {
			slave.Coils[requestPDU.StartingAddress+uint16(i)] = false
		}
	}

	responsePDU := NewWriteMultipleCoilsResponse(requestPDU.StartingAddress, requestPDU.Quantity)
	return responsePDU.ToBytes(), nil
}

func handleWriteSingleRegister(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU PDUWriteSingleRegister
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := NewPDUErrorResponse(WriteSingleRegister, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != WriteSingleRegister {
		Logger.Sugar().Errorf("invalid function code for write single register: %v", err)
		responsePDU := NewPDUErrorResponse(WriteSingleRegister, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	_, ok := slave.HoldingRegisters[requestPDU.Address]
	if !ok {
		Logger.Sugar().Errorf("holding register not found: %v", err)
		responsePDU := NewPDUErrorResponse(WriteSingleRegister, IllegalDataAddress)
		return responsePDU.ToBytes(), nil
	}

	slave.HoldingRegisters[requestPDU.Address] = requestPDU.Value

	responsePDU := NewWriteSingleRegisterResponse(requestPDU.Address, requestPDU.Value)
	return responsePDU.ToBytes(), nil
}

func handleWriteMultipleRegisters(slave *Slave, requestPDUBytes []byte) ([]byte, error) {

	var requestPDU PDUWriteMultipleRegisters
	err := requestPDU.FromBytes(requestPDUBytes)
	if err != nil {
		Logger.Sugar().Errorf("failed to parse request: %v", err)
		responsePDU := NewPDUErrorResponse(WriteMultipleRegisters, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if requestPDU.FunctionCode != WriteMultipleRegisters {
		Logger.Sugar().Errorf("invalid function code for write multiple registers: %v", err)
		responsePDU := NewPDUErrorResponse(WriteMultipleRegisters, IllegalFunction)
		return responsePDU.ToBytes(), nil
	}

	if len(requestPDU.Values) != int(requestPDU.Quantity)*2 {
		Logger.Sugar().Errorf("invalid values length for write multiple registers, expected %d, got %d", requestPDU.Quantity*2, len(requestPDU.Values))
		responsePDU := NewPDUErrorResponse(WriteMultipleRegisters, IllegalDataValue)
		return responsePDU.ToBytes(), nil
	}

	for i := 0; i < int(requestPDU.Quantity); i++ {
		if _, ok := slave.HoldingRegisters[requestPDU.StartingAddress+uint16(i)]; !ok {
			Logger.Sugar().Errorf("holding register not found: %v", err)
			responsePDU := NewPDUErrorResponse(WriteMultipleRegisters, IllegalDataAddress)
			return responsePDU.ToBytes(), nil
		}
		slave.HoldingRegisters[requestPDU.StartingAddress+uint16(i)] = requestPDU.Values[2*i : 2*i+2]
	}

	responsePDU := NewWriteMultipleRegistersResponse(requestPDU.StartingAddress, requestPDU.Quantity)
	return responsePDU.ToBytes(), nil
}

