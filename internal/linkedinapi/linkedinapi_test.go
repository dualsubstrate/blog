package linkedinapi

import (
	"net/http"
	"testing"
)

func TestContentTypeForPath(t *testing.T) {
	got := contentTypeForPath("x.webp")
	if got == "" {
		t.Fatalf("empty")
	}
}

func TestEncodeRestLiResourceKey(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"urn:li:share:7462984940427677697", "urn%3Ali%3Ashare%3A7462984940427677697"},
		{"urn:li:ugcPost:68447855235931240", "urn%3Ali%3AugcPost%3A68447855235931240"},
	}
	for _, tc := range tests {
		if got := encodeRestLiResourceKey(tc.in); got != tc.want {
			t.Fatalf("encodeRestLiResourceKey(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestGetPostRequestPathEncodesURN(t *testing.T) {
	postURN := "urn:li:share:7462984940427677697"
	u := "https://api.linkedin.com/rest/posts/" + encodeRestLiResourceKey(postURN)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		t.Fatal(err)
	}
	got := req.URL.EscapedPath()
	want := "/rest/posts/urn%3Ali%3Ashare%3A7462984940427677697"
	if got != want {
		t.Fatalf("EscapedPath = %q, want %q (URL.String=%q)", got, want, req.URL.String())
	}
}
