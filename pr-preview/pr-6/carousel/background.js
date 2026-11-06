import { EDITORIAL_STRIPE_COLORS } from './palettes.js';

/**
 * @typedef {Object} GradientStop
 * @property {number} at 0–1
 * @property {string} color `#rrggbb`
 */

/**
 * @typedef {Object} BackgroundLayer
 * @property {'linear' | 'stripe-left'} type
 * @property {number} [angle] CSS degrees for linear layers
 * @property {GradientStop[]} stops
 * @property {number} [alpha] 0–1 layer opacity (default 1)
 * @property {number} [width] Stripe width as fraction of canvas (stripe-left)
 * @property {GlobalCompositeOperation} [compositeOperation]
 */

/**
 * @typedef {Object} BackgroundGradient
 * @property {'linear' | 'composite'} [kind]
 * @property {number} [angle] Legacy single linear layer
 * @property {GradientStop[]} [stops] Legacy single linear layer
 * @property {BackgroundLayer[]} [layers] Multi-layer / mixed gradients
 */

/** @type {{ id: string, label: string }[]} */
export const BACKGROUND_GRADIENT_PRESETS = [
  { id: 'none', label: 'Solid (no gradient)' },
  { id: 'soft-warm', label: 'Soft warm diagonal' },
  { id: 'soft-cool', label: 'Soft cool diagonal' },
  { id: 'depth', label: 'Depth (lighter top)' },
  { id: 'accent-glow', label: 'Accent glow (top-right)' },
  { id: 'duo-accent', label: 'Dual accent wash' },
  { id: 'left-glow', label: 'Left accent glow' },
  { id: 'mixed-wash', label: 'Mixed wash (crossed diagonals)' },
  { id: 'editorial-stripe', label: 'Editorial stripe (mixed wash + tri-color edge)' },
];

/** @param {string} hex */
function parseHex(hex) {
  let raw = hex.replace('#', '').trim();
  if (raw.length === 3) {
    raw = raw[0] + raw[0] + raw[1] + raw[1] + raw[2] + raw[2];
  }
  if (raw.length !== 6) return { r: 26, g: 30, b: 38 };
  return {
    r: parseInt(raw.slice(0, 2), 16),
    g: parseInt(raw.slice(2, 4), 16),
    b: parseInt(raw.slice(4, 6), 16),
  };
}

/**
 * @param {string} hexA
 * @param {string} hexB
 * @param {number} t 0–1
 */
export function mixHex(hexA, hexB, t) {
  const a = parseHex(hexA);
  const b = parseHex(hexB);
  const clampT = Math.max(0, Math.min(1, t));
  const r = Math.round(a.r + (b.r - a.r) * clampT);
  const g = Math.round(a.g + (b.g - a.g) * clampT);
  const bl = Math.round(a.b + (b.b - a.b) * clampT);
  return `#${[r, g, bl].map((v) => v.toString(16).padStart(2, '0')).join('')}`;
}

/** @param {string} hex @param {number} amount -1 to 1 (negative = darken) */
export function shiftHex(hex, amount) {
  if (amount >= 0) return mixHex(hex, '#ffffff', amount);
  return mixHex(hex, '#000000', -amount);
}

/** @param {number} r @param {number} g @param {number} b */
function rgbToHsl(r, g, b) {
  const rn = r / 255;
  const gn = g / 255;
  const bn = b / 255;
  const max = Math.max(rn, gn, bn);
  const min = Math.min(rn, gn, bn);
  const l = (max + min) / 2;
  if (max === min) {
    return { h: 0, s: 0, l };
  }
  const d = max - min;
  const s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
  let h = 0;
  if (max === rn) {
    h = ((gn - bn) / d + (gn < bn ? 6 : 0)) / 6;
  } else if (max === gn) {
    h = ((bn - rn) / d + 2) / 6;
  } else {
    h = ((rn - gn) / d + 4) / 6;
  }
  return { h: h * 360, s, l };
}

/** @param {number} h @param {number} s @param {number} l */
function hslToRgb(h, s, l) {
  const hue = ((h % 360) + 360) % 360;
  if (s === 0) {
    const gray = Math.round(l * 255);
    return { r: gray, g: gray, b: gray };
  }
  const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
  const p = 2 * l - q;
  const hk = hue / 360;
  const hueToRgb = (t) => {
    let tt = t;
    if (tt < 0) tt += 1;
    if (tt > 1) tt -= 1;
    if (tt < 1 / 6) return p + (q - p) * 6 * tt;
    if (tt < 1 / 2) return q;
    if (tt < 2 / 3) return p + (q - p) * (2 / 3 - tt) * 6;
    return p;
  };
  return {
    r: Math.round(hueToRgb(hk + 1 / 3) * 255),
    g: Math.round(hueToRgb(hk) * 255),
    b: Math.round(hueToRgb(hk - 1 / 3) * 255),
  };
}

