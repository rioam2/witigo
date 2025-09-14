package abi_test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Synthetic variant mirroring code generation pattern
type SampleVariantType int

const (
	SampleVariantTypeA SampleVariantType = 0
	SampleVariantTypeB SampleVariantType = 1
	SampleVariantTypeC SampleVariantType = 2
)

type SampleVariant struct {
	Type SampleVariantType
	A    struct{}
	B    uint32
	C    string
}

func TestWriteParameterVariant(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	v := SampleVariant{Type: SampleVariantTypeB, B: 0xDEADBEEF}
	params, free, err := abi.WriteParameters(opts, v)
	require.NoError(t, err)
	defer free()
	// Flattened shape should be: discriminant + 2 slots (string case pointer,len) => 3 params
	require.Len(t, params, 3)
	assert.Equal(t, uint64(SampleVariantTypeB), params[0])
	assert.Equal(t, uint64(0xDEADBEEF), params[1]) // B payload in first payload slot
	// Third slot unused for this case => zero
	assert.Equal(t, uint64(0), params[2])
}

func TestVariantRoundTripEmptyCase(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	original := SampleVariant{Type: SampleVariantTypeA}
	ptr, free, err := abi.Write(opts, original, nil)
	require.NoError(t, err)
	defer free()
	var decoded SampleVariant
	err = abi.Read(opts, ptr, &decoded)
	require.NoError(t, err)
	assert.Equal(t, original.Type, decoded.Type)
}

func TestVariantRoundTripPayloadCases(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	cases := []SampleVariant{
		{Type: SampleVariantTypeB, B: 42},
		{Type: SampleVariantTypeC, C: "hello"},
		{Type: SampleVariantTypeC, C: ""},                                                     // empty string case
		{Type: SampleVariantTypeB, B: 0},                                                      // zero value case
		{Type: SampleVariantTypeB, B: 0xFFFFFFFF},                                             // max uint32
		{Type: SampleVariantTypeC, C: "A longer string to test string handling in variants."}, // long string
		{Type: SampleVariantTypeA},                                                            // empty case again
	}
	for _, c := range cases {
		ptr, free, err := abi.Write(opts, c, nil)
		require.NoError(t, err)
		var decoded SampleVariant
		err = abi.Read(opts, ptr, &decoded)
		require.NoError(t, err)
		assert.Equal(t, c.Type, decoded.Type)
		switch c.Type {
		case SampleVariantTypeB:
			assert.Equal(t, c.B, decoded.B)
		case SampleVariantTypeC:
			assert.Equal(t, c.C, decoded.C)
		}
		free()
	}
}

func TestVariantErrors(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	// Invalid discriminant when reading
	// Write bytes manually with out-of-range discriminant (e.g., 5)
	opts.Memory.Write(0, []byte{5, 0, 0, 0, 0, 0, 0, 0}) // discriminant stored as 8-byte int
	var v SampleVariant
	err := abi.ReadVariant(opts, 0, &v)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "out of range")

	// Writing with invalid case index (manually craft struct) not possible via type-safe API,
	// but we can simulate by setting discriminant > cases and calling WriteVariant.
	bad := SampleVariant{Type: SampleVariantType(99)}
	_, _, err = abi.WriteVariant(opts, bad, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "out of range")
}
