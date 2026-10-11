package linkedinapi

import (
	"regexp"
	"strings"
)

// hashtagToken matches #Tag after start-of-string or whitespace (linkedin.txt hashtag lines).
var hashtagToken = regexp.MustCompile(`(^|\s)#([A-Za-z][A-Za-z0-9_-]*)`)

// EncodeCommentaryForPostsAPI prepares commentary for POST /rest/posts.
//
// LinkedIn stores commentary in "little text" format. Bare # tokens are parsed as hashtag
// elements and can break parsing of following lines (for example URL blocks). Hashtags are
// encoded as {hashtag|\#|Tag} per Microsoft Learn. Reserved characters outside those
// templates are backslash-escaped.
//
// See https://learn.microsoft.com/en-us/linkedin/marketing/community-management/shares/little-text-format
func EncodeCommentaryForPostsAPI(raw string) string {
	if raw == "" {
		return raw
	}
	s := encodeHashtagTemplates(raw)
	return escapeLittleTextReserved(s)
}

func encodeHashtagTemplates(s string) string {
	return hashtagToken.ReplaceAllStringFunc(s, func(match string) string {
		subs := hashtagToken.FindStringSubmatch(match)
		if len(subs) < 3 {
			return match
		}
		return subs[1] + "{hashtag|\\#|" + subs[2] + "}"
	})
}

// reservedLittleTextRunes must be escaped when not already escaped (Microsoft little text spec).
const reservedLittleText = "|{}@[]()<>#\\*_~"

func escapeLittleTextReserved(s string) string {
	var b strings.Builder
	b.Grow(len(s) + 16)
	for i := 0; i < len(s); i++ {
		if strings.HasPrefix(s[i:], "{hashtag|") {
			end := strings.Index(s[i:], "}")
			if end < 0 {
				end = len(s) - i - 1
			}
			b.WriteString(s[i : i+end+1])
			i += end
			continue
		}
		r := s[i]
		if r == '\\' && i+1 < len(s) {
			b.WriteByte(r)
			b.WriteByte(s[i+1])
			i++
			continue
		}
		if strings.ContainsRune(reservedLittleText, rune(r)) {
			b.WriteByte('\\')
		}
		b.WriteByte(r)
	}
	return b.String()
}
