package test

import (
	"context"

	"github.com/rioam2/witigo/pkg/abi"
)

type FakeMemory struct {
	bytes []byte
}

func (m *FakeMemory) Read(ptr uint32, size uint32) ([]byte, bool) {
	if int(ptr)+int(size) > len(m.bytes) {
		return nil, false
	}
	return m.bytes[ptr : ptr+size], true
}

func (m *FakeMemory) ReadUint32Le(ptr uint32) (uint32, bool) {
	if int(ptr)+4 > len(m.bytes) {
		return 0, false
	}
	return uint32(m.bytes[ptr]) | uint32(m.bytes[ptr+1])<<8 | uint32(m.bytes[ptr+2])<<16 | uint32(m.bytes[ptr+3])<<24, true
}

func (m *FakeMemory) Write(ptr uint32, data []byte) bool {
	diff := int(ptr) + len(data) - len(m.bytes)
	if diff > 0 {
		m.bytes = append(m.bytes, make([]byte, diff)...)
	}
	copy(m.bytes[ptr:], data)
	return true
}

func (m *FakeMemory) WriteUint32Le(ptr uint32, value uint32) bool {
	diff := int(ptr) + 4 - len(m.bytes)
	if diff > 0 {
		m.bytes = append(m.bytes, make([]byte, diff)...)
	}
	m.bytes[ptr] = byte(value)
	m.bytes[ptr+1] = byte(value >> 8)
	m.bytes[ptr+2] = byte(value >> 16)
	m.bytes[ptr+3] = byte(value >> 24)
	return true
}

func (m *FakeMemory) Size() uint32 {
	return uint32(len(m.bytes))
}

func createMemory(bytes []byte) abi.RuntimeMemory {
	return &FakeMemory{
		bytes: bytes,
	}
}

func createMemoryFromMap(data map[uint32][]byte) abi.RuntimeMemory {
	capacity := 0
	for startAddr, b := range data {
		maxAddr := int(startAddr) + int(len(b))
		if maxAddr > capacity {
			capacity = maxAddr
		}
	}
	bytes := make([]byte, capacity)
	for startAddr, b := range data {
		copy(bytes[startAddr:], b)
	}
	return createMemory(bytes)
}

func createAbiOptionsFromMemoryMap(data map[uint32][]byte) abi.AbiOptions {
	allocPtr := uint64(0x10000)
	allocIncr := uint64(0x1000)
	call := func(ctx context.Context, name string, params ...uint64) ([]uint64, error) {
		if name == "cabi_realloc" {
			allocPtr += allocIncr
			return []uint64{allocPtr}, nil
		}
		return []uint64{0}, nil
	}

	mem := createMemoryFromMap(data)
	return abi.AbiOptions{
		Memory:         mem,
		StringEncoding: abi.StringEncodingUTF8, // Default encoding
		Call:           call,
		Context:        context.Background(),
	}
}
