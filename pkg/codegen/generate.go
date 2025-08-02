package codegen

import (
	"fmt"
	"os"

	"github.com/rioam2/witigo/pkg/wasmtools"
	"github.com/rioam2/witigo/pkg/wit"
)

func GenerateFromFile(componentPath string, outDir string) error {
	componentWitJson, componentName, err := wasmtools.ExtractComponentWitJson(componentPath)
	if err != nil {
		return err
	}

	witDefinition, err := wit.NewFromJson(componentWitJson, componentName)
	if err != nil {
		return err
	}

	codeGen := GenerateFromWorld(witDefinition.Worlds()[0], witDefinition.Name())
	code, err := codeGen.EnableSyntaxChecking().Gofmt().Generate(0)
	if err != nil {
		return err
	}

	outputFile := fmt.Sprintf("%s/%s.go", outDir, componentName)
	err = os.WriteFile(outputFile, []byte(code), 0644)
	if err != nil {
		return fmt.Errorf("error writing generated code to file %s: %w", outputFile, err)
	}
	fmt.Printf("Generated code written to %s\n", outputFile)

	coreModule, err := wasmtools.ExtractComponentCoreModule(componentPath)
	if err != nil {
		return fmt.Errorf("error extracting core module: %w", err)
	}

	outputCoreModuleFile := fmt.Sprintf("%s/%s_core.wasm", outDir, componentName)
	err = os.WriteFile(outputCoreModuleFile, coreModule, 0644)
	if err != nil {
		return fmt.Errorf("error writing core module to file %s: %w", outputCoreModuleFile, err)
	}
	fmt.Printf("Core module written to %s\n", outputCoreModuleFile)

	return nil
}
