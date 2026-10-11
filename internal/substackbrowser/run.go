package substackbrowser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// Config controls one browser session that navigates to a URL and pastes HTML.
type Config struct {
	TargetURL  string
	HTML       string
	Title      string
	Subtitle   string
	ButtonText string
	ButtonURL  string
	AutoButton bool
	// InsertSubscribeButton runs after body paste: opens editor **Button** menu and chooses **Subscribe**.
	InsertSubscribeButton bool
	Headless              bool
	UserDataDir           string
	WaitLogin             time.Duration
	PasteTimeout          time.Duration
	KeepOpen              time.Duration
	// AutoCreateArticle clicks Create -> Article when TargetURL is /publish/home.
	// This is intended to make runs hands-off once the Chrome profile is logged in.
	AutoCreateArticle bool
	// PublishHomeSuffix is used to detect the publish home URL, default: "/publish/home".
	PublishHomeSuffix string
	// LoginTitleKeywords are case-insensitive substrings used to detect a login page.
	// Defaults to: ["sign in", "log in"].
	LoginTitleKeywords []string
	// CreateButtonText and ArticleMenuText are text fragments used for the Create -> Article flow.
	// Defaults to "Create" and "Article".
	CreateButtonText string
	ArticleMenuText  string
	// LandingURL, when set, is visited before TargetURL. This is optional and
	// intended for cases where you prefer the browser to open on the public site first.
	LandingURL string
	// NewArticleURL, when set and AutoCreateArticle is true, is navigated to directly
	// after reaching publish home. This avoids brittle menu clicks.
	NewArticleURL string
	// NavigationDelayMS inserts a delay between navigations and major steps.
	NavigationDelayMS int
	// ScheduleEnabled runs Continue plus publish settings (section, tags, delivery, schedule time).
	ScheduleEnabled bool
	// Schedule carries tags, section label, local datetime-local value, and confirmation flags.
	Schedule ScheduleAfterContinueOptions
	// ScheduleAfterPasteConfirm, when non-nil and ScheduleEnabled is true, runs after a successful paste
	// and before any schedule automation. If it returns proceed false, Run returns without clicking
	// Substack Continue or filling publish settings (the browser session still ends).
	ScheduleAfterPasteConfirm func() (proceed bool, err error)
	// ConfirmDismiss, when non-nil, runs after a successful paste and optional schedule step, before the
	// browser session ends. Receives the chromedp context (same tab as paste/schedule) so the callback may
	// inject clicks (for example Substack's purple publish footer). When set, KeepOpen is ignored.
	ConfirmDismiss func(ctx context.Context) error
	// ScheduleDebugDOM, when true and ScheduleEnabled, writes a JSON DOM snapshot if the schedule step fails
	// (or if chromedp cannot evaluate the schedule script). See WriteScheduleDebugSnapshot and substack.json / SUBSTACK_* env.
	ScheduleDebugDOM bool
	// ScheduleDebugDOMPath overrides the output file; empty uses tmp/substack-schedule-debug-<timestamp>.json under cwd.
	ScheduleDebugDOMPath string
	// ScheduleMaxAttempts is how many times to run the publish-settings automation (each run clicks Continue from the editor gate).
	// Values below 1 are treated as 1. After a failed attempt, recoverPublishFlow runs before the next attempt
	// (reload or navigate to TargetURL, then RecoverPublishFlowScheduleRetrySyncJS: Escape and small OK only).
	ScheduleMaxAttempts int
}

