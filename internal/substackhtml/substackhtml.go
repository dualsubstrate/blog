// Package substackhtml converts repository Markdown into HTML intended for
// pasting into Substack's rich-text editor. Markdown remains the source of truth;
// this output is a lossy, compatibility-focused view.
package substackhtml

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

// TableMode controls how GFM tables are emitted after conversion.
type TableMode string

const (
	// TableHTML keeps a minimal semantic table (no class/style). Many rich-text
	// editors accept simple tables on paste; Substack support may vary by client.
	TableHTML TableMode = "html"
	// TableList replaces each table row with a single list item so paste survives
	// even when the editor drops table markup.
	TableList TableMode = "list"
)

// ParagraphMode controls how Markdown paragraphs are emitted after conversion.
type ParagraphMode string

const (
	// ParagraphP keeps semantic paragraphs (<p>).
	ParagraphP ParagraphMode = "p"
	// ParagraphBR flattens paragraphs into inline content separated by <br><br>.
	// This can reduce the "airy" look in some rich-text editors (including Substack).
	ParagraphBR ParagraphMode = "br"
)

// QuoteMode controls how Markdown blockquotes are emitted after conversion.
type QuoteMode string

const (
	// QuoteBlockquote keeps semantic blockquotes (<blockquote>).
	QuoteBlockquote QuoteMode = "blockquote"
	// QuoteMonospace replaces blockquotes with <pre><code>…</code></pre>.
	QuoteMonospace QuoteMode = "monospace"
	// QuotePullquoteMonospace keeps a blockquote container but emits monospace text:
	// <blockquote><code>…</code></blockquote>.
	QuotePullquoteMonospace QuoteMode = "pullquote_monospace"
)

// Options configures Markdown preprocessing and HTML normalization.
type Options struct {
	// TableMode selects table output strategy (see TableMode constants).
	TableMode TableMode
	// ParagraphMode selects how paragraphs are represented (see ParagraphMode constants).
	ParagraphMode ParagraphMode
	// ParagraphBreakBRCount controls how many <br> tags are inserted between paragraphs
	// when ParagraphMode is "br". Default is 2.
	ParagraphBreakBRCount int
	// QuoteMode selects how blockquotes are represented.
	QuoteMode QuoteMode
	// IncludeFrontMatterLead prepends front matter lead content (description, sowhat,
	// youtube_id) before the Markdown body. Useful for Hugo posts where the visible
	// intro is template-driven.
	IncludeFrontMatterLead bool
	// IncludeFeaturedImageLead prepends only the centered featured image (and optional
	// YouTube watch line) from front matter. Used with facebook-en/es body copy.
	IncludeFeaturedImageLead bool
	// DemoteHeadings shifts every heading down one level (h1 to h2, etc.). h6 stays h6.
	// Useful when the post title is entered separately in Substack and body must not
	// compete with the editor title styling.
	DemoteHeadings bool
	// SourcePath is the Markdown file path when known (for example from -in). Used to
	// pick Spanish video lead headings when the basename ends with .es.md, or when
	// front matter lang is es (see buildLeadMarkdown).
	SourcePath string
	// PagePermalink is the published page URL for this Markdown file (trailing slash ok).
	// When set, bundle-relative images from front matter resolve to absolute https URLs
	// so they survive sanitization and display in Substack.
	PagePermalink string
	// ImageResolveOrigin, when set with PagePermalink, replaces only scheme, userinfo, and host
	// (including port) from PagePermalink when resolving bundle-relative featured images, so
	// pasted HTML can load images from a local Hugo server (for example http://localhost:1313)
	// before the post exists on the production site.
	ImageResolveOrigin string

	// bundleImageBase is the URL prefix used to absolutize bundle-relative image paths in the
	// Markdown body (including diagram webp). Convert sets it from PagePermalink and
	// ImageResolveOrigin; callers should not set it.
	bundleImageBase string
}

// DefaultOptions returns conservative defaults for editor paste.
func DefaultOptions() Options {
	return Options{
		TableMode:              TableHTML,
		ParagraphMode:          ParagraphP,
		ParagraphBreakBRCount:  2,
		QuoteMode:              QuoteBlockquote,
		IncludeFrontMatterLead: false,
		DemoteHeadings:         true,
	}
}

