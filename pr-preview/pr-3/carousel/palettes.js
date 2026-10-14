/**
 * Coordinated deck palettes: base colors + default gradient preset.
 * @typedef {import('./theme.js').DeckPalette} DeckPalette
 * @typedef {Object} CarouselPalette
 * @property {string} id
 * @property {string} label
 * @property {DeckPalette} palette
 * @property {string|null} [backgroundGradient] Default gradient preset id
 */

/** @type {CarouselPalette[]} */
export const CAROUSEL_PALETTES = [
  {
    id: 'factory-warm',
    label: 'Factory warm',
    palette: {
      background: '#1a1e26',
      text: '#f5f5f0',
      muted: '#988c81',
      accent1: '#df9311',
      accent2: '#e77218',
    },
    backgroundGradient: 'duo-accent',
  },
  {
    id: 'editorial-trio',
    label: 'Editorial trio',
    palette: {
      background: '#262030',
      text: '#f3f0f5',
      muted: '#95a090',
      accent1: '#9d72b8',
      accent2: '#5c9eb0',
    },
    backgroundGradient: null,
  },
  {
    id: 'ember-peach',
    label: 'Ember peach',
    palette: {
      background: '#1a1e26',
      text: '#f5f0ea',
      muted: '#9a8f88',
      accent1: '#d9937e',
      accent2: '#e6c09a',
    },
    backgroundGradient: null,
  },
  {
    id: 'slate-cool',
    label: 'Slate cool',
    palette: {
      background: '#171c26',
      text: '#edf2f7',
      muted: '#8493a6',
      accent1: '#9ebde0',
      accent2: '#6d94b8',
    },
    backgroundGradient: null,
  },
  {
    id: 'ocean-depth',
    label: 'Ocean depth',
    palette: {
      background: '#1a3352',
      text: '#f2f6fa',
      muted: '#8fa8be',
      accent1: '#e8b060',
      accent2: '#6eb8cf',
    },
    backgroundGradient: null,
  },
  {
    id: 'paper-light',
    label: 'Paper light',
    palette: {
      background: '#f5f5f0',
      text: '#1a1e26',
      muted: '#616a74',
      accent1: '#2f5d82',
      accent2: '#9a7b4f',
    },
    backgroundGradient: null,
  },
  {
    id: 'sage-light',
    label: 'Sage light',
    palette: {
      background: '#eef3eb',
      text: '#1a2420',
      muted: '#5c6b62',
      accent1: '#6d8f5f',
      accent2: '#4a8fa3',
    },
    backgroundGradient: null,
  },
  {
    id: 'wine-depth',
    label: 'Wine depth',
    palette: {
      background: '#130101',
      text: '#dbd4c7',
      muted: '#f49325',
      accent1: '#e18a5b',
      accent2: '#c64601',
    },
    backgroundGradient: 'soft-warm',
  },
  {
    id: 'forest-green',
    label: 'Forest green',
    palette: {
      background: '#1a2822',
      text: '#f5f5f0',
      muted: '#8a9e92',
      accent1: '#df9311',
      accent2: '#6d8f5f',
    },
    backgroundGradient: 'soft-cool',
  },
  {
    id: 'warm-linen',
    label: 'Warm linen',
    palette: {
      background: '#f5efe3',
      text: '#2c251c',
      muted: '#8a7868',
      accent1: '#b85c38',
      accent2: '#c49358',
    },
    backgroundGradient: null,
  },
];

/** @typedef {{ key: keyof DeckPalette, label: string, hex: string }} PaletteColorEntry */

/**
 * Base deck colors for palette cards and chips (background gradient is separate).
 * @param {DeckPalette} palette
 * @returns {PaletteColorEntry[]}
 */
export function paletteBaseColorEntries(palette) {
  return [
    { key: 'background', label: 'Base', hex: palette.background },
    { key: 'text', label: 'Text', hex: palette.text },
    { key: 'muted', label: 'Muted', hex: palette.muted },
    { key: 'accent1', label: 'Accent 1', hex: palette.accent1 },
    { key: 'accent2', label: 'Accent 2', hex: palette.accent2 || palette.accent1 },
  ];
}

/** @typedef {{ key: 'background' | 'muted' | 'accent1' | 'accent2', label: string, hex: string }} WavePaletteColorEntry */

/**
 * Wash-only colors (no text); used when wave palette is unlinked from text palette.
 * @param {import('./theme.js').WavePalette} palette
 * @returns {WavePaletteColorEntry[]}
 */
export function wavePaletteColorEntries(palette) {
  return [
    { key: 'background', label: 'Base', hex: palette.background },
    { key: 'muted', label: 'Muted', hex: palette.muted },
    { key: 'accent1', label: 'Accent 1', hex: palette.accent1 },
    { key: 'accent2', label: 'Accent 2', hex: palette.accent2 || palette.accent1 },
  ];
}

/** Site But Why card left stripe (sage → plum → teal). */
export const EDITORIAL_STRIPE_COLORS = {
  sage: '#6d8f5f',
  plum: '#5b1273',
  teal: '#4a8fa3',
};

/**
 * @param {string|null|undefined} paletteId
 * @returns {CarouselPalette|null}
 */
export function paletteById(paletteId) {
  if (!paletteId) return null;
  return CAROUSEL_PALETTES.find((p) => p.id === paletteId) ?? null;
}

/**
 * Prefer an explicit studio selection when colors still match that preset.
 * @param {import('./theme.js').DeckTheme} theme
 * @param {string|null} [preferredId]
 * @returns {string|null}
 */
export function resolvePaletteId(theme, preferredId = null) {
  if (preferredId) {
    const preset = paletteById(preferredId);
    if (preset && paletteColorsMatch(theme, preset.palette)) {
      return preferredId;
    }
  }
  return matchPaletteId(theme);
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @returns {string|null}
 */
export function matchPaletteId(theme) {
  /** @type {CarouselPalette[]} */
  const colorMatches = [];
  for (const preset of CAROUSEL_PALETTES) {
    if (paletteColorsMatch(theme, preset.palette)) {
      colorMatches.push(preset);
    }
  }
  if (colorMatches.length === 0) return null;
  if (colorMatches.length === 1) return colorMatches[0].id;

  const activeGradient = theme.backgroundGradientPreset || null;
  const gradientMatch = colorMatches.find(
    (preset) => (preset.backgroundGradient ?? null) === activeGradient,
  );
  return (gradientMatch || colorMatches[0]).id;
}

/**
 * @param {import('./theme.js').DeckTheme} theme
 * @param {DeckPalette} palette
 */
function paletteColorsMatch(theme, palette) {
  return (
    toColorInputValue(theme.background) === toColorInputValue(palette.background)
    && toColorInputValue(theme.text) === toColorInputValue(palette.text)
    && toColorInputValue(theme.muted) === toColorInputValue(palette.muted)
    && toColorInputValue(theme.accent1) === toColorInputValue(palette.accent1)
    && toColorInputValue(theme.accent2 || theme.accent1) === toColorInputValue(palette.accent2 || palette.accent1)
  );
}

/** @param {string|null|undefined} value */
function toColorInputValue(value) {
  if (!value || typeof value !== 'string') return '';
  const trimmed = value.trim();
  if (/^#[0-9a-fA-F]{6}$/.test(trimmed)) return trimmed.toLowerCase();
  if (/^[0-9a-fA-F]{6}$/.test(trimmed)) return `#${trimmed.toLowerCase()}`;
  return trimmed.toLowerCase();
}
