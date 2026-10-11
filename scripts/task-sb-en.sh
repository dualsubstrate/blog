#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

if [[ $# -ge 1 && -n "${1:-}" ]]; then
  export POST
  POST="$(scripts/task-normalize-post.sh "$1")"
  exec scripts/task-substack-draft.sh
fi

if [[ -n "${POST:-}" ]]; then
  exec scripts/task-substack-draft.sh
fi

export SUBSTACK_ACTION=pick-draft-en
exec scripts/task-substack-draft.sh
