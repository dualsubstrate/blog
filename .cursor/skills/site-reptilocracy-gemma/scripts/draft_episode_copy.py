#!/usr/bin/env python3
"""Draft Reptilocracy description/tldr/fluff via local Gemma 4 with thinking mode.

Prints labeled candidates. Does not edit Hugo index.md.
"""

from __future__ import annotations

import argparse
import json
import os
import re
import sys
import urllib.error
import urllib.request
from pathlib import Path

_SKILLS = Path(__file__).resolve().parents[2]
if str(_SKILLS) not in sys.path:
    sys.path.insert(0, str(_SKILLS))
from site_local_eval_common.common import (  # noqa: E402
    DEFAULT_BASE,
    DEFAULT_MODEL,
    die,
)

SCRIPT_DIR = Path(__file__).resolve().parent
SKILL_DIR = SCRIPT_DIR.parent
PACK_PATH = SKILL_DIR / "packs" / "episode_copy.md"

LABELS = (
    "DESCRIPTION_1",
    "DESCRIPTION_2",
    "DESCRIPTION_3",
    "TLDR",
    "FLUFF",
)
EM_DASH = "\u2014"


def load_pack_system() -> str:
    if not PACK_PATH.is_file():
        die(f"Missing pack: {PACK_PATH}")
    text = PACK_PATH.read_text(encoding="utf-8")
    cut = text.find("## User message")
    return text[:cut].strip() if cut != -1 else text.strip()


def chat_complete_thinking(
    base_url: str,
    model: str,
    api_key: str | None,
    system: str,
    user: str,
    *,
    max_tokens: int = 4096,
    temperature: float = 0.7,
) -> tuple[str, str]:
    """Return (content, reasoning). Thinking enabled for Gemma 4 chat template."""
    url = base_url.rstrip("/") + "/chat/completions"
    payload: dict = {
        "model": model,
        "messages": [
            {"role": "system", "content": system},
            {"role": "user", "content": user},
        ],
        "temperature": temperature,
        "max_tokens": max_tokens,
        "chat_template_kwargs": {"enable_thinking": True},
    }
    data = json.dumps(payload).encode("utf-8")
    headers = {
        "Content-Type": "application/json",
        "Accept": "application/json",
    }
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"
    req = urllib.request.Request(url, data=data, headers=headers, method="POST")
    try:
        with urllib.request.urlopen(req, timeout=300) as resp:
            raw = resp.read().decode("utf-8")
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        die(f"Gateway HTTP {e.code} for model {model!r}:\n{body}")
    except urllib.error.URLError as e:
        die(f"Gateway unreachable at {url}: {e.reason}")

    try:
        obj = json.loads(raw)
    except json.JSONDecodeError:
        die(f"Non-JSON response from gateway:\n{raw[:800]}")

    if "error" in obj:
        err = obj["error"]
        if isinstance(err, dict):
            die(f"Gateway error: {err.get('message', err)}")
        die(f"Gateway error: {err}")

    choices = obj.get("choices")
    if not choices:
        die(f"No choices in response:\n{raw[:800]}")

    msg = choices[0].get("message") or {}
    content = (msg.get("content") or "").strip()
    reasoning = (msg.get("reasoning") or msg.get("reasoning_content") or "").strip()
    if not content and not reasoning:
        die(f"Empty message content:\n{raw[:800]}")
    return content, reasoning


