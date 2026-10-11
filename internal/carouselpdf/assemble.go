package carouselpdf

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"io"
	"os"
	"strings"

	_ "golang.org/x/image/webp"
)

// AssembleOptions is reserved for future encode settings. Pages are lossless Flate RGB.
type AssembleOptions struct{}

// WriteFile encodes slides as a full-bleed PDF (one image per page, page size = pixel size).
func WriteFile(outPath string, slides []SlideFile, _ AssembleOptions) error {
	if strings.TrimSpace(outPath) == "" {
		return fmt.Errorf("output path is empty")
	}
	pages, err := encodePages(slides)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := writePDF(&buf, pages); err != nil {
		return err
	}
	if err := os.WriteFile(outPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", outPath, err)
	}
	return nil
}

type pdfPage struct {
	width  int
	height int
	flate  []byte
}

func encodePages(slides []SlideFile) ([]pdfPage, error) {
	if len(slides) == 0 {
		return nil, fmt.Errorf("no slides to assemble")
	}
	pages := make([]pdfPage, 0, len(slides))
	for _, slide := range slides {
		page, err := encodeSlide(slide.Path)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}
	return pages, nil
}

func encodeSlide(path string) (pdfPage, error) {
	f, err := os.Open(path)
	if err != nil {
		return pdfPage{}, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return pdfPage{}, fmt.Errorf("decode %s: %w", path, err)
	}
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w < 1 || h < 1 {
		return pdfPage{}, fmt.Errorf("%s has empty dimensions", path)
	}
	flate, err := flateRGB(img)
	if err != nil {
		return pdfPage{}, fmt.Errorf("flate encode %s: %w", path, err)
	}
	return pdfPage{width: w, height: h, flate: flate}, nil
}

func flateRGB(img image.Image) ([]byte, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(rgba, rgba.Bounds(), image.NewUniform(color.Black), image.Point{}, draw.Src)
	draw.Draw(rgba, rgba.Bounds(), img, bounds.Min, draw.Over)

	rgb := make([]byte, w*h*3)
	o := 0
	for y := 0; y < h; y++ {
		row := rgba.Pix[y*rgba.Stride : y*rgba.Stride+w*4]
		for x := 0; x < w; x++ {
			i := x * 4
			rgb[o] = row[i]
			rgb[o+1] = row[i+1]
			rgb[o+2] = row[i+2]
			o += 3
		}
	}

	var buf bytes.Buffer
	zw := zlib.NewWriter(&buf)
	if _, err := zw.Write(rgb); err != nil {
		_ = zw.Close()
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func writePDF(w io.Writer, pages []pdfPage) error {
	if len(pages) == 0 {
		return fmt.Errorf("no pages")
	}
	var b pdfBuf
	b.writeString("%PDF-1.4\n%\xff\xff\xff\xff\n")

	catalogID := 1
	pagesID := 2
	pageIDs := make([]int, len(pages))
	contentIDs := make([]int, len(pages))
	imageIDs := make([]int, len(pages))
	nextID := 3
	for i := range pages {
		contentIDs[i] = nextID
		imageIDs[i] = nextID + 1
		pageIDs[i] = nextID + 2
		nextID += 3
	}

	b.startObj(catalogID)
	b.writeString(fmt.Sprintf("<< /Type /Catalog /Pages %d 0 R >>\n", pagesID))
	b.endObj()

	b.startObj(pagesID)
	b.writeString(fmt.Sprintf("<< /Type /Pages /Count %d /Kids [", len(pages)))
	for i, id := range pageIDs {
		if i > 0 {
			b.writeString(" ")
		}
		b.writeString(fmt.Sprintf("%d 0 R", id))
	}
	b.writeString("] >>\n")
	b.endObj()

	for i, page := range pages {
		content := fmt.Sprintf("q %d 0 0 %d 0 0 cm /Im1 Do Q\n", page.width, page.height)
		b.startObj(contentIDs[i])
		b.writeString(fmt.Sprintf("<< /Length %d >>\nstream\n", len(content)))
		b.writeString(content)
		b.writeString("endstream\n")
		b.endObj()

		b.startObj(imageIDs[i])
		b.writeString(fmt.Sprintf(
			"<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /FlateDecode /Length %d >>\nstream\n",
			page.width, page.height, len(page.flate),
		))
		b.write(page.flate)
		b.writeString("endstream\n")
		b.endObj()

		b.startObj(pageIDs[i])
		b.writeString(fmt.Sprintf(
			"<< /Type /Page /Parent %d 0 R /MediaBox [0 0 %d %d] /Resources << /XObject << /Im1 %d 0 R >> >> /Contents %d 0 R >>\n",
			pagesID, page.width, page.height, imageIDs[i], contentIDs[i],
		))
		b.endObj()
	}

	startxref := b.buf.Len()
	b.writeString(fmt.Sprintf("xref\n0 %d\n", nextID))
	b.writeString("0000000000 65535 f \n")
	for id := 1; id < nextID; id++ {
		off, ok := b.offset[id]
		if !ok {
			return fmt.Errorf("missing pdf object %d", id)
		}
		b.writeString(fmt.Sprintf("%010d 00000 n \n", off))
	}
	b.writeString(fmt.Sprintf("trailer << /Size %d /Root %d 0 R >>\nstartxref\n%d\n%%%%EOF\n", nextID, catalogID, startxref))
	_, err := w.Write(b.buf.Bytes())
	return err
}

type pdfBuf struct {
	buf    bytes.Buffer
	offset map[int]int
}

func (p *pdfBuf) writeString(s string) {
	p.buf.WriteString(s)
}

func (p *pdfBuf) write(b []byte) {
	p.buf.Write(b)
}

func (p *pdfBuf) startObj(id int) {
	if p.offset == nil {
		p.offset = make(map[int]int)
	}
	p.offset[id] = p.buf.Len()
	p.writeString(fmt.Sprintf("%d 0 obj\n", id))
}

func (p *pdfBuf) endObj() {
	p.writeString("endobj\n")
}
