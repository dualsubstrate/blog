package substackbrowser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Grouped JSON layout (optional): top-level objects group related settings.
// Example: docs/substack-html/substack-config.example.json
// Flat root keys (legacy or long names) still work when no section object is present.

var groupedLocalConfigSectionKeys = []string{
	"substack_browser",
	"markdown_export",
	"substack_subtitle",
	"html_footer",
	"site",
	"substack_editor",
	"substack_publish",
}

func rootJSONUsesGroupedSections(m map[string]json.RawMessage) bool {
	for _, k := range groupedLocalConfigSectionKeys {
		raw, ok := m[k]
		if !ok {
			continue
		}
		s := strings.TrimSpace(string(raw))
		if strings.HasPrefix(s, "{") {
			return true
		}
	}
	return false
}

type groupedSubstackBrowser struct {
	PublicationSubdomain        string   `json:"publication_subdomain"`
	PostEditorURL               string   `json:"post_editor_url"`
	ChromiumUserDataDirectory   string   `json:"chromium_user_data_directory"`
	PublishHomeURLSuffix        string   `json:"publish_home_url_suffix"`
	LoginPageTitleKeywords      []string `json:"login_page_title_keywords"`
	NewPostCreateButtonText     string   `json:"new_post_create_button_text"`
	NewPostArticleMenuItemText  string   `json:"new_post_article_menu_item_text"`
	BrowserInitialNavigationURL string   `json:"browser_initial_navigation_url"`
	WriterDashboardURL          string   `json:"writer_home_url_override"`
	NewArticleDirectURL         string   `json:"new_article_direct_url"`
	StepDelayMilliseconds       int      `json:"step_delay_milliseconds"`
}

type groupedMarkdownExport struct {
	TableMode                     string `json:"table_mode"`
	ParagraphMode                 string `json:"paragraph_mode"`
	ParagraphLineBreakRepeatCount int    `json:"paragraph_line_break_repeat_count"`
	BlockquoteMode                string `json:"blockquote_mode"`
	IncludeFrontMatterLeadBlock   bool   `json:"include_front_matter_lead_block"`
	DemoteHeadingLevelsOneStep    *bool  `json:"demote_heading_levels_one_step"`
	LeadImageResolveOrigin        string `json:"lead_image_resolve_origin"`
}

type groupedSubstackSubtitle struct {
	UseCategoryTypeHierarchyLine bool `json:"use_category_type_hierarchy_line"`
	MaxCategoriesInTypeLine      int  `json:"max_categories_in_type_line"`
	MaxLengthCharacters          int  `json:"max_length_characters"`
}

type groupedHTMLFooter struct {
	IncludeArticleTags                   bool  `json:"include_article_tags"`
	IncludeCategoryBrowseLink            bool  `json:"include_category_browse_link"`
	CategoryLinkListIndex                int   `json:"category_link_list_index"`
	IncludeReadOnSiteLink                bool  `json:"include_read_on_site_link"`
	IncludeCognitiveMemeticsProjectAbout *bool `json:"include_cognitive_memetics_project_about"`
}

type groupedSite struct {
	CanonicalBaseURL string `json:"canonical_base_url"`
}

type groupedSubstackEditor struct {
	InsertCategoryBrowseButtonAfterPaste bool   `json:"insert_category_browse_button_after_paste"`
	InsertSubscribeButtonAfterPaste      *bool  `json:"insert_subscribe_button_after_paste"`
	CategoryBrowseButtonLabel            string `json:"category_browse_button_label"`
	CategoryBrowseButtonURL              string `json:"category_browse_button_url"`
}

type groupedSubstackPublish struct {
	// AfterContinueSectionLabel is ignored by cmd/substack-draft (section comes from post categories[0]).
	AfterContinueSectionLabel string `json:"after_continue_section_label"`
	SchedulePushDeliveries    *bool  `json:"schedule_push_deliveries"`
	// FillScheduleDatetimeFromFrontMatter is deprecated: unmarshaled for backward compatibility; ignored.
	// paste-schedule always fills Substack schedule from Hugo front matter `date`.
	FillScheduleDatetimeFromFrontMatter *bool `json:"fill_schedule_datetime_from_front_matter"`
	ScheduleMaxAttempts                 int   `json:"schedule_max_attempts"`
}

