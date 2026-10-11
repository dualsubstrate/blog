package tagregister

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// ScanMarkdownTags walks root for *.md files, parses front matter, and counts how many
// files contain each tag (at most once per file per distinct tag).
func ScanMarkdownTags(root string) (map[string]int, error) {
	counts := make(map[string]int)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(strings.ToLower(path), ".md") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fm, err := FrontMatterYAML(b)
		if err != nil {
			// Pages without Hugo-style front matter: skip, do not fail the whole scan.
			return nil
		}
		tags, err := TagsFromYAML(fm)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		tags = UniqueStrings(tags)
		for _, t := range tags {
			counts[t]++
		}
		return nil
	})
	return counts, err
}
