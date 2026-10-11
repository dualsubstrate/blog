package substackhtml

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestReplaceHTTPOrigin(t *testing.T) {
	got, err := ReplaceHTTPOrigin("https://blog.dualsubstrate.com/cognitive-memetics/foo/", "http://localhost:1313")
	if err != nil {
		t.Fatal(err)
	}
	want := "http://localhost:1313/cognitive-memetics/foo/"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
	if _, err := ReplaceHTTPOrigin("/relative-only/", "http://localhost:1313"); err == nil {
		t.Fatal("expected error for relative page URL")
	}
}

func TestResolveImageReference(t *testing.T) {
	if got := ResolveImageReference("https://x/y.png", ""); got != "https://x/y.png" {
		t.Fatalf("absolute: %q", got)
	}
	if got := ResolveImageReference("z.png", "https://site.com/post/"); got != "https://site.com/post/z.png" {
		t.Fatalf("join: %q", got)
	}
	if got := ResolveImageReference("/z.png", "https://site.com/post/"); got != "https://site.com/post/z.png" {
		t.Fatalf("strip slash: %q", got)
	}
	if got := ResolveImageReference("z.png", ""); got != "z.png" {
		t.Fatalf("no permalink keeps relative: %q", got)
	}
}

func TestLookupPermalinkCowPost(t *testing.T) {
	if _, err := exec.LookPath("hugo"); err != nil {
		t.Skip("hugo not in PATH")
	}
	repo := filepath.Clean(filepath.Join("..", ".."))
	md := filepath.Join(repo, "content", "cognitive-memetics", "cows", "2026-02-26-cow-w01", "index.md")
	if _, err := os.Stat(md); err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	rel := "content/cognitive-memetics/cows/2026-02-26-cow-w01/index.md"
	p, err := LookupPermalink(repo, rel)
	if err != nil {
		t.Fatalf("lookup: %v", err)
	}
	want := "https://blog.dualsubstrate.com/cognitive-memetics/cows/2026-02-26-cow-w01/"
	if p != want {
		t.Fatalf("permalink: got %q want %q", p, want)
	}
}

func TestResolveBundleImageJoinBaseSpanishSiblingUsesEnglishPermalink(t *testing.T) {
	if _, err := exec.LookPath("hugo"); err != nil {
		t.Skip("hugo not in PATH")
	}
	repo, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	md := filepath.Join(repo, "content", "social-protocols", "2026-03-10-hard-times-cycle", "index.es.md")
	if _, err := os.Stat(md); err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	got, err := ResolveBundleImageJoinBase(md)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	want := "https://blog.dualsubstrate.com/social-protocols/2026-03-10-hard-times-cycle/"
	if got != want {
		t.Fatalf("join base: got %q want %q", got, want)
	}
}

func TestResolveBundleImageJoinBaseRelativePathUsesEnglishPermalink(t *testing.T) {
	if _, err := exec.LookPath("hugo"); err != nil {
		t.Skip("hugo not in PATH")
	}
	repo, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}
	rel := "content/social-protocols/2026-03-10-hard-times-cycle/index.es.md"
	if _, err := os.Stat(rel); err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	got, err := ResolveBundleImageJoinBase(rel)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	want := "https://blog.dualsubstrate.com/social-protocols/2026-03-10-hard-times-cycle/"
	if got != want {
		t.Fatalf("join base: got %q want %q", got, want)
	}
}

func TestResolvePagePermalinkForMarkdownCowPost(t *testing.T) {
	if _, err := exec.LookPath("hugo"); err != nil {
		t.Skip("hugo not in PATH")
	}
	repo := filepath.Clean(filepath.Join("..", ".."))
	md := filepath.Join(repo, "content", "cognitive-memetics", "cows", "2026-02-26-cow-w01", "index.md")
	if _, err := os.Stat(md); err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	p, err := ResolvePagePermalinkForMarkdown(md)
	if err != nil {
		t.Fatalf("resolve: %v", err)
	}
	if !strings.Contains(p, "2026-02-26-cow-w01") {
		t.Fatalf("unexpected permalink: %q", p)
	}
}
