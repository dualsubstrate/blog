package linkedinapi

import "testing"

func TestExtractAndPickCanonicalURL(t *testing.T) {
	txt := "🧷 Full post →\n- ES: https://blog.dualsubstrate.com/es/cognitive-memetics/x/\n- EN: https://blog.dualsubstrate.com/cognitive-memetics/x/\n"
	urls := ExtractSiteURLs(txt)
	if len(urls) != 2 {
		t.Fatalf("urls=%v", urls)
	}
	if got := PickCanonicalURL(urls); got != "https://blog.dualsubstrate.com/cognitive-memetics/x/" {
		t.Fatalf("got %q", got)
	}
}
