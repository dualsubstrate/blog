// Command substack-draft opens Chromium, converts Markdown to HTML, then pastes
// that HTML into a rich-text editor surface (local fixture or Substack draft UI).
// Optional -confirm-dismiss (paste-schedule) can click Substack's purple footer after you explicitly choose Publish in the terminal (automation never clicks the modal's final Send by itself).
// After a successful paste-schedule run, an optional prompt offers to append social-published for the current -publish-target (or pick default); default is skip unless the user confirms.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dualsubstrate/blog/internal/cliout"
	"github.com/dualsubstrate/blog/internal/substackbrowser"
	"github.com/dualsubstrate/blog/internal/substackhtml"
	"github.com/dualsubstrate/blog/internal/substackpublishstate"
)

func main() {
	log.SetFlags(0)

	action := flag.String("action", "paste", `action: "paste" (default), "paste-schedule", "pick-draft-en" / "pick-draft-es" (interactive list of unpublished bundles, then paste-schedule), "login", "cc", "list-unpublished", or "mark-published"`)
	mdPath := flag.String("in", "", "path to Markdown source (or use -html)")
	htmlPath := flag.String("html", "", "path to HTML fragment (skip Markdown conversion)")
	pub := flag.String("pub", "", "Substack publication subdomain (example: mypub for mypub.substack.com)")
	targetURL := flag.String("url", "", "full URL to open (overrides -pub default when set)")
	fixture := flag.Bool("fixture", false, "use local data: URL with a fake ProseMirror editor (no Substack login)")
	headless := flag.Bool("headless", false, "run Chrome headless (default: headed for login and inspection)")
	userDataDir := flag.String("chrome-user-data-dir", "", "Chrome user data dir so Substack login survives between runs (optional)")
	configPath := flag.String("config", substackbrowser.DefaultLocalConfigPath(), "optional local config path (default: substack.json at repo root)")
	configGlobal := flag.String("config-global", "", "optional shared base config JSON merged under -config (overlay wins); same grouped vs flat shape as -config")
	waitLogin := flag.Duration("wait-login", 0, "after navigation, sleep this long so you can finish Substack login or open New post (example: 2m)")
	pasteTimeout := flag.Duration("paste-timeout", 3*time.Minute, "max wait for editor surface before paste")
	keepOpen := flag.Duration("keep-open", 0, "after paste, keep the browser open this long (example: 30s); ignored when -confirm-dismiss is used")
	confirmDismiss := flag.Bool("confirm-dismiss", false, "headed only: paste-schedule prompts after paste (Continue default) then after automation (huh: Publish default vs Keep open); paste-only prompts before close (Continue default) (huh when stdin+stdout are TTYs, else plain on stderr)")
	tables := flag.String("tables", "html", `Markdown table mode for conversion: "html" or "list"`)
	markdownLeadImageResolveOrigin := flag.String("markdown-lead-image-resolve-origin", "", `optional origin for featured image URLs (e.g. http://localhost:1313); overrides substack.json; empty keeps Hugo permalink host; see docs/substack-html/README.md`)
	noDemote := flag.Bool("no-demote-headings", false, "pass through to HTML converter")
	dryRun := flag.Bool("dry-run", false, "print target URL and HTML stats only; do not launch Chrome")
	title := flag.String("title", "", "optional title override (default: Markdown front matter title)")
	subtitle := flag.String("subtitle", "", "optional subtitle override (default: front matter sowhat, else description first line)")
	contentRoot := flag.String("content-root", "content", "for list-unpublished: Hugo content directory (list skips draft posts and missing channel files; substack-es also respects draft in index.es.md)")
	substackMarker := flag.String("substack-marker", substackpublishstate.DefaultMarker, "marker file name inside each bundle (multi-channel publish log)")
	publishTarget := flag.String("publish-target", "", "channel id for list/mark/pick (e.g. substack-en, substack-es, linkedin); default substack-en, except pick-draft-es uses substack-es")
	scheduleMaxAttempts := flag.Int("schedule-max-attempts", 0, "paste-schedule: max publish-settings automation passes (1=no retry; 2=one recovery+retry after failure; 0=use substack.json schedule_max_attempts or 1)")
	flag.Parse()

	localCfg, localCfgFound, localCfgErr := substackbrowser.LoadLocalConfigWithGlobal(*configGlobal, *configPath)
	if localCfgErr != nil {
		log.Fatal(localCfgErr)
	}

	aEarly := strings.ToLower(strings.TrimSpace(*action))
	if aEarly == "pick-draft-en" || aEarly == "pick-draft-es" {
		tgt := resolvedPublishTarget(aEarly, *publishTarget)
		rel, err := pickUnpublishedBundleRel(*contentRoot, *substackMarker, tgt)
		if err != nil {
			exitIfAborted(err)
			if handlePublishQueueError(os.Stderr, err) {
				os.Exit(1)
			}
			log.Fatal(err)
		}
		if aEarly == "pick-draft-es" {
			*mdPath = filepath.Join(*contentRoot, rel, "index.es.md")
			if _, err := os.Stat(*mdPath); err != nil {
				cliout.PrintSubstackPickMissingSpanish(os.Stderr, *mdPath)
				os.Exit(1)
			}
		} else {
			*mdPath = filepath.Join(*contentRoot, rel, "index.md")
		}
		*action = "paste-schedule"
	}
	if aEarly == "list-unpublished" {
		tgt := resolvedPublishTarget(aEarly, *publishTarget)
		list, err := substackpublishstate.ListUnpublished(*contentRoot, *substackMarker, tgt)
		if err != nil {
			log.Fatal(err)
		}
		if len(list) == 0 {
			printPublishQueueEmpty(os.Stdout, *contentRoot, tgt, *substackMarker)
			return
		}
		printPublishQueueList(os.Stdout, list, tgt, *substackMarker)
		return
	}
	if aEarly == "mark-published" {
		p := strings.TrimSpace(*mdPath)
		if p == "" {
			cliout.PrintSubstackMarkPublishedUsage(os.Stderr)
			os.Exit(1)
		}
		if strings.HasSuffix(strings.ToLower(filepath.ToSlash(p)), "/index.md") || strings.EqualFold(filepath.Base(p), "index.md") {
			p = filepath.Dir(p)
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			log.Fatal(err)
		}
		tgt := resolvedPublishTarget(aEarly, *publishTarget)
		if err := substackpublishstate.MarkPublished(abs, *substackMarker, tgt); err != nil {
			log.Fatal(err)
		}
		cliout.PrintSubstackMarkPublished(os.Stdout, abs, *substackMarker, tgt)
		return
	}

	if localCfgFound {
		// Resolve pub first so later defaults (like new_article_url) can use it.
		if strings.TrimSpace(*pub) == "" && strings.TrimSpace(*targetURL) == "" && strings.TrimSpace(localCfg.Pub) != "" {
			*pub = localCfg.Pub
		}

		// Treat localCfg.URL as the default editor target for paste and paste-schedule.
		// For login/cc actions, prefer the command center URL or pub-derived /publish/home.
		a0 := strings.TrimSpace(*action)
		if a0 == "" || strings.EqualFold(a0, "paste") || strings.EqualFold(a0, "paste-schedule") {
			if strings.TrimSpace(*targetURL) == "" && strings.TrimSpace(localCfg.URL) != "" {
				*targetURL = localCfg.URL
			}
			// If we have a pub and a configured new-article URL, go directly there.
			// This avoids extra navigations and reduces the chance of rate limiting.
			if strings.TrimSpace(*targetURL) == "" && strings.TrimSpace(*pub) != "" && strings.TrimSpace(localCfg.NewArticleURL) != "" {
				*targetURL = localCfg.NewArticleURL
			}
		}
		if strings.TrimSpace(*userDataDir) == "" && strings.TrimSpace(localCfg.ChromeUserDataDir) != "" {
			*userDataDir = localCfg.ChromeUserDataDir
		}
		// If caller didn't override -tables, allow local config to pick a better default for Substack paste.
		if strings.TrimSpace(*tables) == "html" && strings.TrimSpace(localCfg.TableMode) != "" {
			*tables = localCfg.TableMode
		}
	}

	if *dryRun {
		html, cfg, openCfg, err := build(buildOptions{
			Action: *action, MDPath: *mdPath, HTMLPath: *htmlPath, Fixture: *fixture,
			Pub: *pub, TargetURL: *targetURL, LocalCfg: localCfg, LocalCfgFound: localCfgFound,
			Tables: *tables, NoDemote: *noDemote, TitleOverride: *title, SubtitleOverride: *subtitle,
			ImageResolveOrigin: *markdownLeadImageResolveOrigin,
		})
		if err != nil {
			log.Fatal(err)
		}
		fields := map[string]string{
			"headless":                fmt.Sprintf("%v", *headless),
			"insert_subscribe_button": fmt.Sprintf("%v", cfg.InsertSubscribeButton),
		}
		if strings.TrimSpace(openCfg.URL) != "" {
			fields["url"] = openCfg.URL
		} else {
			fields["url"] = cfg.TargetURL
		}
		if html != "" {
			fields["html_bytes"] = fmt.Sprintf("%d", len(html))
		}
		if localCfgFound {
			fields["include_project_about"] = fmt.Sprintf("%v", substackbrowser.EffectiveIncludeCognitiveMemeticsProjectAbout(localCfg))
			if strings.TrimSpace(*configGlobal) != "" {
				fields["config"] = *configGlobal + " + " + *configPath
			} else {
				fields["config"] = *configPath
			}
		}
		if cfg.ScheduleEnabled {
			maxA := effectiveScheduleMaxAttempts(*scheduleMaxAttempts, localCfg, localCfgFound)
			fields["schedule"] = fmt.Sprintf("section=%q tags=%v datetime_local=%q date_display=%q push=%v max_attempts=%d",
				cfg.Schedule.SectionLabel, cfg.Schedule.Tags, cfg.Schedule.DateTimeLocal, cfg.Schedule.DateDisplay,
				cfg.Schedule.TickEmailSubstack, maxA)
			if cfg.ScheduleDebugDOM {
				fields["schedule"] += fmt.Sprintf(" debug_dom=%q", cfg.ScheduleDebugDOMPath)
			}
		}
		cliout.PrintSubstackDryRun(os.Stdout, fields)
		return
	}

	html, pasteCfg, openCfg, err := build(buildOptions{
		Action: *action, MDPath: *mdPath, HTMLPath: *htmlPath, Fixture: *fixture,
		Pub: *pub, TargetURL: *targetURL, LocalCfg: localCfg, LocalCfgFound: localCfgFound,
		Tables: *tables, NoDemote: *noDemote, TitleOverride: *title, SubtitleOverride: *subtitle,
		ImageResolveOrigin: *markdownLeadImageResolveOrigin,
	})
	if err != nil {
		log.Fatal(err)
	}
	pasteCfg.ScheduleMaxAttempts = effectiveScheduleMaxAttempts(*scheduleMaxAttempts, localCfg, localCfgFound)
	if strings.TrimSpace(openCfg.URL) != "" {
		openCfg.Headless = *headless
		openCfg.UserDataDir = *userDataDir
		openCfg.Timeout = *pasteTimeout
		openCfg.KeepOpen = *keepOpen
		if localCfgFound && localCfg.NavigationDelayMS > 0 {
			openCfg.NavigationDelayMS = localCfg.NavigationDelayMS
		}
		if localCfgFound && len(localCfg.LoginTitleKeywords) > 0 {
			openCfg.LoginTitleKeywords = localCfg.LoginTitleKeywords
		}
		cliout.PrintSubstackBrowserOpen(os.Stdout, strings.TrimSpace(*action), openCfg.URL)
		if err := substackbrowser.Open(context.Background(), openCfg); err != nil {
			log.Fatal(err)
		}
		return
	}

	_ = html
	pasteCfg.Headless = *headless
	pasteCfg.UserDataDir = *userDataDir
	pasteCfg.WaitLogin = *waitLogin
	pasteCfg.PasteTimeout = *pasteTimeout
	pasteCfg.KeepOpen = *keepOpen
	if *confirmDismiss && !*headless {
		pasteCfg.KeepOpen = 0
		if strings.EqualFold(strings.TrimSpace(*action), "paste-schedule") {
			pasteCfg.ScheduleAfterPasteConfirm = confirmScheduleAfterPaste
			pasteCfg.ConfirmDismiss = confirmAfterPublishSettingsReview
		} else {
			pasteCfg.ConfirmDismiss = confirmDraftReview
		}
	}
	if localCfgFound {
		if strings.TrimSpace(localCfg.PublishHomeSuffix) != "" {
			pasteCfg.PublishHomeSuffix = localCfg.PublishHomeSuffix
		}
		if len(localCfg.LoginTitleKeywords) > 0 {
			pasteCfg.LoginTitleKeywords = localCfg.LoginTitleKeywords
		}
		if strings.TrimSpace(localCfg.CreateButtonText) != "" {
			pasteCfg.CreateButtonText = localCfg.CreateButtonText
		}
		if strings.TrimSpace(localCfg.ArticleMenuText) != "" {
			pasteCfg.ArticleMenuText = localCfg.ArticleMenuText
		}
		if strings.TrimSpace(localCfg.LandingURL) != "" {
			pasteCfg.LandingURL = localCfg.LandingURL
		}
		if strings.TrimSpace(localCfg.NewArticleURL) != "" {
			pasteCfg.NewArticleURL = localCfg.NewArticleURL
		}
		if localCfg.NavigationDelayMS > 0 {
			pasteCfg.NavigationDelayMS = localCfg.NavigationDelayMS
		}
	}

	cliout.PrintSubstackBrowserOpen(os.Stdout, strings.TrimSpace(*action), pasteCfg.TargetURL)
	if err := substackbrowser.Run(context.Background(), pasteCfg); err != nil {
		exitIfAborted(err)
		log.Fatal(err)
	}
	actionNorm := strings.ToLower(strings.TrimSpace(*action))
	marked := false
	var markedBundle, markedTarget string
	if actionNorm == "paste-schedule" && !*fixture && strings.TrimSpace(*mdPath) != "" {
		bundleDir, err := bundleDirFromMarkdownPath(*mdPath)
		if err != nil {
			log.Printf("substack-draft: skip social-published prompt: %v", err)
		} else {
			tgt := resolvedPublishTarget(aEarly, *publishTarget)
			display := bundleDir
			if rel, err := filepath.Rel(wdOrDot(), bundleDir); err == nil && rel != "" && !strings.HasPrefix(rel, "..") {
				display = rel
			}
			doMark, err := confirmMarkPublishedOptional(display, tgt)
			if err != nil {
				exitIfAborted(err)
				log.Fatal(err)
			}
			if doMark {
				if err := substackpublishstate.MarkPublished(bundleDir, *substackMarker, tgt); err != nil {
					log.Fatalf("substack-draft: mark published: %v", err)
				}
				marked = true
				markedBundle = bundleDir
				markedTarget = tgt
			}
		}
	}
	if actionNorm == "paste" || actionNorm == "paste-schedule" {
		cliout.PrintSubstackPasteFinished(os.Stdout, actionNorm, marked, markedBundle, *substackMarker, markedTarget)
	}
}

