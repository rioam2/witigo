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
