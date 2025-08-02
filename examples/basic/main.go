package main

import (
	"fmt"
	"os"

	"github.com/rioam2/witigo/pkg/codegen"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <input> <outDir>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outDir := os.Args[2]

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	err := codegen.GenerateFromFile(inputFile, outDir)
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		os.Exit(1)
	}
}
