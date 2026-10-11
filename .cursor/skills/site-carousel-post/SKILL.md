---
name: site-carousel-post
description: >-
  Generates editorial social-media carousel decks for Hugo post bundles: writes
  carousel.json (narrative arc, slide copy, hierarchy, variants) and carousel.preview (HTML preview stub)
  next to index.md, using the shared renderer at static/carousel/. Requires swipe-order arc and studio
  visual check (claim_proof footers, verticalAlign). Use when the user asks for a carousel, Instagram
  slides, LinkedIn carousel, 1:1 slide deck, or carousel export from a site post.
---

# Carousel from Hugo post

## Goal

Turn a Hugo post into a **coherent 1:1 slide deck** with preview and PNG export.

- **Data** lives next to the post: `carousel.json`, `carousel.preview`
- **Renderer** is shared: `static/carousel/` (served at `/carousel/` under Hugo)
- Preview via **`make serve`**, then open the bundle's `carousel.preview` URL

## When to use

- User asks for carousel, slide deck, Instagram/LinkedIn slides, or 1:1 export from a post
- User provides source notes, essay, or points at an existing `index.md` bundle

## Output files (in the post bundle)

| File | Role |
|------|------|
| `carousel.json` | Deck spec (slides, variants, theme). See **`reference.md`**. |
| `carousel.preview` | Thin HTML stub; copy from **`.cursor/skills/site-carousel-post/templates/carousel.preview`**. **`make serve`**: Hugo serves it as HTML. **`make build`**: renamed to **`carousel.preview.html`** in `public/` for GitHub Pages MIME (same URL without `.html`). |

Do **not** duplicate renderer JS into the bundle.

## Workflow

1. **Read source**: English `index.md` (and body). Optionally `index.es.md` for a separate `carousel.es.json` only if the user asks for Spanish.
2. **Plan the deck** using the editorial rules in **Editorial system** below.
3. **Write `carousel.json`**:
   - Set `title`, `slug` (kebab-case from bundle folder), `source` path
   - Keep one `deck` theme for the full series
   - Set **`"backgroundGradient": "solid"`** unless the author asks for a gradient wash
   - Usually **6 to 9 slides**; **7** is a good default arc
   - Provide **2 variants** (`a`, `b`) per slide; add `c` only on hook or closing if useful
   - Variants share the **same slide message**; only typography, line breaks, archetype, alignment, and emphasis change (see **Variants** below)
4. **Write `carousel.preview`** from **`.cursor/skills/site-carousel-post/templates/carousel.preview`** unchanged unless deck JSON filename differs.
5. **Preview in studio** (see **Studio visual check**): open the bundle `carousel.preview` with **`make serve`**, scan every slide thumbnail at 1080, fix layout before handoff.
6. **Tell the user** the preview URL: `/<section>/<bundle-folder>/carousel.preview` (run `make serve`).

## Editorial system

Apply these rules when writing slide copy and hierarchy.

### Primary goal

Slides must be visually coherent, mobile-readable, emotionally punchy, intellectually clear, and varied without looking random.

### Design philosophy

- Deck is a **system first**, variation second
- Clarity beats decoration; hierarchy beats novelty
- Emphasis guides reading; do not interrupt it
- Surprising at first glance, obvious after one second

### Global style (encode in `deck` + variants)

- **2 fonts only**: one display (headlines, punch phrases), one body (support, labels)
- **1 background system** per deck; default **`"backgroundGradient": "solid"`** (flat fill, no wash)
- Palette: 1 background, 1 text, 1 muted, up to 2 accents
- Consistent margin and left alignment unless center clearly wins (often closing slide)

### Header zone (opener only)

The renderer treats **`header`** as a **separate top band**, not part of the main **`body`** cluster. In the studio, header and body have different placement controls for a reason.

| `header` is for | `header` is not for |
|-----------------|---------------------|
| Mouth-opener: a **question**, gap, or tension line the body then answers | The slide's main claim (belongs in **`body`** punch) |
| Curiosity pull: `What are the game mechanics?` | Neutral **`##` section titles** copied from the post (`Inside the simulation`, `The game mechanics`) |
| A turn that reframes before the punch | Continuity from the prior slide (use a muted first **`body`** line, e.g. `So you keep playing:`) |
| Optional; omit when the punch carries the slide alone | Essay prose, proof, or a second idea competing with the punch |

**Read order:** opener (`header`) → argument (`body` cluster) → optional one-line sting (`footer`). The opener must not repeat the punch in different words.

**Omit `header` on most slides.** Use it sparingly (often one mechanism or turn slide per deck). Slide 1 hook and closing slides usually keep everything in **`body`** or **`footer`**.

