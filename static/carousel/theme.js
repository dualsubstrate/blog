import {
  parseBackgroundGradientMode,
  parseBackgroundWaveConfig,
} from './background-panorama.js';
import { buildBackgroundGradient, normalizeGradientPresetId, shiftHex } from './background.js';
import { CAROUSEL_SLIDE_WIDTH_PX, PANORAMA_SLIDE_WIDTH_PX } from './slide-constants.js';

/** @typedef {'display' | 'body'} FontRole */
/** @typedef {'text' | 'muted' | 'accent1' | 'accent2'} ColorToken */
/** @typedef {'header' | 'body' | 'footer' | 'grid'} Section */
/** @typedef {'normal' | 'punch'} BodyEmphasis */

/**
 * @typedef {Object} DeckTheme
 * @property {string} displayFont
 * @property {string} bodyFont
 * @property {string} background
 * @property {string|null} [backgroundGradientPreset]
 * @property {import('./background.js').BackgroundGradient|null} [backgroundGradient]
 * @property {'panoramic-wave'} [backgroundGradientMode]
 * @property {import('./background-panorama.js').BackgroundWaveConfig} [backgroundWave]
 * @property {string} text
 * @property {string} muted
 * @property {string} accent1
 * @property {string|null} [accent2]
 * @property {import('./theme.js').WavePalette|null} [wavePalette] Separate wash colors when unlinked from text palette
 * @property {boolean} [wavePaletteLinked] When true (default), wave uses text palette base + accents for washes
 * @property {number} marginHorizontal
 * @property {number} marginVertical
 * @property {number} [margin]
 * @property {number} size Canvas width in px (design width scales from this).
 * @property {number} sizeHeight Canvas height in px (derived from width and aspect ratio).
 * @property {number} aspectRatioWidth Aspect ratio width unit (default 4).
 * @property {number} aspectRatioHeight Aspect ratio height unit (default 5).
 * @property {number} contentMaxWidth
 * @property {number} [previewMaxPx]
 * @property {Record<'header'|'body'|'footer', number>} [fontSizes] Base px at 1080 per section
 * @property {Record<'normal'|'punch', number>} [emphasisScale] Multiplier on body base (`punch` default 1.25)
 * @property {Record<string, number>} [lineHeights] Multiplier on orange band for blue line slot (header, footer, normal, punch)
 * @property {number} [emphasisGap] Textless gap at normal↔punch seams, fraction of body base
 */

export const SLIDE_SIZE_MIN = PANORAMA_SLIDE_WIDTH_PX;
export const SLIDE_SIZE_MAX = CAROUSEL_SLIDE_WIDTH_PX;
/** Default export aspect ratio (portrait 4:5, e.g. 1080×1350). */
export const DEFAULT_ASPECT_RATIO_WIDTH = 4;
export const DEFAULT_ASPECT_RATIO_HEIGHT = 5;
/** Design canvas width for font sizes, margins, and horizontal scale (`theme.size`). */
export const DESIGN_CANVAS_SIZE = SLIDE_SIZE_MAX;
/** Design canvas height at default aspect ratio (1350 when width is 1080). */
export const DESIGN_CANVAS_HEIGHT = Math.round(
  DESIGN_CANVAS_SIZE * DEFAULT_ASPECT_RATIO_HEIGHT / DEFAULT_ASPECT_RATIO_WIDTH,
);
export const DEFAULT_MARGIN_PX = 112;
/** Default left/right inset when no `marginHorizontal` / `margin` (~10.4% at 1080). */
export const DEFAULT_MARGIN_FRACTION = DEFAULT_MARGIN_PX / DESIGN_CANVAS_SIZE;
/** Default top/bottom inset when no `marginVertical` / `margin`. */
export const DEFAULT_MARGIN_VERTICAL_FRACTION = 0.1;
export const DEFAULT_PREVIEW_MAX_PX = 300;
export const DEFAULT_EMPHASIS_GAP = 0.25;
export const DEFAULT_EMPHASIS_SCALE_PUNCH = 1.25;

export const DEFAULT_THEME = {
  displayFont: 'Playfair Display',
  bodyFont: 'Source Sans 3',
  background: '#1a1e26',
  text: '#f5f5f0',
  muted: '#9aa3ad',
  accent1: '#d69a80',
  accent2: null,
  size: CAROUSEL_SLIDE_WIDTH_PX,
  previewMaxPx: DEFAULT_PREVIEW_MAX_PX,
  emphasisGap: DEFAULT_EMPHASIS_GAP,
  fontSizes: { header: 60, body: 80, footer: 52 },
  emphasisScale: { normal: 1, punch: DEFAULT_EMPHASIS_SCALE_PUNCH },
  lineHeights: { header: 1, footer: 1, normal: 1, punch: 1 },
};

/** @param {number} value @param {number} min @param {number} max */
function clamp(value, min, max) {
  return Math.max(min, Math.min(max, value));
}

/**
 * @typedef {{ width: number, height: number }} AspectRatio
 */

