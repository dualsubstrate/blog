package substackhtml

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ResolveImageReference returns ref unchanged if it is empty or an absolute http(s) URL.
// Otherwise, if pagePermalink is set, joins the page URL with the relative image path
// (bundle resource file name or path fragment).
func ResolveImageReference(ref, pagePermalink string) string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return ""
	}
	if isAbsHTTPURL(ref) {
		return ref
	}
	pagePermalink = strings.TrimSpace(pagePermalink)
	if pagePermalink == "" {
		return ref
	}
	base := strings.TrimRight(pagePermalink, "/")
	ref = strings.TrimPrefix(strings.TrimSpace(ref), "/")
	return base + "/" + ref
}

func isAbsHTTPURL(s string) bool {
	ls := strings.ToLower(s)
	return strings.HasPrefix(ls, "http://") || strings.HasPrefix(ls, "https://")
}

// ReplaceHTTPOrigin returns fullURL with scheme, host, port, and userinfo taken from origin.
// Path, query, and fragment from fullURL are kept. origin is typically "http://localhost:1313"
// so bundle image URLs load from a local Hugo server before the production deploy.
func ReplaceHTTPOrigin(fullURL, origin string) (string, error) {
	fullURL = strings.TrimSpace(fullURL)
	origin = strings.TrimSpace(origin)
	if origin == "" {
		return fullURL, nil
	}
	u, err := url.Parse(fullURL)
	if err != nil {
		return "", fmt.Errorf("substackhtml: parse page URL %q: %w", fullURL, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("substackhtml: page URL %q must include scheme and host", fullURL)
	}
	base, err := url.Parse(origin)
	if err != nil {
		return "", fmt.Errorf("substackhtml: parse image resolve origin %q: %w", origin, err)
	}
	if base.Host == "" {
		return "", fmt.Errorf("substackhtml: image resolve origin %q must include host", origin)
	}
	scheme := base.Scheme
	if scheme == "" {
		scheme = "https"
	}
	u2 := *u
	u2.Scheme = scheme
	u2.Host = base.Host
	u2.User = base.User
	return u2.String(), nil
}

// FindHugoSiteRoot walks upward from startPath (file or directory) until a directory
// containing hugo.toml is found.
func FindHugoSiteRoot(startPath string) (string, error) {
	startPath = filepath.Clean(startPath)
	dir := startPath
	if st, err := os.Stat(dir); err == nil && !st.IsDir() {
		dir = filepath.Dir(dir)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "hugo.toml")); err == nil {
			abs, absErr := filepath.Abs(dir)
			if absErr != nil {
				return dir, nil
			}
			return abs, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("substackhtml: no hugo.toml above %s", startPath)
		}
		dir = parent
	}
}

// LookupPermalink runs `hugo list all` in siteDir and returns the permalink for the
// content file path relative to the site root (for example content/sec/post/index.md).
func LookupPermalink(siteDir, contentPathRel string) (string, error) {
	siteDir = filepath.Clean(siteDir)
	contentPathRel = filepath.ToSlash(strings.TrimSpace(contentPathRel))
	if contentPathRel == "" {
		return "", errors.New("substackhtml: empty content path")
	}
	cmd := exec.Command("hugo", "list", "all")
	cmd.Dir = siteDir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("substackhtml: hugo list all: %w", err)
	}
	r := csv.NewReader(bytes.NewReader(out))
	r.LazyQuotes = true
	r.ReuseRecord = true
	header, err := r.Read()
	if err != nil {
		return "", fmt.Errorf("substackhtml: hugo list header: %w", err)
	}
	pathIdx, permIdx := -1, -1
	for i, h := range header {
		switch strings.TrimSpace(h) {
		case "path":
			pathIdx = i
		case "permalink":
			permIdx = i
		}
	}
	if pathIdx < 0 || permIdx < 0 {
		return "", errors.New("substackhtml: hugo list missing path or permalink column")
	}
	for {
		rec, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		if len(rec) <= pathIdx || len(rec) <= permIdx {
			continue
		}
		if strings.TrimSpace(rec[pathIdx]) == contentPathRel {
			return strings.TrimSpace(rec[permIdx]), nil
		}
	}
	return "", fmt.Errorf("substackhtml: hugo list has no row for path %q", contentPathRel)
}

// ResolvePagePermalinkForMarkdown returns the site permalink for a Markdown file using hugo list all.
func ResolvePagePermalinkForMarkdown(mdAbsPath string) (string, error) {
	mdAbsPath = filepath.Clean(mdAbsPath)
	if abs, err := filepath.Abs(mdAbsPath); err == nil {
		mdAbsPath = abs
	}
	root, err := FindHugoSiteRoot(mdAbsPath)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(root, mdAbsPath)
	if err != nil {
		return "", err
	}
	rel = filepath.ToSlash(rel)
	if !strings.HasPrefix(rel, "content/") {
		return "", fmt.Errorf("substackhtml: file not under content/: %s", rel)
	}
	return LookupPermalink(root, rel)
}

// ResolveBundleImageJoinBase returns the page URL prefix for bundle-relative images
// (featured image, body figures). Spanish locale siblings publish raster assets under
// the default-language page path on this site, so index.es.md and substack.es.md use
// index.md in the same bundle when present.
func ResolveBundleImageJoinBase(mdAbsPath string) (string, error) {
	mdAbsPath = filepath.Clean(mdAbsPath)
	if abs, err := filepath.Abs(mdAbsPath); err == nil {
		mdAbsPath = abs
	}
	root, err := FindHugoSiteRoot(mdAbsPath)
	if err != nil {
		return "", err
	}
	bundleDir, err := BundleDirFromMarkdownPath(mdAbsPath)
	if err != nil {
		return "", err
	}
	base := strings.ToLower(filepath.Base(mdAbsPath))
	if base == "index.es.md" || base == "substack.es.md" {
		enIndex := filepath.Join(bundleDir, "index.md")
		if _, err := os.Stat(enIndex); err == nil {
			rel, err := filepath.Rel(root, enIndex)
			if err != nil {
				return "", err
			}
			rel = filepath.ToSlash(rel)
			if strings.HasPrefix(rel, "content/") {
				if perm, err := LookupPermalink(root, rel); err == nil && strings.TrimSpace(perm) != "" {
					return strings.TrimSpace(perm), nil
				}
			}
		}
	}
	rel, err := filepath.Rel(root, mdAbsPath)
	if err != nil {
		return "", err
	}
	rel = filepath.ToSlash(rel)
	if !strings.HasPrefix(rel, "content/") {
		return "", fmt.Errorf("substackhtml: file not under content/: %s", rel)
	}
	if base == "substack.es.md" || base == "substack.md" {
		localeIndex := "index.md"
		if base == "substack.es.md" {
			localeIndex = "index.es.md"
		}
		indexPath := filepath.Join(bundleDir, localeIndex)
		if _, err := os.Stat(indexPath); err == nil {
			rel, err = filepath.Rel(root, indexPath)
			if err != nil {
				return "", err
			}
			rel = filepath.ToSlash(rel)
		}
	}
	return LookupPermalink(root, rel)
}
