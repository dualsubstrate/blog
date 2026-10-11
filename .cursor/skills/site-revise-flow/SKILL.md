---
name: site-revise-flow
description: >-
  Revises prose through a five-pass scoring loop (clarity, compression, flow,
  voice, purpose) with trade-off guardrails and convergence rounds. Use when the
  user asks for prose pass, clarity pass, compression pass, read-aloud revision,
  cadence repair, revise-flow, or a multi-pass edit that must not flatten voice.
  Works on any markdown or plain text (site posts, substack.md, social sidecars).
  For Hugo posts under content/, run this before revise-score or as a phase of revise-post.
  In **revise-post**, use **heavy** flow at the start and optional **fine** flow after post
  filters (before **revise-format**). See **`.cursor/skills/site-revise-post/SKILL.md`** modes.
---

# Revise flow

## Purpose

Make prose **clear, then tight, then natural, then personal, then suited to purpose**. Never accept a gain in neatness that makes the text sound dead.

**Compression is the highest-risk pass.** The goal is **good enough**, not **as short as possible**. If the draft already reads aloud well, **do not keep tightening**.

This skill handles **line-level prose** (grammar, filler, cadence, voice, format fit). It does **not** replace Hugo hook audits, type compliance, or site formatting rules. For `content/**/*.md` site posts, run **`revise-flow`** first, then **`revise-score`** and/or **`revise-post`**.

**AI voice and negation fluff** are **not** owned here. After flow passes, MUST run **`.cursor/skills/site-revise-post/SKILL.md`** → **Steps 2, 3, and 5** on all prose fields (body, list copy, sidecars).

**Compression guardrail:** Satisficing cuts **MUST NOT** introduce clever filler, caption pivots, staccato punch stacks, **Metaphor-shell restack**, **abstract social metonymy** ("the room" for a group), or **Industry-verb shells** (`sells`, `steer`, `assembling`, `harden`, fake scene with no people; see **revise-post** → **Step 2** and **`.cursor/skills/site-revise-post/reference.md`**). If tightening creates those patterns, **revert** and merge into one spoken sentence instead.

## Voice lock (global)

- MUST preserve the author's stance, concrete examples, contrast that carries information, emphasis, and oral cadence.
- MUST NOT replace vivid wording with generic corporate phrasing.
- For **claims**, **video**, and **Substack** long prose: MUST rewrite **revelation stacks**, rhetorical noun fragments, and poetic thesis remixes into explanatory paragraphs (see **revise-post** → **Step 2** → **Explanatory prose**). Aphoristic lines are **not** protected when they fail that bar.
- For **cognitive-memetics** (panel / sayings): MUST follow **Hybrid prose** in **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`**. MAY keep street/cartoon punch; MUST NOT keep empty punch closers or poetic thesis remixes.
- MUST store an **`original` snapshot** at the start (full file or scoped selection). Pass 4 (Voice) compares **`current`** to **`original`**.

## Scope

| Request | Behavior |
|---------|----------|
| Full loop (default) | All five passes in order, then rounds until stop condition |
| Single pass only | User names pass (e.g. "compression pass only"): run that pass only; still apply guardrails |
| Round only | User names round (e.g. "round 2 rhythm repair"): run that round's rules across passes as needed |
| Selection | User points at paragraphs or a file path: scope state to that text only |

## Unit of analysis

1. **Paragraph** first: score and edit per paragraph.
2. **Sentence** second: only when a paragraph fails a pass threshold or guardrail triggers a revert.

## Pass order (fixed)

Run passes **in this order only**. One pass may modify a paragraph at a time.

| # | Pass | Threshold (pass if ≥) |
|---|------|------------------------|
| 1 | Clarity | 4.5 |
| 2 | Compression | 4.0 (floor only; do not maximize) |
| 3 | Flow | 4.5 |
| 4 | Voice | 4.5 |
| 5 | Purpose | 4.0 |

Allowed and blocked edits per pass: **`.cursor/skills/site-revise-flow/reference.md`**.

## Scoring

- Score each paragraph on **Clarity (C)**, **Compression (K)**, **Flow (F)**, **Voice (V)**, **Purpose (P)** from 0 to 5 using the rubric in **`reference.md`**.
- **Document score** (weighted):

  `D = 0.25·C + 0.20·K + 0.25·F + 0.20·V + 0.10·P`

  Use paragraph averages for C, K, F, V, P unless the user asks for sentence-level detail.

- **Good enough:** every dimension at or above its threshold **and** no guardrail reverts pending.

## Compression discipline (MUST — anti over-tightening)

Pass 2 has the **lowest threshold (4.0)**. Agents often over-cut and **destroy fluency and naturalness**. Treat compression as **satisficing**: meet the floor, then stop.

### When Pass 2 proposes nothing (valid outcome)

- Paragraph already **K ≥ 4.0** and **F ≥ 4.5** and **V ≥ 4.5** → **MUST NOT** propose further cuts for that paragraph.
- Whole document already meets all thresholds → Round 1 Pass 2: **high-confidence filler only** (see reference). Round 2+ is **rhythm and voice**, not more shortening.
- User text is **social sidecar** or already edited in session → bias to **zero** compression edits unless a line is objectively bloated.

### MUST NOT (over-compression)

- Merge **intentional parallel hooks** (rule of three, *En X… En Y…*, matched refrains).
- Strip **timing words** (`when`, `then`, `still`, `so`, `now`, `just`) unless read-aloud proves they are dead weight.
- Replace **spoken syntax** with noun-heavy or telegraphic compression.
- **Combine two sentences** that each carry a separate beat into one abstract line.
- Run **a second compression pass** on the same paragraph in one session after guardrails passed.
- Chase **higher K** or **ΔD** when **F** or **V** would drop (trade-off stop).

### After every compression candidate (MUST)

1. **Read the paragraph aloud** (before vs after).
2. IF it sounds **choppier, flatter, or more translated** → **revert** (even if K improved).
3. IF meaning is unchanged but **cadence lost** → **revert** (rhythm words and short parallels often earn their place).

### Signs the draft was cut too hard

- Telegraphic stack: every sentence same length, no breath.
- Punch lines from **`original`** disappeared without a clearer replacement.
- Reads like outline notes, not spoken prose.
- **Round 3 (identity):** restore from **`original`** until V ≥ 4.5 again.

For Spanish bundles, pair with **`.cursor/skills/site-revise-spanish/SKILL.md`** after flow if copy still sounds translated.

## Trade-off guardrails (MUST)

After any edit in **Compression** (pass 2):

- IF **F < 4.5** OR **V < 4.5** after the edit → **revert** that edit and note why in the edit log.

After any edit in **Clarity** (pass 1):

- IF **V < 4.5** and meaning is unchanged → prefer revert; clarity must not flatten voice.

**Stop accepting edits** when the next change would raise one score but lower **F** or **V** below threshold (trade-off stop).

**MUST NOT** treat "Compression: Fix needed" as a mandate to shorten. Fix means **remove dead weight only**, not **shrink working prose**.

### Canonical example (keep original)

| | Text |
|---|------|
| Original | The trap is when you forget you made that cut. |
| Compression candidate | The trap is forgetting you made that cut. |
| Result | **Keep original.** Compression improves K; Flow and Voice drop. |

## Rounds

| Round | Focus |
|-------|--------|
| **1** | Coarse: grammar, structural clarity, **high-confidence filler only**; flag clipped sentences; **no aggressive tightening** |
| **2** | Rhythm: read aloud; **restore** timing words (`when`, `then`, `still`, `so`, `now`); split only where reading stumbles; **merge only when stumble is objective** |
| **3** | Identity: diff vs `original`; **restore** phrases where meaning holds but personality or cadence dropped |
| **4+** | Convergence: **no new compression** unless K still below 4.0; Pareto edits only |

Repeat full pass sequence within each round until the round's focus is addressed, then advance.

## Stop conditions

Stop the loop when **any** of:

1. A full round produces **no paragraph changes**.
2. **ΔD < 0.1** between rounds and all thresholds met.
3. Further edits only trade one dimension against **F** or **V** (trade-off stop).
4. User says stop.

## Purpose pass and site sidecars

- **`substack.md`** / **`substack.es.md`**: Purpose pass MUST respect **`.cursor/skills/site-substack-post/SKILL.md`** (scanability, `##` hooks, emphasis). Do not strip bold or blockquotes for minimalism.
- **`content/`** Hugo posts**: Purpose pass notes `type` and audience; do not rewrite front matter hooks here (use **`revise-post`** Step 1).
- **Social sidecars** (`linkedin.txt`, `facebook-*.txt`): Purpose pass checks platform density and CTA shape; voice pass keeps platform templates intact.

