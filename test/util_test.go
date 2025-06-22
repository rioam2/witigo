package test

import (
	"testing"

	witigo "github.com/rioam2/witigo/pkg"
)

func TestAlignTo(t *testing.T) {
	tests := []struct {
		name      string
		ptr       int
		alignment int
		want      int
	}{
		{
			name:      "already aligned",
			ptr:       100,
			alignment: 10,
			want:      100,
		},
		{
			name:      "needs alignment",
			ptr:       105,
			alignment: 10,
			want:      110,
		},
		{
			name:      "zero alignment",
			ptr:       105,
			alignment: 0,
			want:      105,
		},
		{
			name:      "negative alignment",
			ptr:       105,
			alignment: -5,
			want:      105,
		},
		{
			name:      "both zero",
			ptr:       0,
			alignment: 0,
			want:      0,
		},
		{
			name:      "large alignment",
			ptr:       1000,
			alignment: 512,
			want:      1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := witigo.AlignTo(tt.ptr, tt.alignment)
			if got != tt.want {
				t.Errorf("AlignTo(%d, %d) = %d, want %d",
					tt.ptr, tt.alignment, got, tt.want)
			}
		})
	}
}
