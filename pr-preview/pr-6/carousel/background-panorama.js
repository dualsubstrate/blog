import { mixHex, rotateHexHue, shiftHex } from './background.js';

/**
 * @typedef {'none' | 'drift' | 'mesh-corners'} BackgroundWaveStyle
 */

/** Studio + parser bounds for `backgroundWave` numeric fields. */
export const BACKGROUND_WAVE_LIMITS = {
  intensity: { min: 0, max: 0.72, step: 0.01, default: 0.32 },
  color: { min: 0, max: 1, step: 0.01, default: 0.55 },
  variety: { min: 0, max: 1, step: 0.01, default: 0.62 },
  blur: { min: 0.15, max: 1.2, step: 0.01, default: 0.55 },
  radius: { min: 0.35, max: 1.85, step: 0.01, default: 1 },
  lobes: { min: 2, max: 64, step: 1 },
  phase: { min: -6.29, max: 6.29, step: 0.05, default: 0 },
  hueShift: { min: -180, max: 180, step: 1, default: 0 },
};

/** @type {{ id: BackgroundWaveStyle, label: string }[]} */
export const BACKGROUND_WAVE_STYLES = [
  { id: 'none', label: 'None (solid)' },
  { id: 'drift', label: 'Drift (horizontal wash)' },
  { id: 'mesh-corners', label: 'Mesh corners' },
];

/**
 * @typedef {Object} BackgroundWaveConfig
 * @property {BackgroundWaveStyle} style `none` (flat), `drift` (horizontal), or `mesh-corners` (corner pools)
 * @property {number|null} lobes Extra mid-strip pools (`drift`: count across deck; `mesh-corners`: bottom travelers)
 * @property {number} intensity 0–0.72 wash presence (opacity / coverage)
 * @property {number} color 0–1 accent richness in each lobe (mix toward palette accents)
 * @property {number} variety 0–1 extended background palette steps (blends between base, accents, muted)
 * @property {number} blur 0.2–1.2 soft spread factor in base lobe sizing
 * @property {number} radius 0.35–1.85 multiplier on lobe size (default 1)
 * @property {number} phase Radians offset for lobe layout
 * @property {number} hueShift Degrees to rotate background wash hues (-180 to 180)
 */

/**
 * @typedef {Object} BackgroundPanoramaContext
 * @property {number} slideIndex 0-based index in `deck.slides`
 * @property {number} slideCount
 */

/**
 * @typedef {Object} PanoramaPalette
 * @property {string} background
 * @property {string} accent1
 * @property {string} accent2
 * @property {string} muted
 */

/** @param {string} hex @param {number} alpha */
function hexToRgba(hex, alpha) {
  let raw = hex.replace('#', '').trim();
  if (raw.length === 3) {
    raw = raw[0] + raw[0] + raw[1] + raw[1] + raw[2] + raw[2];
  }
  const r = parseInt(raw.slice(0, 2), 16) || 0;
  const g = parseInt(raw.slice(2, 4), 16) || 0;
  const b = parseInt(raw.slice(4, 6), 16) || 0;
  return `rgba(${r}, ${g}, ${b}, ${Math.max(0, Math.min(1, alpha))})`;
}

/**
 * Map studio intensity slider (0–0.72) to wash presence (0–1): opacity and coverage only.
 * @param {number} rawIntensity
 * @returns {number}
 */
export function waveVisualStrength(rawIntensity) {
  const { min, max } = BACKGROUND_WAVE_LIMITS.intensity;
  const t = clamp((rawIntensity - min) / (max - min), 0, 1);
  return Math.pow(t, 0.78) * 0.98;
}

/**
 * Map studio color slider (0–1) to accent mix strength (0–1): how saturated each lobe reads.
 * @param {number} [rawColor]
 * @returns {number}
 */
export function waveColorStrength(rawColor) {
  const { min, max, default: fallback } = BACKGROUND_WAVE_LIMITS.color;
  const raw = Number.isFinite(rawColor) ? rawColor : fallback;
  const t = clamp((raw - min) / (max - min), 0, 1);
  return Math.pow(t, 0.82);
}

/**
 * Map studio variety slider (0–1): how many blended steps from the deck palette feed the wash.
 * @param {number} [rawVariety]
 * @returns {number}
 */
