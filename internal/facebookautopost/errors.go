package facebookautopost

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// GraphHTTPError is a non-2xx Graph API HTTP response.
type GraphHTTPError struct {
	Endpoint   string
	StatusCode int
	Body       string
}

func (e *GraphHTTPError) Error() string {
	return fmt.Sprintf("facebook %s: status %d body %s", e.Endpoint, e.StatusCode, strings.TrimSpace(e.Body))
}

// Transient reports whether the Graph response should be retried.
func (e *GraphHTTPError) Transient() bool {
	if e.StatusCode == 429 || e.StatusCode >= 500 {
		return true
	}
	lower := strings.ToLower(e.Body)
	return graphErrorCodeIs(e.Body, 1) && strings.Contains(lower, "retry")
}

// graphErrorCodeIs reports whether body contains a Graph error object with the exact code.
func graphErrorCodeIs(body string, code int) bool {
	var payload struct {
		Error *struct {
			Code int `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(body)), &payload); err != nil {
		return false
	}
	return payload.Error != nil && payload.Error.Code == code
}

func newGraphHTTPError(endpoint string, statusCode int, body string) error {
	return &GraphHTTPError{
		Endpoint:   endpoint,
		StatusCode: statusCode,
		Body:       body,
	}
}

// IsTransientGraphError reports whether err is a retryable Graph or transport failure.
func IsTransientGraphError(err error) bool {
	if err == nil {
		return false
	}
	var graphErr *GraphHTTPError
	if errors.As(err, &graphErr) {
		return graphErr.Transient()
	}
	return isTransientTransportError(err)
}

func isTransientTransportError(err error) bool {
	lower := strings.ToLower(err.Error())
	return strings.Contains(lower, "context deadline exceeded") ||
		strings.Contains(lower, "client.timeout exceeded") ||
		strings.Contains(lower, "i/o timeout") ||
		strings.Contains(lower, "connection reset") ||
		strings.Contains(lower, "connection refused") ||
		strings.Contains(lower, "tls handshake timeout") ||
		strings.Contains(lower, "unexpected eof")
}
