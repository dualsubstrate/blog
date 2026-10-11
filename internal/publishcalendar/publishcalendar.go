// Package publishcalendar builds a JSON snapshot of planned publish dates and
// per-channel file presence from Hugo page bundles.
package publishcalendar

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dualsubstrate/blog/internal/substackpublishstate"
	"github.com/dualsubstrate/blog/internal/tagregister"
	"gopkg.in/yaml.v3"
)

const (
	StatusMissing = "missing"
	StatusPresent = "present"
)

// DefaultChannels is the channel list shown in the calendar UI.
var DefaultChannels = []string{
	substackpublishstate.TargetSiteEN,
	substackpublishstate.TargetSiteES,
	substackpublishstate.TargetLinkedIn,
	substackpublishstate.TargetFacebook,
	substackpublishstate.TargetSubstackEN,
}

// ChannelState is one channel row for a bundle.
type ChannelState struct {
	Status string `json:"status"`
}

// Entry is one Hugo page bundle on the calendar.
type Entry struct {
	Bundle     string                  `json:"bundle"`
	Section    string                  `json:"section"`
	Title      string                  `json:"title"`
	Type       string                  `json:"type,omitempty"`
	Categories []string                `json:"categories,omitempty"`
	Planned    string                  `json:"planned"`
	Draft      bool                    `json:"draft"`
	Channels   map[string]ChannelState `json:"channels"`
}

// Calendar is the root JSON document.
type Calendar struct {
	GeneratedAt string   `json:"generatedAt"`
	Channels    []string `json:"channels"`
	Entries     []Entry  `json:"entries"`
}

type frontMatterDoc struct {
	Title      string   `yaml:"title"`
	Date       string   `yaml:"date"`
	Draft      bool     `yaml:"draft"`
	Type       string   `yaml:"type"`
	Categories []string `yaml:"categories"`
}

// Build scans contentRoot for leaf bundles with index.md and returns a calendar snapshot.
func Build(contentRoot string, channels []string) (*Calendar, error) {
	contentRoot = strings.TrimSpace(contentRoot)
	if contentRoot == "" {
		contentRoot = "content"
	}
	if len(channels) == 0 {
		channels = append([]string(nil), DefaultChannels...)
	}
	absRoot, err := filepath.Abs(contentRoot)
	if err != nil {
		return nil, fmt.Errorf("publishcalendar: content root: %w", err)
	}

	var entries []Entry
	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || d.Name() != "index.md" {
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

		entry, err := entryFromBundle(bundleDir, filepath.ToSlash(rel), parts[0], channels)
		if err != nil {
			return err
		}
		if entry == nil {
			return nil
		}
		entries = append(entries, *entry)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Planned != entries[j].Planned {
			return entries[i].Planned < entries[j].Planned
		}
		return entries[i].Bundle < entries[j].Bundle
	})

	return &Calendar{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Channels:    append([]string(nil), channels...),
		Entries:     entries,
	}, nil
}

func entryFromBundle(bundleDir, rel, section string, channels []string) (*Entry, error) {
	indexPath := filepath.Join(bundleDir, "index.md")
	raw, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}
	fmBytes, err := tagregister.FrontMatterYAML(raw)
	if err != nil {
		return nil, nil
	}
	var doc frontMatterDoc
	if err := yaml.Unmarshal(fmBytes, &doc); err != nil {
		return nil, nil
	}

	planned := plannedDate(doc.Date, rel)
	if planned == "" {
		return nil, nil
	}

	title := strings.TrimSpace(doc.Title)
	if title == "" {
		title = filepath.Base(bundleDir)
	}

	chState := make(map[string]ChannelState, len(channels))
	for _, ch := range channels {
		chState[ch] = channelState(bundleDir, ch)
	}

	return &Entry{
		Bundle:     rel,
		Section:    section,
		Title:      title,
		Type:       strings.TrimSpace(doc.Type),
		Categories: append([]string(nil), doc.Categories...),
		Planned:    planned,
		Draft:      doc.Draft,
		Channels:   chState,
	}, nil
}

func plannedDate(dateStr, bundleRel string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr != "" {
		if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
			return t.Format("2006-01-02")
		}
	}
	base := filepath.Base(bundleRel)
	if len(base) >= 10 && base[4] == '-' && base[7] == '-' {
		candidate := base[:10]
		if _, err := time.Parse("2006-01-02", candidate); err == nil {
			return candidate
		}
	}
	return ""
}

func channelState(bundleDir string, targetKey string) ChannelState {
	present, err := channelFilePresent(bundleDir, targetKey)
	if err != nil || !present {
		return ChannelState{Status: StatusMissing}
	}
	return ChannelState{Status: StatusPresent}
}

func channelFilePresent(bundleDir, targetKey string) (bool, error) {
	switch targetKey {
	case substackpublishstate.TargetSiteEN:
		return hugoPageFilePresent(filepath.Join(bundleDir, "index.md"))
	case substackpublishstate.TargetSiteES:
		return hugoPageFilePresent(filepath.Join(bundleDir, "index.es.md"))
	case substackpublishstate.TargetLinkedIn:
		return sidecarFilePresent(filepath.Join(bundleDir, "linkedin.txt"))
	case substackpublishstate.TargetFacebook:
		return sidecarFilePresent(filepath.Join(bundleDir, "facebook-es.txt"))
	case substackpublishstate.TargetSubstackEN:
		return sidecarFilePresent(filepath.Join(bundleDir, "substack.md"))
	default:
		return false, nil
	}
}

func sidecarFilePresent(path string) (bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(string(b)) != "", nil
}

func hugoPageFilePresent(path string) (bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if strings.TrimSpace(string(b)) == "" {
		return false, nil
	}
	fmBytes, err := tagregister.FrontMatterYAML(b)
	if err != nil {
		return true, nil
	}
	var doc struct {
		Draft bool `yaml:"draft"`
	}
	if err := yaml.Unmarshal(fmBytes, &doc); err == nil && doc.Draft {
		return false, nil
	}
	return true, nil
}

// WriteJSON encodes cal to path with a trailing newline.
func WriteJSON(path string, cal *Calendar) error {
	if cal == nil {
		return fmt.Errorf("publishcalendar: nil calendar")
	}
	data, err := json.MarshalIndent(cal, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("publishcalendar: mkdir: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("publishcalendar: write: %w", err)
	}
	return nil
}
