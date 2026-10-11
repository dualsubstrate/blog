---
name: site-revise-post
description: >-
  Default Hugo revise conductor: modes rough, standard, polish, format-only,
  checklist. Onion workflow: write a lot plan first (optional list-scope Gemma
  4), then iterate phases until the merged apply list. Prompt to skip Gemma.
  Use for "revise this post", full publish pass, or step-only cold read. Tokens:
  with gemma, skip gemma, gemma4, plan only. Step 2 covers AI voice; Step 4
  runs revise-format. Focused cadence / hooks / format: revise-flow, revise-hooks,
  or revise-format.
---

# Revise post (conductor + filters)

## Purpose

One skill runs the **publish lot** in a fixed order. Child skills stay separate for focused passes.

This skill **orchestrates** phases and owns **Steps 1–5** (post filters). MUST read and follow linked skills for flow, hooks, and format; MUST NOT duplicate their rubrics.

## When to use

| User intent | Behavior |
|-------------|----------|
| Full lot / publish-ready / "revise this post" | **This skill**, default mode **`rough`** |
| Cadence only | `.cursor/skills/site-revise-flow/SKILL.md` |
| Hooks only | `.cursor/skills/site-revise-hooks/SKILL.md` |
| Site checklist only (no heavy/fine flow) | **This skill**, mode **`checklist`** |
| Em dash + emphasis only | `.cursor/skills/site-revise-format/SKILL.md` |
| Score gate | `.cursor/skills/site-revise-score/SKILL.md` (optional phase) |
| Spanish sibling | `.cursor/skills/site-revise-post-es/SKILL.md` (after English) |
| Cold read / accessibility / Step 3 only | Run **Step 3** only on the text they point at |

## Modes

| Mode | Phases | When |
|------|--------|------|
| **`rough`** (default) | Lot plan → Heavy flow → Hooks → Steps 1–3, 5 → Fine flow → **Format last** | New or messy drafts; revise-post without a mode |
| **`standard`** | Lot plan → Heavy flow → Hooks → Steps 1–5 → **Format last** | Draft already fairly tight; skip fine flow |
| **`polish`** | Lot plan → Hooks → Steps 1–3, 5 → Fine flow → **Format last** | Prose already cut; title/body mostly frozen |
| **`format-only`** | **Format only** (no plan Gemma) | Last sweep before ship |
| **`checklist`** | Lot plan → Steps 1–5 only | Site filters without heavy/fine flow; Step 1 still reads revise-hooks |

User MAY say: `rough`, `standard`, `polish`, `format-only`, or `checklist`.

## Onion: plan first, then iterate

Do **not** dump the whole body into Gemma before any editorial work. Peel layers:

1. **Big picture** (cheap): agent cold-read of title, list fields, headings, and body spine. If Gemma is on, run **`--scope list` only** (title, Claim/`description`, `grounding`, `sowhat` when present). A few units, not one call per body paragraph.
2. **Write the lot plan** in chat (ranked issues, which phases will fire, which body lines look risky). MUST NOT treat Gemma wording as apply-ready text.
3. **Iterate the plan** through the mode's remaining phases. Fold plan items into flow, hooks, and Steps 1–3, 5. One apply ask at the end.

**`plan only`:** stop after the written plan; do not run later phases until the user says continue.

**Body Gemma (optional, never blocking):** MAY eval a **flagged** paragraph with `--scope body` (or a temp excerpt) after the plan names it. If the gateway 500s or times out, note it and keep iterating. MUST NOT run full-post progressive (17+ units) as Phase 0.

This is **not** Cursor Plan mode. Do not switch modes for a normal revise-post run.

## Gemma 4 (default, list-scope)

Default: include list-scope Gemma in the lot plan, via **`.cursor/skills/site-revise-prose/SKILL.md`**.

```bash
python3 .cursor/skills/site-revise-prose/scripts/evaluate_prose.py \
  content/<section>/<slug>/index.md \
  --scope list
```

