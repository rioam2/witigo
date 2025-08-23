package abi

import (
	"errors"
	"fmt"
	"reflect"
)

// ReadList reads a list from linear memory at the specified pointer into the result.
func ReadList(opts AbiOptions, ptr uint64, result any) error {
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
	for i := range uint64(listLength) {
		elemPtr := uint64(listDataPtr) + i*elemSize
		elemVal := reflect.New(elemType).Interface()
		err := Read(opts, elemPtr, elemVal)
		if err != nil {
			return fmt.Errorf("failed to read element %d at %d: %w", i, elemPtr, err)
		}
		elemRv := reflect.ValueOf(elemVal).Elem()
		newSlice.Index(int(i)).Set(elemRv)
	}

	// Set the new slice to the result
	rv.Set(newSlice)
	return nil
}

func WriteList(opts AbiOptions, value any, ptrHint *uint64) (ptr uint64, free AbiFreeCallback, err error) {
	// Initialize return values
	ptr = 0
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() || rv.Kind() != reflect.Slice {
		return ptr, free, errors.New("must pass a valid slice pointer value")
	}

	// Extract ABI properties of intrinsic type
	size := SizeOf(value)
	alignment := AlignmentOf(value)

	// Allocate memory if ptrHint is not provided or is zero
	if ptrHint != nil && *ptrHint != 0 {
		ptr = AlignTo(*ptrHint, alignment)
	} else {
		var freeList AbiFreeCallback
		ptr, freeList, err = abi_malloc(opts, size, alignment)
		if err != nil {
			return ptr, free, fmt.Errorf("failed to allocate memory for list: %w", err)
		}
		freeCallbacks = append(freeCallbacks, freeList)
	}

	listDataArgs, freeList, err := WriteParameterList(opts, value)
	if err != nil {
		return ptr, free, fmt.Errorf("failed to write list data: %w", err)
	}
	freeCallbacks = append(freeCallbacks, freeList)
	listDataPtr := listDataArgs[0].Value
	listLength := listDataArgs[1].Value

	// Write list header (data pointer and length)
	if !opts.Memory.WriteUint32Le(ptr, uint32(listDataPtr)) {
		return ptr, free, fmt.Errorf("failed to write list data pointer at %d", ptr)
	}

	if !opts.Memory.WriteUint32Le(ptr+4, uint32(listLength)) {
		return ptr, free, fmt.Errorf("failed to write list length at %d", ptr+4)
	}

	return ptr, free, nil
}

func WriteParameterList(opts AbiOptions, value any) (params []Parameter, free AbiFreeCallback, err error) {
	// Initialize return values
	params = []Parameter{}
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() || rv.Kind() != reflect.Slice {
		return params, free, errors.New("must pass a valid slice pointer value")
	}

	// Allocate memory for the list data
	listLength := uint64(rv.Len())
	elemType := rv.Type().Elem()
	elemSize := SizeOf(reflect.Zero(elemType).Interface())
	elemAlignment := AlignmentOf(reflect.Zero(elemType).Interface())
	listDataPtr, listDataFree, err := abi_malloc(opts, elemSize*listLength, elemAlignment)
	if err != nil {
		return params, free, fmt.Errorf("failed to allocate memory for list data: %w", err)
	}
	freeCallbacks = append(freeCallbacks, listDataFree)

	// Write each element to memory
	for i := range listLength {
		elemPtr := listDataPtr + i*elemSize
		_, elemFree, err := Write(opts, rv.Index(int(i)).Interface(), &elemPtr)
		freeCallbacks = append(freeCallbacks, elemFree)

		if err != nil {
			return params, free, fmt.Errorf("failed to write element %d: %w", i, err)
		}
	}

	params = append(params, Parameter{
		Value:     listDataPtr,
		Size:      4,
		Alignment: 4,
	})
	params = append(params, Parameter{
		Value:     listLength,
		Size:      4,
		Alignment: 4,
	})

	return params, free, nil
}
