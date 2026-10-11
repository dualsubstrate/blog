package outbound

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/circuitbreaker"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
)

// RetryClass controls whether HTTP-layer retries are allowed.
type RetryClass int

const (
	// ClassIdempotent retries classified transport errors and HTTP 429/5xx.
	ClassIdempotent RetryClass = 1
	// ClassWrite does not retry HTTP calls; circuit breaker still records failures.
	ClassWrite RetryClass = 2
)

// Config configures an outbound HTTP hop.
type Config struct {
	HTTP  *http.Client
	Name  string
	Class RetryClass
}

const maxErrorBodyBytes = 1 << 20

var (
	breakerMu sync.Mutex
	breakers  = map[string]circuitbreaker.CircuitBreaker[*http.Response]{}
)

// ResetBreakersForTest clears the process-wide breaker map. Call from tests that exercise breakers.
func ResetBreakersForTest() {
	breakerMu.Lock()
	breakers = make(map[string]circuitbreaker.CircuitBreaker[*http.Response])
	breakerMu.Unlock()
}

func breakerFor(name string) circuitbreaker.CircuitBreaker[*http.Response] {
	breakerMu.Lock()
	defer breakerMu.Unlock()
	if b, ok := breakers[name]; ok {
		return b
	}
	b := circuitbreaker.NewBuilder[*http.Response]().
		HandleIf(func(_ *http.Response, err error) bool { return isTransientForRetry(err) }).
		WithFailureThreshold(5).
		WithDelay(30 * time.Second).
		Build()
	breakers[name] = b
	return b
}

func retryPolicyFor(class RetryClass) retrypolicy.RetryPolicy[*http.Response] {
	maxRetries := 0
	if class == ClassIdempotent {
		maxRetries = 2
	}
	return retrypolicy.NewBuilder[*http.Response]().
		HandleIf(func(_ *http.Response, err error) bool {
			return class == ClassIdempotent && isTransientForRetry(err)
		}).
		WithBackoff(100*time.Millisecond, time.Second).
		WithJitterFactor(0.2).
		WithMaxRetries(maxRetries).
		Build()
}

func isTransientForRetry(err error) bool {
	if err == nil {
		return false
	}
	var status *StatusError
	if errors.As(err, &status) {
		return IsRetryableHTTPStatus(status.StatusCode)
	}
	lower := strings.ToLower(err.Error())
	return strings.Contains(lower, "context deadline exceeded") ||
		strings.Contains(lower, "client.timeout exceeded") ||
		strings.Contains(lower, "i/o timeout") ||
		strings.Contains(lower, "connection reset") ||
		strings.Contains(lower, "connection refused") ||
		strings.Contains(lower, "tls handshake timeout") ||
		strings.Contains(lower, "unexpected eof")
}

// Do executes an HTTP request under failsafe retry and circuit breaker policies.
func Do(ctx context.Context, cfg Config, newReq func() (*http.Request, error)) (*http.Response, error) {
	if ctx == nil {
		return nil, ErrNilContext
	}
	if _, ok := ctx.Deadline(); !ok {
		return nil, ErrMissingDeadline
	}
	if cfg.HTTP == nil {
		return nil, ErrNilHTTPClient
	}
	name := strings.TrimSpace(cfg.Name)
	if name == "" {
		name = "outbound"
	}

	breaker := breakerFor(name)
	retry := retryPolicyFor(cfg.Class)

	resp, err := failsafe.With(breaker, retry).
		WithContext(ctx).
		GetWithExecution(func(exec failsafe.Execution[*http.Response]) (*http.Response, error) {
			req, err := newReq()
			if err != nil {
				return nil, err
			}
			if req.Context() != exec.Context() {
				req = req.Clone(exec.Context())
			}
			resp, err := cfg.HTTP.Do(req)
			if err != nil {
				return nil, fmt.Errorf("outbound %s: %w", name, err)
			}
			if resp == nil {
				return nil, fmt.Errorf("outbound %s: nil response", name)
			}
			if cfg.Class == ClassIdempotent && IsRetryableHTTPStatus(resp.StatusCode) {
				body, _ := readLimited(resp.Body, maxErrorBodyBytes)
				resp.Body.Close()
				return nil, &StatusError{StatusCode: resp.StatusCode, Body: body}
			}
			return resp, nil
		})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func readLimited(r io.Reader, limit int64) (string, error) {
	if r == nil {
		return "", nil
	}
	b, err := io.ReadAll(io.LimitReader(r, limit+1))
	if err != nil {
		return "", err
	}
	if len(b) > int(limit) {
		b = b[:limit]
	}
	return string(b), nil
}
