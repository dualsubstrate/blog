# Substack HTML export (paste preview)

This folder holds **sample HTML** produced by the translation layer so you can inspect paste behavior before wiring browser automation.

- **Regenerate:** `go tool task sb-html` or `go tool task sb-html` (writes `sample-free-energy.html` and `sample-free-energy.es.html`).
- **CLI:** `go run ./cmd/substack-html -in <path-to-index.md> [-out file] [-tables html|list] [-paragraphs p|br] [-quotes blockquote|monospace]`.
- **Full behavior notes:** `internal/substackhtml/README.md`.

Markdown under `content/` remains the source of truth; the HTML is a compatibility-oriented export only.

## Browser paste test (draft UI)

To exercise the same HTML inside Chromium (local fake editor or a Substack tab you already opened), use `cmd/substack-draft` and see `internal/substackbrowser/README.md`. Quick local fixture (no Substack): `go run ./cmd/substack-draft -fixture -in content/mind-infrastructure/2026-04-29-free-energy-principle-hallucination-machine/index.md -keep-open 25s` (requires Chrome or Chromium on PATH).

### Make targets (`substack-draft`, `sb-en`, `sb-es`)

From the repo root (with direnv loaded or Chrome profile set):

- **English (default publication):** `go tool task substack-draft` or **`go tool task sb-en`** (same target)
  Uses root **`substack.json`** only. Optional: `SUBSTACK_URL=...` (omit if **`post_editor_url`** is set for paste and paste-schedule), `SUBSTACK_PUB=...`, `SUBSTACK_IN=...`. **`go tool task substack-draft`** defaults **`SUBSTACK_ACTION=paste-schedule`** (clicks **Continue**, then fills section, tags, delivery, schedule from the Markdown file). Use **`SUBSTACK_ACTION=paste`** for body paste only. **`CONFIRM_DISMISS=1`** (default) passes **`-confirm-dismiss`**: **`paste-schedule`** uses two prompts (after paste: **Continue** default vs **Close**; after automation: **Publish** default vs **Keep open**), or plain **c/q** then **p/k** on stderr in IDE terminals (empty line after the second prompt means **p**, click purple footer). **`paste`** uses one end prompt (**Continue** default to close). **`CONFIRM_DISMISS=0`** uses **`KEEP_OPEN`** (default `30s`) instead. After a successful **`paste-schedule`** run, a third optional prompt asks whether to append **`social-published`** for **`PUBLISH_TARGET`** (default **substack-en**; **`pick-draft-es`** defaults **substack-es**); default is **skip** (empty line or **No**); choose **Yes** only after the post is live. **Featured images:** by default, **`img src`** uses the Hugo permalink (production **`baseURL`**), which works once the post and bundle assets are deployed. Substack cover and social previews need a **URL their systems can reach**; **`localhost`** in pasted HTML usually shows broken. Set **`markdown_lead_image_resolve_origin`** (flat), **`markdown_export.lead_image_resolve_origin`** (grouped), or **`SUBSTACK_MARKDOWN_LEAD_IMAGE_RESOLVE_ORIGIN`** only when you deliberately point at a dev server on your machine. **`site_base_url_for_generated_links`** affects footer and browse links only, not image hosts.
  Bundle path: **`POST`** defaults to the sample **`mind-infrastructure/2026-04-29-free-energy-principle-hallucination-machine`** bundle (under `content/`). Optional `content/` prefix is stripped. Override with **`POST=section/slug`** or a **second Make goal** after the draft target (folder path, must include `/`), for example `go tool task sb-en human-condition/2026-05-01-ego-as-game` or `go tool task sb-en content/human-condition/2026-05-01-ego-as-game`. If the second goal ends with **`index.md`** or **`index.es.md`**, that file segment is stripped so the bundle folder is used. **`SUBSTACK_IN`** and **`SUBSTACK_IN_ES`** become `…/index.md` and `…/index.es.md` unless you set those from the command line or environment (full path override). **`POST=`** from the command line or environment still wins over the positional goal.

  **Body copy:** **`substack.md`** / **`substack.es.md`** is required (see **`.cursor/skills/site-substack-post/SKILL.md`**). No fallback. Front matter, schedule, tags, and featured image still come from the index file.

