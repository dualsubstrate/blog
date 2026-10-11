package substackbrowser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

// ScheduleDebugSnapshotJS returns synchronous browser JS that evaluates to a JSON string
// describing the publish dialog, cmdk rows, and portaled lists. Used when schedule automation fails.
func ScheduleDebugSnapshotJS() string {
	return `(function(){
    function clip(s, max) {
      if (s == null || s === undefined) return '';
      s = String(s);
      if (s.length <= max) return s;
      return s.slice(0, max) + '\n... [clipped ' + (s.length - max) + ' chars]';
    }
    function visible(el) {
      if (!el) return false;
      const r = el.getBoundingClientRect();
      return r.width > 0 && r.height > 0;
    }
    function publishDialog() {
      const dialogs = Array.from(document.querySelectorAll('[role="dialog"]')).filter(visible);
      let best = null;
      let bestScore = -1;
      for (const d of dialogs) {
        const txt = (d.innerText || '').toLowerCase();
        const looks =
          txt.includes('add tags') ||
          txt.includes('this post belongs') ||
          txt.includes('delivery') ||
          (txt.includes('scheduling') && txt.includes('email'));
        if (!looks) continue;
        let score = 0;
        if (txt.includes('this post belongs')) score += 3;
        if (txt.includes('add tags')) score += 2;
        if (txt.includes('delivery')) score += 2;
        if (txt.includes('publish')) score += 1;
        if (score > bestScore) {
          bestScore = score;
          best = d;
        }
      }
      return best;
    }
    const ae = document.activeElement;
    const active = ae
      ? {
          tag: ae.tagName,
          id: clip(ae.id, 160),
          className: clip(ae.className, 240),
          outerHTML: clip(ae.outerHTML, 2400),
        }
      : null;
    const dlg = publishDialog();
    const listboxOptions = [];
    document.querySelectorAll('[role="listbox"] [role="option"]').forEach((el, i) => {
      if (i > 80) return;
      if (!visible(el)) return;
      listboxOptions.push({
        i: i,
        text: clip(String((el.innerText || '').replace(/\s+/g, ' ').trim()), 220),
        ariaSelected: el.getAttribute('aria-selected'),
        dataHeadlessUI: el.getAttribute('data-headlessui-state'),
        outerHTML: clip(el.outerHTML, 1400),
      });
    });
    const cmdk = [];
    document.querySelectorAll('[cmdk-item]').forEach((el, i) => {
      if (i > 100) return;
      if (!visible(el)) return;
      cmdk.push({
        i: i,
        text: clip(String((el.innerText || '').replace(/\s+/g, ' ').trim()), 220),
        dataSelected: el.getAttribute('data-selected'),
        ariaSelected: el.getAttribute('aria-selected'),
        dataHighlighted: el.getAttribute('data-highlighted'),
        outerHTML: clip(el.outerHTML, 1400),
      });
    });
    const portals = [];
    document.querySelectorAll('[data-radix-popper-content-wrapper], [data-radix-portal]').forEach((el, i) => {
      if (i > 16) return;
      if (!visible(el)) return;
      const r = el.getBoundingClientRect();
      portals.push({ i: i, w: r.width, h: r.height, outerHTML: clip(el.outerHTML, 9000) });
    });
    const lists = [];
    document.querySelectorAll('[cmdk-list], [cmdk-root]').forEach((el, i) => {
      if (i > 10) return;
      if (!visible(el)) return;
      lists.push({
        tag: el.tagName,
        id: clip(el.id, 100),
        className: clip(el.className, 200),
        outerHTML: clip(el.outerHTML, 14000),
      });
    });
    const tagInputs = [];
    const roots = dlg ? [dlg] : [document.body];
    roots.forEach((root) => {
      if (!root || !root.querySelectorAll) return;
      root.querySelectorAll('input[type="text"], input[type="search"], input:not([type]), input[role="combobox"]').forEach((inp, j) => {
        if (j > 20) return;
        if (!visible(inp) && inp !== document.activeElement) return;
        const ph = (inp.getAttribute('placeholder') || '').toLowerCase();
        const al = (inp.getAttribute('aria-label') || '').toLowerCase();
        const role = (inp.getAttribute('role') || '').toLowerCase();
        const id = String(inp.id || '');
        const parCls = String((inp.parentElement && inp.parentElement.className) || '');
        const looksTag =
          ph.includes('tag') ||
          (ph.includes('select') && ph.includes('create')) ||
          al.includes('tag') ||
          role === 'combobox' ||
          id.indexOf('headlessui-combobox-input') >= 0 ||
          parCls.indexOf('inputBox') >= 0 ||
          parCls.indexOf('hasChips') >= 0;
        if (!looksTag) return;
        tagInputs.push({
          placeholder: clip(inp.getAttribute('placeholder'), 120),
          ariaLabel: clip(inp.getAttribute('aria-label'), 120),
          role: clip(role, 40),
          id: clip(id, 120),
          value: clip(inp.value, 200),
          ariaExpanded: clip(inp.getAttribute('aria-expanded'), 20),
          outerHTML: clip(inp.outerHTML, 1600),
        });
      });
    });
    return JSON.stringify({
      capturedAtISO: new Date().toISOString(),
      href: String(location.href || ''),
      title: String(document.title || ''),
      activeElement: active,
      publishDialogFound: !!dlg,
      publishDialogOuterHTML: dlg ? clip(dlg.outerHTML, 200000) : '',
      listboxOptions: listboxOptions,
      cmdkItems: cmdk,
      radixPortalsSample: portals,
      cmdkListsSample: lists,
      tagInputsInDialog: tagInputs,
    });
  })()`
}

// DefaultScheduleDebugSnapshotPath returns a file path under tmp/ in the current working directory.
func DefaultScheduleDebugSnapshotPath() string {
	return filepath.Join("tmp", fmt.Sprintf("substack-schedule-debug-%s.json", time.Now().Format("20060102-150405")))
}

// WriteScheduleDebugSnapshot runs ScheduleDebugSnapshotJS in the browser and writes a JSON file
// that includes scheduleFailureReason and the snapshot object. DestPath directories are created as needed.
func WriteScheduleDebugSnapshot(ctx context.Context, destPath, scheduleFailureReason string) error {
	var raw string
	if err := chromedp.Run(ctx, chromedp.Evaluate(ScheduleDebugSnapshotJS(), &raw)); err != nil {
		return fmt.Errorf("substackbrowser: schedule debug DOM evaluate: %w", err)
	}
	var inner map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &inner); err != nil {
		return fmt.Errorf("substackbrowser: schedule debug DOM inner JSON: %w", err)
	}
	out := struct {
		WrittenAt             string          `json:"writtenAt"`
		ScheduleFailureReason string          `json:"scheduleFailureReason"`
		Snapshot              json.RawMessage `json:"snapshot"`
	}{
		WrittenAt:             time.Now().UTC().Format(time.RFC3339),
		ScheduleFailureReason: scheduleFailureReason,
		Snapshot:              json.RawMessage(raw),
	}
	b, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return fmt.Errorf("substackbrowser: schedule debug DOM marshal: %w", err)
	}
	dir := filepath.Dir(destPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("substackbrowser: schedule debug DOM mkdir: %w", err)
		}
	}
	if err := os.WriteFile(destPath, b, 0o644); err != nil {
		return fmt.Errorf("substackbrowser: schedule debug DOM write: %w", err)
	}
	return nil
}
