---
name: site-video-content
description: >-
  Authors and edits Hugo posts with type video: YouTube embed via youtube_id or
  body shortcode, description as novelty lead, sowhat payoff, list-row
  embed plus fullPost CTA, categories (often Mind-Infrastructure,
  Human-Condition, or Social-Protocols by topic), tags, optional featured image for cards, and a
  thesis-first TLDR body (so-what article) for text-first readers and feed skims, while the lead still invites a full watch. Use when editing or adding video picks, when
  the user mentions type video, youtube_id, video archetype, curated videos, or
  chapter notes and summaries.
---

# Video content type (`type: video`)

## What this type is for

**Video** posts pair a **YouTube embed** with a **stand-alone article** below it. **Primary goal:** the talk is worth watching; **`description`** novelty open loops should **sell the play button**. **Parallel goal:** many people **read first** (feeds, quick tabs, preview-then-decide) and some **never** press play, so the **body** still delivers a **thesis-first** argument, mechanism, and **so what** in text. Voice follows the site-wide **commute test** in **`.cursor/rules/site-content-markdown-writing.mdc`** → **Voice while drafting** (and essay pack **`shared_guide`**). The embed remains the full experience you recommend.

**Essay MCP / `content_form=video`:** The same band split (DESCRIPTION novelty, SOWHAT payoff, thesis-first TLDW, optional CHAPTER_GUIDE jump map) is the host pack in **`data/essay-content-forms.yaml`** (`forms.video`). MCP spine/arc/composition/polish inject that guide. Composition Gate review still uses **Unknown yet / Will know / Hook / blocks / Close**; those fields map to DESCRIPTION / SOWHAT / TLDW per the pack. Keep this skill and that YAML aligned; do not put per-piece band rules only in WIP notes.

**UI (this repo):** `layouts/video/single.html` renders **title** → **subtitle** (optional) → **meta** → **featured image** only if there is **no** `youtube_id` → **tags** → optional **TOC** → **lead** (see below) → **embed** (when `youtube_id` set) → **TL;DW** (body TLDR) → **Chapter Guide** (when present) → optional **Dig deeper** (`reader_landing` or legacy `research`) → **piece footer** (updated date, tags + shares on one row, prev/next) → optional **Related** shelf (`related` or shared-tag fallback; RELATED band above the cards) → **Back | Home**.

**Lead on the single (order):** If **`sowhat`** is set, **`description`** is shown under fixed **`h3`** **“What you probably do not know yet”** (🎬), then **`sowhat`** under **`h3`** **“What you will know after”** (🎯). If **`sowhat`** is omitted, **`description`** has no extra heading above it. Both fields are markdownified.

**After the embed:** The template adds **`h2` TL;DW** (⏩, “too long, didn’t watch”) above the TLDR body, then **`h2` Chapter Guide** (📑) when the bundle has a chapter table. Body hook sections stay **`###`** (`h3` under TL;DW). Legacy body **`##`** is demoted to **`h3`** in the layout. Do **not** rely on a body heading for **TL;DW** or **Chapter Guide**; the template supplies those bands.

**Section list rows** (`layouts/partials/seven-style-row.html`): **Left column** — same **YouTube** embed as on the single (not a still), then a pill link to the post using **`{{ T "fullPost" }}`** (English string in **`i18n/en.toml`**, e.g. “Read the article”; the iframe is not a link to the single). **Right column** — **`description`** and optional **`sowhat`** with labels **“What you probably do not know yet”** / **“What you will know after”** (fixed in the partial; not per-page front matter). **Note:** In the default list aside, **“What you probably do not know yet”** appears above **`description`** even when **`sowhat`** is empty; the **single** only adds that heading when **`sowhat`** is set—so write **`description`** so it still reads well with that list label, or add **`sowhat`** when you want single and list to match.

**Home feed** (`layouts/index.html` → **`layouts/_default/home-tile.html`**): **`type: video`** with **`youtube_id`** (or **`youtube`**) uses the **same YouTube embed** in the card (iframes cannot wrap in the card’s outer link, so the template uses a **split card**: embed block, then a **text link** to the single with summary + optional **`fullPost`** pill). Featured image is **not** shown as the card hero when the embed is used. **`assets/css/_custom.scss`**: prose blurbs use **line-clamp + ellipsis**; **`description`** with a **top-level list** uses block layout, a **taller cap**, a **bottom fade**, and at most **four** list items so bullets do not collide with `-webkit-line-clamp` or clip mid-line.