// recoverPublishFlow runs after a failed schedule attempt. Substack often navigates the tab during the
// long schedule script (-32000). Restoring the editor uses Navigate to a stable draft URL, or Reload when
// TargetURL is a generic /publish/post?type=… composer (Navigate would start a new blank draft).
func recoverPublishFlow(browserCtx context.Context, cfg Config) {
	u := strings.TrimSpace(cfg.TargetURL)
	if u != "" {
		recWait := 2 * time.Minute
		recCtx, cancel := context.WithTimeout(browserCtx, recWait)
		defer cancel()
		delay := chromedp.Tasks{}
		if cfg.NavigationDelayMS > 0 {
			delay = chromedp.Tasks{chromedp.Sleep(time.Duration(cfg.NavigationDelayMS) * time.Millisecond)}
		}
		waitPM := chromedp.WaitVisible(`div.ProseMirror`, chromedp.ByQuery)
		genericComposer := LooksLikeSubstackGenericNewPostEditorURL(u)
		var primary, fallback chromedp.Tasks
		if genericComposer {
			log.Printf("substackbrowser: publish flow recovery: TargetURL is a generic Substack composer; reloading this tab first (Navigate to the same URL can open another blank draft)")
			primary = chromedp.Tasks{
				chromedp.Reload(),
				chromedp.Sleep(800 * time.Millisecond),
			}
			primary = append(primary, delay...)
			primary = append(primary, waitPM)
			fallback = chromedp.Tasks{chromedp.Navigate(u)}
			fallback = append(fallback, delay...)
			fallback = append(fallback, waitPM)
		} else {
			primary = chromedp.Tasks{chromedp.Navigate(u)}
			primary = append(primary, delay...)
			primary = append(primary, waitPM)
			fallback = chromedp.Tasks{
				chromedp.Reload(),
				chromedp.Sleep(800 * time.Millisecond),
			}
			fallback = append(fallback, delay...)
			fallback = append(fallback, waitPM)
		}
		if err := chromedp.Run(recCtx, primary); err != nil {
			log.Printf("substackbrowser: publish flow recovery: primary editor restore failed: %v; trying fallback", err)
			if err2 := chromedp.Run(recCtx, fallback); err2 != nil {
				log.Printf("substackbrowser: publish flow recovery: fallback editor restore failed: %v", err2)
			} else {
				log.Printf("substackbrowser: publish flow recovery: editor back after fallback")
			}
		} else {
			if genericComposer {
				log.Printf("substackbrowser: publish flow recovery: editor back after reload (same tab)")
			} else {
				log.Printf("substackbrowser: publish flow recovery: editor back after navigate to draft URL")
			}
		}
	}
	var raw string
	if err := chromedp.Run(browserCtx, chromedp.Evaluate(RecoverPublishFlowScheduleRetrySyncJS(), &raw)); err != nil {
		log.Printf("substackbrowser: publish flow recovery: dialog cleanup evaluate: %v", err)
	} else {
		log.Printf("substackbrowser: publish flow recovery: dialog cleanup %s", strings.TrimSpace(raw))
	}
	d := cfg.NavigationDelayMS * 2
	if d < 600 {
		d = 600
	}
	_ = chromedp.Run(browserCtx, chromedp.Sleep(time.Duration(d)*time.Millisecond))
}

