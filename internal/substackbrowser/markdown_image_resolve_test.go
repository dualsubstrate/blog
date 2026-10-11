package substackbrowser

import "testing"

func TestEffectiveMarkdownLeadImageResolveOrigin(t *testing.T) {
	if got := EffectiveMarkdownLeadImageResolveOrigin(LocalConfig{
		MarkdownLeadImageResolveOrigin: "http://example.com:9999",
		SiteBaseURL:                    "http://localhost:1313/",
	}); got != "http://example.com:9999" {
		t.Fatalf("explicit wins: %q", got)
	}
	if got := EffectiveMarkdownLeadImageResolveOrigin(LocalConfig{
		SiteBaseURL: "http://localhost:1313/",
	}); got != "" {
		t.Fatalf("site_base_url alone must not set image origin: %q", got)
	}
	if got := EffectiveMarkdownLeadImageResolveOrigin(LocalConfig{
		SiteBaseURL: "https://blog.dualsubstrate.com/",
	}); got != "" {
		t.Fatalf("production site base: want empty, got %q", got)
	}
}
