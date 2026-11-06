const DOWNLOAD_ICON = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M12 3v12"/><path d="m7 11 5 5 5-5"/><path d="M5 21h14"/></svg>`;

/**
 * @param {{ id?: string, className?: string, label: string, title?: string }} options
 * @returns {HTMLButtonElement}
 */
export function createDownloadButton(options) {
  const btn = document.createElement('button');
  btn.type = 'button';
  const classes = ['carousel-download-btn', 'carousel-icon-button', options.className].filter(Boolean);
  btn.className = classes.join(' ');
  btn.innerHTML = DOWNLOAD_ICON;
  btn.dataset.downloadLabel = options.label;
  btn.setAttribute('aria-label', options.label);
  btn.title = options.title ?? options.label;
  if (options.id) btn.id = options.id;
  return btn;
}
