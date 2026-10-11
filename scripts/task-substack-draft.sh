#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

SUBSTACK_ACTION="${SUBSTACK_ACTION:-paste-schedule}"
PUBLISH_TARGET="${PUBLISH_TARGET:-substack-en}"
SUBSTACK_LANG="${SUBSTACK_LANG:-}"
SUBSTACK_OVERLAY_EN="${SUBSTACK_OVERLAY_EN:-substack.overlay.en.json}"
SUBSTACK_OVERLAY_ES="${SUBSTACK_OVERLAY_ES:-substack.overlay.es.json}"
WAIT_LOGIN="${WAIT_LOGIN:-0s}"
KEEP_OPEN="${KEEP_OPEN:-30s}"
CONFIRM_DISMISS="${CONFIRM_DISMISS:-1}"
SUBSTACK_DRAFT_FLAGS="${SUBSTACK_DRAFT_FLAGS:-}"
PASTE_TIMEOUT="${PASTE_TIMEOUT:-12m}"
CHROME_PROFILE="${CHROME_PROFILE:-$HOME/.cache/substack-chrome-profile}"

post="${POST:-mind-infrastructure/2026-04-29-free-energy-principle-hallucination-machine}"
if [[ $# -ge 1 && -n "${1:-}" ]]; then
  post="$(scripts/task-normalize-post.sh "$1")"
fi
post="$(scripts/task-normalize-post.sh "$post")"
SUBSTACK_IN="${SUBSTACK_IN:-content/$post/index.md}"
SUBSTACK_IN_ES="${SUBSTACK_IN_ES:-content/$post/index.es.md}"

merge=()
if [[ "$SUBSTACK_LANG" == "es" ]]; then
  merge=(-config-global substack.json -config "$SUBSTACK_OVERLAY_ES")
elif [[ "$SUBSTACK_LANG" == "en" ]]; then
  merge=(-config-global substack.json -config "$SUBSTACK_OVERLAY_EN")
fi

dismiss=()
if [[ "$CONFIRM_DISMISS" == "1" ]]; then
  dismiss=(-confirm-dismiss -keep-open 0)
else
  dismiss=(-keep-open "$KEEP_OPEN")
fi

# shellcheck disable=SC2206
extra_flags=($SUBSTACK_DRAFT_FLAGS)

run_draft() {
  local in_path="$1"
  if [[ -n "${SUBSTACK_URL:-}" ]]; then
    go run ./cmd/substack-draft "${merge[@]}" -action "$SUBSTACK_ACTION" -url "$SUBSTACK_URL" \
      -in "$in_path" -chrome-user-data-dir "$CHROME_PROFILE" -paste-timeout "$PASTE_TIMEOUT" \
      -publish-target "$PUBLISH_TARGET" "${extra_flags[@]}" "${dismiss[@]}"
  elif [[ -n "${SUBSTACK_PUB:-}" ]]; then
    go run ./cmd/substack-draft "${merge[@]}" -action "$SUBSTACK_ACTION" -pub "$SUBSTACK_PUB" \
      -in "$in_path" -chrome-user-data-dir "$CHROME_PROFILE" -wait-login "$WAIT_LOGIN" \
      -paste-timeout "$PASTE_TIMEOUT" -publish-target "$PUBLISH_TARGET" "${extra_flags[@]}" "${dismiss[@]}"
  else
    go run ./cmd/substack-draft "${merge[@]}" -action "$SUBSTACK_ACTION" -in "$in_path" \
      -chrome-user-data-dir "$CHROME_PROFILE" -wait-login "$WAIT_LOGIN" -paste-timeout "$PASTE_TIMEOUT" \
      -publish-target "$PUBLISH_TARGET" "${extra_flags[@]}" "${dismiss[@]}"
  fi
}

if [[ "$SUBSTACK_LANG" == "es" ]]; then
  run_draft "$SUBSTACK_IN_ES"
else
  run_draft "$SUBSTACK_IN"
fi
