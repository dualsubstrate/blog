---
name: site-tag-register
description: >-
  Uses the generated site tag inventory (data/tag-register.txt) and deprecations
  (data/tag-deprecations.toml) so agents reuse existing Hugo tags when they fit and avoid
  accidental near-duplicates. New tags are allowed when nothing matches. Use when editing
  front matter tags, choosing taxonomy hooks, or before adding content; pair with
  tag-unify for consolidation work.
---

# Tag register (inventory before you tag)

## Files

| File | Role |
|------|------|
| **`data/tag-register.txt`** | **Generated.** Active tags with per-file counts; Deprecated rows (from + to + files still using `from`). Do **not** edit by hand. |
| **`data/tag-deprecations.toml`** | **Source for Deprecated.** Add `[[deprecated]]` rows with `from` / `to`. AI may propose edits here; humans merge like any content. |

Regenerate the register with **`make tag-register`** or any **`make build`** / CI Hugo build (build runs tag-register first).

## Workflow (MUST)

1. Open **`data/tag-register.txt`** and skim **Active tags** (and **Deprecated tags** if present).
2. Prefer an **existing** active tag when it matches the post’s angle (see **`.cursor/skills/site-claims-content/SKILL.md`** *Tag voice* where it applies).
3. If nothing fits, add a **new** tag string in front matter. The next build refreshes the register so later posts can reuse it.
4. For a legacy string listed under **Deprecated**, use **`to`** for **new** work instead of **`from`**. Old pages may still carry **`from`** until someone runs a unify pass.

## Deprecations (editing `data/tag-deprecations.toml`)

Use TOML tables in order (duplicate `from` values are rejected by the generator):

```toml
[[deprecated]]
from = "OldHashtagStyle"
to = "NewCanonicalStyle"
```

Then run **`make tag-register`** and commit both the TOML change and the updated **`data/tag-register.txt`**.

## Limits

- The register lists **exact strings** in use; it does not guess synonyms. Near-duplicates need human or **tag-unify** judgment.
- Hugo tag URLs derive from the string; renaming tags everywhere changes listing URLs (usually acceptable here; see **tag-unify** if you need a gentle migration).
