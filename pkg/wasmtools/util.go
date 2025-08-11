package wasmtools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ExtractComponentWitJson extracts the WIT definition in JSON format from a component file, as well as the name of the component.
func ExtractComponentWitJson(path string) (json.RawMessage, string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, "", fmt.Errorf("error getting absolute path: %w", err)
	}
	dirname := filepath.Dir(absolutePath)
	baseName := absolutePath
	if idx := bytes.LastIndexByte([]byte(absolutePath), '/'); idx != -1 {
		baseName = absolutePath[idx+1:]
	}
	name := baseName[:len(baseName)-(len(".wit")+1)]
	ctx := context.Background()
	runtime, err := New(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("error creating wasmtools runtime: %w", err)
	}
	defer runtime.Close(ctx)

	// Run the wasm-tools command to get the WIT definition in JSON format.
	fsMap := map[string]string{dirname: dirname}
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}
	err = runtime.Run(ctx, os.Stdin, stdoutBuffer, stderrBuffer, fsMap, "component", "wit", "-j", "--all-features", absolutePath)
	if err != nil {
		return nil, "", fmt.Errorf("error running wasm-tools: %w\n%s", err, stderrBuffer.String())
	}

	witJson := stdoutBuffer.String()
	return []byte(witJson), name, nil
}

// ExtractComponentCoreModule extracts a core WebAssembly module from a component file
func ExtractComponentCoreModule(path string) ([]byte, error) {
	componentAbsolutePath, err := filepath.Abs(path)
	componentDirname := filepath.Dir(componentAbsolutePath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %w", err)
	}
	ctx := context.Background()
	runtime, err := New(ctx)
	if err != nil {
		return nil, fmt.Errorf("error creating wasmtools runtime: %w", err)
	}
	defer runtime.Close(ctx)

	// Create a temporary directory for the core module output
	coreModuleDir, err := os.MkdirTemp("", "witigo-core-module-*")
	if err != nil {
		return nil, fmt.Errorf("error creating temporary directory for core module: %w", err)
	}

	fsMap := map[string]string{
		componentDirname: componentDirname,
		coreModuleDir:    coreModuleDir,
	}
	stdoutBuffer := &bytes.Buffer{}
	stderrBuffer := &bytes.Buffer{}
	err = runtime.Run(
		ctx, os.Stdin, stdoutBuffer, stderrBuffer, fsMap,
		"component", "unbundle",
		"--module-dir", coreModuleDir,
		componentAbsolutePath,
		"-t",
	)
	if err != nil {
		return nil, fmt.Errorf("error running wasm-tools: %w\n%s\n%s", err, stderrBuffer.String(), stdoutBuffer.String())
	}

	// Read the core module output from the temporary directory
	coreModuleFiles, err := os.ReadDir(coreModuleDir)
	if err != nil {
		return nil, fmt.Errorf("error reading core module directory: %w", err)
	}
	coreModuleFile := coreModuleFiles[0]
	coreModulePath := filepath.Join(coreModuleDir, coreModuleFile.Name())
	coreModuleData, err := os.ReadFile(coreModulePath)
	if err != nil {
		return nil, fmt.Errorf("error reading core module file: %w", err)
	}

	// Clean up the temporary directory
	if err := os.RemoveAll(coreModuleDir); err != nil {
		return nil, fmt.Errorf("error removing core module directory: %w", err)
	}

	// Return the core module data
	return coreModuleData, nil
}