**`date`:** MUST use the site **`date`** form in **`.cursor/rules/site-content-markdown-writing.mdc`** → **Publish `date`**.

## Field roles

| Field | Role |
|-------|------|
| **`type`** | Must be **`video`**. |
| **`description`** | **Novelty lead** above the player (Hugo: **“What you probably do not know yet”** when paired with **`sowhat`**): 3–4 cold-readable **open loops** that sell pressing play. Not topic-gap lines (“why tools fail”); not a thesis paragraph. Also drives the **list row** aside. |
| **`sowhat`** | **Required payoff** in practice (essay MCP requires the SOWHAT band): one short paragraph on what the viewer gains (**“What you will know after”**). On single and list, rendered after **`description`**. When set, the layout adds **“What you probably do not know yet”** above **`description`**. |
| **`youtube_id`** | The id from `https://www.youtube.com/watch?v=THIS` (or the `v=` value in a short URL). The layout injects Hugo’s **`youtube`** shortcode into **`.video-page__embed`**. Alias front matter key **`youtube`** is accepted. |
| **Body** | **Thesis-first TLDR / “so what” article** **below** the embed: one accruing claim with evidence, mechanisms, examples, and your angle. Use **`###`** hook sections so skimmers can scan (see **Body headings**; do not use **`##`** in the site body). **Not** talk chronology and **not** a novelty stack (novelty lives in **`description`**). Optional **Chapter Guide** table at the end for jump-in after reading. Default intent: **substantive** text backup, because plenty of traffic never plays the file. |
| **`subtitle`** | Optional second line under the **`title`** on the **single** page only (not list rows or home tiles). See **Subtitle (optional)** below. |
| **`categories`** | Taxonomy hubs. Often **`Mind-Infrastructure`**, **`Human-Condition`**, **`Social-Protocols`**, or **`X-Minds`** (or **`Reality-Protocols`** under claims) by topic; pick what matches the post (see other posts in that section). For **`content/human-condition/`**, you MAY use **`Human-Condition`** plus **exactly one** theme hub (second term), same pattern as Cognitive-Memetics umbrellas (**`.cursor/skills/site-claims-content/SKILL.md`** → **Human-Condition theme hubs**). For **`content/x-minds/`**, use **`X-Minds`** first. |
| **`related`** | Optional **keep-reading** paths (Hugo `GetPage` strings). Layout: **one** claims/video banner on top, then up to **two** sayings/panel on the next row. Empty or omitted: layout fills from **shared tags**. MUST **not** dump related links into the TLDR body. See **`layouts/partials/related-keep-reading.html`**. |
| **`reader_landing`** | Optional **Dig deeper** after TL;DW / Chapter Guide: `why_it_matters`, `takeaways`, `explore` (related / source / prefilled Perplexity queries / optional completed thread). See **`layouts/partials/reader-landing.html`**. Legacy **`research`** markdown still works when `reader_landing` is absent. |

## Subtitle (optional)

**UI:** `layouts/video/single.html` renders **`subtitle`** as `<h2 class="single-subtitle">` directly under the **`title`**. Feeds and list rows still use **`title`**, **`description`**, and **`sowhat`** only.

**Split of labor**

| Field | Role |
|-------|------|
| **`title`** | Your **hook** or thesis (may use leading emoji under allowed sections). |
| **`subtitle`** | Plain **attribution + topic** (who spoke, or what the pick is about) when the title cannot carry that without weakening the hook. |
| **`description`** / **`sowhat`** | Capture and stay payoff; do **not** move hook bullets into **`subtitle`**. |

**Use `subtitle` when**

