package linkedinapi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/dualsubstrate/blog/internal/outbound"
)

var youtubeIDRe = regexp.MustCompile(`(?:youtube\.com/watch\?v=|youtu\.be/)([A-Za-z0-9_-]{11})`)

// YouTubeWatchURL returns a canonical watch URL for a video id.
func YouTubeWatchURL(videoID string) string {
	videoID = strings.TrimSpace(videoID)
	if videoID == "" {
		return ""
	}
	return "https://www.youtube.com/watch?v=" + videoID
}

// YouTubeThumbnailURL returns the standard hqdefault JPEG URL (no local file required).
func YouTubeThumbnailURL(videoID string) string {
	videoID = strings.TrimSpace(videoID)
	if videoID == "" {
		return ""
	}
	return "https://img.youtube.com/vi/" + videoID + "/hqdefault.jpg"
}

// ExtractYouTubeVideoID returns the first YouTube video id in text, if any.
func ExtractYouTubeVideoID(text string) string {
	m := youtubeIDRe.FindStringSubmatch(text)
	if len(m) < 2 {
		return ""
	}
	return m[1]
}

// FetchYouTubeThumbnail downloads the hqdefault image for videoID.
func FetchYouTubeThumbnail(ctx context.Context, httpClient *http.Client, videoID string) ([]byte, error) {
	u := YouTubeThumbnailURL(videoID)
	if u == "" {
		return nil, fmt.Errorf("youtube: empty video id")
	}
	if ctx == nil {
		return nil, outbound.ErrNilContext
	}
	if _, ok := ctx.Deadline(); !ok {
		return nil, outbound.ErrMissingDeadline
	}
	if httpClient == nil {
		return nil, outbound.ErrNilHTTPClient
	}
	resp, err := outbound.Do(ctx, outbound.Config{
		HTTP:  httpClient,
		Name:  "youtube-img",
		Class: outbound.ClassIdempotent,
	}, func() (*http.Request, error) {
		return http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	})
	if err != nil {
		return nil, fmt.Errorf("youtube thumbnail GET: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("youtube thumbnail GET %s: status %d", u, resp.StatusCode)
	}
	if len(body) == 0 {
		return nil, fmt.Errorf("youtube thumbnail GET %s: empty body", u)
	}
	return body, nil
}