export function waveVarietyStrength(rawVariety) {
  const { min, max, default: fallback } = BACKGROUND_WAVE_LIMITS.variety;
  const raw = Number.isFinite(rawVariety) ? rawVariety : fallback;
  const t = clamp((raw - min) / (max - min), 0, 1);
  return Math.pow(t, 0.88);
}

/**
 * Full smooth ring of wash colors derived only from deck base + accent1 + accent2 + muted.
 * @param {PanoramaPalette} palette
 * @returns {string[]}
 */
export function buildExtendedBackgroundRing(palette) {
  if (!palette) return [];
  const base = palette.background || '#ffffff';
  const a1 = palette.accent1 || '#808080';
  const a2 = palette.accent2 || a1;
  const muted = palette.muted || '#808080';

  /** @param {string} from @param {string} to @param {number[]} mixSteps */
  const segment = (from, to, mixSteps) => {
    const out = [];
    for (const t of mixSteps) out.push(mixHex(from, to, t));
    out.push(to);
    return out;
  };

  return [
    ...segment(base, a1, [0.2, 0.4, 0.62, 0.82]),
    ...segment(a1, a2, [0.22, 0.45, 0.68]),
    ...segment(a2, muted, [0.22, 0.45, 0.68]),
    ...segment(muted, base, [0.2, 0.4, 0.62]),
  ];
}

/**
 * Subset of {@link buildExtendedBackgroundRing} used by drift/mesh washes.
 * Low variety keeps the three deck accents; high variety adds in-between blends for smoother fields.
 * @param {PanoramaPalette} palette
 * @param {number} [rawVariety]
 * @returns {string[]}
 */
export function extendedBackgroundPalette(palette, rawVariety) {
  const full = buildExtendedBackgroundRing(palette);
  const v = waveVarietyStrength(rawVariety);
  const minSteps = 3;
  if (v >= 0.995) return full;
  const count = Math.max(minSteps, Math.round(minSteps + v * (full.length - minSteps)));
  if (count >= full.length) return full;
  /** @type {string[]} */
  const picked = [];
  for (let i = 0; i < count; i += 1) {
    const idx = Math.round((i / Math.max(1, count - 1)) * (full.length - 1));
    picked.push(full[idx]);
  }
  return picked;
}

/**
 * @param {string[]} extendedColors
 * @param {number} index
 * @param {number} t
 */
function pickDriftAccent(extendedColors, index, t) {
  if (!extendedColors.length) return '#808080';
  const rawIndex = index * 2 + Math.floor(t * 5);
  return extendedColors[rawIndex % extendedColors.length];
}

/**
 * @param {string[]} extendedColors
 * @param {number} fraction 0–1 position around the extended ring
 */
function pickExtendedAt(extendedColors, fraction) {
  if (!extendedColors.length) return '#808080';
  const idx = Math.round(clamp(fraction, 0, 1) * (extendedColors.length - 1));
  return extendedColors[idx];
}

/** Relative luminance (sRGB); matches studio palette swatch threshold. */
function isLightBackground(hex) {
  if (typeof hex !== 'string') return false;
  let raw = hex.replace('#', '').trim();
  if (raw.length === 3) {
    raw = raw[0] + raw[0] + raw[1] + raw[1] + raw[2] + raw[2];
  }
  if (raw.length !== 6) return false;
  const r = (parseInt(raw.slice(0, 2), 16) || 0) / 255;
  const g = (parseInt(raw.slice(2, 4), 16) || 0) / 255;
  const b = (parseInt(raw.slice(4, 6), 16) || 0) / 255;
  const luminance = 0.2126 * r + 0.7152 * g + 0.0722 * b;
  return luminance > 0.72;
}

/** @param {number} value @param {number} min @param {number} max */
function clamp(value, min, max) {
  return Math.max(min, Math.min(max, value));
}

/** Carousel backgrounds are always panoramic across slides. */
export function parseBackgroundGradientMode() {
  return 'panoramic-wave';
}

/** @param {unknown} value @returns {BackgroundWaveStyle} */
function parseBackgroundWaveStyle(value) {
  if (typeof value !== 'string') return 'drift';
  const norm = value.trim().toLowerCase().replace(/_/g, '-');
  if (norm === 'none' || norm === 'solid' || norm === 'off') return 'none';
  if (norm === 'mesh-corners' || norm === 'mesh' || norm === 'corners') return 'mesh-corners';
  return 'drift';
}

