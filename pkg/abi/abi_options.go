package abi

import (
	"context"
)

type StringEncoding string

const (
	StringEncodingUTF8  StringEncoding = "utf8"
	StringEncodingUTF16 StringEncoding = "utf16"
)

func (e StringEncoding) CodeUnitSize() uint64 {
	switch e {
	case StringEncodingUTF8:
		return 1
	case StringEncodingUTF16:
		return 2
	default:
		return 1
	}
}

func (e StringEncoding) Alignment() uint64 {
	switch e {
	case StringEncodingUTF8:
		return 1
	case StringEncodingUTF16:
		return 2
	default:
		return 1
	}
}

type RuntimeCall func(ctx context.Context, name string, params ...uint64) ([]uint64, error)

type RuntimeMemory interface {
	Size() uint64
	Read(offset, byteCount uint64) ([]byte, bool)
	ReadUint32Le(offset uint64) (uint32, bool)
	Write(offset uint64, v []byte) bool
	WriteUint32Le(offset uint64, v uint32) bool
}

type AbiOptions struct {
	StringEncoding StringEncoding
	Memory         RuntimeMemory
	Call           RuntimeCall
	Context        context.Context
}
