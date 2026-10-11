package outbound

import "time"

// JobBudget returns a context timeout budget for a job with several per-attempt timeouts and extra slack.
func JobBudget(perAttempt time.Duration, attempts int, extra time.Duration) time.Duration {
	if attempts < 1 {
		attempts = 1
	}
	if perAttempt < 0 {
		perAttempt = 0
	}
	if extra < 0 {
		extra = 0
	}
	return perAttempt*time.Duration(attempts) + extra
}
