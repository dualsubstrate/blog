---
name: site-format-chat-neurodivergent
description: >-
  Rewrites a prior chat reply (or user-pasted text) into the skim-first chat
  template in .cursor/rules/always-rules-01-human-interaction.mdc. Preserves every
  fact; does not re-debug or re-run tools. Use when the user asks to format for
  neurodivergent, reformat neurodivergent, format chat, make readable, or after
  a dense incident dump.
disable-model-invocation: true
---

# Format chat (neurodivergent)

## Purpose

Turn a dense or wall-of-text **chat** reply into the tiered layout in **`.cursor/rules/always-rules-01-human-interaction.mdc`** (Chat shape). This is a **presentation** pass only.

## MUST NOT

- Re-run tools, re-read logs, or change conclusions.
- Drop HTTP codes, error names, file paths, flags, URNs, scopes, or numeric limits that appeared in the source.
- Add investigation, options, or steps that were not in the source.
- Edit files under `content/` unless the user explicitly asks.

## MUST

1. Follow **`.cursor/rules/always-rules-01-human-interaction.mdc`** Chat shape (gate, section order, voice).
2. Use the **immediately prior assistant message** as the source unless the user pasted replacement text.
3. Put a plain-English **outcome** in **Direct answer** (one line; two only if a caveat is essential).
4. Use **Most likely cause** for a single mechanism when the source states one.
5. Use **Do this next** for numbered actions only when the source lists actions or clear next steps.
6. Use **Uncertainty / assumptions** only for genuine unknowns from the source.
7. Move leftover identifiers (paths, `package.Func`, CLI flags, API codes) into **Details** (no icon heading) as one fact per bullet, or attach them to the relevant **Do this next** step.

## Output template

```markdown
✅ Direct answer
<one plain-English outcome line>

🔍 Most likely cause
<one line, only if source states a cause>

➡️ Do this next
1. ...
2. ...

❓ Uncertainty / assumptions
- ... (only if needed)

Details
- ... (only if identifiers do not fit above)
```

Omit empty sections. Do not repeat the direct answer in later sections.

## Jargon

- Replace or gloss jargon when the source allows a plain phrase without losing accuracy.
- Keep exact strings (paths, flags, error codes, URNs) literal inside **Details** or step lines.

## Example (incident-shaped source)

**Source (dense):** Post published; CI failed on post-verify. POST 201 and urn:li:share:…. GET used unencoded URN; 400 ILLEGAL_ARGUMENT. Fix: encodeRestLiResourceKey in internal/linkedinapi/restli.go. Day 20 content already on LinkedIn. Merge fix; token needs read scope. Until then -no-verify-commentary. Commentary 1208 bytes vs ~800 limit; check URLs on live post.

**Formatted:**

✅ Direct answer
Your LinkedIn post is live; CI failed only on verify because the follow-up GET used an unencoded share URN in the path.

🔍 Most likely cause
`GetPost` used `url.PathEscape`, which leaves `:` unencoded; LinkedIn Rest.li expects `urn:li:share:…` encoded as `urn%3Ali%3Ashare%3A…` in the path.

➡️ Do this next
1. Merge the fix (`encodeRestLiResourceKey` in `internal/linkedinapi/restli.go` plus tests).
2. Re-run CI; verify should pass if the token has read scope (`r_member_social` / `r_organization_social`).
3. Until merged, use `-no-verify-commentary` to avoid a false verify failure (posting is already done).
4. On LinkedIn, check the live post for truncated links (commentary was ~1208 UTF-8 bytes vs ~800 image-post limit).

Details
- POST `/rest/posts` returned 201 and `urn:li:share:7462984940427677697`.
- GET error: 400 `ILLEGAL_ARGUMENT`, "Syntax exception in path variables".
- Day 20 psych-fitness content: no repost needed unless verify matters for another reason.