// UnmarshalJSON accepts the current key schedule_push_deliveries and the legacy
// schedule_enable_email_and_app_delivery; the new key wins when both are set.
func (p *groupedSubstackPublish) UnmarshalJSON(data []byte) error {
	var aux struct {
		AfterContinueSectionLabel           string `json:"after_continue_section_label"`
		SchedulePushDeliveries              *bool  `json:"schedule_push_deliveries"`
		FillScheduleDatetimeFromFrontMatter *bool  `json:"fill_schedule_datetime_from_front_matter"`
		ScheduleEnableEmailAndAppDelivery   *bool  `json:"schedule_enable_email_and_app_delivery"`
		ScheduleMaxAttempts                 int    `json:"schedule_max_attempts"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.AfterContinueSectionLabel = aux.AfterContinueSectionLabel
	p.ScheduleMaxAttempts = aux.ScheduleMaxAttempts
	p.FillScheduleDatetimeFromFrontMatter = aux.FillScheduleDatetimeFromFrontMatter
	if aux.SchedulePushDeliveries != nil {
		p.SchedulePushDeliveries = aux.SchedulePushDeliveries
	} else {
		p.SchedulePushDeliveries = aux.ScheduleEnableEmailAndAppDelivery
	}
	return nil
}

type groupedLocalConfigFile struct {
	SubstackBrowser  *groupedSubstackBrowser  `json:"substack_browser"`
	MarkdownExport   *groupedMarkdownExport   `json:"markdown_export"`
	SubstackSubtitle *groupedSubstackSubtitle `json:"substack_subtitle"`
	HTMLFooter       *groupedHTMLFooter       `json:"html_footer"`
	Site             *groupedSite             `json:"site"`
	SubstackEditor   *groupedSubstackEditor   `json:"substack_editor"`
	SubstackPublish  *groupedSubstackPublish  `json:"substack_publish"`
}

func decodeGroupedLocalConfig(b []byte) (LocalConfig, error) {
	var g groupedLocalConfigFile
	if err := json.Unmarshal(b, &g); err != nil {
		return LocalConfig{}, fmt.Errorf("grouped config: %w", err)
	}
	return g.toLocalConfig(), nil
}

func (g groupedLocalConfigFile) toLocalConfig() LocalConfig {
	var c LocalConfig
	if b := g.SubstackBrowser; b != nil {
		c.Pub = b.PublicationSubdomain
		c.URL = b.PostEditorURL
		c.ChromeUserDataDir = b.ChromiumUserDataDirectory
		c.PublishHomeSuffix = b.PublishHomeURLSuffix
		c.LoginTitleKeywords = b.LoginPageTitleKeywords
		c.CreateButtonText = b.NewPostCreateButtonText
		c.ArticleMenuText = b.NewPostArticleMenuItemText
		c.LandingURL = b.BrowserInitialNavigationURL
		c.CommandCenterURL = b.WriterDashboardURL
		c.NewArticleURL = b.NewArticleDirectURL
		c.NavigationDelayMS = b.StepDelayMilliseconds
	}
	if m := g.MarkdownExport; m != nil {
		c.TableMode = m.TableMode
		c.ParagraphMode = m.ParagraphMode
		c.ParagraphBreakBRCount = m.ParagraphLineBreakRepeatCount
		c.QuoteMode = m.BlockquoteMode
		c.IncludeFrontMatterLead = m.IncludeFrontMatterLeadBlock
		c.DemoteHeadings = m.DemoteHeadingLevelsOneStep
		c.MarkdownLeadImageResolveOrigin = m.LeadImageResolveOrigin
	}
	if s := g.SubstackSubtitle; s != nil {
		c.SubtitleIncludeCategories = s.UseCategoryTypeHierarchyLine
		c.SubtitleCategoriesMax = s.MaxCategoriesInTypeLine
		c.SubtitleMaxChars = s.MaxLengthCharacters
	}
	if f := g.HTMLFooter; f != nil {
		c.FooterIncludeTags = f.IncludeArticleTags
		c.FooterIncludeCategoryLink = f.IncludeCategoryBrowseLink
		c.FooterCategoryLinkIndex = f.CategoryLinkListIndex
		c.FooterIncludeSiteLink = f.IncludeReadOnSiteLink
		c.IncludeCognitiveMemeticsProjectAbout = f.IncludeCognitiveMemeticsProjectAbout
	}
	if s := g.Site; s != nil {
		c.SiteBaseURL = s.CanonicalBaseURL
	}
	if e := g.SubstackEditor; e != nil {
		c.AutoButton = e.InsertCategoryBrowseButtonAfterPaste
		c.InsertSubscribeButtonAfterPaste = e.InsertSubscribeButtonAfterPaste
		c.ButtonText = e.CategoryBrowseButtonLabel
		c.ButtonURL = e.CategoryBrowseButtonURL
	}
	if p := g.SubstackPublish; p != nil {
		c.ScheduleSubstackSection = p.AfterContinueSectionLabel
		c.SchedulePushDeliveries = p.SchedulePushDeliveries
		c.ScheduleFillDatetimeFromPost = p.FillScheduleDatetimeFromFrontMatter
		if p.ScheduleMaxAttempts > 0 {
			c.ScheduleMaxAttempts = p.ScheduleMaxAttempts
		}
	}
	return c
}
