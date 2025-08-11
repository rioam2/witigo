package codegen

import (
	"fmt"
	"math"

	"github.com/golang-cz/textcase"
	"github.com/moznion/gowrtr/generator"
	witigo "github.com/rioam2/witigo/pkg"
	"github.com/rioam2/witigo/pkg/wit"
)

const emptyStructGolangTypename = "struct{}"

func GenerateTypenameFromType(w wit.WitType) string {
	if w == nil {
		return emptyStructGolangTypename
	}

	kind := w.Kind()

	if generatePrimitiveTypenameFromType(w) != "" {
		return generatePrimitiveTypenameFromType(w)
	}

	switch kind {
	case witigo.AbiTypeTuple:
		return generateTupleTypenameFromType(w)
	case witigo.AbiTypeOption:
		return generateOptionTypenameFromType(w)
	case witigo.AbiTypeList:
		return generateListTypenameFromType(w)
	case witigo.AbiTypeRecord:
		return generateRecordTypenameFromType(w)
	case witigo.AbiTypeResult:
		return generateResultTypenameFromType(w)
	case witigo.AbiTypeEnum:
		return generateEnumTypenameFromType(w)
	case witigo.AbiTypeVariant:
		return generateVariantTypenameFromType(w)
	case witigo.AbiTypeHandle:
		return generateHandleTypenameFromType(w)
	default:
		if w.Name() != "" && w.Name() != "(none)" {
			return w.Name()
		}
		panic(fmt.Sprintf("Unknown WIT type kind: %s", kind))
	}
}

func discriminantSize(n int) int {
	bits := math.Ceil(math.Log2(float64(n))/8) * 8
	return int(bits)
}

func generatePrimitiveTypenameFromType(w wit.WitType) string {
	switch w.Kind() {
	case witigo.AbiTypeString:
		return "string"
	case witigo.AbiTypeBool:
		return "bool"
	case witigo.AbiTypeS8:
		return "int8"
	case witigo.AbiTypeS16:
		return "int16"
	case witigo.AbiTypeS32:
		return "int32"
	case witigo.AbiTypeS64:
		return "int64"
	case witigo.AbiTypeU8:
		return "uint8"
	case witigo.AbiTypeU16:
		return "uint16"
	case witigo.AbiTypeU32:
		return "uint32"
	case witigo.AbiTypeU64:
		return "uint64"
	case witigo.AbiTypeF32:
		return "float32"
	case witigo.AbiTypeF64:
		return "float64"
	case witigo.AbiTypeChar:
		return "rune"
	default:
		return ""
	}
}

func generateTupleTypenameFromType(w wit.WitType) string {
	subTypes := w.SubTypes()
	if len(subTypes) == 0 {
		return "EmptyTuple"
	}
	result := ""
	for _, t := range subTypes {
		result += "-" + GenerateTypenameFromType(t.Type())
	}
	return textcase.PascalCase(result) + "Tuple"
}

func generateOptionTypenameFromType(w wit.WitType) string {
	subType := w.SubType()
	if subType == nil {
		panic("Option type must have a subtype")
	}
	return "*" + GenerateTypenameFromType(subType.Type())
}

func generateListTypenameFromType(w wit.WitType) string {
	subType := w.SubType()
	if subType == nil {
		return "[]" + emptyStructGolangTypename
	}
	return "[]" + GenerateTypenameFromType(subType.Type())
}

func generateRecordTypenameFromType(w wit.WitType) string {
	return textcase.PascalCase(w.Name()) + "Record"
}

func generateResultTypenameFromType(w wit.WitType) string {
	subTypes := w.SubTypes()
	if len(subTypes) != 2 {
		panic(fmt.Sprintf("Expected 2 subtypes for Result type, got %d", len(subTypes)))
	}
	okType := GenerateTypenameFromType(subTypes[0].Type())
	errType := GenerateTypenameFromType(subTypes[1].Type())
	return textcase.PascalCase(okType+"-"+errType) + "Result"
}

func generateEnumTypenameFromType(w wit.WitType) string {
	return textcase.PascalCase(w.Name()) + "Enum"
}

func generateVariantTypenameFromType(w wit.WitType) string {
	return textcase.PascalCase(w.Name()) + "Variant"
}

