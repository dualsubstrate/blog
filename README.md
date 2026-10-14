# Dual Substrate System (site)

Hugo/LoveIt host for the **main Dual Substrate marketing site** (hero, surfaces, advantages, proof) plus editorial sections **Platform**, **Lineage**, **Proof**, and **Notes**, and static **About**.

Until custom DNS is configured, the public URL is [https://dualsubstrate.github.io/blog/](https://dualsubstrate.github.io/blog/). When you point **`www.dualsubstrate.com`** (or apex) at GitHub Pages for this repo, set `baseURL` in `hugo.toml` and add `static/CNAME`; do not use a separate `blog.` subdomain for this site.

## Prerequisites

- Hugo Extended **0.152.2** (CI pin)
- Go **1.27** (see `go.mod`)
- Sibling checkouts linked by `go run ./cmd/link-providers`:
  - `../../../ai/content-pipelines`
  - `../../../ai/content-pipelines-mcp`
- Git submodules: `themes/LoveIt`, `.cursor/packs/shared`

```bash
git submodule update --init --recursive
GOWORK=off go run ./cmd/link-providers
git config core.hooksPath .githooks
```

## Build and preview

```bash
GOWORK=off go tool task calendar
GOWORK=off go tool task hugo
GOWORK=off go tool task serve   # process-compose TUI; supervisor on tmp/process-compose.sock (not TCP :8080)
```

Local site URL when using `serve`: **http://127.0.0.1:28492/** (Hugo default :1313 is used by the Behaviour Engineering site). Carousel-save and essays MCP still use **:3848** and **:3849**; run only one full `serve` stack at a time if both repos need MCP.

Hugo only (no MCP/carousel):

```bash
hugo server --buildDrafts --buildFuture --bind 127.0.0.1 --port 28492
```

## GitHub Pages

**GitHub Pages (default):** [https://dualsubstrate.github.io/blog/](https://dualsubstrate.github.io/blog/). Source branch **`gh-pages`** / root (`.github/workflows/hugo-pages.yml`).

**Custom domain (main site):** In the `dualsubstrate.com` zone, point **`www`** (recommended) or apex at GitHub Pages, add `static/CNAME` with that hostname, set `baseURL` in `hugo.toml` to match, and enable the custom domain under repo **Settings → Pages**.

## Explore further / Perplexity

Explore-further skills may still reference the Behaviour Engineering Perplexity project until a **Dual Substrate** project is created. Use a DSS project for new research threads when available.

## License

Site content: CC BY-SA 4.0 (see footer). Tooling follows each dependency's license.
