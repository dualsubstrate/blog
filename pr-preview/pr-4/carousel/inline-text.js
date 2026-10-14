import {
  ASCENDER_PROBE_CHARS,
  assertProbeSetsDisjoint,
  CAP_HEIGHT_PROBE_CHARS,
  DESCENDER_PROBE_CHARS,
  EM_ORANGE_INSET_PX,
  EM_STRUT_CHARS,
  EM_STRUT_STRINGS,
  isXHeightProbeChar,
  textHasLowercase,
  validateFontLineMetrics,
  X_HEIGHT_PROBE_CHARS,
} from './typo-probes.js';
import { columnEdgeBleedPx } from './theme.js';

export {
  ASCENDER_PROBE_CHARS,
  assertProbeSetsDisjoint,
  CAP_HEIGHT_PROBE_CHARS,
  DESCENDER_PROBE_CHARS,
  EM_ORANGE_INSET_PX,
  EM_STRUT_CHARS,
  EM_STRUT_STRINGS,
  findProbeSetOverlaps,
  isAscenderProbeChar,
  isDescenderProbeChar,
  isXHeightProbeChar,
  TYPO_PROBE_BY_ROLE,
  validateFontLineMetrics,
  baselineYFromSymmetricOrange,
  lineSlotHeightPx,
  orangeBandHeightPx,
  symmetricOrangeInsetPx,
  typoOrangeLinesInEmBox,
  X_HEIGHT_PROBE_CHARS,
} from './typo-probes.js';

/**
 * Minimal inline Markdown in block `text`: `**bold**`, `*italic*`, `***both***`.
 * Optional palette tokens: `<accent1>…</accent1>`, `<color accent1>…</color>`, optional
 * `brightness="+10"` or `brightness='-20'` (percent lighten/darken vs that token's base hex).
 * Unclosed markers stay literal. No links, code, or nested spans.
 *
 * @typedef {Object} InlineRun
 * @property {string} text
 * @property {boolean} [bold]
 * @property {boolean} [italic]
 * @property {string} [color] `text`, `muted`, `accent1`, `accent2`; falls back to block color
 * @property {number} [brightness] -100..100; applied at render in `theme.resolveInlineRunColor`
 */

