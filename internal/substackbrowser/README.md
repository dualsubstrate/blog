# substackbrowser: Chromium paste helper for Substack drafts

This package drives **Chromium via chromedp**, navigates to a URL, waits for an editor surface, then runs `document.execCommand("insertHTML", …)` so you can validate the HTML export in a real rich-text editor.

Nothing here clicks **Publish**. Substack may autosave drafts when the editor accepts changes.

## How to test without Substack (recommended first)

Use the **fixture** mode. It opens a `data:` URL with a fake `div.ProseMirror` editor so you can confirm paste behavior with zero credentials.

From the repository root (Chrome or Chromium must be installed):

```bash
go run ./cmd/substack-draft -fixture -in content/mind-infrastructure/2026-04-29-free-energy-principle-hallucination-machine/index.md -keep-open 20s
```

Flags:

- `-headless` if you truly do not need a visible window (harder to debug).
- `-keep-open 20s` keeps the tab open after a successful paste so you can inspect formatting.

## How to test against Substack

Substack’s **post editor URL varies** and the writer dashboard (`…/publish/home`) often **does not** contain the body editor until you open a draft. Reliable flows:

1. **Manual URL (most reliable):** In your browser, start a new Text post until the editor loads. Copy the address bar URL, then run:

   ```bash
   go run ./cmd/substack-draft -url 'PASTE_EDITOR_URL_HERE' -in path/to/index.md -chrome-user-data-dir "$HOME/.cache/substack-chrome-profile" -paste-timeout 10m
   ```

2. **Dashboard + wait:** Open the default home URL, then **during** `-wait-login` click **New post** so the editor appears before `-paste-timeout` elapses:

   ```bash
   go run ./cmd/substack-draft -pub YOURSUBDOMAIN -in path/to/index.md \
     -chrome-user-data-dir "$HOME/.cache/substack-chrome-profile" \
     -wait-login 3m -paste-timeout 12m -keep-open 30s
   ```

`-chrome-user-data-dir` is optional but strongly recommended so you log in once and reuse the profile.

### Local config (optional)

Discovery order: **`SUBSTACK_CONFIG`**, then **`~/.config/behaviour-engineering/substack.json`** (`go tool task sb-config-init`), then repo **`substack.json`** (often committed; avoid secrets), or pass **`-config`**. Optional root **`.envrc`** (direnv) exports **`SUBSTACK_*`** overrides; see `docs/substack-html/README.md`. If the resolved file is missing, the loader tries **`substack.config`**, then **`.substack/config.json`**, then **`.substack/substack.json`** for migration.

Prefer the grouped layout in `docs/substack-html/substack-config.example.json` (sections `substack_browser`, `markdown_export`, …). Flat root keys and legacy keys (`pub`, …) still load; see `internal/substackbrowser/localconfig.go`.

**`-config-global`:** `cmd/substack-draft` and `cmd/substack-html` merge a shared base JSON under `-config` (overlay wins; nested objects merge). Spanish-only deltas live in `docs/substack-html/substack-overlay.publish.example.json` (merge with `substack-config.example.json`).

**`SUBSTACK_*` env vars:** after JSON loads, optional environment variables override the same fields (see the table in `docs/substack-html/README.md`). `cmd/substack-draft` flags still win over config where the command sets them after load.

`cmd/substack-draft` reads it automatically unless you override with flags. From the repo Makefile: **`go tool task substack-draft`** or **`go tool task sb-en`** (English, root `substack.json` only) and **`go tool task sb-es`** (Spanish overlay merge + default `index.es.md`); see `docs/substack-html/README.md`.

When `substack_editor.insert_category_browse_button_after_paste` is true and `category_browse_button_url` is empty, the draft command fills the button URL from the post’s first Hugo `categories` value and `site.canonical_base_url` (same as `…/categories/<slug>/` under that base).

## Dry run (no browser)

```bash
go run ./cmd/substack-draft -dry-run -fixture -in path/to/index.md
```

## Run.Config extras