func exitIfAborted(err error) {
	if errors.Is(err, substackbrowser.ErrAbortedBeforePublish) {
		cliout.PrintSubstackAborted(os.Stderr)
		os.Exit(130)
	}
}

func wdOrDot() string {
	wd, err := os.Getwd()
	if err != nil || strings.TrimSpace(wd) == "" {
		return "."
	}
	return wd
}

// bundleDirFromMarkdownPath returns the absolute bundle directory for a path to index.md, index.es.md,
// or an existing bundle directory.
func bundleDirFromMarkdownPath(mdIn string) (string, error) {
	p := strings.TrimSpace(mdIn)
	if p == "" {
		return "", fmt.Errorf("empty markdown path")
	}
	base := strings.ToLower(filepath.Base(p))
	if base == "index.md" || base == "index.es.md" {
		p = filepath.Dir(p)
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return abs, nil
}

// buildOptions groups inputs for build and buildHTMLAndURL so call sites stay readable and new flags can extend the struct.
type buildOptions struct {
	Action             string
	MDPath             string
	HTMLPath           string
	Fixture            bool
	Pub                string
	TargetURL          string
	LocalCfg           substackbrowser.LocalConfig
	LocalCfgFound      bool
	Tables             string
	NoDemote           bool
	TitleOverride      string
	SubtitleOverride   string
	ImageResolveOrigin string
}

func build(opts buildOptions) (string, substackbrowser.Config, substackbrowser.OpenConfig, error) {
	var openCfg substackbrowser.OpenConfig
	a := strings.TrimSpace(opts.Action)
	if a == "" {
		a = "paste"
	}
	if a != "login" && a != "cc" && a != "paste" && a != "paste-schedule" {
		return "", substackbrowser.Config{}, openCfg, fmt.Errorf("unknown action %q (use paste, paste-schedule, login, cc, list-unpublished, or mark-published)", a)
	}

	// login/cc actions do not require Markdown or HTML.
	if a == "login" || a == "cc" {
		// For login, prefer the public Substack home if configured. This makes it
		// easy to establish a session without immediately hitting a private dashboard URL.
		if a == "login" && opts.LocalCfgFound && strings.TrimSpace(opts.LocalCfg.LandingURL) != "" {
			openCfg.URL = opts.LocalCfg.LandingURL
			openCfg.WaitUntilLoggedIn = true
			return "", substackbrowser.Config{}, openCfg, nil
		}

		u := strings.TrimSpace(opts.TargetURL)
		if u == "" && strings.TrimSpace(opts.Pub) != "" {
			u = substackbrowser.DefaultPublishHomeURL(opts.Pub)
		}
		if u == "" && opts.LocalCfgFound && strings.TrimSpace(opts.LocalCfg.CommandCenterURL) != "" {
			u = opts.LocalCfg.CommandCenterURL
		}
		if u == "" && opts.LocalCfgFound && strings.TrimSpace(opts.LocalCfg.Pub) != "" {
			u = substackbrowser.DefaultPublishHomeURL(opts.LocalCfg.Pub)
		}
		if strings.TrimSpace(u) == "" {
			return "", substackbrowser.Config{}, substackbrowser.OpenConfig{}, fmt.Errorf("need pub/url or substack.json for action %q", a)
		}
		openCfg.URL = u
		if opts.LocalCfgFound && strings.TrimSpace(opts.LocalCfg.LandingURL) != "" {
			openCfg.LandingURL = opts.LocalCfg.LandingURL
		}
		openCfg.WaitUntilLoggedIn = a == "login"
		return "", substackbrowser.Config{}, openCfg, nil
	}

	o := opts
	o.Action = a
	return buildHTMLAndURL(o)
}

func buildHTMLAndURL(opts buildOptions) (string, substackbrowser.Config, substackbrowser.OpenConfig, error) {
	action := strings.TrimSpace(opts.Action)
	mdPath := opts.MDPath
	htmlPath := opts.HTMLPath
	fixture := opts.Fixture
	pub := opts.Pub
	targetURL := opts.TargetURL
	tables := opts.Tables
	localCfg := opts.LocalCfg
	localCfgFound := opts.LocalCfgFound
	noDemote := opts.NoDemote
	titleOverride := opts.TitleOverride
	subtitleOverride := opts.SubtitleOverride
	imageResolveOrigin := opts.ImageResolveOrigin
	var cfg substackbrowser.Config
	var html string
	var meta substackhtml.FrontMatterMeta

	if htmlPath != "" {
		b, err := os.ReadFile(htmlPath)
		if err != nil {
			return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("read html: %w", err)
		}
		html = string(b)
	} else {
		raw, err := os.ReadFile(mdPath)
		if err != nil {
			return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("read markdown: %w", err)
		}
		indexRaw := raw
		bodyRes, fbErr := substackhtml.ResolveSubstackBody(raw, mdPath)
		if fbErr != nil {
			return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("substack body: %w", fbErr)
		}
		usedSubstackSidecar := bodyRes.Source == substackhtml.SubstackBodyFromSidecarMD
		raw = bodyRes.IndexRaw
		if strings.TrimSpace(titleOverride) != "" {
			cfg.Title = strings.TrimSpace(titleOverride)
		}
		if strings.TrimSpace(subtitleOverride) != "" {
			cfg.Subtitle = strings.TrimSpace(subtitleOverride)
		}
		meta = substackhtml.ExtractFrontMatterMeta(indexRaw)
		if strings.TrimSpace(cfg.Title) == "" && strings.TrimSpace(meta.Title) != "" {
			cfg.Title = meta.Title
		}
		if strings.TrimSpace(cfg.Subtitle) == "" {
			// Substack reuses the subtitle as the default social preview line; keep it short by
			// preferring the Type: category hierarchy when configured, not the full description.
			if localCfgFound && localCfg.SubtitleIncludeCategories && len(meta.Categories) > 0 {
				max := localCfg.SubtitleCategoriesMax
				if max <= 0 {
					max = 3
				}
				if max > len(meta.Categories) {
					max = len(meta.Categories)
				}
				if line := subtitleTypeCategoryLine(meta.Categories, max, meta.Type); line != "" {
					cfg.Subtitle = line
				}
			}
		}
		if strings.TrimSpace(cfg.Subtitle) == "" {
			// Subtitle fallback: prefer sowhat, else first non-empty line of description (list-style).
			sub := strings.TrimSpace(meta.SoWhat)
			if sub == "" {
				for _, ln := range strings.Split(meta.Description, "\n") {
					t := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(ln), "-"))
					if t != "" {
						sub = t
						break
					}
				}
			}
			cfg.Subtitle = sub
		}
		if localCfgFound {
			maxChars := localCfg.SubtitleMaxChars
			if maxChars <= 0 {
				maxChars = 120
			}
			cfg.Subtitle = truncateRunes(strings.TrimSpace(cfg.Subtitle), maxChars)
		}
		opt := substackhtml.DefaultOptions()
		if localCfgFound && localCfg.DemoteHeadings != nil {
			opt.DemoteHeadings = *localCfg.DemoteHeadings
		}
		if noDemote {
			opt.DemoteHeadings = false
		}
		if localCfgFound && strings.TrimSpace(localCfg.ParagraphMode) != "" {
			switch strings.TrimSpace(localCfg.ParagraphMode) {
			case "p":
				opt.ParagraphMode = substackhtml.ParagraphP
			case "br":
				opt.ParagraphMode = substackhtml.ParagraphBR
			default:
				return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("unknown paragraph_mode %q (use p or br)", localCfg.ParagraphMode)
			}
		}
		if localCfgFound && localCfg.ParagraphBreakBRCount > 0 {
			opt.ParagraphBreakBRCount = localCfg.ParagraphBreakBRCount
		}
		if localCfgFound && localCfg.IncludeFrontMatterLead && !usedSubstackSidecar {
			opt.IncludeFrontMatterLead = true
		}
		if usedSubstackSidecar {
			opt.IncludeFeaturedImageLead = true
		}
		if localCfgFound && strings.TrimSpace(localCfg.QuoteMode) != "" {
			switch strings.TrimSpace(localCfg.QuoteMode) {
			case "blockquote":
				opt.QuoteMode = substackhtml.QuoteBlockquote
			case "monospace":
				opt.QuoteMode = substackhtml.QuoteMonospace
			case "pullquote_monospace":
				opt.QuoteMode = substackhtml.QuotePullquoteMonospace
			default:
				return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("unknown quote_mode %q (use blockquote, monospace, or pullquote_monospace)", localCfg.QuoteMode)
			}
		}
		switch tables {
		case "html":
			opt.TableMode = substackhtml.TableHTML
		case "list":
			opt.TableMode = substackhtml.TableList
		default:
			return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("unknown -tables %q", tables)
		}
		opt.SourcePath = strings.TrimSpace(mdPath)
		if opt.SourcePath != "" {
			if perm, err := substackhtml.ResolvePagePermalinkForMarkdown(opt.SourcePath); err == nil && strings.TrimSpace(perm) != "" {
				opt.PagePermalink = strings.TrimSpace(perm)
			}
		}
		if strings.TrimSpace(imageResolveOrigin) != "" {
			opt.ImageResolveOrigin = strings.TrimSpace(imageResolveOrigin)
		} else {
			opt.ImageResolveOrigin = substackbrowser.EffectiveMarkdownLeadImageResolveOrigin(localCfg)
		}
		out, err := substackhtml.Convert(raw, opt)
		if err != nil {
			return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("convert markdown: %w", err)
		}
		html = out
		if localCfgFound && substackbrowser.EffectiveIncludeCognitiveMemeticsProjectAbout(localCfg) {
			spanish := substackhtml.IsSpanishSiteLocale(meta, mdPath)
			var err error
			html, err = substackhtml.AppendCognitiveMemeticsProjectAboutHTML(html, mdPath, meta, spanish)
			if err != nil {
				return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("cognitive-memetics project about: %w", err)
			}
		}
		if localCfgFound {
			html = appendFooterHTML(html, mdPath, meta, localCfg)
		}
	}

	if strings.EqualFold(strings.TrimSpace(action), "paste-schedule") {
		if fixture {
			return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("action paste-schedule is not compatible with -fixture")
		}
		if strings.TrimSpace(htmlPath) != "" {
			return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("action paste-schedule requires -in markdown for date, tags, and categories")
		}
		cfg.ScheduleEnabled = true
		sec := scheduleSectionLabel(meta, mdPath)
		// Substack schedule fields always use the Hugo post `date` (publishing time). Substack's date
		// picker often blocks past dates; automation fills both datetime-local and the visible text field.
		dtLocal := scheduleDateTimeLocalFromYAMLDate(meta.Date)
		dateDisp := scheduleSubstackDateDisplayFromYAMLDate(meta.Date)
		cfg.Schedule = substackbrowser.ScheduleAfterContinueOptions{
			Tags:              append([]string(nil), meta.Tags...),
			SectionLabel:      sec,
			DateTimeLocal:     dtLocal,
			DateDisplay:       dateDisp,
			TickEmailSubstack: schedulePushDeliveriesDefault(localCfgFound, localCfg),
		}
		if localCfgFound {
			cfg.ScheduleDebugDOM = localCfg.ScheduleDebugDOMOnFailure
			if strings.TrimSpace(localCfg.ScheduleDebugDOMFile) != "" {
				cfg.ScheduleDebugDOMPath = strings.TrimSpace(localCfg.ScheduleDebugDOMFile)
			}
		}
	}

	if fixture {
		cfg.TargetURL = substackbrowser.FixtureDataURL()
		cfg.HTML = html
		cfg.InsertSubscribeButton = false
		return html, cfg, substackbrowser.OpenConfig{}, nil
	}

	u := strings.TrimSpace(targetURL)
	if u == "" {
		if strings.TrimSpace(pub) == "" {
			return "", cfg, substackbrowser.OpenConfig{}, fmt.Errorf("need -url or -pub for non-fixture runs")
		}
		u = substackbrowser.DefaultPublishHomeURL(pub)
		cfg.AutoCreateArticle = true
	}
	cfg.TargetURL = u
	cfg.HTML = html
	if localCfgFound && localCfg.AutoButton {
		cfg.AutoButton = true
		cfg.ButtonText = localCfg.ButtonText
		cfg.ButtonURL = strings.TrimSpace(localCfg.ButtonURL)
		if cfg.ButtonURL == "" {
			cfg.ButtonURL = categoryBrowseButtonURL(localCfg, meta, mdPath)
		}
	}
	cfg.InsertSubscribeButton = true
	if localCfgFound && localCfg.InsertSubscribeButtonAfterPaste != nil {
		cfg.InsertSubscribeButton = *localCfg.InsertSubscribeButtonAfterPaste
	}
	return html, cfg, substackbrowser.OpenConfig{}, nil
}

