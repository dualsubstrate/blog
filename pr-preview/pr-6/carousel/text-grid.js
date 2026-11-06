import { EM_LINE_GAP_PX, measureGridCellInkExtent, measurePreparedInkHeight } from './layout-metrics.js';
import { resolveWidth } from './theme.js';

/**
 * @typedef {import('./theme.js').DeckTheme} DeckTheme
 * @typedef {import('./theme.js').BodyEmphasis} BodyEmphasis
 * @typedef {import('./theme.js').ColorToken} ColorToken
 * @typedef {import('./theme.js').FontRole} FontRole
 * @typedef {Object} TextBlock
 * @property {string} text
 * @property {'header'|'body'|'footer'|'grid'} section
 * @property {BodyEmphasis} [emphasis]
 * @property {ColorToken} [color]
 * @property {FontRole} [font]
 * @property {string} [weight]
 * @property {number} [fontSize]
 * @property {number} [lineHeight]
 */

/**
 * @typedef {Object} PreparedBlock
 * @property {TextBlock} block
 * @property {string} section
 * @property {number} width
 * @property {number} height
 * @property {number} fontSizePx
 * @property {number} lineSlotPx
 * @property {number} emBoxHeight
 * @property {import('./inline-text.js').InlineRun[][]} inlineLines
 * @property {string} [color]
 */

/**
 * @typedef {Object} GridCellSpec
 * @property {number} row 0-based
 * @property {number} col 0-based
 * @property {number} [rowSpan] default 1
 * @property {number} [colSpan] default 1
 * @property {string} text
 * @property {import('./theme.js').BodyEmphasis} [emphasis]
 * @property {import('./theme.js').ColorToken} [color]
 * @property {import('./theme.js').FontRole} [font]
 * @property {string} [weight]
 * @property {number} [fontSize]
 * @property {number} [lineHeight]
 * @property {number} [maxLines] Override grid `cellMaxLines` for this cell
 */

/**
 * @typedef {Object} GridBlock
 * @property {'grid'} section
 * @property {number} [columns]
 * @property {number} [rows]
 * @property {number|string} [gap] Column gap if `columnGap` omitted; legacy fallback for `rowGap`
 * @property {number|string} [columnGap] Horizontal gap between columns (% or px at 1080)
 * @property {number|string} [rowGap] Vertical gap between rows (% or px); default matches body cluster (`2px`)
 * @property {number[]} [columnWidths] Optional fractions summing to 1
 * @property {'center'|'top'|'bottom'} [cellAlign] Vertical align inside each cell (default `top` when `rows` > 1)
 * @property {GridCellSpec[]} cells
 * @property {number} [cellMaxLines] Default max wrapped lines per cell (often `1`)
 */

/**
 * @typedef {Object} PreparedGridCell
 * @property {number} row
 * @property {number} col
 * @property {number} rowSpan
 * @property {number} colSpan
 * @property {PreparedBlock} prepared
 */

/**
 * @typedef {Object} PreparedGridSection
 * @property {'grid'} section
 * @property {GridBlock} block
 * @property {number} columns
 * @property {number} rows
 * @property {number} colGapPx
 * @property {number} rowGapPx
 * @property {number[]} colWidths
 * @property {number[]} rowHeights
 * @property {number} totalWidth
 * @property {number} totalHeight
 * @property {PreparedGridCell[]} cells
 */

/**
 * @param {GridBlock} block
 * @returns {GridBlock}
 */
export function normalizeGridBlock(block) {
  const cells = Array.isArray(block.cells) ? block.cells : [];
  return {
    ...block,
    section: 'grid',
    cells: cells.map((cell) => {
      const row = Math.max(0, Math.round(cell.row ?? 0));
      const col = Math.max(0, Math.round(cell.col ?? 0));
      const rowSpan = Math.max(1, Math.round(cell.rowSpan ?? 1));
      const colSpan = Math.max(1, Math.round(cell.colSpan ?? 1));
      return {
        emphasis: 'normal',
        ...cell,
        row,
        col,
        rowSpan,
        colSpan,
      };
    }),
  };
}

/**
 * @param {GridCellSpec[]} cells
 * @returns {{ columns: number, rows: number }}
 */
export function inferGridDimensions(cells) {
  let maxCol = 0;
  let maxRow = 0;
  for (const cell of cells) {
    const colEnd = cell.col + (cell.colSpan ?? 1);
    const rowEnd = cell.row + (cell.rowSpan ?? 1);
    maxCol = Math.max(maxCol, colEnd);
    maxRow = Math.max(maxRow, rowEnd);
  }
  return { columns: Math.max(1, maxCol), rows: Math.max(1, maxRow) };
}

/**
 * @param {number} columns
 * @param {number} contentWidth
 * @param {number} gapPx
 * @param {number[]|undefined} fractions
 * @returns {number[]}
 */