- **`InsertSubscribeButton`:** when true (default in **`cmd/substack-draft`** for non-fixture runs), **`Run`** evaluates **`InsertSubscribeButton`** after body paste and optional category-browse button: opens the editor **Button** menu and selects **Subscribe** (native Substack block). Disable via **`substack_insert_subscribe_button_after_paste`** in **`substack.json`** or **`SUBSTACK_INSERT_SUBSCRIBE_BUTTON_AFTER_PASTE=0`**.
- **`ScheduleEnabled` runs Continue plus publish settings (section, tags, delivery, schedule time).** Tags use **CDP mouse clicks** after the schedule script (Substack Headless UI combobox ignores programmatic `el.click`). Prefer an existing listbox match; if Substack only offers **Create "…"**, wait briefly then click Create so new Hugo tags become chips. When the tag combobox goes **readonly** (chip row full), remaining Hugo tags are skipped with a log line so section and schedule still run. The schedule JS still accepts `tags` for fixture/legacy paths; `Run` clears them for the JS pass and applies Hugo tags via `AddPublishTagsCDP`.
- **`ScheduleAfterPasteConfirm`:** optional callback invoked after a successful paste when **`ScheduleEnabled`** is true, before any publish-settings automation. If it returns **`proceed` false**, **`Run`** returns without clicking Substack Continue or filling section/tags/schedule.
- **`ConfirmDismiss`:** optional **`func(ctx context.Context) error`** invoked after a successful paste and optional schedule step, before the browser session ends. The context is the same chromedp tab as paste/schedule. When non-nil, **`KeepOpen`** timed sleep is skipped so the caller can block on a TUI (see **`cmd/substack-draft`** **`-confirm-dismiss`**: **`paste`** uses one end prompt; **`paste-schedule`** uses **`ScheduleAfterPasteConfirm`** after paste, then **`ConfirmDismiss`** after schedule; accepting the default **Publish** (or choosing **Publish**) in the TUI runs **`ClickPublishModalFooter`** so Substack’s purple footer button is clicked before Chrome exits).
- **`ScheduleDebugDOM` / `ScheduleDebugDOMPath`:** when **`ScheduleDebugDOM`** is true and the schedule step fails (or chromedp cannot evaluate the schedule script), **`Run`** writes a JSON file with **`scheduleFailureReason`**, clipped publish dialog HTML, **`cmdk-item`** rows, and radix portal samples. Enable via **`substack_publish_schedule_debug_dom_on_failure`** in **`substack.json`**, or env **`SUBSTACK_PUBLISH_SCHEDULE_DEBUG_DOM_ON_FAILURE=1`**. Optional path: **`substack_schedule_debug_dom_file`** or **`SUBSTACK_SCHEDULE_DEBUG_DOM_FILE`**; default **`tmp/substack-schedule-debug-<timestamp>.json`** under the current working directory. Treat snapshots as sensitive (draft content); do not commit them.
- **`ScheduleMaxAttempts`:** when **`ScheduleEnabled`**, **`Run`** logs a **publish flow stage** line before each pass (editor ProseMirror, **Continue**, publish settings text, dialog count). On failure with attempts remaining, **`recoverPublishFlow`** restores the editor: **Navigate** to **`TargetURL`** when it looks like a stable **`/publish/post/<id>`** draft URL, otherwise **Reload** first when **`TargetURL`** is a generic **`/publish/post?type=…`** composer (Navigating to the latter can spawn a **new** blank draft). Then **`RecoverPublishFlowScheduleRetrySyncJS`** runs (Escape, small **OK** dialogs only; no header **Publish** / **Next** clicks, so recovery does not drive the publish flow). This addresses CDP **-32000** when Substack navigates the tab during the long async schedule script. For reliable retries, set **`substack_post_editor_url`** (grouped **`substack_browser.url`**, or pass **`-url`**) to the **current draft** URL from the address bar, not only **`substack_new_article_direct_url`**. Set attempts via **`substack_publish_schedule_max_attempts`**, grouped **`substack_publish.schedule_max_attempts`**, env **`SUBSTACK_PUBLISH_SCHEDULE_MAX_ATTEMPTS`**, or **`cmd/substack-draft -schedule-max-attempts=N`**. Values below **1** are treated as **1** inside **`Run`**; above **6** cap at **6**.
- **Duplicate tag `alert`:** Substack may call **`window.alert`** with copy such as **"Tag already set"** or **"Tag already exists"** if automation tries to create a tag that is already on the publication. **`Run`** enables the CDP Page domain and listens for those messages (see **`substackTagConflictAlert`**), then accepts the dialog in a background **`chromedp.Run`** so the session does not hang. Other alert text is left for the user.

## Limitations

- **Schedule exit code vs Substack state:** if the schedule script hits CDP **-32000** (tab navigated or execution context destroyed) **after** Substack has already accepted **Schedule** / **Publish** in the modal, the run can still exit **non-zero** while the post is on its way. A later retry may then fail on tags or other steps even though the send already happened. Check Substack **Posts** when the log shows **-32000** or the hint line about Chrome losing the tab mid-run.
- `insertHTML` is **not a public Substack API**. If Substack or Chromium changes behavior, paste may fail; the command exits non-zero with the injected script’s reason string.
- Title and subtitle setting is **best-effort** and depends on Substack’s DOM. If fields cannot be found, the command still pastes the body HTML. Substack often reuses the subtitle as the default **social preview** description; when **`subtitle_include_categories`** is on and the post has categories, `substack-draft` fills the subtitle with the short **`Type: …`** hierarchy line first, then falls back to **`sowhat`** or the first non-empty **`description`** line if the subtitle is still empty.
- Cover image selection is **out of scope** for this package. **`ScheduleAfterContinueJS`** fills the publish modal only; final send is manual in Substack or via **`ClickPublishModalFooter`** from **`cmd/substack-draft`** when you confirm in the TUI. Substack UI changes can still break automation.
