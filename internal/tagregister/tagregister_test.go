package tagregister

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFrontMatterYAML_andTags(t *testing.T) {
	const md = `---
title: "x"
tags: ["A", "B", "A"]
---
body
`
	fm, err := FrontMatterYAML([]byte(md))
	if err != nil {
		t.Fatal(err)
	}
	tags, err := TagsFromYAML(fm)
	if err != nil {
		t.Fatal(err)
	}
	got := UniqueStrings(tags)
	if len(got) != 2 || got[0] != "A" || got[1] != "B" {
		t.Fatalf("tags: %#v", got)
	}
}

func TestTagsFromYAML_string(t *testing.T) {
	const y = `tags: SingleTag
`
	tags, err := TagsFromYAML([]byte(y))
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 || tags[0] != "SingleTag" {
		t.Fatalf("%#v", tags)
	}
}

func TestScanMarkdownTags(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "a"), 0o755); err != nil {
		t.Fatal(err)
	}
	one := filepath.Join(dir, "a", "one.md")
	two := filepath.Join(dir, "two.md")
	if err := os.WriteFile(one, []byte("---\ntags: [X, Y]\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(two, []byte("---\ntags: [X]\n---\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	counts, err := ScanMarkdownTags(dir)
	if err != nil {
		t.Fatal(err)
	}
	if counts["X"] != 2 || counts["Y"] != 1 {
		t.Fatalf("%v", counts)
	}
}

func TestLoadDeprecations_duplicateFrom(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "d.toml")
	body := `[[deprecated]]
from = "a"
to = "b"
[[deprecated]]
from = "a"
to = "c"
`
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := LoadDeprecations(p)
	if err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("want duplicate error, got %v", err)
	}
}

func TestRenderRegister(t *testing.T) {
	counts := map[string]int{"Old": 2, "New": 5}
	deps := []Deprecation{{From: "Old", To: "New"}}
	s := RenderRegister(counts, deps)
	if !strings.Contains(s, "Old\tNew\t2") {
		t.Fatalf("deprecated line missing: %q", s)
	}
	if !strings.Contains(s, "New\t5") {
		t.Fatalf("active line missing: %q", s)
	}
}
