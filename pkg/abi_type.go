package witigo

import "reflect"

type AbiType int

const (
	AbiTypeBool = iota
	AbiTypeS8
	AbiTypeS16
	AbiTypeS32
	AbiTypeS64
	AbiTypeU8
	AbiTypeU16
	AbiTypeU32
	AbiTypeU64
	AbiTypeF32
	AbiTypeF64
	AbiTypeChar
	AbiTypeString
	AbiTypeErrorContext
	AbiTypeList
	AbiTypeRecord
	AbiTypeVariant
	AbiTypeFlags
	AbiTypeOwn
	AbiTypeBorrow
	AbiTypeStream
	AbiTypeFuture
	AbiTypeTuple
	AbiTypeEnum
	AbiTypeOption
	AbiTypeResult
)

func (a AbiType) String() string {
	switch a {
	case AbiTypeBool:
		return "bool"
	case AbiTypeS8:
		return "s8"
	case AbiTypeS16:
		return "s16"
	case AbiTypeS32:
		return "s32"
	case AbiTypeS64:
		return "s64"
	case AbiTypeU8:
		return "u8"
	case AbiTypeU16:
		return "u16"
	case AbiTypeU32:
		return "u32"
	case AbiTypeU64:
		return "u64"
	case AbiTypeF32:
		return "f32"
	case AbiTypeF64:
		return "f64"
	case AbiTypeChar:
		return "char"
	case AbiTypeString:
		return "string"
	case AbiTypeErrorContext:
		return "error_context"
	case AbiTypeList:
		return "list"
	case AbiTypeRecord:
		return "record"
	case AbiTypeVariant:
		return "variant"
	case AbiTypeFlags:
		return "flags"
	case AbiTypeOwn:
		return "own"
	case AbiTypeBorrow:
		return "borrow"
	case AbiTypeStream:
		return "stream"
	case AbiTypeFuture:
		return "future"
	case AbiTypeOption:
		return "option"
	case AbiTypeTuple:
		return "tuple"
	case AbiTypeEnum:
		return "enum"
	case AbiTypeResult:
		return "result"
	default:
		return "UnknownAbiType(" + string(rune(a)) + ")"
	}
}

type AbiTypeProperties struct {
	SubTypes    []AbiTypeDefinition
	Length      *int
	ReflectType reflect.Kind
}

type AbiTypeDefinition interface {
	Type() AbiType
	Properties() AbiTypeProperties
	Alignment() int
	SizeInBytes() int
	String() string
}
