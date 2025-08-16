package abi

import (
	"errors"
	"fmt"
	"reflect"

	"golang.org/x/text/encoding/unicode"
)

// ReadString reads a string from linear memory at the specified pointer into the result.
func ReadString(opts AbiOptions, ptr uint32, result any) error {
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
	if strPtr != AlignTo(strPtr, strAlignment) {
		return fmt.Errorf("string pointer %d is not aligned to %d bytes", strPtr, strAlignment)
	}

	// Validate that the string pointer is within bounds
	strByteLength := taggedCodeUnits * taggedCodeUnitSize
	if strPtr+strByteLength > opts.Memory.Size() {
		return fmt.Errorf("string pointer %d with length %d exceeds memory size %d", strPtr, strByteLength, opts.Memory.Size())
	}

	// Read the string data from memory
	strData, ok := opts.Memory.Read(strPtr, strByteLength)
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
