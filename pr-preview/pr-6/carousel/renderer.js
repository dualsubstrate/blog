import {
  drawImageContainCentered,
  drawImageContainInBox,
  fitImageBox,
  loadPostCtaAssets,
  mergePostCtaConfig,
  qrModuleColorSpec,
} from './assets.js';
import {
  baselineYFromSymmetricOrange,
  buildInlineFont,
  drawInlineLineInColumn,
  lineInkBoundsInColumn,
  layoutInlineLinesFromInk,
  fontMetricsCacheKey,
  lineSlotHeightPx,
  measureFontLineMetrics,
  measureInlineLine,
  measureLineXHeightBandAscent,
  typoOrangeLinesInEmBox,
  wrapInlineText,
} from './inline-text.js';
import {
  EM_LINE_GAP_PX,
  measurePreparedInkHeight,
} from './layout-metrics.js';
import {
  cellWidthForColumnSpan,
  drawGridSection,
  measurePreparedGridSection,
  prepareGridSection,
  refreshPreparedGridMetrics,
} from './text-grid.js';
import {
  drawQrCodeInBox,
  ensureQrCodegenLoaded,
  isQrBackgroundTransparent,
  makeQrCode,
  parseQrSizePercent,
} from './qr.js';
import { rotateHexHue } from './background.js';
import { paintPanoramicWaveSlice, shouldPaintPanoramicWave } from './background-panorama.js';
import { wavePaletteFromTheme } from './theme.js';
import { isCarouselCtaRole } from './slide-constants.js';
import {
  loadMotifStripImage,
  motifStripEnabledForSlide,
  paintMotifStripSlice,
  parseMotifStripConfig,
} from './motif-strip.js';
import {
  buildFontPx,
  innerWidthFor,
  loadDeckFonts,
  contentColumnBounds,
  mergeTheme,
  SLIDE_SIZE_MAX,
  SLIDE_SIZE_MIN,
  DEFAULT_ASPECT_RATIO_WIDTH,
  DEFAULT_ASPECT_RATIO_HEIGHT,
  canvasHeightForWidth,
  resolveColor,
  resolveInlineRunColor,
  resolveThemedColor,
  resolveDesignFontSize,
  resolveFontSizePx,
  resolveLineHeightMult,
  resolveWidth,
  scaleCanvasPx,
  scaleFontSize,
} from './theme.js';

export { parseQrSizePercent } from './qr.js';

/**
 * @typedef {import('./theme.js').Section} Section
 * @typedef {import('./theme.js').BodyEmphasis} BodyEmphasis
 * @typedef {Object} TextBlock
 * @property {string} text
 * @property {Section} section
 * @property {BodyEmphasis} [emphasis] For `section: body`: normal or punch (accent)
 * @property {import('./theme.js').ColorToken} [color]
 * @property {import('./theme.js').FontRole} [font]
 * @property {string} [weight]
 * @property {number|string} [maxWidth] Canvas px or `%` of deck content column
 * @property {number} [fontSize] Explicit size in px at 1080 canvas (overrides section base + emphasis scale)
 * @property {number} [lineHeight] Optional line-height multiplier (default from section / emphasis)
 * @property {number} [maxLines] Cap wrapped lines (grid cells often use `1`)
 */

/**
 * @typedef {'left' | 'center' | 'right'} HorizontalAlign
 * @typedef {Partial<Record<Section, HorizontalAlign>>} SectionAlignmentMap
 */

/**
 * @typedef {Object} SlideVariant
 * @property {string} archetype
 * @property {SectionAlignmentMap | HorizontalAlign} [alignment] Per-section horizontal alignment; legacy string applies to all sections
 * @property {'top' | 'center' | 'bottom'} [verticalAlign] Vertical body placement (default `top`). Footer blocks stay pinned to the bottom.
 * @property {TextBlock[]} blocks
 * @property {string} [logo] Site or bundle image path for `post_cta` (default `/images/head.svg`)
 * @property {string} [logoColor] Theme color token or `#hex` for SVG logo tint
 * @property {string} [featuredImage] Bundle-relative or absolute URL for `post_cta` hero image
 * @property {string} [brandImage] Full brand lockup at slide bottom (e.g. `/images/og-logo.webp`); skips top `logo`
 * @property {string} [postUrl] Display URL on CTA slide (not clickable in export)
 */

/**
 * @typedef {import('./assets.js').PostCtaAssets} PostCtaAssets
 */

/**
 * @typedef {Object} RenderOptions
 * @property {boolean} [grain] Apply background grain (default true; ignored when `studioPreview` is true).
 * @property {boolean} [studioPreview] Studio strip / grid preview: no grain overlay.
 * @property {HTMLCanvasElement} [targetCanvas] Draw into this canvas instead of creating one.
 * @property {number} [outputSize] Canvas width in px (default merged `theme.size`; height follows aspect ratio).
 * @property {number} [supersample] Render at `outputSize * supersample`, then downscale (default 1).
 * @property {boolean} [showLineBoxes] Debug: fluo-green canvas edge, per-line blue/orange/red guides, post_cta brand lockup box.
 * @property {string} [assetBaseUrl] Hugo bundle directory URL for resolving `featuredImage`
 * @property {string} [bundleBaseUrl] Public bundle URL (from preview path or `deck.source`)
 * @property {PostCtaAssets} [deckCta] Deck-level defaults for `post_cta` slides
 * @property {import('./assets.js').LoadedPostCtaAssets} [loadedCtaAssets] Preloaded images (set by renderSlideToCanvas)
 * @property {import('./background-panorama.js').BackgroundPanoramaContext} [backgroundPanoramaContext] Panoramic wave slice index
 * @property {import('./motif-strip.js').MotifStripContext} [motifStripContext] Panoramic motif strip slice index
 * @property {HTMLImageElement|ImageBitmap} [loadedMotifStrip] Preloaded motif strip asset
 * @property {boolean} [skipBackground] Omit background paint (strip composites a shared panoramic layer first)
 * @property {boolean} [skipMotifStrip] Omit motif strip paint (strip composites a shared panoramic motif layer first)
 * @property {{ deck?: Record<string, unknown> }} [deck] Full deck (studio export context)
 * @property {string} [slideRole] Normalized slide role for motif gating
 */

/** @type {RenderOptions|null} */
let activeRenderOptions = null;

/** @type {Map<string, import('./inline-text.js').FontLineMetrics>|null} */
let activeFontMetricsCache = null;

/**
 * @typedef {Object} LineInkMetrics
 * @property {number} ascent
 * @property {number} descent
 * @property {number} inkHeight
 * @property {number} width
 * @property {number} inkWidth Measured ink box width (blue slot span)
 * @property {number} inkLeft Ink offset left of draw origin
 * @property {number} inkRight Ink offset right of draw origin
 */

/**
 * @typedef {Object} PreparedBlock
 * @property {TextBlock} block
 * @property {Section} section
 * @property {import('./theme.js').FontRole} role
 * @property {string} weight
 * @property {number} fontSizePx
 * @property {number} lineHeightMult Multiplier on orange band for blue line slot height
 * @property {number} emBoxAscent Font strut ascent (typo probes; debug)
 * @property {number} emBoxDescent Font strut descent (typo probes; debug)
 * @property {number} emBoxHeight Blue line slot height (orange band × lineHeightMult); stack advance
 * @property {number} lineSlotPx Same as emBoxHeight; used by typo overlays
 * @property {number} ascenderLinePx Ascender line above baseline (`ASCENDER_PROBE_CHARS`)
 * @property {number} descenderLinePx Descender line below baseline (`DESCENDER_PROBE_CHARS`)
 * @property {number} xHeightPx Lowercase x-height (`X_HEIGHT_PROBE_CHARS`)
 * @property {number} capHeightPx Cap height (baseline to cap top) for this block
 * @property {string[]} lines Plain-text preview per line (markers stripped)
 * @property {import('./inline-text.js').InlineRun[][]} inlineLines Styled tokens per line
 * @property {LineInkMetrics[]} lineMetrics Measured ink box per line
 * @property {number[]} lineAdvances Ink height + optional leading per line (canvas px)
 * @property {number} lineHeight First line advance (cluster stacking)
 * @property {number} height Sum of line advances
 * @property {number} width
 */

function drawBackground(ctx, theme, options = {}) {
  if (options.skipBackground) return;
  const applyGrain = options.grain !== false && options.studioPreview !== true;
  const { size, sizeHeight } = theme;
  const canvasHeight = sizeHeight ?? size;

  if (isCarouselCtaRole(options.slideRole)) {
    ctx.fillStyle = resolveColor(theme, theme.background);
    ctx.fillRect(0, 0, size, canvasHeight);
    if (applyGrain) {
      drawGrainOverlay(ctx, size, canvasHeight, 0.22);
    }
    return;
  }

  const waveColors = wavePaletteFromTheme(theme);
  const panorama = options.backgroundPanoramaContext;
  const waveConfig = theme.backgroundWave ?? {
    style: 'drift',
    lobes: null,
    intensity: 0.32,
    color: 0.55,
    variety: 0.62,
    blur: 0.55,
    radius: 1,
    phase: 0,
  };
  const usePanoramicPaint =
    panorama && panorama.slideCount > 0 && shouldPaintPanoramicWave(waveConfig);

  if (usePanoramicPaint) {
    paintPanoramicWaveSlice(
      ctx,
      size,
      canvasHeight,
      waveColors,
      panorama,
      waveConfig,
    );
  } else {
    const hueShift = waveConfig.hueShift ?? 0;
    ctx.fillStyle = rotateHexHue(waveColors.background, hueShift);
    ctx.fillRect(0, 0, size, canvasHeight);
  }
  if (applyGrain) {
    const grainStrength = usePanoramicPaint ? 0.14 : 0.22;
    drawGrainOverlay(ctx, size, canvasHeight, grainStrength);
  }
}