// underCognitiveMemeticsContent is true when the Markdown file lives under content/cognitive-memetics/
// (any lane: panel, sayings, reptilocracy, t-shirt-art, …).
func underCognitiveMemeticsContent(mdPath string) bool {
	s := strings.ToLower(filepath.ToSlash(filepath.Clean(strings.TrimSpace(mdPath))))
	return strings.Contains(s, "content/cognitive-memetics/")
}

func scheduleSectionLabel(meta substackhtml.FrontMatterMeta, mdPath string) string {
	// Substack section label must match a publication section in Substack. Cognitive-memetics posts use the
	// specialised lane (second category) when present: either the usual umbrella in categories[0] or the
	// file path under content/cognitive-memetics/. All other posts use the first category (e.g. Human-Condition
	// plus a theme hub in categories[1] still maps to Human-Condition for Substack).
	// Spanish drafts (index.es.md) map Hugo category ids through i18n/es.toml (same labels as the public site).
	if len(meta.Categories) == 0 {
		return ""
	}
	first := strings.TrimSpace(meta.Categories[0])
	raw := first
	useLane := (strings.EqualFold(first, "Cognitive-Memetics") || underCognitiveMemeticsContent(mdPath)) &&
		len(meta.Categories) >= 2
	if useLane {
		if s := strings.TrimSpace(meta.Categories[1]); s != "" {
			raw = s
		}
	}
	return substackhtml.CategorySubstackSectionLabel(raw, substackhtml.IsSpanishSiteLocale(meta, mdPath))
}

