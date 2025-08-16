package abi

import (
	"errors"
	"fmt"
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
