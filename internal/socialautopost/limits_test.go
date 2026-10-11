package socialautopost

import (
	"strings"
	"testing"
)

func TestValidateLinkedInTxt_withImage_noHardByteCap(t *testing.T) {
	over := strings.Repeat("a", ObservedLinkedInImagePostByteWarn+100)
	if err := ValidateLinkedInTxt(over, true); err != nil {
		t.Fatalf("image posts should not hard-fail on bytes: %v", err)
	}
}

func TestValidateLinkedInTxt_textOnly_runeLimit(t *testing.T) {
	ok := strings.Repeat("a", MaxLinkedInTxtRunesTextOnly)
	if err := ValidateLinkedInTxt(ok, false); err != nil {
		t.Fatalf("at limit: %v", err)
	}
	over := ok + "x"
	if err := ValidateLinkedInTxt(over, false); err == nil {
		t.Fatal("expected error over rune limit text-only")
	}
}

func TestWarnCommentaryLimits_imageBytes(t *testing.T) {
	s := strings.Repeat("a", ObservedLinkedInImagePostByteWarn+1)
	w := WarnCommentaryLimits(s, true)
	if len(w) == 0 {
		t.Fatal("expected warning over observed byte threshold")
	}
}