- **Spanish (overlay publication):** `go tool task sb-es`
  Runs with **`-config-global substack.json`** and **`-config substack.overlay.es.json`**, using the same **`POST`** default or override for **`index.es.md`**, unless **`SUBSTACK_IN_ES`** is set from the command line or environment.
  The matching English overlay is **`substack.overlay.en.json`**. Run `go tool task substack-draft SUBSTACK_LANG=en`.

### direnv (`.envrc`)

The repo root **`.envrc`** exports **`SUBSTACK_*`** variables (for example **`SUBSTACK_CHROMIUM_USER_DATA_DIRECTORY`**) so they override **`substack.json`** after JSON is loaded. Install [direnv](https://direnv.net/), then from this directory run **`direnv allow`** once per clone. Direnv keeps its cache under **`.direnv/`** (gitignored).

### `substack.json` and the Substack button

When `substack_insert_category_browse_button_after_paste` is true and `substack_category_browse_button_url` is empty, the draft command sets the button link from the Markdown front matter: first `categories` entry plus `site_base_url_for_generated_links`, as `https://…/categories/<slug>/`. Set `substack_category_browse_button_url` only when you want a fixed link for every post.

The default example uses **grouped objects** (`substack_browser`, `markdown_export`, …). A **flat** file with long keys at the root still works, as do **legacy** short keys (`pub`, `site_base_url`, …); see `legacyLocalConfigJSONKeys` and `rootJSONUsesGroupedSections` in `internal/substackbrowser/localconfig.go`.

**Global plus overlay:** point **`-config-global`** at the shared JSON (for example `substack-config.example.json`) and **`-config`** at a small overlay (see `substack-overlay.publish.example.json` for the Spanish publication: subdomain, writer home, new-post URL, browse button label). Overlay keys **win** on conflicts; nested objects are **deep-merged** (same top-level layout in both files, typically grouped). Example: `go run ./cmd/substack-draft -config-global docs/substack-html/substack-config.example.json -config docs/substack-html/substack-overlay.publish.example.json -in …/index.es.md`. **`paste-schedule` Substack section** comes from each post’s Hugo **`categories`**: the **second** entry when the post is **under cognitive memetics** (Markdown path under `content/cognitive-memetics/`, **or** first category is **`Cognitive-Memetics`**) and a non-empty second category exists; otherwise the **first** entry (so pillar + theme hubs such as `Human-Condition` + `Mental-Processes` still use the pillar). For **`index.es.md`** (or `lang: es`), that category id is mapped through **`i18n/es.toml`** (same display names as the Hugo site, e.g. `Social-Protocols` → `Protocolos-sociales`) so the Spanish Substack publication section picker can match. English drafts keep the raw Hugo category id. **Schedule send time** always comes from the post’s Hugo **`date`** (not optional). **`after_continue_section_label`** in JSON is **ignored** (kept only for backward-compatible parsing).

### `SUBSTACK_*` environment overrides

After JSON is loaded (including merge), `internal/substackbrowser` applies **non-empty** `SUBSTACK_*` variables. **Invalid booleans or integers are ignored** (no error). For `cmd/substack-draft`, flags such as **`-pub`** still override config **after** this step.

| Variable | Type | Maps to |
|----------|------|---------|
| `SUBSTACK_PUBLICATION_SUBDOMAIN` | string | publication subdomain |
| `SUBSTACK_POST_EDITOR_URL` | string | post editor URL |
| `SUBSTACK_CHROMIUM_USER_DATA_DIRECTORY` | string | Chrome user data dir |
| `SUBSTACK_PUBLISH_HOME_URL_SUFFIX` | string | publish home path suffix |
| `SUBSTACK_BROWSER_INITIAL_NAVIGATION_URL` | string | first navigation URL |
| `SUBSTACK_WRITER_HOME_URL_OVERRIDE` | string | writer dashboard override |
| `SUBSTACK_NEW_ARTICLE_DIRECT_URL` | string | new article URL |
| `SUBSTACK_SITE_BASE_URL_FOR_GENERATED_LINKS` | string | site base for footer and category-browse links (not used for featured image `src`) |
| `SUBSTACK_MARKDOWN_LEAD_IMAGE_RESOLVE_ORIGIN` | string | optional `http(s)://host:port` origin; overrides production permalink host for bundle-relative image `src` values (featured image, body images such as diagram `.webp`, and `{{< mermaidfile >}}` when a sibling `.webp` exists) |
| `SUBSTACK_CATEGORY_BROWSE_BUTTON_LABEL` | string | category browse button label |
| `SUBSTACK_CATEGORY_BROWSE_BUTTON_URL` | string | category browse button URL |
| `SUBSTACK_NEW_POST_CREATE_BUTTON_TEXT` | string | Create menu label |
| `SUBSTACK_NEW_POST_ARTICLE_MENU_ITEM_TEXT` | string | Article menu label |
| `SUBSTACK_MARKDOWN_TABLE_MODE` | string | table mode |
| `SUBSTACK_MARKDOWN_PARAGRAPH_MODE` | string | paragraph mode |
| `SUBSTACK_MARKDOWN_BLOCKQUOTE_MODE` | string | blockquote mode |
| `SUBSTACK_CHROMEDP_STEP_DELAY_MILLISECONDS` | int | step delay ms |
| `SUBSTACK_MARKDOWN_PARAGRAPH_LINE_BREAK_REPEAT_COUNT` | int | BR repeat count |
| `SUBSTACK_HTML_FOOTER_CATEGORY_LINK_LIST_INDEX` | int | category link list index |
| `SUBSTACK_SUBTITLE_MAX_CATEGORIES_IN_TYPE_LINE` | int | subtitle max categories |
| `SUBSTACK_SUBTITLE_MAX_LENGTH_CHARACTERS` | int | subtitle max length |
| `SUBSTACK_MARKDOWN_INCLUDE_FRONT_MATTER_LEAD_BLOCK` | bool | include lead block |
| `SUBSTACK_HTML_FOOTER_INCLUDE_ARTICLE_TAGS` | bool | footer tags |
| `SUBSTACK_HTML_FOOTER_INCLUDE_CATEGORY_BROWSE_LINK` | bool | footer category link |
| `SUBSTACK_HTML_FOOTER_INCLUDE_READ_ON_SITE_LINK` | bool | footer read-on-site |
| `SUBSTACK_HTML_FOOTER_INCLUDE_COGNITIVE_MEMETICS_PROJECT_ABOUT` | bool | append cognitive-memetics "But why" blocks (same copy as site partials); unset in JSON defaults to on |
| `SUBSTACK_SUBTITLE_USE_CATEGORY_TYPE_HIERARCHY_LINE` | bool | subtitle category line |
| `SUBSTACK_INSERT_CATEGORY_BROWSE_BUTTON_AFTER_PASTE` | bool | insert browse button |
| `SUBSTACK_PUBLISH_SCHEDULE_PUSH_DELIVERIES` | bool | schedule email and app delivery |
| `SUBSTACK_PUBLISH_SCHEDULE_MAX_ATTEMPTS` | int | max publish-settings automation passes (recovery between passes); 0 uses JSON or default 1 |
| `SUBSTACK_PUBLISH_SCHEDULE_DEBUG_DOM_ON_FAILURE` | bool | when true, `paste-schedule` writes a JSON DOM snapshot under `tmp/` if publish-settings automation fails (see `internal/substackbrowser/schedule_debug_dom.go`) |
| `SUBSTACK_SCHEDULE_DEBUG_DOM_FILE` | string | optional output path for that snapshot (default: `tmp/substack-schedule-debug-<timestamp>.json` in the process working directory) |
| `SUBSTACK_MARKDOWN_DEMOTE_HEADING_LEVELS_ONE_STEP` | bool | demote headings |

**Booleans:** `1`, `t`, `true`, `y`, `yes` for true; `0`, `f`, `false`, `n`, `no` for false (case-insensitive).

**Schedule debug snapshots** may contain draft body text from the open editor dialog. Keep them local; do not commit them (repo `tmp/` is gitignored).
