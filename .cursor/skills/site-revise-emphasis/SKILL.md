---
name: site-revise-emphasis
description: >-
  Sets default Markdown **bold** for site content under content/: bold **ideas** (constructs,
  mechanisms, story lines, metaphors), not lone verbs or filler. Restrained density (a few hooks per
  block), not wall-to-wall bold. Use when authoring or editing any post type, when the user mentions
  bold for ideas, negrita en ideas, too much bold, scan-style emphasis, list cards, fixing emphasis
  density, or revise-emphasis. At ship time, emphasis is audited with punctuation in revise-format
  (revise-format skill).
---

# Revise emphasis (Markdown bold)

**Ship-time sweep:** run **`.cursor/skills/site-revise-format/SKILL.md`** (emphasis audit + em dash search). Use **this skill** while **authoring** or when only bold rules are needed.

## Scope

- Applies to **any** Hugo content under `content/` where **`**bold**`** appears or would be added: **`type: claims`**, **`video`**, **`panel`**, **`sayings`**, and markdown **body** (Thoughts, essays, optional sections).
- Bold is **styling**. It must not change meaning, add claims, or replace clear prose with shouted keywords.

## Default: restrained emphasis (MUST)

Goal: **guide the eye** to the twist, risk, or anchor term. If almost every content word is bold, nothing wins.

## Bold for ideas (MUST)

**Primary rule:** use **`**bold**`** to mark **ideas** the reader should carry (mechanism, construct, metaphor, story line, stake), not to decorate every strong word.

| Bold the **idea** | Do **not** bold |
|-------------------|-----------------|
| Named constructs (**perceived control**, **control percibido**, **pre-set defaults**) | Bare verbs (**se aprovechan**, *hack*, *remove*) unless the verb **is** the coined idea |
| Story lines and agency hooks (**story of agency**, **I am still in control**, **sigo mandando yo**) | Articles, conjunctions, light filler |
| Mechanism metaphors (**brakes**, **frenos**, **numbness**, **anestesia**) | Every noun in a sentence |
| Turns and payoffs (**relief from anxiety**, **estrés**, **handing the choice off**) | The same stressed word in **title**, **Claim**, and the next sentence unless repetition is deliberate |

- Each bold span must pass the **standalone concept test** below: read only the bold text; it must still name an **idea**, not a grammar fragment.
- In bilingual bundles, bold the **same role** in EN and ES (construct ↔ construct, metaphor ↔ metaphor). Do **not** bold different concepts on each side.
- When the author says **bold for ideas** or **negrita en ideas**, apply this section; still keep **restrained** density (most sentences stay plain).

**Example (claims, EN / ES aligned):** Claim stresses the payoff (**relief from anxiety** / **estrés**); Grounding anchors **perceived control** / **control percibido** and **pre-set defaults** / **lo que ya viene marcado**; Thoughts bold **story of agency**, **brakes** / **frenos**, **numbness** / **anestesia**, not every verb in the psychopathy paragraph.

### Global rules

- Prefer **a few** bold spans per short block (rough guide: about **2–5** for a one-paragraph card field; scale up slightly for long paragraphs, but keep most sentences plain).
- Bold **ideas** first; MAY also bold **contrast / turn** words (*only*, *but*, *not*) when they complete the idea—not instead of it.
- **Standalone concept test (MUST):** If you bold only the surrounding words and read the bold span alone, it must still mean something (*flat ideals*, **nested control**, **the people who eat the cost**). **Fail:** lone *layers*, *rank*, *near money*, *delete button* without their frame. Prefer a **short phrase** (2–5 words) or a **full punch sentence** (*Perception is layered.*).
- Do **not** bold articles, conjunctions, and light filler (*the*, *and*, *each*) unless the author deliberately stresses them.
- Do **not** repeat the same stressed words in **`title`**, **`description`**, and the next line unless repetition is the point (card stacks get noisy).

### By type

| Context | Notes |
|---------|--------|
| **`type: claims`** **`description` (Claim)** | One or two **ideas** (thesis turn or stake); see **Bold for ideas**. Must read as a sentence, not a tag cloud. |
| **`type: claims`** **`grounding`** | Construct names and one mechanism anchor; see **Bold for ideas** and **`.cursor/skills/site-claims-content/SKILL.md`** for digest length. |
| **`type: claims`** **body (Thoughts)** | Bold **ideas** per section (delegation, agency story, metaphor, consequence); each span passes **standalone concept test**. Use blockquotes for full lines that already work alone. |
| **`type: video`** **`description`** (lead) | Light; lead reads as prose above the embed. |
| **`type: panel`** **`description`** (Teaser) | Same restraint as other teasers. |
| **`type: sayings`** **`description`**, **`tldr`**, **`fluff`** | Follow **Sayings** below and **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`** (Sayings card teaser for **`description`** derives only from **`title`**, **`tldr`**, **`fluff`**). |

### Sayings (`type: sayings`)

| Block | Typical bold spans (not a hard cap) |
|-------|-------------------------------------|
| **`description`** (teaser) | About **2–4** short spans. |
| **`tldr`** | About **2–4** spans (often one on the lead hook, one on the punch). |
| **`fluff`** | About **3–6** spans if the block is long; keep most sentences plain. |

**Order of work:** Draft **`tldr`** / **`fluff`** (or preserve author wording) → add minimal bold → write **`description`** per Sayings card teaser → light bold on the teaser if helpful.

**Legacy:** Older **`content/cognitive-memetics/`** sayings may use **heavier** bold. Do **not** rewrite for density unless the user asks. For **new** posts, default **restrained**. Use **dense** “scan” style only when the user explicitly asks to match an old example.

## Editing existing posts (MUST)

- When the user asks to **fix bold**, **bold for ideas**, **reduce emphasis**, or **match site style**, MAY **add, move, strip, or trim** **`**bold**`** on existing words **without** treating that as a prose rewrite (markup only), as long as each span passes **Bold for ideas**.
- SHOULD trim **wall-to-wall** bold toward this skill’s default for the relevant type.
- When a bundle has **`index.md`** and **`index.es.md`**, after emphasis edits MUST re-check **idea parity** (same roles bolded in each language).

## English spelling (default-language `content/`)

- MUST use **US English** (**generalize**, **behavior**). See **`.cursor/rules/site-content-markdown-writing.mdc`** → **Language**.

## Spanish (`*.es.md`)

- Same **Bold for ideas** and **restraint** as English: a few idea-level hooks per block, not saturation. Follow **`.cursor/skills/site-spanish-translation-content/SKILL.md`** for idiom; follow **this skill** for which **ideas** get bold and how many spans per block.

## Related

- **`.cursor/rules/site-content-markdown-writing.mdc`**
- **`.cursor/skills/site-claims-content/SKILL.md`**
- **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`**
