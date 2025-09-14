package abi

import (
	"errors"
	"fmt"
	"reflect"
)

// ReadEnum reads an enum value from linear memory at the specified pointer into the result.
// Enums are represented in generated code as named integer types with the suffix `Enum` whose
// underlying type is the smallest unsigned integer capable of holding all cases (u8/u16/u32/u64).
// We treat them identically to their underlying integer representation while validating the type.
func ReadEnum(opts AbiOptions, ptr uint64, result any) error {
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()
	if !isEnumType(rv) {
		return fmt.Errorf("result must be an enum pointer, got %s", rv.Type().Name())
	}
	if err := ReadInt(opts, ptr, result); err != nil {
		return err
	}
	return nil
}

// WriteEnum writes an enum value to linear memory (or ptrHint if provided) and returns the pointer.
func WriteEnum(opts AbiOptions, value any, ptrHint *uint64) (ptr uint64, free AbiFreeCallback, err error) {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return 0, AbiFreeCallbackNoop, errors.New("must pass a valid enum value")
	}
	if !isEnumType(rv) {
		return 0, AbiFreeCallbackNoop, fmt.Errorf("value must be an enum, got %s", rv.Type().Name())
	}
	return WriteInt(opts, value, ptrHint)
}

// WriteParameterEnum flattens an enum value to a single parameter containing its discriminant.
func WriteParameterEnum(opts AbiOptions, value any) (params []Parameter, free AbiFreeCallback, err error) {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return nil, AbiFreeCallbackNoop, errors.New("must pass a valid enum value")
	}
	if !isEnumType(rv) {
		return nil, AbiFreeCallbackNoop, fmt.Errorf("value must be an enum, got %s", rv.Type().Name())
	}
	return WriteParameterInt(opts, value)
}
