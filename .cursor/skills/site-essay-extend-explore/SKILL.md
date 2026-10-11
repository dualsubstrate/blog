---
name: site-essay-extend-explore
description: >-
  Ships extension research as Explore further links on a published post
  (reader_landing.explore: perplexity_thread, perplexity_query, Gemma hooks).
  Gemma thinking drafts Perplexity research prompts from the essay when asked.
  Use for follow-up searches, Perplexity threads, Dig deeper, or extension
  research that should become clickable links, not a new page.
---

# Site essay extend explore

**Moral:** The default extension deliverable is a row in **Explore further** on the source post, not a new Hugo page or a pasted essay in chat.

## Where this lands on the site

- **Section:** Reader landing panel on the post (i18n H3 **Explore further**), below **Go further** takeaways.
- **Template:** `layouts/partials/reader-landing.html` renders `Params.reader_landing.explore`.
- **Keep:** Existing `type: related` and `type: source` rows unless the operator asks to remove them.
- **Add or replace:** `perplexity_thread` (completed research URL) and/or `perplexity_query` (prefilled new search). Each perplexity row MUST have `label`, `hook`, and `url` or `query`.

MUST NOT create a sibling `content/` bundle for extension research unless the operator explicitly asked for a published companion page. MUST NOT dump research only in chat when the goal is to extend the article for readers.

CORRECT:
```yaml
reader_landing:
  explore:
    - type: related
      label: "..."
      path: "/x-minds/..."
    - type: perplexity_thread
      label: "..."
      hook: "One-line why click"
      url: "https://www.perplexity.ai/search/<real-id>"
    - type: perplexity_query
      label: "..."
      hook: "One-line teaser"
      query: "Full pasteable question?"
```

PROHIBITED:
```text
Save research as extension.md in tmp and never patch index.md explore
```

**Moral (explore rows):** Explore links are new research doors, not SEO keywords and not a summary of the piece.

Load when the operator wants follow-up searches, Perplexity queries, or Dig deeper explore items for a Hugo `reader_landing`.

## When to load

- User says follow-up searches, Perplexity, Dig deeper questions, explore queries
- User asks for a **new Perplexity prompt**, **research prompt**, or pasteable Perplexity question for an essay
- User points at `index.md` `reader_landing.explore`

MUST NOT load this skill to rewrite the essay body. MUST NOT treat a prefilled question as completed research.

## Perplexity research prompt (Gemma)

**CONSTRAINT:** When the operator asks for a new Perplexity research prompt (paste into Perplexity Computer or browser), the agent MUST NOT author that prompt. MUST run Polypus Gemma 4 on `gemini_free/gemma-4-26b-a4b-it` after `GET /health` (see `.cursor/skills/ask-polypus/SKILL.md`).

