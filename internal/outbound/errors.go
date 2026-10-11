package outbound

import (
	"errors"
	"fmt"
)

// ErrNilContext is returned when ctx is nil.
var ErrNilContext = errors.New("outbound: context is required")

// ErrMissingDeadline is returned when ctx has no deadline.
var ErrMissingDeadline = errors.New("outbound: missing deadline")

// ErrNilHTTPClient is returned when the HTTP client is nil.
var ErrNilHTTPClient = errors.New("outbound: http client is required")

// StatusError is a non-success HTTP response surfaced for retry classification.
type StatusError struct {
	StatusCode int
	Body       string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("outbound: http status %d", e.StatusCode)
}

// IsRetryableHTTPStatus reports whether an HTTP status should be retried on idempotent hops.
func IsRetryableHTTPStatus(code int) bool {
	return code == 429 || code >= 500
}
