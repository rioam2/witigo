package witigo

import "fmt"

const (
	ErrInvalidOptionElement = "invalid option element: must not be nil"
)

func NewAbiTypeDefinitionOption(elemType AbiTypeDefinition) AbiTypeDefinition {
	if elemType == nil {
		panic(fmt.Errorf("NewAbiTypeDefinitionOption: %s", ErrInvalidOptionElement))
	}
	caseTypes := make([]AbiTypeDefinition, 2)
	caseTypes[0] = nil      // None case
	caseTypes[1] = elemType // Some case
	return AbiTypeDefinitionVariant{caseTypes}
}
