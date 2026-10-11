#!/usr/bin/env python3
"""Draft LinkedIn carousel captions from carousel.json via local Gemma 4.

Writes numbered candidates next to the deck. Does not overwrite linkedin.txt.
"""

from __future__ import annotations

import argparse
import json
import os
import re
import sys
from pathlib import Path

_SKILLS = Path(__file__).resolve().parents[2]
if str(_SKILLS) not in sys.path:
    sys.path.insert(0, str(_SKILLS))
from site_local_eval_common.common import (  # noqa: E402
    DEFAULT_BASE,
    DEFAULT_MODEL,
    chat_complete,
    die,
)

SCRIPT_DIR = Path(__file__).resolve().parent
SKILL_DIR = SCRIPT_DIR.parent
PACK_PATH = SKILL_DIR / "packs" / "caption.md"

ACCENT_TAG = re.compile(r"</?accent[12][^>]*>", re.I)
VARIANT_SPLIT = re.compile(r"(?m)^(?:#{0,3}\s*)?(\d+)[.)]\s+")
EM_DASH = "\u2014"


def load_pack_system() -> str:
    if not PACK_PATH.is_file():
        die(f"Missing pack: {PACK_PATH}")
    text = PACK_PATH.read_text(encoding="utf-8")
    cut = text.find("## User message")
    return text[:cut].strip() if cut != -1 else text.strip()


def clean_slide_text(raw: str) -> str:
    text = ACCENT_TAG.sub("", raw)
    text = text.replace("**", "")
    text = text.replace("\\n", " ")
    text = text.replace("\n", " ")
    return " ".join(text.split()).strip()


def extract_beats(deck: dict) -> list[str]:
    beats: list[str] = []
    title = str(deck.get("title") or "").strip()
    if title:
        beats.append(f"Title: {title}")
    cta = deck.get("deck") if isinstance(deck.get("deck"), dict) else {}
    cta = cta.get("cta") if isinstance(cta, dict) else {}
    if isinstance(cta, dict):
        short_url = str(cta.get("shortUrl") or "").strip()
        post_url = str(cta.get("postUrl") or "").strip()
        if short_url:
            beats.append(f"Short URL: {short_url}")
        if post_url:
            beats.append(f"Post URL: {post_url}")
    slides = deck.get("slides")
    if not isinstance(slides, list):
        return beats
    for slide in slides:
        if not isinstance(slide, dict):
            continue
        number = slide.get("number", "?")
        role = str(slide.get("role") or "slide").strip()
        variants = slide.get("variants")
        if not isinstance(variants, list) or not variants:
            continue
        first = variants[0]
        if not isinstance(first, dict):
            continue
        blocks = first.get("blocks")
        if not isinstance(blocks, list):
            continue
        parts: list[str] = []
        for block in blocks:
            if not isinstance(block, dict):
                continue
            chunk = clean_slide_text(str(block.get("text") or ""))
            if chunk:
                parts.append(chunk)
        if parts:
            beats.append(f"Slide {number} ({role}): " + " / ".join(parts))
    return beats


def essay_url(deck: dict) -> str:
    cta = deck.get("deck") if isinstance(deck.get("deck"), dict) else {}
    cta = cta.get("cta") if isinstance(cta, dict) else {}
    if not isinstance(cta, dict):
        return ""
    short_url = str(cta.get("shortUrl") or "").strip()
    post_url = str(cta.get("postUrl") or "").strip()
    return short_url or post_url


def resolve_carousel_path(raw: str) -> Path:
    path = Path(raw).expanduser().resolve()
    if path.is_dir():
        candidate = path / "carousel.json"
        if candidate.is_file():
            return candidate
        die(f"No carousel.json in directory: {path}")
    if path.name == "carousel.json" and path.is_file():
        return path
    if path.is_file():
        sibling = path.parent / "carousel.json"
        if sibling.is_file():
            return sibling
    die(f"Could not find carousel.json from: {raw}")


def split_variants(raw: str) -> list[str]:
    matches = list(VARIANT_SPLIT.finditer(raw))
    if not matches:
        text = raw.strip()
        return [text] if text else []
    variants: list[str] = []
    for i, match in enumerate(matches):
        start = match.end()
        end = matches[i + 1].start() if i + 1 < len(matches) else len(raw)
        body = raw[start:end].strip()
        if body:
            variants.append(body)
    return variants


