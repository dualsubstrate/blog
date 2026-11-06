/**
 * Typographic probe alphabets for carousel line-box layout and debug.
 *
 * RULES (do not break these when editing layout or debug code):
 * 1. X_HEIGHT_PROBE_CHARS → red meanline / layout centering only. Never bdfhklt or gjpqy.
 * 2. ASCENDER_PROBE_CHARS → orange top line only. Never used for x-height or em box.
 * 3. DESCENDER_PROBE_CHARS → orange bottom line only. Never used for x-height or em box.
 * 4. EM_STRUT_* → blue em box only. Uses fontBoundingBox (full em), not glyph ink.
 * 5. Orange lines use actualBoundingBox ink on ascender/descender probes only.
 * 6. At `lineHeights` 1, blue line slot hugs orange (ascender to descender); mult > 1 pads blue symmetrically.
 * 7. Orange band height = ascenderLinePx + descenderLinePx; blue slot = band × lineHeight mult; stack advances by slot height.
 *
 * @typedef {'xHeight'|'ascender'|'descender'|'emStrut'|'capHeight'} TypoProbeRole
 */

/** Lowercase at x-height only (no ascenders, no descenders). Red meanline band. */
export const X_HEIGHT_PROBE_CHARS = 'acemnorsuvxz';

/** Lowercase stems above x-height. Orange ascender line (baseline to top). */
export const ASCENDER_PROBE_CHARS = 'bdfhklt';

/** Lowercase tails below baseline. Orange descender line (baseline to bottom). */
export const DESCENDER_PROBE_CHARS = 'gjpqy';

/** Cap height probes (all-caps lines). */
export const CAP_HEIGHT_PROBE_CHARS = 'HEP';

/** Full em strut: caps, ascenders, descenders (blue em box). */
export const EM_STRUT_CHARS = 'MHÉÓÁgjpqy';
export const EM_STRUT_STRINGS = ['M', 'Hgjpqy', 'Éy'];

/** Minimum px between blue em edge and orange ascender/descender lines. */
export const EM_ORANGE_INSET_PX = 4;

/** @type {Record<TypoProbeRole, string>} */
export const TYPO_PROBE_BY_ROLE = {
  xHeight: X_HEIGHT_PROBE_CHARS,
  ascender: ASCENDER_PROBE_CHARS,
  descender: DESCENDER_PROBE_CHARS,
  capHeight: CAP_HEIGHT_PROBE_CHARS,
  emStrut: EM_STRUT_CHARS,
};

/**
 * Fail fast if probe alphabets overlap (prevents mixing roles in one string).
 * @returns {string[]} overlap errors; empty when valid
 */
export function findProbeSetOverlaps() {
  /** @type {string[]} */
  const errors = [];
  const pairs = [
    ['xHeight', X_HEIGHT_PROBE_CHARS, 'ascender', ASCENDER_PROBE_CHARS],
    ['xHeight', X_HEIGHT_PROBE_CHARS, 'descender', DESCENDER_PROBE_CHARS],
    ['ascender', ASCENDER_PROBE_CHARS, 'descender', DESCENDER_PROBE_CHARS],
  ];

  for (const [aName, aChars, bName, bChars] of pairs) {
    for (const ch of aChars) {
      if (bChars.includes(ch)) {
        errors.push(`probe overlap: "${ch}" is in both ${aName} and ${bName}`);
      }
    }
  }
  return errors;
}

/** @throws {Error} when probe sets overlap */
export function assertProbeSetsDisjoint() {
  const errors = findProbeSetOverlaps();
  if (errors.length > 0) {
    throw new Error(`typo-probes: ${errors.join('; ')}`);
  }
}

/** @param {string} ch */
export function isXHeightProbeChar(ch) {
  return X_HEIGHT_PROBE_CHARS.includes(ch);
}

/** @param {string} ch */
export function isAscenderProbeChar(ch) {
  return ASCENDER_PROBE_CHARS.includes(ch);
}

/** @param {string} ch */
export function isDescenderProbeChar(ch) {
  return DESCENDER_PROBE_CHARS.includes(ch);
}

/** @param {string} text */
export function textHasLowercase(text) {
  return /[a-z]/.test(text);
}

/** @param {string} text */
export function textHasXHeightProbeChar(text) {
  for (const ch of text) {
    if (isXHeightProbeChar(ch)) return true;
  }
  return false;
}

/**
 * @typedef {Object} FontLineMetricsShape
 * @property {number} emBoxAscent
 * @property {number} emBoxDescent
 * @property {number} emBoxHeight
 * @property {number} ascenderLinePx
 * @property {number} descenderLinePx
 * @property {number} xHeightPx
 */

