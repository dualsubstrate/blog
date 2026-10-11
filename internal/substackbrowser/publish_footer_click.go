package substackbrowser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// publishFooterSharedJS is inlined in sync snippets: visible, publishModalRoot, confirmSchedule.
const publishFooterSharedJS = `
  function visible(el) {
    if (!el) return false;
    const r = el.getBoundingClientRect();
    return r.width > 0 && r.height > 0;
  }

  function publishModalRoot() {
    const direct = document.querySelector('[data-testid="publish-modal"]');
    if (direct && visible(direct)) return direct;
    const dialogs = Array.from(document.querySelectorAll('[role="dialog"]')).filter(visible);
    let best = null;
    let bestScore = -1;
    for (const d of dialogs) {
      const txt = (d.innerText || '').toLowerCase();
      const looksPublish =
        txt.includes('add tags') ||
        txt.includes('this post belongs') ||
        txt.includes('delivery') ||
        (txt.includes('scheduling') && (txt.includes('email') || txt.includes('schedule time')));
      if (!looksPublish) continue;
      let score = 0;
      if (txt.includes('this post belongs')) score += 3;
      if (txt.includes('add tags')) score += 2;
      if (txt.includes('delivery')) score += 2;
      if (txt.includes('publish')) score += 1;
      const r = d.getBoundingClientRect();
      score += Math.min(4, Math.floor((r.width * r.height) / 250000));
      if (score > bestScore) {
        bestScore = score;
        best = d;
      }
    }
    return best;
  }

  function confirmSchedule() {
    const modal = publishModalRoot();
    const footer = modal && modal.querySelector ? modal.querySelector('[data-modal-role="footer"]') : null;
    const scope = footer || modal || document;
    const skip = new Set(['cancel', 'close', 'back']);
    const hintEl = scope.querySelector && scope.querySelector('input[data-track-input="publish_button_text"]');
    const hintVal = (hintEl && hintEl.value && String(hintEl.value).trim().toLowerCase()) || '';
    if (hintVal) {
      const btns = Array.from(scope.querySelectorAll('button,[role="button"]')).filter(visible);
      for (const el of btns) {
        const t = (el.innerText || '').trim().toLowerCase().replace(/\s+/g, ' ');
        if (!t || skip.has(t)) continue;
        if (t === hintVal || t.includes(hintVal) || hintVal.includes(t)) {
          try {
            el.click();
            return true;
          } catch (e) {}
        }
      }
    }
    const cands = Array.from(scope.querySelectorAll('button,[role="button"]')).filter(el => {
      if (!visible(el)) return false;
      const t = (el.innerText || '').trim().toLowerCase();
      if (!t || t.length > 120) return false;
      if (skip.has(t)) return false;
      if (t === 'publish now' || t === 'publish') return true;
      if (t.includes('send to everyone')) return true;
      if (t.includes('publish') && (t.includes('second') || t.includes('now'))) return true;
      if (t === 'schedule' || (t.startsWith('schedule') && t.length < 80)) return true;
      return false;
    });
    if (cands.length === 0) return false;
    cands.sort((a, b) => {
      const ra = a.getBoundingClientRect();
      const rb = b.getBoundingClientRect();
      return (rb.width * rb.height) - (ra.width * ra.height);
    });
    try {
      cands[0].click();
      return true;
    } catch (e) {
      return false;
    }
  }
`

// clickPublishModalFooterSyncJS clicks the publish modal footer once (no Promise, no long waits).
// Substack may navigate immediately after click; CDP must not await a long async handler on the same target.
func clickPublishModalFooterSyncJS() string {
	return `(function() {
` + publishFooterSharedJS + `
  try {
    if (!publishModalRoot()) {
      return JSON.stringify({ ok: false, reason: 'publish modal not found; open Publish in Substack first' });
    }
    if (!confirmSchedule()) {
      return JSON.stringify({ ok: false, reason: 'footer primary button not found in publish modal' });
    }
    return JSON.stringify({ ok: true, reason: '' });
  } catch (e) {
    return JSON.stringify({ ok: false, reason: String(e) });
  }
})()`
}

// publishModalFooterPollStepJS tries to dismiss the email nudge once, then reports whether the publish modal is gone.
func publishModalFooterPollStepJS() string {
	return `(function() {
` + publishFooterSharedJS + `
  function tryDismissEmailNudge(wantEmailDelivery) {
    const body = (document.body && document.body.innerText || '').toLowerCase();
    const isNudge = body.includes('send this post via email') || body.includes('less than 1% of subscribers');
    if (!isNudge) return;
    const wantWeb = !wantEmailDelivery;
    const btns = Array.from(document.querySelectorAll('button,[role="button"]')).filter(visible);
    for (const btn of btns) {
      const t = (btn.innerText || '').trim().toLowerCase().replace(/\s+/g, ' ');
      if (wantWeb && t.includes('publish on web only')) {
        try { btn.click(); } catch (e) {}
        return;
      }
      if (!wantWeb && t.includes('also send via email')) {
        try { btn.click(); } catch (e) {}
        return;
      }
    }
  }
  try {
    tryDismissEmailNudge(false);
    const m = publishModalRoot();
    let gone = true;
    if (m && visible(m)) {
      const st = (m.getAttribute && m.getAttribute('data-state')) || '';
      gone = st === 'closed';
      if (!gone) gone = false;
    }
    return JSON.stringify({ gone: gone });
  } catch (e) {
    return JSON.stringify({ gone: true, err: String(e) });
  }
})()`
}

// isChromeTargetGone reports CDP errors when the tab navigates or closes mid-command (common after Publish).
func isChromeTargetGone(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "-32000") ||
		strings.Contains(s, "Inspected target navigated") ||
		strings.Contains(s, "target closed") ||
		strings.Contains(s, "execution context was destroyed")
}

type publishModalPollResult struct {
	Gone bool `json:"gone"`
}

// ClickPublishModalFooter clicks Substack's purple publish/schedule footer button in the open tab.
// Uses short synchronous Evaluates plus Go-side polling so a post-click navigation does not fail the run.
func ClickPublishModalFooter(ctx context.Context) error {
	var raw string
	err := chromedp.Run(ctx, chromedp.Evaluate(clickPublishModalFooterSyncJS(), &raw))
	if err != nil && isChromeTargetGone(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("substackbrowser: click publish footer: %w", err)
	}
	r, err := ParsePasteResult(raw)
	if err != nil {
		return fmt.Errorf("substackbrowser: click publish footer result: %w", err)
	}
	if !r.OK {
		return fmt.Errorf("substackbrowser: click publish footer: %s", r.Reason)
	}

	deadline := time.Now().Add(24 * time.Second)
	for time.Now().Before(deadline) {
		var pollRaw string
		err := chromedp.Run(ctx, chromedp.Evaluate(publishModalFooterPollStepJS(), &pollRaw))
		if err != nil && isChromeTargetGone(err) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("substackbrowser: click publish footer poll: %w", err)
		}
		var pr publishModalPollResult
		if jerr := json.Unmarshal([]byte(pollRaw), &pr); jerr != nil {
			return fmt.Errorf("substackbrowser: click publish footer poll decode: %w", jerr)
		}
		if pr.Gone {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(400 * time.Millisecond):
		}
	}
	return fmt.Errorf("substackbrowser: click publish footer: publish modal still open after footer click (timeout)")
}
