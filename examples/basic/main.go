package main

import (
	"fmt"
	"os"

	"github.com/rioam2/witigo/pkg/codegen"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <file>\n", os.Args[0])
		os.Exit(1)
	}

	inputFile := os.Args[1]
	code, err := codegen.GenerateFromFile(inputFile)
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(code)
}
