---
name: site-revise-score
description: >-
  Scores a Hugo post against the site's editorial standards. Use when the user asks to
  "score this post", "rate this content", "check against the rules", or revise-score.
  It evaluates hooks, voice, formatting, structural integrity (including standalone-worth /
  zero-decoration), accessibility, and type compliance.
---

# Revise score (editorial standards)

## Purpose

Evaluate a post (`content/**/*.md`) against the site's strict editorial rules. Output a score and actionable feedback with exact rewrites.

## Voice Lock (Global Constraint)

This score measures editorial mechanics **without altering the author's voice**. Technical terms, metaphors, and domain-specific phrasing should be flagged for clarification (glosses), NOT replaced with bland paraphrases.

**Do/Don't Examples:**
- Do (flag for): "hippocampal system (needs gloss)"
- Don't (score deduction): Replacing with "brain's navigation centers" (term lost)

## Scoring Criteria (Max 100 points)

### 1. Hook Strength (20 points)

- **Title (8):** Does it create tension, paradox, or state a clear thesis? Deduct for bland labels, corporate cadence (e.g., "unlock", "elevate"), fake urgency, or mystery boxes that withhold what the post is about.
- **Lead / list copy (8):** Does it fulfill the title's hook with concrete, non-obvious promises? For `type: video`, does it sell the play button? For `type: claims`, does the `description` stand alone as a clear narrative assertion without the body? Deduct if list-facing copy duplicates the body rather than acting as hooks. Deduct if **`description`**, Claim, **`grounding`**, or **`sowhat`** fail the **cold-read** test (stranger cannot decode the line without the body or video). See **`.cursor/skills/site-revise-hooks/SKILL.md`** → **Capture vs pass** and **Cold-read**; for **`type: video`**, also **`.cursor/skills/site-video-content/SKILL.md`** → **Cold-read gate**.
- **Body headings (4):** When the post has markdown body sections, do **`##`** / **`###`** headings **pull** (tension, question, punch), not only label the topic? Deduct **2 points** per seminar-style label that could be a hook (`## Ideas that travel`, `## The Takeaway`, `## Tools shape the metaphor`). Neutral nav labels (for example **`## Chapter Guide`**) do not deduct. See **`.cursor/skills/site-revise-hooks/SKILL.md`** → **Body headings (`##` and `###`) — hooks, not labels**.

### 2. Voice and Tone (25 points)

- **Active Voice (10):** Are subjects performing actions? Deduct for passive voice pile-ups ("it was determined that...") where the subject is known.
- **No Fluff (10):** Is it free of AI rhetoric, filler, vague intensifiers, and corporate jargon?

  **Banned Patterns (deduct for any occurrence):**
  | Category | Examples |
  |----------|----------|
  | AI rhetoric | "This changes how we understand...", "This framework maps...", "Classic experiments reveal..." |
  | Filler hedges | "One live possibility is that...", "Some studies suggest...", "It is important to note that..." |
  | Vague intensifiers | "very", "really", "truly", "fundamentally", "categorically" |
  | Adverb stacks | "precisely exactly", "quite fundamentally", "very significantly" |
  | Corporate jargon | "unlock", "elevate", "drive outcomes", "synergize" |
  | Metaphor-shell restack | After a mechanism sentence: "The blend is your experienced reality…" then "If the self were a sealed file, you would have no steering wheel." Deduct **2 points** per shell sentence (cap with other No Fluff hits inside this 10). Sequence test: **`.cursor/skills/site-revise-post/reference.md`**. |
  | Industry-verb shells | "sells the blend as 'me'"; "steer the rewrite"; "keeps assembling it"; "next assembly"; "the same room"; "template can harden"; "name the box… and the box." Deduct **2 points** per shell clause. Clause test: **`.cursor/skills/site-revise-post/reference.md`**. |

  **Do/Don't Examples:**
  - Do: "Place cells track location."
  - Don't: "One live possibility is that place cells track location." (-2 points for hedge)
  - Do: "The model predicts behavior."
  - Don't: "This framework fundamentally maps the predictive architecture." (-2 points for intensifier + jargon)
  - Do: "The brain mixes past expectations with what the senses pick up now."
  - Don't: "The blend is your experienced reality, including the sense of who you are." (-2 points for Metaphor-shell restack)
  - Don't: "Identity feels solid because the brain keeps assembling it." / "then sells the blend as 'me.'" (-2 points for Industry-verb shells)
- **No Condescension (5):** Does it state ideas directly without teacher-like framing? Deduct 1 point per instance.

  **Teacher Framing Patterns:**
  | Pattern | Example |
  |---------|---------|
  | Future tense setup | "You will learn...", "You will see...", "You will discover..." |
  | Importance signaling | "It is important to note that...", "What is crucial here is..." |
  | Comprehension directive | "What you need to understand is...", "The key thing to grasp is..." |

  **Do/Don't Examples:**
  - Do: "See where this model fits..."
  - Don't: "You will see where this model fits..."
  - Do: "Place cells fire when you enter a room."
  - Don't: "What you need to understand is that place cells fire when you enter a room."

### 3. Accessibility and Concrete Language (20 points)

