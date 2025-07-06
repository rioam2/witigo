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
	Fields() []WitTypeReference
}

type WitTypeImpl struct {
	Raw  json.RawMessage
	Root WitDefinition
}

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
			List   *any             `json:"list"`
			Option *any             `json:"option"`
			Record *json.RawMessage `json:"record"`
			Tuple  *json.RawMessage `json:"tuple"`
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
	panic(fmt.Sprintf("Unknown WIT type kind: %v", data.Kind))
}

func (w *WitTypeImpl) String() string {
	base := w.Kind().String()
	switch w.Kind() {
	case witigo.AbiTypeRecord:
		fields := w.Fields()
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
	}
	return base

}

func (w *WitTypeImpl) Fields() []WitTypeReference {
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
