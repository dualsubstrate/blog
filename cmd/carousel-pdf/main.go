// Command carousel-pdf assembles studio slide rasters into a full-bleed LinkedIn PDF.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dualsubstrate/blog/internal/carouselpdf"
	"github.com/dualsubstrate/blog/internal/cliout"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dir := flag.String("dir", "", "directory of studio slide WebP exports")
	slug := flag.String("slug", "", "deck slug filter (required when -dir contains more than one deck)")
	variant := flag.String("variant", "", "variant letter when several exist per slide (a, b, ...)")
	outPath := flag.String("o", "", "output PDF path (default: DIR/SLUG-linkedin.pdf)")
	flag.Parse()

	slides, err := collectSlides(*dir, *slug, *variant, flag.Args())
	if err != nil {
		log.Fatalf("slides: %v", err)
	}
	out, err := resolveOutPath(*outPath, *dir, slides)
	if err != nil {
		log.Fatalf("out: %v", err)
	}
	if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
		log.Fatalf("mkdir: %v", err)
	}
	if err := carouselpdf.WriteFile(out, slides, carouselpdf.AssembleOptions{}); err != nil {
		log.Fatalf("write: %v", err)
	}
	absOut, err := filepath.Abs(out)
	if err != nil {
		absOut = out
	}
	cliout.PrintFileWritten(os.Stdout, "carousel-pdf", absOut, fmt.Sprintf("%d pages", len(slides)))
	cliout.Hint(os.Stdout, "LinkedIn", "upload this PDF as a Document / Carousel post")
	cliout.Blank(os.Stdout)
}

func collectSlides(dir, slug, variant string, args []string) ([]carouselpdf.SlideFile, error) {
	if strings.TrimSpace(dir) != "" {
		if len(args) > 0 {
			return nil, fmt.Errorf("pass either -dir or image paths, not both")
		}
		return carouselpdf.CollectFromDir(carouselpdf.CollectOptions{
			Dir:     dir,
			Slug:    slug,
			Variant: variant,
		})
	}
	if len(args) == 0 {
		return nil, fmt.Errorf("pass -dir of studio WebPs, or image paths")
	}
	return carouselpdf.CollectFromPaths(args)
}

func resolveOutPath(out, dir string, slides []carouselpdf.SlideFile) (string, error) {
	if strings.TrimSpace(out) != "" {
		return out, nil
	}
	slug := ""
	if len(slides) > 0 {
		slug = slides[0].Slug
	}
	if slug == "" {
		slug = "carousel"
	}
	base := strings.TrimSpace(dir)
	if base == "" && len(slides) > 0 {
		base = filepath.Dir(slides[0].Path)
	}
	if base == "" {
		return "", fmt.Errorf("need -o when output directory cannot be inferred")
	}
	return filepath.Join(base, slug+"-linkedin.pdf"), nil
}