def extract_fields(text: str) -> dict[str, str]:
    """Parse DESCRIPTION_*/TLDR/FLUFF from model content (or reasoning fallback).

    Prefers a ===FINAL=== fence when present. Otherwise keeps the *last* value
    for each label (Gemma often drafts labels mid-thinking, then revises).
    """
    cut = text.rfind("===FINAL===")
    if cut != -1:
        text = text[cut + len("===FINAL===") :]

    pattern = re.compile(
        r"(?ms)^\s*(DESCRIPTION_[123]|TLDR|FLUFF)\s*:\s*(.*?)(?=^\s*(?:DESCRIPTION_[123]|TLDR|FLUFF)\s*:|\Z)"
    )
    found: dict[str, str] = {}
    for match in pattern.finditer(text):
        key = match.group(1)
        val = match.group(2).strip()
        # Drop trailing thinking notes after the field body.
        val = re.split(r"(?m)^\s*\*", val, maxsplit=1)[0].strip()
        val = val.replace(EM_DASH, ",")
        # Un-indent accidental leading spaces on each line.
        val = "\n".join(line.strip() for line in val.splitlines()).strip()
        if val:
            found[key] = val
    return found


def build_user_message(title: str, brief: str, same_character: bool) -> str:
    parts = [
        f"Episode title (working or final): {title.strip()}",
        "",
        "Art / mechanism brief:",
        brief.strip(),
        "",
    ]
    if same_character:
        parts.append(
            "HARD FACT: it is ONE character shown on a good day and a bad day "
            "(same face, two moods). Never write as if there are two leaders."
        )
        parts.append("")
    parts.append(
        "Write three competing satirical card teasers (DESCRIPTION_1..3), plus TLDR and FLUFF. "
        "Knife in each teaser: mood talk is cheap; psychological fitness for leaders is not even discussed; "
        "authority theater still covers both moods."
    )
    return "\n".join(parts)


def format_report(fields: dict[str, str], reasoning_len: int) -> str:
    lines = [
        f"# Reptilocracy Gemma candidates (thinking_chars={reasoning_len})",
        "",
    ]
    for key in LABELS:
        lines.append(f"{key}:")
        lines.append(fields.get(key, "(missing)"))
        lines.append("")
    if any(k not in fields for k in LABELS):
        lines.append("WARNING: one or more labels missing; re-run if truncated.")
        lines.append("")
    return "\n".join(lines).rstrip() + "\n"


def main() -> None:
    parser = argparse.ArgumentParser(
        description="Draft Reptilocracy episode copy via Gemma 4 (thinking on)."
    )
    parser.add_argument("--title", required=True, help="Episode title")
    parser.add_argument(
        "--brief",
        required=True,
        help="Art + institutional knife (not a full prop inventory)",
    )
    parser.add_argument(
        "--same-character",
        action=argparse.BooleanOptionalAction,
        default=True,
        help="Tell Gemma the art is one character, two moods (default: true)",
    )
    parser.add_argument(
        "--out",
        type=Path,
        default=None,
        help="Optional path for gemma-candidates.txt",
    )
    parser.add_argument("--model", default=None)
    parser.add_argument("--base-url", default=None)
    parser.add_argument("--temperature", type=float, default=0.7)
    args = parser.parse_args()

    base = args.base_url or os.environ.get("LOCAL_LLM_BASE_URL") or DEFAULT_BASE
    model = args.model or os.environ.get("LOCAL_LLM_MODEL") or DEFAULT_MODEL
    api_key = os.environ.get("LOCAL_LLM_API_KEY") or os.environ.get("OPENAI_API_KEY")

    system = load_pack_system()
    user = build_user_message(args.title, args.brief, args.same_character)
    content, reasoning = chat_complete_thinking(
        base,
        model,
        api_key,
        system,
        user,
        temperature=args.temperature,
    )

    fields = extract_fields(content)
    if len(fields) < 3 and reasoning:
        # Thinking sometimes ate the answer channel; salvage labeled finals from reasoning.
        fields = extract_fields(reasoning) or fields

    report = format_report(fields, len(reasoning))
    print(report)
    if args.out:
        args.out.parent.mkdir(parents=True, exist_ok=True)
        args.out.write_text(report, encoding="utf-8")
        print(f"Wrote {args.out}", file=sys.stderr)

    missing = [k for k in LABELS if k not in fields]
    if missing:
        die(f"Missing labels after parse: {', '.join(missing)}")


if __name__ == "__main__":
    main()
