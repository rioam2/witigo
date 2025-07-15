package main

import (
	"fmt"
	"os"

	"github.com/moznion/gowrtr/generator"
	"github.com/rioam2/witigo/pkg/wit"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

	// Extract JSON representation of WebAssembly Interface Types (WIT) from the given file
	inputFile := os.Args[1]
	wit, err := wit.NewFromFile(inputFile)
	if err != nil {
		fmt.Printf("Error extracting WIT: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(wit.String())
	fmt.Println("---")
	fmt.Println("Types:")
	for idx, t := range wit.Types() {
		fmt.Printf("  %d: %s: %s\n", idx, t.Name(), t.String())
	}

	fmt.Println("---")

	for _, t := range wit.Types() {
		typeGenerator := t.CodegenGolangTypedef()
		if typeGenerator == nil {
			continue
		}
		codeGenerator := generator.NewRoot(typeGenerator)
		code, err := codeGenerator.Generate(0)
		if err != nil {
			fmt.Printf("Error generating code for type %s: %v\n", t.Name(), err)
			continue
		}
		fmt.Println(code)
	}

	fmt.Println()

	for _, f := range wit.Worlds()[0].ExportedFunctions() {
		funcGenerator := f.Codegen()
		codeGenerator := generator.NewRoot(funcGenerator)
		code, err := codeGenerator.Generate(0)
		if err != nil {
			fmt.Printf("Error generating code for function %s: %v\n", f.Name(), err)
			continue
		}
		fmt.Println(code)
	}
}
