package facebookautopost

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRecentlyPostedURLErrorOmitsAccessToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"message":"bad request"}}`))
	}))
	t.Cleanup(srv.Close)

	c := &Client{
		HTTP:    srv.Client(),
		BaseURL: strings.TrimSuffix(srv.URL, "/"),
	}
	ctx := testCtx(t)
	_, err := c.RecentlyPostedURL(ctx, "page1", "secret-token-xyz", "https://example.com/p/", 5)
	if err == nil {
		t.Fatal("expected error")
	}
	msg := err.Error()
	if strings.Contains(msg, "secret-token-xyz") {
		t.Fatalf("error leaks access_token: %s", msg)
	}
	if !strings.Contains(msg, "/page1/feed") {
		t.Fatalf("error should name endpoint path: %s", msg)
	}
}

func TestReadGraphResponseRejectsOversizedBody(t *testing.T) {
	oversized := strings.Repeat("x", maxGraphResponseBytes+1)
	_, err := readGraphResponse(strings.NewReader(oversized))
	if err == nil {
		t.Fatal("expected error for oversized body")
	}
	if !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("unexpected error: %v", err)
	}
}
