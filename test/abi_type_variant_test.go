package test

import (
	"testing"

	witigo "github.com/rioam2/witigo/pkg"
)

func TestAbiTypeDefinitionVariant(t *testing.T) {
	t.Run("NewAbiTypeDefinitionVariant with valid cases", func(t *testing.T) {
		variant := witigo.NewAbiTypeDefinitionVariant([]witigo.AbiTypeDefinition{
			witigo.NewAbiTypeDefinitionU8(),
			witigo.NewAbiTypeDefinitionU16(),
			witigo.NewAbiTypeDefinitionU32(),
		})
		variantLength := *variant.Properties().Length
		if variantLength != 3 {
			t.Fatalf("Expected 3 case, got %d", variantLength)
		}
	})
}

func TestAbiTypeDefinitionVariantAlignmentAndSize(t *testing.T) {
	testCases := []struct {
		name                string
		caseTypes           []witigo.AbiTypeDefinition
		expectedAlign       int
		expectedSizeInBytes int
		expectedString      string
	}{
		{
			name: "Simple (U8) Variant",
			caseTypes: []witigo.AbiTypeDefinition{
				witigo.NewAbiTypeDefinitionU8(),
			},
			expectedAlign:       1,
			expectedSizeInBytes: 2,
			expectedString:      "variant{u8}",
		},
		{
			name: "Simple (U8+U16) Variant",
			caseTypes: []witigo.AbiTypeDefinition{
				witigo.NewAbiTypeDefinitionU8(),
				witigo.NewAbiTypeDefinitionU16(),
			},
			expectedAlign:       2,
			expectedSizeInBytes: 4,
			expectedString:      "variant{u8; u16}",
		},
		{
			name: "Simple (U8+U32) Variant",
			caseTypes: []witigo.AbiTypeDefinition{
				witigo.NewAbiTypeDefinitionU8(),
				witigo.NewAbiTypeDefinitionU32(),
			},
			expectedAlign:       4,
			expectedSizeInBytes: 8,
			expectedString:      "variant{u8; u32}",
		},
		{
			name: "(U8+U32+S16+Variant) Variant",
			caseTypes: []witigo.AbiTypeDefinition{
				witigo.NewAbiTypeDefinitionU8(),
				witigo.NewAbiTypeDefinitionU32(),
				witigo.NewAbiTypeDefinitionS16(),
				witigo.NewAbiTypeDefinitionVariant([]witigo.AbiTypeDefinition{
					witigo.NewAbiTypeDefinitionU8(),
					witigo.NewAbiTypeDefinitionU16(),
					witigo.NewAbiTypeDefinitionU32(),
				}),
			},
			expectedAlign:       4,
			expectedSizeInBytes: 12,
			expectedString:      "variant{u8; u32; s16; variant{u8; u16; u32}}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			enumType := witigo.NewAbiTypeDefinitionVariant(tc.caseTypes)
			alignment := enumType.Alignment()
			if alignment != tc.expectedAlign {
				t.Fatalf("Expected alignment of %d, got %d", tc.expectedAlign, alignment)
			}
			sizeInBytes := enumType.SizeInBytes()
			if sizeInBytes != tc.expectedSizeInBytes {
				t.Fatalf("Expected size in bytes of %d, got %d", tc.expectedSizeInBytes, sizeInBytes)
			}
			if enumType.String() != tc.expectedString {
				t.Fatalf("Expected string representation of '%s', got '%s'", tc.expectedString, enumType.String())
			}
		})
	}
}
