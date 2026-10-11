#!/usr/bin/env bash
# Sync Continuereadable copies from Cursor canonical trees.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

SKILLS=(
  site-claims-content site-video-content site-cognitive-memetics-content site-spanish-translation-content
  site-revise-post site-revise-flow site-revise-hooks site-revise-format site-revise-emphasis
  site-curiosity-title site-revise-score site-revise-post-es site-revise-spanish site-revise-prose
  site-linkedin-post site-facebook-post site-substack-post site-carousel-post
  site-tag-register site-tag-unify site-short-link-register site-ai-for-general-audience
)
RULES=(
  always-rules-0-ai.mdc
  always-rules-01-human-interaction.mdc
  site-always-rules-01-human-interaction.mdc
  site-content-placement.mdc
  site-content-markdown-writing.mdc
)

mkdir -p content-ai/skills content-ai/rules
for s in "${SKILLS[@]}"; do
  src=".cursor/skills/$s"
  dst="content-ai/skills/$s"
  [[ -d "$src" ]] || { echo "skip missing $src"; continue; }
  mkdir -p "$dst"
  rsync -a --delete \
    --include='SKILL.md' --include='reference.md' --include='examples.md' \
    --include='*/' --exclude='*' \
    "$src/" "$dst/"
done
for r in "${RULES[@]}"; do
  cp ".cursor/rules/$r" "content-ai/rules/$r"
done

# Continue-safe paths: .cursor -> content-ai inside copies only
find content-ai -type f \( -name '*.md' -o -name '*.mdc' \) -print0 | while IFS= read -r -d '' f; do
  perl -i -pe 's/`\.cursor\//`content-ai\//g; s{(?<![`/\w])\.cursor/(skills|rules)/}{content-ai/$1/}g' "$f"
done

echo "Synced content-ai/ from .cursor/ (Continue-readable copies)."
