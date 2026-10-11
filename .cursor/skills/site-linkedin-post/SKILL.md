---
name: site-linkedin-post
description: >-
  Writes a standalone LinkedIn post from an English Hugo site post. The post
  delivers the full argument on-platform so readers can like without clicking
  away. LinkedIn posts are ALWAYS in English, even if the user asks for a Spanish
  version (they should use facebook-post for Spanish or English Facebook friends copy). When index.es.md exists, linkedin.txt dual-language site links MUST use
  hugo list all permalinks (Spanish URLs may use translated section segments and
  title slugs). Use when the user wants to share a site post on LinkedIn, asks to
  write a LinkedIn post, or mentions sharing, cross-posting, or social copy for
  LinkedIn. Output must not use the em dash (U+2014). For type video with a site
  Chapter Guide, includes a What's in the talk bullet outline (Label · topic, no
  clock times; no pipe/space tables). For cognitive-memetics: **Street Wisdom**, **Reptilocracy**, **Sm(art)**, **Pawtropolis**, and **Cube-Cows** use fold-first layouts (**🟣** after the hook; **`❓`** / **`AND:`** when the episode has separate beats); **Psych-Fitness-28** uses quoted hook line one plus **`{heading_code}: {project}`** on line two; then series blocks, hashtags, **`🧷`** or **`🖋️`**, **`🔗`** hub.
---

# LinkedIn post from site content

## Goal

The post delivers the complete argument on LinkedIn (thesis, ladder, close),
not a compressed copy of every site section. The site link is a citation, not
required reading. A reader who never clicks should still get the full point and
feel comfortable hitting Like.

**Reptilocracy exception:** Deliver the **same scene** as the site (`tldr` /
`fluff` front matter) plus a short **BUT WHY** block that ends with the **Change.org petition**
link (see Reptilocracy section). Do **not** retell the full article or add essay
depth that belongs on the site. The site holds analytic phrasing; LinkedIn
holds the punch.

## Section label: Context (Hugo field `fluff`)

Templates that use labeled blocks (**Por-Estas-Calles**, **Sm(art) T-Shirt Art**, **Reptilocracy**, **Psych-Fitness-28**) map Hugo front matter **`fluff`** to the LinkedIn label **`➕ Context:`** (never **`➕ FLUFF:`**). The YAML field name stays **`fluff`** on the site; only the social section heading changed.

## Source fields

Read from the English `index.md`:

