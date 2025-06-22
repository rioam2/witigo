package witigo

import (
	"fmt"
	"math"
)

const (
	ErrInvalidVariantCases = "invalid variant cases: at least one case type is required"
)

type AbiTypeDefinitionVariant struct {
	caseTypes []AbiTypeDefinition
}

func NewAbiTypeDefinitionVariant(caseTypes []AbiTypeDefinition) AbiTypeDefinition {
	if len(caseTypes) <= 0 {
		panic(fmt.Errorf("NewAbiTypeDefinitionVariant: %s", ErrInvalidVariantCases))
	}
	return AbiTypeDefinitionVariant{caseTypes: caseTypes}
}

func (a AbiTypeDefinitionVariant) Type() AbiType {
	return AbiTypeVariant
}

func (a AbiTypeDefinitionVariant) Properties() AbiTypeProperties {
	length := len(a.caseTypes)
	return AbiTypeProperties{
		SubTypes: append([]AbiTypeDefinition{}, a.caseTypes...),
		Length:   &length,
	}
}

func (a AbiTypeDefinitionVariant) DiscriminantType() AbiTypeDefinition {
	n := len(a.caseTypes)
	if n <= 8 {
		return AbiTypeDefinitionPrimitive{AbiTypeU8}
	} else if n <= 16 {
		return AbiTypeDefinitionPrimitive{AbiTypeU16}
	} else {
		return AbiTypeDefinitionPrimitive{AbiTypeU32}
	}
}

func (a AbiTypeDefinitionVariant) MaxCaseAlignment() int {
	n := len(a.caseTypes)
	if n > math.MaxInt {
		panic("invalid number of variant cases: " + fmt.Sprint(n))
	}
	alignment := int(1)
	for _, caseType := range a.caseTypes {
		if caseType != nil {
			caseAlignment := caseType.Alignment()
			if caseAlignment > alignment {
				alignment = caseAlignment
			}
		}
	}
	return alignment
}

func (a AbiTypeDefinitionVariant) Alignment() int {
	discriminantAlignment := a.DiscriminantType().Alignment()
	maxCaseAlignment := a.MaxCaseAlignment()
	if discriminantAlignment > maxCaseAlignment {
		return discriminantAlignment
	}
	return maxCaseAlignment
}

func (a AbiTypeDefinitionVariant) SizeInBytes() int {
	size := a.DiscriminantType().SizeInBytes()
	size = AlignTo(size, a.MaxCaseAlignment())
	maxCaseSize := int(0)
	for _, caseType := range a.caseTypes {
		if caseType != nil {
			caseSize := caseType.SizeInBytes()
			if caseSize > maxCaseSize {
				maxCaseSize = caseSize
			}
		}
	}
	size += maxCaseSize
	return AlignTo(size, a.Alignment())
}

func (a AbiTypeDefinitionVariant) String() string {
	cases := ""
	for _, caseType := range a.caseTypes {
		if cases != "" {
			cases += "; "
		}
		if caseType == nil {
			cases += "none"
			continue
		}
		cases += caseType.String()
	}
	return fmt.Sprintf("variant{%s}", cases)
}