// categoryBrowseButtonURL returns site_base_url[/es/]/categories/<slug>/ from the first Hugo category, or "".
func categoryBrowseButtonURL(lc substackbrowser.LocalConfig, meta substackhtml.FrontMatterMeta, mdPath string) string {
	base := strings.TrimSpace(lc.SiteBaseURL)
	if base == "" {
		base = "https://blog.dualsubstrate.com/"
	}
	base = strings.TrimSuffix(base, "/") + "/"
	if len(meta.Categories) == 0 {
		return ""
	}
	name := strings.TrimSpace(meta.Categories[0])
	if name == "" {
		return ""
	}
	pre := hugoLangURLPathPrefix(base, substackhtml.IsSpanishSiteLocale(meta, mdPath))
	return base + pre + "categories/" + slugify(name) + "/"
}

// hugoLangURLPathPrefix returns "es/" for Spanish content when base is not already under /es/.
func hugoLangURLPathPrefix(siteBase string, spanish bool) string {
	if !spanish {
		return ""
	}
	b := strings.TrimSuffix(strings.TrimSpace(siteBase), "/")
	if strings.HasSuffix(strings.ToLower(b), "/es") {
		return ""
	}
	return "es/"
}

func truncateRunes(s string, max int) string {
	if max <= 0 {
		return strings.TrimSpace(s)
	}
	rs := []rune(s)
	if len(rs) <= max {
		return strings.TrimSpace(s)
	}
	if max <= 1 {
		return strings.TrimSpace(string(rs[:max]))
	}
	return strings.TrimSpace(string(rs[:max-1])) + "…"
}

