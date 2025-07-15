package wit

import (
	"encoding/json"
	"fmt"

	"github.com/golang-cz/textcase"
	"github.com/moznion/gowrtr/generator"
	witigo "github.com/rioam2/witigo/pkg"
)

const emptyStructGolangTypename = "struct{}"

type WitType interface {
	Name() string
	Kind() witigo.AbiType
	String() string
	SubType() WitTypeReference
	SubTypes() []WitTypeReference
	IsPrimitive() bool
	CodegenGolangTypename() string
	CodegenGolangTypedef() *generator.Root
}

type WitTypeImpl struct {
	Raw  json.RawMessage
	Root WitDefinition
}

var _ WitType = &WitTypeImpl{}

func (w *WitTypeImpl) Name() string {
	var data struct {
		Name *string `json:"name"`
	}
	json.Unmarshal(w.Raw, &data)
	if data.Name == nil {
		return "(none)"
	}
	return *data.Name
}

func (w *WitTypeImpl) Kind() witigo.AbiType {
	var data struct {
		Kind *struct {
			List    *json.RawMessage `json:"list"`
			Option  *json.RawMessage `json:"option"`
			Record  *json.RawMessage `json:"record"`
			Tuple   *json.RawMessage `json:"tuple"`
			Result  *json.RawMessage `json:"result"`
			Variant *json.RawMessage `json:"variant"`
			Enum    *json.RawMessage `json:"enum"`
		} `json:"kind"`
		Type *string `json:"type"`
	}
	json.Unmarshal(w.Raw, &data)
	if data.Type != nil {
		switch *data.Type {
		case "string":
			return witigo.AbiTypeString
		case "bool":
			return witigo.AbiTypeBool
		case "s8":
			return witigo.AbiTypeS8
		case "s16":
			return witigo.AbiTypeS16
		case "s32":
			return witigo.AbiTypeS32
		case "s64":
			return witigo.AbiTypeS64
		case "u8":
			return witigo.AbiTypeU8
		case "u16":
			return witigo.AbiTypeU16
		case "u32":
			return witigo.AbiTypeU32
		case "u64":
			return witigo.AbiTypeU64
		case "f32":
			return witigo.AbiTypeF32
		case "f64":
			return witigo.AbiTypeF64
		case "char":
			return witigo.AbiTypeChar
		}
	}
	if data.Kind.List != nil {
		return witigo.AbiTypeList
	}
	if data.Kind.Option != nil {
		return witigo.AbiTypeOption
	}
	if data.Kind.Record != nil {
		return witigo.AbiTypeRecord
	}
	if data.Kind.Tuple != nil {
		return witigo.AbiTypeTuple
	}
	if data.Kind.Result != nil {
		return witigo.AbiTypeResult
	}
	if data.Kind.Variant != nil {
		return witigo.AbiTypeVariant
	}
	if data.Kind.Enum != nil {
		return witigo.AbiTypeEnum
	}
	panic(fmt.Sprintf("Unknown WIT type kind: %v", data.Kind))
}