/**
 * @param {BackgroundWaveConfig|undefined|null} wave
 * @returns {boolean}
 */
export function shouldPaintPanoramicWave(wave) {
  return (wave?.style ?? 'drift') !== 'none';
}

/**
 * @param {unknown} spec
 * @returns {BackgroundWaveConfig}
 */

export function parseBackgroundWaveConfig(spec) {
  /** @type {BackgroundWaveConfig} */
  const defaults = {
    style: 'drift',
    lobes: null,
    intensity: BACKGROUND_WAVE_LIMITS.intensity.default,
    color: BACKGROUND_WAVE_LIMITS.color.default,
    variety: BACKGROUND_WAVE_LIMITS.variety.default,
    blur: BACKGROUND_WAVE_LIMITS.blur.default,
    radius: BACKGROUND_WAVE_LIMITS.radius.default,
    phase: BACKGROUND_WAVE_LIMITS.phase.default,
    hueShift: BACKGROUND_WAVE_LIMITS.hueShift.default,
  };
  if (!spec || typeof spec !== 'object') return defaults;
  const raw = /** @type {Record<string, unknown>} */ (spec).backgroundWave;
  if (!raw || typeof raw !== 'object') return defaults;
  const src = /** @type {Record<string, unknown>} */ (raw);
  return {
    style: parseBackgroundWaveStyle(src.style),
    lobes: Number.isFinite(Number(src.lobes))
      ? clamp(Math.round(Number(src.lobes)), BACKGROUND_WAVE_LIMITS.lobes.min, BACKGROUND_WAVE_LIMITS.lobes.max)
      : null,
    intensity: clamp(
      Number(src.intensity ?? defaults.intensity),
      BACKGROUND_WAVE_LIMITS.intensity.min,
      BACKGROUND_WAVE_LIMITS.intensity.max,
    ),
    color: clamp(
      Number(src.color ?? defaults.color),
      BACKGROUND_WAVE_LIMITS.color.min,
      BACKGROUND_WAVE_LIMITS.color.max,
    ),
    variety: clamp(
      Number(src.variety ?? src.colorVariety ?? defaults.variety),
      BACKGROUND_WAVE_LIMITS.variety.min,
      BACKGROUND_WAVE_LIMITS.variety.max,
    ),
    blur: clamp(
      Number(src.blur ?? defaults.blur),
      BACKGROUND_WAVE_LIMITS.blur.min,
      BACKGROUND_WAVE_LIMITS.blur.max,
    ),
    radius: clamp(
      Number(src.radius ?? defaults.radius),
      BACKGROUND_WAVE_LIMITS.radius.min,
      BACKGROUND_WAVE_LIMITS.radius.max,
    ),
    phase: Number.isFinite(Number(src.phase)) ? Number(src.phase) : defaults.phase,
    hueShift: clamp(
      Number(src.hueShift ?? defaults.hueShift),
      BACKGROUND_WAVE_LIMITS.hueShift.min,
      BACKGROUND_WAVE_LIMITS.hueShift.max,
    ),
  };
}

/**
 * Apply studio hue shift to palette colors used for background painting only.
 * @param {PanoramaPalette} palette
 * @param {BackgroundWaveConfig|undefined|null} wave
 * @returns {PanoramaPalette}
 */
export function paletteWithBackgroundHueShift(palette, wave) {
  const shift = wave?.hueShift ?? BACKGROUND_WAVE_LIMITS.hueShift.default;
  if (!Number.isFinite(shift) || Math.abs(shift) < 0.001) return palette;
  return {
    background: rotateHexHue(palette.background, shift),
    accent1: rotateHexHue(palette.accent1, shift),
    accent2: rotateHexHue(palette.accent2, shift),
    muted: rotateHexHue(palette.muted, shift),
  };
}

/**
 * @param {{ slides?: { number: number }[], deck?: Record<string, unknown> }|null|undefined} deck
 * @param {number} slideNumber
 * @returns {BackgroundPanoramaContext|null}
 */