/**
 * @param {unknown} value
 * @returns {AspectRatio}
 */
export function parseAspectRatio(value) {
  if (value && typeof value === 'object' && !Array.isArray(value)) {
    const obj = /** @type {Record<string, unknown>} */ (value);
    const w = Number(obj.width ?? obj.w);
    const h = Number(obj.height ?? obj.h);
    if (Number.isFinite(w) && Number.isFinite(h) && w > 0 && h > 0) {
      return { width: w, height: h };
    }
  }
  if (typeof value === 'string') {
    const match = /^(\d+(?:\.\d+)?)\s*:\s*(\d+(?:\.\d+)?)$/.exec(value.trim());
    if (match) {
      const w = parseFloat(match[1]);
      const h = parseFloat(match[2]);
      if (Number.isFinite(w) && Number.isFinite(h) && w > 0 && h > 0) {
        return { width: w, height: h };
      }
    }
  }
  return { width: DEFAULT_ASPECT_RATIO_WIDTH, height: DEFAULT_ASPECT_RATIO_HEIGHT };
}

/** @param {number} width @param {AspectRatio} aspect */
export function canvasHeightForWidth(width, aspect) {
  return Math.max(1, Math.round(width * aspect.height / aspect.width));
}

/** @param {number|string|null|undefined} value @param {number} basisPx @param {number} [defaultFraction] */
export function resolveWidth(value, basisPx, defaultFraction = 1) {
  if (value == null || value === '') {
    return Math.round(basisPx * defaultFraction);
  }
  if (typeof value === 'number' && Number.isFinite(value)) {
    return Math.round(value);
  }
  if (typeof value === 'string') {
    const trimmed = value.trim();
    const pctMatch = /^(\d+(?:\.\d+)?)\s*%$/.exec(trimmed);
    if (pctMatch) {
      const pct = parseFloat(pctMatch[1]);
      if (Number.isFinite(pct)) {
        return Math.round(basisPx * (pct / 100));
      }
    }
    const pxMatch = /^(\d+(?:\.\d+)?)\s*px$/i.exec(trimmed);
    if (pxMatch) {
      const designPx = parseFloat(pxMatch[1]);
      if (Number.isFinite(designPx)) {
        return Math.round(designPx * (basisPx / DESIGN_CANVAS_SIZE));
      }
    }
    const asNum = Number(trimmed);
    if (Number.isFinite(asNum)) {
      return Math.round(asNum);
    }
  }
  return Math.round(basisPx * defaultFraction);
}

/** @typedef {number|string} MarginSpec */

/** @param {Record<string, unknown>} overrides @param {'horizontal' | 'vertical'} axis */
function marginSpecForAxis(overrides, axis) {
  const keys = axis === 'horizontal'
    ? ['marginHorizontal', 'marginX']
    : ['marginVertical', 'marginY'];
  for (const key of keys) {
    if (Object.prototype.hasOwnProperty.call(overrides, key)) {
      return /** @type {MarginSpec|null} */ (overrides[key]);
    }
  }
  if (Object.prototype.hasOwnProperty.call(overrides, 'margin')) {
    return /** @type {MarginSpec|null} */ (overrides.margin);
  }
  return null;
}

/** @param {Record<string, unknown>} overrides @param {number} size @param {'horizontal' | 'vertical'} axis */
function resolveMarginPx(overrides, size, axis) {
  const maxPx = Math.round(size * 0.16);
  const spec = marginSpecForAxis(overrides, axis);
  if (spec != null && spec !== '') {
    const px = resolveWidth(spec, size, 0);
    return clamp(px, 0, maxPx);
  }
  const defaultFraction = axis === 'vertical'
    ? DEFAULT_MARGIN_VERTICAL_FRACTION
    : DEFAULT_MARGIN_FRACTION;
  const px = resolveWidth(null, size, defaultFraction);
  return clamp(px, Math.min(64, maxPx), maxPx);
}

/** Margins resolved on the 1080 design canvas (then scaled with {@link scaleCanvasPx}). */
function resolveMarginPxDesign(overrides, axis) {
  return resolveMarginPx(overrides, DESIGN_CANVAS_SIZE, axis);
}

/** @param {DeckTheme} theme */
export function innerWidthFor(theme) {
  return theme.size - theme.marginHorizontal * 2;
}

/** Design px reserved inside the bitmap edge so glyphs are not clipped (scaled to render size). */
export const BITMAP_SAFE_INSET_DESIGN = 20;

/** Inset from the physical canvas edge when the text column touches it (margin 0). */
export function bitmapSafeInsetPx(canvasSize) {
  const scaled = scaleCanvasPx(BITMAP_SAFE_INSET_DESIGN, canvasSize);
  // Small preview rasters round more aggressively, so keep a stronger minimum inset
  // to avoid display serif right-edge cuts.
  return Math.max(8, scaled);
}

/**
 * Extra inset on the aligned edge when `marginHorizontal` is 0 (display serifs overshoot metrics).
 * @param {number} fontSizePx
 * @param {number} canvasSize
 * @param {number} marginHorizontal
 */
