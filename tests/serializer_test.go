package tests

import (
	"testing"

	"github.com/ulfaric/gomodbus"
)

func TestSerializer_SerializeInt16(t *testing.T) {
	var value_int16 int16 = 12345
	registers, err := gomodbus.Serializer(value_int16, gomodbus.BigEndian, gomodbus.LittleEndian)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}
	t.Logf("Registers: %x", registers)
}

func TestSerializer_SerializeInt32(t *testing.T) {
	var value_int32 int32 = 123456
	registers, err := gomodbus.Serializer(value_int32, gomodbus.BigEndian, gomodbus.LittleEndian)
	if err != nil {
		t.Fatalf("Failed to serialize: %v", err)
	}
	t.Logf("Registers: %x", registers)
}