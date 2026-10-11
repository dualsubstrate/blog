---
name: site-pr-review-comments
description: >-
  Triage GitHub pull request review comments, apply scoped fixes, push, and
  resolve review threads. Use when the user asks to review PR comments, fix
  review feedback, address review comments, resolve PR comments, or resolve
  review threads after fixes. Pair with babysit for full merge-ready loops
  (conflicts, CI, and comments together).
---

# PR review comments

## Scope

This skill covers **review feedback only**: inline review threads, review summaries, and actionable bot comments (for example Bugbot, Gemini Code Assist).

For merge conflicts, unrelated CI failures, or a full merge-ready loop, use **`~/.cursor/skills-cursor/babysit/SKILL.md`** instead.

## Workflow

Copy and track progress:

```
- [ ] 1. Find the PR for the current branch
- [ ] 2. List unresolved review threads (skip resolved)
- [ ] 3. Triage each comment (fix valid, explain invalid)
- [ ] 4. Commit and push scoped fixes
- [ ] 5. Resolve fixed threads on GitHub
- [ ] 6. Re-check CI
```

### 1. Find the PR

Run from the repo root (use `gh` for all GitHub tasks):

```bash
git branch --show-current
gh pr list --head "$(git branch --show-current)" --json number,title,url,state
```

If the user gave a PR number or URL, use that instead.

### 2. List unresolved review threads

Prefer **unresolved inline threads** first. Use GraphQL and filter `isResolved: false`:

```bash
gh api graphql -f query='
query($owner: String!, $name: String!, $number: Int!) {
  repository(owner: $owner, name: $name) {
    pullRequest(number: $number) {
      reviewThreads(first: 100) {
        nodes {
          id
          isResolved
          path
          line
          comments(first: 3) {
            nodes {
              body
              author { login }
            }
          }
        }
      }
    }
  }
}' -f owner=OWNER -f name=REPO -F number=PR_NUMBER
```

Read only each thread's `path`, `line`, first comment `body`, and `author.login`. Do not dump full JSON into chat.

Also skim:

```bash
gh pr view PR_NUMBER --json reviews,comments
gh api repos/OWNER/REPO/pulls/PR_NUMBER/comments --jq '.[] | {path, line, body, user: .user.login}'
```

Skip:

- Preview/deploy bot noise (for example pr-preview-action)
- Already resolved threads
- Comments that are questions or style nits with no clear fix unless the user asked to address them

### 3. Triage

For each unresolved thread:

| Verdict | Action |
|---------|--------|
| **Valid** | Implement the smallest correct fix in scope |
| **Partially valid** | Fix the real issue; note what you declined |
| **Invalid / out of scope** | Do not change code; explain why in chat (and optionally reply on the thread) |

Rules:

- MUST fix only what the comment requests; no drive-by refactors
- MUST NOT weaken CI, skip hooks, or change workflows just to silence checks
- MUST NOT resolve a thread until the fix is pushed (or you documented why no fix is needed)

### 4. Commit and push

Follow the user's git commit rules:

1. `git status`, `git diff`, `git log -1` (parallel)
2. Stage only files related to the review fixes
3. Commit with a message focused on **why** (HEREDOC)
4. `git push -u origin HEAD`

Only commit when the user asked to fix comments or the workflow implies ship (for example "review and fix", "resolve comments").

### 5. Resolve threads on GitHub

After fixes are pushed, resolve each addressed thread:

```bash
gh api graphql -f query='
mutation($threadId: ID!) {
  resolveReviewThread(input: { threadId: $threadId }) {
    thread { isResolved }
  }
}' -f threadId=PRRT_...
```

Resolve only threads that are fixed or explicitly declined with a reply.

Optional: reply before resolving so reviewers see the commit:

```bash
gh api repos/OWNER/REPO/pulls/PR_NUMBER/comments/COMMENT_ID/replies \
  -f body='Fixed in COMMIT_SHA: brief summary.'
```

### 6. Re-check CI

```bash
gh pr checks PR_NUMBER
```

If checks fail on the new commit, fix failures in PR scope or report blockers. If failures look unrelated, check whether the branch is behind base and merge or rebase latest `main` before chasing unrelated fixes.

## Reporting back

Summarize for the user:

1. How many threads were open
2. What you fixed (file + one line each)
3. What you skipped and why
4. Commit SHA and whether threads are resolved
5. CI status

Keep chat concise; link the PR URL.

## Examples

**User:** "Review the PR comments and fix them"

1. Find PR for current branch
2. Fix valid threads
3. Commit + push
4. Resolve threads
5. Report CI

**User:** "Resolve comments?"

Assume fixes are already pushed. List unresolved threads; if any still need code changes, fix first, then resolve all addressed threads.

**User:** "Address the Gemini review on PR 72"

Use PR 72 directly; same workflow.
