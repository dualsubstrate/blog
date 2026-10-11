#!/usr/bin/env bash
# Stop process-compose and free loopback ports used by make serve.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

ports=(3848 3849)

process-compose down -f process-compose.yaml 2>/dev/null || true

for port in "${ports[@]}"; do
	pids="$(lsof -nP -iTCP:"${port}" -sTCP:LISTEN -t 2>/dev/null || true)"
	if [[ -z "${pids}" ]]; then
		continue
	fi
	kill -TERM ${pids} 2>/dev/null || true
	sleep 0.3
	pids="$(lsof -nP -iTCP:"${port}" -sTCP:LISTEN -t 2>/dev/null || true)"
	if [[ -n "${pids}" ]]; then
		kill -KILL ${pids} 2>/dev/null || true
	fi
done

if [[ "${1:-}" == "--preflight" ]]; then
	exit 0
fi

if [[ -t 1 ]]; then
	echo "serve: stopped process-compose and freed ports ${ports[*]}"
fi
