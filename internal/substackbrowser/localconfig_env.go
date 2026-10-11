package substackbrowser

import (
	"os"
	"strconv"
	"strings"
)

// applyLocalConfigFromEnv applies SUBSTACK_* overrides after JSON is loaded.
// Precedence: JSON merge, then env (this function), then CLI flags in cmd/substack-draft
// where the command applies flags after loading. See docs/substack-html/README.md for names.

func finalizeLoadedLocalConfig(cfg *LocalConfig) {
	applyLocalConfigFromEnv(cfg)
	trimLocalConfigFields(cfg)
}

func applyLocalConfigFromEnv(cfg *LocalConfig) {
	if v := getenvTrim("SUBSTACK_PUBLICATION_SUBDOMAIN"); v != "" {
		cfg.Pub = v
	}
	if v := getenvTrim("SUBSTACK_POST_EDITOR_URL"); v != "" {
		cfg.URL = v
	}
	if v := getenvTrim("SUBSTACK_CHROMIUM_USER_DATA_DIRECTORY"); v != "" {
		cfg.ChromeUserDataDir = v
	}
	if v := getenvTrim("SUBSTACK_PUBLISH_HOME_URL_SUFFIX"); v != "" {
		cfg.PublishHomeSuffix = v
	}
	if v := getenvTrim("SUBSTACK_BROWSER_INITIAL_NAVIGATION_URL"); v != "" {
		cfg.LandingURL = v
	}
	if v := getenvTrim("SUBSTACK_WRITER_HOME_URL_OVERRIDE"); v != "" {
		cfg.CommandCenterURL = v
	}
	if v := getenvTrim("SUBSTACK_NEW_ARTICLE_DIRECT_URL"); v != "" {
		cfg.NewArticleURL = v
	}
	if v := getenvTrim("SUBSTACK_SITE_BASE_URL_FOR_GENERATED_LINKS"); v != "" {
		cfg.SiteBaseURL = v
	}
	if v := getenvTrim("SUBSTACK_MARKDOWN_LEAD_IMAGE_RESOLVE_ORIGIN"); v != "" {
		cfg.MarkdownLeadImageResolveOrigin = v
	}
	if v := getenvTrim("SUBSTACK_CATEGORY_BROWSE_BUTTON_LABEL"); v != "" {
		cfg.ButtonText = v
	}
	if v := getenvTrim("SUBSTACK_CATEGORY_BROWSE_BUTTON_URL"); v != "" {
		cfg.ButtonURL = v
	}
	if v := getenvTrim("SUBSTACK_NEW_POST_CREATE_BUTTON_TEXT"); v != "" {
		cfg.CreateButtonText = v
	}
	if v := getenvTrim("SUBSTACK_NEW_POST_ARTICLE_MENU_ITEM_TEXT"); v != "" {
		cfg.ArticleMenuText = v
	}
	if v := getenvTrim("SUBSTACK_MARKDOWN_TABLE_MODE"); v != "" {
		cfg.TableMode = v
	}
	if v := getenvTrim("SUBSTACK_MARKDOWN_PARAGRAPH_MODE"); v != "" {
		cfg.ParagraphMode = v
	}
	if v := getenvTrim("SUBSTACK_MARKDOWN_BLOCKQUOTE_MODE"); v != "" {
		cfg.QuoteMode = v
	}
	if n, ok := getenvInt("SUBSTACK_CHROMEDP_STEP_DELAY_MILLISECONDS"); ok {
		cfg.NavigationDelayMS = n
	}
	if n, ok := getenvInt("SUBSTACK_MARKDOWN_PARAGRAPH_LINE_BREAK_REPEAT_COUNT"); ok {
		cfg.ParagraphBreakBRCount = n
	}
	if n, ok := getenvInt("SUBSTACK_HTML_FOOTER_CATEGORY_LINK_LIST_INDEX"); ok {
		cfg.FooterCategoryLinkIndex = n
	}
	if n, ok := getenvInt("SUBSTACK_SUBTITLE_MAX_CATEGORIES_IN_TYPE_LINE"); ok {
		cfg.SubtitleCategoriesMax = n
	}
	if n, ok := getenvInt("SUBSTACK_SUBTITLE_MAX_LENGTH_CHARACTERS"); ok {
		cfg.SubtitleMaxChars = n
	}
	if b, ok := getenvBool("SUBSTACK_MARKDOWN_INCLUDE_FRONT_MATTER_LEAD_BLOCK"); ok {
		cfg.IncludeFrontMatterLead = b
	}
	if b, ok := getenvBool("SUBSTACK_HTML_FOOTER_INCLUDE_ARTICLE_TAGS"); ok {
		cfg.FooterIncludeTags = b
	}
	if b, ok := getenvBool("SUBSTACK_HTML_FOOTER_INCLUDE_CATEGORY_BROWSE_LINK"); ok {
		cfg.FooterIncludeCategoryLink = b
	}
	if b, ok := getenvBool("SUBSTACK_HTML_FOOTER_INCLUDE_READ_ON_SITE_LINK"); ok {
		cfg.FooterIncludeSiteLink = b
	}
	if b, ok := getenvBool("SUBSTACK_HTML_FOOTER_INCLUDE_COGNITIVE_MEMETICS_PROJECT_ABOUT"); ok {
		v := b
		cfg.IncludeCognitiveMemeticsProjectAbout = &v
	}
	if b, ok := getenvBool("SUBSTACK_SUBTITLE_USE_CATEGORY_TYPE_HIERARCHY_LINE"); ok {
		cfg.SubtitleIncludeCategories = b
	}
	if b, ok := getenvBool("SUBSTACK_INSERT_CATEGORY_BROWSE_BUTTON_AFTER_PASTE"); ok {
		cfg.AutoButton = b
	}
	if b, ok := getenvBool("SUBSTACK_INSERT_SUBSCRIBE_BUTTON_AFTER_PASTE"); ok {
		v := b
		cfg.InsertSubscribeButtonAfterPaste = &v
	}
	if b, ok := getenvBool("SUBSTACK_PUBLISH_SCHEDULE_PUSH_DELIVERIES"); ok {
		v := b
		cfg.SchedulePushDeliveries = &v
	}
	if n, ok := getenvInt("SUBSTACK_PUBLISH_SCHEDULE_MAX_ATTEMPTS"); ok && n > 0 {
		cfg.ScheduleMaxAttempts = n
	}
	if b, ok := getenvBool("SUBSTACK_PUBLISH_SCHEDULE_DEBUG_DOM_ON_FAILURE"); ok {
		cfg.ScheduleDebugDOMOnFailure = b
	}
	if v := getenvTrim("SUBSTACK_SCHEDULE_DEBUG_DOM_FILE"); v != "" {
		cfg.ScheduleDebugDOMFile = v
	}
	if b, ok := getenvBool("SUBSTACK_MARKDOWN_DEMOTE_HEADING_LEVELS_ONE_STEP"); ok {
		v := b
		cfg.DemoteHeadings = &v
	}
}

func getenvTrim(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func getenvBool(key string) (value bool, ok bool) {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return false, false
	}
	switch strings.ToLower(s) {
	case "1", "t", "true", "y", "yes":
		return true, true
	case "0", "f", "false", "n", "no":
		return false, true
	default:
		return false, false
	}
}

func getenvInt(key string) (int, bool) {
	s := strings.TrimSpace(os.Getenv(key))
	if s == "" {
		return 0, false
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	return n, true
}
