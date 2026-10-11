#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SRC="${DSS_HOPS_CANVAS:-$HOME/.cursor/projects/Users-hector-Xynova-dss/canvases/DSS-core-topology.canvas.tsx}"
DEST="$ROOT/web/dss-hops/src/DSS-core-topology.canvas.tsx"

if [[ ! -f "$SRC" ]]; then
  echo "sync-dss-hops-canvas: missing source canvas: $SRC" >&2
  exit 1
fi

mkdir -p "$(dirname "$DEST")"
cp "$SRC" "$DEST"
echo "sync-dss-hops-canvas: copied to $DEST"
