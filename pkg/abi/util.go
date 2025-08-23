package abi

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

// AbiFreeCallback is a returned function from write operations that can be used to free resources.
type AbiFreeCallback func() error

// AbiFreeCallbackNoop is a no-operation callback that does nothing.
var AbiFreeCallbackNoop = func() error {
	return nil
}

// Read reads a value from memory at the specified pointer into the result.
func Read(opts AbiOptions, ptr uint64, result any) error {
	// Validate input and retrieve element type of result
	rv := reflect.ValueOf(result)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return errors.New("must pass a non-nil pointer result")
	}
	rv = rv.Elem()

	// Read based on the kind of the result
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ReadInt(opts, ptr, result)
	case reflect.Bool:
		return ReadBool(opts, ptr, result)
	case reflect.Float32, reflect.Float64:
		return ReadFloat(opts, ptr, result)
	case reflect.String:
		return ReadString(opts, ptr, result)
	case reflect.Slice:
		return ReadList(opts, ptr, result)
	case reflect.Struct:
		structName := rv.Type().Name()
		if len(structName) >= 6 && structName[len(structName)-6:] == "Record" {
			return ReadRecord(opts, ptr, result)
		} else if len(structName) >= 6 && structName[:6] == "Option" {
			return ReadOption(opts, ptr, result)
		} else {
			return fmt.Errorf("reading struct %s is not implemented", structName)
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
}

// Write writes a value to memory at the specified pointer from the result.
func Write(opts AbiOptions, value any, ptrHint *uint64) (ptr uint64, free AbiFreeCallback, err error) {
	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return ptr, free, errors.New("must pass a valid value")
	}

	// Write based on the kind of the value
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return WriteInt(opts, value, ptrHint)
	case reflect.Bool:
		return WriteBool(opts, value, ptrHint)
	case reflect.Float32, reflect.Float64:
		return WriteFloat(opts, value, ptrHint)
	case reflect.String:
		return WriteString(opts, value, ptrHint)
	case reflect.Slice:
		return WriteList(opts, value, ptrHint)
	case reflect.Struct:
		structName := rv.Type().Name()
		if len(structName) >= 6 && structName[len(structName)-6:] == "Record" {
			return WriteRecord(opts, value, ptrHint)
		} else if len(structName) >= 6 && structName[:6] == "Option" {
			return WriteOption(opts, value, ptrHint)
		} else {
			return 0, AbiFreeCallbackNoop, fmt.Errorf("writing struct %s is not implemented", structName)
		}
	default:
		return 0, AbiFreeCallbackNoop, fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
}

// AlignTo aligns a pointer to the nearest multiple of the specified alignment.
func AlignTo(ptr uint64, alignment uint64) uint64 {
	if alignment <= 0 {
		return ptr
	}
	return uint64(math.Ceil(float64(ptr)/float64(alignment)) * float64(alignment))
}

// SizeOf returns the size in bytes of the given value type as defined in the Canonical ABI.
func SizeOf(value any) uint64 {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		return 1
	case reflect.Int16, reflect.Uint16:
		return 2
	case reflect.Int, reflect.Uint, reflect.Int32, reflect.Uint32, reflect.Float32:
		return 4
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.String, reflect.Slice:
		return 8
	case reflect.Struct:
		structName := rv.Type().Name()
		if len(structName) >= 6 && structName[len(structName)-6:] == "Record" {
			size := uint64(0)
			for i := 0; i < rv.NumField(); i++ {
				field := rv.Field(i)
				fieldSize := SizeOf(field.Interface())
				fieldAlignment := AlignmentOf(field.Interface())
				size = AlignTo(size, fieldAlignment)
				size += fieldSize
			}
			recordAlignment := AlignmentOf(value)
			return AlignTo(size, recordAlignment)
		} else if len(structName) >= 6 && structName[:6] == "Option" {
			numFields := rv.NumField()
			if numFields != 2 {
				panic(fmt.Errorf("Option type must contain only discriminant and value fields"))
			}
			discriminantRv := rv.Field(0)
			valueRv := rv.Field(1)
			valueAlignment := AlignmentOf(valueRv.Interface())
			totalSize := SizeOf(discriminantRv.Interface()) + SizeOf(valueRv.Interface())
			return AlignTo(totalSize, valueAlignment)
		} else {
			panic(fmt.Errorf("size of struct %s is not implemented", structName))
		}
	default:
		panic("unsupported type for SizeOf: " + rv.Kind().String())
	}
}

// AlignmentOf returns the alignment in bytes of the given value type as defined in the Canonical ABI.
func AlignmentOf(value any) uint64 {
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	switch rv.Kind() {
	case reflect.Bool, reflect.Int8, reflect.Uint8:
		return 1
	case reflect.Int16, reflect.Uint16:
		return 2
	case reflect.Int, reflect.Uint, reflect.Int32, reflect.Uint32, reflect.Float32, reflect.String, reflect.Slice:
		return 4
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 8
	case reflect.Struct:
		structName := rv.Type().Name()
		if len(structName) >= 6 && structName[len(structName)-6:] == "Record" {
			alignment := uint64(1)
			for i := 0; i < rv.NumField(); i++ {
				field := rv.Field(i)
				fieldAlignment := AlignmentOf(field.Interface())
				if fieldAlignment > alignment {
					alignment = fieldAlignment
				}
			}
			return alignment
		} else if len(structName) >= 6 && structName[:6] == "Option" {
			numFields := rv.NumField()
			if numFields != 2 {
				panic(fmt.Errorf("Option type must contain only discriminant and value fields"))
			}
			valueRv := rv.Field(1)
			alignment := AlignmentOf(valueRv.Interface())
			return alignment
		} else {
			panic(fmt.Errorf("alignment of struct %s is not implemented", structName))
		}
	default:
		panic("unsupported type for AlignmentOf: " + rv.Kind().String())
	}
}

func wrapFreeCallbacks(freeCallbacks *[]AbiFreeCallback) AbiFreeCallback {
	return func() error {
		for _, cb := range *freeCallbacks {
			if err := cb(); err != nil {
				return err
			}
		}
		return nil
	}
}
