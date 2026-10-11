package substackbrowser

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// OpenConfig controls a browser session that only navigates (no paste).
type OpenConfig struct {
	// URL is the primary page to navigate to.
	URL string
	// LandingURL, when set, is visited before URL.
	LandingURL string

	Headless    bool
	UserDataDir string

	// WaitUntilLoggedIn waits until the page is not detected as a login page.
	WaitUntilLoggedIn  bool
	LoginTitleKeywords []string

	// Timeout bounds the entire navigation and optional login wait.
	Timeout time.Duration
	// KeepOpen keeps the browser open after navigation.
	KeepOpen time.Duration

	// NavigationDelayMS inserts a delay between navigations and major steps.
	NavigationDelayMS int
}

func Open(parent context.Context, cfg OpenConfig) error {
	if strings.TrimSpace(cfg.URL) == "" {
		return fmt.Errorf("substackbrowser: empty URL")
	}
	if len(cfg.LoginTitleKeywords) == 0 {
		cfg.LoginTitleKeywords = []string{"sign in", "log in"}
	}

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.WindowSize(1400, 900),
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

	runCtx := ctx
	cancelRun := func() {}
	if cfg.Timeout > 0 {
		runCtx, cancelRun = context.WithTimeout(ctx, cfg.Timeout)
	}
	defer cancelRun()

	tasks := chromedp.Tasks{}
	if strings.TrimSpace(cfg.LandingURL) != "" {
		tasks = append(tasks, chromedp.Navigate(cfg.LandingURL))
		if cfg.NavigationDelayMS > 0 {
			tasks = append(tasks, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond))
		}
	}
	tasks = append(tasks, chromedp.Navigate(cfg.URL))
	if cfg.NavigationDelayMS > 0 {
		tasks = append(tasks, chromedp.Sleep(time.Duration(cfg.NavigationDelayMS)*time.Millisecond))
	}
	if cfg.WaitUntilLoggedIn {
		tasks = append(tasks, waitUntilNotOnLogin(cfg.LoginTitleKeywords))
	}
	if cfg.KeepOpen > 0 {
		tasks = append(tasks, chromedp.Sleep(cfg.KeepOpen))
	} else {
		// Keep the window open until the user stops the command (Ctrl+C) or
		// the browser disconnects.
		tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		}))
	}

	if err := chromedp.Run(runCtx, tasks); err != nil {
		return fmt.Errorf("substackbrowser: chromedp: %w", err)
	}
	return nil
}
