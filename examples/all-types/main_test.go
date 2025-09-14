package main_test

import (
	"context"
	"testing"

	all_types_example_component "github.com/rioam2/witigo/examples/all-types/generated"
)

func TestEnum(t *testing.T) {
	instance, err := all_types_example_component.New(context.Background())
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}

	tests := []struct {
		name     string
		input    all_types_example_component.ColorEnum
		expected all_types_example_component.ColorEnum
	}{
		{
			name:     "Test with NavyBlue enum value",
			input:    all_types_example_component.ColorEnumNavyBlue,
			expected: all_types_example_component.ColorEnumNavyBlue,
		},
		{
			name:     "Test with LimeGreen enum value",
			input:    all_types_example_component.ColorEnumLimeGreen,
			expected: all_types_example_component.ColorEnumLimeGreen,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := instance.EnumFunc(tt.input)
			if err != nil {
				t.Fatalf("Enum operation failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
