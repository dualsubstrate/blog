---
name: site-extension-pipeline
description: >-
  Blog-side workflow for essay extensions: Gemma research prompts, Perplexity
  thread, Gemma explore hooks, optional companion draft. Use when formalizing
  extension work before content-pipelines MCP owns it.
---

# Site extension pipeline (blog operator)

**Moral:** Extension research ships as **Explore further** links on the source post; companion prose is optional.

Human-driven loop on the Hugo host (future: content-pipelines MCP verbs).

## Default deliverable: Explore further

| Output | Where | When |
| --- | --- | --- |
| Completed Perplexity research | `reader_landing.explore` → `type: perplexity_thread` + `url` + Gemma `label`/`hook` | Operator pasted a real thread URL |
| Follow-up question door | `reader_landing.explore` → `type: perplexity_query` + `query` + Gemma `label`/`hook` | One novel hinge still open |
| Companion essay | New draft or bundle | Operator explicitly asked for essay/sequel |

The on-page section title is **Explore further** (EN) / explore copy in ES. Template: `layouts/partials/reader-landing.html`.

## Loop

1. **Source:** Read English `index.md`; name settled thesis and one unpaid hinge.
2. **Research packet:** Assemble essay path, 2 or 3 target hinges, Gemma PRIMARY/BACKUP prompts, existing explore excludes, and candidate rows (`id`, `url`, `hinge`, export or summary). See **Research packet** in `site-essay-extend-explore` and `site-explore-research-review`.
3. **Research prompt (when operator asks):** MUST NOT draft the Perplexity paste prompt in the agent voice. Polypus Gemma 4 with `enable_thinking: true`. Deliver **PRIMARY_PROMPT**, optional **BACKUP**, and **NO:** lines. Operator runs the shared Perplexity project and returns result URLs.
4. **Research:** Perplexity **project** only ([PERPLEXITY-PROJECT.md](PERPLEXITY-PROJECT.md)). Operator starts each thread from https://www.perplexity.ai/projects/72fa0438-7055-414e-8c7d-2c7642bb8ed5 (Instructions + Skills). MUST NOT use MCP `perplexity_research` for Explore batches until project-scoped research exists. Bullets, plain English, real citations. No essay unless asked.
5. **Thread URLs:** Operator supplies one URL per candidate (2 or 3 total for a full set). Each thread MUST be **shared public** in Perplexity (**Anyone with the link**) before ship; see **Public share (Perplexity)** below.
6. **Gemma research review (mandatory):** `site-explore-research-review` on each candidate before hooks. Reject weak or duplicate research.
7. **Explore ship set:** Approved candidates only; `site-essay-extend-explore` mechanical dedupe and counter-cost check.
8. **Gemma link hooks (mandatory):** `site-explore-link-hooks` Call A + Call B. Agent MUST NOT write labels or hooks.
9. **Proposal or apply:** Default `scripts/explore_proposal.py` (proposal only). Patch `reader_landing.explore` only after operator accepts the proposal.
10. **Companion (optional):** `site-essay-extend` only on explicit request.

## Proposal runner (no Hugo write)

**Runner:** `go tool task explore-proposal` → [`scripts/explore_proposal.py`](/Users/hector/Xynova/ai/behaviourengineering/site/scripts/explore_proposal.py) (implementation under `.cursor/skills/site-extension-pipeline/scripts/`).

```bash
go tool task explore-proposal ESSAY=content/x-minds/2026-09-26-the-octopus-advantage/index.md \
  CANDIDATES=tmp/explore-proposals/2026-09-26-the-octopus-advantage/candidates.json \
  ALLOW_SUMMARY=1
```

By default the runner **syncs `research_export` paths** and **requires** markdown export files. Use `ALLOW_SUMMARY=1` only when the operator has not yet exported threads.

### Perplexity exports (MCP, before proposal)

1. `go tool task explore-fetch-exports ESSAY=... CANDIDATES=... SYNC=1`
2. Read `tmp/explore-proposals/<slug>/export-manifest.json`
3. **Agent:** for each pending row, call `user-perplexity-browser` **`perplexity_export`** with `thread_id`, `save_dir=<repo>/<save_dir>`, `format=markdown` (one attempt per row; no retry loop)
4. Ensure each file exists at `research_export` path, then re-run `go tool task explore-proposal` (without `ALLOW_SUMMARY=1` when exports are present)

