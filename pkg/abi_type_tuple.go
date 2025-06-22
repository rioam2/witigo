package witigo

import "fmt"

const (
	ErrInvalidTupleLength = "invalid tuple length: must contain at least one element"
)

func NewAbiTypeDefinitionTuple(elemTypes []AbiTypeDefinition) AbiTypeDefinition {
	if len(elemTypes) <= 0 {
		panic(fmt.Errorf("NewAbiTypeDefinitionTuple: %s", ErrInvalidTupleLength))
	}
	return AbiTypeDefinitionRecord{fieldTypes: elemTypes}
}
