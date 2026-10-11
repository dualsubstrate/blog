---
name: site-substack-post
description: >-
  Authors required Substack newsletter bodies in substack.md / substack.es.md
  page-bundle sidecars. Over-coffee voice, ## hook headings, neurodivergent-friendly
  typography (**bold**, *italics*, blockquotes, quoted hooks for scanability). Depth by
  Hugo type (claims, sayings, panel, video). For type video, MUST copy the index Chapter
  Guide timestamp table into the sidecar before watch CTAs. When the site body has inline figures
  (Hugo image shortcode or Markdown image), MUST mirror them in the sidecar with bundle-relative
  ![alt](file.webp). No hashtags or read-more blocks. Triggers: substack post, substack.md,
  substack.es.md, make sb-en, make sb-es, newsletter body, chapter guide, inline image,
  neurodivergent substack.
---

# Substack post (`substack.md`)

## Required files

| Locale | Index | Substack body (required) |
|--------|-------|---------------------------|
| English | `index.md` | **`substack.md`** |
| Spanish | `index.es.md` | **`substack.es.md`** |

**`make sb-en`** / **`make sb-es`** fail if the sidecar is missing or empty. There is **no fallback** to `index.md` body, `facebook-en.txt`, or `linkedin.txt`.

Metadata (title, date, tags, featured image) still comes from the index file. The pipeline merges index front matter with the sidecar body only. The **featured image** is prepended on publish; **inline body figures** live only in the sidecar (see **Inline images**).

## Philosophy

- **Adaptation over copy-paste:** Use `index.md` (and type-specific fields) as source material, not as the paste target.
- **Newsletter voice:** Warm, direct, spoken prose (same café register as friends Facebook copy in **`.cursor/skills/site-facebook-post/SKILL.md`**, but structured for email). For claims and video newsletters, MUST follow **Explanatory prose** in **`.cursor/rules/site-content-markdown-writing.mdc`** (mechanism paragraphs, not revelation stacks; audit: **revise-post** Step 2).
- **One review surface for Substack:** Maintain **`substack.md`** only; do not also maintain `facebook-en.md` for Substack.
- **Video picks:** Mirror the site **Chapter guide** table in the sidecar so newsletter readers can jump by timestamp (see **Chapter guide (`type: video` only)**).
- **Neurodivergent-friendly scan:** Guide the eye with **typography** (`**bold**`, `*italics*`, `>` quotes). **MUST NOT** rely on chopping every idea into tiny paragraphs as the main readability fix.

## Format (all types)

- Markdown body only (no YAML front matter in the sidecar).
- Short **`##`** hook headings (not seminar labels like `### The Takeaway`). A skimmer should grasp each section from the heading alone.
- **Blank line between every paragraph.**
- **No em dash (U+2014).** Use comma, semicolon, colon, or parentheses.
- **No hashtags** and **no** bilingual *If you want to read more* / *Si quieren leer más* blocks (Substack footer and tags cover that).
- **Optional** closing citations as normal Markdown link lines (from `grounding` or casual sources).
- **MUST NOT** copy `linkedin.txt` structure (✜ blocks, hub lines, hashtag line).
- **MUST NOT** use tables except the **Chapter guide** on **`type: video`** (see below). No other tables or multi-level outlines.

## Neurodivergent-friendly typography (MUST for sidecars)

Applies **only** to **`substack.md`** and **`substack.es.md`**, not to `index.md`, `linkedin.txt`, `facebook-*.txt`, or Hugo site body. Pair with **`.cursor/skills/site-revise-emphasis/SKILL.md`** for density (restrained, not wall-to-wall).

### What to use

| Tool | Use for | Example |
|------|---------|---------|
| **`**bold**`** | Turn words (*never*, *not*, *forget*), stakes (*trap*), named constructs (**glia**, **LLM**, **AGI**), and one anchor per idea | The **trap** is when you **forget** you made that cut |
| **`*italics*`** | Coined phrases, metaphor labels, fashion terms, contrast pairs | *like* becomes *is*; the *"prediction engine"* story |
| **`*"quoted phrase"'*`** or **`"quoted phrase"`** | Spoken hooks, slogans, rhetorical questions | *"How does the brain work?"* sounds simple |
| **`>` blockquote** | One punch line per section when it lands (site pull quotes work well here) | > Proving the world is simple was **never** the job. |

