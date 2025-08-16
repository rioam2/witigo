package abi

import (
	"errors"
	"fmt"
	"reflect"
)

// ReadList reads a list from linear memory at the specified pointer into the result.
func ReadList(opts AbiOptions, ptr uint32, result any) error {
	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()

	// Validate that the result is a slice pointer
	if rv.Kind() != reflect.Slice {
		return errors.New("result must be a slice pointer")
	}

	// Ensure the result is settable
	if !rv.CanSet() {
		return errors.New("result must be a settable pointer")
	}

	// Extract ABI properties of intrinsic type
	alignment := AlignmentOf(result)
	ptr = AlignTo(ptr, alignment)

	// Extract the list data pointer from memory
	listDataPtr, ok := opts.Memory.ReadUint32Le(ptr)
	if !ok {
		return fmt.Errorf("failed to read list data pointer at %d", ptr)
	}

	// Extract the list length from memory
	listLength, ok := opts.Memory.ReadUint32Le(ptr + 4)
	if !ok {
		return fmt.Errorf("failed to read list length at %d", ptr+4)
	}

	// Create a new slice of the appropriate type
	elemType := rv.Type().Elem()
	elemSize := SizeOf(reflect.Zero(elemType).Interface())
	newSlice := reflect.MakeSlice(rv.Type(), int(listLength), int(listLength))

	// Read each element from memory and populate the new slice
	for i := 0; i < int(listLength); i++ {
		elemPtr := listDataPtr + uint32(i)*elemSize
		elemVal := reflect.New(elemType).Interface()
		err := Read(opts, elemPtr, elemVal)
		if err != nil {
			return fmt.Errorf("failed to read element %d at %d: %w", i, elemPtr, err)
		}
		elemRv := reflect.ValueOf(elemVal).Elem()
		newSlice.Index(i).Set(elemRv)
	}

	// Set the new slice to the result
	rv.Set(newSlice)
	return nil
}
