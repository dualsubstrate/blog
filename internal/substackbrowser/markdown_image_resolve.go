package substackbrowser

import "strings"

// EffectiveMarkdownLeadImageResolveOrigin returns the optional HTTP(S) origin used only
// for resolving bundle-relative featured image URLs in Substack HTML export.
// It is set only from markdown_lead_image_resolve_origin (JSON), grouped
// markdown_export.lead_image_resolve_origin, SUBSTACK_MARKDOWN_LEAD_IMAGE_RESOLVE_ORIGIN,
// or -markdown-lead-image-resolve-origin. Empty means keep the Hugo permalink host
// (production baseURL from hugo list all), which is what you want for Substack cover
// images once the post is live. Do not infer from site_base_url_for_generated_links:
// that setting is for footer and button links and is often localhost while images
// must stay on a public URL Substack can fetch.
func EffectiveMarkdownLeadImageResolveOrigin(cfg LocalConfig) string {
	return strings.TrimSpace(cfg.MarkdownLeadImageResolveOrigin)
}
