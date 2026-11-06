import { loadImage, resolveAssetUrl } from './assets.js';
import { CAROUSEL_SLIDE_WIDTH_PX, isCarouselCtaRole } from './slide-constants.js';
import { resolveColor, resolveWidth, scaleCanvasPx } from './theme.js';

/**
 * @typedef {import('./background-panorama.js').BackgroundPanoramaContext} MotifStripContext
 */

/**
 * @typedef {Object} MotifStripConfig
 * @property {boolean} [enabled] When false, motif strip is not loaded or drawn (default true).
 * @property {string} src Bundle-relative SVG or raster strip (one slide width per deck slide).
 * @property {number} [slideWidth] Width of one slide slice in strip asset coords (default 1080).
 * @property {number} [stripHeight] Height of strip asset coords (default from image or 300).
 * @property {string} [anchor] `bottom` (default) or `top`.
 * @property {string|number} [marginBottom] Gap above the slide bottom edge (`3%` or px at 1080).
 * @property {string|number} [marginTop] Gap below the slide top edge when `anchor` is `top`.
 * @property {number} [offsetX] Horizontal nudge in design px at 1080 (positive = right), applied once on the seamless band.
 * @property {number} [offsetY] Vertical nudge in design px at 1080 (positive = down).
 * @property {string|number} [inset] Deprecated alias for `marginBottom` (or `marginTop` when anchored top).
 * @property {string|number} [bandWidth] Uniform size of the panoramic motif (`100%` default). Percent scales width and height together from the bottom-left; the strip stays continuous across slides.
 * @property {string|number} [bandHeight] Optional height override; when omitted, height follows `bandWidth` (locked aspect).
 * @property {string} [color] Palette token or `#hex` for SVG tint (default `accent1`).
 * @property {string} [keyColor] Raster only: knock out this color (e.g. `#000000` export backdrop).
 * @property {number} [keyTolerance] Per-channel tolerance for `keyColor` (default 32).
 * @property {number} [opacity] 0–1 (default 1).
 * @property {string[]} [excludeRoles] Skip motif on slides with these roles.
 */

/**
 * @typedef {Object} MotifSeamlessBand
 * @property {HTMLCanvasElement} band
 * @property {number} pad Horizontal padding baked into the band for negative offsetX
 * @property {number} destHeight
 * @property {number} destWidth Output width of one slide crop (full canvas width)
 * @property {number} destY
 * @property {number} slideCanvasWidth Output slide width
 * @property {number} bandSlideWidth Width of one slide slice inside the seamless band
 */

/** @type {Map<string, HTMLImageElement|ImageBitmap>} */
const motifStripCache = new Map();

/** @type {Map<string, MotifSeamlessBand>} */
const motifBandCache = new Map();
const MOTIF_BAND_CACHE_MAX = 32;

/**
 * @param {unknown} raw
 * @returns {MotifStripConfig|null}
 */
export function parseMotifStripConfig(raw) {
  if (!raw || typeof raw !== 'object' || Array.isArray(raw)) return null;
  const spec = /** @type {Record<string, unknown>} */ (raw);
  const src = typeof spec.src === 'string' ? spec.src.trim() : '';
  if (!src) return null;
  if (spec.enabled === false) return null;
  return {
    src,
    enabled: true,
    slideWidth: Number.isFinite(Number(spec.slideWidth)) ? Number(spec.slideWidth) : CAROUSEL_SLIDE_WIDTH_PX,
    stripHeight: Number.isFinite(Number(spec.stripHeight)) ? Number(spec.stripHeight) : null,
    anchor: spec.anchor === 'top' ? 'top' : 'bottom',
    marginBottom: spec.marginBottom ?? (spec.anchor === 'top' ? null : (spec.inset ?? '3%')),
    marginTop: spec.marginTop ?? (spec.anchor === 'top' ? (spec.inset ?? '3%') : null),
    offsetX: Number.isFinite(Number(spec.offsetX)) ? Number(spec.offsetX) : 0,
    offsetY: Number.isFinite(Number(spec.offsetY)) ? Number(spec.offsetY) : 0,
    bandWidth: spec.bandWidth ?? '100%',
    bandHeight: spec.bandHeight != null && spec.bandHeight !== '' ? spec.bandHeight : null,
    color: typeof spec.color === 'string' ? spec.color : 'accent1',
    keyColor: typeof spec.keyColor === 'string' ? spec.keyColor.trim() : '',
    keyTolerance: Number.isFinite(Number(spec.keyTolerance)) ? Number(spec.keyTolerance) : 32,
    opacity: Number.isFinite(Number(spec.opacity)) ? Number(spec.opacity) : 1,
    excludeRoles: Array.isArray(spec.excludeRoles)
      ? spec.excludeRoles.map((role) => String(role).trim().toLowerCase()).filter(Boolean)
      : [],
  };
}