export function computeGridColumnWidths(columns, contentWidth, gapPx, fractions) {
  const usable = Math.max(1, contentWidth - gapPx * Math.max(0, columns - 1));
  if (Array.isArray(fractions) && fractions.length === columns) {
    const sum = fractions.reduce((a, b) => a + b, 0) || 1;
    return fractions.map((f) => Math.round((usable * f) / sum));
  }
  const each = Math.floor(usable / columns);
  const widths = Array.from({ length: columns }, () => each);
  let remainder = usable - each * columns;
  for (let i = 0; widths[i] !== undefined && remainder > 0; i += 1) {
    widths[i] += 1;
    remainder -= 1;
  }
  return widths;
}

/**
 * @param {number} col
 * @param {number} colSpan
 * @param {number[]} colWidths
 * @param {number} gapPx
 */
export function cellWidthForColumnSpan(col, colSpan, colWidths, gapPx) {
  let w = 0;
  for (let c = col; c < col + colSpan; c += 1) {
    w += colWidths[c] ?? 0;
    if (c > col) w += gapPx;
  }
  return Math.max(1, w);
}

/**
 * @param {PreparedGridCell[]} cells
 * @param {number} rows
 * @param {number} rowGapPx
 * @param {(prepared: PreparedBlock) => number} measureHeight
 * @returns {number[]}
 */
export function computeGridRowHeights(cells, rows, rowGapPx, measureHeight) {
  /** @type {number[]} */
  const rowHeights = Array.from({ length: rows }, () => 0);

  for (const entry of cells) {
    if (entry.rowSpan === 1) {
      const h = measureHeight(entry.prepared);
      rowHeights[entry.row] = Math.max(rowHeights[entry.row], h);
    }
  }

  for (const entry of cells) {
    if (entry.rowSpan <= 1) continue;
    const span = entry.rowSpan;
    const cellH = measureHeight(entry.prepared);
    const start = entry.row;
    let current = 0;
    for (let r = start; r < start + span; r += 1) {
      current += rowHeights[r] ?? 0;
      if (r > start) current += rowGapPx;
    }
    if (cellH > current) {
      const extra = cellH - current;
      const perRow = extra / span;
      for (let r = start; r < start + span; r += 1) {
        rowHeights[r] += perRow;
      }
    }
  }

  return rowHeights;
}

/**
 * @param {number[]} rowHeights
 * @param {number} rowGapPx
 */
