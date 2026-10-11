# Revise flow — reference

## Pass allowed / blocked edits

### Pass 1: Clarity

| Allowed | Blocked |
|---------|---------|
| Add missing prepositions | Cuts made only for elegance |
| Fix tense and agreement | Metaphor rewrites |
| Remove accidental duplication (doubled words) | Style rewrites |
| Resolve vague pronouns when referent is unclear | |

**Ask:** Is every sentence immediately parseable?

### Pass 2: Compression

**Default:** satisfice (K ≥ 4.0), do **not** maximize brevity. If the paragraph already reads aloud well, propose **no** Pass 2 edits.

| Allowed | Blocked |
|---------|---------|
| Remove filler (see list below) when read-aloud shows dead weight | Any cut when K ≥ 4.0 and F, V already pass |
| Remove **same idea twice in one paragraph** (echo redundancy) | Remove `when` / `then` / `still` / `now` / `just` unless clearly dead weight |
| Shorten tails that **only** restate the metaphor | Convert spoken syntax to noun-heavy abstraction |
| | Merge parallel hooks or refrains for "tightness" |
| | Cut beat markers, contrast, or emphasis |
| | Remove lines that carry oral force for a shorter paraphrase |
| | Second compression pass on a paragraph already guardrail-clean |

**Ask:** Does this phrase fail the paragraph, or am I just making it shorter?

**High-confidence filler (remove when dead weight):** really, very, actually, some, quite, basically, kind of, sort of, in order to, the fact that, it is important to note that.

**Not automatic filler:** `just`, `when`, `then`, `still`, `so`, `now` often carry rhythm. Cut only after read-aloud proves they add nothing.

### Over-compression smells (revert)

| After cut, you hear… | Action |
|----------------------|--------|
| Telegraphic / outline voice | Revert; restore spoken syntax |
| Same rhythm every sentence | Revert merge; vary length |
| Undefined noun + tautology after a mechanism sentence ("The blend is your experienced reality") | Revert; next sentence must add how, not a label. **Metaphor-shell restack:** **`.cursor/skills/site-revise-post/reference.md`** |
| Extra metaphor pair after a mechanism sentence ("sealed file" / "steering wheel") | Revert; state the implication in plain words or cut |
| Industry-verb shell after a cut (`sells the blend`, `steer the rewrite`, `next assembly`, `template can harden`, abstract "the room") | Revert; spoken verb + real scene. **Industry-verb shells:** **`.cursor/skills/site-revise-post/reference.md`** |
| Flat vs `original` | Round 3 restore |
| Choppier read-aloud | Revert even if K improved |

### Pass 3: Flow

| Allowed | Blocked |
|---------|---------|
| Keep or restore beat words | Over-standardizing |
| Alternate short and medium sentence lengths | Back-to-back compressed sentences with the same rhythm |
| One bridge sentence between dense ideas | Over-tightening |
| Prefer spoken syntax if equally clear | |

**Ask:** Does it read aloud naturally?

### Pass 4: Voice

| Allowed | Blocked |
|---------|---------|
| Keep high-value phrases even if not minimal | Generic corporate phrasing |
| Keep emotionally exact hedges: mostly, roughly, can | Flattening contrast |
| Keep contrastive structures the author uses often | Replacing signature metaphors |
| Restore lines from `original` where personality dropped | |

**Ask:** Does it still sound like the author?

### Pass 5: Purpose

| Allowed | Blocked |
|---------|---------|
| Tune density to format (article, newsletter, caption, script) | Generic polish that weakens intent |
| Surface "so what" earlier when audience needs it | Stripping platform-required structure |
| Add signposts or headings for scannability | "Best writing" with no job fit |

**Ask:** Does it serve the job and audience?

## Paragraph rubric (0–5)

### Clarity (C)

| Score | Meaning |
|-------|---------|
| 5 | Understood on first read |
| 3 | Understandable but one stumble |
| 1 | Confusing or overloaded |

### Compression (K)

| Score | Meaning |
|-------|---------|
| 5 | No obvious dead weight |
| 3 | One redundant phrase |
| 1 | Multiple phrases could vanish without loss |

### Flow (F)

| Score | Meaning |
|-------|---------|
| 5 | Reads aloud smoothly |
| 3 | A bit choppy or too dense |
| 1 | Robotic or halting |

### Voice (V)

| Score | Meaning |
|-------|---------|
| 5 | Unmistakably the author's |
| 3 | Partly generic |
| 1 | Could be anyone's text |

### Purpose (P)

| Score | Meaning |
|-------|---------|
| 5 | Fully suited to audience and format |
| 3 | Solid but misweighted |
| 1 | Wrong density, tone, or frame |

## Document score

\[
D = 0.25 C + 0.20 K + 0.25 F + 0.20 V + 0.10 P
\]

Use paragraph-level averages for C, K, F, V, P unless the user requests sentence-level scoring.

## Minimal operating procedure

1. Read the paragraph aloud.
2. Fix only grammar and obvious repetition (pass 1).
3. Cut filler and repeated logic (pass 2); apply guardrails.
4. Read aloud again (pass 3).
5. Restore any small word whose loss hurts cadence.
6. Restore any phrase whose loss makes it sound less like the author (pass 4).
7. Tune for format and audience (pass 5).
8. Stop when the next edit improves concision but hurts flow or voice.

## Round cheat sheet

| Round | Do | Do not |
|-------|-----|--------|
| 1 | Obvious grammar; high-confidence fluff; mark clipped lines | Aggressive tightening; metaphor swaps |
| 2 | Aloud read; **restore** timing words; split on stumble only | Wall-to-wall shortening; merge hooks for brevity |
| 3 | Compare to `original`; restore personality **and cadence** | Re-import all fluff |
| 4+ | Pareto edits only; **no new compression** if K passes | Cosmetic churn; chasing ΔD with cuts |