/** @type {{ canvas: HTMLCanvasElement, width: number, height: number } | null} */
let grainOverlayCache = null;

/** @param {CanvasRenderingContext2D} ctx @param {number} width @param {number} height @param {number} [strength] */
function drawGrainOverlay(ctx, width, height, strength = 0.22) {
  if (
    !grainOverlayCache
    || grainOverlayCache.width !== width
    || grainOverlayCache.height !== height
  ) {
    const scale = 2;
    const grainW = Math.ceil(width / scale);
    const grainH = Math.ceil(height / scale);
    const grainCanvas = document.createElement('canvas');
    grainCanvas.width = grainW;
    grainCanvas.height = grainH;
    const grainCtx = grainCanvas.getContext('2d');
    if (!grainCtx) return;

    const imageData = grainCtx.getImageData(0, 0, grainW, grainH);
    const data = imageData.data;
    for (let i = 0; i < data.length; i += 4) {
      const n = 128 + (Math.random() - 0.5) * 18;
      data[i] = n;
      data[i + 1] = n;
      data[i + 2] = n;
      data[i + 3] = 255;
    }
    grainCtx.putImageData(imageData, 0, 0);
    grainOverlayCache = { canvas: grainCanvas, width, height };
  }

  ctx.save();
  ctx.globalAlpha = strength;
  ctx.globalCompositeOperation = 'overlay';
  ctx.imageSmoothingEnabled = true;
  ctx.imageSmoothingQuality = 'high';
  ctx.drawImage(grainOverlayCache.canvas, 0, 0, width, height);
  ctx.restore();
}

/** @param {TextBlock} block */
function maxLinesForBlock(block) {
  if (block.maxLines != null && Number.isFinite(block.maxLines)) {
    return Math.max(1, Math.round(block.maxLines));
  }
  if (block.section === 'header') return 2;
  if (block.section === 'footer') return 4;
  if (block.emphasis === 'punch') return 4;
  return 6;
}

/** @param {TextBlock} block @param {import('./theme.js').DeckTheme} theme */
function defaultRoleForBlock(block, theme) {
  if (block.font) return block.font;
  if (block.section === 'body' && block.emphasis === 'punch') return 'display';
  return 'body';
}

/** @param {TextBlock} block */
function defaultWeightForBlock(block) {
  if (block.weight) return block.weight;
  if (block.section === 'header') return '600';
  if (block.section === 'body' && block.emphasis === 'punch') return '700';
  return '400';
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: import('./inline-text.js').InlineRun) => string} fontFor
 * @param {import('./theme.js').FontRole} role
 * @param {number} fontSizePx
 * @param {string} weight
 * @param {boolean[]} boldProbes
 */
function cachedFontLineMetrics(ctx, fontFor, role, fontSizePx, weight, boldProbes) {
  const cache = activeFontMetricsCache;
  const key = fontMetricsCacheKey(role, fontSizePx, weight, boldProbes);
  if (cache) {
    const hit = cache.get(key);
    if (hit) return hit;
  }
  const metrics = measureFontLineMetrics(ctx, fontFor, fontSizePx, boldProbes);
  if (cache) cache.set(key, metrics);
  return metrics;
}

