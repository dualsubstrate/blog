const STORAGE_PREFIX = 'carousel-studio-strip-included';
const VARIANT_STORAGE_PREFIX = 'carousel-studio-strip-variant';

/** @param {string} slug */
export function inclusionStorageKey(slug) {
  return `${STORAGE_PREFIX}:${slug}`;
}

/** @param {string} slug */
export function variantStorageKey(slug) {
  return `${VARIANT_STORAGE_PREFIX}:${slug}`;
}

/**
 * @typedef {Object} CarouselDeckRef
 * @property {Record<string, unknown>} [deck]
 * @property {Array<{ number: number, role?: string, variants?: unknown[] }>} slides
 */

/**
 * Default included slide numbers: every slide in the deck.
 * @param {CarouselDeckRef} deck
 * @returns {Set<number>}
 */
export function defaultIncludedSlideNumbers(deck) {
  return new Set((deck.slides ?? []).map((slide) => slide.number));
}

/**
 * @param {string} slug
 * @param {CarouselDeckRef} deck
 * @returns {Set<number>}
 */
export function readIncludedSlideNumbers(slug, deck) {
  const defaults = defaultIncludedSlideNumbers(deck);
  const deckNumbers = new Set((deck.slides ?? []).map((slide) => slide.number));

  try {
    const raw = localStorage.getItem(inclusionStorageKey(slug));
    if (raw == null) return defaults;

    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) return defaults;

    /** @type {Set<number>} */
    const saved = new Set();
    for (const item of parsed) {
      const num = Number(item);
      if (Number.isFinite(num) && deckNumbers.has(num)) {
        saved.add(num);
      }
    }
    return saved;
  } catch {
    return defaults;
  }
}

/**
 * @param {string} slug
 * @param {Set<number>|number[]} included
 */
export function writeIncludedSlideNumbers(slug, included) {
  const numbers = [...included].sort((a, b) => a - b);
  try {
    localStorage.setItem(inclusionStorageKey(slug), JSON.stringify(numbers));
  } catch {
    // ignore private mode / quota
  }
}

/**
 * Slide numbers excluded from the studio strip preview.
 * @param {CarouselDeckRef} deck
 * @param {string} slug
 * @returns {number[]}
 */
export function studioExtraExcludeSlideNumbers(deck, slug) {
  const included = readIncludedSlideNumbers(slug, deck);
  return (deck.slides ?? [])
    .map((slide) => slide.number)
    .filter((number) => !included.has(number));
}

/**
 * @param {string} slug
 * @param {CarouselDeckRef} deck
 */
export function resetIncludedSlideNumbers(slug, deck) {
  try {
    localStorage.removeItem(inclusionStorageKey(slug));
  } catch {
    // ignore
  }
  return defaultIncludedSlideNumbers(deck);
}

/**
 * @param {CarouselDeckRef} deck
 * @returns {Map<number, number>}
 */
export function defaultStripVariantIndices(deck) {
  /** @type {Map<number, number>} */
  const map = new Map();
  for (const slide of deck.slides ?? []) {
    map.set(slide.number, 0);
  }
  return map;
}

/**
 * @param {string} slug
 * @param {CarouselDeckRef} deck
 * @returns {Map<number, number>}
 */
export function readStripVariantIndices(slug, deck) {
  const defaults = defaultStripVariantIndices(deck);
  try {
    const raw = localStorage.getItem(variantStorageKey(slug));
    if (raw == null) return defaults;

    const parsed = JSON.parse(raw);
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) return defaults;

    for (const slide of deck.slides ?? []) {
      const rawIndex = /** @type {Record<string, unknown>} */ (parsed)[String(slide.number)];
      const idx = Number(rawIndex);
      const maxIndex = Math.max(0, (slide.variants?.length ?? 1) - 1);
      if (Number.isFinite(idx) && idx >= 0 && idx <= maxIndex) {
        defaults.set(slide.number, idx);
      }
    }
    return defaults;
  } catch {
    return defaults;
  }
}

/**
 * @param {string} slug
 * @param {Map<number, number>|ReadonlyMap<number, number>} indices
 */
export function writeStripVariantIndices(slug, indices) {
  /** @type {Record<string, number>} */
  const payload = {};
  for (const [slideNumber, variantIndex] of indices) {
    payload[String(slideNumber)] = variantIndex;
  }
  try {
    localStorage.setItem(variantStorageKey(slug), JSON.stringify(payload));
  } catch {
    // ignore private mode / quota
  }
}

/**
 * @param {string} slug
 * @param {CarouselDeckRef} deck
 * @param {number} slideNumber
 * @param {Map<number, number>} indices
 * @returns {number}
 */
export function stripVariantIndexFor(slug, deck, slideNumber, indices) {
  const slide = (deck.slides ?? []).find((entry) => entry.number === slideNumber);
  const maxIndex = Math.max(0, (slide?.variants?.length ?? 1) - 1);
  const saved = indices.get(slideNumber);
  if (Number.isFinite(saved) && saved >= 0 && saved <= maxIndex) {
    return saved;
  }
  return readStripVariantIndices(slug, deck).get(slideNumber) ?? 0;
}

/**
 * @param {string} slug
 * @param {CarouselDeckRef} deck
 */
export function resetStripVariantIndices(slug, deck) {
  try {
    localStorage.removeItem(variantStorageKey(slug));
  } catch {
    // ignore
  }
  return defaultStripVariantIndices(deck);
}