Rough guide: about **2–5** bold spans per paragraph; italics on phrases, not whole sentences unless the phrase is the hook.

### What not to do

- **MUST NOT** bold every noun and verb (tag-cloud paragraphs).
- **MUST NOT** split prose into one-sentence paragraphs everywhere unless the block is genuinely unreadable.
- **MUST NOT** drop typography and rely only on extra `##` headings or bullet lists.
- **MUST NOT** use LinkedIn plain-text conventions (no `**bold**` in `linkedin.txt`; Substack sidecars **do** use Markdown emphasis).

### Paragraph length

- Normal paragraph length is fine (often 2–4 sentences).
- Split a paragraph only when it packs unrelated beats or becomes a wall on mobile.
- Gloss jargon in parentheses when you keep a technical term: **glia** (support cells in the brain, not the firing neurons).

### Spanish (`substack.es.md`)

- Same typography rules: **negrita**, *cursiva*, `>` citas, comillas en preguntas retóricas.
- Match emphasis placement to the English sidecar when bilingual, but keep idiomatic Spanish wording per **`.cursor/skills/site-spanish-translation-content/SKILL.md`**.
- **MUST** mirror hierarchy and algorithm verbs from **`index.es.md`** (**promueve/promueven**, **escalan jerarquías**); **MUST NOT** slip back to **suben**, **deja subir**, or **impulsa** in the sidecar for the same beat.

### Reference bundle

**`content/mind-infrastructure/2026-05-21-brain-is-not-a-computer/`** (`substack.md`, `substack.es.md`): video pick with typography pass, full **Chapter guide** table (YouTube `t=` links), then primer + watch CTA.

## Depth by `type`

| `type` | Target length | Source | Include |
|--------|---------------|--------|---------|
| **`claims`** | Full argument (~same depth as site body) | `index.md` body + `description` / `grounding` ideas | Capture pipeline, examples, mechanism, **inline body images** when present; convert `###` to `##`; apply typography pass |
| **`sayings`** | Medium (~150–300 words) | `title`, `tldr`, `fluff`, optional body | Street hook; skip Hugo "Article" section labels; bold key beats from `tldr` / `fluff` |
| **`panel`** | Medium | `description` teaser + body if present | Cube-cow / panel payoff; optional BUT WHY tone in prose, not template blocks |
| **`video`** | Medium + watch CTA | `description`, `sowhat`, body TLDR | Invite a full watch; do not paste Claim/Grounding cards verbatim; bold/italic from `description` bullets; **include Chapter guide table** (see below) |

## Chapter guide (`type: video` only)

When the bundle has **`youtube_id`** and a **Chapter Guide** table in **`index.md`** body:

- **MUST** copy that table into **`substack.md`** / **`substack.es.md`** (same rows, timestamps, and YouTube `t=` links). Use the locale’s table from **`index.es.md`** for Spanish.
- **Heading:** `## Chapter guide` (English) or `## Guía de capítulos` (Spanish). Hook-style headings are fine if they stay scannable (for example `## Jump to a moment in the video`).
- **Placement:** after the prose argument, **before** primer links and the YouTube CTA line.
- **Why:** Newsletter readers often skim text first; timestamped jumps are high value and match the site pick (see **`.cursor/skills/site-video-content/SKILL.md`**).
- **Pipeline:** `internal/substackhtml` keeps Chapter guide tables as HTML on paste (time links in the first column get bold). If paste drops tables, retry with `go run ./cmd/substack-html -in … -tables list`.

**MUST NOT** add a chapter table for non-video types unless the user explicitly asks.

## Inline images (body figures)

When the locale **`index.md`** / **`index.es.md`** body includes an inline figure, **MUST** carry it into **`substack.md`** / **`substack.es.md`** at the **same position** in the argument (not only the featured card image from front matter).

