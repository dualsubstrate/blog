#!/usr/bin/env bash
# Normalize POST (section/slug or content/.../index.md) to section/slug without content/ prefix.
set -euo pipefail
post="${1:-}"
post="${post#content/}"
post="${post%/}"
post="${post%/index.md}"
post="${post%/index.es.md}"
printf '%s' "$post"
