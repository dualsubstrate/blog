# Social autopost (GitHub Actions)

This repo can publish daily posts to:

- Facebook Page (Graph API)
- LinkedIn (Posts API, optional image upload)

The system is designed to be stateless in CI:

- **Facebook:** resolves **every** `content/` bundle whose front matter **`date`** matches the target day (same rules as `internal/contentbundle`), reads each **`linkedin.txt`**, then checks recent Page posts for that bundle’s site URL before posting. If the URL is already on the Page, that item is skipped. When the bundle has a **`featuredImage`** file next to **`index.md`**, the tool publishes a **Page photo** with **`linkedin.txt`** as the **caption** (multipart upload to Graph **`/{page-id}/photos`**). Otherwise it posts a **link** preview (`message` + **`link`** to the canonical post URL).
- **LinkedIn:** uses the **same bundle list** as Facebook (sorted paths). Optional **URL idempotency** calls LinkedIn’s “recent posts by author” API; that needs **`r_member_social`** or **`r_organization_social`**. Many member apps never get read scope, so CI often sets **`LINKEDIN_DISABLE_IDEMPOTENCY=1`** or passes **`-disable-idempotency`** (see `.github/workflows/social-autopost.yml`). Re-enable the scan when your token has the matching read scope.

**Dry-run** (`-dry-run` / `DRY_RUN=1` in `make`) does **not** call Facebook or LinkedIn; it only prints which bundles would be posted from the repo for that **`-date`**.

**Local interactive mode:** when stdin is a TTY and the process is not in CI (`CI` / `GITHUB_ACTIONS`), `facebook-autopost` and `linkedin-autopost` **prompt once per bundle**: **publish** (API + `social-published`), **tag-as-published** (`social-published` only), or **quit**. **`DRY_RUN=1`** (make default) shows the menu but writes **no** files and calls **no** API; **publish** only prints the preview payload. Use this when several posts share a date and some were already published manually. CI and non-TTY runs post every bundle that passes idempotency without prompting. Flags: **`-ask`** / **`SOCIAL_AUTOPOST_ASK=1`** force prompts; **`-no-ask`** / **`SOCIAL_AUTOPOST_NO_ASK=1`** disable them.

**`social-published` markers:** unless **`-dry-run`** or **`-no-mark-published`** / **`SOCIAL_AUTOPOST_NO_MARK=1`**, the tool updates `social-published` in the page bundle: target **`linkedin`** or **`facebook`** plus a UTC timestamp (same format as Substack: `go tool task sb-mark-published POST=section/slug PUBLISH_TARGET=linkedin`). Written when you choose **tag-as-published**, when **API idempotency** skips (URL already on the network), and after a **successful publish**. Bundles that already list the channel in `social-published` are skipped on later runs without prompting. Commit the file when you want the repo to remember state.

**Closed loop (local, same file as Substack):** **tag-as-published** or a successful **publish** appends **`linkedin`** or **`facebook`** in that bundle’s `social-published`; the next autopost run skips bundles that already have that line. To list gaps manually: `go tool task sb-list-unpublished PUBLISH_TARGET=linkedin` (or `facebook`). One bundle can hold several targets (`substack-en`, `linkedin`, `facebook`), each with its own timestamp.

## GitHub Actions workflows

- `.github/workflows/hugo-pages.yml`: deploy the site (scheduled)
- `.github/workflows/social-autopost.yml`: publish to Facebook + LinkedIn (scheduled)

Both workflows accept a manual date override through `workflow_dispatch`.

### Required GitHub secrets

Facebook:

- `FACEBOOK_PAGE_ID` (example: `1071261842742219`)
- `FACEBOOK_PAGE_ACCESS_TOKEN` (Page token)

LinkedIn:

- `LINKEDIN_AUTHOR_URN` (example: `urn:li:person:...` or `urn:li:organization:...`)
- `LINKEDIN_ACCESS_TOKEN` (OAuth access token)

