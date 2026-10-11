import { parseBackgroundWaveConfig } from './background-panorama.js';
import { variantIdFromIndex } from './export.js';
import { parseQrSizePercent } from './qr.js';
import { paletteById, matchPaletteId } from './palettes.js';
import {
  defaultIncludedSlideNumbers,
  defaultStripVariantIndices,
} from './studio-slide-inclusion.js';
import { mergeTheme, deckPaletteFromTheme, deckWavePaletteFromTheme, isWavePaletteLinked, DESIGN_CANVAS_SIZE } from './theme.js';

/**
 * @typedef {Object} FileStripSnapshot
 * @property {number[]} included
 * @property {Record<string, number>} variants
 */

/**
 * @typedef {Object} StateRow
 * @property {string} label
 * @property {string} fileValue
 * @property {string} browserValue
 * @property {boolean} differs
 */

/**
 * @typedef {Object} StudioStatePanelContext
 * @property {Record<string, unknown>} fileDeckSnapshot
 * @property {FileStripSnapshot} fileStripSnapshot
 * @property {() => { deck?: Record<string, unknown>, slides?: Array<{ number: number, variants?: unknown[] }> }} getDeck
 * @property {() => import('./theme.js').DeckTheme} getTheme
 * @property {() => string|null} getPaletteId
 * @property {() => boolean} getShowLineBoxes
 * @property {() => Set<number>} getIncludedSlideNumbers
 * @property {() => Map<number, number>} getStripVariantIndices
 * @property {() => Promise<void>} [onSave]
 */

/** @param {{ deck?: Record<string, unknown>, slides?: Array<{ number: number, variants?: unknown[] }> }} deckRef */
export function captureFileStripSnapshot(deckRef) {
  /** @type {Record<string, number>} */
  const variants = {};
  for (const [slideNumber, variantIndex] of defaultStripVariantIndices(deckRef)) {
    variants[String(slideNumber)] = variantIndex;
  }
  return {
    included: [...defaultIncludedSlideNumbers(deckRef)].sort((a, b) => a - b),
    variants,
  };
}

/**
 * @param {Record<string, unknown>} deckSpec
 * @returns {Record<string, unknown>}
 */
export function captureFileDeckSnapshot(deckSpec) {
  return structuredClone(deckSpec ?? {});
}

/** @param {string|null|undefined} paletteId */
function formatPaletteLabel(paletteId) {
  if (!paletteId) return 'Custom colors';
  return paletteById(paletteId)?.label ?? paletteId;
}

/** @param {import('./theme.js').WavePalette} palette */
function formatWavePaletteColors(palette) {
  return `base ${palette.background}, a1 ${palette.accent1}`;
}

/** @param {import('./background-panorama.js').BackgroundWaveConfig|undefined|null} wave */
function formatWave(wave) {
  if (!wave || wave.style === 'none') {
    if (wave?.hueShift) return `Flat · hue ${Math.round(wave.hueShift)}°`;
    return 'None (flat)';
  }
  const styleLabels = {
    drift: 'Drift',
    'mesh-corners': 'Mesh corners',
  };
  const style = styleLabels[wave.style] ?? wave.style;
  const lobes = wave.lobes == null ? 'auto' : String(wave.lobes);
  const parts = [
    style,
    `I ${roundNum(wave.intensity, 2)}`,
    `C ${roundNum(wave.color ?? 0.55, 2)}`,
    `V ${roundNum(wave.variety ?? 0.62, 2)}`,
    `blur ${roundNum(wave.blur, 2)}`,
    `R ${roundNum(wave.radius ?? 1, 2)}`,
    `phase ${roundNum(wave.phase, 2)}`,
    `lobes ${lobes}`,
  ];
  if (wave.hueShift) parts.push(`hue ${Math.round(wave.hueShift)}°`);
  return parts.join(' · ');
}

/** @param {number|null|undefined} value @param {number} digits */
function roundNum(value, digits) {
  if (!Number.isFinite(value)) return '—';
  const factor = 10 ** digits;
  return String(Math.round(value * factor) / factor);
}

