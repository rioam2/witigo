package wit

import (
	"encoding/json"
	"fmt"
	"math"

	witigo "github.com/rioam2/witigo/pkg"
)

type WitType interface {
	Name() string
	Kind() witigo.AbiType
	Owner() *string
	String() string
	SubType() WitTypeReference
	SubTypes() []WitTypeReference
	IsPrimitive() bool
	IsExported() bool
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
	var dataVer1 struct {
		Kind *string `json:"kind"`
	}
	err := json.Unmarshal(w.Raw, &dataVer1)
	if err == nil && dataVer1.Kind != nil {
		switch *dataVer1.Kind {
		case "resource":
			return witigo.AbiTypeResource
		}
	}

	var dataVer2 struct {
		Kind *struct {
			List    *json.RawMessage `json:"list"`
			Option  *json.RawMessage `json:"option"`
			Record  *json.RawMessage `json:"record"`
			Tuple   *json.RawMessage `json:"tuple"`
			Result  *json.RawMessage `json:"result"`
			Variant *json.RawMessage `json:"variant"`
			Enum    *json.RawMessage `json:"enum"`
			Type    *json.RawMessage `json:"type"`
			Handle  *json.RawMessage `json:"handle"`
		} `json:"kind"`
		Type *string `json:"type"`
	}
	json.Unmarshal(w.Raw, &dataVer2)
	if dataVer2.Type != nil {
		switch *dataVer2.Type {
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
	if dataVer2.Kind.List != nil {
		return witigo.AbiTypeList
	}
	if dataVer2.Kind.Option != nil {
		return witigo.AbiTypeOption
	}
	if dataVer2.Kind.Record != nil {
		return witigo.AbiTypeRecord
	}
	if dataVer2.Kind.Tuple != nil {
		return witigo.AbiTypeTuple
	}
	if dataVer2.Kind.Result != nil {
		return witigo.AbiTypeResult
	}
	if dataVer2.Kind.Variant != nil {
		return witigo.AbiTypeVariant
	}
	if dataVer2.Kind.Enum != nil {
		return witigo.AbiTypeEnum
	}
	if dataVer2.Kind.Type != nil {
		rawTypeRef, err := json.Marshal(map[string]any{
			"type": dataVer2.Kind.Type,
			"name": w.Name(),
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal type reference: %v", err))
		}
		ref := &WitTypeReferenceImpl{Raw: rawTypeRef, Root: w.Root}
		return ref.Type().Kind()
	}
	if dataVer2.Kind.Handle != nil {
		return witigo.AbiTypeHandle
	}
	panic(fmt.Sprintf("Unknown WIT type kind: %v", dataVer2.Kind))
}

func (w *WitTypeImpl) Owner() *string {
	var data struct {
		Owner *struct {
			Interface *float64 `json:"interface"`
			World     *float64 `json:"world"`
		} `json:"owner"`
	}
	err := json.Unmarshal(w.Raw, &data)
	if err != nil || data.Owner == nil {
		return nil
	}
	if data.Owner.World != nil {
		name := w.Root.Worlds()[int(*data.Owner.World)].Name()
		return &name
	}
	if data.Owner.Interface != nil {
		name := fmt.Sprintf("interface %d", int(*data.Owner.Interface))
		return &name
	}
	return nil
}

func (w *WitTypeImpl) SubType() WitTypeReference {
	var data struct {
		Kind struct {
			List   *any `json:"list"`
			Option *any `json:"option"`
			Handle *struct {
				Own    *json.RawMessage `json:"own"`
				Borrow *json.RawMessage `json:"borrow"`
			}
		} `json:"kind"`
	}
	err := json.Unmarshal(w.Raw, &data)
	if err != nil {
		return nil
	}
	var subTypeRef any = nil
	if data.Kind.List != nil {
		subTypeRef = data.Kind.List
	} else if data.Kind.Option != nil {
		subTypeRef = data.Kind.Option
	} else if data.Kind.Handle != nil {
		if data.Kind.Handle.Own != nil {
			subTypeRef = data.Kind.Handle.Own
		} else if data.Kind.Handle.Borrow != nil {
			subTypeRef = data.Kind.Handle.Borrow
		}
	}
	rawTypeRef, err := json.Marshal(map[string]any{
		"type": subTypeRef,
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
	err := json.Unmarshal(w.Raw, &data)
	if err != nil {
		return nil
	}
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

func discriminantSize(n int) int {
	bits := math.Ceil(math.Log2(float64(n))/8) * 8
	return int(bits)
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
	discriminantType := fmt.Sprintf("u%d", discriminantSize(len(enum.Cases)))
	for _, c := range enum.Cases {
		remappedType, err := json.Marshal(map[string]any{
			"type": discriminantType,
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
		witigo.AbiTypeF32, witigo.AbiTypeF64, witigo.AbiTypeChar,
		witigo.AbiTypeResource:
		return true
	default:
		return false
	}
}

func (w *WitTypeImpl) IsExported() bool {
	for _, world := range w.Root.Worlds() {
		for _, export := range world.ExportedFunctions() {
			if export.Returns().String() == w.String() {
				return true
			}
			for _, arg := range export.Params() {
				if arg.String() == w.String() {
					return true
				}
			}
		}
	}
	return false
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
