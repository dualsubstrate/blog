package substackbrowser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CreateButton returns JS that opens Substack's editor "More" menu, selects "Create a button",
// fills the dialog, and confirms. Best-effort: returns ok=true even if the UI isn't found.
func CreateButton(text string, url string) (string, error) {
	tEnc, err := json.Marshal(strings.TrimSpace(text))
	if err != nil {
		return "", fmt.Errorf("substackbrowser: encode button text: %w", err)
	}
	uEnc, err := json.Marshal(strings.TrimSpace(url))
	if err != nil {
		return "", fmt.Errorf("substackbrowser: encode button url: %w", err)
	}

	js := `(function () {
  const text = ` + string(tEnc) + `;
  const url = ` + string(uEnc) + `;

  function clickByText(root, want) {
    const w = (want || '').trim().toLowerCase();
    if (!root || !w) return false;
    const els = root.querySelectorAll('button, [role="menuitem"], [role="button"], a, div');
    for (const el of els) {
      const t = (el.textContent || '').trim().toLowerCase();
      if (t === w) {
        el.click();
        return true;
      }
    }
    return false;
  }

  function pickMoreButton() {
    // Try common toolbar patterns.
    const candidates = [
      'button[aria-label="More"]',
      'button:has(svg[aria-label="More"])',
      'button:has(span):contains("More")'
    ];
    for (const sel of candidates) {
      try {
        const el = document.querySelector(sel);
        if (el) return el;
      } catch (e) {}
    }
    // Fallback: exact text match.
    const btns = Array.from(document.querySelectorAll('button'));
    for (const b of btns) {
      if ((b.textContent || '').trim() === 'More') return b;
    }
    return null;
  }

  function pickDialog() {
    const dlg = document.querySelector('[role="dialog"]');
    if (dlg) return dlg;
    return null;
  }

  const more = pickMoreButton();
  if (!more) {
    return JSON.stringify({ ok: true, reason: 'more button not found' });
  }
  more.click();

  // Menu item: Create a button.
  let clicked = false;
  const menus = document.querySelectorAll('[role="menu"], [data-radix-popper-content-wrapper], .menu, .dropdown');
  for (const m of menus) {
    if (clickByText(m, 'Create a button')) { clicked = true; break; }
  }
  if (!clicked) {
    // Broad fallback.
    clicked = clickByText(document.body, 'Create a button');
  }
  if (!clicked) {
    return JSON.stringify({ ok: true, reason: 'create a button item not found' });
  }

  // Fill dialog.
  const dlg = pickDialog();
  if (!dlg) {
    return JSON.stringify({ ok: true, reason: 'button dialog not found' });
  }
  const textEl = dlg.querySelector('input[placeholder="Enter text..."], textarea[placeholder="Enter text..."], input');
  const urlEl = dlg.querySelector('input[placeholder="Enter URL..."], textarea[placeholder="Enter URL..."]');
  if (textEl) {
    textEl.focus();
    textEl.value = text;
    textEl.dispatchEvent(new Event('input', { bubbles: true }));
  }
  if (urlEl) {
    urlEl.focus();
    urlEl.value = url;
    urlEl.dispatchEvent(new Event('input', { bubbles: true }));
  }
  // Confirm.
  const okBtn = Array.from(dlg.querySelectorAll('button')).find(b => (b.textContent || '').trim().toLowerCase() === 'ok');
  if (okBtn) okBtn.click();

  return JSON.stringify({ ok: true, reason: '' });
})()`
	return js, nil
}

// InsertSubscribeButton returns async JS that uses Substack's editor **Button** menu and selects
// **Subscribe** (native subscribe block). Best-effort: returns ok:true with a reason when the UI is not found.
func InsertSubscribeButton() (string, error) {
	return `(async function() {
  function sleep(ms) { return new Promise(r => setTimeout(r, ms)); }

  function visible(el) {
    if (!el) return false;
    const r = el.getBoundingClientRect();
    return r.width > 0 && r.height > 0;
  }

  function clickByTextExact(root, want) {
    const w = (want || '').trim().toLowerCase();
    if (!root || !w) return false;
    const els = root.querySelectorAll('button, [role="menuitem"], [role="option"], [role="button"], a, div');
    for (const el of els) {
      const t = (el.textContent || '').trim().toLowerCase();
      if (t === w) {
        el.click();
        return true;
      }
    }
    return false;
  }

  function editorRoot() {
    return document.querySelector('[data-testid="post-editor"]') ||
      document.querySelector('[data-testid="editor"]') ||
      document.querySelector('[class*="post-editor"]');
  }

  function findButtonMenuTrigger() {
    const root = editorRoot() || document.body;
    const trySelectors = [
      'button[data-testid="button-toolbar-dropdown"]',
      'button[data-testid*="button"]',
      '[role="combobox"][aria-label="Button"]',
      'button[aria-label="Button"]',
      'button[aria-label="Insert button"]',
      '[aria-label="Button menu"]',
    ];
    for (const sel of trySelectors) {
      try {
        const el = root.querySelector(sel);
        if (el && visible(el)) return el;
      } catch (e) {}
    }
    const candidates = Array.from(root.querySelectorAll('button, [role="button"], [role="combobox"]'));
    for (const b of candidates) {
      if (!visible(b)) continue;
      const t = (b.textContent || '').replace(/\s+/g, ' ').trim().toLowerCase();
      if (t === 'button' || t === 'button\u25be' || t.startsWith('button ')) return b;
      const al = (b.getAttribute('aria-label') || '').toLowerCase();
      if (al === 'button') return b;
      if (al.includes('insert') && al.includes('button') && al.length < 64) return b;
      if (al.includes('button') && al.length < 48 && !al.includes('more') && !al.includes('template')) return b;
    }
    return null;
  }

  await sleep(420);
  const trigger = findButtonMenuTrigger();
  if (!trigger) {
    return JSON.stringify({ ok: true, reason: 'Button toolbar control not found; skipped subscribe insert' });
  }
  try { trigger.scrollIntoView({ block: 'center' }); } catch (e) {}
  trigger.click();
  await sleep(400);

  let clicked = false;
  for (let attempt = 0; attempt < 18; attempt++) {
    const roots = document.querySelectorAll('[role="menu"], [role="listbox"], [data-radix-popper-content-wrapper]');
    for (const m of roots) {
      if (clickByTextExact(m, 'Subscribe')) {
        clicked = true;
        break;
      }
    }
    if (clicked) break;
    if (clickByTextExact(document.body, 'Subscribe')) {
      clicked = true;
      break;
    }
    await sleep(220);
  }
  if (!clicked) {
    return JSON.stringify({ ok: true, reason: 'Subscribe menu item not found' });
  }
  await sleep(200);
  return JSON.stringify({ ok: true, reason: '' });
})()`, nil
}
