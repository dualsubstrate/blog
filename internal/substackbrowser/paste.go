package substackbrowser

import (
	"encoding/json"
	"fmt"
)

// PasteHTMLIntoEditor runs in the browser context. It targets Substack-style
// ProseMirror roots first, then falls back to any contenteditable inside the
// post editor. It uses insertHTML, which is not guaranteed on every build of
// Chromium or every editor revision, but is the most reliable single step for
// rich paste in automation without OS clipboard integration.
func PasteHTMLIntoEditor(html string) (string, error) {
	expr, err := pasteExpression(html)
	if err != nil {
		return "", err
	}
	return expr, nil
}

func pasteExpression(html string) (string, error) {
	enc, err := json.Marshal(html)
	if err != nil {
		return "", err
	}
	// enc is a JSON string literal safe to splice into JS as a string value.
	js := `(function () {
  const html = ` + string(enc) + `;
  const selectors = [
    'div.ProseMirror[contenteditable="true"]',
    'div.ProseMirror[contenteditable=""]',
    '[contenteditable="true"].ProseMirror'
  ];
  let el = null;
  for (let i = 0; i < selectors.length; i++) {
    el = document.querySelector(selectors[i]);
    if (el) break;
  }
  if (!el) {
    const root = document.querySelector('[data-testid="post-editor"]') || document.querySelector('[data-testid="editor"]');
    if (root) {
      el = root.querySelector('[contenteditable="true"]');
    }
  }
  if (!el) {
    el = document.querySelector('[contenteditable="true"]');
  }
  if (!el) {
    return JSON.stringify({ ok: false, reason: 'no contenteditable editor matched' });
  }
  el.focus();
  const sel = window.getSelection();
  const range = document.createRange();
  range.selectNodeContents(el);
  range.collapse(false);
  sel.removeAllRanges();
  sel.addRange(range);
  let inserted = false;
  try {
    inserted = document.execCommand('insertHTML', false, html);
  } catch (e) {
    return JSON.stringify({ ok: false, reason: String(e) });
  }
  return JSON.stringify({ ok: inserted, reason: inserted ? '' : 'insertHTML returned false' });
})()`
	return js, nil
}

// PasteResult is the JSON shape returned by the injected paste script.
type PasteResult struct {
	OK     bool   `json:"ok"`
	Reason string `json:"reason"`
}

func ParsePasteResult(raw string) (PasteResult, error) {
	var r PasteResult
	if err := json.Unmarshal([]byte(raw), &r); err != nil {
		return PasteResult{}, fmt.Errorf("substackbrowser: paste result: %w", err)
	}
	return r, nil
}
