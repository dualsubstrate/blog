package facebookautopost

import "testing"

func TestGraphHTTPErrorTransient(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		want       bool
	}{
		{name: "500", statusCode: 500, body: `{"error":{"code":1,"message":"retry"}}`, want: true},
		{name: "429", statusCode: 429, body: "rate limit", want: true},
		{name: "400", statusCode: 400, body: `{"error":{"code":190}}`, want: false},
		{name: "code1_retry", statusCode: 400, body: `{"error":{"code":1,"message":"Please retry your request"}}`, want: true},
		{name: "code100_retry_message", statusCode: 400, body: `{"error":{"code":100,"message":"Please retry your request"}}`, want: false},
		{name: "code10_retry_message", statusCode: 400, body: `{"error":{"code":10,"message":"Please retry your request"}}`, want: false},
		{name: "code1_spaced_json", statusCode: 400, body: `{"error": {"code": 1, "message": "Please retry your request"}}`, want: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := &GraphHTTPError{Endpoint: "https://example/feed", StatusCode: tc.statusCode, Body: tc.body}
			if got := err.Transient(); got != tc.want {
				t.Fatalf("Transient() = %v, want %v", got, tc.want)
			}
			if got := IsTransientGraphError(err); got != tc.want {
				t.Fatalf("IsTransientGraphError() = %v, want %v", got, tc.want)
			}
		})
	}
}
