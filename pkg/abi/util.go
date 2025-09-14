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
		if isEnumType(rv) {
			return ReadEnum(opts, ptr, result)
		}
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
		// Treat empty struct (including anonymous struct{} payload cases) as zero-sized with alignment 1
		if rv.NumField() == 0 {
			return nil
		}
		if isStructVariantType(rv) {
			return ReadVariant(opts, ptr, result)
		} else if isStructRecordType(rv) {
			return ReadRecord(opts, ptr, result)
		} else if isStructOptionType(rv) {
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
		if isEnumType(rv) {
			return WriteEnum(opts, value, ptrHint)
		}
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
		if rv.NumField() == 0 { // empty struct{} case payload
			return 0, AbiFreeCallbackNoop, nil
		}
		if isStructVariantType(rv) {
			return WriteVariant(opts, value, ptrHint)
		} else if isStructRecordType(rv) {
			return WriteRecord(opts, value, ptrHint)
		} else if isStructOptionType(rv) {
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
		if rv.NumField() == 0 && structName == "" { // anonymous empty struct
			return 0
		}
		if isStructVariantType(rv) {
			// Variant size = size(discriminant) + max(size(case_i)) with inner value aligned, then
			// aligned to max alignment of discriminant and cases.
			if rv.NumField() == 0 {
				panic("Variant struct must contain at least a discriminant field")
			}
			discriminantSize := SizeOf(rv.Field(0).Interface())
			maxCaseSize := uint64(0)
			maxAlignment := AlignmentOf(rv.Field(0).Interface())
			for i := 1; i < rv.NumField(); i++ {
				field := rv.Field(i)
				fieldSize := SizeOf(field.Interface())
				fieldAlignment := AlignmentOf(field.Interface())
				if fieldSize > maxCaseSize {
					maxCaseSize = fieldSize
				}
				if fieldAlignment > maxAlignment {
					maxAlignment = fieldAlignment
				}
			}
			// Place case value immediately after discriminant and align whole to maxAlignment
			total := discriminantSize + maxCaseSize
			return AlignTo(total, maxAlignment)
		} else if isStructRecordType(rv) {
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
		} else if isStructOptionType(rv) {
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
		if rv.NumField() == 0 && structName == "" { // anonymous empty struct
			return 1
		}
		if isStructVariantType(rv) {
			if rv.NumField() == 0 {
				panic("Variant struct must contain at least a discriminant field")
			}
			alignment := AlignmentOf(rv.Field(0).Interface())
			for i := 1; i < rv.NumField(); i++ {
				field := rv.Field(i)
				fieldAlignment := AlignmentOf(field.Interface())
				if fieldAlignment > alignment {
					alignment = fieldAlignment
				}
			}
			return alignment
		} else if isStructRecordType(rv) {
			alignment := uint64(1)
			for i := 0; i < rv.NumField(); i++ {
				field := rv.Field(i)
				fieldAlignment := AlignmentOf(field.Interface())
				if fieldAlignment > alignment {
					alignment = fieldAlignment
				}
			}
			return alignment
		} else if isStructOptionType(rv) {
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

func isStructRecordType(rv reflect.Value) bool {
	if rv.Kind() != reflect.Struct {
		return false
	}
	structName := rv.Type().Name()
	return len(structName) >= 6 && structName[len(structName)-6:] == "Record"
}

func isStructVariantType(rv reflect.Value) bool {
	if rv.Kind() != reflect.Struct {
		return false
	}
	name := rv.Type().Name()
	if len(name) < 7 || name[len(name)-7:] != "Variant" { // suffix Variant
		return false
	}
	// Minimal structural check: first field named Type
	if rv.NumField() == 0 {
		return false
	}
	return rv.Type().Field(0).Name == "Type"
}

func isStructOptionType(rv reflect.Value) bool {
	if rv.Kind() != reflect.Struct {
		return false
	}
	structName := rv.Type().Name()
	return len(structName) >= 6 && structName[:6] == "Option"
}

// isEnumType returns true if the reflected value is a named integer type whose
// Go typename ends with the canonical "Enum" suffix produced by the code
// generator (see generateEnumTypedefFromType). Enums are represented as the
// smallest unsigned integer type capable of holding the discriminant as per
// the Canonical ABI; here we simply treat them as integers for load/store and
// parameter flattening but add this predicate so that future specialized logic
// (e.g. bounds checking) can hook in without changing call sites.
func isEnumType(rv reflect.Value) bool {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return false
	}
	kind := rv.Kind()
	switch kind {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		// continue
	default:
		return false
	}
	t := rv.Type()
	// Only user-defined (named) types should be considered – primitive
	// aliases like plain "uint8" should not match.
	if t.PkgPath() == "" { // builtin / unnamed
		return false
	}

	const enumSuffix = "Enum"
	const enumSuffixLen = len(enumSuffix)
	name := t.Name()
	return len(name) >= enumSuffixLen && name[len(name)-enumSuffixLen:] == enumSuffix
}
