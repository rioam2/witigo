package abi

import (
	"context"
	"fmt"

	"github.com/tetratelabs/wazero/api"
)

type WazeroMemory struct {
	module api.Module
}

func (m WazeroMemory) Size() uint64 {
	return uint64(m.module.Memory().Size())
}

func (m WazeroMemory) Read(offset, byteCount uint64) ([]byte, bool) {
	// Check if offset or byteCount exceed uint32 max
	if offset > uint64(^uint32(0)) || byteCount > uint64(^uint32(0)) {
		return nil, false
	}
	return m.module.Memory().Read(uint32(offset), uint32(byteCount))
}

func (m WazeroMemory) Write(offset uint64, data []byte) bool {
	// Check if offset exceeds uint32 max
	if offset > uint64(^uint32(0)) {
		return false
	}
	return m.module.Memory().Write(uint32(offset), data)
}

func (m WazeroMemory) ReadUint32Le(offset uint64) (uint32, bool) {
	// Check if offset exceeds uint32 max
	if offset > uint64(^uint32(0)) {
		return 0, false
	}
	return m.module.Memory().ReadUint32Le(uint32(offset))
}

func (m WazeroMemory) WriteUint32Le(offset uint64, value uint32) bool {
	// Check if offset exceeds uint32 max
	if offset > uint64(^uint32(0)) {
		return false
	}
	return m.module.Memory().WriteUint32Le(uint32(offset), value)
}

// GetRuntimeMemoryFromWazero converts a Wazero module into a RuntimeMemory.
func GetRuntimeMemoryFromWazero(module api.Module) RuntimeMemory {
	return WazeroMemory{module: module}
}

// GetRuntimeCallFromWazero wraps a Wazero module's exported function call into a RuntimeCall.
func GetRuntimeCallFromWazero(module api.Module) RuntimeCall {
	return func(ctx context.Context, name string, params ...uint64) ([]uint64, error) {
		fn := module.ExportedFunction(name)
		if fn == nil {
			return nil, fmt.Errorf("function %s not found in module", name)
		}
		return fn.Call(ctx, params...)
	}
}
