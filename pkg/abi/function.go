package abi

import (
	"fmt"
)

func Call(opts AbiOptions, name string, params ...uint64) (ret uint32, err error) {
	if opts.Func == nil {
		return 0, fmt.Errorf("module is not set in AbiOptions")
	}
	fn := opts.Func(name)
	if fn == nil {
		return 0, fmt.Errorf("function %s not found in module", name)
	}
	results, err := fn.Call(opts.Context, params...)
	if err != nil {
		return 0, fmt.Errorf("failed to call function %s: %w", name, err)
	}
	if len(results) > 0 {
		ret = uint32(results[0])
	}
	postFn := opts.Func("cabi_post_" + name)
	if postFn != nil {
		defer func() {
			_, postErr := postFn.Call(opts.Context, uint64(ret))
			if postErr != nil {
				ret = 0
				err = fmt.Errorf("post function call failed: %w", postErr)
			}
		}()
	}
	return ret, err
}
