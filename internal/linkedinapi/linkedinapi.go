package linkedinapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dualsubstrate/blog/internal/outbound"
)

// DefaultLinkedInVersion is the Linkedin-Version header value (YYYYMM).
// Do not default to the current calendar month: LinkedIn activates versions on their own schedule,
// which can return HTTP 426 NONEXISTENT_VERSION for months that are not live yet.
const DefaultLinkedInVersion = "202604"

// Client calls LinkedIn REST APIs with required version/protocol headers.
type Client struct {
	HTTP          *http.Client
	AccessToken   string
	LinkedinVer   string // YYYYMM
	RestliProto   string // usually 2.0.0
	APIBase       string // https://api.linkedin.com
	UserAgent     string
	RequestLogger func(method, url string, status int, bodyPreview string)
}

func NewClient(timeout time.Duration, accessToken, linkedinVer string) *Client {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ver := strings.TrimSpace(linkedinVer)
	if ver == "" {
		ver = DefaultLinkedInVersion
	}
	_ = timeout // per-hop bounds come from caller context deadlines (outbound.Do).
	return &Client{
		HTTP:        &http.Client{},
		AccessToken: strings.TrimSpace(accessToken),
		LinkedinVer: ver,
		RestliProto: "2.0.0",
		APIBase:     "https://api.linkedin.com",
		UserAgent:   "behaviour-engineering-bot",
	}
}

func (c *Client) addHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("Linkedin-Version", c.LinkedinVer)
	req.Header.Set("X-Restli-Protocol-Version", c.RestliProto)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
}

func (c *Client) do(ctx context.Context, class outbound.RetryClass, newReq func() (*http.Request, error)) (*http.Response, []byte, error) {
	var method string
	resp, err := outbound.Do(ctx, outbound.Config{
		HTTP:  c.HTTP,
		Name:  "linkedin-rest",
		Class: class,
	}, func() (*http.Request, error) {
		req, err := newReq()
		if err != nil {
			return nil, err
		}
		method = req.Method
		c.addHeaders(req)
		return req, nil
	})
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("read body: %w", err)
	}
	if c.RequestLogger != nil {
		preview := strings.TrimSpace(string(body))
		if len(preview) > 600 {
			preview = preview[:600] + "..."
		}
		c.RequestLogger(method, resp.Request.URL.String(), resp.StatusCode, preview)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, nil, errWithLinkedInAPIHints(fmt.Errorf("linkedin %s %s: status %d body %s", method, resp.Request.URL.String(), resp.StatusCode, strings.TrimSpace(string(body))))
	}
	return resp, body, nil
}

// IsAccessDenied reports whether err is a LinkedIn 403 ACCESS_DENIED response.
func IsAccessDenied(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "status 403") && strings.Contains(s, "ACCESS_DENIED")
}

// errWithLinkedInAPIHints appends short guidance for common REST failures (403 read scope, 426 version).
func errWithLinkedInAPIHints(err error) error {
	if err == nil {
		return nil
	}
	s := err.Error()
	if strings.Contains(s, "status 403") && strings.Contains(s, "ACCESS_DENIED") {
		return fmt.Errorf("%w\nhint: list-by-author needs OAuth scope r_member_social (person) or r_organization_social (org); see docs/social-autopost/README.md; or use -disable-idempotency if read scope is unavailable", err)
	}
	if strings.Contains(s, "status 426") && strings.Contains(s, "NONEXISTENT_VERSION") {
		return fmt.Errorf("%w\nhint: Linkedin-Version month is not active yet; bump DefaultLinkedInVersion or set LINKEDIN_VERSION / -linkedin-version to an active YYYYMM (see docs/social-autopost/README.md)", err)
	}
	return err
}

// ---------------- Posts: list by author (idempotency) ----------------

