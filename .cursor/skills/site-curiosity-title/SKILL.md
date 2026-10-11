---
name: site-curiosity-title
description: >-
  Drafts or refines Hugo post titles with legitimate curiosity tension (real open
  loops the body pays off), not fake-urgency clickbait. Mandatory Polypus Gemma 4
  (thinking) for gap + candidates + cold-read review; agent formats the report and
  applies title only after user confirms. Example gold standard: "Why does every
  group organize into a pyramid?" Triggers: curiosity title, clickbait title (real
  gap + payoff), title ideas, retitling; pair with revise-hooks to audit. User may
  say skip gemma when the gateway is down.
---

# Curiosity title (legitimate tension, not bait)

## Purpose

Find **`title`** lines that make a stranger think *I don't know why that happens, and I want the answer* — then **deliver** that answer in the Claim, lead, or body.

**If the user says "clickbait":** they mean **this** skill (real open loop + payoff), **not** fake urgency, mystery boxes, or bait-and-switch. **Not** "You won't believe…" patterns.

**Legitimate curiosity:** a **specific gap** the post actually closes.

**Audit existing titles:** **`.cursor/skills/site-revise-hooks/SKILL.md`**.  
**This skill:** **generate and compare** candidates before you commit.

**Moral:** Title candidates and the recommended pick MUST come from **Gemma 4 (thinking on Polypus)**. The agent assembles the post packet, runs Gemma, and renders the **Required output format** table. The agent MUST NOT invent curiosity titles in its own voice except when the user opts out (**skip gemma**) or the gateway fails after retry.

## When to use

| Situation | Skill |
|-----------|--------|
| New post needs a title | **This skill** → then **revise-hooks** cold-read |
| Retitle after the essay is written | **This skill** (read body + Claim first) |
| Title already good; check list row | **revise-hooks** only |
| Full publish pass | **revise-post** (hooks phase uses **revise-hooks**) |

## Gold-standard example (this repo)

**Post:** `content/human-condition/2026-05-28-why-humans-keep-building-pyramids/`

**Title:** `🏢🔺 Why does every group organize into a pyramid?`

**Why it works:**

| Test | Pass? |
|------|-------|
| **Prediction error** | Reader assumes groups can stay flat; title says they become a pyramid anyway. |
| **Concrete image** | *pyramid* (not "hierarchy" alone). |
| **Universal scope** | *every group* (family, company, forum) without listing them. |
| **Plain verb** | *organize into* (readable cold; not "map complexity flat"). |
| **Question opens a loop** | *Why* — gap the essay closes. |
| **Payoff exists** | Claim + body explain layers, rank, attractor, counter-design. |
| **Not a mystery box** | You know the topic (groups, pyramid shape); you don't know the **why**. |

**Pairing:** Question in **`title`**; **`description` (Claim)** stays **assertive** (answers the why). Do not put the only question in the Claim.

## Gemma 4 (Polypus, thinking)

**CONSTRAINT:** MUST run Gemma before presenting title candidates to the human. MUST use **`.cursor/skills/ask-polypus/SKILL.md`** (health check, loopback `:1320`, prefixed model id).

| Setting | Value |
|---------|--------|
| Model | `gemini_free/gemma-4-26b-a4b-it` |
| Thinking | `"chat_template_kwargs": {"enable_thinking": true}` |
| `max_tokens` | `8192` minimum on thinking calls (thinking can consume the budget) |

**Skip without asking:** user says **`skip gemma`**, **`no gemma`**, or **`without gemma`**.

**Gateway down / empty reply:** Probe health; one retry with a shorter artifact (title, `description`, body spine bullets, `###` headings only). If still empty or `finish_reason: length`, report failure and continue with **agent-only** candidates only after stating Gemma did not run.

### Call A — Gap + candidates (create)

1. **Health:** `GET ${POLYPUS_BASE_URL:-http://127.0.0.1:1320}/health`.
2. **Payload:** Current `title`, full `description` / Claim / lead, `type`, section folder, body spine (agent one-liner), body `###` headings (lines only), gold-standard reminder (pyramids question title), pattern table **A–E**, emoji rules for the section, banned bait list from **Legitimate vs bait**.
3. **Ask Gemma:** State **Reader believes / Article shows / Gap** in plain English; draft **5** titles using **at least three different patterns** (A–E); tag each with pattern letter; note optional **1–2** leading emoji only when section allows.
4. **Output shape (MUST):** Ask Gemma to end with a **paste-ready block** only (no prose after it). Prefer line-oriented keys (thinking models often miss fenced YAML):

