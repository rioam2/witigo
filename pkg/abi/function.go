package abi

import (
	"fmt"
)

// Call invokes the function specified by name in the provided WASM module with the given parameters.
// It returns the result of the function call and a post-return function to handle memory cleanup.
func Call(opts AbiOptions, name string, params ...uint64) (ret uint64, postReturn AbiFreeCallback, err error) {
	if opts.Call == nil {
		return 0, AbiFreeCallbackNoop, fmt.Errorf("call function is not defined in AbiOptions")
	}
	results, err := opts.Call(opts.Context, name, params...)
	if err != nil {
		return 0, AbiFreeCallbackNoop, fmt.Errorf("function call %s failed: %w", name, err)
	}
	if len(results) > 0 {
		ret = results[0]
	}
	postReturn = func() error {
		_, err := opts.Call(opts.Context, "cabi_post_"+name, ret)
		return err
	}
	return ret, postReturn, err
}

// realloc reallocates memory at the specified pointer with the given size and alignment.
func abi_realloc(opts AbiOptions, oldPtr uint64, oldSize uint64, alignment uint64, newSize uint64) (ptr uint64, free AbiFreeCallback, err error) {
	return Call(opts, "cabi_realloc", oldPtr, oldSize, alignment, newSize)
}

// free releases memory at the specified pointer.
func abi_free(opts AbiOptions, ptr uint64) error {
	_, _, err := abi_realloc(opts, ptr, 0, 0, 0)
	return err
}

// abi_malloc allocates memory of the specified size and alignment.
func abi_malloc(opts AbiOptions, size uint64, alignment uint64) (ptr uint64, free AbiFreeCallback, err error) {
	ptr, _, err = abi_realloc(opts, 0, 0, alignment, size)
	return ptr, func() error {
		return abi_free(opts, ptr)
	}, err
}
