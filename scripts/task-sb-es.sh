#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

export SUBSTACK_LANG=es
export PUBLISH_TARGET=substack-es

if [[ $# -ge 1 && -n "${1:-}" ]]; then
  export POST
  POST="$(go run ./cmd/normalize-post "$1")"
  export SUBSTACK_IN_ES="content/$POST/index.es.md"
  exec scripts/task-substack-draft.sh
fi

if [[ -n "${POST:-}" ]]; then
  export SUBSTACK_IN_ES="content/$(go run ./cmd/normalize-post "$POST")/index.es.md"
  exec scripts/task-substack-draft.sh
fi

export SUBSTACK_ACTION=pick-draft-es
exec scripts/task-substack-draft.sh