export function gridTotalHeight(rowHeights, rowGapPx) {
  if (rowHeights.length === 0) return 0;
  const sum = rowHeights.reduce((a, b) => a + b, 0);
  return sum + rowGapPx * Math.max(0, rowHeights.length - 1);
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {DeckTheme} theme
 * @param {GridBlock} block
 * @param {number} contentWidth
 * @param {(block: TextBlock, maxWidth: number) => PreparedBlock} prepareBlock
 * @returns {PreparedGridSection}
 */
export function prepareGridSection(ctx, theme, block, contentWidth, prepareBlock) {
  const grid = normalizeGridBlock(block);
  const inferred = inferGridDimensions(grid.cells);
  const columns = Math.max(1, grid.columns ?? 0, inferred.columns);
  const rows = Math.max(1, grid.rows ?? 0, inferred.rows);
  const colGapPx = Math.round(resolveWidth(grid.columnGap ?? grid.gap ?? '5%', contentWidth, 0));
  const rowGapPx = grid.rowGap != null
    ? Math.round(resolveWidth(grid.rowGap, theme.size, 0))
    : (grid.gap != null && grid.columnGap == null
      ? Math.round(resolveWidth(grid.gap, theme.size, 0))
      : EM_LINE_GAP_PX);
  const colWidths = computeGridColumnWidths(columns, contentWidth, colGapPx, grid.columnWidths);

  /** @type {PreparedGridCell[]} */
  const cells = grid.cells.map((cell) => {
    const colSpan = cell.colSpan ?? 1;
    const cellWidth = cellWidthForColumnSpan(cell.col, colSpan, colWidths, colGapPx);
    /** @type {TextBlock} */
    const maxLines = cell.maxLines ?? grid.cellMaxLines;
    /** @type {TextBlock} */
    const textBlock = {
      section: 'body',
      text: cell.text,
      emphasis: cell.emphasis,
      color: cell.color,
      font: cell.font,
      weight: cell.weight,
      fontSize: cell.fontSize,
      lineHeight: cell.lineHeight,
      ...(maxLines != null ? { maxLines } : {}),
    };
    return {
      row: cell.row,
      col: cell.col,
      rowSpan: cell.rowSpan ?? 1,
      colSpan,
      prepared: prepareBlock(textBlock, cellWidth),
    };
  });

  const rowHeights = computeGridRowHeights(cells, rows, rowGapPx, measureGridCellInkExtent);

  /** @type {PreparedGridSection} */
  const section = {
    section: 'grid',
    block: grid,
    columns,
    rows,
    colGapPx,
    rowGapPx,
    colWidths,
    rowHeights,
    totalWidth: contentWidth,
    totalHeight: 0,
    cells,
  };
  refreshPreparedGridMetrics(section);
  return section;
}

/**
 * @param {PreparedGridSection} grid
 */
export function measurePreparedGridSection(grid) {
  return grid.totalHeight;
}

/**
 * @param {PreparedGridSection} grid
 */
function resolveGridCellAlign(grid) {
  return grid.block.cellAlign ?? (grid.rows > 1 ? 'top' : 'center');
}

/**
 * Bottom extent of all grid cells (matches drawGridSection placement).
 * @param {PreparedGridSection} grid
 */
export function measureGridSectionDrawHeight(grid) {
  const cellAlign = resolveGridCellAlign(grid);
  let maxBottom = 0;
  for (const entry of grid.cells) {
    const box = cellBox(grid, entry);
    const ink = measureGridCellInkExtent(entry.prepared);
    let drawY = box.top;
    if (cellAlign === 'center') {
      drawY += Math.max(0, (box.height - ink) / 2);
    } else if (cellAlign === 'bottom') {
      drawY += Math.max(0, box.height - ink);
    }
    maxBottom = Math.max(maxBottom, drawY + ink);
  }
  return Math.max(0, maxBottom);
}

/**
 * Recompute row heights and total height after cell font sizes change.
 * @param {PreparedGridSection} grid
 */
export function refreshPreparedGridMetrics(grid) {
  const rowGapPx = Number.isFinite(grid.rowGapPx) ? grid.rowGapPx : EM_LINE_GAP_PX;
  grid.rowGapPx = rowGapPx;
  if (!Number.isFinite(grid.colGapPx)) {
    grid.colGapPx = rowGapPx;
  }
  grid.rowHeights = computeGridRowHeights(
    grid.cells,
    grid.rows,
    rowGapPx,
    measureGridCellInkExtent,
  );
  grid.totalHeight = measureGridSectionDrawHeight(grid);
}

/**
 * @param {PreparedGridSection} grid
 * @param {number} col
 * @param {number} row
 */
function cellOrigin(grid, col, row) {
  let x = 0;
  for (let c = 0; c < col; c += 1) {
    x += grid.colWidths[c] + grid.colGapPx;
  }
  let y = 0;
  for (let r = 0; r < row; r += 1) {
    y += grid.rowHeights[r] + grid.rowGapPx;
  }
  return { x, y };
}

/**
 * @param {PreparedGridSection} grid
 * @param {PreparedGridCell} entry
 */
function cellBox(grid, entry) {
  const origin = cellOrigin(grid, entry.col, entry.row);
  const width = cellWidthForColumnSpan(entry.col, entry.colSpan, grid.colWidths, grid.colGapPx);
  let height = 0;
  for (let r = entry.row; r < entry.row + entry.rowSpan; r += 1) {
    height += grid.rowHeights[r];
    if (r > entry.row) height += grid.rowGapPx;
  }
  return { left: origin.x, top: origin.y, width, height };
}

/**
 * @param {CanvasRenderingContext2D} ctx
 * @param {DeckTheme} theme
 * @param {object} variant
 * @param {PreparedGridSection} grid
 * @param {{ columnLeft: number, contentWidth: number, columnRight?: number, canvasSize: number, marginHorizontal: number }} column
 * @param {number} top
 * @param {{
 *   drawPreparedBlock: Function,
 *   measureBlockInkHeight: Function,
 *   sectionAlignmentFor: Function,
 * }} deps
 */
export function drawGridSection(ctx, theme, variant, grid, column, top, deps) {
  const { drawPreparedBlock, measureBlockInkHeight, sectionAlignmentFor } = deps;
  const cellAlign = resolveGridCellAlign(grid);
  const gridLeft = column.columnLeft + Math.round((column.contentWidth - grid.totalWidth) / 2);

  for (const entry of grid.cells) {
    const box = cellBox(grid, entry);
    const inkExtent = measureGridCellInkExtent(entry.prepared);
    const inkW = entry.prepared.width ?? column.contentWidth;
    const boxLeft = gridLeft + box.left;
    const boxTop = top + box.top;

    let drawY = boxTop;
    const stackH = inkExtent;
    if (cellAlign === 'center') {
      drawY = boxTop + Math.max(0, (box.height - stackH) / 2);
    } else if (cellAlign === 'bottom') {
      drawY = boxTop + Math.max(0, box.height - stackH);
    }

    const subColumn = {
      columnLeft: boxLeft + Math.max(0, Math.round((box.width - inkW) / 2)),
      contentWidth: Math.min(box.width, Math.max(inkW, 1)),
      columnRight: boxLeft + box.width,
      canvasSize: column.canvasSize,
      marginHorizontal: 0,
    };

    const alignment = sectionAlignmentFor(variant, 'body');
    drawPreparedBlock(ctx, theme, entry.prepared, alignment, drawY, subColumn);
  }

  return top + measureGridSectionDrawHeight(grid);
}
