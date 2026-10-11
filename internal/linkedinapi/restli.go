package linkedinapi

import (
	"net/url"
	"strings"
)

// encodeRestLiResourceKey URL-encodes a Rest.li 2.0 resource key for use in a path segment.
// Url.PathEscape does not encode ":" (valid in URL paths per RFC 3986), but LinkedIn requires
// colons in URNs to be encoded as %3A (see protocol version 2.0 and Posts API docs).
func encodeRestLiResourceKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return key
	}
	enc := url.QueryEscape(key)
	return strings.ReplaceAll(enc, "+", "%20")
}
