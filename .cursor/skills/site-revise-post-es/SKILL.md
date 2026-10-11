---
name: site-revise-post-es
description: >-
  Systematically revises Spanish Hugo posts to ensure they flow naturally for a native Spanish-speaking audience. Use when the user asks to "revise the Spanish", "fix the translation", or wants step-by-step Spanish refinement. Focuses on adaptation to the audience, not fidelity to English source. Enforces hard native-Spanish constraints (conceptual anglicisms, active syntax, colloquial `###` hooks) per spanish-translation-content.
---

# Systematic Spanish Post Revision Workflow

## Purpose

Produce publish-ready Spanish content that works for its audience: clear, engaging, and culturally natural. This skill does NOT enforce translation fidelity; it validates that the Spanish text succeeds on its own terms.

Work through five sequential steps. Do not skip steps. The COT (Chain of Thought) in Step 0 is mandatory before evaluating.

## Philosophy: Adaptation Over Fidelity

| Translation-First (OLD) | Adaptation-First (THIS SKILL) |
|------------------------|-------------------------------|
| "Does this match the English?" | "Does this work for Spanish speakers?" |
| Flag "He aquí" as archaic | Check if the opener creates momentum |
| Enforce "nos" instead of impersonal | Either works; check consistency and flow |
| Ban "Lamentablemente" as filler | Judge if the emotional beat serves the piece |
| Word-for-word metaphor preservation | Spanish-equivalent imagery that carries the idea |
| Same `###` titles as English | Spanish hooks that **fit the paragraph** (MAY be colloquial; MUST NOT be empty labels) |
| Sentence-for-sentence mirror | Same thesis, facts, and ladder; Spanish prose and bridges |

**Golden rule:** The Spanish version should feel like it was written *for* Spanish speakers, not *translated from* English.

**Hard constraints (MUST):** Before Steps 1–5, internalize **`.cursor/skills/site-spanish-translation-content/SKILL.md`** → **Hard constraints (native Spanish, not translation-shaped)** and **Conceptual anglicisms (MUST check)**. Glance at English only for meaning gaps, not to copy syntax or section titles.

## Pass/Fail Definitions

Each step MUST be evaluated against these strict definitions:

- **Pass**: The Spanish text achieves the step's goal for its audience. Either no issues, or issues are stylistic preferences (author's choice).
- **Fix Needed**: Objective problems that hinder comprehension or flow. List them explicitly.

**A step may NOT be marked Pass while objective problems remain.**

## The 5-Step Spanish Revision Process

---

### Step 0: COT Pre-Reading (Internal Reasoning)

**Goal:** Understand what the piece is trying to achieve before judging how it does it.

For each section of the Spanish text (front matter and body), perform this mental process:

1. **Core intent:** What is the key idea this section delivers? Ignore English sentence structure entirely.
2. **Audience check:** Does the Spanish phrasing assume context that Spanish readers won't have?
3. **Flow audit:** Does each sentence lead naturally to the next? Are transitions logical in Spanish?
4. **Rhythm test:** Read aloud. Does it sound like written Spanish (not spoken, not translated)?
5. **Metaphor fit:** Do any metaphors feel forced or alien to Spanish ears? Would a different image work better?

**Only after this internal process, begin the evaluation.**

---

### Step 1: Hook and Flow Audit (Título + Lead)

**Goal:** Ensure the Spanish title and lead pull readers in with natural force. Evaluate the Spanish hook on its own terms, not against the English.

**Read first:** `.cursor/skills/site-revise-hooks/SKILL.md` for hook psychology. Adapt for Spanish cultural context, not word-for-word.

**Checks:**
- [ ] Title creates tension, paradox, or states a clear thesis in natural Spanish
- [ ] Title avoids corporate / management cadence in Spanish (e.g., "potenciar", "impulsar resultados")
- [ ] `description`/`sowhat` deliver concrete hooks in Spanish (evaluated as Spanish copy, not compared to English)
- [ ] The lead paragraph establishes momentum (reader wants to continue)
- [ ] For `type: video`: `description` sells the watch in Spanish; compelling as standalone hook
- [ ] For `type: claims`: `description` stands alone as the Claim in Spanish without needing the body

