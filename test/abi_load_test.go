package test

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"

	witigo "github.com/rioam2/witigo/pkg"
)

type FakeMemory struct {
	bytes []byte
}

func (m *FakeMemory) Read(ptr int32, size int32) ([]byte, error) {
	if ptr < 0 || int(ptr)+int(size) > len(m.bytes) {
		return nil, fmt.Errorf("memory out of bounds: ptr=%d, size=%d, memory size=%d", ptr, size, len(m.bytes))
	}
	return m.bytes[ptr : ptr+size], nil
}

func (m *FakeMemory) Write(ptr int32, data []byte) error {
	if ptr < 0 || int(ptr)+len(data) > len(m.bytes) {
		return fmt.Errorf("memory out of bounds: ptr=%d, data size=%d, memory size=%d", ptr, len(data), len(m.bytes))
	}
	copy(m.bytes[ptr:], data)
	return nil
}

func (m *FakeMemory) Size() int32 {
	return int32(len(m.bytes))
}

func createMemory(bytes []byte) witigo.RuntimeMemory {
	return &FakeMemory{
		bytes: bytes,
	}
}

func TestAbiLoadInt(t *testing.T) {
	tests := []struct {
		name    string
		typeDef witigo.AbiTypeDefinition
		value   any
		offset  int
	}{
		{
			name:    "int8(4) offset 0",
			typeDef: witigo.NewAbiTypeDefinitionS8(),
			value:   int8(4),
			offset:  1,
		},
		{
			name:    "int8(5) offset 34",
			typeDef: witigo.NewAbiTypeDefinitionS8(),
			value:   int8(5),
			offset:  34,
		},
		{
			name:    "int16(27) offset 8",
			typeDef: witigo.NewAbiTypeDefinitionS16(),
			value:   int16(27),
			offset:  8,
		},
		{
			name:    "math.MaxInt32 offset 32",
			typeDef: witigo.NewAbiTypeDefinitionS32(),
			value:   int32(math.MaxInt32),
			offset:  32,
		},
		{
			name:    "math.MaxInt64 offset 64",
			typeDef: witigo.NewAbiTypeDefinitionS64(),
			value:   int64(math.MaxInt64),
			offset:  64,
		},
		{
			name:    "uint8(7) offset 5",
			typeDef: witigo.NewAbiTypeDefinitionU8(),
			value:   uint8(7),
			offset:  5,
		},
		{
			name:    "uint16(100) offset 10",
			typeDef: witigo.NewAbiTypeDefinitionU16(),
			value:   uint16(100),
			offset:  10,
		},
		{
			name:    "uint32(200) offset 20",
			typeDef: witigo.NewAbiTypeDefinitionU32(),
			value:   uint32(200),
			offset:  20,
		},
		{
			name:    "uint64(300) offset 40",
			typeDef: witigo.NewAbiTypeDefinitionU64(),
			value:   uint64(300),
			offset:  40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memoryBuffer := make([]byte, 1024)
			abiOptions := witigo.AbiOptions{
				StringEncoding: witigo.StringEncodingUTF8,
				Memory:         createMemory(memoryBuffer),
			}

			binary.Encode(memoryBuffer[tt.offset:], binary.LittleEndian, tt.value)

			var result any
			var err error
			switch v := tt.value.(type) {
			case int8:
				result, err = witigo.AbiLoadInt[int8](abiOptions, int32(tt.offset), tt.typeDef)
			case int16:
				result, err = witigo.AbiLoadInt[int16](abiOptions, int32(tt.offset), tt.typeDef)
			case int32:
				result, err = witigo.AbiLoadInt[int32](abiOptions, int32(tt.offset), tt.typeDef)
			case int64:
				result, err = witigo.AbiLoadInt[int64](abiOptions, int32(tt.offset), tt.typeDef)
			case uint8:
				result, err = witigo.AbiLoadInt[uint8](abiOptions, int32(tt.offset), tt.typeDef)
			case uint16:
				result, err = witigo.AbiLoadInt[uint16](abiOptions, int32(tt.offset), tt.typeDef)
			case uint32:
				result, err = witigo.AbiLoadInt[uint32](abiOptions, int32(tt.offset), tt.typeDef)
			case uint64:
				result, err = witigo.AbiLoadInt[uint64](abiOptions, int32(tt.offset), tt.typeDef)
			default:
				t.Fatalf("Unsupported type %T for test case %s", v, tt.name)
			}

			if err != nil {
				t.Fatalf("AbiLoadInt failed: %v", err)
			}
			if result != tt.value {
				t.Fatalf("Expected result to be %d, got %d", tt.value, result)
			}
		})
	}
}

