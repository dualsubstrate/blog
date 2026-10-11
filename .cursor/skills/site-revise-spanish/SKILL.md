---
name: site-revise-spanish
description: >-
  Audits Spanish site copy via local Gemma (evaluate-only) for native naturalness
  (not translation-shaped prose) and conceptual fidelity (do not soften contrasts
  or drop meaning load when “improving”). Loads index.es.md and ES sidecars;
  scores calques, read-aloud, cross-file sync, meaning load. Use when the user
  asks revise spanish, revise-spanish, revisa español,
  naturalidad del español, suena a traducción, calques, español nativo, or a
  quick Spanish idiom pass. For full Hugo Spanish revision (hooks, format,
  URLs), use revise-post-es.
---

# Revise Spanish (native naturalness, local Gemma)

## Purpose

Answer two questions:

1. **Does it read like Spanish written for Spanish speakers, or like translated English?**
2. **Do proposed “improvements” keep the meaning load** (contrasts, clinical precision, subject/object), or do they soften the idea to sound smoother?

**Default path:** evaluate-only via local gateway (Gemma 4), same stack as **revise-prose**. Does **not** auto-edit `content/`. Apply fixes after user confirms, or use **revise-post-es** for a full publish workflow.

## When to use

| Use **revise-spanish** | Use another skill |
|------------------------|-------------------|
| After new or edited `index.es.md` | **spanish-translation-content** to draft the sibling (that skill already runs Gemma 4 on the first draft) |
| "Sounds translated", "revisa la naturalidad" | **revise-post-es** for full 5-step Hugo + format + URLs |
| ES sidecars idiom pass | **revise-flow** on English only |
| Quick calque audit before publish | **revise-prose** for English human feel |

## Defaults (local eval)

| Setting | Value |
|---------|--------|
| Base URL | `LOCAL_LLM_BASE_URL` or `http://127.0.0.1:1320/v1` |
| Model | `LOCAL_LLM_MODEL` or `gemini_free/gemma-4-26b-a4b-it` |
| Scope | `full` (all ES files in the bundle that exist) |
| Mode | `progressive` (default): each field/paragraph with **prior units as context**, then overall synthesis. `--mode whole` for one shot. |

Rubrics: **`reference.md`**, pack: **`packs/evaluate.md`**.

## Scope (default files)

**Priority:** fix **`index.es.md` first** (source of truth). Sidecars are reflections; audit them after the index is Ready, or with `--scope site` / `--scope body` when the user is mid-index.

Same bundle (when `--scope full`):

1. **`index.es.md`** (always first in progressive order)
2. **`facebook-es.txt`** if present
3. **`linkedin.es.txt`** if present
4. **`substack.es.md`** if present

Script also attaches a short **English meaning context** from `index.md` (list/lead only; not syntax to copy).

## Index is source; sidecars are reflections

- **MUST** treat `index.es.md` as the only Spanish thesis/wording source; sync sidecars after apply.
- **MUST** keep `title` / `subtitle` as one short line (no body bleed into scalars).
- **MUST** preserve index lexicon (*persona*, *patrón*, *sobre-diagnostican*).
- **MUST NOT** expand sidecars beyond the index or soften clinical precision to “sound nicer.”
- **MUST NOT** edit the index to match a sidecar.

## Procedure (caller)

1. Resolve bundle (path to `index.es.md`, a sidecar, or the folder). Prefer focused Spanish file or sibling of focused EN post.
2. Prefer **`--scope site`** or **`--scope body`** when repairing the index; use **`--scope full`** only after the index is stable (or when the user asks for a full-bundle audit).
3. Run evaluate-only (do not hand-roll the API call unless the script fails):

```bash
python3 .cursor/skills/site-revise-spanish/scripts/evaluate_spanish.py \
  content/<section>/<slug>/index.es.md \
  --scope site
```

Scopes: `full` | `site` | `list` | `body` | `facebook` | `linkedin` | `substack`.  
Optional: `--mode whole`, `--model`, `--out docs/research/YYYY-MM-DD-local-es-<slug>.md`, `--no-en-context`.

Progress on stderr (`Evaluating unit k/n…`). Surface the full stdout report (per-unit + **Overall synthesis**).

4. Surface the **full Gemma report** to the user.
5. **Do not edit files** until the user confirms (`y`, `aplica`, `apply all`), unless they already asked to apply in the same message.
6. On apply: fix **`index.es.md` first**; then sync refrains into ES sidecars; re-read changed lines.

### Progressive unit safety (MUST)

When any rewrite/eval unit is applied by an agent or LLM:

- **title / subtitle:** exactly **one** plain line; if output has `###`, newlines, or body paste → **discard**.
- **description / sowhat:** keep block shape; **MUST NOT** inject body paragraphs. For **`description`** list bullets (video novelty lead), **MUST NOT** use pitch-deck colon formulas (*X: de A a B*), bullets starting with *Cómo…* (*How to* calque), or *mentes diseñadas para* (*brains built to*); use active Spanish claims with finite verbs.
- **body unit:** rewrite only that unit; **MUST NOT** repeat prior paragraphs or emit meta labels (*UNIDAD ACTUAL*, *Aquí tienes*).
- **Reverse-translation test:** if ES remounts into the EN sentence almost word-for-word → fail `anti_calque` and rewrite syntax.

### Apply / “mejora” guardrails (MUST)

- **MUST** treat author micro-edits and exact phrases as binding (*sobre-diagnostican a la persona*, *se enfoca en el patrón*).
- **MUST** keep conceptual contrasts (e.g. clinic labels fall short or over-diagnose the person **vs** this name focuses on the pattern).
- **MUST** when a sentence stacks cause / effect / payoff, clarify causal order; do not compress into prettier opacity.
- **MUST NOT** swap precision terms for softer generics (*clínicas* → *de consultorio*; *patrón* → *hábito*) just to “sound native”.
- **MUST NOT** drop verb objects that carry ethical or clinical load (*a la persona*).
- **MUST NOT** prioritize fluency over `conceptual_fidelity` when both conflict.
- **MUST NOT** treat “reads smoothly” as native if clause order still mirrors English (esqueleto EN + léxico ES). Rewrite syntax, not only vocabulary.
- **MUST NOT** apply cosmetic word-level patches (swapping 1 or 2 words) when a paragraph fails the reverse-translation test or carries an English sentence skeleton. Swapping *impuesto de carrera* for *lastre* while keeping the English clause order leaves the prose sounding like dubbing.
- **MUST** re-adapt the full paragraph when syntax fails reverse-translation: extract the claim and mechanism, close the English source, and reconstruct the paragraph in natural Spanish using organic hypotaxis, active verbs, and natural connectors.

### Author micro-edits bind

When the user pastes line fixes (or has already locked wording in `index.es.md`):

- **MUST** treat that wording as law for the line, even if it matches a calque sniff (*se les cae el piso*).
- **MUST** apply a **surgical patch** only; do not rewrite the surrounding paragraph for “flow.”
- **MUST** prefer clinical/conflict precision and dry manifesto tone over literary elegance.
- **MUST NOT** “clean up” raw author voice into professional/consultorio Spanish.
- **MUST NOT** call Gemma to “fix” an intentional author calque (Gemma may “correct” it wrongly).
- **MUST NOT** soft-synonym precision terms (*patrón* → *hábito*).

**Fix-direction example:** soft abstract *Sienten una gran decepción cuando…* → author-direct *Se les cae el piso cuando…* (keep if author chose it).

### Workflow: line-by-line micro-edits

1. Identify intent of the patch (force/tension, clinical precision, or own-voice calque).
2. Patch **only** the affected span in **`index.es.md`**.
3. Sync sidecars **after** the index line is stable (same refrain/thesis); do not invent sidecar-only wording.
4. Call Gemma (`--scope site`) only when you need a consistency audit across the index or before a full-bundle pass; **skip Gemma** for pure author wording swaps.

**Progressive rewrite bans:** no body bleed into `title`/`subtitle`; no duplicated sections; no meta labels (*UNIDAD ACTUAL*); no softening for “professional” tone.

## Gateway down (no fallback)

MUST report that Gemma did not run, name the error, and stop. There is **no** agent-only substitute audit.

- MUST NOT run a manual read-aloud, calque, grammar, or cross-file pass in place of the model report.
- MUST NOT edit `index.es.md` or ES sidecars on your own judgment after a gateway failure.
- MAY retry the same command once when the failure looks transient, then report either way.
- The user decides next: retry later, ship unaudited, or hand-review.

## Required output (after apply decision)

When applying or summarizing for chat after the model report, use:

```markdown
## Revise Spanish: [bundle]

**Verdict:** [from report]

### Top fixes applied or proposed
- [quote] → [fix]

### Cross-file sync
- [...]

Apply changes? (`y` / `cancel`)  # if not yet applied
```

**MUST NOT** paste full rewritten files unless the user asks for **`[Texto final]`** or the diff is tiny.

## Pass / fail

| Result | Meaning |
|--------|---------|
| **Ready** | No objective calques; clean read-aloud; sidecars aligned; meaning load intact |
| **Almost** | 1–3 fixable lines (style or clarity without thinning the thesis) |
| **Re-adapt section** | English skeleton, or “smooth” rewrite that lost contrasts; use **revise-post-es** or **spanish-translation-content** |

## Related

- **English local twin:** `.cursor/skills/site-revise-prose/SKILL.md`
- **Authoring siblings:** `.cursor/skills/site-spanish-translation-content/SKILL.md`
- **Full Spanish pipeline:** `.cursor/skills/site-revise-post-es/SKILL.md`

## Reference bundle

`content/human-condition/2026-06-11-think-with-elbow/` (`index.es.md`, sidecars).
