package substackbrowser

import (
	"net/url"
	"strings"
)

// LooksLikeSubstackGenericNewPostEditorURL is true for Substack writer URLs that open the blank
// composer (for example /publish/post?type=newsletter) rather than /publish/post/<draft-id>.
// Navigating to these during schedule recovery starts a new draft in a new document session, so
// recovery should reload the current tab instead of calling Navigate with this URL.
func LooksLikeSubstackGenericNewPostEditorURL(raw string) bool {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || u.Host == "" {
		return false
	}
	host := strings.ToLower(u.Hostname())
	if !strings.HasSuffix(host, ".substack.com") {
		return false
	}
	segs := strings.FieldsFunc(strings.Trim(u.Path, "/"), func(r rune) bool { return r == '/' })
	if len(segs) < 2 {
		return false
	}
	if segs[0] != "publish" || segs[1] != "post" {
		return false
	}
	// /publish/post or /publish/post?… (no draft id path segment)
	return len(segs) == 2
}
