package test

import (
	"encoding/binary"
	"fmt"
	"math"
	"testing"

	witigo "github.com/rioam2/witigo/pkg"
	"golang.org/x/text/encoding/unicode"
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
func TestAbiLoadString(t *testing.T) {
	tests := []struct {
		name           string
		typeDef        witigo.AbiTypeDefinition
		value          string
		ptrOffset      int
		dataOffset     int
		stringEncoding witigo.StringEncoding
	}{
		{
			name:           "string(\"hello\") ptr offset 0, data offset 100, UTF8",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "hello",
			ptrOffset:      0,
			dataOffset:     100,
			stringEncoding: witigo.StringEncodingUTF8,
		},
		{
			name:           "string(\"世界\") ptr offset 8, data offset 120, UTF8",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "世界", // Unicode CJK characters
			ptrOffset:      8,
			dataOffset:     120,
			stringEncoding: witigo.StringEncodingUTF8,
		},
		{
			name:           "string(\"😀👍🚀\") ptr offset 16, data offset 140, UTF8",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "😀👍🚀", // Unicode emojis
			ptrOffset:      16,
			dataOffset:     140,
			stringEncoding: witigo.StringEncodingUTF8,
		},
		{
			name:           "string(\"Café ñandú\") ptr offset 24, data offset 160, UTF8",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "Café ñandú", // Unicode with diacritics
			ptrOffset:      24,
			dataOffset:     160,
			stringEncoding: witigo.StringEncodingUTF8,
		},
		{
			name:           "empty string ptr offset 32, data offset 200, UTF8",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "",
			ptrOffset:      32,
			dataOffset:     200,
			stringEncoding: witigo.StringEncodingUTF8,
		},
		{
			name:           "long string ptr offset 40, data offset 220, UTF8",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "This is a longer string to test handling of strings with more characters",
			ptrOffset:      40,
			dataOffset:     220,
			stringEncoding: witigo.StringEncodingUTF8,
		},
		{
			name:           "string(\"hello\") ptr offset 300, data offset 400, UTF16",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "hello",
			ptrOffset:      300,
			dataOffset:     400,
			stringEncoding: witigo.StringEncodingUTF16,
		},
		{
			name:           "string(\"世界\") ptr offset 308, data offset 420, UTF16",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "世界", // Unicode CJK characters
			ptrOffset:      308,
			dataOffset:     420,
			stringEncoding: witigo.StringEncodingUTF16,
		},
		{
			name:           "string(\"😀👍🚀\") ptr offset 316, data offset 440, UTF16",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "😀👍🚀", // Unicode emojis
			ptrOffset:      316,
			dataOffset:     440,
			stringEncoding: witigo.StringEncodingUTF16,
		},
		{
			name:           "string(\"Café ñandú\") ptr offset 324, data offset 460, UTF16",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "Café ñandú", // Unicode with diacritics
			ptrOffset:      324,
			dataOffset:     460,
			stringEncoding: witigo.StringEncodingUTF16,
		},
		{
			name:           "empty string ptr offset 332, data offset 500, UTF16",
			typeDef:        witigo.NewAbiTypeDefinitionString(),
			value:          "",
			ptrOffset:      332,
			dataOffset:     500,
			stringEncoding: witigo.StringEncodingUTF16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memoryBuffer := make([]byte, 1024)
			abiOptions := witigo.AbiOptions{
				StringEncoding: tt.stringEncoding,
				Memory:         createMemory(memoryBuffer),
			}

			// Write string data to memory based on encoding
			var taggedCodeUnits int
			if tt.stringEncoding == witigo.StringEncodingUTF16 {
				// Convert tt.value to UTF-16 and write to memory
				encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
				utf16String, err := encoder.String(tt.value)
				taggedCodeUnits = len(utf16String) / 2
				if err != nil {
					t.Fatalf("Failed to encode string to UTF-16: %v", err)
				}
				copy(memoryBuffer[tt.dataOffset:], utf16String)
			} else {
				// Write UTF-8 string data to memory
				taggedCodeUnits = len(tt.value)
				copy(memoryBuffer[tt.dataOffset:], []byte(tt.value))
			}

			// Write string pointer structure: [ptr:4bytes][len:4bytes]
			binary.LittleEndian.PutUint32(memoryBuffer[tt.ptrOffset:tt.ptrOffset+4], uint32(tt.dataOffset))
			binary.LittleEndian.PutUint32(memoryBuffer[tt.ptrOffset+4:tt.ptrOffset+8], uint32(taggedCodeUnits))

			result, err := witigo.AbiLoadString(abiOptions, int32(tt.ptrOffset), tt.typeDef)
			if err != nil {
				t.Fatalf("AbiLoadString failed: %v", err)
			}
			if result != tt.value {
				t.Fatalf("Expected result to be %q, got %q", tt.value, result)
			}
		})
	}
}
