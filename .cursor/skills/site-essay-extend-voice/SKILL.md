---
name: site-essay-extend-voice
description: >-
  Maps celebrity or documentary-narrator requests (including Morgan Freeman)
  onto the host patient-narrator catalog for essay extensions, and applies that
  voice only at polish. Use when the user asks for Morgan Freeman voice,
  documentary narrator, unhurried gravitas, or patient-narrator on a companion
  piece.
---

# Site essay extend voice

**Moral:** The site already has the unhurried narrator. Name the catalog file. Do not impersonate a celebrity.

Load when the operator asks for Morgan Freeman, a documentary narrator, gravitas, or patient-narrator polish on an extension.

## When to load

- User says Morgan Freeman, documentary voice, National Geographic narrator, unhurried, gravitas
- User says polish this extension / companion in the usual essay voice

MUST NOT load this skill to invent a new `data/voices/` celebrity id. MUST NOT apply voice during composition.

## Core constraints

**CONSTRAINT:** For `human-condition` and `x-minds` extensions, polish MUST use voice id `patient-narrator` and file `data/voices/patient-narrator.md`. MUST NOT put a celebrity name, impression notes, or catchphrase list in `essay_objective_set` polish focus.

- Enforcement: Polish focus starts with `patient-narrator`; grep focus for Freeman / impersonat / "in a low voice"
- Violation: STOP, rewrite focus from the generator contract on that file

CORRECT:
```text
patient-narrator: unhurried structural diagnosis; uneven cadence; diagnose the gap; no Look,; no pep talk
```

PROHIBITED:
```text
Morgan Freeman: warm, wise, "let me tell you about the universe"
```

**CONSTRAINT:** MUST pass `data/voices/<id>.md` through `notes_files` (or `essay_polish_voice_set`) before polish generate. MUST run `make audit-voice POST=<section/slug> VOICE=patient-narrator` when a Hugo English `index.md` exists. Exit 1 is a fail.

- Enforcement: Board notes include the voice path; audit exit code is 0
- Violation: STOP, add the file and re-audit; do not agree polish over a failing audit

CORRECT:
```text
notes_files includes data/voices/patient-narrator.md
make audit-voice POST=x-minds/2026-09-26-the-octopus-advantage VOICE=patient-narrator
```

PROHIBITED:
```text
Polish from chat memory of a movie narrator
```

**CONSTRAINT:** Cadence MUST match the catalog generator: uneven sentence length, long landscape then short fact, structural verbs, diagnose not soothe. Frame extensions in plain English using concrete physical analogies. Diagnose the structural gap without moralizing friction or treating human limits as personal defects. MUST NOT stack more than `max_staccato_run` (2) similar short sentences. MUST NOT use banned_patterns from the voice file.

- Enforcement: `make audit-voice`; also grep banned strings as writer wording
- Violation: STOP, revise prose or (only if the operator asked) the catalog file in `data/voices/`

CORRECT:
```text
A long sentence carries how interest gates effort. The calendar line does not.
```

PROHIBITED:
```text
Look, here's the thing. You're not broken. You're an octopus. Own it.
```

**CONSTRAINT:** MUST NOT mimic a real person's vocal tics, catchphrases, or biographical jokes. The mapping is **tone contract**, not impersonation.

- Enforcement: Scan draft for celebrity name, "I have been narrating", penguin / Shawshank jokes
- Violation: STOP, delete impersonation; keep structural diagnosis

CORRECT:
```text
Prose a cold reader would attribute to the site, not to a film actor
```

PROHIBITED:
```text
And that's when the universe, in all its mystery, invented the to-do list.
```

## What the operator meant (translate, then discard the celebrity label)

Keep these as internal checks. MUST NOT paste them into polish focus as a persona.

| Wanted effect | Catalog move |
| --- | --- |
| Unhurried, low heat | Long sentence for the mechanism; short sentence for the fact |
| Authority without coaching | Diagnose the gap; no pep close |
| Documentary clarity | Name what the mechanism does; gloss jargon in the same breath |
| Warmth | Still no soothe; respect is in precision, not in consolation |

## Steps

1. **Route:** Confirm section. `x-minds` / `human-condition` → `patient-narrator`. Other sections → `.cursor/skills/site-essay-voices/SKILL.md` table.
2. **Translate:** If the user said Morgan Freeman, reply once that polish uses `patient-narrator`, then drop the celebrity name from subsequent tool calls.
3. **Attach:** `essay_polish_voice_set` or `notes_files` with `data/voices/patient-narrator.md`.
4. **Focus:** One sentence: id + generator contract (cadence, verbs, bans).
5. **Audit:** `make audit-voice` on the English post or checkout body.
6. **ES:** If a Spanish companion exists, keep the same structural diagnosis; do not calque English cadence. Load `site-revise-spanish` when the operator asks.

## Pre-completion checklist

- [ ] **Catalog id:** Polish focus starts with `patient-narrator` (or the table id)
      Method: Read the focus string
      Pass: File id present
      Fail: STOP, replace celebrity text
- [ ] **File attached:** Voice markdown is in notes
      Method: Board notes / notes_files
      Pass: Path present
      Fail: STOP, add it
- [ ] **Audit:** `make audit-voice` exit 0 when a post path exists
      Method: Run the target
      Pass: Exit 0
      Fail: STOP, fix prose
- [ ] **No impersonation:** Celebrity name absent from body and focus
      Method: Grep Freeman / impersonat
      Pass: Zero
      Fail: STOP, delete
