package substackbrowser

import (
	"context"
	"log"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// substackTagConflictAlert matches duplicate-tag copy Substack shows in window.alert.
func substackTagConflictAlert(msg string) bool {
	m := strings.ToLower(strings.TrimSpace(msg))
	if strings.Contains(m, "tag already set") || strings.Contains(m, "tag already exists") {
		return true
	}
	if strings.Contains(m, "already exists") && strings.Contains(m, "tag") {
		return true
	}
	return false
}

// installDuplicateTagAlertDismiss registers a target listener that accepts Substack's
// duplicate-tag alerts (e.g. "Tag already set", "Tag already exists") so the tab does not
// stall chromedp. The listener callback must not block; HandleJavaScriptDialog runs in a
// new goroutine per chromedp docs.
func installDuplicateTagAlertDismiss(ctx context.Context) {
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		e, ok := ev.(*page.EventJavascriptDialogOpening)
		if !ok {
			return
		}
		msg := strings.ToLower(strings.TrimSpace(e.Message))
		if !substackTagConflictAlert(msg) {
			return
		}
		c := ctx
		go func() {
			if err := chromedp.Run(c, page.HandleJavaScriptDialog(true)); err != nil {
				log.Printf("substackbrowser: duplicate-tag alert dismiss: %v", err)
			}
		}()
	})
}
