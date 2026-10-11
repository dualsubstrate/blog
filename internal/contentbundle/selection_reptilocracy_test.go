package contentbundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dualsubstrate/blog/internal/tagregister"
	"gopkg.in/yaml.v3"
)

func TestReptilocracyMay10FrontMatterParsesForDate(t *testing.T) {
	repoRoot := findRepoRoot(t)
	path := filepath.Join(repoRoot, "content", "cognitive-memetics", "reptilocracy", "2026-05-10-principles-with-escape-hatches", "index.md")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Skip("reptilocracy may 10 bundle not in checkout:", err)
	}
	fm, err := tagregister.FrontMatterYAML(raw)
	if err != nil {
		t.Fatalf("FrontMatterYAML: %v", err)
	}
	var doc struct {
		Date string `yaml:"date"`
	}
	if err := yaml.Unmarshal(fm, &doc); err != nil {
		t.Fatalf("yaml.Unmarshal: %v", err)
	}
	tm, err := time.Parse(time.RFC3339, strings.TrimSpace(doc.Date))
	if err != nil {
		t.Fatalf("time.Parse(%q): %v", doc.Date, err)
	}
	if got := tm.Format("2006-01-02"); got != "2026-05-10" {
		t.Fatalf("wall date got %q want 2026-05-10", got)
	}
}

func TestPublishedBundleRelsForDate20260510TwoCognitiveHubs(t *testing.T) {
	repoRoot := findRepoRoot(t)
	reptPath := filepath.Join(repoRoot, "content", "cognitive-memetics", "reptilocracy", "2026-05-10-principles-with-escape-hatches", "index.md")
	if _, err := os.Stat(reptPath); err != nil {
		t.Skip("reptilocracy 2026-05-10 bundle not in checkout; cannot assert two-bundle day")
	}
	rels, err := PublishedBundleRelsForDate(repoRoot, "2026-05-10", "")
	if err != nil {
		t.Fatal(err)
	}
	var hasPsych, hasRept bool
	for _, r := range rels {
		if strings.Contains(r, "psych-fitness-28") && strings.Contains(r, "2026-05-10-day-09") {
			hasPsych = true
		}
		if strings.Contains(r, "reptilocracy") && strings.Contains(r, "principles-with-escape-hatches") {
			hasRept = true
		}
	}
	if !hasPsych {
		t.Fatalf("missing psych bundle, rels=%v", rels)
	}
	if !hasRept {
		t.Fatalf("missing reptilocracy bundle, rels=%v (if file exists on disk, selection or YAML is wrong)", rels)
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for d := dir; d != "" && d != "/"; d = filepath.Dir(d) {
		if st, err := os.Stat(filepath.Join(d, "go.mod")); err == nil && !st.IsDir() {
			if _, err := os.Stat(filepath.Join(d, "content")); err == nil {
				return d
			}
		}
	}
	t.Fatal("could not find repo root (go.mod + content/)")
	return ""
}
