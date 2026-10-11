package substackhtml

import (
	"bytes"
	"strings"
)

// FrontMatterMeta is the minimal metadata we care about for Substack draft fields.
// It is intentionally tiny and YAML-lite so we avoid pulling a full YAML parser.
type FrontMatterMeta struct {
	Title       string
	Date        string
	Lang        string
	Description string
	SoWhat      string
	YouTubeID   string
	Type        string
	ImageURL    string
	Tags        []string
	Categories  []string
	Grounding   string
	TLDR        string
	Fluff       string
}

// ExtractFrontMatterMeta tries to read a leading YAML front matter block and return
// a few known keys (title, description, sowhat, youtube_id, type, images[0]).
// It never errors; it falls back to empty.
func ExtractFrontMatterMeta(src []byte) FrontMatterMeta {
	s := bytes.TrimPrefix(src, []byte{0xEF, 0xBB, 0xBF}) // UTF-8 BOM
	// Allow leading whitespace/newlines before front matter.
	s = bytes.TrimLeft(s, " \t\r\n")
	if !bytes.HasPrefix(s, []byte("---")) {
		return FrontMatterMeta{}
	}
	_, rest, ok := cutLine(s)
	if !ok {
		return FrontMatterMeta{}
	}

	lines := [][]byte{}
	for {
		line, after, ok := cutLine(rest)
		if !ok {
			return FrontMatterMeta{}
		}
		if bytes.Equal(bytes.TrimSpace(line), []byte("---")) {
			break
		}
		lines = append(lines, line)
		rest = after
	}

	getScalar := func(key string) string {
		prefix := []byte(key + ":")
		for i := 0; i < len(lines); i++ {
			ln := lines[i]
			trim := bytes.TrimSpace(ln)
			if !bytes.HasPrefix(trim, prefix) {
				continue
			}
			after := strings.TrimSpace(string(bytes.TrimSpace(trim[len(prefix):])))
			// Block scalar: key: |
			if after == "|" || after == "|-" || after == "|+" {
				var out []string
				for j := i + 1; j < len(lines); j++ {
					raw := string(lines[j])
					if strings.TrimSpace(raw) == "" {
						// keep blank lines inside blocks
						out = append(out, "")
						continue
					}
					if !strings.HasPrefix(raw, " ") && !strings.HasPrefix(raw, "\t") {
						break
					}
					out = append(out, strings.TrimRight(raw, "\r\n"))
				}
				block := strings.Join(out, "\n")
				block = strings.TrimSpace(dedent(block))
				return block
			}
			return unquote(after)
		}
		return ""
	}

	getFirstListItem := func(key string) string {
		needle := []byte(key + ":")
		for i := 0; i < len(lines); i++ {
			trim := bytes.TrimSpace(lines[i])
			if !bytes.HasPrefix(trim, needle) {
				continue
			}
			for j := i + 1; j < len(lines); j++ {
				raw := string(lines[j])
				if strings.TrimSpace(raw) == "" {
					continue
				}
				// End of the list when indentation stops.
				if !strings.HasPrefix(raw, " ") && !strings.HasPrefix(raw, "\t") {
					break
				}
				t := strings.TrimSpace(raw)
				if strings.HasPrefix(t, "-") {
					item := strings.TrimSpace(strings.TrimPrefix(t, "-"))
					return unquote(item)
				}
			}
			return ""
		}
		return ""
	}

	getListItems := func(key string) []string {
		needle := []byte(key + ":")
		for i := 0; i < len(lines); i++ {
			trim := bytes.TrimSpace(lines[i])
			if !bytes.HasPrefix(trim, needle) {
				continue
			}
			var out []string
			for j := i + 1; j < len(lines); j++ {
				raw := string(lines[j])
				if strings.TrimSpace(raw) == "" {
					continue
				}
				// End of the list when indentation stops.
				if !strings.HasPrefix(raw, " ") && !strings.HasPrefix(raw, "\t") {
					break
				}
				t := strings.TrimSpace(raw)
				if strings.HasPrefix(t, "-") {
					item := strings.TrimSpace(strings.TrimPrefix(t, "-"))
					item = unquote(item)
					if item != "" {
						out = append(out, item)
					}
				}
			}
			return out
		}
		return nil
	}

	tags := getInlineStringListFromLines(lines, "tags")
	if len(tags) == 0 {
		tags = getListItems("tags")
	}
	categories := getInlineStringListFromLines(lines, "categories")
	if len(categories) == 0 {
		categories = getListItems("categories")
	}

	imageURL := strings.TrimSpace(getFirstListItem("images"))
	if imageURL == "" {
		imageURL = strings.TrimSpace(getScalar("featuredImage"))
	}

	return FrontMatterMeta{
		Title:       strings.TrimSpace(getScalar("title")),
		Date:        strings.TrimSpace(getScalar("date")),
		Lang:        strings.TrimSpace(getScalar("lang")),
		Description: strings.TrimSpace(getScalar("description")),
		SoWhat:      strings.TrimSpace(getScalar("sowhat")),
		YouTubeID:   strings.TrimSpace(getScalar("youtube_id")),
		Type:        strings.TrimSpace(getScalar("type")),
		ImageURL:    imageURL,
		Tags:        tags,
		Categories:  categories,
		Grounding:   strings.TrimSpace(getScalar("grounding")),
		TLDR:        strings.TrimSpace(getScalar("tldr")),
		Fluff:       strings.TrimSpace(getScalar("fluff")),
	}
}

