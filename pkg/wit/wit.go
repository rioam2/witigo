package wit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/rioam2/witigo/pkg/wasmtools"
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
	ctx := context.Background()
	runtime, err := wasmtools.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating wasmtools runtime: %w", err)
	}
	defer runtime.Close(ctx)

	// Run the wasm-tools command to get the WIT definition in JSON format.
	fsMap := map[string]fs.FS{".": os.DirFS(".")}
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}
	err = runtime.Run(ctx, os.Stdin, stdoutBuffer, stderrBuffer, fsMap, "component", "wit", "-j", "--all-features", filePath)
	if err != nil {
		return nil, fmt.Errorf("error running wasm-tools: %w\n%s", err, stderrBuffer.String())
	}

	return &WitDefinitionImpl{Raw: json.RawMessage(stdoutBuffer.String())}, nil
}
