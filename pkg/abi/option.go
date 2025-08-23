package abi

import (
	"errors"
	"fmt"
	"reflect"
)

type Option[T any] struct {
	IsSome bool
	Value  T
}

func ReadOption(opts AbiOptions, ptr uint64, result any) error {
	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()

	// Check if the result is an Option type
	structName := rv.Type().Name()
	if rv.Kind() != reflect.Struct || len(structName) < 6 || structName[:6] != "Option" {
		return fmt.Errorf("expected Option type, got %s", structName)
	}

	// Read the discriminant
	var isSome bool
	if err := Read(opts, ptr, &isSome); err != nil {
		return err
	}
	rv.Field(0).SetBool(isSome)

	// If the discriminant indicates None, set the value to nil
	if !isSome {
		return nil
	}

	fieldRv := rv.Field(1)
	valueAlignment := AlignmentOf(fieldRv.Interface())
	valuePtr := AlignTo(ptr+1, valueAlignment)

	// Read the value into the second field of the Option struct
	return Read(opts, valuePtr, fieldRv.Addr().Interface())
}

func WriteOption(opts AbiOptions, value any, ptrHint *uint64) (ptr uint64, free AbiFreeCallback, err error) {
	// Initialize return values
	ptr = 0
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Struct && (value == nil || rv.IsZero() || rv.IsNil()) {
		return ptr, free, errors.New("must pass a non-nil pointer value")
	}
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	// Check if the result is an Option type
	structName := rv.Type().Name()
	if rv.Kind() != reflect.Struct || len(structName) < 6 || structName[:6] != "Option" {
		return ptr, free, fmt.Errorf("expected Option type, got %s", structName)
	}

	// Extract ABI properties of intrinsic type
	size := SizeOf(value)
	alignment := AlignmentOf(value)

	// Allocate memory if ptrHint is not provided or is zero
	if ptrHint != nil && *ptrHint != 0 {
		ptr = AlignTo(*ptrHint, alignment)
	} else {
		var freeOption AbiFreeCallback
		ptr, freeOption, err = abi_malloc(opts, size, alignment)
		if err != nil {
			return ptr, free, err
		}
		freeCallbacks = append(freeCallbacks, freeOption)
	}

	discriminant := rv.Field(0).Bool()
	discriminantBytes := []byte{0x00}
	if discriminant {
		discriminantBytes[0] = 1
	}
	valueInterface := rv.Field(1).Interface()
	valuePtr := AlignTo(ptr+1, alignment)

	// Write string descriptor to linear memory
	if ok := opts.Memory.Write(ptr, discriminantBytes); !ok {
		return ptr, free, fmt.Errorf("failed to write string data pointer at %d", ptr)
	}

	_, valueFree, err := Write(opts, valueInterface, &valuePtr)
	freeCallbacks = append(freeCallbacks, valueFree)
	if err != nil {
		return ptr, free, err
	}

	return ptr, free, nil
}

func WriteParameterOption(opts AbiOptions, value any) (params []Parameter, free AbiFreeCallback, err error) {
	// Initialize return values
	params = []Parameter{}
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Struct && (value == nil || rv.IsZero() || rv.IsNil()) {
		return params, free, errors.New("must pass a non-nil pointer value")
	}
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	// Check if the result is an Option type
	structName := rv.Type().Name()
	if rv.Kind() != reflect.Struct || len(structName) < 6 || structName[:6] != "Option" {
		return params, free, fmt.Errorf("expected Option type, got %s", structName)
	}

	discriminant := rv.Field(0).Bool()
	discriminantUint := uint64(0)
	if discriminant {
		discriminantUint = 1
	}
	valueInterface := rv.Field(1).Interface()
	valueParams, valueFree, err := WriteParameter(opts, valueInterface)
	freeCallbacks = append(freeCallbacks, valueFree)
	if err != nil {
		return params, free, err
	}

	discriminantAlignment := uint64(1)
	if len(valueParams) > 0 && valueParams[0].Alignment > discriminantAlignment {
		discriminantAlignment = valueParams[0].Alignment
	}

	params = append(params, Parameter{
		Value:     discriminantUint,
		Size:      1,
		Alignment: discriminantAlignment,
	})
	params = append(params, valueParams...)

	return params, free, nil
}