export function columnEdgeBleedPx(fontSizePx, canvasSize, marginHorizontal) {
  if (marginHorizontal > 0) return 0;
  const fromFont = Math.max(2, Math.round(fontSizePx * 0.07));
  const fromCanvas = bitmapSafeInsetPx(canvasSize);
  return Math.max(fromFont, fromCanvas);
}

/** @param {DeckTheme} theme */
export function contentMaxWidthPx(theme) {
  return theme.contentMaxWidth ?? innerWidthFor(theme);
}

/**
 * Text column for layout and draw (`marginHorizontal`, `contentMaxWidth`).
 * @param {DeckTheme} theme
 */
export function contentColumnBounds(theme) {
  const inner = innerWidthFor(theme);
  const rawContentWidth = Math.min(contentMaxWidthPx(theme), inner);
  const columnLeft = theme.marginHorizontal + Math.max(0, Math.round((inner - rawContentWidth) / 2));
  const columnRight = columnLeft + rawContentWidth;
  return {
    columnLeft,
    contentWidth: columnRight - columnLeft,
    columnRight,
    canvasSize: theme.size,
    marginHorizontal: theme.marginHorizontal,
    bitmapSafe: bitmapSafeInsetPx(theme.size),
  };
}

/**
 * Deck spec (%, strings) plus live studio palette/typography for `mergeTheme` at render size.
 * Omits resolved `size` and px margins so preview/export rescale correctly.
 * @param {Partial<DeckTheme> & Record<string, unknown>} [deckSpec]
 * @param {Partial<DeckTheme>} [liveTheme]
 */
export function buildThemeRenderOverrides(deckSpec = {}, liveTheme = {}) {
  const spec = deckSpec && typeof deckSpec === 'object' ? deckSpec : {};
  const live = liveTheme && typeof liveTheme === 'object' ? liveTheme : {};
  const {
    marginHorizontal: _mh,
    marginVertical: _mv,
    margin: _marginDrop,
    contentMaxWidth: _cmw,
    size: _size,
    palette: _paletteDrop,
    background: _bg,
    text: _text,
    muted: _muted,
    accent1: _a1,
    accent2: _a2,
    wavePalette: _wavePaletteDrop,
    wavePaletteLinked: _wavePaletteLinked,
    backgroundWave: _backgroundWave,
    ...deckRest
  } = spec;

  const palette = deckPaletteFromTheme(live);
  /** @type {Record<string, unknown>} */
  const out = {
    ...deckRest,
    ...palette,
    wavePaletteLinked: live.wavePaletteLinked ?? spec.wavePaletteLinked,
    displayFont: live.displayFont ?? spec.displayFont,
    bodyFont: live.bodyFont ?? spec.bodyFont,
    fontSizes: live.fontSizes ?? spec.fontSizes,
    lineHeights: live.lineHeights ?? spec.lineHeights,
    emphasisGap: live.emphasisGap ?? spec.emphasisGap,
    emphasisScale: live.emphasisScale ?? spec.emphasisScale,
    previewMaxPx: live.previewMaxPx ?? spec.previewMaxPx,
    ...mergeBackgroundGradientFields(live, spec),
    ...mergeBackgroundWaveFields(live, spec),
  };

  const liveWavePalette = deckWavePaletteFromTheme(live);
  if (liveWavePalette) {
    out.wavePalette = liveWavePalette;
    out.wavePaletteLinked = false;
  } else if (Object.prototype.hasOwnProperty.call(live, 'wavePaletteLinked')) {
    out.wavePaletteLinked = live.wavePaletteLinked !== false;
  } else if (spec.wavePalette && typeof spec.wavePalette === 'object') {
    out.wavePalette = normalizeWavePalette(/** @type {Partial<WavePalette>} */ (spec.wavePalette), palette);
    out.wavePaletteLinked = false;
  }

  if (Object.prototype.hasOwnProperty.call(spec, 'marginHorizontal')) {
    out.marginHorizontal = spec.marginHorizontal;
  } else if (Object.prototype.hasOwnProperty.call(spec, 'marginX')) {
    out.marginHorizontal = spec.marginX;
  }
  if (Object.prototype.hasOwnProperty.call(spec, 'marginVertical')) {
    out.marginVertical = spec.marginVertical;
  } else if (Object.prototype.hasOwnProperty.call(spec, 'marginY')) {
    out.marginVertical = spec.marginY;
  }
  if (Object.prototype.hasOwnProperty.call(spec, 'margin')) {
    out.margin = spec.margin;
  }
  if (Object.prototype.hasOwnProperty.call(spec, 'contentMaxWidth')) {
    out.contentMaxWidth = spec.contentMaxWidth;
  }
  if (Object.prototype.hasOwnProperty.call(spec, 'size')) {
    out.size = spec.size;
  }

  return out;
}

