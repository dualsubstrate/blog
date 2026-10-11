#!/usr/bin/env bash
# Render a Mermaid file to PNG, WebP, or SVG using the official mermaid-cli Docker image (no local Node).
# WebP: mmdc only writes PNG; this script renders a temp PNG then runs cwebp (or magick) on the host.
# Usage: ./scripts/mermaid-docker.sh <input.mmd> <output.png|.webp|.svg|.pdf> [docker-image]
# Paths may be absolute or relative to the repository root (parent of scripts/).
# Styling: passes -C scripts/mermaid-site-export.css when present (mirrors site rules; mmdc injects
# that CSS inside the SVG, so the file scopes with #my-svg, not `.mermaid`). Override with MERMAID_CSS=.
# Raster background (-b) applies to .png and .webp; defaults to #2d2a32 (site panel); set MERMAID_BACKGROUND= to override.
# WebP: MERMAID_WEBP_QUALITY (default 80), or set MERMAID_WEBP_LOSSLESS=1 for cwebp -lossless.
# Bundle shorthand: make mermaid human-condition/... (EN+ES WebP) or make mermaid-render-en path
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
IN="${1:?first arg: path to .mmd file}"
OUT="${2:?second arg: path to output (.png, .webp, .svg, or .pdf)}"
IMAGE="${3:-minlag/mermaid-cli:11.12.0}"

case "$IN" in
/*) IN_ABS="$IN" ;;
*) IN_ABS="$ROOT/$IN" ;;
esac
case "$OUT" in
/*) OUT_ABS="$OUT" ;;
*) OUT_ABS="$ROOT/$OUT" ;;
esac

if ! command -v docker >/dev/null 2>&1; then
	echo "mermaid-docker: docker not found in PATH" >&2
	exit 127
fi

case "$IN_ABS" in
"$ROOT"/*) ;;
*)
	echo "mermaid-docker: input must be under repo root: $ROOT" >&2
	exit 2
	;;
esac
case "$OUT_ABS" in
"$ROOT"/*) ;;
*)
	echo "mermaid-docker: output must be under repo root: $ROOT" >&2
	exit 2
	;;
esac

IN_REL="${IN_ABS#"$ROOT"/}"
OUT_REL="${OUT_ABS#"$ROOT"/}"

WEBP_TMP_REL=""
MMD_OUT_REL="$OUT_REL"
case "$OUT_ABS" in
*.webp|*.WEBP)
	WEBP_TMP_REL="$(dirname "$OUT_REL")/.mermaid-docker-$$.png"
	MMD_OUT_REL="$WEBP_TMP_REL"
	;;
esac

MDC=()

CSS_REL="${MERMAID_CSS:-scripts/mermaid-site-export.css}"
CSS_ABS="$ROOT/$CSS_REL"
if [ -f "$CSS_ABS" ]; then
	MDC+=( -C "/data/$CSS_REL" )
	echo "mermaid-docker: css=/data/$CSS_REL" >&2
elif [ -n "${MERMAID_CSS:-}" ]; then
	echo "mermaid-docker: warning: MERMAID_CSS=$MERMAID_CSS not found under repo ($CSS_ABS), skipping -C" >&2
fi

if [ -n "${MERMAID_WIDTH:-}" ]; then MDC+=( -w "${MERMAID_WIDTH}" ); fi
if [ -n "${MERMAID_HEIGHT:-}" ]; then MDC+=( -H "${MERMAID_HEIGHT}" ); fi

case "$OUT_ABS" in
*.png|*.PNG|*.webp|*.WEBP)
	if [ -n "${MERMAID_BACKGROUND:-}" ]; then
		MDC+=( -b "${MERMAID_BACKGROUND}" )
		echo "mermaid-docker: raster background=$MERMAID_BACKGROUND" >&2
	else
		MDC+=( -b "#2d2a32" )
		echo "mermaid-docker: raster background=#2d2a32 (default; set MERMAID_BACKGROUND= to override)" >&2
	fi
	;;
esac

mkdir -p "$(dirname "$OUT_ABS")"

echo "mermaid-docker: image=$IMAGE" >&2
echo "mermaid-docker: in=/data/$IN_REL out=/data/$MMD_OUT_REL" >&2

MDC+=( -i "/data/$IN_REL" -o "/data/$MMD_OUT_REL" )

run_docker() {
	docker run --rm \
		-u "$(id -u):$(id -g)" \
		-v "$ROOT:/data" \
		-w /data \
		"$IMAGE" \
		"${MDC[@]}"
}

png_to_webp() {
	local png="$1"
	local webp="$2"
	if command -v cwebp >/dev/null 2>&1; then
		if [ -n "${MERMAID_WEBP_LOSSLESS:-}" ]; then
			echo "mermaid-docker: webp=cwebp -lossless" >&2
			cwebp -lossless "$png" -o "$webp"
		else
			local q="${MERMAID_WEBP_QUALITY:-80}"
			echo "mermaid-docker: webp=cwebp -q $q (set MERMAID_WEBP_QUALITY or MERMAID_WEBP_LOSSLESS=1)" >&2
			cwebp -quiet -q "$q" "$png" -o "$webp"
		fi
		return 0
	fi
	if command -v magick >/dev/null 2>&1; then
		if [ -n "${MERMAID_WEBP_LOSSLESS:-}" ]; then
			echo "mermaid-docker: webp=magick lossless" >&2
			magick "$png" -define webp:lossless=true "$webp"
		else
			local q="${MERMAID_WEBP_QUALITY:-80}"
			echo "mermaid-docker: webp=magick -quality $q" >&2
			magick "$png" -quality "$q" "$webp"
		fi
		return 0
	fi
	echo "mermaid-docker: WebP needs host 'cwebp' (e.g. brew install webp) or ImageMagick 'magick'" >&2
	return 1
}

if [ -n "$WEBP_TMP_REL" ]; then
	run_docker
	TMP_ABS="$ROOT/$WEBP_TMP_REL"
	if ! png_to_webp "$TMP_ABS" "$OUT_ABS"; then
		rm -f "$TMP_ABS"
		exit 1
	fi
	rm -f "$TMP_ABS"
else
	exec docker run --rm \
		-u "$(id -u):$(id -g)" \
		-v "$ROOT:/data" \
		-w /data \
		"$IMAGE" \
		"${MDC[@]}"
fi