// Run opens Chromium, navigates to TargetURL, waits for an editor surface, then
// injects HTML. When ScheduleEnabled, it fills the publish modal (Continue, section, tags, schedule)
// but never clicks the modal's final Publish or Schedule; finish send in Substack yourself.
// Substack still autosaves drafts when the editor accepts changes.
func Run(parent context.Context, cfg Config) error {
	if strings.TrimSpace(cfg.TargetURL) == "" {
		return fmt.Errorf("substackbrowser: empty TargetURL")
	}
	if cfg.HTML == "" {
		return fmt.Errorf("substackbrowser: empty HTML")
	}
	if cfg.PasteTimeout <= 0 {
		cfg.PasteTimeout = 3 * time.Minute
	}
	cfg = withDefaults(cfg)

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.WindowSize(1400, 900),
		// chromedp often ends the browser without a full graceful shutdown; without this,
		// the next launch shows "Restore pages? Chrome didn't shut down correctly."
		chromedp.Flag("hide-crash-restore-bubble", true),
	}
	if cfg.Headless {
		opts = append(opts, chromedp.Headless)
	} else {
		opts = append(opts,
			chromedp.Flag("headless", false),
			chromedp.Flag("disable-gpu", false),
			chromedp.Flag("mute-audio", false),
		)
	}
	if strings.TrimSpace(cfg.UserDataDir) != "" {
		opts = append(opts, chromedp.UserDataDir(cfg.UserDataDir))
	}

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(parent, opts...)
	defer cancelAlloc()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	installDuplicateTagAlertDismiss(ctx)

	pasteExpr, err := PasteHTMLIntoEditor(cfg.HTML)
	if err != nil {
		return err
	}
	titleExpr, err := SetTitleAndSubtitle(cfg.Title, cfg.Subtitle)
	if err != nil {
		return err
	}
	btnExpr, err := CreateButton(cfg.ButtonText, cfg.ButtonURL)
	if err != nil {
		return err
	}
	var subscribeExpr string
	if cfg.InsertSubscribeButton {
		subscribeExpr, err = InsertSubscribeButton()
		if err != nil {
			return err
		}
	}

	pasteCtx, cancelPaste := context.WithTimeout(ctx, cfg.PasteTimeout)
	defer cancelPaste()

	var pasteJSON string
	var titleJSON string
	var btnJSON string
	var subscribeJSON string
	tasks := chromedp.Tasks{}
	if strings.TrimSpace(cfg.LandingURL) != "" {
		tasks = append(tasks, chromedp.Navigate(cfg.LandingURL))
		if cfg.NavigationDelayMS > 0 {
			tasks = append(tasks, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond))
		}
	}
	tasks = append(tasks, chromedp.Navigate(cfg.TargetURL))
	if cfg.NavigationDelayMS > 0 {
		tasks = append(tasks, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond))
	}
	if cfg.WaitLogin > 0 {
		tasks = append(tasks, chromedp.Sleep(cfg.WaitLogin))
	}
	if cfg.AutoCreateArticle && looksLikePublishHome(cfg.TargetURL, cfg.PublishHomeSuffix) {
		tasks = append(tasks, waitUntilNotOnLogin(cfg.LoginTitleKeywords))
		if strings.TrimSpace(cfg.NewArticleURL) != "" {
			tasks = append(tasks, chromedp.Navigate(cfg.NewArticleURL))
			if cfg.NavigationDelayMS > 0 {
				tasks = append(tasks, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond))
			}
		} else {
			tasks = append(tasks, clickCreateThenArticle(cfg.CreateButtonText, cfg.ArticleMenuText))
		}
	}
	if cfg.NavigationDelayMS > 0 {
		tasks = append(tasks, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond))
	}
	if strings.TrimSpace(cfg.Title) != "" || strings.TrimSpace(cfg.Subtitle) != "" {
		// Wait for the editor shell before trying to set title/subtitle.
		// Substack often mounts subtitle later than title.
		tasks = append(tasks, chromedp.WaitVisible(`div.ProseMirror`, chromedp.ByQuery))
		tasks = append(tasks, chromedp.Evaluate(titleExpr, &titleJSON, awaitPromiseEvaluate))
		if cfg.NavigationDelayMS > 0 {
			tasks = append(tasks, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond))
		}
	}
	tasks = append(tasks,
		chromedp.WaitVisible(`div.ProseMirror`, chromedp.ByQuery),
		chromedp.Evaluate(pasteExpr, &pasteJSON, awaitPromiseEvaluate),
	)
	if cfg.AutoButton && strings.TrimSpace(cfg.ButtonText) != "" && strings.TrimSpace(cfg.ButtonURL) != "" {
		if cfg.NavigationDelayMS > 0 {
			tasks = append(tasks, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond))
		}
		tasks = append(tasks, chromedp.Evaluate(btnExpr, &btnJSON, awaitPromiseEvaluate))
	}
	if cfg.InsertSubscribeButton {
		if cfg.NavigationDelayMS > 0 {
			tasks = append(tasks, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond))
		}
		tasks = append(tasks, chromedp.Evaluate(subscribeExpr, &subscribeJSON, awaitPromiseEvaluate))
	}
	if cfg.ConfirmDismiss == nil && cfg.KeepOpen > 0 {
		tasks = append(tasks, chromedp.Sleep(cfg.KeepOpen))
	}

	tasks = append(chromedp.Tasks{
		chromedp.ActionFunc(func(c context.Context) error {
			return page.Enable().Do(c)
		}),
	}, tasks...)

	if err := chromedp.Run(pasteCtx, tasks); err != nil {
		return fmt.Errorf("substackbrowser: chromedp: %w", err)
	}
	if strings.TrimSpace(titleJSON) != "" {
		var titleMeta struct {
			TitleSet    bool `json:"title_set"`
			SubtitleSet bool `json:"subtitle_set"`
		}
		if err := json.Unmarshal([]byte(strings.TrimSpace(titleJSON)), &titleMeta); err == nil {
			if strings.TrimSpace(cfg.Title) != "" && !titleMeta.TitleSet {
				log.Printf("substackbrowser: title field not set (Substack UI may use new placeholders; body paste continues)")
			}
			if strings.TrimSpace(cfg.Subtitle) != "" && !titleMeta.SubtitleSet {
				log.Printf("substackbrowser: subtitle field not set (Substack UI may use new placeholders; body paste continues)")
			}
		}
	}
	res, err := ParsePasteResult(pasteJSON)
	if err != nil {
		return err
	}
	if !res.OK {
		return fmt.Errorf("substackbrowser: paste failed: %s", res.Reason)
	}
	if cfg.InsertSubscribeButton && strings.TrimSpace(subscribeJSON) != "" {
		if sr, err := ParsePasteResult(subscribeJSON); err == nil && strings.TrimSpace(sr.Reason) != "" {
			log.Println("substackbrowser: insert subscribe:", strings.TrimSpace(sr.Reason))
		}
	}
	if cfg.ScheduleEnabled {
		if cfg.ScheduleAfterPasteConfirm != nil {
			proceed, err := cfg.ScheduleAfterPasteConfirm()
			if err != nil {
				if errors.Is(err, ErrAbortedBeforePublish) {
					return err
				}
				return fmt.Errorf("substackbrowser: schedule after-paste confirm: %w", err)
			}
			if !proceed {
				return ErrAbortedBeforePublish
			}
		}
		maxAttempts := cfg.ScheduleMaxAttempts
		if maxAttempts < 1 {
			maxAttempts = 1
		}
		if maxAttempts > 6 {
			maxAttempts = 6
		}
		// Headless UI ignores JS el.click on tag options; select tags with CDP mouse after schedule JS.
		schedOpt := cfg.Schedule
		publishTags := append([]string(nil), schedOpt.Tags...)
		schedOpt.Tags = nil
		scheduleExpr, schedErr := ScheduleAfterContinueJS(schedOpt)
		if schedErr != nil {
			return schedErr
		}
		var scheduleJSON string
		var sres PasteResult
		scheduleHadTabNavError := false
		for attempt := 1; attempt <= maxAttempts; attempt++ {
			if attempt > 1 {
				recoverPublishFlow(ctx, cfg)
			}
			var snapRaw string
			_ = chromedp.Run(pasteCtx, chromedp.Evaluate(PublishFlowStageSnapshotJS(), &snapRaw))
			LogPublishFlowSnapshot(attempt, maxAttempts, snapRaw)
			var snap PublishFlowStageSnapshot
			if err := json.Unmarshal([]byte(strings.TrimSpace(snapRaw)), &snap); err == nil {
				if attempt == 1 && !snap.EditorProseMirror {
					return fmt.Errorf("substackbrowser: schedule aborted: editor surface not detected (no visible ProseMirror); keep the draft tab in the post editor and retry")
				}
			}

			if cfg.NavigationDelayMS > 0 {
				if err := chromedp.Run(pasteCtx, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond)); err != nil {
					return fmt.Errorf("substackbrowser: schedule pre-delay: %w", err)
				}
			}
			log.Println("substackbrowser: running publish settings (Substack Continue, section, tags, schedule)...")
			scheduleJSON = ""
			if err := chromedp.Run(pasteCtx, chromedp.Evaluate(scheduleExpr, &scheduleJSON, awaitPromiseEvaluate)); err != nil {
				if isChromeTargetGone(err) {
					scheduleHadTabNavError = true
				}
				if attempt < maxAttempts {
					log.Printf("substackbrowser: schedule chromedp error on attempt %d/%d: %v", attempt, maxAttempts, err)
					if isChromeTargetGone(err) {
						log.Printf("substackbrowser: schedule: tab navigated or context destroyed mid-script; next attempt reloads TargetURL then retries publish settings")
					}
					continue
				}
				if cfg.ScheduleDebugDOM {
					p := strings.TrimSpace(cfg.ScheduleDebugDOMPath)
					if p == "" {
						p = DefaultScheduleDebugSnapshotPath()
					}
					if werr := WriteScheduleDebugSnapshot(pasteCtx, p, "chromedp schedule evaluate: "+err.Error()); werr != nil {
						log.Printf("substackbrowser: schedule debug DOM not written: %v", werr)
					} else {
						log.Printf("substackbrowser: wrote schedule debug DOM after evaluate error: %s", p)
					}
				}
				if isChromeTargetGone(err) {
					log.Printf("substackbrowser: schedule: hint: CDP lost the tab while the schedule script was running; Substack may already have scheduled or sent the post; confirm in Substack Posts before you edit or re-run")
				}
				return fmt.Errorf("substackbrowser: schedule step: %w", err)
			}
			sres, err = ParsePasteResult(scheduleJSON)
			if err != nil {
				if attempt < maxAttempts {
					log.Printf("substackbrowser: schedule result parse error on attempt %d/%d: %v", attempt, maxAttempts, err)
					continue
				}
				return err
			}
			if !sres.OK {
				if attempt < maxAttempts {
					log.Printf("substackbrowser: schedule failed on attempt %d/%d: %s; retrying after recovery", attempt, maxAttempts, sres.Reason)
					continue
				}
				if cfg.ScheduleDebugDOM {
					p := strings.TrimSpace(cfg.ScheduleDebugDOMPath)
					if p == "" {
						p = DefaultScheduleDebugSnapshotPath()
					}
					if werr := WriteScheduleDebugSnapshot(pasteCtx, p, sres.Reason); werr != nil {
						log.Printf("substackbrowser: schedule debug DOM not written: %v", werr)
					} else {
						log.Printf("substackbrowser: wrote schedule debug DOM: %s", p)
					}
				}
				if scheduleHadTabNavError {
					log.Printf("substackbrowser: schedule: hint: Chrome lost the tab mid-run (-32000) on an earlier attempt; Substack may already have scheduled or sent the post; check recent posts before re-running publish settings")
				}
				return fmt.Errorf("substackbrowser: schedule step failed: %s", sres.Reason)
			}
			if strings.TrimSpace(sres.Reason) != "" {
				log.Println("substackbrowser: schedule notes:", strings.TrimSpace(sres.Reason))
			}
			if len(publishTags) > 0 {
				LogTagCDPStart(len(publishTags))
				if tagErr := AddPublishTagsCDP(pasteCtx, publishTags); tagErr != nil {
					if attempt < maxAttempts {
						log.Printf("substackbrowser: CDP tag select failed on attempt %d/%d: %v; retrying after recovery", attempt, maxAttempts, tagErr)
						continue
					}
					if cfg.ScheduleDebugDOM {
						p := strings.TrimSpace(cfg.ScheduleDebugDOMPath)
						if p == "" {
							p = DefaultScheduleDebugSnapshotPath()
						}
						if werr := WriteScheduleDebugSnapshot(pasteCtx, p, tagErr.Error()); werr != nil {
							log.Printf("substackbrowser: schedule debug DOM not written: %v", werr)
						} else {
							log.Printf("substackbrowser: wrote schedule debug DOM after tag failure: %s", p)
						}
					}
					return fmt.Errorf("substackbrowser: schedule tags: %w", tagErr)
				}
			}
			log.Printf("substackbrowser: publish settings automation returned OK (section=%q, tags=%d, email_delivery=%v, schedule_fields=%v, attempts_used=%d)",
				strings.TrimSpace(cfg.Schedule.SectionLabel),
				len(publishTags),
				cfg.Schedule.TickEmailSubstack,
				strings.TrimSpace(cfg.Schedule.DateTimeLocal) != "" || strings.TrimSpace(cfg.Schedule.DateDisplay) != "",
				attempt,
			)
			break
		}
	}
	if cfg.ConfirmDismiss != nil {
		if err := cfg.ConfirmDismiss(pasteCtx); err != nil {
			if errors.Is(err, ErrAbortedBeforePublish) {
				return err
			}
			return fmt.Errorf("substackbrowser: confirm dismiss: %w", err)
		}
	}
	return nil
}