/** @param {Record<string, unknown>} spec @param {'horizontal'|'vertical'} axis */
function marginFromSpec(spec, axis) {
  if (axis === 'horizontal') {
    return spec.marginHorizontal ?? spec.marginX ?? spec.margin ?? null;
  }
  return spec.marginVertical ?? spec.marginY ?? spec.margin ?? null;
}

/** @param {import('./theme.js').DeckTheme} theme @param {'horizontal'|'vertical'} axis */
function marginPercentFromTheme(theme, axis) {
  const px = axis === 'horizontal' ? theme.marginHorizontal : theme.marginVertical;
  return `${Math.round((px / DESIGN_CANVAS_SIZE) * 100)}%`;
}

/** @param {Record<string, unknown>} spec @param {'horizontal'|'vertical'} axis */
function formatMarginFromSpec(spec, axis) {
  const raw = marginFromSpec(spec, axis);
  if (typeof raw === 'string' && raw.trim()) return raw.trim();
  if (typeof raw === 'number' && Number.isFinite(raw)) return `${raw}%`;
  const theme = mergeTheme(spec);
  return marginPercentFromTheme(theme, axis);
}

/** @param {Record<string, unknown>} spec */
function formatLineHeights(spec) {
  const theme = mergeTheme(spec);
  const map = theme.lineHeights || {};
  const keys = ['normal', 'punch', 'header', 'footer'];
  return keys.map((key) => `${key} ${map[key] ?? 1}`).join(', ');
}

/** @param {Record<string, unknown>|undefined|null} cta */
function formatCtaLayout(cta) {
  if (!cta || typeof cta !== 'object') return '—';
  const src = /** @type {Record<string, unknown>} */ (cta);
  const parts = [];
  if (Number.isFinite(Number(src.featuredMaxHeight))) {
    parts.push(`img ${Number(src.featuredMaxHeight)}`);
  }
  const qrSizeRaw = (src.qr && typeof src.qr === 'object' && !Array.isArray(src.qr)
    ? /** @type {Record<string, unknown>} */ (src.qr).size
    : undefined);
  if (qrSizeRaw != null && qrSizeRaw !== '') {
    parts.push(`qr ${parseQrSizePercent(qrSizeRaw)}%`);
  }
  if (Number.isFinite(Number(src.brandMaxHeight))) {
    parts.push(`logo ${Number(src.brandMaxHeight)}`);
  }
  return parts.length > 0 ? parts.join(', ') : '—';
}

/**
 * @param {number[]} included
 * @param {Record<string, number>} variants
 * @param {Array<{ number: number, variants?: unknown[] }>|undefined} slides
 */
function formatStripPicks(included, variants, slides) {
  const slideList = slides ?? [];
  const parts = included.map((slideNumber) => {
    const slide = slideList.find((entry) => entry.number === slideNumber);
    const variantCount = slide?.variants?.length ?? 1;
    if (variantCount <= 1) return String(slideNumber);
    const variantIndex = variants[String(slideNumber)] ?? 0;
    return `${slideNumber}:${variantIdFromIndex(variantIndex)}`;
  });
  return parts.length > 0 ? parts.join(', ') : '—';
}

/** @param {Record<string, unknown>} spec */
function formatMotifStrip(spec) {
  const motif = spec?.motifStrip;
  if (!motif || typeof motif !== 'object' || Array.isArray(motif)) return '—';
  const src = typeof motif.src === 'string' ? motif.src.trim() : '';
  if (!src) return '—';
  const parts = [];
  parts.push(/** @type {Record<string, unknown>} */ (motif).enabled === false ? 'Off' : 'On');
  const offsetX = Number(motif.offsetX);
  const offsetY = Number(motif.offsetY);
  if (Number.isFinite(offsetX) && offsetX !== 0) parts.push(`X ${offsetX}`);
  if (Number.isFinite(offsetY) && offsetY !== 0) parts.push(`Y ${offsetY}`);
  const bandWidth = typeof motif.bandWidth === 'string' ? motif.bandWidth.trim() : '';
  const bandPct = /^(\d+(?:\.\d+)?)\s*%$/.exec(bandWidth);
  if (bandPct) {
    const delta = Math.round(Number(bandPct[1]) - 100);
    if (delta !== 0) parts.push(delta > 0 ? `+${delta}%` : `${delta}%`);
  } else if (bandWidth && bandWidth !== '100%') {
    parts.push(bandWidth);
  }
  return parts.join(', ');
}

