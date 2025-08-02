package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rioam2/witigo/pkg/codegen"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <input> <outDir>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Printf("Error resolving input file path: %v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Input file does not exist: %s\n", inputFile)
		os.Exit(1)
	}

	outDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		fmt.Printf("Error resolving output directory path: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			fmt.Printf("Error creating output directory: %v\n", err)
			os.Exit(1)
		}
	}

	err = codegen.GenerateFromFile(inputFile, outDir)
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		os.Exit(1)
	}
}