/**
 * Motif panorama context: CTA slides are omitted so adding a QR frame does not stretch the strip.
 * @param {{ slides?: { number: number, role?: string }[], deck?: Record<string, unknown> }|null|undefined} deck
 * @param {number} slideNumber
 * @returns {import('./background-panorama.js').BackgroundPanoramaContext|null}
 */
export function buildMotifStripContext(deck, slideNumber) {
  if (!deck?.slides?.length) return null;
  const config = parseMotifStripConfig(deck.deck?.motifStrip);
  const eligible = deck.slides.filter((slide) => motifStripEnabledForSlide(config, slide.role));
  const index = eligible.findIndex((slide) => slide.number === slideNumber);
  if (index < 0) return null;
  return { slideIndex: index, slideCount: Math.max(1, eligible.length) };
}

/**
 * @param {string} svg
 */
function ensureSvgPixelSize(svg) {
  if (/\bwidth="/i.test(svg) && /\bheight="/i.test(svg)) return svg;
  const viewBoxMatch = svg.match(/viewBox="\s*([\d.]+)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)\s*"/i);
  if (!viewBoxMatch) return svg;
  const width = viewBoxMatch[3];
  const height = viewBoxMatch[4];
  return svg.replace(/<svg\b/i, `<svg width="${width}" height="${height}"`);
}

async function loadTintedMotifSvg(url, fillHex) {
  const response = await fetch(url, { cache: 'no-store' });
  if (!response.ok) {
    throw new Error(`Failed to load SVG: ${url}`);
  }
  let svg = await response.text();
  svg = svg.replace(/fill="currentColor"/gi, `fill="${fillHex}"`);
  svg = svg.replace(/stroke="currentColor"/gi, `stroke="${fillHex}"`);
  svg = ensureSvgPixelSize(svg);
  const blob = new Blob([svg], { type: 'image/svg+xml' });
  const objectUrl = URL.createObjectURL(blob);
  try {
    return await new Promise((resolve, reject) => {
      const el = new Image();
      el.decoding = 'async';
      el.onload = () => resolve(el);
      el.onerror = () => reject(new Error(`Failed to decode SVG: ${url}`));
      el.src = objectUrl;
    });
  } finally {
    URL.revokeObjectURL(objectUrl);
  }
}

/** @param {string} hex */
function parseHexRgb(hex) {
  if (typeof hex !== 'string') return { r: 0, g: 0, b: 0 };
  let raw = hex.replace('#', '').trim();
  if (raw.length === 3) {
    raw = raw[0] + raw[0] + raw[1] + raw[1] + raw[2] + raw[2];
  }
  if (raw.length !== 6) return { r: 0, g: 0, b: 0 };
  return {
    r: parseInt(raw.slice(0, 2), 16),
    g: parseInt(raw.slice(2, 4), 16),
    b: parseInt(raw.slice(4, 6), 16),
  };
}

/**
 * @param {HTMLImageElement|ImageBitmap} image
 * @param {string} keyHex
 * @param {number} tolerance
 */
