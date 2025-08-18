package test

import (
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
	"github.com/stretchr/testify/assert"
)

type Test1Record struct {
	ID   int
	Name string
}

func TestRead_ValidData(t *testing.T) {
	data := map[uint32][]byte{0x00: {1, 0, 0, 0, 'T', 'e', 's', 't'}}
	memory := createMemoryFromMap(data)
	opts := abi.AbiOptions{Memory: memory}
	record := &Test1Record{}
	err := abi.Read(opts, 0x00, record)

	assert.NoError(t, err)
	assert.Equal(t, 1, record.ID)
	assert.Equal(t, "Test", record.Name)
}
