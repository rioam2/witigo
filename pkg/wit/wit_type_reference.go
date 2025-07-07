package wit

import (
	"encoding/json"
)

type WitTypeReference interface {
	Name() string
	Type() WitType
	String() string
}

type WitTypeReferenceImpl struct {
	Raw  json.RawMessage
	Root WitDefinition
}

var _ WitTypeReference = &WitTypeReferenceImpl{}

func (w *WitTypeReferenceImpl) Name() string {
	var data struct {
		Name *string `json:"name"`
	}
	json.Unmarshal(w.Raw, &data)
	if data.Name == nil {
		return "(none)"
	}
	return *data.Name
}

func (w *WitTypeReferenceImpl) Type() WitType {
	var data struct {
		Type *any `json:"type"`
	}
	json.Unmarshal(w.Raw, &data)
	if data.Type == nil {
		// Default to u32 if no type is specified
		remappedType, err := json.Marshal(map[string]any{
			"type": "u32",
			"name": w.Name(),
		})
		if err != nil {
			return nil
		}
		return &WitTypeImpl{Raw: remappedType, Root: w.Root}
	}
	switch t := (*data.Type).(type) {
	case string:
		return &WitTypeImpl{w.Raw, w.Root}
	case float64:
		return w.Root.Types()[int(t)]
	default:
		return nil
	}
}

func (w *WitTypeReferenceImpl) String() string {
	return w.Type().String()
}
