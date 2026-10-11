---
name: site-reptilocracy-gemma
description: >-
  Drafts Reptilocracy episode description, tldr, and fluff via local Gemma 4
  with thinking mode, satire-first briefs, and a psychological-fitness knife.
  Use when creating or rewriting Reptilocracy copy, Gemma Reptilocracy text,
  good-day/bad-day emergency-powers satire, or when the user asks Gemma to
  write Reptilocracy fields (not LinkedIn; not full bundle layout).
---

# Reptilocracy episode copy (Gemma 4 + thinking)

## Goal

Get **local Gemma 4** to draft **`description`**, **`tldr`**, and **`fluff`** for posts under **`content/cognitive-memetics/reptilocracy/`**, using the brief that works in practice: **satire first**, then the knife that **psychological fitness** is almost never on the leadership agenda.

This skill owns **Gemma prompting and pick flow** only. Bundle layout, tags, dates, images, LinkedIn, and Spanish naturalness audits stay with sibling skills (see **Related**).

## When to use

| Use this skill | Use another skill |
|----------------|-------------------|
| New Reptilocracy episode copy via Gemma | **`site-cognitive-memetics-content`** for folder, front matter shell, categories, tags, images |
| Rewrite / punch up Reptilocracy `description` / `tldr` / `fluff` with Gemma | **`site-linkedin-post`** for `linkedin.txt` |
| User says Gemma, thinking mode, satire, psych fitness for Reptilocracy | **`site-revise-spanish`** after ES copy exists |
| | Generic **Gemma teaser** in cognitive-memetics for Cube-Cows / Raymond / Pawtropolis / T-Shirt Art (not this hub) |

## Defaults

| Setting | Value |
|---------|--------|
| Gateway | `LOCAL_LLM_BASE_URL` or `http://127.0.0.1:1320/v1` |
| Model | `LOCAL_LLM_MODEL` or `gemini_free/gemma-4-26b-a4b-it` |
| Thinking | **MUST** set `chat_template_kwargs.enable_thinking: true` |
| Pack | **`packs/episode_copy.md`** |
| Script | **`scripts/draft_episode_copy.py`** |

## What “done” means for Gemma

**CONSTRAINT:** For new or rewritten Reptilocracy episode fields, MUST obtain **`description`**, **`tldr`**, and **`fluff`** from Gemma under this skill’s brief. MUST NOT ship agent-invented stamp poetry or policy-memo prose and claim Gemma wrote it.

- Enforcement: candidates file or chat shows Gemma output; pick recorded before apply
- Violation: STOP, re-run script or Gemma call; do not apply agent-only fields as final

CORRECT:
```text
Run draft_episode_copy.py → surface 3 descriptions + tldr + fluff → user picks → apply
```

PROHIBITED:
```text
Agent writes clinical fluff about the stamp, pastes a Gemma teaser only for description
```

## Brief rules (what to send Gemma)

**CONSTRAINT:** System and user briefs MUST be **satire first** (dark café humor, Latin American tragicomic shrug: sad and funny). MUST land the knife that institutions gossip about **mood / weather / good day vs bad day** and almost never put **psychological fitness** on the table for leaders (honesty, accountability, courage, reality-contact, impulse control, empathy under power). Emergency / continuity theater covers both moods because fitness was never a gate.

- Enforcement: scan the prompt for satire + psych-fitness knife before calling Gemma
- Violation: STOP, rewrite the brief; do not call with a white-paper thesis stack alone

**CONSTRAINT:** MUST NOT brief Gemma as a policy memo, LinkedIn think-piece, MBA fog deck, or activist petition sermon. Petition CTA belongs in LinkedIn / footer i18n, not in episode `description` / `tldr` / `fluff`.

CORRECT:
```text
Joke spine: same leader, smile and shout; clerk stamps emergency powers on both.
People track temperament like humidity. Fitness audit never enters the briefing.
```

