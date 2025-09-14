package abi_test

import (
	"math"
	"testing"

	"github.com/rioam2/witigo/pkg/abi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Synthetic variant mirroring code generation pattern
type SampleVariantType int

const (
	SampleVariantTypeA SampleVariantType = 0
	SampleVariantTypeB SampleVariantType = 1
	SampleVariantTypeC SampleVariantType = 2
)

type SampleVariant struct {
	Type SampleVariantType
	A    struct{}
	B    uint32
	C    string
}

func TestWriteParameterVariant(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	v := SampleVariant{Type: SampleVariantTypeB, B: 0xDEADBEEF}
	params, free, err := abi.WriteParameters(opts, v)
	require.NoError(t, err)
	defer free()
	// Flattened shape should be: discriminant + 2 slots (string case pointer,len) => 3 params
	require.Len(t, params, 3)
	assert.Equal(t, uint64(SampleVariantTypeB), params[0])
	assert.Equal(t, uint64(0xDEADBEEF), params[1]) // B payload in first payload slot
	// Third slot unused for this case => zero
	assert.Equal(t, uint64(0), params[2])
}

func TestVariantRoundTripEmptyCase(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	original := SampleVariant{Type: SampleVariantTypeA}
	ptr, free, err := abi.Write(opts, original, nil)
	require.NoError(t, err)
	defer free()
	var decoded SampleVariant
	err = abi.Read(opts, ptr, &decoded)
	require.NoError(t, err)
	assert.Equal(t, original.Type, decoded.Type)
}

func TestVariantRoundTripPayloadCases(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	cases := []SampleVariant{
		{Type: SampleVariantTypeB, B: 42},
		{Type: SampleVariantTypeC, C: "hello"},
		{Type: SampleVariantTypeC, C: ""},                                                     // empty string case
		{Type: SampleVariantTypeB, B: 0},                                                      // zero value case
		{Type: SampleVariantTypeB, B: 0xFFFFFFFF},                                             // max uint32
		{Type: SampleVariantTypeC, C: "A longer string to test string handling in variants."}, // long string
		{Type: SampleVariantTypeA},                                                            // empty case again
	}
	for _, c := range cases {
		ptr, free, err := abi.Write(opts, c, nil)
		require.NoError(t, err)
		var decoded SampleVariant
		err = abi.Read(opts, ptr, &decoded)
		require.NoError(t, err)
		assert.Equal(t, c.Type, decoded.Type)
		switch c.Type {
		case SampleVariantTypeB:
			assert.Equal(t, c.B, decoded.B)
		case SampleVariantTypeC:
			assert.Equal(t, c.C, decoded.C)
		}
		free()
	}
}

func TestVariantErrors(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	// Invalid discriminant when reading
	// Write bytes manually with out-of-range discriminant (e.g., 5)
	opts.Memory.Write(0, []byte{5, 0, 0, 0, 0, 0, 0, 0}) // discriminant stored as 8-byte int
	var v SampleVariant
	err := abi.ReadVariant(opts, 0, &v)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "out of range")

	// Writing with invalid case index (manually craft struct) not possible via type-safe API,
	// but we can simulate by setting discriminant > cases and calling WriteVariant.
	bad := SampleVariant{Type: SampleVariantType(99)}
	_, _, err = abi.WriteVariant(opts, bad, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "out of range")
}

// Record type for nested payload usage
type NestedRecord struct {
	X int16
	Y uint64
}

// ComplexVariant tests joining across integer, float, string, list and empty cases.
// Case ordering chosen to force widening (i32 vs f32 -> i32, i32 vs i64 -> i64) behavior.
type ComplexVariantType int

const (
	ComplexVariantTypeEmpty  ComplexVariantType = iota
	ComplexVariantTypeNumber                    // int32 payload
	ComplexVariantTypeFloat                     // float32 payload (joins with int32 -> i32 slot)
	ComplexVariantTypeBig                       // uint64 payload (forces widening to i64)
	ComplexVariantTypeText                      // string payload (adds two i32 slots)
	ComplexVariantTypeList                      // []uint8 payload (pointer+len two i32 slots)
	ComplexVariantTypeRecord                    // nestedRecord payload (two fields -> i32 + i64 -> widens first to i64 total slots maybe 2)
)

