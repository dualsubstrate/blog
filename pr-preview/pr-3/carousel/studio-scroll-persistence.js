const STORAGE_PREFIX = 'carousel-studio-scroll';

/** @param {string} slug */
export function scrollStorageKey(slug) {
  return `${STORAGE_PREFIX}:${slug}`;
}

/**
 * @typedef {Object} StudioScrollState
 * @property {number} [windowX]
 * @property {number} [windowY]
 * @property {number} [leftFloatTop]
 * @property {number} [rightFloatTop]
 * @property {number} [stripLeft]
 */

/** @param {unknown} value @returns {number} */
function finiteScroll(value) {
  const num = Number(value);
  return Number.isFinite(num) && num >= 0 ? num : 0;
}

/**
 * @param {string} slug
 * @returns {StudioScrollState|null}
 */
export function readStudioScrollState(slug) {
  try {
    const raw = sessionStorage.getItem(scrollStorageKey(slug));
    if (raw == null) return null;
    const parsed = JSON.parse(raw);
    if (!parsed || typeof parsed !== 'object') return null;
    return {
      windowX: finiteScroll(parsed.windowX),
      windowY: finiteScroll(parsed.windowY),
      leftFloatTop: finiteScroll(parsed.leftFloatTop),
      rightFloatTop: finiteScroll(parsed.rightFloatTop),
      stripLeft: finiteScroll(parsed.stripLeft),
    };
  } catch {
    return null;
  }
}

/**
 * @param {string} slug
 * @param {StudioScrollState} state
 */
export function writeStudioScrollState(slug, state) {
  try {
    sessionStorage.setItem(scrollStorageKey(slug), JSON.stringify({
      windowX: finiteScroll(state.windowX),
      windowY: finiteScroll(state.windowY),
      leftFloatTop: finiteScroll(state.leftFloatTop),
      rightFloatTop: finiteScroll(state.rightFloatTop),
      stripLeft: finiteScroll(state.stripLeft),
    }));
  } catch {
    // ignore private mode / quota
  }
}

/** @param {ParentNode} root @returns {HTMLElement|null} */
export function studioLeftFloatBody(root) {
  const el = root.querySelector('.carousel-float--left .carousel-float-body');
  return el instanceof HTMLElement ? el : null;
}

/** @param {ParentNode} root @returns {HTMLElement|null} */
export function studioRightFloat(root) {
  const el = root.querySelector('.carousel-float--right');
  return el instanceof HTMLElement ? el : null;
}

/** @param {ParentNode} root @returns {HTMLElement|null} */
export function studioStripMount(root) {
  const el = root.querySelector('.carousel-vision-strip-mount:not(.carousel-vision-strip-mount--hidden)');
  return el instanceof HTMLElement ? el : null;
}

/**
 * @param {ParentNode} root
 * @param {StudioScrollState|null} [previous]
 * @returns {StudioScrollState}
 */
export function captureStudioScrollState(root, previous = null) {
  const leftBody = studioLeftFloatBody(root);
  const rightFloat = studioRightFloat(root);
  const strip = studioStripMount(root);
  const stripHasCanvas = Boolean(strip?.querySelector('canvas'));
  return {
    windowX: window.scrollX,
    windowY: window.scrollY,
    leftFloatTop: leftBody ? leftBody.scrollTop : 0,
    rightFloatTop: rightFloat ? rightFloat.scrollTop : 0,
    stripLeft: stripHasCanvas && strip
      ? strip.scrollLeft
      : finiteScroll(previous?.stripLeft),
  };
}

/**
 * @param {ParentNode} root
 * @param {StudioScrollState|null} state
 * @param {{ window?: boolean, leftFloat?: boolean, rightFloat?: boolean, strip?: boolean }} [parts]
 */
export function applyStudioScrollState(root, state, parts = {}) {
  if (!state) return;
  const include = {
    window: parts.window !== false,
    leftFloat: parts.leftFloat !== false,
    rightFloat: parts.rightFloat !== false,
    strip: parts.strip !== false,
  };

  if (include.window) {
    window.scrollTo(finiteScroll(state.windowX), finiteScroll(state.windowY));
  }
  if (include.leftFloat) {
    const leftBody = studioLeftFloatBody(root);
    if (leftBody) leftBody.scrollTop = finiteScroll(state.leftFloatTop);
  }
  if (include.rightFloat) {
    const rightFloat = studioRightFloat(root);
    if (rightFloat) rightFloat.scrollTop = finiteScroll(state.rightFloatTop);
  }
  if (include.strip) {
    const strip = studioStripMount(root);
    if (strip) strip.scrollLeft = finiteScroll(state.stripLeft);
  }
}

/**
 * Persist scroll while editing; sessionStorage survives refresh in this tab.
 * @param {string} slug
 * @param {ParentNode} root
 * @returns {() => void} cleanup
 */
export function bindStudioScrollPersistence(slug, root) {
  let persistTimer = 0;
  const persist = () => {
    writeStudioScrollState(slug, captureStudioScrollState(root, readStudioScrollState(slug)));
  };
  const onScroll = () => {
    window.clearTimeout(persistTimer);
    persistTimer = window.setTimeout(persist, 120);
  };

  const leftBody = studioLeftFloatBody(root);
  const rightFloat = studioRightFloat(root);
  const strip = studioStripMount(root);

  window.addEventListener('scroll', onScroll, { passive: true });
  leftBody?.addEventListener('scroll', onScroll, { passive: true });
  rightFloat?.addEventListener('scroll', onScroll, { passive: true });
  strip?.addEventListener('scroll', onScroll, { passive: true });
  window.addEventListener('pagehide', persist);

  return () => {
    window.clearTimeout(persistTimer);
    window.removeEventListener('scroll', onScroll);
    leftBody?.removeEventListener('scroll', onScroll);
    rightFloat?.removeEventListener('scroll', onScroll);
    strip?.removeEventListener('scroll', onScroll);
    window.removeEventListener('pagehide', persist);
  };
}
