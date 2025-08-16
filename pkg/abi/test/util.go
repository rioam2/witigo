package test

import (
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
	if int(ptr)+len(data) > len(m.bytes) {
		return false
	}
	copy(m.bytes[ptr:], data)
	return true
}

func (m *FakeMemory) WriteUint32Le(ptr uint32, value uint32) bool {
	if int(ptr)+4 > len(m.bytes) {
		return false
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
