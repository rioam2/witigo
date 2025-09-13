package abi_test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
)

type SizeOfTest1Record struct {
	A int32
	B float64
}

func TestSizeOf(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected uint64
	}{
		{
			name:     "int32",
			input:    int32(0),
			expected: 4,
		},
		{
			name:     "float64",
			input:    float64(0),
			expected: 8,
		},
		{
			name:     "string",
			input:    "test",
			expected: 8, // Size of string header (pointer + length)
		},
		{
			name:     "struct",
			input:    SizeOfTest1Record{},
			expected: 16, // Assuming alignment and padding
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			size := abi.SizeOf(tt.input)
			if size != tt.expected {
				t.Errorf("SizeOf(%T) = %d, want %d", tt.input, size, tt.expected)
			}
		})
	}
}
