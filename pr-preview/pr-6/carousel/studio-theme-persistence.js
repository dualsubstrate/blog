import { parseBackgroundWaveConfig } from './background-panorama.js';
import { parseQrSizePercent } from './qr.js';
import { deckPaletteFromTheme, normalizeWavePalette, panoramaPaletteFromTextPalette } from './theme.js';

const STORAGE_PREFIX = 'carousel-studio-theme';

/** @param {string} slug */
export function themeStorageKey(slug) {
  return `${STORAGE_PREFIX}:${slug}`;
}

/**
 * @typedef {Object} StudioThemeState
 * @property {string|null} [paletteId]
 * @property {import('./theme.js').DeckPalette} [palette]
 * @property {import('./theme.js').WavePalette|null} [wavePalette]
 * @property {boolean} [wavePaletteLinked]
 * @property {import('./background-panorama.js').BackgroundWaveConfig} [backgroundWave]
 * @property {string} [marginHorizontal]
 * @property {string} [marginVertical]
 * @property {Record<string, number>} [lineHeights]
 * @property {{ featuredMaxHeight?: number, qrSize?: number, brandMaxHeight?: number }} [cta]
 * @property {{ enabled?: boolean, offsetX?: number, offsetY?: number, bandWidth?: string }} [motifStrip]
 * @property {boolean} [showLineBoxes]
 */

/**
 * @param {string} slug
 * @returns {StudioThemeState|null}
 */
export function readStudioThemeState(slug) {
  try {
    const raw = localStorage.getItem(themeStorageKey(slug));
    if (raw == null) return null;
    const parsed = JSON.parse(raw);
    return parsed && typeof parsed === 'object' ? /** @type {StudioThemeState} */ (parsed) : null;
  } catch {
    return null;
  }
}

/**
 * @param {string} slug
 * @param {{ deck?: Record<string, unknown> }} deck
 * @param {import('./theme.js').DeckTheme} theme
 * @returns {{ paletteId: string|null, showLineBoxes: boolean|null }}
 */
export function applyStudioThemeState(slug, deck, theme) {
  const saved = readStudioThemeState(slug);
  if (!saved) {
    return { paletteId: null, showLineBoxes: null };
  }

  if (!deck.deck) deck.deck = {};

  if (saved.palette && typeof saved.palette === 'object') {
    const palette = saved.palette;
    if (palette.background) theme.background = palette.background;
    if (palette.text) theme.text = palette.text;
    if (palette.muted) theme.muted = palette.muted;
    if (palette.accent1) theme.accent1 = palette.accent1;
    if (palette.accent2) theme.accent2 = palette.accent2;
    deck.deck.palette = deckPaletteFromTheme(theme);
  }

  if (saved.wavePalette && typeof saved.wavePalette === 'object') {
    theme.wavePalette = normalizeWavePalette(
      saved.wavePalette,
      panoramaPaletteFromTextPalette(deckPaletteFromTheme(theme)),
    );
    theme.wavePaletteLinked = false;
    deck.deck.wavePalette = { ...theme.wavePalette };
    deck.deck.wavePaletteLinked = false;
  } else if (saved.wavePaletteLinked === false) {
    theme.wavePaletteLinked = false;
    theme.wavePalette = normalizeWavePalette(
      {},
      panoramaPaletteFromTextPalette(deckPaletteFromTheme(theme)),
    );
    deck.deck.wavePaletteLinked = false;
  } else if (saved.wavePaletteLinked === true) {
    theme.wavePaletteLinked = true;
    theme.wavePalette = null;
    delete deck.deck.wavePalette;
    deck.deck.wavePaletteLinked = true;
  }

  if (saved.backgroundWave && typeof saved.backgroundWave === 'object') {
    theme.backgroundWave = parseBackgroundWaveConfig({ backgroundWave: saved.backgroundWave });
    deck.deck.backgroundWave = { ...theme.backgroundWave };
  }

  if (typeof saved.marginHorizontal === 'string' && saved.marginHorizontal.trim()) {
    deck.deck.marginHorizontal = saved.marginHorizontal;
  }
  if (typeof saved.marginVertical === 'string' && saved.marginVertical.trim()) {
    deck.deck.marginVertical = saved.marginVertical;
  }

  if (saved.lineHeights && typeof saved.lineHeights === 'object') {
    deck.deck.lineHeights = {
      ...(deck.deck.lineHeights && typeof deck.deck.lineHeights === 'object'
        ? deck.deck.lineHeights
        : {}),
      ...saved.lineHeights,
    };
    theme.lineHeights = {
      ...(theme.lineHeights || {}),
      ...saved.lineHeights,
    };
  }

  if (saved.cta && typeof saved.cta === 'object') {
    const existing = deck.deck.cta && typeof deck.deck.cta === 'object' ? deck.deck.cta : {};
    const { qrSize, ...ctaRest } = /** @type {Record<string, unknown>} */ (saved.cta);
    deck.deck.cta = { ...existing, ...ctaRest };
    if (qrSize != null && Number.isFinite(Number(qrSize))) {
      const formatted = `${Math.round(Number(qrSize))}%`;
      const existingQr = deck.deck.cta.qr && typeof deck.deck.cta.qr === 'object' && !Array.isArray(deck.deck.cta.qr)
        ? { .../** @type {Record<string, unknown>} */ (deck.deck.cta.qr) }
        : {};
      existingQr.size = formatted;
      deck.deck.cta.qr = existingQr;
      delete deck.deck.cta.qrSize;
    }
  }

  if (saved.motifStrip && typeof saved.motifStrip === 'object') {
    const existing = deck.deck.motifStrip;
    if (existing && typeof existing === 'object' && !Array.isArray(existing)) {
      const motif = /** @type {Record<string, unknown>} */ (existing);
      if (typeof saved.motifStrip.enabled === 'boolean') {
        motif.enabled = saved.motifStrip.enabled;
      }
      const offsetX = Number(saved.motifStrip.offsetX);
      if (Number.isFinite(offsetX) && offsetX !== 0) {
        motif.offsetX = offsetX;
      } else {
        delete motif.offsetX;
      }
      const offsetY = Number(saved.motifStrip.offsetY);
      if (Number.isFinite(offsetY) && offsetY !== 0) {
        motif.offsetY = offsetY;
      } else {
        delete motif.offsetY;
      }
      const bandWidth = typeof saved.motifStrip.bandWidth === 'string'
        ? saved.motifStrip.bandWidth.trim()
        : '';
      if (bandWidth) {
        motif.bandWidth = bandWidth;
      }
    }
  }

  return {
    paletteId: typeof saved.paletteId === 'string' ? saved.paletteId : null,
    showLineBoxes: typeof saved.showLineBoxes === 'boolean' ? saved.showLineBoxes : null,
  };
}

