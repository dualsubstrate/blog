package contentbundle

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeBundleRel(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"human-condition/2026-05-01-ego-as-game", "human-condition/2026-05-01-ego-as-game"},
		{"/content/human-condition/2026-05-01-ego-as-game/", "human-condition/2026-05-01-ego-as-game"},
	}
	for _, tc := range cases {
		if got := normalizeBundleRel(tc.in); got != tc.want {
			t.Fatalf("normalizeBundleRel(%q) = %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestPublishedBundleRelsForDateEmptyIsNotError(t *testing.T) {
	root := t.TempDir()
	content := filepath.Join(root, "content", "human-condition", "2099-01-01-placeholder")
	if err := os.MkdirAll(content, 0o755); err != nil {
		t.Fatal(err)
	}
	index := []byte("---\ndate: '2099-01-01T01:00:00+11:00'\ntitle: placeholder\n---\n")
	if err := os.WriteFile(filepath.Join(content, "index.md"), index, 0o644); err != nil {
		t.Fatal(err)
	}
	rels, err := PublishedBundleRelsForDate(root, "1990-01-01", "")
	if err != nil {
		t.Fatalf("expected nil error on empty date, got %v", err)
	}
	if len(rels) != 0 {
		t.Fatalf("expected empty rels, got %v", rels)
	}
}
