package tagregister

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// FrontMatterYAML returns the YAML document between the first and second "---" lines.
func FrontMatterYAML(content []byte) ([]byte, error) {
	br := bytes.NewReader(content)
	sc := bufio.NewScanner(br)
	if !sc.Scan() {
		return nil, fmt.Errorf("empty file")
	}
	if strings.TrimSpace(sc.Text()) != "---" {
		return nil, fmt.Errorf("missing opening front matter delimiter")
	}
	var fm []string
	for sc.Scan() {
		line := sc.Text()
		if line == "---" {
			return []byte(strings.Join(fm, "\n")), nil
		}
		fm = append(fm, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("unclosed front matter")
}

// TagsFromYAML extracts the tags key from a YAML front matter blob.
// Each tag string is trimmed; empty strings are dropped. Order is preserved for dedup per file.
func TagsFromYAML(yamlBytes []byte) ([]string, error) {
	if len(bytes.TrimSpace(yamlBytes)) == 0 {
		return nil, nil
	}
	var fm struct {
		Tags any `yaml:"tags"`
	}
	if err := yaml.Unmarshal(yamlBytes, &fm); err != nil {
		return nil, fmt.Errorf("yaml: %w", err)
	}
	return normalizeTagsList(fm.Tags), nil
}

func normalizeTagsList(raw any) []string {
	switch v := raw.(type) {
	case nil:
		return nil
	case string:
		s := strings.TrimSpace(v)
		if s == "" {
			return nil
		}
		return []string{s}
	case []any:
		var out []string
		for _, item := range v {
			for _, t := range normalizeTagsList(item) {
				if t != "" {
					out = append(out, t)
				}
			}
		}
		return out
	default:
		return nil
	}
}

// UniqueStrings preserves first-seen order, drops duplicates and empty strings.
func UniqueStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	var out []string
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}