/**
 * @param {string} slug
 * @param {{ deck?: Record<string, unknown> }} deck
 * @param {import('./theme.js').DeckTheme} theme
 * @param {string|null|undefined} paletteId
 * @param {boolean} showLineBoxes
 */
export function writeStudioThemeState(slug, deck, theme, paletteId, showLineBoxes) {
  /** @type {StudioThemeState} */
  const payload = {
    paletteId: paletteId ?? null,
    palette: deckPaletteFromTheme(theme),
    wavePaletteLinked: theme.wavePaletteLinked !== false,
    showLineBoxes,
  };

  const wavePalette = theme.wavePalette && typeof theme.wavePalette === 'object'
    ? normalizeWavePalette(theme.wavePalette, deckPaletteFromTheme(theme))
    : null;
  if (wavePalette && theme.wavePaletteLinked === false) {
    payload.wavePalette = wavePalette;
    payload.wavePaletteLinked = false;
  }

  if (theme.backgroundWave) {
    payload.backgroundWave = { ...theme.backgroundWave };
  }

  const spec = deck.deck || {};
  if (typeof spec.marginHorizontal === 'string' && spec.marginHorizontal.trim()) {
    payload.marginHorizontal = spec.marginHorizontal;
  }
  if (typeof spec.marginVertical === 'string' && spec.marginVertical.trim()) {
    payload.marginVertical = spec.marginVertical;
  }
  if (spec.lineHeights && typeof spec.lineHeights === 'object') {
    payload.lineHeights = { .../** @type {Record<string, number>} */ (spec.lineHeights) };
  }
  if (spec.cta && typeof spec.cta === 'object') {
    const cta = /** @type {Record<string, unknown>} */ (spec.cta);
    /** @type {StudioThemeState['cta']} */
    const savedCta = {};
    if (Number.isFinite(Number(cta.featuredMaxHeight))) {
      savedCta.featuredMaxHeight = Number(cta.featuredMaxHeight);
    }
    const qrSizeRaw = (cta.qr && typeof cta.qr === 'object' && !Array.isArray(cta.qr)
      ? /** @type {Record<string, unknown>} */ (cta.qr).size
      : undefined);
    if (qrSizeRaw != null && qrSizeRaw !== '') {
      savedCta.qrSize = parseQrSizePercent(qrSizeRaw);
    }
    if (Number.isFinite(Number(cta.brandMaxHeight))) {
      savedCta.brandMaxHeight = Number(cta.brandMaxHeight);
    }
    if (Object.keys(savedCta).length > 0) {
      payload.cta = savedCta;
    }
  }

  const motif = spec.motifStrip;
  if (motif && typeof motif === 'object' && !Array.isArray(motif)) {
    const motifSpec = /** @type {Record<string, unknown>} */ (motif);
    const src = typeof motifSpec.src === 'string' ? motifSpec.src.trim() : '';
    if (src) {
      /** @type {NonNullable<StudioThemeState['motifStrip']>} */
      const savedMotif = {
        enabled: motifSpec.enabled !== false,
      };
      const offsetX = Number(motifSpec.offsetX);
      if (Number.isFinite(offsetX) && offsetX !== 0) {
        savedMotif.offsetX = offsetX;
      }
      const offsetY = Number(motifSpec.offsetY);
      if (Number.isFinite(offsetY) && offsetY !== 0) {
        savedMotif.offsetY = offsetY;
      }
      const bandWidth = typeof motifSpec.bandWidth === 'string' ? motifSpec.bandWidth.trim() : '';
      if (bandWidth) {
        savedMotif.bandWidth = bandWidth;
      }
      payload.motifStrip = savedMotif;
    }
  }

  try {
    localStorage.setItem(themeStorageKey(slug), JSON.stringify(payload));
  } catch {
    // ignore private mode / quota
  }
}

/** @param {string} slug */
export function clearStudioThemeState(slug) {
  try {
    localStorage.removeItem(themeStorageKey(slug));
  } catch {
    // ignore
  }
}
