package content

import "testing"

func TestNormalizePostSlug(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"platform/foo", "platform/foo"},
		{"content/platform/foo", "platform/foo"},
		{"content/platform/foo/", "platform/foo"},
		{"content/platform/foo/index.md", "platform/foo"},
		{"content/platform/foo/index.es.md", "platform/foo"},
	}
	for _, tc := range tests {
		if got := NormalizePostSlug(tc.in); got != tc.want {
			t.Errorf("NormalizePostSlug(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
