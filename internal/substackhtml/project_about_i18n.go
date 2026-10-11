package substackhtml

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/BurntSushi/toml"
)

// projectAboutCopy holds the "But why" strings used for Substack paste HTML,
// loaded from the same Hugo i18n tables as layouts/partials/*-project-about.html.
type projectAboutCopy struct {
	SayingsTitle, SayingsP1, SayingsP2 string
	CowsTitle, CowsBody                string
	ReptoTitle, ReptoBody              string
	ReptoCtaTitle, ReptoCtaButton      string
	PawTitle, PawBody                  string
}

// hugoI18nProjectAbout matches the [table] / other = "..." blocks in i18n/en.toml and i18n/es.toml.
// Extra keys in those files are ignored.
type hugoI18nProjectAbout struct {
	SayingsProjectAboutTitle struct {
		Other string `toml:"other"`
	} `toml:"sayingsProjectAboutTitle"`
	SayingsProjectAboutP1 struct {
		Other string `toml:"other"`
	} `toml:"sayingsProjectAboutP1"`
	SayingsProjectAboutP2 struct {
		Other string `toml:"other"`
	} `toml:"sayingsProjectAboutP2"`
	CowsProjectAboutTitle struct {
		Other string `toml:"other"`
	} `toml:"cowsProjectAboutTitle"`
	CowsProjectAboutBody struct {
		Other string `toml:"other"`
	} `toml:"cowsProjectAboutBody"`
	ReptilocracyProjectAboutTitle struct {
		Other string `toml:"other"`
	} `toml:"reptilocracyProjectAboutTitle"`
	ReptilocracyProjectAboutBody struct {
		Other string `toml:"other"`
	} `toml:"reptilocracyProjectAboutBody"`
	ReptilocracyProjectAboutCtaTitle struct {
		Other string `toml:"other"`
	} `toml:"reptilocracyProjectAboutCtaTitle"`
	ReptilocracyProjectAboutCtaButton struct {
		Other string `toml:"other"`
	} `toml:"reptilocracyProjectAboutCtaButton"`
	PawtropolisProjectAboutTitle struct {
		Other string `toml:"other"`
	} `toml:"pawtropolisProjectAboutTitle"`
	PawtropolisProjectAboutBody struct {
		Other string `toml:"other"`
	} `toml:"pawtropolisProjectAboutBody"`
}

var (
	projectAboutCache struct {
		en, es projectAboutCopy
	}
	projectAboutOnce sync.Once
	projectAboutErr  error
)

func findModuleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		st, err := os.Stat(filepath.Join(dir, "go.mod"))
		if err == nil && !st.IsDir() {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found when searching upward from %q", dir)
		}
		dir = parent
	}
}

func decodeProjectAboutFile(path string) (projectAboutCopy, error) {
	var raw hugoI18nProjectAbout
	if _, err := toml.DecodeFile(path, &raw); err != nil {
		return projectAboutCopy{}, fmt.Errorf("decode %s: %w", path, err)
	}
	out := projectAboutCopy{
		SayingsTitle:   raw.SayingsProjectAboutTitle.Other,
		SayingsP1:      raw.SayingsProjectAboutP1.Other,
		SayingsP2:      raw.SayingsProjectAboutP2.Other,
		CowsTitle:      raw.CowsProjectAboutTitle.Other,
		CowsBody:       raw.CowsProjectAboutBody.Other,
		ReptoTitle:     raw.ReptilocracyProjectAboutTitle.Other,
		ReptoBody:      raw.ReptilocracyProjectAboutBody.Other,
		ReptoCtaTitle:  raw.ReptilocracyProjectAboutCtaTitle.Other,
		ReptoCtaButton: raw.ReptilocracyProjectAboutCtaButton.Other,
		PawTitle:       raw.PawtropolisProjectAboutTitle.Other,
		PawBody:        raw.PawtropolisProjectAboutBody.Other,
	}
	if err := validateProjectAboutCopy(out, path); err != nil {
		return projectAboutCopy{}, err
	}
	return out, nil
}

func validateProjectAboutCopy(c projectAboutCopy, path string) error {
	check := func(name, s string) error {
		if s == "" {
			return fmt.Errorf("%s: missing or empty i18n string for %q", path, name)
		}
		return nil
	}
	for _, p := range []struct {
		name, val string
	}{
		{"sayingsProjectAboutTitle", c.SayingsTitle},
		{"sayingsProjectAboutP1", c.SayingsP1},
		{"sayingsProjectAboutP2", c.SayingsP2},
		{"cowsProjectAboutTitle", c.CowsTitle},
		{"cowsProjectAboutBody", c.CowsBody},
		{"reptilocracyProjectAboutTitle", c.ReptoTitle},
		{"reptilocracyProjectAboutBody", c.ReptoBody},
		{"reptilocracyProjectAboutCtaTitle", c.ReptoCtaTitle},
		{"reptilocracyProjectAboutCtaButton", c.ReptoCtaButton},
		{"pawtropolisProjectAboutTitle", c.PawTitle},
		{"pawtropolisProjectAboutBody", c.PawBody},
	} {
		if err := check(p.name, p.val); err != nil {
			return err
		}
	}
	return nil
}

// ensureProjectAboutI18n loads i18n/en.toml and i18n/es.toml once (from the module root that contains go.mod).
func ensureProjectAboutI18n() error {
	projectAboutOnce.Do(func() {
		root, err := findModuleRoot()
		if err != nil {
			projectAboutErr = err
			return
		}
		i18n := filepath.Join(root, "i18n")
		enPath := filepath.Join(i18n, "en.toml")
		esPath := filepath.Join(i18n, "es.toml")
		projectAboutCache.en, projectAboutErr = decodeProjectAboutFile(enPath)
		if projectAboutErr != nil {
			return
		}
		projectAboutCache.es, projectAboutErr = decodeProjectAboutFile(esPath)
	})
	return projectAboutErr
}

func projectAboutForLang(spanish bool) *projectAboutCopy {
	if spanish {
		return &projectAboutCache.es
	}
	return &projectAboutCache.en
}
