package linkedinapi

import (
	"strings"
	"testing"
)

func TestEncodeCommentaryForPostsAPI_hashtags(t *testing.T) {
	in := "#CubeCows #CorporateJargon\n\nhttps://blog.dualsubstrate.com/x/\n"
	got := EncodeCommentaryForPostsAPI(in)
	if strings.Contains(got, " #CubeCows") || strings.HasPrefix(got, "#CubeCows") {
		t.Fatalf("bare hashtag should be encoded: %q", got)
	}
	if !strings.Contains(got, "{hashtag|\\#|CubeCows}") || !strings.Contains(got, "{hashtag|\\#|CorporateJargon}") {
		t.Fatalf("expected hashtag templates: %q", got)
	}
	if !strings.Contains(got, "https://blog.dualsubstrate.com/x/") {
		t.Fatalf("URL should remain: %q", got)
	}
}

func TestEncodeCommentaryForPostsAPI_escapesParens(t *testing.T) {
	in := "TS5: Sm(art)\n"
	got := EncodeCommentaryForPostsAPI(in)
	if !strings.Contains(got, `Sm\(art\)`) {
		t.Fatalf("expected escaped parens: %q", got)
	}
}

func TestEncodeCommentaryForPostsAPI_hashtagNoHyphen(t *testing.T) {
	in := "#ArepaContigo #StreetWisdom\n"
	got := EncodeCommentaryForPostsAPI(in)
	if !strings.Contains(got, "{hashtag|\\#|ArepaContigo}") {
		t.Fatalf("expected ArepaContigo template without hyphen: %q", got)
	}
}

func TestEncodeCommentaryForPostsAPI_preservesURLBlock(t *testing.T) {
	in := "#CubeCows\n\n🧷 Full post (site) →\n- EN: https://blog.dualsubstrate.com/cognitive-memetics/cows/w13/\n"
	got := EncodeCommentaryForPostsAPI(in)
	if !strings.Contains(got, "https://blog.dualsubstrate.com/cognitive-memetics/cows/w13/") {
		t.Fatalf("URL block should survive encoding: %q", got)
	}
}

func TestVerifyCommentaryURLs(t *testing.T) {
	source := "see https://blog.dualsubstrate.com/cognitive-memetics/cows/w13/\n"
	stored := "see https://blog.dualsubstrate.com/cognitive-memetics/cows/w13/\n"
	if err := VerifyCommentaryURLs(source, stored); err != nil {
		t.Fatal(err)
	}
	storedShort := "see https://blog.dualsubstrate.com/cognitive-memetics/\n"
	if err := VerifyCommentaryURLs(source, storedShort); err == nil {
		t.Fatal("expected missing URL error")
	}
}
