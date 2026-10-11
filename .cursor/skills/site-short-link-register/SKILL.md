---
name: site-short-link-register
description: >-
  Tracks root-level short URLs for carousel CTAs and print-friendly links.
  Source of truth is data/short-links.toml; inventory in data/short-link-register.txt.
  Use before adding deck.cta.shortUrl, Hugo aliases for /slug/ paths, or carousel post_cta slides.
---

# Short-link register

## Files

| File | Role |
|------|------|
| **`data/short-links.toml`** | **Source of truth.** One `[[link]]` per root short path. |
| **`data/short-link-register.txt`** | **Inventory.** Tab-separated list for quick scanning; update when TOML changes. |

## What counts as a short link

- A **single-segment** site-root alias: `/war-game/` (not `/social-protocols/.../`).
- Implemented as Hugo **`aliases`** on the **English** `index.md` for the target bundle.
- Printed on carousel **`post_cta`** slides; encoded in the on-slide **QR** via `deck.cta.shortUrl`.

## Workflow (MUST)

1. Open **`data/short-link-register.txt`** and confirm the slug is **unused**.
2. Add a `[[link]]` row to **`data/short-links.toml`** (`path`, `bundle`, optional `carousel`, `note`).
3. Mirror the row in **`data/short-link-register.txt`** (Active short paths table).
4. Add the same `path` to **`aliases`** on that bundle's English **`index.md`**.
5. Set **`deck.cta.shortUrl`** in `carousel.json`:
   - QR: full URL with scheme, e.g. `https://blog.dualsubstrate.com/war-game/`
   - Printed line: may omit `https://`, e.g. `behaviorengineering.ai/war-game`
6. Keep **`deck.cta.postUrl`** as the canonical long permalink (traceability; not printed when QR is used).
7. Slide footer: short URL + label such as "Scan for full essay + sources"; do **not** print the long path.

## Collision rules

- MUST NOT reuse a `path` for a different post. Retire or rename the old link first.
- MUST NOT add a short path that matches an existing Hugo **section** (`/social-protocols/`, `/human-condition/`, `/x-minds/`, etc.).
- Spanish short links: only if needed; add `aliases` on **`index.es.md`** under `/es/...` (separate register row optional).

## Related

- Carousel CTA authoring: **`.cursor/skills/site-carousel-post/SKILL.md`** → **Post CTA short link**
- Tag inventory (similar pattern): **`.cursor/skills/site-tag-register/SKILL.md`**
