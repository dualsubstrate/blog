---
name: site-tag-unify
description: >-
  Consolidates duplicate or deprecated Hugo tags across content: pick a canonical tag,
  add it to affected posts while optionally keeping legacy tags temporarily, update
  data/tag-deprecations.toml, and regenerate data/tag-register.txt. Keeps EN/ES sibling
  tags aligned. Use when the user asks to unify, dedupe, migrate, or balance tags.
---

# Tag unify (consolidation workflow)

## Preconditions

- Read **`.cursor/skills/site-tag-register/SKILL.md`** and the current **`data/tag-register.txt`**.
- Confirm **`data/tag-deprecations.toml`** does not already define a conflicting `from` (the generator rejects duplicate `from`).

## Steps (MUST)

1. **Pick canonical:** Choose one **`to`** string that new and updated posts should use (often an existing high-count tag from the register).
2. **Record deprecation:** Add a `[[deprecated]]` row with `from` = legacy tag, `to` = canonical (edit **`data/tag-deprecations.toml`**). AI may propose the row; it must land in this file so the Deprecated block in **`data/tag-register.txt`** stays honest.
3. **Gentle migration (default):** On every markdown file under **`content/`** that should roll forward, **append** the canonical tag to **`tags`** if missing. **Keep** the legacy tag in place so old `/tags/<legacy>/` pages stay populated until you strip it later.
4. **Bilingual:** For bundles with **`index.es.md`**, MUST mirror the same **`tags`** list as **`index.md`** (see **`.cursor/skills/site-spanish-translation-content/SKILL.md`**).
5. **Regenerate:** Run **`make tag-register`** (or **`make build`**) and commit **`data/tag-register.txt`** with the content and TOML edits.

## Optional hard cleanup

- After every post dropped the legacy tag, you can remove its `[[deprecated]]` row or leave it (register will show **files still using `from`** = 0).
- Full removal of a tag from all files can make its taxonomy URL disappear; add redirects only if you care about old tag URLs (often skipped).

## YAML hygiene

- Preserve valid YAML; keep **`tags`** as a flow or block list; dedupe within a file if the canonical was already present.