/** @param {import('./inline-text.js').InlineRun[][]} inlineLines @param {string} blockWeight */
function fontWeightProbes(inlineLines, blockWeight) {
  /** @type {boolean[]} */
  const probes = [false];
  if (blockWeight === '700') {
    return probes;
  }
  for (const line of inlineLines) {
    for (const token of line) {
      if (token.bold) {
        probes.push(true);
        return probes;
      }
    }
  }
  return probes;
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {TextBlock} block
 * @param {number} maxWidth
 */
function prepareBlock(ctx, theme, block, maxWidth) {
  const normalized = applyBlockDefaults(block);
  const role = defaultRoleForBlock(normalized, theme);
  const weight = defaultWeightForBlock(normalized);
  const designPx = resolveDesignFontSize(normalized, theme);
  const budget = maxLinesForBlock(normalized);
  let fontSizePx = resolveFontSizePx(theme, designPx);
  /** @type {import('./inline-text.js').InlineRun[][]} */
  let inlineLines = [];

  /** @param {import('./inline-text.js').InlineRun} token */
  const fontFor = (token) => buildInlineFont(theme, role, fontSizePx, weight, token);

  while (true) {
    inlineLines = wrapInlineText(ctx, fontFor, normalized.text, maxWidth);
    if (inlineLines.length <= budget || fontSizePx <= scaleFontSize(30, theme.size)) break;
    fontSizePx -= Math.max(1, Math.round(scaleFontSize(4, theme.size)));
  }

  const lineHeightMult = resolveLineHeightMult(theme, normalized);
  const boldProbes = fontWeightProbes(inlineLines, weight);
  const fontMetrics = cachedFontLineMetrics(ctx, fontFor, role, fontSizePx, weight, boldProbes);
  const lineSlotHeight = lineSlotHeightPx({ ...fontMetrics, lineHeightMult }, lineHeightMult);
  const layout = layoutInlineLinesFromInk(ctx, inlineLines, fontFor, lineSlotHeight, fontSizePx);
  const lines = inlineLines.map((tokens) => tokens.map((token) => token.text).join(''));

  return {
    block: normalized,
    section: normalized.section,
    role,
    weight,
    fontSizePx,
    lineHeightMult,
    emBoxAscent: fontMetrics.emBoxAscent,
    emBoxDescent: fontMetrics.emBoxDescent,
    emBoxHeight: lineSlotHeight,
    lineSlotPx: lineSlotHeight,
    ascenderLinePx: fontMetrics.ascenderLinePx,
    descenderLinePx: fontMetrics.descenderLinePx,
    xHeightPx: fontMetrics.xHeightPx,
    capHeightPx: fontMetrics.capHeightPx,
    lines,
    inlineLines,
    lineMetrics: layout.lineMetrics,
    lineAdvances: layout.lineAdvances,
    lineHeight: layout.lineHeight,
    height: layout.height,
    width: layout.width,
  };
}

/** @param {TextBlock} block */
function applyBlockDefaults(block) {
  const section = block.section;
  if (section === 'body') {
    const emphasis = block.emphasis ?? 'normal';
    if (emphasis === 'punch') {
      return {
        ...block,
        emphasis,
        color: block.color ?? 'accent1',
        font: block.font ?? 'display',
        weight: block.weight ?? '700',
      };
    }
    return {
      ...block,
      emphasis,
      color: block.color ?? 'text',
      font: block.font ?? 'body',
      weight: block.weight ?? '400',
    };
  }
  if (section === 'header') {
    return {
      ...block,
      color: block.color ?? 'muted',
      font: block.font ?? 'body',
      weight: block.weight ?? '600',
    };
  }
  if (section === 'footer') {
    return {
      ...block,
      color: block.color ?? 'muted',
      font: block.font ?? 'body',
      weight: block.weight ?? '400',
    };
  }
  return block;
}

/** @param {SlideVariant} variant */
function normalizeBlocks(variant) {
  return variant.blocks.map((block) => {
    if (block.section === 'grid') {
      return /** @type {TextBlock} */ ({ ...block, section: 'grid' });
    }
    return applyBlockDefaults({ ...block });
  });
}

/**
 * @param {PreparedBlock[]} prepared
 * @param {import('./theme.js').DeckTheme} theme
 */
function measureSlideStackHeight(prepared, theme) {
  const headerPad = gapFor(theme, 0.028);
  const footerPad = gapFor(theme, 0.018);
  const { marginVertical } = theme;
  const header = pick(prepared, 'header')[0];
  const footer = pick(prepared, 'footer')[0];
  const clusterHeight = measureMessageStackHeight(prepared, theme);

  let total = marginVertical * 2 + clusterHeight;
  if (header) total += measureBlockInkHeight(header) + headerPad;
  if (footer) total += measureBlockInkHeight(footer) + footerPad;
  return total;
}

/**
 * Body + grid sections in document order (excludes header/footer).
 * @param {PreparedBlock[]} prepared
 */
function messageStackItems(prepared) {
  return prepared.filter((item) => item.section === 'body' || item.section === 'grid');
}

/**
 * @param {PreparedBlock[]} prepared
 */
function refreshAllGridMetrics(prepared) {
  for (const item of prepared) {
    if (item.section === 'grid') {
      refreshPreparedGridMetrics(/** @type {import('./text-grid.js').PreparedGridSection} */ (item));
    }
  }
}

/**
 * Middle-band height available for body + grid message stack.
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock[]} prepared
 */
function messageStackFitBudget(theme, prepared) {
  const canvasHeight = theme.sizeHeight ?? theme.size;
  let budget = canvasHeight - theme.marginVertical * 2;
  const headerPad = gapFor(theme, 0.028);
  const footerPad = gapFor(theme, 0.018);
  const header = pick(prepared, 'header')[0];
  const footer = pick(prepared, 'footer')[0];
  if (header) budget -= measureBlockInkHeight(header) + headerPad;
  if (footer) budget -= measureBlockInkHeight(footer) + footerPad;
  return Math.max(1, budget);
}

/**
 * Gap between a body cluster and a grid block (or two grid/body segments).
 * @param {import('./theme.js').DeckTheme} theme
 */
function messageStackSectionGap(theme) {
  return Math.round(theme.fontSizes.body * (theme.emphasisGap ?? 0.25));
}

/**
 * Body blocks run together; grid is its own segment.
 * @param {PreparedBlock[]} items
 */
function messageStackSegments(items) {
  /** @type {Array<{ kind: 'body', blocks: PreparedBlock[] } | { kind: 'grid', grid: import('./text-grid.js').PreparedGridSection }>} */
  const segments = [];
  /** @type {PreparedBlock[]} */
  let bodyRun = [];

  const flushBody = () => {
    if (bodyRun.length === 0) return;
    segments.push({ kind: 'body', blocks: bodyRun });
    bodyRun = [];
  };

  for (const item of items) {
    if (item.section === 'grid') {
      flushBody();
      segments.push({
        kind: 'grid',
        grid: /** @type {import('./text-grid.js').PreparedGridSection} */ (item),
      });
    } else {
      bodyRun.push(item);
    }
  }
  flushBody();
  return segments;
}

/**
 * @param {PreparedBlock[]} prepared
 * @param {import('./theme.js').DeckTheme} theme
 */
function measureMessageStackHeight(prepared, theme) {
  const segments = messageStackSegments(messageStackItems(prepared));
  if (segments.length === 0) return 0;
  const sectionGap = messageStackSectionGap(theme);
  let height = 0;
  for (let i = 0; i < segments.length; i += 1) {
    const segment = segments[i];
    if (segment.kind === 'body') {
      height += measureClusterHeight(segment.blocks);
    } else {
      height += measurePreparedGridSection(segment.grid);
    }
    if (i < segments.length - 1) height += sectionGap;
  }
  return height;
}

/**
 * Re-wrap and re-measure a prepared block after `fontSizePx` changed.
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock} prepared
 * @param {number} maxWidth
 */
function relayoutPreparedBlock(ctx, theme, prepared, maxWidth) {
  const { block, role, weight, fontSizePx } = prepared;
  /** @param {import('./inline-text.js').InlineRun} token */
  const fontFor = (token) => buildInlineFont(theme, role, fontSizePx, weight, token);
  const inlineLines = wrapInlineText(ctx, fontFor, block.text, maxWidth);
  const lineHeightMult = resolveLineHeightMult(theme, block);
  prepared.lineHeightMult = lineHeightMult;
  const boldProbes = fontWeightProbes(inlineLines, weight);
  const fontMetrics = cachedFontLineMetrics(ctx, fontFor, role, fontSizePx, weight, boldProbes);
  const lineSlotHeight = lineSlotHeightPx({ ...fontMetrics, lineHeightMult }, lineHeightMult);
  const layout = layoutInlineLinesFromInk(ctx, inlineLines, fontFor, lineSlotHeight, fontSizePx);
  prepared.emBoxAscent = fontMetrics.emBoxAscent;
  prepared.emBoxDescent = fontMetrics.emBoxDescent;
  prepared.emBoxHeight = lineSlotHeight;
  prepared.lineSlotPx = lineSlotHeight;
  prepared.ascenderLinePx = fontMetrics.ascenderLinePx;
  prepared.descenderLinePx = fontMetrics.descenderLinePx;
  prepared.xHeightPx = fontMetrics.xHeightPx;
  prepared.capHeightPx = fontMetrics.capHeightPx;
  prepared.inlineLines = inlineLines;
  prepared.lines = inlineLines.map((tokens) => tokens.map((token) => token.text).join(''));
  prepared.lineMetrics = layout.lineMetrics;
  prepared.lineAdvances = layout.lineAdvances;
  prepared.lineHeight = layout.lineHeight;
  prepared.height = layout.height;
  prepared.width = layout.width;
}

/**
 * Shrink all blocks until header, body cluster, and footer fit on the slide.
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock[]} prepared
 * @param {number[]} maxWidths
 */
function fitPreparedToCanvas(ctx, theme, prepared, maxWidths) {
  const minPx = scaleFontSize(30, theme.size);
  const step = Math.max(1, Math.round(scaleFontSize(4, theme.size)));

  const stackBudget = messageStackFitBudget(theme, prepared);
  while (measureMessageStackHeight(prepared, theme) > stackBudget) {
    let reduced = false;
    for (let i = 0; i < prepared.length; i += 1) {
      const item = prepared[i];
      if (item.section === 'grid') {
        const grid = /** @type {import('./text-grid.js').PreparedGridSection} */ (item);
        let gridReduced = false;
        for (const entry of grid.cells) {
          if (entry.prepared.fontSizePx <= minPx) continue;
          entry.prepared.fontSizePx -= step;
          const cellWidth = cellWidthForColumnSpan(
            entry.col,
            entry.colSpan,
            grid.colWidths,
            grid.colGapPx,
          );
          relayoutPreparedBlock(ctx, theme, entry.prepared, cellWidth);
          gridReduced = true;
        }
        if (gridReduced) {
          refreshPreparedGridMetrics(grid);
          reduced = true;
        }
        continue;
      }
      if (item.fontSizePx <= minPx) continue;
      item.fontSizePx -= step;
      relayoutPreparedBlock(ctx, theme, item, maxWidths[i]);
      reduced = true;
    }
    refreshAllGridMetrics(prepared);
    if (!reduced) break;
  }

  refreshAllGridMetrics(prepared);
  return prepared;
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {SlideVariant} variant
 */
function prepareVariant(ctx, theme, variant) {
  const { contentWidth } = contentColumnBounds(theme);
  const blocks = normalizeBlocks(variant);
  const maxWidths = blocks.map((block) => (
    block.maxWidth != null
      ? resolveWidth(block.maxWidth, contentWidth, 1)
      : contentWidth
  ));

  activeFontMetricsCache = new Map();
  try {
    const prepared = blocks.map((block, index) => {
      if (block.section === 'grid') {
        return /** @type {PreparedBlock} */ (/** @type {unknown} */ (
          prepareGridSection(ctx, theme, /** @type {import('./text-grid.js').GridBlock} */ (block), contentWidth, (b, w) => prepareBlock(ctx, theme, b, w))
        ));
      }
      return prepareBlock(ctx, theme, block, maxWidths[index]);
    });
    const fitted = fitPreparedToCanvas(ctx, theme, prepared, maxWidths);
    refreshAllGridMetrics(fitted);
    return fitted;
  } finally {
    activeFontMetricsCache = null;
  }
}

/**
 * @param {PreparedBlock[]} group
 */
function clusterLineCount(group) {
  return group.reduce((count, block) => count + block.inlineLines.length, 0);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {PreparedBlock} prepared
 * @param {number} lineIndex
 */
function lineXHeightBandAscent(ctx, fontFor, prepared, lineIndex) {
  return measureLineXHeightBandAscent(ctx, fontFor, prepared.inlineLines[lineIndex], {
    xHeightPx: prepared.xHeightPx,
    capHeightPx: prepared.capHeightPx,
    emBoxAscent: prepared.emBoxAscent,
  });
}


/**
 * Vertical rhythm for one line: blue line slot height (orange band × lineHeights).
 * @param {PreparedBlock} prepared
 */
function lineAdvancePx(prepared) {
  if (prepared.lineSlotPx != null && Number.isFinite(prepared.lineSlotPx)) {
    return prepared.lineSlotPx;
  }
  return lineSlotHeightPx(prepared, prepared.lineHeightMult);
}

/** @param {number} emTop @param {PreparedBlock} prepared */
function emTopAfterLine(emTop, prepared) {
  return typoOrangeLinesInEmBox(emTop, prepared, prepared.lineHeightMult).emBottom;
}

/**
 * @param {number|null|'all'} seamlessAfter
 * @param {number} blockIndex
 * @param {number} lineIndex
 * @param {number} linesInBlock
 * @param {number} blockCount
 */
function clusterSeamlessAfterBlock(seamlessAfter, blockIndex, lineIndex, linesInBlock, blockCount) {
  if (lineIndex !== linesInBlock - 1 || blockIndex >= blockCount - 1) {
    return false;
  }
  if (seamlessAfter === 'all') return true;
  return seamlessAfter === blockIndex;
}

/**
 * @param {PreparedBlock[]} group
 * @param {{ seamlessAfterBlockIndex?: number|null|'all' }} [options]
 */
function measureClusterHeight(group, options = {}) {
  const seamlessAfter = options.seamlessAfterBlockIndex;
  let height = 0;
  let lineOrdinal = 0;
  const totalLines = clusterLineCount(group);

  for (let b = 0; b < group.length; b += 1) {
    const block = group[b];
    for (let i = 0; i < block.inlineLines.length; i += 1) {
      height += lineAdvancePx(block);
      lineOrdinal += 1;
      if (lineOrdinal < totalLines) {
        const seamless = clusterSeamlessAfterBlock(
          seamlessAfter,
          b,
          i,
          block.inlineLines.length,
          group.length,
        );
        if (!seamless) height += EM_LINE_GAP_PX;
      }
    }
  }
  return height;
}

/**
 * @param {PreparedBlock} block
 */
function measureBlockInkHeight(block) {
  return measurePreparedInkHeight(block);
}

/**
 * @param {PreparedBlock[]} blocks
 * @param {number} gapBetween
 * @param {number} [gapAfterLast]
 */
function measureBlocksStackHeight(blocks, gapBetween, gapAfterLast = 0) {
  if (!blocks.length) return 0;
  let h = 0;
  for (let i = 0; i < blocks.length; i += 1) {
    h += measureBlockInkHeight(blocks[i]);
    if (i < blocks.length - 1) h += gapBetween;
    else h += gapAfterLast;
  }
  return h;
}

/** Pixels from last line descender to em slot bottom (padding below footer glyphs). */
function measureBlockTrailingEmPad(block) {
  const lineCount = block.inlineLines.length;
  if (lineCount <= 0) return 0;
  let emTop = 0;
  for (let i = 0; i < lineCount - 1; i += 1) {
    emTop = emTopAfterLine(emTop, block) + EM_LINE_GAP_PX;
  }
  const orange = typoOrangeLinesInEmBox(emTop, block, block.lineHeightMult);
  return Math.max(0, orange.emBottom - orange.descenderY);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {PreparedBlock} prepared
 * @param {number} emTop
 * @param {number} lineIndex
 */
function baselineAtEmTop(_ctx, _fontFor, prepared, emTop, _lineIndex) {
  return baselineYFromSymmetricOrange(emTop, prepared);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} width
 */
function strokeDebugHorizontal(ctx, x, y, width) {
  ctx.beginPath();
  ctx.moveTo(x, y);
  ctx.lineTo(x + width, y);
  ctx.stroke();
}

/**
 * Stroked rect fully inside the given bounds (avoids right/bottom clip at canvas edge).
 * Canvas strokes are centered on the path; without inset, outer half is clipped on max edges.
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} width
 * @param {number} height
 */
function strokeDebugRect(ctx, x, y, width, height) {
  const lw = ctx.lineWidth;
  const insetW = Math.max(0, width - lw);
  const insetH = Math.max(0, height - lw);
  ctx.strokeRect(x + lw / 2, y + lw / 2, insetW, insetH);
}

/**
 * Blue dashed slot outline (matches per-line text debug boxes).
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} x
 * @param {number} y
 * @param {number} width
 * @param {number} height
 * @param {number} canvasSize
 */
function drawImageSlotDebug(ctx, x, y, width, height, canvasSize) {
  const ref = scaleCanvasPx(48, canvasSize);
  const blueW = Math.max(2, Math.round(ref * 0.03));
  const dash = [
    Math.max(3, Math.round(ref * 0.035)),
    Math.max(3, Math.round(ref * 0.035)),
  ];

  ctx.save();
  ctx.setLineDash(dash);
  ctx.lineWidth = blueW;
  ctx.strokeStyle = '#5ec8ff';
  strokeDebugRect(ctx, x, y, Math.max(1, width), Math.max(1, height));
  ctx.restore();
}

/** @param {number} px */
function clampSlideEdge(px) {
  return Math.max(SLIDE_SIZE_MIN, Math.min(SLIDE_SIZE_MAX, Math.round(px)));
}

/** @param {string} text */
function looksLikeUrlText(text) {
  if (!text) return false;
  const t = String(text).trim();
  if (t.length < 10) return false;
  if (/^https?:\/\//i.test(t)) return true;
  // "No scheme" URLs like blog.dualsubstrate.com/foo/bar
  if (/^[a-z0-9.-]+\.[a-z]{2,}(\/\S*)?$/i.test(t)) return true;
  return false;
}

/** @param {string} text */
function looksLikeLongUrlText(text) {
  if (!looksLikeUrlText(text)) return false;
  const t = String(text).trim();
  // If it includes a deep path, it tends to be illegible on-image.
  const slashes = (t.match(/\//g) || []).length;
  return t.length >= 34 || slashes >= 3;
}

/**
 * @param {PreparedBlock[]} bodyBlocks
 */
function splitTitleAndUrlBlocks(bodyBlocks) {
  /** @type {PreparedBlock[]} */
  const titles = [];
  /** @type {PreparedBlock[]} */
  const urls = [];
  for (const block of bodyBlocks) {
    if (looksLikeUrlText(block.text) && !looksLikeLongUrlText(block.text)) {
      urls.push(block);
    } else {
      titles.push(block);
    }
  }
  return { titles, urls };
}

/**
 * @param {PreparedBlock[]} footers
 */
function splitPostCtaFooters(footers) {
  /** @type {PreparedBlock[]} */
  const urlLines = [];
  /** @type {PreparedBlock[]} */
  const scanLines = [];
  for (const block of footers) {
    if (looksLikeUrlText(block.text) && !looksLikeLongUrlText(block.text)) {
      urlLines.push(block);
    } else {
      scanLines.push(block);
    }
  }
  return { urlLines, scanLines };
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock[]} blocks
 * @param {SlideVariant} variant
 * @param {number} startY
 * @param {import('./theme.js').ContentColumn} column
 * @param {number} gap
 */
function drawBlocksTopDown(ctx, theme, blocks, variant, startY, column, gap, options = {}) {
  const gapAfterLast = options.gapAfterLast ?? gap;
  let y = startY;
  for (let i = 0; i < blocks.length; i += 1) {
    const block = blocks[i];
    drawAt(ctx, theme, block, variant, y, column);
    const after = i < blocks.length - 1 ? gap : gapAfterLast;
    y += measureBlockInkHeight(block) + after;
  }
  return y;
}

/**
 * Single canvas bitmap outline (drawn behind text). Uses the real bitmap size, not theme.size.
 * @param {CanvasRenderingContext2D} ctx
 */
function drawCanvasBoundsDebug(ctx) {
  const canvasW = ctx.canvas.width;
  const canvasH = ctx.canvas.height;
  const edge = Math.min(canvasW, canvasH);
  const edgeW = Math.max(2, Math.round(edge * 0.006));

  ctx.save();
  ctx.setLineDash([]);
  ctx.lineWidth = edgeW;
  ctx.strokeStyle = '#39ff14';
  strokeDebugRect(ctx, 0, 0, canvasW, canvasH);
  ctx.restore();
}

/** Skip orange draw when it would still sit on the em box edge (safety net). */
const DEBUG_ORANGE_DRAW_INSET_PX = 2;

/** @param {LineInkMetrics|undefined} lineMetric */
function lineExtentsFromMetric(lineMetric) {
  if (!lineMetric) {
    return { inkLeft: 0, inkRight: 0, inkWidth: 0 };
  }
  const inkWidth = lineMetric.inkWidth ?? lineMetric.width ?? 0;
  return {
    inkLeft: lineMetric.inkLeft ?? 0,
    inkRight: lineMetric.inkRight ?? inkWidth,
    inkWidth,
  };
}

/**
 * Line-box overlay: draw after text so blue em edges stay visible.
 * Blue slot spans this line's measured ink width and x (not the full content column).
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: import('./inline-text.js').InlineRun) => string} fontFor
 * @param {PreparedBlock} prepared
 * @param {{ columnLeft: number, contentWidth: number, columnRight?: number, canvasSize?: number, marginHorizontal?: number }} column
 * @param {'left' | 'center' | 'right'} alignment
 * @param {number} emTop
 * @param {number} lineIndex
 * @param {import('./inline-text.js').InlineRun[]} lineTokens
 */
function drawLineDebugAtEmTop(ctx, fontFor, prepared, column, alignment, emTop, lineIndex, lineTokens) {
  const cachedExtents = lineExtentsFromMetric(prepared.lineMetrics[lineIndex]);
  const inkBounds = lineInkBoundsInColumn(
    ctx,
    lineTokens,
    fontFor,
    column,
    alignment,
    prepared.fontSizePx,
    cachedExtents,
  );
  const slotLeft = inkBounds.x + inkBounds.inkLeft;
  const slotWidth = Math.max(1, Math.round(inkBounds.width));
  const bandAscent = lineXHeightBandAscent(ctx, fontFor, prepared, lineIndex);
  const baselineY = baselineAtEmTop(ctx, fontFor, prepared, emTop, lineIndex);
  const meanlineY = baselineY - bandAscent;
  const orangeLines = typoOrangeLinesInEmBox(emTop, prepared);
  const emBottomY = orangeLines.emBottom;
  const slotHeight = Math.max(1, emBottomY - emTop);

  // Debug colors map 1:1 to typo-probes roles (see static/carousel/typo-probes.js RULES).

  const thin = Math.max(1, Math.round(prepared.fontSizePx * 0.012));
  const blueW = Math.max(2, Math.round(prepared.fontSizePx * 0.03));
  const dash = [
    Math.max(3, Math.round(prepared.fontSizePx * 0.035)),
    Math.max(3, Math.round(prepared.fontSizePx * 0.035)),
  ];

  ctx.save();

  ctx.setLineDash([]);
  ctx.lineWidth = blueW;
  ctx.strokeStyle = '#5ec8ff';
  strokeDebugRect(ctx, slotLeft, emTop, slotWidth, slotHeight);

  ctx.lineWidth = thin;
  ctx.strokeStyle = 'rgba(255, 168, 64, 0.95)';
  ctx.setLineDash(dash);
  if (orangeLines.ascenderY - emTop >= DEBUG_ORANGE_DRAW_INSET_PX) {
    strokeDebugHorizontal(ctx, slotLeft, orangeLines.ascenderY, slotWidth);
  }
  if (emBottomY - orangeLines.descenderY >= DEBUG_ORANGE_DRAW_INSET_PX) {
    strokeDebugHorizontal(ctx, slotLeft, orangeLines.descenderY, slotWidth);
  }

  ctx.setLineDash([2, Math.max(2, Math.round(prepared.fontSizePx * 0.02))]);
  ctx.strokeStyle = 'rgba(255, 96, 96, 0.95)';
  strokeDebugHorizontal(ctx, slotLeft, meanlineY, slotWidth);
  strokeDebugHorizontal(ctx, slotLeft, baselineY, slotWidth);

  ctx.restore();
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock[]} group
 * @param {SlideVariant} variant
 * @param {number} bandTop
 * @param {{ seamlessAfterBlockIndex?: number|null|'all' }} [options]
 * @returns {number} Y at bottom of last line slot (blue box)
 */
function drawClusterLines(ctx, theme, group, variant, bandTop, column, options = {}) {
  if (group.length === 0) return bandTop;

  const seamlessAfter = options.seamlessAfterBlockIndex;
  const blockCount = group.length;
  let emTop = bandTop;
  const showDebug = Boolean(activeRenderOptions?.showLineBoxes);
  const totalLines = clusterLineCount(group);
  let lineOrdinal = 0;

  for (let b = 0; b < group.length; b += 1) {
    const prepared = group[b];
    const blockColor = prepared.block.color || 'text';
    ctx.fillStyle = resolveColor(theme, blockColor);
    ctx.textAlign = 'left';
    /** @param {import('./inline-text.js').InlineRun} token */
    const fontFor = (token) => buildInlineFont(theme, prepared.role, prepared.fontSizePx, prepared.weight, token);
    /** @param {import('./inline-text.js').InlineRun} token */
    const colorFor = (token) => resolveInlineRunColor(theme, token, blockColor);
    const alignment = sectionAlignmentFor(variant, prepared.section);

    for (let i = 0; i < prepared.inlineLines.length; i += 1) {
      const line = prepared.inlineLines[i];
      const baselineY = baselineAtEmTop(ctx, fontFor, prepared, emTop, i);
      const lineExtents = lineExtentsFromMetric(prepared.lineMetrics[i]);
      drawInlineLineInColumn(
        ctx,
        line,
        fontFor,
        column,
        alignment,
        baselineY,
        prepared.fontSizePx,
        lineExtents,
        colorFor,
      );
      if (showDebug) {
        drawLineDebugAtEmTop(ctx, fontFor, prepared, column, alignment, emTop, i, line);
      }
      lineOrdinal += 1;
      if (lineOrdinal < totalLines) {
        const seamless = clusterSeamlessAfterBlock(
          seamlessAfter,
          b,
          i,
          prepared.inlineLines.length,
          blockCount,
        );
        emTop = emTopAfterLine(emTop, prepared) + (seamless ? 0 : EM_LINE_GAP_PX);
      }
    }
  }

  return emTopAfterLine(emTop, group[group.length - 1]);
}

function drawPreparedBlock(ctx, theme, prepared, alignment, y, column) {
  const blockColor = prepared.block.color || 'text';
  ctx.fillStyle = resolveColor(theme, blockColor);
  ctx.textAlign = 'left';

  let emTop = y;
  const showDebug = Boolean(activeRenderOptions?.showLineBoxes);
  /** @param {import('./inline-text.js').InlineRun} token */
  const colorFor = (token) => resolveInlineRunColor(theme, token, blockColor);

  for (let i = 0; i < prepared.inlineLines.length; i += 1) {
    /** @param {import('./inline-text.js').InlineRun} token */
    const fontFor = (token) => buildInlineFont(theme, prepared.role, prepared.fontSizePx, prepared.weight, token);
    const line = prepared.inlineLines[i];
    const baselineY = baselineAtEmTop(ctx, fontFor, prepared, emTop, i);
    const lineExtents = lineExtentsFromMetric(prepared.lineMetrics[i]);
    drawInlineLineInColumn(
      ctx,
      line,
      fontFor,
      column,
      alignment,
      baselineY,
      prepared.fontSizePx,
      lineExtents,
      colorFor,
    );
    if (showDebug) {
      drawLineDebugAtEmTop(ctx, fontFor, prepared, column, alignment, emTop, i, line);
    }
    if (i < prepared.inlineLines.length - 1) {
      emTop = emTopAfterLine(emTop, prepared) + EM_LINE_GAP_PX;
    }
  }
}

/**
 * @param {PreparedBlock[]} items
 * @param {Section} section
 */
function pick(items, section) {
  return items.filter((item) => item.section === section);
}

/** @param {unknown} value @returns {value is HorizontalAlign} */
function isHorizontalAlign(value) {
  return value === 'left' || value === 'center' || value === 'right';
}

/** @param {unknown} value @returns {HorizontalAlign} */
function coerceHorizontalAlign(value) {
  return isHorizontalAlign(value) ? value : 'left';
}

/** @param {SlideVariant} variant @returns {Record<Section, HorizontalAlign>} */
export function normalizeSectionAlignment(variant) {
  const raw = variant.alignment;
  if (typeof raw === 'string') {
    const align = coerceHorizontalAlign(raw);
    return { header: align, body: align, footer: align };
  }
  if (raw && typeof raw === 'object' && !Array.isArray(raw)) {
    return {
      header: coerceHorizontalAlign(raw.header),
      body: coerceHorizontalAlign(raw.body),
      footer: coerceHorizontalAlign(raw.footer),
    };
  }
  return { header: 'left', body: 'left', footer: 'left' };
}

/** @param {SlideVariant} variant @param {Section} section */
export function sectionAlignmentFor(variant, section) {
  return normalizeSectionAlignment(variant)[section] || 'left';
}

/** @param {SlideVariant} variant @returns {SectionAlignmentMap} */
export function compactSectionAlignment(variant) {
  const normalized = normalizeSectionAlignment(variant);
  /** @type {SectionAlignmentMap} */
  const compact = {};
  for (const section of /** @type {Section[]} */ (['header', 'body', 'footer'])) {
    if (normalized[section] !== 'left') {
      compact[section] = normalized[section];
    }
  }
  return compact;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {number} [ratio]
 */
function gapFor(theme, ratio = 0.04) {
  return Math.round(theme.size * ratio);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock} item
 * @param {SlideVariant} variant
 * @param {number} y
 */
function drawAt(ctx, theme, item, variant, y, column) {
  const alignment = sectionAlignmentFor(variant, item.section);
  drawPreparedBlock(ctx, theme, item, alignment, y, column);
}


/**
 * @param {SlideVariant} variant
 * @returns {'top' | 'center' | 'bottom'}
 */
function resolveVerticalAlign(variant) {
  const v = variant.verticalAlign;
  if (v === 'center' || v === 'bottom') return v;
  return 'top';
}

/**
 * Stack message blocks in a content band.
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock[]} group
 * @param {SlideVariant} variant
 * @param {number} [bandTop]
 * @param {number} [bandBottom]
 * @param {{ align?: 'top' | 'center' | 'bottom' }} [options]
 */
function layoutGroupedStack(ctx, theme, group, variant, bandTop, bandBottom, column, options = {}) {
  const { sizeHeight, marginVertical } = theme;
  const canvasHeight = sizeHeight ?? theme.size;
  const top = bandTop ?? marginVertical;
  const bottom = bandBottom ?? canvasHeight - marginVertical;
  const align = options.align ?? 'top';

  const groupHeight = measureClusterHeight(group);

  let clusterTop = top;
  if (align === 'bottom') {
    clusterTop = bottom - groupHeight;
  } else if (align === 'center') {
    clusterTop = top + (bottom - top - groupHeight) / 2;
  }

  drawClusterLines(ctx, theme, group, variant, clusterTop, column);
}

/**
 * @param {PreparedBlock[]} prepared
 */
function bodyCluster(prepared) {
  return pick(prepared, 'body');
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock[]} prepared
 * @param {SlideVariant} variant
 * @param {number} top
 * @param {{ columnLeft: number, contentWidth: number, columnRight?: number, canvasSize: number, marginHorizontal: number }} column
 */
function drawMessageStack(ctx, theme, prepared, variant, top, column) {
  const segments = messageStackSegments(messageStackItems(prepared));
  const sectionGap = messageStackSectionGap(theme);
  let y = top;

  for (let i = 0; i < segments.length; i += 1) {
    const segment = segments[i];
    if (segment.kind === 'body') {
      y = drawClusterLines(ctx, theme, segment.blocks, variant, y, column);
    } else {
      y = drawGridSection(ctx, theme, variant, segment.grid, column, y, {
        drawPreparedBlock,
        measureBlockInkHeight,
        sectionAlignmentFor,
      });
    }
    if (i < segments.length - 1) y += sectionGap;
  }
}

/**
 * Header top, footer bottom, body + grid stack in the middle band.
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock[]} prepared
 * @param {SlideVariant} variant
 */
function layoutLabeledMessage(ctx, theme, prepared, variant) {
  const { sizeHeight, marginVertical } = theme;
  const canvasHeight = sizeHeight ?? theme.size;
  const headerPad = gapFor(theme, 0.028);
  const footerPad = gapFor(theme, 0.018);

  const header = pick(prepared, 'header')[0];
  const footer = pick(prepared, 'footer')[0];
  const verticalAlign = resolveVerticalAlign(variant);
  const column = contentColumnBounds(theme);

  let bandTop = marginVertical;
  if (header) {
    drawAt(ctx, theme, header, variant, bandTop, column);
    bandTop += measureBlockInkHeight(header) + headerPad;
  }

  refreshAllGridMetrics(prepared);
  const stackHeight = measureMessageStackHeight(prepared, theme);
  let stackTop = bandTop;
  const bandBottom = footer
    ? canvasHeight - marginVertical - measureBlockInkHeight(footer) - footerPad
    : canvasHeight - marginVertical;

  if (verticalAlign === 'center') {
    stackTop = bandTop + Math.max(0, (bandBottom - bandTop - stackHeight) / 2);
  } else if (verticalAlign === 'bottom') {
    stackTop = bandBottom - stackHeight;
  }

  drawMessageStack(ctx, theme, prepared, variant, stackTop, column);

  if (footer) {
    const bottomY = canvasHeight - marginVertical - measureBlockInkHeight(footer);
    drawAt(ctx, theme, footer, variant, bottomY, column);
  }
}

function layoutStackedRhythm(ctx, theme, prepared, variant) {
  layoutLabeledMessage(ctx, theme, prepared, variant);
}

function layoutHeroPunch(ctx, theme, prepared, variant) {
  layoutLabeledMessage(ctx, theme, prepared, variant);
}

function layoutClaimProof(ctx, theme, prepared, variant) {
  layoutLabeledMessage(ctx, theme, prepared, variant);
}

function layoutKeywordAnchor(ctx, theme, prepared, variant) {
  layoutLabeledMessage(ctx, theme, prepared, variant);
}

function layoutClosingThesis(ctx, theme, prepared, variant) {
  layoutLabeledMessage(ctx, theme, prepared, variant);
}

/** @param {SlideVariant} variant */
function variantWithFooterLeft(variant) {
  const base = compactSectionAlignment(variant);
  return {
    ...variant,
    alignment: { ...base, footer: 'left' },
  };
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {import('./assets.js').PostCtaAssets} ctaConfig
 */
function qrPaintColors(theme, ctaConfig) {
  const qr = ctaConfig.qr && typeof ctaConfig.qr === 'object' ? ctaConfig.qr : {};
  const colorSpec = qrModuleColorSpec(qr);
  const lightSpec = qr.light ?? 'transparent';
  return {
    dark: resolveThemedColor(theme, colorSpec, 'accent2'),
    light: isQrBackgroundTransparent(lightSpec)
      ? 'transparent'
      : resolveThemedColor(theme, lightSpec, 'muted'),
  };
}

/**
 * QR left; scan CTA + brand lockup stacked on the right.
 * @returns {boolean} Whether the brand image was drawn here
 */
function layoutPostCtaQrSplit(ctx, theme, variant, options, layout) {
  const {
    ctaConfig,
    column,
    scanLines,
    cursorY,
    sizeHeight,
    marginVertical,
    assets,
    hasBrand,
    brandMaxH,
    footerPad,
  } = layout;

  const bandTop = cursorY;
  const bandBottom = sizeHeight - marginVertical;
  const bandHeight = Math.max(0, bandBottom - bandTop);
  if (bandHeight <= 0) return false;

  const splitGap = scaleCanvasPx(16, theme.size);
  const scanLineGap = scaleCanvasPx(8, theme.size);
  const scanBrandGap = scaleCanvasPx(4, theme.size);
  const rawRatio = Number(ctaConfig.qr?.columnRatio);
  const columnRatio = Number.isFinite(rawRatio) && rawRatio > 0.2 && rawRatio <= 0.55
    ? rawRatio
    : 0.5;
  const leftMaxW = Math.round(column.contentWidth * columnRatio);

  const scanStackH = measureBlocksStackHeight(scanLines, scanLineGap, 0);
  const scanTrailingEm = scanLines.length
    ? measureBlockTrailingEmPad(scanLines[scanLines.length - 1])
    : 0;
  const scanInkH = Math.max(0, scanStackH - scanTrailingEm);

  let rightBrandBoxH = 0;
  let brandDrawW = 0;
  let brandDrawH = 0;
  const rightWidthEstimate = Math.max(
    1,
    column.contentWidth - Math.round(column.contentWidth * columnRatio) - splitGap,
  );
  if (hasBrand && assets.brand) {
    const bnw = assets.brand.naturalWidth || assets.brand.width || 0;
    const bnh = assets.brand.naturalHeight || assets.brand.height || 0;
    if (bnw > 0 && bnh > 0) {
      const fit = fitImageBox(bnw, bnh, rightWidthEstimate, brandMaxH);
      brandDrawW = fit.width;
      brandDrawH = fit.height;
      rightBrandBoxH = fit.height;
    }
  }

  let rightStackH = scanInkH + (scanLines.length && rightBrandBoxH > 0 ? scanBrandGap : 0) + rightBrandBoxH;

  const qrPercent = parseQrSizePercent(ctaConfig.qr?.size);
  const qrMin = scaleCanvasPx(96, theme.size);
  const maxSquare = Math.min(bandHeight, leftMaxW);
  let qrBox = (maxSquare * qrPercent) / 100;
  if (maxSquare >= qrMin) {
    qrBox = Math.max(qrMin, qrBox);
  }
  qrBox = Math.min(qrBox, maxSquare, leftMaxW, bandHeight);

  const rightLeft = column.columnLeft + qrBox + splitGap;
  const rightWidth = Math.max(1, column.contentWidth - qrBox - splitGap);
  const rightColumn = {
    columnLeft: rightLeft,
    contentWidth: rightWidth,
    columnRight: rightLeft + rightWidth,
  };

  if (hasBrand && assets.brand) {
    const bnw = assets.brand.naturalWidth || assets.brand.width || 0;
    const bnh = assets.brand.naturalHeight || assets.brand.height || 0;
    if (bnw > 0 && bnh > 0) {
      const fit = fitImageBox(bnw, bnh, rightWidth, brandMaxH);
      brandDrawW = fit.width;
      brandDrawH = fit.height;
      rightBrandBoxH = fit.height;
      rightStackH = scanInkH + (scanLines.length ? scanBrandGap : 0) + rightBrandBoxH;
    }
  }

  const rightStackTop = bandBottom - rightStackH;

  const splitVariant = variantWithFooterLeft(variant);
  if (scanLines.length) {
    drawBlocksTopDown(ctx, theme, scanLines, splitVariant, rightStackTop, rightColumn, scanLineGap, {
      gapAfterLast: 0,
    });
  }

  if (hasBrand && assets.brand && rightBrandBoxH > 0) {
    const brandY = rightStackTop + scanInkH + (scanLines.length ? scanBrandGap : 0);
    drawImageContainInBox(ctx, assets.brand, rightLeft, brandY, brandDrawW, brandDrawH);
    if (Boolean(activeRenderOptions?.showLineBoxes)) {
      drawImageSlotDebug(ctx, rightLeft, brandY, brandDrawW, brandDrawH, theme.size);
    }
  }

  const qrX = column.columnLeft;
  let qrY = bandBottom - qrBox;
  if (qrY < bandTop) {
    qrBox = Math.max(1, bandBottom - bandTop);
    qrY = bandBottom - qrBox;
  }
  drawQrCodeInBox(ctx, options.qrCode, qrX, qrY, qrBox, {
    marginModules: 2,
    ...qrPaintColors(theme, ctaConfig),
  });
  if (Boolean(activeRenderOptions?.showLineBoxes) && qrBox > 0) {
    drawImageSlotDebug(ctx, qrX, qrY, qrBox, qrBox, theme.size);
  }

  return hasBrand && rightBrandBoxH > 0;
}

/**
 * CTA slide: logo, featured image, title/url text (funnel to full post).
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {PreparedBlock[]} prepared
 * @param {SlideVariant} variant
 * @param {RenderOptions} options
 */
function layoutPostCta(ctx, theme, prepared, variant, options) {
  const { sizeHeight, marginVertical } = theme;
  const canvasHeight = sizeHeight ?? theme.size;
  const column = contentColumnBounds(theme);
  const assets = options.loadedCtaAssets;
  const ctaConfig = mergePostCtaConfig(variant, options.deckCta);
  const footerPad = gapFor(theme, 0.02);
  const sectionGap = scaleCanvasPx(20, theme.size);

  const logoMaxH = scaleCanvasPx(64, theme.size);
  const logoMaxW = scaleCanvasPx(220, theme.size);
  const hasBrand = Boolean(assets?.brand);
  const brandDesignMax = ctaConfig.brandMaxHeight ?? 180;
  const brandMaxH = hasBrand ? scaleCanvasPx(brandDesignMax, theme.size) : 0;
  let brandBoxH = brandMaxH;
  if (hasBrand && assets.brand) {
    const bnw = assets.brand.naturalWidth || assets.brand.width || 0;
    const bnh = assets.brand.naturalHeight || assets.brand.height || 0;
    if (bnw > 0 && bnh > 0) {
      brandBoxH = fitImageBox(bnw, bnh, column.contentWidth, brandMaxH).height;
    }
  }

  const qrText = options.qrText || ctaConfig.shortUrl || ctaConfig.postUrl || '';
  const showQr = typeof qrText === 'string' && qrText.trim().length > 0 && options.qrCode;
  const defaultFeaturedMax = showQr ? (hasBrand ? 200 : 220) : (hasBrand ? 300 : 360);
  const featuredDesignMax = ctaConfig.featuredMaxHeight ?? defaultFeaturedMax;
  const featuredMaxH = scaleCanvasPx(featuredDesignMax, theme.size);
  let featuredBoxH = featuredMaxH;
  if (assets?.featured) {
    const fnw = assets.featured.naturalWidth || assets.featured.width || 0;
    const fnh = assets.featured.naturalHeight || assets.featured.height || 0;
    if (fnw > 0 && fnh > 0) {
      featuredBoxH = fitImageBox(fnw, fnh, column.contentWidth, featuredMaxH).height;
    }
  }

  const qrLayoutRaw = ctaConfig.qr?.layout;
  const useSplitQr = showQr && hasBrand
    && (qrLayoutRaw === 'split' || (qrLayoutRaw !== 'stack' && qrLayoutRaw !== 'stacked'));
  const brandY = hasBrand && !useSplitQr
    ? canvasHeight - marginVertical - brandBoxH
    : canvasHeight - marginVertical;
  const brandGap = 0;
  const contentBottom = hasBrand && !useSplitQr ? brandY - brandGap : canvasHeight - marginVertical;
  let brandDrawnInSplit = false;

  let y = marginVertical;

  if (assets?.logo) {
    drawImageContainCentered(
      ctx,
      assets.logo,
      column.columnLeft + column.contentWidth / 2,
      y,
      logoMaxW,
      logoMaxH,
    );
    y += logoMaxH + sectionGap;
  }

  if (assets?.featured) {
    const fnw = assets.featured.naturalWidth || assets.featured.width || 0;
    const fnh = assets.featured.naturalHeight || assets.featured.height || 0;
    if (fnw > 0 && fnh > 0) {
      drawImageContainInBox(
        ctx,
        assets.featured,
        column.columnLeft,
        y,
        column.contentWidth,
        featuredBoxH,
      );
      y += featuredBoxH + sectionGap;
    } else {
      console.warn('[carousel] post_cta featured image has no dimensions');
    }
  } else if (variant.featuredImage || options.deckCta?.featuredImage) {
    console.warn('[carousel] post_cta featured image did not load');
  }

  const header = pick(prepared, 'header')[0];
  const allFooters = pick(prepared, 'footer');
  const footers = showQr
    ? allFooters.filter((b) => !looksLikeLongUrlText(b.text))
    : allFooters;
  const { urlLines, scanLines } = splitPostCtaFooters(footers);
  const cluster = bodyCluster(prepared);
  const verticalAlign = resolveVerticalAlign(variant);

  let bandTop = y;
  if (header) {
    drawAt(ctx, theme, header, variant, bandTop, column);
    bandTop += measureBlockInkHeight(header) + footerPad;
  }

  if (showQr) {
    const urlToQrGap = scaleCanvasPx(18, theme.size);
    const qrGapAfter = scaleCanvasPx(14, theme.size);

    const { titles, urls: bodyUrls } = splitTitleAndUrlBlocks(cluster);
    const urlBlocks = [...bodyUrls, ...urlLines];
    const headlineBlocks = [...titles, ...urlBlocks];
    const seamlessAfter = titles.length > 0 && urlBlocks.length > 0 ? titles.length - 1 : null;

    let cursorY = bandTop;
    if (headlineBlocks.length) {
      cursorY = drawClusterLines(ctx, theme, headlineBlocks, variant, bandTop, column, {
        seamlessAfterBlockIndex: seamlessAfter,
      });
      cursorY += urlToQrGap;
    }

    if (useSplitQr) {
      brandDrawnInSplit = layoutPostCtaQrSplit(ctx, theme, variant, options, {
        ctaConfig,
        column,
        scanLines,
        cursorY,
        sizeHeight: canvasHeight,
        marginVertical,
        assets,
        hasBrand,
        brandMaxH,
        footerPad,
      });
    } else {
    const scanStackH = measureBlocksStackHeight(scanLines, footerPad, 0);
    const scanEmBottom = scanLines.length
      ? brandY + measureBlockTrailingEmPad(scanLines[scanLines.length - 1])
      : contentBottom;

    // Bottom stack: tuck scan footer descender into og-logo top padding (no brandGap).
    let stackBottom = scanEmBottom;
    if (scanLines.length) {
      stackBottom -= scanStackH;
      drawBlocksTopDown(ctx, theme, scanLines, variant, stackBottom, column, footerPad, {
        gapAfterLast: 0,
      });
      stackBottom -= qrGapAfter;
    } else {
      stackBottom -= qrGapAfter;
    }

    const qrPercent = parseQrSizePercent(ctaConfig.qr?.size);
    const availableH = Math.max(0, stackBottom - cursorY);
    const qrMin = scaleCanvasPx(96, theme.size);
    const maxSquare = Math.min(availableH, column.contentWidth);
    let qrBox = (maxSquare * qrPercent) / 100;
    if (maxSquare >= qrMin) {
      qrBox = Math.max(qrMin, qrBox);
    }
    qrBox = Math.min(qrBox, maxSquare, availableH);
    let qrY = cursorY + (availableH - qrBox) / 2;
    const minQrY = cursorY;
    if (qrY < minQrY) {
      qrY = minQrY;
      qrBox = Math.max(1, Math.min(stackBottom - qrY, maxSquare));
    }

    const qrX = column.columnLeft + (column.contentWidth - qrBox) / 2;
    drawQrCodeInBox(ctx, options.qrCode, qrX, qrY, qrBox, {
      marginModules: 2,
      ...qrPaintColors(theme, ctaConfig),
    });
    if (Boolean(activeRenderOptions?.showLineBoxes) && qrBox > 0) {
      drawImageSlotDebug(ctx, qrX, qrY, qrBox, qrBox, theme.size);
    }
    }
  } else {
    let bandBottom = contentBottom;
    if (hasBrand && footers.length) {
      bandBottom += measureBlockTrailingEmPad(footers[0]);
    }
    for (let i = footers.length - 1; i >= 0; i -= 1) {
      const footer = footers[i];
      bandBottom -= measureBlockInkHeight(footer);
      drawAt(ctx, theme, footer, variant, bandBottom, column);
      bandBottom -= footerPad;
    }

    layoutGroupedStack(ctx, theme, cluster, variant, bandTop, bandBottom, column, { align: verticalAlign });
  }

  if (hasBrand && assets.brand && !brandDrawnInSplit) {
    drawImageContainInBox(
      ctx,
      assets.brand,
      column.columnLeft,
      brandY,
      column.contentWidth,
      brandBoxH,
    );
    if (Boolean(activeRenderOptions?.showLineBoxes)) {
      const bnw = assets.brand.naturalWidth || assets.brand.width || 0;
      const bnh = assets.brand.naturalHeight || assets.brand.height || 0;
      if (bnw > 0 && bnh > 0) {
        const fit = fitImageBox(bnw, bnh, column.contentWidth, brandBoxH);
        const boxX = column.columnLeft + (column.contentWidth - fit.width) / 2;
        const boxY = brandY + (brandBoxH - fit.height) / 2;
        drawImageSlotDebug(ctx, boxX, boxY, fit.width, fit.height, theme.size);
      }
    }
  } else if (variant.brandImage || options.deckCta?.brandImage) {
    console.warn('[carousel] post_cta brand image did not load');
  }
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {SlideVariant} variant
 * @param {RenderOptions} [options]
 */
function renderVariant(ctx, theme, variant, options = {}) {
  activeRenderOptions = options;
  try {
    if (options.skipBackground) {
      ctx.clearRect(0, 0, ctx.canvas.width, ctx.canvas.height);
    }
    drawBackground(ctx, theme, options);

    const motifConfig = parseMotifStripConfig(
      /** @type {Record<string, unknown>} */ (options.deck?.deck)?.motifStrip,
    );
    if (
      motifConfig
      && !options.skipMotifStrip
      && options.motifStripContext
      && options.loadedMotifStrip
      && motifStripEnabledForSlide(motifConfig, options.slideRole)
    ) {
      paintMotifStripSlice(
        ctx,
        theme.size,
        theme.sizeHeight ?? theme.size,
        options.loadedMotifStrip,
        options.motifStripContext,
        motifConfig,
        theme,
      );
    }

    if (options.showLineBoxes) {
      drawCanvasBoundsDebug(ctx);
    }

    const prepared = prepareVariant(ctx, theme, variant);
    const archetype = variant.archetype || 'stacked_rhythm';

    switch (archetype) {
    case 'hero_punch':
      layoutHeroPunch(ctx, theme, prepared, variant);
      break;
    case 'claim_proof':
      layoutClaimProof(ctx, theme, prepared, variant);
      break;
    case 'keyword_anchor':
      layoutKeywordAnchor(ctx, theme, prepared, variant);
      break;
    case 'closing_thesis':
      layoutClosingThesis(ctx, theme, prepared, variant);
      break;
    case 'post_cta':
      layoutPostCta(ctx, theme, prepared, variant, options);
      break;
    case 'labeled_message':
      layoutLabeledMessage(ctx, theme, prepared, variant);
      break;
    case 'stacked_rhythm':
    default:
      layoutStackedRhythm(ctx, theme, prepared, variant);
      break;
    }

  } finally {
    activeRenderOptions = null;
  }
}

/**
 * Debug: grid row metrics after prepare (for studio / Playwright QA).
 * @param {CanvasRenderingContext2D} ctx
 * @param {import('./theme.js').DeckTheme} theme
 * @param {SlideVariant} variant
 */
export function inspectVariantPrepare(ctx, theme, variant) {
  const prepared = prepareVariant(ctx, theme, variant);
  const grid = /** @type {import('./text-grid.js').PreparedGridSection|undefined} */ (
    prepared.find((item) => item.section === 'grid')
  );
  if (!grid) {
    return { preparedSections: prepared.map((item) => item.section) };
  }
  return {
    stackHeight: measureMessageStackHeight(prepared, theme),
    stackBudget: messageStackFitBudget(theme, prepared),
    rowHeights: grid.rowHeights,
    rowGapPx: grid.rowGapPx,
    colWidths: grid.colWidths,
    totalHeight: grid.totalHeight,
    cells: grid.cells.map((entry) => ({
      row: entry.row,
      col: entry.col,
      rowSpan: entry.rowSpan,
      fontSizePx: entry.prepared.fontSizePx,
      lineSlotPx: entry.prepared.lineSlotPx,
      lineCount: entry.prepared.inlineLines?.length ?? 0,
      inkExtent: measurePreparedInkHeight(entry.prepared),
      text: entry.prepared.block?.text,
    })),
  };
}

export async function renderSlideToCanvas(variant, themeOverrides = {}, renderOptions = {}) {
  const { grain = true, targetCanvas, outputSize, supersample = 1 } = renderOptions;
  const spec = /** @type {Partial<import('./theme.js').DeckTheme> & Record<string, unknown>} */ (themeOverrides);
  const baseTheme = mergeTheme(spec);
  const finalWidth = clampSlideEdge(outputSize ?? baseTheme.size);
  const aspect = {
    width: baseTheme.aspectRatioWidth,
    height: baseTheme.aspectRatioHeight,
  };
  const finalHeight = canvasHeightForWidth(finalWidth, aspect);
  const sampleScale = Math.max(1, supersample);
  const renderWidth = Math.round(finalWidth * sampleScale);
  const renderHeight = Math.round(finalHeight * sampleScale);

  /** @type {RenderOptions} */
  const renderOpts = { ...renderOptions };

  let theme = renderWidth === baseTheme.size && renderHeight === baseTheme.sizeHeight
    ? baseTheme
    : mergeTheme({ ...spec, size: renderWidth }, { canvasSizeMax: renderWidth });

  if (variant.archetype === 'post_cta') {
    const ctaConfig = mergePostCtaConfig(variant, renderOpts.deckCta);
    if (!renderOpts.slideRole) {
      renderOpts.slideRole = 'cta';
    }
    theme = mergeTheme(
      {
        ...spec,
        size: theme.size,
        background: ctaConfig.background,
        backgroundGradient: 'solid',
        backgroundGradientPreset: null,
        backgroundWave: { style: 'none' },
      },
      { canvasSizeMax: renderWidth },
    );
    renderOpts.grain = false;
    renderOpts.loadedCtaAssets = await loadPostCtaAssets(
      ctaConfig,
      renderOpts.assetBaseUrl,
      (token) => resolveColor(theme, token),
      renderOpts.bundleBaseUrl,
    );
    const qrTarget = ctaConfig.shortUrl || ctaConfig.postUrl;
    if (qrTarget) {
      renderOpts.qrText = qrTarget;
      await ensureQrCodegenLoaded();
      renderOpts.qrCode = makeQrCode(qrTarget, 'MEDIUM');
    }
  }

  await loadDeckFonts(theme);

  const motifConfig = parseMotifStripConfig(
    /** @type {Record<string, unknown>} */ (renderOpts.deck?.deck)?.motifStrip,
  );
  if (motifConfig?.src && renderOpts.motifStripContext) {
    try {
      renderOpts.loadedMotifStrip = renderOpts.loadedMotifStrip
        ?? await loadMotifStripImage(
          motifConfig,
          renderOpts.assetBaseUrl ?? '',
          (token) => resolveColor(theme, token),
        );
    } catch (error) {
      console.warn('[carousel] Motif strip failed to load:', error);
    }
  }

  /** @type {HTMLCanvasElement} */
  let rasterCanvas;
  if (sampleScale > 1) {
    rasterCanvas = document.createElement('canvas');
    rasterCanvas.width = renderWidth;
    rasterCanvas.height = renderHeight;
  } else if (targetCanvas) {
    rasterCanvas = targetCanvas;
    rasterCanvas.width = finalWidth;
    rasterCanvas.height = finalHeight;
  } else {
    rasterCanvas = document.createElement('canvas');
    rasterCanvas.width = finalWidth;
    rasterCanvas.height = finalHeight;
  }

  const ctx = rasterCanvas.getContext('2d');
  if (!ctx) throw new Error('Canvas 2D context unavailable');

  ctx.imageSmoothingEnabled = true;
  ctx.imageSmoothingQuality = 'high';
  renderVariant(ctx, theme, variant, renderOpts);

  if (sampleScale > 1) {
    const output = targetCanvas ?? document.createElement('canvas');
    output.width = finalWidth;
    output.height = finalHeight;
    const outCtx = output.getContext('2d');
    if (!outCtx) throw new Error('Canvas 2D context unavailable');
    outCtx.imageSmoothingEnabled = true;
    outCtx.imageSmoothingQuality = 'high';
    outCtx.drawImage(
      rasterCanvas,
      0,
      0,
      renderWidth,
      renderHeight,
      0,
      0,
      finalWidth,
      finalHeight,
    );
    return output;
  }

  return rasterCanvas;
}

/** @param {HTMLElement} frame @param {number} [fallback] */
export function previewPixelSize(frame, fallback = 360) {
  const dpr = Math.min(window.devicePixelRatio || 1, 2);
  const displayWidth = frame.clientWidth || frame.offsetWidth || fallback;
  return Math.max(1, Math.round(displayWidth * dpr * 2));
}

/**
 * Preview raster dimensions matching deck aspect ratio (default 4:5).
 * @param {HTMLElement} frame
 * @param {import('./theme.js').AspectRatio} [aspect]
 * @param {number} [fallback]
 */
export function previewPixelDimensions(frame, aspect, fallback = 360) {
  const width = previewPixelSize(frame, fallback);
  const ratio = aspect ?? {
    width: DEFAULT_ASPECT_RATIO_WIDTH,
    height: DEFAULT_ASPECT_RATIO_HEIGHT,
  };
  return {
    width,
    height: canvasHeightForWidth(width, ratio),
  };
}

/**
 * Downscale a full-size slide canvas into a preview element (legacy blit path).
 * Prefer rendering with `renderSlideToCanvas(..., { targetCanvas, outputSize })` instead.
 */
export function paintPreviewCanvas(preview, source) {
  const pixelWidth = previewPixelSize(preview);
  const pixelHeight = Math.max(
    1,
    Math.round(pixelWidth * (source.height / source.width)),
  );

  preview.width = pixelWidth;
  preview.height = pixelHeight;

  const ctx = preview.getContext('2d');
  if (!ctx) return;

  ctx.imageSmoothingEnabled = true;
  ctx.imageSmoothingQuality = 'high';
  ctx.clearRect(0, 0, pixelWidth, pixelHeight);
  ctx.drawImage(source, 0, 0, pixelWidth, pixelHeight);
}