/**
 * Rotate hex color hue by degrees; saturation and lightness preserved.
 * @param {string} hex
 * @param {number} degrees
 * @returns {string}
 */
export function rotateHexHue(hex, degrees) {
  if (typeof hex !== 'string' || !Number.isFinite(degrees) || Math.abs(degrees) < 0.001) return hex;
  const { r, g, b } = parseHex(hex);
  const { h, s, l } = rgbToHsl(r, g, b);
  const rotated = hslToRgb(h + degrees, s, l);
  return `#${[rotated.r, rotated.g, rotated.b].map((v) => v.toString(16).padStart(2, '0')).join('')}`;
}

/** Stronger lighten/darken for visible depth on saturated base colors. */
function depthShift(hex, amount) {
  const boosted = amount >= 0 ? amount * 2.4 : amount * 2.2;
  return shiftHex(hex, boosted);
}

/**
 * @param {string|null|undefined} value
 * @returns {string|null}
 */
export function normalizeGradientPresetId(value) {
  if (value == null || value === '' || value === 'none' || value === 'solid') return null;
  if (typeof value === 'string') {
    const id = value.trim();
    if (BACKGROUND_GRADIENT_PRESETS.some((p) => p.id === id && p.id !== 'none')) return id;
    return null;
  }
  return null;
}

/**
 * @param {{ background: string, accent1: string, accent2?: string|null }} colors
 * @returns {BackgroundLayer[]}
 */
function buildMixedWashLayers(colors) {
  const base = colors.background || '#1a1e26';
  const accent1 = colors.accent1 || '#d69a80';
  const accent2 = colors.accent2 || accent1;
  const plum = EDITORIAL_STRIPE_COLORS.plum;
  const teal = EDITORIAL_STRIPE_COLORS.teal;
  const sage = EDITORIAL_STRIPE_COLORS.sage;

  return [
    {
      type: 'linear',
      angle: 145,
      stops: [
        { at: 0, color: mixHex(base, plum, 0.32) },
        { at: 0.48, color: mixHex(base, plum, 0.18) },
        { at: 1, color: mixHex(base, teal, 0.22) },
      ],
    },
    {
      type: 'linear',
      angle: 315,
      alpha: 0.88,
      compositeOperation: 'source-over',
      stops: [
        { at: 0, color: mixHex(base, sage, 0.34) },
        { at: 0.38, color: mixHex(base, accent1, 0.22) },
        { at: 1, color: base },
      ],
    },
    {
      type: 'linear',
      angle: 200,
      alpha: 0.72,
      stops: [
        { at: 0, color: mixHex(base, accent2, 0.26) },
        { at: 0.55, color: mixHex(base, accent2, 0.08) },
        { at: 1, color: base },
      ],
    },
  ];
}

/** @returns {BackgroundLayer} */
function buildEditorialStripeLayer() {
  const { sage, plum, teal } = EDITORIAL_STRIPE_COLORS;
  return {
    type: 'stripe-left',
    width: 0.014,
    stops: [
      { at: 0, color: sage },
      { at: 0.52, color: plum },
      { at: 1, color: teal },
    ],
  };
}

/**
 * @param {{ background: string, accent1: string, accent2?: string|null }} colors
 * @param {string|null} presetId
 * @returns {BackgroundGradient|null}
 */
