package substackbrowser

import (
	"encoding/json"
	"log"
	"strings"
)

// PublishFlowStageSnapshot is the JSON shape from PublishFlowStageSnapshotJS.
type PublishFlowStageSnapshot struct {
	EditorProseMirror      bool `json:"editorProseMirror"`
	ContinueButton         bool `json:"continueButton"`
	PublishSettingsVisible bool `json:"publishSettingsVisible"`
	PublishModalLikely     bool `json:"publishModalLikely"`
	DialogCount            int  `json:"dialogCount"`
}

// PublishFlowStageSnapshotJS returns synchronous JS that reports whether key Substack publish-flow
// UI pieces are present (editor, Continue gate, publish settings body).
func PublishFlowStageSnapshotJS() string {
	return `(function() {
  function visible(el) {
    if (!el) return false;
    const r = el.getBoundingClientRect();
    return r.width > 0 && r.height > 0;
  }
  function editorPM() {
    const el = document.querySelector('div.ProseMirror');
    return !!(el && visible(el));
  }
  function hasContinue() {
    for (const b of document.querySelectorAll('button,[role="button"]')) {
      if (!visible(b)) continue;
      const t = (b.innerText || '').trim().toLowerCase();
      if (t === 'continue' || t === 'continuar') return true;
    }
    return false;
  }
  function bodyLow() {
    return (document.body && document.body.innerText || '').toLowerCase();
  }
  function settingsBodyVisible() {
    const t = bodyLow();
    return t.includes('this post belongs') ||
      (t.includes('add tags') && (t.includes('audience') || t.includes('delivery'))) ||
      (t.includes('añadir etiquetas') && (t.includes('audiencia') || t.includes('entrega'))) ||
      t.includes('esta publicación pertenece') ||
      t.includes('este artículo pertenece');
  }
  function modalLikelyFlag() {
    const direct = document.querySelector('[data-testid="publish-modal"]');
    if (direct && visible(direct)) return true;
    const t = bodyLow();
    return t.includes('add tags') && (t.includes('delivery') || t.includes('scheduling')) ||
      t.includes('añadir etiquetas') && (t.includes('entrega') || t.includes('program'));
  }
  function dialogCount() {
    return document.querySelectorAll('[role="dialog"]').length;
  }
  return JSON.stringify({
    editorProseMirror: editorPM(),
    continueButton: hasContinue(),
    publishSettingsVisible: settingsBodyVisible(),
    publishModalLikely: modalLikelyFlag(),
    dialogCount: dialogCount()
  });
})()`
}

// RecoverPublishFlowScheduleRetrySyncJS dismisses overlays before a schedule-script retry.
// It does not click header Publish / Next / Send (that can open a new composer or advance the wrong flow).
func RecoverPublishFlowScheduleRetrySyncJS() string {
	return `(function() {
  function visible(el) {
    if (!el) return false;
    const r = el.getBoundingClientRect();
    return r.width > 0 && r.height > 0;
  }
  const steps = [];
  try {
    const tgt = document.body || document.documentElement;
    if (tgt) {
      for (let i = 0; i < 2; i++) {
        tgt.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', code: 'Escape', keyCode: 27, bubbles: true }));
        tgt.dispatchEvent(new KeyboardEvent('keyup', { key: 'Escape', code: 'Escape', keyCode: 27, bubbles: true }));
      }
      steps.push('escape');
    } else {
      steps.push('escape_skip_no_root');
    }
  } catch (e) {
    steps.push('escape_err');
  }
  const nibbles = ['ok', 'got it', 'dismiss', 'not now'];
  for (const el of document.querySelectorAll('button,[role="button"]')) {
    if (!visible(el)) continue;
    const t = (el.innerText || '').trim().toLowerCase();
    if (!t || t.length > 48) continue;
    if (nibbles.includes(t)) {
      try {
        el.click();
        steps.push('click_' + t.replace(/\s+/g, '_'));
      } catch (e) {}
    }
  }
  return JSON.stringify({ ok: true, reason: steps.join(',') });
})()`
}

// RecoverPublishFlowSyncJS dismisses stray overlays (Escape, small OK dialogs) and tries to reopen
// the publish flow from the editor header (Next / Publish) when the primary gate is missing.
func RecoverPublishFlowSyncJS() string {
	return `(function() {
  function visible(el) {
    if (!el) return false;
    const r = el.getBoundingClientRect();
    return r.width > 0 && r.height > 0;
  }
  const steps = [];
  try {
    const tgt = document.body || document.documentElement;
    if (tgt) {
      for (let i = 0; i < 2; i++) {
        tgt.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', code: 'Escape', keyCode: 27, bubbles: true }));
        tgt.dispatchEvent(new KeyboardEvent('keyup', { key: 'Escape', code: 'Escape', keyCode: 27, bubbles: true }));
      }
      steps.push('escape');
    } else {
      steps.push('escape_skip_no_root');
    }
  } catch (e) {
    steps.push('escape_err');
  }
  const nibbles = ['ok', 'got it', 'dismiss', 'not now'];
  for (const el of document.querySelectorAll('button,[role="button"]')) {
    if (!visible(el)) continue;
    const t = (el.innerText || '').trim().toLowerCase();
    if (!t || t.length > 48) continue;
    if (nibbles.includes(t)) {
      try {
        el.click();
        steps.push('click_' + t.replace(/\s+/g, '_'));
      } catch (e) {}
    }
  }
  const hdrSelectors = 'header button, [data-testid="header"] button, nav button, [class*="header"] button';
  for (const el of document.querySelectorAll(hdrSelectors)) {
    if (!visible(el)) continue;
    const t = (el.innerText || '').trim().toLowerCase();
    if (t === 'next' || t === 'publish' || t === 'send' || (t.length < 24 && t.includes('publish'))) {
      try {
        el.scrollIntoView({ block: 'center' });
        el.click();
        steps.push('header_' + t.replace(/\s+/g, '_'));
        return JSON.stringify({ ok: true, reason: steps.join(',') });
      } catch (e) {
        steps.push('header_err');
      }
    }
  }
  return JSON.stringify({ ok: true, reason: steps.join(',') });
})()`
}

// LogPublishFlowSnapshot decodes snapshot JSON from PublishFlowStageSnapshotJS and logs one line.
func LogPublishFlowSnapshot(attempt, maxAttempts int, raw string) {
	raw = strings.TrimSpace(raw)
	var s PublishFlowStageSnapshot
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		head := raw
		if len(head) > 160 {
			head = head[:160] + "…"
		}
		log.Printf("substackbrowser: publish flow snapshot attempt %d/%d: decode error: %v raw=%q", attempt, maxAttempts, err, head)
		return
	}
	log.Printf("substackbrowser: publish flow preflight pass %d/%d (before schedule script; publish_settings is usually false until Continue opens the modal): editor=%v continue=%v publish_settings=%v publish_modal_est=%v dialogs=%d",
		attempt, maxAttempts, s.EditorProseMirror, s.ContinueButton, s.PublishSettingsVisible, s.PublishModalLikely, s.DialogCount)
}