- **`title`** is a **thesis or metaphor** and the reader needs **speaker credit** or a **plain topic line** on the article page (example: **`content/social-protocols/2026-05-30-the-prediction-business/`**).
- One **named speaker** is the main source and is **not** already obvious from **`title`** (example: Max Bennett, Barbara Tversky).
- The pick is a **curated multi-speaker** page: list the voices and the shared topic in one line (example: Patric Gagné, M.E. Thomas, and James Fallon on sociopathy and scaffolding).
- The embed is a **multi-speaker compilation**: use a **topic** subtitle, not one name (example: On brain metaphors, prediction machines, and what counts as understanding).

**Skip `subtitle` when**

- **`title`** already names the talk, game, or principle (seminar-style title).
- **`subtitle`** would mostly **repeat** the first **`description`** bullet or **`sowhat`**.
- The source is a **channel explainer** with no person to credit on the page header (Kurzgesagt-style picks).

**MUST**

- Keep **`subtitle`** to **one short line** (speaker + topic, or topic only for compilations).
- Mirror **`subtitle`** on **`index.es.md`** when the English post has one (idiomatic Spanish, same attribution).
- Verify speaker names against the **YouTube source** (or the body’s credited interview) before publishing.

**MUST NOT**

- Add **`subtitle`** to every video by default.
- Put decorative emoji in **`subtitle`**.
| **`tags`** | Punchy hooks; align with **`.cursor/skills/site-claims-content/SKILL.md`** *Tag voice* for shape (PascalCase, reusable, no duplicates). Read **`data/tag-register.txt`** and **`.cursor/skills/site-tag-register/SKILL.md`** before minting new tag strings. |
| **Featured image** | Prefer **`featuredImage: "file.ext"`** and optional **`featuredImagePreview: "file.ext"`** (page-bundle local resource). Use this for **cards and social previews**. If **`youtube_id`** is set, the **large hero image is hidden** on the single page so the page is not image + player; the image can still help list thumbnails depending on theme partials, so keep it if you want a strong card. If you need advanced resource metadata, you MAY instead use `resources` with `name: "featured-image"` and `src: "file.ext"`. |

## Chapter Guide (optional)

For long videos (>15m) or dense talks, include a **Chapter Guide** at the bottom of the body. Job: **timestamp jump map** so readers can land in the talk after the TL;DW. **MUST NOT** write novelty slogans or second-teaser lines as chapter titles; keep plain, orienting labels.
- Use **`### Chapter Guide`** (Spanish: **`### Guía de capítulos`**) in the markdown source so the layout can split the table into its own **`h2`** band after **TL;DW**. Do **not** add a duplicate **`## Chapter Guide`** heading for display; the template renders **`h2` Chapter Guide** (📑).
- Use a **2-column Markdown table** so the timecodes align.
- Link directly to the YouTube timestamp (e.g., `&t=694`).
- Keep chapter titles descriptive but concise (what that segment is about, not a curiosity hook).
- This helps readers navigate the infrastructure of the talk without watching the whole long video.

### Chapter Guide formatting (this repo)

- Use a short heading: **`### Chapter Guide`** (site body).
- Use this table shape:

| Time | Chapter |
| --- | --- |
| `[MM:SS](https://www.youtube.com/watch?v=VIDEO_ID&t=SECONDS)` | `**Label** Chapter title` |

- Bold only the first label word or phrase (for scanability), not the whole chapter title.
- In this repo, the visual separator after the label is handled by CSS for video pages, so the chapter text should not include a literal em dash character.
- Prefer the `&t=` seconds parameter, since it is easy to generate from a timestamp (for example 1:36 is `t=96`).

### Substack (`substack.md` / `substack.es.md`)

When you add a **Chapter Guide** on the site, **MUST** copy the same table into the bundle sidecar for Substack (English from `index.md`, Spanish from `index.es.md` under `## Guía de capítulos`). Place it after the newsletter prose and **before** primer links and the YouTube watch line. Full rules: **`.cursor/skills/site-substack-post/SKILL.md`** → **Chapter guide (`type: video` only)**.

### LinkedIn (`linkedin.txt`)

When the site has a **Chapter Guide**, **MUST** add **`📑 What's in the talk`** to **`linkedin.txt`**: one line per row, `• Label · chapter title`, **no clock times**, after the close and before hashtags. No pipe or space-aligned tables (LinkedIn uses a proportional font). Full rules: **`.cursor/skills/site-linkedin-post/SKILL.md`** → **Video chapter guide (TOC)**. Reference: **`content/mind-infrastructure/2026-05-21-brain-is-not-a-computer/linkedin.txt`**.

