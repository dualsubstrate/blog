// Package contentbundle resolves Hugo page bundles under content/ by publish date,
// using the same rules as facebook-autopost (two- and three-level section paths).
package contentbundle

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dualsubstrate/blog/internal/tagregister"
	"gopkg.in/yaml.v3"
)

// PublishedBundleRelsForDate returns bundle directory paths relative to content/
// (for example human-condition/2026-05-01-ego-as-game or cognitive-memetics/reptilocracy/2026-05-10-slug),
// sorted lexicographically by that relative path. Each bundle has index.md whose front matter date
// matches dateYYYYMMDD (wall date in the timestamp's offset).
//
// If no bundles match the date, it returns an empty slice and logs a warning (not an error), so
// scheduled social autopost can exit cleanly on empty calendar days.
//
// If postPath is non-empty, it must be a single bundle path under content/ (with or without content/ prefix);
// that path is returned as a one-element slice without date filtering (same as facebook-autopost -post).
func PublishedBundleRelsForDate(repoRoot, dateYYYYMMDD, postPath string) ([]string, error) {
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, err
	}
	contentRoot := filepath.Join(absRoot, "content")
	postPath = strings.TrimSpace(postPath)
	if postPath != "" {
		return []string{normalizeBundleRel(postPath)}, nil
	}

	pattern := filepath.Join(contentRoot, "*", "*", "index.md")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	deeper, err := filepath.Glob(filepath.Join(contentRoot, "*", "*", "*", "index.md"))
	if err != nil {
		return nil, err
	}
	matches = append(matches, deeper...)

	type cand struct {
		bundleDir string
	}
	var cands []cand
	for _, indexPath := range matches {
		raw, err := os.ReadFile(indexPath)
		if err != nil {
			return nil, err
		}
		relLog := indexPath
		if r, err := filepath.Rel(contentRoot, indexPath); err == nil {
			relLog = filepath.ToSlash(r)
		}
		fm, err := tagregister.FrontMatterYAML(raw)
		if err != nil {
			log.Printf("contentbundle: skip %s: front matter: %v", relLog, err)
			continue
		}
		var doc struct {
			Date string `yaml:"date"`
		}
		if err := yaml.Unmarshal(fm, &doc); err != nil {
			log.Printf("contentbundle: skip %s: YAML date block: %v", relLog, err)
			continue
		}
		dateStr := strings.TrimSpace(doc.Date)
		t, err := time.Parse(time.RFC3339, dateStr)
		if err != nil {
			log.Printf("contentbundle: skip %s: date %q is not RFC3339: %v", relLog, dateStr, err)
			continue
		}
		if t.Format("2006-01-02") != dateYYYYMMDD {
			continue
		}
		cands = append(cands, cand{bundleDir: filepath.Dir(indexPath)})
	}

	if len(cands) == 0 {
		log.Printf("contentbundle: warning: no content bundle found for date %s (nothing to post)", dateYYYYMMDD)
		return []string{}, nil
	}

	type bundlePick struct {
		rel string
	}
	var picks []bundlePick
	for _, c := range cands {
		rel, err := filepath.Rel(contentRoot, c.bundleDir)
		if err != nil {
			return nil, err
		}
		picks = append(picks, bundlePick{rel: filepath.ToSlash(rel)})
	}
	sort.Slice(picks, func(i, j int) bool { return picks[i].rel < picks[j].rel })

	out := make([]string, 0, len(picks))
	for _, b := range picks {
		out = append(out, b.rel)
	}
	return out, nil
}

func normalizeBundleRel(postPath string) string {
	s := strings.Trim(postPath, "/\\")
	s = strings.TrimPrefix(s, "content/")
	s = strings.Trim(s, "/\\")
	return filepath.ToSlash(filepath.Clean(s))
}
