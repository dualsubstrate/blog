package socialbundle

import (
	"path/filepath"
	"testing"
)

func TestLoadBundleVideoArticleFields(t *testing.T) {
	root := findRepoRoot(t)
	b, err := LoadBundle(filepath.Join(root, "content"), "proof/2026-05-21-sample-video")
	if err != nil {
		t.Fatal(err)
	}
	if b.YouTubeID != "pO0WZsN8Oiw" {
		t.Fatalf("YouTubeID=%q", b.YouTubeID)
	}
	if !b.UseLinkedInArticlePost() {
		t.Fatal("expected article post without local featured image")
	}
	if b.FeaturedImagePath != "" {
		t.Fatalf("unexpected local image %q", b.FeaturedImagePath)
	}
	if b.ArticleTitle == "" {
		t.Fatal("empty ArticleTitle")
	}
}
