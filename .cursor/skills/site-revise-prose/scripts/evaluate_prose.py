#!/usr/bin/env python3
"""Evaluate-only prose critique via OpenAI-compatible local gateway.

Default: progressive unit analysis (each paragraph/list field with prior as context).
Does not edit Hugo content. Diagnose only.
"""

from __future__ import annotations

import argparse
import os
import re
import sys
from datetime import date
from pathlib import Path
from typing import Any

# Shared progressive helpers (sibling skill folder)
_SKILLS = Path(__file__).resolve().parents[2]
if str(_SKILLS) not in sys.path:
    sys.path.insert(0, str(_SKILLS))
from site_local_eval_common.common import (  # noqa: E402
    DEFAULT_BASE,
    DEFAULT_MODEL,
    MAX_TOKENS_WHOLE,
    Unit,
    chat_complete,
    die,
    run_progressive,
    split_paragraphs,
)

LIST_KEYS = ("title", "description", "grounding", "sowhat", "subtitle")

SCRIPT_DIR = Path(__file__).resolve().parent
SKILL_DIR = SCRIPT_DIR.parent
PACK_PATH = SKILL_DIR / "packs" / "evaluate.md"
REFERENCE_PATH = SKILL_DIR / "reference.md"


def split_front_matter(text: str) -> tuple[dict[str, Any], str]:
    """Parse simple YAML-ish front matter. Multiline | blocks preserved as strings."""
    if not text.startswith("---"):
        return {}, text
    parts = text.split("---", 2)
    if len(parts) < 3:
        return {}, text
    raw_fm, body = parts[1], parts[2]
    fm: dict[str, Any] = {}
    lines = raw_fm.strip("\n").split("\n")
    i = 0
    while i < len(lines):
        line = lines[i]
        if not line.strip() or line.strip().startswith("#"):
            i += 1
            continue
        m = re.match(r"^([A-Za-z0-9_]+):\s*(.*)$", line)
        if not m:
            i += 1
            continue
        key, val = m.group(1), m.group(2)
        if val in ("|", ">"):
            block: list[str] = []
            i += 1
            while i < len(lines):
                nxt = lines[i]
                if nxt and not nxt[0].isspace() and re.match(r"^[A-Za-z0-9_]+:", nxt):
                    break
                block.append(nxt)
                i += 1
            cleaned = "\n".join(block).strip("\n")
            if cleaned:
                indents = [
                    len(x) - len(x.lstrip(" "))
                    for x in cleaned.split("\n")
                    if x.strip()
                ]
                if indents:
                    n = min(indents)
                    cleaned = "\n".join(
                        x[n:] if len(x) >= n else x for x in cleaned.split("\n")
                    )
            fm[key] = cleaned.strip()
            continue
        val = val.strip()
        if (val.startswith('"') and val.endswith('"')) or (
            val.startswith("'") and val.endswith("'")
        ):
            val = val[1:-1]
        fm[key] = val
        i += 1
    return fm, body.lstrip("\n")


def load_pack_system() -> str:
    if not PACK_PATH.is_file():
        die(f"Missing pack: {PACK_PATH}")
    text = PACK_PATH.read_text(encoding="utf-8")
    cut = text.find("## User message template")
    system_from_pack = text[:cut].strip() if cut != -1 else text.strip()
    ref = ""
    if REFERENCE_PATH.is_file():
        ref = "\n\n---\n\n# Rubric detail\n\n" + REFERENCE_PATH.read_text(
            encoding="utf-8"
        )
    return system_from_pack + ref


def build_excerpt(path: Path, scope: str, fm: dict[str, Any], body: str) -> str:
    title = str(fm.get("title") or path.stem)
    lines: list[str] = []
    if scope in ("full", "list"):
        lines.append(f"title: {title}")
        for k in LIST_KEYS:
            if k == "title":
                continue
            if k in fm and str(fm[k]).strip():
                lines.append(f"{k}:\n{fm[k]}\n")
    if scope in ("full", "body"):
        if scope == "body":
            lines.append(f"(title for context only: {title})\n")
        lines.append("## Body\n")
        lines.append(body.strip())
    if scope == "list" and not any(k in fm for k in LIST_KEYS if k != "title"):
        lines.append("(no list-facing description/grounding/sowhat found)")
    return "\n".join(lines).strip() + "\n"


def build_units(scope: str, fm: dict[str, Any], body: str) -> list[Unit]:
    units: list[Unit] = []
    title = str(fm.get("title") or "").strip()
    if scope in ("full", "list") and title:
        units.append(Unit("title", title))
    if scope in ("full", "list"):
        for k in LIST_KEYS:
            if k == "title":
                continue
            if k in fm and str(fm[k]).strip():
                units.append(Unit(f"front-matter:{k}", str(fm[k]).strip()))
    if scope in ("full", "body"):
        for i, para in enumerate(split_paragraphs(body), start=1):
            units.append(Unit(f"body paragraph {i}", para))
    return units


def build_user_message(path: Path, scope: str, title: str, excerpt: str) -> str:
    return f"""## Target

Path: {path}
Scope: {scope}
Title: {title}

## Excerpt (evaluate this text only)

{excerpt}

## Output format (required)

Structure your answer as:

1. **Cold-read verdict** (1 short paragraph: keep reading? where do you bounce?)
2. **Scores** (table or list): for each of clarity_readability, anti_fluff_compliance, conversational_tone, no_therapist_voice, anti_ai_feel: score 0–2 and one line why
3. **Keep** (2–5 quoted lines that work; one phrase each on why)
4. **Failures** table: Quote (exact) | Criterion or pattern | Why it fails for humans | Fix direction (not full rewrite)
5. **Top 5 fixes** ranked (one actionable line each)
6. **Open questions** for the author (voice choices, not facts)
7. **Compression smell** (if any): quote + short note; else "none"
"""