### Facebook (`facebook-en.txt` / `facebook-es.txt`)

When the site has a **Chapter Guide**, **MUST** add the same outline shape to friends Facebook copy: **`📑 What's in the talk`** (EN) or **`📑 De qué va el video`** (ES), `• Label · chapter title`, **no clock times**, after the close and before hashtags. Write fresh friends prose; do not copy **`linkedin.txt`**. Full rules: **`.cursor/skills/site-facebook-post/SKILL.md`** → **Video chapter outline**. Reference: **`content/mind-infrastructure/2026-05-21-brain-is-not-a-computer/facebook-en.txt`** and **`facebook-es.txt`**.

## Embed: one source of truth

- **Preferred:** set **`youtube_id`** in front matter and put **only** commentary in the **body** (no duplicate `{{< youtube >}}` in the body).
- **Alternative:** omit **`youtube_id`** and place **`{{< youtube VIDEO_ID >}}`** in the **body** where you want it (still responsive; no `.video-page__embed` wrapper unless you add your own markup).

Do **not** set **`youtube_id`** **and** repeat the same id in a body shortcode (double embed).

## Lead, TLDR body, and optional Chapter Guide

**Typical paths:** some readers **watch first** or right after the lead; others **read the TLDR** and only then play, or **never** play. The page should **invite the watch** and still deliver a clear **“so what”** in text: what the talk says, why you care, and how it connects elsewhere.

**Division of labor**

| Part | Role |
|------|------|
| **`description`** | **Novelty open loops** for list and card views: concrete surprises or stakes that make **pressing play feel worthwhile**. Bullets are fine (about 3–4). MUST read well **above** the player and in **feeds**; it is **not** the full article. **MUST pass the cold-read gate** (below). |
| **`sowhat`** | **One-paragraph** umbrella payoff (list + single). Required for essay MCP video form; strongly preferred on published pages. The **body** still does the real TLDR work. **MUST pass the cold-read gate**. |
| **Body** | **Thesis-first** so-what article / TLDR after the embed: structured prose (**`###`** hook sections), one accruing claim, mechanisms, examples, and your angle. A reader who never presses play should leave with the gist. Use **active, plain** explanations; name constructs when they matter. Short **speaker quotes** MAY punctuate a section. **`###` headings MAY stay punchy.** Body paragraphs MUST NOT restack the section hook in metaphor shells or industry verbs (fail: "sells the blend as 'me'", "steer the rewrite", "keeps assembling it"). Patterns: **Metaphor-shell restack** and **Industry-verb shells** in **`.cursor/skills/site-revise-post/reference.md`**. |
| **Chapter Guide** | Optional **end** table: timestamps + plain chapter labels for jumping into the video after the summary. Not novelty slogans. |

## Cold-read gate (`description` + `sowhat`) (MUST)

### Objective (what you are optimizing for)

**`description`** and **`sowhat`** decide whether you **capture** the reader or they **pass**.

- On **feeds and list rows**, most people never open the post. They only see **`title`**, **`description`**, and (when set) **`sowhat`**. If a line does not make sense in a few seconds, they scroll on. The body and the video never get a chance.
- **`description`** = **capture the click** with **novelty open loops** (why open this page or press play).
- **`sowhat`** = **capture the stay** (what they get if they read or watch; why it is worth their time).
- The **body** is for people who already said yes. It must stand alone as a **thesis-first** TLDR for text-first readers, but it cannot rescue opaque list copy.

**Success:** A stranger understands each line and feels a reason to open or keep going. **Failure:** They need the essay, the video, or your private metaphors to decode the line.

Many visitors have **not** read the body and **not** watched the video. Treat **`description`** and **`sowhat`** as two separate front-matter sections that **always** get a **naive-reader (cold) test** before publish or before marking a revision Pass.

**How to test:** Read each field alone. Ask: *Would someone with no context know what this claims and why it matters?* If any line needs the essay, the video, or an in-joke to decode, it **fails**.

**`description`**