func (w *WitTypeImpl) SubType() WitTypeReference {
	var data struct {
		Kind struct {
			List   *any `json:"list"`
			Option *any `json:"option"`
		} `json:"kind"`
	}
	json.Unmarshal(w.Raw, &data)
	var listOrOptionType any
	if data.Kind.List != nil {
		listOrOptionType = data.Kind.List
	} else if data.Kind.Option != nil {
		listOrOptionType = data.Kind.Option
	}
	rawTypeRef, err := json.Marshal(map[string]any{
		"type": listOrOptionType,
		"name": w.Name(),
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal list type reference: %v", err))
	}
	return &WitTypeReferenceImpl{Raw: rawTypeRef, Root: w.Root}
}

func (w *WitTypeImpl) SubTypes() []WitTypeReference {
	var data struct {
		Kind struct {
			Record  *json.RawMessage `json:"record"`
			Variant *json.RawMessage `json:"variant"`
			Enum    *json.RawMessage `json:"enum"`
			Tuple   *json.RawMessage `json:"tuple"`
			Result  *json.RawMessage `json:"result"`
		} `json:"kind"`
	}
	json.Unmarshal(w.Raw, &data)
	if data.Kind.Record != nil {
		return w.handleRecordType(data.Kind.Record)
	}
	if data.Kind.Variant != nil {
		return w.handleVariantType(data.Kind.Variant)
	}
	if data.Kind.Enum != nil {
		return w.handleEnumType(data.Kind.Enum)
	}
	if data.Kind.Tuple != nil {
		return w.handleTupleType(data.Kind.Tuple)
	}
	if data.Kind.Result != nil {
		return w.handleResultType(data.Kind.Result)
	}
	return nil
}

func (w *WitTypeImpl) handleRecordType(rawRecord *json.RawMessage) []WitTypeReference {
	var subTypes []WitTypeReference
	var record struct {
		Fields []json.RawMessage `json:"fields"`
	}
	json.Unmarshal(*rawRecord, &record)
	for _, field := range record.Fields {
		subTypes = append(subTypes, &WitTypeReferenceImpl{Raw: field, Root: w.Root})
	}
	return subTypes
}

func (w *WitTypeImpl) handleVariantType(rawVariant *json.RawMessage) []WitTypeReference {
	var subTypes []WitTypeReference
	var variant struct {
		Cases []json.RawMessage `json:"cases"`
	}
	json.Unmarshal(*rawVariant, &variant)
	for _, c := range variant.Cases {
		subTypes = append(subTypes, &WitTypeReferenceImpl{Raw: c, Root: w.Root})
	}
	return subTypes
}

func (w *WitTypeImpl) handleEnumType(rawEnum *json.RawMessage) []WitTypeReference {
	var subTypes []WitTypeReference
	var enum struct {
		Cases []struct {
			Name string `json:"name"`
		} `json:"cases"`
	}
	json.Unmarshal(*rawEnum, &enum)
	for _, c := range enum.Cases {
		remappedType, err := json.Marshal(map[string]any{
			// TODO: I believe this is dependent on number of cases
			"type": "u32",
			"name": c.Name,
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal enum type reference: %v", err))
		}
		subTypes = append(subTypes, &WitTypeReferenceImpl{Raw: remappedType, Root: w.Root})
	}
	return subTypes
}

func (w *WitTypeImpl) handleTupleType(rawTuple *json.RawMessage) []WitTypeReference {
	var subTypes []WitTypeReference
	var tuple struct {
		Types []any `json:"types"`
	}
	json.Unmarshal(*rawTuple, &tuple)
	for _, t := range tuple.Types {
		remappedType, err := json.Marshal(map[string]any{
			"type": t,
			"name": nil,
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal tuple type reference: %v", err))
		}
		subTypes = append(subTypes, &WitTypeReferenceImpl{Raw: remappedType, Root: w.Root})
	}
	return subTypes
}

func (w *WitTypeImpl) handleResultType(rawResult *json.RawMessage) []WitTypeReference {
	var result struct {
		Ok  *any `json:"ok"`
		Err *any `json:"err"`
	}
	json.Unmarshal(*rawResult, &result)
	okType, err := json.Marshal(map[string]any{
		"type": result.Ok,
		"name": "ok",
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal result ok type reference: %v", err))
	}
	errType, err := json.Marshal(map[string]any{
		"type": result.Err,
		"name": "error",
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal result err type reference: %v", err))
	}
	return []WitTypeReference{
		&WitTypeReferenceImpl{Raw: okType, Root: w.Root},
		&WitTypeReferenceImpl{Raw: errType, Root: w.Root},
	}
}

func (w *WitTypeImpl) IsPrimitive() bool {
	switch w.Kind() {
	case witigo.AbiTypeString, witigo.AbiTypeBool, witigo.AbiTypeS8, witigo.AbiTypeS16,
		witigo.AbiTypeS32, witigo.AbiTypeS64, witigo.AbiTypeU8,
		witigo.AbiTypeU16, witigo.AbiTypeU32, witigo.AbiTypeU64,
		witigo.AbiTypeF32, witigo.AbiTypeF64, witigo.AbiTypeChar:
		return true
	default:
		return false
	}
}

func (w *WitTypeImpl) String() string {
	base := w.Kind().String()
	switch w.Kind() {
	case witigo.AbiTypeList, witigo.AbiTypeOption:
		base = w.formatSingleTypeContainer(base)
	case witigo.AbiTypeRecord, witigo.AbiTypeVariant, witigo.AbiTypeEnum:
		base = w.formatNamedTypes(base)
	case witigo.AbiTypeTuple, witigo.AbiTypeResult:
		base = w.formatUnnamedTypes(base)
	}
	return base
}

func (w *WitTypeImpl) formatSingleTypeContainer(base string) string {
	listType := w.SubType()
	if listType != nil {
		return base + "<" + listType.String() + ">"
	}
	return base
}

func (w *WitTypeImpl) formatNamedTypes(base string) string {
	cases := w.SubTypes()
	if len(cases) == 0 {
		return base
	}
	result := base + "{ "
	for i, c := range cases {
		if i > 0 {
			result += ", "
		}
		result += c.Name() + ": " + c.String()
	}
	return result + " }"
}

func (w *WitTypeImpl) formatUnnamedTypes(base string) string {
	types := w.SubTypes()
	if len(types) == 0 {
		return base
	}
	result := base + "<"
	for i, t := range types {
		if i > 0 {
			result += ", "
		}
		result += t.String()
	}
	return result + ">"
}

func (w *WitTypeImpl) CodegenGolangTypename() string {
	kind := w.Kind()

	if w.codegenGolangPrimitiveTypename() != "" {
		return w.codegenGolangPrimitiveTypename()
	}

	switch kind {
	case witigo.AbiTypeTuple:
		return w.codegenTupleGolangTypename()
	case witigo.AbiTypeOption:
		return w.codegenOptionGolangTypename()
	case witigo.AbiTypeList:
		return w.codegenListGolangTypename()
	case witigo.AbiTypeRecord:
		return w.codegenRecordGolangTypename()
	case witigo.AbiTypeResult:
		return w.codegenResultGolangTypename()
	case witigo.AbiTypeEnum:
		return w.codegenEnumGolangTypename()
	case witigo.AbiTypeVariant:
		return w.codegenVariantGolangTypename()
	default:
		if w.Name() != "" && w.Name() != "(none)" {
			return w.Name()
		}
		panic(fmt.Sprintf("Unknown WIT type kind: %s", kind))
	}
}

func (w *WitTypeImpl) codegenGolangPrimitiveTypename() string {
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

func (w *WitTypeImpl) codegenTupleGolangTypename() string {
	subTypes := w.SubTypes()
	if len(subTypes) == 0 {
		return "EmptyTuple"
	}
	result := ""
	for _, t := range subTypes {
		result += "-" + t.Type().CodegenGolangTypename()
	}
	return textcase.PascalCase(result) + "Tuple"
}

func (w *WitTypeImpl) codegenOptionGolangTypename() string {
	subType := w.SubType()
	if subType == nil {
		panic("Option type must have a subtype")
	}
	return "*" + subType.Type().CodegenGolangTypename()
}

func (w *WitTypeImpl) codegenListGolangTypename() string {
	subType := w.SubType()
	if subType == nil {
		return "[]" + emptyStructGolangTypename
	}
	return "[]" + subType.Type().CodegenGolangTypename()
}

func (w *WitTypeImpl) codegenRecordGolangTypename() string {
	return textcase.PascalCase(w.Name()) + "Record"
}

func (w *WitTypeImpl) codegenResultGolangTypename() string {
	subTypes := w.SubTypes()
	if len(subTypes) != 2 {
		panic(fmt.Sprintf("Expected 2 subtypes for Result type, got %d", len(subTypes)))
	}
	okType := subTypes[0].Type().CodegenGolangTypename()
	errType := subTypes[1].Type().CodegenGolangTypename()
	return textcase.PascalCase(okType+"-"+errType) + "Result"
}

func (w *WitTypeImpl) codegenEnumGolangTypename() string {
	return textcase.PascalCase(w.Name()) + "Enum"
}

func (w *WitTypeImpl) codegenVariantGolangTypename() string {
	return textcase.PascalCase(w.Name()) + "Variant"
}

func (w *WitTypeImpl) CodegenGolangTypedef() *generator.Root {
	var baseType WitType = w
	for baseType != nil && baseType.SubType().Type() != nil {
		baseType = baseType.SubType().Type()
	}

	// Primitive base types do not need a typedef
	if baseType.IsPrimitive() {
		return nil
	}

	switch baseType.Kind() {
	case witigo.AbiTypeRecord:
		fields := make([]*generator.FuncSignature, len(baseType.SubTypes()))
		for i, field := range baseType.SubTypes() {
			fields[i] = generator.
				NewFuncSignature(textcase.PascalCase(field.Name())).
				AddReturnTypes(field.Type().CodegenGolangTypename())
		}
		return generator.NewRoot(
			generator.NewInterface(
				w.CodegenGolangTypename(),
				fields...,
			),
		)
	case witigo.AbiTypeResult:
		okType := baseType.SubTypes()[0].Type().CodegenGolangTypename()
		errType := baseType.SubTypes()[1].Type().CodegenGolangTypename()
		return generator.NewRoot(
			generator.NewInterface(
				w.CodegenGolangTypename(),
				generator.NewFuncSignature("Ok").AddReturnTypes(okType),
				generator.NewFuncSignature("Error").AddReturnTypes(errType),
			),
		)
	case witigo.AbiTypeTuple:
		subTypes := baseType.SubTypes()
		fields := make([]*generator.FuncSignature, len(subTypes))
		for i, subType := range subTypes {
			fields[i] = generator.
				NewFuncSignature(textcase.PascalCase(fmt.Sprintf("Field%d", i+1))).
				AddReturnTypes(subType.Type().CodegenGolangTypename())
		}
		return generator.NewRoot(
			generator.NewInterface(
				w.CodegenGolangTypename(),
				fields...,
			),
		)
	case witigo.AbiTypeEnum:
		root := generator.NewRoot().
			AddStatements(
				generator.NewRawStatement(fmt.Sprintf("type %s int", w.CodegenGolangTypename())))

		for i, c := range baseType.SubTypes() {
			root = root.AddStatements(
				generator.NewRawStatement(
					fmt.Sprintf(
						"const %s = %d",
						w.CodegenGolangTypename()+textcase.PascalCase(c.Name()),
						i,
					),
				),
			)
		}
		return root
	default:
		return nil
	}
}