function applyColorKey(image, keyHex, tolerance) {
  if (!image) return image;
  const width = image.naturalWidth || image.width;
  const height = image.naturalHeight || image.height;
  if (!width || !height) return image;

  const canvas = document.createElement('canvas');
  canvas.width = width;
  canvas.height = height;
  const ctx = canvas.getContext('2d');
  if (!ctx) return image;

  ctx.drawImage(image, 0, 0, width, height);
  const { r: kr, g: kg, b: kb } = parseHexRgb(keyHex);
  const tol = Math.max(0, Math.min(255, tolerance));
  const data = ctx.getImageData(0, 0, width, height);
  for (let i = 0; i < data.data.length; i += 4) {
    const r = data.data[i];
    const g = data.data[i + 1];
    const b = data.data[i + 2];
    if (Math.abs(r - kr) <= tol && Math.abs(g - kg) <= tol && Math.abs(b - kb) <= tol) {
      data.data[i + 3] = 0;
    }
  }
  ctx.putImageData(data, 0, 0);
  return canvas;
}

/**
 * @param {MotifStripConfig} config
 * @param {string} assetBaseUrl
 * @param {(token: string) => string} resolveColorFn
 */
export async function loadMotifStripImage(config, assetBaseUrl, resolveColorFn) {
  const url = resolveAssetUrl(config.src, assetBaseUrl);
  const fill = config.color?.startsWith('#')
    ? config.color
    : resolveColorFn(config.color || 'accent1');
  const keyColor = config.keyColor
    ? (config.keyColor.startsWith('#') ? config.keyColor : resolveColorFn(config.keyColor))
    : '';
  const cacheKey = `${url}#${fill}#${keyColor}#${config.keyTolerance ?? 32}`;
  const cached = motifStripCache.get(cacheKey);
  if (cached) return cached;

  const isSvg = /\.svg(?:$|\?)/i.test(url);
  let image = isSvg
    ? await loadTintedMotifSvg(url, fill)
    : await loadImage(url);
  if (!isSvg && keyColor) {
    image = applyColorKey(image, keyColor, config.keyTolerance ?? 32);
  }
  motifStripCache.set(cacheKey, image);
  return image;
}

/**
 * @param {HTMLImageElement|ImageBitmap|HTMLCanvasElement} image
 * @param {MotifStripConfig} config
 * @param {number} slideCount
 * @param {number} canvasWidth
 * @param {number} canvasHeight
 */
function motifBandCacheKey(image, config, slideCount, canvasWidth, canvasHeight) {
  const naturalW = image.naturalWidth || image.width || 0;
  const naturalH = image.naturalHeight || image.height || 0;
  return [
    naturalW,
    naturalH,
    slideCount,
    canvasWidth,
    canvasHeight,
    config.offsetX ?? 0,
    config.offsetY ?? 0,
    config.bandWidth ?? '',
    config.bandHeight ?? '',
    config.marginBottom ?? '',
    config.marginTop ?? '',
    config.anchor ?? '',
    config.slideWidth ?? '',
    config.stripHeight ?? '',
  ].join('|');
}

/**
 * Size percent from `bandWidth` (`80%` → 0.8). Applied to width and height together.
 * @param {string|number|null|undefined} bandWidth
 */
function motifBandScale(bandWidth) {
  if (bandWidth == null || bandWidth === '') return 1;
  if (typeof bandWidth === 'string') {
    const trimmed = bandWidth.trim();
    const pct = /^(\d+(?:\.\d+)?)\s*%$/.exec(trimmed);
    if (pct) {
      return Math.max(0.2, Math.min(2, Number(pct[1]) / 100));
    }
  }
  const asNum = Number(bandWidth);
  if (Number.isFinite(asNum) && asNum > 0 && asNum <= 200) {
    return Math.max(0.2, Math.min(2, asNum / 100));
  }
  const px = resolveWidth(bandWidth, CAROUSEL_SLIDE_WIDTH_PX, 1);
  return Math.max(0.2, Math.min(2, px / CAROUSEL_SLIDE_WIDTH_PX));
}