/** @param {DeckTheme} theme */
export function innerHeightFor(theme) {
  const canvasHeight = theme.sizeHeight ?? canvasHeightForWidth(theme.size, {
    width: theme.aspectRatioWidth ?? DEFAULT_ASPECT_RATIO_WIDTH,
    height: theme.aspectRatioHeight ?? DEFAULT_ASPECT_RATIO_HEIGHT,
  });
  return canvasHeight - theme.marginVertical * 2;
}

/** @typedef {Object} DeckPalette
 * @property {string} background
 * @property {string} text
 * @property {string} muted
 * @property {string} accent1
 * @property {string|null} [accent2]
 */

/**
 * @param {Partial<DeckPalette>|null|undefined} input
 * @returns {DeckPalette}
 */
export function normalizePalette(input = {}) {
  const src = input && typeof input === 'object' ? input : {};
  return {
    background: src.background ?? DEFAULT_THEME.background,
    text: src.text ?? DEFAULT_THEME.text,
    muted: src.muted ?? DEFAULT_THEME.muted,
    accent1: src.accent1 ?? DEFAULT_THEME.accent1,
    accent2: src.accent2 ?? DEFAULT_THEME.accent2,
  };
}

/** @param {DeckTheme} theme @returns {DeckPalette} */
export function deckPaletteFromTheme(theme) {
  return {
    background: theme.background,
    text: theme.text,
    muted: theme.muted,
    accent1: theme.accent1,
    accent2: theme.accent2,
  };
}

/** @typedef {Object} WavePalette
 * @property {string} background
 * @property {string} muted
 * @property {string} accent1
 * @property {string|null} [accent2]
 */

/** @typedef {import('./background-panorama.js').PanoramaPalette} PanoramaPalette */

/**
 * @param {Partial<WavePalette>|null|undefined} input
 * @param {Partial<DeckPalette>|null|undefined} [fallback]
 * @returns {WavePalette}
 */
export function normalizeWavePalette(input = {}, fallback = {}) {
  const src = input && typeof input === 'object' ? input : {};
  const fb = fallback && typeof fallback === 'object' ? fallback : {};
  return {
    background: src.background ?? fb.background ?? DEFAULT_THEME.background,
    muted: src.muted ?? fb.muted ?? DEFAULT_THEME.muted,
    accent1: src.accent1 ?? fb.accent1 ?? DEFAULT_THEME.accent1,
    accent2: src.accent2 ?? fb.accent2 ?? DEFAULT_THEME.accent2,
  };
}

/** @param {DeckPalette} textPalette @returns {PanoramaPalette} */
export function panoramaPaletteFromTextPalette(textPalette) {
  return {
    background: textPalette.background,
    accent1: textPalette.accent1,
    accent2: textPalette.accent2 || textPalette.accent1,
    muted: textPalette.muted,
  };
}

/** @param {Partial<DeckTheme> & Record<string, unknown>} theme */
export function isWavePaletteLinked(theme) {
  if (theme.wavePaletteLinked === false) return false;
  if (theme.wavePalette && typeof theme.wavePalette === 'object') return false;
  return true;
}

/**
 * Colors for panoramic wave painting (may differ from text palette when unlinked).
 * @param {Partial<DeckTheme> & Record<string, unknown>} theme
 * @returns {PanoramaPalette}
 */
export function wavePaletteFromTheme(theme) {
  const text = deckPaletteFromTheme(/** @type {DeckTheme} */ (theme));
  if (isWavePaletteLinked(theme)) {
    return panoramaPaletteFromTextPalette(text);
  }
  const wave = theme.wavePalette && typeof theme.wavePalette === 'object'
    ? theme.wavePalette
    : {};
  return panoramaPaletteFromTextPalette(normalizeWavePalette(wave, text));
}

/** @param {Partial<DeckTheme> & Record<string, unknown>} theme @returns {WavePalette|null} */
export function deckWavePaletteFromTheme(theme) {
  if (isWavePaletteLinked(theme)) return null;
  const text = deckPaletteFromTheme(/** @type {DeckTheme} */ (theme));
  const wave = theme.wavePalette && typeof theme.wavePalette === 'object'
    ? theme.wavePalette
    : {};
  return normalizeWavePalette(wave, text);
}

/** @type {(keyof DeckPalette)[]} */
const PALETTE_COLOR_KEYS = ['background', 'text', 'muted', 'accent1', 'accent2'];

/**
 * Read palette from nested `palette` or flat color keys (already-merged studio theme).
 * @param {Record<string, unknown>} overrides
 * @returns {DeckPalette}
 */
function extractPaletteFromOverrides(overrides) {
  if (overrides.palette && typeof overrides.palette === 'object') {
    return normalizePalette(/** @type {Partial<DeckPalette>} */ (overrides.palette));
  }
  if (PALETTE_COLOR_KEYS.some((key) => Object.prototype.hasOwnProperty.call(overrides, key))) {
    return normalizePalette({
      background: /** @type {string|undefined} */ (overrides.background),
      text: /** @type {string|undefined} */ (overrides.text),
      muted: /** @type {string|undefined} */ (overrides.muted),
      accent1: /** @type {string|undefined} */ (overrides.accent1),
      accent2: /** @type {string|undefined} */ (overrides.accent2),
    });
  }
  return normalizePalette({});
}

