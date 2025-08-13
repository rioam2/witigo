package wit

import (
	"encoding/json"
)

type WitWorldDefinition interface {
	Name() string
	ExportedFunctions() []WitFunction
	Types() []WitType
	String() string
	ReferencesType(w WitType) bool
}

type WitWorldDefinitionImpl struct {
	Raw  json.RawMessage
	Root WitDefinition
}

var _ WitWorldDefinition = &WitWorldDefinitionImpl{}

func (w *WitWorldDefinitionImpl) Name() string {
	var data struct {
		Name string `json:"name"`
	}
	json.Unmarshal(w.Raw, &data)
	return data.Name
}

func (w *WitWorldDefinitionImpl) ExportedFunctions() []WitFunction {
	var data struct {
		Exports map[string]struct {
			Function *json.RawMessage `json:"function"`
		} `json:"exports"`
	}
	json.Unmarshal(w.Raw, &data)
	var functions []WitFunction
	for _, export := range data.Exports {
		if export.Function == nil {
			continue
		}
		functions = append(functions, &WitFunctionImpl{*export.Function, w.Root})
	}
	return functions
}

func (w *WitWorldDefinitionImpl) Types() []WitType {
	types := make([]WitType, 0)
	for _, t := range w.Root.Types() {
		for _, function := range w.ExportedFunctions() {
			if function.ReferencesType(t) {
				types = append(types, t)
				break
			}
		}
	}
	return types
}

func (w *WitWorldDefinitionImpl) String() string {
	base := "world " + w.Name() + " {"
	for _, function := range w.ExportedFunctions() {
		base += "\n  export " + function.String()
	}
	base += "\n}"
	return base
}

func (w *WitWorldDefinitionImpl) ReferencesType(t WitType) bool {
	for _, function := range w.ExportedFunctions() {
		if function.ReferencesType(t) {
			return true
		}
	}
	return false
}