PROHIBITED:
```text
Write clinical institutional analysis of continuity vs early structural change
with board-deck vocabulary and no satire beat.
```

**CONSTRAINT:** When the art shows **one** character on a good day and a bad day, MUST tell Gemma it is the **same** person twice. MUST NOT allow two-lizard readings.

**CONSTRAINT:** `DESCRIPTION_1` / `DESCRIPTION_2` / `DESCRIPTION_3` MUST be three **complete alternative card teasers** (each 1–2 sentences, standalone). MUST NOT split smile / shout / stamp across three panel captions.

## Procedure

1. Resolve episode **`title`** (working or final) and a **mechanism brief** (gag + institutional knife; not a full prop inventory). Prefer art caption + one sentence of what the stamp / system move means.
2. From repo root, run (do not hand-roll the API unless the script fails):

```bash
python3 .cursor/skills/site-reptilocracy-gemma/scripts/draft_episode_copy.py \
  --title "Good Days, Bad Days" \
  --brief "Same lizard leader: GOOD DAY smile and BAD DAY shout in one dossier. Clerk stamps EMERGENCY POWERS APPROVED across both. Mood gossip replaces psychological fitness; stamp covers both."
```

Optional: `--out content/cognitive-memetics/reptilocracy/<slug>/gemma-candidates.txt`

3. Surface the **three numbered descriptions** plus proposed **`tldr`** and **`fluff`** from stdout (and the candidates file if written).
4. Wait for a numbered **description** pick (and confirmation that TLDR/FLUFF are accepted, or a rewrite pass).
5. Apply the pick to English **`index.md`** with light mechanical cleanup only (em dash, YAML safety, restrained `**bold**` per **`.cursor/skills/site-revise-emphasis/SKILL.md`**). MUST NOT flatten the joke into a caption inventory.
6. If **`index.es.md`** exists or is being created: invoke Gemma again (thinking on) to adapt the chosen EN **`description`** into three Spanish candidates; pick or apply the best mirror of the EN joke spine; adapt **`tldr`** / **`fluff`** into native Spanish (not calque). Then run **`.cursor/skills/site-revise-spanish/SKILL.md`** when the user asks.
7. Bundle shell (date, `heading_code`, tags, image, calendar): **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`**. LinkedIn: **`.cursor/skills/site-linkedin-post/SKILL.md`**.

### Thinking-mode output hygiene

**CONSTRAINT:** Final Gemma **content** MUST contain only the labeled fields (`DESCRIPTION_1`…`FLUFF`). If thinking spills into `content` and truncates fields, MUST re-run with the pack’s final-answer-only instruction (script already sets this). MUST NOT paste raw reasoning into Hugo front matter.

- Enforcement: parse labels before apply
- Violation: re-run script; discard truncated content

### Gateway down

**CONSTRAINT:** If the gateway is unreachable, MUST say so. MUST NOT invent teasers and claim Gemma wrote them. MAY leave a clearly labeled stub only if the user accepts shipping later.

## Gold tone (energy, not to copy)

- `content/cognitive-memetics/reptilocracy/2026-08-02-photo-op-readiness/` (dry satire punch)
- `content/cognitive-memetics/reptilocracy/2026-09-13-the-alarm-worked/` (short scene knife)
- `content/cognitive-memetics/reptilocracy/2026-09-20-good-days-bad-days/` (mood weather vs psychological fitness)

Sister psych-fitness close (tone only; do not paste petition into episode fields): episodes that end without soft MBA fog, e.g. `2026-08-16-insert-vote-receive-same` fluff beat on fitness standards.

## Related

- Bundle / hub rules: **`.cursor/skills/site-cognitive-memetics-content/SKILL.md`**
- LinkedIn: **`.cursor/skills/site-linkedin-post/SKILL.md`**
- Spanish audit: **`.cursor/skills/site-revise-spanish/SKILL.md`**
- Shared gateway helpers: **`.cursor/skills/site_local_eval_common/common.py`**
