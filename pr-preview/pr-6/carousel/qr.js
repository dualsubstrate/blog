/**
 * Lightweight QR rendering glue for the carousel renderer.
 *
 * Uses Nayuki's MIT-licensed QR encoder vendored as a classic script:
 * `static/carousel/vendor/qrcodegen-v1.8.0-es6.js`.
 */

let qrLibPromise = null;

/**
 * Load the QR encoder into the page global scope (classic script).
 * @returns {Promise<any>} Resolves to global `qrcodegen` namespace.
 */
export async function ensureQrCodegenLoaded() {
  // eslint-disable-next-line no-undef
  if (globalThis.qrcodegen && globalThis.qrcodegen.QrCode) return globalThis.qrcodegen;
  if (qrLibPromise) return qrLibPromise;

  const src = new URL('./vendor/qrcodegen-v1.8.0-es6.js', import.meta.url).href;
  qrLibPromise = new Promise((resolve, reject) => {
    const existing = document.querySelector(`script[data-carousel-qr="1"][src="${src}"]`);
    if (existing) {
      existing.addEventListener('load', () => resolve(globalThis.qrcodegen));
      existing.addEventListener('error', () => reject(new Error('Failed to load QR encoder')));
      return;
    }

    const script = document.createElement('script');
    script.src = src;
    script.async = true;
    script.defer = true;
    script.dataset.carouselQr = '1';
    script.onload = () => resolve(globalThis.qrcodegen);
    script.onerror = () => reject(new Error('Failed to load QR encoder'));
    document.head.appendChild(script);
  });

  return qrLibPromise;
}

/**
 * @param {string} text
 * @param {'LOW'|'MEDIUM'|'QUARTILE'|'HIGH'} [ecc]
 */
export function makeQrCode(text, ecc = 'MEDIUM') {
  // eslint-disable-next-line no-undef
  const lib = globalThis.qrcodegen;
  if (!lib?.QrCode) throw new Error('QR encoder not loaded');
  const QrCode = lib.QrCode;
  const level = ecc === 'LOW'
    ? QrCode.Ecc.LOW
    : ecc === 'QUARTILE'
      ? QrCode.Ecc.QUARTILE
      : ecc === 'HIGH'
        ? QrCode.Ecc.HIGH
        : QrCode.Ecc.MEDIUM;
  return QrCode.encodeText(String(text || ''), level);
}

/**
 * QR size as percent of the vertical slot between URL and scan footer (1–100).
 * Accepts `100`, `"100%"`, or legacy px values above 100 at 1080.
 * @param {number|string|null|undefined} raw
 */
export function parseQrSizePercent(raw) {
  if (typeof raw === 'string') {
    const trimmed = raw.trim();
    const withPct = trimmed.match(/^(\d+(?:\.\d+)?)\s*%$/);
    if (withPct) {
      return Math.min(100, Math.max(1, Number(withPct[1])));
    }
    const bare = Number(trimmed);
    if (Number.isFinite(bare)) {
      if (bare > 100) return 100;
      if (bare <= 0) return 100;
      return bare;
    }
    return 100;
  }
  if (!Number.isFinite(raw) || raw <= 0) return 100;
  if (raw > 100) return 100;
  return raw;
}

/**
 * @param {string|{color?: string}|null|undefined} light
 */
export function isQrBackgroundTransparent(light) {
  if (light == null || light === '') return false;
  if (typeof light === 'object') {
    return isQrBackgroundTransparent(light.color);
  }
  const t = String(light).trim().toLowerCase();
  return t === 'transparent' || t === 'none';
}

/**
 * Draw a QR code onto a canvas.
 * Scales modules to fill the target box edge to edge (quiet zone included).
 *
 * @param {CanvasRenderingContext2D} ctx
 * @param {{ size: number, getModule: (x: number, y: number) => boolean }} qr
 * @param {number} x Left of target box (px)
 * @param {number} y Top of target box (px)
 * @param {number} boxSize Target square size (px)
 * @param {{ marginModules?: number, dark?: string, light?: string, radius?: number }} [opts]
 */
export function drawQrCodeInBox(ctx, qr, x, y, boxSize, opts = {}) {
  const marginModules = Number.isFinite(opts.marginModules) ? opts.marginModules : 2;
  const dark = opts.dark || '#000000';
  const light = opts.light ?? '#ffffff';

  const moduleCount = (qr?.size ?? 0) + 2 * marginModules;
  if (!qr || !qr.size || moduleCount <= 0) return;

  const modulePx = boxSize / moduleCount;

  ctx.save();

  if (!isQrBackgroundTransparent(light)) {
    ctx.fillStyle = light;
    ctx.fillRect(x, y, boxSize, boxSize);
  }

  ctx.fillStyle = dark;
  for (let my = 0; my < qr.size; my += 1) {
    for (let mx = 0; mx < qr.size; mx += 1) {
      if (!qr.getModule(mx, my)) continue;
      const px = x + (mx + marginModules) * modulePx;
      const py = y + (my + marginModules) * modulePx;
      const px2 = x + (mx + marginModules + 1) * modulePx;
      const py2 = y + (my + marginModules + 1) * modulePx;
      ctx.fillRect(px, py, px2 - px, py2 - py);
    }
  }

  ctx.restore();
}
