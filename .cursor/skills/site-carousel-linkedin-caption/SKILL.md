---
name: site-carousel-linkedin-caption
description: >-
  Drafts a short LinkedIn document caption from carousel.json via local Gemma 4
  and writes it to linkedin-carousel.txt. Use when the user asks for carousel
  LinkedIn text, a LinkedIn document caption, tap-to-explore copy, Gemma 4
  carousel caption, linkedin-carousel.txt, or a short post that sits on a
  LinkedIn carousel (not the full linkedin.txt essay).
---

# LinkedIn carousel caption (Gemma 4)

## Goal

Turn **`carousel.json` slide copy** into a **short LinkedIn document caption** (hook, one mechanism paragraph, `Tap to explore`, essay URL). Write it to **`linkedin-carousel.txt`** in the same bundle.

This is **not** **`linkedin.txt`**. That file stays the full standalone LinkedIn essay (**`.cursor/skills/site-linkedin-post/SKILL.md`**).

## When to use

- "caption for the carousel", "LinkedIn document text", "tap to explore", "ask Gemma to write the carousel post"

Need a deck first? Use **`.cursor/skills/site-carousel-post/SKILL.md`**.

## Output files

| File | Role |
|------|------|
| `linkedin-carousel.candidates.txt` | Three numbered Gemma drafts |
| `linkedin-carousel.txt` | One paste-ready caption after the user picks |

MUST NOT overwrite **`linkedin.txt`**.

## Caption shape

Plain text. No markdown. No em dash (U+2014).

```
[Curiosity question?] [optional one emoji]

[One short paragraph: mechanism + one construct from the deck]

Tap to explore [the mechanism the slides pay off].

Essay:
<shortUrl or postUrl from deck.cta>
```

Rhythm sample (different topic; copy shape only):

```
Why do we hide our true motives, even from ourselves? 🎭

We often use an acceptable outward story as a "deniability hood." It lowers conflict and protects our reputation, while the real push: spite, belonging, or fear, keeps driving us from the shadows.

Tap to explore the cold-hot empathy gap and why you always have two reasons for what you do.
```

## Defaults

| Setting | Value |
|---------|--------|
| Gateway | `LOCAL_LLM_BASE_URL` or `http://127.0.0.1:1320/v1` |
| Model | `LOCAL_LLM_MODEL` or `gemini_free/gemma-4-26b-a4b-it` |
| Source | First variant (`a`) of each slide in `carousel.json` |
| Essay URL | `deck.cta.shortUrl`, else `deck.cta.postUrl` |

Pack: **`packs/caption.md`**.

## Procedure

1. Resolve the bundle: path the user names, focused `carousel.json`, or sibling of the focused post.
2. Run from repo root (do not hand-roll the API unless the script fails):

```bash
python3 .cursor/skills/site-carousel-linkedin-caption/scripts/draft_caption.py \
  content/<section>/<slug>/carousel.json
```

3. Surface the **three numbered captions** from stdout (also written to `linkedin-carousel.candidates.txt`).
4. Light check only: drop an em dash if Gemma added one; reject slips that are not in the deck (*wired to*, *mathematical attractor*) by asking the user or picking another variant. MUST NOT rewrite the argument.
5. Stop and ask which variant to keep (`1` / `2` / `3` / cancel).
6. After they pick, write the paste file:

```bash
python3 .cursor/skills/site-carousel-linkedin-caption/scripts/draft_caption.py \
  content/<section>/<slug>/carousel.json \
  --pick 1
```

`--pick` reads the candidates file (no second Gemma call). `--out` overrides the final path.

If they already said `write 1` / `use 2` in the same request, skip the wait and pick immediately.

## Gateway down

MUST say the gateway is down. MUST NOT invent captions and pretend Gemma wrote them.

Optional debug (no model):

```bash
python3 .cursor/skills/site-carousel-linkedin-caption/scripts/draft_caption.py \
  content/<section>/<slug>/carousel.json \
  --dump-beats
```

## Related

- Deck authoring: **`.cursor/skills/site-carousel-post/SKILL.md`**
- Full LinkedIn essay: **`.cursor/skills/site-linkedin-post/SKILL.md`**
- Local Gemma stack: **`.cursor/skills/site-revise-prose/SKILL.md`**