**Action:** Identify where the Spanish hook is weak (not where it differs from English). Suggest improvements that serve Spanish readers.

---

### Step 2: Natural Flow and Idiom Check

**Goal:** Ensure the Spanish reads as natural written content, not translated text. Focus on flow and comprehension, not literal accuracy.

**Also run** the **standalone worth** test from **`.cursor/skills/site-revise-post/SKILL.md`** → **Step 2** on every Spanish sentence and paragraph (cut decoration and spine restatements, not only calques).

**Hard constraints + calque tables (run first in Step 2):**

1. **`.cursor/skills/site-spanish-translation-content/SKILL.md`** → **Hard constraints (native Spanish, not translation-shaped)** and **Clause rebuild** (no EN relatives, *thus/making*, *You hear it when*, *in a vacuum*).
2. Same skill → **Conceptual anglicisms (MUST check)** (e.g. *viajar gratis* → *salir gratis*; *dispara límites* → *activa límites*; **sesgo a la verdad**; *defecto de diseño*).
3. Same skill → **Calcos frecuentes en contenido de ciencia e IA**.
4. Same skill → **Organizational, hierarchy, and workplace English (MUST check)** (e.g. *Sobre el papel*, *demostró*, *en estado puro*, *a conciencia*, *El precio que pagas son*, *a quién echarle la culpa*).

**Native opinion-editor audit (run inside Step 2):** Same skill → **Native opinion-editor audit (five directives)**. In one pass, check:

1. **Calques / false friends** (literal *en papel*, *pagas* abstract costs, dubbed metaphors).
2. **Anti-staccato** (true redundancy and robotic chains only; see **Anti-staccato guardrail** in **spanish-translation-content**; do not merge intentional parallel hooks).
3. **Collocations** (verb–noun pairs organic in Spanish).
4. **Tone** (do not soften punch; ban corporate/AI moralizing).

**Idiom-only mode (user asks for auditoría idiomática / native editor pass / revisa español / revise-spanish):** Use **`.cursor/skills/site-revise-spanish/SKILL.md`** instead of re-running this full workflow. For full **revise-post-es**, keep the step tables and confirmation gate below.

**What to flag (objective problems):**

| Problem | Spanish Example | Why It Fails |
|---------|-----------------|--------------|
| Incomprehensible without English | "El sistema fue re-ponderado" | Word doesn't exist; reader can't parse it |
| Broken Spanish syntax | "Disparan cuando entras a una habitación" | Missing reflexive; grammatically incomplete |
| Missing technical explanation | "Usan distancia KL" | Technical term unexplained; reader lost |
| Incomplete sentence openers | "En este, comparan..." | "En este" has no noun; reader stumbles |
| Staccato echo (objective) | *Diez personas… Diez personas debaten…* (same number twice) | Accidental English echo; re-point (*de ese tamaño*) without merging the whole paragraph |
| Staccato translation rhythm | Five generic S+V+O beats with **no** parallel job | Robotic; group or vary syntax |
| English skeleton / clause-map | *El X que Y carga Z*; *y así se hace posible*; *La oyes cuando*; *nace en el vacío* | Same relatives and glue as the English sentence; rebuild per **Clause rebuild** |
| Conceptual anglicism | "deja de viajar gratis", "dispara límites" | English idiom pasted into Spanish; use **Conceptual anglicisms** table |
| `###` hook mismatch | Rótulo que no describe el párrafo debajo | Hook must match the section beat (read paragraph, then fix or drop the heading) |
| Tone softening | "Es importante señalar que la jerarquía puede resultar problemática" | Corporate/AI smoothing; restore sharp manifesto voice |

**What NOT to flag (stylistic choices):**

