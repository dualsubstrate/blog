# Carousel deck JSON schema

Version **1** deck file consumed by `static/carousel/studio.js`.

## Top level

| Field | Required | Notes |
|-------|----------|-------|
| `version` | yes | Must be `1` |
| `title` | yes | Deck title (preview header) |
| `slug` | yes | Used in export filenames (kebab-case) |
| `source` | optional | Path to Hugo post for traceability |
| `deck` | optional | Theme overrides (see below) |
| `slides` | yes | Array of slide objects |

## Theme (`deck`)

| Field | Default |
|-------|---------|
| `displayFont` | `Playfair Display` |
| `bodyFont` | `Source Sans 3` |
| `palette` | Base colors object (see below) |
| `backgroundWave` | Panoramic color field across all slides (each slide is a slice). `style`: `none`, `drift`, or `mesh-corners`; `lobes`, `intensity` (0–0.72, presence), `color` (0–1, accent richness), `variety` (0–1, palette color steps), `blur`, optional `phase` |
| `backgroundGradient` | Legacy; keep **`solid`**. Per-slide gradient presets are no longer used for carousel frames. |
| `margin` | Shorthand for both axes: `%` or px resolved on the **1080 design canvas**, then scaled to `size` (same as fonts). Default ~10.4% |
| `marginHorizontal` | Left/right inset; default ~10.4% when omitted. `"0%"` / `0` = full bleed. Aliases: `marginX` |
| `marginVertical` | Top/bottom inset; default **10%** when omitted. Aliases: `marginY` |
| `size` | Export/render canvas px (clamped 600–1080); margins and content width scale from 1080 |
| `contentMaxWidth` | `%` of inner width at 1080 design, then scaled; default `88%` if omitted |
| `previewMaxPx` | `300` (studio preview thumbnail max width; min 100) |

### Palette (`deck.palette`)

Text and slide copy colors. Panoramic wash uses the same base and accents unless `deck.wavePalette` is set.

```json
"palette": {
  "background": "#1a1e26",
  "text": "#f5f5f0",
  "muted": "#9aa3ad",
  "accent1": "#df9311",
  "accent2": "#e77218"
},
"wavePalette": {
  "background": "#eef3eb",
  "muted": "#5c6b62",
  "accent1": "#6d8f5f",
  "accent2": "#4a8fa3"
},
"backgroundGradient": "solid",
"backgroundWave": {
  "style": "mesh-corners",
  "intensity": 0.4,
  "blur": 0.6
}
```

| Key | Role |
|-----|------|
| `background` | Slide base fill when wave palette is linked; CTA flat fill |
| `text` | Default body text |
| `muted` | Header, footer, quiet labels |
| `accent1` | Primary punch accent (text) |
| `accent2` | Second punch accent (optional; falls back to `accent1`) |

### Wave palette (`deck.wavePalette`, optional)

Wash-only colors when unlinked from text palette in studio (**Separate wave palette**). Omit to keep wave tied to `palette` base + accents (default).

| Key | Role |
|-----|------|
| `background` | Wash base fill |
| `muted` | Wash muted pools |
| `accent1` | Primary wash accent |
| `accent2` | Second wash accent (optional) |

### Font sizes (`deck.fontSizes`)

Optional defaults at 1080 canvas (block-level `fontSize` overrides):

```json
"fontSizes": {
  "header": 60,
  "body": 80,
  "footer": 52
},
"emphasisScale": {
  "normal": 1,
  "punch": 1.25
}
```

- **`header`**, **`body`**, **`footer`**: base px at 1080 canvas for each section.
- **`emphasisScale.punch`**: multiplier on **`body`** base for `section: body` + `emphasis: punch` (default `1.25` → 100px when body is 80).
- Block-level **`fontSize`** overrides section base + emphasis scale.

### CTA assets (`deck.cta`)

Optional defaults for **`post_cta`** slides (funnel to the Hugo post):

