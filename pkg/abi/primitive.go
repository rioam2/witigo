package abi

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
)

// Read reads a value from linear memory at the specified pointer into the result.
func ReadInt(opts AbiOptions, ptr uint32, result any) error {
	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()

	// Validate that the result is an integer or unsigned integer type
	if !rv.CanUint() && !rv.CanInt() {
		return fmt.Errorf("cannot read int/uint into: %s", rv.Kind())
	}

	// Extract ABI properties of intrinsic type
	size := SizeOf(result)
	alignment := AlignmentOf(result)
	ptr = AlignTo(ptr, alignment)

	// Read the bytes from memory
	bytes, ok := opts.Memory.Read(ptr, size)
	if !ok {
		return fmt.Errorf("failed to read %d bytes at pointer %d", size, ptr)
	}

	// Convert bytes to the appropriate type
	if rv.CanUint() {
		uint64Value := uint64(0)
		for i, b := range bytes {
			uint64Value |= (uint64(b) << (8 * i))
		}
		rv.SetUint(uint64Value)
	} else if rv.CanInt() {
		int64Value := int64(0)
		for i, b := range bytes {
			int64Value |= (int64(b) << (8 * i))
		}
		rv.SetInt(int64Value)
	}

	return nil
}

func ReadBool(opts AbiOptions, ptr uint32, result any) error {
	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()

	// Validate that the result is a boolean type
	if rv.Kind() != reflect.Bool {
		return fmt.Errorf("cannot read bool into: %s", rv.Kind())
	}

	// Read a single byte from memory
	bytes, ok := opts.Memory.Read(ptr, 1)
	if !ok {
		return fmt.Errorf("failed to read boolean at pointer %d", ptr)
	}

	// Set the boolean value based on the byte read
	rv.SetBool(bytes[0] != 0)

	return nil
}

var CANONICAL_FLOAT32_NAN = []byte{0x7f, 0xc0, 0x00, 0x00}
var CANONICAL_FLOAT64_NAN = []byte{0x7f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

func ReadFloat(opts AbiOptions, ptr uint32, result any) error {
	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()

	// Validate that the result is a float type
	if rv.Kind() != reflect.Float32 && rv.Kind() != reflect.Float64 {
		return fmt.Errorf("cannot read float into: %s", rv.Kind())
	}

	// Extract ABI properties of intrinsic type
	size := SizeOf(result)
	alignment := AlignmentOf(result)
	ptr = AlignTo(ptr, alignment)

	// Read the floatBytes from memory
	floatBytes, ok := opts.Memory.Read(ptr, size)
	if !ok {
		return fmt.Errorf("failed to read %d bytes at pointer %d", size, ptr)
	}

	// Handle NaN values for float types
	if rv.Kind() == reflect.Float32 && bytes.Equal(floatBytes, CANONICAL_FLOAT32_NAN) ||
		rv.Kind() == reflect.Float64 && bytes.Equal(floatBytes, CANONICAL_FLOAT64_NAN) {
		rv.SetFloat(math.NaN())
		return nil
	}

	// Convert bytes to the appropriate float type
	if rv.Kind() == reflect.Float32 {
		var value float32
		bytesDecoded, err := binary.Decode(floatBytes, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("failed to decode float32: %w", err)
		}
		if bytesDecoded != int(size) {
			return fmt.Errorf("byte size mismatch: expected %d bytes, got %d bytes", size, bytesDecoded)
		}
		rv.SetFloat(float64(value))
	} else if rv.Kind() == reflect.Float64 {
		var value float64
		bytesDecoded, err := binary.Decode(floatBytes, binary.LittleEndian, &value)
		if err != nil {
			return fmt.Errorf("failed to decode float64: %w", err)
		}
		if bytesDecoded != int(size) {
			return fmt.Errorf("byte size mismatch: expected %d bytes, got %d bytes", size, bytesDecoded)
		}
		rv.SetFloat(value)
	}

	return nil
}