/**
 * @param {StudioStatePanelContext} context
 * @returns {StateRow[]}
 */
export function buildStudioStateRows(context) {
  const fileSpec = context.fileDeckSnapshot;
  const fileTheme = mergeTheme(fileSpec);
  const fileWave = parseBackgroundWaveConfig({ backgroundWave: fileSpec.backgroundWave });
  const filePalette = deckPaletteFromTheme(fileTheme);
  const filePaletteId = matchPaletteId(fileTheme);

  const deckRef = context.getDeck();
  const liveSpec = deckRef.deck || {};
  const liveTheme = context.getTheme();
  const livePalette = deckPaletteFromTheme(liveTheme);
  const livePaletteId = context.getPaletteId() ?? matchPaletteId(liveTheme);

  /** @type {StateRow[]} */
  const rows = [];

  const pushRow = (label, fileValue, browserValue) => {
    rows.push({
      label,
      fileValue,
      browserValue,
      differs: fileValue !== browserValue,
    });
  };

  pushRow('Panoramic wave', formatWave(fileWave), formatWave(liveTheme.backgroundWave));

  const filePaletteLabel = `${formatPaletteLabel(filePaletteId)} (${formatWavePaletteColors(filePalette)})`;
  const browserPaletteLabel = `${formatPaletteLabel(livePaletteId)} (${formatWavePaletteColors(livePalette)})`;
  pushRow('Text palette', filePaletteLabel, browserPaletteLabel);

  const fileWaveLinked = isWavePaletteLinked(fileTheme);
  const liveWaveLinked = isWavePaletteLinked(liveTheme);
  const fileWavePalette = deckWavePaletteFromTheme(fileTheme);
  const liveWavePalette = deckWavePaletteFromTheme(liveTheme);
  const formatWavePaletteRow = (linked, wavePalette, textPalette) => {
    if (linked) return 'Linked to text palette';
    if (wavePalette) return formatWavePaletteColors(wavePalette);
    return formatWavePaletteColors({
      background: textPalette.background,
      accent1: textPalette.accent1,
      accent2: textPalette.accent2,
      muted: textPalette.muted,
    });
  };
  pushRow(
    'Wave palette',
    formatWavePaletteRow(fileWaveLinked, fileWavePalette, filePalette),
    formatWavePaletteRow(liveWaveLinked, liveWavePalette, livePalette),
  );

  pushRow(
    'Margins',
    `H ${formatMarginFromSpec(fileSpec, 'horizontal')}, V ${formatMarginFromSpec(fileSpec, 'vertical')}`,
    `H ${formatMarginFromSpec(liveSpec, 'horizontal')}, V ${formatMarginFromSpec(liveSpec, 'vertical')}`,
  );

  pushRow('Line heights', formatLineHeights(fileSpec), formatLineHeights(liveSpec));

  if (formatMotifStrip(fileSpec) !== '—' || formatMotifStrip(liveSpec) !== '—') {
    pushRow('Motif strip', formatMotifStrip(fileSpec), formatMotifStrip(liveSpec));
  }

  if (liveSpec.cta || fileSpec.cta) {
    pushRow(
      'CTA layout',
      formatCtaLayout(fileSpec.cta),
      formatCtaLayout(liveSpec.cta),
    );
  }

  pushRow('Line boxes debug', 'Off', context.getShowLineBoxes() ? 'On' : 'Off');

  const liveIncluded = [...context.getIncludedSlideNumbers()].sort((a, b) => a - b);
  const liveVariants = Object.fromEntries(context.getStripVariantIndices());
  pushRow(
    'Strip slides',
    formatStripPicks(context.fileStripSnapshot.included, context.fileStripSnapshot.variants, deckRef.slides),
    formatStripPicks(liveIncluded, liveVariants, deckRef.slides),
  );

  return rows;
}

