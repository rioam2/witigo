package main_test

import (
	"context"
	"testing"

	all_types_example_component "github.com/rioam2/witigo/examples/all-types/generated"
	"github.com/stretchr/testify/assert"
)

const createInstanceErrFmt = "Failed to create instance: %v"

func TestEnum(t *testing.T) {
	instance, err := all_types_example_component.New(context.Background())
	if err != nil {
		t.Fatalf(createInstanceErrFmt, err)
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
		t.Fatalf(createInstanceErrFmt, err)
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
			name:     "Test with None variant type",
			input:    all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeNone},
			expected: all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeNone},
		},
		{
			name:     "Test with Restricted variant type",
			input:    all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeRestricted, Restricted: []string{"123-456-7890"}},
			expected: all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeRestricted, Restricted: []string{"123-456-7890 - modified by C++"}},
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

func TestComplexVariant(t *testing.T) {
	instance, err := all_types_example_component.New(context.Background())
	if err != nil {
		t.Fatalf(createInstanceErrFmt, err)
	}

	// Table covers all cases of ComplexUnionVariant transformation rules
	cases := []struct {
		name     string
		input    all_types_example_component.ComplexUnionVariant
		expected all_types_example_component.ComplexUnionVariant
	}{
		{
			name: "empty->empty",
			input: all_types_example_component.ComplexUnionVariant{
				Type: all_types_example_component.ComplexUnionVariantTypeEmpty,
			},
			expected: all_types_example_component.ComplexUnionVariant{
				Type: all_types_example_component.ComplexUnionVariantTypeEmpty,
			},
		},
		{
			name: "number inc",
			input: all_types_example_component.ComplexUnionVariant{
				Type:   all_types_example_component.ComplexUnionVariantTypeNumber,
				Number: 41,
			},
			expected: all_types_example_component.ComplexUnionVariant{
				Type:   all_types_example_component.ComplexUnionVariantTypeNumber,
				Number: 42,
			},
		},
		{
			name: "float double",
			input: all_types_example_component.ComplexUnionVariant{
				Type:     all_types_example_component.ComplexUnionVariantTypeFloating,
				Floating: 3.5,
			},
			expected: all_types_example_component.ComplexUnionVariant{
				Type:     all_types_example_component.ComplexUnionVariantTypeFloating,
				Floating: 7.0,
			},
		},
		{
			name: "big decrement",
			input: all_types_example_component.ComplexUnionVariant{
				Type: all_types_example_component.ComplexUnionVariantTypeBig,
				Big:  100,
			},
			expected: all_types_example_component.ComplexUnionVariant{
				Type: all_types_example_component.ComplexUnionVariantTypeBig,
				Big:  99,
			},
		},
		{
			name: "text append !",
			input: all_types_example_component.ComplexUnionVariant{
				Type: all_types_example_component.ComplexUnionVariantTypeText,
				Text: "hi",
			},
			expected: all_types_example_component.ComplexUnionVariant{
				Type: all_types_example_component.ComplexUnionVariantTypeText,
				Text: "hi!",
			},
		},
		{
			name: "bytes append ff",
			input: all_types_example_component.ComplexUnionVariant{
				Type:  all_types_example_component.ComplexUnionVariantTypeBytes,
				Bytes: []uint8{0x01, 0x02},
			},
			expected: all_types_example_component.ComplexUnionVariant{
				Type:  all_types_example_component.ComplexUnionVariantTypeBytes,
				Bytes: []uint8{0x01, 0x02, 0xff},
			},
		},
		{
			name: "pair transform",
			input: all_types_example_component.ComplexUnionVariant{
				Type: all_types_example_component.ComplexUnionVariantTypePair,
				Pair: all_types_example_component.SmallRecordRecord{X: 5, Y: 10},
			},
			expected: all_types_example_component.ComplexUnionVariant{
				Type: all_types_example_component.ComplexUnionVariantTypePair,
				Pair: all_types_example_component.SmallRecordRecord{X: 4, Y: 11},
			},
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			out, err := instance.ComplexVariantFunc(c.input)
			assert.NoError(t, err)
			assert.Equal(t, c.expected, out)
		})
	}
}