```json
"cta": {
  "background": "#20111c",
  "backgroundGradient": "solid",
  "featuredImage": "cowboys.jpg",
  "brandImage": "/images/og-logo.webp",
  "shortUrl": "https://blog.dualsubstrate.com/war-game/",
  "postUrl": "https://blog.dualsubstrate.com/social-protocols/your-post-slug/",
  "qr": {
    "size": "85%",
    "layout": "split",
    "color": "accent2",
    "brightness": "+55",
    "light": "transparent"
  }
}
```

| Key | Notes |
|-----|-------|
| `logo` | Site path (`/images/head.svg`) or bundle-relative file |
| `logoColor` | Theme color token or `#hex` (tints SVG `currentColor`) |
| `featuredImage` | Bundle-relative to the post folder (same dir as `carousel.json`) |
| `brandImage` | Full-width lockup at slide bottom (e.g. `/images/og-logo.webp`); when set, top `logo` is skipped |
| `background` | Flat slide fill for `post_cta` (often matches OG lockup background) |
| `backgroundGradient` | Use `solid` for CTA slides |
| `shortUrl` | **QR target** and printed short link; register path in `data/short-links.toml` + Hugo `aliases` (see **short-link-register** skill) |
| `postUrl` | Canonical long permalink; not printed on-slide when QR is used; keep for traceability |
| `featuredMaxHeight` | Max featured image height in px at 1080 canvas (default ~200 when QR + brand lockup). Tunable in studio **CTA slide (1080px)** panel. |
| `qr` | Nested QR section. Variant `qr` keys override `deck.cta.qr` for that slide only. |
| `qr.size` | QR square as **percent** of the vertical slot between the URL cluster and scan footer (`"100%"` or `100`; **100%** = largest square that fits, bounded by column width). Tunable in studio **CTA slide (1080px)** panel. Legacy values above 100 (old px at 1080) are treated as 100%. |
| `qr.color` | Module color. Same tokens as slide text (`accent1`, `accent2`, `muted`, `text`) or `#hex`. To lighten or darken, add sibling **`qr.brightness`** (`"+55"` or `55`, range **-100** to **+100**), or write the color as a text tag: `"<accent2 brightness='+55'>"`. Defaults to **`accent2`**. |
| `qr.light` | Tile behind the modules (`transparent` default so the slide shows through, or a token / `#hex`). Use `"light": "muted"` for a solid cream tile. |
| `qr.layout` | **`split`**: QR on the left; `footer` scan line + `brandImage` lockup stacked on the right (default when QR + brand). **`stack`**: centered QR, full-width brand at bottom. |
| `qr.columnRatio` | Left column width fraction for `split` (default `0.5`). |
| `brandMaxHeight` | Bottom brand lockup max height in px at 1080 canvas (default ~180). Tunable in studio **Logo max H**. |

Variant fields with the same keys override `deck.cta` for that slide only.

### Line heights (`deck.lineHeights`)

Multiplier on the **orange band** (ascender probe line to descender probe line). Sets **blue line slot** height; slots stack edge to edge (plus **2px** between consecutive lines in one block). Block-level `"lineHeight"` overrides.

```json
"lineHeights": {
  "header": 1,
  "normal": 1,
  "punch": 1,
  "footer": 1.2
}
```

Keys: **`header`**, **`footer`**, **`normal`**, **`punch`** (`normal` / `punch` apply to `section: body` only).

- **`1`**: blue top/bottom align with orange ascender and descender (no extra padding inside the slot).
- **`2`**: blue slot height is **2×** the orange band; extra space is split evenly above and below the orange lines.
- **Stacking**: each line’s advance equals its blue slot height; blue boxes touch as before.

Probe alphabets (`static/carousel/typo-probes.js`): `ASCENDER_PROBE_CHARS`, `DESCENDER_PROBE_CHARS` (orange), `X_HEIGHT_PROBE_CHARS` (red meanline), `EM_STRUT_CHARS` (font strut, debug only). Line boxes debug: **blue** = line slot; **orange** = ascender + descender; **red** = x-height + layout baseline.

### Body↔punch gap (`deck.emphasisGap`)

A **textless gap box** at body↔punch seams only. Sized as a **fraction of `fontSizes.body`** (same px gap for body→punch and punch→body):

