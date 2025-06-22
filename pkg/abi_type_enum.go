package witigo

import "fmt"

const (
	ErrInvalidEnumLength = "invalid enum length: must contain at least one case"
)

func NewAbiTypeDefinitionEnum(length int) AbiTypeDefinition {
	if length <= 0 {
		panic(fmt.Errorf("NewAbiTypeDefinitionEnum: %s", ErrInvalidEnumLength))
	}
	caseTypes := make([]AbiTypeDefinition, length)
	return NewAbiTypeDefinitionVariant(caseTypes)
}
