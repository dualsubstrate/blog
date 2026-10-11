# Perplexity threads: public share before Hugo ship

Explore further `perplexity_thread` rows MUST open for **cold readers** (logged out or incognito), not only inside your account.

## Critical: default is private

New searches (including MCP `perplexity_research` outside the blog project) start as **Only people with access can view**. Explore batch threads MUST be created inside the project first; see [PERPLEXITY-PROJECT.md](PERPLEXITY-PROJECT.md).

## Do not test inside Perplexity “incognito”

If the thread shows **Exit incognito**, you are in Perplexity’s private session mode, not a cold reader test.

1. Click **Exit incognito** (or open the thread in a normal logged-in tab).
2. **Share** → **Anyone with the link can view** → **Copy link**.
3. Verify in **Safari/Chrome private window** where you are **not** logged into Perplexity.

Pass: full answer visible without signing in. Fail: sign-in prompt or empty thread.

## Operator steps (each completed thread)

1. Open the thread in Perplexity (browser or app).
2. Use **Share** (or thread menu → share).
3. Under **General access**, select **Anyone with the link can view** (globe icon, not “Only people with access” or project-only).
4. Click **Copy link** and confirm the URL is `https://www.perplexity.ai/search/<uuid>`.
5. Paste that URL into `candidates.json` / Hugo YAML only if it matches the thread you shared.

Perplexity Computer project: run this on every thread before you paste the URL into the blog packet.

## Share manifest (required before `explore-apply`)

`explore-apply` (review-only or `APPLY=1`) runs a **fail-closed** preflight against `tmp/explore-proposals/<slug>/share-manifest.json`. Each approved proposal URL must have:

- `share_setting: anyone_with_link` and `share_confirmed_at` (after you set **Anyone with the link** in the UI)
- `cold_verified: true` and a fresh `cold_verified_at` (incognito pass, optionally plus HTTP probe)

```bash
go tool task explore-share-prepare CANDIDATES=tmp/explore-proposals/<slug>/candidates.json
# Share each thread in Perplexity UI (globe / anyone with link)
go tool task explore-share-confirm CANDIDATES=... CANDIDATE_ID=blame-blindspot
go tool task explore-share-verify-cold CANDIDATES=... CANDIDATE_ID=blame-blindspot PROBE=1 OPERATOR=1
go tool task explore-share-check PROPOSAL=tmp/explore-proposals/<slug>.proposal.yaml CANDIDATES=...
```

Cursor browser steps: [PERPLEXITY-SHARE-BROWSER.md](PERPLEXITY-SHARE-BROWSER.md). Scripts do **not** click Share for you; they record and gate on operator-confirmed public access.

## Verify (cold reader)

- Open the URL in a **private/incognito** window (or another browser where you are not logged into Perplexity).
- Pass: the research answer and citations load.
- Fail: login wall, blank thread, or “no access” → re-share as public, then update YAML if the URL changed.

After Hugo ship, list URLs on the post:

```bash
go tool task verify-explore-links POST=x-minds/2026-09-26-the-octopus-advantage
```

HTTP probe (`PROBE=1`) may false-fail on bot blocking; `OPERATOR=1` records your incognito pass and is required for a green `explore-share-check`.
