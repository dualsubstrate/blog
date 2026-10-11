---
name: site-explore-research-review
description: >-
  Gemma thinking gate for completed Perplexity threads before Explore further
  ship: score evidence, novelty, fit, counter-cost; approve or reject each
  candidate URL. Use after operator returns shared project links, before
  site-explore-link-hooks.
---

# Site explore research review

**Moral:** A completed Perplexity URL is a **candidate**, not a ship row. Gemma reviews research quality and hinge fit before any `label`/`hook` work.

## When to load

- Operator pasted one or more **shared Perplexity result links** from the essay project
- Before `site-explore-link-hooks` Call A
- When `explore_proposal.py` reports candidates needing human review

MUST NOT skip this gate because link copy already exists on a prior draft row. Re-review when research or explore inventory changed.

## Research packet (inputs)

Build one packet per essay extension batch. The agent or `explore_proposal.py` assembles:

| Field | Source |
| --- | --- |
| `essay_path` | English `index.md` (front matter + body) |
| `settled_thesis` | 2 bullets from agent read (not Gemma paraphrase of whole essay) |
| `target_count` | 2 or 3 approved `perplexity_thread` rows (pipeline default) |
| `excludes` | Existing `reader_landing.explore` labels, URLs, queries |
| `prompts` | Gemma-generated PRIMARY/BACKUP research prompts (if used) |
| `candidates[]` | Each: `id`, `url`, `hinge` (one sentence), `research_export` or `research_summary` |

MUST NOT review a candidate without readable research text (export markdown or operator summary bullets). MUST NOT invent thread content from the URL alone.

## Gemma review call (Polypus, thinking)

1. **Health:** `GET /health` (`.cursor/skills/ask-polypus/SKILL.md`).
2. **Model:** `gemini_free/gemma-4-26b-a4b-it` (optional `"chat_template_kwargs": {"enable_thinking": true}`; gemini_free may return `reasoning_details` without CF kwargs).
3. **Payload:** Essay excerpt (title + thesis bullets + close), full exclude list, editorial filters from `tmp/essay-extension-skills/perplexity-computer/PROJECT-INSTRUCTIONS.md`, and each candidate's hinge + research text + URL.
4. **Ask:** One structured block per candidate (YAML only in final output). Score 1-5: `evidence`, `novelty`, `source_fit`, `counter_cost_value`. Set `duplication_risk: low|medium|high`. Set `verdict: approve|reject` and `reason` (one sentence). Set `ship_priority: 1|2|3` for approved only (1 = ship first).
5. **Batch rule:** Approve at most `target_count` candidates. If more than `target_count` would approve, reject the weakest with reason `capacity`.

CORRECT:
```yaml
candidate_id: mgr-era
url: "https://www.perplexity.ai/search/aee56cb3-690e-4a6c-abd4-b3c5f41bf4f5"
hinge: "Permissionless execution vs managerial gatekeeping"
scores:
  evidence: 4
  novelty: 4
  source_fit: 5
  counter_cost_value: 3
duplication_risk: low
verdict: approve
ship_priority: 1
reason: "Cited mechanisms on cognitive offloading and post-industrial execution match the octopus close."
```

PROHIBITED:
```yaml
verdict: approve
reason: "Looks fine"
```

**CONSTRAINT:** Reject candidates that restate the essay thesis, recycle first-order clichés (GPS, pomodoro, generic productivity), or duplicate an excluded explore hinge.

- Enforcement: `verdict: reject` with explicit duplication or cliché reason
- Violation: STOP, re-run review with banned-topic list

## Handoff to link hooks

Only candidates with `verdict: approve` enter `site-explore-link-hooks` Call A/B. Rejected URLs MUST NOT appear in `reader_landing.explore`.

## Pre-completion checklist

- [ ] **Research text present:** Every candidate has export or summary
- [ ] **Structured verdict:** All four scores + verdict + reason per candidate
- [ ] **Capacity:** Approved count <= target_count (2 or 3)
- [ ] **At least one counter-cost:** When source assumed pure upside, an approved row OR a paired approved row must score counter_cost_value >= 3 OR review notes a dedicated counter-cost candidate
- [ ] **No agent override:** Verdict text came from Gemma (agent may parse YAML only)