**When the opener and punch must read as one thought** (thesis, closing): put the opener as the first muted **`body`** line and use **`verticalAlign: "center"`** on the variant. Do not use **`header`** if preview shows a large void between the top band and the body cluster (same failure mode as `claim_proof` + center).

### Typography on slides

- One core idea per slide; **two lines in the body cluster** beats four zones of micro-copy
- **Block budget:** default **2 `body` blocks** (punch + muted support). Add **`header`** only for an opener (see above). Add **`footer`** only for a one-line sting or rare bridge; MUST NOT stack opener + two body ideas + footer bridge on one slide unless preview proves it still reads clean at 1080
- **`keyword_anchor`:** only when the sentence still reads aloud in order; MUST NOT fragment grammar across three tiny blocks (e.g. "Brain builds" / "**me**" / "tracks approval…")
- **`\n` line breaks:** break at phrase boundaries (comma, clause), not before a single trailing word (e.g. avoid `Just to avoid\nlosing`). Preview the variant; orphan words look like a layout bug.
- **Two-line argument slides (2–7):** use **`verticalAlign: "center"`** and **`alignment: { "body": "center" }`** so the cluster sits in the frame; left-only body on a sparse slide leaves a dead right half.
- **Section sizes**: set `deck.fontSizes` (`header`, `body`, `footer`) and `emphasisScale.punch` for accent lines
- **Explicit sizes**: optional per-block `fontSize` (px at 1080 canvas) when deck defaults are too coarse for a variant
- Emphasize only **2 to 5 words or phrases** per slide
- Prefer less text and stronger hierarchy over shrinking to fit
- If crowded, split into two slides

### Layout archetypes (rotate across deck)

Map to `archetype` in JSON:

1. **stacked_rhythm** (`stacked_rhythm`) — default when a slide needs **two or more** support lines plus a punch
2. **hero_punch** (`hero_punch`)
3. **claim_proof** (`claim_proof`) — only when punch is **short** and footer is **one quiet line** (see **Studio visual check**)
4. **keyword_anchor** (`keyword_anchor`)
5. **closing_thesis** (`closing_thesis`)
6. **post_cta** (`post_cta`) — logo, bundle featured image, post title, display URL (last slide)

### Studio visual check (MUST)

JSON validates; **layout is judged in the studio**, not in the editor. After writing `carousel.json`, open **`carousel.preview`** and scan **every** slide at 1080 (thumbnails are enough for a first pass; open variants that look cramped).

| Issue | Cause | Fix |
|-------|--------|-----|
| Huge empty band between header and footer | `claim_proof` + `verticalAlign: "center"` | Drop `verticalAlign` (default top) or switch to **`stacked_rhythm`** |
| Footer reads like essay prose | Two sentences or 15+ words in **`footer`** | **One** footer line (~12 words max); move other ideas to **`body`** blocks or the **next slide** |
| Punch feels marooned in the middle | Long punch in `claim_proof` with tiny header/footer | Shorten punch, use `\n` for one break, or use **`stacked_rhythm`** |
| Slide crowded top to bottom | Four body blocks + header + footer | Split into two slides (keeps arc) or cut a line |

**`claim_proof` rules**

- **`header` (optional):** opener question or tension only (see **Header zone**), not a topic label.
- **Body punch:** one claim line; prefer ≤8 words, or one intentional `\n` (two lines max).
- **Footer:** **one** payoff or bridge line. MUST NOT be two sentences.
- **Avoid** `verticalAlign: "center"` on `claim_proof` unless you have checked the preview and want deliberate vertical centering (rare).

**`stacked_rhythm` rules**

- Put setup lines in **`body`** (`normal` then `punch`). Reserve **`footer`** for a single bridge or sting.
- **Hook slides (slide 1):** keep the full hook in the **`body`** cluster (including the sting line). Avoid a lone **`footer`** pinned to the bottom while the punch sits at the top; use **`verticalAlign: "center"`** on the variant so the cluster sits in the frame.
- **Thesis / mechanism slides (2+):** default top alignment; punch + one muted support in **`body`**. **`header`** only when an opener question earns its own band (e.g. game mechanics).
- **Example / scene slides (after hijack):** MUST NOT be a centered two-liner with no context. Open with a muted **`body`** bridge (e.g. `So you keep playing:`), then the example punch. Do not put that bridge in **`header`**.
- Both variants of the same slide number should use the **same facts**; only archetype, breaks, and accent placement differ.

### Pacing arc (narrative, not topics)

The deck must read as **one argument in order** when swiped start to finish. Each slide is the **next causal beat**, not a themed poster or a bullet from the essay.