export function buildBackgroundPanoramaContext(deck, slideNumber) {
  if (!deck?.slides?.length) return null;
  const index = deck.slides.findIndex((slide) => slide.number === slideNumber);
  if (index < 0) return null;
  return { slideIndex: index, slideCount: deck.slides.length };
}

/** @type {Map<string, HTMLCanvasElement>} */
const panoramaCache = new Map();
const PANORAMA_CACHE_MAX = 16;
/** Render mesh/drift washes at 2× then downsample to reduce gradient banding. */
const PANORAMA_SUPERSAMPLE = 2;

/**
 * @param {string} key
 * @param {HTMLCanvasElement} canvas
 */
function storePanoramaCache(key, canvas) {
  if (panoramaCache.has(key)) {
    panoramaCache.delete(key);
  }
  panoramaCache.set(key, canvas);
  while (panoramaCache.size > PANORAMA_CACHE_MAX) {
    const oldest = panoramaCache.keys().next().value;
    if (oldest) panoramaCache.delete(oldest);
  }
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} totalWidth
 * @param {number} height
 * @param {number} cx
 * @param {number} cy
 * @param {number} radius
 * @param {PanoramaPalette} palette
 * @param {string} accent
 * @param {BackgroundWaveConfig} wave
 * @param {{ presenceScale?: number, gain?: number }} [modifiers]
 */
function paintColorLobe(ctx, totalWidth, height, cx, cy, radius, palette, accent, wave, modifiers = {}) {
  const presenceScale = modifiers.presenceScale ?? 1;
  const gain = modifiers.gain ?? 1;
  const presence = waveVisualStrength(wave.intensity) * presenceScale;
  const chroma = waveColorStrength(wave.color);
  const variety = waveVarietyStrength(wave.variety);
  const light = isLightBackground(palette.background);
  let core;
  let mid;
  let compositeOp;
  let alphaCore;
  let alphaMid;
  let alphaTail;

  if (light) {
    // Screen blend washes out on light bases; direct mix matches studio wave thumbnails.
    const mixT = 0.06 + chroma * 0.76;
    core = mixHex(palette.background, accent, Math.min(mixT + 0.12, 0.88));
    mid = mixHex(palette.background, accent, mixT * 0.82);
    compositeOp = 'source-over';
    alphaCore = (0.12 + presence * 1.12) * gain;
    alphaMid = (0.07 + presence * 0.62) * gain;
    alphaTail = (0.01 + presence * 0.09) * gain;
  } else {
    const mixT = 0.1 + chroma * 0.74;
    core = mixHex(palette.background, accent, Math.min(mixT + 0.16, 0.9));
    mid = mixHex(palette.background, accent, mixT * 0.78);
    compositeOp = 'screen';
    alphaCore = (0.08 + presence * 0.82) * gain;
    alphaMid = (0.04 + presence * 0.44) * gain;
    alphaTail = (0.01 + presence * 0.08) * gain;
  }

  const gradient = ctx.createRadialGradient(cx, cy, 0, cx, cy, radius);
  const midStop = 0.26 + (1 - variety) * 0.14;
  const tailStop = 0.68 + (1 - variety) * 0.18;
  gradient.addColorStop(0, hexToRgba(core, alphaCore));
  gradient.addColorStop(midStop, hexToRgba(mid, alphaMid));
  gradient.addColorStop(tailStop, hexToRgba(mid, alphaTail));
  gradient.addColorStop(1, hexToRgba(palette.background, 0));
  ctx.save();
  ctx.globalCompositeOperation = compositeOp;
  ctx.fillStyle = gradient;
  ctx.beginPath();
  ctx.arc(cx, cy, radius, 0, Math.PI * 2);
  ctx.fill();
  ctx.restore();
}

const PHI_FRAC = 0.6180339887498949;

/**
 * Mixed incommensurate waves along panoramic t (0..1). Avoids a single repeating period.
 * @param {number} t
 * @param {number} phase
 */
function driftMixX(t, phase) {
  return (
    Math.sin(t * Math.PI * 2.07 + phase) * 0.4
    + Math.sin(t * Math.PI * 3.29 + phase * 0.41) * 0.27
    + Math.cos(t * Math.PI * 1.31 - phase * 0.53) * 0.17
    + Math.sin(t * Math.PI * 5.17 + phase * 0.17) * 0.11
    + Math.cos(t * Math.PI * 0.67 + phase * 0.89) * 0.08
  );
}