/**
 * @param {StudioStatePanelContext} context
 */
export function renderStudioStatePanel(context) {
  const section = document.createElement('section');
  section.className = 'carousel-studio-state';
  section.setAttribute('aria-label', 'carousel.json vs browser settings');

  const head = document.createElement('div');
  head.className = 'carousel-studio-state-head';

  const title = document.createElement('h2');
  title.textContent = 'Settings: file vs browser';

  const hint = document.createElement('p');
  hint.className = 'carousel-studio-state-hint';
  hint.textContent = 'Highlighted rows differ from carousel.json. Save writes browser settings back to the source file.';

  const headText = document.createElement('div');
  headText.className = 'carousel-studio-state-head-text';
  headText.append(title, hint);

  head.appendChild(headText);

  if (typeof context.onSave === 'function') {
    const saveBtn = document.createElement('button');
    saveBtn.type = 'button';
    saveBtn.className = 'carousel-button carousel-button--compact carousel-studio-save';
    saveBtn.textContent = 'Save';
    saveBtn.title = 'Write browser settings to carousel.json';
    saveBtn.addEventListener('click', async () => {
      if (saveBtn.disabled) return;
      saveBtn.disabled = true;
      saveBtn.setAttribute('aria-busy', 'true');
      const previous = saveBtn.textContent;
      saveBtn.textContent = 'Saving…';
      try {
        await context.onSave();
        saveBtn.textContent = 'Saved';
        window.setTimeout(() => {
          saveBtn.textContent = previous;
        }, 1600);
      } catch (error) {
        if (!(error instanceof Error && error.name === 'AbortError')) {
          console.error('[carousel] save failed:', error);
          window.alert(error instanceof Error ? error.message : String(error));
        }
        saveBtn.textContent = previous;
      }
      saveBtn.disabled = false;
      saveBtn.removeAttribute('aria-busy');
    });
    head.appendChild(saveBtn);
  }

  const grid = document.createElement('div');
  grid.className = 'carousel-studio-state-grid';
  grid.setAttribute('role', 'grid');
  grid.setAttribute('aria-label', 'Settings comparison');

  const gridHead = document.createElement('div');
  gridHead.className = 'carousel-studio-state-grid-head';
  gridHead.setAttribute('role', 'row');
  gridHead.innerHTML = `
    <span class="carousel-studio-state-cell carousel-studio-state-cell--head" role="columnheader">Setting</span>
    <span class="carousel-studio-state-cell carousel-studio-state-cell--head" role="columnheader">carousel.json</span>
    <span class="carousel-studio-state-cell carousel-studio-state-cell--head" role="columnheader">Browser</span>
  `;

  const gridBody = document.createElement('div');
  gridBody.className = 'carousel-studio-state-grid-body';

  grid.append(gridHead, gridBody);
  section.append(head, grid);

  const refresh = () => {
    gridBody.replaceChildren();
    const rows = buildStudioStateRows(context);
    let diffCount = 0;
    for (const row of rows) {
      if (row.differs) diffCount += 1;
      const rowEl = document.createElement('div');
      rowEl.className = row.differs
        ? 'carousel-studio-state-row carousel-studio-state-row--diff'
        : 'carousel-studio-state-row';
      rowEl.setAttribute('role', 'row');
      rowEl.innerHTML = `
        <span class="carousel-studio-state-cell carousel-studio-state-cell--label" role="rowheader">${escapeHtml(row.label)}</span>
        <span class="carousel-studio-state-cell carousel-studio-state-cell--file" role="gridcell">${escapeHtml(row.fileValue)}</span>
        <span class="carousel-studio-state-cell carousel-studio-state-cell--browser" role="gridcell">${escapeHtml(row.browserValue)}</span>
      `;
      gridBody.appendChild(rowEl);
    }
    section.dataset.diffCount = String(diffCount);
    section.classList.toggle('carousel-studio-state--has-diff', diffCount > 0);
  };

  refresh();

  return { element: section, refresh };
}

/** @param {string} text */
function escapeHtml(text) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}
