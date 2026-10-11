# substackhtml: Markdown to Substack-oriented HTML

This package turns repository Markdown (including YAML front matter) into a **small, allow-listed HTML fragment** meant for **paste into Substack’s rich-text editor**. It does **not** call Substack APIs or drive a browser.

Design goals: **predictable output**, **no raw Markdown in the clipboard**, and **no assumption that Substack accepts arbitrary HTML**. Treat the output as “probably safe paste”; always spot-check in a draft.

## Preserved formatting

After conversion and sanitization, these constructs are typically preserved (when present in the source):

- Headings `h1`–`h6` (optional demotion by one level; see `Options.DemoteHeadings`)
- Paragraphs implied by blank lines in Markdown (or `<br>` breaks if `Options.ParagraphMode = "br"`)
- **Bold** and *italic* (`strong` / `em`, and `b` / `i` if emitted)
- Links with `http`, `https`, or `mailto` URLs only
- Ordered and unordered lists, including nested lists
- Blockquotes
- Fenced code blocks (`pre` + `code`; language classes are stripped)
- Horizontal rules (`hr`)
- GitHub-flavored Markdown tables (see table strategy below)
- Strikethrough (`s` / `del`) when authored via GFM strikethrough syntax
- Task list checkboxes become plain `[ ]` / `[x]` text markers (GFM task lists)

## Explicit handling strategies

### Tables (`Options.TableMode`)

- `html` (default): keep a **minimal** `<table>` tree (`thead` / `tbody` / `tr` / `th` / `td`) with **no** `class`, `style`, or `id`. This is the best fidelity when the editor keeps tables.
- `list`: replace each table with an **ordered list** (`<ol><li>…</li></ol>`). Each row becomes one item; cells on a row are joined with ` — `. Use this when paste drops table markup.

### Paragraphs (`Options.ParagraphMode`)

- `p` (default): keep semantic paragraphs (`<p>`).
- `br`: flatten paragraphs into inline content separated by `<br><br>`. This can look tighter in Substack’s editor. When normalizing newline runs inside top-level text nodes, only newline characters are trimmed from each segment so a space before inline markup (for example `la <strong>…</strong>`) is never removed.

### Code blocks

Goldmark renders fenced code as `<pre><code>…</code></pre>`. Sanitization **strips** `class` and other attributes from `code` so paste is not tied to site-specific highlighter classes. Content is HTML-escaped on output.

### Images

- Allowed only when `src` is `http://` or `https://` after resolution.
- Output keeps **`src` and `alt` only** (no `width`, `loading`, etc.).
- Bundle-relative **`src`** values in the Markdown body (for example `![Diagram](diagram.webp)`) are joined with the same base URL as the featured image when **`Options.PagePermalink`** is set (see CLI below). If `src` stays non-HTTP(S) after that step, the `img` is dropped. If `alt` text exists, **alt is emitted as plain text** so you do not silently lose the description.
- Bundle-relative featured image paths are joined with **`Options.PagePermalink`** from `hugo list all` (production **`baseURL`**). Optional **`Options.ImageResolveOrigin`** swaps only scheme and host for local **`hugo server`** preview; leave empty so Substack and readers use the public site URL once the post is live. Set from **`substack.json`**, **`SUBSTACK_MARKDOWN_LEAD_IMAGE_RESOLVE_ORIGIN`**, or **`-markdown-lead-image-resolve-origin`** (see **`docs/substack-html/README.md`**). **`site_base_url_for_generated_links`** does not set image origin.

### Embeds and Hugo shortcodes