| Site source | Sidecar form |
|-------------|--------------|
| `{{< image src="diagram.webp" alt="…" >}}` | `![…](diagram.webp)` using the shortcode **`alt`** (or a tight paraphrase) |
| `![Caption](diagram.webp)` in index body | Same line in the sidecar |
| `{{< mermaidfile >}}` with sibling **`diagram.webp`** | `![…](diagram.webp)` (pipeline expands mermaid to image when a webp sibling exists) |

Rules:

- **MUST** use **bundle-relative** paths (`factory.webp`, not `/social-protocols/...`). **`substackhtml`** resolves them to `https://…` on paste when the post is live (`PagePermalink` from `hugo list all`).
- **MUST** copy **`alt`** from the **locale** index (English alt in **`substack.md`**, Spanish alt in **`substack.es.md`**).
- **MUST NOT** rely on Hugo shortcodes in the sidecar; they are stripped on export.
- **MUST NOT** duplicate **`featuredImage`** / **`featuredImagePreview`** in the sidecar body; publish prepends that lead image separately.
- Blank line before and after the image line, same as between paragraphs.

**Reference bundle:** **`content/social-protocols/2026-05-26-developed-countries-are-factories/`** (`factory.webp` after the Nordics / citizen-check section).

## Execution

1. Read **`index.md`** / **`index.es.md`** and note **`type`**.
2. Draft **`substack.md`** (and **`substack.es.md`** when `index.es.md` exists) using the depth table above.
3. **`type: video`:** append the **Chapter guide** table from the index body when present.
4. **Inline figures:** for each body image in the locale index, add `![alt](file.webp)` at the matching position (see **Inline images**).
5. **Typography pass (MUST):** add **`**bold**`**, **`*italics*`**, **`>`** punch lines, and quoted hooks; verify ~2–5 bold spans per paragraph, not saturation.
6. Preview: `go run ./cmd/substack-html -in content/<section>/<slug>/substack.md` (use `substack.es.md` for Spanish)
7. Publish: `make sb-en POST=<section>/<slug>` or **`make sb-es`** for Spanish.

## Skeleton (`substack.md`)

```markdown
## [Hook heading: tension or question]

You **[verb anchor]** to think. … The **trap** is when you **forget** …

> [One punch line with optional **bold** on the turn word.]

## When *[contrast A]* becomes *[contrast B]*

… *metaphor label* … **named construct** (short gloss in parentheses) …

The comparison **stops being a comparison**. …

## [Next hook heading]

*"Rhetorical question?"* sounds simple. **It is not one question.**

… **mechanism**, **story**, or **reassurance** … *"quoted slogan"* …

![Alt text from the locale index](figure.webp)

## [Watch / close heading]

… *mind-as-software* … **prediction** in AI is **not** the same as **understanding** …

## Chapter guide

| Time | Chapter |
| --- | --- |
| [0:00](https://www.youtube.com/watch?v=VIDEO_ID&t=0) | **Speaker** Topic label |
| … | … |

[Primer link sentence with Markdown link.]

Watch on YouTube: [Video title](https://www.youtube.com/watch?v=VIDEO_ID)
```

## Related files (do not duplicate for Substack)

| File | Role |
|------|------|
| `index.md` | Site canonical |
| `linkedin.txt` | LinkedIn + Facebook Page autopost (plain text, no Markdown emphasis) |
| `facebook-en.txt` | Friends Facebook only (plain text) |

## Related skills

- **Video picks (site Chapter guide):** `.cursor/skills/site-video-content/SKILL.md`
- **Emphasis density:** `.cursor/skills/site-revise-emphasis/SKILL.md`
- **Friends Facebook:** `.cursor/skills/site-facebook-post/SKILL.md`
- **LinkedIn:** `.cursor/skills/site-linkedin-post/SKILL.md`
- **Spanish bundles:** `.cursor/skills/site-spanish-translation-content/SKILL.md`
- **Pipeline safety:** `.cursor/skills/site-substack-pipeline-safety/SKILL.md`
- **Image paste / URL resolve:** `internal/substackhtml/README.md` (Images section)
