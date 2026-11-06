/** WebP quality for slide export (0–1). */
export const EXPORT_WEBP_QUALITY = 0.92;

/**
 * Trigger a browser download from an in-memory blob.
 * @param {Blob} blob
 * @param {string} filename
 */
export function downloadBlob(blob, filename) {
  const url = URL.createObjectURL(blob);
  const anchor = document.createElement('a');
  anchor.href = url;
  anchor.download = filename;
  anchor.click();
  window.setTimeout(() => {
    URL.revokeObjectURL(url);
  }, 100);
}

/**
 * @param {HTMLCanvasElement} canvas
 * @param {string} filename Must end in `.webp` (or any name; MIME is WebP).
 */
export function downloadCanvasWebp(canvas, filename) {
  return new Promise((resolve, reject) => {
    canvas.toBlob(
      (blob) => {
        if (!blob) {
          reject(new Error('Canvas export produced empty blob'));
          return;
        }
        downloadBlob(blob, filename);
        window.setTimeout(() => {
          resolve(undefined);
        }, 100);
      },
      'image/webp',
      EXPORT_WEBP_QUALITY,
    );
  });
}

/**
 * Variant label for export filenames: 0 → `a`, 1 → `b`, … (26+ → `v27`, …).
 * @param {number} index 0-based index in `slide.variants`
 */
export function variantIdFromIndex(index) {
  const i = Math.max(0, Math.floor(index));
  if (i < 26) {
    return String.fromCharCode(97 + i);
  }
  return `v${i + 1}`;
}

/**
 * @param {string} slug
 * @param {number|string} slideNumber
 * @param {number} variantIndex 0-based index in `slide.variants`
 */
export function exportFilename(slug, slideNumber, variantIndex) {
  const num = String(slideNumber).padStart(2, '0');
  return `${slug}-slide-${num}-${variantIdFromIndex(variantIndex)}.webp`;
}

/**
 * @param {string} slug
 */
export function exportStripFilename(slug) {
  return `${slug}-panorama.webp`;
}

/**
 * @param {string} slug
 */
export function exportPdfFilename(slug) {
  return `${slug}-linkedin.pdf`;
}
