package wit

import (
	"encoding/json"
	"fmt"

	witigo "github.com/rioam2/witigo/pkg"
)

type WitType interface {
	Name() string
	Kind() witigo.AbiType
	String() string
	RecordFields() []WitTypeReference
	ListType() WitTypeReference
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

func (w *WitTypeImpl) String() string {
	base := w.Kind().String()
	switch w.Kind() {
	case witigo.AbiTypeRecord:
		fields := w.RecordFields()
		if len(fields) > 0 {
			base += "{ "
			for i, field := range fields {
				if i > 0 {
					base += ", "
				}
				base += field.Name() + ":" + field.String()
			}
			base += " }"
		}
	case witigo.AbiTypeList:
		listType := w.ListType()
		if listType != nil {
			base += "<" + listType.String() + ">"
		}
	}
	return base

}

func (w *WitTypeImpl) RecordFields() []WitTypeReference {
	var data struct {
		Kind struct {
			Record struct {
				Fields []json.RawMessage `json:"fields"`
			} `json:"record"`
		} `json:"kind"`
	}
	json.Unmarshal(w.Raw, &data)
	var fields []WitTypeReference
	for _, field := range data.Kind.Record.Fields {
		fields = append(fields, &WitTypeReferenceImpl{Raw: field, Root: w.Root})
	}
	return fields
}

func (w *WitTypeImpl) ListType() WitTypeReference {
	var data struct {
		Kind struct {
			List *any `json:"list"`
		} `json:"kind"`
	}
	json.Unmarshal(w.Raw, &data)
	if data.Kind.List == nil {
		return nil
	}
	var typeRef struct {
		Name *string `json:"name"`
		Type any     `json:"type"`
	}
	typeRef.Type = *data.Kind.List
	typeRef.Name = nil
	rawTypeRef, err := json.Marshal(typeRef)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal list type reference: %v", err))
	}
	return &WitTypeReferenceImpl{Raw: rawTypeRef, Root: w.Root}
}
