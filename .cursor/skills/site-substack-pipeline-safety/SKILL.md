---
name: site-substack-pipeline-safety
description: Keep the Substack Markdown to HTML and browser automation pipeline safe and reliable. Use when working on substackhtml, substackbrowser, cmd/substack-html, cmd/substack-draft, or Make targets related to Substack drafts. Prevent committing secrets or draft URLs, avoid Go toolchain version bumps from dependencies, and prefer fixture-based testing before touching Substack.
---

# Substack pipeline safety

## Guardrails (must follow)

- **Never commit secrets**: no tokens, cookies, Chrome profiles, `.env`, or any local auth state.
- **Never hardcode Substack URLs**: do not bake `SUBSTACK_URL` values into code, docs, Make targets, or tests. Treat draft links as sensitive.
- **Keep Markdown as source of truth**: edits happen in `content/`; HTML is derived output.
- **Substack body** requires **`substack.md`** / **`substack.es.md`** (locale paired with **`index.md`** / **`index.es.md`**). No fallback to index body or **`facebook-*`** sidecars. Front matter (date, tags, cover) always comes from the index file. Authoring: **`.cursor/skills/site-substack-post/SKILL.md`**.
- **Optimize for reliability**: prefer explicit allow-lists and predictable fallbacks over clever HTML tricks.

## Git hygiene

- **Schedule debug DOM snapshots**: when `SUBSTACK_PUBLISH_SCHEDULE_DEBUG_DOM_ON_FAILURE` or `substack_publish_schedule_debug_dom_on_failure` is on, failed `paste-schedule` runs write JSON under `tmp/` (or `SUBSTACK_SCHEDULE_DEBUG_DOM_FILE`). Those files can include draft HTML and dialog copy. Do not commit them; root `tmp/` is gitignored.
- **`.gitignore`**: ensure Chrome profile path copies (`substack-chrome-profile/`, `.cache/substack-chrome-profile/`) and legacy `.substack/` local state stay ignored. Root `substack.json` may be committed if it contains no secrets; prefer `SUBSTACK_CHROMIUM_USER_DATA_DIRECTORY` for machine-specific paths.
- **Review new files**: if any file looks like browser state (SQLite, cookies, profile folders), stop and add to `.gitignore`.
- **Layered config**: optional `-config-global` merges a shared JSON under `-config` (overlay wins; deep-merge objects). Prefer grouped layout in both files; shared templates may live under `docs/substack-html/`.
- **Env overrides**: `SUBSTACK_*` variables override JSON after merge (see `docs/substack-html/README.md`); invalid bool or int values are ignored. Do not put secrets in env in CI logs. Optional repo **`.envrc`** (direnv) exports `SUBSTACK_*`; run `direnv allow` after clone.

## Go toolchain and dependencies

- **Do not bump the module Go version** to satisfy dependencies unless the user explicitly requests it.
- When adding dependencies that might force a toolchain upgrade, prefer:
  - **Pinning an older compatible version**, or
  - **Replacing the dependency** with something lighter.
- When running module commands in this repo, prefer:
  - `GOTOOLCHAIN=local go mod tidy`
  - `GOTOOLCHAIN=local go test ./...`

## Browser automation testing workflow

1. **Start with the fixture**: run `go run ./cmd/substack-draft -fixture -in <path-to-index.md> -keep-open 25s` to validate paste and HTML shape without credentials.
2. **Then test Substack** with `make substack-draft` (optional `SUBSTACK_URL` when `post_editor_url` is empty in `substack.json`) or `SUBSTACK_PUB='...'`. The Makefile defaults **`SUBSTACK_ACTION=paste-schedule`** so the flow includes **Continue** plus section, tags, and schedule fields from front matter; override with **`SUBSTACK_ACTION=paste`** for paste-only. Default **`CONFIRM_DISMISS=1`** passes **`-confirm-dismiss`**: for **`paste-schedule`**, first prompt after paste (**Continue** default vs **Close**), then after automation a second prompt (**Publish** default vs **Keep open**); for **`paste`**, one prompt before close (**Continue** default). Plain stderr prompts use **c/q** and **p/k** when the terminal is not a full TTY (empty line after the second prompt acts like **p**). Set **`CONFIRM_DISMISS=0`** and **`KEEP_OPEN=30s`** for a timed wait instead.
3. **Publishing**: automation fills the publish modal (section, tags, delivery, schedule fields) but does **not** click Substack's final **Schedule** / **Publish** in that modal; you finish send in Substack. **Cancel before publish:** choosing **Close** after paste, **q** in plain prompts, **q** in the bundle picker, or **Ctrl+C** on any huh prompt exits immediately with **`substack-draft: aborted before publish`** (exit code 130). No publish-settings step, no **Publish** prompt, no **social-published** marker prompt. After a **completed** **`paste-schedule`** run only, **`substack-draft`** may ask whether to append **`social-published`** (default **No**). You can still run **`make sb-mark-published`** manually. **Featured images** use the Hugo permalink host unless **`markdown_lead_image_resolve_origin`** (flat), **`markdown_export.lead_image_resolve_origin`** (grouped), or **`SUBSTACK_MARKDOWN_LEAD_IMAGE_RESOLVE_ORIGIN`** is set. Substack cover and list thumbnails need a **public** **`img src`**; **`localhost`** usually breaks after paste. **`site_base_url_for_generated_links`** is for footer and browse links only, not image hosts. Section from Hugo `categories` uses the **second** entry when the post is under cognitive memetics (path under `content/cognitive-memetics/` **or** first category `Cognitive-Memetics` with a second category), else **first** entry. **`index.es.md`**: section label is mapped via **`i18n/es.toml`** (must match Substack ES section names). **Schedule time** uses Hugo **`date`**. Optional: `schedule_push_deliveries` (grouped) / `substack_publish_schedule_push_deliveries` (flat). Keys `fill_schedule_datetime_from_front_matter` / `substack_publish_schedule_fill_datetime_from_post` are **deprecated** (ignored). Config `after_continue_section_label` is ignored. Env override: `SUBSTACK_PUBLISH_SCHEDULE_PUSH_DELIVERIES`.

## What to preserve in HTML

Preserve only: headings, paragraphs, emphasis, links, lists, blockquotes, code blocks, and minimal tables (or table-to-list mode). Strip scripts, styles, custom widgets, and unsupported embeds.

