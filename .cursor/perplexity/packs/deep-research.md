# Pack — Deep Research (behaviorengineering.ai content)

Workflow: `perplexity-browser-research` overlay `.cursor/perplexity/`. Mode: **deep** (use `search` only for a quick single-fact lookup).

Fill every `[bracket]` section from the target post and repo files before submit. Copy the **whole** block into `perplexity_research` `prompt`.

---

## Context

External research for **behaviorengineering.ai**: a Hugo (LoveIt) site publishing essays, curated video picks, and cognitive-memetics strips. Research supports **accurate public prose**, not Go tooling or theme internals.

**Repo:** `github.com/xynova/behaviour-engineering` (local folder: `site`)

**Site sections and types** (pick what applies):

| Section | Typical `type` | Research usually needed for |
| --- | --- | --- |
| `social-protocols/` | `claims`, `video` | Norms, institutions, game theory, media, coordination |
| `human-condition/` | `claims`, `video` | Psychology, morality, identity, person-level mechanisms |
| `x-minds/` | `claims`, `video` | Mixed-wiring lives, parents, community (public name **X-Minds**; do not use neurodivergence in titles or teasers) |
| `mind-infrastructure/` | `video` (often) | Models, tools, cross-topic picks |
| `cognitive-memetics/` | `panel`, `sayings` | Rarely Perplexity; satire/sayings unless fact-checking a cited study |

**Relevant paths reviewed before this pack** (check all that apply):

- Target bundle: `[content/<section>/<slug>/index.md]`
- Spanish sibling: `[index.es.md if present]`
- Sidecars: `[substack.md, linkedin.txt, tmp/articles/<topic>.md, etc.]`
- Type skill: `[.cursor/skills/site-claims-content/SKILL.md | site-video-content/SKILL.md | …]`
- Placement: `.cursor/rules/site-content-placement.mdc`
- Prior research: `[docs/research/*.md if any]`

---

## Research intent

Pick **one** primary intent (delete the others):

- [ ] **Fact-check** list copy and body claims before publish (video pick or claims post)
- [ ] **Grounding hunt** for a new or revised `type: claims` post (papers, definitions, dates)
- [ ] **Source map** for a long interview/podcast (attribute claims to primary literature)
- [ ] **Topic explore** for a post not yet drafted (mechanism, debates, best citations)

**Target post (if any):** `[path, title, type, section]`

**Why now:** `[one sentence: what decision this research unlocks]`

---

## Question

[State the research question in plain English. Be specific: name the speaker, study, construct, or controversy.]

**Claims or hooks to validate** (one per line; paste from `description`, `grounding`, `sowhat`, or body):

1. [claim]
2. [claim]
3. [claim]

**Optional follow-ups:**

- [e.g. best primary citation for each supported claim]
- [e.g. what is consensus vs speaker's framework vs podcast rhetoric]
- [e.g. human vs animal vs in-vitro evidence]

---

## Constraints

**Evidence quality**

- Prefer **peer-reviewed papers**, preprints from named labs, and **official institutional pages** over podcast summaries, YouTube descriptions, or secondary explainers.
- Distinguish **established science**, **the speaker's interpretive framework**, and **interview rhetoric**.
- Note study design: in vitro, animal, observational, RCT, postmortem, meta-analysis.
- Do not invent facts about this codebase or post; flag gaps explicitly.
- Prefer material from **2010–2026** unless the question is historical.

**Editorial (how findings will be used)**

- Prose is **US English**, accessible mechanism, **no em dashes** (U+2014).
- Do not overstate causality; correlation and small-n studies need explicit caveats.
- For `type: claims`: Grounding needs a **Source** line with Markdown link (`[title](url)`), not bare URLs. See `.cursor/skills/site-claims-content/SKILL.md`.
- For `type: video`: `description` is the list lead; body is a TLDR for text-first readers. See `.cursor/skills/site-video-content/SKILL.md`.
- Perplexity output is **research input only**; a human validates before edits ship.

**Out of scope** (unless the question says otherwise)

- Hugo theme internals, `go.mod`, Substack pipeline, social autopost tooling
- Spanish translation (research in English; adapt later via `revise-post-es`)
- Prose quality / AI voice (use `.cursor/perplexity/packs/prose-review.md` instead)

---

## Repo excerpts (paste below)

### Excerpt A — List-facing copy (Claim / lead / sowhat)

```text
[Paste description, grounding, sowhat, and/or title from front matter]
```

### Excerpt B — Body claims or draft spine

```text
[Paste relevant body paragraphs, bullet claims, or tmp/articles/ notes]
```

### Excerpt C — Source already in hand (optional)

```text
[Paste any URL, paper title, youtube_id, or timestamp the post already cites]
```

---

## Output format (request from Perplexity)

Please structure your answer as:

1. **Summary** (5–10 bullets)
2. **Claim-by-claim verdict table** with columns: Claim | Verdict (supported / partial / unsupported / unclear) | Best primary source | Nuance for editors
3. **Recommendation** for this post type:
   - **Keep** (well supported hooks)
   - **Soften** (frame as hypothesis, correlation, or speaker model)
   - **Cut or fix** (likely errors, unsourced numbers)
4. **Grounding candidates** (if `type: claims`): 1–3 Source lines ready for front matter (`[Paper or article title](url)` plus one-line digest)
5. **Risks / tradeoffs** (overstatement, missing context, conflicting studies)
6. **Sources** (links; label peer-reviewed vs preprint vs blog/podcast)
7. **Open questions** for the author

---

## After export (repo handling)

1. Export via `perplexity_export`; note path under `~/.perplexity-browser-mcp/exports/`.
2. Create `docs/research/YYYY-MM-DD-perplexity-<slug>.md` with: pack path, thread URL, export path, summary bullets, edit priorities, open questions.
3. After **human review**, apply edits via the matching type skill and **`.cursor/skills/site-revise-post/SKILL.md`** (full lot; or focused `revise-format` alone).
4. Do **not** commit raw Perplexity exports or auto-merge research prose into `content/` without validation.
