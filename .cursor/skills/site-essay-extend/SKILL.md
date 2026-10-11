---
name: site-essay-extend
description: >-
  Drafts a short companion extension to a published Hugo essay or video pick
  (stanceful 400–500 words, one open hinge, catalog voice at polish). Use when
  the user asks to extend this essay, write a follow-up essay, companion piece,
  sequel, or deeper cut under 500 words.
---

# Site essay extend

**Moral:** An extension is a new argument on one hinge the original named and left unpaid. It is not a restack, a longer TL;DW, or a seminar afterword.

Load when the operator asks for an extension, companion, sequel, or follow-up essay off an existing `content/` piece.

## When to load

- User says extend this essay, write a companion, follow-up piece, deeper cut, or sequel
- User caps length at 500 words or companion band
- Source is an existing Hugo bundle (`index.md`) or an approved essay checkout

MUST NOT load this skill for Dig deeper YAML only (use `site-essay-extend-explore`). MUST NOT load it to restyle the original post.

## Core constraints

**CONSTRAINT:** The companion MUST stay in the companion essay band: about 400–500 words of hook + blocks + close. MUST NOT exceed 500 words in that band. `reader_landing` is extra and MUST NOT count toward the 500.

- Enforcement: Word-count hook + body sections + close only; exclude title, front matter, and reader_landing
- Violation: STOP, cut until the band is inside 400–500

CORRECT:
```text
Counted: hook, three blocks, close = 465 words
Uncounted: why_it_matters + takeaways + explore
```

PROHIBITED:
```text
1200-word "fuller version" of the same thesis
```

**CONSTRAINT:** The companion MUST pick **one** unpaid hinge from the source: focus on a novel, non-obvious friction, structural trade-off, or second-order mechanism. MUST NOT reprint the source thesis paragraph, TL;DW, or takeaways as the new body. Strictly avoid consensus clichés (e.g., generic productivity advice, moralizing friction as laziness, or surface-level summaries).

- Enforcement: First operator sentence names the novel hinge; body never pastes the source close
- Violation: STOP, name a single novel hinge and delete restack

CORRECT:
```text
Hinge: execution without permission: AI automates routine logistics so the non-linear mind can ship directly.
```

PROHIBITED:
```text
Recapping the source thesis with new adjectives
```

**CONSTRAINT:** Composition MUST use the shared commute voice in `data/essay-content-forms.yaml`. Polish MUST use the section voice from `.cursor/skills/site-essay-voices/SKILL.md` (`patient-narrator` for `human-condition` and `x-minds`). MUST NOT set polish focus to a celebrity name.

- Enforcement: Composition focus is structural; polish focus starts with a `data/voices/` id
- Violation: STOP, move celebrity wording out of polish; load `site-essay-extend-voice`

CORRECT:
```text
essay_objective_set step=polish focus=patient-narrator: unhurried structural diagnosis
```

PROHIBITED:
```text
polish focus: Morgan Freeman reads this
```

**CONSTRAINT:** MUST emit `reader_landing` after close (why_it_matters, 2–4 takeaways, 2–4 explore items). Takeaways MUST NOT restack the close. Explore MUST prefer full Perplexity questions. MUST NOT invent source or thread URLs.

- Enforcement: YAML or labeled blocks present; no `https://www.perplexity.ai/search/` unless a real completed thread exists in notes
- Violation: STOP, add landing or drop invented URLs

**CONSTRAINT:** MUST NOT use the em dash character (U+2014) in the companion or sidecars.

- Enforcement: Search the draft for U+2014
- Violation: STOP, replace with comma, colon, semicolon, or parentheses

## Steps

1. **Source:** Read the English `index.md` (and `index.es.md` if extending both). Copy title, permalink, and the unpaid hinge in one sentence.
2. **Hinge pick:** Name one of: mechanism (why interest beats importance), combinatorial pattern (beyond the source's example), practical first step, or counter-cost. If the operator did not pick, recommend the counter-cost and wait.
3. **Band:** Draft hook + blocks + close only, commute voice, 600–800 words.
4. **Landing:** Add why_it_matters, takeaways, explore. For explore, load `site-essay-extend-explore` if queries are the job.
5. **Polish:** Load `site-essay-extend-voice`. Run `make audit-voice POST=<section/slug> VOICE=patient-narrator` when a Hugo path exists.
6. **Stop:** Do not export into `content/` unless the operator asked to publish.

## Pre-completion checklist

- [ ] **Band:** Hook + blocks + close are 600–800 words
      Method: Count those sections only
      Pass: 600–800
      Fail: STOP, cut or expand
- [ ] **One hinge:** Body pays a gap the source left open
      Method: Compare close of source vs first claim of companion
      Pass: New claim
      Fail: STOP, rewrite
- [ ] **Voice route:** Polish names `patient-narrator` (or the section table), not a celebrity
      Method: Read polish focus
      Pass: Catalog id
      Fail: STOP, swap
- [ ] **Landing:** why / takeaways / explore present; no invented thread URLs
      Method: Scan labeled blocks
      Pass: All three jobs
      Fail: STOP, fill or delete fake URLs
- [ ] **No em dash:** U+2014 absent
      Method: Search draft
      Pass: Zero hits
      Fail: STOP, replace