/**
 * @param {number} t
 * @param {number} phase
 * @param {number} seed
 */
function driftMixY(t, phase, seed) {
  return (
    0.5
    + Math.sin(t * Math.PI * 1.59 + phase * 0.62 + seed * 0.37) * 0.17
    + Math.cos(t * Math.PI * 2.73 - phase * 0.28 + seed * 0.19) * 0.13
    + Math.sin(t * Math.PI * 0.83 + seed * 0.61) * 0.09
    + Math.cos(t * Math.PI * 4.03 + phase * 0.11) * 0.06
    + Math.sin(t * Math.PI * 6.21 - seed * 0.23 + phase * 0.05) * 0.04
  );
}

/**
 * @param {number} t
 * @param {number} phase
 * @param {number} seed
 */
function driftRadiusScale(t, phase, seed) {
  return (
    0.76
    + 0.2 * Math.sin(t * Math.PI * 2.41 + seed * 0.53 + phase)
    + 0.11 * Math.cos(t * Math.PI * 1.19 - phase * 0.37 + seed * 0.29)
    + 0.07 * Math.sin(t * Math.PI * 3.83 + seed * 0.17)
  );
}

/**
 * Golden-scramble position so lobes are not evenly spaced.
 * @param {number} index
 * @param {number} phase
 * @param {number} salt
 */
function driftScrambleT(index, phase, salt = 0) {
  const raw = (index * PHI_FRAC + 0.23 + salt * 0.137 + phase * 0.03) % 1;
  const wobble = driftMixX(raw, phase + salt) * 0.06;
  return Math.max(0.02, Math.min(0.98, raw + wobble));
}

/** @param {BackgroundWaveConfig} wave */
function waveRadiusScale(wave) {
  const raw = wave.radius ?? BACKGROUND_WAVE_LIMITS.radius.default;
  return clamp(raw, BACKGROUND_WAVE_LIMITS.radius.min, BACKGROUND_WAVE_LIMITS.radius.max);
}

/**
 * @param {number} slideWidth
 * @param {number} height
 * @param {PanoramaPalette} palette
 * @param {number} slideCount
 * @param {BackgroundWaveConfig} wave
 */
function buildDriftWaveCanvas(slideWidth, height, palette, slideCount, wave) {
  const totalWidth = slideWidth * slideCount;
  const lobeCount = wave.lobes ?? Math.max(4, slideCount + 2);
  const blurScale = wave.blur;
  const radiusScale = waveRadiusScale(wave);
  const phase = wave.phase;
  const extendedColors = extendedBackgroundPalette(palette, wave.variety);

  const canvas = document.createElement('canvas');
  canvas.width = totalWidth;
  canvas.height = height;
  const ctx = canvas.getContext('2d');
  if (!ctx) throw new Error('Canvas 2D context unavailable');

  ctx.fillStyle = palette.background;
  ctx.fillRect(0, 0, totalWidth, height);

  const span = Math.max(slideWidth, height);
  const baseRadius = span * (0.4 + blurScale * 0.52) * radiusScale;

  for (let i = 0; i < lobeCount; i += 1) {
    const t = driftScrambleT(i, phase, 0);
    const cx = t * totalWidth + driftMixX(t, phase) * slideWidth * 0.17;
    const cy = height * driftMixY(t, phase, i + 1);
    const radius = baseRadius * driftRadiusScale(t, phase, i);
    const accent = pickDriftAccent(extendedColors, i, t);
    paintColorLobe(ctx, totalWidth, height, cx, cy, radius, palette, accent, wave);
  }

  const shimmerCount = Math.round(lobeCount * 1.5) + 1;
  for (let i = 0; i < shimmerCount; i += 1) {
    const t = driftScrambleT(i, phase + 1.4, 7);
    const cx = t * totalWidth + driftMixX(t, phase + 2.1) * slideWidth * 0.11;
    const cy = height * driftMixY(t, phase + 0.9, i + 20);
    const radius = baseRadius * 0.48 * driftRadiusScale(t, phase + 0.5, i + 3);
    const accent = pickDriftAccent(extendedColors, i + 11, t);
    paintColorLobe(ctx, totalWidth, height, cx, cy, radius, palette, accent, wave, {
      presenceScale: 0.72,
      gain: 0.55,
    });
  }

  return canvas;
}

