package motifeditor

import (
	"image"
	"testing"
)

func TestStripPanoramaSeparatorsImage(t *testing.T) {
	t.Parallel()
	img := image.NewRGBA(image.Rect(0, 0, 4230, 10))
	for x := 0; x < 4230; x++ {
		if isGapColumn(x, 7, 600, 5) {
			img.Set(x, 0, image.Black)
		} else {
			img.Set(x, 0, image.White)
		}
	}

	out, changed := stripPanoramaSeparatorsImage(img, 7, 600, 5)
	if !changed {
		t.Fatal("expected gaps to be stripped")
	}
	if out.Bounds().Dx() != 4200 {
		t.Fatalf("width = %d, want 4200", out.Bounds().Dx())
	}
	if _, _, _, a := out.At(599, 0).RGBA(); a>>8 != 255 {
		t.Fatal("expected opaque pixel at end of first slice")
	}
	if _, _, _, a := out.At(600, 0).RGBA(); a>>8 != 255 {
		t.Fatal("expected opaque pixel at start of second slice (no gap)")
	}
}

func TestStripPanoramaSeparatorsImageAlreadySeamless(t *testing.T) {
	t.Parallel()
	img := image.NewRGBA(image.Rect(0, 0, 4200, 10))
	_, changed := stripPanoramaSeparatorsImage(img, 7, 600, 5)
	if changed {
		t.Fatal("already seamless input should not change")
	}
}

func isGapColumn(x, slideCount, slideWidth, gapPx int) bool {
	seamless := slideCount * slideWidth
	expectedWithGaps := seamless + (slideCount-1)*gapPx
	if x < 0 || x >= expectedWithGaps {
		return false
	}
	pos := 0
	for i := 0; i < slideCount; i++ {
		if x >= pos && x < pos+slideWidth {
			return false
		}
		pos += slideWidth
		if i < slideCount-1 {
			if x >= pos && x < pos+gapPx {
				return true
			}
			pos += gapPx
		}
	}
	return false
}
