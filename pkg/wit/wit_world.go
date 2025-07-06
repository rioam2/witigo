package wit

import "encoding/json"

type WitWorldDefinition interface {
	Name() string
	ExportedFunctions() []WitFunction
	String() string
}

type WitWorldDefinitionImpl struct {
	Raw  json.RawMessage
	Root WitDefinition
}

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

func (w *WitWorldDefinitionImpl) String() string {
	base := "world " + w.Name() + " {"
	for _, function := range w.ExportedFunctions() {
		base += "\n  export " + function.String()
	}
	base += "\n}"
	return base
}
