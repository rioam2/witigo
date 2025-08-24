package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	basic_example_component "github.com/rioam2/witigo/examples/basic/generated"
)

func checkErr(err error) {
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func toPrettyJson(obj interface{}) string {
	bytes, _ := json.MarshalIndent(obj, "", "  ")
	return string(bytes)
}

func main() {
	fmt.Printf("--- Basic Example ---\n\n")

	instance, err := basic_example_component.New(context.Background())
	checkErr(err)

	doubleOperation := basic_example_component.DoubleOperationRecord{
		DoubleList:   []float64{1.1, 2.2, 3.3},
		DoubleString: "Hello from Witigo! ",
	}

	fmt.Printf("Doubling input: \n%s\n\n", toPrettyJson(doubleOperation))

	doubleOperationResult, err := instance.Double(doubleOperation)
	checkErr(err)

	fmt.Printf("Doubled result: \n%s\n\n", toPrettyJson(doubleOperationResult))
	fmt.Println("---------------------")
}
