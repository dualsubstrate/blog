---
name: site-essay-voices
description: >-
  Routes Hugo essay sections to host voice files under data/voices and tells
  operators how to apply a voice at polish time. Use when choosing a narrator
  for human-condition, x-minds, social-protocols, or mind-infrastructure, when
  editing a voice file, when running make audit-voice, or when setting essay
  polish focus from a voice profile.
---

# Site essay voices

Host catalog: `data/voices/`. Portable checker: strop `pkg/evaluation/voice`. Pipeline pin: content-pipelines essay composition and polish.

Composition stays on the shared commute voice in `data/essay-content-forms.yaml`. A section voice is a polish-stage contract.

## Section routing

| Section | Voice id | File |
| --- | --- | --- |
| `human-condition` | `patient-narrator` | `data/voices/patient-narrator.md` |
| `x-minds` | `patient-narrator` | `data/voices/patient-narrator.md` |
| `social-protocols` | `clint-eastwood` | `data/voices/clint-eastwood.md` |
| `mind-infrastructure` | `cognitive-pragmatist` | `data/voices/cognitive-pragmatist.md` |

Technical pieces that must stay technical take no voice id.

**CONSTRAINT:** A new or edited voice file MUST live under `data/voices/` in this repository, with YAML front matter (`id`, `label`, `sections`, `example`, `tone_notes`, `banned_patterns`, `max_staccato_run`; optional `required_verb_classes`) and two Markdown sections: Generator contract and Evaluator contract. `example` MUST be two short paragraphs rewriting the shared seed in `data/voices/picker.yaml` in that voice (setup, then turn; polish picker hover and polish `voice_profile` demonstration). `tone_notes` and `example` MUST be non-empty: polish empties notes and only sends front-matter JSON to the model, so missing either leaves a voice-named pass without a usable cadence contract.
- Enforcement: File path is `data/voices/<id>.md`; front matter parses; both headings are present; `example` is two paragraphs and clearly a rewrite of the picker seed; `tone_notes` and `example` are present and non-empty.
- Violation: STOP, move the file into the catalog or add the missing heading, example, or `tone_notes`. Do not put phrase banks in strop.

CORRECT:
```text
data/voices/patient-narrator.md
id matches the filename
Evaluator contract repeats the front matter checks in prose
```

PROHIBITED:
```text
Phrase list committed under providers/strop
Voice rules pasted only into a chat, with no catalog file
```

**CONSTRAINT:** MUST NOT apply a section voice during composition. Composition focus stays structural. Polish focus names the voice id.
- Enforcement: `essay_objective_set` for step `composition` does not name a voice id. Step `polish` does, when a section voice applies.
- Violation: STOP, remove the voice id from composition focus and set it on polish.

CORRECT:
```text
essay_objective_set step=polish focus=patient-narrator: unhurried structural diagnosis; verbs collapse, rebuild, populate, absorb; no Look,; no therapy labels
```

PROHIBITED:
```text
essay_objective_set step=composition focus=write this in the patient narrator voice
```

**CONSTRAINT:** Before polish generate, the operator MUST pass the voice file through `notes_files` on `essay_create` or a later notes update, and MUST set polish focus to that voice id plus the generator contract in one sentence.
- Enforcement: Board notes contain `data/voices/<id>.md`. Polish step focus starts with the voice id.
- Violation: STOP, add the file and set focus. Do not polish from memory of a persona.

CORRECT:
```json
{"notes_files": ["data/voices/patient-narrator.md"]}
```

PROHIBITED:
```text
Polish focus: "make it sound like Morgan Freeman"
```

**CONSTRAINT:** A voice change MUST be checked with `make audit-voice POST=<section/slug> VOICE=<id>` before the polish Gate is agreed. Exit 0 is pass. Exit 1 is a failed contract.
- Enforcement: Run the Make target on the English `index.md`. Read each violation line.
- Violation: STOP, revise the prose or the profile. Do not agree the polish Gate over a failing audit.

CORRECT:
```text
make audit-voice POST=human-condition/2026-09-22-same-words-different-worlds VOICE=patient-narrator
```

PROHIBITED:
```text
Agree polish because the paragraph "sounds right" while audit-voice exits 1
```

## Operator runbook

1. Pick the voice id from the routing table. If the piece stays technical, skip the rest.
2. `essay_create` with `notes_files` including `data/voices/<id>.md` plus the source notes.
3. Run spine, arc, and composition with structural focus only.
4. `essay_objective_set` with `step` `polish` and `focus` set to the voice id plus the generator contract (cadence, verbs, bans).
5. `essay_polish_generate`, then `make audit-voice` on the checkout or the Hugo post.
6. Polish Gate reviews the diff. A failing audit is a disagree until the prose or the profile changes.

## Checklist

- [ ] **Catalog file:** `data/voices/<id>.md` exists and `id` matches the filename
      Method: Open the file; compare front matter `id` to the filename
      Pass: They match
      Fail: STOP, rename or fix `id`
- [ ] **Both contracts:** Generator contract and Evaluator contract headings are present
      Method: Search the file for those headings
      Pass: Both present
      Fail: STOP, add the missing section
- [ ] **tone_notes:** Front matter `tone_notes` is present and non-empty
      Method: Read YAML front matter
      Pass: A plain tone brief the model can match
      Fail: STOP, add `tone_notes` before polish
- [ ] **Composition untouched:** composition focus has no voice id
      Method: Read the composition step focus on the board
      Pass: Structural only
      Fail: STOP, move the voice id to polish
- [ ] **Audit:** `make audit-voice` exits 0 on the English post
      Method: Run the target
      Pass: Exit 0
      Fail: STOP, fix prose or profile
