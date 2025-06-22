package test

import (
	"testing"

	witigo "github.com/rioam2/witigo/pkg"
)

func TestAbiTypeDefinitionEnum(t *testing.T) {
	t.Run("NewAbiTypeDefinitionEnum with valid cases", func(t *testing.T) {
		enumType := witigo.NewAbiTypeDefinitionEnum(3)
		if _, ok := enumType.(witigo.AbiTypeDefinitionVariant); !ok {
			t.Fatalf("Expected AbiTypeDefinitionVariant, got %T", enumType)
		}
		variant := enumType.(witigo.AbiTypeDefinitionVariant)
		variantLength := *variant.Properties().Length
		if variantLength != 3 {
			t.Fatalf("Expected 3 cases, got %d", variantLength)
		}
		variantString := "variant{none; none; none}"
		if variant.String() != variantString {
			t.Fatalf("Expected variant string '%s', got '%s'", variantString, variant.String())
		}
	})
}

func TestAbiTypeDefinitionEnumAlignmentAndSize(t *testing.T) {
	testCases := []struct {
		name                string
		size                int
		expectedAlign       int
		expectedSizeInBytes int
	}{
		{
			name:                "Enum of size 4",
			size:                4,
			expectedAlign:       1,
			expectedSizeInBytes: 1,
		},
		{
			name:                "Enum of size 16",
			size:                16,
			expectedAlign:       2,
			expectedSizeInBytes: 2,
		},
		{
			name:                "Enum of size 32",
			size:                32,
			expectedAlign:       4,
			expectedSizeInBytes: 4,
		},
		{
			name:                "Enum of size 1024",
			size:                1024,
			expectedAlign:       4,
			expectedSizeInBytes: 4,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enumType := witigo.NewAbiTypeDefinitionEnum(tc.size)
			alignment := enumType.Alignment()
			if alignment != tc.expectedAlign {
				t.Fatalf("Expected alignment of %d, got %d", tc.expectedAlign, alignment)
			}
			sizeInBytes := enumType.SizeInBytes()
			if sizeInBytes != tc.expectedSizeInBytes {
				t.Fatalf("Expected size in bytes of %d, got %d", tc.expectedSizeInBytes, sizeInBytes)
			}
		})
	}
}
