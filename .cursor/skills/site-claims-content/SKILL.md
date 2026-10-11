---
name: site-claims-content
description: >-
  Authors and edits Hugo posts with type claims: Claim (description), concise Grounding
  (short digest plus Source line with Markdown link: article or paper title as anchor text, not a bare URL),
  optional primary-source quote blocks (`###` heading + blockquote), prose style, and
  preserving the author’s direct voice when they supply finished copy (no fluff rewrites).
  Use when editing or adding content under social-protocols, human-condition, x-minds, or
  other sections using claims, when the user mentions Claim, Grounding, categories,
  tags, or claims archetype, or when shaping list-view copy for those posts.
---

# Claims content type (`type: claims`)

## What this type is for

Posts use front matter fields that split **narrative** from **intellectual anchor**. The **list view** (cards, summaries) leans heavily on **title**, **`description`**, and **`grounding`**; polish those before spending time on body-only formatting.

**UI (this repo):** Any page with **`type: claims`** uses the same list and single labels (**Claim**, **Grounding**, **Dig deeper** where applicable), regardless of section. Section only affects URL path and default category display when categories are omitted.

**`date`:** MUST use the site **`date`** form in **`.cursor/rules/site-content-markdown-writing.mdc`** → **Publish `date`**.

## Field roles

| Field | Role |
|-------|------|
| **`description`** | **Claim**: the narrative assertion in your voice. What the reader should believe or notice. Often shown as the **Claim** block in the UI. MUST **cold-read** on cards: no **Claim fog** (abstract jargon a stranger cannot picture). See **`.cursor/skills/site-revise-hooks/SKILL.md`** → **Claim fog**. |
| **Body** (markdown under front matter) | **Thoughts** on the detail page: essay, metaphor, examples, optional **primary-source quote** section. Rendered **after** Claim, **before** Grounding. |
| **`grounding`** | **Support for the Claim**: definitions, named constructs, dates, and **links** (paper, overview) so the Claim is tied to sources or standard terminology. **Source** lines use **Markdown** `[title or short label](url)` so readers see a **title**, not a raw URL string. **Keep it short** (see **Grounding** below); it is not a second Claim and not a duplicate of the full article. |
| **`related`** | Optional **keep-reading** paths (Hugo `GetPage` strings). Layout: **one** claims/video banner on top, then up to **two** sayings/panel on the next row. Empty or omitted: layout fills from **shared tags**. MUST **not** dump related links into the Thoughts body. |
| **`reader_landing`** | Optional **Dig deeper** block after Grounding: `why_it_matters` (reader world + interest), `takeaways` (2–4 carry-forward lines), `explore` (related paths, titled sources, prefilled Perplexity questions, optional completed thread URLs). Rendered by `layouts/partials/reader-landing.html`. Legacy freeform **`research`** still works when `reader_landing` is absent. |
| **`image_credit`** | Optional **Markdown** for attribution. On the **detail** page it appears in **post metadata** (with date / word count / time, above the hero image). On **section list** rows it appears **under the thumbnail** (before tags). Requires a featured image in both cases. |

**Relationship:** Grounding **supports** the Claim. The Claim says *what you are arguing in plain language*; Grounding says *what ideas or citations that maps onto*.

**Detail page order** (`layouts/claims/single.html`): **Post meta** (title block, then optional **`image_credit`** when there is a hero image) → **Hero image** (if any) → **Claim** → **Thoughts** (body) → **Grounding** → optional **Dig deeper** (`reader_landing` or legacy `research`) → **piece footer** (updated date, tags/shares when applicable, prev/next) → optional **Related** shelf (`related` or shared-tag fallback via `layouts/partials/related-keep-reading.html`; RELATED band above the cards) → **Back | Home**.

**Thoughts subsections (`###`):** Use **plain titles** with **no leading emoji**. Prefer **hook-style** titles (tension, question, punch), not seminar labels only; see **`.cursor/skills/site-revise-hooks/SKILL.md`** → **Body headings**. The template already adds icons for the main bands (**Claim**, **Thoughts**, **Grounding**, **Dig deeper**); emoji on every body subsection tends to look busy. Subsections still appear in **Contents** nested under **Thoughts**.

**No opening `##` in the body (MUST):** The single-page template renders a **Thoughts** band title above the body markdown. MUST **not** start the body with a top-level **`##`** heading (article title, chapter label, or repeat of the Hugo **`title`**). That double-heading clashes in the UI (two competing section titles). Open with **prose** (lead paragraph), then the first structural heading MUST be **`###`** when you need a hook. Mid-essay **`##`** bands are rare; prefer **`###`** throughout so all body headings nest under **Thoughts** in **Contents**. Put the essay hook in front matter **`title`** and the thesis in **`description`** (Claim), not a body **`##`**.

