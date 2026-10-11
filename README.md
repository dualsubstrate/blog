# Dual Substrate blog

Hugo/LoveIt host for [blog.dualsubstrate.com](https://blog.dualsubstrate.com/). Editorial sections: **Platform**, **Lineage**, **Proof**, **Notes**, plus static **About**.

## Prerequisites

- Hugo Extended **0.152.2** (CI pin)
- Go **1.27** (see `go.mod`)
- Sibling checkouts linked by `./scripts/link-providers.sh`:
  - `../../../ai/content-pipelines`
  - `../../../ai/content-pipelines-mcp`
- Git submodules: `themes/LoveIt`, `.cursor/packs/shared`

```bash
git submodule update --init --recursive
./scripts/link-providers.sh
git config core.hooksPath .githooks
```

## Build and preview

```bash
GOWORK=off go tool task calendar
GOWORK=off go tool task hugo
GOWORK=off go tool task serve   # Hugo + carousel-save + essays MCP
```

## DNS (GitHub Pages)

In the `dualsubstrate.com` zone, add:

- **Type:** CNAME  
- **Name:** `blog`  
- **Target:** `dualsubstrate.github.io`

Repo **Settings → Pages:** custom domain `blog.dualsubstrate.com`, enforce HTTPS after the certificate provisions. Source branch: **`gh-pages`** / root (deployed by `.github/workflows/hugo-pages.yml`).

## Explore further / Perplexity

Explore-further skills may still reference the Behaviour Engineering Perplexity project until a **Dual Substrate** project is created. Use a DSS project for new research threads when available.

## License

Site content: CC BY-SA 4.0 (see footer). Tooling follows each dependency's license.