```json
"fontSizes": { "body": 80 },
"emphasisScale": { "punch": 1.25 },
"emphasisGap": 0.4
```

- **`0`**: no textless seam gap (default for stacks that rely on line height only).
- **`0.4`**: optional extra band at normal↔punch block seams only (fraction of `fontSizes.body`); omit or set `0` when line-height alone should control rhythm.

Scales with body size and canvas export. Block seams otherwise follow the same line-advance rules as `\n` inside one block.

### Panoramic background (`deck.backgroundWave`)

Carousel frames always use one **panoramic color field** across the deck; each slide renders a horizontal slice. Studio panel: **Panoramic wave**.

| `style` | Effect |
|---------|--------|
| `none` | Flat `palette.background` on every slide (no wash) |
| `drift` | Mixed-wave lobes (non-uniform spacing, shimmer layer; not a single repeating ramp) |
| `mesh-corners` | Corner color pools (warm bottom, cool top), dark center vignette |

| Field | Notes |
|-------|-------|
| `intensity` | Wash presence: opacity and coverage (0–0.72; default ~0.32; `0` = flat base) |
| `color` | Accent richness: how far lobes mix toward palette accents (0–1; default ~0.55) |
| `variety` | Extended background palette: blends between base, accent1, accent2, and muted (low = 3 deck colors; high = full smooth ring of in-between steps; default ~0.62) |
| `blur` | Spread factor in base lobe sizing (default ~0.55; max 1.2) |
| `radius` | Lobe size multiplier (default `1`; lower = tighter pools, higher = wider wash) |
| `lobes` | Extra mid-strip pools (`drift`: across deck; `mesh-corners`: bottom travelers) |
| `phase` | Optional phase offset (radians) |

`post_cta` slides stay flat (`palette.background` only). Omit `backgroundWave` for renderer defaults (`drift`).

### Color palettes (studio)

The preview toolbar has two panels:

- **Text palette:** preset cards + five chips (text, accents, base for CTA). Copy palette JSON → `deck.palette`.
- **Background wave:** optional **Separate wave palette** toggle + wash chips (Base, Muted, Accents). When linked, wash follows text palette base and accents. Copy wave JSON → `backgroundGradient`, `backgroundWave`, and `wavePalette` when unlinked.
- **Wave style:** `none`, `drift`, or `mesh-corners` tiles; modifiers (`hue`, `intensity`, `color`, `variety`, `blur`, `radius`, `phase`, optional `lobes`). Modifiers hide when style is `none`.

The fixed **Line boxes** panel (top right) toggles metric overlays and holds **Line heights** fields (`normal`, `punch`, `header`, `footer`) plus **Copy lineHeights JSON**. Use **Save** in **Settings: file vs browser** to write live studio settings back to `carousel.json`.

Studio presets (`static/carousel/palettes.js`): `factory-warm`, `editorial-trio`, `ember-peach`, `slate-cool`, `ocean-depth`, `paper-light`, `sage-light`, `warm-linen` (light bases), `wine-depth` (burgundy base), `forest-green` (dark green base). Wave style is chosen separately.

### Motif strip (`deck.motifStrip`)

Optional **panoramic artwork** sliced across slides: one wide SVG/PNG/WebP in the bundle, equal-width horizontal slices per slide. Drawn **under** text at the bottom (or top when `anchor: "top"`).

```json
"motifStrip": {
  "enabled": true,
  "src": "motifs/pyramids-scribble.webp",
  "bandWidth": "100%",
  "marginBottom": "3%"
}
```

| Field | Default | Notes |
|-------|---------|-------|
| `enabled` | `true` when `src` is set | Set `false` to skip load/draw |
| `src` | required | Bundle-relative strip asset (one slide width per deck slide) |
| `bandWidth` | `100%` | Uniform size of the panoramic motif (width and height locked). Scales from the **bottom-left**. Studio **Size** writes this as a percent of native. The strip stays continuous across slides (not inset per slide). |
| `marginBottom` | `3%` | Gap above slide bottom when `anchor` is `bottom` |
| `marginTop` | `3%` | Gap below slide top when `anchor` is `top` |
| `offsetX` | `0` | Horizontal pan in px at 1080 (positive = right). Studio range is ±12 slide widths. |
| `offsetY` | `0` | Vertical nudge in px at 1080 (positive = down). Studio range is ± one slide. |
| `anchor` | `bottom` | `bottom` or `top` |
| `color` | `accent1` | SVG tint (palette token or `#hex`) |
| `keyColor` / `keyTolerance` | off | Raster chroma-key for black-backdrop exports |
| `opacity` | `1` | 0–1 |
| `excludeRoles` | `[]` | Skip motif on slides with matching `role` |

