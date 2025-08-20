package abi

import (
	"errors"
	"fmt"
	"reflect"
)

// ReadRecord reads a record from memory at the specified pointer into the result.
func ReadRecord(opts AbiOptions, ptr uint32, result any) error {
	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()

	// Validate that the result is a struct pointer
	if rv.Kind() != reflect.Struct {
		return errors.New("result must be a struct pointer")
	}

	// Ensure the result is settable
	if !rv.CanSet() {
		return errors.New("result must be a settable pointer")
	}

	alignment := AlignmentOf(result)
	ptr = AlignTo(ptr, alignment)

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if !field.CanSet() {
			return errors.New("field is not settable")
		}
		fieldType := field.Type()
		fieldVal := reflect.New(fieldType).Interface()

		fieldSize := SizeOf(field.Interface())
		fieldAlignment := AlignmentOf(field.Interface())
		fieldPtr := AlignTo(ptr, fieldAlignment)

		err := Read(opts, fieldPtr, fieldVal)
		if err != nil {
			return fmt.Errorf("failed to read field %d at %d: %w", i, fieldPtr, err)
		}

		fieldRv := reflect.ValueOf(fieldVal).Elem()
		field.Set(fieldRv)

		ptr = fieldPtr + fieldSize
	}
	return nil
}

func WriteRecord(opts AbiOptions, value any, ptrHint *uint32) (ptr uint32, free AbiFreeCallback, err error) {
	// Initialize return values
	ptr = 0
	freeCallbacks := []AbiFreeCallback{AbiFreeCallbackNoop}
	free = func() error {
		for _, cb := range freeCallbacks {
			if err := cb(); err != nil {
				return err
			}
		}
		return nil
	}

	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return ptr, free, errors.New("must pass a valid struct pointer value")
	}

	// Write based on the kind of the value
	if rv.Kind() != reflect.Struct {
		return ptr, free, fmt.Errorf("value must be a struct, got %s", rv.Kind())
	}

	// Allocate memory if ptrHint is not provided or is zero
	size := SizeOf(value)
	alignment := AlignmentOf(value)
	if ptrHint != nil && *ptrHint != 0 {
		ptr = AlignTo(*ptrHint, alignment)
	} else {
		var freeRecord AbiFreeCallback
		ptr, freeRecord, err = malloc(opts, size, alignment)
		if err != nil {
			return ptr, free, err
		}
		freeCallbacks = append(freeCallbacks, freeRecord)
	}

	// Write fields of the struct to linear memory
	fieldPtr := ptr
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if !field.CanSet() {
			return ptr, free, fmt.Errorf("field %d is not settable", i)
		}

		fieldSize := SizeOf(field.Interface())
		fieldAlignment := AlignmentOf(field.Interface())
		fieldPtr = AlignTo(fieldPtr, fieldAlignment)

		_, fieldFree, err := Write(opts, field.Interface(), &fieldPtr)
		freeCallbacks = append(freeCallbacks, fieldFree)

		if err != nil {
			return ptr, free, fmt.Errorf("failed to write field %d at %d: %w", i, fieldPtr, err)
		}

		fieldPtr = fieldPtr + fieldSize
	}

	return ptr, free, nil
}