```text
GAP_BELIEVES=
GAP_SHOWS=
GAP=
CANDIDATE_1_PATTERN=A
CANDIDATE_1_TITLE=
CANDIDATE_2_PATTERN=
CANDIDATE_2_TITLE=
...
CANDIDATE_5_PATTERN=
CANDIDATE_5_TITLE=
```

5. **Parse:** Strip leaked thinking wrappers; take the **last complete** `CANDIDATE_*` set in the response. If missing, re-call Call A once with “output ONLY the KEY= lines.”

### Call B — Cold-read + recommend (review)

1. Same model and thinking flag; `max_tokens` `8192`.
2. **Payload:** Gap lines + all Call A titles + current `description` / Claim (first screen).
3. **Ask Gemma:** For each candidate, Pass/Fail on cold-read (subject + image, open gap, no bait trick); score **legitimate curiosity** vs **mystery box**; pick **recommended** `#1`–`#5` with one-line why; **`claim_pairing`:** does the lead answer the recommended title? (`yes` / `no` + one fix direction if `no`).
4. **Output shape:**

```text
REVIEW_1=Pass|Fail — …
...
RECOMMENDED_NUM=3
RECOMMENDED_WHY=
CLAIM_PAIRING=yes|no — …
APPROVE=yes|no
```

5. **Apply revisions:** If `APPROVE: no`, merge Gemma’s revised titles from `REVIEW_*` lines and re-run Call B once, or stop and show blockers.

### Agent duties (not Gemma)

- Read the post and build the packet (workflow step 1).
- Run Call A then Call B; map Gemma keys into **Required output format** (table + Recommended + Claim pairing).
- Change **`title`** in front matter only after user confirms (`y` / pick `#` / `cancel`).
- For **`index.es.md`**, after EN title is chosen, MAY run a **third** short Gemma call (thinking optional) for a native Spanish `title` line only; MUST NOT calque the English question word-for-word if Spanish idiom needs a different shape.

CORRECT:
```text
User: /site-curiosity-title on content/.../index.md
→ health OK → Call A KEY= block → Call B APPROVE=yes → render markdown table → ask to apply
```

PROHIBITED:
```text
Agent drafts five titles without Polypus because “the essay is clear”
```

## Workflow (MUST follow)

### 1. Read the post first

- **`type`** skill (`claims-content`, `video-content`, `cognitive-memetics-content`, etc.).
- **`description`** / Claim / lead: what answer does the page **owe** the reader?
- Body **spine** (one sentence): what mechanism or scene is the real payoff?
- Section **`##`** hooks: do not copy verbatim; title can rhyme with them.

### 2. Name the curiosity gap (one line)

Use **Gemma Call A** `GAP_*` lines in the report. If Gemma cannot state a gap, stop curiosity mode and recommend a **direct thesis title** (pattern **C**) via Call A with that instruction.

```text
Reader believes: …
Article shows: …
Gap: …
```

**Fail** if you cannot state the gap in plain English. No gap → no curiosity title; use a **direct thesis title** instead (see **revise-hooks**).

### 3. Draft 3–5 title candidates (Gemma Call A)

**CONSTRAINT:** Candidates MUST come from **Call A**, not agent improvisation (unless **skip gemma** / gateway failure path above).

Use **different patterns** (not five rewrites of the same joke):

| Pattern | Template | Use when |
|---------|----------|----------|
| **A. Gap question** | Why does [concrete subject] [surprising verb] [concrete image]? | Universal pattern, strong cold-read (pyramids model). |
| **B. Paradox / contrast** | [X looks flat] until [Y] | Two beats; good for claims. |
| **C. Direct thesis** | [Actor] [verb] [mechanism] | When question would be vague. |
| **D. Stakes** | What [cost] when [mechanism] | Only if body proves the cost. |
| **E. Scene hook** | [Concrete scene] — [turn] | When one image carries the post. |

**MUST:**

- One **concrete noun** from the piece (pyramid, org chart, mod, feed, etc.).
- **Active verb** where possible (*organize into*, *redraw*, *track*, not *engagement dynamics*).
- **US English** in `content/` default pages.
- **Truthful:** title claim must be defended in the first screen (Claim or lead).

