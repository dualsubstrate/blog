package linkedinapi

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// minSourceRunesForLengthCheck skips the length ratio test on very short linkedin.txt files.
const minSourceRunesForLengthCheck = 350

// minStoredRatioOfSource is the minimum stored/source rune ratio before we suspect truncation.
const minStoredRatioOfSource = 0.45

// minStoredRatioOfEncoded applies when encodedSent is non-empty (little text payload we POSTed).
const minStoredRatioOfEncoded = 0.45

// VerifyCommentary checks stored commentary from GET /rest/posts/{id} after publish.
// SourceRaw is linkedin.txt; encodedSent is the POST body when little text encoding ran (else "").
func VerifyCommentary(sourceRaw, stored, encodedSent string) error {
	sourceRaw = strings.TrimSpace(sourceRaw)
	stored = strings.TrimSpace(stored)
	if stored == "" {
		return fmt.Errorf("stored commentary is empty")
	}
	if err := VerifyCommentaryURLs(sourceRaw, stored); err != nil {
		return err
	}
	if err := verifyCommentaryFooters(sourceRaw, stored); err != nil {
		return err
	}
	if err := verifyCommentaryHashtagTags(sourceRaw, stored); err != nil {
		return err
	}
	if err := verifyCommentaryBilingualBullets(sourceRaw, stored); err != nil {
		return err
	}
	if err := verifyCommentaryLength(sourceRaw, stored, encodedSent); err != nil {
		return err
	}
	if err := verifyCommentaryLittleTextTemplates(stored); err != nil {
		return err
	}
	if err := verifyCommentaryLittleTextLeak(sourceRaw, stored, encodedSent); err != nil {
		return err
	}
	return nil
}

func verifyCommentaryFooters(source, stored string) error {
	type check struct {
		needle string
		label  string
	}
	checks := []check{
		{"🧷", "pin footer (🧷)"},
		{"🔗", "hub footer (🔗)"},
	}
	for _, c := range checks {
		if !strings.Contains(source, c.needle) {
			continue
		}
		if !strings.Contains(stored, c.needle) {
			return fmt.Errorf("stored commentary missing %s from linkedin.txt", c.label)
		}
	}
	if strings.Contains(source, "Full post") && !strings.Contains(stored, "Full post") {
		return fmt.Errorf("stored commentary missing Full post label (footer likely truncated at parentheses)")
	}
	if strings.Contains(source, "Por-Estas-Calles") && !strings.Contains(stored, "Por-Estas-Calles") {
		return fmt.Errorf("stored commentary missing Por-Estas-Calles hub label")
	}
	return nil
}

func verifyCommentaryHashtagTags(source, stored string) error {
	tags := extractHashtagNames(source)
	if len(tags) == 0 {
		return nil
	}
	var missing []string
	for _, tag := range tags {
		if hashtagPresentInStored(tag, stored) {
			continue
		}
		missing = append(missing, tag)
	}
	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf(
		"stored commentary missing %d hashtag tag name(s) from linkedin.txt: %s",
		len(missing),
		strings.Join(missing, ", "),
	)
}

func extractHashtagNames(text string) []string {
	var out []string
	seen := map[string]struct{}{}
	for _, match := range hashtagToken.FindAllStringSubmatch(text, -1) {
		if len(match) < 3 {
			continue
		}
		tag := match[2]
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
	}
	return out
}

func hashtagPresentInStored(tag, stored string) bool {
	if strings.Contains(stored, tag) {
		return true
	}
	// Encoded little text form: {hashtag|\#|Tag}
	if strings.Contains(stored, "{hashtag|\\#|"+tag+"}") {
		return true
	}
	return false
}

func verifyCommentaryBilingualBullets(source, stored string) error {
	storedLower := strings.ToLower(stored)
	if strings.Contains(source, "- EN:") {
		if !strings.Contains(storedLower, "en:") {
			return fmt.Errorf("stored commentary missing EN: site link line")
		}
	}
	if strings.Contains(source, "- ES:") {
		if !strings.Contains(storedLower, "es:") {
			return fmt.Errorf("stored commentary missing ES: site link line")
		}
	}
	return nil
}

func verifyCommentaryLength(source, stored, encodedSent string) error {
	sourceRunes := utf8.RuneCountInString(source)
	storedRunes := utf8.RuneCountInString(stored)
	if sourceRunes < minSourceRunesForLengthCheck {
		return nil
	}
	if storedRunes < int(float64(sourceRunes)*minStoredRatioOfSource) {
		return fmt.Errorf(
			"stored commentary too short (%d runes vs %d in linkedin.txt); likely truncation after hashtags or parentheses",
			storedRunes, sourceRunes,
		)
	}
	encodedSent = strings.TrimSpace(encodedSent)
	if encodedSent == "" {
		return nil
	}
	encodedRunes := utf8.RuneCountInString(encodedSent)
	if encodedRunes < minSourceRunesForLengthCheck {
		return nil
	}
	if storedRunes < int(float64(encodedRunes)*minStoredRatioOfEncoded) {
		return fmt.Errorf(
			"stored commentary too short vs encoded POST payload (%d runes stored vs %d sent); likely API truncation",
			storedRunes, encodedRunes,
		)
	}
	return nil
}

func verifyCommentaryLittleTextTemplates(stored string) error {
	search := 0
	for {
		i := strings.Index(stored[search:], "{hashtag|")
		if i < 0 {
			return nil
		}
		start := search + i
		end := strings.Index(stored[start:], "}")
		if end < 0 {
			return fmt.Errorf("stored commentary has unclosed {hashtag|...} template (parse may have failed)")
		}
		search = start + end + 1
	}
}

// verifyCommentaryLittleTextLeak flags cases where encoded little text was sent but stored
// commentary lost footer material while still containing raw template/escape artifacts.
func verifyCommentaryLittleTextLeak(source, stored, encodedSent string) error {
	encodedSent = strings.TrimSpace(encodedSent)
	if encodedSent == "" || !strings.Contains(encodedSent, "{hashtag|") {
		return nil
	}
	hasTemplates := strings.Contains(stored, "{hashtag|")
	hasEscapes := strings.Contains(stored, `\(`) || strings.Contains(stored, `\)`)
	if !hasTemplates && !hasEscapes {
		return nil
	}
	// If templates/escapes are present but footers or URLs failed, other checks already errored.
	// When templates are present, require that at least one footer path segment from source appears.
	if strings.Contains(source, "categories/") {
		if idx := strings.Index(source, "categories/"); idx >= 0 {
			rest := source[idx:]
			end := strings.IndexAny(rest, " \n\r\t→)")
			if end > 0 {
				rest = rest[:end]
			}
			if rest != "" && !strings.Contains(stored, rest) && !strings.Contains(stored, strings.TrimPrefix(rest, "categories/")) {
				return fmt.Errorf(
					"stored commentary has little text templates/escapes but is missing hub path %q (footer likely dropped)",
					rest,
				)
			}
		}
	}
	return nil
}
