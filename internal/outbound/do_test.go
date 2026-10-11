package outbound

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/failsafe-go/failsafe-go/circuitbreaker"
)

func testCtx(t *testing.T) context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func TestDoMissingDeadline(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	_, err := Do(context.Background(), Config{
		HTTP:  srv.Client(),
		Name:  "test",
		Class: ClassIdempotent,
	}, func() (*http.Request, error) {
		return http.NewRequest(http.MethodGet, srv.URL, nil)
	})
	if !errors.Is(err, ErrMissingDeadline) {
		t.Fatalf("err = %v, want ErrMissingDeadline", err)
	}
	if hits.Load() != 0 {
		t.Fatalf("hits = %d, want 0", hits.Load())
	}
}

func TestDoNilContext(t *testing.T) {
	_, err := Do(nil, Config{HTTP: http.DefaultClient, Name: "test", Class: ClassIdempotent}, func() (*http.Request, error) {
		return nil, nil
	})
	if !errors.Is(err, ErrNilContext) {
		t.Fatalf("err = %v, want ErrNilContext", err)
	}
}

func TestDoGETRetriesOn500Then200(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if n == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte("fail"))
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	t.Cleanup(srv.Close)

	resp, err := Do(testCtx(t), Config{
		HTTP:  srv.Client(),
		Name:  "test-retry-" + t.Name(),
		Class: ClassIdempotent,
	}, func() (*http.Request, error) {
		return http.NewRequestWithContext(testCtx(t), http.MethodGet, srv.URL, nil)
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if calls.Load() != 2 {
		t.Fatalf("calls = %d, want 2", calls.Load())
	}
}

func TestDoPOST500SingleAttempt(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	resp, err := Do(testCtx(t), Config{
		HTTP:  srv.Client(),
		Name:  "test-write-" + t.Name(),
		Class: ClassWrite,
	}, func() (*http.Request, error) {
		return http.NewRequestWithContext(testCtx(t), http.MethodPost, srv.URL, nil)
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1", calls.Load())
	}
}

func TestDoGET400NoRetry(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	t.Cleanup(srv.Close)

	resp, err := Do(testCtx(t), Config{
		HTTP:  srv.Client(),
		Name:  "test-400-" + t.Name(),
		Class: ClassIdempotent,
	}, func() (*http.Request, error) {
		return http.NewRequestWithContext(testCtx(t), http.MethodGet, srv.URL, nil)
	})
	if err != nil {
		t.Fatalf("Do: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1", calls.Load())
	}
}

func TestBreakerOpensAfterFailures(t *testing.T) {
	ResetBreakersForTest()
	const dep = "test-breaker-shared"
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)

	ctx := testCtx(t)
	cfg := Config{HTTP: srv.Client(), Name: dep, Class: ClassIdempotent}
	for i := 0; i < 5; i++ {
		_, err := Do(ctx, cfg, func() (*http.Request, error) {
			return http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
		})
		if err == nil {
			t.Fatalf("attempt %d: expected error", i+1)
		}
	}
	before := calls.Load()
	_, err := Do(ctx, cfg, func() (*http.Request, error) {
		return http.NewRequestWithContext(ctx, http.MethodGet, srv.URL, nil)
	})
	if !errors.Is(err, circuitbreaker.ErrOpen) {
		t.Fatalf("err = %v, want ErrOpen", err)
	}
	if calls.Load() != before {
		t.Fatalf("expected no new HTTP call when breaker open; calls before=%d after=%d", before, calls.Load())
	}
}

func TestJobBudget(t *testing.T) {
	got := JobBudget(30*time.Second, 3, 20*time.Second)
	want := 110 * time.Second
	if got != want {
		t.Fatalf("JobBudget = %v, want %v", got, want)
	}
}