| Choice | Examples | Why It's Valid |
|--------|----------|----------------|
| Formal register | "He aquí", "Conviene notar" | Author may want elevated tone |
| Emotional openers | "Lamentablemente", "Por desgracia" | Tone-setting is authorial choice |
| "Nosotros" vs impersonal | Either can work | Depends on desired intimacy/distance |
| Spanish-appropriate metaphors | Different from English | "Te explota en la cara" vs "it backfires" |
| Intentional short parallels | *En familia… En clase… En la empresa…*; *Abajo… arriba…* | Manifesto rhythm; merging into one sentence usually weakens the hook |
| Brief closers with meaning | *para poder entenderlo* | Not filler if it completes the thought; cut only when empty |

**Checks:**
- [ ] No invented words or non-existent technical terms
- [ ] No grammatically broken sentences (missing subjects, wrong verb forms)
- [ ] Technical terms explained on first use
- [ ] Transitions between paragraphs are logical in Spanish
- [ ] Consistent address form (tú/usted/nosotros) throughout
- [ ] Negation-first constructions are purposeful, not default

**Action:** Flag only objective comprehension barriers. Leave stylistic choices to author judgment.

---

### Step 3: Accessibility and Clarity Pass

**Goal:** Ensure the post makes sense to a reader who encounters it cold (no English source, no video, no prior context).

**Technical Term Rule:**
- Do NOT remove or rename technical terms
- Keep the exact term and add a gloss in parentheses on first use
- **Fail this step if any technical term is missing or would confuse a cold reader**

**Read if applicable:** `.cursor/skills/site-ai-for-general-audience/SKILL.md` when the post covers AI, neuroscience, or cognitive science.