type PostElement struct {
	ID            string `json:"id"`
	Lifecycle     string `json:"lifecycleState"`
	PublishedAt   int64  `json:"publishedAt"`
	LastModified  int64  `json:"lastModifiedAt"`
	Commentary    string `json:"commentary"`
	Visibility    string `json:"visibility"`
	Author        string `json:"author"`
	CreatedAt     int64  `json:"createdAt"`
	Distribution  any    `json:"distribution"`
	Content       any    `json:"content"`
	LifecycleInfo any    `json:"lifecycleStateInfo"`
}

type FindPostsResponse struct {
	Elements []PostElement `json:"elements"`
}

// FindRecentPostsByAuthor returns up to count recent posts for an author URN.
// Requires r_member_social or r_organization_social.
func (c *Client) FindRecentPostsByAuthor(ctx context.Context, authorURN string, count int) ([]PostElement, error) {
	if strings.TrimSpace(authorURN) == "" {
		return nil, fmt.Errorf("authorURN is empty")
	}
	if count <= 0 {
		count = 20
	}
	if count > 100 {
		count = 100
	}
	u, err := url.Parse(c.APIBase + "/rest/posts")
	if err != nil {
		return nil, err
	}
	q := u.Query()
	// author value must be URL-encoded, but url.Values will encode on Encode(). We pass the raw URN.
	q.Set("author", authorURN)
	q.Set("q", "author")
	q.Set("count", fmt.Sprintf("%d", count))
	q.Set("sortBy", "LAST_MODIFIED")
	u.RawQuery = q.Encode()

	_, body, err := c.do(ctx, outbound.ClassIdempotent, func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-RestLi-Method", "FINDER")
		return req, nil
	})
	if err != nil {
		return nil, err
	}
	var r FindPostsResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return nil, fmt.Errorf("json: %w", err)
	}
	return r.Elements, nil
}

// ---------------- Images: initializeUpload + upload bytes ----------------

type InitializeUploadResponse struct {
	Value struct {
		UploadURL string `json:"uploadUrl"`
		ImageURN  string `json:"image"`
	} `json:"value"`
}

func (c *Client) InitializeImageUpload(ctx context.Context, ownerURN string) (uploadURL, imageURN string, err error) {
	u := c.APIBase + "/rest/images?action=initializeUpload"
	payload := map[string]any{
		"initializeUploadRequest": map[string]any{
			"owner": ownerURN,
		},
	}
	b, _ := json.Marshal(payload)
	_, body, err := c.do(ctx, outbound.ClassWrite, func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	})
	if err != nil {
		return "", "", err
	}
	var r InitializeUploadResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return "", "", fmt.Errorf("json: %w", err)
	}
	if strings.TrimSpace(r.Value.UploadURL) == "" || strings.TrimSpace(r.Value.ImageURN) == "" {
		return "", "", fmt.Errorf("initializeUpload returned empty uploadUrl or image URN")
	}
	return r.Value.UploadURL, r.Value.ImageURN, nil
}

func contentTypeForPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == "" {
		return "application/octet-stream"
	}
	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}
	switch ext {
	case ".webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}

// UploadToURL uploads raw bytes to the LinkedIn-provided uploadUrl.
// This request does not use api.linkedin.com, so it does not include LinkedIn version headers.
func (c *Client) UploadToURL(ctx context.Context, uploadURL string, data []byte, contentType string) error {
	if strings.TrimSpace(uploadURL) == "" {
		return fmt.Errorf("uploadURL is empty")
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	resp, err := outbound.Do(ctx, outbound.Config{
		HTTP:  c.HTTP,
		Name:  "linkedin-upload",
		Class: outbound.ClassWrite,
	}, func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", contentType)
		return req, nil
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("linkedin upload PUT %s: status %d body %s", uploadURL, resp.StatusCode, strings.TrimSpace(string(body)))
	}
	return nil
}

// UploadImage reads a local image file and uploads it to LinkedIn, returning the resulting image URN.
func (c *Client) UploadImage(ctx context.Context, ownerURN, localPath string) (string, error) {
	b, err := osReadFile(localPath)
	if err != nil {
		return "", err
	}
	uploadURL, imageURN, err := c.InitializeImageUpload(ctx, ownerURN)
	if err != nil {
		return "", err
	}
	ct := contentTypeForPath(localPath)
	if err := c.UploadToURL(ctx, uploadURL, b, ct); err != nil {
		return "", err
	}
	return imageURN, nil
}

// osReadFile is a tiny seam for tests.
var osReadFile = func(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ---------------- Posts: create ----------------

type createPostRequest struct {
	Author                    string `json:"author"`
	Commentary                string `json:"commentary"`
	Visibility                string `json:"visibility"`
	Distribution              any    `json:"distribution"`
	Content                   any    `json:"content,omitempty"`
	LifecycleState            string `json:"lifecycleState"`
	IsReshareDisabledByAuthor bool   `json:"isReshareDisabledByAuthor"`
}

// UploadImageBytes uploads image bytes via initializeUpload + PUT, returning the image URN.
func (c *Client) UploadImageBytes(ctx context.Context, ownerURN string, data []byte, contentType string) (string, error) {
	if len(data) == 0 {
		return "", fmt.Errorf("image data is empty")
	}
	uploadURL, imageURN, err := c.InitializeImageUpload(ctx, ownerURN)
	if err != nil {
		return "", err
	}
	if contentType == "" {
		contentType = "image/jpeg"
	}
	if err := c.UploadToURL(ctx, uploadURL, data, contentType); err != nil {
		return "", err
	}
	return imageURN, nil
}

// CreatePost publishes text-only, single-image media, or an article link card.
// Returns the post URN from x-restli-id header.
func (c *Client) CreatePost(ctx context.Context, authorURN, commentary string, opts PostOptions) (string, error) {
	u := c.APIBase + "/rest/posts"
	reqBody := createPostRequest{
		Author:     authorURN,
		Commentary: commentary,
		Visibility: "PUBLIC",
		Distribution: map[string]any{
			"feedDistribution":               "MAIN_FEED",
			"targetEntities":                 []any{},
			"thirdPartyDistributionChannels": []any{},
		},
		LifecycleState:            "PUBLISHED",
		IsReshareDisabledByAuthor: false,
	}
	switch {
	case opts.hasArticle():
		a := opts.Article
		reqBody.Content = map[string]any{
			"article": map[string]any{
				"source":      strings.TrimSpace(a.Source),
				"thumbnail":   strings.TrimSpace(a.ThumbnailURN),
				"title":       strings.TrimSpace(a.Title),
				"description": strings.TrimSpace(a.Description),
			},
		}
	case opts.hasMedia():
		reqBody.Content = map[string]any{
			"media": map[string]any{
				"id":      strings.TrimSpace(opts.ImageURN),
				"altText": strings.TrimSpace(opts.AltText),
			},
		}
	}
	b, _ := json.Marshal(reqBody)
	resp, _, err := c.do(ctx, outbound.ClassWrite, func() (*http.Request, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(b))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		return req, nil
	})
	if err != nil {
		return "", err
	}
	id := resp.Header.Get("x-restli-id")
	if strings.TrimSpace(id) == "" {
		// Some gateways return the id in X-RestLi-Id casing.
		id = resp.Header.Get("X-RestLi-Id")
	}
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("create post succeeded but missing x-restli-id header")
	}
	return id, nil
}

// GetPost fetches a single post by URN (value from x-restli-id after create).
// Requires read scope (r_member_social or r_organization_social).
func (c *Client) GetPost(ctx context.Context, postURN string) (PostElement, error) {
	postURN = strings.TrimSpace(postURN)
	if postURN == "" {
		return PostElement{}, fmt.Errorf("postURN is empty")
	}
	u := c.APIBase + "/rest/posts/" + encodeRestLiResourceKey(postURN)
	_, body, err := c.do(ctx, outbound.ClassIdempotent, func() (*http.Request, error) {
		return http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	})
	if err != nil {
		return PostElement{}, err
	}
	var el PostElement
	if err := json.Unmarshal(body, &el); err != nil {
		return PostElement{}, fmt.Errorf("json: %w", err)
	}
	return el, nil
}