**Skip without asking:** `format-only`.

**Prompt (MUST unless already chosen):** At setup, stop and ask. Do **not** start the lot until they answer.

| Token already in the invoke | Behavior |
|-----------------------------|----------|
| `with gemma` / `gemma4` / `gemma` | List-scope Gemma in the plan; do not re-ask |
| `skip gemma` / `no gemma` / `without gemma` | Agent-only lot plan; do not re-ask |
| `plan only` | Write the plan, then stop |
| (none) | Ask, then wait. Default: Gemma 4 (list) |

Ask text: `Local Gemma 4 list eval is on by default (title and card fields only). Reply y / with gemma to include it, or skip gemma for an agent-only plan.`

If **AskQuestion** is available: Gemma 4 list (Recommended), Skip local eval.

Gateway 500 / timeout: say so; continue the lot plan **without** Gemma unless the user asks to retry. MUST NOT silently pretend Gemma ran.

Optional later: user says `with score` → **revise-score** after format (unchanged).

## Fixed phase order (MUST NOT reorder within a mode)

### `rough` and `polish` (format always last)

```text
Phase 0   Lot plan                agent skim + list Gemma (default); write the plan
Phase 1   revise-flow (heavy)     cuts, compression, standalone worth, full passes
Phase 2   revise-hooks            title, list fields, ## / ### headings
Phase 3   post filters            Steps 1–3 and 5 (Step 4 deferred; see below)
Phase 4   revise-flow (fine)      rhythm + voice only
Phase 5   revise-format           em dash search + emphasis audit (MANDATORY)
Phase 6   revise-score            OPTIONAL: user asks "with score"
```

### `standard`

```text
Phase 0   Lot plan                skip Gemma if user opts out
Phase 1   revise-flow (heavy)
Phase 2   revise-hooks
Phase 3   post filters            Steps 1–5 (Step 4 uses revise-format searches)
Phase 4   revise-format           final gate (catches regressions)
Phase 5   revise-score            OPTIONAL
```

### `checklist`

```text
Phase 0   Lot plan                skip Gemma if user opts out
Phase 1   post filters            Steps 1–5 in order
```

### `format-only`

```text
Phase 1   revise-format only
```

**Why format is last:** Fine flow and post edits often reintroduce em dashes (`—`) and extra `**bold**`. A single final **`revise-format`** pass catches them.

**Why heavy flow is first:** Big cuts before hook rewrites and glosses. **Why fine flow is late:** Read-aloud polish on text that includes Step 3 glosses and scenes.

## Heavy vs fine `revise-flow`

| | **Heavy** | **Fine** (`rough` / `polish` only) |
|---|-----------|-------------------------------------|
| **Passes** | Full loop per **revise-flow** (clarity → compression → flow → voice → purpose) | **Round 2–3** focus: rhythm, timing words, voice vs `original` |
| **Compression** | Allowed with guardrails; **satisfice K ≥ 4.0**, do not maximize brevity | **MUST NOT** run Pass 2 tightening; high-confidence filler only if objectively dead |
| **Scope** | Body + sidecars in scope | Body + sidecars only; **MUST NOT** rewrite `title`, Claim, `grounding` |
| **Stop** | Normal revise-flow stop conditions | One round unless user asks for more; trade-off stop applies |

**`standard` / `checklist`:** no fine flow.

## Step 4 deferral (post filters phase)

- **`rough` / `polish`:** **Defer Step 4** to the final **revise-format** phase. In the post-filters report, note `Step 4: deferred to revise-format`. Still run Step 4 **non-format** items if needed: US spelling spot-check, title emoji rules (no em dash/bold tables until format phase).
- **`standard` / `checklist`:** Run Steps 1–5; **`standard`** still runs **revise-format** again as final gate.
- **Step 1 dedup:** If the hooks phase already audited hooks, do not repeat the same checklist; carry Pass / Fix Needed from that phase.
- **Step 2 after heavy flow:** Light pass only (confirm standalone worth; do not re-cut what heavy flow already cut).
- **Step 2 scope:** Voice scrub on **all prose** in the bundle (body, `description`, `sowhat`, `grounding` digest, sidecars when in scope), not body-only.

