// Package carouselpdf assembles studio slide rasters into a full-bleed LinkedIn PDF.
package carouselpdf

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var slideFileRe = regexp.MustCompile(`(?i)^(.+)-slide-(\d+)-([a-z][a-z0-9]*)\.(webp|png|jpe?g)$`)

// SlideFile is one studio export chosen for a deck slot.
type SlideFile struct {
	Path     string
	Slug     string
	Number   int
	Variant  string
	Ext      string
	BaseName string
}

// CollectOptions controls directory scans of studio WebP exports.
type CollectOptions struct {
	Dir     string
	Slug    string
	Variant string
}

// CollectFromDir finds `{slug}-slide-NN-{variant}.{ext}` files, sorted by slide number.
// Panorama strips are ignored. If several variants exist for one number, Variant must be set
// (or every duplicate group must share a single variant).
func CollectFromDir(opts CollectOptions) ([]SlideFile, error) {
	dir := strings.TrimSpace(opts.Dir)
	if dir == "" {
		return nil, fmt.Errorf("directory is empty")
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir %s: %w", dir, err)
	}
	wantSlug := strings.TrimSpace(opts.Slug)
	wantVar := strings.ToLower(strings.TrimSpace(opts.Variant))

	var found []SlideFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		parsed, ok := parseSlideFileName(name)
		if !ok {
			continue
		}
		if wantSlug != "" && parsed.Slug != wantSlug {
			continue
		}
		parsed.Path = filepath.Join(dir, name)
		found = append(found, parsed)
	}
	if len(found) == 0 {
		return nil, emptyDirError(dir, wantSlug, entries)
	}

	slugs := uniqueSlugs(found)
	if wantSlug == "" && len(slugs) > 1 {
		return nil, fmt.Errorf("multiple slugs in %s (%s); pass -slug", dir, strings.Join(slugs, ", "))
	}

	grouped := groupByNumber(found)
	out := make([]SlideFile, 0, len(grouped))
	numbers := make([]int, 0, len(grouped))
	for n := range grouped {
		numbers = append(numbers, n)
	}
	sort.Ints(numbers)
	for _, n := range numbers {
		chosen, err := pickVariant(grouped[n], wantVar, n)
		if err != nil {
			return nil, err
		}
		out = append(out, chosen)
	}
	return out, nil
}

// CollectFromPaths keeps explicit file order after validating each path exists.
func CollectFromPaths(paths []string) ([]SlideFile, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no image paths")
	}
	out := make([]SlideFile, 0, len(paths))
	for _, raw := range paths {
		path := strings.TrimSpace(raw)
		if path == "" {
			return nil, fmt.Errorf("empty image path")
		}
		st, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("stat %s: %w", path, err)
		}
		if st.IsDir() {
			return nil, fmt.Errorf("%s is a directory; use -dir", path)
		}
		parsed, ok := parseSlideFileName(filepath.Base(path))
		if !ok {
			parsed = SlideFile{BaseName: filepath.Base(path)}
		}
		parsed.Path = path
		out = append(out, parsed)
	}
	return out, nil
}

func emptyDirError(dir, wantSlug string, entries []os.DirEntry) error {
	example := "deck-slide-01-a.webp"
	if wantSlug != "" {
		example = wantSlug + "-slide-01-a.webp"
	}
	otherSlugs := slugsInDir(entries)
	msg := fmt.Sprintf("no studio slide exports in %s (expected %s)", dir, example)
	if wantSlug != "" && len(otherSlugs) > 0 {
		msg += fmt.Sprintf("; found other slugs: %s", strings.Join(otherSlugs, ", "))
	} else {
		msg += "; export from carousel.preview with the Slides button first (Chrome may ask to Allow multiple downloads)"
	}
	return fmt.Errorf("%s", msg)
}

func slugsInDir(entries []os.DirEntry) []string {
	var files []SlideFile
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		parsed, ok := parseSlideFileName(entry.Name())
		if !ok {
			continue
		}
		files = append(files, parsed)
	}
	return uniqueSlugs(files)
}

func parseSlideFileName(name string) (SlideFile, bool) {
	m := slideFileRe.FindStringSubmatch(name)
	if m == nil {
		return SlideFile{}, false
	}
	num, err := strconv.Atoi(m[2])
	if err != nil || num < 1 {
		return SlideFile{}, false
	}
	return SlideFile{
		Slug:     m[1],
		Number:   num,
		Variant:  strings.ToLower(m[3]),
		Ext:      strings.ToLower(m[4]),
		BaseName: name,
	}, true
}

func uniqueSlugs(files []SlideFile) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, f := range files {
		if f.Slug == "" {
			continue
		}
		if _, ok := seen[f.Slug]; ok {
			continue
		}
		seen[f.Slug] = struct{}{}
		out = append(out, f.Slug)
	}
	sort.Strings(out)
	return out
}

func groupByNumber(files []SlideFile) map[int][]SlideFile {
	out := make(map[int][]SlideFile)
	for _, f := range files {
		out[f.Number] = append(out[f.Number], f)
	}
	return out
}

func pickVariant(group []SlideFile, wantVar string, number int) (SlideFile, error) {
	if len(group) == 0 {
		return SlideFile{}, fmt.Errorf("no files for slide %02d", number)
	}
	if wantVar != "" {
		for _, f := range group {
			if f.Variant == wantVar {
				return f, nil
			}
		}
		return SlideFile{}, fmt.Errorf("slide %02d has no variant %q", number, wantVar)
	}
	if len(group) == 1 {
		return group[0], nil
	}
	vars := make([]string, 0, len(group))
	for _, f := range group {
		vars = append(vars, f.Variant)
	}
	sort.Strings(vars)
	return SlideFile{}, fmt.Errorf("slide %02d has multiple variants (%s); pass -variant", number, strings.Join(vars, ", "))
}
