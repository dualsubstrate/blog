// Package socialbundle loads Hugo page bundles for Facebook and LinkedIn autopost tools
// (linkedin.txt, canonical site URL, optional featured image from index front matter).
package socialbundle

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dualsubstrate/blog/internal/contentbundle"
	"github.com/dualsubstrate/blog/internal/linkedinapi"
	"github.com/dualsubstrate/blog/internal/tagregister"
	"gopkg.in/yaml.v3"
)

// Human-readable labels for how each network post is sourced (dry-run previews and logs).
const (
	FacebookPostMode = "Page photo + caption from linkedin.txt when featured image exists; else link post"
)

// Bundle is one content bundle ready for social posting.
type Bundle struct {
	RelUnderContent    string // path relative to content/, forward slashes
	BundleDir          string // absolute directory
	IndexPath          string
	LinkedInPath       string
	Message            string // trimmed linkedin.txt body
	PostURL            string // canonical blog.dualsubstrate.com URL from Message
	FeaturedImagePath  string // absolute path if file exists, else empty
	AltText            string // index title, for LinkedIn image alt
	Type               string // Hugo type (e.g. video)
	YouTubeID          string // from front matter or linkedin.txt
	YouTubeURL         string // watch URL when YouTubeID set
	ArticleTitle       string // article card title
	ArticleDescription string // article card description
}

// LoadBundlesForPublishDate resolves bundle paths for the date (same rules as facebook-autopost).
// Then loads each bundle. Rel is under content/ without leading "content/".
func LoadBundlesForPublishDate(repoRoot, dateYYYYMMDD, postPath string) ([]*Bundle, error) {
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, err
	}
	contentRoot := filepath.Join(absRoot, "content")
	rels, err := contentbundle.PublishedBundleRelsForDate(repoRoot, dateYYYYMMDD, postPath)
	if err != nil {
		return nil, err
	}
	var out []*Bundle
	for _, rel := range rels {
		b, err := LoadBundle(contentRoot, rel)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

// LoadBundle reads index.md front matter and linkedin.txt under content/<rel>/.
func LoadBundle(contentRoot, relUnderContent string) (*Bundle, error) {
	rel := strings.Trim(relUnderContent, "/\\")
	rel = strings.TrimPrefix(rel, "content/")
	rel = strings.Trim(rel, "/\\")
	rel = filepath.ToSlash(filepath.Clean(rel))

	bundleDir := filepath.Join(contentRoot, filepath.FromSlash(rel))
	indexPath := filepath.Join(bundleDir, "index.md")
	if st, err := os.Stat(indexPath); err != nil || st.IsDir() {
		return nil, fmt.Errorf("post bundle not found: content/%s (missing index.md)", rel)
	}
	liPath := filepath.Join(bundleDir, "linkedin.txt")
	b, err := os.ReadFile(liPath)
	if err != nil {
		return nil, fmt.Errorf("content/%s: missing linkedin.txt", rel)
	}
	msg := strings.TrimSpace(string(b))
	if msg == "" {
		return nil, fmt.Errorf("content/%s/linkedin.txt: social post text is empty", rel)
	}
	urls := linkedinapi.ExtractSiteURLs(msg)
	postURL := linkedinapi.PickCanonicalURL(urls)
	if postURL == "" {
		return nil, fmt.Errorf("content/%s/linkedin.txt: no blog.dualsubstrate.com URL found (needed for idempotency + post link)", rel)
	}

	raw, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}
	fmYAML, err := tagregister.FrontMatterYAML(raw)
	if err != nil {
		return nil, fmt.Errorf("content/%s: front matter: %w", rel, err)
	}
	var doc struct {
		Type          string `yaml:"type"`
		Title         string `yaml:"title"`
		Subtitle      string `yaml:"subtitle"`
		Description   string `yaml:"description"`
		FeaturedImage string `yaml:"featuredImage"`
		YouTubeID     string `yaml:"youtube_id"`
	}
	if err := yaml.Unmarshal(fmYAML, &doc); err != nil {
		return nil, fmt.Errorf("content/%s: yaml: %w", rel, err)
	}

	var imgPath string
	if strings.TrimSpace(doc.FeaturedImage) != "" {
		candidate := filepath.Join(bundleDir, filepath.FromSlash(strings.TrimSpace(doc.FeaturedImage)))
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			imgPath = candidate
		}
	}

	ytID := strings.TrimSpace(doc.YouTubeID)
	if ytID == "" {
		ytID = linkedinapi.ExtractYouTubeVideoID(msg)
	}
	articleDesc := strings.TrimSpace(doc.Subtitle)
	if articleDesc == "" {
		articleDesc = firstDescriptionLine(doc.Description)
	}

	return &Bundle{
		RelUnderContent:    rel,
		BundleDir:          bundleDir,
		IndexPath:          indexPath,
		LinkedInPath:       liPath,
		Message:            msg,
		PostURL:            postURL,
		FeaturedImagePath:  imgPath,
		AltText:            strings.TrimSpace(doc.Title),
		Type:               strings.TrimSpace(doc.Type),
		YouTubeID:          ytID,
		YouTubeURL:         linkedinapi.YouTubeWatchURL(ytID),
		ArticleTitle:       truncateRunes(doc.Title, 200),
		ArticleDescription: truncateRunes(articleDesc, 300),
	}, nil
}
