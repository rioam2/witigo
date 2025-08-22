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
		return errors.New("must pass a non-nil int/uint pointer result")
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

func WriteInt(opts AbiOptions, value any, ptrHint *uint32) (ptr uint32, free AbiFreeCallback, err error) {
	// Initialize return values
	ptr = 0
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return ptr, free, errors.New("must pass a valid int/uint pointer value")
	}

	// Validate that the value is an integer or unsigned integer type
	if !rv.CanUint() && !rv.CanInt() {
		return ptr, free, fmt.Errorf("cannot write int/uint from: %s", rv.Kind())
	}

	// Extract ABI properties of intrinsic type
	size := SizeOf(value)
	alignment := AlignmentOf(value)

	// Allocate memory if ptrHint is not provided or is zero
	if ptrHint != nil && *ptrHint != 0 {
		ptr = AlignTo(*ptrHint, alignment)
	} else {
		var freeInt AbiFreeCallback
		ptr, freeInt, err = abi_malloc(opts, size, alignment)
		if err != nil {
			return ptr, free, err
		}
		freeCallbacks = append(freeCallbacks, freeInt)
	}

	// Prepare bytes to write
	bytes := make([]byte, size)
	if rv.CanUint() {
		uint64Value := rv.Uint()
		for i := range size {
			bytes[i] = byte(uint64Value >> (8 * i))
		}
	} else if rv.CanInt() {
		int64Value := rv.Int()
		for i := range size {
			bytes[i] = byte(int64Value >> (8 * i))
		}
	}

	// Write bytes to memory
	if !opts.Memory.Write(ptr, bytes) {
		return ptr, free, fmt.Errorf("failed to write %d bytes at int/uint pointer %d", size, ptr)
	}

	return ptr, free, nil
}

func WriteParameterInt(opts AbiOptions, value any) (args []uint32, free AbiFreeCallback, err error) {
	// Initialize return values
	args = []uint32{}
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return args, free, errors.New("must pass a valid int/uint pointer value")
	}

	// Validate that the value is an integer or unsigned integer type
	if !rv.CanUint() && !rv.CanInt() {
		return args, free, fmt.Errorf("cannot write int/uint from: %s", rv.Kind())
	}

	if rv.CanUint() {
		// TODO: Fix trucation of byte size
		args = append(args, uint32(rv.Uint()))
	} else if rv.CanInt() {
		// TODO: Fix trucation of byte size
		args = append(args, uint32(rv.Int()))
	}
	return args, free, nil
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

func WriteBool(opts AbiOptions, value any, ptrHint *uint32) (ptr uint32, free AbiFreeCallback, err error) {
	// Initialize return values
	ptr = 0
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return ptr, free, errors.New("must pass a valid boolean value")
	}

	// Validate that the value is a boolean type
	if rv.Kind() != reflect.Bool {
		return ptr, free, fmt.Errorf("cannot write bool from: %s", rv.Kind())
	}

	// Extract ABI properties of intrinsic type
	size := SizeOf(value)
	alignment := AlignmentOf(value)

	// Allocate memory if ptrHint is not provided or is zero
	if ptrHint != nil && *ptrHint != 0 {
		ptr = AlignTo(*ptrHint, alignment)
	} else {
		var freeBool AbiFreeCallback
		ptr, freeBool, err = abi_malloc(opts, size, alignment)
		if err != nil {
			return ptr, free, err
		}
		freeCallbacks = append(freeCallbacks, freeBool)
	}

	// Prepare bytes to write
	bytes := []byte{0}
	if rv.Bool() {
		bytes[0] = 1
	}

	// Write bytes to memory
	if !opts.Memory.Write(ptr, bytes) {
		return ptr, free, fmt.Errorf("failed to write %d bytes at boolean pointer %d", size, ptr)
	}

	return ptr, free, nil
}

func WriteParameterBool(opts AbiOptions, value any) (args []uint32, free AbiFreeCallback, err error) {
	// Initialize return values
	args = []uint32{}
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return args, free, errors.New("must pass a valid boolean pointer value")
	}

	// Validate that the value is a boolean type
	if rv.Kind() != reflect.Bool {
		return args, free, fmt.Errorf("cannot write bool from: %s", rv.Kind())
	}

	// Append the boolean value as an argument
	if rv.Bool() {
		args = append(args, 1)
	} else {
		args = append(args, 0)
	}
	return args, free, nil
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

func WriteFloat(opts AbiOptions, value any, ptrHint *uint32) (ptr uint32, free AbiFreeCallback, err error) {
	// Initialize return values
	ptr = 0
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return ptr, free, errors.New("must pass a valid float value")
	}

	// Validate that the value is a float type
	if rv.Kind() != reflect.Float32 && rv.Kind() != reflect.Float64 {
		return ptr, free, fmt.Errorf("cannot write float from: %s", rv.Kind())
	}

	// Extract ABI properties of intrinsic type
	size := SizeOf(value)
	alignment := AlignmentOf(value)

	// Allocate memory if ptrHint is not provided or is zero
	if ptrHint != nil && *ptrHint != 0 {
		ptr = AlignTo(*ptrHint, alignment)
	} else {
		var freeFloat AbiFreeCallback
		ptr, freeFloat, err = abi_malloc(opts, size, alignment)
		if err != nil {
			return ptr, free, err
		}
		freeCallbacks = append(freeCallbacks, freeFloat)
	}

	// Prepare bytes to write
	floatBytes := make([]byte, size)
	if rv.Kind() == reflect.Float32 {
		binary.LittleEndian.PutUint32(floatBytes, math.Float32bits(float32(rv.Float())))
	} else if rv.Kind() == reflect.Float64 {
		binary.LittleEndian.PutUint64(floatBytes, math.Float64bits(rv.Float()))
	}

	// Write bytes to memory
	if !opts.Memory.Write(ptr, floatBytes) {
		return ptr, free, fmt.Errorf("failed to write %d bytes at float pointer %d", size, ptr)
	}

	return ptr, free, nil
}

func WriteParameterFloat(opts AbiOptions, value any) (args []uint32, free AbiFreeCallback, err error) {
	// Initialize return values
	args = []uint32{}
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return args, free, errors.New("must pass a valid float pointer value")
	}

	// Validate that the value is a float type
	if rv.Kind() != reflect.Float32 && rv.Kind() != reflect.Float64 {
		return args, free, fmt.Errorf("cannot write float from: %s", rv.Kind())
	}

	// Append the float value as an argument
	if rv.Kind() == reflect.Float32 {
		args = append(args, math.Float32bits(float32(rv.Float())))
	} else if rv.Kind() == reflect.Float64 {
		// TODO: Fix trucation of byte size
		args = append(args, uint32(math.Float64bits(rv.Float())))
	}

	return args, free, nil
}
