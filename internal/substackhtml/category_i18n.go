package substackhtml

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

// i18nOtherEntry matches Hugo LoveIt [CategoryName] / other = "..." blocks in i18n/*.toml.
type i18nOtherEntry struct {
	Other string `toml:"other"`
}

var (
	categoryI18nES   map[string]string
	categoryI18nOnce sync.Once
	categoryI18nErr  error
)

func loadCategoryI18nMaps() error {
	categoryI18nOnce.Do(func() {
		root, err := findModuleRoot()
		if err != nil {
			categoryI18nErr = err
			return
		}
		categoryI18nES, categoryI18nErr = decodeCategoryI18nFile(filepath.Join(root, "i18n", "es.toml"))
	})
	return categoryI18nErr
}

func decodeCategoryI18nFile(path string) (map[string]string, error) {
	var raw map[string]i18nOtherEntry
	if _, err := toml.DecodeFile(path, &raw); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(raw))
	for key, entry := range raw {
		label := strings.TrimSpace(entry.Other)
		if label == "" {
			continue
		}
		out[key] = label
	}
	return out, nil
}

// CategorySubstackSectionLabel returns the Substack publish-modal section name for a Hugo
// categories[] value. For Spanish drafts (index.es.md or lang es), it uses i18n/es.toml
// display names (e.g. Social-Protocols → Protocolos-sociales). For English, it returns
// the category id unchanged so it can match the EN publication section labels.
func CategorySubstackSectionLabel(hugoCategory string, spanish bool) string {
	name := strings.TrimSpace(hugoCategory)
	if name == "" {
		return ""
	}
	if !spanish {
		return name
	}
	if err := loadCategoryI18nMaps(); err != nil {
		return name
	}
	m := categoryI18nES
	if label, ok := m[name]; ok && label != "" {
		return label
	}
	// Case-insensitive fallback (Hugo taxonomy ids are usually canonical-cased).
	lower := strings.ToLower(name)
	for k, v := range m {
		if strings.ToLower(k) == lower && v != "" {
			return v
		}
	}
	return name
}