def build_unit_user(
    path: Path,
    scope: str,
    title: str,
) -> Any:
    def _fn(step: int, n: int, prior: list[Unit], unit: Unit) -> str:
        from site_local_eval_common.common import format_prior_block

        return f"""## Target

Path: {path}
Scope: {scope}
Post title: {title}
Progressive step: {step} of {n}
Current unit label: {unit.label}

## PRIOR units (context only; do not re-score in full)

{format_prior_block(prior)}

## CURRENT unit (evaluate this)

{unit.text}

## Output format (this unit only)

1. **Unit verdict** (1–2 sentences: keep going? any local bounce?)
2. **Scores** 0–2 for this unit only: clarity_readability | anti_fluff_compliance | conversational_tone | no_therapist_voice | anti_ai_feel (one line why each)
3. **Keep** (0–3 quotes from CURRENT only)
4. **Failures** from CURRENT: Quote | Criterion/pattern | Why | Fix direction
5. **Continuity** vs PRIOR (echo, setup unpaid, voice drift) or "none"
6. **Top fixes for this unit** (at most 3)
"""

    return _fn


def build_synthesis_user(
    path: Path,
    scope: str,
    title: str,
) -> Any:
    def _fn(unit_notes: list[tuple[Unit, str]]) -> str:
        digest_parts = []
        for unit, note in unit_notes:
            digest_parts.append(f"### {unit.label}\n\n{note}\n")
        digest = "\n".join(digest_parts)
        return f"""## Target

Path: {path}
Scope: {scope}
Title: {title}

You already judged each unit in order. Below are the per-unit notes (and unit labels).
Synthesize a **piece-level** evaluate-only report. Prefer patterns that recur across units.
Do not invent new quotes that never appeared in the notes or units.

## Per-unit notes

{digest}

## Output format (required: overall)

1. **Cold-read verdict** (1 short paragraph for the whole piece)
2. **Scores** 0–2 overall: clarity_readability | anti_fluff_compliance | conversational_tone | no_therapist_voice | anti_ai_feel
3. **Keep** (2–5 best quotes from the notes)
4. **Failures** consolidating the worst issues across units (quote | criterion | fix direction)
5. **Top 5 fixes** ranked for the author
6. **Open questions** for the author
7. **Compression smell** (if any) or "none"
8. **Progress map** (bullets): unit label → drag / momentum / note
"""

    return _fn


def write_out(
    out_path: Path,
    path: Path,
    model: str,
    scope: str,
    mode: str,
    report: str,
) -> None:
    out_path.parent.mkdir(parents=True, exist_ok=True)
    header = f"""# Local prose review

- **date:** {date.today().isoformat()}
- **path:** `{path}`
- **scope:** {scope}
- **mode:** {mode}
- **model:** `{model}`

---

"""
    out_path.write_text(header + report + "\n", encoding="utf-8")
    print(f"Wrote {out_path}", file=sys.stderr)


def main() -> None:
    p = argparse.ArgumentParser(
        description="Evaluate-only local prose review (no content edits)."
    )
    p.add_argument("path", type=Path, help="Path to Hugo markdown (e.g. index.md)")
    p.add_argument(
        "--scope",
        choices=("full", "list", "body"),
        default="full",
        help="Text scope to evaluate (default: full)",
    )
    p.add_argument(
        "--mode",
        choices=("progressive", "whole"),
        default="progressive",
        help="progressive=one unit at a time with prior context (default); whole=single pass",
    )
    p.add_argument(
        "--model",
        default=os.environ.get("LOCAL_LLM_MODEL", DEFAULT_MODEL),
        help=f"Model id (default: {DEFAULT_MODEL})",
    )
    p.add_argument(
        "--base-url",
        default=os.environ.get("LOCAL_LLM_BASE_URL", DEFAULT_BASE),
        help=f"OpenAI-compatible base URL (default: {DEFAULT_BASE})",
    )
    p.add_argument(
        "--out",
        type=Path,
        default=None,
        help="Optional path to write report markdown (e.g. docs/research/...)",
    )
    args = p.parse_args()

    path: Path = args.path
    if not path.is_file():
        die(f"File not found: {path}")

    text = path.read_text(encoding="utf-8")
    fm, body = split_front_matter(text)
    title = str(fm.get("title") or path.stem)
    system = load_pack_system()
    api_key = os.environ.get("LOCAL_LLM_API_KEY") or None

    if args.mode == "whole":
        excerpt = build_excerpt(path, args.scope, fm, body)
        if not excerpt.strip():
            die("Empty excerpt after scope load; nothing to evaluate.")
        user = build_user_message(path, args.scope, title, excerpt)
        report = chat_complete(
            args.base_url,
            args.model,
            api_key,
            system,
            user,
            max_tokens=MAX_TOKENS_WHOLE,
        )
    else:
        units = build_units(args.scope, fm, body)
        if not units:
            die("No units to evaluate after split.")
        report = run_progressive(
            units=units,
            system=system,
            build_unit_user=build_unit_user(path, args.scope, title),
            build_synthesis_user=build_synthesis_user(path, args.scope, title),
            base_url=args.base_url,
            model=args.model,
            api_key=api_key,
            language="en",
        )

    print(report)

    if args.out is not None:
        write_out(args.out, path, args.model, args.scope, args.mode, report)


if __name__ == "__main__":
    main()