- MUST cold-read **every** bullet (or paragraph) in **`description`**.
- MUST give **extra scrutiny to the first two bullets** when using a list; they are where compressed metaphors most often fail (shorthand, *like*/*is*, “bench/kit,” unnamed “tools”).
- Each bullet MUST name the subject (**brain metaphors**, **LLM**, **simplify**) and MAY use **one short parenthetical example** before a punch line.
- MUST NOT rely on metaphors introduced only in the body (“yesterday’s kit,” “mouse trap,” “spherical cow”) unless you gloss them in the same line.

**`sowhat`**

- MUST cold-read **`sowhat`** when the key is set, same standard as **`description`**.
- MUST NOT use body-only shorthand (“bench,” “trap,” “slice”) without a plain anchor in the same sentence.
- MUST apply **`.cursor/skills/site-revise-post/SKILL.md`** → **Step 2** (AI voice) and **Step 5** (negation fluff) on **`sowhat`** and **`description`**; full banned-pattern table lives there, not in this skill.

**`sowhat` fail examples (rewrite):**

| Fail | Why |
|------|-----|
| "This long interview bridges felt fatigue and cellular allocation: not more fuel, better aim." | Bridge verb + semicolon aphorism |
| "Watch if you want the mechanism behind why exhaustion and meaning can feel like the same problem." | Meta CTA; reader needs the video to decode |
| "Picard's hook is a fixed energy budget." | Editor-speak; label without same-sentence gloss |

**Pass:** A stranger could quote back the point of each line. **Fail:** You have to say “well, earlier in the post…” or “after you watch…”

See **`.cursor/skills/site-revise-post/SKILL.md`** Step 1 and Step 3 for the revision workflow.

**MUST**

- Make **watching** feel like the natural next step (strong lead, concrete payoff, why this speaker or framing).
- Write the **body** so it **stands alone** for text-first visitors, without implying “you failed if you did not watch.”
- Keep **`description`** as **novelty hooks**, not a duplicate of the whole body.
- Run the **cold-read gate** on **`description`** and **`sowhat`** on every new or updated video post (not optional).
- Prefer a filled **`sowhat`** on every video pick (essay MCP `content_form=video` requires the SOWHAT band).
- Use **educational, active voice**. Explain terms in plain language when you introduce them. Body, `description`, and `sowhat` MUST follow **Explanatory prose** in **`.cursor/rules/site-content-markdown-writing.mdc`** (claim before interpretation; 2–4 sentence paragraphs; no rhetorical fragment stacks; no **Metaphor-shell restack**; no **Industry-verb shells**).

**MUST NOT**

- Paste a **verbatim transcript** or scene-by-scene narration of the video.
- Replace the body with **only** a timestamp table (put long maps in **Chapter Guide** at the bottom, after the TLDR prose).
- Dump **course-rubric** boilerplate (rubrics, week-by-week syllabi) meant for **`courses/**`** onto a normal video pick.

**MAY**

- Use a **shorter** body when the talk is narrow and **`description`** + **`sowhat`** already carry the full gist (rare; default is still a real TLDR below the fold).
- Move **archive-grade** notes (raw outline dumps) to a **separate** page and **link once** from the body.
- Use **`cognitive-science-methodology.md`** (repo root) only as **light** editing discipline (chunking, cut fluff). Do **not** paste full course templates onto **`type: video`** pages.

## Section paths

Posts are normal pages under a section folder, for example:

- `content/social-protocols/<slug>/index.md`
- `content/mind-infrastructure/<slug>/index.md`
- `content/human-condition/<slug>/index.md`
- `content/x-minds/<slug>/index.md`

**Section** comes from the **folder**, not from `type`. Pick the folder by **main job** (same intent as **`claims`**): see **`.cursor/rules/site-content-placement.mdc`** → **`type: video` (YouTube picks)**. Use **`social-protocols/`** when the pick is really about norms, reciprocity, coordination, or institutions; **`human-condition/`** when the core is person-level psychology or identity; **`x-minds/`** when the core is mixed-wiring lives, parents, or community; **`mind-infrastructure/`** as the **default lane** for general-interest picks when the other lanes do not fit. Taxonomy **`categories`** are independent of the folder. `hugo.toml` lists **`social-protocols`**, **`mind-infrastructure`**, **`human-condition`**, and **`x-minds`** in **`params.home.contentSections`** so new video posts in those folders can appear on the home feed.

## Authoring workflow

1. Create with **`hugo new content/<section>/<slug>/index.md --kind video`** (archetype: **`archetypes/video.md`**).
2. Set **`title`**, **`description`** (novelty lead), **`sowhat`** (payoff; add the key in front matter; the archetype may omit it by default), **`youtube_id`**, **`categories`**, **`tags`**, **`draft`**, optional **`subtitle`** when **Subtitle (optional)** applies, and optional featured image resource.
3. **Recommend `related`:** after tags, search for up to **two** **cognitive-memetics** sayings or panels that share the mechanism or a punchy parallel; MAY add **one** other claim or video. Write Hugo paths into **`related`**. If none fit, leave empty for tag fallback. MUST **not** paste related links into the TLDR body.
4. Write the **body** as the **thesis-first TLDR / so-what article** below the embed (see **Lead, TLDR body, and optional Chapter Guide** above). Use **hook `###` headings** per **`.cursor/skills/site-revise-hooks/SKILL.md`**. Add **Chapter Guide** (orienting labels, not novelty slogans) when the talk is long or dense.
5. If you maintain **`substack.md`**, copy the **Chapter Guide** table into the sidecar per **`.cursor/skills/site-substack-post/SKILL.md`**.
6. If you maintain **`linkedin.txt`**, add the chapter outline per **`.cursor/skills/site-linkedin-post/SKILL.md`**.
7. If you maintain **`facebook-en.txt`** / **`facebook-es.txt`**, add the locale outline per **`.cursor/skills/site-facebook-post/SKILL.md`**.
8. **Cold-read** **`description`** (especially the **first two bullets**) and **`sowhat`**; rewrite until each line passes without the body or video.
9. Run a local **`hugo`** build before publishing.

## Style and repo rules

- MUST follow **`.cursor/rules/site-content-markdown-writing.mdc`** for English, punctuation, and prose habits.
- Styling hooks: **`assets/css/_custom.scss`** — single: `.single-video`, `.video-page__embed`, `.video-page__lead`; section list video row: `.seven-list__figure--video`, `.seven-list__video-embed`, `.seven-list__video-detail-link`.

## Body headings (`###` on site pages)

- **`layouts/video/single.html`** outline: **`h1`** (`title`), optional **`h2`** (`subtitle`), optional **`h3`** lead/payoff before the embed, **`h2` TL;DW** after the embed, **`h3`** TLDR hook sections in the body, optional **`h2` Chapter Guide** at the end (template bands; split from body markdown at **`### Chapter Guide`** / **`### Guía de capítulos`**).
- MUST follow **`.cursor/skills/site-revise-hooks/SKILL.md`** → **Body headings (`##` and `###`) — hooks, not labels** for hook style, cold-read, and do/don't examples.
- Put the timestamp table immediately below **`### Chapter Guide`** (or **`### Guía de capítulos`** on **`index.es.md`**) in the source file; the template strips that heading and renders the table under **`h2` Chapter Guide**.
- **`substack.md`** / **`substack.es.md`** sidecars MAY keep **`## Chapter guide`** / **`## Guía de capítulos`** per **`.cursor/skills/site-substack-post/SKILL.md`** (newsletter layout, not the video single template).
- Older video bundles may still use body **`##`** until revised; default for new and updated posts is **`###`**.

## References in this repo

- Archetype: **`archetypes/video.md`**
- Layouts: **`layouts/video/single.html`**, **`layouts/partials/video-body-sections.html`**, **`layouts/partials/video-table-of-contents.html`**, **`layouts/partials/related-keep-reading.html`**, **`layouts/partials/seven-style-row.html`** (list row embed + CTA + aside), **`layouts/_default/home-tile.html`** + **`layouts/partials/home-tile-body.html`** (home grid)
- List CTA string: **`i18n/en.toml`** → **`[fullPost]`**
- Built-in embed: LoveIt docs *Theme Documentation - Built-in Shortcodes* → **youtube** (see `themes/LoveIt/exampleSite/content/posts/theme-documentation-built-in-shortcodes/index.en.md`).
