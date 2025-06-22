package witigo

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

func AbiLoadInt[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64](opts AbiOptions, ptr int32, t AbiTypeDefinition) (T, error) {
	if t.Type() != AbiTypeS8 && t.Type() != AbiTypeS16 && t.Type() != AbiTypeS32 && t.Type() != AbiTypeS64 &&
		t.Type() != AbiTypeU8 && t.Type() != AbiTypeU16 && t.Type() != AbiTypeU32 && t.Type() != AbiTypeU64 {
		return T(0), fmt.Errorf("invalid type %s for AbiLoadInt", t.String())
	}

	if reflect.TypeOf(T(0)).Kind() != t.Properties().ReflectType {
		return T(0), fmt.Errorf("type mismatch: expected %s, got %s", reflect.TypeOf(T(0)).Kind().String(), t.Properties().ReflectType.String())
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
		return value, fmt.Errorf("expected %d bytes, got %d", byteSize, bytesDecoded)
	}
	return value, nil

}
