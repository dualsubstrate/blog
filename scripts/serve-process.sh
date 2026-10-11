#!/usr/bin/env bash
# Run a long-lived dev command under process-compose and tear down the whole
# process tree on exit (wgo/go run often leave listeners behind without this).
set -euo pipefail

cleanup() {
	trap - EXIT INT TERM
	kill -TERM 0 2>/dev/null || true
	sleep 0.5
	kill -KILL 0 2>/dev/null || true
}

(
	trap cleanup EXIT INT TERM
	"$@"
)
