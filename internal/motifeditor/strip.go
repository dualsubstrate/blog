package motifeditor

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"time"
)

// StripPanoramaSeparators removes inter-slide gaps from a studio panorama export.
// Returns true when gaps were stripped and written to outputPath.
func StripPanoramaSeparators(inputPath, outputPath string, slideCount, slideWidth, gapPx int) (bool, error) {
	if slideCount < 1 || slideWidth < 1 || gapPx < 0 {
		return false, nil
	}

	f, err := os.Open(inputPath)
	if err != nil {
		return false, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return false, err
	}

	stripped, changed := stripPanoramaSeparatorsImage(img, slideCount, slideWidth, gapPx)
	if !changed {
		return false, nil
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return false, err
	}
	defer out.Close()
	if err := png.Encode(out, stripped); err != nil {
		return false, err
	}
	return true, out.Close()
}

func stripPanoramaSeparatorsImage(img image.Image, slideCount, slideWidth, gapPx int) (image.Image, bool) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	seamless := slideCount * slideWidth
	expectedWithGaps := seamless + (slideCount-1)*gapPx

	if width == seamless {
		return img, false
	}
	if width != expectedWithGaps {
		return img, false
	}

	if rgba, ok := img.(*image.RGBA); ok {
		out := image.NewRGBA(image.Rect(0, 0, seamless, height))
		dx := 0
		for i := 0; i < slideCount; i++ {
			sx := bounds.Min.X + i*(slideWidth+gapPx)
			for y := 0; y < height; y++ {
				srcStart := rgba.PixOffset(sx, bounds.Min.Y+y)
				srcEnd := srcStart + slideWidth*4
				dstStart := out.PixOffset(dx, y)
				copy(out.Pix[dstStart:dstStart+slideWidth*4], rgba.Pix[srcStart:srcEnd])
			}
			dx += slideWidth
		}
		return out, true
	}
	if nrgba, ok := img.(*image.NRGBA); ok {
		out := image.NewNRGBA(image.Rect(0, 0, seamless, height))
		dx := 0
		for i := 0; i < slideCount; i++ {
			sx := bounds.Min.X + i*(slideWidth+gapPx)
			for y := 0; y < height; y++ {
				srcStart := nrgba.PixOffset(sx, bounds.Min.Y+y)
				srcEnd := srcStart + slideWidth*4
				dstStart := out.PixOffset(dx, y)
				copy(out.Pix[dstStart:dstStart+slideWidth*4], nrgba.Pix[srcStart:srcEnd])
			}
			dx += slideWidth
		}
		return out, true
	}

	out := image.NewRGBA(image.Rect(0, 0, seamless, height))
	dx := 0
	for i := 0; i < slideCount; i++ {
		sx := bounds.Min.X + i*(slideWidth+gapPx)
		for y := 0; y < height; y++ {
			for x := 0; x < slideWidth; x++ {
				out.Set(dx+x, y, img.At(sx+x, bounds.Min.Y+y))
			}
		}
		dx += slideWidth
	}
	return out, true
}

// MaybeStripPanoramaSeparators strips gaps when dimensions match; otherwise returns inputPath unchanged.
func MaybeStripPanoramaSeparators(inputPath string, slideCount, slideWidth, gapPx int) (string, func(), error) {
	if slideCount < 1 {
		return inputPath, func() {}, nil
	}

	strippedPath := filepath.Join(filepath.Dir(inputPath), fmt.Sprintf("upscale-stripped-%d.png", time.Now().UnixNano()))
	stripped, err := StripPanoramaSeparators(inputPath, strippedPath, slideCount, slideWidth, gapPx)
	if err != nil {
		os.Remove(strippedPath)
		return inputPath, func() {}, err
	}
	if !stripped {
		return inputPath, func() {}, nil
	}
	return strippedPath, func() { os.Remove(strippedPath) }, nil
}