| Site field | Maps to |
|---|---|
| `title` | Hook material; **line 1** SHOULD repeat the full English **`title`** string (including any leading emoji per **Optional leading emoji in `title`** in `.cursor/rules/site-content-markdown-writing.mdc` for `social-protocols`, `human-condition`, `mind-infrastructure`, `x-minds`). Add a blank line, then the opening hook paragraphs if they are not the same text. For **`type: video`** or plain titles without emoji, still use **line 1** = full **`title`** before the body. **Do not** use this line-one rule for **cognitive-memetics** bundles that use the [two-line series header](#cognitive-memetics-two-line-header-all-series). Strip emoji **only** if the author asks or the fold leaves no room (rare). |
| `description` | The Claim; strong hook candidate |
| `grounding` | Citation line at the close |
| Body `###` sections | The argument and examples to compress |
| Body **`## Chapter Guide`** table | **`type: video` only:** **`📑 What's in the talk`** outline (see [Video chapter guide (TOC)](#video-chapter-guide-toc)); copy **label + topic** per row; **omit clock times** on LinkedIn |
| `tags` | Hashtag line after the body (not counted in body word target; see [Hashtags from site `tags`](#hashtags-from-site-tags)) |

For **`type: sayings`** under **`content/cognitive-memetics/sayings/`**, use **[Sayings / Street Wisdom](#sayings--street-wisdom-cognitive-memetics-type-sayings)** instead of the default structure below. Source mapping:

| Site field | Maps to |
|---|---|
| `heading_code` | After **IN ENGLISH**, on the line under **`🟣`** (e.g. `W21: Street-Wisdom 💬🇻🇪`) |
| `project` | Same line as `heading_code` (canonical **`Street-Wisdom 💬🇻🇪`**) |
| `title` | Spanish (or source) saying, in straight double quotes on **line one** only |
| `tldr` | Body under standalone **`❓`** (no `TLDR:` label; strip `**bold**`) |
| `fluff` | Body under standalone **`AND:`** (not `Context:` / `FLUFF:`; strip Markdown) |
| Body `### In English` | **🔤 IN ENGLISH:** block: one **`• `** line per closest equivalent from the blockquote list; plain text only (strip `*` / `**` / backticks). Omit this block if that section is missing. Do **not** paste later `###` body subsections (essay under the equivalents). |
| `tags` | All entries → hashtag line before the hub link; see [Hashtags from site `tags`](#hashtags-from-site-tags) and **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`** |
| `categories` | Second term picks the **hub** for the closing link (see [English category URL](#english-behaviour-engineering-category-link)) |
| `sayingsProjectAboutP1` + `sayingsProjectAboutP2` (`i18n/en.toml`) | Two paragraphs under **`🟣`** (plain text; strip Markdown) |

For **`type: panel`** under **`content/cognitive-memetics/cows/`** (Tales from the Cube Farm), use **[Cube cows / Tales from the Cube Farm](#cube-cows--tales-from-the-cube-farm-cognitive-memetics-type-panel)**. Source mapping:

| Site field | Maps to |
|---|---|
| `heading_code` | After the teaser, on the line under **`🟣`** (e.g. `W7: Cube-Cows 🐮📈`) |
| `project` | Same line as `heading_code` (canonical **`Cube-Cows 🐮📈`**) |
| `description` | Quoted strip teaser on **line one** only (strip `**bold**` and smart quotes as needed for plain text) |
| `cowsProjectAboutBody` (`i18n/en.toml`) | Three paragraphs under **`🟣`** (plain text; strip Markdown from the i18n string) |
| `tags` | All entries → one hashtag line before the category link; see [Hashtags from site `tags`](#hashtags-from-site-tags) |
| `categories` | **`Cube-Cows`** hub → closing link `/categories/cube-cows/` (see [English category URL](#english-behaviour-engineering-category-link)) |

For **`type: sayings`** under **`content/cognitive-memetics/reptilocracy/`**, use **Reptilocracy** below. Source mapping:

| Site field | Maps to |
|---|---|
| `heading_code` | After **AND:** block, on the line under **`🟣`** (e.g. `W6: Reptilocracy 🦎🏛️`) |
| `project` | Same line as `heading_code` (canonical **`Reptilocracy 🦎🏛️`**) |
| `title` | Episode name in straight double quotes on **line one** only |
| `tldr` | Body after **`❓ `** on the same line (no `TLDR:` label) |
| `fluff` | Body under **`AND:`** (not `Context:` / `FLUFF:`) |
| `tags` | Full hashtag line before site and hub links |
| `categories` | **`Reptilocracy`** hub → **`🔗 Reptilocracy (English) →`** `/categories/reptilocracy/` |
| `reptilocracyProjectAboutBody` (`i18n/en.toml`) | Series lines under **`🟣`** (plain text from i18n; MAY compress for LinkedIn length/tone; strip Markdown) |
| `reptilocracyProjectAboutCtaTitle` / `CtaButton` (`i18n/en.toml`) | Petition lead-in wording (adapt to LinkedIn plain text; keep Change.org URL block below) |

## LinkedIn constraints

- **Plain text only.** No Markdown renders: `**bold**` appears as literal asterisks.
- No headers, no code blocks.
- **Exception:** the closing citations may use hyphen bullets for the dual-language site permalinks (see [Site link format (EN + ES)](#site-link-format-en--es)).
- **Exception:** the closing citations may use a two-line external URL format (label line ending with `→`, then `- https://...`) for example YouTube (see [Citation close format](#citation-close-format)).
- Do not use bullet lists for the main argument.
- Emoji and Unicode symbols (✜, ℹ, →, •) are the only visual tools.
- Line breaks and short paragraphs create all spacing and hierarchy.
- Posts truncate at roughly 210 characters before a "see more" button. The
  opening lines must work as a standalone hook.
- **Social autopost (`linkedin-autopost`):** write full `linkedin.txt` including **`🧷` / `🔗` URL blocks**; do not drop links to save space. Autopost encodes [little text](https://learn.microsoft.com/en-us/linkedin/marketing/community-management/shares/little-text-format) for the API (hashtag templates, reserved-char escapes), logs byte/rune counts, warns above ~800 UTF-8 bytes on image posts, and **verifies** stored `commentary` after publish via `GET /rest/posts/{id}` (site URLs, footer markers, hashtag tag names, length vs truncation; needs read scope). See `docs/social-autopost/README.md`.
- **No em dash (U+2014).** MUST NOT use `—` in generated `linkedin.txt` (body, hooks, or citation label lines). Prefer comma, semicolon, colon, parentheses, or a spaced hyphen (` - `) where a sentence needs a break.

## Hashtags from site `tags`

When the source `index.md` has a non-empty `tags` array:

- **Include every tag** in the same order as front matter. Do not drop tags for brevity.
- **Format:** `#` plus the tag string exactly as stored (for example `ArepaContigo` → `#ArepaContigo`). Space-separated on one line; a second line is fine only if the list is very long.
- **Placement (default posts):** After the close / main argument, **before** the site-link and research lines. Skip the hashtag block only when `tags` is missing or empty.

**Cognitive-memetics sayings (Por-Estas-Calles / Street Wisdom):** Use the [fold-first layout](#sayings--street-wisdom-cognitive-memetics-type-sayings): after the **`🟣`** series block (and optional `⁒` kicker), hashtags, then **`🖋️ Full post (site) →`** (**`- ES:`** first, **`- EN:`** second when bilingual), then **`🔗 Por-Estas-Calles (English) →`** hub URL.

**Cube cows:** After the **`🟣`** series block, hashtags, then **`🧷 Full post (site) →`** with **`- EN:`** first and **`- ES:`** second when **`index.es.md`** exists, blank line, then **`🔗 Cube-Cows (English) →`** hub URL.

**Reptilocracy:** Petition close lives inside the **`🟣`** block (before hashtags). Then hashtags, **`🧷 Full post`**, then **`🔗 Reptilocracy (English) →`** hub.

## Title emoji alignment (site ↔ LinkedIn)

- When `index.md` **`title`** begins with one or two emoji (allowed sections above), **`linkedin.txt` line 1** MUST repeat the **same emoji and order**, then a space, then the **rest of the title without extra emoji**.
- Do **not** add title emoji to LinkedIn if the English `title` has none.

## Structure

```text
[Hook: line 1 often repeats site title with any leading emoji; then ≤ 2 more short lines if needed]

[Blank line]

[Logical ladder: each layer/step as a short block, 2-4 lines each]
[Use a symbol prefix (✜, ℹ, →) to mark distinct layers when useful]

[Blank line]

[Close: punchy, 1-3 short sentences; ideally a reframe of the hook]

[Blank line]

[Video only, when index has ## Chapter Guide:]
[📑 What's in the talk]
[Blank line]
[• Label · chapter title  (one line per site table row; no MM:SS)]
[Blank line]

[#Tag1 #Tag2 ... from front matter `tags`; omit only if `tags` is empty]

[Blank line]

[Site link lines: emoji + short label + arrow]
[- EN: <English permalink>]
[- ES: <Spanish permalink>]
[Blank line]
[Optional external link: label line ending with →]
[- <external URL>]
[Blank line]
[Optional research citation: one line: author/title → URL]
```

## Site link format (EN then ES)

The user wants a ready-to-copy dual-language link block so they do not have to
copy permalinks after publishing.

### Rules

- Always output **both** URLs when the bundle has both `index.md` and `index.es.md`.
- Format exactly as two hyphen-bullet lines: **`- EN:`** first (full English permalink from `hugo list all`), **`- ES:`** second (Spanish permalink). Same order for every bundle, including `type: video` and claims.
- Put **no** blank lines between the site-link label line (the line ending with `→`), `- EN: ...`, and `- ES: ...`.
- Put one blank line after `- ES: ...`, then the next citation line.
- For an external video link (for example YouTube), use a **two-line** format: a label line ending with `→`, then `- https://...` on the next line. Do not put trailing spaces after `→`.
- Use **final permalinks**, not guessed paths. Spanish permalinks may differ by section (see `hugo.toml` language permalinks).

### How to get the permalinks

- Run `hugo list all` and copy the **Permalink** values for:
  - the English page (`index.md`) for the **`- EN:`** line
  - the Spanish page (`index.es.md`) for the **`- ES:`** line
- Alternatively run **`make list`** at the repo root (same output).

#### Spanish URL pitfalls (MUST)

- The **`- EN:`** line uses the **`index.md`** permalink (often matches the English folder slug in the path, but still copy from `hugo list all`).
- The **`- ES:`** line MUST use the **`index.es.md`** permalink from `hugo list all`, **not** a path built from the English folder name under `/es/social-protocols/`. For **`content/social-protocols/`** bundles, Spanish routes use **`/es/protocolos-sociales/<slug>/`** per **`hugo.toml`**; **`<slug>`** is usually title-derived for Spanish.
- Guessed URLs break dual-language blocks. See **`.cursor/skills/site-spanish-translation-content/SKILL.md`** → **Spanish URLs, permalinks, and aliases**.

### Template

```text
🧷 Full post (site) →
- EN: https://blog.dualsubstrate.com/...
- ES: https://blog.dualsubstrate.com/es/...

▶️ Example video (YouTube) →
- https://www.youtube.com/watch?v=VIDEO_ID
```

## Video chapter guide (TOC)

When **`type: video`** and the English **`index.md`** body has **`## Chapter Guide`** with a timestamp table:

### Goal

Give a **topic outline** on LinkedIn (who speaks, what beat). **Not** a jump map: clock times and clickable timestamps belong on the **site** and in **`substack.md`**, not in `linkedin.txt`.

### Rules (MUST)

- **Placement:** after the closing punch, **before** the hashtag line.
- **Heading:** exactly `📑 What's in the talk` (plain text; no Markdown). Do not say "timestamps" in the heading.
- **Rows:** one line per site table row, same order. Format: `• Label · chapter title` (Unicode bullet `•`, middle dot ` · ` between label and title).
- **Omit `MM:SS` (and `t=` links)** on every LinkedIn row. Strip `**` from site labels; keep the scan label (for example `Friston`, `Metaphor`, `AGI`).
- **Title half:** chapter text after the label; **MAY** shorten for plain English; **MUST NOT** drop rows unless the author asks for a short list.
- **Exception:** this bullet list is allowed even though the main argument must not use bullet lists.
- **Not counted** toward the [Length](#length) body word target.
- **Autopost:** keep the full outline plus **`🧷` / `▶️` blocks** even if total bytes exceed the ~800-byte image-post warn, unless the author asks to trim.

### Forbidden on LinkedIn (proportional font)

- Pipe tables (`Time | Speaker | Chapter`)
- Box-drawing separator lines
- Space-padded columns or two-line time/title stacks meant to mimic tables
- Clock times on each bullet (tested: adds noise without click-through)

### Where timestamps stay

| Surface | Chapter data |
|---------|----------------|
| Site `index.md` | Full **`## Chapter Guide`** table with YouTube `t=` links |
| `substack.md` / `substack.es.md` | Same table with timestamps ( **`.cursor/skills/site-substack-post/SKILL.md`** ) |
| `linkedin.txt` | Label · topic bullets only (this section) |

Do **not** paste per-row YouTube URLs in the outline; the closing **`▶️ … (YouTube) →`** block is enough.

### Reference bundle

**`content/mind-infrastructure/2026-05-21-brain-is-not-a-computer/linkedin.txt`** (canonical). Site table source: same bundle **`index.md`** → **`## Chapter Guide`**.

```text
📑 What's in the talk

• Friston · free energy and the spherical cow
• Chirimuuta · why scientists abstract
• Debate · Simplicius vs Ignorantio on simplicity
…
```

See **`.cursor/skills/site-video-content/SKILL.md`** for site **`## Chapter Guide`** authoring.

### Cognitive-memetics two-line header (all series)

When a bundle has **`heading_code`** (or day code) plus **`project`** / series line:

- **Line one:** quoted episode hook (`title`, strip **`description`** teaser, or petition headline) in straight double quotes.
- **Line two:** `{heading_code}: {project}` (for example `W19: Street-Wisdom 💬🇻🇪`, `TS5: Sm(art)`, `W7: Cube-Cows 🐮📈`, `D5: Psych-Fitness-28 🙏`, `P1: Pawtropolis (Under Fire) 🐾`).
- Blank line, then the first labeled block (**✔️ TLDR:**, **❓ BUT WHY:**, etc.).

**MUST** use this order on **Psych-Fitness-28** only. **Por-Estas-Calles / Street Wisdom**, **Reptilocracy**, **Sm(art)**, **Pawtropolis**, and **Cube-Cows** use fold-first layouts (quoted hook on line one only; **`🟣`** + **`{heading_code}: {project}`** after the scene blocks, or after the teaser when there is no separate **`❓`** / **`AND:`**; see each series section). If `heading_code` is missing, put **`{project}`** on line two only; if `project` is missing, line two is `heading_code` + quoted **`title`** from the site (legacy only).

### Sayings / Street Wisdom (`cognitive-memetics`, `type: sayings`)

Use this **fold-first layout** for Venezuelan sayings (**`Por-Estas-Calles`** hub). Optimized for LinkedIn "see more": the saying and meaning lead; the series badge comes after the English equivalents, not at the top.

**Rules**

- **Plain text only:** remove all Markdown (`**bold**`, backticks) from `tldr` and `fluff` when composing.
- **Line one:** quoted **`title`** (the saying) only. **MUST NOT** put **`{heading_code}: {project}`** on line two (unlike other cognitive-memetics series).
- **Meaning block:** `❓ ` immediately before the adapted **`tldr`** on the **same line** (no line break after the emoji). **MUST NOT** use `✔️ TLDR:`.
- **Nuance block:** `AND:` on its own line, then the adapted **`fluff`**. **MUST NOT** use `➕ Context:` or `➕ FLUFF:`.
- **🔤 IN ENGLISH:** then one **`• `** line per equivalent from **`### In English`** in the English `index.md` only (strip Markdown). Omit the block if that section is missing.
- **Series block:** `🟣 ` immediately before **`{heading_code}: Street-Wisdom 💬🇻🇪`** on the **same line** (no line break after the emoji). **`project`** MUST be the canonical **`Street-Wisdom 💬🇻🇪`**. Blank line, then paste **`sayingsProjectAboutP1`** and **`sayingsProjectAboutP2`** from **`i18n/en.toml`** as plain text (two paragraphs; strip Markdown). **MUST NOT** invent or paraphrase series copy. **MUST NOT** use `❓ BUT WHY:` as a label; the ❓ icon is reserved for the meaning block above.
- Optional one-line kicker with trailing `⁒` after the series copy when it fits.
- **Hashtags:** one line with **all** `tags` from front matter.
- **Closing:** **`🖋️ Full post (site) →`** (not **`🧷`** on Street Wisdom posts). When bilingual, **`- ES:`** first, **`- EN:`** second (matches live Por-Estas-Calles posts). Blank line; then **`🔗 Por-Estas-Calles (English) →`** hub category URL.

**Template** (series body: paste from **`sayingsProjectAboutP1`** / **`sayingsProjectAboutP2`** in **`i18n/en.toml`**; do not hardcode it here)

```
"Juntos, pero no revueltos."

❓ [Plain-text tldr: meaning, image, when you use it]

AND:
[Plain-text nuance from site fluff]

🔤 IN ENGLISH:
• [Closest equivalent phrase 1]
• [Closest equivalent phrase 2]

🟣 W21: Street-Wisdom 💬🇻🇪

{sayingsProjectAboutP1 from i18n/en.toml, plain text}

{sayingsProjectAboutP2 from i18n/en.toml, plain text}

[Optional one-liner ending with ⁒]

#StreetWisdom #CulturalStopwatch #TakeBackYourMcDonaldsCulture #ArepaContigo #VenezuelanSayings #DistinctRoles

🖋️ Full post (site) →
- ES: https://blog.dualsubstrate.com/es/cognitive-memetics/sayings/2026-06-01-saying-21/
- EN: https://blog.dualsubstrate.com/cognitive-memetics/sayings/2026-06-01-saying-21/

🔗 Por-Estas-Calles (English) → https://blog.dualsubstrate.com/categories/por-estas-calles/
```

Use the **hub category** URL for the series, not the episode URL. Map from `categories` (second term): **`Por-Estas-Calles`** → `/categories/por-estas-calles/`; **`T-Shirt Art`** → `/categories/t-shirt-art/`; **`Sm-Art`** → `/categories/sm-art/`; **`Cube-Cows`** → `/categories/cube-cows/`; **`Raymond`** → `/categories/raymond/`; **`Reptilocracy`** → `/categories/reptilocracy/`; **`Pawtropolis-Under-Fire`** → `/categories/pawtropolis-under-fire/`. Prefix with **`baseURL`** from **`hugo.toml`**. Confirm with `hugo list all` on `content/categories/<slug>/_index.md` (English row, not `*.es.md`).

**Reference bundle** for structure and tone: `content/cognitive-memetics/sayings/2026-06-01-saying-21/` (W21; matches live LinkedIn layout).

### Sm(art) / T-Shirt Art (`cognitive-memetics`, `type: sayings`, **`content/cognitive-memetics/t-shirt-art/`**)

**Fold-first layout** (same read as [Reptilocracy](#reptilocracy-contentcognitive-memeticsreptilocracy-type-sayings)): quoted episode **`title`** on line one only; meaning and nuance before the series badge.

**Rules**

- **Plain text only:** strip all Markdown (`**bold**`, backticks) from **`tldr`** and **`fluff`**; move bare URLs to optional **`🔖`** lines at the close when needed.
- **Line one:** quoted **`title`** (English episode name). **MUST NOT** put **`{heading_code}: {project}`** on line two.
- **Meaning:** `❓ ` immediately before adapted **`tldr`** on the **same line**. **MUST NOT** use `✔️ TLDR:`.
- **Nuance:** `AND:` on its own line, then adapted **`fluff`**. **MUST NOT** use `➕ Context:` or `➕ FLUFF:`.
- **Series badge:** `🟣 ` immediately before **`{heading_code}: Sm(art)`** on the **same line** (for example `TS5: Sm(art)`). **Omit** a series essay block under **`🟣`** (no Por-Estas-Calles / Reptilocracy framing on these bundles).
- **Omit** **`🔤 IN ENGLISH:`** when the English `index.md` has no **`### In English`** block. **Omit** **`❓ BUT WHY:`**.
- **Hashtags:** one line with **all** `tags` after the **`🟣`** line, before **`🧷`**.
- **Closing:** **`🧷 Full post (site) →`** with **`- EN:`** first and **`- ES:`** second when **`index.es.md`** exists; blank line; then **`🔗`** plus the **English** hub URL from the second **`categories`** term:
  - **`Sm-Art`** → **`🔗 Sm(art) (English) →`** `https://blog.dualsubstrate.com/categories/sm-art/` (confirm with `hugo list all` on `content/categories/sm-art/_index.md`).
  - **`T-Shirt Art`** (when used as the project hub term without **`Sm-Art`**) → **`🔗 T-Shirt Art (English) →`** `https://blog.dualsubstrate.com/categories/t-shirt-art/`.

**Compression:** Prefer concrete verbs and tight phrasing when you compress **`tldr`** / **`fluff`** for LinkedIn (drop empty softeners like "usually" unless the site wording depends on them). Keep the same ideas as the site, not a different thesis.

**Template**

```
"When the model feels threatened"

❓ [Plain-text tldr]

AND:
[Plain-text fluff]

🟣 TS5: Sm(art)

#TShirtArt #Denial #Agency #RealityCheck

🧷 Full post (site) →
- EN: https://blog.dualsubstrate.com/cognitive-memetics/t-shirt-art/2026-05-16-smart-cover-your-ears/
- ES: https://blog.dualsubstrate.com/es/cognitive-memetics/t-shirt-art/2026-05-16-smart-cover-your-ears/

🔗 Sm(art) (English) → https://blog.dualsubstrate.com/categories/sm-art/
```

**Reference bundles:** `content/cognitive-memetics/t-shirt-art/2026-05-04-tshirt-not-protection-escape/` (TS4); `content/cognitive-memetics/t-shirt-art/2026-05-16-smart-cover-your-ears/` (TS5, with **`🔖`** citations).

### Pawtropolis (Under Fire) (`content/cognitive-memetics/pawtropolis/`, `type: panel` or `type: sayings`)

**Fold-first layout** (same fold as [Sm(art)](#smart--t-shirt-art-cognitive-memetics-type-sayings-contentcognitive-memetics-t-shirt-art) and [Reptilocracy](#reptilocracy-contentcognitive-memeticsreptilocracy-type-sayings)): quoted hook on line one; scene before the series badge; series framing under **`🟣`**.

**Rules**

- **Plain text only:** strip Markdown from **`description`**, body copy, and i18n-derived series text.
- **Line one:** quoted hook. For **`type: panel`**, use **`title`** in quotes (episode name). **MUST NOT** put **`{heading_code}: {project}`** on line two.
- **Meaning:** `❓ ` immediately before adapted **`description`** (or opening beat) on the **same line**. **MUST NOT** use `✔️ TLDR:` or `❓ BUT WHY:` as a label on this line.
- **Scene:** `AND:` on its own line, then the first substantive body paragraph after **`<!--more-->`** in English **`index.md`** (or equivalent opening scene). **MUST NOT** use `➕ Context:`.
- **Series block:** `🟣 ` immediately before **`{heading_code}: Pawtropolis (Under Fire) 🐾`** on the **same line** (canonical **`project`** string). Blank line, then paste **`pawtropolisProjectAboutBody`** from **`i18n/en.toml`** as plain text (`\n\n` = paragraph break; strip Markdown). **MUST NOT** invent or paraphrase series copy. **MUST NOT** duplicate the **`AND:`** scene in the series block.
- **Omit** **`🔤 IN ENGLISH:`** and petition blocks (no Change.org CTA on Pawtropolis).
- **Hashtags:** one line with **all** `tags` after the **`🟣`** block, before **`🧷`**.
- **Closing:** **`🧷 Full post (site) →`** with **`- EN:`** first and **`- ES:`** second when **`index.es.md`** exists; blank line; **`🔗 Pawtropolis (Under Fire) (English) →`** hub URL (`/categories/pawtropolis-under-fire/`).

**Reference bundle:** `content/cognitive-memetics/pawtropolis/2026-05-13-01-i-want-this-to-be-over/linkedin.txt` (P1). Series body under **`🟣`**: paste from **`pawtropolisProjectAboutBody`** in **`i18n/en.toml`** (do not hardcode in this skill).

### Cube cows / Tales from the Cube Farm (`cognitive-memetics`, `type: panel`)

**Fold-first layout** for **Tales from the Cube Farm** strips (**`Cube-Cows`** hub): quoted strip teaser on line one only; series badge and framing after the fold.

**Rules**

- **Plain text only:** strip `**bold**` from `description` when quoting the teaser; strip Markdown (`*italics*`, `**bold**`, curly quotes as needed) from **`cowsProjectAboutBody`** when pasting the series block.
- **Line one:** the teaser in straight double quotes. Use the **`description`** field (the strip line), not the browser-tab **`title`**, unless you deliberately want the episode name instead.
- **MUST NOT** put **`{heading_code}: {project}`** on line two.
- **Omit** **`❓`** and **`AND:`** on standard image-only strips (no body below **`<!--more-->`**; the teaser is the whole joke).
- **Series block:** `🟣 ` immediately before **`{heading_code}: Cube-Cows 🐮📈`** on the **same line** (for example `W7: Cube-Cows 🐮📈`). Blank line, then paste **`cowsProjectAboutBody`** from **`i18n/en.toml`** as plain text (`\n\n` = paragraph break; strip Markdown). **MUST NOT** invent or paraphrase series copy. **MUST NOT** use `❓ BUT WHY:` as a label.
- **Hashtags:** one line with **all** entries from front matter `tags` (see [Hashtags from site `tags`](#hashtags-from-site-tags)), after the **`🟣`** block, before site and hub links.
- **Closing:** **`🧷 Full post (site) →`** with **`- EN:`** first and **`- ES:`** second when **`index.es.md`** exists; blank line; then **`🔗 Cube-Cows (English) →`** hub URL (see below).

### Raymond (`content/cognitive-memetics/raymond/`, `type: panel`)

**Fold-first layout** (same image-only fold as Cube-Cows): quoted **`description`** on line one; series badge after the fold.

**Rules**

- **Plain text only:** strip `**bold**` from `description`; strip Markdown from **`raymondProjectAboutBody`**.
- **Line one:** quoted **`description`** teaser. **MUST NOT** put **`{heading_code}: {project}`** on line two.
- **Omit** **`❓`** and **`AND:`** on standard image-only strips.
- **Series block:** `🟣 ` immediately before **`{heading_code}: Everyone ❤️ Raymond 🐕📓`** on the **same line** (canonical **`project`**; heart replaces “loves” in display). Blank line, then paste **`raymondProjectAboutTitle`** and **`raymondProjectAboutBody`** from **`i18n/en.toml`** as plain text. **MUST NOT** invent or paraphrase series copy.
- **Hashtags:** all front matter `tags` (Raymond episodes MUST include **`EveryoneLovesRaymond`**, **`PetLife`**, **`TheCutestDog`**, **`OfficeCulture`**, and **`OfficeSatire`**; do **not** add bare **`Raymond`**, **`Dog`**, **`Dogs`**, or **`DogLife`** when the series pet tags are present). Pet tags = character discovery; office tags = subject matter. Prefer order: series + pet, then office anchors, then episode hooks. Do **not** claim tags “trick” a platform classifier.
- **Quoted teaser:** MUST keep the dog signal from site **`description`** (prefer **Dawg Raymond** in EN). MUST also make the episode’s workplace tension clear in that same line-one teaser (office joke first; dog signal still required).
- **Closing:** **`🧷 Full post (site) →`** (**`- EN:`** then **`- ES:`** when bilingual); blank line; **`🔗 Everyone ❤️ Raymond (English) →`** hub URL (`/categories/raymond/`).

**Reference bundle:** `content/cognitive-memetics/raymond/2026-08-28-collaborating-with-myself/linkedin.txt` (R1).

**Template** (series body: paste from **`cowsProjectAboutBody`** in **`i18n/en.toml`**; do not hardcode it here)

```
"This week, the bell curve decides who 'exceeds expectations' and who just needs more cowbell."

🟣 W7: Cube-Cows 🐮📈

{cowsProjectAboutBody from i18n/en.toml, plain text}

#CubeCows #BellCurve #CowbellJoke #PerformanceReview

🧷 Full post (site) →
- EN: https://blog.dualsubstrate.com/cognitive-memetics/cows/2026-04-09-cow-w07/
- ES: https://blog.dualsubstrate.com/es/cognitive-memetics/cows/2026-04-09-cow-w07/

🔗 Cube-Cows (English) → https://blog.dualsubstrate.com/categories/cube-cows/
```

**Reference bundle** for structure and tone: `content/cognitive-memetics/cows/2026-04-09-cow-w07/` (W7).

### Reptilocracy (`content/cognitive-memetics/reptilocracy/`, `type: sayings`)

**Fold-first layout** (same read as [Street Wisdom](#sayings--street-wisdom-cognitive-memetics-type-sayings)): quoted **`title`** only on line one; scene before the series badge; petition stays under the **`🟣`** block.

**Petition URL (canonical):** `https://www.change.org/p/stronger-checks-and-balances-psychological-fitness-for-australia-s-top-leaders`

**Rules**

- **Plain text only:** strip all Markdown (`**bold**`, backticks) from `tldr` and `fluff`.
- **Line one:** episode **`title`** in straight double quotes. **MUST NOT** put **`{heading_code}: {project}`** on line two.
- **Meaning:** `❓ ` immediately before adapted **`tldr`** on the **same line**. **MUST NOT** use `✔️ TLDR:`.
- **Scene block:** `AND:` on its own line, then adapted **`fluff`** (one short paragraph, about 3-5 lines). **MUST NOT** use `➕ Context:` or `➕ FLUFF:`.
- **Scene, not essay:** Adapt **`tldr`** and **`fluff`** so the reader gets the **same scene** (image, move, stakes). **MUST NOT** paste the full site **`fluff`** or retell the whole article.
- **Series + petition block:** `🟣 ` immediately before **`{heading_code}: Reptilocracy 🦎🏛️`** on the **same line**. Blank line, then:
  1. Series prose from **`reptilocracyProjectAboutBody`** in **`i18n/en.toml`** (plain text; strip Markdown links to plain wording). **MAY** compress for LinkedIn length/tone rules below; **MUST NOT** invent a different series thesis. **MUST NOT** stack site jargon (`performance incentives`, `psychological fitness filters`, `behavioural patterns`, `status games`) unless rewritten in plain English.
  2. **One short line** pointing to the concrete step: the **Change.org** petition (align with petition mention in i18n / **`reptilocracyProjectAboutCta*`**; plain wording; no long policy title required on LinkedIn).
  3. Petition lead-in from **`reptilocracyProjectAboutCtaTitle`** / button sense in **`i18n/en.toml`** (or one equivalent short sentence such as **`We need your say. If you want to help, use the link below.`**).
  4. **Petition link** (still under the series block, before the blank line before hashtags):

     ```text
     Petition (Change.org) →
     - https://www.change.org/p/stronger-checks-and-balances-psychological-fitness-for-australia-s-top-leaders
     ```

  **MUST NOT** move the petition to the **`❓`** or **`AND:`** blocks, or after hashtags. **MUST NOT** omit the petition when the bundle standard includes it. **MUST NOT** use `❓ BUT WHY:` as a label (❓ is only for the meaning line).
- **Omit** **`🔤 IN ENGLISH:`**.
- **Hashtags:** one line with **all** `tags` after the **`🟣`** block (including petition lines), before **`🧷`**.
- **Closing:** **`🧷 Full post (site) →`** with **`- EN:`** first, **`- ES:`** second when **`index.es.md`** exists; blank line; **`🔗 Reptilocracy (English) →`** hub URL.

**Length (Reptilocracy only)**

- **Body prose only:** about **200-300 words** (hard ceiling **320**); count **`❓`** / **`AND:`** / **`🟣`** block prose only, not hashtags, **`🧷`**, **`🔗`**, or petition URL lines.
- **`❓` line:** 2-3 lines (~40-60 words).
- **`AND:` block:** one paragraph, 3-5 lines (~80-120 words).
- **`🟣` block:** pattern lines (~60-80 words) plus petition lead-in and two-line petition block (~40-50 words) when required.

**Tone (Reptilocracy only)**

- Street-smart, conversational English. Cooler and plainer than the site post.
- Prefer concrete phrases (`pledges as props`, `story managers over stewards`) over academic stacks (`narrative control wrapped in performance incentives`, `audience attention span and fragmented responsibility`).
- Site copy may stay analytic; **LinkedIn MUST compress** those lines into what a reader would say out loud.

**Template** (episode **`❓`** / **`AND:`** are tone samples; series body: paste/compress from **`reptilocracyProjectAboutBody`** in **`i18n/en.toml`**; do not hardcode series copy here)

```
"Now You See It, Now You Don't"

❓ We watch the campaign promise go into the hat and never come back out. The trick isn't hiding the lie; it's making us feel silly for expecting follow-through.

AND:
Leaders gain status from bold promises, then redefine, delay, or forget them once in office. News cycles and scattered responsibility make pledges props, not commitments. Break the spirit of the promise and you can still defend the wording or point to a small gesture.

🟣 W6: Reptilocracy 🦎🏛️

{reptilocracyProjectAboutBody from i18n/en.toml, plain text; MAY compress for LinkedIn}

{petition lead-in from reptilocracyProjectAboutCta* / equivalent}

Petition (Change.org) →
- https://www.change.org/p/stronger-checks-and-balances-psychological-fitness-for-australia-s-top-leaders

#Reptilocracy #NarrativeControl #PromiseManagement #PerceptionEngineering #StatusDrivenDeception

🧷 Full post (site) →
- EN: https://blog.dualsubstrate.com/cognitive-memetics/reptilocracy/2026-05-17-the-vanishing-promise-act/
- ES: https://blog.dualsubstrate.com/es/cognitive-memetics/reptilocracy/2026-05-17-the-vanishing-promise-act/

🔗 Reptilocracy (English) → https://blog.dualsubstrate.com/categories/reptilocracy/
```

**Reference bundles**

- **Petition inside series block:** `content/cognitive-memetics/reptilocracy/2026-04-26-the-process-is-blamed/` (W3) for placement.
- **Layout + tone:** template in this section; series paragraphs from **`i18n/en.toml`**; **MUST NOT** drop the petition block from the **`🟣`** section.

### English Behaviour Engineering category link

- **Site `baseURL`:** read from **`hugo.toml`** at the repo root (currently `https://blog.dualsubstrate.com/`).
- **Cognitive-memetics hubs:** closing link targets the **English taxonomy term** for the project hub (**`Cube-Cows`**, **`Raymond`**, **`Por-Estas-Calles`**, **`T-Shirt Art`**, **`Sm-Art`**, **`Reptilocracy`**, or **`Pawtropolis-Under-Fire`**), **not** the individual post URL.
- **Must use** the **default-language** term page (path **without** `/es/`).
- If unsure of the slug, run `hugo list all` and copy the **permalink** for `content/categories/<hub-slug>/_index.md` (English).

## What to keep and what to compress

- **Keep:** the concrete hook, the logical ladder, the close, all **`tags`**
  as a hashtag line when present, and both citations.
- **Reptilocracy:** use the fixed template in **Reptilocracy** below, not the default hook + logical ladder. Keep the scene and one system punch, not a full essay ladder.
- **Also keep emotional stakes:** when compressing, preserve **one** visceral
  detail that makes the argument feel human (e.g., "four years and six figures,"
  the fear of a "terrifying hole" in identity). Do not sterilize into abstract
  mechanics. Do **not** stack multiple vignettes or repeat the same beat.
- **Compress or drop:** diagrams (Mermaid), section headings, and anything
  that only works as formatted text. Re-express the idea in prose.
- **Default drop (claims / video / `social-protocols` / `human-condition` /
  `mind-infrastructure` / `x-minds`):** second workplace vignette, "behind the curtain"
  restatements already covered by the ladder, extra `###` subsections, and
  mechanism lists the Claim already states. One example max unless the author
  asks for long form.
- **Do NOT:** write a teaser that withholds the point and ends with "read more
  on my site". The post must stand alone.
- **Carousel document caption:** that short "Tap to explore" text lives in
  **`linkedin-carousel.txt`**, not here. Use
  **`.cursor/skills/site-carousel-linkedin-caption/SKILL.md`**. Do not replace
  **`linkedin.txt`** with a teaser.

## Length

### Default posts (claims, video, `social-protocols`, `human-condition`, `mind-infrastructure`, `x-minds`)

- **Target:** **150-300 words** for the **body only**.
- **Body** runs from **line 1** (usually the site **`title`**) through the
  closing punch or question. Stop counting **before** the hashtag line.
- **Do not count** toward the target: the hashtag line, the **`📑 What's in the talk`**
  block on **`type: video`** picks, **`🧷 Full post`**
  permalinks, **`🔗`** hub links, **`🔖`** / YouTube citation lines, or
  Reptilocracy petition URL lines.
- **Default to the low end** (~150-220) when the site already has a clear
  ladder (e.g. three ✜ roles). **MAY** use up to **300** when one narrative
  hook carries the whole post (e.g. degree-and-defensiveness openers).
- Do not pad. Do not repeat. If a layer can be one line, keep it one line.
- **Tight reference (claims):** `content/social-protocols/2026-05-19-be-the-captain-not-the-vessel/linkedin.txt`.
- **Video + chapter TOC:** `content/mind-infrastructure/2026-05-21-brain-is-not-a-computer/linkedin.txt`.

**Reptilocracy:** about **200-300 words** for labeled blocks (**TLDR** /
**Context** / **BUT WHY** text only; same exclusions for hashtags, **`🧷`**,
**`🔗`**, and petition URL lines). Per-section caps in **Reptilocracy** above.
Shorter beats a mini article.

## Tone

Same register as the site post. Informal-standard English. Short declarative
sentences. Direct address ("you", "we") is fine for stakes.

**Reptilocracy:** street-smart and conversational; **cooler and plainer** than
the site. Rewrite analytic site lines for LinkedIn; do not copy them verbatim.

## Citation close format

End the post with the site permalinks (when both languages exist), then any
extra citations (video, paper, etc.).

```text
[Emoji] [short description of site post] →
- EN: [URL]
- ES: [URL]

[Emoji] [Video title] (YouTube) →
- [URL]

[Emoji or 🔖] [Author, year, paper title] → [URL]
```

Example from `content/human-condition/2026-04-10-empathy-levels/index.md` (hashtag line uses that page’s `tags`; site block uses **`🧷`** with both permalinks, not a single-arrow one-language line):

```text
#SystemsCanHijackYou #NotEveryoneFeelsTheSame #EmpathyIsAGreyArea #LetsDefineBad

🧷 Full post (site) →
- EN: https://blog.dualsubstrate.com/human-condition/2026-04-10-empathy-levels/
- ES: https://blog.dualsubstrate.com/es/human-condition/2026-04-10-empathy-levels/

🔖 Children's understanding of moral emotions (JSTOR) Nunner-Winkler & Sodian (1988, Child Development) → https://...
```

## Output

Save the finished post as `linkedin.txt` in the same bundle folder as the
source article. Example: if the source is
`content/human-condition/2026-04-10-empathy-levels/index.md`, save to
`content/human-condition/2026-04-10-empathy-levels/linkedin.txt`.

For **sayings**, same rule: `content/cognitive-memetics/sayings/<bundle>/linkedin.txt`
next to `index.md`.

For **cube cows**, same rule: `content/cognitive-memetics/cows/<bundle>/linkedin.txt`
next to `index.md`.

For **T-Shirt Art / Sm(art)** sayings under **`t-shirt-art/`**, same rule: `content/cognitive-memetics/t-shirt-art/<bundle>/linkedin.txt` next to `index.md`.

For **Reptilocracy**, same rule: `content/cognitive-memetics/reptilocracy/<bundle>/linkedin.txt` next to `index.md`.

The file is plain text: no front matter, no Markdown, just the post exactly
as it would be pasted into LinkedIn.

## Reference example

The post the user wrote for the empathy-levels article is the canonical
example for this skill. It is at
`content/human-condition/2026-04-10-empathy-levels/index.md`.

Key choices made in that post:

- Hook ends just before "see more" on a tension ("cannot fathom why")
- ✜ symbols mark each empathy layer without needing headers
- Developmental arc (kids → 7-8 year-olds → adults) is the payload
- "Internal courtroom" is a concrete image that replaces an abstract definition
- Close uses three short punched lines ("Same feelings. Different target.")
- Two citations close the post as sourcing, not as a read-more prompt
