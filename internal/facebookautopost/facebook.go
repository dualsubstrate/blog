package facebookautopost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dualsubstrate/blog/internal/outbound"
)

const (
	graphAPIVersion       = "v20.0"
	maxGraphResponseBytes = 1 << 20 // 1 MiB cap on Graph API JSON/error bodies
)

func readGraphResponse(r io.Reader) ([]byte, error) {
	limited := io.LimitReader(r, maxGraphResponseBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, err
	}
	if len(body) > maxGraphResponseBytes {
		return nil, fmt.Errorf("response body exceeds %d bytes", maxGraphResponseBytes)
	}
	return body, nil
}

// Client calls the Facebook Graph API. Per-hop timeouts come from caller context deadlines.
type Client struct {
	HTTP           *http.Client
	BaseURL        string
	OperationDelay func(attempt int) time.Duration
}

// NewClient returns a Graph API client. Timeout is kept for API compatibility; use context deadlines at call sites.
func NewClient(timeout time.Duration) *Client {
	if timeout < time.Second {
		timeout = time.Second
	}
	_ = timeout
	return &Client{
		HTTP:    &http.Client{},
		BaseURL: "https://graph.facebook.com/" + graphAPIVersion,
	}
}

func (c *Client) graphURL(path string) string {
	path = strings.TrimLeft(path, "/")
	return c.BaseURL + "/" + path
}

func (c *Client) postForm(ctx context.Context, endpoint string, form url.Values) error {
	encoded := form.Encode()
	return graphStatus(ctx, c, outbound.ClassWrite, endpoint, func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(encoded))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return req, nil
	})
}

type pagePost struct {
	Message string `json:"message"`
	Story   string `json:"story"`
}

type pageFeedResponse struct {
	Data []pagePost `json:"data"`
}

// RecentlyPostedURL reports whether any of the most recent Page feed items already contains urlStr
// in the message or story text.
func (c *Client) RecentlyPostedURL(ctx context.Context, pageID, accessToken, urlStr string, limit int) (bool, error) {
	if strings.TrimSpace(urlStr) == "" {
		return false, fmt.Errorf("url is empty")
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	endpoint := c.graphURL(pageID + "/feed")
	u, err := url.Parse(endpoint)
	if err != nil {
		return false, err
	}
	q := u.Query()
	q.Set("fields", "message,story")
	q.Set("limit", fmt.Sprintf("%d", limit))
	q.Set("access_token", accessToken)
	u.RawQuery = q.Encode()

	body, err := graphRead(ctx, c, outbound.ClassIdempotent, endpoint, func() (*http.Request, error) {
		return http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	})
	if err != nil {
		return false, err
	}

	var r pageFeedResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return false, fmt.Errorf("facebook %s: json: %w", endpoint, err)
	}
	for _, p := range r.Data {
		if strings.Contains(p.Message, urlStr) || strings.Contains(p.Story, urlStr) {
			return true, nil
		}
	}
	return false, nil
}

// PostPhoto publishes a Page photo post (image URL + caption).
func (c *Client) PostPhoto(ctx context.Context, pageID, accessToken, imageURL, caption string) error {
	endpoint := c.graphURL(pageID + "/photos")
	form := url.Values{}
	form.Set("url", imageURL)
	form.Set("caption", caption)
	form.Set("published", "true")
	form.Set("access_token", accessToken)
	return c.postForm(ctx, endpoint, form)
}

// PostPhotoFromFile uploads a local image as multipart form field `source` and publishes it with caption.
func (c *Client) PostPhotoFromFile(ctx context.Context, pageID, accessToken, localPath, caption string) error {
	localPath = strings.TrimSpace(localPath)
	if localPath == "" {
		return fmt.Errorf("local image path is empty")
	}
	f, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("open image: %w", err)
	}
	defer f.Close()

	st, err := f.Stat()
	if err != nil {
		return fmt.Errorf("stat image: %w", err)
	}
	if st.IsDir() {
		return fmt.Errorf("image path is a directory: %s", localPath)
	}

	imageBytes, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read image: %w", err)
	}

	endpoint := c.graphURL(pageID + "/photos")
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	q := parsed.Query()
	q.Set("access_token", accessToken)
	parsed.RawQuery = q.Encode()

	return graphStatus(ctx, c, outbound.ClassWrite, endpoint, func() (*http.Request, error) {
		var buf bytes.Buffer
		mp := multipart.NewWriter(&buf)
		if err := mp.WriteField("published", "true"); err != nil {
			return nil, err
		}
		if err := mp.WriteField("caption", caption); err != nil {
			return nil, err
		}
		part, err := mp.CreateFormFile("source", filepath.Base(localPath))
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(part, bytes.NewReader(imageBytes)); err != nil {
			return nil, err
		}
		if err := mp.Close(); err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, parsed.String(), &buf)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", mp.FormDataContentType())
		return req, nil
	})
}

// PostLink publishes a feed post with message and link preview.
func (c *Client) PostLink(ctx context.Context, pageID, accessToken, message, link string) error {
	endpoint := c.graphURL(pageID + "/feed")
	form := url.Values{}
	form.Set("message", message)
	form.Set("link", link)
	form.Set("published", "true")
	form.Set("access_token", accessToken)
	return c.postForm(ctx, endpoint, form)
}