/** @type {RegExp} */
const INLINE_COLOR_SEGMENT_RE =
  /<(?:color\s+)?(accent1|accent2|text|muted)(?:\s+brightness=["']([+-]?\d{1,3})["'])?\s*>([\s\S]*?)<\/(?:\1|color)>/gi;

const INLINE_COLOR_TAG_PROBE_RE = /<(?:color\s+)?(?:accent1|accent2|text|muted)\b/i;

/** @type {Set<string>} */
const INLINE_COLOR_TOKENS = new Set(['accent1', 'accent2', 'text', 'muted']);

/**
 * @param {string|undefined} raw
 * @returns {number|undefined}
 */
export function parseInlineBrightness(raw) {
  if (raw == null || raw === '') return undefined;
  const value = Number.parseInt(raw, 10);
  if (!Number.isFinite(value)) return undefined;
  return Math.max(-100, Math.min(100, value));
}

/**
 * @param {string} text
 * @returns {InlineRun[]}
 */
export function parseInlineMarkdown(text) {
  if (!text) return [{ text: '' }];
  if (!text.includes('*')) {
    return [{ text }];
  }

  /** @type {InlineRun[]} */
  const runs = [];
  let plain = '';
  let i = 0;

  /** @param {string} chunk @param {boolean} [bold] @param {boolean} [italic] */
  function pushRun(chunk, bold = false, italic = false) {
    if (!chunk) return;
    runs.push(bold || italic ? { text: chunk, bold, italic } : { text: chunk });
  }

  function flushPlain() {
    if (plain) {
      pushRun(plain);
      plain = '';
    }
  }

  /** @param {string} delimiter @param {number} markerLen @param {boolean} bold @param {boolean} italic */
  function tryDelimited(delimiter, markerLen, bold, italic) {
    const start = i + markerLen;
    const end = text.indexOf(delimiter, start);
    if (end === -1) return false;
    flushPlain();
    pushRun(text.slice(start, end), bold, italic);
    i = end + markerLen;
    return true;
  }

  while (i < text.length) {
    if (text.startsWith('***', i) && tryDelimited('***', 3, true, true)) continue;
    if (text.startsWith('**', i) && tryDelimited('**', 2, true, false)) continue;
    if (text[i] === '*' && text[i + 1] !== '*' && tryDelimited('*', 1, false, true)) continue;
    plain += text[i];
    i += 1;
  }

  flushPlain();
  return runs.length > 0 ? runs : [{ text }];
}

/**
 * Markdown plus inline color tags (see {@link parseInlineMarkdown}).
 * @param {string} text
 * @returns {InlineRun[]}
 */
export function parseInlineRichText(text) {
  if (!text) return [{ text: '' }];
  if (!INLINE_COLOR_TAG_PROBE_RE.test(text)) {
    return parseInlineMarkdown(text);
  }

  INLINE_COLOR_SEGMENT_RE.lastIndex = 0;
  /** @type {InlineRun[]} */
  const runs = [];
  let last = 0;
  let match = INLINE_COLOR_SEGMENT_RE.exec(text);
  while (match) {
    if (match.index > last) {
      runs.push(...parseInlineMarkdown(text.slice(last, match.index)));
    }
    const colorToken = match[1].toLowerCase();
    const brightness = parseInlineBrightness(match[2]);
    if (INLINE_COLOR_TOKENS.has(colorToken)) {
      const inner = parseInlineMarkdown(match[3]);
      for (const run of inner) {
        runs.push({
          ...run,
          color: colorToken,
          ...(brightness !== undefined ? { brightness } : {}),
        });
      }
    } else {
      runs.push(...parseInlineMarkdown(match[0]));
    }
    last = match.index + match[0].length;
    match = INLINE_COLOR_SEGMENT_RE.exec(text);
  }
  if (last < text.length) {
    runs.push(...parseInlineMarkdown(text.slice(last)));
  }
  return runs.length > 0 ? runs : [{ text: '' }];
}

/**
 * @param {InlineRun[]} runs
 * @returns {InlineRun[]}
 */
export function runsToWordTokens(runs) {
  /** @type {InlineRun[]} */
  const tokens = [];
  for (const run of runs) {
    if (!run.text) continue;
    // Keep whitespace between styled runs (e.g. "**candor** and" → "candor" + " " + "and").
    const chunks = run.text.match(/\S+\s*|\s+/g);
    if (!chunks) continue;
    for (const chunk of chunks) {
      tokens.push({
        text: chunk,
        bold: run.bold,
        italic: run.italic,
        ...(run.color ? { color: run.color } : {}),
        ...(run.brightness !== undefined ? { brightness: run.brightness } : {}),
      });
    }
  }
  return tokens;
}

/**
 * Split on `\n`. Trim leading/trailing empty paragraphs; keep internal blank lines.
 * @param {string} text
 * @returns {string[]}
 */
export function splitTextParagraphs(text) {
  const parts = text.split('\n');
  while (parts.length > 0 && parts[parts.length - 1].trim() === '') {
    parts.pop();
  }
  while (parts.length > 0 && parts[0].trim() === '') {
    parts.shift();
  }
  return parts;
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {string} text
 * @param {number} maxWidth
 * @returns {InlineRun[][]}
 */
export function wrapInlineText(ctx, fontFor, text, maxWidth) {
  const paragraphs = splitTextParagraphs(text);
  /** @type {InlineRun[][]} */
  const lines = [];

  for (const paragraph of paragraphs) {
    if (paragraph.trim() === '') {
      lines.push([]);
      continue;
    }
    const parsed = parseInlineRichText(paragraph);
    const tokens = runsToWordTokens(parsed);
    if (tokens.length === 0) {
      lines.push([]);
      continue;
    }
    lines.push(...wrapInlineWords(ctx, fontFor, tokens, maxWidth));
  }

  return lines;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {import('./theme.js').FontRole} role
 * @param {number} fontSizePx
 * @param {string} baseWeight
 * @param {InlineRun} token
 */
export function buildInlineFont(theme, role, fontSizePx, baseWeight, token) {
  const family = role === 'display' ? theme.displayFont : theme.bodyFont;
  const weight = token.bold ? '700' : baseWeight;
  const style = token.italic ? 'italic' : 'normal';
  return `${style} ${weight} ${fontSizePx}px "${family}", serif`;
}

/** Letters used to sample ascender ink when actualBoundingBox is unavailable. */
const EM_INK_SAMPLE_CHARS = 'ÉHÓbdfhkltgjpqy';

/** @param {string} sample */
function sampleText(sample) {
  return sample.trim() === '' ? ' ' : sample;
}

/**
 * Ascent for positioning: actual glyph ink only (never fontBoundingBox, which overshoots cap ink).
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {InlineRun} token
 * @param {number} fontSizePx
 */
function inkAscentForToken(ctx, fontFor, token, fontSizePx) {
  ctx.font = fontFor(token);
  const text = sampleText(token.text);
  const metrics = ctx.measureText(text);
  let ascent = metrics.actualBoundingBoxAscent;
  if (Number.isFinite(ascent) && ascent > 0) return ascent;

  for (const ch of text) {
    if (/\s/.test(ch)) continue;
    const m = ctx.measureText(ch);
    if (Number.isFinite(m.actualBoundingBoxAscent) && m.actualBoundingBoxAscent > 0) {
      ascent = Math.max(ascent || 0, m.actualBoundingBoxAscent);
    }
  }
  if (Number.isFinite(ascent) && ascent > 0) return ascent;

  for (const ch of EM_INK_SAMPLE_CHARS) {
    const m = ctx.measureText(ch);
    if (Number.isFinite(m.actualBoundingBoxAscent) && m.actualBoundingBoxAscent > 0) {
      return m.actualBoundingBoxAscent;
    }
  }
  return Math.round(fontSizePx * 0.72);
}

/**
 * Descent for ink sizing: actual glyph ink only.
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {InlineRun} token
 * @param {number} fontSizePx
 */
function inkDescentForToken(ctx, fontFor, token, fontSizePx) {
  ctx.font = fontFor(token);
  const text = sampleText(token.text);
  const metrics = ctx.measureText(text);
  let descent = metrics.actualBoundingBoxDescent;
  if (Number.isFinite(descent) && descent > 0) return descent;

  for (const ch of text) {
    if (/\s/.test(ch)) continue;
    const m = ctx.measureText(ch);
    if (Number.isFinite(m.actualBoundingBoxDescent) && m.actualBoundingBoxDescent > 0) {
      descent = Math.max(descent || 0, m.actualBoundingBoxDescent);
    }
  }
  if (Number.isFinite(descent) && descent > 0) return descent;

  for (const ch of 'gjpqy') {
    const m = ctx.measureText(ch);
    if (Number.isFinite(m.actualBoundingBoxDescent) && m.actualBoundingBoxDescent > 0) {
      return m.actualBoundingBoxDescent;
    }
  }
  return Math.round(fontSizePx * 0.2);
}

/**
 * Glyph ink from actualBoundingBox only (width + ink box for metrics).
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {InlineRun} token
 * @param {number} [fontSizePx]
 */
function measureTokenInkActual(ctx, fontFor, token, fontSizePx = 80) {
  ctx.font = fontFor(token);
  const text = sampleText(token.text);
  const metrics = ctx.measureText(text);
  const ascent = inkAscentForToken(ctx, fontFor, token, fontSizePx);
  const descent = inkDescentForToken(ctx, fontFor, token, fontSizePx);
  return {
    ascent,
    descent,
    width: metrics.width,
  };
}

/**
 * Horizontal ink span when tokens are drawn left-to-right from x=0 (alphabetic baseline).
 * Uses actualBoundingBox so right/center alignment does not clip serif or display ink.
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {InlineRun[]} tokens
 * @param {number} [fontSizePx]
 * @returns {{ inkLeft: number, inkRight: number, inkWidth: number }}
 */
export function measureInlineLineInkExtents(ctx, fontFor, tokens, fontSizePx = 80) {
  if (tokens.length === 0) {
    return { inkLeft: 0, inkRight: 0, inkWidth: 0 };
  }

  let cursor = 0;
  let inkLeft = 0;
  let inkRight = 0;

  for (const token of tokens) {
    ctx.font = fontFor(token);
    const text = sampleText(token.text);
    const metrics = ctx.measureText(text);
    const advance = metrics.width;
    let boxLeft = Number.isFinite(metrics.actualBoundingBoxLeft) ? metrics.actualBoundingBoxLeft : 0;
    let boxRight = Number.isFinite(metrics.actualBoundingBoxRight) ? metrics.actualBoundingBoxRight : advance;
    inkLeft = Math.min(inkLeft, cursor - boxLeft);
    inkRight = Math.max(inkRight, cursor + boxRight);
    cursor += advance;

    let charCursor = cursor - advance;
    for (const ch of text) {
      const chMetrics = ctx.measureText(ch);
      const chAdvance = chMetrics.width;
      const chLeft = Number.isFinite(chMetrics.actualBoundingBoxLeft) ? chMetrics.actualBoundingBoxLeft : 0;
      const chRight = Number.isFinite(chMetrics.actualBoundingBoxRight) ? chMetrics.actualBoundingBoxRight : chAdvance;
      inkLeft = Math.min(inkLeft, charCursor - chLeft);
      inkRight = Math.max(inkRight, charCursor + chRight);
      charCursor += chAdvance;
    }
  }

  // Small preview rasters and font hinting can under-report ink bounds by a couple pixels.
  // Add extra safety pad so the right-aligned placement stays inside the bitmap.
  const pad = Math.max(4, Math.round(fontSizePx * 0.06));
  return {
    inkLeft: inkLeft - pad,
    inkRight: inkRight + pad,
    inkWidth: inkRight - inkLeft + pad * 2,
  };
}

/**
 * Full em box (actual ink expanded to font bounding box). Used for line-height strut only.
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {InlineRun} token
 */
function measureTokenInkFontBox(ctx, fontFor, token) {
  ctx.font = fontFor(token);
  const sample = token.text.trim() === '' ? ' ' : token.text;
  const metrics = ctx.measureText(sample);
  let ascent = metrics.actualBoundingBoxAscent;
  let descent = metrics.actualBoundingBoxDescent;
  const fontAscent = metrics.fontBoundingBoxAscent;
  const fontDescent = metrics.fontBoundingBoxDescent;
  if (!Number.isFinite(ascent) || !Number.isFinite(descent) || ascent + descent <= 0) {
    ascent = fontAscent;
    descent = fontDescent;
  }
  if (Number.isFinite(fontAscent)) ascent = Math.max(ascent, fontAscent);
  if (Number.isFinite(fontDescent)) descent = Math.max(descent, fontDescent);
  if (!Number.isFinite(ascent)) ascent = 0;
  if (!Number.isFinite(descent)) descent = 0;
  return {
    ascent,
    descent,
    width: metrics.width,
  };
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {InlineRun[]} tokens
 */
export function measureInlineLine(ctx, fontFor, tokens, fontSizePx = 80) {
  if (tokens.length === 0) {
    const ink = measureTokenInkActual(ctx, fontFor, { text: ' ' }, fontSizePx);
    const inkHeight = Math.max(1, Math.round(ink.ascent + ink.descent));
    return {
      ascent: ink.ascent,
      descent: ink.descent,
      inkHeight,
      width: 0,
      inkWidth: 0,
      inkLeft: 0,
      inkRight: 0,
    };
  }

  let maxAscent = 0;
  let maxDescent = 0;

  for (const token of tokens) {
    const ink = measureTokenInkActual(ctx, fontFor, token, fontSizePx);
    maxAscent = Math.max(maxAscent, ink.ascent);
    maxDescent = Math.max(maxDescent, ink.descent);
  }

  const { inkLeft, inkRight, inkWidth } = measureInlineLineInkExtents(ctx, fontFor, tokens, fontSizePx);
  const inkHeight = Math.max(1, Math.round(maxAscent + maxDescent));
  return {
    ascent: maxAscent,
    descent: maxDescent,
    inkHeight,
    width: Math.ceil(inkWidth),
    inkWidth,
    inkLeft,
    inkRight,
  };
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {InlineRun} token
 * @param {{ ascent: number, descent: number }} acc
 */
function accumulateEmBoxMetrics(ctx, fontFor, token, acc) {
  ctx.font = fontFor(token);
  const text = sampleText(token.text);
  const metrics = ctx.measureText(text);
  const ascent = Math.max(glyphInkMetric(metrics, 'ascent'), fontBoxMetric(metrics, 'ascent'));
  const descent = Math.max(glyphInkMetric(metrics, 'descent'), fontBoxMetric(metrics, 'descent'));
  if (Number.isFinite(ascent) && ascent > 0) {
    acc.ascent = Math.max(acc.ascent, ascent);
  }
  if (Number.isFinite(descent) && descent > 0) {
    acc.descent = Math.max(acc.descent, descent);
  }
  for (const ch of text) {
    if (/\s/.test(ch)) continue;
    ctx.font = fontFor({ ...token, text: ch });
    const m = ctx.measureText(ch);
    const inkAsc = glyphInkMetric(m, 'ascent');
    const inkDesc = glyphInkMetric(m, 'descent');
    const boxAsc = fontBoxMetric(m, 'ascent');
    const boxDesc = fontBoxMetric(m, 'descent');
    const chAsc = Math.max(inkAsc, boxAsc);
    const chDesc = Math.max(inkDesc, boxDesc);
    if (chAsc > 0) {
      acc.ascent = Math.max(acc.ascent, chAsc);
    }
    if (chDesc > 0) {
      acc.descent = Math.max(acc.descent, chDesc);
    }
  }
}

/**
 * Font em square for a block (~1em at this block's fontSize).
 * Probes caps, ascenders, and descenders so every line reserves full alphabet depth.
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} fontSizePx
 * @param {boolean[]} [boldProbes] Which inline weights to measure (default both).
 */
export function measureFontStrut(ctx, fontFor, fontSizePx, boldProbes = [false, true]) {
  /** @type {{ ascent: number, descent: number }} */
  const acc = { ascent: 0, descent: 0 };

  for (const bold of boldProbes) {
    for (const text of EM_STRUT_STRINGS) {
      accumulateEmBoxMetrics(ctx, fontFor, { text, bold }, acc);
    }
    for (const ch of EM_STRUT_CHARS) {
      accumulateEmBoxMetrics(ctx, fontFor, { text: ch, bold }, acc);
    }
  }

  if (acc.ascent + acc.descent <= 0) {
    acc.ascent = Math.round(fontSizePx * 0.8);
    acc.descent = Math.round(fontSizePx * 0.2);
  }

  const inkHeight = Math.max(1, Math.round(acc.ascent + acc.descent));
  return { ascent: acc.ascent, descent: acc.descent, inkHeight };
}

/**
 * @param {TextMetrics} metrics
 * @param {'ascent'|'descent'} side
 */
function fontBoxMetric(metrics, side) {
  let v = side === 'ascent' ? metrics.fontBoundingBoxAscent : metrics.fontBoundingBoxDescent;
  if (!Number.isFinite(v) || v <= 0) {
    v = side === 'ascent' ? metrics.actualBoundingBoxAscent : metrics.actualBoundingBoxDescent;
  }
  return Number.isFinite(v) && v > 0 ? v : 0;
}

/**
 * Glyph ink only (actualBoundingBox). fontBoundingBoxAscent is font-wide and
 * would place the ascender line on the em box top.
 * @param {TextMetrics} metrics
 * @param {'ascent'|'descent'} side
 */
function glyphInkMetric(metrics, side) {
  const v = side === 'ascent' ? metrics.actualBoundingBoxAscent : metrics.actualBoundingBoxDescent;
  return Number.isFinite(v) && v > 0 ? v : 0;
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {boolean[]} boldProbes
 * @param {'ascent'|'descent'} side
 */
function maxInkFromXHeightProbes(ctx, fontFor, boldProbes, side) {
  return maxGlyphInkFromProbeChars(ctx, fontFor, X_HEIGHT_PROBE_CHARS, boldProbes, side);
}

function maxInkFromAscenderProbes(ctx, fontFor, boldProbes, side) {
  return maxGlyphInkFromProbeChars(ctx, fontFor, ASCENDER_PROBE_CHARS, boldProbes, side);
}

function maxInkFromDescenderProbes(ctx, fontFor, boldProbes, side) {
  return maxGlyphInkFromProbeChars(ctx, fontFor, DESCENDER_PROBE_CHARS, boldProbes, side);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {string} probeChars Must come from typo-probes.js (one role only).
 * @param {boolean[]} boldProbes
 * @param {'ascent'|'descent'} side
 */
function maxGlyphInkFromProbeChars(ctx, fontFor, probeChars, boldProbes, side) {
  let max = 0;
  for (const bold of boldProbes) {
    for (const ch of probeChars) {
      ctx.font = fontFor({ text: ch, bold });
      max = Math.max(max, glyphInkMetric(ctx.measureText(ch), side));
    }
  }
  return max;
}

/**
 * Distance from baseline to ascender line (ASCENDER_PROBE_CHARS only).
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} fontSizePx
 * @param {boolean[]} [boldProbes]
 */
export function measureFontAscenderPx(ctx, fontFor, fontSizePx, boldProbes = [false]) {
  const max = maxInkFromAscenderProbes(ctx, fontFor, boldProbes, 'ascent');
  return max > 0 ? max : Math.round(fontSizePx * 0.75);
}

/**
 * Distance from baseline to descender line (DESCENDER_PROBE_CHARS only).
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} fontSizePx
 * @param {boolean[]} [boldProbes]
 */
export function measureFontDescenderPx(ctx, fontFor, fontSizePx, boldProbes = [false]) {
  const max = maxInkFromDescenderProbes(ctx, fontFor, boldProbes, 'descent');
  return max > 0 ? max : Math.round(fontSizePx * 0.25);
}

/**
 * Keep orange lines strictly inside the em box (never on blue top/bottom).
 * @param {number} linePx
 * @param {number} emSidePx
 */
function clampTypoLineInsideEm(linePx, emSidePx) {
  const max = Math.max(1, emSidePx - EM_ORANGE_INSET_PX);
  return Math.min(Math.max(1, linePx), max);
}

/**
 * @typedef {Object} FontLineMetrics
 * @property {number} emBoxAscent
 * @property {number} emBoxDescent
 * @property {number} emBoxHeight
 * @property {number} ascenderLinePx Baseline to ascender line (`ASCENDER_PROBE_CHARS` only)
 * @property {number} descenderLinePx Baseline to descender line (`DESCENDER_PROBE_CHARS` only)
 * @property {number} xHeightPx Baseline to meanline (`X_HEIGHT_PROBE_CHARS` only)
 * @property {number} capHeightPx
 */

/**
 * @param {import('./theme.js').FontRole} role
 * @param {number} fontSizePx
 * @param {string} weight
 * @param {boolean[]} boldProbes
 */
export function fontMetricsCacheKey(role, fontSizePx, weight, boldProbes) {
  const boldKey = boldProbes.map((b) => (b ? 'b' : 'n')).join('');
  return `${role}|${fontSizePx}|${weight}|${boldKey}`;
}

/**
 * Font-level em box and ascender/descender lines (not per-line ink).
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} fontSizePx
 * @param {boolean[]} boldProbes
 * @returns {FontLineMetrics}
 */
export function measureFontLineMetrics(ctx, fontFor, fontSizePx, boldProbes) {
  const emBox = measureFontStrut(ctx, fontFor, fontSizePx, boldProbes);
  const xHeightPx = measureFontXHeight(ctx, fontFor, fontSizePx, boldProbes);
  const capHeightPx = measureFontCapHeight(ctx, fontFor, fontSizePx, boldProbes);
  let ascenderLinePx = measureFontAscenderPx(ctx, fontFor, fontSizePx, boldProbes);
  let descenderLinePx = measureFontDescenderPx(ctx, fontFor, fontSizePx, boldProbes);

  if (ascenderLinePx >= emBox.ascent - EM_ORANGE_INSET_PX) {
    ascenderLinePx = capHeightPx > 0 ? capHeightPx : Math.round(emBox.ascent * 0.78);
  }
  ascenderLinePx = clampTypoLineInsideEm(ascenderLinePx, emBox.ascent);
  descenderLinePx = clampTypoLineInsideEm(descenderLinePx, emBox.descent);

  const metrics = {
    emBoxAscent: emBox.ascent,
    emBoxDescent: emBox.descent,
    emBoxHeight: emBox.inkHeight,
    ascenderLinePx,
    descenderLinePx,
    xHeightPx,
    capHeightPx,
  };

  if (typeof console !== 'undefined' && console.warn) {
    for (const msg of validateFontLineMetrics(metrics)) {
      console.warn(`[carousel typo] ${msg}`);
    }
  }

  return metrics;
}

/**
 * @param {TextMetrics} metrics
 */
function actualInkAscent(metrics) {
  if (Number.isFinite(metrics.actualBoundingBoxAscent) && metrics.actualBoundingBoxAscent > 0) {
    return metrics.actualBoundingBoxAscent;
  }
  const fm = metrics.fontMetrics;
  if (fm && Number.isFinite(fm.xHeight) && fm.xHeight > 0) {
    return fm.xHeight;
  }
  return 0;
}

/**
 * x-height for placement: actual ink only (matches fillText).
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} fontSizePx
 */
export function measureFontXHeight(ctx, fontFor, fontSizePx, boldProbes = [false]) {
  let maxX = 0;

  for (const bold of boldProbes) {
    for (const ch of X_HEIGHT_PROBE_CHARS) {
      ctx.font = fontFor({ text: ch, bold });
      const m = ctx.measureText(ch);
      const xHeight = actualInkAscent(m);
      if (xHeight > 0) maxX = Math.max(maxX, xHeight);
    }
  }

  if (maxX <= 0) {
    maxX = Math.round(fontSizePx * 0.52);
  }
  return maxX;
}

/**
 * Cap height fallback from probe letters (actual ink only).
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} fontSizePx
 */
export function measureFontCapHeight(ctx, fontFor, fontSizePx, boldProbes = [false]) {
  let maxCap = 0;

  for (const bold of boldProbes) {
    for (const ch of CAP_HEIGHT_PROBE_CHARS) {
      ctx.font = fontFor({ text: ch, bold });
      const m = ctx.measureText(ch);
      const ascent = actualInkAscent(m);
      if (ascent > 0) maxCap = Math.max(maxCap, ascent);
    }
  }

  if (maxCap <= 0) {
    maxCap = Math.round(fontSizePx * 0.72);
  }
  return maxCap;
}

/**
 * Ascent band used to center this line in its em box (must match fillText ink).
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {InlineRun[]} tokens
 * @param {{ xHeightPx: number, capHeightPx: number, emBoxAscent: number }} metrics
 */
export function measureLineXHeightBandAscent(ctx, fontFor, tokens, metrics) {
  const hasLower = tokens.some((token) => textHasLowercase(token.text));
  if (hasLower) {
    let maxX = 0;
    for (const token of tokens) {
      ctx.font = fontFor(token);
      for (const ch of token.text) {
        if (!isXHeightProbeChar(ch)) continue;
        const m = ctx.measureText(ch);
        if (Number.isFinite(m.actualBoundingBoxAscent) && m.actualBoundingBoxAscent > 0) {
          maxX = Math.max(maxX, m.actualBoundingBoxAscent);
        }
      }
    }
    if (maxX > 0) return maxX;
    return metrics.xHeightPx;
  }

  let maxAscent = 0;
  for (const token of tokens) {
    ctx.font = fontFor(token);
    const sample = token.text.trim() === '' ? ' ' : token.text;
    const m = ctx.measureText(sample);
    if (m.actualBoundingBoxAscent > 0) {
      maxAscent = Math.max(maxAscent, m.actualBoundingBoxAscent);
      continue;
    }
    for (const ch of sample) {
      if (/\s/.test(ch)) continue;
      const cm = ctx.measureText(ch);
      if (cm.actualBoundingBoxAscent > 0) {
        maxAscent = Math.max(maxAscent, cm.actualBoundingBoxAscent);
      }
    }
  }

  if (maxAscent > 0) return maxAscent;
  return metrics.capHeightPx || metrics.emBoxAscent;
}

/** @deprecated use measureLineXHeightBandAscent */
export const measureLineBandAscent = measureLineXHeightBandAscent;

/**
 * Orange ascender/descender lines from the layout baseline (same Y as fillText).
 * @param {number} baselineY
 * @param {FontLineMetrics} metrics
 */
export function typoOrangeLinesAtBaseline(baselineY, metrics) {
  return {
    ascenderY: baselineY - metrics.ascenderLinePx,
    descenderY: baselineY + metrics.descenderLinePx,
  };
}

/**
 * Font metric baseline and ascender/descender lines inside the em box (canvas Y).
 * Orange Y values come from ASCENDER_PROBE_CHARS / DESCENDER_PROBE_CHARS only.
 * Prefer typoOrangeLinesAtBaseline for debug when text uses a centered layout baseline.
 * @param {number} emTop
 * @param {FontLineMetrics} metrics
 */
export function fontMetricLinesAtEmTop(emTop, metrics) {
  const metricBaselineY = emTop + metrics.emBoxAscent;
  return {
    metricBaselineY,
    ascenderY: metricBaselineY - metrics.ascenderLinePx,
    descenderY: metricBaselineY + metrics.descenderLinePx,
    emBottom: emTop + metrics.emBoxHeight,
  };
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {(token: InlineRun) => string} fontFor
 * @param {InlineRun[]} tokens
 * @param {number} maxWidth
 * @returns {InlineRun[][]}
 */
export function wrapInlineWords(ctx, fontFor, tokens, maxWidth) {
  /** @type {InlineRun[][]} */
  const lines = [];
  /** @type {InlineRun[]} */
  let current = [];

  for (const token of tokens) {
    const candidate = [...current, token];
    const width = measureInlineLine(ctx, fontFor, candidate).width;
    if (current.length > 0 && width > maxWidth) {
      lines.push(current);
      current = [token];
    } else {
      current = candidate;
    }
  }

  if (current.length > 0) {
    lines.push(current);
  }

  return lines;
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {InlineRun[][]} inlineLines
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} lineSlotPx Blue line slot height (orange band × lineHeights); stack step per line
 * @param {number} fontSizePx
 */
export function layoutInlineLinesFromInk(ctx, inlineLines, fontFor, lineSlotPx, fontSizePx) {
  const slotStep = Math.max(1, Math.round(lineSlotPx));
  /** @type {{ width: number, lineAscent: number, lineDescent: number, lineInkHeight: number }[]} */
  const lineMetrics = [];

  for (const tokens of inlineLines) {
    const ink = measureInlineLine(ctx, fontFor, tokens, fontSizePx);
    lineMetrics.push({
      width: ink.width,
      inkWidth: ink.width,
      inkLeft: ink.inkLeft,
      inkRight: ink.inkRight,
      lineAscent: ink.ascent,
      lineDescent: ink.descent,
      lineInkHeight: ink.inkHeight,
    });
  }

  /** @type {number[]} */
  const lineAdvances = [];
  for (let i = 0; i < lineMetrics.length; i += 1) {
    lineAdvances.push(slotStep);
  }

  const lineCount = lineMetrics.length;
  const height = lineCount <= 0 ? 0 : lineCount * slotStep;
  return {
    lineMetrics,
    lineAdvances,
    lineHeight: slotStep,
    height,
    width: Math.max(...lineMetrics.map((line) => line.width), 0),
  };
}

/**
 * Max actual ink ascent for a line (matches drawInlineLineAtTop placement).
 * @param {CanvasRenderingContext2D} ctx
 * @param {InlineRun[]} tokens
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} fontSizePx
 */
export function measureLineInkAscent(ctx, tokens, fontFor, fontSizePx) {
  let maxAscent = 0;
  for (const token of tokens) {
    maxAscent = Math.max(maxAscent, inkAscentForToken(ctx, fontFor, token, fontSizePx));
  }
  return maxAscent;
}

/**
 * Draw on an alphabetic baseline (mixed weights on one line stay aligned).
 * @param {CanvasRenderingContext2D} ctx
 * @param {InlineRun[]} tokens
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} x
 * @param {number} baselineY
 * @param {number} [fontSizePx]
 */
/**
 * @param {(token: InlineRun) => string} [colorFor] When set, called before each token draw.
 */
export function drawInlineLine(ctx, tokens, fontFor, x, baselineY, fontSizePx = 80, colorFor = null) {
  let cursorX = x;
  for (const token of tokens) {
    if (colorFor) {
      ctx.fillStyle = colorFor(token);
    }
    ctx.font = fontFor(token);
    ctx.textBaseline = 'alphabetic';
    ctx.fillText(token.text, cursorX, baselineY);
    const ink = measureTokenInkActual(ctx, fontFor, token, fontSizePx);
    cursorX += ink.width;
  }
}

/**
 * Draw a line inside a text column. Right alignment uses canvas `textAlign: 'right'`
 * so the painted edge matches the column (avoids under-measured display-serif ink).
 * @param {CanvasRenderingContext2D} ctx
 * @param {InlineRun[]} tokens
 * @param {(token: InlineRun) => string} fontFor
 * @param {{ columnLeft: number, contentWidth: number, columnRight?: number }} column
 * @param {'left' | 'center' | 'right'} alignment
 * @param {number} baselineY
 * @param {number} [fontSizePx]
 */
/**
 * Draw column shrunk on the aligned edge when margin is 0 (bitmap safe inset).
 * @param {{ columnLeft: number, contentWidth: number, columnRight?: number, canvasSize?: number, marginHorizontal?: number, bitmapSafe?: number }} column
 * @param {number} fontSizePx
 * @param {'left' | 'center' | 'right'} alignment
 */
/** @param {{ columnLeft: number, contentWidth: number, columnRight?: number, canvasSize?: number, marginHorizontal?: number }} column @param {number} fontSizePx @param {'left'|'center'|'right'} alignment */
export function columnDrawBounds(column, fontSizePx, alignment) {
  return drawColumnForLine(column, fontSizePx, alignment);
}

function drawColumnForLine(column, fontSizePx, alignment) {
  const canvasSize = column.canvasSize
    ?? (column.columnRight != null ? column.columnRight : column.contentWidth);
  const marginH = column.marginHorizontal ?? 0;
  const rawLeft = column.columnLeft;
  const rawRight = column.columnRight ?? column.columnLeft + column.contentWidth;
  const bleed = columnEdgeBleedPx(fontSizePx, canvasSize, marginH);

  let columnLeft = rawLeft;
  let columnRight = rawRight;
  if (bleed > 0) {
    if (alignment === 'left') {
      columnLeft = rawLeft + bleed;
    } else if (alignment === 'right') {
      columnRight = rawRight - bleed;
    } else {
      columnLeft = rawLeft + bleed;
      columnRight = rawRight - bleed;
    }
  }

  return {
    columnLeft,
    columnRight,
    contentWidth: Math.max(1, columnRight - columnLeft),
    canvasSize,
    marginHorizontal: marginH,
  };
}

/**
 * Keep measured ink inside the canvas bitmap.
 * @param {number} x
 * @param {{ inkLeft: number, inkRight: number }} extents
 * @param {number} canvasSize
 */
function clampLineX(x, extents, canvasSize) {
  const pad = Math.max(2, Math.round(canvasSize * 0.002));
  let nx = x;
  if (nx + extents.inkRight > canvasSize - pad) {
    nx = canvasSize - pad - extents.inkRight;
  }
  if (nx + extents.inkLeft < pad) {
    nx = pad - extents.inkLeft;
  }
  return nx;
}

/**
 * @typedef {{ inkLeft: number, inkRight: number, inkWidth: number }} LineInkExtents
 */

/**
 * Place a measured ink box inside the text column (L/C/R uses the box width, not the column).
 * @param {LineInkExtents} extents
 * @param {{ columnLeft: number, contentWidth: number, columnRight?: number, canvasSize?: number, marginHorizontal?: number }} column
 * @param {'left' | 'center' | 'right'} alignment
 * @param {number} fontSizePx
 * @returns {{ x: number, width: number, inkLeft: number, inkRight: number }}
 */
export function lineInkBoundsFromExtents(extents, column, alignment, fontSizePx) {
  if (!extents || extents.inkWidth <= 0) {
    return { x: column.columnLeft, width: 0, inkLeft: 0, inkRight: 0 };
  }

  const drawColumn = drawColumnForLine(column, fontSizePx, alignment);
  const canvasSize = drawColumn.canvasSize;
  let x;
  if (alignment === 'right') {
    x = drawColumn.columnRight - extents.inkRight;
  } else if (alignment === 'center') {
    const boxLeft = drawColumn.columnLeft + (drawColumn.contentWidth - extents.inkWidth) / 2;
    x = boxLeft - extents.inkLeft;
  } else {
    x = drawColumn.columnLeft;
  }

  x = clampLineX(x, extents, canvasSize);
  return {
    x,
    width: extents.inkWidth,
    inkLeft: extents.inkLeft,
    inkRight: extents.inkRight,
  };
}

/**
 * Horizontal ink bounds for one line inside the text column (matches {@link drawInlineLineInColumn}).
 * @param {CanvasRenderingContext2D} ctx
 * @param {InlineRun[]} tokens
 * @param {(token: InlineRun) => string} fontFor
 * @param {{ columnLeft: number, contentWidth: number, columnRight?: number, canvasSize?: number, marginHorizontal?: number }} column
 * @param {'left' | 'center' | 'right'} alignment
 * @param {number} [fontSizePx]
 * @param {LineInkExtents|null} [cachedExtents] Layout-time ink box (from lineMetrics)
 * @returns {{ x: number, width: number, inkLeft: number, inkRight: number }}
 */
export function lineInkBoundsInColumn(
  ctx,
  tokens,
  fontFor,
  column,
  alignment,
  fontSizePx = 80,
  cachedExtents = null,
) {
  const extents = cachedExtents ?? measureInlineLineInkExtents(ctx, fontFor, tokens, fontSizePx);
  if (tokens.length === 0 && !cachedExtents) {
    return { x: column.columnLeft, width: 0, inkLeft: 0, inkRight: 0 };
  }
  return lineInkBoundsFromExtents(extents, column, alignment, fontSizePx);
}

export function drawInlineLineInColumn(
  ctx,
  tokens,
  fontFor,
  column,
  alignment,
  baselineY,
  fontSizePx = 80,
  cachedExtents = null,
  colorFor = null,
) {
  if (tokens.length === 0) return;

  const { x } = lineInkBoundsInColumn(
    ctx,
    tokens,
    fontFor,
    column,
    alignment,
    fontSizePx,
    cachedExtents,
  );

  ctx.textAlign = 'left';
  ctx.textBaseline = 'alphabetic';
  drawInlineLine(ctx, tokens, fontFor, x, baselineY, fontSizePx, colorFor);
}

/**
 * Top of painted ink sits on lineTop; mixed weights share one alphabetic baseline.
 * @param {CanvasRenderingContext2D} ctx
 * @param {InlineRun[]} tokens
 * @param {(token: InlineRun) => string} fontFor
 * @param {number} x
 * @param {number} lineTop
 * @param {number} fontSizePx
 */
export function drawInlineLineAtTop(ctx, tokens, fontFor, x, lineTop, fontSizePx, colorFor = null) {
  const baselineY = lineTop + measureLineInkAscent(ctx, tokens, fontFor, fontSizePx);
  drawInlineLine(ctx, tokens, fontFor, x, baselineY, fontSizePx, colorFor);
}