Local **`go tool task linkedin-autopost`** defaults **`LINKEDIN_DISABLE_IDEMPOTENCY=1`** (no recent-post API scan) and **`LINKEDIN_NO_VERIFY_COMMENTARY=1`** (no GET-after-post; needs the same read scopes). Set either to **`0`** only if your token has **`r_member_social`** or **`r_organization_social`**. A **403** on verify after **201** publish means the post went live; verify was skipped.

### GitHub Environment for the workflow

The **Social autopost** workflow sets **`jobs.<job>.environment`** to a named GitHub Environment (see **`environment:`** in `.github/workflows/social-autopost.yml`, for example `social-media-production`). Create that name under **Settings → Environments**, then add **environment variables** there (for example **`LINKEDIN_DISABLE_IDEMPOTENCY`**, **`LINKEDIN_VERSION`**) so you can change behavior without editing YAML. Repository secrets (`LINKEDIN_*`, `FACEBOOK_*`) can stay at repo scope; they still resolve for this job.

If that environment uses **required reviewers** or a **wait timer**, scheduled runs pause until someone approves, which is usually wrong for a daily cron. Prefer no protection rules for this environment unless you intend manual approval.

## Facebook token minting

You need a **Page access token** (not an App token).

### Required permissions for this repo

For posting and idempotency checks, mint the token with:

- `pages_show_list`
- `pages_read_engagement`
- `pages_read_user_content`
- `pages_manage_posts`

### Get the Page ID

The Page ID is visible in Meta business tools, or you can query it if you know it already.

Example Page ID in this repo: `1071261842742219` (BehaviorEngineering.ai).

### Mint a Page access token (quick, may be short-lived)

1. Open Meta for Developers **Graph API Explorer**.
2. Select your app.
3. In **User or Page**, pick **User Token**.
4. Add permissions:
   - `pages_show_list`
   - `pages_read_engagement`
   - `pages_read_user_content`
   - `pages_manage_posts`
5. Click **Generate Access Token** and complete login/approval.
6. Fetch the Page token (either method):
   - `GET /me/accounts?fields=name,id,access_token`
     - Find the row with your Page ID and copy its `access_token`.
   - `GET /<PAGE_ID>?fields=access_token`
     - Copy the `access_token`.

Store it as the GitHub secret `FACEBOOK_PAGE_ACCESS_TOKEN`.

### Mint a long-lived Page access token (recommended for automation)

The reliable flow is:

1. Generate a short-lived **User token** (Graph API Explorer) with the 4 permissions above.
2. Exchange it for a long-lived **User token**:

   - `GET /oauth/access_token?grant_type=fb_exchange_token&client_id={APP_ID}&client_secret={APP_SECRET}&fb_exchange_token={SHORT_LIVED_USER_TOKEN}`

3. Use the long-lived User token to fetch a long-lived **Page token**:

   - `GET /me/accounts?fields=name,id,access_token&access_token={LONG_LIVED_USER_TOKEN}`

4. Store the resulting Page `access_token` as `FACEBOOK_PAGE_ACCESS_TOKEN`.

In Access Token Debugger, this Page token often shows:

- **Expires**: `Never`
- **Data Access Expires**: a date in ~90 days (you must re-mint before this)

### Check whether the token expires

In Meta for Developers, open **Tools -> Access Token Debugger**.

Paste the Page token and click **Debug**.

Two expiry fields may appear:

- **Expires**: when the token itself stops working
- **Data Access Expires**: when the app's granted access window ends (you must re-mint)

### Re-mint the token (when Data Access Expires is near)

Repeat the same steps:

1. Generate a fresh User token in Graph API Explorer with the three permissions above.
2. Fetch a fresh Page access token from `/me/accounts` or `/<PAGE_ID>?fields=access_token`.
3. Replace the GitHub secret value.

