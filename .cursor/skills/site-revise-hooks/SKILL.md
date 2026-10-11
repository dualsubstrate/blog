---
name: site-revise-hooks
description: >-
  Shapes Hugo front matter so links and list rows pull readers in: curiosity, tension,
  and stakes, with preference for direct active wording and plain mechanism (not
  corporate headline cadence). Covers title vs description, coordination with claims
  and video types, tags as attitude, and forbidden clickbait patterns. Use when
  rewriting titles, card copy, teasers, list-view hooks, emotional pull, direct hooks,
  active voice, body ## hook headings, capture vs pass cold-read, revise-hooks, or
  “make people want to click” for site content; apply after the relevant type skill
  (claims-content, video-content, cognitive-memetics-content).
---

# Revise hooks (titles, list copy, and body headings)

## Purpose

Increase **pull** on links, list rows, and in-page scans: copy should trigger **curiosity** (open loop, tension, stakes) and still **deliver** on the next screen. This skill governs **hook psychology**, **capture vs pass**, **body headings**, and **cold-read** tests. Type-specific field rules live in **`.cursor/skills/site-claims-content/SKILL.md`**, **`.cursor/skills/site-video-content/SKILL.md`**, and **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`**.

## Order of operations

1. Open the **content type** skill for the page (`claims`, `video`, `panel`, `sayings`, etc.).
2. Apply **this** skill to **`title`** and any **teaser** fields the type uses, without breaking the type’s obligations (for example the Claim must stay a real assertion for `type: claims`).

**Full lot (flow + hooks + post in one go):** **`.cursor/skills/site-revise-post/SKILL.md`** (default mode **`rough`**). Use **this skill alone** when you only need list copy and heading hooks.

**Draft or retitle (curiosity gap, 3–5 candidates):** **`.cursor/skills/site-curiosity-title/SKILL.md`** first (Gemma 4 thinking on Polypus for candidates + cold-read), then **this skill** to audit the pick.

## How the theme uses fields (this repo)

- **List rows** (`layouts/partials/seven-style-row.html`): **`title`** is the main link line; **`description`** is the rich summary beside or below it (markdown).
- **`subtitle`** (`type: video`): optional second line under the title on the single page; use for a clarifying phrase if the title is very hook-forward.

There is no separate “card title” field in normal posts: **`title`** is the primary curiosity lever everywhere.

## Field split (default)

| Field | Hook job |
|-------|----------|
| **`title`** | **Pull**: tension, contrast, stakes, paradox, or a sharp question. MAY be **direct** (name the mechanism in plain language) or **suspenseful** (withhold one beat) as long as the page pays it off. Optimized for **scanning**. Still truthful and specific enough to stand alone. |
| **`description`** | **Promise and precision** (type-dependent): what the page actually gives. For **`type: claims`**, **`description` is the Claim** — see **Claims** below; do not replace it with pure hype. |
| **`tags`** | **Attitude and theme** (PascalCase hooks); see claims-content *Tag voice*. They complement the title, not repeat it. Check **`data/tag-register.txt`** and **`.cursor/skills/site-tag-register/SKILL.md`** before inventing new tag strings. |

## Curiosity patterns (use what the piece supports)

- **Prediction error:** imply the reader’s default story is incomplete (“when X looks like Y until Z”).
- **Tension:** two motives or forces in conflict (comfort vs clarity, speed vs truth).
- **Stakes:** what gets worse if ignored, or what becomes possible — only if the body pays this off.
- **Specific oddity:** one concrete noun, number, or scene from the piece in the title (avoids generic labels).

## Directness and voice (titles)

**Direct** titles state **who does what** or **what follows from what** without hiding behind vague setup. They often work better than faux-suspense (“what X does to Y”) when the goal is **clarity plus punch**.

- **Prefer active voice:** named subjects with strong verbs (*feeds reward*, *memes jump*, *you sync*). SHOULD NOT lean on stative filler (*you are not on the same…*) unless brevity wins; if you use *be*, pair it with a sharp predicate.
- **Say the mechanism:** name the force (feed, meme, algorithm, time together) and the outcome in plain words the essay supports.
- **Keep parallel rhythm** when stacking clauses (*same X, different Y*; *hard times, strong men, and…*) so list rows scan fast.

### MUST NOT (tone)

- **Corporate / management headline cadence:** empty multipliers (*unlock*, *elevate*, *drive outcomes*), HR-tech abstractions (*engagement*, *performers*, *systems rot* as a slogan), or strategy-deck cause chains that sound like LinkedIn.
- **Fake sophistication:** long nominal strings where a short verb would do.

Suspense and directness can coexist: a **short paradox** (*same words, different worlds*) or a **clear if-then** (*feeds reward noise over skill*) both beat a title that sounds strategic but says little.

## Capture vs pass (list-facing copy)

**Objective:** List-facing fields decide whether you **capture** the reader or they **pass**. If a line does not make sense in a few seconds, they scroll on; the body never loads.

| `type` | Fields that MUST cold-read (no body required) |
|--------|-----------------------------------------------|
| **`claims`** | **`title`**, **`description`** (Claim), **`grounding`** (short digest on many list rows) |
| **`video`** | **`title`**, **`description`**, **`sowhat`** when set (see **`.cursor/skills/site-video-content/SKILL.md`** → **Cold-read gate** for bullet rules) |
| **`panel`**, **`sayings`** | **`title`**, **`description`** (teaser) |

**Success:** A stranger understands each line and feels a reason to open. **Failure:** They need the essay, video, or private metaphors to decode the line.

**Voice (MUST):** List-facing fields are also subject to **`.cursor/skills/site-revise-post/SKILL.md`** → **Step 2** (AI voice) and **Step 5** (negation fluff). Cold-read alone is not enough if the line sounds polished but hollow.

## Body headings (`##` and `###`) — hooks, not labels (MUST)