Export art at **N × 1080** width (N = slide count) so each slice aligns to one slide. The studio strip panel (`#carousel-vision-strip`) previews the deck; use **PDF** for a LinkedIn document, or **Panorama** / **Slides** for WebPs.

Drawn after background, before type. Avoid long footers on the same band if they collide visually.

## Slide

| Field | Required | Notes |
|-------|----------|-------|
| `number` | yes | 1-based slide index |
| `role` | yes | hook, mechanism, evidence, scene, trade-off, closing thesis, etc. |
| `variants` | yes | 1 to 3 **typographic** treatments of the same slide message |

## Variants

Each slide number has one **message beat**. Every entry in `variants` must express that same beat.

- **Variant `a` / `b` / `c`:** different archetype, line breaks, emphasis split, alignment, or accent placement
- **Not allowed:** different claims, alternate summaries, or different examples between variants

Pick the winning variant at export time; the deck arc does not branch.

## Variant

| Field | Required | Notes |
|-------|----------|-------|
| `archetype` | yes | Layout label (stacked_rhythm, hero_punch, claim_proof, keyword_anchor, closing_thesis, post_cta) |
| `alignment` | optional | Per-section horizontal map: `{ "header"?: "left"|"center"|"right", "body"?: …, "footer"?: … }`. Omit keys that default to `left`. Legacy string value applies the same alignment to every section. |
| `verticalAlign` | optional | `top` (default), `center`, or `bottom`. Controls body cluster placement in the middle band. `header` stays top; `footer` stays bottom. |
| `blocks` | yes | Text blocks, top to bottom |

## Text block

| Field | Required | Notes |
|-------|----------|-------|
| `text` | yes | Supports `\n` for manual line breaks. Inline `**bold**`, `*italic*`, `***both***`, and palette color tags (see below). |
| `section` | yes | `header`, `body`, `footer`, or `grid` (see **Grid block** below) |
| `emphasis` | for `body` | `normal` (default body text) or `punch` (accent headline). Sets default color, font, and size scale. |
| `color` | optional | `text`, `muted`, `accent1`, `accent2`, or `#hex` (per-block override) |
| `font` | optional | `display` or `body` |
| `weight` | optional | CSS weight, e.g. `400`, `600`, `700` |
| `fontSize` | optional | Explicit px at **1080 canvas** (overrides section base + emphasis scale). Scales with `deck.size`. |
| `lineHeight` | optional | Multiplier on orange band for blue line slot (`1` = blue hugs orange) |
| `maxWidth` | optional | Narrow one block: `"80%"` of deck content column, or canvas px |

### Section defaults

| `section` | Default color | Default font | Base size key |
|-----------|---------------|--------------|---------------|
| `header` | `muted` | body, weight 600 | `fontSizes.header` |
| `body` + `normal` | `text` | body, weight 400 | `fontSizes.body` |
| `body` + `punch` | `accent1` | display, weight 700 | `fontSizes.body × emphasisScale.punch` |
| `footer` | `muted` | body, weight 400 | `fontSizes.footer` |

Override `color`, `font`, `weight`, or `fontSize` on any block when needed.

### Inline emphasis in `text`

Inside any block `text` string:

