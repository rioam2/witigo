package witigo

import "math"

// AlignTo aligns a pointer to the nearest multiple of the specified alignment.
func AlignTo(ptr int, alignment int) int {
	if alignment <= 0 {
		return ptr
	}
	return int(math.Ceil(float64(ptr)/float64(alignment)) * float64(alignment))
}
