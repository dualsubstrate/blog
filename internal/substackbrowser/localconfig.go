package substackbrowser

import (
	"encoding/json"
	"strings"
)

// legacyLocalConfigJSONKeys maps historical short JSON keys to the current descriptive keys.
// LoadLocalConfig copies a legacy value into the new key only when the new key is absent.
var legacyLocalConfigJSONKeys = map[string]string{
	"pub":                          "substack_publication_subdomain",
	"url":                          "substack_post_editor_url",
	"chrome_user_data_dir":         "chromium_user_data_directory",
	"publish_home_suffix":          "substack_publish_home_url_suffix",
	"login_title_keywords":         "substack_login_page_title_keywords",
	"create_button_text":           "substack_new_post_create_button_text",
	"article_menu_text":            "substack_new_post_article_menu_item_text",
	"landing_url":                  "browser_initial_navigation_url",
	"command_center_url":           "substack_writer_home_url_override",
	"new_article_url":              "substack_new_article_direct_url",
	"navigation_delay_ms":          "chromedp_step_delay_milliseconds",
	"table_mode":                   "markdown_table_mode",
	"paragraph_mode":               "markdown_paragraph_mode",
	"paragraph_break_br_count":     "markdown_paragraph_line_break_repeat_count",
	"quote_mode":                   "markdown_blockquote_mode",
	"include_frontmatter_lead":     "markdown_include_front_matter_lead_block",
	"subtitle_include_categories":  "substack_subtitle_use_category_type_hierarchy_line",
	"subtitle_categories_max":      "substack_subtitle_max_categories_in_type_line",
	"subtitle_max_chars":           "substack_subtitle_max_length_characters",
	"footer_include_tags":          "html_footer_include_article_tags",
	"footer_include_category_link": "html_footer_include_category_browse_link",
	"footer_category_link_index":   "html_footer_category_link_list_index",
	"footer_include_site_link":     "html_footer_include_read_on_site_link",
	"footer_include_cognitive_memetics_project_about": "html_footer_include_cognitive_memetics_project_about",
	"site_base_url":               "site_base_url_for_generated_links",
	"auto_button":                 "substack_insert_category_browse_button_after_paste",
	"button_text":                 "substack_category_browse_button_label",
	"button_url":                  "substack_category_browse_button_url",
	"demote_headings":             "markdown_demote_heading_levels_one_step",
	"schedule_substack_section":   "substack_publish_after_continue_section_label",
	"schedule_email_substack_app": "substack_publish_schedule_push_deliveries",
	"substack_publish_schedule_enable_email_and_app_delivery": "substack_publish_schedule_push_deliveries",
	"schedule_debug_dom_on_failure":                           "substack_publish_schedule_debug_dom_on_failure",
	"schedule_max_attempts":                                   "substack_publish_schedule_max_attempts",
}