type ComplexVariant struct {
	Type   ComplexVariantType
	Empty  struct{}
	Number int32
	Float  float32
	Big    uint64
	Text   string
	List   []uint8
	Record NestedRecord
}

// buildComplexCases returns instances of ComplexVariant covering all payload kinds.
func buildComplexCases() []ComplexVariant {
	return []ComplexVariant{
		{Type: ComplexVariantTypeEmpty},
		{Type: ComplexVariantTypeNumber, Number: 12345},
		{Type: ComplexVariantTypeFloat, Float: 3.14},
		{Type: ComplexVariantTypeBig, Big: math.MaxUint64},
		{Type: ComplexVariantTypeText, Text: "hi"},
		{Type: ComplexVariantTypeText, Text: ""}, // empty string
		{Type: ComplexVariantTypeList, List: []uint8{1, 2, 3, 4}},
		{Type: ComplexVariantTypeList, List: []uint8{}},
		{Type: ComplexVariantTypeRecord, Record: NestedRecord{X: -7, Y: 99}},
	}
}

func TestComplexVariantParameterFlatteningShapeConsistency(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	cases := buildComplexCases()
	var expectedLen int = -1
	for _, c := range cases {
		flat, free, err := abi.WriteParameters(opts, c)
		require.NoError(t, err)
		if expectedLen == -1 {
			expectedLen = len(flat)
		} else {
			assert.Equal(t, expectedLen, len(flat), "all complex variant param lists must have same length")
		}
		// Discriminant must always be first
		assert.Equal(t, uint64(c.Type), flat[0])
		free()
	}
	assert.GreaterOrEqual(t, expectedLen, 2) // Should have at least discriminant + one payload slot
}

func TestComplexVariantRoundTrip(t *testing.T) {
	opts := createAbiOptionsFromMemoryMap(nil)
	for _, c := range buildComplexCases() {
		ptr, free, err := abi.Write(opts, c, nil)
		require.NoError(t, err)
		var decoded ComplexVariant
		err = abi.Read(opts, ptr, &decoded)
		require.NoError(t, err)
		assert.Equal(t, c.Type, decoded.Type)
		switch c.Type {
		case ComplexVariantTypeNumber:
			assert.Equal(t, c.Number, decoded.Number)
		case ComplexVariantTypeFloat:
			assert.InDelta(t, c.Float, decoded.Float, 1e-6)
		case ComplexVariantTypeBig:
			assert.Equal(t, c.Big, decoded.Big)
		case ComplexVariantTypeText:
			assert.Equal(t, c.Text, decoded.Text)
		case ComplexVariantTypeList:
			assert.Equal(t, c.List, decoded.List)
		case ComplexVariantTypeRecord:
			assert.Equal(t, c.Record, decoded.Record)
		}
		free()
	}
}

func TestComplexVariantInvalidDiscriminantRead(t *testing.T) {
	// Craft memory with invalid discriminant for ComplexVariant (e.g., 99)
	opts := createAbiOptionsFromMemoryMap(map[uint64][]byte{
		0: {99, 0, 0, 0, 0, 0, 0, 0},
	})
	var cv ComplexVariant
	err := abi.ReadVariant(opts, 0, &cv)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "out of range")
}

func TestComplexVariantParameterMismatchAfterTypeChange(t *testing.T) {
	// Simulate user mistakenly changing discriminant after preparing payload
	opts := createAbiOptionsFromMemoryMap(nil)
	v := ComplexVariant{Type: ComplexVariantTypeText, Text: "hello"}
	flat, free, err := abi.WriteParameters(opts, v)
	require.NoError(t, err)
	// Mutate discriminant to different case that expects scalar only
	flat[0] = uint64(ComplexVariantTypeNumber)
	// We cannot directly call into guest here, but we can assert shape still coherent
	// and differs from naive encoding that would omit slots.
	require.GreaterOrEqual(t, len(flat), 3)
	free()
}
