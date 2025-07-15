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

func NewFromFile(filePath string) (WitDefinition, error) {
	baseName := filePath
	if idx := bytes.LastIndexByte([]byte(filePath), '/'); idx != -1 {
		baseName = filePath[idx+1:]
	}
	name := baseName[:len(baseName)-(len(".wit")+1)]
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

	witJson := stdoutBuffer.String()
	return &WitDefinitionImpl{Raw: json.RawMessage(witJson), name: name}, nil
}