/** @param {BackgroundWaveConfig} wave @param {number} slideCount */
function meshTravelerCount(wave, slideCount) {
  const fallback = Math.max(2, Math.min(4, slideCount));
  if (wave.lobes == null) return fallback;
  // Keep mesh pools sparse and off the slide grid; high `lobes` values are for drift only.
  return clamp(Math.round(wave.lobes), 2, Math.min(8, slideCount + 2));
}

/**
 * Corner-anchored mesh: warm bottom pools, cool top accents, dark center (panoramic strip).
 * @param {number} slideWidth
 * @param {number} height
 * @param {PanoramaPalette} palette
 * @param {number} slideCount
 * @param {BackgroundWaveConfig} wave
 */
function buildMeshCornersCanvas(slideWidth, height, palette, slideCount, wave) {
  const totalWidth = slideWidth * slideCount;
  const blurScale = wave.blur;
  const radiusScale = waveRadiusScale(wave);
  const phase = wave.phase;
  const travelerCount = meshTravelerCount(wave, slideCount);
  const span = Math.max(slideWidth, height);
  const cornerRadius = span * (0.78 + blurScale * 0.42) * radiusScale;
  const topRadius = span * (0.46 + blurScale * 0.28) * radiusScale;
  const travelerRadius = span * (0.44 + blurScale * 0.32) * radiusScale;

  const canvas = document.createElement('canvas');
  canvas.width = totalWidth;
  canvas.height = height;
  const ctx = canvas.getContext('2d');
  if (!ctx) throw new Error('Canvas 2D context unavailable');

  ctx.fillStyle = palette.background;
  ctx.fillRect(0, 0, totalWidth, height);

  const coolTop = shiftHex(palette.background, -0.12);
  const extendedColors = extendedBackgroundPalette(palette, wave.variety);

  paintColorLobe(ctx, totalWidth, height, totalWidth * 0.04, height * 0.1, topRadius, palette, coolTop, wave, {
    presenceScale: 0.75,
    gain: 0.85,
  });
  paintColorLobe(
    ctx,
    totalWidth,
    height,
    totalWidth * 0.96,
    height * 0.08,
    topRadius * 1.05,
    palette,
    pickExtendedAt(extendedColors, 0.72),
    wave,
    { presenceScale: 0.9, gain: 0.95 },
  );
  paintColorLobe(
    ctx,
    totalWidth,
    height,
    totalWidth * 0.06,
    height * 0.94,
    cornerRadius * 1.08,
    palette,
    pickExtendedAt(extendedColors, 0.22),
    wave,
    { gain: 1.15 },
  );
  paintColorLobe(
    ctx,
    totalWidth,
    height,
    totalWidth * 0.94,
    height * 0.92,
    cornerRadius,
    palette,
    pickExtendedAt(extendedColors, 0.48),
    wave,
    { gain: 1.05 },
  );

  for (let i = 0; i < travelerCount; i += 1) {
    const t = driftScrambleT(i, phase, 19);
    const cx = t * totalWidth + Math.sin(t * Math.PI * 2.4 + phase) * slideWidth * 0.08;
    const cy = height * (0.72 + 0.14 * Math.sin(t * Math.PI * 1.5 + phase * 0.5));
    paintColorLobe(
      ctx,
      totalWidth,
      height,
      cx,
      cy,
      travelerRadius,
      palette,
      pickExtendedAt(extendedColors, (0.55 + i * 0.11) % 1),
      wave,
      {
        presenceScale: 0.82,
        gain: 0.75,
      },
    );
  }

  return canvas;
}

/**
 * @param {number} slideWidth
 * @param {number} height
 * @param {PanoramaPalette} palette
 * @param {number} slideCount
 * @param {BackgroundWaveConfig} wave
 */
export function buildPanoramaWaveCanvas(slideWidth, height, palette, slideCount, wave) {
  if (wave.style === 'mesh-corners') {
    return buildMeshCornersCanvas(slideWidth, height, palette, slideCount, wave);
  }
  return buildDriftWaveCanvas(slideWidth, height, palette, slideCount, wave);
}

