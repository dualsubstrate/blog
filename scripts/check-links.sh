#!/usr/bin/env bash
# Build the Hugo site, serve ./public over HTTP, and crawl it with muffet (pinned in go.mod).
# Requires: hugo (extended), go >= 1.24 (for `go tool`), python3 (3.7+).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

PORT="${LINK_CHECK_PORT:-8765}"
BASE_URL="http://127.0.0.1:${PORT}/"

# Remove stale HTML so a prior `hugo server` or partial build cannot leave dev-only references in public/.
rm -rf public

# Keep data/tag-register.txt aligned with content (same as hugo task prerequisite).
go run ./cmd/tag-register -content content -deprecations data/tag-deprecations.toml -out data/tag-register.txt

# Absolute links must target this server; otherwise muffet would hit production baseURL from hugo.toml.
hugo --minify --gc --baseURL "$BASE_URL"

python3 "$ROOT/scripts/serve-public-threaded.py" "$PORT" >/dev/null 2>&1 &
SERVER_PID=$!

cleanup() {
  kill "$SERVER_PID" 2>/dev/null || true
}
trap cleanup EXIT

for _ in $(seq 1 30); do
  if curl -sf -o /dev/null "$BASE_URL"; then
    break
  fi
  sleep 1
done
if ! curl -sf -o /dev/null "$BASE_URL"; then
  echo "check-links: static server did not become ready at ${BASE_URL}" >&2
  exit 1
fi

# Only crawl the built site at 127.0.0.1 (avoids flaky paywalls, embeds, and oversized third-party headers).
# Add patterns here if the theme references assets not yet in static/.
go tool github.com/raviqqe/muffet/v2 \
  "$BASE_URL" \
  --include='^http://127\.0\.0\.1:' \
  --max-connections=32 \
  --timeout=30 \
  --buffer-size=65536 \
  --exclude='^mailto:' \
  --exclude='^https://(www\.)?facebook\.com' \
  --exclude='/apple-touch-icon\.png$' \
  --exclude='/safari-pinned-tab\.svg$' \
  --exclude='livereload'
