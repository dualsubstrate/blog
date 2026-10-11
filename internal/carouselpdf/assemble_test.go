package carouselpdf

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseSlideFileName(t *testing.T) {
	t.Parallel()
	got, ok := parseSlideFileName("why-humans-keep-building-pyramids-slide-09-a.webp")
	if !ok {
		t.Fatal("expected match")
	}
	if got.Slug != "why-humans-keep-building-pyramids" || got.Number != 9 || got.Variant != "a" {
		t.Fatalf("got %+v", got)
	}
	if _, ok := parseSlideFileName("why-humans-keep-building-pyramids-panorama.webp"); ok {
		t.Fatal("panorama should not match")
	}
}

func TestCollectFromDirPicksSortedSlides(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mustWritePNG(t, filepath.Join(dir, "deck-slide-02-a.png"))
	mustWritePNG(t, filepath.Join(dir, "deck-slide-01-a.png"))
	mustWritePNG(t, filepath.Join(dir, "deck-panorama.png"))

	files, err := CollectFromDir(CollectOptions{Dir: dir, Slug: "deck"})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 2 || files[0].Number != 1 || files[1].Number != 2 {
		t.Fatalf("got %+v", files)
	}
}

func TestCollectFromDirRequiresVariantWhenDuplicates(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	mustWritePNG(t, filepath.Join(dir, "deck-slide-01-a.png"))
	mustWritePNG(t, filepath.Join(dir, "deck-slide-01-b.png"))
	if _, err := CollectFromDir(CollectOptions{Dir: dir, Slug: "deck"}); err == nil {
		t.Fatal("expected error for multiple variants")
	}
	files, err := CollectFromDir(CollectOptions{Dir: dir, Slug: "deck", Variant: "b"})
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 || files[0].Variant != "b" {
		t.Fatalf("got %+v", files)
	}
}

func TestWriteFileFullBleedPages(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p1 := filepath.Join(dir, "deck-slide-01-a.png")
	p2 := filepath.Join(dir, "deck-slide-02-a.png")
	mustWritePNGSize(t, p1, 8, 10, color.RGBA{R: 200, A: 255})
	mustWritePNGSize(t, p2, 8, 10, color.RGBA{B: 200, A: 255})
	out := filepath.Join(dir, "deck-linkedin.pdf")
	slides, err := CollectFromPaths([]string{p1, p2})
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteFile(out, slides, AssembleOptions{}); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(raw, []byte("%PDF-1.4")) {
		t.Fatalf("missing PDF header, got %q", raw[:min(len(raw), 16)])
	}
	if !bytes.Contains(raw, []byte("%%EOF")) {
		t.Fatal("missing EOF")
	}
	if c := bytes.Count(raw, []byte("/Type /Page ")); c != 2 {
		t.Fatalf("page objects = %d, want 2", c)
	}
	if !bytes.Contains(raw, []byte("/MediaBox [0 0 8 10]")) {
		t.Fatal("expected full-bleed MediaBox matching pixel size")
	}
	if !bytes.Contains(raw, []byte("/Filter /FlateDecode")) {
		t.Fatal("expected lossless FlateDecode image stream")
	}
	if bytes.Contains(raw, []byte("/DCTDecode")) {
		t.Fatal("JPEG DCTDecode should not be present")
	}
}

func mustWritePNG(t *testing.T, path string) {
	t.Helper()
	mustWritePNGSize(t, path, 4, 4, color.RGBA{G: 180, A: 255})
}

func mustWritePNGSize(t *testing.T, path string, w, h int, c color.RGBA) {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}

func TestCollectFromDirEmpty(t *testing.T) {
	t.Parallel()
	_, err := CollectFromDir(CollectOptions{Dir: t.TempDir()})
	if err == nil || !strings.Contains(err.Error(), "no studio slide exports") {
		t.Fatalf("got %v", err)
	}
}
