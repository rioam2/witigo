package witigo

import "reflect"

type AbiTypeDefinitionPrimitive struct {
	t AbiType
}

func NewAbiTypeDefinitionBool() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeBool}
}

func NewAbiTypeDefinitionS8() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeS8}
}

func NewAbiTypeDefinitionU8() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeU8}
}

func NewAbiTypeDefinitionS16() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeS16}
}

func NewAbiTypeDefinitionU16() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeU16}
}

func NewAbiTypeDefinitionS32() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeS32}
}

func NewAbiTypeDefinitionU32() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeU32}
}

func NewAbiTypeDefinitionS64() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeS64}
}

func NewAbiTypeDefinitionU64() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeU64}
}

func NewAbiTypeDefinitionF32() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeF32}
}

func NewAbiTypeDefinitionF64() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeF64}
}

func NewAbiTypeDefinitionChar() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeChar}
}

func NewAbiTypeDefinitionString() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeString}
}

func NewAbiTypeDefinitionErrorContext() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeErrorContext}
}

func NewAbiTypeDefinitionOwn() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeOwn}
}

func NewAbiTypeDefinitionBorrow() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeBorrow}
}

func NewAbiTypeDefinitionStream() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeStream}
}

func NewAbiTypeDefinitionFuture() AbiTypeDefinition {
	return AbiTypeDefinitionPrimitive{AbiTypeFuture}
}

func (a AbiTypeDefinitionPrimitive) Type() AbiType {
	return a.t
}

func (a AbiTypeDefinitionPrimitive) Properties() AbiTypeProperties {
	var reflectType reflect.Kind
	switch a.t {
	case AbiTypeBool:
		reflectType = reflect.Bool
	case AbiTypeU8:
		reflectType = reflect.Uint8
	case AbiTypeS8:
		reflectType = reflect.Int8
	case AbiTypeU16:
		reflectType = reflect.Uint16
	case AbiTypeS16:
		reflectType = reflect.Int16
	case AbiTypeU32:
		reflectType = reflect.Uint32
	case AbiTypeS32:
		reflectType = reflect.Int32
	case AbiTypeF32:
		reflectType = reflect.Float32
	case AbiTypeU64:
		reflectType = reflect.Uint64
	case AbiTypeS64:
		reflectType = reflect.Int64
	case AbiTypeF64:
		reflectType = reflect.Float64
	case AbiTypeChar:
		reflectType = reflect.Int32
	case AbiTypeString, AbiTypeErrorContext:
		reflectType = reflect.String

	}
	return AbiTypeProperties{
		SubTypes:    nil,
		Length:      nil,
		ReflectType: reflectType,
	}
}

func (a AbiTypeDefinitionPrimitive) Alignment() int {
	// Primitive types have fixed alignments based on their intrinsic size
	switch a.t {
	case AbiTypeBool, AbiTypeS8, AbiTypeU8:
		return 1
	case AbiTypeS16, AbiTypeU16:
		return 2
	case AbiTypeS32, AbiTypeU32, AbiTypeF32, AbiTypeChar,
		AbiTypeString, AbiTypeErrorContext, AbiTypeOwn,
		AbiTypeBorrow, AbiTypeStream, AbiTypeFuture:
		return 4
	case AbiTypeS64, AbiTypeU64, AbiTypeF64:
		return 8
	default:
		panic("unknown ABI type for alignment: " + a.t.String())
	}
}

func (a AbiTypeDefinitionPrimitive) SizeInBytes() int {
	switch a.t {
	case AbiTypeBool, AbiTypeS8, AbiTypeU8:
		return 1
	case AbiTypeS16, AbiTypeU16:
		return 2
	case AbiTypeS32, AbiTypeU32, AbiTypeF32, AbiTypeChar,
		AbiTypeErrorContext, AbiTypeOwn, AbiTypeBorrow,
		AbiTypeStream, AbiTypeFuture:
		return 4
	case AbiTypeS64, AbiTypeU64, AbiTypeF64, AbiTypeString:
		return 8
	default:
		panic("unknown ABI type for calculating size in bytes: " + a.t.String())
	}
}

func (a AbiTypeDefinitionPrimitive) String() string {
	return a.t.String()
}