/**
 * @typedef {Object} MergeThemeOptions
 * @property {number} [canvasSizeMax] Upper bound for `theme.size` (defaults to `SLIDE_SIZE_MAX`). Export supersample passes `renderSize` here so layout fills the raster buffer.
 */

/** @param {Partial<DeckTheme> & Record<string, unknown>} [overrides] @param {MergeThemeOptions} [options] */
export function mergeTheme(overrides = {}, options = {}) {
  const palette = extractPaletteFromOverrides(overrides);
  const { palette: _paletteDrop, wavePalette: _wavePaletteDrop, wavePaletteLinked: _wpl, background, text, muted, accent1, accent2, ...deckRest } = overrides;
  const theme = { ...DEFAULT_THEME, ...palette, ...deckRest };
  if (typeof background === 'string') theme.background = background;
  if (typeof text === 'string') theme.text = text;
  if (typeof muted === 'string') theme.muted = muted;
  if (typeof accent1 === 'string') theme.accent1 = accent1;
  if (typeof accent2 === 'string') theme.accent2 = accent2;
  if (overrides.wavePalette && typeof overrides.wavePalette === 'object') {
    theme.wavePalette = normalizeWavePalette(/** @type {Partial<WavePalette>} */ (overrides.wavePalette), palette);
    theme.wavePaletteLinked = false;
  } else {
    theme.wavePalette = null;
    theme.wavePaletteLinked = overrides.wavePaletteLinked !== false;
  }
  const sizeMax = Number.isFinite(options.canvasSizeMax)
    ? Math.max(SLIDE_SIZE_MAX, Math.round(options.canvasSizeMax))
    : SLIDE_SIZE_MAX;
  theme.size = clamp(Math.round(theme.size), SLIDE_SIZE_MIN, sizeMax);
  const aspect = parseAspectRatio(overrides.aspectRatio);
  theme.aspectRatioWidth = aspect.width;
  theme.aspectRatioHeight = aspect.height;
  theme.sizeHeight = canvasHeightForWidth(theme.size, aspect);

  const marginHorizontalDesign = resolveMarginPxDesign(overrides, 'horizontal');
  const marginVerticalDesign = resolveMarginPxDesign(overrides, 'vertical');
  theme.marginHorizontal = scaleCanvasPx(marginHorizontalDesign, theme.size);
  theme.marginVertical = scaleCanvasPx(marginVerticalDesign, theme.size);

  theme.fontSizes = {
    ...DEFAULT_THEME.fontSizes,
    ...(/** @type {Record<string, number>} */ (overrides.fontSizes || theme.fontSizes || {})),
  };
  theme.emphasisScale = {
    ...DEFAULT_THEME.emphasisScale,
    ...(/** @type {Record<string, number>} */ (overrides.emphasisScale || theme.emphasisScale || {})),
  };
  theme.lineHeights = {
    ...DEFAULT_THEME.lineHeights,
    ...(/** @type {Record<string, number>} */ (overrides.lineHeights || theme.lineHeights || {})),
  };

  const innerWidth = innerWidthFor(theme);
  const innerWidthDesign = DESIGN_CANVAS_SIZE - marginHorizontalDesign * 2;
  const contentSpec = Object.prototype.hasOwnProperty.call(overrides, 'contentMaxWidth')
    ? overrides.contentMaxWidth
    : null;
  const contentWidthDesign = resolveWidth(contentSpec, innerWidthDesign, 0.88);
  const minContentWidth = Math.min(scaleCanvasPx(420, theme.size), innerWidth);
  theme.contentMaxWidth = clamp(
    scaleCanvasPx(contentWidthDesign, theme.size),
    minContentWidth,
    innerWidth,
  );
  theme.previewMaxPx = clamp(
    Math.round(theme.previewMaxPx ?? DEFAULT_PREVIEW_MAX_PX),
    100,
    480,
  );

  const emphasisGapDesign = Number.isFinite(overrides.emphasisGap)
    ? overrides.emphasisGap
    : (theme.emphasisGap ?? DEFAULT_EMPHASIS_GAP);
  theme.emphasisGap = clamp(Number(emphasisGapDesign) || DEFAULT_EMPHASIS_GAP, 0, 3);

  theme.backgroundGradientPreset = resolveBackgroundGradientPreset(overrides, theme);
  theme.backgroundGradient = buildThemeBackgroundGradient(theme);
  theme.backgroundGradientMode = parseBackgroundGradientMode(overrides);
  theme.backgroundWave = parseBackgroundWaveConfig(overrides);

  return theme;
}

/**
 * Live studio wins over deck spec; explicit `null` preset means solid (not `??` fallback).
 * @param {Partial<DeckTheme> & Record<string, unknown>} live
 * @param {Partial<DeckTheme> & Record<string, unknown>} spec
 */
