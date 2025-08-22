package test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
)

func TestAlignTo(t *testing.T) {
	tests := []struct {
		name      string
		ptr       uint64
		alignment uint64
		expected  uint64
	}{
		{
			name:      "No alignment needed",
			ptr:       16,
			alignment: 4,
			expected:  16,
		},
		{
			name:      "Align to next multiple",
			ptr:       17,
			alignment: 4,
			expected:  20,
		},
		{
			name:      "Already aligned",
			ptr:       32,
			alignment: 8,
			expected:  32,
		},
		{
			name:      "Align to larger alignment",
			ptr:       33,
			alignment: 8,
			expected:  40,
		},
		{
			name:      "Zero alignment",
			ptr:       10,
			alignment: 0,
			expected:  10,
		},
		{
			name:      "Alignment of 1",
			ptr:       25,
			alignment: 1,
			expected:  25,
		},
		{
			name:      "Align to larger",
			ptr:       1,
			alignment: 4,
			expected:  4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := abi.AlignTo(tt.ptr, tt.alignment)
			if result != tt.expected {
				t.Errorf("AlignTo(%d, %d) = %d; want %d", tt.ptr, tt.alignment, result, tt.expected)
			}
		})
	}
}
