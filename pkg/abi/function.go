package abi

import (
	"fmt"
)

// Call invokes the function specified by name in the provided WASM module with the given parameters.
// It returns the result of the function call and a post-return function to handle memory cleanup.
func Call(opts AbiOptions, name string, params ...uint64) (ret uint32, postReturn func() error, err error) {
	if opts.Call == nil {
		return 0, nil, fmt.Errorf("call function is not defined in AbiOptions")
	}
	results, err := opts.Call(opts.Context, name, params...)
	if err != nil {
		return 0, nil, fmt.Errorf("function call %s failed: %w", name, err)
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
