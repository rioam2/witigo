package abi

import (
	"context"
	"fmt"

	"github.com/tetratelabs/wazero/api"
)

// GetRuntimeCallFromWazero wraps a Wazero module's exported function call into a RuntimeCall.
func GetRuntimeCallFromWazero(module api.Module) RuntimeCall {
	return func(ctx context.Context, name string, params ...uint32) ([]uint32, error) {
		params64 := make([]uint64, len(params))
		for i, p := range params {
			params64[i] = uint64(p)
		}
		fn := module.ExportedFunction(name)
		if fn == nil {
			return nil, fmt.Errorf("function %s not found in module", name)
		}
		ret64, err := fn.Call(ctx, params64...)
		if err != nil {
			return nil, fmt.Errorf("call to function %s failed: %w", name, err)
		}
		ret := make([]uint32, len(ret64))
		for i, r := range ret64 {
			if r > uint64(^uint32(0)) {
				return nil, fmt.Errorf("function %s returned value %d which exceeds uint32 limit", name, r)
			}
			ret[i] = uint32(r)
		}
		return ret, nil
	}
}
