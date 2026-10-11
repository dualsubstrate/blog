---
name: site-revise-format
description: >-
  Ship-time mechanics for content/ and sidecars: MUST search for em dashes (U+2014)
  and run a full emphasis (Markdown bold) audit per revise-emphasis. Punctuation +
  emphasis in one pass. Use for revise-format, em dashes, too much bold, formatting
  only, or last step before publish. While authoring new bold, see revise-emphasis.
---

# Revise format (punctuation + emphasis)

## Purpose

One **last-mile** pass before publish:

1. **Punctuation:** zero em dashes (`—`).
2. **Emphasis:** restrained `**bold**` per type (rulebook: **`.cursor/skills/site-revise-emphasis/SKILL.md`**).

Authors often still call out dashes and bold because agents skip passive checks; this skill **requires searches**, not a vague “check formatting.”

**Also covers** quick Step 4 mechanics: US English spot-check, title/Claim emoji rules (see **`.cursor/skills/site-revise-post/SKILL.md`** → Step 4).

| When | Skill |
|------|--------|
| **Shipping** / last sweep | **This skill** (punctuation + emphasis) |
| **Adding or fixing bold while drafting** | **`.cursor/skills/site-revise-emphasis/SKILL.md`** |
| Full checklist | **`revise-post`** (format phase = this skill) |

## Scope

- Default: one file or page bundle under **`content/`** (`index.md`, `*.es.md`, `substack.md`, `linkedin.txt`, `facebook-*.txt` when user says whole bundle).
- **MUST** scan **front matter** and **body** (and sidecars in scope).
- **MUST** read **`revise-emphasis`** before the emphasis audit (type-specific span counts: Claim, grounding, sayings `tldr`/`fluff`, etc.).

## Execution (MANDATORY — both pillars)

### 1. Em dash (MUST NOT skip)

Search target path(s) for U+2014 (`—`), not hyphen `-`.

- Tool: **`Grep`** for `—`, or `rg -n $'\u2014'` on the path.
- **Fail** if any hit in scope.
- **Fix:** comma, semicolon, colon, or parentheses (site rule).

```markdown
### Em dash hits
| Line | Excerpt | Proposed fix |
|------|---------|--------------|
```

Zero hits → **`Em dash: 0 hits (Pass)`**.

### 2. Emphasis audit (MUST NOT skip)

**Read first:** **`.cursor/skills/site-revise-emphasis/SKILL.md`** (restrained default; per-type tables).

1. **Grep** `\*\*` on every file in scope.
2. Per **block** (Claim/`description`, `grounding`, each body paragraph, each sidecar paragraph, sayings `tldr`/`fluff` when present):
   - **Count** bold spans.
   - Compare to **revise-emphasis** limits (~2–5 per short card block; sayings tables if `type: sayings`).
   - **Flag Trim:** wall-to-wall bold, decorative bold on filler, or span count over type limit.
   - **Flag Review:** legacy cognitive-memetics heavy bold (normalize only if user asked).

**Fix:** markup-only (remove or move `**`); MUST NOT rewrite author prose.

```markdown
### Emphasis audit
| Block | Spans | Limit (type) | Verdict | Action |
|-------|-------|--------------|---------|--------|
```

**Verdict:** Pass / Trim / Review.

### 3. Other mechanics (quick)

- [ ] US English in default-language `content/` (spot UK: behaviour, centre, generalise).
- [ ] No decorative emoji in Claim / grounding.
- [ ] Title leading emoji: max 2 (social-protocols, human-condition, mind-infrastructure, x-minds); **none** in cognitive-memetics `title`.

## Output format

```markdown
## Revise format: [path(s)]

**Em dash:** Pass (0) | Fix Needed (N)
**Emphasis:** Pass | Fix Needed (blocks listed)
**Other mechanics:** Pass | Fix Needed

| # | Location | Pillar | Issue | Before | After |
|---|----------|--------|-------|--------|-------|

Apply all? (y / cancel / apply em dash only / apply emphasis only)
```

**Default:** analysis only; write after user confirms.

## Pass / Fail

- **Pass:** 0 em dashes; every block **Pass** on emphasis; other mechanics pass.
- **Fix Needed:** any em dash OR any block **Trim**.

**MUST NOT** mark Pass without **both** em dash search and emphasis audit.

## After apply

- **`hugo build`** (or `make build`) when `content/**/*.md` changed.

## Related

- **`.cursor/skills/site-revise-emphasis/SKILL.md`** — rulebook while authoring
- **`.cursor/skills/site-revise-post/SKILL.md`** → Step 4
- **`.cursor/skills/site-revise-post/SKILL.md`** → format phase (last)
- **`.cursor/rules/site-content-markdown-writing.mdc`**
- **`.cursor/rules/always-rules-0-ai.mdc`** (no `—` in generated text)
