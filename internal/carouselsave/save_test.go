package carouselsave

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePathAcceptsBundleCarousel(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	got, err := ResolvePath(root, "content/human-condition/demo/index.md", "carousel.json")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "content", "human-condition", "demo", "carousel.json")
	if got != want {
		t.Fatalf("got %s want %s", got, want)
	}
}

func TestResolvePathRejectsTraversalAndWrongName(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if _, err := ResolvePath(root, "content/../secrets/index.md", "carousel.json"); err == nil {
		t.Fatal("expected traversal error")
	}
	if _, err := ResolvePath(root, "static/carousel.json", "carousel.json"); err == nil {
		t.Fatal("expected content/ prefix error")
	}
	if _, err := ResolvePath(root, "content/human-condition/demo/index.md", "index.md"); err == nil {
		t.Fatal("expected filename error")
	}
}

func TestWriteBodyWritesPrettyJSON(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	bundle := filepath.Join(root, "content", "human-condition", "demo")
	if err := os.MkdirAll(bundle, 0o755); err != nil {
		t.Fatal(err)
	}
	body := "{\n  \"deck\": {},\n  \"slides\": []\n}"
	path, err := WriteBody(root, Request{
		Source:   "content/human-condition/demo/index.md",
		Filename: "carousel.json",
		Body:     body,
	})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(string(raw), "\n") {
		t.Fatal("expected trailing newline")
	}
	if !strings.Contains(string(raw), `"slides"`) {
		t.Fatalf("unexpected file: %s", raw)
	}
}

func TestWriteBodyRejectsNonObject(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	if _, err := WriteBody(root, Request{
		Source:   "content/human-condition/demo/index.md",
		Filename: "carousel.json",
		Body:     `["nope"]`,
	}); err == nil {
		t.Fatal("expected error")
	}
}
