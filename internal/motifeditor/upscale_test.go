package motifeditor

import "testing"

func TestRealesrganTileSize(t *testing.T) {
	t.Parallel()
	tests := []struct {
		w, h, want int
	}{
		{4000, 150, 150},
		{7560, 200, 200},
		{800, 600, 512},
		{400, 300, 300},
		{20, 10, 32},
	}
	for _, tc := range tests {
		got := realesrganTileSize(tc.w, tc.h)
		if got != tc.want {
			t.Errorf("realesrganTileSize(%d, %d) = %d, want %d", tc.w, tc.h, got, tc.want)
		}
	}
}

func TestPreferWidthResizeOnly(t *testing.T) {
	t.Parallel()
	tests := []struct {
		w, h int
		want bool
	}{
		{4000, 150, true},
		{7560, 200, true},
		{800, 600, false},
		{2000, 400, false},
		{1000, 321, false},
	}
	for _, tc := range tests {
		got := preferWidthResizeOnly(tc.w, tc.h)
		if got != tc.want {
			t.Errorf("preferWidthResizeOnly(%d, %d) = %v, want %v", tc.w, tc.h, got, tc.want)
		}
	}
}