export function buildBackgroundGradient(colors, presetId) {
  const base = colors.background || '#1a1e26';
  const accent1 = colors.accent1 || '#d69a80';
  const accent2 = colors.accent2 || accent1;

  /** @type {BackgroundLayer[]|null} */
  let layers = null;

  switch (presetId) {
    case 'soft-warm':
      layers = [{
        type: 'linear',
        angle: 135,
        stops: [
          { at: 0, color: mixHex(base, accent1, 0.18) },
          { at: 0.45, color: mixHex(base, accent1, 0.28) },
          { at: 1, color: mixHex(base, accent2, 0.22) },
        ],
      }];
      break;
    case 'soft-cool':
      layers = [{
        type: 'linear',
        angle: 145,
        stops: [
          { at: 0, color: shiftHex(base, 0.1) },
          { at: 0.45, color: mixHex(base, '#8aa3c4', 0.22) },
          { at: 1, color: shiftHex(base, -0.18) },
        ],
      }];
      break;
    case 'depth':
      layers = [{
        type: 'linear',
        angle: 180,
        stops: [
          { at: 0, color: depthShift(base, 0.1) },
          { at: 0.32, color: mixHex(depthShift(base, 0.04), base, 0.35) },
          { at: 0.68, color: base },
          { at: 1, color: depthShift(base, -0.14) },
        ],
      }];
      break;
    case 'accent-glow':
      layers = [{
        type: 'linear',
        angle: 225,
        stops: [
          { at: 0, color: mixHex(base, accent1, 0.42) },
          { at: 0.32, color: mixHex(base, accent1, 0.2) },
          { at: 0.62, color: mixHex(base, accent1, 0.06) },
          { at: 1, color: base },
        ],
      }];
      break;
    case 'duo-accent':
      layers = [{
        type: 'linear',
        angle: 120,
        stops: [
          { at: 0, color: mixHex(base, accent1, 0.32) },
          { at: 0.4, color: base },
          { at: 1, color: mixHex(base, accent2, 0.28) },
        ],
      }];
      break;
    case 'left-glow':
      layers = [{
        type: 'linear',
        angle: 90,
        stops: [
          { at: 0, color: mixHex(base, accent1, 0.55) },
          { at: 0.1, color: mixHex(base, accent1, 0.28) },
          { at: 0.22, color: mixHex(base, accent1, 0.08) },
          { at: 0.38, color: base },
          { at: 1, color: base },
        ],
      }];
      break;
    case 'mixed-wash':
      layers = buildMixedWashLayers(colors);
      break;
    case 'editorial-stripe':
      layers = [...buildMixedWashLayers(colors), buildEditorialStripeLayer()];
      break;
    default:
      return null;
  }

  if (!layers?.length) return null;

  if (layers.length === 1 && layers[0].type === 'linear') {
    const layer = layers[0];
    return {
      kind: 'linear',
      angle: layer.angle,
      stops: layer.stops,
      layers,
    };
  }

  return { kind: 'composite', layers };
}

/** @param {number} width @param {number} height @param {number} angleDeg */
function linearGradientEndpoints(width, height, angleDeg) {
  const rad = ((angleDeg - 90) * Math.PI) / 180;
  const cx = width / 2;
  const cy = height / 2;
  const len = Math.sqrt(width * width + height * height) / 2;
  return {
    x0: cx + Math.cos(rad) * len,
    y0: cy + Math.sin(rad) * len,
    x1: cx - Math.cos(rad) * len,
    y1: cy - Math.sin(rad) * len,
  };
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} width
 * @param {number} height
 * @param {BackgroundLayer} layer
 */
function paintBackgroundLayer(ctx, width, height, layer) {
  if (layer.type === 'stripe-left') {
    const stripeWidth = Math.max(4, Math.round(width * (layer.width ?? 0.012)));
    const g = ctx.createLinearGradient(0, 0, 0, height);
    for (const stop of layer.stops) {
      g.addColorStop(Math.max(0, Math.min(1, stop.at)), stop.color);
    }
    ctx.fillStyle = g;
    ctx.fillRect(0, 0, stripeWidth, height);
    return;
  }

  const { x0, y0, x1, y1 } = linearGradientEndpoints(width, height, layer.angle ?? 180);
  const g = ctx.createLinearGradient(x0, y0, x1, y1);
  for (const stop of layer.stops) {
    g.addColorStop(Math.max(0, Math.min(1, stop.at)), stop.color);
  }
  ctx.fillStyle = g;
  ctx.fillRect(0, 0, width, height);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {number} width
 * @param {number} height
 * @param {BackgroundGradient} gradient
 */
export function paintBackgroundGradient(ctx, width, height, gradient) {
  const layers = gradient.layers?.length
    ? gradient.layers
    : gradient.stops?.length >= 2
      ? [{ type: 'linear', angle: gradient.angle ?? 180, stops: gradient.stops }]
      : [];

  for (const layer of layers) {
    ctx.save();
    if (layer.alpha != null && layer.alpha < 1) {
      ctx.globalAlpha = layer.alpha;
    }
    if (layer.compositeOperation) {
      ctx.globalCompositeOperation = layer.compositeOperation;
    }
    paintBackgroundLayer(ctx, width, height, layer);
    ctx.restore();
  }
}

/**
 * @param {BackgroundGradient|null} gradient
 * @param {string} [fallbackColor]
 */
export function gradientToCss(gradient, fallbackColor = '#1a1e26') {
  if (!gradient) return fallbackColor;

  const layers = gradient.layers?.length
    ? gradient.layers
    : gradient.stops?.length >= 2
      ? [{ type: 'linear', angle: gradient.angle ?? 180, stops: gradient.stops }]
      : [];

  if (!layers.length) return fallbackColor;

  const cssLayers = layers
    .filter((layer) => layer.type === 'linear')
    .map((layer) => {
      const stops = layer.stops
        .map((s) => `${s.color} ${Math.round(s.at * 100)}%`)
        .join(', ');
      return `linear-gradient(${layer.angle ?? 180}deg, ${stops})`;
    });

  if (cssLayers.length === 0) return fallbackColor;
  if (cssLayers.length === 1) return cssLayers[0];
  return cssLayers.join(', ');
}
