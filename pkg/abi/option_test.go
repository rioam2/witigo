package abi_test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadOption_NilPointer(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	err := abi.ReadOption(opts, 0, nil)
	assert.EqualError(t, err, "must pass a non-nil pointer result")
}

func TestReadOption_NonPointerResult(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	var result int
	err := abi.ReadOption(opts, 0, result)
	assert.EqualError(t, err, "must pass a non-nil pointer result")
}

func TestReadOption_InvalidType(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	var result int
	err := abi.ReadOption(opts, 0, &result)
	assert.EqualError(t, err, "expected Option type, got int")
}

func TestReadOption_None(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(map[uint64][]byte{
		0x00: {0x00},
	})
	var result abi.Option[int]
	err := abi.ReadOption(opts, 0, &result)
	require.NoError(t, err)
}

func TestReadOption_Some(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(map[uint64][]byte{
		0x00: {0x01},
		0x04: {0x42, 0x00, 0x00, 0x00},
	})
	var result abi.Option[int]
	err := abi.ReadOption(opts, 0, &result)
	require.NoError(t, err)
	assert.NotNil(t, result.Value)
	assert.Equal(t, 0x42, result.Value)
}
func TestWriteOption_NilPointer(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	_, _, err := abi.WriteOption(opts, nil, nil)
	assert.EqualError(t, err, "must pass a non-nil pointer value")
}

func TestWriteOption_NonPointerValue(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	var value int
	_, _, err := abi.WriteOption(opts, value, nil)
	assert.EqualError(t, err, "must pass a non-nil pointer value")
}

func TestWriteOption_InvalidType(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	var value int
	_, _, err := abi.WriteOption(opts, &value, nil)
	assert.EqualError(t, err, "expected Option type, got int")
}

func TestWriteOption_None(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	value := abi.Option[int]{IsSome: false}
	ptr, free, err := abi.WriteOption(opts, &value, nil)
	require.NoError(t, err)
	defer free()
	discriminantValue, ok := opts.Memory.Read(ptr, 1)
	assert.True(t, ok)
	assert.Equal(t, []byte{0x00}, discriminantValue)
}

func TestWriteOption_Some(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	value := abi.Option[int]{IsSome: true, Value: 0x42}
	ptr, free, err := abi.WriteOption(opts, &value, nil)
	require.NoError(t, err)
	defer free()
	discriminantValue, ok := opts.Memory.Read(ptr, 1)
	assert.True(t, ok)
	assert.Equal(t, []byte{0x01}, discriminantValue)
	intValue, ok := opts.Memory.Read(ptr+4, 4)
	assert.True(t, ok)
	assert.Equal(t, []byte{0x42, 0x00, 0x00, 0x00}, intValue)
}

func TestWriteParameterOption_NilPointer(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	_, _, err := abi.WriteParameterOption(opts, nil)
	assert.EqualError(t, err, "must pass a non-nil pointer value")
}

func TestWriteParameterOption_NonPointerValue(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	var value int
	_, _, err := abi.WriteParameterOption(opts, value)
	assert.EqualError(t, err, "must pass a non-nil pointer value")
}

func TestWriteParameterOption_InvalidType(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	var value int
	_, _, err := abi.WriteParameterOption(opts, &value)
	assert.EqualError(t, err, "expected Option type, got int")
}

func TestWriteParameterOption_None(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	value := abi.Option[int]{IsSome: false}
	args, free, err := abi.WriteParameterOption(opts, &value)
	require.NoError(t, err)
	defer free()
	assert.Equal(t, []abi.Parameter{
		{Value: uint64(0), Size: 0x01, Alignment: 0x04}, // Discriminant
		{Value: uint64(0), Size: 0x04, Alignment: 0x04}, // Value
	}, args)
}

func TestWriteParameterOption_Some(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	value := abi.Option[int]{IsSome: true, Value: 0x42}
	args, free, err := abi.WriteParameterOption(opts, &value)
	require.NoError(t, err)
	defer free()
	assert.Equal(t, 2, len(args))
	assert.Equal(t, abi.Parameter{Value: 0x01, Size: 0x01, Alignment: 0x04}, args[0]) // Discriminant
	assert.Equal(t, abi.Parameter{Value: 0x42, Size: 0x04, Alignment: 0x04}, args[1]) // Value
}

func TestWriteThenReadOptionNone(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	value := abi.Option[int]{IsSome: false}
	ptr, free, err := abi.WriteOption(opts, &value, nil)
	require.NoError(t, err)
	defer free()

	var result abi.Option[int]
	err = abi.ReadOption(opts, ptr, &result)
	require.NoError(t, err)
	assert.False(t, result.IsSome)
}

func TestWriteThenReadOptionSome(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	value := abi.Option[int]{IsSome: true, Value: 0x42}
	ptr, free, err := abi.WriteOption(opts, &value, nil)
	require.NoError(t, err)
	defer free()

	var result abi.Option[int]
	err = abi.ReadOption(opts, ptr, &result)
	require.NoError(t, err)
	assert.True(t, result.IsSome)
	assert.Equal(t, 0x42, result.Value)
}

func TestWriteThenReadOptionSomeSlice(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	value := abi.Option[[]byte]{IsSome: true, Value: []byte{0x01, 0x02, 0x03}}
	ptr, free, err := abi.WriteOption(opts, &value, nil)
	require.NoError(t, err)
	defer free()

	var result abi.Option[[]byte]
	err = abi.ReadOption(opts, ptr, &result)
	require.NoError(t, err)
	assert.True(t, result.IsSome)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, result.Value)
}
