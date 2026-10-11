# Revise-post reference (Step 2 detail)

Full banned-pattern tables and examples for **Step 2** in **`SKILL.md`**. Agents MUST apply these when running voice scrub.

## Banned patterns (delete or rewrite)

| Pattern | Fail example | Fix direction |
|---------|--------------|---------------|
| Teacher framing | "You will learn...", "What you need to understand is...", "The useful move is smaller" | State the point directly |
| Editor/meta openers | "X's hook is...", "This long interview bridges..." | Name the claim; drop editor-speak |
| Bridge verbs | "bridges felt fatigue and cellular allocation" | Say what connects: subject + mechanism |
| Semicolon aphorisms | "Stress response drains; purpose concentrates." | One sentence with a subject and verb, or two full sentences |
| Staccato thesis stacks | Same-length punch lines every sentence; "That is not abstract. That is..." | Vary length; merge or add a setup sentence |
| Meta-commentary | "It frames the boundary as...", "This makes the point clear..." | Direct impact on the reader |
| Watch/read CTAs in standalone copy | "Watch if you want the mechanism behind..." | State payoff; no homework framing |
| Opaque coined labels | "stuck on vigilance", "fixed energy budget" alone | Gloss in the **same sentence** (see Step 3) |
| Filler hedges | "One live possibility is that...", "It is important to note that..." | Cut or state directly |
| Vague intensifiers | "very", "really", "truly", "fundamentally" | Cut or replace with a concrete detail |
| Adverb stacks | "precisely exactly", "quite fundamentally" | One modifier or none |
| Classic AI rhetoric | "This changes how we understand...", "This framework maps...", "Classic experiments reveal..." | Name who found what |
| Clever filler metaphors | "not a label on your calendar", "the meter is running", "bodies buffer and compensates", "sells the blend as 'me'", "steer the rewrite" | State the qualification in plain words; one metaphor thread per paragraph |
| Caption pivots | "Picard is careful here.", "The pattern matters.", "Still," opening a punch line | Fold the caveat into the sentence before it; use concrete symptoms or numbers |
| Telegraphic compression | "Cortisol rises. Shoulders tighten. You replay..." (five same-length beats) | One spoken sentence with commas or a subordinate clause |
| Parallel contrast stacks | "Some people run a cold two-track story for the audience. Others sell themselves the good reason until it feels true." | Say the two moves in plain verbs (lie to the audience / buy your own story); avoid matched Some/Others thesis pairs |
| Abstract metonymy | "The room learns to freeze." / "You can lie to the room." / "read the room" | Name the actual people, bystanders, or physical setting. Literal physical rooms are allowed. |
| Construct-labeling meta | "**Self-deception** names that move:", "X names the pattern:" (label is the sentence's only job) | State the behavior in spoken English first; MAY gloss a named construct in the same sentence (`**cold-hot empathy gap**: when you are calm…`) |
| Academic accessibility jargon | "partly available to reflection", "cognitively inaccessible", "phenomenally present but report-poor" | Plain access words: hard to see, hard to say, foggy even to you |
| Mechanism jargon without a scene | "cold two-track story", "dual-process cover", "strategic opacity" with no picture | Who does what in a concrete scene (kitchen, meeting): quit a job, pick a side, nod along |
| Rhetorical noun fragments | "The mirror." / "The wound." / "The truth." as standalone beats | Fold into a full sentence with subject + verb, or cut |
| Revelation stacks | Several one-line poetic restatements of the same idea | One claim, then mechanism, scope, or evidence |
| Poetic thesis remix | Same claim restated in prettier metaphors without new information | Keep one clear statement; cut the echo |
| Metaphor-shell restack | After a live-mix claim: "The blend is your experienced reality…" then "If the self were a sealed file, you would have no steering wheel." | Next sentences MUST add how (inputs, evidence, example, scope). Cut undefined nouns (`the blend`) and extra metaphor pairs (`sealed file` / `steering wheel`). |
| Industry-verb shells | "then sells the blend as 'me.'"; "let you steer the rewrite."; "identity feels solid because the brain keeps assembling it."; "can change the next assembly."; "what you remember from the same room" (abstract group, no people); "that template can harden"; "You can name the box parents or peers put you in, and the box you put yourself in." | Drop sales, factory, driving, and hardening verbs when they are the clause's only job. Use spoken verbs and a real scene (kitchen, named colleagues). Parallel "name the box… and the box" fails; say the story they taught you and the story you now tell yourself. |

## Explanatory prose (claims, video, Substack)

**Scope:** `type: claims` Thoughts body; `type: video` body, `description`, and `sowhat`; `substack.md` / `substack.es.md`. Drafting constraint also in **`.cursor/rules/site-content-markdown-writing.mdc`** → **Explanatory prose**.

**Goal:** Write like a clear human explaining a difficult idea. Polished fragment stacks that sound deep but do not name a mechanism or a falsifiable claim are fluff.

**MUST:**

- State the claim directly before interpreting it.
- Default to paragraphs of **2–4 sentences** (not one-line fragments).
- Prefer concrete verbs (`people deny`, `the mind protects`).
- Keep qualifications the evidence needs (`can`, `may`, `in some cases`).
- Explain the mechanism or evidence, then let it carry the point.
- Stay psychologically careful and morally restrained: do not accuse the reader or promise secret insight.

**MUST NOT:**

- Ship rhetorical noun fragments or revelation stacks (table above).
- Restack the same claim in altered poetic language.
- Ship **Metaphor-shell restack** (table above): complete sentences that only relabel or metaphorize the previous claim.
- Ship **Industry-verb shells** (table above): sales, factory, driving, or hardening verbs (`sells`, `steer`, `assembling`, `assembly`, `harden`) or a parallel "name the box… and the box" restack when that verb or restack is the clause's only job.
- Use metaphors as the default texture; use them rarely, and only when they clarify.

**Sequence test (Metaphor-shell restack):** After sentence N, sentence N+1 **fails** if it only (a) relabels N with an undefined noun (`the blend`), (b) restates N as a tautology (`experienced reality`), or (c) restates the implication of N with a new metaphor pair (`sealed file` / `steering wheel`). Consecutive punches that do not add mechanism, evidence, example, or scope fail **even if each line is a complete sentence**.

**Clause test (Industry-verb shells):** A clause **fails even alone** (not only as N+1) if its payload is a clever industry verb for a mental process, or a fake scene (vague group container with no kitchen, people, or objects; includes abstract "the room" for an audience). "Identity feels solid because the brain keeps assembling it" fails this test: assembling is the factory verb, not the mechanism (expectations + current input).

**`cognitive-memetics`:** MUST follow **Hybrid prose** in **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`** (emotion + claim; punch OK when the next clause names mechanism/scene; ban empty closers like "The long view matters."). Not full essay Explanatory prose.

**Do/Don't (explanatory):**

| Do | Don't |
|----|-------|
| "People sometimes react strongly to traits in others that they find hard to acknowledge in themselves. Psychologists call one version **defensive projection**: a person may deny an unwanted trait, then become more likely to perceive it in someone else." | "We call it seeing. But it is a mirror. The other becomes the container. Judgment becomes a shield." |
| "This does not mean every criticism is projection. Strong moral judgment can sometimes serve self-protection as well as truth-seeking." | "The wound. The shield. The truth beneath the verdict." |
| Open with the central question in plain words, then define the construct and its limits. | Open with a stack of oracular one-liners that never name who does what. |
| "The brain mixes what past experience taught you to expect with what the senses pick up now. Change what you attend to, and you take a different mix as who you are." | "Identity feels solid because the brain keeps assembling it. The blend is your experienced reality, including the sense of who you are. If the self were a sealed file, you would have no steering wheel." |
| "You take that mix as who you are." / "Attention and small wins change what gets rebuilt." / "Mood decides what you notice, even when the kitchen has not changed." / "As an adult you still use that same self-criticism." / "Say the story they taught you, and the story you now tell yourself." | "then sells the blend as 'me.'" / "let you steer the rewrite." / "can change the next assembly." / "what you remember from the same room" / "that template can harden" / "You can name the box parents or peers put you in, and the box you put yourself in." |

**Additional checks (this scope):**

- [ ] Opening states the claim before interpretation
- [ ] Default paragraph shape is 2–4 sentences, not fragment stacks
- [ ] Zero rhetorical noun fragments / revelation stacks / poetic thesis remixes / **Metaphor-shell restack** / **Industry-verb shells**
- [ ] Sequence test: each sentence after the first adds mechanism, evidence, example, or scope
- [ ] Clause test: no sales, factory, driving, or hardening verb as the clause's only job; no fake scene (vague group container; abstract "the room")
- [ ] Metaphors are rare clarifiers, not the paragraph's main job
- [ ] Evidence-scope qualifiers appear where overclaim would mislead

## Clever filler and mixed metaphors (MUST NOT)

After compression or a "polish" pass, agents often add **smooth-sounding filler** that fails cold-read. **MUST NOT** ship these even when they feel literary.

| Fail | Why | Fix |
|------|-----|-----|
| "not a label on your calendar" | Opaque cleverness; reader must decode | "That figure is from dish biology, not your whole daily burn." |
| "A whole body buffers and compensates." | Vague abstraction; no picture | "Your whole body does not jump sixty percent on a stressed day, but heart rate, tight muscles, and rumination still drain the **budget**." |
| "the meter is running" (when **budget** is the spine metaphor) | Mixed metaphor | Stay on **budget** / **spend** / **pay** language |
| "Picard is careful here." | Editor caption; no new fact | Drop, or merge: "Picard applies that figure to cultured cells, not…" |
| "The pattern matters." / "The awe holds either way." | Thesis-stack throat-clear | Cut or replace with the pattern in plain words |
| Industry-verb shells (`sells`, `steer`, `assembling`, `harden`, "name the box… and the box") | Sales/factory/driving payload; see table above | Spoken verb + real scene. Full row: **Industry-verb shells**. |

**Rule:** One primary metaphor per paragraph (budget **or** meter **or** ledger, not all three). Qualifications use **concrete nouns** (dish cells, whole daily burn, 2 a.m. rumination), not witty substitutes.

**Do/Don't examples:**

| Do | Don't |
|----|-------|
| "Picard argues your body runs on a **fixed energy budget**: a limited daily pool of cellular spend." | "Picard's hook is a fixed energy budget." |
| "Exhaustion and lost meaning can show up together when the **stress response** keeps spending after the moment passes." | "…when the budget stays stuck on vigilance." |
| "Picard applies that figure to dish biology, not your whole daily burn, yet faster heart rate, tight muscles, and rumination at 2 a.m. still drain the **budget**." | "not a label on your calendar. A whole body buffers and compensates. Still, …the meter is running." |
| "Place cells track location." (direct) | "One live possibility is that place cells track location." (hedged) |
| "This interview ties fatigue to where cellular **budget** goes." | "This interview bridges fatigue and allocation: not more fuel, better aim." |
| "You can lie to the audience on purpose. You can also buy your own story." | "Some people run a cold two-track story for the audience. Others sell themselves the good reason until it feels true. **Self-deception** names that move…" |
