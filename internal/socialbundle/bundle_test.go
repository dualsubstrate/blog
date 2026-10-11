package socialbundle

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadBundleMissingLinkedIn(t *testing.T) {
	root := findRepoRoot(t)
	contentRoot := filepath.Join(root, "content")
	_, err := LoadBundle(contentRoot, "this-bundle-does-not-exist-xyz")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "missing index.md") {
		t.Fatalf("unexpected: %v", err)
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
	t.Fatal("repo root not found")
	return ""
}
