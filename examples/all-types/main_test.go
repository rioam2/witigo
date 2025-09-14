package main_test

import (
	"context"
	"testing"

	all_types_example_component "github.com/rioam2/witigo/examples/all-types/generated"
	"github.com/stretchr/testify/assert"
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
			assert.NoError(t, err)
			assert.Equal(t, result, tt.expected)
		})
	}
}

func TestVariant(t *testing.T) {
	instance, err := all_types_example_component.New(context.Background())
	if err != nil {
		t.Fatalf("Failed to create instance: %v", err)
	}

	tests := []struct {
		name     string
		input    all_types_example_component.AllowedDestinationsVariant
		expected all_types_example_component.AllowedDestinationsVariant
	}{
		{
			name:     "Test with Any variant type",
			input:    all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeAny},
			expected: all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeAny},
		},
		{
			name:     "Test with Email variant type",
			input:    all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeNone},
			expected: all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeNone},
		},
		{
			name:     "Test with Phone variant type",
			input:    all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeRestricted, Restricted: []string{"123-456-7890"}},
			expected: all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeRestricted, Restricted: []string{"123-456-7890"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := instance.VariantFunc(tt.input)
			assert.NoError(t, err)
			assert.Equal(t, result, tt.expected)
		})
	}
}
