package abi

import (
	"errors"
	"fmt"
	"reflect"
)

// Read reads a value from memory at the specified pointer into the result.
func Read(opts AbiOptions, ptr uint32, result any) error {
	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()

	// Read based on the kind of the result
	switch rv.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ReadInt(opts, ptr, result)
	case reflect.Bool:
		return ReadBool(opts, ptr, result)
	case reflect.Float32, reflect.Float64:
		return ReadFloat(opts, ptr, result)
	case reflect.String:
		return ReadString(opts, ptr, result)
	case reflect.Slice:
		return ReadList(opts, ptr, result)
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
}
