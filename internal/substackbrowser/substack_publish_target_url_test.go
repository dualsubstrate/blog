package substackbrowser

import "testing"

func TestLooksLikeSubstackGenericNewPostEditorURL(t *testing.T) {
	cases := []struct {
		u    string
		want bool
	}{
		{"https://behaviorengineering.substack.com/publish/post?type=newsletter", true},
		{"https://behaviorengineering.substack.com/publish/post/", true},
		{"https://behaviorengineering.substack.com/publish/post", true},
		{"https://behaviorengineering.substack.com/publish/post/draft-uuid-here", false},
		{"https://example.com/publish/post?type=newsletter", false},
		{"", false},
		{"not-a-url", false},
	}
	for _, tc := range cases {
		if got := LooksLikeSubstackGenericNewPostEditorURL(tc.u); got != tc.want {
			t.Errorf("LooksLikeSubstackGenericNewPostEditorURL(%q) = %v want %v", tc.u, got, tc.want)
		}
	}
}
