package substackpublishstate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListUnpublishedAndMarkPublished(t *testing.T) {
	root := t.TempDir()
	content := filepath.Join(root, "content")
	sec := filepath.Join(content, "human-condition")
	b1 := filepath.Join(sec, "post-a")
	b2 := filepath.Join(sec, "post-b")
	if err := os.MkdirAll(b1, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(b2, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b1, "index.md"), []byte("---\ntitle: A\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b1, "substack.md"), []byte("## A\n\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b2, "index.md"), []byte("---\ntitle: B\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b2, DefaultMarker), []byte("2026-05-01T12:00:00Z\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	list, err := ListUnpublished(content, DefaultMarker, LegacyDefaultTarget)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || string(list[0]) != "human-condition/post-a" {
		t.Fatalf("got %v want [human-condition/post-a]", list)
	}

	if err := MarkPublished(b1, DefaultMarker, LegacyDefaultTarget); err != nil {
		t.Fatal(err)
	}
	list2, err := ListUnpublished(content, DefaultMarker, LegacyDefaultTarget)
	if err != nil {
		t.Fatal(err)
	}
	if len(list2) != 0 {
		t.Fatalf("after mark: got %v want []", list2)
	}
}

func TestListUnpublishedPerTarget(t *testing.T) {
	root := t.TempDir()
	content := filepath.Join(root, "content")
	b := filepath.Join(content, "human-condition", "post-x")
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.md"), []byte("---\ntitle: X\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.es.md"), []byte("---\ntitle: X\n---\ncuerpo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "substack.es.md"), []byte("## X\n\nCuerpo.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "substack.md"), []byte("## X\n\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, DefaultMarker), []byte("substack-en 2026-05-01T12:00:00Z\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	es, err := ListUnpublished(content, DefaultMarker, "substack-es")
	if err != nil {
		t.Fatal(err)
	}
	if len(es) != 1 || string(es[0]) != "human-condition/post-x" {
		t.Fatalf("substack-es unpublished: got %v", es)
	}
	en, err := ListUnpublished(content, DefaultMarker, "substack-en")
	if err != nil {
		t.Fatal(err)
	}
	if len(en) != 0 {
		t.Fatalf("substack-en should be done: got %v", en)
	}
}

func TestMarkPublishedMergesTargets(t *testing.T) {
	root := t.TempDir()
	b := filepath.Join(root, "content", "sec", "slug")
	if err := os.MkdirAll(b, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(b, "index.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := MarkPublished(b, DefaultMarker, "substack-en"); err != nil {
		t.Fatal(err)
	}
	if err := MarkPublished(b, DefaultMarker, "linkedin"); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join(b, DefaultMarker))
	if err != nil {
		t.Fatal(err)
	}
	m := ParseMarkerContent(raw)
	if len(m) != 2 || m["substack-en"] == "" || m["linkedin"] == "" {
		t.Fatalf("map: %#v", m)
	}
}

func TestListUnpublishedSkipsSectionIndexOnly(t *testing.T) {
	root := t.TempDir()
	content := filepath.Join(root, "content")
	sec := filepath.Join(content, "human-condition")
	if err := os.MkdirAll(sec, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sec, "_index.md"), []byte("---\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	list, err := ListUnpublished(content, DefaultMarker, LegacyDefaultTarget)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("got %v want []", list)
	}
}
