package abi

import (
	"fmt"
)

// Call invokes the function specified by name in the provided WASM module with the given parameters.
// It returns the result of the function call and a post-return function to handle memory cleanup.
func Call(opts AbiOptions, name string, params ...uint64) (ret uint32, postReturn AbiFreeCallback, err error) {
	if opts.Call == nil {
		return 0, AbiFreeCallbackNoop, fmt.Errorf("call function is not defined in AbiOptions")
	}
	results, err := opts.Call(opts.Context, name, params...)
	if err != nil {
		return 0, AbiFreeCallbackNoop, fmt.Errorf("function call %s failed: %w", name, err)
	}
	if len(results) > 0 {
		ret = uint32(results[0])
	}
	postReturn = func() error {
		_, err := opts.Call(opts.Context, "cabi_post_"+name, uint64(ret))
		return err
	}
	return ret, postReturn, err
}

// realloc reallocates memory at the specified pointer with the given size and alignment.
func realloc(opts AbiOptions, oldPtr uint32, oldSize uint32, alignment uint32, newSize uint32) (ptr uint32, free AbiFreeCallback, err error) {
	return Call(opts, "cabi_realloc", uint64(oldPtr), uint64(oldSize), uint64(alignment), uint64(newSize))
}

// free releases memory at the specified pointer.
func free(opts AbiOptions, ptr uint32) error {
	_, _, err := realloc(opts, ptr, 0, 0, 0)
	return err
}

// malloc allocates memory of the specified size and alignment.
func malloc(opts AbiOptions, size uint32, alignment uint32) (ptr uint32, free AbiFreeCallback, err error) {
	return realloc(opts, 0, 0, alignment, size)
}
