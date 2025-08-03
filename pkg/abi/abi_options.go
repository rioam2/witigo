package abi

import (
	"context"

	"github.com/tetratelabs/wazero/api"
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

type RuntimeMemory interface {
	Size() uint32
	ReadByte(offset uint32) (byte, bool)
	ReadUint16Le(offset uint32) (uint16, bool)
	ReadUint32Le(offset uint32) (uint32, bool)
	ReadFloat32Le(offset uint32) (float32, bool)
	ReadUint64Le(offset uint32) (uint64, bool)
	ReadFloat64Le(offset uint32) (float64, bool)
	Read(offset, byteCount uint32) ([]byte, bool)
	WriteByte(offset uint32, v byte) bool
	WriteUint16Le(offset uint32, v uint16) bool
	WriteUint32Le(offset, v uint32) bool
	WriteFloat32Le(offset uint32, v float32) bool
	WriteUint64Le(offset uint32, v uint64) bool
	WriteFloat64Le(offset uint32, v float64) bool
	Write(offset uint32, v []byte) bool
}

type AbiOptions struct {
	StringEncoding StringEncoding
	Context        context.Context
	Memory         RuntimeMemory
	Func           func(name string) api.Function
}