/**
 * Build one continuous motif band for the full deck. Size scales from the bottom-left
 * (left edge + bottom anchor); offsetX pans the band so slide seams stay aligned.
 *
 * @param {HTMLImageElement|ImageBitmap|HTMLCanvasElement} image
 * @param {MotifStripConfig} config
 * @param {number} slideCount
 * @param {number} canvasWidth
 * @param {number} canvasHeight
 * @returns {MotifSeamlessBand}
 */
function buildMotifSeamlessBand(image, config, slideCount, canvasWidth, canvasHeight) {
  const key = motifBandCacheKey(image, config, slideCount, canvasWidth, canvasHeight);
  const cached = motifBandCache.get(key);
  if (cached) return cached;

  const slices = Math.max(1, slideCount);
  const naturalW = image.naturalWidth || image.width || 1;
  const naturalH = image.naturalHeight || image.height || 1;
  const stripHeight = config.stripHeight ?? naturalH;

  const configuredSpan = (config.slideWidth ?? CAROUSEL_SLIDE_WIDTH_PX) * slices;
  const useEqualSlices = Math.abs(configuredSpan - naturalW) > 2;
  const sliceWidth = useEqualSlices ? naturalW / slices : (config.slideWidth ?? CAROUSEL_SLIDE_WIDTH_PX);
  const scaleX = useEqualSlices ? 1 : naturalW / configuredSpan;
  const sw = sliceWidth * scaleX;
  const scaleY = stripHeight > 0 ? naturalH / stripHeight : 1;
  const sh = stripHeight * scaleY;
  const scale = motifBandScale(config.bandWidth);

  const destWidth = canvasWidth;
  const destHeight = config.bandHeight != null && config.bandHeight !== ''
    ? Math.max(8, Math.round(resolveWidth(config.bandHeight, canvasHeight, 0.22)))
    : Math.max(1, Math.round(destWidth * (sh / Math.max(1, sw)) * scale));

  const anchor = config.anchor === 'top' ? 'top' : 'bottom';
  const marginSpec = anchor === 'bottom'
    ? (config.marginBottom ?? '3%')
    : (config.marginTop ?? '3%');
  const margin = Math.round(resolveWidth(marginSpec, canvasHeight, 0.03));
  const baseY = anchor === 'bottom'
    ? canvasHeight - margin - destHeight
    : margin;
  const destY = baseY + scaleCanvasPx(config.offsetY ?? 0, canvasWidth);

  const bandSlideWidth = destWidth;
  const deckWidth = slices * bandSlideWidth;
  const contentWidth = Math.max(1, Math.round(deckWidth * scale));
  const nudgeX = scaleCanvasPx(config.offsetX ?? 0, canvasWidth);
  const contentX = nudgeX;
  const pad = Math.max(0, Math.ceil(-contentX));
  const padRight = Math.max(0, Math.ceil(contentX + contentWidth - deckWidth));

  const band = document.createElement('canvas');
  band.width = pad + deckWidth + padRight;
  band.height = destHeight;
  const bctx = band.getContext('2d');
  if (!bctx) throw new Error('Canvas 2D context unavailable');

  bctx.imageSmoothingEnabled = true;
  bctx.imageSmoothingQuality = 'high';
  bctx.drawImage(
    image,
    0,
    0,
    naturalW,
    naturalH,
    pad + contentX,
    0,
    contentWidth,
    destHeight,
  );

  /** @type {MotifSeamlessBand} */
  const entry = {
    band,
    pad,
    destHeight,
    destWidth,
    destY,
    slideCanvasWidth: canvasWidth,
    bandSlideWidth,
  };
  motifBandCache.set(key, entry);
  while (motifBandCache.size > MOTIF_BAND_CACHE_MAX) {
    const oldest = motifBandCache.keys().next().value;
    if (oldest) motifBandCache.delete(oldest);
  }
  return entry;
}

/**
 * @param {MotifSeamlessBand} seamless
 * @param {number} slideIndex
 */
function bandSliceRect(seamless, slideIndex) {
  const bandStride = seamless.bandSlideWidth ?? seamless.destWidth;
  return {
    sx: seamless.pad + slideIndex * bandStride,
    sy: 0,
    sw: bandStride,
    sh: seamless.destHeight,
    dx: 0,
    dy: seamless.destY,
    dw: bandStride,
    dh: seamless.destHeight,
  };
}