// Convert reads full Markdown file bytes (including optional YAML front matter),
// strips front matter and Hugo-oriented noise, renders Markdown to HTML, then
// normalizes HTML for Substack-oriented paste. The returned string is a document
// fragment (no html/head/body wrapper).
func Convert(source []byte, opt Options) (string, error) {
	if len(bytes.TrimSpace(source)) == 0 {
		return "", errors.New("substackhtml: empty source")
	}
	meta := ExtractFrontMatterMeta(source)
	joinBase, err := bundleImageJoinBase(opt)
	if err != nil {
		return "", fmt.Errorf("substackhtml: image resolve origin: %w", err)
	}
	meta.ImageURL = ResolveImageReference(meta.ImageURL, joinBase)
	opt.bundleImageBase = joinBase
	body, err := StripFrontMatter(source)
	if err != nil {
		return "", err
	}
	body = PreprocessMarkdown(body, strings.TrimSpace(opt.SourcePath))
	if opt.IncludeFeaturedImageLead {
		feat, err := renderFeaturedImageMarkdown(meta, opt)
		if err != nil {
			return "", fmt.Errorf("substackhtml: featured image lead: %w", err)
		}
		if strings.TrimSpace(feat) != "" {
			body = append(append(bytes.TrimSpace([]byte(feat)), []byte("\n\n")...), body...)
		}
	}
	if opt.IncludeFrontMatterLead {
		body, err = assembleMarkdownByType(meta, body, opt)
		if err != nil {
			return "", fmt.Errorf("substackhtml: assemble lead markdown: %w", err)
		}
	}
	mdHTML, err := markdownToHTML(body)
	if err != nil {
		return "", fmt.Errorf("substackhtml: markdown: %w", err)
	}
	if opt.TableMode == "" {
		opt.TableMode = TableHTML
	}
	if opt.ParagraphMode == "" {
		opt.ParagraphMode = ParagraphP
	}
	if opt.ParagraphBreakBRCount <= 0 {
		opt.ParagraphBreakBRCount = 2
	}
	if opt.QuoteMode == "" {
		opt.QuoteMode = QuoteBlockquote
	}
	out, err := normalizeHTML(mdHTML, opt)
	if err != nil {
		return "", fmt.Errorf("substackhtml: normalize: %w", err)
	}
	return out, nil
}

// bundleImageJoinBase returns the page URL used to join bundle-relative image paths
// (featured image in front matter and body Markdown images). When ImageResolveOrigin
// is set with a non-empty PagePermalink, scheme and host are taken from ImageResolveOrigin.
func bundleImageJoinBase(opt Options) (string, error) {
	pagePerm := strings.TrimSpace(opt.PagePermalink)
	if src := strings.TrimSpace(opt.SourcePath); src != "" {
		if base, err := ResolveBundleImageJoinBase(src); err == nil && strings.TrimSpace(base) != "" {
			pagePerm = strings.TrimSpace(base)
		}
	}
	if imgOrig := strings.TrimSpace(opt.ImageResolveOrigin); imgOrig != "" && pagePerm != "" {
		return ReplaceHTTPOrigin(pagePerm, imgOrig)
	}
	return pagePerm, nil
}

func isSpanishLeadLocale(meta FrontMatterMeta, opt Options) bool {
	return IsSpanishSiteLocale(meta, opt.SourcePath)
}

// IsSpanishSiteLocale reports whether site links should use the Spanish locale
// (Hugo language prefix /es/ when defaultContentLanguageInSubdir is false).
// Rules match video lead headings: meta.lang "es", or sourcePath basename ending in .es.md.
func IsSpanishSiteLocale(meta FrontMatterMeta, sourcePath string) bool {
	if strings.EqualFold(strings.TrimSpace(meta.Lang), "es") {
		return true
	}
	p := strings.TrimSpace(sourcePath)
	if p == "" {
		return false
	}
	return strings.HasSuffix(strings.ToLower(filepath.Base(p)), ".es.md")
}
