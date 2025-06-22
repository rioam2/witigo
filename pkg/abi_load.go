package witigo

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

var CANONICAL_FLOAT32_NAN = []byte{0x7f, 0xc0, 0x00, 0x00}
var CANONICAL_FLOAT64_NAN = []byte{0x7f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}

const (
	ErrAbiLoadIntTypeMismatch = "type mismatch: expected %s, got %s"
	ErrByteSizeMismatch       = "byte size mismatch: expected %d bytes, got %d bytes"
)

func AbiLoadInt[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](opts AbiOptions, ptr int32, t AbiTypeDefinition) (T, error) {
	if t.Type() != AbiTypeS8 && t.Type() != AbiTypeS16 && t.Type() != AbiTypeS32 && t.Type() != AbiTypeS64 &&
		t.Type() != AbiTypeU8 && t.Type() != AbiTypeU16 && t.Type() != AbiTypeU32 && t.Type() != AbiTypeU64 {
		return T(0), fmt.Errorf("invalid type %s for AbiLoadInt", t.String())
	}

	if reflect.TypeOf(T(0)).Kind() != t.Properties().ReflectType {
		return T(0), fmt.Errorf(ErrAbiLoadIntTypeMismatch, reflect.TypeOf(T(0)).Kind().String(), t.Properties().ReflectType.String())
	}

	byteSize := int32(t.SizeInBytes())
	var value T
	data, err := opts.Memory.Read(ptr, byteSize)
	if err != nil {
		return value, err
	}

	bytesDecoded, err := binary.Decode(data, binary.LittleEndian, &value)
	if err != nil {
		return value, err
	}
	if int32(bytesDecoded) != byteSize {
		return value, fmt.Errorf(ErrByteSizeMismatch, byteSize, bytesDecoded)
	}
	return value, nil
}

func AbiLoadFloat[T float32 | float64](opts AbiOptions, ptr int32, t AbiTypeDefinition) (T, error) {
	if t.Type() != AbiTypeF32 && t.Type() != AbiTypeF64 {
		return T(0), fmt.Errorf("invalid type %s for AbiLoadFloat", t.String())
	}

	if reflect.TypeOf(T(0)).Kind() != t.Properties().ReflectType {
		return T(0), fmt.Errorf(ErrAbiLoadIntTypeMismatch, reflect.TypeOf(T(0)).Kind().String(), t.Properties().ReflectType.String())
	}

	byteSize := int32(t.SizeInBytes())
	var value T
	data, err := opts.Memory.Read(ptr, byteSize)
	if err != nil {
		return value, err
	}

	if t.Type() == AbiTypeF32 && bytes.Equal(data, CANONICAL_FLOAT32_NAN) {
		return T(math.NaN()), nil
	}
	if t.Type() == AbiTypeF64 && bytes.Equal(data, CANONICAL_FLOAT64_NAN) {
		return T(math.NaN()), nil
	}

	bytesDecoded, err := binary.Decode(data, binary.LittleEndian, &value)
	if err != nil {
		return value, err
	}
	if int32(bytesDecoded) != byteSize {
		return value, fmt.Errorf(ErrByteSizeMismatch, byteSize, bytesDecoded)
	}
	return value, nil
}

func AbiLoadBool(opts AbiOptions, ptr int32, t AbiTypeDefinition) (bool, error) {
	if t.Type() != AbiTypeBool {
		return false, fmt.Errorf("invalid type %s for AbiLoadBool", t.String())
	}

	byteSize := int32(t.SizeInBytes())
	var value bool
	data, err := opts.Memory.Read(ptr, byteSize)
	if err != nil {
		return value, err
	}

	bytesDecoded, err := binary.Decode(data, binary.LittleEndian, &value)
	if err != nil {
		return value, err
	}
	if int32(bytesDecoded) != byteSize {
		return value, fmt.Errorf(ErrByteSizeMismatch, byteSize, bytesDecoded)
	}
	return value, nil
}

func AbiLoadChar(opts AbiOptions, ptr int32, t AbiTypeDefinition) (rune, error) {
	if t.Type() != AbiTypeChar {
		return 0, fmt.Errorf("invalid type %s for AbiLoadChar", t.String())
	}

	byteSize := int32(t.SizeInBytes())
	var value rune
	data, err := opts.Memory.Read(ptr, byteSize)
	if err != nil {
		return value, err
	}

	bytesDecoded, err := binary.Decode(data, binary.LittleEndian, &value)
	if err != nil {
		return value, err
	}
	if int32(bytesDecoded) != byteSize {
		return value, fmt.Errorf(ErrByteSizeMismatch, byteSize, bytesDecoded)
	}
	return value, nil
}
