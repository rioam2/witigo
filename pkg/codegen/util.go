package codegen

import (
	"math"
)

// discriminantSize returns the size in bits of the discriminant needed to represent n cases.
func discriminantSize(n int) int {
	bytesNeeded := int(math.Ceil(math.Log2(float64(n+1)) / 8.0))
	switch bytesNeeded {
	case 0, 1:
		return 8
	case 2:
		return 16
	default:
		return 32
	}
}
