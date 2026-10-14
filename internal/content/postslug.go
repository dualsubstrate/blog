// Package content holds small helpers for Hugo content paths.
package content

import (
	"strings"
)

// NormalizePostSlug strips content/ prefix, index.md suffixes, and trailing slashes.
func NormalizePostSlug(post string) string {
	post = strings.TrimPrefix(post, "content/")
	post = strings.TrimSuffix(post, "/")
	post = strings.TrimSuffix(post, "/index.md")
	post = strings.TrimSuffix(post, "/index.es.md")
	return post
}