## Execution (analysis first)

1. **Read** the full target (or selection). Save **`original`**.
2. **Score** each paragraph (C, K, F, V, P) before edits.
3. **Round 1:** Run passes 1→5 in order. **Skip Pass 2** on paragraphs that already meet K, F, and V floors. Log every proposed edit; **zero compression proposals is valid**.
4. **Rounds 2–4+** as needed until a stop condition hits.
5. **Output** using the format below. **Do not write to the file yet.**
6. **Ask:** "Apply all changes? Or specify round/pass (e.g. 'apply round 1 only', 'apply pass 2 only')?"
7. **Only after explicit confirmation**, apply selected edits.
8. For `content/**/*.md` after apply: suggest **`revise-score`** or **`revise-post`** for site checklist.

Default is **analysis-only** unless the user says "apply until stable" or "apply all rounds without asking."

## Required output format

### Opening

```
## Revise flow: [file or selection]

**Document score:** D = X.XX (before) → Y.YY (after proposed)
**Thresholds:** C≥4.5 K≥4.0 F≥4.5 V≥4.5 P≥4.0
**Round:** N (stop reason if ending)
**Compression edits proposed:** N (0 is valid when draft already flows)
```

### Per paragraph (only paragraphs with changes or sub-threshold scores)

```
### Paragraph N

| C | K | F | V | P |
|---|---|---|---|---|
| before | … | … | … | … |
| after (proposed) | … | … | … | … |

**Pass:** [name] — [allowed edit summary]

| # | Before | After | Reverted? |
|---|--------|-------|-----------|
| 1 | … | … | yes/no + reason |

**Guardrail:** [none | reverted pass 2 edit: flow dropped]
```

### Closing

```
## Loop summary

| Pass | Status | Notes |
|------|--------|-------|
| 1 Clarity | Pass / Fix needed | … |
| … | … | … |

**Trade-offs blocked:** [count] — [one-line list]

Apply all changes? Or specify round/pass (e.g. "apply round 1", "apply pass 3 only", "cancel").
```

If no edits proposed, state that and report scores only.

## Selective application (after confirmation)

| User says | Action |
|-----------|--------|
| apply all / yes | Write all proposed edits |
| apply round N | Write only that round's edits |
| apply pass N only | Write only that pass's edits |
| cancel / no | Write nothing |

## Related skills

- **`.cursor/skills/site-revise-flow/reference.md`** — rubric, allowed/blocked tables, filler list
- **`.cursor/skills/site-revise-score/SKILL.md`** — Hugo editorial score (after prose loop on site posts)
- **`.cursor/skills/site-revise-post/SKILL.md`** — full-lot conductor (heavy/fine flow phases) and post filters
