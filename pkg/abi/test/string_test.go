package test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadString(t *testing.T) {
	tests := []struct {
		name           string
		memoryMap      map[uint32][]byte
		ptr            uint32
		encoding       abi.StringEncoding
		expectedString string
		expectError    bool
	}{
		{
			name: "read UTF8 string",
			memoryMap: map[uint32][]byte{
				0:   {100, 0, 0, 0}, // String pointer at offset 100
				4:   {5, 0, 0, 0},   // String length 5 code units
				100: []byte("hello"),
			},
			ptr:            0,
			encoding:       abi.StringEncodingUTF8,
			expectedString: "hello",
			expectError:    false,
		},
		{
			name: "read UTF16 string",
			memoryMap: map[uint32][]byte{
				0:   {100, 0, 0, 0},                           // String pointer at offset 100
				4:   {5, 0, 0, 0},                             // String length 5 code units (10 bytes for UTF16)
				100: {'h', 0, 'e', 0, 'l', 0, 'l', 0, 'o', 0}, // UTF16-LE encoding of "hello"
			},
			ptr:            0,
			encoding:       abi.StringEncodingUTF16,
			expectedString: "hello",
			expectError:    false,
		},
		{
			name: "unaligned string pointer",
			memoryMap: map[uint32][]byte{
				0: {101, 0, 0, 0}, // String pointer at offset 101 (unaligned for UTF16)
				4: {5, 0, 0, 0},   // String length 5 code units
			},
			ptr:         0,
			encoding:    abi.StringEncodingUTF16,
			expectError: true,
		},
		{
			name: "out of bounds string pointer",
			memoryMap: map[uint32][]byte{
				0: {0x00, 0x04, 0x00, 0x00}, // String pointer beyond memory bounds (1024)
				4: {5, 0, 0, 0},             // String length 5 code units
			},
			ptr:         0,
			encoding:    abi.StringEncodingUTF8,
			expectError: true,
		},
		{
			name: "out of bounds string length",
			memoryMap: map[uint32][]byte{
				0: {100, 0, 0, 0},     // String pointer at offset 100
				4: {0xFF, 0x03, 0, 0}, // String length too large (0x3FF = 1023)
			},
			ptr:         0,
			encoding:    abi.StringEncodingUTF8,
			expectError: true,
		},
		{
			name: "empty string",
			memoryMap: map[uint32][]byte{
				0:   {100, 0, 0, 0}, // String pointer at offset 100
				4:   {0, 0, 0, 0},   // String length 0 code units
				100: {},             // No data at string pointer
			},
			ptr:            0,
			encoding:       abi.StringEncodingUTF8,
			expectedString: "",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := createMemoryFromMap(tt.memoryMap)

			opts := abi.AbiOptions{
				Memory:         mem,
				StringEncoding: tt.encoding,
			}

			var result string
			err := abi.Read(opts, tt.ptr, &result)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedString, result)
			}
		})
	}
}

func TestReadString_InvalidArgs(t *testing.T) {
	mem := createMemoryFromMap(nil) // Empty memory
	opts := abi.AbiOptions{Memory: mem, StringEncoding: abi.StringEncodingUTF8}

	t.Run("nil pointer", func(t *testing.T) {
		err := abi.Read(opts, 0, nil)
		assert.Error(t, err)
	})

	t.Run("non-pointer result", func(t *testing.T) {
		var str string
		err := abi.Read(opts, 0, str)
		assert.Error(t, err)
	})

	t.Run("wrong pointer type", func(t *testing.T) {
		var num int
		err := abi.Read(opts, 0, &num)
		assert.Error(t, err)
	})
}
