package main

import (
	"context"
	"fmt"
	"os"

	all_types_example_component "github.com/rioam2/witigo/examples/all-types/generated"
)

func main() {
	instance, err := all_types_example_component.New(context.Background())
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

	recordFuncResult, err := instance.RecordFunc(all_types_example_component.CustomerRecord{Id: 1, Picture: all_types_example_component.Option[[]uint8]{IsSome: true, Value: []uint8{1, 2}}, Name: "John Doe", Age: 30})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling RecordFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of RecordFunc: %+v\n", recordFuncResult)

	nestedRecordFuncResult, err := instance.NestedRecordFunc(all_types_example_component.NestedRecord{Level: 1, Color: all_types_example_component.ColorEnumNavyBlue, Customer: all_types_example_component.CustomerRecord{Id: 1, Picture: all_types_example_component.Option[[]uint8]{IsSome: true, Value: []uint8{1, 2}}, Name: "John Doe", Age: 30}})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling NestedRecordFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of NestedRecordFunc: %+v\n", nestedRecordFuncResult)

	simpleRecordFuncResult, err := instance.SimpleRecordFunc(all_types_example_component.SimpleRecordRecord{Id: 42})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling SimpleRecordFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of SimpleRecordFunc: %+v\n", simpleRecordFuncResult)

	listFuncResult, err := instance.ListFunc([]uint64{1, 2, 3, 4, 5})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling ListFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of ListFunc: %+v\n", listFuncResult)

	optionFuncResult, err := instance.OptionFunc(all_types_example_component.Option[uint64]{IsSome: true, Value: 42})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling OptionFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of OptionFunc: %+v\n", optionFuncResult)

	resultFuncResult, err := instance.ResultFunc(all_types_example_component.Uint64StringResult{Ok: 42, Error: "Success"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling ResultFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of ResultFunc: %+v\n", resultFuncResult)

	variantFuncResult, err := instance.VariantFunc(all_types_example_component.AllowedDestinationsVariant{Type: all_types_example_component.AllowedDestinationsVariantTypeAny})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling VariantFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of VariantFunc: %+v\n", variantFuncResult)

	enumFuncResult, err := instance.EnumFunc(all_types_example_component.ColorEnumLimeGreen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling EnumFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of EnumFunc: %+v\n", enumFuncResult)

	tupleFuncResult, err := instance.TupleFunc(all_types_example_component.StringUint32Tuple{Elem0: "example", Elem1: 12345})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error calling TupleFunc: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Result of TupleFunc: %+v\n", tupleFuncResult)
}
