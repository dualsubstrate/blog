package substackbrowser

import (
	"encoding/json"
	"fmt"
)

// ScheduleAfterContinueOptions is embedded as JSON into ScheduleAfterContinueJS.
type ScheduleAfterContinueOptions struct {
	Tags          []string `json:"tags"`
	SectionLabel  string   `json:"sectionLabel"`
	DateTimeLocal string   `json:"dateTimeLocal"`
	// DateDisplay is Substack's text schedule field (e.g. "29/04/2026, 08:40 am"), not ISO datetime-local.
	DateDisplay       string `json:"dateDisplay"`
	TickEmailSubstack bool   `json:"tickEmailSubstack"`
}

// ScheduleAfterContinueJS returns browser JS that clicks Continue, then fills the
// publish settings screen (section, optional email delivery off, schedule time then tags, delivery, confirm).
// Substack markup changes often; this uses broad heuristics and returns ok:false with a reason on miss.
func ScheduleAfterContinueJS(opt ScheduleAfterContinueOptions) (string, error) {
	raw, err := json.Marshal(opt)
	if err != nil {
		return "", fmt.Errorf("substackbrowser: marshal schedule options: %w", err)
	}
	js := `(async function(){
  const cfg = ` + string(raw) + `;
  const reasons = [];
  function sleep(ms) { return new Promise(r => setTimeout(r, ms)); }

  function visible(el) {
    if (!el) return false;
    const r = el.getBoundingClientRect();
    return r.width > 0 && r.height > 0;
  }

  function inputReachable(el) {
    if (!el || el.tagName !== 'INPUT') return false;
    try {
      if (document.activeElement === el) return true;
    } catch (e) {}
    const r = el.getBoundingClientRect();
    const st = window.getComputedStyle(el);
    if (st.display === 'none' || st.visibility === 'hidden') return false;
    return r.width > 0 && r.height > 0;
  }

  async function blurTagComboboxIntoModal() {
    const modal = publishModalRoot();
    if (!modal) return;
    const nodes = Array.from(modal.querySelectorAll('h2, h3, h4, label, div, span, p')).filter(visible);
    for (const el of nodes) {
      const raw = (el.innerText || '').trim();
      const low = raw.toLowerCase();
      if (low.length < 6 || low.length > 120) continue;
      if (low === 'add tags' || low.startsWith('add tags')) continue;
      if (low === 'cancel' || low.includes('close dialog')) continue;
      if (el.tagName === 'BUTTON' || (el.getAttribute('role') || '').toLowerCase() === 'button') continue;
      if (el.querySelector && el.querySelector('input, textarea, [role="combobox"]')) continue;
      if (low.includes('social preview') || low === 'audience' || low.includes('this post belongs')) {
        try { el.click(); await sleep(120); return; } catch (e) {}
      }
    }
    try {
      if (document.activeElement && document.activeElement !== document.body) {
        document.activeElement.blur();
      }
    } catch (e) {}
    await sleep(100);
  }

  async function focusHeadlessTagCombobox() {
    const modal = publishModalRoot();
    if (!modal) return;
    for (const inp of modal.querySelectorAll('input[role="combobox"]')) {
      if (!isHeadlessTagComboboxInput(inp)) continue;
      const host = inp.parentElement;
      if (!host || !visible(host)) continue;
      try {
        inp.focus();
        syntheticPointerClick(inp);
        return;
      } catch (e) {}
    }
  }

  async function waitForTagInputAgain(maxMs) {
    const end = Date.now() + maxMs;
    let n = 0;
    while (Date.now() < end) {
      let el = tagInput();
      if (el) return el;
      if (n % 2 === 0) {
        await focusHeadlessTagCombobox();
        await sleep(120);
        el = tagInput();
        if (el) return el;
      }
      if (n % 3 !== 2) await sleep(240);
      else await blurTagComboboxIntoModal();
      n++;
    }
    return tagInput();
  }

  function clickExactButton(text) {
    const want = text.trim().toLowerCase();
    const cands = Array.from(document.querySelectorAll('button,[role="button"],a')).filter(el => {
      if (!visible(el)) return false;
      const t = (el.innerText || '').trim().toLowerCase();
      return t === want;
    });
    if (cands.length === 0) return false;
    cands.sort((a, b) => {
      const ra = a.getBoundingClientRect();
      const rb = b.getBoundingClientRect();
      return (rb.width * rb.height) - (ra.width * ra.height);
    });
    try { cands[0].scrollIntoView({ block: 'center' }); } catch (e) {}
    try { cands[0].click(); return true; } catch (e) { reasons.push(String(e)); return false; }
  }


  function clickExactButtonAny(labels) {
    for (const text of labels) {
      if (clickExactButton(text)) return true;
    }
    return false;
  }

  function publishSettingsBodyVisible(t) {
    if (t.includes('this post belongs')) return true;
    if (t.includes('esta publicación pertenece') || t.includes('este artículo pertenece')) return true;
    const hasTags = t.includes('add tags') || t.includes('añadir etiquetas') || t.includes('agregar etiquetas');
    const hasAudience = t.includes('audience') || t.includes('audiencia');
    if (hasTags && hasAudience) return true;
    if (t.includes('delivery') || t.includes('entrega')) return true;
    if (t.includes('send via email') && t.includes('substack')) return true;
    if (t.includes('enviar por correo') && t.includes('substack')) return true;
    return false;
  }


  if (!clickExactButtonAny(['Continue', 'Continuar'])) {
    return JSON.stringify({ ok: false, reason: 'Continue button not found' });
  }

  async function waitForPublishSettings(timeoutMs) {
    const deadline = Date.now() + timeoutMs;
    while (Date.now() < deadline) {
      const t = (document.body && document.body.innerText || '').toLowerCase();
      if (publishSettingsBodyVisible(t)) return true;
      await sleep(200);
    }
    return false;
  }
  if (!(await waitForPublishSettings(120000))) {
    return JSON.stringify({ ok: false, reason: 'Publish settings screen did not appear after Continue (timeout 120s)' });
  }

  async function pickSection(label) {
    if (!label || !label.trim()) return true;
    const want = label.trim().toLowerCase();
    const selects = Array.from(document.querySelectorAll('select')).filter(visible);
    for (const sel of selects) {
      for (let i = 0; i < sel.options.length; i++) {
        const o = sel.options[i];
        const tx = (o.text || '').trim().toLowerCase();
        const vl = (o.value || '').trim().toLowerCase();
        if (tx === want || vl === want || tx.includes(want) || want.includes(tx)) {
          sel.selectedIndex = i;
          sel.dispatchEvent(new Event('change', { bubbles: true }));
          return true;
        }
      }
    }
    const triggers = Array.from(document.querySelectorAll('[role="combobox"],button')).filter(el => {
      if (!visible(el)) return false;
      const t = (el.innerText || '').trim().toLowerCase();
      return t.length > 0 && t.length < 120;
    });
    triggers.sort((a, b) => {
      const ra = a.getBoundingClientRect();
      const rb = b.getBoundingClientRect();
      return (rb.width * rb.height) - (ra.width * ra.height);
    });
    for (const tr of triggers.slice(0, 12)) {
      try { tr.scrollIntoView({ block: 'center' }); } catch (e) {}
      try { tr.click(); } catch (e) { continue; }
      const deadline = Date.now() + 3000;
      while (Date.now() < deadline) {
        const opts = Array.from(document.querySelectorAll('[role="option"],[role="menuitem"],li,button,div,span')).filter(o => {
          if (!visible(o)) return false;
          const x = (o.innerText || '').trim().toLowerCase();
          return x === want || x.includes(want) || want.includes(x);
        });
        if (opts.length > 0) {
          try { opts[0].click(); return true; } catch (e) {}
        }
        await sleep(80);
      }
    }
    reasons.push('section not matched: ' + label);
    return false;
  }

  function tagInputRoots() {
    const dlg = publishModalRoot();
    if (dlg) return [dlg];
    return [document];
  }

  function isHeadlessTagComboboxInput(inp) {
    if (!inp || inp.tagName !== 'INPUT') return false;
    if ((inp.getAttribute('role') || '').toLowerCase() !== 'combobox') return false;
    const par = inp.parentElement;
    if (!par) return false;
    const cls = String(par.className || '');
    if (cls.indexOf('inputBox') >= 0 || cls.indexOf('hasChips') >= 0) return true;
    return false;
  }

  function tagFieldInputUsable(inp) {
    if (!inp) return false;
    try {
      if (document.activeElement === inp) return true;
    } catch (e) {}
    if (isHeadlessTagComboboxInput(inp)) {
      const host = inp.parentElement;
      if (host && visible(host)) return true;
    }
    return visible(inp) || inputReachable(inp);
  }

  function tagInput() {
    const qs = [
      'input[placeholder*="tag" i]',
      'input[placeholder*="Tag" i]',
      'input[placeholder*="create tag" i]',
      'input[placeholder*="select or create" i]',
      'input[placeholder*="select" i]',
      'input[aria-label*="tag" i]',
      'input[name*="tag" i]',
      '[role="combobox"] input[type="text"]',
      '[role="combobox"] input[type="search"]',
      '[role="combobox"] input:not([type])',
    ];
    for (const root of tagInputRoots()) {
      if (!root || !root.querySelector) continue;
      for (const q of qs) {
        const el = root.querySelector(q);
        if (!el || !tagFieldInputUsable(el)) continue;
        const ty = (el.type || 'text').toLowerCase();
        if (el.tagName === 'INPUT' && ty !== 'checkbox' && ty !== 'radio' && ty !== 'hidden' && ty !== 'file') {
          return el;
        }
      }
      const fromPh = Array.from(root.querySelectorAll('input[type="text"],input[type="search"],input:not([type])')).filter(inp => {
        if (!tagFieldInputUsable(inp)) return false;
        const ty = (inp.type || 'text').toLowerCase();
        if (ty === 'checkbox' || ty === 'radio' || ty === 'hidden') return false;
        const ph = (inp.getAttribute('placeholder') || '').toLowerCase();
        const al = (inp.getAttribute('aria-label') || '').toLowerCase();
        return ph.includes('tag') || (ph.includes('select') && ph.includes('create')) || al.includes('tag');
      });
      if (fromPh.length > 0) return fromPh[0];
      for (const inp of root.querySelectorAll('input[role="combobox"]')) {
        if (!isHeadlessTagComboboxInput(inp)) continue;
        const ty = (inp.type || 'text').toLowerCase();
        if (ty === 'checkbox' || ty === 'radio' || ty === 'hidden' || ty === 'file') continue;
        if (tagFieldInputUsable(inp)) return inp;
      }
      const combos = Array.from(root.querySelectorAll('[role="combobox"]')).filter(el => visible(el) || (el.tagName === 'INPUT' && tagFieldInputUsable(el)));
      for (const box of combos) {
        if (box.tagName === 'INPUT') {
          const ty = (box.type || 'text').toLowerCase();
          if (ty === 'checkbox' || ty === 'radio' || ty === 'hidden' || ty === 'file') continue;
          if (isHeadlessTagComboboxInput(box) && tagFieldInputUsable(box)) return box;
          continue;
        }
        const inner = box.querySelector('input[type="text"],input[type="search"],input:not([type])');
        if (!inner || !tagFieldInputUsable(inner)) continue;
        const ph = (inner.getAttribute('placeholder') || '').toLowerCase();
        const al = (inner.getAttribute('aria-label') || '').toLowerCase();
        if (ph.includes('tag') || (ph.includes('select') && ph.includes('create')) || al.includes('tag')) {
          return inner;
        }
      }
    }
    return null;
  }

  function alnumKey(s) {
    return String(s || '').toLowerCase().replace(/[^a-z0-9]/g, '');
  }

  function looksLikeCreateTagLine(txLower, needleLower) {
    if (!txLower.includes('create')) return false;
    if (txLower.includes(needleLower)) return true;
    const a = alnumKey(txLower);
    const b = alnumKey(needleLower);
    return b.length > 1 && a.includes(b);
  }

  function cleanTagRowText(s) {
    return String(s || '').split('\n')[0].trim();
  }

  function syntheticPointerClick(el) {
    if (!el) return;
    const o = { bubbles: true, cancelable: true, view: window };
    try {
      el.dispatchEvent(new MouseEvent('mousedown', o));
      el.dispatchEvent(new MouseEvent('mouseup', o));
      el.dispatchEvent(new MouseEvent('click', o));
    } catch (e) {}
  }

  function tagListRowContextOk(el) {
    if (!el) return false;
    if (tagSuggestionElementOk(el)) return true;
    if (el.matches && el.matches('[role="option"]')) {
      const lb = el.closest('[role="listbox"]');
      if (lb && visible(lb)) {
        const lr = lb.getBoundingClientRect();
        if (lr.width > 4 && lr.height > 4) {
          const inp = tagInput();
          if (inp) {
            const ir = inp.getBoundingClientRect();
            if (lr.top < ir.bottom + 480 && lr.bottom > ir.top - 120 && lr.left < ir.right + 420 && lr.right > ir.left - 280) return true;
          }
          const modal = publishModalRoot();
          if (modal) {
            const mr = modal.getBoundingClientRect();
            if (lr.top < mr.bottom + 480 && lr.bottom > mr.top - 80 && lr.left < mr.right + 240 && lr.right > mr.left - 240) return true;
          }
        }
      }
    }
    const list = el.closest('[cmdk-list],[cmdk-root]');
    if (!list) return false;
    const r = list.getBoundingClientRect();
    if (r.width < 4 || r.height < 4) return false;
    const re = el.getBoundingClientRect();
    if (re.width < 2 || re.height < 2) return false;
    const inp = tagInput();
    if (inp) {
      const ir = inp.getBoundingClientRect();
      const nearInp = r.top < ir.bottom + 480 && r.bottom > ir.top - 120 && r.left < ir.right + 420 && r.right > ir.left - 280;
      if (nearInp) return true;
    }
    const modal = publishModalRoot();
    if (modal) {
      const mr = modal.getBoundingClientRect();
      const nearModal = r.top < mr.bottom + 480 && r.bottom > mr.top - 80 && r.left < mr.right + 240 && r.right > mr.left - 240;
      if (nearModal) return true;
    }
    return false;
  }

  function listOptionShowsAlreadyApplied(el) {
    if (!el || !el.querySelector) return false;
    if (el.querySelector('svg.lucide-check')) return true;
    if (el.querySelector('svg path[d*="M20 6"]')) return true;
    const t = el.innerText || '';
    if (t.indexOf('\u2713') >= 0 || t.indexOf('\u2714') >= 0) return true;
    return false;
  }

  function pressTagComboboxEnter(inp) {
    if (!inp) return;
    try { inp.focus(); } catch (e) {}
    const opts = { bubbles: true, cancelable: true, key: 'Enter', code: 'Enter', keyCode: 13, which: 13 };
    try {
      inp.dispatchEvent(new KeyboardEvent('keydown', opts));
      inp.dispatchEvent(new KeyboardEvent('keypress', opts));
      inp.dispatchEvent(new KeyboardEvent('keyup', opts));
    } catch (e) {}
  }

  function rowMatchesTagNeedle(line, needle, want) {
    if (!line) return false;
    if (looksLikeCreateTagLine(line, needle) || line.startsWith('create')) return true;
    return line === needle || alnumKey(line) === want;
  }

  function exactTagOptionVisible(tagText) {
    const needle = String(tagText || '').trim().toLowerCase();
    const want = alnumKey(tagText);
    if (!needle || want.length < 2) return false;
    const rows = document.querySelectorAll('[cmdk-item], [data-radix-collection-item], [role="listbox"] [role="option"]');
    for (const el of rows) {
      if (!visible(el) || !tagListRowContextOk(el)) continue;
      const line = cleanTagRowText(el.innerText || '').toLowerCase();
      if (!line || line.startsWith('create')) continue;
      if (line === needle || alnumKey(line) === want) return true;
    }
    return false;
  }

  async function commitOneTagFromFullQuery(tagText, inp) {
    const needle = tagText.trim().toLowerCase();
    const want = alnumKey(tagText);
    if (!needle || want.length < 2) return false;
    let createOnlyWaits = 0;
    for (let attempt = 0; attempt < 14; attempt++) {
      if (hasCommittedTagPill(tagText)) return true;
      const rows = [];
      const seen = new Set();
      for (const el of document.querySelectorAll('[cmdk-item], [data-radix-collection-item], [role="listbox"] [role="option"]')) {
        if (seen.has(el) || !visible(el) || !tagListRowContextOk(el)) continue;
        const raw = (el.innerText || '').trim();
        const line = cleanTagRowText(raw).toLowerCase();
        if (!line) continue;
        seen.add(el);
        rows.push({ el, line, tx: raw.toLowerCase() });
      }
      let exactHit = null;
      let createHit = null;
      for (const r of rows) {
        if (looksLikeCreateTagLine(r.line, needle) || r.line.startsWith('create')) {
          if (!createHit) createHit = r;
          continue;
        }
        if (r.line === needle || alnumKey(r.line) === want) {
          if (!exactHit || r.line.length < exactHit.line.length) exactHit = r;
        }
      }
      if (!exactHit && createHit && createOnlyWaits < 5) {
        createOnlyWaits++;
        await sleep(220);
        continue;
      }
      if (exactHit) {
        const el = exactHit.el;
        // List checkmark is not a chip: only succeed when hasCommittedTagPill is true.
        if (listOptionShowsAlreadyApplied(el)) {
          await sleep(200);
          if (hasCommittedTagPill(tagText)) return true;
        }
        if (el.matches && el.matches('[role="option"]') && el.getAttribute('aria-selected') === 'true') {
          await sleep(160);
          if (hasCommittedTagPill(tagText)) return true;
        }
        if (hasCommittedTagPill(tagText)) return true;
        try { el.scrollIntoView({ block: 'center' }); el.click(); } catch (e) {}
        syntheticPointerClick(el);
        await sleep(240);
        if (hasCommittedTagPill(tagText)) return true;
        pressTagComboboxEnter(inp);
        await sleep(200);
        if (hasCommittedTagPill(tagText)) return true;
        await sleep(90);
        continue;
      }
      if (createHit) {
        const el = createHit.el;
        try { el.scrollIntoView({ block: 'center' }); el.click(); } catch (e) {}
        syntheticPointerClick(el);
        await sleep(240);
        if (hasCommittedTagPill(tagText)) return true;
        pressTagComboboxEnter(inp);
        await sleep(200);
        if (hasCommittedTagPill(tagText)) return true;
        await sleep(90);
        continue;
      }
      await sleep(100);
    }
    return hasCommittedTagPill(tagText);
  }

  function clickTagListRowForNeedle(tagText) {
    const needle = tagText.trim().toLowerCase();
    const want = alnumKey(tagText);
    if (!needle || want.length < 2) return false;
    if (hasCommittedTagPill(tagText)) return false;
    const lists = Array.from(document.querySelectorAll(
      '[cmdk-list], [cmdk-root], [data-radix-popper-content-wrapper], [data-radix-portal], [data-radix-menu-content], [role="listbox"]'
    )).filter(el => {
      if (!el || !el.getBoundingClientRect) return false;
      const r = el.getBoundingClientRect();
      return r.width > 4 && r.height > 4;
    });
    const roots = lists.length > 0 ? lists : [document.body];
    for (const root of roots) {
      const rows = Array.from(root.querySelectorAll('[cmdk-item], [data-radix-collection-item], [role="option"]')).filter(el => {
        if (!el.getBoundingClientRect) return false;
        const r = el.getBoundingClientRect();
        return r.width > 2 && r.height > 2;
      });
      const hits = [];
      for (const el of rows) {
        if (!tagListRowContextOk(el)) continue;
        const line = cleanTagRowText(el.innerText || '').toLowerCase();
        if (!line || line.startsWith('create')) continue;
        if (line === needle || alnumKey(line) === want) hits.push({ el, len: line.length });
      }
      hits.sort((a, b) => a.len - b.len);
      if (hits.length > 0) {
        const target = hits[0].el;
        // aria-selected highlight is not a chip; still click when no pill.
        try {
          target.scrollIntoView({ block: 'center' });
          target.click();
        } catch (e) {}
        syntheticPointerClick(target);
        const inner = target.querySelector('span,div');
        if (inner && inner !== target) syntheticPointerClick(inner);
        return true;
      }
    }
    return false;
  }

  function highlightedTagListRow(preferTagText) {
    const prefer = document.querySelector(
      '[cmdk-item][data-selected], [cmdk-item][aria-selected="true"], [cmdk-item][data-highlighted]'
    );
    if (prefer && tagListRowContextOk(prefer)) return prefer;
    let needle = String(preferTagText || '').trim().toLowerCase();
    let want = alnumKey(preferTagText);
    if (!needle) {
      const inp = tagInput();
      if (inp) {
        needle = String(inp.value || '').trim().toLowerCase();
        want = alnumKey(needle);
      }
    }
    const items = Array.from(document.querySelectorAll('[cmdk-item]')).filter(el => tagListRowContextOk(el));
    if (items.length > 0) {
      if (needle && want.length >= 2) {
        for (const el of items) {
          const line = cleanTagRowText(el.innerText || '').toLowerCase();
          if (line && !line.startsWith('create') && (line === needle || alnumKey(line) === want)) return el;
        }
      }
      return items[0];
    }
    const opts = Array.from(document.querySelectorAll('[role="listbox"] [role="option"]')).filter(visible);
    let firstOk = null;
    let exactOk = null;
    for (const o of opts) {
      if (!tagListRowContextOk(o)) continue;
      if (!firstOk) firstOk = o;
      const line = cleanTagRowText(o.innerText || '').toLowerCase();
      if (needle && want.length >= 2 && line && !line.startsWith('create') && (line === needle || alnumKey(line) === want)) {
        exactOk = o;
      }
      const aria = (o.getAttribute('aria-selected') || '').toLowerCase();
      const st = (o.getAttribute('data-headlessui-state') || '').toLowerCase();
      const cls = String(o.className || '').toLowerCase();
      if (aria === 'true' || st.includes('active') || st.includes('selected') || cls.includes('active') || cls.includes('highlight')) {
        return o;
      }
    }
    if (exactOk) return exactOk;
    return firstOk;
  }

  function highlightedTagRowMatchesNeedle(tagText) {
    const hi = highlightedTagListRow(tagText);
    if (!hi) return false;
    const needle = tagText.trim().toLowerCase();
    const want = alnumKey(tagText);
    const line = cleanTagRowText(hi.innerText || '').toLowerCase();
    return rowMatchesTagNeedle(line, needle, want);
  }

  async function navigateTagListWithArrows(tagText, inp) {
    const want = alnumKey(tagText);
    const needle = tagText.trim().toLowerCase();
    if (!want || want.length < 2) return false;
    if (!tagQueryLooksApplied(inp, tagText)) return false;
    for (let step = 0; step < 48; step++) {
      if (hasCommittedTagPill(tagText)) return true;
      if (!tagQueryLooksApplied(inp, tagText)) return false;
      const hi = highlightedTagListRow(tagText);
      if (hi) {
        const line = cleanTagRowText(hi.innerText || '').toLowerCase();
        if (line && (looksLikeCreateTagLine(line, needle) || (!line.startsWith('create') && (line === needle || alnumKey(line) === want)))) {
          if (hasCommittedTagPill(tagText)) return true;
          try { hi.scrollIntoView({ block: 'center' }); hi.click(); } catch (e) {}
          syntheticPointerClick(hi);
          await sleep(260);
          if (hasCommittedTagPill(tagText)) return true;
          pressTagComboboxEnter(inp);
          await sleep(200);
          if (hasCommittedTagPill(tagText)) return true;
        }
      }
      try { inp.focus(); } catch (e) {}
      inp.dispatchEvent(new KeyboardEvent('keydown', { bubbles: true, key: 'ArrowDown', code: 'ArrowDown', key: 'ArrowDown', keyCode: 40 }));
      inp.dispatchEvent(new KeyboardEvent('keyup', { bubbles: true, key: 'ArrowDown', code: 'ArrowDown', key: 'ArrowDown', keyCode: 40 }));
      await sleep(85);
    }
    if (highlightedTagRowMatchesNeedle(tagText) || exactTagOptionVisible(tagText)) {
      pressTagComboboxEnter(inp);
      await sleep(240);
    }
    return hasCommittedTagPill(tagText);
  }

  function isLikelyTagChipButton(el) {
    if (!el || el.tagName !== 'BUTTON') return false;
    if ((el.getAttribute('aria-haspopup') || '').toLowerCase() === 'listbox') return false;
    const id = String(el.id || '');
    if (id.indexOf('headlessui-combobox-button') >= 0) return false;
    if (el.closest('[role="listbox"]')) return false;
    const cls = String(el.className || '');
    if (cls.indexOf('chip') >= 0) return true;
    const par = el.parentElement;
    if (par) {
      const pc = String(par.className || '');
      if (pc.indexOf('hasChips') >= 0 || pc.indexOf('inputBox') >= 0) return true;
    }
    return false;
  }

  function tagComboboxChipHost(inp) {
    if (!inp) return null;
    const p = inp.parentElement;
    if (!p || !p.querySelector) return null;
    if (p.querySelector('input[role="combobox"]') !== inp) return null;
    return p;
  }

  function hasCommittedTagPill(tagText) {
    const want = alnumKey(tagText);
    if (!want || want.length < 2) return false;
    function pillInChipHost(host) {
      if (!host) return false;
      const btns = Array.from(host.querySelectorAll('button')).filter(visible);
      for (const el of btns) {
        if (!isLikelyTagChipButton(el)) continue;
        const raw = (el.innerText || '').trim();
        if (raw.length < 2 || raw.length > 80) continue;
        const low = raw.toLowerCase();
        if (low.includes('select or create')) continue;
        if (alnumKey(raw) === want) return true;
      }
      return false;
    }
    const inp = tagInput();
    if (inp) {
      const host = tagComboboxChipHost(inp);
      if (pillInChipHost(host)) return true;
      const combo = inp.closest('[role="combobox"]');
      if (combo && combo !== inp) {
        const inner = Array.from(combo.querySelectorAll('span, div, button, a')).filter(visible);
        for (const el of inner) {
          if (el === inp || el.contains(inp)) continue;
          const raw = (el.innerText || '').trim();
          if (raw.length < 2 || raw.length > 72) continue;
          const low = raw.toLowerCase();
          if (low.includes('select or create')) continue;
          if (alnumKey(raw) === want) return true;
        }
      }
      const modal = publishModalRoot();
      const shell = modal || inp.closest('form');
      const roots = [];
      if (modal) roots.push(modal);
      if (inp.parentElement) roots.push(inp.parentElement);
      if (shell && shell !== modal) roots.push(shell);
      const seen = new Set();
      const noise = new Set(['delivery', 'scheduling', 'audience', 'publish', 'edit', 'cancel', 'save', 'continue', 'schedule']);
      for (const root of roots) {
        if (!root || seen.has(root)) continue;
        seen.add(root);
        const chips = Array.from(root.querySelectorAll(
          '[aria-label*="Remove" i], [aria-label*="remove tag" i], button[title*="emove" i], [data-testid*="tag"], button, [role="button"], a, span'
        )).filter(visible);
        for (const el of chips) {
          if (el === inp || el.contains(inp) || inp.contains(el)) continue;
          const raw = (el.innerText || '').trim();
          if (!raw || raw.length > 72) continue;
          const low = raw.toLowerCase();
          if (low.includes('select or create') || low.includes('add tags')) continue;
          if (low.startsWith('create')) continue;
          if (noise.has(low)) continue;
          if (alnumKey(raw) === want) return true;
        }
      }
      return false;
    }
    const modalOnly = publishModalRoot();
    if (!modalOnly) return false;
    for (const inp2 of modalOnly.querySelectorAll('input[role="combobox"]')) {
      if (!isHeadlessTagComboboxInput(inp2)) continue;
      const h = tagComboboxChipHost(inp2);
      if (pillInChipHost(h)) return true;
    }
    return false;
  }

  function tagSuggestionElementOk(el) {
    if (!el) return false;
    const modal = publishModalRoot();
    if (!modal) return true;
    if (modal.contains(el)) return true;
    const portal = el.closest(
      '[data-radix-popper-content-wrapper],[data-radix-portal],[data-radix-menu-content],[role="listbox"],[cmdk-root],[cmdk-list]'
    );
    if (portal && visible(portal)) return true;
    const list = el.closest('[cmdk-list],[cmdk-root]');
    if (list) {
      const lr = list.getBoundingClientRect();
      if (lr.width > 4 && lr.height > 4) {
        const inp = tagInput();
        if (inp) {
          const ir = inp.getBoundingClientRect();
          if (lr.top < ir.bottom + 480 && lr.bottom > ir.top - 120 && lr.left < ir.right + 420 && lr.right > ir.left - 280) return true;
        }
        const modal = publishModalRoot();
        if (modal) {
          const mr = modal.getBoundingClientRect();
          if (lr.top < mr.bottom + 480 && lr.bottom > mr.top - 80 && lr.left < mr.right + 240 && lr.right > mr.left - 240) return true;
        }
      }
    }
    const r = el.getBoundingClientRect();
    const r2 = modal.getBoundingClientRect();
    const cx = r.left + r.width / 2;
    const cy = r.top + r.height / 2;
    return cx >= r2.left && cx <= r2.right && cy >= r2.top && cy <= r2.bottom;
  }

  async function clickCreateTagSuggestion(tagText) {
    const needle = tagText.trim().toLowerCase();
    if (!needle) return false;
    const deadline = Date.now() + 5500;
    while (Date.now() < deadline) {
      if (hasCommittedTagPill(tagText)) return true;
      if (clickTagListRowForNeedle(tagText)) {
        await sleep(200);
        if (hasCommittedTagPill(tagText)) return true;
      }
      const hi = document.querySelector(
        '[cmdk-item][data-selected="true"], [cmdk-item][data-selected], [cmdk-item][aria-selected="true"], [cmdk-item][data-highlighted], [data-highlighted][role="option"]'
      );
      if (hi && tagListRowContextOk(hi)) {
        const line = cleanTagRowText(hi.innerText || '').toLowerCase();
        if (line && !line.startsWith('create') && (line === needle || alnumKey(line) === alnumKey(tagText))) {
          try { hi.scrollIntoView({ block: 'center' }); hi.click(); } catch (e) {}
          syntheticPointerClick(hi);
          await sleep(200);
          if (hasCommittedTagPill(tagText)) return true;
        }
      }
      const sel = '[role="option"],[role="menuitem"],[data-radix-collection-item],[cmdk-item],button,div[role="button"],li,span';
      const cands = Array.from(document.querySelectorAll(sel)).filter(el => {
        if (!visible(el)) return false;
        return tagSuggestionElementOk(el) || tagListRowContextOk(el);
      });
      for (const el of cands) {
        const raw = (el.innerText || '').trim();
        const tx = raw.toLowerCase();
        const role = (el.getAttribute('role') || '').toLowerCase();
        const isCmdk = el.matches && el.matches('[cmdk-item], [data-radix-collection-item]');
        if (looksLikeCreateTagLine(tx, needle)) {
          try { el.scrollIntoView({ block: 'center' }); el.click(); } catch (e) {}
          syntheticPointerClick(el);
          await sleep(200);
          if (hasCommittedTagPill(tagText)) return true;
        }
        if (isCmdk && raw.length > 0 && raw.length < 120) {
          const line = cleanTagRowText(raw).toLowerCase();
          if (!line.startsWith('create') && (line === needle || alnumKey(line) === alnumKey(tagText))) {
            try { el.scrollIntoView({ block: 'center' }); el.click(); } catch (e) {}
            syntheticPointerClick(el);
            await sleep(200);
            if (hasCommittedTagPill(tagText)) return true;
          }
        }
        if ((role === 'option' || role === 'menuitem') && (tx === needle || tx.replace(/^#/, '') === needle)) {
          try { el.scrollIntoView({ block: 'center' }); el.click(); } catch (e) {}
          syntheticPointerClick(el);
          await sleep(200);
          if (hasCommittedTagPill(tagText)) return true;
        }
        if ((role === 'option' || role === 'menuitem') && alnumKey(raw) === alnumKey(tagText)) {
          try { el.scrollIntoView({ block: 'center' }); el.click(); } catch (e) {}
          syntheticPointerClick(el);
          await sleep(200);
          if (hasCommittedTagPill(tagText)) return true;
        }
      }
      await sleep(100);
    }
    return false;
  }

  function setNativeInputValue(inp, val) {
    if (!inp) return;
    const value = val == null ? '' : String(val);
    try {
      const proto = window.HTMLInputElement.prototype;
      const desc = Object.getOwnPropertyDescriptor(proto, 'value');
      const own = Object.getOwnPropertyDescriptor(inp, 'value');
      if (desc && desc.set && (!own || own.set !== desc.set)) {
        desc.set.call(inp, value);
      } else if (own && own.set) {
        own.set.call(inp, value);
      } else {
        inp.value = value;
      }
    } catch (e) {
      try { inp.value = value; } catch (e2) {}
    }
  }

  function dispatchTagInputEvent(inp, data, inputType) {
    if (!inp) return;
    try {
      inp.dispatchEvent(new InputEvent('input', { bubbles: true, cancelable: true, data: data, inputType: inputType || 'insertText' }));
    } catch (e) {
      inp.dispatchEvent(new Event('input', { bubbles: true }));
    }
  }

  function tagQueryLooksApplied(inp, tagText) {
    if (!inp) return false;
    const want = String(tagText || '').trim().toLowerCase();
    if (!want) return false;
    const got = String(inp.value || '').trim().toLowerCase();
    return got === want || got.replace(/\s+/g, '') === want.replace(/\s+/g, '');
  }

  function tagSuggestionReady(tagText) {
    const needle = String(tagText || '').trim().toLowerCase();
    const want = alnumKey(tagText);
    if (!needle || want.length < 2) return false;
    const rows = document.querySelectorAll('[cmdk-item], [data-radix-collection-item], [role="listbox"] [role="option"]');
    for (const el of rows) {
      if (!visible(el) || !tagListRowContextOk(el)) continue;
      const line = cleanTagRowText(el.innerText || '').toLowerCase();
      if (!line) continue;
      if (looksLikeCreateTagLine(line, needle) || line.startsWith('create')) return true;
      if (line === needle || alnumKey(line) === want) return true;
    }
    return false;
  }

  async function clearTagComboboxQuery(inp) {
    if (!inp) return;
    try { inp.focus(); } catch (e) {}
    await sleep(30);
    try { inp.select(); } catch (e) {}
    setNativeInputValue(inp, '');
    dispatchTagInputEvent(inp, '', 'deleteContentBackward');
    try {
      document.execCommand('selectAll', false, null);
      document.execCommand('delete', false, null);
    } catch (e) {}
    try {
      inp.dispatchEvent(new KeyboardEvent('keydown', { bubbles: true, key: 'Backspace', code: 'Backspace', keyCode: 8 }));
      inp.dispatchEvent(new KeyboardEvent('keyup', { bubbles: true, key: 'Backspace', code: 'Backspace', keyCode: 8 }));
    } catch (e) {}
    await sleep(40);
  }

  async function typeTagQuery(inp, tagText) {
    const text = String(tagText || '').trim();
    if (!inp || !text) return false;
    try { inp.scrollIntoView({ block: 'center' }); } catch (e) {}
    await clearTagComboboxQuery(inp);
    try { inp.focus(); } catch (e) {}
    await sleep(40);
    let inserted = false;
    try {
      inserted = !!document.execCommand('insertText', false, text);
    } catch (e) {
      inserted = false;
    }
    if (inserted && tagQueryLooksApplied(inp, text)) {
      try { inp.dispatchEvent(new Event('change', { bubbles: true })); } catch (e) {}
      return true;
    }
    setNativeInputValue(inp, text);
    dispatchTagInputEvent(inp, text, 'insertText');
    try { inp.dispatchEvent(new Event('change', { bubbles: true })); } catch (e) {}
    if (tagQueryLooksApplied(inp, text)) return true;
    await clearTagComboboxQuery(inp);
    try { inp.focus(); } catch (e) {}
    let built = '';
    for (let i = 0; i < text.length; i++) {
      const ch = text.charAt(i);
      built += ch;
      try {
        inp.dispatchEvent(new KeyboardEvent('keydown', { bubbles: true, key: ch, code: 'Key' + ch.toUpperCase(), keyCode: ch.charCodeAt(0) }));
      } catch (e) {}
      let charOk = false;
      try {
        charOk = !!document.execCommand('insertText', false, ch);
      } catch (e) {
        charOk = false;
      }
      if (!charOk || String(inp.value || '') !== built) {
        setNativeInputValue(inp, built);
        dispatchTagInputEvent(inp, ch, 'insertText');
      }
      try {
        inp.dispatchEvent(new KeyboardEvent('keyup', { bubbles: true, key: ch, code: 'Key' + ch.toUpperCase(), keyCode: ch.charCodeAt(0) }));
      } catch (e) {}
      await sleep(18);
    }
    try { inp.dispatchEvent(new Event('change', { bubbles: true })); } catch (e) {}
    return tagQueryLooksApplied(inp, text);
  }

  async function waitForTagSuggestions(tagText, inp, maxMs) {
    const deadline = Date.now() + (maxMs || 2500);
    while (Date.now() < deadline) {
      if (hasCommittedTagPill(tagText)) return true;
      if (tagSuggestionReady(tagText)) return true;
      if (inp && !tagQueryLooksApplied(inp, tagText)) return false;
      await sleep(80);
    }
    return tagSuggestionReady(tagText) || hasCommittedTagPill(tagText);
  }

  async function addTags(tags) {
    if (!tags || !tags.length) return true;
    const waitInpUntil = Date.now() + 22000;
    while (Date.now() < waitInpUntil && !tagInput()) {
      await sleep(200);
    }
    if (!tagInput()) {
      reasons.push('tag input not found (waited for Add tags field in Publish modal)');
      return false;
    }
    for (const tag of tags) {
      const t = String(tag || '').trim();
      if (!t) continue;
      let inp = tagInput();
      if (!inp) {
        inp = await waitForTagInputAgain(16000);
      }
      if (!inp) {
        reasons.push('tag input disappeared mid-run (could not recover after blur and wait)');
        return false;
      }
      if (hasCommittedTagPill(t)) {
        await sleep(40);
        continue;
      }
      let typed = await typeTagQuery(inp, t);
      if (!typed) {
        await sleep(120);
        typed = await typeTagQuery(inp, t);
      }
      if (!typed) {
        reasons.push('tag query did not stick in combobox for: ' + t);
        return false;
      }
      await waitForTagSuggestions(t, inp, 2800);
      let added = await commitOneTagFromFullQuery(t, inp);
      if (!added) {
        added = await clickCreateTagSuggestion(t);
      }
      if (!added && tagQueryLooksApplied(inp, t)) {
        added = await navigateTagListWithArrows(t, inp);
      }
      // Enter only when an exact matching option is visible (or highlighted for that needle), not any open list.
      if (!added && tagQueryLooksApplied(inp, t) && (highlightedTagRowMatchesNeedle(t) || exactTagOptionVisible(t))) {
        pressTagComboboxEnter(inp);
        await sleep(220);
        added = hasCommittedTagPill(t);
      }
      await sleep(280);
      if (!added && hasCommittedTagPill(t)) added = true;
      if (!added) {
        await sleep(400);
        if (hasCommittedTagPill(t)) added = true;
      }
      if (!added) {
        for (let poll = 0; poll < 20; poll++) {
          await sleep(200);
          if (hasCommittedTagPill(t)) {
            added = true;
            break;
          }
          if (poll === 10) {
            try { await blurTagComboboxIntoModal(); } catch (e) {}
            await sleep(180);
            if (hasCommittedTagPill(t)) {
              added = true;
              break;
            }
          }
        }
      }
      if (!added) {
        reasons.push('tag create UI not confirmed for: ' + t);
        return false;
      }
      await sleep(120);
    }
    return true;
  }

  function publishModalRoot() {
    const dialogs = Array.from(document.querySelectorAll('[role="dialog"]')).filter(visible);
    let best = null;
    let bestScore = -1;
    for (const d of dialogs) {
      const txt = (d.innerText || '').toLowerCase();
      const looksPublish =
        txt.includes('add tags') ||
        txt.includes('this post belongs') ||
        txt.includes('delivery') ||
        (txt.includes('scheduling') && txt.includes('email'));
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

  function scopeForPublishUI() {
    return publishModalRoot() || document;
  }

  function isControlChecked(cb) {
    if (!cb) return false;
    if (cb.type === 'checkbox') return !!cb.checked;
    const role = (cb.getAttribute && cb.getAttribute('role')) || '';
    if (role === 'checkbox' || role === 'switch') {
      return cb.getAttribute('aria-checked') === 'true';
    }
    return false;
  }

  function rowBlobAround(el) {
    let blob = ((el.getAttribute && el.getAttribute('aria-label')) || '') + ' ';
    let p = el;
    for (let d = 0; d < 10 && p; d++, p = p.parentElement) {
      blob += (p.innerText || '').slice(0, 400) + ' ';
    }
    return blob.toLowerCase();
  }

  // shallowRowContext limits text to the control's label and a few ancestors (narrow slices).
  // rowBlobAround walks the whole publish modal, so "email" + "substack" falsely match the schedule checkbox.
  function shallowRowContext(cb) {
    if (!cb) return '';
    let t = ((cb.getAttribute && cb.getAttribute('aria-label')) || '') + ' ';
    const lab = cb.closest('label');
    if (lab) t += (lab.innerText || '').slice(0, 240) + ' ';
    let p = cb.parentElement;
    for (let d = 0; d < 5 && p; d++, p = p.parentElement) {
      t += (p.innerText || '').slice(0, 240) + ' ';
    }
    return t.toLowerCase();
  }

  function isSchedulePublishRowControl(cb) {
    if (!cb) return false;
    if (cb.getAttribute && cb.getAttribute('data-track-input') === 'schedule_post') return true;
    if (cb.getAttribute && cb.getAttribute('data-testid') === 'scheduled-at') return true;
    const s = shallowRowContext(cb);
    if (s.includes('schedule time to email and publish')) return true;
    if (s.includes('schedule time to email')) return true;
    return s.includes('schedule time to publish') || (s.includes('scheduling') && s.includes('schedule time'));
  }

  function isSendViaEmailDeliveryRowControl(cb) {
    if (!cb) return false;
    if (isSchedulePublishRowControl(cb)) return false;
    const s = shallowRowContext(cb);
    if (!s.includes('email')) return false;
    return s.includes('send via email') || (s.includes('substack') && (s.includes('app') || s.includes('subscribers')));
  }

  function querySendEmailTrackControl(scope) {
    const tryOne = (s) => {
      if (!s || !s.querySelectorAll) return null;
      const hits = Array.from(s.querySelectorAll('[data-track-input="send_email"]'));
      const vis = hits.filter(visible);
      if (vis.length > 0) return vis[0];
      if (hits.length > 0) return hits[0];
      return null;
    };
    let el = tryOne(scope);
    if (!el && scope !== document) el = tryOne(document);
    return el;
  }

  function querySchedulePostTrackControl(scope) {
    const tryOne = (s) => {
      if (!s || !s.querySelectorAll) return null;
      const hits = Array.from(s.querySelectorAll('[data-track-input="schedule_post"]'));
      const vis = hits.filter(visible);
      if (vis.length > 0) return vis[0];
      if (hits.length > 0) return hits[0];
      return null;
    };
    let el = tryOne(scope);
    if (!el && scope !== document) el = tryOne(document);
    return el;
  }

  function tryClickSendEmailTrackButton(scope) {
    const el = querySendEmailTrackControl(scope);
    if (!el) return false;
    if (isControlChecked(el)) return true;
    try { el.scrollIntoView({ block: 'center' }); el.click(); return true; } catch (e) { return false; }
  }

  function emailDeliveryAlreadyOn(scope) {
    const tracked = querySendEmailTrackControl(scope);
    if (tracked && isControlChecked(tracked)) return true;
    const boxes = Array.from(scope.querySelectorAll('input[type="checkbox"],[role="checkbox"]')).filter(visible);
    for (const cb of boxes) {
      if (!isSendViaEmailDeliveryRowControl(cb)) continue;
      if (isControlChecked(cb)) return true;
    }
    return false;
  }

  function clickCheckboxNearLabelSubstring(needle, root) {
    const n = needle.toLowerCase();
    if (!n) return false;
    const scope = root || document;
    const labels = Array.from(scope.querySelectorAll('label,div,span,p')).filter(el => {
      if (!visible(el)) return false;
      const t = (el.innerText || '').toLowerCase();
      return t.includes(n) && t.length < 500;
    });
    for (const lab of labels) {
      let rootLabel = lab.closest('label');
      if (!rootLabel) rootLabel = lab;
      let cb = rootLabel.querySelector && rootLabel.querySelector('input[type="checkbox"],[role="checkbox"]');
      if (!cb && lab.parentElement) {
        cb = lab.parentElement.querySelector('input[type="checkbox"],[role="checkbox"]');
      }
      if (cb && visible(cb)) {
        if (isControlChecked(cb)) return true;
        try { cb.scrollIntoView({ block: 'center' }); cb.click(); return true; } catch (e) {}
      }
    }
    const boxes = Array.from(scope.querySelectorAll('input[type="checkbox"],[role="checkbox"]')).filter(visible);
    for (const cb of boxes) {
      const lab = cb.closest('label') || (cb.id ? document.querySelector('label[for="' + cb.id + '"]') : null);
      const t = ((lab && lab.innerText) || (cb.getAttribute && cb.getAttribute('aria-label')) || '').toLowerCase();
      if (!t.includes(n)) continue;
      if (isControlChecked(cb)) return true;
      try { cb.scrollIntoView({ block: 'center' }); cb.click(); return true; } catch (e) {}
    }
    return false;
  }

  function clickDeliveryRowInModal(scope) {
    const wantA = 'send via email';
    const wantB = 'substack';
    const candidates = Array.from(scope.querySelectorAll('label, button, div, span')).filter(el => {
      if (!visible(el)) return false;
      const t = (el.innerText || '').trim().toLowerCase();
      if (t.length < 8 || t.length > 200) return false;
      return t.includes(wantA) && t.includes(wantB);
    });
    candidates.sort((a, b) => (a.innerText || '').length - (b.innerText || '').length);
    for (const el of candidates) {
      const cb = el.querySelector && el.querySelector('input[type="checkbox"],[role="checkbox"]');
      if (cb && visible(cb)) {
        if (isControlChecked(cb)) return true;
        try { cb.scrollIntoView({ block: 'center' }); cb.click(); return true; } catch (e) {}
      }
      try { el.scrollIntoView({ block: 'center' }); el.click(); return true; } catch (e) {}
    }
    return false;
  }

  function tryTickSubstackEmailDelivery() {
    const scope = scopeForPublishUI();
    if (tryClickSendEmailTrackButton(scope)) return true;
    if (emailDeliveryAlreadyOn(scope)) return true;
    const needles = [
      'send via email and the substack app',
      'send via email',
      'email and the substack app',
      'the substack app',
      'substack app',
      'notify subscribers',
      'send to your subscribers',
      'email your subscribers',
      'send as an email',
      'deliver this post',
    ];
    for (const needle of needles) {
      if (clickCheckboxNearLabelSubstring(needle, scope)) return true;
    }
    const boxes = Array.from(scope.querySelectorAll('input[type="checkbox"],[role="checkbox"]')).filter(visible);
    for (const cb of boxes) {
      if (isControlChecked(cb)) continue;
      if (!isSendViaEmailDeliveryRowControl(cb)) continue;
      try {
        cb.scrollIntoView({ block: 'center' });
        cb.click();
        return true;
      } catch (e) {}
    }
    if (clickDeliveryRowInModal(scope)) return true;
    return false;
  }

  async function tryUntickSubstackEmailDelivery() {
    const scope = scopeForPublishUI();
    for (let round = 0; round < 8; round++) {
      if (!emailDeliveryAlreadyOn(scope)) return;
      const el = querySendEmailTrackControl(scope);
      if (el && isControlChecked(el)) {
        try { el.scrollIntoView({ block: 'center' }); } catch (e) {}
        try { el.focus(); } catch (e) {}
        try { el.click(); } catch (e) {}
        const lab = el.closest('label');
        if (lab) {
          try { lab.click(); } catch (e) {}
        }
        try {
          el.dispatchEvent(new MouseEvent('mousedown', { bubbles: true, cancelable: true, view: window }));
          el.dispatchEvent(new MouseEvent('mouseup', { bubbles: true, cancelable: true, view: window }));
          el.dispatchEvent(new MouseEvent('click', { bubbles: true, cancelable: true, view: window }));
        } catch (e) {}
        await sleep(240);
        continue;
      }
      for (const cb of Array.from(scope.querySelectorAll('[role="checkbox"],input[type="checkbox"]')).filter(visible)) {
        if (!isSendViaEmailDeliveryRowControl(cb)) continue;
        if (!isControlChecked(cb)) continue;
        try { cb.scrollIntoView({ block: 'center' }); } catch (e) {}
        try { cb.click(); } catch (e) {}
        await sleep(200);
      }
      await sleep(180);
    }
    if (emailDeliveryAlreadyOn(scope)) {
      reasons.push('send via email still on after automation; untick manually if you want it off');
    }
  }

  async function spinWait(ms) { await sleep(ms); }

  function schedulingRegionRoot(scope) {
    const s = scope || document;
    const heads = Array.from(s.querySelectorAll('h1,h2,h3,h4,h5')).filter(el => {
      if (!visible(el)) return false;
      const t = (el.innerText || '').trim().toLowerCase();
      return t === 'scheduling' || (t.startsWith('scheduling') && t.length < 48);
    });
    for (const h of heads) {
      let best = null;
      let bestLen = 1e9;
      let p = h.parentElement;
      for (let d = 0; d < 8 && p; d++, p = p.parentElement) {
        const tx = (p.innerText || '') || '';
        const low = tx.toLowerCase();
        if (!low.includes('scheduling')) continue;
        if (tx.length > 1400) continue;
        if (tx.length < bestLen) {
          bestLen = tx.length;
          best = p;
        }
      }
      if (best) return best;
    }
    return null;
  }

  function isEmailDeliveryTextBlob(low) {
    if (!low) return false;
    if (low.includes('send via email')) return true;
    if (low.includes('substack app') && low.includes('email')) return true;
    if (low.includes('notify subscribers')) return true;
    if (low.includes('deliver this post')) return true;
    if (low.includes('email and the substack')) return true;
    return false;
  }

  function isEmailDeliveryCheckboxBlob(blob) {
    if (!blob) return false;
    const low = blob.toLowerCase();
    return isEmailDeliveryTextBlob(low);
  }

  function schedulePublishCheckboxInScope(scope) {
    const regions = [];
    const s = scope || document;
    const trackedSched = querySchedulePostTrackControl(s);
    if (trackedSched && visible(trackedSched)) return trackedSched;
    const r0 = schedulingRegionRoot(scope);
    if (r0) regions.push(r0);
    regions.push(s);
    const needles = [
      'schedule time to publish',
      'schedule time to email and publish',
      'schedule time to email',
      'schedule to publish',
      'time to publish',
      'schedule this post',
      'schedule for later',
      'schedule post',
    ];
    for (const region of regions) {
      for (const n of needles) {
        const low = n.toLowerCase();
        const labels = Array.from(region.querySelectorAll('label,div,span,p')).filter(el => {
          if (!visible(el)) return false;
          const t = (el.innerText || '').toLowerCase();
          return t.includes(low) && t.length < 520;
        });
        for (const lab of labels) {
          let root = lab.closest('label') || lab;
          let cb = root.querySelector && root.querySelector('input[type="checkbox"],[role="checkbox"]');
          if (!cb && lab.parentElement) {
            cb = lab.parentElement.querySelector('input[type="checkbox"],[role="checkbox"]');
          }
          if (!cb && lab.nextElementSibling) {
            cb = lab.nextElementSibling.querySelector && lab.nextElementSibling.querySelector('input[type="checkbox"],[role="checkbox"]');
          }
          if (cb && visible(cb)) {
            if (isSendViaEmailDeliveryRowControl(cb)) continue;
            return cb;
          }
        }
      }
      const boxes = Array.from(region.querySelectorAll('input[type="checkbox"],[role="checkbox"]')).filter(visible);
      for (const cb of boxes) {
        const blob = rowBlobAround(cb);
        if (isSendViaEmailDeliveryRowControl(cb)) continue;
        if (!blob.includes('schedul')) continue;
        if (!(blob.includes('publish') || blob.includes('post') || blob.includes('later') || blob.includes('time'))) continue;
        if (blob.includes('email') && blob.includes('substack') && !blob.includes('schedule time')) continue;
        if (!blob.includes('schedule time') && !blob.includes('scheduling')) continue;
        return cb;
      }
      const switches = Array.from(region.querySelectorAll('[role="switch"],button[aria-checked]')).filter(visible);
      for (const sw of switches) {
        const blob = rowBlobAround(sw);
        if (isSendViaEmailDeliveryRowControl(sw)) continue;
        if (!blob.includes('schedul')) continue;
        if (!(blob.includes('publish') || blob.includes('post') || blob.includes('time'))) continue;
        if (!blob.includes('schedule time') && !blob.includes('scheduling')) continue;
        return sw;
      }
    }
    return null;
  }

  function scheduleTimeControlsVisible(scope) {
    const s = scope || document;
    const dtl = Array.from(s.querySelectorAll('input[type="datetime-local"]')).filter(visible);
    if (dtl.length > 0) return true;
    const tagEl = tagInput();
    const texts = Array.from(s.querySelectorAll('input[type="text"],input[type="search"],input:not([type])')).filter(inp => {
      if (!visible(inp)) return false;
      if (tagEl && inp === tagEl) return false;
      if (inputLooksLikeTagField(inp)) return false;
      const ty = (inp.type || 'text').toLowerCase();
      if (ty === 'checkbox' || ty === 'radio' || ty === 'hidden' || ty === 'file') return false;
      const r = inp.getBoundingClientRect();
      if (r.width < 100) return false;
      const blob = rowBlobAround(inp);
      if (blob.includes('add tags') && !blob.includes('scheduling')) return false;
      return blob.includes('scheduling') || (blob.includes('schedule time') && (blob.includes('publish') || blob.includes('email and publish')));
    });
    return texts.length > 0;
  }

  function isScheduleTimeToPublishChecked(scope) {
    if (scheduleTimeControlsVisible(scope)) return true;
    const cb = schedulePublishCheckboxInScope(scope);
    return cb && isControlChecked(cb);
  }

  async function clickSchedulePublishRow(scope) {
    const region = schedulingRegionRoot(scope) || scope || document;
    const want = [
      'schedule time to publish',
      'schedule time to email and publish',
      'schedule time to email',
      'schedule to publish',
      'time to publish',
      'schedule this post',
    ];
    const candidates = Array.from(region.querySelectorAll('label, button, div[role="button"], span, div')).filter(el => {
      if (!visible(el)) return false;
      const t = (el.innerText || '').trim().toLowerCase();
      if (t.length < 6 || t.length > 220) return false;
      if (isEmailDeliveryTextBlob(t)) return false;
      for (const w of want) {
        if (t.includes(w)) return true;
      }
      return t.includes('schedul') && t.includes('publish') && t.length < 120;
    });
    candidates.sort((a, b) => (a.innerText || '').length - (b.innerText || '').length);
    for (const el of candidates) {
      const cb = el.querySelector && el.querySelector('input[type="checkbox"],[role="checkbox"],[role="switch"]');
      if (cb && visible(cb)) {
        if (!isControlChecked(cb)) {
          try {
            cb.scrollIntoView({ block: 'center' });
            cb.click();
          } catch (e) {}
        } else {
          try {
            el.scrollIntoView({ block: 'center' });
            el.click();
          } catch (e) {}
        }
        await sleep(350);
        return true;
      }
      try {
        el.scrollIntoView({ block: 'center' });
        el.click();
        await sleep(350);
        return true;
      } catch (e) {}
    }
    return false;
  }

  async function ensureScheduleTimeToPublishEnabled(scope) {
    if (isScheduleTimeToPublishChecked(scope)) return true;
    const trackedSched = querySchedulePostTrackControl(scope);
    if (trackedSched && !isControlChecked(trackedSched)) {
      try {
        trackedSched.scrollIntoView({ block: 'center' });
        trackedSched.click();
      } catch (e) {}
      await sleep(450);
      if (isScheduleTimeToPublishChecked(scope)) return true;
    }
    const schedRegion = schedulingRegionRoot(scope);
    const narrow = schedRegion || scope;
    const needles = [
      'schedule time to publish',
      'schedule time to email and publish',
      'schedule time to email',
      'schedule to publish',
      'time to publish',
      'schedule this post',
      'schedule for later',
    ];
    for (const needle of needles) {
      if (clickCheckboxNearLabelSubstring(needle, narrow)) {
        await sleep(450);
        if (isScheduleTimeToPublishChecked(scope)) return true;
      }
    }
    await clickSchedulePublishRow(scope);
    if (isScheduleTimeToPublishChecked(scope)) return true;
    const cb = schedulePublishCheckboxInScope(scope);
    if (cb && !isControlChecked(cb)) {
      try {
        cb.scrollIntoView({ block: 'center' });
        cb.click();
      } catch (e) {}
      await sleep(450);
    }
    if (isScheduleTimeToPublishChecked(scope)) return true;
    if (schedRegion) {
      for (const sw of Array.from(schedRegion.querySelectorAll('[role="switch"]')).filter(visible)) {
        const blob = rowBlobAround(sw);
        if (isEmailDeliveryCheckboxBlob(blob)) continue;
        if (!blob.includes('schedul')) continue;
        if (!blob.includes('schedule time') && !blob.includes('scheduling')) continue;
        if (isControlChecked(sw)) return true;
        try {
          sw.scrollIntoView({ block: 'center' });
          sw.click();
        } catch (e) {}
        await sleep(400);
        if (isScheduleTimeToPublishChecked(scope)) return true;
      }
    }
    return isScheduleTimeToPublishChecked(scope);
  }

  function inputLooksLikeTagField(inp) {
    if (!inp) return false;
    const ph = ((inp.getAttribute('placeholder') || '') + ' ' + (inp.getAttribute('aria-label') || '')).toLowerCase();
    if (ph.includes('tag')) return true;
    const role = (inp.getAttribute('role') || '').toLowerCase();
    if (role === 'combobox') return true;
    const host = inp.closest && inp.closest('[role="combobox"]');
    if (host && visible(host)) {
      const htxt = ((host.innerText || '') + ph).toLowerCase();
      if (htxt.includes('tag') || htxt.includes('add tags')) return true;
    }
    return false;
  }

  function scheduleDateTextCandidates(scope) {
    const tagEl = tagInput();
    return Array.from(scope.querySelectorAll('input')).filter(inp => {
      if (!visible(inp)) return false;
      if (tagEl && inp === tagEl) return false;
      if (inputLooksLikeTagField(inp)) return false;
      const ty = (inp.type || 'text').toLowerCase();
      if (ty === 'checkbox' || ty === 'radio' || ty === 'hidden' || ty === 'file') return false;
      if (ty === 'datetime-local' || ty === 'date' || ty === 'time') return false;
      const r = inp.getBoundingClientRect();
      if (r.width < 100) return false;
      const blob = rowBlobAround(inp);
      if (blob.includes('add tags') && !blob.includes('schedule time')) return false;
      return (
        blob.includes('schedule time to publish') ||
        blob.includes('schedule time to email and publish') ||
        blob.includes('schedule time to email') ||
        (blob.includes('scheduling') && blob.includes('publish'))
      );
    });
  }

  async function fillScheduleDateTime(dtLocal, dateDisplay) {
    const hasLocal = dtLocal && dtLocal.trim();
    const hasDisp = dateDisplay && dateDisplay.trim();
    if (!hasLocal && !hasDisp) return true;
    const scope = scopeForPublishUI();
    if (!(await ensureScheduleTimeToPublishEnabled(scope))) {
      reasons.push('could not enable schedule checkbox (Schedule time to publish / email and publish)');
      return false;
    }
    await sleep(200);
    if (hasDisp) {
      for (let attempt = 0; attempt < 30; attempt++) {
        const cands = scheduleDateTextCandidates(scope);
        let best = null;
        if (cands.length === 1) {
          best = cands[0];
        } else if (cands.length > 1) {
          cands.sort((a, b) => {
            const ra = a.getBoundingClientRect();
            const rb = b.getBoundingClientRect();
            return (rb.width * rb.height) - (ra.width * ra.height);
          });
          best = cands[0];
        }
        if (best) {
          try {
            best.scrollIntoView({ block: 'center' });
          } catch (e) {}
          best.focus();
          try {
            best.value = dateDisplay.trim();
          } catch (e) {
            reasons.push('schedule text field set error');
            return false;
          }
          best.dispatchEvent(new Event('input', { bubbles: true }));
          best.dispatchEvent(new Event('change', { bubbles: true }));
          try {
            best.blur();
          } catch (e) {}
          return true;
        }
        await spinWait(120);
      }
    }
    if (hasLocal) {
      const dtlHits = Array.from(scope.querySelectorAll('input[type="datetime-local"]')).filter(visible);
      let dtl = null;
      if (dtlHits.length === 1) {
        dtl = dtlHits[0];
      } else if (dtlHits.length > 1) {
        for (const cand of dtlHits) {
          const blob = rowBlobAround(cand);
          if (blob.includes('schedule') || blob.includes('publish')) {
            dtl = cand;
            break;
          }
        }
        if (!dtl) dtl = dtlHits[0];
      }
      if (dtl) {
        try {
          dtl.value = dtLocal.trim();
        } catch (e) {
          reasons.push('datetime-local set error');
          return false;
        }
        dtl.dispatchEvent(new Event('input', { bubbles: true }));
        dtl.dispatchEvent(new Event('change', { bubbles: true }));
        return true;
      }
      const di = scope.querySelector('input[type="date"]');
      const ti = scope.querySelector('input[type="time"]');
      if (di && ti && visible(di) && visible(ti)) {
        const parts = dtLocal.trim().split('T');
        if (parts.length >= 2) {
          try {
            di.value = parts[0];
            ti.value = parts[1].slice(0, 5);
            di.dispatchEvent(new Event('change', { bubbles: true }));
            ti.dispatchEvent(new Event('change', { bubbles: true }));
            return true;
          } catch (e) {}
        }
      }
    }
    reasons.push('schedule datetime control not found');
    return false;
  }

  if (cfg.sectionLabel && cfg.sectionLabel.trim()) {
    if (!(await pickSection(cfg.sectionLabel.trim()))) {
      return JSON.stringify({ ok: false, reason: reasons.join('; ') });
    }
  }
  if (!cfg.tickEmailSubstack) {
    await tryUntickSubstackEmailDelivery();
  }
  // Schedule before tags: Substack only mounts the time field after the schedule row is on.
  // Tags are normally applied from Go via CDP mouse clicks (Headless UI ignores JS el.click);
  // when cfg.tags is empty this no-ops. Filling schedule before any tag typing avoids date bleed.
  if ((cfg.dateTimeLocal && cfg.dateTimeLocal.trim()) || (cfg.dateDisplay && cfg.dateDisplay.trim())) {
    if (!(await fillScheduleDateTime(cfg.dateTimeLocal || '', cfg.dateDisplay || ''))) {
      return JSON.stringify({ ok: false, reason: reasons.join('; ') });
    }
  }
  if (!(await addTags(cfg.tags || []))) {
    return JSON.stringify({ ok: false, reason: reasons.length ? reasons.join('; ') : 'tags failed' });
  }
  if (cfg.tickEmailSubstack) {
    if (!tryTickSubstackEmailDelivery()) {
      reasons.push('email/app delivery checkbox not found or copy changed; skipped (tick manually if needed)');
    }
  } else {
    await tryUntickSubstackEmailDelivery();
  }
  return JSON.stringify({ ok: true, reason: reasons.join('; ') });
})()`
	return js, nil
}
