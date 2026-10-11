---
name: site-revise-prose
description: >-
  Evaluate-only local prose critique for Hugo posts via OpenAI-compatible gateway
  (default http://127.0.0.1:1320/v1). Scores human-friendliness, bar-test tone, and
  AI artifacts; does not rewrite content. Triggers: revise-prose, local prose
  review, local evaluate, bar-test critique, human-friendliness score,
  evaluate-only.
---

# Revise prose (evaluate only)

**Workflow id:** `revise-prose`

Diagnose readability and human voice for a Hugo post using a **local LLM**. **Must not** edit `content/` from this skill. Apply fixes later via **revise-post** and **revise-flow**.

## When to use

- local prose review / local evaluate / revise-prose
- bar-test critique / human-friendliness score
- evaluate-only (no rewrite)
- Offline twin of Perplexity prose review when the gateway is up

## Peers

| Skill | Role |
|-------|------|
| **This skill** | Local evaluate-only: voice, bar test, AI feel |
| `.cursor/skills/perplexity-browser-research/` + `.cursor/perplexity/packs/prose-review.md` | Same diagnostic intent via Perplexity |
| `.cursor/skills/site-revise-score/SKILL.md` | Editorial mechanics / composite scorecard |
| `.cursor/skills/site-revise-post/SKILL.md` + `revise-flow` | Apply fixes after critique. **revise-post** Phase 0 uses **`--scope list`** (not full-body progressive); user can `skip gemma` |

## Rubrics

Scored criteria and AI artifact shortlist: **`reference.md`** in this folder.  
Prompt pack filled by the CLI: **`packs/evaluate.md`**.

## Defaults

| Setting | Value |
|---------|--------|
| Base URL | `LOCAL_LLM_BASE_URL` or `http://127.0.0.1:1320/v1` |
| Model | `LOCAL_LLM_MODEL` or `gemini_free/gemma-4-26b-a4b-it` |
| Second opinion | `--model @cf/zai-org/glm-4.7-flash` |
| Scope | `full` (title + list fields + body) |
| Mode | `progressive` (default): each list field / body paragraph judged alone with **prior units as context**, then overall synthesis. Use `--mode whole` for a single pass. |

## Procedure (caller)

1. Resolve target path (user path, or focused file under `content/`). Prefer English `index.md` for v1.
2. Pick scope: `full` (default) | `list` | `body`. Selection paste is agent-side only: write temp text and pass that file, or scope to body and note the range in the user message after the report.
3. Run the script (do not hand-roll the API call unless the script fails):

```bash
python3 .cursor/skills/site-revise-prose/scripts/evaluate_prose.py \
  content/<section>/<slug>/index.md \
  --scope full
```

Optional: `--mode whole`, `--model <id>`, `--out docs/research/YYYY-MM-DD-local-prose-<slug>.md`, `LOCAL_LLM_API_KEY` if the gateway requires a key.

Progress shows on stderr (`Evaluating unit k/n…`). Surface the full stdout report (per-unit sections + **Overall synthesis**).

4. Surface the **full report** to the user. Do not silently rewrite the draft from the critique.
5. Hand off: apply ranked fixes with **revise-post** Steps 2, 3, 5 and **revise-flow**. Do not paste wholesale rewrites if the model invents them.

## Constraints

- MUST NOT edit `content/` as part of this skill.
- MUST quote evidence from the post; prefer fix **direction** over full rewrites.
- MUST fail loudly if the gateway is down or the model is disallowed.
- Spanish (`index.es.md`) and fact-checking are out of scope for v1.
