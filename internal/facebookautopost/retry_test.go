package facebookautopost

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func testCtx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func zeroOperationDelay() func(int) time.Duration {
	return func(int) time.Duration { return 0 }
}

func TestIsTransientGraphError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "graph 500 typed",
			err: &GraphHTTPError{
				Endpoint:   "https://graph.facebook.com/v20.0/page/photos",
				StatusCode: 500,
				Body:       `{"error":{"code":1,"message":"Please reduce the amount of data you're asking for, then retry your request"}}`,
			},
			want: true,
		},
		{
			name: "graph 429 typed",
			err:  &GraphHTTPError{Endpoint: "https://graph.facebook.com/v20.0/page/feed", StatusCode: 429, Body: "rate limit"},
			want: true,
		},
		{
			name: "graph 400 auth typed",
			err:  &GraphHTTPError{Endpoint: "https://graph.facebook.com/v20.0/page/photos", StatusCode: 400, Body: `{"error":{"code":190}}`},
			want: false,
		},
		{
			name: "network timeout",
			err:  errors.New("Post \"https://graph.facebook.com/v20.0/page/photos\": dial tcp: i/o timeout"),
			want: true,
		},
		{
			name: "client timeout",
			err:  errors.New(`Post "https://graph.facebook.com/v20.0/page/photos": context deadline exceeded (Client.Timeout exceeded while awaiting headers)`),
			want: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsTransientGraphError(tc.err); got != tc.want {
				t.Fatalf("IsTransientGraphError() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDoWithRetryRecoversFromTransientError(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":{"code":1,"message":"retry"}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"123"}`))
	}))
	t.Cleanup(srv.Close)

	c := &Client{
		HTTP:           srv.Client(),
		BaseURL:        strings.TrimSuffix(srv.URL, "/"),
		OperationDelay: zeroOperationDelay(),
	}
	ctx := testCtx(t)
	err := c.DoWithRetry(ctx, 3, func() error {
		return c.PostLink(ctx, "page1", "tok", "hello", "https://example.com/p/")
	})
	if err != nil {
		t.Fatalf("DoWithRetry: %v", err)
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("calls = %d, want 2", got)
	}
}

func TestDoWithRetryDoesNotRetryPermanentError(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":{"code":100,"message":"Invalid parameter"}}`))
	}))
	t.Cleanup(srv.Close)

	c := &Client{
		HTTP:           srv.Client(),
		BaseURL:        strings.TrimSuffix(srv.URL, "/"),
		OperationDelay: zeroOperationDelay(),
	}
	ctx := testCtx(t)
	err := c.DoWithRetry(ctx, 3, func() error {
		return c.PostLink(ctx, "page1", "tok", "hello", "https://example.com/p/")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if got := calls.Load(); got != 1 {
		t.Fatalf("calls = %d, want 1", got)
	}
}

func TestPublishWithRetrySkipsSecondPostWhenURLOnFeed(t *testing.T) {
	const siteURL = "https://blog.dualsubstrate.com/example-post/"
	var postCalls atomic.Int32
	var feedCalls atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/feed"):
			postCalls.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":{"code":1,"message":"retry"}}`))
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/feed"):
			feedCalls.Add(1)
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"data":[{"message":"caption with ` + siteURL + ` link"}]}`))
		default:
			t.Fatalf("unexpected %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(srv.Close)

	c := &Client{
		HTTP:           srv.Client(),
		BaseURL:        strings.TrimSuffix(srv.URL, "/"),
		OperationDelay: zeroOperationDelay(),
	}
	ctx := testCtx(t)
	err := c.PublishWithRetry(ctx, PublishRequest{
		PageID:               "page1",
		AccessToken:          "tok",
		PostURL:              siteURL,
		FeedLimit:            DefaultFeedScanLimit,
		MaxAttempts:          3,
		CheckFeedBeforeRetry: true,
		Post: func() error {
			return c.PostLink(ctx, "page1", "tok", "hello", siteURL)
		},
	})
	if err != nil {
		t.Fatalf("PublishWithRetry: %v", err)
	}
	if got := postCalls.Load(); got != 1 {
		t.Fatalf("post calls = %d, want 1", got)
	}
	if got := feedCalls.Load(); got != 1 {
		t.Fatalf("feed calls = %d, want 1", got)
	}
}

func TestRecentlyPostedURLWithRetryRecoversFromTransientFeedError(t *testing.T) {
	const siteURL = "https://blog.dualsubstrate.com/other/"
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":{"code":1,"message":"retry"}}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"message":"` + siteURL + `"}]}`))
	}))
	t.Cleanup(srv.Close)

	c := &Client{
		HTTP:           srv.Client(),
		BaseURL:        strings.TrimSuffix(srv.URL, "/"),
		OperationDelay: zeroOperationDelay(),
	}
	ctx := testCtx(t)
	already, err := c.RecentlyPostedURLWithRetry(ctx, "page1", "tok", siteURL, 10, 3)
	if err != nil {
		t.Fatalf("RecentlyPostedURLWithRetry: %v", err)
	}
	if !already {
		t.Fatal("expected URL to be found on feed")
	}
	if got := calls.Load(); got != 2 {
		t.Fatalf("calls = %d, want 2", got)
	}
}