func appendFooterHTML(html string, mdPath string, meta substackhtml.FrontMatterMeta, cfg substackbrowser.LocalConfig) string {
	var b strings.Builder
	if !cfg.FooterIncludeSiteLink && !cfg.FooterIncludeTags && !cfg.FooterIncludeCategoryLink {
		return html
	}
	b.WriteString(strings.TrimSpace(html))
	// Keep the HR, then render the footer as a compact callout (blockquote) for Substack.
	b.WriteString("<hr><blockquote>")
	base := strings.TrimSpace(cfg.SiteBaseURL)
	if base == "" {
		base = "https://blog.dualsubstrate.com/"
	}
	base = strings.TrimSuffix(base, "/") + "/"
	spanish := substackhtml.IsSpanishSiteLocale(meta, mdPath)
	langPre := hugoLangURLPathPrefix(base, spanish)
	postURL := ""
	if cfg.FooterIncludeSiteLink {
		if rel := contentPathToRelURL(mdPath); rel != "" {
			postURL = base + langPre + strings.TrimPrefix(rel, "/")
		}
	}

	// First line: "+ <Category> on <Site>" (matches the desired footer style).
	if cfg.FooterIncludeCategoryLink && len(meta.Categories) > 0 {
		idx := cfg.FooterCategoryLinkIndex
		if idx == 0 {
			idx = -1
		}
		if idx < 0 {
			idx = len(meta.Categories) + idx
		}
		if idx < 0 {
			idx = 0
		}
		if idx >= len(meta.Categories) {
			idx = len(meta.Categories) - 1
		}
		name := strings.TrimSpace(meta.Categories[idx])
		if name != "" {
			u := base + langPre + "categories/" + slugify(name) + "/"
			b.WriteString(`+ <strong><a href="` + u + `">` + htmlEscape(name) + `</a></strong>`)
			if postURL != "" {
				b.WriteString(` on <strong><a href="` + postURL + `">blog.dualsubstrate.com</a></strong>`)
			}
			b.WriteString("<br>")
		}
	} else if postURL != "" {
		b.WriteString(`Read on <strong><a href="` + postURL + `">blog.dualsubstrate.com</a></strong><br>`)
	}
	if cfg.FooterIncludeTags && len(meta.Tags) > 0 {
		b.WriteString("Tags: ")
		b.WriteString(htmlEscape(strings.Join(meta.Tags, ", ")))
		b.WriteString("<br>")
	}
	out := b.String()
	out = strings.TrimSuffix(out, "<br>")
	out = strings.TrimSuffix(out, "<br/>")
	out = strings.TrimSuffix(out, "<br />")
	out += "</blockquote>"
	return out
}

