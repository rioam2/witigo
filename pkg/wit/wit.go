package wit

import (
	"encoding/json"
)

type WitDefinition interface {
	Name() string
	Worlds() []WitWorldDefinition
	Types() []WitType
	String() string
}

type WitDefinitionImpl struct {
	name string
	Raw  json.RawMessage
}

var _ WitDefinition = &WitDefinitionImpl{}

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

func (w *WitDefinitionImpl) String() string {
	var base string
	for _, world := range w.Worlds() {
		base += world.String() + "\n"
	}
	return base
}

func (w *WitDefinitionImpl) Name() string {
	return w.name
}

func NewFromJson(json json.RawMessage, name string) (WitDefinition, error) {

	return &WitDefinitionImpl{Raw: json, name: name}, nil
}