- **Payload:** Full English `index.md` (front matter + body); list of existing `reader_landing.explore` labels, thread URLs, and queries to exclude; post-industrial / novelty filter from `tmp/essay-extension-skills/perplexity-computer/PROJECT-INSTRUCTIONS.md` (research mode, no essay).
- **Ask Gemma:** One **PRIMARY_PROMPT** (2-4 sentences, hands off settled thesis, opens one unpaid hinge, requests bulleted cited research, plain US English) plus a **NO:** line (consensus clichés to reject). Optional **BACKUP_PROMPT** on a second hinge.
- **Deliver:** Paste Gemma blocks to the operator only. Thread URL and explore **label**/**hook** still follow `site-explore-link-hooks` after research completes.

- Enforcement: Turn log shows Polypus call before the prompt appears in chat; prompt text matches Gemma output (agent may strip leaked thinking, not rewrite substance)
- Violation: STOP, re-run Gemma; delete agent-authored prompt prose

CORRECT:
```text
Operator: new Perplexity prompt for Octopus Advantage
→ Gemma PRIMARY_PROMPT + BACKUP + NO lines
→ Operator runs Perplexity; later supplies thread URL
```

PROHIBITED:
```text
Agent writes a clever Perplexity paragraph without Gemma
```

## Research packet (extension batch)

Use one JSON or operator checklist per essay when running 2 or 3 Perplexity project results into Explore further.

| Field | Purpose |
| --- | --- |
| `essay_path` | English bundle `index.md` |
| `target_count` | 2 or 3 approved threads (default 3) |
| `excludes` | Labels, URLs, queries already on `reader_landing.explore` |
| `prompts.primary` / `prompts.backup` | Gemma-generated research prompts sent to the shared project |
| `candidates[]` | `id`, `url`, `hinge`, plus `research_export` (markdown path) or `research_summary` |
| `pipeline_calibration` | Optional `true`: score quality even when URL is already on explore (proposal dry run) |

Example: [tmp/explore-proposals/2026-09-26-the-octopus-advantage/candidates.json](/Users/hector/Xynova/ai/behaviourengineering/site/tmp/explore-proposals/2026-09-26-the-octopus-advantage/candidates.json)

**CONSTRAINT:** MUST NOT patch Hugo until the operator accepts a proposal from `scripts/explore_proposal.py` or explicitly says to apply. MUST NOT ship `perplexity_thread` URLs until the operator shared each thread in Perplexity as **Anyone with the link** and verified in incognito (`go tool task verify-explore-links`).

## Core constraints

**CONSTRAINT:** Each explore query MUST be a full question a reader can paste into Perplexity (one interrogative, concrete nouns from the piece). MUST NOT be a keyword string. MUST use plain US English and simple analogies instead of academic jargon.

- Enforcement: Every `query` ends with `?` or is a complete interrogative clause; jargon is replaced
- Violation: STOP, rewrite as a full plain-English question

CORRECT:
```yaml
query: "Does offloading executive function tasks to AI tools create a strategy migration where a person stops maintaining internal planning models?"
```

PROHIBITED:
```yaml
query: "metacognitive calibration failure in intention offloading"
```

**CONSTRAINT:** Each explore query MUST target a novel, non-obvious mechanism, asymmetric cost, or second-order effect for this post's hinge. MUST NOT be a generic search for the author's primary thesis. Strictly avoid consensus clichés (e.g., GPS navigation, generic productivity tips, surface summaries, or pop-psychology buzzwords).

- Enforcement: Every ship item explores a new hinge; zero consensus clichés
- Violation: STOP, add one novel door

CORRECT:
```yaml
label: "Premature curiosity closure"
query: "Does solving an organizational puzzle instantly with AI extinguish the dopamine incubation chamber needed to start the work?"
```

PROHIBITED:
```yaml
label: "Google effect on memory"
query: "How does the Google effect impact long-term memory for offloaded facts?"
```

**CONSTRAINT:** Each shipped `perplexity_thread` and `perplexity_query` MUST include Gemma-authored `label` and `hook` from `.cursor/skills/site-explore-link-hooks/SKILL.md` (Polypus Gemma 4, thinking). The agent MUST NOT draft link titles or sublines. MUST NOT patch explore YAML until Gemma Call B returns `approve: yes` (or the operator explicitly overrides).

- Enforcement: Diff shows Gemma-sourced label/hook only; review pass logged
- Violation: STOP, load `site-explore-link-hooks`; delete agent copy

**CONSTRAINT:** Each `label` MUST be 4–8 words with a curiosity gap (Gemma-owned). MUST NOT duplicate an existing `label` or `query` already on the page.

- Enforcement: Diff against current `reader_landing.explore`; word-count the label
- Violation: STOP, drop duplicates; shorten labels

CORRECT:
```yaml
label: "ADHD interest nervous system"
```

PROHIBITED:
```yaml
label: "How does divergent thinking interact with AI-driven workflows"
```

**CONSTRAINT:** A ship set for Hugo Explore further MUST be **2 or 3** `perplexity_thread` rows when using the Perplexity Explore pipeline (shared project batch). MAY add at most **one** `perplexity_query` follow-up door if a novel hinge remains open after threads ship. Legacy cap: 2–4 perplexity rows total if operator mixes queries and threads without the batch pipeline.

- Enforcement: Count approved threads; pipeline batch prefers threads over queries
- Violation: STOP, add counter-cost thread or trim via Gemma review

CORRECT:
```text
Ship 4: mechanism, generalized example, practical framework, offloading risk
```

PROHIBITED:
```text
Four restatements of "use AI for executive function"
```

**CONSTRAINT:** Explore further research MUST run inside the Perplexity project **behaviorengineering blog** (https://www.perplexity.ai/projects/72fa0438-7055-414e-8c7d-2c7642bb8ed5) so Instructions and Skills apply. MUST NOT use MCP `user-perplexity-browser` **`perplexity_research`** for Explore batch threads until the MCP can start threads in that project. Operator pastes Gemma prompts in the project ask box; agent MAY use **`perplexity_export`** on operator-supplied URLs. MUST attach the source thesis and existing explore queries in the Gemma prompt so the model does not duplicate them. MUST NOT invent `perplexity_thread` URLs.

- Enforcement: Prompt contains source hinge + existing queries; YAML `type: perplexity_thread` only with a real `url` from the operator or export
- Violation: STOP, strip invented threads; re-query with the exclusion list

CORRECT:
```yaml
- type: perplexity_query
  label: "Slow hunch incubation"
  query: "What is Steven Johnson's concept of the slow hunch and how do side projects turn into breakthroughs?"
```

PROHIBITED:
```yaml
- type: perplexity_thread
  url: "https://www.perplexity.ai/search/made-up-id"
```

**CONSTRAINT:** Spanish `index.es.md` MUST get native questions, not calques of the English interrogative, when a Spanish explore set is requested.

- Enforcement: Read ES queries for English syntax (`How does` calques, stacked infinitives)
- Violation: STOP, rewrite in native Spanish; keep the same hinge

CORRECT:
```text
¿Qué dice la neurociencia del sistema nervioso basado en el interés en el TDAH?
```

PROHIBITED:
```text
¿Cómo interactúa el pensamiento divergente con workflows de IA?
```

## Steps

1. **Inventory:** List existing `explore` labels and queries on EN (and ES if present).
2. **Thesis packet:** One short block: claim, examples, close, rows to exclude.
3. **Research prompt (if asked):** Gemma thinking pass per **Perplexity research prompt (Gemma)**; hand prompts to operator to paste in the **Perplexity project** ([PERPLEXITY-PROJECT.md](../site-extension-pipeline/PERPLEXITY-PROJECT.md)); stop until they return thread URL(s).
4. **Research packet:** Fill `candidates.json` (or equivalent) with 2 or 3 URLs plus export or summary text per candidate.
5. **Gemma research review:** Load `site-explore-research-review`; reject weak candidates before hooks.
6. **Proposal:** Run `go tool task explore-proposal` or `python3 scripts/explore_proposal.py`; operator reviews YAML proposal.
7. **Gemma ship set:** Load `site-explore-link-hooks` Call A (approved rows only) and Call B. Agent does not substitute its own link copy.
8. **Mechanical check:** Dedupe labels/queries, real thread URL, counter-cost row when needed.
9. **Apply (explicit):** Run `go tool task explore-apply` (applicability + operator question). Patch only after operator confirms and `APPLY=1 CONFIRM=1`.

## Themes that pay (use as a checklist, not a dump)

- Mechanism the piece named but did not explain
- Second example that tests whether the metaphor generalizes
- Practical framework for the close (side quest, first three steps)
- Counter-cost of the proposed fix
- Organizational / hiring scale (not only the individual)
- Tooling that matches interest-based attention (named systems, not "use ChatGPT")

## Pre-completion checklist

- [ ] **Gemma create + review:** Call A and Call B from `site-explore-link-hooks` completed; `approve: yes`
      Method: Confirm Polypus responses in turn log
      Pass: All label/hook from Gemma
      Fail: STOP, no YAML apply
- [ ] **Full questions:** Every ship `query` is pasteable (agent or research source, not Gemma paraphrase in `query` field)
      Method: Read each query
      Pass: Interrogative, specific
      Fail: STOP, fix query text
- [ ] **No dupes:** None match existing explore queries
      Method: Diff labels and queries
      Pass: Zero overlap
      Fail: STOP, replace
- [ ] **Counter-cost:** At least one ship item challenges the upside when the source assumed one
      Method: Mark the cost question or thread hinge
      Pass: Present when required
      Fail: STOP, add one
- [ ] **Ship count:** 2–4 perplexity rows total (`thread` + `query`)
      Method: Count perplexity types in explore
      Pass: 2–4
      Fail: STOP, trim or add via Gemma
- [ ] **No fake threads:** No invented Perplexity result URLs
      Method: Scan `type: perplexity_thread`
      Pass: URL from operator
      Fail: STOP, remove or convert to `perplexity_query`