// DefaultPublishHomeURL returns the writer home URL for a publication subdomain.
func DefaultPublishHomeURL(pub string) string {
	p := strings.TrimSpace(pub)
	p = strings.TrimSuffix(p, "/")
	p = strings.TrimSuffix(p, ".substack.com")
	p = strings.TrimSpace(p)
	return "https://" + p + ".substack.com/publish/home"
}

func looksLikePublishHome(u string, suffix string) bool {
	uu := strings.TrimSpace(u)
	uu = strings.TrimSuffix(uu, "/")
	s := strings.TrimSpace(suffix)
	if s == "" {
		s = "/publish/home"
	}
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	return strings.HasSuffix(uu, s)
}

func waitUntilNotOnLogin(keywords []string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		tick := time.NewTicker(1 * time.Second)
		defer tick.Stop()
		for {
			var title string
			if err := chromedp.Title(&title).Do(ctx); err != nil {
				return err
			}
			if !looksLikeLoginTitle(title, keywords) {
				return nil
			}
			select {
			case <-ctx.Done():
				return fmt.Errorf("substackbrowser: still on sign-in page; log in using CHROME_PROFILE then re-run")
			case <-tick.C:
			}
		}
	})
}

func clickCreateThenArticle(createText string, articleText string) chromedp.Tasks {
	c := strings.TrimSpace(createText)
	if c == "" {
		c = "Create"
	}
	a := strings.TrimSpace(articleText)
	if a == "" {
		a = "Article"
	}

	return chromedp.Tasks{
		clickCreateThenMenuItem(c, a),
	}
}