// LocalConfig is project-local Substack automation settings (default file: substack.json at repo root).
type LocalConfig struct {
	Pub string `json:"substack_publication_subdomain"`
	URL string `json:"substack_post_editor_url"`

	ChromeUserDataDir string `json:"chromium_user_data_directory"`

	PublishHomeSuffix     string   `json:"substack_publish_home_url_suffix"`
	LoginTitleKeywords    []string `json:"substack_login_page_title_keywords"`
	CreateButtonText      string   `json:"substack_new_post_create_button_text"`
	ArticleMenuText       string   `json:"substack_new_post_article_menu_item_text"`
	LandingURL            string   `json:"browser_initial_navigation_url"`
	CommandCenterURL      string   `json:"substack_writer_home_url_override"`
	NewArticleURL         string   `json:"substack_new_article_direct_url"`
	NavigationDelayMS     int      `json:"chromedp_step_delay_milliseconds"`
	TableMode             string   `json:"markdown_table_mode"`
	ParagraphMode         string   `json:"markdown_paragraph_mode"`
	DemoteHeadings        *bool    `json:"markdown_demote_heading_levels_one_step"`
	ParagraphBreakBRCount int      `json:"markdown_paragraph_line_break_repeat_count"`
	QuoteMode             string   `json:"markdown_blockquote_mode"`

	IncludeFrontMatterLead bool `json:"markdown_include_front_matter_lead_block"`

	SubtitleIncludeCategories bool `json:"substack_subtitle_use_category_type_hierarchy_line"`
	SubtitleCategoriesMax     int  `json:"substack_subtitle_max_categories_in_type_line"`
	SubtitleMaxChars          int  `json:"substack_subtitle_max_length_characters"`

	FooterIncludeTags         bool `json:"html_footer_include_article_tags"`
	FooterIncludeCategoryLink bool `json:"html_footer_include_category_browse_link"`
	FooterCategoryLinkIndex   int  `json:"html_footer_category_link_list_index"`
	FooterIncludeSiteLink     bool `json:"html_footer_include_read_on_site_link"`
	// IncludeCognitiveMemeticsProjectAbout: nil means true (append "But why" blocks for matching cognitive-memetics posts).
	IncludeCognitiveMemeticsProjectAbout *bool  `json:"html_footer_include_cognitive_memetics_project_about"`
	SiteBaseURL                          string `json:"site_base_url_for_generated_links"`
	// MarkdownLeadImageResolveOrigin overrides the scheme and host (only) of the Hugo
	// permalink when turning bundle-relative featured images into absolute URLs for Substack paste.
	// Example: "http://localhost:1313" for local preview only. Leave empty for production
	// image URLs (default). Not inferred from site_base_url_for_generated_links.
	MarkdownLeadImageResolveOrigin string `json:"markdown_lead_image_resolve_origin"`

	AutoButton bool   `json:"substack_insert_category_browse_button_after_paste"`
	ButtonText string `json:"substack_category_browse_button_label"`
	// When empty, cmd/substack-draft sets URL from first Hugo category and site_base_url_for_generated_links.
	ButtonURL string `json:"substack_category_browse_button_url"`
	// InsertSubscribeButtonAfterPaste: when nil, cmd/substack-draft defaults to true (Substack editor **Button** then **Subscribe**).
	InsertSubscribeButtonAfterPaste *bool `json:"substack_insert_subscribe_button_after_paste"`

	// ScheduleSubstackSection is loaded from JSON for compatibility but cmd/substack-draft ignores it:
	// paste-schedule derives the section from categories and content path (see cmd/substack-draft scheduleSectionLabel).
	ScheduleSubstackSection string `json:"substack_publish_after_continue_section_label"`
	SchedulePushDeliveries  *bool  `json:"substack_publish_schedule_push_deliveries"`
	// ScheduleFillDatetimeFromPost is deprecated: paste-schedule always sets Substack schedule time from the
	// Hugo post `date`. The JSON key still unmarshals for backward compatibility but is ignored.
	ScheduleFillDatetimeFromPost *bool `json:"substack_publish_schedule_fill_datetime_from_post"`
	// ScheduleDebugDOMOnFailure, when true, writes a JSON DOM snapshot under tmp/ (or ScheduleDebugDOMFile) if publish-settings automation fails.
	ScheduleDebugDOMOnFailure bool `json:"substack_publish_schedule_debug_dom_on_failure"`
	// ScheduleDebugDOMFile is an optional absolute or relative path for the snapshot JSON (default: tmp/substack-schedule-debug-<timestamp>.json in cwd).
	ScheduleDebugDOMFile string `json:"substack_schedule_debug_dom_file"`
	// ScheduleMaxAttempts is the max number of publish-settings automation passes (recovery runs between passes). 0 or unset means 1 in cmd/substack-draft unless -schedule-max-attempts is set.
	ScheduleMaxAttempts int `json:"substack_publish_schedule_max_attempts"`
}

func mergeLegacyLocalConfigJSON(raw map[string]json.RawMessage) {
	for oldKey, newKey := range legacyLocalConfigJSONKeys {
		if _, hasNew := raw[newKey]; hasNew {
			continue
		}
		if v, ok := raw[oldKey]; ok {
			raw[newKey] = v
		}
	}
}

func trimLocalConfigFields(cfg *LocalConfig) {
	cfg.Pub = strings.TrimSpace(cfg.Pub)
	cfg.URL = strings.TrimSpace(cfg.URL)
	cfg.ChromeUserDataDir = strings.TrimSpace(cfg.ChromeUserDataDir)
	cfg.PublishHomeSuffix = strings.TrimSpace(cfg.PublishHomeSuffix)
	cfg.CreateButtonText = strings.TrimSpace(cfg.CreateButtonText)
	cfg.ArticleMenuText = strings.TrimSpace(cfg.ArticleMenuText)
	cfg.LandingURL = strings.TrimSpace(cfg.LandingURL)
	cfg.CommandCenterURL = strings.TrimSpace(cfg.CommandCenterURL)
	cfg.NewArticleURL = strings.TrimSpace(cfg.NewArticleURL)
	cfg.TableMode = strings.TrimSpace(cfg.TableMode)
	cfg.ParagraphMode = strings.TrimSpace(cfg.ParagraphMode)
	cfg.QuoteMode = strings.TrimSpace(cfg.QuoteMode)
	cfg.SiteBaseURL = strings.TrimSpace(cfg.SiteBaseURL)
	cfg.MarkdownLeadImageResolveOrigin = strings.TrimSpace(cfg.MarkdownLeadImageResolveOrigin)
	cfg.ButtonText = strings.TrimSpace(cfg.ButtonText)
	cfg.ButtonURL = strings.TrimSpace(cfg.ButtonURL)
	cfg.ScheduleSubstackSection = strings.TrimSpace(cfg.ScheduleSubstackSection)
	cfg.ScheduleDebugDOMFile = strings.TrimSpace(cfg.ScheduleDebugDOMFile)
}

// EffectiveIncludeCognitiveMemeticsProjectAbout reports whether to append cognitive-memetics "But why"
// blocks to Substack HTML. A nil JSON pointer means true (on by default).
func EffectiveIncludeCognitiveMemeticsProjectAbout(c LocalConfig) bool {
	if c.IncludeCognitiveMemeticsProjectAbout == nil {
		return true
	}
	return *c.IncludeCognitiveMemeticsProjectAbout
}

func LoadLocalConfig(path string) (LocalConfig, bool, error) {
	return LoadLocalConfigWithGlobal("", path)
}
