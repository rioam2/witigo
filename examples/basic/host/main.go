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

	fmt.Printf("Result of StringFunc: %+v\n", instance.StringFunc("Hello, World!"))
	fmt.Printf("Result of RecordFunc: %+v\n", instance.RecordFunc(basic_example_component.CustomerRecord{Name: "John Doe", Age: 30}))
	fmt.Printf("Result of ListFunc: %+v\n", instance.ListFunc([]uint64{1, 2, 3, 4, 5}))
	fmt.Printf("Result of OptionFunc: %+v\n", instance.OptionFunc(nil))
	fmt.Printf("Result of ResultFunc: %+v\n", instance.ResultFunc(basic_example_component.Uint64StringResult{Ok: 42, Error: "Success"}))
	fmt.Printf("Result of VariantFunc: %+v\n", instance.VariantFunc(basic_example_component.AllowedDestinationsVariant{Type: basic_example_component.AllowedDestinationsVariantTypeAny}))
	fmt.Printf("Result of EnumFunc: %+v\n", instance.EnumFunc(basic_example_component.ColorEnumLimeGreen))
	fmt.Printf("Result of TupleFunc: %+v\n", instance.TupleFunc(basic_example_component.StringUint32Tuple{Elem0: "example", Elem1: 12345}))
}