| Syntax | Effect |
|--------|--------|
| `**phrase**` | Bold (weight 700) |
| `*phrase*` | Italic (loads italic font face) |
| `***phrase***` | Bold + italic |
| `<accent1>phrase</accent1>` | Palette color on that span (same size as the block; does not switch to punch scale) |
| `<accent1 brightness="+10">phrase</accent1>` | Lighten that token's base hex by 10% (toward white). Single or double quotes both work. |
| `<accent2 brightness="-20">phrase</accent2>` | Darken by 20% (toward black); range **-100** to **+100** |
| `<color accent2>phrase</color>` | Same as `<accent2>phrase</accent2>` |

Tokens: `accent1`, `accent2`, `text`, `muted`. Combine with bold: `<accent1>**game**</accent1>`. Use **brightness** for tone steps on one accent (e.g. base punch + darker/lighter spans) without switching to `accent2`.

Example (mechanics slide, accent tone steps):

```json
{
  "section": "body",
  "emphasis": "punch",
  "text": "<accent1 brightness=\"+10\">Avatar</accent1>, <accent1 brightness=\"-15\">score,</accent1>\n<accent2 brightness=\"+8\">predicted moves</accent2>"
}
```

Example (one block, two hues, manual line break):

```json
{
  "section": "body",
  "emphasis": "punch",
  "text": "The game is virtual.\n<accent2 brightness=\"+5\">The emotion is not.</accent2>"
}
```

Example footer:

```json
{
  "section": "footer",
  "text": "Money is a tool in service of life, **not the goal**."
}
```

Unclosed `*`, `**`, or color tags stay literal on the slide. No links, code spans, or nested color tags. Prefer separate `blocks` for punch-sized headlines; use inline color when one wrapped line needs two hues.

### Grid block

Use `section: "grid"` in `blocks[]` between normal `body` blocks (opener above, payoff below). Each cell is its own punch-sized line group with the same fields as a `body` block (`emphasis`, `color`, `fontSize`, inline tags).

| Field | Required | Notes |
|-------|----------|-------|
| `section` | yes | `"grid"` |
| `cells` | yes | Array of cell specs (see below) |
| `columns` | optional | Column count; inferred from cells if omitted |
| `rows` | optional | Row count; inferred from cells if omitted |
| `gap` | optional | Legacy column gap when `columnGap` omitted |
| `columnGap` | optional | Horizontal gap between columns (`%` or px at 1080; default `5%`) |
| `rowGap` | optional | Vertical gap between rows (`%` or px); default `2px` (same rhythm as stacked body blocks) |
| `columnWidths` | optional | Fractions summing to ~1 (e.g. `[0.52, 0.48]`) |
| `cellAlign` | optional | `center` (default), `top`, or `bottom` inside each cell box |
| `cellMaxLines` | optional | Max wrapped lines per cell (often `1`; renderer shrinks font to fit) |

**Cell spec:** `row`, `col` (0-based), optional `rowSpan` / `colSpan` (default 1), plus `text` and the same optional styling fields as `body`.

Example (thesis slide: left stack + tall right cell):

```json
{
  "section": "grid",
  "columns": 2,
  "rows": 2,
  "columnGap": "4%",
  "rowGap": "2px",
  "columnWidths": [0.62, 0.38],
  "cellMaxLines": 1,
  "cells": [
    { "row": 0, "col": 0, "emphasis": "punch", "text": "<accent1 brightness=\"+8\">Your brain</accent1>" },
    { "row": 1, "col": 0, "emphasis": "punch", "text": "<accent1 brightness=\"-5\">provides the</accent1>" },
    { "row": 0, "col": 1, "rowSpan": 2, "emphasis": "punch", "text": "<accent1 brightness=\"-12\">world</accent1>" }
  ]
}
```

The renderer measures the grid as one stack item (gap to adjacent `body` blocks uses `emphasisGap`). Grid cells shrink together when the slide overflows the canvas.

## Slide structure

Three zones on most slides:

| Zone | Section | Role |
|------|---------|------|
| Top | `header` | **Opener only** (question, gap, tension). Pinned to top margin; not part of the body cluster. Omit on most slides. Not for `##` topic labels or swipe bridges. |
| Middle | `body` | **Main content:** clustered punch + support. Bridges from the prior slide start here (muted first line). |
| Bottom | `footer` | Quiet payoff or rare one-line bridge. Optional. |

### Body cluster

