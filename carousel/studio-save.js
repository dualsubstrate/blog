import { downloadBlob } from './export.js';

export const DEFAULT_CAROUSEL_SAVE_URL = 'http://127.0.0.1:3848/api/carousel-json';

/**
 * @param {string} [deckUrl]
 * @returns {string}
 */
export function carouselFilenameFromDeckUrl(deckUrl) {
  const raw = String(deckUrl || './carousel.json').split('?')[0];
  const base = raw.split('/').pop() || 'carousel.json';
  if (base === 'carousel.es.json') return 'carousel.es.json';
  return 'carousel.json';
}

/**
 * @param {Record<string, unknown>} deck
 * @returns {string}
 */
export function serializeCarouselDeck(deck) {
  return `${JSON.stringify(deck, null, 2)}\n`;
}

/**
 * @param {string} saveUrl
 * @param {{ source: string, filename: string, body: string }} payload
 * @returns {Promise<{ path?: string }>}
 */
export async function postCarouselJson(saveUrl, payload) {
  const response = await fetch(saveUrl, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  });
  const text = await response.text();
  let parsed = {};
  if (text) {
    try {
      parsed = JSON.parse(text);
    } catch {
      parsed = { error: text };
    }
  }
  if (!response.ok) {
    const message = typeof parsed.error === 'string' && parsed.error
      ? parsed.error
      : `Save failed (${response.status})`;
    throw new Error(message);
  }
  return parsed;
}

/**
 * @param {string} body
 * @param {string} filename
 * @returns {Promise<'picker'|'download'>}
 */
export async function saveCarouselJsonFallback(body, filename) {
  if (typeof window.showSaveFilePicker === 'function') {
    const handle = await window.showSaveFilePicker({
      suggestedName: filename,
      types: [{ description: 'Carousel JSON', accept: { 'application/json': ['.json'] } }],
    });
    const writable = await handle.createWritable();
    await writable.write(body);
    await writable.close();
    return 'picker';
  }
  downloadBlob(new Blob([body], { type: 'application/json' }), filename);
  return 'download';
}

/**
 * Write live deck JSON to the bundle file (local save server) or a file picker.
 * @param {Record<string, unknown>} deck
 * @param {{ deckUrl?: string, saveUrl?: string }} [options]
 * @returns {Promise<'server'|'picker'|'download'>}
 */
export async function saveCarouselDeckToSource(deck, options = {}) {
  const filename = carouselFilenameFromDeckUrl(options.deckUrl);
  const source = typeof deck.source === 'string' ? deck.source.trim() : '';
  if (!source) {
    throw new Error('Deck is missing source (path to index.md). Cannot save to carousel.json.');
  }
  const body = serializeCarouselDeck(deck);
  const saveUrl = options.saveUrl || DEFAULT_CAROUSEL_SAVE_URL;
  try {
    await postCarouselJson(saveUrl, { source, filename, body });
    return 'server';
  } catch (error) {
    if (error instanceof Error && error.name === 'AbortError') {
      throw error;
    }
    const networkFail = error instanceof TypeError
      || (error instanceof Error && /Failed to fetch|NetworkError|Load failed/i.test(error.message));
    if (!networkFail) {
      throw error;
    }
    return saveCarouselJsonFallback(body, filename);
  }
}
