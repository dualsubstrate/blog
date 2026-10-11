# Perplexity project for Explore further research

Extension research for **Explore further** MUST run inside the shared Perplexity **Project**, not from global Search, Computer-only threads, or MCP `perplexity_research` (which today opens a non-project thread).

## Canonical project

| Field | Value |
| --- | --- |
| Name | Dual Substrate blog (create a DSS Perplexity project; legacy BE project name was behaviorengineering blog) |
| URL | https://www.perplexity.ai/projects/72fa0438-7055-414e-8c7d-2c7642bb8ed5 |
| Context | Instructions + Skills (essay-extend, explore hooks, etc.) uploaded in Project **Settings** |

Paste sources for Instructions and skill zips: `tmp/essay-extension-skills/perplexity-computer/` (local; not in git).

## Blog Explore links (`perplexity_query`)

Hugo `reader_landing.explore` rows with `type: perplexity_query` link to `site.Params.perplexityProjectURL` in `hugo.toml` (same URL as this project). The full question is shown under the link for paste into **Start a Session**. MUST NOT use `?q=` on the project URL: that does not attach threads to the project Overview.

## Orphan threads (not on project Overview)

Threads started from global **New** / Search, `/search/new?q=`, MCP `perplexity_research`, or the top search bar without the project composer do **not** appear under the project session list. Recovery: **History** (library) → ⋯ on the thread → **Add it to a project** → **behaviorengineering blog**. Prefer re-running from the project ask box when project Instructions must apply from turn one.

Shipped completed research uses `type: perplexity_thread` with public share URLs.

## Operator workflow (each new research thread)

1. Open the **project URL** above (Overview tab).
2. Use the project **Ask anything about this project** field (or **New thread** inside this project). MUST NOT use the top-level **New** search outside the project.
3. Paste the Gemma **PRIMARY_PROMPT** (from `site-essay-extend-explore` / Polypus). One hinge per thread for a 2–3 thread batch.
4. When the answer finishes: **Share** → **Anyone with the link can view** (see [PERPLEXITY-SHARE.md](PERPLEXITY-SHARE.md)) → copy URL into `candidates.json` → record in `share-manifest.json` via `go tool task explore-share-prepare` / `explore-share-confirm` / `explore-share-verify-cold` (browser helper: [PERPLEXITY-SHARE-BROWSER.md](PERPLEXITY-SHARE-BROWSER.md)).
5. Confirm the thread appears under this project’s session list (same Overview), not only under generic Sessions.

Pass: thread listed on the project Overview with project Instructions/Skills active.  
Fail: thread exists only as a standalone `/search/<uuid>` you opened from global Search or MCP; re-run from the project ask box.

## Agent / MCP rules

| Tool | Explore batch |
| --- | --- |
| `perplexity_research` | **PROHIBITED** until MCP supports starting a thread inside `72fa0438-7055-414e-8c7d-2c7642bb8ed5` |
| `perplexity_export` | Allowed on real `thread_id` URLs the operator created **from the project** |
| `perplexity_continue` | Only if the active thread was started in the project |

CORRECT:

```text
Gemma prompt → operator pastes in project Ask → public share → URL in candidates.json → export + proposal
```

PROHIBITED:

```text
Agent calls perplexity_research for Explore further (no project skills)
```

## Project access vs blog readers

Project footer **Anyone in this project can view** is not enough for cold readers on behaviorengineering.ai. Each shipped thread still needs thread-level **Anyone with the link can view** per PERPLEXITY-SHARE.md.
