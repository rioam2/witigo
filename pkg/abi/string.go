package abi

import (
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/text/encoding/unicode"
)

// ReadString reads a string from linear memory at the specified pointer into the result.
func ReadString(opts AbiOptions, ptr uint64, result any) error {
	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()

	// Validate that the result is a string pointer
	if rv.Kind() != reflect.String {
		return errors.New("result must be a string pointer")
	}

	// Ensure the result is settable
	if !rv.CanSet() {
		return errors.New("result must be a settable pointer")
	}

	// Extract ABI properties of intrinsic type
	alignment := AlignmentOf(result)
	ptr = AlignTo(ptr, alignment)
	strEncoding := opts.StringEncoding
	strAlignment := strEncoding.Alignment()
	taggedCodeUnitSize := strEncoding.CodeUnitSize()

	// Read location of string data
	strPtr, ok := opts.Memory.ReadUint32Le(ptr)
	if !ok {
		return fmt.Errorf("failed to read string pointer at %d", ptr)
	}

	// Read the number of tagged code units in the string
	taggedCodeUnits, ok := opts.Memory.ReadUint32Le(ptr + 4)
	if !ok {
		return fmt.Errorf("failed to read tagged code units at %d", ptr+4)
	}

	// Validate alignment of string data pointer
	if uint64(strPtr) != AlignTo(uint64(strPtr), uint64(strAlignment)) {
		return fmt.Errorf("string pointer %d is not aligned to %d bytes", strPtr, strAlignment)
	}

	// Validate that the string pointer is within bounds
	strByteLength := uint64(taggedCodeUnits) * taggedCodeUnitSize
	if uint64(strPtr)+strByteLength > opts.Memory.Size() {
		return fmt.Errorf("string pointer %d with length %d exceeds memory size %d", strPtr, strByteLength, opts.Memory.Size())
	}

	// Read the string data from memory
	strData, ok := opts.Memory.Read(uint64(strPtr), strByteLength)
	if !ok {
		return fmt.Errorf("failed to read string data at %d with length %d", strPtr, strByteLength)
	}

	// Convert the string data based on the encoding
	switch strEncoding {
	case StringEncodingUTF8:
		rv.SetString(string(strData))
		return nil
	case StringEncodingUTF16:
		decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
		str, err := decoder.String(string(strData))
		if err != nil {
			return err
		}
		rv.SetString(str)
		return nil
	default:
		return fmt.Errorf("unsupported string encoding: %s", strEncoding)
	}
}

func WriteString(opts AbiOptions, value any, ptrHint *uint64) (ptr uint64, free AbiFreeCallback, err error) {
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
		return ptr, free, errors.New("must pass a valid string pointer value")
	}

	// Validate that the value is a string type
	if rv.Kind() != reflect.String {
		return ptr, free, fmt.Errorf("cannot write string from: %s", rv.Kind())
	}

	// Extract ABI properties of intrinsic type
	size := SizeOf(value)
	alignment := AlignmentOf(value)

	// Allocate memory if ptrHint is not provided or is zero
	if ptrHint != nil && *ptrHint != 0 {
		ptr = AlignTo(*ptrHint, alignment)
	} else {
		var freeString AbiFreeCallback
		ptr, freeString, err = abiMalloc(opts, size, alignment)
		if err != nil {
			return ptr, free, err
		}
		freeCallbacks = append(freeCallbacks, freeString)
	}

	params, strFree, err := WriteParameterString(opts, value)
	freeCallbacks = append(freeCallbacks, strFree)
	if err != nil {
		return ptr, free, err
	}
	strDataPtr := params[0].Value
	strDataLen := params[1].Value

	// Write string descriptor to linear memory
	if ok := opts.Memory.WriteUint32Le(ptr, uint32(strDataPtr)); !ok {
		return ptr, free, fmt.Errorf("failed to write string data pointer at %d", ptr)
	}
	if ok := opts.Memory.WriteUint32Le(ptr+4, uint32(strDataLen)); !ok {
		return ptr, free, fmt.Errorf("failed to write string length at %d", ptr+4)
	}

	return ptr, free, nil
}

func WriteParameterString(opts AbiOptions, value any) (params []Parameter, free AbiFreeCallback, err error) {
	// Initialize return values
	params = []Parameter{}
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return params, free, errors.New("must pass a valid string pointer value")
	}

	// Validate that the value is a string type
	if rv.Kind() != reflect.String {
		return params, free, fmt.Errorf("cannot write string from: %s", rv.Kind())
	}

	// Get the byte representation of the string
	strEncoding := opts.StringEncoding
	var strData []byte

	switch strEncoding {
	case StringEncodingUTF8:
		strData = []byte(rv.String())
	case StringEncodingUTF16:
		encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
		strData, err = encoder.Bytes([]byte(rv.String()))
		if err != nil {
			return params, free, fmt.Errorf("failed to encode string to UTF-16: %w", err)
		}
	default:
		return params, free, fmt.Errorf("unsupported string encoding: %s", strEncoding)
	}

	// Get the alignment and size for the string data
	strCodeUnits := uint64(len(rv.String()))
	strAlignment := strEncoding.Alignment()
	taggedCodeUnitSize := strEncoding.CodeUnitSize()
	strByteLength := uint64(len(strData))

	if strByteLength%taggedCodeUnitSize != 0 {
		return params, free, fmt.Errorf("string data length %d is not a multiple of tagged code unit size %d", strByteLength, taggedCodeUnitSize)
	}

	// Allocate memory for the string data
	strDataPtr, strFree, err := abiMalloc(opts, strByteLength, strAlignment)
	if err != nil {
		return params, free, err
	}
	freeCallbacks = append(freeCallbacks, strFree)

	// Write the string data to memory
	if !opts.Memory.Write(strDataPtr, strData) {
		return params, free, fmt.Errorf("failed to write string data at %d", strDataPtr)
	}

	params = append(params, Parameter{
		Value:     strDataPtr,
		Size:      4,
		Alignment: 4,
	})
	params = append(params, Parameter{
		Value:     strCodeUnits,
		Size:      4,
		Alignment: 4,
	})

	return params, free, nil
}
