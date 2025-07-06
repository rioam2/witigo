package wit

import (
	"encoding/json"
)

type WitFunction interface {
	Name() string
	Params() []WitTypeReference
	Returns() WitType
	String() string
}

type WitFunctionImpl struct {
	Raw  json.RawMessage
	Root WitDefinition
}

func (w *WitFunctionImpl) Name() string {
	var data struct {
		Name string `json:"name"`
	}
	json.Unmarshal(w.Raw, &data)
	return data.Name
}

func (w *WitFunctionImpl) Params() []WitTypeReference {
	var data struct {
		Params []json.RawMessage `json:"params"`
	}
	json.Unmarshal(w.Raw, &data)
	var params []WitTypeReference
	for _, param := range data.Params {
		params = append(params, &WitTypeReferenceImpl{Raw: param, Root: w.Root})
	}
	return params
}

func (w *WitFunctionImpl) Returns() WitType {
	var data struct {
		Result any `json:"result"`
	}
	json.Unmarshal(w.Raw, &data)
	switch t := data.Result.(type) {
	case string:
		return &WitTypeImpl{[]byte("{\"type\":\"" + t + "\"}"), w.Root}
	case float64:
		return w.Root.Types()[int(t)]
	default:
		return nil
	}
}

func (w *WitFunctionImpl) String() string {
	params := ""
	for idx, param := range w.Params() {
		if idx > 0 {
			params += ", "
		}
		params += param.Name() + ": " + param.String()
	}
	return w.Name() + ": func (" + params + ") -> " + w.Returns().String()
}
