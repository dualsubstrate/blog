package carousel

import "testing"

func TestMotifStripWidthPx(t *testing.T) {
	t.Parallel()
	if got := MotifStripWidthPx(7); got != 7*PanoramaSlideWidthPx {
		t.Fatalf("MotifStripWidthPx(7) = %d, want %d", got, 7*PanoramaSlideWidthPx)
	}
	if got := MotifStripWidthPx(0); got != PanoramaSlideWidthPx {
		t.Fatalf("zero slides should clamp to 1 slice, got %d", got)
	}
}

func TestPanoramaWidthWithGapsPx(t *testing.T) {
	t.Parallel()
	if got := PanoramaWidthWithGapsPx(7, PanoramaSlideWidthPx); got != 4230 {
		t.Fatalf("PanoramaWidthWithGapsPx(7) = %d, want 4230", got)
	}
	if got := PanoramaGapPx(PanoramaSlideWidthPx); got != PanoramaGapPxDefault {
		t.Fatalf("PanoramaGapPx(600) = %d, want %d", got, PanoramaGapPxDefault)
	}
}