Helper: `go tool task explore-fetch-exports` / [`scripts/explore_fetch_exports.py`](/Users/hector/Xynova/ai/behaviourengineering/site/scripts/explore_fetch_exports.py).

## Share manifest (before apply)

**Runner:** `go tool task explore-share-prepare` / `explore-share-confirm` / `explore-share-verify-cold` / `explore-share-check` → [`scripts/explore_share.py`](/Users/hector/Xynova/ai/behaviourengineering/site/scripts/explore_share.py).

1. After `candidates.json` has URLs: `go tool task explore-share-prepare CANDIDATES=...`
2. Operator sets **Anyone with the link** in Perplexity (see [PERPLEXITY-SHARE-BROWSER.md](PERPLEXITY-SHARE-BROWSER.md)); `go tool task explore-share-confirm CANDIDATES=... CANDIDATE_ID=...`
3. Incognito pass (+ optional `PROBE=1`): `go tool task explore-share-verify-cold ... OPERATOR=1`
4. `go tool task explore-share-check PROPOSAL=... CANDIDATES=...` must pass before `explore-apply`.

## Apply runner (operator gate)

**Runner:** `go tool task explore-apply` → [`scripts/explore_apply.py`](/Users/hector/Xynova/ai/behaviourengineering/site/scripts/explore_apply.py).

1. Share preflight (manifest) runs automatically unless `SKIP_SHARE_PREFLIGHT=1` (tests only).
2. `go tool task explore-apply ESSAY=... PROPOSAL=tmp/.../slug.proposal.yaml CANDIDATES=.../candidates.json`  
   Gemma **applicability** review writes `<proposal>.applicability.yaml` and prints `operator_question`.
2. **Agent MUST** surface that question to the operator (AskQuestion or plain ask). MUST NOT pass `APPLY=1` in the same turn as the review.
3. Only if the operator confirms:  
   `go tool task explore-apply ESSAY=... PROPOSAL=... APPLY=1 CONFIRM=1`

MUST NOT patch `index.md` / `index.es.md` without step 2 and `CONFIRM=1`.

## Public share (Perplexity)

**CONSTRAINT:** Every `perplexity_thread` URL on the live site MUST be **viewable by anyone with the link** (Perplexity Share → public / anyone with link). MUST NOT ship workspace-private threads that show a login wall to readers.

- **Before explore-apply:** `go tool task explore-share-check` (manifest) for proposal URLs; after ship, `go tool task verify-explore-links POST=section/slug`.
- **Reference:** [PERPLEXITY-SHARE.md](PERPLEXITY-SHARE.md), [PERPLEXITY-SHARE-BROWSER.md](PERPLEXITY-SHARE-BROWSER.md)

CORRECT:
```text
Share thread → Anyone with the link → incognito opens full research → paste URL in YAML
```

PROHIBITED:
```text
Copy URL from address bar while thread is still private to your account
```

## Perplexity Computer assets

Project instructions and skill zips: `tmp/essay-extension-skills/perplexity-computer/` (paste into project Settings; upload skills).

**Explore research home:** [PERPLEXITY-PROJECT.md](PERPLEXITY-PROJECT.md) (mandatory project URL; MCP `perplexity_research` not for Explore batch).

## Pre-ship checklist

- [ ] **Project scope:** Every candidate thread was started from [PERPLEXITY-PROJECT.md](PERPLEXITY-PROJECT.md) (not global Search / MCP `perplexity_research`)
- [ ] Gemma research review: every candidate has verdict (`.cursor/skills/site-explore-research-review/SKILL.md`)
- [ ] Only approved URLs in ship set (2 or 3 `perplexity_thread` target for full pipeline batch)
- [ ] Gemma Call A + Call B complete for hooks (`approve: yes`)
- [ ] Real thread URL only (operator supplied)
- [ ] Every perplexity row has Gemma `label` + `hook`
- [ ] No em dash in hooks or labels
- [ ] **Public share:** `share-manifest.json` passes `go tool task explore-share-check` (confirm + cold verify per URL)
      After Hugo apply: `go tool task verify-explore-links POST=...`
      Pass: Cold reader sees research
      Fail: STOP, re-share in Perplexity before deploy
- [ ] ES hooks native, not calque
