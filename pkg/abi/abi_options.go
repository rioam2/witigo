package abi

import (
	"context"
)

type StringEncoding string

const (
	StringEncodingUTF8  StringEncoding = "utf8"
	StringEncodingUTF16 StringEncoding = "utf16"
)

func (e StringEncoding) CodeUnitSize() uint32 {
	switch e {
	case StringEncodingUTF8:
		return 1
	case StringEncodingUTF16:
		return 2
	default:
		return 1
	}
}

func (e StringEncoding) Alignment() uint32 {
	switch e {
	case StringEncodingUTF8:
		return 1
	case StringEncodingUTF16:
		return 2
	default:
		return 1
	}
}

type RuntimeCall func(ctx context.Context, name string, params ...uint32) ([]uint32, error)

type RuntimeMemory interface {
	Size() uint32
	Read(offset, byteCount uint32) ([]byte, bool)
	ReadUint32Le(offset uint32) (uint32, bool)
	Write(offset uint32, v []byte) bool
	WriteUint32Le(offset, v uint32) bool
}

type AbiOptions struct {
	StringEncoding StringEncoding
	Memory         RuntimeMemory
	Call           RuntimeCall
	Context        context.Context
}