You do not need to publish the Meta app to do this for your own Page, as long as the token user is an admin of the Page and is a tester/admin/developer of the app.

## LinkedIn token minting (high level)

LinkedIn uses OAuth. You need:

- the author URN you want to post as (`LINKEDIN_AUTHOR_URN`)
- an OAuth access token with the right scope (`LINKEDIN_ACCESS_TOKEN`)

Posting scopes:

- Member posting: `w_member_social`
- Organization posting: `w_organization_social` (and your user must have the required org role)

Idempotency (the recent-post scan so retries do not duplicate) **also** needs a read scope on the same token:

- Member author (`urn:li:person:...`): `r_member_social`
- Organization author (`urn:li:organization:...`): `r_organization_social` (and a qualifying company page role)

If the token only has `w_*` scopes, `linkedin-autopost` fails with HTTP **403** and `ACCESS_DENIED` on `GET .../rest/posts?q=author` (LinkedIn names this `partnerApiPostsExternal.FINDER-author` in the error body). Re-authorize after adding the matching **read** scope in your LinkedIn app.

LinkedIn documents `r_member_social` as **restricted** (not every developer app gets it). If you cannot obtain it for a personal author, options are: post as an **organization** with `r_organization_social` + `w_organization_social`, or set `LINKEDIN_DISABLE_IDEMPOTENCY=1` / pass `-disable-idempotency` (posting still works; duplicates are possible if the job retries).

The **Linkedin-Version** header defaults to a pinned month in code (`internal/linkedinapi.DefaultLinkedInVersion`). Bump that constant when LinkedIn deprecates it, or override with env `LINKEDIN_VERSION` or `linkedin-autopost -linkedin-version` (YYYYMM). Using the current calendar month often returns **426 NONEXISTENT_VERSION** before LinkedIn activates that version.

## Local commands (from repo root)

Shared logic for both commands lives in **`internal/contentbundle`** (which bundles match a publish date) and **`internal/socialbundle`** (load `linkedin.txt`, canonical URL, featured image path, and identical dry-run print layout).

Facebook:

```bash
go run ./cmd/facebook-autopost -date "YYYY-MM-DD" -dry-run
```

Notes:

- Resolves every matching bundle under `content/<section>/<slug>/` and `content/<section>/<hub>/<slug>/` whose front matter **`date`** matches `YYYY-MM-DD` (sorted by path). Each bundle must have **`linkedin.txt`** with a **`blog.dualsubstrate.com`** URL.
- **Empty day:** if no bundle matches the date, both tools **warn and exit 0** (scheduled CI stays green on quiet calendar days).
- If more than one bundle shares that date, it posts **each** in order. Use **`-post section/slug`** (or **`SOCIAL_POST=...`** with `go tool task facebook-autopost`) to target a single bundle.
- It extracts the canonical site URL from `linkedin.txt` for idempotency and for the Facebook link post.
- **Image:** when **`featuredImage`** resolves to a file in the bundle, Facebook autopost uploads that file as a **published Page photo** and sets **`linkedin.txt`** as the photo **caption** (so the post card shows your art). If Graph rejects a format (rare for **`.webp`**), fix the raster or temporarily clear **`featuredImage`** to fall back to a link-only post.
- **Caption limits:** link-only / text-only posts must be at most **3000 characters** (`facebook-autopost` fails before calling the API). Image captions are not hard-capped in bytes (LinkedIn UI may accept more than the API stored in some cases).
- **Retries:** Graph idempotency reads and publish calls retry transient failures (HTTP 5xx, 429, timeouts) up to **3** times by default. Override with **`-http-retries`** or env **`FACEBOOK_HTTP_RETRIES`**. After a transient publish error, the tool re-scans the Page feed for the site URL before posting again (avoids duplicates when Meta accepts the post but returns an error). Failed bundles are skipped and the run continues; exit code is non-zero if any bundle failed. Invalid **`linkedin.txt`** or broken **`social-published`** files still stop the whole run.
- **CI:** GitHub **Social autopost** runs Facebook first, then LinkedIn even when Facebook exits non-zero (LinkedIn step uses `if: always()`).

