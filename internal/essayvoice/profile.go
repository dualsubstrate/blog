// Package essayvoice loads host voice files and runs the portable strop audit.
package essayvoice

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/behaviorengineering/strop/pkg/evaluation/voice"
	"gopkg.in/yaml.v3"
)

// File is one catalog entry: the portable profile plus the sections it serves.
type File struct {
	Profile  voice.Profile
	Sections []string
}

type frontMatter struct {
	ID                  string   `yaml:"id"`
	Label               string   `yaml:"label"`
	Sections            []string `yaml:"sections"`
	BannedPatterns      []string `yaml:"banned_patterns"`
	RequiredVerbClasses []string `yaml:"required_verb_classes"`
	MaxStaccatoRun      int      `yaml:"max_staccato_run"`
}

// Load reads data/voices/<id>.md and maps YAML front matter onto a strop profile.
func Load(ctx context.Context, dir, id string) (File, error) {
	if err := ctx.Err(); err != nil {
		return File{}, err
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return File{}, fmt.Errorf("voice id is empty")
	}
	path := filepath.Join(dir, id+".md")
	raw, err := os.ReadFile(path)
	if err != nil {
		return File{}, fmt.Errorf("read voice %s: %w", id, err)
	}
	meta, _, err := splitFrontMatter(raw)
	if err != nil {
		return File{}, fmt.Errorf("voice %s: %w", id, err)
	}
	var fm frontMatter
	if err := yaml.Unmarshal(meta, &fm); err != nil {
		return File{}, fmt.Errorf("parse voice %s: %w", id, err)
	}
	if fm.ID == "" {
		fm.ID = id
	}
	if fm.ID != id {
		return File{}, fmt.Errorf("voice file %s declares id %s", id, fm.ID)
	}
	return File{
		Profile: voice.Profile{
			ID:             fm.ID,
			Label:          fm.Label,
			BannedPatterns: fm.BannedPatterns,
			RequiredVerbs:  fm.RequiredVerbClasses,
			MaxStaccatoRun: fm.MaxStaccatoRun,
		},
		Sections: fm.Sections,
	}, nil
}

func splitFrontMatter(raw []byte) (meta, body []byte, err error) {
	raw = bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF})
	if bytes.HasPrefix(raw, []byte("---\n")) || bytes.HasPrefix(raw, []byte("---\r\n")) {
		return splitDelimited(raw, []byte("---"))
	}
	if bytes.HasPrefix(raw, []byte("+++\n")) || bytes.HasPrefix(raw, []byte("+++\r\n")) {
		return splitDelimited(raw, []byte("+++"))
	}
	return nil, raw, nil
}

func splitDelimited(raw, delim []byte) (meta, body []byte, err error) {
	rest := raw[len(delim):]
	rest = bytes.TrimPrefix(rest, []byte("\r\n"))
	rest = bytes.TrimPrefix(rest, []byte("\n"))
	idx := bytes.Index(rest, append(append([]byte("\n"), delim...), '\n'))
	idxCR := bytes.Index(rest, append(append([]byte("\r\n"), delim...), []byte("\r\n")...))
	cut := idx
	sepLen := 1 + len(delim) + 1
	if idxCR >= 0 && (cut < 0 || idxCR < cut) {
		cut = idxCR
		sepLen = 2 + len(delim) + 2
	}
	if cut < 0 {
		return nil, nil, fmt.Errorf("front matter is missing the closing delimiter")
	}
	return rest[:cut], rest[cut+sepLen:], nil
}
