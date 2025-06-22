package witigo

import "fmt"

const (
	ErrNilListElemType   = "invalid list element type: cannot be nil"
	ErrInvalidListLength = "invalid list length: a fixed-length list must be non-negative"
)

type AbiTypeDefinitionList struct {
	elemType    AbiTypeDefinition
	maybeLength *int
}

func NewAbiTypeDefinitionList(elemType AbiTypeDefinition, maybeLength *int) AbiTypeDefinition {
	if elemType == nil {
		panic(fmt.Errorf("NewAbiTypeDefinitionList: %s", ErrNilListElemType))
	}
	if maybeLength != nil && *maybeLength < 0 {
		panic(fmt.Errorf("NewAbiTypeDefinitionList: %s", ErrInvalidListLength))
	}
	return AbiTypeDefinitionList{elemType: elemType, maybeLength: maybeLength}
}

func (a AbiTypeDefinitionList) Type() AbiType {
	return AbiTypeList
}

func (a AbiTypeDefinitionList) Properties() AbiTypeProperties {
	return AbiTypeProperties{
		SubTypes: []AbiTypeDefinition{a.elemType},
		Length:   a.maybeLength,
	}
}

func (a AbiTypeDefinitionList) Alignment() int {
	// When a length is given to a list, elements are packed into linear memory
	if a.maybeLength != nil {
		return a.elemType.Alignment()
	}
	// Otherwise elements are U32 pointers
	return 4
}

func (a AbiTypeDefinitionList) SizeInBytes() int {
	// When a length is given to a list, elements are packed into linear memory
	// Size is equal to the length times the size of the element type
	if a.maybeLength != nil {
		return *a.maybeLength * a.elemType.SizeInBytes()
	}
	// Otherwise, the list is an indirect array with a pointer to the elements
	return 8 // Size of ptr to indirect array (U32) + length (U32)
}

func (a AbiTypeDefinitionList) String() string {
	return fmt.Sprintf("list<%s, length: %v>", a.elemType.String(), a.maybeLength)
}
