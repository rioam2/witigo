package abi

import (
	"math"
	"reflect"
)

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
	case reflect.Int32, reflect.Uint32, reflect.Float32:
		return 4
	case reflect.Int64, reflect.Uint64, reflect.Float64, reflect.String, reflect.Slice:
		return 8
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
	case reflect.Int32, reflect.Uint32, reflect.Float32, reflect.String, reflect.Slice:
		return 4
	case reflect.Int64, reflect.Uint64, reflect.Float64:
		return 8
	default:
		panic("unsupported type for AlignmentOf: " + rv.Kind().String())
	}
}