## Section, categories, and tags (Hugo)

- **Section** comes from the **content path** (for example `content/social-protocols/...` → section `social-protocols`). It is not set by `categories`.
- **`categories`** are a **taxonomy**: optional labels for grouping and taxonomy pages. They do not have to mirror the folder name.
- If **`categories`** is omitted, list partials in this project can fall back to the **section** name for display; if set, the **first** category is often what themes show in list meta—check `layouts/partials/seven-style-row.html` when in doubt.
- Hugo **categories** are **flat** (no built-in parent/child). Use **tags** for hooks and subtopics (see **Tag voice** below).
- **Recommendation:** one **primary category** aligned with the main site area (for example `Reality-Protocols` under `content/social-protocols/`, `Human-Condition` under `content/human-condition/`, or `X-Minds` under `content/x-minds/`), and **several tags** for topic detail. Add a second category only if you need another taxonomy axis and accept how the theme surfaces the first term.

### Human-Condition theme hubs (same idea as Cognitive-Memetics “subcats”)

**Cognitive-Memetics** uses **two** `categories` terms: an **umbrella** plus **exactly one project hub** (**`Cube-Cows`**, **`Raymond`**, **`Por-Estas-Calles`**, **`T-Shirt Art`**, **`Reptilocracy`**, **`Pawtropolis-Under-Fire`**). See **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`** → **Front matter conventions** → **`categories`**.

**Human-Condition** can use the **same pattern**: taxonomy is still **flat**, but you treat the second term as a **theme hub** (not a Hugo parent category).

**Publish days (site calendar):**

- **`content/x-minds/`:** new posts MUST use a **Thursday** calendar day in the bundle folder name and in **`date`**.
- **Tuesday** is **freestyle**: any non-reserved section may publish that day (including **`human-condition/`**, **`mind-infrastructure/`**, **`social-protocols/`**, and video picks there). Prefer Tuesday when the author asks for a mid-week freestyle slot.
- **`content/human-condition/`:** no fixed weekday lock. Older Thursday- or Friday-dated episodes stay as-is unless the author asks to reschedule.

- **Umbrella:** always **`Human-Condition`** first (so list meta and habits stay consistent).
- **Theme (pick one):** add a **second** term only when the post clearly belongs under that theme, for example **`Mental-Processes`** (choice, mental models, development of empathy or morals), **`Social-Protocols`** (reciprocity rules, iterated strategies, norms of response), **`Social-Behaviour`** (what people do together in the wild), **`Cooperation`** (helping, collective outcomes), **`Social-Trust`** (expectations, reputation, repair), **`Dark-Triad`**, **`Present-Moment`**. Do **not** stack several theme hubs on one post; use **tags** for extra angles. **`Social-Protocols`** is **not** **`Reality-Protocols`** (different section and hub: shared belief and large-scale coordination). Mixed-wiring community posts (parents, school, finding others like you) belong in **`content/x-minds/`**, not as a Human-Condition theme.
- **Hub pages:** when you introduce or reuse a theme term, add or maintain **`content/categories/<slug>/_index.md`** (and **`_index.es.md`** if the site localizes that hub) so `/categories/<slug>/` reads like the Cognitive-Memetics project pages.

Example:

```yaml
categories: ["Human-Condition", "Social-Protocols"]
```

### X-Minds (`content/x-minds/`)

Public section **X-Minds** (`/x-minds/`): mixed-wiring lives, parents, and community. Usual types are **`claims`** and **`video`**.

- **Umbrella:** always **`X-Minds`** first.
- **Publish day:** MUST use a **Thursday** calendar day in the bundle folder name and in **`date`** (see **Publish days** above).
- MUST NOT use **neurodivergence** (or clinical “off the typical” labels) in **`title`**, **`description`**, or hub `_index` copy. Body MAY name a diagnosis when a source or a person uses that word; prefer **mixed wiring** in the author’s frame.
- Hub page: **`content/categories/x-minds/_index.md`** (and **`_index.es.md`**).

Example:

```yaml
categories: ["X-Minds"]
```

### Tag register (site-wide inventory)

- MUST read **`data/tag-register.txt`** before assigning tags (it updates on **`make tag-register`** and before every **`make build`**). Follow **`.cursor/skills/site-tag-register/SKILL.md`**: prefer an existing tag when it fits; new tags are allowed when none do. For consolidating duplicates, use **`.cursor/skills/site-tag-unify/SKILL.md`**.

### Tag voice (punchy, hashtag-like)

- **Purpose:** Tags carry **attitude, provocation, and memory**. They are **not** a second abstract: do **not** lift jargon from **Grounding** (for example paper constructs) into tags unless you want that label on purpose.
- **Prefer:** questions, tribal / social stakes, irony about fads, time and attention (`SameButDifferentWorlds`, `DoYouHaveATribe`, `GoBackToOfficeFad`, `TimeThatCounts`, `RealityCheck`, `KnowYourPolitics`). Do not stack two tags that ask the same question (for example drop one of two near-duplicate tribe hooks).
- **Draft in your head as hashtags**, then put them in front matter **without** the `#` character (YAML and Hugo taxonomy terms work better as plain strings; the theme shows them as tag titles).
- **Shape:** Compact **PascalCase** fused phrases; readable aloud like a bumper sticker.
- **Count:** Aim for roughly **three to five** tags per post; quality over volume.
- **No repetition:** Each tag must add a **different** angle (provocation, target, or question). Do **not** use two tags that restate the same idea in different words, and do **not** repeat a hook you already used on another post unless you mean a running theme. If two tags feel like one joke twice, delete one.
- **Generalizable vs post-specific:** Aim for a **mix**. One or two tags can name the **essay’s spine** (for example layers of empathy, outcome vs guilt) if you want list readers to see that hook; the rest should skew toward **themes that could tag a future post** (systems that steer behaviour, moral definition, variation between people). **Title** and **Claim** already carry the spine; tags do not have to repeat them. When choosing between two phrases, prefer the one that still reads true if you strip this article’s proper nouns.
- **Reusability:** Prefer tags that can **show up again** when another post hits the same theme (a recurring series or the same fault line in the wild). Avoid one-off in-jokes unless that is the point. If a phrase only fits a single paragraph in a single article, it is usually too narrow for a tag.
- **Example** (savings, sovereign debt, moral cost of funding):

  ```yaml
  tags: ["RealityCheck", "KnowYourPolitics", "ThinkHardThisTime"]
  ```