func getInlineStringListFromLines(lines [][]byte, key string) []string {
	prefix := []byte(key + ":")
	for i := 0; i < len(lines); i++ {
		trim := bytes.TrimSpace(lines[i])
		if !bytes.HasPrefix(trim, prefix) {
			continue
		}
		after := strings.TrimSpace(string(bytes.TrimSpace(trim[len(prefix):])))
		if !strings.HasPrefix(after, "[") || !strings.HasSuffix(after, "]") {
			return nil
		}
		body := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(after, "["), "]"))
		if body == "" {
			return nil
		}
		return splitInlineYAMLList(body)
	}
	return nil
}

func splitInlineYAMLList(body string) []string {
	// Parse a YAML/JSON-style inline list body without splitting on commas inside quotes.
	// Examples:
	// - `"a", "b"` -> ["a","b"]
	// - `"tag, with comma", "tag2"` -> ["tag, with comma","tag2"]
	var out []string
	var buf strings.Builder
	inSingle := false
	inDouble := false
	escape := false

	flush := func() {
		t := strings.TrimSpace(buf.String())
		t = unquote(t)
		if t != "" {
			out = append(out, t)
		}
		buf.Reset()
	}

	for i := 0; i < len(body); i++ {
		ch := body[i]
		if escape {
			buf.WriteByte(ch)
			escape = false
			continue
		}
		if inDouble && ch == '\\' {
			// YAML/JSON double-quoted strings can escape quotes and commas.
			buf.WriteByte(ch)
			escape = true
			continue
		}
		if ch == '"' && !inSingle {
			inDouble = !inDouble
			buf.WriteByte(ch)
			continue
		}
		if ch == '\'' && !inDouble {
			inSingle = !inSingle
			buf.WriteByte(ch)
			continue
		}
		if ch == ',' && !inSingle && !inDouble {
			flush()
			continue
		}
		buf.WriteByte(ch)
	}
	flush()
	return out
}

func unquote(s string) string {
	ss := strings.TrimSpace(s)
	if len(ss) >= 2 {
		if (ss[0] == '"' && ss[len(ss)-1] == '"') || (ss[0] == '\'' && ss[len(ss)-1] == '\'') {
			return ss[1 : len(ss)-1]
		}
	}
	return ss
}

func dedent(s string) string {
	lines := strings.Split(s, "\n")
	min := -1
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		n := 0
		for n < len(ln) && (ln[n] == ' ' || ln[n] == '\t') {
			n++
		}
		if min == -1 || n < min {
			min = n
		}
	}
	if min <= 0 {
		return s
	}
	for i, ln := range lines {
		if len(ln) >= min {
			lines[i] = ln[min:]
		}
	}
	return strings.Join(lines, "\n")
}