LinkedIn:

```bash
go run ./cmd/linkedin-autopost -date "YYYY-MM-DD" -dry-run
```

Notes:

- Same bundle selection as **facebook-autopost** (see `internal/contentbundle`). Optional **`-post`** / env **`SOCIAL_POST`** (same as **`go tool task linkedin-autopost`** with **`SOCIAL_POST=section/slug`**).
- If more than one bundle shares that date, the command posts **each** in order unless **`-post`** narrows to one bundle.
- **LinkedIn guardrails (autopost):**
  - **Video bundles (`youtube_id`, no local `featuredImage`):** `linkedin-autopost` builds a **link card** via `content.article`: downloads `https://img.youtube.com/vi/{id}/hqdefault.jpg`, uploads it with the Images API, sets `source` to the YouTube watch URL, and uses `title` / `subtitle` (or first `description` line) from `index.md`. No local thumb file required. Local `featuredImage` still wins (image + caption post).
  - **Little text encoding (default):** before `POST /rest/posts`, `linkedin-autopost` encodes hashtags as `{hashtag|\#|Tag}` and escapes reserved characters per [little text format](https://learn.microsoft.com/en-us/linkedin/marketing/community-management/shares/little-text-format). Plain `#` in API commentary can break parsing and drop trailing URL blocks. Disable with **`-skip-little-text-encode`** (not recommended).
  - **Commentary logging:** logs UTF-8 byte and rune counts for `linkedin.txt` and the encoded payload.
  - **Soft warning:** logs a warning when an image post exceeds **800 UTF-8 bytes** (observed threshold where link blocks were lost); does not block the post.
  - **Post-verify (default):** after publish, `GET /rest/posts/{urn}` and **`VerifyCommentary`** fails when stored `commentary` is missing site URLs, footer markers (`🧷` / `🔗` / `Full post`), hashtag tag names from `linkedin.txt`, bilingual `- EN:` / `- ES:` lines when present, or is much shorter than sent (truncation after `#` or `(`). Also flags unclosed `{hashtag|...}` templates. Requires read scope (`r_member_social` / `r_organization_social`). Skipped with a log line if the token lacks read scope; disable with **`-no-verify-commentary`**.
- **Text-only cap:** link-only Facebook posts and text-only LinkedIn posts (no `featuredImage`) must be at most **3000 characters** (hard fail before API).

### Post text looks truncated (for example `TS5: Sm(` or missing `🧷` / `🔗` links)

**`linkedin-autopost`** and **`facebook-autopost`** read **`linkedin.txt`** with Go from disk. They do **not** pipe the body through **`sh`**. Copy such as **`TS5: Sm(art)`** is safe for **`go tool task linkedin-autopost`** in this repository.

If links are missing on LinkedIn but paste in the UI works:

1. Compare **dry-run** output (`----- BEGIN POST TEXT -----`) with the live post.
2. Re-run with verify enabled and a token that has **read** scope; a failed verify means LinkedIn stored a shorter `commentary` than `linkedin.txt`.
3. Ensure little text encoding is on (default); bare `#hashtags` before URL lines are a common API parse failure.

If you see a cut right after **`Sm(`** when testing outside this command:

1. **GNU Make:** In a **`Makefile`**, the substring **`$(art)`** expands the Make variable **`art`**. Never inline post bodies into Makefiles.
2. **Shell:** Unquoted **`echo ... Sm(art)`** treats **`(`** as a subshell. Use **`printf '%s\n' '...'`** or a here-doc.
3. **Sanity check:** **`DRY_RUN=1 SOCIAL_POST=section/slug go tool task linkedin-autopost DATE=YYYY-MM-DD`** prints the exact API payload (little text encoded when enabled).
