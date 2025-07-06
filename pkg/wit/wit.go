package wit

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type WitDefinition interface {
	Worlds() []WitWorldDefinition
	Types() []WitType
}

type WitDefinitionImpl struct {
	Raw json.RawMessage
}

func (w *WitDefinitionImpl) Worlds() []WitWorldDefinition {
	var data struct {
		Worlds []json.RawMessage `json:"worlds"`
	}
	json.Unmarshal(w.Raw, &data)
	var worlds []WitWorldDefinition
	for _, world := range data.Worlds {
		worlds = append(worlds, &WitWorldDefinitionImpl{world, w})
	}
	return worlds
}

func (w *WitDefinitionImpl) Types() []WitType {
	var data struct {
		Types []json.RawMessage `json:"types"`
	}
	json.Unmarshal(w.Raw, &data)
	var types []WitType
	for _, t := range data.Types {
		types = append(types, &WitTypeImpl{t, w})
	}
	return types
}

func NewFromFile(filePath string) (WitDefinition, error) {
	cmd := exec.Command("wasm-tools", "component", "wit", "-j", "--all-features", filePath)
	witJsonStr, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error running wasm-tools: %w", err)
	}
	return &WitDefinitionImpl{Raw: json.RawMessage(witJsonStr)}, nil
}
