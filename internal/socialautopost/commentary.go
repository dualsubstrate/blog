package socialautopost

import (
	"fmt"
	"unicode/utf8"
)

// ObservedLinkedInImagePostByteWarn is a soft threshold from API image posts that lost link blocks.
// Not documented by LinkedIn; use VerifyCommentaryURLs after publish when possible.
const ObservedLinkedInImagePostByteWarn = 800

// MaxLinkedInTxtRunesTextOnly is a conservative cap for link-only / text-only posts (Facebook + LinkedIn).
const MaxLinkedInTxtRunesTextOnly = 3000

// CommentaryStats summarizes linkedin.txt size for logs.
type CommentaryStats struct {
	UTF8Bytes int
	Runes     int
}

func (s CommentaryStats) String() string {
	return fmt.Sprintf("%d UTF-8 bytes, %d runes", s.UTF8Bytes, s.Runes)
}

// MeasureCommentary returns size stats for s.
func MeasureCommentary(s string) CommentaryStats {
	return CommentaryStats{
		UTF8Bytes: len(s),
		Runes:     utf8.RuneCountInString(s),
	}
}

// WarnCommentaryLimits returns non-fatal warnings (logged before publish).
func WarnCommentaryLimits(message string, withImage bool) []string {
	st := MeasureCommentary(message)
	var out []string
	if withImage && st.UTF8Bytes > ObservedLinkedInImagePostByteWarn {
		out = append(out, fmt.Sprintf(
			"%d UTF-8 bytes exceeds observed %d-byte image-post threshold; links may be dropped unless post-verify passes",
			st.UTF8Bytes, ObservedLinkedInImagePostByteWarn,
		))
	}
	if !withImage && st.Runes > MaxLinkedInTxtRunesTextOnly {
		out = append(out, fmt.Sprintf(
			"%d runes exceeds %d rune limit for text-only / link posts",
			st.Runes, MaxLinkedInTxtRunesTextOnly,
		))
	}
	return out
}

// ValidateLinkedInTxt returns an error only for hard limits (text-only rune cap).
// Image posts are checked after publish via linkedinapi.VerifyCommentaryURLs when enabled.
func ValidateLinkedInTxt(message string, withImage bool) error {
	if withImage {
		return nil
	}
	st := MeasureCommentary(message)
	if st.Runes > MaxLinkedInTxtRunesTextOnly {
		return fmt.Errorf(
			"linkedin.txt is %d characters (max %d for link-only Facebook posts and text-only LinkedIn posts)",
			st.Runes, MaxLinkedInTxtRunesTextOnly,
		)
	}
	return nil
}