func clickCreateThenMenuItem(createText string, itemText string) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var raw string
		if err := chromedp.Evaluate(clickCreateThenMenuItemJS(createText, itemText), &raw, awaitPromiseEvaluate).Do(ctx); err != nil {
			return err
		}
		r, err := ParsePasteResult(raw)
		if err != nil {
			return fmt.Errorf("substackbrowser: click create/menu: %w", err)
		}
		if !r.OK {
			return fmt.Errorf("substackbrowser: click create/menu: %s", r.Reason)
		}
		return nil
	})
}

func clickCreateThenMenuItemJS(createText string, itemText string) string {
	c := strings.TrimSpace(createText)
	i := strings.TrimSpace(itemText)
	esc := func(s string) string {
		s = strings.ReplaceAll(s, "\\", "\\\\")
		s = strings.ReplaceAll(s, `"`, `\"`)
		return s
	}
	return `(function(){` +
		`const create="` + esc(c) + `".trim().toLowerCase();` +
		`const item="` + esc(i) + `".trim().toLowerCase();` +
		`function visible(el){if(!el) return false; const r=el.getBoundingClientRect(); return r.width>0 && r.height>0;}` +
		`function bestExact(target){` +
		`  const els=Array.from(document.querySelectorAll('button,[role="button"],a,div,span')).filter(el=>{` +
		`    if(!visible(el)) return false; const t=(el.innerText||"").trim().toLowerCase(); return t===target;` +
		`  });` +
		`  if(els.length===0) return null;` +
		`  els.sort((a,b)=>{const ra=a.getBoundingClientRect();const rb=b.getBoundingClientRect();return (rb.width*rb.height)-(ra.width*ra.height);});` +
		`  return els[0];` +
		`}` +
		`function openMenus(){` +
		`  const menus=Array.from(document.querySelectorAll('[role="menu"],[data-radix-popper-content-wrapper],.dropdown-menu,div[aria-label*="Create"],div[aria-label*="create"]')).filter(visible);` +
		`  return menus;` +
		`}` +
		`function findInMenu(menus,target){` +
		`  const hits=[];` +
		`  for(const m of menus){` +
		`    const els=Array.from(m.querySelectorAll('button,[role="menuitem"],[role="menuitemradio"],a,div,span')).filter(el=>visible(el));` +
		`    for(const el of els){const t=(el.innerText||"").trim().toLowerCase(); if(t===target) hits.push(el);}` +
		`  }` +
		`  if(hits.length===0) return null;` +
		`  hits.sort((a,b)=>{const ra=a.getBoundingClientRect();const rb=b.getBoundingClientRect();return (rb.width*rb.height)-(ra.width*ra.height);});` +
		`  return hits[0];` +
		`}` +
		`const btn=bestExact(create);` +
		`if(!btn) return JSON.stringify({ok:false,reason:"Create button not found"});` +
		`try{btn.scrollIntoView({block:"center"});}catch(e){}` +
		`try{btn.click();}catch(e){return JSON.stringify({ok:false,reason:String(e)});}` +
		`// Give the menu a tick to appear (sync spin with small upper bound)` +
		`let menus=[];` +
		`for(let k=0;k<50;k++){menus=openMenus(); if(menus.length>0) break;}` +
		`if(menus.length===0){return JSON.stringify({ok:false,reason:"Create menu did not open"});}` +
		`const mi=findInMenu(menus,item);` +
		`if(!mi){return JSON.stringify({ok:false,reason:"Menu item not found: "+item});}` +
		`try{mi.scrollIntoView({block:"center"});}catch(e){}` +
		`try{mi.click();}catch(e){return JSON.stringify({ok:false,reason:String(e)});}` +
		`return JSON.stringify({ok:true,reason:""});` +
		`})()`
}