func contentPathToRelURL(mdPath string) string {
	p := filepath.ToSlash(strings.TrimSpace(mdPath))
	i := strings.Index(p, "/content/")
	if i >= 0 {
		p = p[i+len("/content/"):]
	} else {
		p = strings.TrimPrefix(p, "content/")
	}
	p = strings.TrimSuffix(p, "/index.md")
	p = strings.TrimSuffix(p, "/index.es.md")
	p = strings.Trim(p, "/")
	if p == "" {
		return ""
	}
	return "/" + p + "/"
}

func htmlEscape(s string) string {
	// Minimal escape for a plain text footer.
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&#34;")
	return s
}

// subtitleTypeCategoryLine builds "Type: slugified-category[:more...]:slugified-type" for the
// Substack subtitle when subtitle_include_categories is enabled. Content type is the last segment.
func subtitleTypeCategoryLine(categories []string, max int, contentType string) string {
	if len(categories) == 0 {
		return ""
	}
	if max <= 0 {
		max = 3
	}
	if max > len(categories) {
		max = len(categories)
	}
	categories = categories[:max]
	var segs []string
	for _, c := range categories {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		segs = append(segs, slugify(c))
	}
	if len(segs) == 0 {
		return ""
	}
	if t := strings.TrimSpace(contentType); t != "" {
		segs = append(segs, slugify(t))
	}
	return "Type: " + strings.Join(segs, ":")
}

