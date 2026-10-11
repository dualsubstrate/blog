package substackbrowser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SetTitleAndSubtitle returns JS that tries to populate the Substack title and subtitle fields
// when present. It is best-effort and returns ok=true even if no fields are found, so paste
// can still proceed.
func SetTitleAndSubtitle(title string, subtitle string) (string, error) {
	tEnc, err := json.Marshal(strings.TrimSpace(title))
	if err != nil {
		return "", fmt.Errorf("substackbrowser: encode title: %w", err)
	}
	sEnc, err := json.Marshal(strings.TrimSpace(subtitle))
	if err != nil {
		return "", fmt.Errorf("substackbrowser: encode subtitle: %w", err)
	}
	js := `(async function () {
  const title = ` + string(tEnc) + `;
  const subtitle = ` + string(sEnc) + `;
  function sleep(ms) { return new Promise(r => setTimeout(r, ms)); }

  function setValue(el, val) {
    if (!el || !val) return false;
    try {
      el.focus();
      if (el.isContentEditable) {
        // Use execCommand insertText when possible. Substack sometimes ignores
        // direct textContent assignment for editor-managed fields.
        const sel = window.getSelection();
        const range = document.createRange();
        range.selectNodeContents(el);
        range.collapse(true);
        sel.removeAllRanges();
        sel.addRange(range);
        let ok = false;
        try {
          ok = document.execCommand('insertText', false, val);
        } catch (e) {
          ok = false;
        }
        if (!ok) {
          el.textContent = val;
        }
        el.dispatchEvent(new Event('input', { bubbles: true }));
        return true;
      }
      if ('value' in el) {
        el.value = val;
        el.dispatchEvent(new Event('input', { bubbles: true }));
        el.dispatchEvent(new Event('change', { bubbles: true }));
        return true;
      }
    } catch (e) {
      return false;
    }
    return false;
  }

  function pickTitleEl() {
    const candidates = [
      'textarea[placeholder="Title"]',
      'input[placeholder="Title"]',
      'textarea[placeholder="Título"]',
      'input[placeholder="Título"]',
      'textarea[placeholder="Titulo"]',
      'input[placeholder="Titulo"]',
      '[data-testid="post-title"] textarea',
      '[data-testid="post-title"] input',
      '[data-testid="post-title"] [contenteditable="true"]'
    ];
    for (const sel of candidates) {
      const el = document.querySelector(sel);
      if (el) return el;
    }
    const textareas = Array.from(document.querySelectorAll('textarea, input, [contenteditable="true"]'));
    for (const el of textareas) {
      const ph = (el.getAttribute('placeholder') || el.getAttribute('data-placeholder') || el.getAttribute('aria-label') || '').toLowerCase();
      if (ph === 'title' || ph === 'título' || ph === 'titulo') return el;
    }
    return null;
  }

  function looksLikeSubtitleHint(ph) {
    const p = (ph || '').toLowerCase();
    if (!p) return false;
    if (p.includes('subtitle') || p.includes('subtítulo') || p.includes('subtitulo')) return true;
    if (p.includes('add a subtitle')) return true;
    if (p.includes('añade un subtítulo') || p.includes('anade un subtitulo')) return true;
    if (p.includes('agrega un subtítulo') || p.includes('agrega un subtitulo')) return true;
    if (p.includes('agregar subtítulo') || p.includes('agregar subtitulo')) return true;
    return false;
  }

  function pickSubtitleEl() {
    const candidates = [
      'textarea[placeholder="Add a subtitle..."]',
      'input[placeholder="Add a subtitle..."]',
      'textarea[placeholder="Añade un subtítulo..."]',
      'input[placeholder="Añade un subtítulo..."]',
      'textarea[placeholder="Agrega un subtítulo..."]',
      'input[placeholder="Agrega un subtítulo..."]',
      '[data-placeholder="Add a subtitle..."]',
      '[data-placeholder="Añade un subtítulo..."]',
      '[aria-label="Add a subtitle..."]',
      '[aria-label="Añade un subtítulo..."]',
      '[data-testid="post-subtitle"] textarea',
      '[data-testid="post-subtitle"] input',
      '[data-testid="post-subtitle"] [contenteditable="true"]',
      '[data-testid="subtitle"] [contenteditable="true"]',
      'textarea[placeholder="Subtitle"]',
      'input[placeholder="Subtitle"]',
      'textarea[placeholder="Subtítulo"]',
      'input[placeholder="Subtítulo"]',
      'textarea[placeholder="Subheading"]',
      'input[placeholder="Subheading"]'
    ];
    for (const sel of candidates) {
      const el = document.querySelector(sel);
      if (el) return el;
    }
    const textareas = Array.from(document.querySelectorAll('textarea, input'));
    for (const el of textareas) {
      const ph = el.getAttribute('placeholder') || '';
      if (looksLikeSubtitleHint(ph)) return el;
    }
    const editables = Array.from(document.querySelectorAll('[contenteditable="true"]'));
    for (const el of editables) {
      const ph = el.getAttribute('data-placeholder') || el.getAttribute('aria-label') || el.getAttribute('placeholder') || '';
      if (looksLikeSubtitleHint(ph)) return el;
    }
    return null;
  }

  let didTitle = false;
  let didSub = false;

  const deadline = Date.now() + 5000;
  while (Date.now() < deadline) {
    const titleEl = pickTitleEl();
    const subEl = pickSubtitleEl();
    if (!didTitle && titleEl && title) didTitle = setValue(titleEl, title);
    if (!didSub && subEl && subtitle) didSub = setValue(subEl, subtitle);
    if ((didTitle || !title) && (didSub || !subtitle)) break;
    await sleep(60);
  }

  return JSON.stringify({ ok: true, reason: '', title_set: didTitle, subtitle_set: didSub });
})()`
	return js, nil
}