**Checks:**
- [ ] Every technical term has a plain-language gloss in parentheses on first use in Spanish
- [ ] No assumed context from English source (e.g., "as we saw in the video" when reader hasn't watched)
- [ ] Concrete visual language: actions and scenes, not abstract noun stacks
- [ ] Short paragraphs (1-3 lines) for scanability
- [ ] Primary-source quotes in English may stay in English; surrounding commentary must be natural Spanish

**Action:** Add glosses where missing. Rewrite abstract stacks as visual actions. Chunk dense paragraphs.

---

### Step 4: Formatting Sweep

**Goal:** Enforce strict mechanical rules.

**Checks:**
- [ ] ZERO em dashes (`—`). Replace with commas, semicolons, colons, or parentheses.
- [ ] Bold (`**`) is restrained: ~2-5 spans per short block (not wall-to-wall).
- [ ] `translationKey` matches the English file exactly.
- [ ] `categories` and `tags` are identical to the English post (do not translate taxonomy terms).
- [ ] `date`, `youtube_id`, `images`, and structural front matter are mirrored from the English file ( **`date`** format: **`.cursor/rules/site-content-markdown-writing.mdc`** → **Publish `date`**).
- [ ] No decorative emoji in `description` or `grounding`.

**Action:** Search for `—` and replace. Audit `**` density. Verify front matter field parity with the English file.

---

### Step 5: URL and Social File Check

**Goal:** Ensure Spanish URLs are correct and social files use the right permalinks.

**Checks:**
- [ ] Run `hugo list all` and confirm the Spanish permalink for this post.
- [ ] `facebook-es.txt` uses the correct ES and EN permalinks from `hugo list all` (not constructed from the English folder name).
- [ ] If `linkedin.es.txt` exists, same check.
- [ ] Internal links in the Spanish body use the Spanish canonical path where a Spanish page exists.
- [ ] For `content/social-protocols/`: `aliases` on `index.es.md` mirror the pattern from sibling bundles.

**Reference:** `.cursor/skills/site-spanish-translation-content/SKILL.md` → **Spanish URLs, permalinks, and aliases**.

**Action:** Correct any URL that does not match `hugo list all`. Do not guess paths.

---

## Execution

1. **Read the Spanish post** fully before starting. Glance at English only if meaning is unclear.
2. **Run Step 0 (COT) internally** for each section before evaluating.
3. **Work through Steps 1-5 in order.**
4. **For each step, first output a findings list:**
   - For objective problems: quote exact phrases and explain why they fail
   - For stylistic observations: note them as "Author choice:" not as violations
   - If no issues found, state "No issues found"
5. **Then produce a numbered table with Before/After columns** showing changes that *would* address objective problems. **DO NOT write changes to the file yet.**
6. **Stop and wait for user confirmation.** Present the full analysis with all findings and Before/After tables, then ask: "Apply all changes? Or specify which steps (e.g., 'apply Steps 1-3' or 'apply only Step 2')?"
7. **Only after explicit confirmation**, apply the selected changes to the file.
8. **Final check:** Run `hugo build` to ensure no syntax errors and confirm the Spanish URL renders.

## Required Output Format (Analysis Phase)

### For each step, produce:

1. **Step header** with number and name
2. **Findings list:**
   - **Objective problems** (must fix): quote exact phrases and explain failure
   - **Stylistic observations** (author's choice): note but don't mandate changes
   - Or: "No issues found"
3. **Numbered table** with four columns: `#`, `Ubicación`, `Antes`, `Después`
4. **Status line** indicating Pass or Fix Needed

### Example:

```
### Step 1: Hook and Flow Audit

**Objective problems:**
- "Cognición Incorporada" (título genérico, no genera tensión ni curiosidad)

**Stylistic observations:**
- Author uses formal "He aquí" opener (valid choice for elevated tone)

| # | Ubicación | Antes | Después |
|---|-----------|-------|---------|
| 1 | title | Cognición Incorporada | 🖐️🧠 Señalas primero, luego piensas |

**Status:** Fix Needed (Title needs stronger hook for Spanish readers)
```

### At the end of all steps:

1. **Summary header:** "## Resumen de Revisión"
2. **Status overview:** List each step with Pass or Fix Needed status
3. **User prompt:** "¿Aplicar todos los cambios? O especifica qué pasos (por ejemplo, 'aplicar pasos 1-3', 'aplicar solo paso 2', 'revertir paso 4')."
4. **Proposed revised post:** Show the full post as it *would* appear (for reference only; do not write yet)

## Selective Application (After User Confirmation)

Only after the user explicitly responds:

- "aplicar todo" or "sí" / "yes" → Write all changes from all steps to the file
- "aplicar paso 2 solo" / "apply Step 2 only" → Write only Step 2's changes
- "aplicar pasos 3-4" / "apply Steps 3-4" → Write only those steps' changes
- "revertir paso 1" / "revert Step 1" → Skip Step 1's changes; apply all others
- "no" or "cancelar" / "cancel" → Do not write any changes

## Optional: idiom audit output (chat)

When the user requests a **native editor** or **idiom audit** without the full step-by-step revision:

```
[Texto final]
(full corrected Spanish)

[Changelog]
- (3–4 bullets: calques / redundancy only; why the original sounded translated)
```

Then wait for **apply** before editing files unless the user said implement directly.

## Related skills

- **Rules reference:** `.cursor/skills/site-spanish-translation-content/SKILL.md` (**Hard constraints**, **Conceptual anglicisms**, **Native opinion-editor audit**, calque tables)
- **Validated bundle example:** `content/human-condition/2026-06-11-deception-truth-bias-incentive-gaps/index.es.md` (colloquial `###`, idiomatic fixes, not EN title mirror)
- **English revision workflow:** `.cursor/skills/site-revise-post/SKILL.md`
- **AI/neuroscience plain language:** `.cursor/skills/site-ai-for-general-audience/SKILL.md`
- **Facebook post (ES/EN):** `.cursor/skills/site-facebook-post/SKILL.md`
- **LinkedIn post (bilingual links):** `.cursor/skills/site-linkedin-post/SKILL.md`