## Scope

- **Default target:** one Hugo page under `content/` (`index.md` or path user gave).
- **MAY include** bundle sidecars (`substack.md`, `linkedin.txt`, `facebook-*.txt`) when user says whole bundle.
- **MUST** read the page **`type`** skill first.

## Voice Lock (Global Constraint)

**Preserve the author's stance, concrete examples, technical terms, and headings.** Change wording when a check requires it (title fails Hook Audit, revelation-stack fluff, negation-first pattern).

For **claims**, **video**, and **Substack** long prose, polished aphoristic fragments are **not** protected voice. Prefer clear mechanism over literary restatement (see **Step 2** → **Explanatory prose** in **`reference.md`**).

**Do not generalize domain-specific wording into bland paraphrases.** Examples:
- Do: "hippocampal system (the brain's navigation and memory hub)"
- Don't: "brain's navigation centers" (dropping the technical term)
- Do: Keep "place cells", "occipital cortex" with glosses on first use
- Don't: Replace with vague approximations like "brain location cells"

## Pass/Fail Definitions

Each step MUST be evaluated against these strict definitions:

- **Pass**: All violations for this step are fixed in the proposed revised post. No forbidden patterns remain.
- **Fix Needed**: At least one violation from the checklist remains in the text. List them explicitly.

**A step may NOT be marked Pass while any banned pattern is still present.**

## Execution (single session)

### 1. Setup

- Read target(s). Save **`original`** snapshot.
- Pick **mode** (`rough` unless user names another; cold-read-only → Step 3 only).
- Resolve **Gemma 4** per **Gemma 4 (default, list-scope)** above. If unset, **stop and ask**; do not start Phase 0.
- Open matching **type** skill.

### 2. Run phases in mode order

For each phase, collect proposed edits; **do not write to disk** until final apply.

- **Lot plan (Phase 0):** write ranked issues + phase intent. List-scope Gemma when on. If `plan only`, stop here.
- **Heavy / fine flow:** **revise-flow** output format (condensed in merged report for fine).
- **Hooks:** **revise-hooks** cold-read + tables. For **`type: claims`**, MUST run **Claim fog** on **`description`** (sentence-by-sentence).
- **Post filters:** step tables below (Step 4 deferred per mode).
- **Format:** **revise-format** MUST run Grep for `—` and emphasis audit; **Fail** if searches skipped.

### 3. Merge and present (one report)

```markdown
## Revise post: [path]

**Mode:** rough | standard | polish | format-only | checklist
**Gemma 4:** list | skipped (reason) | failed (continued)
**Phases run:** …
**Write policy:** Analysis only until you confirm apply.

### Lot plan
Ranked issues; phases that will fire; body lines to watch.

### Phase — Gemma 4 (list)
… or Skipped / failed, continued without it

### Phase — Flow (heavy)
…

### Phase — Hooks
…

### Phase — Post filters
| Step | Status |
|------|--------|
| 4 Formatting | deferred → format phase | (rough/polish only)

### Phase — Flow (fine)
… or Skipped

### Phase — Revise format
**Em dash:** 0 hits | N hits
**Bold:** …

### Merged apply list
| # | Location | Before | After | Phase |

Apply all changes? Or: `apply Phase N only`, `cancel`.
```

**Dedup:** Same line in heavy flow and post → one row, prefer voice. Hook edits + Step 1 → one row. Fine flow MUST NOT undo format fixes (apply order fixes that).

For **`checklist`** mode, use the per-step violations + Before/After tables (see **Post filter output** below) instead of multi-phase headers when only Steps 1–5 ran.

### 4. Apply (after user confirms)

Write **Merged apply list** in this order:

```text
heavy flow → hooks → post filters (Steps per mode) → fine flow (if run) → revise-format
```