func clickByExactVisibleText(text string) chromedp.Action {
	t := strings.TrimSpace(text)
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var res string
		if err := chromedp.Evaluate(clickByExactTextJS(t), &res, awaitPromiseEvaluate).Do(ctx); err != nil {
			return err
		}
		r, err := ParsePasteResult(res)
		if err != nil {
			// Reuse ParsePasteResult shape: ok/reason JSON.
			return fmt.Errorf("substackbrowser: clickByExactVisibleText: %w", err)
		}
		if !r.OK {
			// Do not fail hard, the caller may fall back to XPath.
			return nil
		}
		return nil
	})
}

func clickByExactTextJS(text string) string {
	// Returns JSON {ok, reason}.
	// Strategy: scan clickable elements, match innerText exactly (case-insensitive),
	// pick the one with the largest area (usually the real button), and click it.
	esc := strings.ReplaceAll(strings.ReplaceAll(text, "\\", "\\\\"), `"`, `\"`)
	return `(function(){` +
		`const target = "` + esc + `".trim().toLowerCase();` +
		`if(!target){return JSON.stringify({ok:false,reason:"empty target"});}` +
		`const candidates = Array.from(document.querySelectorAll('button,[role="button"],a,div,span')).filter(el=>{` +
		`  if(!el || !el.innerText) return false;` +
		`  const t = el.innerText.trim().toLowerCase();` +
		`  if(t !== target) return false;` +
		`  const r = el.getBoundingClientRect();` +
		`  return r.width>0 && r.height>0;` +
		`});` +
		`if(candidates.length===0){return JSON.stringify({ok:false,reason:"no exact match"});}` +
		`candidates.sort((a,b)=>{const ra=a.getBoundingClientRect();const rb=b.getBoundingClientRect();return (rb.width*rb.height)-(ra.width*ra.height);});` +
		`const el=candidates[0];` +
		`try{el.scrollIntoView({block:"center"});}catch(e){}` +
		`try{el.click();}catch(e){return JSON.stringify({ok:false,reason:String(e)});}` +
		`return JSON.stringify({ok:true,reason:""});` +
		`})()`
}

