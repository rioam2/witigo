package main

import (
	"context"
	"testing"

	basic_example_component "github.com/rioam2/witigo/examples/basic/generated"
)

func TestDoubleOperation(t *testing.T) {
	instance, err := basic_example_component.New(context.Background())
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}

	tests := []struct {
		name     string
		input    basic_example_component.DoubleOperationRecord
		expected basic_example_component.DoubleResultRecord
	}{
		{
			name: "Basic test with positive numbers",
			input: basic_example_component.DoubleOperationRecord{
				DoubleList:   []float64{1.1, 2.2, 3.3},
				DoubleString: "Hello ",
			},
			expected: basic_example_component.DoubleResultRecord{
				DoubledList:   []float64{2.2, 4.4, 6.6},
				DoubledString: "Hello Hello ",
			},
		},
		{
			name: "Test with empty input",
			input: basic_example_component.DoubleOperationRecord{
				DoubleList:   []float64{},
				DoubleString: "",
			},
			expected: basic_example_component.DoubleResultRecord{
				DoubledList:   []float64{},
				DoubledString: "",
			},
		},
		{
			name: "Test with negative numbers",
			input: basic_example_component.DoubleOperationRecord{
				DoubleList:   []float64{-1.5, -2.5},
				DoubleString: "Negative ",
			},
			expected: basic_example_component.DoubleResultRecord{
				DoubledList:   []float64{-3.0, -5.0},
				DoubledString: "Negative Negative ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := instance.Double(tt.input)
			if err != nil {
				t.Fatalf("Double operation failed: %v", err)
			}

			if len(result.DoubledList) != len(tt.expected.DoubledList) {
				t.Errorf("Expected DoubleList length %d, got %d", len(tt.expected.DoubledList), len(result.DoubledList))
			}

			for i, v := range result.DoubledList {
				if v != tt.expected.DoubledList[i] {
					t.Errorf("Expected DoubleList[%d] = %f, got %f", i, tt.expected.DoubledList[i], v)
				}
			}

			if result.DoubledString != tt.expected.DoubledString {
				t.Errorf("Expected DoubleString = %q, got %q", tt.expected.DoubledString, result.DoubledString)
			}
		})
	}
}
