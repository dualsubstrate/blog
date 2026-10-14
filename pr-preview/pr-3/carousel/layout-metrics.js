import { typoOrangeLinesInEmBox } from './typo-probes.js';

/** Gap between stacked em boxes (canvas px at export size). Matches body cluster rhythm. */
export const EM_LINE_GAP_PX = 2;

/**
 * @typedef {Object} PreparedInkShape
 * @property {import('./inline-text.js').InlineRun[][]} inlineLines
 * @property {number} [lineHeightMult]
 * @property {number} [lineSlotPx]
 * @property {number} [emBoxHeight]
 * @property {number} [fontSizePx]
 * @property {number} [ascenderLinePx]
 * @property {number} [descenderLinePx]
 * @property {number} [xHeightPx]
 * @property {number} [capHeightPx]
 * @property {number} [emBoxAscent]
 */

/**
 * @param {number} emTop
 * @param {PreparedInkShape} prepared
 */
function emBottomAfterLine(emTop, prepared) {
  return typoOrangeLinesInEmBox(emTop, prepared, prepared.lineHeightMult).emBottom;
}

/**
 * Total vertical ink from emTop 0 through all inline lines (matches drawPreparedBlock / drawClusterLines).
 * @param {PreparedInkShape} prepared
 */
export function measurePreparedInkHeight(prepared) {
  const lines = prepared.inlineLines?.length ?? 0;
  if (lines <= 0) return 0;
  let emTop = 0;
  for (let i = 0; i < lines; i += 1) {
    if (i > 0) {
      emTop = emBottomAfterLine(emTop, prepared) + EM_LINE_GAP_PX;
    }
  }
  const emStack = emBottomAfterLine(emTop, prepared);
  const slot = prepared.lineSlotPx ?? prepared.emBoxHeight;
  const slotStack = slot != null && Number.isFinite(slot)
    ? lines * slot + Math.max(0, lines - 1) * EM_LINE_GAP_PX
    : 0;
  const fallback = Number.isFinite(prepared.fontSizePx) ? prepared.fontSizePx : 0;
  return Math.max(emStack, slotStack, fallback);
}

/**
 * Grid row sizing: never below the blue line slot stack (matches drawPreparedBlock).
 * @param {PreparedInkShape} prepared
 */
export function measureGridCellInkExtent(prepared) {
  return measurePreparedInkHeight(prepared);
}
