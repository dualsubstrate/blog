# Voice catalog

Host-owned narrator contracts for essay polish. Composition stays on the shared commute voice in `data/essay-content-forms.yaml`. Polish applies one file from this directory when the operator names it.

The portable checker lives in strop (`pkg/evaluation/voice`). These files are the phrase banks and cadence rules. Strop does not store section names or persona labels.

Polish does **not** send the Markdown body of a voice file to the model (notes are emptied to avoid truncation). The only voice channel at generate time is compact JSON from YAML front matter: `banned_patterns`, `required_verb_classes`, `max_staccato_run`, `tone_notes`, and `example`. That is why `tone_notes` and `example` are required for polish: rules plus a two-paragraph cadence demonstration.

## Section routing

| Section | Voice file |
| --- | --- |
| `human-condition` | `patient-narrator.md` |
| `x-minds` | `patient-narrator.md` |
| `social-protocols` | `clint-eastwood.md` (or `lean-realist.md`) |
| `mind-infrastructure` | `cognitive-pragmatist.md` |

Technical posts that should stay technical do not take a narrator file. Leave polish focus empty.

## Two contracts in each file

Generator contract: how the polish model should write (cadence, verbs, what to refuse). The Markdown body still helps humans and `polish_focus` paste; the model sees `tone_notes` and bans via `voice_profile` JSON. All catalog voices refuse abstract "the room" metonymy (group/audience); literal physical rooms stay allowed. Do not add "the room" to `banned_patterns` (substring audits cannot distinguish abstract vs literal use).

Evaluator contract: what a deterministic audit and the `voice_fidelity` role can fail closed (banned phrases, staccato run, required verbs, tone). Those checks are the YAML front matter.

## Front matter

```yaml
---
id: patient-narrator
label: Patient Narrator
sections: ["human-condition", "x-minds"]
banned_patterns:
  - "Look,"
tone_notes: >-
  Unhurried structural diagnosis in plain English. Diagnose the gap.
max_staccato_run: 2
---
```

| Field | Required | Meaning |
| --- | --- | --- |
| `id` | yes | File id used by `make audit-voice VOICE=` and polish focus |
| `label` | yes | Human name |
| `sections` | yes | Hugo section folder names this voice is for |
| `example` | yes for polish picker and polish generate | Two short paragraphs rewriting the shared seed in `picker.yaml` in that voice (setup, then turn). Also shipped in `voice_profile` JSON so the polish model can match the cadence. |
| `tone_notes` | yes for polish | Plain-language tone the model must match when `voice_profile` is active. Without this, polish stays a light tidy of composition diction. |
| `banned_patterns` | no | Case-insensitive substrings. A hit fails the audit. Include pep-talk / flattery the voice refuses. |
| `required_verb_classes` | no | Whole words. Optional list; when present, a body paragraph with none of them fails. Standard catalog profiles omit this so they generalize across different topics |
| `max_staccato_run` | no | Longest allowed run of similar short sentences. Default 2. A longer run fails |

`picker.yaml` holds the shared two-paragraph seed. The polish command-center shows that seed once, then each chip’s `example` rewrite on hover.

`cmd/audit-essay-voice` reads this front matter and the post body. It skips the README and `picker.yaml`.