- `{{< youtube ID >}}` is rewritten **before** Markdown rendering to a short plain-text line containing the canonical `https://www.youtube.com/watch?v=…` URL (no iframe).
- `{{< mermaidfile >}}` or `{{< mermaidfile "name.mmd" >}}` is expanded **before** other shortcode removal when **`Options.SourcePath`** points at the Markdown file. If a sibling **`.webp`** exists with the same basename as the `.mmd` file (for example **`diagram.mmd`** next to **`diagram.webp`**, or **`diagram.es.mmd`** next to **`diagram.es.webp`**), the shortcode becomes a Markdown image `![Diagram](…webp)` so Substack gets a normal **`img`** instead of a mermaid code block. Otherwise the named `.mmd` file (default **`diagram.mmd`**) is read from the **same directory** and inlined as a fenced mermaid block (Goldmark emits **`pre`/`code`**). Same filenames as **`layouts/shortcodes/mermaidfile.html`** on the Hugo site.
- Other `{{< … >}}` / `{{% … %}}` shortcodes are **removed** (no placeholder inside the article body). Prefer adding a normal Markdown link in the source when an embed must survive export.
- Raw `<iframe>`, `<object>`, and `<embed>` are not kept as embeds. When a `src` URL is HTTP(S), the sanitizer emits a **simple text link** (`<a href="…">…</a>`) instead.

### Unsupported or site-specific HTML

Tags outside the allow-list are **unwrapped** (children are kept; the container is dropped). `script` / `style` nodes are removed entirely. HTML comments are removed during Markdown preprocessing.

### Boilerplate lines

Lines matching injected editor hints (for example a line starting with `The following cursor rule files are relevant`) are stripped from Markdown **before** render. This keeps accidental tooling noise out of exports.

## Video lead headings (English vs Spanish)

When `Options.IncludeFrontMatterLead` is true and `type: video`, the converter prepends the same lead structure as the site templates. Heading text is **English** by default. It switches to **Spanish** when either:

- `Options.SourcePath` basename ends with **`.es.md`** (Hugo sibling locale files), or
- front matter includes **`lang: es`**.

Set `SourcePath` from the real `-in` path so `go run ./cmd/substack-html` and `substack-draft` pick the right strings. The centered YouTube line uses **Watch on** vs **Ver en** accordingly.

The same locale rules are exposed as **`IsSpanishSiteLocale`**: `substack-draft` uses it for footer and auto-button URLs so Spanish drafts link under **`/es/`** on the live site (when `site_base_url` is the site root, not already ending in `/es`).

## Friends copy as Substack body

**`ResolveSubstackBody`** requires **`substack.md`** or **`substack.es.md`** beside the bundle (locale paired with **`index.md`** / **`index.es.md`**). There is no fallback to index body or **`facebook-*`** files. Missing or empty sidecars return an error. Authoring: **`.cursor/skills/site-substack-post/SKILL.md`**.

Site template leads (Claim, Teaser, TLDR, and so on) are not prepended when a substack sidecar is used; the sidecar is the full paste body.

## Cognitive-memetics "But why" paste block

When **`substack.json`** is loaded, **`cmd/substack-draft`** and **`cmd/substack-html`** call **`AppendCognitiveMemeticsProjectAboutHTML`** after **`Convert`** for posts under **`content/cognitive-memetics/`** that match the same rules as the Hugo partials (**Cube-Cows**, **Por-Estas-Calles** sayings, **Reptilocracy**, **Pawtropolis-Under-Fire**). The appended HTML is a compact **`<hr>` + `<blockquote>`** section (English or Spanish from the path and **`IsSpanishSiteLocale`**). Copy is defined in **`project_about_paste.go`** and must stay aligned with **`i18n/en.toml`** / **`i18n/es.toml`** for the `*ProjectAbout*` keys. Toggle with **`html_footer.include_cognitive_memetics_project_about`** (grouped JSON), **`html_footer_include_cognitive_memetics_project_about`** (flat), or **`SUBSTACK_HTML_FOOTER_INCLUDE_COGNITIVE_MEMETICS_PROJECT_ABOUT`**. Omit the key for default **on**.

## CLI

From the repository root:

```bash
go run ./cmd/substack-html -in path/to/post/index.md -out /tmp/post.html
go run ./cmd/substack-html -in path/to/post/index.md   # stdout
go run ./cmd/substack-html -in path/to/post/index.md -tables list
```

## Library

```go
opt := substackhtml.DefaultOptions()
opt.SourcePath = "/path/to/index.es.md"
html, err := substackhtml.Convert(sourceBytes, opt)
```

`Convert` strips YAML front matter, preprocesses shortcodes, renders with Goldmark (GFM), then normalizes HTML.