Each `body` block is one line group (supports `\n`). Stack multiple blocks for normal + punch in the same cluster:

```json
{ "section": "body", "emphasis": "normal", "text": "Family and tribe" },
{ "section": "body", "emphasis": "punch", "text": "outrank productivity." }
```

Per block:

```json
{ "section": "body", "emphasis": "punch", "fontSize": 100, "text": "life at the center." }
```

Variants change **line breaks and emphasis split**, not the slide meaning.

## Layout engine (renderer)

| Zone | Placement |
|------|-----------|
| `header` | Top margin |
| `footer` | Bottom margin |
| `body` blocks | **Top-aligned** by default. Set `verticalAlign: "center"` or `"bottom"` to place the body cluster in the middle band (between header and footer, or in the full content area). |

Archetype names (`labeled_message`, `stacked_rhythm`, etc.) all use this three-zone model.

| Archetype | Typical blocks |
|-----------|----------------|
| `stacked_rhythm` | optional **`header`** opener; body normal + punch; optional **one-line** footer. Prefer when the slide has 2+ support lines. |
| `hero_punch` | single body punch block, optional footer |
| `claim_proof` | optional **`header`** opener; body punch claim (short); footer proof (**one line only**) |
| `keyword_anchor` | body normal setup, body punch keyword, body normal payoff |
| `closing_thesis` | restrained body stack; one accent word or line |
| `post_cta` | Logo + featured image + title/URL text (slide 8 funnel to full post) |

Keep footer to **one short line** on argument slides (~12 words max; **one sentence**). MUST NOT use the footer as a paragraph or to carry the next slide's full argument. **`post_cta`** may use two footer blocks (tagline + URL).

### `claim_proof` and `verticalAlign`

The renderer pins **`header`** to the top margin and **`footer`** to the bottom. Body blocks default to **top-aligned** in the middle band.

- **`claim_proof` + `verticalAlign: "center"`** often leaves a large empty band and makes the punch look floating. Avoid unless you checked the studio preview and want that look.
- If the slide needs more than one support line, use **`stacked_rhythm`** instead of a long footer under `claim_proof`.

Authoring checklist for each `claim_proof` variant: **`header` only if an opener question/tension**; short punch (≤8 words or one `\n`); one-line footer; preview at 1080.

## Export filenames

`{slug}-slide-{NN}-{variant}.webp` at 1080×1080 (WebP, quality ~0.92).

## Preview URL (Hugo)

With `make serve`, open:

`/social-protocols/<bundle-folder>/carousel.preview`

Example:

`/social-protocols/2026-05-26-developed-countries-are-factories/carousel.preview`

Author **`carousel.preview`** in the bundle (Hugo does **not** publish arbitrary `.html` sidecars from bundles). **`make build`** renames each output to **`carousel.preview.html`** so GitHub Pages serves **`text/html`** at the same URL path (extensionless **`/.../carousel.preview`** maps to the `.html` file). **`make serve`** serves the sidecar directly as HTML without the rename.

Shared assets load from `/carousel/*` (served from `static/carousel/`). Under `make serve`, Hugo sends `Cache-Control: no-store` for `/carousel/**`; the preview stub also cache-busts JS/CSS/JSON on each load.

### Studio palette panel

The preview toolbar **Palette** grid sets base colors; the **Gradient** grid below it picks the background wash. Copy hex from chips or **Copy palette JSON**. Changes update all slide thumbnails and exports in that session.

### Studio placement controls

Each variant card toolbar: **Top** / **Center** / **Bottom** on the left; **Copy placement** (`verticalAlign` when not top, plus compact `alignment` object) and **Export** on the right. Section **L** / **C** / **R** controls sit in the side rail. Click **Save** to write placement and theme settings to `carousel.json`.

### Studio floating panel (debug)

Fixed top-right panel: **Line boxes** toggle; optional **CTA slide** sizes when the deck has `post_cta`; **Margins** (`marginHorizontal` / `marginVertical` as % of the 1080 design canvas, 0–16%); **Line heights**; each section has **Copy … JSON**. **Save** writes the live browser settings back to `carousel.json`.
