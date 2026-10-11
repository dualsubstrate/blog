package linkedinapi

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dualsubstrate/blog/internal/outbound"
)

func TestYouTubeURLs(t *testing.T) {
	id := "pO0WZsN8Oiw"
	if got := YouTubeWatchURL(id); got != "https://www.youtube.com/watch?v=pO0WZsN8Oiw" {
		t.Fatalf("watch: %q", got)
	}
	if got := YouTubeThumbnailURL(id); got != "https://img.youtube.com/vi/pO0WZsN8Oiw/hqdefault.jpg" {
		t.Fatalf("thumb: %q", got)
	}
}

func TestFetchYouTubeThumbnailRequiresContext(t *testing.T) {
	_, err := FetchYouTubeThumbnail(nil, http.DefaultClient, "pO0WZsN8Oiw")
	if !errors.Is(err, outbound.ErrNilContext) {
		t.Fatalf("err = %v, want ErrNilContext", err)
	}
	ctx := context.Background()
	_, err = FetchYouTubeThumbnail(ctx, http.DefaultClient, "pO0WZsN8Oiw")
	if !errors.Is(err, outbound.ErrMissingDeadline) {
		t.Fatalf("err = %v, want ErrMissingDeadline", err)
	}
}

func TestExtractYouTubeVideoID(t *testing.T) {
	txt := "▶️ watch →\n- https://www.youtube.com/watch?v=pO0WZsN8Oiw\n"
	if got := ExtractYouTubeVideoID(txt); got != "pO0WZsN8Oiw" {
		t.Fatalf("got %q", got)
	}
}