func isVisibleTextPresentJS(substr string) string {
	esc := strings.ReplaceAll(strings.ReplaceAll(substr, "\\", "\\\\"), `"`, `\"`)
	return `(function(){` +
		`const s="` + esc + `".trim().toLowerCase();` +
		`if(!s) return false;` +
		`return document.body && document.body.innerText && document.body.innerText.toLowerCase().includes(s);` +
		`})()`
}

func xpathRoleOrElementContainsTextCI(text string, createButton bool) string {
	// Build a case-insensitive contains() by translating A-Z to a-z in XPath.
	needle := strings.ToLower(strings.TrimSpace(text))
	if needle == "" {
		needle = "create"
	}

	// Prefer a button-like element for "Create", but allow role="button" wrappers.
	var elements string
	if createButton {
		elements = "self::button or self::a or self::div or self::span"
	} else {
		// Menu items in Substack can be buttons, links, or role="menuitem".
		elements = "self::button or self::a or self::div or self::span"
	}

	// normalize-space(.) reduces whitespace differences.
	// translate() lowercases the content for a cheap case-insensitive match.
	return fmt.Sprintf(
		`//*[ (%s or @role="button" or @role="menuitem" or @role="menuitemradio") and contains(translate(normalize-space(.), "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "abcdefghijklmnopqrstuvwxyz"), %q) ]`,
		elements,
		needle,
	)
}

