package abi_test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Define a synthetic enum type mirroring codegen pattern (suffix Enum)
type SampleEnum uint8

const (
	SampleEnumAlpha SampleEnum = 2
	SampleEnumBeta  SampleEnum = 3
	SampleEnumGamma SampleEnum = 4
)

func TestReadEnum(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(map[uint64][]byte{
		0x00: {0x04}, // gamma
	})
	var v SampleEnum
	err := abi.Read(opts, 0x00, &v)
	require.NoError(t, err)
	assert.Equal(t, SampleEnumGamma, v)
}

func TestWriteEnum(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	val := SampleEnumBeta
	ptr, free, err := abi.Write(opts, &val, nil)
	require.NoError(t, err)
	defer free()
	b, ok := opts.Memory.Read(ptr, 1)
	require.True(t, ok)
	assert.Equal(t, []byte{0x03}, b)
}

func TestWriteParameterEnum(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	val := SampleEnumAlpha
	params, free, err := abi.WriteParameter(opts, &val)
	require.NoError(t, err)
	defer free()
	require.Len(t, params, 1)
	assert.Equal(t, uint64(2), params[0].Value)
}

func TestWriteThenReadEnumRoundTrip(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	original := SampleEnumGamma
	ptr, free, err := abi.Write(opts, &original, nil)
	require.NoError(t, err)
	defer free()
	var decoded SampleEnum
	err = abi.Read(opts, ptr, &decoded)
	require.NoError(t, err)
	assert.Equal(t, original, decoded)
}

// Sample type that is NOT an enum (no Enum suffix) to trigger predicate failures
type PlainNumber uint8

func TestReadEnumErrors(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(map[uint64][]byte{})

	t.Run("nil pointer result via generic Read", func(t *testing.T) {
		var ptr *SampleEnum
		err := abi.Read(opts, 0, ptr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "non-nil pointer")
	})

	t.Run("non-enum pointer passed to ReadEnum", func(t *testing.T) {
		var v PlainNumber
		err := abi.ReadEnum(opts, 0, &v)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "enum pointer")
	})

	t.Run("underlying int read failure on valid enum", func(t *testing.T) {
		var v SampleEnum
		err := abi.ReadEnum(opts, 0x9999, &v)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read")
	})
}

func TestWriteEnumErrors(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)

	t.Run("invalid (nil) enum pointer value", func(t *testing.T) {
		var nilPtr *SampleEnum = nil
		_, _, err := abi.WriteEnum(opts, nilPtr, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "valid enum value")
	})

	t.Run("non-enum value passed to WriteEnum", func(t *testing.T) {
		var v PlainNumber = 5
		_, _, err := abi.WriteEnum(opts, &v, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value must be an enum")
	})
}

func TestWriteParameterEnumErrors(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)

	t.Run("invalid (nil) enum pointer", func(t *testing.T) {
		var nilPtr *SampleEnum
		_, _, err := abi.WriteParameterEnum(opts, nilPtr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "valid enum value")
	})

	t.Run("non-enum value passed to WriteParameterEnum", func(t *testing.T) {
		var v PlainNumber = 7
		_, _, err := abi.WriteParameterEnum(opts, &v)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value must be an enum")
	})
}