**Swipe test (MUST before finishing):** Read slide numbers 1 through N aloud. Each line must answer "because of the previous slide" or "so what happens next?" If a slide could be swapped with slide 4 or 6 without breaking the story, reorder or rewrite.

Default arc:

1. **Hook** (question or tension)
2. **Thesis** (what the thing is)
3. **Mechanism / priority** (how it works or what comes first)
4. **Implication** (cost, sting, or trade-off that follows the mechanism)
5. **Hijack or mindset** (why you keep obeying, or what the system demands)
6. **Examples** (ground slide 5; muted **`body`** bridge, e.g. `So you keep playing:`)
7. **Closing thesis** (lie revealed + exit, or forward look)
8. **Post CTA** (`post_cta`, optional): featured image, punch title, short URL + scan label, on-slide QR, brand lockup. Set shared assets in `deck.cta` (see **Post CTA short link** below). Put the full permalink in the social caption.

**Bridge the chain**

- **Swipe continuity:** use a muted first **`body`** line when the slide continues the prior beat (e.g. `So you keep playing:`). That stays inside the main cluster, not in **`header`**.
- **Opener tension:** use **`header`** only when the slide needs a mouth-opener the **`body`** then answers (see **Header zone**). Rare; not on every slide.
- **`footer`** on slide N may tee up slide N+1 in **one short line** only when preview shows no top/body/footer void; prefer a muted opening **`body`** line on the next slide instead.
- When the post has **`##` sections**, follow that **section order** in the deck. Do not front-load the glitch, awakening, or examples before the reader knows what the game is.

**Arc anti-patterns (MUST NOT)**

- Topic pile: mechanism, then unrelated topic, then closing metaphor, with no "so / therefore / so you".
- Re-stating the thesis on slide 2 after the hook already said it, without adding a new fact.
- Putting **examples** before **why the system hijacks you**.
- Variant `a` vs `b` that close with **different arguments** (e.g. Nash on `a`, "not real" on `b`). Same beat, different typography only.

When the author supplies a numbered slide list with roles (Title, Center, Priority, …), **map each item to one slide in that order**. Do not shuffle or front-load abstraction.

### Variants (typography only)

Variants are **layout options for the same slide**, not alternate copy.

| Same across variants | May differ across variants |
|----------------------|----------------------------|
| Slide role and core message | `archetype` (stacked_rhythm vs claim_proof vs hero_punch) |
| Facts, claims, examples | Line breaks, accent placement, emphasis split |
| On-slide meaning | `alignment`, block `section` grouping |

**Do not** write variant `b` as a paraphrase, shorter summary, or different argument. Pick one on-slide message per slide number, then compose 2 (or 3) typographic treatments of that message.

Example (slide 3 Priority): both variants include family/tribe and money-as-tool; variant `a` uses claim_proof (header + punch + footer), variant `b` uses stacked_rhythm with a line break in the punch.

- Middle slides (2–6) should stay on the **same narrative beat** across variants; only the visual treatment changes.

### Voice

- Plain English, direct, conversational
- No corporate jargon or generic inspiration
- Preserve sharp metaphors from the source post
- Each slide adds a new fact, mechanism, example, trade-off, or framing line

### Variation budget (per slide)

- At most 1 oversized phrase cluster
- At most 2 accent-colored phrases
- At most 1 unusual placement move

## Post CTA short link

When the last slide funnels to the Hugo post:

1. MUST read **`data/short-link-register.txt`** and follow **`.cursor/skills/site-short-link-register/SKILL.md`** before minting a new path.
2. Register the path in **`data/short-links.toml`** and add a matching Hugo **`aliases`** entry on the English **`index.md`** (e.g. `/war-game/`).
3. In **`deck.cta`**:
   - **`shortUrl`**: full URL with `https://` for the QR (e.g. `https://blog.dualsubstrate.com/war-game/`)
   - **`postUrl`**: canonical long permalink (traceability; not printed on-slide when QR is used)
4. On-slide stack (top to bottom): post title (`body` punch), short URL (`body` line directly under title), QR, scan label (`footer`, e.g. "Scan to read essay"). MUST NOT print the long essay URL.
5. Renderer draws a QR from **`shortUrl`** (falls back to `postUrl` if `shortUrl` is omitted).
6. In carousel studio, tune **`featuredMaxHeight`**, **`qr.size`**, and **`brandMaxHeight`** in the floating **CTA slide (1080px)** panel (writes `deck.cta`; copy JSON into `carousel.json`).

## JSON authoring

Follow **`reference.md`** for field names and defaults.

Default theme (matches site editorial dark deck):

