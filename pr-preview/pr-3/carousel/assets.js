import { coerceBrightness, parseThemedColorSpec } from './theme.js';

/** @type {Map<string, HTMLImageElement>} */
const imageCache = new Map();

/** Site OG / header bar background (`static/images/og-default.png`, `assets/css/_custom.scss`). */
export const OG_BRAND_BACKGROUND = '#20111c';

/**
 * @param {string} path Bundle-relative, site-root (`/…`), or absolute URL
 * @param {string} [assetBaseUrl] Directory URL for the Hugo page bundle (trailing slash)
 */
export function resolveAssetUrl(path, assetBaseUrl) {
  if (!path) return '';
  if (/^https?:\/\//i.test(path) || path.startsWith('data:')) return path;
  if (path.startsWith('/')) {
    const origin = typeof window !== 'undefined' ? window.location.origin : '';
    return `${origin}${path}`;
  }
  const base = assetBaseUrl || (typeof window !== 'undefined' ? window.location.href : '');
  return new URL(path, base).href;
}

/**
 * @param {string} url
 * @returns {Promise<HTMLImageElement>}
 */
/**
 * @param {string} url
 * @returns {Promise<HTMLImageElement|ImageBitmap>}
 */
export async function loadImage(url) {
  if (!url) {
    throw new Error('Missing image URL');
  }
  const cached = imageCache.get(url);
  if (cached) return cached;

  if (typeof createImageBitmap === 'function') {
    try {
      const response = await fetch(url, { cache: 'no-store' });
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}`);
      }
      const blob = await response.blob();
      const bitmap = await createImageBitmap(blob);
      imageCache.set(url, bitmap);
      return bitmap;
    } catch {
      // Fall back to HTMLImageElement
    }
  }

  const img = await new Promise((resolve, reject) => {
    const el = new Image();
    el.decoding = 'async';
    el.onload = () => resolve(el);
    el.onerror = () => reject(new Error(`Failed to load image: ${url}`));
    el.src = url;
  });
  try {
    if (typeof img.decode === 'function') {
      await img.decode();
    }
  } catch {
    // ignore
  }
  imageCache.set(url, img);
  return img;
}

/**
 * Load an SVG and bake `fill="currentColor"` to a solid color for canvas draw.
 * @param {string} url
 * @param {string} fillHex
 */
export async function loadTintedSvgImage(url, fillHex) {
  const cacheKey = `${url}#${fillHex}`;
  const cached = imageCache.get(cacheKey);
  if (cached) return cached;

  const response = await fetch(url, { cache: 'force-cache' });
  if (!response.ok) {
    throw new Error(`Failed to load SVG: ${url}`);
  }
  let svg = await response.text();
  svg = svg.replace(/fill="currentColor"/gi, `fill="${fillHex}"`);
  svg = svg.replace(/stroke="currentColor"/gi, `stroke="${fillHex}"`);
  if (!/fill="/i.test(svg.split('</svg>')[0])) {
    svg = svg.replace('<svg', `<svg fill="${fillHex}"`);
  }
  const blob = new Blob([svg], { type: 'image/svg+xml' });
  const objectUrl = URL.createObjectURL(blob);
  try {
    const img = await new Promise((resolve, reject) => {
      const el = new Image();
      el.decoding = 'async';
      el.onload = () => resolve(el);
      el.onerror = () => reject(new Error(`Failed to decode SVG: ${url}`));
      el.src = objectUrl;
    });
    imageCache.set(cacheKey, img);
    return img;
  } finally {
    URL.revokeObjectURL(objectUrl);
  }
}

/**
 * Palette token, `#hex`, `transparent`, or `{ color, brightness }` (-100..100).
 * @typedef {string|{color?: string, brightness?: number}} QrColorSpec
 */

/**
 * @typedef {Object} PostCtaQr
 * @property {number|string} [size] Percent of URL–footer slot (`100` or `"100%"`)
 * @property {'split'|'stack'|'stacked'} [layout]
 * @property {number} [columnRatio] Left column width fraction for `split`
 * @property {QrColorSpec} [color] Module color (default `accent2`)
 * @property {number|string} [brightness] Sibling lighten/darken for string `color` (`"+55"` or `55`)
 * @property {QrColorSpec} [light] Tile color (default `transparent`)
 */

/**
 * @typedef {Object} PostCtaAssets
 * @property {string} [logo]
 * @property {string} [logoColor] Theme color token or `#hex` (for SVG tint)
 * @property {string} [featuredImage]
 * @property {string} [brandImage] Full-width brand lockup (e.g. `/images/og-logo.webp`) pinned to slide bottom
 * @property {string} [background] Slide background (default {@link OG_BRAND_BACKGROUND} on `post_cta`)
 * @property {string} [backgroundGradient] `solid` for flat fill matching brand lockup
 * @property {string} [postUrl] Canonical post URL (may be long)
 * @property {string} [shortUrl] Short human URL for printing + QR (recommended)
 * @property {number} [featuredMaxHeight] Featured image max height in px at 1080 canvas (`post_cta`)
 * @property {PostCtaQr} [qr] Nested QR section (`size`, `layout`, `color`, `brightness`, `light`, `columnRatio`)
 * @property {string} [scanLabel] Default scan CTA copy when no `footer` block is present
 * @property {number} [brandMaxHeight] Bottom brand lockup max height in px at 1080 canvas (`post_cta`)
 */

/**
 * @typedef {Object} LoadedPostCtaAssets
 * @property {HTMLImageElement|ImageBitmap|null} logo
 * @property {HTMLImageElement|ImageBitmap|null} featured
 * @property {HTMLImageElement|ImageBitmap|null} brand
 */

/**
 * @param {string} featuredImage
 * @param {string} [assetBaseUrl]
 * @param {string} [bundleBaseUrl] Public Hugo bundle URL (trailing slash)
 */
export function featuredImageCandidateUrls(featuredImage, assetBaseUrl, bundleBaseUrl) {
  if (!featuredImage) return [];
  /** @type {string[]} */
  const urls = [];
  const add = (path, base) => {
    if (!path) return;
    const url = resolveAssetUrl(path, base);
    if (url && !urls.includes(url)) urls.push(url);
  };

  if (/^https?:\/\//i.test(featuredImage) || featuredImage.startsWith('data:')) {
    urls.push(featuredImage);
    return urls;
  }

  if (featuredImage.startsWith('/')) {
    add(featuredImage, assetBaseUrl);
    return urls;
  }

  add(featuredImage, bundleBaseUrl);
  add(featuredImage, assetBaseUrl);
  return urls;
}

/**
 * @param {string} path
 * @param {string[]} urls
 */
async function loadImageFromCandidates(path, urls) {
  let lastError = null;
  for (const url of urls) {
    try {
      return await loadImage(url);
    } catch (error) {
      lastError = error;
    }
  }
  throw lastError ?? new Error(`Failed to load image: ${path}`);
}

/**
 * Preload CTA assets when the studio opens so slide 8 does not race the first paint.
 * @param {PostCtaAssets} config
 * @param {string} assetBaseUrl
 * @param {(token: string) => string} resolveColor
 * @param {string} [bundleBaseUrl]
 */
export async function preloadPostCtaAssets(config, assetBaseUrl, resolveColor, bundleBaseUrl) {
  if (!config?.featuredImage && !config?.logo && !config?.brandImage) return;
  await loadPostCtaAssets(config, assetBaseUrl, resolveColor, bundleBaseUrl);
}

export async function loadPostCtaAssets(config, assetBaseUrl, resolveColor, bundleBaseUrl) {
  /** @type {LoadedPostCtaAssets} */
  const loaded = { logo: null, featured: null, brand: null };

  const jobs = [];
  const useTopLogo = config.logo && !config.brandImage;

  if (useTopLogo) {
    jobs.push((async () => {
      const logoUrl = resolveAssetUrl(config.logo, assetBaseUrl);
      const fill = config.logoColor?.startsWith('#')
        ? config.logoColor
        : resolveColor(config.logoColor || 'text');
      loaded.logo = await loadTintedSvgImage(logoUrl, fill);
    })());
  }

  if (config.featuredImage) {
    jobs.push((async () => {
      const urls = featuredImageCandidateUrls(config.featuredImage, assetBaseUrl, bundleBaseUrl);
      loaded.featured = await loadImageFromCandidates(config.featuredImage, urls);
    })());
  }

  if (config.brandImage) {
    jobs.push((async () => {
      const urls = featuredImageCandidateUrls(config.brandImage, assetBaseUrl, bundleBaseUrl);
      loaded.brand = await loadImageFromCandidates(config.brandImage, urls);
    })());
  }

  const results = await Promise.allSettled(jobs);
  for (const result of results) {
    if (result.status === 'rejected') {
      console.warn('[carousel] CTA asset failed to load:', result.reason);
    }
  }

  return loaded;
}

/**
 * Resolve public bundle base URL for page-bundle images (trailing slash).
 * @param {{ source?: string, slug?: string }} deck
 * @param {string} [deckUrl]
 */
export function resolveBundleBaseUrl(deck, deckUrl) {
  if (typeof window === 'undefined') return '';

  const { origin, pathname } = window.location;
  const previewMatch = pathname.match(/^(.*\/)carousel\.preview(?:\.html)?$/i);
  if (previewMatch) {
    return `${origin}${previewMatch[1]}`;
  }

  if (deckUrl) {
    return new URL('.', new URL(deckUrl, window.location.href)).href;
  }

  const source = deck?.source;
  if (typeof source === 'string') {
    const match = source.match(/content\/([^/]+)\/([^/]+)\//);
    if (match) {
      return `${origin}/${match[1]}/${match[2]}/`;
    }
  }

  if (deck?.slug) {
    return `${origin}/social-protocols/${deck.slug}/`;
  }

  return new URL('.', window.location.href).href;
}

/** @param {unknown} source @returns {PostCtaQr} */
function qrSectionFrom(source) {
  if (!source || typeof source !== 'object' || Array.isArray(source)) return {};
  const qr = /** @type {{ qr?: unknown }} */ (source).qr;
  if (!qr || typeof qr !== 'object' || Array.isArray(qr)) return {};
  return /** @type {PostCtaQr} */ (qr);
}

/** @param {...unknown} vals */
function firstDefined(...vals) {
  for (const value of vals) {
    if (value !== undefined && value !== null) return value;
  }
  return undefined;
}

/**
 * Module-color spec from nested `qr` plus sibling `brightness`.
 * `color` may be a token, `#hex`, `{ color, brightness }`, or `<accent2 brightness='+55'>`.
 * @param {PostCtaQr|Record<string, unknown>|null|undefined} qr
 * @param {string} [fallback]
 * @returns {{ color: string, brightness?: number }}
 */
export function qrModuleColorSpec(qr, fallback = 'accent2') {
  const source = qr && typeof qr === 'object' && !Array.isArray(qr) ? qr : {};
  const parsed = parseThemedColorSpec(source.color ?? fallback) ?? { color: fallback };
  const sibling = coerceBrightness(source.brightness);
  const brightness = sibling !== undefined ? sibling : parsed.brightness;
  return {
    color: parsed.color ?? fallback,
    ...(brightness !== undefined ? { brightness } : {}),
  };
}

/**
 * @param {PostCtaAssets} variant
 * @param {PostCtaAssets} [deckCta]
 */
export function mergePostCtaConfig(variant, deckCta) {
  const deck = deckCta && typeof deckCta === 'object' ? deckCta : {};
  const variantQr = qrSectionFrom(variant);
  const deckQr = qrSectionFrom(deck);
  const size = firstDefined(variantQr.size, deckQr.size) ?? null;
  const layout = firstDefined(variantQr.layout, deckQr.layout) ?? null;
  const columnRatio = firstDefined(variantQr.columnRatio, deckQr.columnRatio) ?? null;
  const color = firstDefined(variantQr.color, deckQr.color) ?? null;
  const brightness = firstDefined(variantQr.brightness, deckQr.brightness) ?? null;
  const light = firstDefined(variantQr.light, deckQr.light) ?? null;
  /** @type {PostCtaQr} */
  const qr = { size, layout, columnRatio, color, brightness, light };
  return {
    logo: variant.logo ?? deck.logo ?? '/images/head.svg',
    logoColor: variant.logoColor ?? deck.logoColor ?? 'text',
    featuredImage: variant.featuredImage ?? deck.featuredImage ?? '',
    brandImage: variant.brandImage ?? deck.brandImage ?? '',
    background: variant.background ?? deck.background ?? OG_BRAND_BACKGROUND,
    backgroundGradient: variant.backgroundGradient ?? deck.backgroundGradient ?? 'solid',
    postUrl: variant.postUrl ?? deck.postUrl ?? '',
    shortUrl: variant.shortUrl ?? deck.shortUrl ?? '',
    featuredMaxHeight: variant.featuredMaxHeight ?? deck.featuredMaxHeight ?? null,
    qr,
    scanLabel: variant.scanLabel ?? deck.scanLabel ?? null,
    brandMaxHeight: variant.brandMaxHeight ?? deck.brandMaxHeight ?? null,
  };
}

/**
 * @param {number} nw
 * @param {number} nh
 * @param {number} maxW
 * @param {number} maxH
 */
export function fitImageBox(nw, nh, maxW, maxH) {
  if (!nw || !nh) {
    return { width: maxW, height: maxH };
  }
  const scale = Math.min(maxW / nw, maxH / nh);
  return {
    width: Math.max(1, Math.round(nw * scale)),
    height: Math.max(1, Math.round(nh * scale)),
  };
}

/**
 * Draw image with aspect ratio preserved (letterbox inside box; box should use slide background).
 * @param {CanvasRenderingContext2D} ctx
 * @param {HTMLImageElement|ImageBitmap} img
 */
export function drawImageContainInBox(ctx, img, dx, dy, dw, dh) {
  const nw = img.naturalWidth || img.width || 0;
  const nh = img.naturalHeight || img.height || 0;
  if (!nw || !nh) return;

  const fit = fitImageBox(nw, nh, dw, dh);
  const x = dx + (dw - fit.width) / 2;
  const y = dy + (dh - fit.height) / 2;
  ctx.drawImage(img, x, y, fit.width, fit.height);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {HTMLImageElement|ImageBitmap} img
 */
export function drawImageContainRounded(ctx, img, dx, dy, dw, dh, radius) {
  ctx.save();
  roundRectPath(ctx, dx, dy, dw, dh, radius);
  ctx.clip();
  drawImageContainInBox(ctx, img, dx, dy, dw, dh);
  ctx.restore();
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} w
 * @param {number} h
 * @param {number} r
 */
export function roundRectPath(ctx, x, y, w, h, r) {
  const radius = Math.min(r, w / 2, h / 2);
  ctx.beginPath();
  ctx.moveTo(x + radius, y);
  ctx.arcTo(x + w, y, x + w, y + h, radius);
  ctx.arcTo(x + w, y + h, x, y + h, radius);
  ctx.arcTo(x, y + h, x, y, radius);
  ctx.arcTo(x, y, x + w, y, radius);
  ctx.closePath();
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {HTMLImageElement} img
 * @param {number} dx
 * @param {number} dy
 * @param {number} dw
 * @param {number} dh
 */
export function drawImageCover(ctx, img, dx, dy, dw, dh) {
  const nw = img.naturalWidth || img.width || 0;
  const nh = img.naturalHeight || img.height || 0;
  if (!nw || !nh) return;

  const ir = nw / nh;
  const dr = dw / dh;
  let sw;
  let sh;
  let sx;
  let sy;

  if (ir > dr) {
    sh = nh;
    sw = sh * dr;
    sx = (nw - sw) / 2;
    sy = 0;
  } else {
    sw = nw;
    sh = sw / dr;
    sx = 0;
    sy = (nh - sh) / 2;
  }

  ctx.drawImage(img, sx, sy, sw, sh, dx, dy, dw, dh);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {HTMLImageElement} img
 * @param {number} dx
 * @param {number} dy
 * @param {number} dw
 * @param {number} dh
 * @param {number} radius
 */
export function drawImageCoverRounded(ctx, img, dx, dy, dw, dh, radius) {
  ctx.save();
  roundRectPath(ctx, dx, dy, dw, dh, radius);
  ctx.clip();
  drawImageCover(ctx, img, dx, dy, dw, dh);
  ctx.restore();
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {HTMLImageElement} img
 * @param {number} cx Center x
 * @param {number} top Top y
 * @param {number} maxW
 * @param {number} maxH
 */
export function drawImageContainCentered(ctx, img, cx, top, maxW, maxH) {
  const nw = img.naturalWidth || img.width || 0;
  const nh = img.naturalHeight || img.height || 0;
  if (!nw || !nh) return;

  const scale = Math.min(maxW / nw, maxH / nh);
  const dw = nw * scale;
  const dh = nh * scale;
  const dx = cx - dw / 2;
  const dy = top + (maxH - dh) / 2;
  ctx.drawImage(img, dx, dy, dw, dh);
}
