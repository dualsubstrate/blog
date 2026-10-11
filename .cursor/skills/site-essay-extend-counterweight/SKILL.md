---
name: site-essay-extend-counterweight
description: >-
  Forces an essay extension to pay one missing counter-hinge (cost, limit, or
  disconfirming evidence) the source treated as pure upside. Use when the user
  asks for a counterargument, steelman the other side, hidden cost, atrophy
  risk, or a more honest follow-up to an essay (stanceful 400-500 words).
---

# Site essay extend counterweight

**Moral:** The useful sequel is often the cost the original skipped. Name it, then keep the original claim if it still stands.

Load when the operator wants a counterargument, hidden cost, steelman, or "what the piece assumed."

## When to load

- User says counterweight, counterargument, hidden cost, steelman, other side, atrophy, risk
- An extension is about to repeat the source's upside with more examples
- User caps length at 500 words or companion band

MUST NOT load this skill to invert the author's stance into a takedown. MUST NOT add a both-sides shrug that cancels the hinge.

## Core constraints

**CONSTRAINT:** The counterweight MUST be one concrete, novel cost or counter-intuitive mechanism tied to the source's proposed move (example: unpriced coordination friction, strategy migration, legibility cost). MUST NOT be a generic "everything has tradeoffs" paragraph. Strictly avoid consensus clichés (such as GPS studies, generic memory decay, or moralizing friction).

- Enforcement: One named novel cost in the operator hinge; body shows a mechanism for that cost
- Violation: STOP, pick a specific non-cliché cost

CORRECT:
```text
Hinge: automating planning can trigger a false sense of completion, extinguishing the dopamine drive needed to execute.
```

PROHIBITED:
```text
Of course, technology has pros and cons, and balance is essential.
```

**CONSTRAINT:** After the cost is on the table, the close MUST still land: keep, bound, or revise the original move. MUST NOT end in pure negation or in a coaching pep talk.

- Enforcement: Last paragraph states keep / bound / revise in plain verbs
- Violation: STOP, write a diagnostic landing

CORRECT:
```text
Keep the offload for the first three admin steps. Bound it: do not outsource the choice of which side quest matters.
```

PROHIBITED:
```text
So maybe don't use AI. Believe in yourself instead.
```

**CONSTRAINT:** Evidence MUST stay honest. MUST NOT invent studies, percentages, or paper URLs. Use simple analogies (e.g. "the strategy one-way gate") instead of academic jargon. If the cost is a live research question, put it in explore as a Perplexity query, not as a fake citation.

- Enforcement: Every empirical claim has a source already in notes, or is phrased as an open question; simple analogies replace jargon
- Violation: STOP, convert invented facts into explore queries; rewrite jargon

CORRECT:
```yaml
query: "Does offloading executive function tasks to AI tools create a strategy migration where a person stops maintaining internal planning models?"
```

PROHIBITED:
```text
A 2024 Harvard study showed a 37% drop in planning skill...
```

**CONSTRAINT:** Word band still belongs to `site-essay-extend` (400–500 counted words). This skill MUST NOT authorize a longer "balanced essay."

- Enforcement: Same count rule as the extend skill
- Violation: STOP, cut

## Hinge menu by domain (pick one; do not stack)

| Domain / Source upside | Counter-hinge |
| --- | --- |
| Cognition & Tools: Offload logistics to AI | Strategy migration (one-way gate), or premature curiosity closure (plan satisfies dopamine before execution) |
| Social Protocols: Explicit rules & boundaries | Destroys plausible deniability, or creates a legibility tax that stiffens informal relationships |
| Attention & Interest: Interest-driven activation | Unfinished loops when dopamine leaves, or coordination friction with clock-bound teams |
| Institutions & Careers: Generalist synthesis | Credential filters and hiring still pay specialists in many lanes |
| Leverage & Automation: Permissionless execution | Distribution without feedback filters, or slow hunches dying without external containers |

## Steps

1. **Name the upside:** Quote the source close in one line.
2. **Pick one novel cost:** Use the menu or the operator's cost. Wait if two costs would split the piece.
3. **Mechanism:** Explain how the cost happens using simple analogies (attention, incentives, strategy).
4. **Bound:** Close with keep / bound / revise.
5. **Explore:** Add the cost as a Perplexity question via `site-essay-extend-explore`.

## Pre-completion checklist

- [ ] **One novel cost:** A single named non-cliché mechanism of harm
      Method: Highlight it in the hook or first block
      Pass: Named
      Fail: STOP, pick one
- [ ] **Landing:** Close keeps, bounds, or revises the original move
      Method: Read last paragraph
      Pass: Diagnostic verb
      Fail: STOP, rewrite close
- [ ] **Simple language:** Analogies replace academic jargon
      Method: Skim for "metacognitive", "transfer cost", etc.
      Pass: Simple English
      Fail: STOP, reword
- [ ] **Band:** Still ≤500 counted words
      Method: Count
      Pass: Inside band
      Fail: STOP, cut
