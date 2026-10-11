package facebookautopost

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
)

const (
	DefaultPostRetries   = 3
	DefaultFeedScanLimit = 25
)

type withRetryOpts struct {
	maxAttempts int
	logLabel    string
	beforeRetry func(attempt int, lastErr error) (done bool, doneErr error)
}

// PublishRequest configures a retried Page publish (photo or link).
type PublishRequest struct {
	PageID               string
	AccessToken          string
	PostURL              string
	FeedLimit            int
	MaxAttempts          int
	CheckFeedBeforeRetry bool
	Post                 func() error
}

// DoWithRetry calls fn up to maxAttempts times when Graph returns transient errors.
func (c *Client) DoWithRetry(ctx context.Context, maxAttempts int, fn func() error) error {
	return c.withRetry(ctx, withRetryOpts{
		maxAttempts: maxAttempts,
		logLabel:    "transient error",
	}, fn)
}

// RecentlyPostedURLWithRetry is RecentlyPostedURL with transient retries on the feed read.
func (c *Client) RecentlyPostedURLWithRetry(ctx context.Context, pageID, accessToken, urlStr string, limit, maxAttempts int) (bool, error) {
	var already bool
	err := c.withRetry(ctx, withRetryOpts{
		maxAttempts: maxAttempts,
		logLabel:    "transient error",
	}, func() error {
		var err error
		already, err = c.RecentlyPostedURL(ctx, pageID, accessToken, urlStr, limit)
		return err
	})
	return already, err
}

// PublishWithRetry posts up to MaxAttempts times. When CheckFeedBeforeRetry is true and PostURL is set,
// a transient publish error triggers a feed scan before the next attempt so a post that landed despite
// an error response does not get published twice.
func (c *Client) PublishWithRetry(ctx context.Context, req PublishRequest) error {
	maxAttempts := req.MaxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	feedLimit := req.FeedLimit
	if feedLimit <= 0 {
		feedLimit = DefaultFeedScanLimit
	}
	return c.withRetry(ctx, withRetryOpts{
		maxAttempts: maxAttempts,
		logLabel:    "transient publish error",
		beforeRetry: func(attempt int, last error) (bool, error) {
			if !req.CheckFeedBeforeRetry || strings.TrimSpace(req.PostURL) == "" {
				return false, nil
			}
			already, err := c.RecentlyPostedURLWithRetry(ctx, req.PageID, req.AccessToken, req.PostURL, feedLimit, maxAttempts)
			if err == nil && already {
				log.Printf("facebook: URL already on Page after transient publish error; treating as success: %s", req.PostURL)
				return true, nil
			}
			if err != nil {
				log.Printf("facebook: feed check before publish retry failed: %v; will retry publish", err)
			}
			return false, nil
		},
	}, req.Post)
}

func (c *Client) withRetry(ctx context.Context, opts withRetryOpts, fn func() error) error {
	maxAttempts := opts.maxAttempts
	if maxAttempts < 1 {
		maxAttempts = 1
	}

	builder := retrypolicy.NewBuilder[struct{}]().
		HandleIf(func(_ struct{}, err error) bool { return IsTransientGraphError(err) }).
		WithMaxRetries(maxAttempts - 1).
		ReturnLastFailure()

	if c.OperationDelay != nil {
		builder = builder.WithRandomDelay(0, 0)
	} else {
		builder = builder.WithBackoff(5*time.Second, 15*time.Second)
	}

	var skipRemaining bool
	var abortErr error

	retry := builder.
		OnRetry(func(e failsafe.ExecutionEvent[struct{}]) {
			if opts.beforeRetry == nil {
				return
			}
			done, doneErr := opts.beforeRetry(e.Attempts(), e.LastError())
			if done {
				skipRemaining = true
				abortErr = doneErr
			}
		}).
		Build()

	return failsafe.With(retry).WithContext(ctx).Run(func() error {
		if skipRemaining {
			return abortErr
		}
		last := fn()
		if last == nil {
			return nil
		}
		if !IsTransientGraphError(last) {
			return last
		}
		log.Printf("facebook: %s attempt failed: %v", opts.logLabel, last)
		return last
	})
}
