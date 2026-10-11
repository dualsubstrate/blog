#!/usr/bin/env bash
# Start process-compose for local dev; run serve-cleanup on exit (Task/go-task safe).
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

SOCK="$root/tmp/process-compose.sock"
mkdir -p "$root/tmp"

cleanup() {
	"$root/scripts/serve-cleanup.sh"
}

trap cleanup EXIT INT TERM
exec process-compose up -f process-compose.yaml -U -u "$SOCK"