**MUST NOT:**

- "You won't believe…", "The truth about…", "Everything you know about…"
- Questions with **no** imaginable answer in the piece ("What if reality is a simulation?" unless the post is about that).
- Jargon the card cannot decode (*predictive processing*, *heterarchy*) without a plain noun.
- **Mystery box:** withhold **what the post is about** (bad: "This changes how we think about groups").
- Duplicate the **Claim** sentence as the title.

### 4. Emoji (section rules)

- **`human-condition`**, **`social-protocols`**, **`mind-infrastructure`**, **`x-minds`:** MAY use **1–2** leading emoji if they signal the hook.
- **`cognitive-memetics`:** **no** leading emoji in **`title`** (use **`heading_code`** when needed).

### 5. Cold-read each candidate (Gemma Call B)

Map **Call B** `REVIEW_*` lines into the table **Cold-read** column. The human-facing tests are still:

1. What is this about? (subject + image)
2. What don't I know yet? (the gap)
3. Would I feel **tricked** after opening? (if yes → reject)

**Pass** if a stranger can answer 1 and 2 and 3 is **no**.

### 6. Pick one and check Claim pairing (Gemma Call B)

| Field | Job |
|-------|-----|
| **`title`** | Opens the gap (often a **question**). |
| **`description` (Claim)** | **Answers** with mechanism + stakes (assertive prose). |

**Reject** a title that forces the Claim to repeat the question without answering.

**Claim fog check:** After picking a title, read the **Claim** per **`.cursor/skills/site-revise-hooks/SKILL.md`** → **Claim fog**. Title can be perfect while a Claim sentence still fails cold-read; fix the Claim before publish.

## Required output format

```markdown
## Curiosity titles: [bundle or path]

**Gap:** Reader believes … / Article shows … / Gap: …

| # | Pattern | Title | Cold-read |
|---|---------|-------|-----------|
| 1 | A Gap question | … | Pass / Fail + why |
| 2 | … | … | … |

**Recommended:** #N — [one line why]

**Claim pairing:** [does current Claim answer the title? yes/no + one fix if no]

Apply recommended title to front matter? (y / pick # / cancel)
```

**Default:** analysis only; change **`title`** only after user confirms (stable URL: do not rename bundle folder unless user asks).

**CONSTRAINT:** The **Gap** row and every **Title** cell in the table MUST trace to Gemma Call A/B output unless **skip gemma** was invoked.

## Pre-completion checklist

- [ ] **Health checked:** Polypus returned 2xx before chat
      Method: Inspect health probe
      Pass: ok
      Fail: STOP or **skip gemma** path documented
- [ ] **Call A ran:** Five patterned candidates with `GAP_*` lines
      Method: Trace titles to Gemma KEY block
      Pass: All five from Gemma
      Fail: STOP, re-run Call A
- [ ] **Call B ran:** Cold-read + `RECOMMENDED_NUM` + `CLAIM_PAIRING`
      Method: Read `APPROVE` / revisions
      Pass: `APPROVE: yes` or operator override stated
      Fail: STOP, re-run Call B or show blockers
- [ ] **Output format:** Human sees the markdown table + apply question
      Method: Skim reply
      Pass: Matches **Required output format**
      Fail: STOP, reformat
- [ ] **No unsolicited apply:** `title` front matter unchanged until user confirms
      Method: `git diff` on target bundle
      Pass: No title edit unless user said `y` / pick `#`
      Fail: STOP, revert title change

## Legitimate vs bait (quick)

| Legitimate | Bait |
|------------|------|
| Specific *why* you can answer in the post | Vague *what* you'll "reveal" |
| Concrete noun (pyramid, VP, mod) | Abstract label (systems, paradigms) |
| Reader learns something true | Reader feels manipulated |
| Claim delivers on the title | Claim is unrelated hype |

## Related skills

- **`.cursor/skills/ask-polypus/SKILL.md`** — health, model id, consult hygiene
- **`.cursor/skills/site-revise-hooks/SKILL.md`** — audit, integrity, body headings
- **`.cursor/skills/site-claims-content/SKILL.md`** — Claim must stay assertion
- **`.cursor/skills/site-revise-post/SKILL.md`** — full lot after title is set (hooks phase)
- **`.cursor/skills/site-curiosity-title/examples.md`** — more before/after pairs (optional read)
