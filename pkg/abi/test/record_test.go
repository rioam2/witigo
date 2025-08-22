package test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
	"github.com/stretchr/testify/assert"
)

func TestRead_ComplexStruct(t *testing.T) {
	type ComplexRecord struct {
		ID      int32
		Active  bool
		Score   float64
		Message string
	}

	data := map[uint64][]byte{
		0x00: {1, 0, 0, 0},              // ID field
		0x04: {1},                       // Active field
		0x08: {0, 0, 0, 0, 0, 0, 0, 64}, // Score field (double: 2.0)
		0x10: {0x18, 0, 0, 0},           // Message field: String pointer
		0x14: {5, 0, 0, 0},              // Message field: String length
		0x18: []byte("world"),
	}
	opts := createAbiOptionsFromMemoryMap(data)
	record := &ComplexRecord{}
	err := abi.Read(opts, 0x00, record)

	assert.NoError(t, err)
	assert.Equal(t, int32(1), record.ID)
	assert.Equal(t, true, record.Active)
	assert.Equal(t, 2.0, record.Score)
	assert.Equal(t, "world", record.Message)
}

func TestRead_NestedStruct(t *testing.T) {
	type DemographicRecord struct {
		Age  int32
		City string
	}

	type NestedRecord struct {
		ID      int32
		Details DemographicRecord
	}

	data := map[uint64][]byte{
		0x00: {1, 0, 0, 0},    // ID field
		0x04: {25, 0, 0, 0},   // Details.Age field
		0x08: {0x10, 0, 0, 0}, // Details.City field: String pointer
		0x0C: {4, 0, 0, 0},    // Details.City field: String length
		0x10: []byte("Rome"),
	}
	opts := createAbiOptionsFromMemoryMap(data)
	record := &NestedRecord{}
	err := abi.Read(opts, 0x00, record)

	assert.NoError(t, err)
	assert.Equal(t, int32(1), record.ID)
	assert.Equal(t, int32(25), record.Details.Age)
	assert.Equal(t, "Rome", record.Details.City)
}

// func TestRead_ArrayField(t *testing.T) {
// 	type ArrayRecord struct {
// 		Values [3]int32
// 	}

// 	data := map[uint32][]byte{
// 		0x00: {1, 0, 0, 0}, // Values[0]
// 		0x04: {2, 0, 0, 0}, // Values[1]
// 		0x08: {3, 0, 0, 0}, // Values[2]
// 	}
// 	memory := createMemoryFromMap(data)
// 	opts := abi.AbiOptions{Memory: memory, StringEncoding: abi.StringEncodingUTF8}
// 	record := &ArrayRecord{}
// 	err := abi.Read(opts, 0x00, record)

// 	assert.NoError(t, err)
// 	assert.Equal(t, [3]int32{1, 2, 3}, record.Values)
// }

func TestRead_StructWithPadding(t *testing.T) {
	type PaddedRecord struct {
		ID   int32
		Flag bool
	}

	data := map[uint64][]byte{
		0x00: {1, 0, 0, 0}, // ID field
		0x04: {1},          // Flag field
	}
	opts := createAbiOptionsFromMemoryMap(data)
	record := &PaddedRecord{}
	err := abi.Read(opts, 0x00, record)

	assert.NoError(t, err)
	assert.Equal(t, int32(1), record.ID)
	assert.Equal(t, true, record.Flag)
}

func TestWriteThenRead_StructWithPadding(t *testing.T) {
	type PaddedRecord struct {
		ID   int32
		Flag bool
	}

	opts := createAbiOptionsFromMemoryMap(nil)
	record := &PaddedRecord{ID: 42, Flag: true}

	ptr, free, err := abi.Write(opts, record, nil)
	assert.NoError(t, err)
	defer free()

	readRecord := &PaddedRecord{}
	err = abi.Read(opts, ptr, readRecord)
	assert.NoError(t, err)

	assert.Equal(t, int32(42), readRecord.ID)
	assert.Equal(t, true, readRecord.Flag)
}