- **Example** (sync, tribe, real time vs performative togetherness):

  ```yaml
  tags: ["SameButDifferentWorlds", "DoYouHaveATribe", "GoBackToOfficeFad", "TimeThatCounts"]
  ```

- **Example** (culture, identity, illusion as infrastructure; tags reusable across posts):

  ```yaml
  tags: ["SharedIllusionsRunTheWorld", "IdeasBecomeIdentity", "TribeBeforeTruth"]
  ```

- **Example** (empathy stack, moral development, systems aiming the beam; mix of reusable themes):

  ```yaml
  tags: ["SystemsCanHijackYou", "NotEveryoneFeelsTheSame", "EmpathyIsAGreyArea", "LetsDefineBad"]
  ```

## Grounding (how to write it)

- **Length (UI):** On the detail layout, **Grounding** often sits **beside** **Claim** in a two-column band. Long text truncates or crowds the column. MUST keep Grounding to **one or two short sentences** (plus an optional **Source:** line with the link). Put interpretation, stakes, and “why this matters for the essay” in **Thoughts**, not in Grounding.
- **Source link format (MUST):** After **`Source:`**, use a **Markdown link** whose **anchor text** is the **page or paper title** (or a very short equivalent if the official title is unusably long). Readers and the UI should see a readable label, for example `Source: [Regulating emotion through distancing (PMC)](https://www.ncbi.nlm.nih.gov/pmc/articles/PMC1234567/)`. MAY append a short hint in parentheses when it helps: **`(PMC)`**, **`(PDF)`**, **`(preprint)`**, journal or venue in a few words. MUST **not** paste the **raw URL** as the visible citation (no `https://` URL as plain text in Grounding).
- **How to pick anchor text:** Prefer the **HTML `<title>`**, PDF metadata title, journal article headline, or book or chapter title. If several sources, one **Source:** line per link or a short comma-separated list of titled links is fine; keep total Grounding length rules.
- **Density:** A **digest**, not a mini-essay: enough to name the **construct** or paper and **one** anchor (e.g. year, mechanism, or method), with key terms in **bold** (sparingly per **`.cursor/skills/site-revise-emphasis/SKILL.md`**), then point to a primary source or stable overview link.
- **Not the same as** an optional **quote exhibit** in the body: Grounding is the **digest**; a blockquote under a small heading is the **exhibit** (exact words). Grounding can be broader than one quote; the quote is optional and narrow.
- **Etymology vs theory:** For coinages (e.g. *meme*), etymology + year + overview link fits. For theory papers, state the **construct** (e.g. shared reality, epistemic companions), what the cited work does, then **`Source:`** with a **titled** Markdown link (mark `(PDF)` when it is a PDF).

