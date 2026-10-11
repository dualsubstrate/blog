package socialbundle

import "strings"

const (
	linkedInModeText    = "Posts API text-only from linkedin.txt"
	linkedInModeMedia   = "Posts API image + caption from linkedin.txt"
	linkedInModeArticle = "Posts API article link card (YouTube thumbnail from img.youtube.com)"
)

// LinkedInPostModeLabel describes how linkedin-autopost will publish this bundle.
func (b *Bundle) LinkedInPostModeLabel() string {
	if b.UseLinkedInArticlePost() {
		return linkedInModeArticle
	}
	if strings.TrimSpace(b.FeaturedImagePath) != "" {
		return linkedInModeMedia
	}
	return linkedInModeText
}

// UseLinkedInArticlePost is true when there is no local featured image but a YouTube id is known.
func (b *Bundle) UseLinkedInArticlePost() bool {
	return strings.TrimSpace(b.FeaturedImagePath) == "" && strings.TrimSpace(b.YouTubeID) != ""
}

// HasLinkedInVisualCard is true for local image posts or article posts (thumbnail on card).
func (b *Bundle) HasLinkedInVisualCard() bool {
	return strings.TrimSpace(b.FeaturedImagePath) != "" || b.UseLinkedInArticlePost()
}

func firstDescriptionLine(desc string) string {
	for _, line := range strings.Split(desc, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimPrefix(line, "- ")
		line = strings.Trim(line, "*")
		if line != "" {
			return line
		}
	}
	return ""
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(strings.TrimSpace(s))
	if len(r) <= max {
		return string(r)
	}
	return string(r[:max])
}
