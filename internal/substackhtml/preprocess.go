package substackhtml

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	reHugoShortcodeAngle   = regexp.MustCompile(`\{\{<[\s\S]*?>\}\}`)
	reHugoShortcodePercent = regexp.MustCompile(`\{\{%[\s\S]*?%\}\}`)
	reHTMLComment          = regexp.MustCompile(`<!--[\s\S]*?-->`)
	// YouTube shortcode common forms: {{< youtube id >}} or {{< youtube "id" >}}.
	reYouTubeShortcode = regexp.MustCompile(`\{\{<\s*youtube\s+["']?([A-Za-z0-9_-]{6,20})["']?\s*>\}\}`)
	// {{< mermaidfile >}} or {{< mermaidfile "diagram.mmd" >}} (same as layouts/shortcodes/mermaidfile.html).
	reMermaidfileShortcode = regexp.MustCompile(`\{\{<\s*mermaidfile(?:\s+["']([^"']+)["'])?\s*>\}\}`)
)

// PreprocessMarkdown removes Hugo shortcodes and HTML comments so goldmark does
// not see raw braces as text. Known media shortcodes become plain paragraphs
// with a normal link so paste stays readable.
// MarkdownSourcePath is the absolute or relative path to the Markdown file (for example -in).
// When set, {{< mermaidfile "name.mmd" >}} is inlined as a fenced mermaid block from the same directory.
func PreprocessMarkdown(src []byte, markdownSourcePath string) []byte {
	s := src
	s = reHTMLComment.ReplaceAll(s, nil)
	s = expandMermaidfileShortcodes(s, markdownSourcePath)
	s = replaceYouTubeShortcodes(s)
	s = reHugoShortcodeAngle.ReplaceAll(s, nil)
	s = reHugoShortcodePercent.ReplaceAll(s, nil)
	s = StripBoilerplateLines(s)
	s = bytes.TrimSpace(s)
	return s
}

func expandMermaidfileShortcodes(src []byte, mdPath string) []byte {
	mdPath = strings.TrimSpace(mdPath)
	if mdPath == "" {
		return src
	}
	return reMermaidfileShortcode.ReplaceAllFunc(src, func(match []byte) []byte {
		sub := reMermaidfileShortcode.FindSubmatch(match)
		fname := "diagram.mmd"
		if len(sub) > 1 && len(bytes.TrimSpace(sub[1])) > 0 {
			fname = string(bytes.TrimSpace(sub[1]))
		}
		full, err := resolveMermaidfilePath(mdPath, fname)
		if err != nil {
			return []byte("\n\n<!-- substackhtml: mermaidfile " + fname + ": invalid filename -->\n\n")
		}
		if webpName := mermaidSourceToWebpFilename(fname); webpName != "" {
			webpFull, werr := resolveMermaidfilePath(mdPath, webpName)
			if werr == nil {
				if st, statErr := os.Stat(webpFull); statErr == nil && !st.IsDir() {
					line := "\n\n![Diagram](" + webpName + ")\n\n"
					return []byte(line)
				}
			}
		}
		data, err := os.ReadFile(full)
		if err != nil {
			return []byte("\n\n<!-- substackhtml: mermaidfile " + fname + ": " + err.Error() + " -->\n\n")
		}
		body := bytes.TrimSpace(data)
		out := make([]byte, 0, len(body)+20)
		out = append(out, []byte("\n```mermaid\n")...)
		out = append(out, body...)
		out = append(out, []byte("\n```\n")...)
		return out
	})
}

// resolveMermaidfilePath returns a diagram path in the same directory as mdPath (no path traversal).
// Fname must be a single path segment (no separators, no "..").
func resolveMermaidfilePath(mdPath, fname string) (string, error) {
	fname = strings.TrimSpace(fname)
	if fname == "" {
		return "", errors.New("empty filename")
	}
	if filepath.IsAbs(fname) {
		return "", errors.New("absolute filename")
	}
	if fname != filepath.Base(fname) {
		return "", errors.New("path separators in filename")
	}
	if fname == "." || fname == ".." {
		return "", errors.New("invalid filename")
	}
	if strings.Contains(fname, "..") {
		return "", errors.New("invalid filename")
	}
	absMD, err := filepath.Abs(mdPath)
	if err != nil {
		return "", err
	}
	absDir := filepath.Clean(filepath.Dir(absMD))
	joined := filepath.Join(absDir, fname)
	cleanJoined := filepath.Clean(joined)
	rel, err := filepath.Rel(absDir, cleanJoined)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", errors.New("path escapes markdown directory")
	}
	return cleanJoined, nil
}

// mermaidSourceToWebpFilename returns a sibling .webp name for a .mmd diagram file
// (for example diagram.mmd -> diagram.webp, diagram.es.mmd -> diagram.es.webp).
// It returns empty when fname does not end with .mmd (case-insensitive).
func mermaidSourceToWebpFilename(fname string) string {
	fname = strings.TrimSpace(fname)
	if fname == "" || !strings.HasSuffix(strings.ToLower(fname), ".mmd") {
		return ""
	}
	return fname[:len(fname)-4] + ".webp"
}

func replaceYouTubeShortcodes(s []byte) []byte {
	return reYouTubeShortcode.ReplaceAllFunc(s, func(m []byte) []byte {
		sub := reYouTubeShortcode.FindSubmatch(m)
		if len(sub) < 2 {
			return nil
		}
		id := string(sub[1])
		url := "https://www.youtube.com/watch?v=" + id
		line := "\n\nWatch on YouTube: " + url + "\n\n"
		return []byte(line)
	})
}

// StripBoilerplateLines removes whole lines that match known editor-injected
// markers (optional second pass for local tooling). It does not parse Markdown.
func StripBoilerplateLines(src []byte) []byte {
	lines := bytes.Split(src, []byte("\n"))
	var out [][]byte
	for _, ln := range lines {
		ls := string(bytes.TrimSpace(ln))
		if strings.HasPrefix(ls, "The following cursor rule files are relevant") {
			break
		}
		if strings.Contains(ls, ".cursor/rules/") && strings.Contains(ls, "relevant to the files you just read") {
			break
		}
		out = append(out, ln)
	}
	return bytes.Join(out, []byte("\n"))
}