func TestAbiLoadFloat(t *testing.T) {
	tests := []struct {
		name    string
		typeDef witigo.AbiTypeDefinition
		value   any
		offset  int
		isNaN   bool
	}{
		{
			name:    "float32(3.14) offset 0",
			typeDef: witigo.NewAbiTypeDefinitionF32(),
			value:   float32(3.14),
			offset:  0,
			isNaN:   false,
		},
		{
			name:    "float32(-42.5) offset 12",
			typeDef: witigo.NewAbiTypeDefinitionF32(),
			value:   float32(-42.5),
			offset:  12,
			isNaN:   false,
		},
		{
			name:    "float32(math.MaxFloat32) offset 24",
			typeDef: witigo.NewAbiTypeDefinitionF32(),
			value:   float32(math.MaxFloat32),
			offset:  24,
			isNaN:   false,
		},
		{
			name:    "float32(NaN) offset 36",
			typeDef: witigo.NewAbiTypeDefinitionF32(),
			value:   float32(math.NaN()),
			offset:  36,
			isNaN:   true,
		},
		{
			name:    "float64(2.71828) offset 4",
			typeDef: witigo.NewAbiTypeDefinitionF64(),
			value:   float64(2.71828),
			offset:  4,
			isNaN:   false,
		},
		{
			name:    "float64(-123.456) offset 16",
			typeDef: witigo.NewAbiTypeDefinitionF64(),
			value:   float64(-123.456),
			offset:  16,
			isNaN:   false,
		},
		{
			name:    "float64(math.MaxFloat64) offset 32",
			typeDef: witigo.NewAbiTypeDefinitionF64(),
			value:   float64(math.MaxFloat64),
			offset:  32,
			isNaN:   false,
		},
		{
			name:    "float64(math.SmallestNonzeroFloat64) offset 48",
			typeDef: witigo.NewAbiTypeDefinitionF64(),
			value:   float64(math.SmallestNonzeroFloat64),
			offset:  48,
			isNaN:   false,
		},
		{
			name:    "float64(NaN) offset 56",
			typeDef: witigo.NewAbiTypeDefinitionF64(),
			value:   math.NaN(),
			offset:  56,
			isNaN:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memoryBuffer := make([]byte, 1024)
			abiOptions := witigo.AbiOptions{
				StringEncoding: witigo.StringEncodingUTF8,
				Memory:         createMemory(memoryBuffer),
			}

			binary.Encode(memoryBuffer[tt.offset:], binary.LittleEndian, tt.value)

			var result any
			var err error
			switch v := tt.value.(type) {
			case float32:
				result, err = witigo.AbiLoadFloat[float32](abiOptions, int32(tt.offset), tt.typeDef)
			case float64:
				result, err = witigo.AbiLoadFloat[float64](abiOptions, int32(tt.offset), tt.typeDef)
			default:
				t.Fatalf("Unsupported type %T for test case %s", v, tt.name)
			}

			if err != nil {
				t.Fatalf("AbiLoadFloat failed: %v", err)
			}

			if tt.isNaN {
				switch v := result.(type) {
				case float32:
					if !math.IsNaN(float64(v)) {
						t.Fatalf("Expected result to be NaN, got %v", v)
					}
				case float64:
					if !math.IsNaN(v) {
						t.Fatalf("Expected result to be NaN, got %v", v)
					}
				}
			} else if result != tt.value {
				t.Fatalf("Expected result to be %v, got %v", tt.value, result)
			}
		})
	}
}

func TestAbiLoadBool(t *testing.T) {
	tests := []struct {
		name    string
		typeDef witigo.AbiTypeDefinition
		value   bool
		offset  int
	}{
		{
			name:    "bool(true) offset 0",
			typeDef: witigo.NewAbiTypeDefinitionBool(),
			value:   true,
			offset:  0,
		},
		{
			name:    "bool(false) offset 10",
			typeDef: witigo.NewAbiTypeDefinitionBool(),
			value:   false,
			offset:  10,
		},
		{
			name:    "bool(true) offset 42",
			typeDef: witigo.NewAbiTypeDefinitionBool(),
			value:   true,
			offset:  42,
		},
		{
			name:    "bool(false) offset 100",
			typeDef: witigo.NewAbiTypeDefinitionBool(),
			value:   false,
			offset:  100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memoryBuffer := make([]byte, 1024)
			abiOptions := witigo.AbiOptions{
				StringEncoding: witigo.StringEncodingUTF8,
				Memory:         createMemory(memoryBuffer),
			}

			binary.Encode(memoryBuffer[tt.offset:], binary.LittleEndian, tt.value)

			result, err := witigo.AbiLoadBool(abiOptions, int32(tt.offset), tt.typeDef)
			if err != nil {
				t.Fatalf("AbiLoadBool failed: %v", err)
			}
			if result != tt.value {
				t.Fatalf("Expected result to be %v, got %v", tt.value, result)
			}
		})
	}
}

func TestAbiLoadChar(t *testing.T) {
	tests := []struct {
		name    string
		typeDef witigo.AbiTypeDefinition
		value   rune
		offset  int
	}{
		{
			name:    "char('A') offset 0",
			typeDef: witigo.NewAbiTypeDefinitionChar(),
			value:   'A',
			offset:  0,
		},
		{
			name:    "char('火') offset 4", // Unicode CJK character
			typeDef: witigo.NewAbiTypeDefinitionChar(),
			value:   '火',
			offset:  4,
		},
		{
			name:    "char('😀') offset 8", // Unicode emoji
			typeDef: witigo.NewAbiTypeDefinitionChar(),
			value:   '😀',
			offset:  8,
		},
		{
			name:    "char('ñ') offset 16", // Unicode with diacritic
			typeDef: witigo.NewAbiTypeDefinitionChar(),
			value:   'ñ',
			offset:  16,
		},
		{
			name:    "char(0) offset 20", // Null character
			typeDef: witigo.NewAbiTypeDefinitionChar(),
			value:   0,
			offset:  20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memoryBuffer := make([]byte, 1024)
			abiOptions := witigo.AbiOptions{
				StringEncoding: witigo.StringEncodingUTF8,
				Memory:         createMemory(memoryBuffer),
			}

			binary.Encode(memoryBuffer[tt.offset:], binary.LittleEndian, tt.value)

			result, err := witigo.AbiLoadChar(abiOptions, int32(tt.offset), tt.typeDef)
			if err != nil {
				t.Fatalf("AbiLoadChar failed: %v", err)
			}
			if result != tt.value {
				t.Fatalf("Expected result to be %v (%U), got %v (%U)", tt.value, tt.value, result, result)
			}
		})
	}
}
