package witigo

import "fmt"

const (
	ErrInvalidRecordLength = "invalid record length: must contain at least one field"
)

type AbiTypeDefinitionRecord struct {
	fieldTypes []AbiTypeDefinition
}

func NewAbiTypeDefinitionRecord(fieldTypes []AbiTypeDefinition) AbiTypeDefinition {
	if len(fieldTypes) <= 0 {
		panic(fmt.Errorf("NewAbiTypeDefinitionRecord: %s", ErrInvalidRecordLength))
	}
	return AbiTypeDefinitionRecord{fieldTypes: fieldTypes}
}

func (a AbiTypeDefinitionRecord) Type() AbiType {
	return AbiTypeRecord
}

func (a AbiTypeDefinitionRecord) Properties() AbiTypeProperties {
	length := len(a.fieldTypes)
	return AbiTypeProperties{
		SubTypes: append([]AbiTypeDefinition{}, a.fieldTypes...),
		Length:   &length,
	}
}

func (a AbiTypeDefinitionRecord) Alignment() int {
	// Record alignment is determined by the largest alignment of its fields.
	alignment := int(1)
	for _, fieldType := range a.fieldTypes {
		fieldAlignment := fieldType.Alignment()
		if fieldAlignment > alignment {
			alignment = fieldAlignment
		}
	}
	return alignment
}

func (a AbiTypeDefinitionRecord) SizeInBytes() int {
	size := int(0)
	for _, fieldType := range a.fieldTypes {
		size = AlignTo(size, fieldType.Alignment())
		size += fieldType.SizeInBytes()
	}
	return AlignTo(size, a.Alignment())
}

func (a AbiTypeDefinitionRecord) String() string {
	fields := ""
	for _, fieldType := range a.fieldTypes {
		if fields != "" {
			fields += "; "
		}
		if fieldType == nil {
			fields += "none"
			continue
		}
		fields += fieldType.String()
	}
	return fmt.Sprintf("record{%s}", fields)
}
