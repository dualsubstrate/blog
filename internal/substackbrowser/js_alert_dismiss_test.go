package substackbrowser

import "testing"

func TestSubstackTagConflictAlert(t *testing.T) {
	for _, c := range []struct {
		msg string
		ok  bool
	}{
		{"Tag already set", true},
		{"tag already exists", true},
		{"Tag already exists", true},
		{"Something about a tag already exists here", true},
		{"Network error", false},
	} {
		if got := substackTagConflictAlert(c.msg); got != c.ok {
			t.Fatalf("substackTagConflictAlert(%q) = %v, want %v", c.msg, got, c.ok)
		}
	}
}
