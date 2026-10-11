# Pack — Prose review (human read, AI artifacts)

Workflow: `perplexity-browser-research` overlay `.cursor/perplexity/`. Mode: **search** (default). Use **deep** only for a very long piece (e.g. full Substack draft).

**Not fact-checking.** Do not verify citations, numbers, or science. This pack judges **readability, flow, voice, and AI slop** only.

Fill every `[bracket]` section from the target post before submit. Copy the **whole** block into `perplexity_research` `prompt`.

---

## Context

Editorial cold-read for **behaviorengineering.ai**: essays and video TLDRs that should read like a sharp human wrote them, not a polished AI draft.

**Repo:** `github.com/xynova/behaviour-engineering` (local folder: `site`)

**Target post:** `[content/<section>/<slug>/index.md — title, type]`

**Scope of this review** (pick one; delete others):

- [ ] **Full post** (list copy + body + sidecar if pasted)
- [ ] **List-facing only** (`description`, `grounding`, `sowhat`, title)
- [ ] **Body only** (Thoughts / TLDR sections)
- [ ] **Selection:** `[paste specific ### sections or paragraph range]`

**Reader persona:** A smart stranger skimming on a phone. They have **not** watched the video or read the paper. They will bounce if the prose feels synthetic, preachy, or hollow.

**After this review:** A human applies fixes via `.cursor/skills/site-revise-post/SKILL.md` (Steps 2, 3, 5) and `.cursor/skills/site-revise-flow/SKILL.md`. Perplexity diagnoses; it does not replace those passes.

---

## Question

How well does this prose read for **humans**? Where does it **lose trust, rhythm, or clarity**? Where does it sound **AI-generated, teacherly, or over-polished**?

**What I need from you:**

1. **Cold-read verdict** in one paragraph: would you keep reading? Where do you stall?
2. **Quoted failures only** (exact phrases from the paste). No vague "some sentences feel off."
3. **What already works** (2–5 lines worth keeping verbatim; say why).
4. **Top 5 fixes** ranked by impact on human trust and flow.
5. **Do not rewrite the whole essay** unless a single sentence is unusable; give **direction**, not a full replacement draft.

---

## Constraints

**Your role**

- Act as a **skeptical magazine editor**, not a writing coach and not a fact-checker.
- Penalize prose that is **smooth but empty** (sounds good, says little).
- Penalize **symmetrical punch-line rhythm** (five short sentences of equal weight).
- Reward **named subjects, concrete mechanisms, and one clear metaphor thread**.
- US English expectations. **No em dashes** (U+2014) in suggested rewrites.

**AI artifacts to flag** (quote when found)

| Pattern | Example shape |
| --- | --- |
| Teacher / editor framing | "You will learn...", "The useful move is...", "X's hook is..." |
| Bridge verbs | "bridges X and Y", "ties together", "maps onto" |
| Semicolon aphorisms | "Stress drains; purpose concentrates." |
| Negation-first fluff | "It is not X; it is Y", "Not A but B" |
| Meta-commentary | "This makes the point clear", "The pattern matters", "Still," as a punch pivot |
| Caption pivots | "Picard is careful here.", "Worth noting:" |
| Classic AI rhetoric | "This changes how we understand...", "Classic experiments reveal..." |
| Clever filler metaphors | Opaque one-liners the reader must decode ("not a label on your calendar") |
| Mixed metaphors | Budget + meter + ledger in one paragraph |
| Staccato thesis stacks | Same-length punch lines every sentence |
| Rhetorical noun fragments | "The mirror." / "The wound." as standalone beats |
| Revelation stacks | Several one-line poetic restatements of one idea, no mechanism |
| Metaphor-shell restack | After a mechanism sentence: "The blend is your experienced reality…" then "sealed file / steering wheel." Complete sentences that only relabel. Full test: **`.cursor/skills/site-revise-post/reference.md`**. |
| Industry-verb shells | "sells the blend"; "steer the rewrite"; "keeps assembling it"; "next assembly"; "the same room"; "template can harden"; "name the box… and the box." Fail even as a single clause. Full row: **`.cursor/skills/site-revise-post/reference.md`**. |
| Throat-clearing | Restates the heading in softer words; wise summary with no new detail |
| Hedge piles | "One possibility is...", "It is important to note..." |
| Watch/read homework | "Watch if you want the mechanism..." in standalone copy |

**Standalone worth test:** For each flagged sentence, ask: *If this stood alone, would it be worth saying?* Flag restatements, bridge-only lines, and decoration.

**Out of scope**

- Factual accuracy, citation quality, grounding sources (use `.cursor/perplexity/packs/deep-research.md` instead)
- SEO, title clickbait, Hugo formatting
- Spanish (`index.es.md`)

**Limitation:** You are also an LLM. Say when you are uncertain a line is AI vs author voice. Prefer **specific quoted evidence** over confidence.

---

## Repo excerpts (paste below)

Paste **plain text only** (strip Hugo front matter keys if you want; keep the prose). Include everything in scope.

### Excerpt A — List-facing copy (if in scope)

```text
[Paste title, description, grounding, sowhat]
```

### Excerpt B — Body (if in scope)

```text
[Paste full body or selected sections]
```

### Excerpt C — Author note (optional)

```text
[e.g. "Voice should feel blunt like Picard, not wellness blog" or "Keep the gray-hair tree-ring metaphor"]
```

---

## Output format (request from Perplexity)

Please structure your answer as:

1. **Cold-read verdict** (1 short paragraph: keep reading? where do you bounce?)
2. **Scores** (1–5 each, one line why): Clarity | Flow (read-aloud) | Human voice | Compelling stakes | Trust (does it sound earned?)
3. **Keep** (2–5 quoted lines that work; one phrase each on why)
4. **AI artifact table** with columns: Quote (exact) | Location | Pattern | Why it fails for humans | Fix direction (not full rewrite)
5. **Flow map** (bullets): paragraph or section → drag / momentum / note
6. **Top 5 fixes** (ranked; one actionable line each)
7. **Paragraphs to cut or merge** (quote opening words; CUT vs MERGE)
8. **Open questions** for the author (voice choices, not facts)

---

## After export (repo handling)

1. Export via `perplexity_export`.
2. Optional note: `docs/research/YYYY-MM-DD-perplexity-prose-<slug>.md` (verdict, top fixes, thread URL).
3. Apply edits with **revise-post** Step 2 (voice) and **revise-flow** (cadence). Do not paste Perplexity rewrites wholesale.
4. Re-run a local cold read or a second prose-review pass only if the draft changed heavily.
