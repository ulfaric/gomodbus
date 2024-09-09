package gomodbus

import (
	"encoding/binary"
	"fmt"
	"math"
)

func reverse(slice [][]byte) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func Serializer(data interface{}, byteOrder string, wordOrder string) ([][]byte, error) {
	switch v := data.(type) {
	case bool:
		if v {
			return [][]byte{{0xff, 0x00}}, nil
		} else {
			return [][]byte{{0x00, 0x00}}, nil
		}
	case int16:
		bytes := make([]byte, 2)
		if byteOrder == BigEndian {
			binary.BigEndian.PutUint16(bytes, uint16(v))
		} else {
			binary.LittleEndian.PutUint16(bytes, uint16(v))
		}
		return [][]byte{bytes}, nil
	case uint16:
		bytes := make([]byte, 2)
		if byteOrder == BigEndian {
			binary.BigEndian.PutUint16(bytes, v)
		} else {
			binary.LittleEndian.PutUint16(bytes, v)
		}
		return [][]byte{bytes}, nil
	case int32:
		bytes := make([]byte, 4)
		if byteOrder == BigEndian {
			binary.BigEndian.PutUint32(bytes, uint32(v))
		} else {
			binary.LittleEndian.PutUint32(bytes, uint32(v))
		}
		registers := make([][]byte, 0)
		for i := 0; i < len(bytes); i = i + 2 {
			registers = append(registers, bytes[i:i+2])
		}
		if wordOrder == LittleEndian {
			reverse(registers)
		}
		return registers, nil
	case uint32:
		bytes := make([]byte, 4)
		if byteOrder == BigEndian {
			binary.BigEndian.PutUint32(bytes, v)
		} else {
			binary.LittleEndian.PutUint32(bytes, v)
		}
		registers := make([][]byte, 0)
		for i := 0; i < len(bytes); i = i + 2 {
			registers = append(registers, bytes[i:i+2])
		}
		if wordOrder == LittleEndian {
			reverse(registers)
		}
		return registers, nil
	case int64:
		bytes := make([]byte, 8)
		if byteOrder == BigEndian {
			binary.BigEndian.PutUint64(bytes, uint64(v))
		} else {
			binary.LittleEndian.PutUint64(bytes, uint64(v))
		}
		registers := make([][]byte, 0)
		for i := 0; i < len(bytes); i = i + 2 {
			registers = append(registers, bytes[i:i+2])
		}
		if wordOrder == LittleEndian {
			reverse(registers)
		}
		return registers, nil
	case uint64:
		bytes := make([]byte, 8)
		if byteOrder == BigEndian {
			binary.BigEndian.PutUint64(bytes, v)
		} else {
			binary.LittleEndian.PutUint64(bytes, v)
		}
		registers := make([][]byte, 0)
		for i := 0; i < len(bytes); i = i + 2 {
			registers = append(registers, bytes[i:i+2])
		}
		if wordOrder == LittleEndian {
			reverse(registers)
		}
		return registers, nil
	case float32:
		bytes := make([]byte, 4)
		if byteOrder == BigEndian {
			binary.BigEndian.PutUint32(bytes, math.Float32bits(v))
		} else {
			binary.LittleEndian.PutUint32(bytes, math.Float32bits(v))
		}
		registers := make([][]byte, 0)
		for i := 0; i < len(bytes); i = i + 2 {
			registers = append(registers, bytes[i:i+2])
		}
		if wordOrder == LittleEndian {
			reverse(registers)
		}
		return registers, nil
	case float64:
		bytes := make([]byte, 8)
		if byteOrder == BigEndian {
			binary.BigEndian.PutUint64(bytes, math.Float64bits(v))
		} else {
			binary.LittleEndian.PutUint64(bytes, math.Float64bits(v))
		}
		registers := make([][]byte, 0)
		for i := 0; i < len(bytes); i = i + 2 {
			registers = append(registers, bytes[i:i+2])
		}
		if wordOrder == LittleEndian {
			reverse(registers)
		}
		return registers, nil
	case string:
		bytes := []byte(v)
		return [][]byte{bytes}, nil
	default:
		return nil, fmt.Errorf("unsupported data type: %T", v)
	}
}

func Deserializer(data [][]byte, dataType string, byteOrder string, wordOrder string) (interface{}, error) {
	if wordOrder == LittleEndian {
		reverse(data)
	}
	switch dataType {
	case "int16":
		if byteOrder == BigEndian {
			return int16(binary.BigEndian.Uint16(data[0])), nil
		} else {
			return int16(binary.LittleEndian.Uint16(data[0])), nil
		}
	case "uint16":
		if byteOrder == BigEndian {
			return binary.BigEndian.Uint16(data[0]), nil
		} else {
			return binary.LittleEndian.Uint16(data[0]), nil
		}
	case "int32":
		bytes := make([]byte, 4)
		for i, b := range data {
			copy(bytes[i*2:], b)
		}
		if byteOrder == BigEndian {
			return int32(binary.BigEndian.Uint32(bytes)), nil
		} else {
			return int32(binary.LittleEndian.Uint32(bytes)), nil
		}
	case "uint32":
		bytes := make([]byte, 4)
		for i, b := range data {
			copy(bytes[i*2:], b)
		}
		if byteOrder == BigEndian {
			return binary.BigEndian.Uint32(bytes), nil
		} else {
			return binary.LittleEndian.Uint32(bytes), nil
		}
	case "int64":
		bytes := make([]byte, 8)
		for i, b := range data {
			copy(bytes[i*2:], b)
		}
		if byteOrder == BigEndian {
			return int64(binary.BigEndian.Uint64(bytes)), nil
		} else {
			return int64(binary.LittleEndian.Uint64(bytes)), nil
		}
	case "uint64":
		bytes := make([]byte, 8)
		for i, b := range data {
			copy(bytes[i*2:], b)
		}
		if byteOrder == BigEndian {
			return binary.BigEndian.Uint64(bytes), nil
		} else {
			return binary.LittleEndian.Uint64(bytes), nil
		}
	case "float32":
		bytes := make([]byte, 4)
		for i, b := range data {
			copy(bytes[i*2:], b)
		}
		if byteOrder == BigEndian {
			return math.Float32frombits(binary.BigEndian.Uint32(bytes)), nil
		} else {
			return math.Float32frombits(binary.LittleEndian.Uint32(bytes)), nil
		}
	case "float64":
		bytes := make([]byte, 8)
		for i, b := range data {
			copy(bytes[i*2:], b)
		}
		if byteOrder == BigEndian {
			return math.Float64frombits(binary.BigEndian.Uint64(bytes)), nil
		} else {
			return math.Float64frombits(binary.LittleEndian.Uint64(bytes)), nil
		}
	case "string":
		bytes := make([]byte, 0)
		for _, b := range data {
			bytes = append(bytes, b...)
		}
		return string(bytes), nil
	default:
		return nil, fmt.Errorf("unsupported data type: %s", dataType)
	}
}