| User says | Action |
|-----------|--------|
| `apply all` / `y` | Full order above |
| `apply Phase N only` / `apply Step N only` | That phase or step's rows only |
| `cancel` | Write nothing |

After `content/**/*.md` changes: **`hugo build`** (or `make build`).

**MUST NOT** auto-apply without confirmation unless user says **apply without asking**.

## Em dash lock

Zero `—` in proposed text; the format phase MUST grep the real file paths.

---

## Post filters (Steps 1–5)

| Step | Name | One-line purpose |
|------|------|------------------|
| **1** | Hook Audit | Pull readers in immediately |
| **2** | Voice Scrub | Remove AI rhetoric and teacher framing |
| **3** | Accessibility Pass | Pass the cold-reader test |
| **4** | Formatting Sweep | Enforce mechanical rules |
| **5** | Structural Integrity | Eliminate anti-patterns |

### Step 1: Hook Audit (Title + Lead)

**Goal:** Ensure the post pulls readers in immediately.

**Read first:** `.cursor/skills/site-revise-hooks/SKILL.md` for hook patterns and forbidden tone. For type-specific lead rules, read the relevant type skill (`.cursor/skills/site-video-content/SKILL.md` or `.cursor/skills/site-claims-content/SKILL.md`) before evaluating.

**Checks:**
- [ ] Title creates tension, paradox, or states a clear thesis (not a bland label)
- [ ] Title avoids corporate cadence ("unlock", "elevate", "drive outcomes")
- [ ] No fake urgency ("last chance", "everyone is talking about")
- [ ] No mystery boxes (title withholds what the post is about)
- [ ] `description`/`sowhat` fulfill the title's promise with concrete hooks
- [ ] For `type: video`: `description` is hooks above the embed, not a body summary
- [ ] For `type: claims`: `description` stands alone as a clear Claim without reading the body

#### `type: claims` — Claim fog gate on `description` (MUST NOT skip)

Before Step 1 can be **Pass** for **`type: claims`**:

- MUST cold-read **`description` (Claim)** alone, **sentence by sentence**.
- MUST apply **`.cursor/skills/site-revise-hooks/SKILL.md`** → **Claim fog (abstract cold-read fail)**.
- **Fail** any sentence a stranger cannot picture (dual-meaning “systems,” “map complexity flat,” “sit above local detail,” construct names without plain outcome).
- **Pass** when each sentence has **direction + outcome** or a **concrete noun** the card can see (layers, rank, upward summary, pyramid, decision).

If this gate fails, Step 1 status is **Fix Needed** even when the **title** hooks.

**Action:** Rewrite title and lead until all checks pass. Stop here if the hook is broken; no amount of polishing fixes a weak premise.

#### `type: video` — cold-read gate on list lead (MUST NOT skip)

**Objective:** **`description`** and **`sowhat`** decide **capture vs pass** on feeds and above the embed. If list copy does not make sense without the body or video, the reader scrolls on and the post loses.

Before Step 1 can be **Pass** for **`type: video`**:

- MUST cold-read **`description`** alone (no body, no video). **Every** bullet must decode for a stranger.
- MUST cold-read **`sowhat`** alone when present.
- MUST **double-check the first two `description` bullets** on list-style leads (highest failure rate).
- **Fail** examples: “brain shorthand,” “yesterday’s kit,” *like*/*is* with no setup, “bench/trap/slice” without a plain noun in the same line.
- **Pass** examples: named subject + optional parenthetical example + punch (“brain metaphors (computer, LLM)… **like** turns into **is**”).

If this gate fails, Step 1 status is **Fix Needed** even when the **title** hooks. See **`.cursor/skills/site-video-content/SKILL.md`** → **Cold-read gate**.

**List copy and body headings (all types with markdown body):**

- [ ] List-facing fields **cold-read** per **`.cursor/skills/site-revise-hooks/SKILL.md`** → **Capture vs pass** and **Cold-read** (claims: Claim + grounding; video: see video-content gate too).
- [ ] Body **`##`** / **`###`** headings are **hooks**, not seminar labels.
- [ ] For **`type: claims`**: body MUST **not** start with **`##`** (Thoughts band already titles the section); see **revise-hooks** → **No opening body `##`**.
- See **`.cursor/skills/site-revise-hooks/SKILL.md`** → **Body headings (`##` and `###`) — hooks, not labels**.

---

### Step 2: Voice Scrub and Standalone Worth (Active, Direct, No Decoration)

**Goal:** Remove AI rhetoric, teacher-like framing, hedges, filler, and sentences or paragraphs that sound smooth but add nothing.

#### Scope (all prose)

MUST run this step on **every prose field** under `content/`, not only the body:

- Body (Thoughts, TLDR, panel copy, sayings context)
- List-facing front matter (`description`, Claim, `grounding` digest, `sowhat`, teasers)
- Bundle sidecars when in scope (`substack.md`, `linkedin.txt`, `facebook-*.txt`)

**revise-flow** handles cadence; **this step** catches AI voice, negation fluff, and opaque labels. See also **Step 5** for negation-first patterns.

**Full banned-pattern tables, Explanatory prose, and clever-filler examples:** **`.cursor/skills/site-revise-post/reference.md`**. MUST apply that reference when scrubbing.

#### Standalone worth test (MANDATORY)

For **every** sentence and paragraph, ask: **If this stood alone, would it be worth saying?**

It must add at least one of: a new fact, example, named mechanism, higher stakes, or a turn the reader did not already have from the title, list copy, or previous paragraph. A **new metaphor or coined label for the same claim** is not a fact, mechanism, stakes, or turn.

**CUT** (quote the full sentence or paragraph) when:
- It restates the section H2 or spine line in softer words.
- It is throat-clearing before the real point ("The way out is...", "This is where it gets interesting...").
- It is a wise-sounding summary with no new concrete detail.
- It only bridges ideas already stated above and below.
- It only relabels the previous claim (`the blend`) or restates its implication with a new metaphor pair (`sealed file` / `steering wheel`). See **Metaphor-shell restack** in **`reference.md`**.
- Its only job is a sales, factory, driving, or hardening verb (`sells the blend`, `steer the rewrite`, `assembling it`, `next assembly`, `template can harden`) or a parallel "name the box… and the box." See **Industry-verb shells** in **`reference.md`**.
- Deleting it changes almost nothing.

**MERGE** when it repeats a nearby beat but has one usable detail; fold that detail into the stronger sentence, then cut the rest.

Apply across body, **`sowhat`**, **`description`**, **`grounding`**, section openers, closers, and sidecars. Do not keep lines because they "sound good."

**Checks:**

- [ ] Every sentence and paragraph passes the **standalone worth** test
- [ ] Every paragraph uses active voice (subjects perform actions)
- [ ] Zero banned patterns from **`reference.md`** (including revelation stacks / rhetorical fragments / **Metaphor-shell restack** / **Industry-verb shells**)
- [ ] Claims, video, and Substack prose also passes **Explanatory prose** in **`reference.md`**
- [ ] No new hedges or intensifiers introduced that were not in the original
- [ ] No stacks of abstract nouns where a concrete verb would do
- [ ] No uniform punch-line rhythm across a whole paragraph (read aloud)
- [ ] No **clever filler** or **mixed metaphors** in the same paragraph (**`reference.md`**)
- [ ] Every coined label or named construct **cold-reads** without the body (Step 3 gloss if needed)

**Action:** List **cut candidates** (full quotes, CUT vs MERGE, one-line why). Rewrite passive sentences. Delete condescending openers, hedges, and failed standalone-worth lines. Replace abstract nominalizations with verbs.

---

### Step 3: Accessibility Pass (Cold Reader Test)

**Goal:** Ensure the post makes sense without reading the source or watching the video. Preserve all technical terms.

#### Named constructs and coined labels (MUST)

Any **label the reader did not bring** fails cold-read unless decoded in the **same sentence** (parenthetical gloss, appositive, or plain clause after a colon). Applies to technical terms **and** author metaphors (`fixed energy budget`, `stress response`, `misallocation`).

| Fail | Pass |
|------|------|
| "Picard's hook is a **fixed energy budget**." | "…runs on a **fixed energy budget**: a limited daily pool of cellular spend." |
| "…stuck on vigilance." | "…when the **stress response** keeps spending after the moment passes." |

**Test:** A stranger can quote back what the line claims without reading the body or watching the video.

**Technical Term Rule:**
- Do NOT remove or rename technical terms (e.g., "hippocampal system", "occipital cortex", "place cells", "vicarious trial and error")
- Keep the exact term and add a gloss in parentheses on first use
- **Fail this step if any technical term from the original is missing or paraphrased away**

**Do/Don't Examples:**
- Do: "hippocampal system (the brain's navigation and memory hub)"
- Don't: "brain's navigation centers" (term dropped)
- Do: "place cells (neurons that track where you are in space)"
- Don't: "location-tracking neurons" (term dropped)

**Note:** When the post covers AI, neuroscience, or cognitive science, also read `.cursor/skills/site-ai-for-general-audience/SKILL.md` for plain-language equivalents of common jargon.

**Checks:**
- [ ] Every technical term has a plain-language gloss in parentheses on first use (or is common knowledge)
- [ ] All original technical terms are preserved exactly; none are missing or paraphrased away
- [ ] No jargon dropped without context (e.g., "hippocampal system", "occipital cortex", "place cells")
- [ ] Concrete visual language replaces abstract descriptions ("simulating the future" not "vicarious trial and error")
- [ ] Short paragraphs (1-3 lines) for scanability

**Action:** Add glosses in parentheses. Chunk dense paragraphs. Use "like" or "for example" to clarify metaphors. **Never delete technical terms.**

#### `type: video` — re-run cold-read on `description` + `sowhat`

- Step 3 is **Fix Needed** for **`type: video`** if **`description`** or **`sowhat`** still fails the naive-reader test after body edits (body glosses do not excuse opaque list bullets).
- Re-apply **`.cursor/skills/site-video-content/SKILL.md`** → **Cold-read gate** before marking Step 3 Pass.

---

### Step 4: Formatting Sweep (Punctuation + Emphasis)

**Goal:** Enforce strict mechanical rules.

**MUST run** **`.cursor/skills/site-revise-format/SKILL.md`** (em dash search + **emphasis** audit per **revise-emphasis**). Step 4 is **Fix Needed** if either mandatory search was skipped or em dash hits &gt; 0.

**Checks:**
- [ ] **Em dash search executed** (Grep or `rg` for U+2014 `—` on target path). ZERO hits in scope.
- [ ] **Bold audit executed** per **`.cursor/skills/site-revise-emphasis/SKILL.md`** (count spans per block; flag wall-to-wall).
- [ ] Default-language `content/` uses **US English** spelling (**generalize**, not **generalise**); see **`.cursor/rules/site-content-markdown-writing.mdc`** → **Language**.
- [ ] No decorative emoji in `description` or `grounding`.
- [ ] For `cognitive-memetics`: no leading emoji in `title`.
- [ ] For other sections: max two leading emoji in `title` if they signal the hook.

**Action:** Report em dash table and bold audit table (see **revise-format**). Replace every `—`. Trim excess `**` (markup only). Strip inappropriate emoji.

**User shortcut:** `/revise-format` on the same path runs only this step.

---

### Step 5: Structural Integrity (Anti-Patterns + Type Compliance)

**Goal:** Eliminate structural flaws and verify type-specific requirements.

**Read first:** The relevant type skill for required fields:
- `type: video` → `.cursor/skills/site-video-content/SKILL.md`
- `type: claims` → `.cursor/skills/site-claims-content/SKILL.md`
- `type: panel` / `type: sayings` → `.cursor/skills/site-cognitive-memetics-content/SKILL.md`

**Banned Anti-Patterns:**

| Pattern | Example | Fix |
|---------|---------|-----|
| Negation-first (hard) | "It is not X; it is Y", "Not A but B" used as a thesis flourish | State Y directly |
| Negation fluff (soft) | "as allocation trouble, not a fuel shortage", "…, not weakness" | Drop the negated tail; state the positive claim |
| Passive piles | "It was determined that..." | Direct statement |
| New hedges | "one possibility", "some argue" (introduced in revision) | Cut or name the source |

**Allowed (claim-scope, not banned):** Narrowing an overclaim with evidence restraint. Example: "This does not mean every criticism is projection." That protects accuracy; it is not poetic "not a mirror; a wound" packaging.

**Do/Don't examples (negation):**

| Do | Don't |
|----|-------|
| "Exhaustion and lost meaning can show up together when the **stress response** keeps spending after the moment passes." | "…as allocation trouble, not a fuel shortage." |
| "This is a behavioral map." | "This is not an algorithm; it is a behavioral map." |
| "Fatigue reads as **misallocation**." | "Fatigue stops feeling like weakness and starts feeling like misallocation." (negation pair; prefer direct when the positive line is clear) |
| "This does not mean every criticism is projection. Strong moral judgment can sometimes serve self-protection as well as truth-seeking." | "It is not seeing; it is a mirror." (oracular negation pair with no scope work) |

**Checks:**

- [ ] Zero negation-first constructions (hard **and** soft forms in the table above), except allowed claim-scope narrowers
- [ ] Zero "It was determined that..." or "it was found that..." passive piles
- [ ] Zero new hedges or filler phrases introduced during revision
- [ ] `type: video`: Has `youtube_id`, has `sowhat`, `description` is not body duplication
- [ ] `type: claims`: Has `description` (Claim), has `grounding`, Claim reads alone
- [ ] `type: panel`/`sayings`: Has `description` (teaser), follows cognitive-memetics field rules

**Action:** Rewrite negation-first sentences as direct statements. Verify required front matter fields per type.

---

## Post filter output (`checklist` or per-step detail)

### For each step, produce:

1. **Step header** with number and name
2. **Violations list** (quoted exact phrases from original that break the rules, or "No violations found")
3. **Numbered table** with four columns: `#`, `Location`, `Before`, `After`
4. **Status line** indicating Pass/Fix Needed per the Pass/Fail Definitions

### Example:

```
### Step 1: Hook Audit

**Violations found:**
- "Embodied Cognition and Language" (title is bland label, no tension)
- "Spatial cognition uses roughly **half the brain**" (description leads with abstraction)

| # | Location | Before | After |
|---|----------|--------|-------|
| 1 | title | Embodied Cognition and Language | 🖐️🧠 You point first, then you think. |
| 2 | description (bullet 1) | Spatial cognition uses roughly **half the brain**... | If you sit on your hands, your ability to explain... |

**Status:** Pass (Title now creates tension. Description leads with concrete hook.)
```

### At the end of checklist-only runs:

1. **Summary header:** "## Revision Summary"
2. **Status overview:** List each step with Pass or Fix Needed status
3. **User prompt:** "Apply all changes? Or specify which steps (e.g., 'apply Steps 1-3', 'apply only Step 2', 'revert Step 4')?"
4. **Proposed revised post:** Show the full post as it *would* appear (for reference only; do not write yet)

## Related skills

- `.cursor/skills/site-revise-flow/SKILL.md`
- `.cursor/skills/site-revise-hooks/SKILL.md`
- `.cursor/skills/site-revise-format/SKILL.md`
- `.cursor/skills/site-revise-prose/SKILL.md` (Gemma 4 default eval)
- `.cursor/skills/site-revise-score/SKILL.md`
- `.cursor/skills/site-revise-post-es/SKILL.md`
- `.cursor/skills/site-revise-post/reference.md` (Step 2 detail)
