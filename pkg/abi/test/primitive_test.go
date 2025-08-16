package test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
	"github.com/stretchr/testify/assert"
)

func TestReadInt(t *testing.T) {
	testCases := []struct {
		name        string
		ptr         uint32
		memoryData  map[uint32][]byte
		valueType   string
		expected    interface{}
		expectError bool
	}{
		{
			name: "read int8",
			ptr:  0x100,
			memoryData: map[uint32][]byte{
				0x100: {0x7F},
			},
			valueType: "int8",
			expected:  int8(127),
		},
		{
			name: "read uint8",
			ptr:  0x100,
			memoryData: map[uint32][]byte{
				0x100: {0xFF},
			},
			valueType: "uint8",
			expected:  uint8(255),
		},
		{
			name: "read int16",
			ptr:  0x200,
			memoryData: map[uint32][]byte{
				0x200: {0x34, 0x12},
			},
			valueType: "int16",
			expected:  int16(0x1234),
		},
		{
			name: "read uint16",
			ptr:  0x200,
			memoryData: map[uint32][]byte{
				0x200: {0x34, 0x12},
			},
			valueType: "uint16",
			expected:  uint16(0x1234),
		},
		{
			name: "read int32",
			ptr:  0x300,
			memoryData: map[uint32][]byte{
				0x300: {0x78, 0x56, 0x34, 0x12},
			},
			valueType: "int32",
			expected:  int32(0x12345678),
		},
		{
			name: "read uint32",
			ptr:  0x300,
			memoryData: map[uint32][]byte{
				0x300: {0x78, 0x56, 0x34, 0x12},
			},
			valueType: "uint32",
			expected:  uint32(0x12345678),
		},
		{
			name: "read int64",
			ptr:  0x400,
			memoryData: map[uint32][]byte{
				0x400: {0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01},
			},
			valueType: "int64",
			expected:  int64(0x0123456789ABCDEF),
		},
		{
			name: "read uint64",
			ptr:  0x400,
			memoryData: map[uint32][]byte{
				0x400: {0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01},
			},
			valueType: "uint64",
			expected:  uint64(0x0123456789ABCDEF),
		},
		{
			name: "invalid memory address",
			ptr:  0x500,
			memoryData: map[uint32][]byte{
				0x100: {0x01},
			},
			valueType:   "int8",
			expectError: true,
		},
		{
			name: "insufficient memory size",
			ptr:  0x200,
			memoryData: map[uint32][]byte{
				0x200: {0x01}, // Only 1 byte, but int16 needs 2
			},
			valueType:   "int16",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			memory := createMemoryFromMap(tc.memoryData)
			opts := abi.AbiOptions{Memory: memory}

			switch tc.valueType {
			case "int8":
				var result int8
				err := abi.Read(opts, tc.ptr, &result)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				}
			case "uint8":
				var result uint8
				err := abi.Read(opts, tc.ptr, &result)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				}
			case "int16":
				var result int16
				err := abi.Read(opts, tc.ptr, &result)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				}
			case "uint16":
				var result uint16
				err := abi.Read(opts, tc.ptr, &result)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				}
			case "int32":
				var result int32
				err := abi.Read(opts, tc.ptr, &result)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				}
			case "uint32":
				var result uint32
				err := abi.Read(opts, tc.ptr, &result)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				}
			case "int64":
				var result int64
				err := abi.Read(opts, tc.ptr, &result)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				}
			case "uint64":
				var result uint64
				err := abi.Read(opts, tc.ptr, &result)
				if tc.expectError {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tc.expected, result)
				}
			}
		})
	}

	// Test nil pointer case
	t.Run("nil result pointer", func(t *testing.T) {
		memory := createMemoryFromMap(map[uint32][]byte{0x100: {0x01}})
		opts := abi.AbiOptions{Memory: memory}
		err := abi.Read(opts, 0x100, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must pass a non-nil pointer result")
	})
}