function mergeBackgroundGradientFields(live, spec) {
  if (Object.prototype.hasOwnProperty.call(live, 'backgroundGradientPreset')) {
    const preset = live.backgroundGradientPreset;
    if (preset == null) {
      return { backgroundGradientPreset: null, backgroundGradient: 'solid' };
    }
    if (typeof live.backgroundGradient === 'string') {
      return { backgroundGradientPreset: preset, backgroundGradient: live.backgroundGradient };
    }
    return { backgroundGradientPreset: preset, backgroundGradient: preset };
  }
  if (Object.prototype.hasOwnProperty.call(live, 'backgroundGradient')) {
    return { backgroundGradient: live.backgroundGradient };
  }
  return {
    backgroundGradientPreset: spec.backgroundGradientPreset,
    backgroundGradient: spec.backgroundGradient,
  };
}

/**
 * Live studio wave style wins over deck spec (same precedence as palette colors).
 * @param {Partial<DeckTheme> & Record<string, unknown>} live
 * @param {Partial<DeckTheme> & Record<string, unknown>} spec
 */
function mergeBackgroundWaveFields(live, spec) {
  if (Object.prototype.hasOwnProperty.call(live, 'backgroundWave') && live.backgroundWave != null) {
    return { backgroundWave: live.backgroundWave };
  }
  if (spec.backgroundWave != null) {
    return { backgroundWave: spec.backgroundWave };
  }
  return {};
}

/** @param {Record<string, unknown>} overrides @param {DeckTheme} theme */
function resolveBackgroundGradientPreset(overrides, theme) {
  if (Object.prototype.hasOwnProperty.call(overrides, 'backgroundGradientPreset')) {
    return normalizeGradientPresetFromInput(overrides.backgroundGradientPreset);
  }
  if (Object.prototype.hasOwnProperty.call(overrides, 'backgroundGradient')) {
    const raw = overrides.backgroundGradient;
    if (typeof raw === 'string') {
      return normalizeGradientPresetFromInput(raw);
    }
  }
  if (theme.backgroundGradientPreset != null) {
    return normalizeGradientPresetFromInput(theme.backgroundGradientPreset);
  }
  return null;
}

/** @param {unknown} value */
function normalizeGradientPresetFromInput(value) {
  if (value == null || value === '' || value === 'none' || value === 'solid') return null;
  if (typeof value === 'string') return normalizeGradientPresetId(value);
  return null;
}

/** @param {DeckTheme} theme */
export function rebuildBackgroundGradient(theme) {
  const preset = resolveBackgroundGradientPreset(theme, theme);
  theme.backgroundGradientPreset = preset;
  theme.backgroundGradient = buildBackgroundGradient(
    { background: theme.background, accent1: theme.accent1, accent2: theme.accent2 },
    preset,
  );
}

/** @param {DeckTheme} theme */
function buildThemeBackgroundGradient(theme) {
  return buildBackgroundGradient(
    { background: theme.background, accent1: theme.accent1, accent2: theme.accent2 },
    theme.backgroundGradientPreset,
  );
}

/** @param {DeckTheme} theme @param {ColorToken|string} token */
export function resolveColor(theme, token) {
  if (typeof token === 'string' && token.startsWith('#')) {
    return token;
  }
  switch (token) {
    case 'muted':
      return theme.muted;
    case 'accent1':
      return theme.accent1;
    case 'accent2':
      return theme.accent2 || theme.accent1;
    case 'text':
    default:
      return theme.text;
  }
}

/**
 * Brightness offset **-100** to **+100** (`"+55"`, `55`, or `undefined`).
 * @param {unknown} raw
 * @returns {number|undefined}
 */
export function coerceBrightness(raw) {
  if (raw == null || raw === '') return undefined;
  if (typeof raw === 'number') {
    if (!Number.isFinite(raw)) return undefined;
    return Math.max(-100, Math.min(100, raw));
  }
  const value = Number.parseInt(String(raw), 10);
  if (!Number.isFinite(value)) return undefined;
  return Math.max(-100, Math.min(100, value));
}

