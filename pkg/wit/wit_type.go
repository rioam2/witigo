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
	SubType() WitTypeReference
	SubTypes() []WitTypeReference
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
	case witigo.AbiTypeList, witigo.AbiTypeOption:
		listType := w.SubType()
		if listType != nil {
			base += "<" + listType.String() + ">"
		}
	case witigo.AbiTypeRecord, witigo.AbiTypeVariant, witigo.AbiTypeEnum:
		cases := w.SubTypes()
		if len(cases) > 0 {
			base += "{ "
			for i, c := range cases {
				if i > 0 {
					base += ", "
				}
				base += c.Name() + ": " + c.String()
			}
			base += " }"
		}
	case witigo.AbiTypeTuple, witigo.AbiTypeResult:
		tupleTypes := w.SubTypes()
		if len(tupleTypes) > 0 {
			base += "<"
			for i, t := range tupleTypes {
				if i > 0 {
					base += ", "
				}
				base += t.String()
			}
			base += ">"
		}
	}
	return base
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
			Record *struct {
				Fields []json.RawMessage `json:"fields"`
			} `json:"record"`
			Variant *struct {
				Cases []json.RawMessage `json:"cases"`
			} `json:"variant"`
			Enum *struct {
				Cases []struct {
					Name string `json:"name"`
				} `json:"cases"`
			} `json:"enum"`
			Tuple *struct {
				Types []any `json:"types"`
			} `json:"tuple"`
			Result *struct {
				Ok  *any `json:"ok"`
				Err *any `json:"err"`
			} `json:"result"`
		} `json:"kind"`
	}
	json.Unmarshal(w.Raw, &data)
	var subTypes []WitTypeReference
	if data.Kind.Record != nil {
		for _, field := range data.Kind.Record.Fields {
			subTypes = append(subTypes, &WitTypeReferenceImpl{Raw: field, Root: w.Root})
		}
	} else if data.Kind.Variant != nil {
		for _, c := range data.Kind.Variant.Cases {
			subTypes = append(subTypes, &WitTypeReferenceImpl{Raw: c, Root: w.Root})
		}
	} else if data.Kind.Enum != nil {
		for _, c := range data.Kind.Enum.Cases {
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
	} else if data.Kind.Tuple != nil {
		for _, t := range data.Kind.Tuple.Types {
			remappedType, err := json.Marshal(map[string]any{
				"type": t,
				"name": nil,
			})
			if err != nil {
				panic(fmt.Sprintf("Failed to marshal tuple type reference: %v", err))
			}
			subTypes = append(subTypes, &WitTypeReferenceImpl{Raw: remappedType, Root: w.Root})
		}
	} else if data.Kind.Result != nil {
		okType, err := json.Marshal(map[string]any{
			"type": data.Kind.Result.Ok,
			"name": "ok",
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal result ok type reference: %v", err))
		}
		errType, err := json.Marshal(map[string]any{
			"type": data.Kind.Result.Err,
			"name": "error",
		})
		if err != nil {
			panic(fmt.Sprintf("Failed to marshal result err type reference: %v", err))
		}
		subTypes = append(subTypes, &WitTypeReferenceImpl{Raw: okType, Root: w.Root})
		subTypes = append(subTypes, &WitTypeReferenceImpl{Raw: errType, Root: w.Root})
	}
	return subTypes
}
