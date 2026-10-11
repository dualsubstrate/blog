/**
 * Carousel raster constants (single source of truth for JS tooling).
 *
 * Keep in sync with `internal/carousel/slide.go`.
 *
 * - `CAROUSEL_SLIDE_WIDTH_PX`: rendered slide canvas width (1080).
 * - `PANORAMA_SLIDE_WIDTH_PX` / `PANORAMA_GAP_PX`: studio panorama export slice + separator.
 */

/** Rendered carousel slide canvas width (px). */
export const CAROUSEL_SLIDE_WIDTH_PX = 1080;

/** Studio panorama export width per slide (default `theme.size` for Download panorama). */
export const PANORAMA_SLIDE_WIDTH_PX = 600;

/** Inter-slide gap in studio panorama export at {@link PANORAMA_SLIDE_WIDTH_PX}. */
export const PANORAMA_GAP_PX = 5;

/** Reference pair for scaling gap when `theme.size` differs from panorama default. */
export const PANORAMA_GAP_PX_AT_REFERENCE = 4;
export const PANORAMA_GAP_REFERENCE_SLIDE_WIDTH_PX = 480;

/**
 * Gap between slides in strip preview and panorama export.
 *
 * @param {number} [slideWidth]
 * @returns {number}
 */
export function stripSlideGapPx(slideWidth = PANORAMA_SLIDE_WIDTH_PX) {
  if (slideWidth === PANORAMA_SLIDE_WIDTH_PX) {
    return PANORAMA_GAP_PX;
  }
  const width = Number.isFinite(slideWidth) && slideWidth > 0 ? slideWidth : PANORAMA_SLIDE_WIDTH_PX;
  return Math.max(2, Math.round(
    width * (PANORAMA_GAP_PX_AT_REFERENCE / PANORAMA_GAP_REFERENCE_SLIDE_WIDTH_PX),
  ));
}

/**
 * Panorama export width including inter-slide gaps.
 *
 * @param {number} slideCount
 * @param {number} [slideWidth]
 * @returns {number}
 */
export function panoramaWidthWithGapsPx(slideCount, slideWidth = PANORAMA_SLIDE_WIDTH_PX) {
  const n = Math.max(1, Math.floor(Number(slideCount) || 0));
  const gap = stripSlideGapPx(slideWidth);
  return n * slideWidth + (n - 1) * gap;
}

/**
 * Seamless motif strip width (gaps removed): slides × panorama slice width.
 *
 * @param {number} slideCount
 * @param {number} [slideWidth]
 * @returns {number}
 */
export function motifStripSeamlessWidthPx(slideCount, slideWidth = PANORAMA_SLIDE_WIDTH_PX) {
  const n = Math.max(1, Math.floor(Number(slideCount) || 0));
  return n * slideWidth;
}

/**
 * @param {number} slideCount
 * @returns {number}
 */
export function motifStripWidthPx(slideCount) {
  return motifStripSeamlessWidthPx(slideCount);
}

/**
 * CTA / QR end frames (`role: cta` or `post_cta`) skip panoramic motif and wave.
 * @param {string} [role]
 */
export function isCarouselCtaRole(role) {
  const normalized = (role || '').trim().toLowerCase();
  return normalized === 'cta' || normalized === 'post_cta';
}

/**
 * Remove studio panorama inter-slide gaps when width matches {@link panoramaWidthWithGapsPx}.
 *
 * @param {number} slideCount
 * @param {number} [slideWidth]
 * @param {number} [gapPx]
 * @returns {boolean}
 */
export function panoramaHasSeparators(width, slideCount, slideWidth = PANORAMA_SLIDE_WIDTH_PX, gapPx = stripSlideGapPx(slideWidth)) {
  const n = Math.max(1, Math.floor(Number(slideCount) || 0));
  return width === panoramaWidthWithGapsPx(n, slideWidth);
}
