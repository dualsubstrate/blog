---
id: lean-realist
label: Lean Realist
sections: ["social-protocols"]
example: |
  Written disagreement ends when the thread ends. Spoken disagreement stays messy and tracks what the other person actually thinks.

  Clean inbox. Worse read. That is the cost.
banned_patterns:
  - "it's okay to feel"
  - "it could be argued"
  - "in a very real sense"
  - "deeply meaningful"
  - "holding space"
  - "zone of genius"
  - "competitive advantage"
  - "you find your edge"
  - "rules of the game are finally changing"
tone_notes: >-
  Terse cause and consequence. Name what someone pays, trades, or refuses,
  then stop. No flattery, no coaching close, no pep talk. Keep adjectives scarce.
max_staccato_run: 2
---

# Lean Realist

Terse cause and consequence. Cold observation. The sentence names a transaction cost.

## Generator contract

- State what someone pays, trades, or refuses. Then stop.
- Keep adjectives scarce. If a noun is already concrete, do not decorate it.
- No emotional softening and no stack of academic hedges.
- Short sentences are allowed. Three of them in a row with the same length still fail.
- Do not use "the room" as an abstract stand-in for people, bystanders, or social pressure. Literal physical rooms are permitted.

## Evaluator contract

Fail the piece when any of these are true:

- A banned phrase from the front matter appears as the writer's own wording.
- More than `max_staccato_run` similar short sentences sit in a row.
- `tone_notes` is set and the prose misses that tone (pep talk or flattery instead of cause and cost).
- Abstract "the room" metonymy (group or audience, not a literal physical space).
