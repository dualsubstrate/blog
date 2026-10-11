package substackhtml

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SubstackBodySource reports which bundle file supplied the newsletter body.
type SubstackBodySource int

const (
	// SubstackBodyFromIndex is unused; ResolveSubstackBody requires substack.md / substack.es.md.
	SubstackBodyFromIndex SubstackBodySource = iota
	SubstackBodyFromSidecarMD
)

// SubstackMarkdownBasename returns the required Substack body sidecar for the locale.
func SubstackMarkdownBasename(markdownPath string) string {
	if isSpanishMarkdownPath(markdownPath) {
		return "substack.es.md"
	}
	return "substack.md"
}

// FacebookCopyBasename returns the friends-Facebook plain-text sidecar for the locale.
func FacebookCopyBasename(markdownPath string) string {
	if isSpanishMarkdownPath(markdownPath) {
		return "facebook-es.txt"
	}
	return "facebook-en.txt"
}

func isSpanishMarkdownPath(markdownPath string) bool {
	return strings.EqualFold(filepath.Base(strings.TrimSpace(markdownPath)), "index.es.md")
}

// BundleDirFromMarkdownPath returns the page bundle directory for index.md, index.es.md,
// or a path that is already the bundle folder.
func BundleDirFromMarkdownPath(mdIn string) (string, error) {
	p := strings.TrimSpace(mdIn)
	if p == "" {
		return "", errors.New("substackhtml: empty markdown path")
	}
	base := strings.ToLower(filepath.Base(p))
	if base == "index.md" || base == "index.es.md" {
		p = filepath.Dir(p)
	}
	return filepath.Abs(p)
}

// StripFacebookSocialTail removes hashtag-only lines and the bilingual read-more URL
// block from friends Facebook copy. Trailing casual citations after the read-more block
// are kept.
func StripFacebookSocialTail(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	var out []string
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if isFacebookReadMoreHeader(line) {
			i++
			for i < len(lines) && isFacebookReadMoreLinkLine(lines[i]) {
				i++
			}
			i--
			continue
		}
		if isFacebookHashtagOnlyLine(line) {
			continue
		}
		out = append(out, line)
	}
	return strings.TrimSpace(strings.Join(out, "\n"))
}

func isFacebookReadMoreHeader(line string) bool {
	s := strings.TrimSpace(strings.ToLower(line))
	return s == "if you want to read more:" ||
		s == "si quieren leer más:" ||
		s == "si quieren leer mas:"
}

func isFacebookReadMoreLinkLine(line string) bool {
	s := strings.TrimSpace(line)
	if s == "" {
		return true
	}
	lower := strings.ToLower(s)
	return strings.HasPrefix(lower, "- en:") || strings.HasPrefix(lower, "- es:")
}

func isFacebookHashtagOnlyLine(line string) bool {
	s := strings.TrimSpace(line)
	if s == "" || !strings.HasPrefix(s, "#") {
		return false
	}
	for _, part := range strings.Fields(s) {
		if !strings.HasPrefix(part, "#") {
			return false
		}
	}
	return true
}

// SubstackBodyResult is the resolved newsletter body for a bundle locale.
type SubstackBodyResult struct {
	IndexRaw []byte
	Source   SubstackBodySource
}

// ResolveSubstackBody merges index front matter with the required substack.md or substack.es.md
// body beside the bundle. There is no fallback to index body or facebook-* sidecars.
func ResolveSubstackBody(indexRaw []byte, markdownPath string) (SubstackBodyResult, error) {
	out := SubstackBodyResult{IndexRaw: indexRaw, Source: SubstackBodyFromIndex}
	bundleDir, err := BundleDirFromMarkdownPath(markdownPath)
	if err != nil {
		return out, err
	}
	name := SubstackMarkdownBasename(markdownPath)
	md, ok, err := readBundleFile(bundleDir, name)
	if err != nil {
		return out, err
	}
	if !ok {
		return out, fmt.Errorf("substackhtml: missing %s (required for Substack; see .cursor/skills/site-substack-post/SKILL.md)", name)
	}
	merged, err := mergeIndexFrontMatterWithBody(indexRaw, md)
	if err != nil {
		return out, err
	}
	out.IndexRaw = merged
	out.Source = SubstackBodyFromSidecarMD
	return out, nil
}

func readBundleFile(bundleDir, name string) (string, bool, error) {
	b, err := os.ReadFile(filepath.Join(bundleDir, name))
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}
	s := strings.TrimSpace(string(b))
	if s == "" {
		return "", false, nil
	}
	return s, true, nil
}

func mergeIndexFrontMatterWithBody(indexRaw []byte, body string) ([]byte, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return indexRaw, nil
	}
	block, _, ok := FrontMatterBlock(indexRaw)
	if !ok {
		return append([]byte(body), '\n'), nil
	}
	var b strings.Builder
	b.Write(block)
	b.WriteString("\n\n")
	b.WriteString(body)
	b.WriteByte('\n')
	return []byte(b.String()), nil
}

// ApplySubstackSidecar merges index front matter with substack.md or substack.es.md when present.
// The bool is true when the sidecar replaced the index body (same as SubstackBodyFromSidecarMD).
func ApplySubstackSidecar(indexRaw []byte, markdownPath string) ([]byte, bool, error) {
	res, err := ResolveSubstackBody(indexRaw, markdownPath)
	if err != nil {
		return indexRaw, false, err
	}
	if res.Source == SubstackBodyFromSidecarMD {
		return res.IndexRaw, true, nil
	}
	return indexRaw, false, nil
}
