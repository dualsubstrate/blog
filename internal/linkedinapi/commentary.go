package linkedinapi

import (
	"fmt"
	"log"
	"strings"

	"github.com/dualsubstrate/blog/internal/socialautopost"
)

// CommentaryPlan is the payload prepared for POST /rest/posts.
type CommentaryPlan struct {
	Raw            string
	Encoded        string
	LittleText     bool
	Stats          socialautopost.CommentaryStats
	Warnings       []string
	SiteURLs       []string
	EncodedChanged bool
}

// PrepareCommentary builds the API commentary string and preflight warnings.
func PrepareCommentary(raw string, withImage, encodeLittleText bool) CommentaryPlan {
	raw = strings.TrimSpace(raw)
	plan := CommentaryPlan{
		Raw:        raw,
		Encoded:    raw,
		LittleText: encodeLittleText,
		Stats:      socialautopost.MeasureCommentary(raw),
		Warnings:   socialautopost.WarnCommentaryLimits(raw, withImage),
		SiteURLs:   ExtractSiteURLs(raw),
	}
	if encodeLittleText {
		plan.Encoded = EncodeCommentaryForPostsAPI(raw)
		plan.EncodedChanged = plan.Encoded != raw
	}
	return plan
}

// LogCommentaryPlan prints byte/rune counts and warnings.
func LogCommentaryPlan(plan CommentaryPlan) {
	log.Printf("commentary: %s", plan.Stats.String())
	if plan.LittleText {
		encStats := socialautopost.MeasureCommentary(plan.Encoded)
		log.Printf("commentary (little text encoded): %s", encStats.String())
		if plan.EncodedChanged {
			log.Printf("commentary: little text encoding changed the payload (hashtags/templates or escapes)")
		}
	}
	for _, w := range plan.Warnings {
		log.Printf("commentary warn: %s", w)
	}
	if len(plan.SiteURLs) > 0 {
		log.Printf("commentary: after publish, verify will check %d site URL(s), footers, hashtags, and length", len(plan.SiteURLs))
	}
}

// VerifyCommentaryURLs ensures LinkedIn stored every blog.dualsubstrate.com URL from the source file.
func VerifyCommentaryURLs(sourceRaw, stored string) error {
	sourceRaw = strings.TrimSpace(sourceRaw)
	stored = strings.TrimSpace(stored)
	if stored == "" {
		return fmt.Errorf("stored commentary is empty")
	}
	want := ExtractSiteURLs(sourceRaw)
	if len(want) == 0 {
		return nil
	}
	var missing []string
	for _, u := range want {
		if commentaryContainsURL(stored, u) {
			continue
		}
		missing = append(missing, u)
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf(
		"stored commentary is missing %d URL(s) from linkedin.txt (API may have truncated or misparsed little text); missing: %s",
		len(missing),
		strings.Join(missing, ", "),
	)
}

func commentaryContainsURL(stored, u string) bool {
	if strings.Contains(stored, u) {
		return true
	}
	// Hashtag templates or encoding may appear in stored form; URL path should still be present.
	if i := strings.Index(u, "blog.dualsubstrate.com"); i >= 0 {
		tail := u[i:]
		return strings.Contains(stored, tail)
	}
	return false
}
