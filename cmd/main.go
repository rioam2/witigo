package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <command>\n", os.Args[0])
		os.Exit(1)
	}

	command := os.Args[1]
	switch command {
	case "hello":
		fmt.Println("Hello, World!")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
