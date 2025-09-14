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
