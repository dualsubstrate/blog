# Perplexity share automation (Cursor browser)

Repository scripts cannot call Perplexity APIs to change thread visibility. When the operator is logged into Perplexity in **Cursor IDE browser**, the agent may run this UI sequence per thread, then record results in `share-manifest.json`.

## Preconditions

- Thread was started from [behaviorengineering blog](PERPLEXITY-PROJECT.md).
- `go tool task explore-share-prepare CANDIDATES=tmp/explore-proposals/<slug>/candidates.json` already ran.

## Agent sequence (one thread)

1. `browser_navigate` to `https://www.perplexity.ai/search/<thread_id>`.
2. `browser_snapshot` and confirm the answer loaded (not a login wall).
3. Click **Share** (thread header).
4. Under **General access**, select **Anyone with the link can view** (globe icon). If already selected, leave it.
5. **Copy link** and confirm it matches the `candidates.json` URL (same UUID).
6. Tell the operator to run (or run via make after operator confirms UI):
   ```bash
   go tool task explore-share-confirm CANDIDATES=... CANDIDATE_ID=<id>
   ```
7. Cold reader check (required):
   - Operator opens the URL in a **private browser window** where they are **not** logged into Perplexity.
   - Pass: full research visible. Fail: login or empty thread → fix share, re-confirm.
8. Record verification:
   ```bash
   go tool task explore-share-verify-cold CANDIDATES=... CANDIDATE_ID=<id> PROBE=1 OPERATOR=1
   ```
   `PROBE=1` runs an unauthenticated HTTP heuristic (may false-fail on bot blocking). `OPERATOR=1` records the human incognito pass; use both when possible.

## When automation stops

- Share menu does not open or access level cannot be set: operator completes steps 3–7 manually, then `explore-share-confirm` + `explore-share-verify-cold OPERATOR=1`.
- Perplexity **Exit incognito** (private session inside the product) is not a cold-reader test. Exit that mode before sharing; verify in Safari/Chrome private window.

## Gate before Hugo apply

`go tool task explore-apply` refuses to run (review or apply) until `go tool task explore-share-check PROPOSAL=... CANDIDATES=...` passes for every approved URL in the proposal.
