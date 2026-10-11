#!/usr/bin/env bash
# Production-like Hugo build (tag-register runs as a Task dep).
set -euo pipefail
HUGO="${HUGO:-hugo}"
flags=(--minify --gc)
if [[ -n "${HUGO_PREVIEW:-}" ]]; then
  flags+=(--environment preview --buildDrafts --buildFuture --buildExpired)
fi
if [[ -n "${BASE_URL:-}" ]]; then
  "$HUGO" "${flags[@]}" --baseURL "$BASE_URL"
else
  "$HUGO" "${flags[@]}"
fi
while IFS= read -r -d '' f; do
  mv "$f" "$f.html"
done < <(find public -name 'carousel.preview' -type f -print0)
