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

func WriteList(opts AbiOptions, value any, ptrHint *uint32) (ptr uint32, free AbiFreeCallback, err error) {
	// // Validate input and retrieve element type of value
	// rv := reflect.ValueOf(value)
	// if rv.Kind() == reflect.Ptr {
	// 	rv = rv.Elem()
	// }
	// if !rv.IsValid() || rv.Kind() != reflect.Slice {
	// 	return 0, AbiFreeCallbackNoop, errors.New("must pass a valid slice pointer value")
	// }

	// // Extract the slice length and element type
	// length := uint32(rv.Len())
	// if length == 0 {
	// 	return 0, AbiFreeCallbackNoop, nil // Empty slice, nothing to write
	// }
	// elemType := rv.Type().Elem()

	// // Allocate memory for the list data if ptrHint is not provided or is zero
	// elemSize := SizeOf(reflect.Zero(elemType).Interface())
	// elemAlignment := AlignmentOf(reflect.Zero(elemType).Interface())

	// if ptrHint != nil && *ptrHint != 0 {
	// 	ptr = AlignTo(*ptrHint, elemAlignment)
	// 	free = AbiFreeCallbackNoop
	// } else {
	// 	ptr, free, err = malloc(opts, length*elemSize, elemAlignment)
	// 	if err != nil {
	// 		return 0, free, fmt.Errorf("failed to allocate memory for list data: %w", err)
	// 	}
	// }

	// // Write each element to memory
	// for i := uint32(0); i < length; i++ {
	// 	elemPtr := ptr + i*elemSize
	// 	_, elemFree, err := Write(opts, rv.Index(int(i)).Interface(), &elemPtr)

	// 	// Wrap free callback to handle element cleanup
	// 	free = func() error {
	// 		if err := free(); err != nil {
	// 			return err
	// 		}
	// 		return elemFree()
	// 	}

	// 	if err != nil {
	// 		return ptr, free, fmt.Errorf("failed to write element %d: %w", i, err)
	// 	}
	// }

	// // Write the list header (data pointer and length)
	// headerPtr := AlignTo(ptr, AlignmentOf(value))
	// if !opts.Memory.WriteUint32Le(headerPtr, ptr) ||
	// 	!opts.Memory.WriteUint32Le(headerPtr+4, uint32(length)) {
	// 	free()
	// 	return 0, AbiFreeCallbackNoop, fmt.Errorf("failed to write list header at %d", headerPtr)
	// }

	// return headerPtr, free, nil
}