def ensure_essay_block(caption: str, url: str) -> str:
    text = caption.strip()
    if EM_DASH in text:
        text = text.replace(EM_DASH, ",")
    if url and "http" not in text:
        text = text.rstrip() + f"\n\nEssay:\n{url}"
    elif url and "Essay:" not in text:
        text = text.rstrip() + f"\n\nEssay:\n{url}"
    return text.strip() + "\n"


def format_candidates(variants: list[str], url: str) -> str:
    chunks: list[str] = []
    for i, variant in enumerate(variants, start=1):
        body = ensure_essay_block(variant, url)
        chunks.append(f"{i}. {body.strip()}\n")
    return "\n".join(chunks).rstrip() + "\n"


def pick_variant(candidates_text: str, index: int, url: str) -> str:
    variants = split_variants(candidates_text)
    if index < 1 or index > len(variants):
        die(f"Pick {index} is out of range (have {len(variants)} variants)")
    return ensure_essay_block(variants[index - 1], url)


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Draft LinkedIn carousel captions from carousel.json"
    )
    parser.add_argument(
        "path",
        help="carousel.json, bundle folder, or any file in the bundle",
    )
    parser.add_argument(
        "--dump-beats",
        action="store_true",
        help="Print extracted slide beats and exit (no Gemma)",
    )
    parser.add_argument(
        "--candidates",
        help="Path for numbered Gemma drafts (default: bundle/linkedin-carousel.candidates.txt)",
    )
    parser.add_argument(
        "--pick",
        type=int,
        metavar="N",
        help="Write variant N from --candidates (or a fresh Gemma run) to linkedin-carousel.txt",
    )
    parser.add_argument(
        "--out",
        help="Final caption path (default: bundle/linkedin-carousel.txt)",
    )
    parser.add_argument("--model", default=os.environ.get("LOCAL_LLM_MODEL", DEFAULT_MODEL))
    args = parser.parse_args()

    carousel_path = resolve_carousel_path(args.path)
    bundle = carousel_path.parent
    deck = json.loads(carousel_path.read_text(encoding="utf-8"))
    if not isinstance(deck, dict):
        die("carousel.json root must be an object")

    beats = extract_beats(deck)
    url = essay_url(deck)
    if args.dump_beats:
        print("\n".join(beats))
        if url:
            print(f"Essay URL: {url}")
        return

    candidates_path = (
        Path(args.candidates).expanduser().resolve()
        if args.candidates
        else bundle / "linkedin-carousel.candidates.txt"
    )
    out_path = (
        Path(args.out).expanduser().resolve()
        if args.out
        else bundle / "linkedin-carousel.txt"
    )

    if args.pick and candidates_path.is_file() and not os.environ.get("FORCE_GEMMA"):
        picked = pick_variant(candidates_path.read_text(encoding="utf-8"), args.pick, url)
        out_path.write_text(picked, encoding="utf-8")
        print(f"Wrote {out_path}")
        return

    if not beats:
        die(f"No slide text found in {carousel_path}")

    base = os.environ.get("LOCAL_LLM_BASE_URL", DEFAULT_BASE)
    api_key = os.environ.get("LOCAL_LLM_API_KEY")
    system = load_pack_system()
    user = (
        "CAROUSEL BEATS (source of truth):\n"
        + "\n".join(beats)
        + "\n\nEssay URL (must appear exactly):\n"
        + (url or "(none in deck.cta; omit Essay block)")
    )
    print(f"model={args.model} base={base}", file=sys.stderr)
    raw = chat_complete(base, args.model, api_key, system, user, max_tokens=1400)
    variants = split_variants(raw)
    if len(variants) < 2:
        die(f"Gemma did not return numbered captions:\n{raw[:800]}")

    formatted = format_candidates(variants, url)
    candidates_path.write_text(formatted, encoding="utf-8")
    print(formatted)
    print(f"Wrote {candidates_path}", file=sys.stderr)

    if args.pick:
        picked = pick_variant(formatted, args.pick, url)
        out_path.write_text(picked, encoding="utf-8")
        print(f"Wrote {out_path}", file=sys.stderr)


if __name__ == "__main__":
    main()