func looksLikeLoginTitle(title string, keywords []string) bool {
	lt := strings.ToLower(strings.TrimSpace(title))
	if lt == "" {
		return false
	}
	ks := keywords
	if len(ks) == 0 {
		ks = []string{"sign in", "log in"}
	}
	for _, k := range ks {
		kk := strings.ToLower(strings.TrimSpace(k))
		if kk == "" {
			continue
		}
		if strings.Contains(lt, kk) {
			return true
		}
	}
	return false
}

// awaitPromiseEvaluate is a chromedp.Evaluate option so expressions that return
// a Promise (for example async IIFEs in title/schedule scripts) are awaited
// before the JSON result is unmarshaled into Go strings.
func awaitPromiseEvaluate(p *runtime.EvaluateParams) *runtime.EvaluateParams {
	if p == nil {
		return nil
	}
	return p.WithAwaitPromise(true)
}

func withDefaults(cfg Config) Config {
	if strings.TrimSpace(cfg.PublishHomeSuffix) == "" {
		cfg.PublishHomeSuffix = "/publish/home"
	}
	if len(cfg.LoginTitleKeywords) == 0 {
		cfg.LoginTitleKeywords = []string{"sign in", "log in"}
	}
	if strings.TrimSpace(cfg.CreateButtonText) == "" {
		cfg.CreateButtonText = "Create"
	}
	if strings.TrimSpace(cfg.ArticleMenuText) == "" {
		cfg.ArticleMenuText = "Article"
	}
	return cfg
}
