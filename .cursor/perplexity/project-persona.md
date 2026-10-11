# Project Perplexity persona (behaviorengineering.ai)

Filled overlay for this Hugo site. Pack skill load order: this file first; pack `reference/project-persona.md` only if this path is missing.

**MCP how-to:** `.cursor/skills/perplexity-browser-research/SKILL.md`

## Who we are

- **Project / product:** behaviorengineering.ai (Hugo LoveIt site; local folder `site`; remote `github.com/xynova/behaviour-engineering`)
- **Audience:** Public site readers (smart strangers on a phone), then the author editing `content/`
- **What Perplexity is for here:** Fact-check posts and grounding sources (deep-research pack); cold-read prose for human flow and AI artifacts (prose-review pack)

## Voice when writing prompts

- Prefer **concrete nouns** and named claims over vague “research this topic.”
- Ask for **primary sources** and study design when facts matter.
- Say what is **out of scope** (Hugo theme internals, Go tooling, Spanish translation, SEO) so the model does not wander.
- **US English** in prompts and requested rewrites.
- **MUST NOT** use the em dash character (U+2014) in suggested rewrites.

## Risk and privacy

- **Sensitive labels:** keep out of `title_hint`. Prompt body may include repo context; cloud upload is the user's risk choice.
- Treat every export as **research input**, not shippable fact.
- **MUST NOT** paste Perplexity prose into `content/` or `docs/research/` without human review and the matching type / revise skill.

## Validate-before-write

1. Compare claims to primary sources and the target post under `content/`.
2. Reject blog-only citations presented as binding authority.
3. Apply edits through the matching type skill and `.cursor/skills/site-revise-post/SKILL.md`, not wholesale paste.
4. Optional research notes: `docs/research/YYYY-MM-DD-perplexity-<slug>.md` after human review. Do not commit raw exports.

## Pack choice

| Need | Pack | Mode |
|------|------|------|
| Facts, sources, grounding | `.cursor/perplexity/packs/deep-research.md` | deep |
| Human voice, AI artifacts | `.cursor/perplexity/packs/prose-review.md` | search |
| Local evaluate-only (no Perplexity) | `.cursor/skills/site-revise-prose/SKILL.md` | local API |

**Pick one pack per run.** Fact-check and prose review are separate submits.

## Prepare before filling a pack

Read these repo paths first:

- Target page bundle under `content/` (`index.md`, optional `index.es.md`, sidecars)
- Matching type skill (`.cursor/skills/site-claims-content/SKILL.md`, `site-video-content/SKILL.md`, etc.)
- `.cursor/rules/site-content-placement.mdc` when section or type is unclear
- Optional: `tmp/articles/`, prior `docs/research/`, `data/tag-register.txt`
- For prose review: `.cursor/skills/site-revise-post/SKILL.md` Step 2 so the artifact list matches site rules
