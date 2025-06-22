package witigo

import "fmt"

const (
	ErrInvalidResultElement = "invalid result element: must not be nil"
	ErrInvalidResultError   = "invalid error element: must not be nil"
)

func NewAbiTypeDefinitionResult(elemType AbiTypeDefinition, errType AbiTypeDefinition) AbiTypeDefinition {
	if elemType == nil {
		panic(fmt.Errorf("NewAbiTypeDefinitionResult: %s", ErrInvalidResultElement))
	}
	if errType == nil {
		panic(fmt.Errorf("NewAbiTypeDefinitionResult: %s", ErrInvalidResultError))
	}
	caseTypes := make([]AbiTypeDefinition, 2)
	caseTypes[0] = elemType // Okay case
	caseTypes[1] = errType  // Error case
	return AbiTypeDefinitionVariant{caseTypes}
}
