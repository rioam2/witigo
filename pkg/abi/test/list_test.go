package test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
	"github.com/stretchr/testify/assert"
)

func TestReadListString(t *testing.T) {
	tests := []struct {
		name        string
		memory      map[uint32][]uint8
		ptr         uint32
		expected    []string
		shouldError bool
	}{
		{
			name: "valid string list",
			memory: map[uint32][]uint8{
				0x00: {0x0C, 0, 0, 0},           // list ptr = 0x0C
				0x04: {3, 0, 0, 0},              // list length = 3
				0x0C: {0x24, 0, 0, 0},           // string 1 pointer = 0x24
				0x10: {5, 0, 0, 0},              // string 1 length = 5
				0x14: {0x29, 0, 0, 0},           // string 2 pointer = 0x29
				0x18: {5, 0, 0, 0},              // string 2 length = 5
				0x1C: {0x2E, 0, 0, 0},           // string 3 pointer = 0x2E
				0x20: {4, 0, 0, 0},              // string 3 length = 4
				0x24: {'h', 'e', 'l', 'l', 'o'}, // string 1 data: "hello"
				0x29: {'w', 'o', 'r', 'l', 'd'}, // string 2 data: "world"
				0x2E: {'t', 'e', 's', 't'},      // string 3 data: "test"
			},
			ptr:         0,
			expected:    []string{"hello", "world", "test"},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := abi.AbiOptions{
				StringEncoding: abi.StringEncodingUTF8,
				Memory:         createMemoryFromMap(tt.memory),
			}
			var result []string
			err := abi.Read(opts, tt.ptr, &result)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestReadListUint32(t *testing.T) {
	tests := []struct {
		name        string
		memory      []byte
		ptr         uint32
		expected    []uint32
		shouldError bool
	}{
		{
			name: "valid uint32 list",
			memory: []byte{
				// 0x00: list ptr = 0x0C
				0x0C, 0, 0, 0,
				// 0x04: list length = 3
				3, 0, 0, 0,
				// 0x08: unused
				0, 0, 0, 0,
				// 0x0C: uint32 data = [1, 2, 3]
				1, 0, 0, 0,
				2, 0, 0, 0,
				3, 0, 0, 0,
			},
			ptr:         0,
			expected:    []uint32{1, 2, 3},
			shouldError: false,
		},
		{
			name:        "nil result pointer",
			memory:      []byte{0, 0, 0, 0},
			ptr:         0,
			shouldError: true,
		},
		{
			name:        "invalid list pointer",
			memory:      []byte{},
			ptr:         0,
			expected:    []uint32{},
			shouldError: true,
		},
		{
			name:        "invalid list length",
			memory:      []byte{0, 0, 0, 0},
			ptr:         0,
			expected:    []uint32{},
			shouldError: true,
		},
		{
			name: "empty list",
			memory: []byte{
				// list ptr at 0: 12
				12, 0, 0, 0,
				// list length at 4: 0
				0, 0, 0, 0,
			},
			ptr:         0,
			expected:    []uint32{},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := abi.AbiOptions{
				StringEncoding: abi.StringEncodingUTF8,
				Memory:         createMemory(tt.memory),
			}
			var result []uint32
			err := abi.Read(opts, tt.ptr, &result)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestReadListOfListStrings(t *testing.T) {
	tests := []struct {
		name        string
		memory      map[uint32][]uint8
		ptr         uint32
		expected    [][]string
		shouldError bool
	}{
		{
			name: "valid list of lists",
			memory: map[uint32][]uint8{
				0x00: {0x08, 0, 0, 0}, // list[*][*] ptr
				0x04: {2, 0, 0, 0},    // list[*][*] length
				0x08: {0x50, 0, 0, 0}, // list[0][*] pointer
				0x0C: {3, 0, 0, 0},    // list[0][*] sublist length
				0x10: {0xA0, 0, 0, 0}, // list[1][*] sublist pointer
				0x14: {2, 0, 0, 0},    // list[1][*] sublist length
				0x50: {0x70, 0, 0, 0}, // list[0][0] pointer
				0x54: {1, 0, 0, 0},    // list[0][0] length
				0x58: {0x71, 0, 0, 0}, // list[0][1] pointer
				0x5C: {1, 0, 0, 0},    // list[0][1] length
				0x60: {0x72, 0, 0, 0}, // list[0][2] pointer
				0x64: {1, 0, 0, 0},    // list[0][2] length
				0x70: {'a', 'b', 'c'}, // list[0][*] data: "abc"
				0xA0: {0xC0, 0, 0, 0}, // list[1][0] pointer
				0xA4: {1, 0, 0, 0},    // list[1][0] length
				0xA8: {0xC1, 0, 0, 0}, // list[1][1] pointer
				0xAC: {1, 0, 0, 0},    // list[1][1] length
				0xC0: {'d', 'e'},      // list[1][*] data: "de"
			},
			ptr: 0,
			expected: [][]string{
				{"a", "b", "c"},
				{"d", "e"},
			},
			shouldError: false,
		},
		{
			name:        "nil result pointer",
			memory:      map[uint32][]uint8{0: {0, 0, 0, 0}},
			ptr:         0,
			shouldError: true,
		},
		{
			name:        "invalid list pointer",
			memory:      map[uint32][]uint8{},
			ptr:         0,
			expected:    [][]string{},
			shouldError: true,
		},
		{
			name:        "invalid list length",
			memory:      map[uint32][]uint8{0: {0, 0, 0, 0}},
			ptr:         0,
			expected:    [][]string{},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := abi.AbiOptions{
				StringEncoding: abi.StringEncodingUTF8,
				Memory:         createMemoryFromMap(tt.memory),
			}
			var result [][]string
			err := abi.Read(opts, tt.ptr, &result)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