```json
{
  "displayFont": "Playfair Display",
  "bodyFont": "Source Sans 3",
  "palette": {
    "background": "#1a1e26",
    "text": "#f5f5f0",
    "muted": "#9aa3ad",
    "accent1": "#d69a80",
    "accent2": "#e4b896"
  },
  "backgroundGradient": "solid",
  "fontSizes": { "header": 60, "body": 80, "footer": 52 },
  "emphasisScale": { "normal": 1, "punch": 1.25 },
  "marginHorizontal": "10%",
  "marginVertical": "10%",
  "size": 1080
}
```

### Variant pattern

Each slide:

```json
{
  "number": 1,
  "role": "hook",
  "variants": [
    { "archetype": "stacked_rhythm", "alignment": "left", "blocks": [] },
    { "archetype": "hero_punch", "alignment": "left", "blocks": [] }
  ]
}
```

Variant labels for export (`…-slide-01-a.webp`, `…-b.webp`) come from **array order** (0 → `a`, 1 → `b`).

Use `\n` in `text` for intentional line breaks.

## Quality check before finishing

- **Swipe test** passes (see **Pacing arc**)
- **Studio visual check** done on every slide number and variant
- Deck feels like one designer made it
- Hierarchy obvious at phone size; no paragraph footers
- Accent color restrained
- No repeated ideas across slides (rephrase ≠ new beat)
- Hook slide works in under 2 seconds
- JSON validates (no trailing commas, valid UTF-8)

## Preview and export

1. Run **`make serve`** (also starts the local save API so studio **Save** can write `carousel.json`)
2. Open **`http://localhost:1313/<section>/<bundle>/carousel.preview`**
3. Tune layout in studio, then click **Save** in **Settings: file vs browser**
4. Set slide size to **1080**, pick winning variants with **In strip**, then click **PDF** (one `{slug}-linkedin.pdf` from this Chrome session)
5. Optional WebPs: click **Slides**, then **`make carousel-pdf DIR=$$HOME/Downloads SLUG=<deck-slug>`** if you already have those files. Optional **`OUT=`**, **`VARIANT=a`**.

Optional: save chosen WebPs to `carousel-slides/` in the bundle. MUST NOT commit exported WebPs or PDFs unless the author asks.

## Typo layout guardrails (renderer)

When editing `static/carousel/renderer.js`, `inline-text.js`, or line-box debug:

- **Single source of probe alphabets:** `static/carousel/typo-probes.js` (loaded at import; throws if sets overlap).
- **Never mix probe roles:** x-height (`acemnorsuvxz`), ascender (`bdfhklt`), descender (`gjpqy`) are measured by separate functions only.
- **Debug colors:** blue = line slot (orange band × `lineHeights`); orange = ascender + descender probes from **layout baseline** (same Y as fillText); red = x-height band + layout baseline.
- **`lineHeights`:** `1` = blue hugs orange; `2` = double orange band height; slots stack edge to edge. Keys: `header`, `footer`, `normal`, `punch` (body).
- **Punch blocks:** `emphasis: punch` → display font, weight `700`, metrics via `fontWeightProbes()`; orange band from punch font size (includes `emphasisScale`).
- **Do not** use `fontBoundingBoxAscent` for orange or red (collapses onto blue em top).
- **Do not** use ascender letters (`bdfhklt`) when measuring x-height band (stretches red to em edges).
- After metric changes, hard-refresh studio and check Line boxes on slide 3 variant **a**.

## Related

- Renderer: `static/carousel/`
- Schema: **`reference.md`**
- HTML stub: **`.cursor/skills/site-carousel-post/templates/carousel.preview`**
- Short URLs: **`.cursor/skills/site-short-link-register/SKILL.md`**, `data/short-links.toml`
- LinkedIn **document caption** (Gemma 4 → `linkedin-carousel.txt`): **`.cursor/skills/site-carousel-linkedin-caption/SKILL.md`**
- Full LinkedIn essay (`linkedin.txt`): **`.cursor/skills/site-linkedin-post/SKILL.md`**
- Example deck (arc + layout): `content/social-protocols/2026-05-26-developed-countries-are-factories/carousel.json`
- Example deck (arc rewrite): `content/human-condition/2026-05-01-ego-as-game/carousel.json`

## Do not

- MUST NOT commit exported WebPs or LinkedIn PDFs unless the user asks
- MUST NOT duplicate renderer code into post bundles
- MUST NOT use em dashes (U+2014) in slide copy
- MUST NOT publish carousel files as Hugo pages; they are bundle sidecars only
- MUST NOT hand off a deck without opening **`carousel.preview`** at least once
- MUST NOT put two sentences in a **`footer`** on argument slides
- MUST NOT use `claim_proof` with a long footer to carry ideas that belong in **`body`** blocks or the next slide
- MUST NOT default `verticalAlign: "center"` on `claim_proof` slides