Applies to **any** post with markdown body sections: **`type: video`** TLDR, **`type: claims`** **Thoughts**, long **`panel`** copy, etc. Headings appear in **TOC** and skims; they must **pull**, not only **name the topic**.

### Objective

Readers who opened the page still **scan**. A label like `## Ideas that travel`, `## The Takeaway`, or `## Tools shape the metaphor` states the theme; it does not give a reason to read that block. Prefer a **hook**: tension, a question, a turn, or a line someone would say out loud.

### MUST

- Write section headings as **hooks** (stakes, question, punch, concrete image), in **active, spoken** voice.
- Match each heading to that section’s **spine line** (the payoff), not the abstract theme name.
- Use **`###`** for **`type: video`** TLDR body sections (and **`### Chapter Guide`**): **`layouts/video/single.html`** already renders **`h1`** title, **`h2`** subtitle, and **`h3`** lead/payoff blocks; body **`##`** would duplicate **`h2`** in the outline. See **`.cursor/skills/site-video-content/SKILL.md`** → **Body headings**.
- Use **`###`** for **Thoughts** subsections on **`type: claims`** (plain text, **no leading emoji** per **`.cursor/skills/site-claims-content/SKILL.md`**). Some older **`claims`** essays still use **`##`**; new claims copy MUST open with prose, not **`##`** (see below).
- MAY quote words the section turns on (for example `### When "like" turns into "is"` on video TLDRs).

### MUST NOT

- Default to **seminar labels** when a hook would work (`## Name the need first`, `## One-shot vs repeated games` as a bare topic label).
- Use **mystery-box** headings that withhold what the section is about.
- Repeat **`title`** or list-facing teaser text verbatim as a heading.

### Neutral exceptions (keep plain)

- **`### Chapter Guide`** on **`type: video`** (navigation; Spanish: **`### Guía de capítulos`**).
- Rare **structural** labels only when the block is pure reference (use sparingly).

### `type: claims` — no opening body `##` (MUST)

The **Thoughts** template already shows a section title above the body. MUST **not** begin claims body markdown with **`##`** (including a repeat of the post **`title`** or a pseudo-article headline). **Fail** hook audit if the first body line is **`## …`**. Open with prose; use **`###`** for the first in-essay hook heading. See **`.cursor/skills/site-claims-content/SKILL.md`** → **No opening `##` in the body**.

### Do / don't (this repo)

| Don't (label) | Do (hook) |
|---------------|-----------|
| `### Ideas that travel` | `### When "like" turns into "is"` |
| `### The brain balances accuracy with other drives` | `### Why corrections struggle` |
| `### Simplification buys one dimension` | `### You simplify, then you forget the cut` |
| `### The Takeaway` | `### What tit-for-tat actually rewards` (example; match the piece) |

Reference: **`type: video`** hook **`###`** lines (for example **`content/social-protocols/2026-06-09-the-prediction-business/`**). Older video bundles may still use **`##`** until touched.

### Cold-read (light)

A heading MAY use a metaphor the **next paragraph** explains. It MUST NOT require a **later section** to decode. If the heading fails alone, rewrite it.

### Standalone worth (body copy)

After headings pass, apply the same test to sentences and paragraphs under **`.cursor/skills/site-revise-post/SKILL.md`** → **Step 2** (standalone worth): cut lines that only decorate or restate the hook without new concrete detail.

## Cold-read (naive reader test)

Before publish or before marking **`.cursor/skills/site-revise-post/SKILL.md`** Step 1 / Step 3 **Pass**:

