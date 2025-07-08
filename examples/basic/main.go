package main

import (
	"fmt"
	"os"

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
}
