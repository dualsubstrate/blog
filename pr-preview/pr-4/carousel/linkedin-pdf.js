/**
 * @typedef {Object} PdfPage
 * @property {number} width
 * @property {number} height
 * @property {Uint8Array} flate
 */

/**
 * Flatten canvas pixels to RGB (alpha blended on black).
 * @param {HTMLCanvasElement} canvas
 * @returns {Uint8Array}
 */
export function canvasToRgbBytes(canvas) {
  const ctx = canvas.getContext('2d');
  if (!ctx) {
    throw new Error('Canvas 2D context unavailable');
  }
  const { width, height } = canvas;
  const src = ctx.getImageData(0, 0, width, height).data;
  const rgb = new Uint8Array(width * height * 3);
  let o = 0;
  for (let i = 0; i < src.length; i += 4) {
    const a = src[i + 3] / 255;
    rgb[o] = Math.round(src[i] * a);
    rgb[o + 1] = Math.round(src[i + 1] * a);
    rgb[o + 2] = Math.round(src[i + 2] * a);
    o += 3;
  }
  return rgb;
}

/**
 * zlib-wrap raw bytes (PDF FlateDecode).
 * @param {Uint8Array} bytes
 * @returns {Promise<Uint8Array>}
 */
export async function flateBytes(bytes) {
  if (typeof CompressionStream !== 'function') {
    throw new Error('Lossless PDF needs CompressionStream (use current Chrome)');
  }
  const stream = new Blob([bytes]).stream().pipeThrough(new CompressionStream('deflate'));
  const buffer = await new Response(stream).arrayBuffer();
  return new Uint8Array(buffer);
}

/**
 * @param {HTMLCanvasElement} canvas
 * @returns {Promise<PdfPage>}
 */
export async function encodeCanvasPage(canvas) {
  const width = canvas.width;
  const height = canvas.height;
  if (width < 1 || height < 1) {
    throw new Error('canvas has empty dimensions');
  }
  const flate = await flateBytes(canvasToRgbBytes(canvas));
  return { width, height, flate };
}

/**
 * Full-bleed LinkedIn PDF: one lossless RGB image per page, MediaBox = pixel size.
 * @param {PdfPage[]} pages
 * @returns {Uint8Array}
 */
export function assembleLinkedInPdf(pages) {
  if (!Array.isArray(pages) || pages.length === 0) {
    throw new Error('no pages');
  }
  const b = new PdfBuf();
  b.writeString('%PDF-1.4\n');
  b.writeBytes(new Uint8Array([0x25, 0xff, 0xff, 0xff, 0xff, 0x0a]));

  const catalogID = 1;
  const pagesID = 2;
  const pageIDs = [];
  const contentIDs = [];
  const imageIDs = [];
  let nextID = 3;
  for (let i = 0; i < pages.length; i += 1) {
    contentIDs[i] = nextID;
    imageIDs[i] = nextID + 1;
    pageIDs[i] = nextID + 2;
    nextID += 3;
  }

  b.startObj(catalogID);
  b.writeString(`<< /Type /Catalog /Pages ${pagesID} 0 R >>\n`);
  b.endObj();

  b.startObj(pagesID);
  b.writeString(`<< /Type /Pages /Count ${pages.length} /Kids [`);
  b.writeString(pageIDs.map((id) => `${id} 0 R`).join(' '));
  b.writeString('] >>\n');
  b.endObj();

  for (let i = 0; i < pages.length; i += 1) {
    const page = pages[i];
    const content = `q ${page.width} 0 0 ${page.height} 0 0 cm /Im1 Do Q\n`;
    const contentBytes = utf8Bytes(content);
    b.startObj(contentIDs[i]);
    b.writeString(`<< /Length ${contentBytes.length} >>\nstream\n`);
    b.writeBytes(contentBytes);
    b.writeString('endstream\n');
    b.endObj();

    b.startObj(imageIDs[i]);
    b.writeString(
      `<< /Type /XObject /Subtype /Image /Width ${page.width} /Height ${page.height} /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /FlateDecode /Length ${page.flate.length} >>\nstream\n`,
    );
    b.writeBytes(page.flate);
    b.writeString('endstream\n');
    b.endObj();

    b.startObj(pageIDs[i]);
    b.writeString(
      `<< /Type /Page /Parent ${pagesID} 0 R /MediaBox [0 0 ${page.width} ${page.height}] /Resources << /XObject << /Im1 ${imageIDs[i]} 0 R >> >> /Contents ${contentIDs[i]} 0 R >>\n`,
    );
    b.endObj();
  }

  const startxref = b.length;
  b.writeString(`xref\n0 ${nextID}\n`);
  b.writeString('0000000000 65535 f \n');
  for (let id = 1; id < nextID; id += 1) {
    const off = b.offset.get(id);
    if (off == null) {
      throw new Error(`missing pdf object ${id}`);
    }
    b.writeString(`${String(off).padStart(10, '0')} 00000 n \n`);
  }
  b.writeString(`trailer << /Size ${nextID} /Root ${catalogID} 0 R >>\nstartxref\n${startxref}\n%%EOF\n`);
  return b.toUint8Array();
}

/**
 * @param {HTMLCanvasElement[]} canvases
 * @returns {Promise<Uint8Array>}
 */
export async function assembleLinkedInPdfFromCanvases(canvases) {
  if (!Array.isArray(canvases) || canvases.length === 0) {
    throw new Error('no slides to assemble');
  }
  /** @type {PdfPage[]} */
  const pages = [];
  for (const canvas of canvases) {
    pages.push(await encodeCanvasPage(canvas));
  }
  return assembleLinkedInPdf(pages);
}

const encoder = new TextEncoder();

/** @param {string} s */
function utf8Bytes(s) {
  return encoder.encode(s);
}

class PdfBuf {
  constructor() {
    /** @type {Uint8Array[]} */
    this.chunks = [];
    this.length = 0;
    /** @type {Map<number, number>} */
    this.offset = new Map();
  }

  /** @param {Uint8Array} bytes */
  writeBytes(bytes) {
    this.chunks.push(bytes);
    this.length += bytes.length;
  }

  /** @param {string} s */
  writeString(s) {
    this.writeBytes(utf8Bytes(s));
  }

  /** @param {number} id */
  startObj(id) {
    this.offset.set(id, this.length);
    this.writeString(`${id} 0 obj\n`);
  }

  endObj() {
    this.writeString('endobj\n');
  }

  /** @returns {Uint8Array} */
  toUint8Array() {
    const out = new Uint8Array(this.length);
    let pos = 0;
    for (const chunk of this.chunks) {
      out.set(chunk, pos);
      pos += chunk.length;
    }
    return out;
  }
}
