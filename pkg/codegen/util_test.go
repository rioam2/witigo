package codegen

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscriminantSize(t *testing.T) {
	tests := []struct {
		name     string
		cases    int
		expected int
	}{
		{
			name:     "zero cases",
			cases:    0,
			expected: 8,
		},
		{
			name:     "one case",
			cases:    1,
			expected: 8,
		},
		{
			name:     "two cases",
			cases:    2,
			expected: 8,
		},
		{
			name:     "255 cases",
			cases:    255,
			expected: 8,
		},
		{
			name:     "256 cases",
			cases:    256,
			expected: 16,
		},
		{
			name:     "65535 cases",
			cases:    65535,
			expected: 16,
		},
		{
			name:     "65536 cases",
			cases:    65536,
			expected: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := discriminantSize(tt.cases)
			assert.Equal(t, tt.expected, result)
		})
	}
}
