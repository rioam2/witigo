package witigo

import "reflect"

type AbiTypeDefinitionFlags struct {
	fieldTypes []AbiTypeDefinition
}

func NewAbiTypeDefinitionFlags(fieldTypes []AbiTypeDefinition) AbiTypeDefinition {
	return NewAbiTypeDefinitionVariant(fieldTypes)
}

func (a AbiTypeDefinitionFlags) Type() AbiType {
	return AbiTypeFlags
}

func (a AbiTypeDefinitionFlags) Properties() AbiTypeProperties {
	length := len(a.fieldTypes)
	var reflectType reflect.Kind
	if length <= 8 {
		reflectType = reflect.Uint8
	} else if length <= 16 {
		reflectType = reflect.Uint16
	} else {
		reflectType = reflect.Uint32
	}
	return AbiTypeProperties{
		SubTypes:    append([]AbiTypeDefinition{}, a.fieldTypes...),
		Length:      &length,
		ReflectType: reflectType,
	}
}

func (a AbiTypeDefinitionFlags) Alignment() int {
	// Flags are packed bit-vectors so alignment is based on the size of entries
	n := len(a.fieldTypes)
	if n <= 8 {
		return 1
	} else if n <= 16 {
		return 2
	} else {
		return 4
	}
}

func (a AbiTypeDefinitionFlags) SizeInBytes() int {
	return a.Alignment() // Flags are packed, so size is equal to alignment
}