/**
 * Paint motif slices across a multi-slide strip (studio panorama preview).
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} slideWidth
 * @param {number} height
 * @param {number} deckSlideCount
 * @param {MotifStripConfig} config
 * @param {HTMLImageElement|ImageBitmap} image
 * @param {import('./theme.js').DeckTheme} theme
 * @param {{ gap?: number, startX?: number, slotSlideIndices?: number[], slideRoles?: string[] }} [layout]
 */
export function paintPanoramicMotifStrip(
  ctx,
  slideWidth,
  height,
  deckSlideCount,
  config,
  image,
  theme,
  layout = {},
) {
  const gap = layout.gap ?? 0;
  const startX = layout.startX ?? 0;
  const slotSlideIndices = layout.slotSlideIndices;
  const slideRoles = layout.slideRoles ?? [];
  const slotCount = slotSlideIndices?.length ?? deckSlideCount;
  const seamless = buildMotifSeamlessBand(image, config, Math.max(1, deckSlideCount), slideWidth, height);

  ctx.save();
  ctx.globalAlpha = Math.max(0, Math.min(1, config.opacity ?? 1));
  ctx.imageSmoothingEnabled = false;
  ctx.imageSmoothingQuality = 'high';

  for (let slot = 0; slot < slotCount; slot += 1) {
    const slideIndex = slotSlideIndices?.[slot] ?? slot;
    const role = slideRoles[slot] ?? '';
    if (slideIndex < 0 || !motifStripEnabledForSlide(config, role)) continue;

    const slotX = startX + slot * (slideWidth + gap);
    const slice = bandSliceRect(seamless, slideIndex);
    ctx.drawImage(
      seamless.band,
      slice.sx,
      slice.sy,
      slice.sw,
      slice.sh,
      slotX + slice.dx,
      slice.dy,
      slice.dw,
      slice.dh,
    );
    const nextRole = slideRoles[slot + 1] ?? '';
    const nextIndex = slotSlideIndices?.[slot + 1] ?? slot + 1;
    if (
      gap > 0
      && slot + 1 < slotCount
      && nextIndex >= 0
      && motifStripEnabledForSlide(config, nextRole)
    ) {
      ctx.drawImage(
        seamless.band,
        slice.sx + slice.sw - 1,
        slice.sy,
        1,
        slice.sh,
        slotX + slice.dx + slice.dw,
        slice.dy,
        gap,
        slice.dh,
      );
    }
  }

  ctx.restore();
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} canvasWidth
 * @param {number} canvasHeight
 * @param {HTMLImageElement|ImageBitmap} image
 * @param {MotifStripContext} context
 * @param {MotifStripConfig} config
 * @param {import('./theme.js').DeckTheme} theme
 */
export function paintMotifStripSlice(
  ctx,
  canvasWidth,
  canvasHeight,
  image,
  context,
  config,
  theme,
) {
  if (!image) return;
  const { slideIndex, slideCount } = context;
  const seamless = buildMotifSeamlessBand(image, config, slideCount, canvasWidth, canvasHeight);
  const slice = bandSliceRect(seamless, slideIndex);

  ctx.save();
  ctx.globalAlpha = Math.max(0, Math.min(1, config.opacity ?? 1));
  ctx.imageSmoothingEnabled = false;
  ctx.imageSmoothingQuality = 'high';
  ctx.drawImage(
    seamless.band,
    slice.sx,
    slice.sy,
    slice.sw,
    slice.sh,
    slice.dx,
    slice.dy,
    slice.dw,
    slice.dh,
  );
  ctx.restore();
}

/**
 * @param {MotifStripConfig|null|undefined} config
 * @param {string} [role]
 */
export function motifStripEnabledForSlide(config, role) {
  if (!config?.src) return false;
  const normalized = (role || '').trim().toLowerCase();
  if (isCarouselCtaRole(normalized)) return false;
  if (!normalized) return true;
  return !(config.excludeRoles ?? []).includes(normalized);
}