const COLOR_TAG_RE =
  /^<(?:color\s+)?(accent1|accent2|text|muted)(?:\s+brightness=["']([+-]?\d{1,3})["'])?\s*>/i;

/**
 * Palette token, `#hex`, `{ color, brightness }`, or the same tag as slide text
 * (`<accent2 brightness='+55'>`).
 * @param {unknown} spec
 * @returns {{ color?: string, brightness?: number }|null}
 */
export function parseThemedColorSpec(spec) {
  if (spec == null || spec === '') return null;
  if (typeof spec === 'object' && !Array.isArray(spec)) {
    const record = /** @type {{ color?: unknown, brightness?: unknown }} */ (spec);
    return {
      color: typeof record.color === 'string' ? record.color : undefined,
      brightness: coerceBrightness(record.brightness),
    };
  }
  if (typeof spec !== 'string') return null;
  const trimmed = spec.trim();
  const tag = COLOR_TAG_RE.exec(trimmed);
  if (tag) {
    return {
      color: tag[1].toLowerCase(),
      brightness: coerceBrightness(tag[2]),
    };
  }
  return { color: trimmed };
}

/**
 * Palette token (or `#hex`) for an inline run, with optional brightness offset.
 * @param {DeckTheme} theme
 * @param {{ color?: string, brightness?: number|string }} run
 * @param {ColorToken|string} [blockColor]
 */
export function resolveInlineRunColor(theme, run, blockColor = 'text') {
  const base = resolveColor(theme, run.color ?? blockColor);
  const brightness = coerceBrightness(run.brightness);
  if (brightness === undefined || brightness === 0) {
    return base;
  }
  const amount = brightness / 100;
  return shiftHex(base, amount);
}

/**
 * Palette token, `#hex`, `{ color, brightness }`, or `<accent2 brightness='+55'>`.
 * @param {DeckTheme} theme
 * @param {string|{color?: string, brightness?: number|string}|null|undefined} spec
 * @param {ColorToken|string} [fallbackToken]
 */
export function resolveThemedColor(theme, spec, fallbackToken = 'text') {
  const parsed = parseThemedColorSpec(spec);
  if (!parsed) {
    return resolveColor(theme, fallbackToken);
  }
  return resolveInlineRunColor(
    theme,
    { color: parsed.color ?? fallbackToken, brightness: parsed.brightness },
    fallbackToken,
  );
}

/** Scale a distance measured on the 1080 design canvas to the render/export canvas. */
export function scaleCanvasPx(designPx, canvasSize) {
  return Math.round(designPx * (canvasSize / DESIGN_CANVAS_SIZE));
}

/** @param {number} px @param {number} canvasSize */
export function scaleFontSize(px, canvasSize) {
  return scaleCanvasPx(px, canvasSize);
}

/** @param {DeckTheme} theme @param {Section} section */
export function sectionBaseDesignPx(theme, section) {
  const map = theme.fontSizes || DEFAULT_THEME.fontSizes;
  if (section === 'header') return map.header ?? 60;
  if (section === 'footer') return map.footer ?? 52;
  return map.body ?? 80;
}

/** @param {DeckTheme} theme @param {BodyEmphasis} [emphasis] */
export function emphasisScaleFor(theme, emphasis) {
  const map = theme.emphasisScale || DEFAULT_THEME.emphasisScale;
  if (emphasis === 'punch') return map.punch ?? DEFAULT_EMPHASIS_SCALE_PUNCH;
  return map.normal ?? 1;
}

/**
 * @param {{ section: Section, emphasis?: BodyEmphasis, fontSize?: number }} block
 * @param {DeckTheme} theme
 */
export function resolveDesignFontSize(block, theme) {
  if (block.fontSize != null && Number.isFinite(block.fontSize)) {
    return block.fontSize;
  }
  const base = sectionBaseDesignPx(theme, block.section);
  if (block.section === 'body') {
    return Math.round(base * emphasisScaleFor(theme, block.emphasis ?? 'normal'));
  }
  return base;
}

/** @param {DeckTheme} theme */
export function resolveEmphasisGapPx(theme) {
  const fraction = theme.emphasisGap ?? DEFAULT_EMPHASIS_GAP;
  const bodyDesign = theme.fontSizes?.body ?? 80;
  return Math.round(scaleFontSize(bodyDesign, theme.size) * fraction);
}

/** @param {DeckTheme} theme @param {number} designPx */
export function resolveFontSizePx(theme, designPx) {
  return scaleFontSize(designPx, theme.size);
}

/**
 * @param {DeckTheme} theme
 * @param {{ section: Section, emphasis?: BodyEmphasis, lineHeight?: number }} block
 */
export function resolveLineHeightMult(theme, block) {
  if (block.lineHeight != null && Number.isFinite(block.lineHeight)) {
    return block.lineHeight;
  }
  const map = theme.lineHeights || DEFAULT_THEME.lineHeights;
  if (block.section === 'body') {
    const key = block.emphasis === 'punch' ? 'punch' : 'normal';
    if (map[key] != null) return map[key];
  }
  if (map[block.section] != null) return map[block.section];
  return 1;
}

/** @param {import('./theme.js').DeckTheme} theme @param {FontRole} role @param {number} px @param {string} [weight] */
export function buildFontPx(theme, role, px, weight = '400') {
  const family = role === 'display' ? theme.displayFont : theme.bodyFont;
  return `${weight} ${px}px "${family}", serif`;
}

/** @type {string|null} */
let loadedFontsCacheKey = null;
/** @type {Promise<void>|null} */
let fontsLoadingPromise = null;
/** @type {string|null} */
let fontsLoadingKey = null;

/** @param {string} href */
async function ensureFontStylesheet(href) {
  let link = document.querySelector(`link[data-carousel-fonts="${href}"]`);
  if (!link) {
    link = document.createElement('link');
    link.rel = 'stylesheet';
    link.href = href;
    link.dataset.carouselFonts = href;
    document.head.appendChild(link);
  }

  if (link.sheet) return link;

  await new Promise((resolve, reject) => {
    link.addEventListener('load', () => resolve(), { once: true });
    link.addEventListener('error', () => reject(new Error('Carousel font stylesheet failed to load')), { once: true });
  });

  return link;
}

/** @param {DeckTheme} theme @returns {number[]} */
function deckFontLoadSizes(theme) {
  /** @type {Set<number>} */
  const sizes = new Set();
  const header = theme.fontSizes?.header ?? 60;
  const body = theme.fontSizes?.body ?? 80;
  const footer = theme.fontSizes?.footer ?? 52;
  const punch = Math.round(body * emphasisScaleFor(theme, 'punch'));
  const minPx = scaleFontSize(30, theme.size);
  const maxDesign = Math.max(header, body, footer, punch);
  const maxPx = scaleFontSize(maxDesign, theme.size);
  const step = Math.max(1, Math.round(scaleFontSize(4, theme.size)));

  for (let px = minPx; px <= maxPx; px += step) {
    sizes.add(px);
  }
  sizes.add(maxPx);
  sizes.add(minPx);

  return [...sizes].sort((a, b) => a - b);
}

/** @param {DeckTheme} theme @param {number[]} sizes */
function deckFontDescriptors(theme, sizes) {
  /** @type {string[]} */
  const descriptors = [];
  for (const px of sizes) {
    for (const weight of ['400', '700']) {
      descriptors.push(`${weight} ${px}px "${theme.displayFont}"`);
      descriptors.push(`italic ${weight} ${px}px "${theme.displayFont}"`);
    }
    for (const weight of ['400', '600']) {
      descriptors.push(`${weight} ${px}px "${theme.bodyFont}"`);
      descriptors.push(`italic ${weight} ${px}px "${theme.bodyFont}"`);
    }
  }
  return descriptors;
}

/** @param {string} descriptor */
async function loadFontDescriptor(descriptor) {
  try {
    await document.fonts.load(descriptor);
  } catch {
    // load() can reject; verification pass below confirms availability.
  }
}

/** @param {string[]} descriptors */
async function verifyFontDescriptors(descriptors) {
  await document.fonts.ready;

  for (let attempt = 0; attempt < 24; attempt += 1) {
    const missing = descriptors.filter((descriptor) => !document.fonts.check(descriptor));
    if (missing.length === 0) return;
    await Promise.all(missing.map((descriptor) => loadFontDescriptor(descriptor)));
    await document.fonts.ready;
    await new Promise((resolve) => setTimeout(resolve, 50));
  }

  const stillMissing = descriptors.filter((descriptor) => !document.fonts.check(descriptor));
  if (stillMissing.length > 0) {
    console.warn(
      `Carousel fonts not ready (${stillMissing.length} descriptors)`,
      stillMissing.slice(0, 3),
    );
  }
}

/** @param {DeckTheme} theme */
async function loadDeckFontsNow(theme) {
  const families = [
    `${theme.displayFont}:ital,wght@0,400;0,700;1,400;1,700`,
    `${theme.bodyFont}:ital,wght@0,400;0,600;1,400;1,600`,
  ];
  const href = `https://fonts.googleapis.com/css2?${families.map((f) => `family=${f.replace(/ /g, '+')}`).join('&')}&display=swap`;

  await ensureFontStylesheet(href);

  const sizes = deckFontLoadSizes(theme);
  const descriptors = deckFontDescriptors(theme, sizes);

  const batchSize = 32;
  for (let i = 0; i < descriptors.length; i += batchSize) {
    await Promise.all(descriptors.slice(i, i + batchSize).map((descriptor) => loadFontDescriptor(descriptor)));
  }
  await verifyFontDescriptors(descriptors);
  await document.fonts.ready;

  loadedFontsCacheKey = `${href}|${sizes.join(',')}|${theme.displayFont}|${theme.bodyFont}`;
}

/** @param {DeckTheme} theme */
export async function loadDeckFonts(theme) {
  const families = [
    `${theme.displayFont}:ital,wght@0,400;0,700;1,400;1,700`,
    `${theme.bodyFont}:ital,wght@0,400;0,600;1,400;1,600`,
  ];
  const href = `https://fonts.googleapis.com/css2?${families.map((f) => `family=${f.replace(/ /g, '+')}`).join('&')}&display=swap`;
  const sizes = deckFontLoadSizes(theme);
  const cacheKey = `${href}|${sizes.join(',')}|${theme.displayFont}|${theme.bodyFont}`;

  if (loadedFontsCacheKey === cacheKey) {
    await document.fonts.ready;
    return;
  }

  if (fontsLoadingPromise && fontsLoadingKey === cacheKey) {
    await fontsLoadingPromise;
    return;
  }

  fontsLoadingKey = cacheKey;
  fontsLoadingPromise = loadDeckFontsNow(theme);

  try {
    await fontsLoadingPromise;
  } catch (error) {
    loadedFontsCacheKey = null;
    throw error;
  } finally {
    if (fontsLoadingKey === cacheKey) {
      fontsLoadingPromise = null;
      fontsLoadingKey = null;
    }
  }
}