## Optional body convention: primary-source quote (“What he said”)

- Use a **`###`** subsection with a **plain title**, then the blockquote. Example: `### What he said`. That gives a heading id, **Contents** entry (nested under **Thoughts**), and anchors consistent with the rest of the outline. Do **not** use a plain **`[WHAT HE SAID]`** line; that was legacy markup, not a project convention.
- Use when a **long primary-source quotation** earns its space. Omit if Grounding + prose are enough.

## Prose style (project preference)

- While drafting, MUST follow **`.cursor/rules/site-content-markdown-writing.mdc`** → **Voice while drafting** and **Explanatory prose** (clear mechanism paragraphs, not revelation stacks; full audit: **revise-post** Step 2).
- **Avoid em dashes (`—`).** Prefer commas, semicolons, colons, or parentheses.
- **English only** for site content (paths, front matter strings, body copy).

## Author-supplied copy (preserve voice)

- When the author already supplied prose in **`description`**, **`grounding`**, or the **body**, MUST **keep their wording** unless they **explicitly** ask for a rewrite, revise, tighten, or similar.
- MAY apply **mechanical** fixes only: YAML safety (quoting, colons), required front matter fields, **`###`** subsection titles, blockquotes, list markup, **`categories`** / **`tags`** / **`date`**, featured image paths, obvious structure mistakes, and turning a bare **`Source:`** URL into **`[title](url)`** when the **title** is known from the page, PDF metadata, or the author (MUST **not** invent anchor text with no basis).
- MAY add Markdown **`**bold**`** around existing words or short phrases **for emphasis** when that matches the author’s preference. MUST NOT pad with filler, swap in a more **packaged** or **marketing** tone, or replace a short direct line with a longer “elevated” paraphrase (for example expanding “not from ancient sources” into a grander disclaimer).
- **`description`** must still function as a **Claim** on cards. If the author’s paste is not a single stand-alone claim, MAY split fields or add minimal framing **only** with their agreement, or ask for one clarifying line; MUST NOT invent a new thesis they did not write.
- **Layout note:** The template shows **Claim** before **Thoughts** on the single page. If the author’s narrative order differs, explain that constraint; MUST NOT reorder ideas by rewriting their copy.

## Authoring workflow

1. Draft **`description`** (Claim) so it reads well alone on a card. **Sentence-by-sentence:** each line must pass **Claim fog** in **revise-hooks** (picture direction + outcome or a concrete noun; not “control systems above local detail” without a scene). If the author supplied Claim text, follow **Author-supplied copy (preserve voice)** above; MAY still flag fog and propose a plainer clause when they ask for revise/tighten.
2. Draft **`grounding`** as a **tight** digest (authors or construct, minimal jargon, **`Source:`** line with **`[title](url)`** Markdown, never a bare URL string). Expand nuance only in **Thoughts** if needed.
3. Write the body (metaphor, examples, optional quote exhibit with a `###` heading before the blockquote).
4. Set **`tags`** (register + tag voice). Then **recommend `related`**: search for up to **two** **cognitive-memetics** sayings or panels that share the mechanism or a punchy parallel; MAY add **one** other claim or video. Write Hugo paths into **`related`**. If none fit, leave empty and rely on tag fallback. MUST **not** paste related links into the Thoughts body.
5. Re-read Claim, then body as **Thoughts**, then Grounding for alignment (same order as the live page). Shorten Grounding if it repeats **Thoughts** or reads like a second Claim.
6. On **new bundle**, **reschedule** (**`date`** or folder rename), or **`draft`** change: run **`make calendar`** (see **`.cursor/rules/site-content-markdown-writing.mdc`** → **Publish calendar**).

## References in this repo

- Archetype: `archetypes/claims.md` (TOML; content bundles may use YAML with the same fields; legacy **`paper`** may still work for grounding-style content).
- Related strip: `layouts/partials/related-keep-reading.html` (curated **`related`** or shared-tag fallback).
- Examples: `content/social-protocols/2026-04-01-memes/index.md` (coinage + quote exhibit), `content/social-protocols/2026-04-03-shared-reality/index.md` (theory + PDF source), `content/human-condition/2026-03-27-born-to-choose/index.md` (short Grounding + PMC source), `content/human-condition/2026-04-10-empathy-levels/index.md` (developmental moral empathy + layered tags).
