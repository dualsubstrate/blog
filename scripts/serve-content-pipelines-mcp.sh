#!/usr/bin/env bash
# Run content-pipelines-mcp from providers/ (gitignored symlink). Used by process-compose.
# Stay in the site module so `go tool wgo` resolves; -cd switches into the provider for run.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
mcp="${root}/providers/content-pipelines-mcp"
cfg="${root}/content-pipelines-mcp.yaml"

if [[ ! -d "${mcp}" ]]; then
	echo "providers/content-pipelines-mcp missing; run ./scripts/link-providers.sh" >&2
	exit 1
fi
if [[ ! -f "${cfg}" ]]; then
	echo "content-pipelines-mcp.yaml missing" >&2
	exit 1
fi

cd "${root}"
exec go tool wgo -cd "${mcp}" -file .go -file .html -file .css go run ./cmd/content-pipelines-mcp --config "${cfg}" "$@"