/**
 * Panorama at target slide size; supersamples washes then filters down to reduce banding.
 * @param {number} slideWidth
 * @param {number} height
 * @param {PanoramaPalette} palette
 * @param {number} slideCount
 * @param {BackgroundWaveConfig} wave
 */
function buildPanoramaWaveCanvasForDisplay(slideWidth, height, palette, slideCount, wave) {
  const ss = PANORAMA_SUPERSAMPLE;
  if (ss <= 1) {
    return buildPanoramaWaveCanvas(slideWidth, height, palette, slideCount, wave);
  }
  const hi = buildPanoramaWaveCanvas(slideWidth * ss, height * ss, palette, slideCount, wave);
  const loW = slideWidth * slideCount;
  const loH = height;
  const lo = document.createElement('canvas');
  lo.width = loW;
  lo.height = loH;
  const ctx = lo.getContext('2d');
  if (!ctx) throw new Error('Canvas 2D context unavailable');
  ctx.imageSmoothingEnabled = true;
  ctx.imageSmoothingQuality = 'high';
  ctx.drawImage(hi, 0, 0, hi.width, hi.height, 0, 0, loW, loH);
  return lo;
}

/**
 * Paint one continuous panoramic background across strip slots (with optional gaps).
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} slideWidth
 * @param {number} height
 * @param {number} slideCount
 * @param {PanoramaPalette} palette
 * @param {BackgroundWaveConfig} wave
 * @param {{ gap?: number, startX?: number, slotSlideIndices?: number[] }} [layout]
 */
export function paintPanoramicStripBackground(ctx, slideWidth, height, deckSlideCount, palette, wave, layout = {}) {
  const gap = layout.gap ?? 0;
  const startX = layout.startX ?? 0;
  const slotSlideIndices = layout.slotSlideIndices;
  const slotCount = slotSlideIndices?.length ?? deckSlideCount;
  const paintPalette = paletteWithBackgroundHueShift(palette, wave);

  if (!shouldPaintPanoramicWave(wave)) {
    ctx.fillStyle = paintPalette.background;
    for (let slot = 0; slot < slotCount; slot += 1) {
      const x = startX + slot * (slideWidth + gap);
      ctx.fillRect(x, 0, slideWidth, height);
    }
    return;
  }

  const cacheKey = JSON.stringify({
    width: slideWidth,
    height,
    slideCount: deckSlideCount,
    palette: paintPalette,
    wave,
    mode: 'strip',
  });

  let panoramaCanvas = panoramaCache.get(cacheKey);
  if (!panoramaCanvas) {
    panoramaCanvas = buildPanoramaWaveCanvasForDisplay(
      slideWidth,
      height,
      paintPalette,
      deckSlideCount,
      wave,
    );
    storePanoramaCache(cacheKey, panoramaCanvas);
  }

  for (let slot = 0; slot < slotCount; slot += 1) {
    const slideIndex = slotSlideIndices?.[slot] ?? slot;
    const x = startX + slot * (slideWidth + gap);
    ctx.drawImage(
      panoramaCanvas,
      slideIndex * slideWidth,
      0,
      slideWidth,
      height,
      x,
      0,
      slideWidth,
      height,
    );
  }
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} width
 * @param {number} height
 * @param {PanoramaPalette} palette
 * @param {BackgroundPanoramaContext} panorama
 * @param {BackgroundWaveConfig} wave
 */
export function paintPanoramicWaveSlice(ctx, width, height, palette, panorama, wave) {
  const { slideIndex, slideCount } = panorama;
  if (slideCount < 1) return;

  const paintPalette = paletteWithBackgroundHueShift(palette, wave);

  const cacheKey = JSON.stringify({
    width,
    height,
    slideCount,
    palette: paintPalette,
    wave,
  });

  let panoramaCanvas = panoramaCache.get(cacheKey);
  if (!panoramaCanvas) {
    panoramaCanvas = buildPanoramaWaveCanvasForDisplay(width, height, paintPalette, slideCount, wave);
    storePanoramaCache(cacheKey, panoramaCanvas);
  }

  const sx = slideIndex * width;
  ctx.drawImage(panoramaCanvas, sx, 0, width, height, 0, 0, width, height);
}
