package linkedinapi

import "github.com/dualsubstrate/blog/internal/socialautopost"

// Deprecated: use socialautopost limits directly. Kept for callers that already use linkedinapi.
const (
	MaxCommentaryBytesWithImage = socialautopost.ObservedLinkedInImagePostByteWarn
	MaxCommentaryRunesTextOnly  = socialautopost.MaxLinkedInTxtRunesTextOnly
)

// ValidateCommentary delegates to socialautopost.ValidateLinkedInTxt (hard limit: text-only runes only).
func ValidateCommentary(commentary string, withImage bool) error {
	return socialautopost.ValidateLinkedInTxt(commentary, withImage)
}