/**
 * Orange band: ascender probe line to descender probe line (layout baseline in the middle).
 * @param {FontLineMetricsShape} metrics
 */
export function orangeBandHeightPx(metrics) {
  return metrics.ascenderLinePx + metrics.descenderLinePx;
}

/**
 * @param {FontLineMetricsShape} metrics
 * @param {number} [lineHeightMult]
 */
export function resolveLineHeightMultForMetrics(metrics, lineHeightMult) {
  if (lineHeightMult != null && Number.isFinite(lineHeightMult)) {
    return lineHeightMult;
  }
  const fromMetrics = /** @type {{ lineHeightMult?: number }} */ (metrics).lineHeightMult;
  if (fromMetrics != null && Number.isFinite(fromMetrics)) {
    return fromMetrics;
  }
  return 1;
}

/**
 * Blue line slot height (stack advance per line): orange band × lineHeight mult (never below band).
 * @param {FontLineMetricsShape} metrics
 * @param {number} [lineHeightMult]
 */
export function lineSlotHeightPx(metrics, lineHeightMult) {
  const mult = resolveLineHeightMultForMetrics(metrics, lineHeightMult);
  const orange = orangeBandHeightPx(metrics);
  return Math.max(orange, Math.round(orange * Math.max(mult, 0.01)));
}

/**
 * Equal inset from blue top to orange ascender and orange descender to blue bottom.
 * @param {FontLineMetricsShape} metrics
 * @param {number} [lineHeightMult]
 */
export function symmetricOrangeInsetPx(metrics, lineHeightMult) {
  const orange = orangeBandHeightPx(metrics);
  const slot = lineSlotHeightPx(metrics, lineHeightMult);
  return Math.max(0, (slot - orange) / 2);
}

/**
 * Orange ascender/descender with equal gap to blue line slot top/bottom.
 * @param {number} emTop
 * @param {FontLineMetricsShape} metrics
 * @param {number} [lineHeightMult]
 */
export function typoOrangeLinesInEmBox(emTop, metrics, lineHeightMult) {
  const mult = resolveLineHeightMultForMetrics(metrics, lineHeightMult);
  const gap = symmetricOrangeInsetPx(metrics, mult);
  const orange = orangeBandHeightPx(metrics);
  const slot = metrics.lineSlotPx != null && Number.isFinite(metrics.lineSlotPx)
    ? metrics.lineSlotPx
    : lineSlotHeightPx(metrics, mult);
  return {
    gap,
    ascenderY: emTop + gap,
    descenderY: emTop + gap + orange,
    emBottom: emTop + slot,
  };
}

/**
 * Layout baseline when orange lines are symmetric inside the blue line slot.
 * @param {number} emTop
 * @param {FontLineMetricsShape} metrics
 * @param {number} [lineHeightMult]
 */
export function baselineYFromSymmetricOrange(emTop, metrics, lineHeightMult) {
  const gap = symmetricOrangeInsetPx(metrics, lineHeightMult);
  return emTop + gap + metrics.ascenderLinePx;
}

/**
 * Sanity-check measured lines so debug colors cannot collapse onto each other.
 * @param {FontLineMetricsShape} metrics
 * @returns {string[]} warnings; empty when ok
 */
export function validateFontLineMetrics(metrics) {
  /** @type {string[]} */
  const warnings = [];
  const {
    emBoxAscent,
    emBoxDescent,
    emBoxHeight,
    ascenderLinePx,
    descenderLinePx,
    xHeightPx,
  } = metrics;

  if (ascenderLinePx >= emBoxAscent) {
    warnings.push('ascenderLinePx >= font em ascent (probe clamp may be wrong)');
  }
  if (descenderLinePx >= emBoxDescent) {
    warnings.push('descenderLinePx >= font em descent (probe clamp may be wrong)');
  }
  if (orangeBandHeightPx({ ascenderLinePx, descenderLinePx, emBoxAscent, emBoxDescent, emBoxHeight, xHeightPx }) <= 0) {
    warnings.push('orange band height is zero');
  }
  if (xHeightPx >= ascenderLinePx) {
    warnings.push('xHeightPx >= ascenderLinePx (red band would reach ascender line)');
  }
  if (emBoxAscent + emBoxDescent !== emBoxHeight && emBoxHeight > 0) {
    warnings.push('emBoxHeight != emBoxAscent + emBoxDescent');
  }
  return warnings;
}

assertProbeSetsDisjoint();