- **Cold Reader Test (10):** Does the post make sense to someone who has not watched the video or read the source? Deduct for technical terms without plain-language glosses on first use. List-facing cold-read is scored under **Hook Strength → Lead / list copy**; do not double-deduct the same line in both blocks.

  **Technical Term Rule:** The term must be PRESERVED with a gloss in parentheses, not paraphrased away.
  - Do (correct): "hippocampal system (the brain's navigation and memory hub)"
  - Don't (deduct -2 points): "brain's navigation centers" (term dropped)
  - Do (correct): "place cells (neurons that track where you are in space)"
  - Don't (deduct -2 points): "location-tracking neurons" (term dropped)

  **Common terms requiring glosses:** "hippocampal system", "occipital cortex", "place cells", "grid cells", "vicarious trial and error"
- **Concrete Verbs (10):** Does it use strong, visual, specific verbs instead of abstract noun stacks? Deduct for constructions like "provides the primary structure of thought" where a simpler verb would do.

### 4. Formatting and Style (20 points)

- **Punctuation (10):** Are there ZERO em dashes (U+2014: `—`)? Must use commas, semicolons, colons, or parentheses instead. One violation = full deduction.
- **Emphasis (10):** Is bolding (`**`) restrained (roughly 2-5 spans per short block)? Deduct for wall-to-wall bold that makes every noun and verb look equally important.

### 5. Structural Integrity (15 points)

- **Standalone Worth (5):** Does every sentence and paragraph earn its place alone (new fact, example, mechanism, stakes, or turn)? A new metaphor or coined label for the same claim is not a turn. Deduct **1 point** per sentence or **2 points** per paragraph that fails (restates spine/H2, throat-clearing, wise summary with no concrete detail, bridge-only, **Metaphor-shell restack**, or deleting it changes almost nothing). Cap this block at 5 points deducted. Quote the **full** failed sentence or paragraph in violations. Do not double-deduct the same sentence in both **No Fluff** and this block; prefer **No Fluff** for the shell, this block for other decoration.

- **No Negation-First (8):** Are there ZERO "negation-first" constructions? Deduct 2 points for each instance (cap this block at 8 points deducted).

  **Definition:** Phrases that define by negation: "It is not X, it is Y" or "Not A but B"

  **Do/Don't Examples:**
  - Do (correct): "This is a behavioral map."
  - Don't (-2 points): "It is not an algorithm; it is a behavioral map."
  - Do (correct): "The brain simulates the future."
  - Don't (-2 points): "It is not random activity; the brain simulates the future."
  - Do (correct): "Gestures build thoughts."
  - Don't (-2 points): "Gestures were not expressing a thought, they were building it."
- **Type Compliance (5):** Check the post's `type` and verify the specific required fields are present and correct:
  - `type: video`: Has `youtube_id`? Has `sowhat`? Does `description` function as hooks above the embed (not a body summary)?
  - `type: claims`: Has `description` (Claim)? Has `grounding`? Does the Claim read clearly without the body?
  - `type: panel` / `type: sayings`: Has `description` (teaser)? Follows cognitive-memetics field rules?

## How to execute the score

1. **Read the post** fully before scoring.
2. **Cold-read** list-facing fields and each body **`##`** / **`###`** heading alone before scoring hooks (per **`.cursor/skills/site-revise-hooks/SKILL.md`**).
3. **Evaluate each criterion block** against the specific banned patterns listed. Note the exact line or phrase that triggered a deduction.
4. **For banned patterns:** When you find a filler hedge, teacher framing pattern, adverb stack, negation-first construction, seminar-label heading, standalone-worth failure, **Metaphor-shell restack**, or **Industry-verb shells**, quote the exact phrase (full sentence or paragraph for decoration) and cite which banned pattern it matches.
5. **Cut candidates:** In **Violations**, group standalone-worth failures under a **Cut candidates** subheading (CUT vs MERGE). Do not propose rewrites for lines marked CUT unless the user asks.
6. **Report** in this exact structure:

```
## Score: XX/100

### 1. Hook Strength: XX/20
### 2. Voice and Tone: XX/25
### 3. Accessibility and Concrete Language: XX/20
### 4. Formatting and Style: XX/20
### 5. Structural Integrity: XX/15

---

### Violations

For each deduction, provide:
- [Criterion] "exact quoted phrase from post" → reason for deduction

**Quote exact phrases.** Do not paraphrase the violation. The quote must match the post text character-for-character so the author can find it.

---

### Fixes

For each violation, provide the exact rewritten text.
```

7. If the post scores 90+, say so clearly and skip rewrites for minor issues.
8. If the post scores below 70, provide a priority order for fixes (highest impact first; cut decoration before polish).

## Important caveat

This score measures **editorial mechanics only**: voice, hooks, formatting, and structure. It does not evaluate the quality of the ideas, the accuracy of claims, or whether the argument holds. A post can score 100/100 and still have a weak thesis. Use this alongside your own editorial judgment.

## Next step

For **clarity, compression, flow, and voice** on the prose itself, run **`revise-flow`** (`.cursor/skills/site-revise-flow/SKILL.md`) before re-scoring.

If the score is below 90, run the **`revise-post`** skill (`.cursor/skills/site-revise-post/SKILL.md`) for a systematic step-by-step fix.

To run **flow → hooks → post** (and optionally score) in one session, use **`.cursor/skills/site-revise-post/SKILL.md`** (ask for score as an optional phase).
