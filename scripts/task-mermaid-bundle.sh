#!/usr/bin/env bash
# Render EN and ES mermaid WebPs for a bundle (parallel when both exist).
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

MERMAID_STEM="${MERMAID_STEM:-diagram}"
MERMAID_DOCKER_IMAGE="${MERMAID_DOCKER_IMAGE:-minlag/mermaid-cli:11.12.0}"

post="${POST:-}"
if [[ $# -ge 1 && -n "${1:-}" ]]; then
  post="$(scripts/task-normalize-post.sh "$1")"
elif [[ -n "$post" ]]; then
  post="$(scripts/task-normalize-post.sh "$post")"
else
  echo "usage: task mermaid -- section/slug" >&2
  exit 2
fi

en_mmd="content/$post/${MERMAID_STEM}.mmd"
es_mmd="content/$post/${MERMAID_STEM}.es.mmd"
test -f "$en_mmd" || {
  printf '\n❌  missing %s\n\n   👉  Example:\n       go tool task mermaid -- human-condition/slug\n\n' "$en_mmd" >&2
  exit 1
}
test -f "$es_mmd" || {
  printf '\n❌  missing %s\n\n   Add Spanish Mermaid or use mermaid-render-en only.\n\n' "$es_mmd" >&2
  exit 1
}

render_one() {
  local lang="$1"
  local in out
  if [[ "$lang" == en ]]; then
    in="$en_mmd"
    out="content/$post/${MERMAID_STEM}.webp"
  else
    in="$es_mmd"
    out="content/$post/${MERMAID_STEM}.es.webp"
  fi
  MERMAID_WIDTH="${MERMAID_WIDTH:-}" MERMAID_HEIGHT="${MERMAID_HEIGHT:-}" \
    MERMAID_BACKGROUND="${MERMAID_BACKGROUND:-}" MERMAID_WEBP_QUALITY="${MERMAID_WEBP_QUALITY:-}" \
    MERMAID_WEBP_LOSSLESS="${MERMAID_WEBP_LOSSLESS:-}" \
    ./scripts/mermaid-docker.sh "$in" "$out" "$MERMAID_DOCKER_IMAGE"
}

render_one en &
pid1=$!
render_one es &
pid2=$!
wait "$pid1"
wait "$pid2"
