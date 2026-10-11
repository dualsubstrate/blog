package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dualsubstrate/blog/internal/substackbrowser"
)

func TestBuildHTMLUsesFacebookENBodyNotClaimLead(t *testing.T) {
	dir := t.TempDir()
	md := filepath.Join(dir, "index.md")
	fm := "---\ntitle: Built\ntype: claims\ndescription: \"Card claim.\"\nfeaturedImage: built.jpg\n---\n\n## Thoughts\n\nSite essay.\n"
	if err := os.WriteFile(md, []byte(fm), 0o644); err != nil {
		t.Fatal(err)
	}
	sb := "## Look: friends hook\n\nPlain paragraph.\n"
	if err := os.WriteFile(filepath.Join(dir, "substack.md"), []byte(sb), 0o644); err != nil {
		t.Fatal(err)
	}
	lc := substackbrowser.LocalConfig{IncludeFrontMatterLead: true}
	html, _, _, err := buildHTMLAndURL(buildOptions{
		Action: "paste", MDPath: md, Fixture: true, Tables: "html",
		LocalCfg: lc, LocalCfgFound: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(html, "Claim") || strings.Contains(html, "Card claim") {
		t.Fatalf("site claim lead should not appear: %s", html)
	}
	if strings.Contains(html, "Site essay") {
		t.Fatalf("index body should not appear: %s", html)
	}
	if !strings.Contains(html, "friends hook") {
		t.Fatalf("substack body missing: %s", html)
	}
}