func generateHandleTypenameFromType(w wit.WitType) string {
	subType := w.SubType()
	if subType == nil {
		return "Handle"
	}
	return textcase.PascalCase(GenerateTypenameFromType(subType.Type()) + "Handle")
}

func GenerateTypedefFromType(w wit.WitType) *generator.Root {
	switch w.Kind() {
	case witigo.AbiTypeRecord:
		return generateRecordTypedefFromType(w)
	case witigo.AbiTypeResult:
		return generateResultTypedefFromType(w)
	case witigo.AbiTypeTuple:
		return generateTupleTypedefFromType(w)
	case witigo.AbiTypeEnum:
		return generateEnumTypedefFromType(w)
	case witigo.AbiTypeVariant:
		return generateVariantTypedefFromType(w)
	case witigo.AbiTypeHandle:
		return generateHandleTypedefFromType(w)
	default:
		// Remaining types are either primitive or do not require a typedef
		return nil
	}
}

func generateRecordTypedefFromType(w wit.WitType) *generator.Root {
	typeDef := generator.NewStruct(GenerateTypenameFromType(w))
	for _, field := range w.SubTypes() {
		typeDef = typeDef.AddField(
			textcase.PascalCase(field.Name()),
			GenerateTypenameFromType(field.Type()),
		)
	}
	return generator.NewRoot(typeDef)
}

func generateResultTypedefFromType(w wit.WitType) *generator.Root {
	okType := GenerateTypenameFromType(w.SubTypes()[0].Type())
	errType := GenerateTypenameFromType(w.SubTypes()[1].Type())
	return generator.NewRoot(
		generator.NewStruct(GenerateTypenameFromType(w)).
			AddField("Ok", okType).
			AddField("Error", errType),
	)
}

func generateTupleTypedefFromType(w wit.WitType) *generator.Root {
	subTypes := w.SubTypes()
	typeDef := generator.NewStruct(GenerateTypenameFromType(w))
	for i, subType := range subTypes {
		typeDef = typeDef.AddField(
			textcase.PascalCase(fmt.Sprintf("Elem%d", i)),
			GenerateTypenameFromType(subType.Type()),
		)
	}
	return generator.NewRoot(typeDef)
}

func generateEnumTypedefFromType(w wit.WitType) *generator.Root {
	root := generator.NewRoot()
	discriminantType := fmt.Sprintf("uint%d", discriminantSize(len(w.SubTypes())))
	enumTypedef := generator.NewRawStatementf("type %s %s", GenerateTypenameFromType(w), discriminantType)
	root = root.AddStatements(enumTypedef)
	for i, c := range w.SubTypes() {
		statement := generator.NewRawStatementf(
			"const %s = %d",
			GenerateTypenameFromType(w)+textcase.PascalCase(c.Name()),
			i,
		)
		root = root.AddStatements(statement)
	}
	return root
}

func generateVariantTypedefFromType(w wit.WitType) *generator.Root {
	root := generator.NewRoot()
	enumTypedefName := GenerateTypenameFromType(w) + "Type"
	enumTypedef := generator.NewRawStatementf("type %s int", enumTypedefName)
	root = root.AddStatements(enumTypedef)
	for i, c := range w.SubTypes() {
		statement := generator.NewRawStatementf(
			"const %s = %d",
			enumTypedefName+textcase.PascalCase(c.Name()),
			i,
		)
		root = root.AddStatements(statement)
	}
	structTypedef := generator.NewStruct(GenerateTypenameFromType(w))
	structTypedef = structTypedef.AddField(
		textcase.PascalCase("Type"),
		enumTypedefName,
	)
	for _, field := range w.SubTypes() {
		fieldType := "struct{}"
		if field.Type() != nil {
			fieldType = GenerateTypenameFromType(field.Type())
		}
		structTypedef = structTypedef.AddField(
			textcase.PascalCase(field.Name()),
			fieldType,
		)
	}

	root = root.AddStatements(structTypedef)
	return root
}

func generateHandleTypedefFromType(w wit.WitType) *generator.Root {
	return generator.NewRoot(
		generator.NewStruct(GenerateTypenameFromType(w)).
			AddField(
				textcase.PascalCase("Type"),
				textcase.PascalCase(GenerateTypenameFromType(w.SubType().Type())),
			),
	)
}
