---
name: site-explore-link-hooks
description: >-
  Mandatory Polypus Gemma 4 (thinking) pass for all Explore further link copy:
  create label/hook, review click-worthiness, and approve ship set before YAML
  apply. Use when shipping or refreshing perplexity_thread or perplexity_query
  rows on reader_landing.explore.
---

# Site explore link hooks

**Moral:** The agent MUST NOT invent Explore further link titles or sublines. Gemma 4 (thinking on Polypus) creates and reviews every `label` and `hook` before YAML lands on the post.

## Explore further (site section)

- **Not** a separate page, Substack sidecar, or chat-only research dump.
- **Is** the bulleted list under **Explore further** in the reader landing panel on the published post.
- **Renders as:** Purple link = `label`; gray subline = `hook` (when set), else generic Perplexity hint strings from i18n.
- **Thread row:** Reader opens your completed Perplexity answer set (`type: perplexity_thread` + real `url`).
- **Query row:** Reader opens a new search with the full question prefilled (`type: perplexity_query` + `query`).

Load after explore targets exist (real thread URL and/or full pasteable query). MUST load only for **research-review-approved** candidates (see `site-explore-research-review`). MUST load before any patch to `reader_landing.explore`.

**CONSTRAINT:** The agent MUST NOT write or edit `label` or `hook` strings without a Gemma thinking pass on Polypus. The agent MAY only paste Gemma-approved YAML into front matter, plus operator-supplied `url` and `query` text.

- Enforcement: Every label/hook in the diff came from the latest Gemma response (or a Gemma revision pass)
- Violation: STOP, run Gemma; delete agent-authored link copy

MUST NOT invent `perplexity_thread` URLs. MUST NOT replace the `query` field with the hook (hook is display copy; query stays the full interrogative for the prefilled URL).

## Front matter shape

```yaml
- type: perplexity_thread
  label: "4-8 words, curiosity gap"
  hook: "One sentence, max ~20 words, names the surprise mechanism"
  url: "https://www.perplexity.ai/search/<real-id>"
- type: perplexity_query
  label: "4-8 words"
  hook: "One sentence teaser (may echo the question, not keyword pile)"
  query: "Full question for Perplexity prefilled search?"
```

Host renders `hook` under the link instead of generic i18n hints (`reader-landing.html`).

## When to load

- Any Explore further row will be added, replaced, or refreshed
- User says ask Gemma for link titles, hooks, review explore links, or ship set
- Finishing `site-essay-extend-explore` (mandatory gate before apply)

## Gemma consult (Polypus, thinking): create then review

Run **two** calls when shipping (or one combined call with two sections in the prompt).

### Call A: Create

1. **Health:** `GET ${POLYPUS_BASE_URL:-http://127.0.0.1:1320}/health` (see `.cursor/skills/ask-polypus/SKILL.md`).
2. **Payload:** Post title, settled thesis (2 bullets), unpaid hinge, thread summary bullets, each full `query` string, existing explore rows to avoid duplicating, banned first-order clichés for this topic.
3. **Model:** `gemini_free/gemma-4-26b-a4b-it` (optional `"chat_template_kwargs": {"enable_thinking": true}`).
4. **Ask:** Propose 2-4 perplexity rows to ship (`thread` and/or `query`). For each row output `label` + `hook`. If `index.es.md` exists, output parallel `es` blocks (native Spanish, not calque).
5. **Parse:** Strip leaked thinking; keep YAML only. Re-call if output is not paste-ready.

### Call B: Review (required before apply)

1. Same model and thinking flag.
2. **Payload:** Paste the proposed EN (and ES) labels, hooks, and the full `query` strings. Ask Gemma to score click-worthiness, novelty, and plain language; reject topic-tag labels and generic Perplexity boilerplate.
3. **Ask:** Return `approve: yes|no`, `revisions:` with replacement label/hook only where weak, or `blockers:` if a row should not ship.
4. **Apply revisions:** If `approve: no`, merge Gemma revisions and re-run Call B until `approve: yes` or the operator overrides.

### Agent duties (not Gemma)

- Attach real `url` for `perplexity_thread` from the operator only.
- Keep full `query` strings accurate for prefilled search URLs.
- Patch `reader_landing.explore` after Gemma approval.

## Hook constraints (enforce in prompt and review)

- Plain language; one concrete surprise (dessert before dinner, false finish, one-way gate).
- No academic jargon walls, no em dash (U+2014), no "Opens Perplexity" boilerplate.
- Label is not the mechanism name alone; it promises a tension the reader might feel.
- Hook is not a summary of the whole post; it names why **this** link exists.

## Steps in the extension workflow

1. Research (Perplexity project): factual bullets; operator supplies real thread URL(s) in the research packet.
2. **`site-explore-research-review`:** Gemma scores and approves or rejects each candidate.
3. **This skill (Call A):** Gemma creates ship set + every `label` + `hook` (EN/ES) for approved rows only.
4. **This skill (Call B):** Gemma reviews and revises until approve.
5. `site-essay-extend-explore`: agent verifies dedupe, counter-cost, real URLs only (no copywriting).
6. Proposal via `explore_proposal.py` or patch YAML after operator OK; verify **Explore further** on Hugo.

## Pre-completion checklist

- [ ] **Gemma created all copy:** No agent-written label/hook in the diff
      Method: Trace each string to Gemma output
      Pass: 100% Gemma-authored link text
      Fail: STOP, run Call A
- [ ] **Gemma review approved:** Call B returned `approve: yes` or operator accepted blockers
      Method: Read review response
      Pass: Explicit approve
      Fail: STOP, revise and re-run Call B
- [ ] **Hooks present:** Every shipped perplexity row has `hook`
- [ ] **Query intact:** Prefilled URL still uses full `query` field
      Method: Read YAML
      Pass: query unchanged
      Fail: STOP, move teaser text to `hook` only
- [ ] **No em dash:** U+2014 absent in label and hook
      Method: Search fields
      Pass: Zero
      Fail: STOP, replace
