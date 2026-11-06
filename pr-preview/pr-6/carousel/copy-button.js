const COPY_ICON = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/></svg>`;

const SUCCESS_ICON = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><polyline points="20 6 9 17 4 12"/></svg>`;

/**
 * @param {{ className?: string, label: string, title?: string }} options
 * @returns {HTMLButtonElement}
 */
export function createCopyButton(options) {
  const btn = document.createElement('button');
  btn.type = 'button';
  const classes = ['carousel-copy-btn', options.className].filter(Boolean);
  btn.className = classes.join(' ');
  btn.innerHTML = COPY_ICON;
  btn.dataset.copyLabel = options.label;
  btn.setAttribute('aria-label', options.label);
  btn.title = options.title ?? options.label;
  return btn;
}

/** @param {HTMLButtonElement} btn */
function restoreCopyButtonDefault(btn) {
  btn.innerHTML = COPY_ICON;
  btn.classList.remove('carousel-color-copied', 'carousel-copy-btn--failed');
  const label = btn.dataset.copyLabel ?? 'Copy';
  btn.setAttribute('aria-label', label);
  btn.title = label;
}

/** @param {HTMLButtonElement} btn */
function showCopyButtonSuccess(btn) {
  btn.innerHTML = SUCCESS_ICON;
  btn.classList.add('carousel-color-copied');
  btn.classList.remove('carousel-copy-btn--failed');
  btn.setAttribute('aria-label', 'Copied');
}

/** @param {HTMLButtonElement} btn */
function showCopyButtonFailed(btn) {
  btn.innerHTML = COPY_ICON;
  btn.classList.add('carousel-copy-btn--failed');
  btn.classList.remove('carousel-color-copied');
  btn.setAttribute('aria-label', 'Copy failed');
}

/**
 * @param {HTMLElement} trigger
 * @param {string} text
 * @param {string} [restoreLabel]
 */
export async function copyText(trigger, text, restoreLabel) {
  if (restoreLabel && trigger instanceof HTMLButtonElement) {
    trigger.dataset.copyLabel = restoreLabel;
  }
  const copied = await writeClipboard(text);
  if (copied) {
    showCopyFeedback(trigger);
    return;
  }
  if (trigger instanceof HTMLButtonElement && trigger.classList.contains('carousel-copy-btn')) {
    showCopyButtonFailed(trigger);
    window.setTimeout(() => restoreCopyButtonDefault(trigger), 1200);
    return;
  }
  const original = restoreLabel ?? trigger.textContent ?? 'Copy';
  trigger.textContent = 'Copy failed';
  trigger.classList.add('carousel-color-copied');
  window.setTimeout(() => {
    trigger.classList.remove('carousel-color-copied');
    trigger.textContent = original;
  }, 1200);
}

/** @param {string} text */
async function writeClipboard(text) {
  try {
    await navigator.clipboard.writeText(text);
    return true;
  } catch {
    try {
      const area = document.createElement('textarea');
      area.value = text;
      area.setAttribute('readonly', '');
      area.style.position = 'fixed';
      area.style.left = '-9999px';
      document.body.appendChild(area);
      area.select();
      const ok = document.execCommand('copy');
      document.body.removeChild(area);
      return ok;
    } catch {
      return false;
    }
  }
}

/** @param {HTMLElement} el */
function showCopyFeedback(el) {
  if (el instanceof HTMLButtonElement && el.classList.contains('carousel-copy-btn')) {
    showCopyButtonSuccess(el);
    window.setTimeout(() => restoreCopyButtonDefault(el), 1200);
    return;
  }
  const original = el.textContent ?? 'Copy';
  el.classList.add('carousel-color-copied');
  el.textContent = 'Copied!';
  window.setTimeout(() => {
    el.classList.remove('carousel-color-copied');
    el.textContent = original;
  }, 1200);
}
