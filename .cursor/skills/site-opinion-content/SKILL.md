---
name: site-opinion-content
description: >-
  Authors and edits Hugo posts with type opinion: editorial position (description),
  full essay body, optional grounding, optional reader_landing, and normal
  categories/tags. Use when adding stanceful essays under social-protocols,
  human-condition, mind-infrastructure, or x-minds without Claim/Grounding
  semantics, when the user mentions type opinion or opinion archetype, or when
  mapping essay MCP content_form opinion_essay to Hugo.
---

# Opinion content type (`type: opinion`)

## What this type is for

**Opinion** is a **stanceful essay** in the author's voice. It is **not** a short cognitive-memetics saying and **not** a **`claims`** page where **`description`** is framed as an evidence-backed **Claim** with mandatory **Grounding** support.

**UI (this repo):** `layouts/opinion/single.html` uses **Opinion**, **Essay**, optional **Grounding**, and optional **Dig deeper** (`reader_landing` or legacy `research`). List rows (`seven-style-row.html`) show **Opinion** plus optional **Grounding** and **Key points** from `reader_landing.why_it_matters`, mirroring claims list behavior without Claim labels.

**Section vs type:** Pick the **folder** by main job (see **`.cursor/rules/site-content-placement.mdc`**). Set **`type: opinion`** when the publish shape is a full essay with an editorial position, not Claim + mandatory Grounding.

| Section folder | When opinion fits |
|----------------|-------------------|
| **`social-protocols/`** | Rooms, norms, blame, accountability, coordination, institutions, media protocols. |
| **`human-condition/`** | Person-level psychology, empathy, morality, identity (not X-Minds community lane). |
| **`mind-infrastructure/`** | General models and tools when the essay's main job is not better served as **`type: video`**. |
| **`x-minds/`** | Mixed-wiring lives and community when the piece is essay-shaped, not a short meme line. |

**MUST NOT** use **`cognitive-memetics/`** for full opinion essays; use **`content_form=cognitive_memetics`** only for short sayings or panel strips.

## Opinion vs claims

| | **`type: opinion`** | **`type: claims`** |
|--|---------------------|-------------------|
| **`description`** | Editorial **position** (card hook + thesis in your voice). | **Claim** (what the reader should believe). |
| **Body** | **Essay** (mechanism, scenes, argument). | **Thoughts** (same mechanical role, different label). |
| **`grounding`** | **Optional** citations or constructs when they help; not the spine of the type. | **Expected** support tying the Claim to sources or standard terms. |
| **Pipeline** | **`content_form=opinion_essay`** (alias **`opinion`**) in **`data/essay-content-forms.yaml`**. | **`content_form=claims`** (CLAIM / THOUGHTS / GROUNDING bands). |

## Field roles

| Field | Role |
|-------|------|
| **`description`** | **Opinion**: the stance in plain language. Cold-read on cards; follow **`.cursor/skills/site-revise-hooks/SKILL.md`**. |
| **Body** | **Essay**: 600–800 word stanceful argument when generated via MCP; hook sections use **`###`** under the template **Essay** band. |
| **`grounding`** | Optional short support (definitions, **Source** lines with `[title](url)`). Omit when the essay stands on scenes and logic alone. |
| **`reader_landing`** | Optional **Dig deeper**: `why_it_matters`, `takeaways`, `explore`. See **`layouts/partials/reader-landing.html`**. |
| **`related`** | Optional keep-reading paths; same rules as claims (**`.cursor/skills/site-claims-content/SKILL.md`** → **`related`**). |
| **`image_credit`** | Same as claims: detail meta and list row under thumbnail when a featured image exists. |

**No opening `##` in the body (MUST):** The single template renders an **Essay** band above the body. Open with prose or **`###`** hooks; do not duplicate **`title`** with a body **`##`**.

**`date`:** MUST use **`.cursor/rules/site-content-markdown-writing.mdc`** → **Publish `date`**.

## Categories and tags

- Prefer one primary **`categories`** term aligned with the section (for example **`Reality-Protocols`** under **`social-protocols/`**, **`Human-Condition`** under **`human-condition/`**).
- Tags: **`.cursor/skills/site-tag-register/SKILL.md`** (punchy, reusable; not a second abstract).

## Essay MCP (`content_form=opinion_essay`)

Host pack: **`data/essay-content-forms.yaml`** → **`forms.opinion_essay`** (aliases **`opinion`**).

- **Generation:** Same long-essay contract as **`essay_argument`** (three essay blocks, 600–800 word band, required **`reader_landing`** after the diagnostic close).
- **Hugo export (manual v1):** Set front matter **`type: opinion`**. Map **`description`** from the editorial position; body from polished essay markdown; copy **`reader_landing`** from pipeline output; add **`grounding`** only when the piece needs explicit citations.
- **CLI / MCP:** Pass **`--content-form opinion_essay`** (or set WIP / create flags your operator UI exposes) when starting composition for this shape. **`content_form`** controls generation; Hugo **`type`** controls rendering.

**Example mapping (blank mirror / outward blame):** Publish under **`content/social-protocols/`** with **`type: opinion`** when the spine is rooms, silence, and accountability. Use **`human-condition/`** only if the essay is reframed primarily as person-level theory.

## Editorial passes

- Voice and commute test: **`.cursor/rules/site-content-markdown-writing.mdc`**
- Full revise conductor: **`.cursor/skills/site-revise-post/SKILL.md`**
- Ship-time em dash / bold: **`.cursor/skills/site-revise-format/SKILL.md`**
