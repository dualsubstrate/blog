// Package substackpublishstate tracks which Hugo page bundles were published to which
// channels using a marker file next to index.md (default name: social-published).
// ListUnpublished only includes bundles that look ready to publish: not draft on the
// relevant Markdown (index.md for all targets; index.es.md too for substack-es), plus
// channel-specific files such as index.es.md or linkedin.txt where applicable.
//
// File format (UTF-8), one record per line:
//   - Blank lines and lines starting with # are ignored.
//   - Each data line: "<target> <RFC3339-UTC>" (whitespace-separated; target has no spaces).
//   - Example targets: substack-en, substack-es, linkedin (any non-empty id you choose).
//
// Legacy: a single line containing only an RFC3339 timestamp is treated as substack-en.
package substackpublishstate

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DefaultMarker is the default marker file name inside each bundle.
const DefaultMarker = "social-published"

// LegacyDefaultTarget is the channel implied by a legacy marker file (one RFC3339 line only).
const LegacyDefaultTarget = "substack-en"

// TargetSubstackEN is the usual id for the English Substack flow (sb-en / pick-draft default).
const TargetSubstackEN = "substack-en"

// TargetSubstackES is the usual id for the Spanish Substack flow (pick-draft-es default).
const TargetSubstackES = "substack-es"

// TargetSiteEN is the Hugo site English page (index.md live on blog.dualsubstrate.com).
const TargetSiteEN = "site-en"

// TargetSiteES is the Hugo site Spanish page (index.es.md).
const TargetSiteES = "site-es"

// BundlePostPath is a slash-separated path relative to the content root, e.g. "human-condition/2026-05-01-ego-as-game".
type BundlePostPath string

// ParseMarkerContent returns target -> timestamp string from marker file bytes.
func ParseMarkerContent(data []byte) map[string]string {
	var lines []string
	for _, ln := range strings.Split(string(data), "\n") {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			lines = append(lines, ln)
		}
	}
	if len(lines) == 0 {
		return nil
	}
	out := make(map[string]string)
	if len(lines) == 1 {
		line := lines[0]
		if !strings.HasPrefix(line, "#") && strings.Count(line, " ") == 0 && isRFC3339Date(line) {
			out[LegacyDefaultTarget] = line
			return out
		}
	}
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" || strings.HasPrefix(ln, "#") {
			continue
		}
		parts := strings.Fields(ln)
		if len(parts) < 2 {
			continue
		}
		key := parts[0]
		ts := strings.Join(parts[1:], " ")
		if key == "" {
			continue
		}
		out[key] = ts
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func isRFC3339Date(s string) bool {
	_, err := time.Parse(time.RFC3339, s)
	return err == nil
}

// MarkerHasTarget reports whether markerPath records an entry for targetKey.
func MarkerHasTarget(markerPath, targetKey string) (bool, error) {
	targetKey = strings.TrimSpace(targetKey)
	if targetKey == "" {
		return false, fmt.Errorf("substackpublishstate: empty target key")
	}
	b, err := os.ReadFile(markerPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	m := ParseMarkerContent(b)
	_, ok := m[targetKey]
	return ok, nil
}

// ListUnpublished scans contentRoot for leaf bundles with index.md. It returns bundles that
// do not record targetKey in markerName (missing file, or file with no line for targetKey)
// and that BundleReadyForPublish reports as ready for targetKey.
func ListUnpublished(contentRoot, markerName, targetKey string) ([]BundlePostPath, error) {
	contentRoot = strings.TrimSpace(contentRoot)
	if contentRoot == "" {
		contentRoot = "content"
	}
	markerName = strings.TrimSpace(markerName)
	if markerName == "" {
		markerName = DefaultMarker
	}
	targetKey = strings.TrimSpace(targetKey)
	if targetKey == "" {
		targetKey = LegacyDefaultTarget
	}
	absRoot, err := filepath.Abs(contentRoot)
	if err != nil {
		return nil, fmt.Errorf("substackpublishstate: content root: %w", err)
	}
	var out []BundlePostPath
	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != "index.md" {
			return nil
		}
		bundleDir := filepath.Dir(path)
		rel, err := filepath.Rel(absRoot, bundleDir)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		parts := strings.Split(rel, string(filepath.Separator))
		if len(parts) < 2 {
			return nil
		}
		markerPath := filepath.Join(bundleDir, markerName)
		has, err := MarkerHasTarget(markerPath, targetKey)
		if err != nil {
			return err
		}
		if has {
			return nil
		}
		ready, err := BundleReadyForPublish(bundleDir, targetKey)
		if err != nil {
			return err
		}
		if !ready {
			return nil
		}
		out = append(out, BundlePostPath(filepath.ToSlash(rel)))
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool { return string(out[i]) < string(out[j]) })
	return out, nil
}

// MarkPublished adds or updates targetKey in the bundle marker file with the current UTC time.
func MarkPublished(bundleDir, markerName, targetKey string) error {
	bundleDir = strings.TrimSpace(bundleDir)
	if bundleDir == "" {
		return fmt.Errorf("substackpublishstate: empty bundle directory")
	}
	markerName = strings.TrimSpace(markerName)
	if markerName == "" {
		markerName = DefaultMarker
	}
	targetKey = strings.TrimSpace(targetKey)
	if targetKey == "" {
		return fmt.Errorf("substackpublishstate: empty publish target (e.g. substack-en, substack-es, linkedin)")
	}
	if strings.Contains(targetKey, " ") {
		return fmt.Errorf("substackpublishstate: publish target must not contain spaces")
	}
	abs, err := filepath.Abs(bundleDir)
	if err != nil {
		return fmt.Errorf("substackpublishstate: bundle dir: %w", err)
	}
	index := filepath.Join(abs, "index.md")
	if st, err := os.Stat(index); err != nil || st.IsDir() {
		return fmt.Errorf("substackpublishstate: need bundle with index.md: %s", index)
	}
	p := filepath.Join(abs, markerName)
	now := time.Now().UTC().Format(time.RFC3339)

	existing := make(map[string]string)
	if b, err := os.ReadFile(p); err == nil {
		existing = ParseMarkerContent(b)
		if existing == nil {
			existing = make(map[string]string)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("substackpublishstate: read marker: %w", err)
	}
	existing[targetKey] = now

	var keys []string
	for k := range existing {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString("# <target> <RFC3339-UTC> - channels published for this bundle\n")
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte(' ')
		sb.WriteString(existing[k])
		sb.WriteByte('\n')
	}
	if err := os.WriteFile(p, []byte(sb.String()), 0o644); err != nil {
		return fmt.Errorf("substackpublishstate: write marker: %w", err)
	}
	return nil
}
