---
name: site-ai-for-general-audience
description: >-
  Writes or rewrites AI and ML technical content for a non-specialist audience.
  Translates jargon into plain-language equivalents without dumbing down the
  accuracy. Use when the user asks to explain an AI concept for a general
  audience, mentions "non-ML people", "lay audience", "plain English", or asks
  to demystify AI/ML terminology in a post, article, or social copy.
---

# Writing AI/tech content for a general audience

## Goal

Deliver the correct technical idea in plain language. The reader should
understand what a system does and why it works that way, without needing a
machine learning background.

## Core rules

- MUST NOT use ML jargon without immediately explaining it in the same sentence.
- MUST NOT treat removing jargon as removing accuracy. Find the plain-English equivalent.
- MUST NOT use dramatic framing ("it sounds like magic", "this is revolutionary") to bridge a gap you have not filled with an explanation.
- MUST prefer concrete behavior over abstract labels. Say what the system actually does, not what it is called.

## Jargon translation patterns

| ML term | Plain-language equivalent |
|---|---|
| frozen weights / frozen model | the model's skills are locked and do not change |
| gradient / backpropagation | the training signal that tells the model what to improve |
| orchestrator | the model that decides what to do and who to delegate to |
| sub-agent | a separate model called like a tool or an API |
| environment observation | information the model receives but cannot train on |
| credit assignment | figuring out which part of the system caused a failure |
| MoE (Mixture of Experts) | a neural network architecture where different internal pathways handle different types of input |
| RL (Reinforcement Learning) | training by reward and penalty based on outcomes, not labeled examples |
| inference | running a trained model to get an answer (not training) |
| checkpoint | a saved snapshot of a model at a point in training |
| decoupled training | training one part of the system while keeping the rest fixed |

## Analogies: use only when they add clarity, drop when they add noise

Good analogies clarify the mechanism. Bad analogies swap one mystery for another.

- Use: "like making an API call to another tool" for sub-agents.
- Use: "training the manager, not the workers" for decoupled orchestrator training.
- Avoid: overextended metaphors that need their own explanation.
- Avoid: "magic", "trick", "under the hood" as substitutes for an explanation.

## Structure for technical posts

1. **State what the system does** in one sentence before any jargon appears.
2. **Name the confusing term** and immediately say what it means in plain English.
3. **Explain the why**: what problem does this design choice solve?
4. **Separate train-time from deploy-time** when a system behaves differently in each phase.
5. **Give one concrete example** (benchmark, real task, or observed behavior).
6. **Hedge headline numbers** clearly: "harness-dependent, verify on your own setup."

## Anti-patterns

- Translating jargon into different jargon ("frozen weights" -> "static gradient-free modules").
- Explaining what a thing is called without explaining what it does.
- Saving the plain-language explanation for the end after paragraphs of jargon.
- Treating analogies as optional decoration rather than load-bearing explanation.
