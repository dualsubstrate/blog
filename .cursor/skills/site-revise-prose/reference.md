# Local prose review: rubrics (evaluate only)

Adapted from n8n `prose_polish` quality criteria (`internal/evaluation/criteria/registry.go`) and site prose-review / revise-post Step 2. Each criterion scores **0–2**.

Process criteria (`instruction_compliance`, `feedback_adherence`) are **out of scope**: no rewrite loop, no preserve_markers contract.

**Single-draft note:** n8n’s compression-vs-voice delta needs two versions. For one draft, flag **compression smell** (over-tight, flattened personality, staccato same-length punches) under notes. Do not invent prior scores.

---

## Criteria (0–2 each)

### clarity_readability

**What:** Ideas are easy to understand; natural flow.

| Score | Meaning |
|-------|---------|
| 2 | Exceptionally clear; natural flow |
| 1 | Mostly clear; awkward structure or muddy phrases |
| 0 | Unclear or hard to follow |

Note: this is communication clarity, not grammar police.

### anti_fluff_compliance

**What:** Avoids abstract, preachy, flowery filler.

**Forbidden shapes:** abstract verbs (*demonstrates*, *illustrates*, *serves as*, *embodies*); abstract nouns (*power of*, *essence of*); meta (*What this teaches…*, *The lesson here is…*); **Metaphor-shell restack**; **Industry-verb shells** (`sells`, `steer`, `assembling`, `harden`, fake scene `the same room`).

| Score | Meaning |
|-------|---------|
| 2 | Direct and concrete |
| 1 | 1–2 fluff instances |
| 0 | 3+ instances |

### conversational_tone (bar test)

**What:** Would you say this to a friend at a bar? Straight talk, not a lesson or textbook frame.

| Score | Meaning |
|-------|---------|
| 2 | Passes bar test; no teaching frames |
| 1 | Mostly conversational; 1–2 formal/teaching phrases |
| 0 | Lesson, textbook, or academic tone |

Fails: *In English you'd say…*, *This concept acknowledges that…*, *When one has already navigated…*

### no_therapist_voice

**What:** No self-help / process jargon that analyzes social dynamics instead of describing action.

**Forbidden shapes:** *navigating*, *maintaining*, *reinforcing*; *validate*, *gentle nudge*, *social dynamics*, *maintain harmony*.

| Score | Meaning |
|-------|---------|
| 2 | None |
| 1 | 1–2 instances |
| 0 | 3+ instances |

### anti_ai_feel

**What:** Sounds like a person, not polished marketing or smooth empty AI prose.

**Forbidden framing:** *That's the vibe…*, *Think of it as…*, *Here's what that means:*, *In other words:* (as empty restatement). Prefer natural colloquial wording over stiff formal synonyms where the register fits the essay.

| Score | Meaning |
|-------|---------|
| 2 | Real person talking |
| 1 | Mostly human; 1–2 over-polished tells |
| 0 | 3+ marketing / over-polished tells |

---

## AI artifact shortlist (quote when found)

Aligns with `.cursor/perplexity/packs/prose-review.md` and revise-post Step 2.

| Pattern | Example shape |
|---------|----------------|
| Teacher / editor framing | *You will learn…*, *The useful move is…* |
| Bridge verbs | *bridges X and Y*, *maps onto* |
| Semicolon aphorisms | *Stress drains; purpose concentrates.* |
| Negation-first fluff | *It is not X; it is Y*, *Not A but B* |
| Meta-commentary | *This makes the point clear*, *The pattern matters* |
| Caption pivots | *Picard is careful here.*, *Worth noting:* |
| Classic AI rhetoric | *This changes how we understand…* |
| Revelation stacks | Several one-line poetic restatements of one idea |
| Rhetorical noun fragments | *The mirror.* / *The wound.* alone |
| Throat-clearing | Soft restatement of the heading with no new detail |
| Hedge piles | *One possibility is…*, *It is important to note…* |
| Staccato thesis stacks | Same-length punch lines every sentence |
| Compression smell | Cuts that kill cadence or voice for tightness alone |
| Metaphor-shell restack | After a mechanism sentence: "The blend is your experienced reality…" then "If the self were a sealed file, you would have no steering wheel." Do not Keep the tautology as punchy. |
| Industry-verb shells | "sells the blend as 'me'"; "steer the rewrite"; "keeps assembling it"; "next assembly"; "the same room"; "template can harden"; "name the box… and the box." Fail even as a single clause. |

**Standalone worth:** For each flagged sentence: *If this stood alone, would it be worth saying?* Flag restatements and decoration. A new metaphor or coined label for the same claim is not a fact, mechanism, or turn.

**Sequence test:** After sentence N, N+1 fails if it only relabels N, tautologizes N, or restates N's implication with a new metaphor pair. Do not Keep those quotes. Owner: **`.cursor/skills/site-revise-post/reference.md`** → **Metaphor-shell restack**.

**Clause test:** A clause fails even alone if its payload is a sales, factory, driving, or hardening verb, or a fake scene. Owner: **`.cursor/skills/site-revise-post/reference.md`** → **Industry-verb shells**.

---

## What is out of scope

- Factual accuracy, citations, grounding quality
- SEO, curiosity-title craft (use revise-hooks / revise-score)
- Spanish body (`index.es.md`)
- Hugo structure compliance beyond prose readability
