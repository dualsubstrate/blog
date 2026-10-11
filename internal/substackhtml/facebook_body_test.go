package substackhtml

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStripFacebookSocialTailKeepsCitationAfterReadMore(t *testing.T) {
	in := strings.Join([]string{
		"Hook line.",
		"",
		"Body paragraph.",
		"",
		"#TagOne #TagTwo",
		"",
		"If you want to read more:",
		"",
		"- EN: https://example.com/en/",
		"",
		"- ES: https://example.com/es/",
		"",
		"Casual read (arXiv): https://arxiv.org/abs/1212.0141",
	}, "\n")
	got := StripFacebookSocialTail(in)
	if strings.Contains(got, "#TagOne") {
		t.Fatalf("hashtags should be stripped: %q", got)
	}
	if strings.Contains(got, "If you want to read more") {
		t.Fatalf("read-more header should be stripped: %q", got)
	}
	if strings.Contains(got, "example.com/en") {
		t.Fatalf("read-more URLs should be stripped: %q", got)
	}
	if !strings.Contains(got, "arxiv.org") {
		t.Fatalf("citation after read-more should remain: %q", got)
	}
}

func TestResolveSubstackBodyRequiresSidecar(t *testing.T) {
	dir := t.TempDir()
	index := filepath.Join(dir, "index.md")
	fm := "---\ntitle: T\ntype: claims\n---\n\nSite.\n"
	if err := os.WriteFile(index, []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ResolveSubstackBody([]byte(fm), index)
	if err == nil || !strings.Contains(err.Error(), "substack.md") {
		t.Fatalf("want missing substack.md error, got %v", err)
	}
}

func TestResolveSubstackBodyMergesSubstackMD(t *testing.T) {
	dir := t.TempDir()
	index := filepath.Join(dir, "index.md")
	fm := "---\ntitle: T\ntype: claims\n---\n\nSite.\n"
	if err := os.WriteFile(index, []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "substack.md"), []byte("## Hook\n\n**Bold** line.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "facebook-en.txt"), []byte("Ignored txt.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := ResolveSubstackBody([]byte(fm), index)
	if err != nil {
		t.Fatal(err)
	}
	if res.Source != SubstackBodyFromSidecarMD {
		t.Fatalf("source=%v want sidecar md", res.Source)
	}
	if !strings.Contains(string(res.IndexRaw), "## Hook") {
		t.Fatalf("expected substack body: %q", res.IndexRaw)
	}
	if strings.Contains(string(res.IndexRaw), "Ignored") {
		t.Fatalf("facebook txt must not be used: %q", res.IndexRaw)
	}
}

func TestResolveSubstackBodySpanishBasename(t *testing.T) {
	dir := t.TempDir()
	index := filepath.Join(dir, "index.es.md")
	fm := "---\ntitle: ES\ntype: sayings\n---\n\nSite body.\n"
	if err := os.WriteFile(index, []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "substack.es.md"), []byte("## Gancho\n\nCuerpo.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	res, err := ResolveSubstackBody([]byte(fm), index)
	if err != nil {
		t.Fatal(err)
	}
	if res.Source != SubstackBodyFromSidecarMD {
		t.Fatal("expected spanish sidecar")
	}
	if !strings.Contains(string(res.IndexRaw), "## Gancho") {
		t.Fatalf("got %q", res.IndexRaw)
	}
}

func TestSubstackMarkdownBasename(t *testing.T) {
	if got := SubstackMarkdownBasename("content/x/index.md"); got != "substack.md" {
		t.Fatalf("en: got %q", got)
	}
	if got := SubstackMarkdownBasename("content/x/index.es.md"); got != "substack.es.md" {
		t.Fatalf("es: got %q", got)
	}
}

func TestFacebookCopyBasename(t *testing.T) {
	if got := FacebookCopyBasename("content/x/index.md"); got != "facebook-en.txt" {
		t.Fatalf("en: got %q", got)
	}
	if got := FacebookCopyBasename("content/x/index.es.md"); got != "facebook-es.txt" {
		t.Fatalf("es: got %q", got)
	}
}
