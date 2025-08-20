package main

import (
	"context"
	"fmt"
	"os"

	basic_example_component "github.com/rioam2/witigo/examples/basic/generated"
)

func main() {
	instance, err := basic_example_component.New(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating instance: %v\n", err)
		os.Exit(1)
	}

	stringFuncResult, err := instance.StringFunc("Hello, World!")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling StringFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of StringFunc: %+v\n", stringFuncResult)

	recordFuncResult, err := instance.RecordFunc(basic_example_component.CustomerRecord{Id: 1, Picture: &[]uint8{1, 2}, Name: "John Doe", Age: 30})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling RecordFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of RecordFunc: %+v\n", recordFuncResult)

	listFuncResult, err := instance.ListFunc([]uint64{1, 2, 3, 4, 5})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling ListFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of ListFunc: %+v\n", listFuncResult)

	optionFuncResult, err := instance.OptionFunc(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling OptionFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of OptionFunc: %+v\n", optionFuncResult)

	resultFuncResult, err := instance.ResultFunc(basic_example_component.Uint64StringResult{Ok: 42, Error: "Success"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling ResultFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of ResultFunc: %+v\n", resultFuncResult)

	variantFuncResult, err := instance.VariantFunc(basic_example_component.AllowedDestinationsVariant{Type: basic_example_component.AllowedDestinationsVariantTypeAny})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling VariantFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of VariantFunc: %+v\n", variantFuncResult)

	enumFuncResult, err := instance.EnumFunc(basic_example_component.ColorEnumLimeGreen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling EnumFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of EnumFunc: %+v\n", enumFuncResult)

	tupleFuncResult, err := instance.TupleFunc(basic_example_component.StringUint32Tuple{Elem0: "example", Elem1: 12345})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling TupleFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of TupleFunc: %+v\n", tupleFuncResult)
}
