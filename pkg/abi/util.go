package abi

import (
	"errors"
	"fmt"
	"math"
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
		} else {
			return fmt.Errorf("reading struct %s is not implemented", structName)
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
}

// AlignTo aligns a pointer to the nearest multiple of the specified alignment.
func AlignTo(ptr uint32, alignment uint32) uint32 {
	if alignment <= 0 {
		return ptr
	}
	return uint32(math.Ceil(float64(ptr)/float64(alignment)) * float64(alignment))
}

// SizeOf returns the size in bytes of the given value type as defined in the Canonical ABI.
func SizeOf(value any) uint32 {
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
			size := uint32(0)
			for i := 0; i < rv.NumField(); i++ {
				field := rv.Field(i)
				fieldSize := SizeOf(field.Interface())
				fieldAlignment := AlignmentOf(field.Interface())
				size = AlignTo(size, fieldAlignment)
				size += fieldSize
			}
			recordAlignment := AlignmentOf(value)
			return AlignTo(size, recordAlignment)
		} else {
			panic(fmt.Errorf("size of struct %s is not implemented", structName))
		}
	default:
		panic("unsupported type for SizeOf: " + rv.Kind().String())
	}
}

// AlignmentOf returns the alignment in bytes of the given value type as defined in the Canonical ABI.
func AlignmentOf(value any) uint32 {
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
			alignment := uint32(1)
			for i := 0; i < rv.NumField(); i++ {
				field := rv.Field(i)
				fieldAlignment := AlignmentOf(field.Interface())
				if fieldAlignment > alignment {
					alignment = fieldAlignment
				}
			}
			return alignment
		} else {
			panic(fmt.Errorf("alignment of struct %s is not implemented", structName))
		}
	default:
		panic("unsupported type for AlignmentOf: " + rv.Kind().String())
	}
}
