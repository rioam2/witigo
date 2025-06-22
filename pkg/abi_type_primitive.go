package witigo

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
	return AbiTypeProperties{
		SubTypes: nil,
		Length:   nil,
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
