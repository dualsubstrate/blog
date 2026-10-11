// Package carousel defines shared raster dimensions for carousel and panorama tooling.
//
// Keep PanoramaSlideWidthPx and PanoramaGapPxDefault in sync with static/carousel/slide-constants.js.
package carousel

import "math"

// SlideWidthPx is the carousel design canvas width (px).
const SlideWidthPx = 1080

// PanoramaSlideWidthPx is the studio panorama export width per slide (theme.size default).
const PanoramaSlideWidthPx = 600

// PanoramaGapPxDefault is the inter-slide gap at PanoramaSlideWidthPx (studio export).
const PanoramaGapPxDefault = 5

const (
	panoramaGapPxAtReference         = 4
	panoramaGapReferenceSlideWidthPx = 480
)

// PanoramaGapPx returns the gap between slides in a studio panorama export.
func PanoramaGapPx(slideWidth int) int {
	if slideWidth < 1 {
		slideWidth = PanoramaSlideWidthPx
	}
	if slideWidth == PanoramaSlideWidthPx {
		return PanoramaGapPxDefault
	}
	gap := int(math.Round(float64(slideWidth) * float64(panoramaGapPxAtReference) / float64(panoramaGapReferenceSlideWidthPx)))
	if gap < 2 {
		return 2
	}
	return gap
}

// PanoramaWidthWithGapsPx returns strip width including inter-slide gaps (studio export).
func PanoramaWidthWithGapsPx(slideCount, slideWidth int) int {
	if slideCount < 1 {
		slideCount = 1
	}
	if slideWidth < 1 {
		slideWidth = PanoramaSlideWidthPx
	}
	gap := PanoramaGapPx(slideWidth)
	return slideCount*slideWidth + (slideCount-1)*gap
}

// MotifStripSeamlessWidthPx returns motif asset width with separators removed (slides × panorama slice).
func MotifStripSeamlessWidthPx(slideCount int) int {
	return MotifStripSeamlessWidthForSlideWidth(slideCount, PanoramaSlideWidthPx)
}

// MotifStripSeamlessWidthForSlideWidth returns seamless width for a custom panorama slice width.
func MotifStripSeamlessWidthForSlideWidth(slideCount, slideWidth int) int {
	if slideCount < 1 {
		slideCount = 1
	}
	if slideWidth < 1 {
		slideWidth = PanoramaSlideWidthPx
	}
	return slideCount * slideWidth
}

// MotifStripWidthPx returns the seamless motif strip width for a deck (motif editor upscale target).
func MotifStripWidthPx(slideCount int) int {
	return MotifStripSeamlessWidthPx(slideCount)
}