1. Read each **list-facing** field alone (see table above).
2. Read each **body section heading** alone (light cold-read).
3. Ask: *Would someone with no context know what this claims and why it matters?*
4. **Fail** if you need the body, video, or in-joke metaphors to decode the line (`brain shorthand`, `yesterday's kit`, bare *like*/*is* without setup).
5. **Pass** if a stranger could quote back the point (named subject + optional example + punch).

**`type: video`:** MUST also follow **`.cursor/skills/site-video-content/SKILL.md`** → **Cold-read gate** (**first two `description` bullets** get extra scrutiny).

## Claims posts (`type: claims`)

- **`description`** MUST remain the **Claim** (narrative assertion). See **`.cursor/skills/site-claims-content/SKILL.md`**.
- Put **pull in `title`** (and tags): the Claim stays precise; the title can be **direct** (thesis in plain language) or **paradox / tension** without changing the logical content of the Claim.
- MUST NOT rewrite the Claim into vague curiosity (“You won’t believe…”) or empty questions.

### Claim fog (abstract cold-read fail) — MUST for `description`

**Objective:** Every **sentence** in the Claim must cold-read on a **card skim** without the body. One foggy clause fails the Claim even when the rest is strong.

**How to test:** Read the Claim alone. For each sentence, ask: *Can a stranger picture who does what to what?* If not, **Fail** that sentence.

| Fail pattern | Example (bad) | Why it fails |
|--------------|---------------|--------------|
| **Dual-meaning jargon** | “control systems that sit above local detail” | Brain or company? “Systems” and “detail” float. |
| **Verb + abstract object** | “map complexity flat” | No image; negation does not fix fog. |
| **Stacked abstractions** | “route attention through nested control structures” | Nouns pile up; no scene. |
| **Metaphor without noun** | “sit above local detail” | Above *what*? Whose detail? |
| **Thesis hidden in jargon** | “predictive processing of social signals” | Construct name replaces the point. |

| Pass pattern | Example (good) | Why it passes |
|--------------|----------------|---------------|
| **Direction + plain outcome** | “route attention upward, so local noise gets summarized before it reaches a decision” | Up, noise, summary, decision. |
| **Concrete nouns** | “layers,” “rank,” “pyramid,” “VP,” “founder” | Reader sees the thing. |
| **One beat per idea** | layers → rank → upward summary → pyramids feel obvious | Card can follow. |

**Fix (keep the thesis):** Replace the foggy clause with **direction + outcome** (up/down, noise/signal, who decides) or one **concrete scene noun** from the body. MUST NOT swap in new jargon.

**Repo example:**

| | Text |
|---|------|
| Fail | …route attention through control systems that sit above local detail. |
| Pass | …route attention upward, so local noise gets summarized before it reaches a decision. |

**`grounding`:** Same fog test on the digest (shorter). Technical terms OK only with a **gloss** or **named construct** the sentence still explains in plain words.

## Video posts (`type: video`)

- Follow **`.cursor/skills/site-video-content/SKILL.md`** for embed, TLDR body, and **Chapter Guide**.
- **`title`** can carry the **watch impulse**; **`description`** lists concrete, honest hooks about what is non-obvious in the talk.
- **MUST** cold-read **`description`** and **`sowhat`** per **Capture vs pass** and **Cold-read** above; extended rules in **`.cursor/skills/site-video-content/SKILL.md`** → **Cold-read gate**.
- Body **`###`** headings: **Body headings** above.

## Cognitive-memetics (`panel`, `sayings`, hubs)

- Follow **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`**. Titles often pair **image + phrase**; keep voice consistent with that skill.
- For **`sayings`** **Teaser** (`description`), including **T-Shirt Art**, avoid **meta openers** (“T-Shirt Art piece,” “this post,” etc.); jump to hooks. Details under that skill’s **T-Shirt Art** and **`type: sayings`** sections.

## MUST NOT (integrity)

- Fake urgency (“last chance,” “everyone is talking about”) without substance.
- Mystery boxes: titles that **withhold** what the page is about so the reader feels tricked after opening.
- Superlatives with no anchor (“best,” “ultimate”) unless the body justifies them factually.
- Emotional manipulation on serious topics (fear/shame) unless the essay earns it and offers grounding.

## SHOULD

- After editing a title, read **`description`** (and type-specific body lead) once: **does the first screen fulfill the hook?**
- Keep URLs stable: change **`title`** for display; do **not** rename bundle folders or slugs unless the user wants new URLs.

## Quick checklist

- [ ] Title creates tension, paradox, or a **clear thesis**, not a bland label.
- [ ] Wording is **direct** where possible; verbs carry the load; no corporate-deck tone.
- [ ] List-facing fields (**description**, Claim, **`grounding`**, **`sowhat`**) **cold-read** without the body.
- [ ] **`type: claims`:** Claim has **no Claim fog** (see **Claim fog** above); each sentence pictures an outcome or names a concrete noun.
- [ ] Body **`###`** headings are **hooks**, not seminar labels (**`type: video`**, claims **Thoughts**; neutral **`### Chapter Guide`** / **`### Guía de capítulos`** on video).
- [ ] Description / Claim / teaser matches the title’s promise.
- [ ] Type skill constraints still satisfied.
- [ ] No forbidden patterns above (integrity or tone).
