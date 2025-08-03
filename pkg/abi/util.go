package abi

import "math"

// AlignTo aligns a pointer to the nearest multiple of the specified alignment.
func AlignTo(ptr uint32, alignment uint32) uint32 {
	if alignment <= 0 {
		return ptr
	}
	return uint32(math.Ceil(float64(ptr)/float64(alignment)) * float64(alignment))
}
