package abi

import (
	"fmt"
)

func Call(opts AbiOptions, name string, params ...uint64) (ret uint32, postReturn func() error, err error) {
	if opts.Func == nil {
		return 0, nil, fmt.Errorf("module is not set in AbiOptions")
	}
	fn := opts.Func(name)
	if fn == nil {
		return 0, nil, fmt.Errorf("function %s not found in module", name)
	}
	results, err := fn.Call(opts.Context, params...)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to call function %s: %w", name, err)
	}
	if len(results) > 0 {
		ret = uint32(results[0])
	}
	postReturn = func() error {
		postFn := opts.Func("cabi_post_" + name)
		if postFn != nil {
			_, err := postFn.Call(opts.Context, uint64(ret))
			return err
		}
		return nil
	}
	return ret, postReturn, err
}
