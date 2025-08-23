package abi

import (
	"errors"
	"fmt"
	"reflect"
)

// ReadRecord reads a record from memory at the specified pointer into the result.
func ReadRecord(opts AbiOptions, ptr uint64, result any) error {
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

	// If the struct has one field that fits into one value, ptr is the value already
	if rv.NumField() == 1 {
		fieldRv := rv.Field(0)
		if fieldRv.CanUint() {
			fieldRv.SetUint(ptr)
			return nil
		} else if fieldRv.CanInt() {
			fieldRv.SetInt(int64(ptr))
			return nil
		} else if fieldRv.CanFloat() {
			fieldRv.SetFloat(float64(ptr))
			return nil
		}
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

func WriteRecord(opts AbiOptions, value any, ptrHint *uint64) (ptr uint64, free AbiFreeCallback, err error) {
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
		ptr, freeRecord, err = abiMalloc(opts, size, alignment)
		if err != nil {
			return ptr, free, err
		}
		freeCallbacks = append(freeCallbacks, freeRecord)
	}

	// Write fields of the struct to linear memory
	fieldPtr := ptr
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
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

func WriteParameterRecord(opts AbiOptions, value any) (params []Parameter, free AbiFreeCallback, err error) {
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
		return params, free, errors.New("must pass a valid struct pointer value")
	}

	// Write based on the kind of the value
	if rv.Kind() != reflect.Struct {
		return params, free, fmt.Errorf("value must be a struct, got %s", rv.Kind())
	}

	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		fieldArgs, fieldFree, err := WriteParameter(opts, field.Interface())
		freeCallbacks = append(freeCallbacks, fieldFree)
		if err != nil {
			return nil, free, fmt.Errorf("failed to write field %d: %w", i, err)
		}
		params = append(params, fieldArgs...)
	}

	return params, free, nil
}
