package gomodbus

import (
	"fmt"
	"math"
)


func CalculateADULength(functionCode byte, quantity int) (int, error) {
	var pduLength int
	const MBAP_HEADER_LENGTH = 7
	switch functionCode {
	case 0x01, 0x02: // Read Coils, Read Discrete Inputs
		// dataLength is the number of coils/inputs
		byteCount := int(math.Ceil(float64(quantity) / 8.0)) // Number of bytes needed
		pduLength = 1 + 1 + byteCount                        // Function Code + Byte Count + Data
	case 0x03, 0x04: // Read Holding Registers, Read Input Registers
		// dataLength is the number of registers
		byteCount := quantity * 2     // 2 bytes per register
		pduLength = 1 + 1 + byteCount // Function Code + Byte Count + Data
	case 0x05: // Write Single Coil
		// The response PDU for Write Single Coil includes the function code, output address, and output value
		pduLength = 1 + 2 + 2 // Function Code + Output Address + Output Value
	case 0x06: // Write Single Register
		// The response PDU for Write Single Register includes the function code, register address, and register value
		pduLength = 1 + 2 + 2 // Function Code + Register Address + Register Value
	case 0x0F: // Write Multiple Coils
		// The response PDU for Write Multiple Coils includes the function code, starting address, and quantity of coils
		pduLength = 1 + 2 + 2 // Function Code + Starting Address + Quantity of Coils
	case 0x10: // Write Multiple Registers
		// The response PDU for Write Multiple Registers includes the function code, starting address, and quantity of registers
		pduLength = 1 + 2 + 2 // Function Code + Starting Address + Quantity of Registers
	default:
		// Other function codes (handle accordingly)
		return 0, fmt.Errorf("function code not implemented")
	}

	return MBAP_HEADER_LENGTH + pduLength, nil
}

func DecodeModbusRegisters(registers []uint16, byte_order, word_order string) interface{} {
	swapBytes := func(value uint16) uint16 {
		return (value>>8)&0xFF | (value&0xFF)<<8
	}

	switch len(registers) {
	case 1:
		// Swap bytes for single uint16
		if byte_order == LittleEndian {
			return swapBytes(registers[0])
		} else {
			return registers[0]
		}
	case 2:
		// Combine two uint16 into a single 32-bit integer
		var combined uint32
		if word_order == LittleEndian {
			for i := 1; i >= 0; i-- {
				if byte_order == LittleEndian {
					swapped := swapBytes(registers[i])
					combined = (combined << 16) | uint32(swapped)
				} else {
					combined = (combined << 16) | uint32(registers[i])
				}
			}
		} else {
			for i := 0; i < 2; i++ {
				if byte_order == LittleEndian {
					swapped := swapBytes(registers[i])
					combined = (combined << 16) | uint32(swapped)
				} else {
					combined = (combined << 16) | uint32(registers[i])
				}
			}
		}
		return combined
	case 4:
		// Combine four uint16 into a single 64-bit integer
		var combined uint64
		if word_order == LittleEndian {
			for i := 3; i >= 0; i-- {
				if byte_order == LittleEndian {
					swapped := swapBytes(registers[i])
					combined = (combined << 16) | uint64(swapped)
				} else {
					combined = (combined << 16) | uint64(registers[i])
				}
			}
		} else {
			for i := 0; i < 4; i++ {
				if byte_order == LittleEndian {
					swapped := swapBytes(registers[i])
					combined = (combined << 16) | uint64(swapped)
				} else {
					combined = (combined << 16) | uint64(registers[i])
				}
			}
		}
		return combined
	default:
		panic("The input must contain exactly 1, 2, or 4 uint16 values.")
	}
}

func EncodeModbusRegisters(value interface{}, byte_order, word_order string) []uint16 {
	splitUint32 := func(value uint32) []uint16 {
		return []uint16{uint16(value >> 16), uint16(value & 0xFFFF)}
	}

	splitUint64 := func(value uint64) []uint16 {
		return []uint16{uint16(value >> 48), uint16(value >> 32), uint16(value >> 16), uint16(value)}
	}

	swapBytes := func(value uint16) uint16 {
		return (value>>8)&0xFF | (value&0xFF)<<8
	}

	reverseSlice := func(slice []uint16) {
		for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
			slice[i], slice[j] = slice[j], slice[i]
		}
	}

	switch value := value.(type) {
	case uint16:
		// Swap bytes for single uint16
		if byte_order == LittleEndian {
			return []uint16{swapBytes(value)}
		} else {
			return []uint16{value}
		}
	case int16:
		// convert int16 to uint16 and swap bytes
		if byte_order == LittleEndian {
			return []uint16{swapBytes(uint16(value))}
		} else {
			return []uint16{uint16(value)}
		}
	case uint32:
		// Split 32-bit integer into two uint16
		registers := splitUint32(value)
		if word_order == LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case int32:
		// convert int32 to uint32 and split into two uint16
		registers := splitUint32(uint32(value))
		if word_order == LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case uint64:
		// Split 64-bit integer into four uint16
		registers := splitUint64(value)
		if word_order == LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case int64:
		// Convert int64 to uint64 and split into four uint16
		registers := splitUint64(uint64(value))
		if word_order == LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case float32:
		// Convert float32 to uint32 and split into two uint16
		bits := math.Float32bits(value)
		registers := splitUint32(uint32(bits))
		if word_order == LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	case float64:
		// Convert float64 to uint64 and split into four uint16
		bits := math.Float64bits(value)
		registers := splitUint64(bits)
		if word_order == LittleEndian {
			reverseSlice(registers)
		}
		if byte_order == LittleEndian {
			for i := range registers {
				registers[i] = swapBytes(registers[i])
			}
		}
		return registers
	default:
		panic("The input must be a uint16, uint32, or uint64 value.")
	}
}

// CheckModbusError checks if the response ADU indicates an error.
func CheckModbusError(response []byte) error {
	// Minimum length for a valid Modbus response ADU
	if len(response) < 9 {
		return fmt.Errorf("invalid response length: %d", len(response))
	}

	// Extract the function code from the response
	functionCode := response[7]

	// Check if the most significant bit of the function code is set
	if functionCode&0x80 != 0 {
		exceptionCode := response[8]
		return fmt.Errorf("modbus exception: %s", getExceptionMessage(exceptionCode))
	}

	return nil
}

// getExceptionMessage returns a human-readable message for a given exception code.
func getExceptionMessage(exceptionCode byte) string {
	switch exceptionCode {
	case 0x01:
		return "Illegal Function"
	case 0x02:
		return "Illegal Data Address"
	case 0x03:
		return "Illegal Data Value"
	case 0x04:
		return "Server Device Failure"
	default:
		return fmt.Sprintf("Unknown exception code: %02X", exceptionCode)
	}
}
