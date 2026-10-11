package substackpublishstate

import (
	"os"
	"path/filepath"
	"strings"
)

// TargetLinkedIn is the usual id for the LinkedIn flow (mark-published / pick by target).
const TargetLinkedIn = "linkedin"

// TargetFacebook is the usual id for the Facebook Page autopost flow (mark-published / pick by target).
const TargetFacebook = "facebook"

// BundleReadyForPublish reports whether a Hugo page bundle is ready to publish to targetKey.
// Draft is always honored: index.md must not set draft to a truthy value. For site-es and
// substack-es, index.es.md must exist with real body text and must not set draft to a truthy value either.
// A future post date in front matter does not block readiness (only draft does, for the files that apply).
// Site-en requires index.md body text; site-es requires index.es.md body text.
// LinkedIn requires linkedin.txt with non-whitespace content (draft still comes from index.md).
// Facebook requires facebook-es.txt; substack-en requires substack.md; substack-es requires substack.es.md (and index.es.md).
//
// Missing or unclosed front matter is treated as not draft.
func BundleReadyForPublish(bundleDir, targetKey string) (bool, error) {
	bundleDir = strings.TrimSpace(bundleDir)
	if bundleDir == "" {
		return false, nil
	}
	targetKey = strings.TrimSpace(targetKey)
	if targetKey == "" {
		targetKey = LegacyDefaultTarget
	}
	indexPath := filepath.Join(bundleDir, "index.md")
	fm, err := extractFrontMatter(indexPath)
	if err != nil {
		return false, err
	}
	if draftFromFrontMatter(fm) {
		return false, nil
	}
	switch targetKey {
	case TargetSiteEN:
		ok, err := markdownBodyHasContent(indexPath)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	case TargetSiteES:
		if ok, err := spanishSitePageReady(bundleDir); err != nil || !ok {
			return ok, err
		}
	case TargetSubstackES:
		if ok, err := spanishSitePageReady(bundleDir); err != nil || !ok {
			return ok, err
		}
		ok, err := fileHasNonSpaceContent(filepath.Join(bundleDir, "substack.es.md"))
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	case TargetLinkedIn:
		ok, err := fileHasNonSpaceContent(filepath.Join(bundleDir, "linkedin.txt"))
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	case TargetFacebook:
		ok, err := fileHasNonSpaceContent(filepath.Join(bundleDir, "facebook-es.txt"))
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	default:
		// substack-en and custom targets: draft on index.md only (already checked).
		if targetKey == TargetSubstackEN || targetKey == LegacyDefaultTarget {
			ok, err := fileHasNonSpaceContent(filepath.Join(bundleDir, "substack.md"))
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
	}
	return true, nil
}

func spanishSitePageReady(bundleDir string) (bool, error) {
	esPath := filepath.Join(bundleDir, "index.es.md")
	b, err := os.ReadFile(esPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if strings.TrimSpace(string(b)) == "" {
		return false, nil
	}
	if draftFromFrontMatter(extractFrontMatterBytes(b)) {
		return false, nil
	}
	body, err := markdownBodyBytes(b)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(body)) != "", nil
}

func markdownBodyHasContent(markdownPath string) (bool, error) {
	b, err := os.ReadFile(markdownPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	body, err := markdownBodyBytes(b)
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(body)) != "", nil
}

func markdownBodyBytes(raw []byte) ([]byte, error) {
	s := strings.TrimPrefix(string(raw), "\ufeff")
	if !strings.HasPrefix(s, "---") {
		return raw, nil
	}
	rest := s[3:]
	if strings.HasPrefix(rest, "\r\n") {
		rest = rest[2:]
	} else if strings.HasPrefix(rest, "\n") {
		rest = rest[1:]
	} else {
		return raw, nil
	}
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return raw, nil
	}
	after := rest[end+4:]
	if strings.HasPrefix(after, "\r\n") {
		after = after[2:]
	} else if strings.HasPrefix(after, "\n") {
		after = after[1:]
	}
	return []byte(after), nil
}

func extractFrontMatter(indexPath string) (string, error) {
	b, err := os.ReadFile(indexPath)
	if err != nil {
		return "", err
	}
	return extractFrontMatterBytes(b), nil
}

func extractFrontMatterBytes(b []byte) string {
	const maxFront = 256 * 1024
	if len(b) > maxFront {
		b = b[:maxFront]
	}
	s := string(b)
	s = strings.TrimPrefix(s, "\ufeff")
	if !strings.HasPrefix(s, "---") {
		return ""
	}
	rest := s[3:]
	if strings.HasPrefix(rest, "\r\n") {
		rest = rest[2:]
	} else if strings.HasPrefix(rest, "\n") {
		rest = rest[1:]
	} else {
		return ""
	}
	end := strings.Index(rest, "\n---")
	if end < 0 {
		return ""
	}
	return rest[:end]
}

func splitYAMLKeyVal(line string) (key, val string, ok bool) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false
	}
	i := strings.IndexByte(line, ':')
	if i < 0 {
		return "", "", false
	}
	key = strings.TrimSpace(line[:i])
	val = strings.TrimSpace(line[i+1:])
	return key, val, true
}

func draftFromFrontMatter(fm string) bool {
	for _, line := range strings.Split(fm, "\n") {
		k, v, ok := splitYAMLKeyVal(line)
		if !ok || k != "draft" {
			continue
		}
		v = strings.Trim(v, `"'`)
		switch strings.ToLower(v) {
		case "true", "yes", "1", "y":
			return true
		default:
			return false
		}
	}
	return false
}

func fileHasNonSpaceContent(path string) (bool, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return strings.TrimSpace(string(b)) != "", nil
}