func effectiveScheduleMaxAttempts(flagValue int, lc substackbrowser.LocalConfig, found bool) int {
	if flagValue > 0 {
		return flagValue
	}
	if found && lc.ScheduleMaxAttempts > 0 {
		return lc.ScheduleMaxAttempts
	}
	return 1
}

func schedulePushDeliveriesDefault(found bool, lc substackbrowser.LocalConfig) bool {
	if !found || lc.SchedulePushDeliveries == nil {
		return true
	}
	return *lc.SchedulePushDeliveries
}

func parseFrontMatterScheduleDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02",
		time.DateOnly,
	}
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unrecognized date layout")
}

// scheduleDateTimeLocalFromYAMLDate formats for HTML datetime-local in the machine local timezone.
func scheduleDateTimeLocalFromYAMLDate(s string) string {
	t, err := parseFrontMatterScheduleDate(s)
	if err != nil {
		return ""
	}
	t = t.In(time.Local)
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

// scheduleSubstackDateDisplayFromYAMLDate formats like Substack's text field: "29/04/2026, 08:40 am".
func scheduleSubstackDateDisplayFromYAMLDate(s string) string {
	t, err := parseFrontMatterScheduleDate(s)
	if err != nil {
		return ""
	}
	return formatSubstackScheduleTextField(t.In(time.Local))
}

func formatSubstackScheduleTextField(t time.Time) string {
	h12, suf := twelveHourSchedule(t.Hour())
	return fmt.Sprintf("%02d/%02d/%04d, %02d:%02d %s",
		t.Day(), int(t.Month()), t.Year(), h12, t.Minute(), strings.ToLower(suf))
}

func twelveHourSchedule(h int) (int, string) {
	if h == 0 {
		return 12, "am"
	}
	if h < 12 {
		return h, "am"
	}
	if h == 12 {
		return 12, "pm"
	}
	return h - 12, "pm"
}

func slugify(s string) string {
	// Hugo's taxonomy paths are usually lowercase with dashes.
	// This is a minimal slugifier for category link construction.
	s = strings.TrimSpace(strings.ToLower(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		isAZ := r >= 'a' && r <= 'z'
		is09 := r >= '0' && r <= '9'
		if isAZ || is09 {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if r == '-' {
			if !lastDash && b.Len() > 0 {
				b.WriteRune('-')
				lastDash = true
			}
			continue
		}
		// Space and other punctuation become a single dash.
		if !lastDash && b.Len() > 0 {
			b.WriteRune('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "category"
	}
	return out
}

// resolvedPublishTarget chooses the channel id for marker list/mark/pick when -publish-target is empty.
func resolvedPublishTarget(cliAction, flagVal string) string {
	if s := strings.TrimSpace(flagVal); s != "" {
		return s
	}
	if strings.EqualFold(strings.TrimSpace(cliAction), "pick-draft-es") {
		return substackpublishstate.TargetSubstackES
	}
	return substackpublishstate.LegacyDefaultTarget
}
