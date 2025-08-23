package abi

import (
	"errors"
	"fmt"
	"reflect"
)

// MAX_FLAT_PARAMS is the canonical ABI-defined constant for the maximum number of “flat” parameters to a wasm function.
// Over this number the heap is used for transferring parameters.
// Reference: https://docs.wasmtime.dev/api/wasmtime_environ/component/constant.MAX_FLAT_PARAMS.html
const MAX_FLAT_PARAMS = 16

type Parameter struct {
	Value     uint64
	Size      uint64
	Alignment uint64
}

// WriteParameter writes a parameter value to memory and returns the arguments for the ABI call.
func WriteParameter(opts AbiOptions, value any) (params []Parameter, free AbiFreeCallback, err error) {
	// Validate input and retrieve element type of value
	rv := reflect.ValueOf(value)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return nil, free, errors.New("must pass a valid value")
	}

	// Write based on the kind of the value
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return WriteParameterInt(opts, value)
	case reflect.Bool:
		return WriteParameterBool(opts, value)
	case reflect.Float32, reflect.Float64:
		return WriteParameterFloat(opts, value)
	case reflect.String:
		return WriteParameterString(opts, value)
	case reflect.Slice:
		return WriteParameterList(opts, value)
	case reflect.Struct:
		structName := rv.Type().Name()
		if isStructRecordType(rv) {
			return WriteParameterRecord(opts, value)
		} else if isStructOptionType(rv) {
			return WriteParameterOption(opts, value)
		} else {
			return nil, AbiFreeCallbackNoop, fmt.Errorf("writing struct %s is not implemented", structName)
		}
	default:
		return nil, AbiFreeCallbackNoop, fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
}

func WriteIndirectParameters(opts AbiOptions, params ...Parameter) (ptr uint64, free AbiFreeCallback, err error) {
	// Initialize return values
	ptr = 0
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	paramListSize := uint64(0)
	paramListAlignment := uint64(1)
	for _, param := range params {
		paramListSize += AlignTo(paramListSize+param.Size, param.Alignment)
		if param.Alignment > paramListAlignment {
			paramListAlignment = param.Alignment
		}
	}

	// Allocate space for the indirect parameters in linear memory
	paramListPtr, paramListFree, err := abiMalloc(opts, paramListSize, paramListAlignment)
	if err != nil {
		return ptr, free, err
	}
	ptr = paramListPtr
	freeCallbacks = append(freeCallbacks, paramListFree)

	offset := uint64(0)
	for _, param := range params {
		offset = AlignTo(offset, param.Alignment)
		paramPtr := ptr + offset

		// Write the param to memory
		bytes := make([]byte, param.Size)
		for i := range param.Size {
			bytes[i] = byte(param.Value >> (8 * i))
		}
		if !opts.Memory.Write(paramPtr, bytes) {
			return ptr, free, fmt.Errorf("failed to write %d bytes at parameter pointer %d", param.Size, paramPtr)
		}

		offset += param.Size
	}

	return ptr, free, nil
}

func WriteParameters(opts AbiOptions, values ...any) (flatParams []uint64, free AbiFreeCallback, err error) {
	// Initialize return values
	flatParams = []uint64{}
	params := []Parameter{}
	freeCallbacks := []AbiFreeCallback{}
	free = wrapFreeCallbacks(&freeCallbacks)

	for i, value := range values {
		currentParams, freeCurrentParams, err := WriteParameter(opts, value)
		if err != nil {
			return flatParams, free, fmt.Errorf("failed to write parameter %d: %w", i, err)
		}
		params = append(params, currentParams...)
		freeCallbacks = append(freeCallbacks, freeCurrentParams)
	}
	if len(params) > MAX_FLAT_PARAMS {
		flatParam, flatParamFree, err := WriteIndirectParameters(opts, params...)
		if err != nil {
			return flatParams, free, fmt.Errorf("failed to write indirect parameters: %w", err)
		}
		freeCallbacks = append(freeCallbacks, flatParamFree)
		flatParams = append(flatParams, flatParam)
		return flatParams, free, nil
	} else {
		flatParams = make([]uint64, len(params))
		for i := range params {
			flatParams[i] = params[i].Value
		}
		return flatParams, free, nil
	}
}
