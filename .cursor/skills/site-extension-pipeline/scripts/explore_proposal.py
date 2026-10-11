#!/usr/bin/env python3
"""Proposal-only Explore further pipeline: Gemma research review + link hooks.

Does not patch Hugo content. Writes a reviewable YAML proposal under tmp/ by default.
See .cursor/skills/site-extension-pipeline/SKILL.md and site-explore-research-review.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
import urllib.error
import urllib.request
from pathlib import Path
from typing import Any

SCRIPT_DIR = Path(__file__).resolve().parent
if str(SCRIPT_DIR) not in sys.path:
    sys.path.insert(0, str(SCRIPT_DIR))

from explore_pipeline_common import (  # noqa: E402
    bundle_slug_from_essay,
    export_file_ready,
    repo_root_from_here,
    sync_export_paths,
    write_json,
)

REPO_ROOT = repo_root_from_here(Path(__file__))
_SKILLS = REPO_ROOT / ".cursor" / "skills"
if str(_SKILLS) not in sys.path:
    sys.path.insert(0, str(_SKILLS))

from site_local_eval_common.common import DEFAULT_BASE, die  # noqa: E402

MODEL = "gemini_free/gemma-4-26b-a4b-it"
MAX_TOKENS = 4096


def health(base: str) -> None:
    url = base.replace("/v1", "").rstrip("/") + "/health"
    try:
        with urllib.request.urlopen(url, timeout=10) as resp:
            if resp.status != 200:
                die(f"Health check failed: HTTP {resp.status}")
    except urllib.error.URLError as e:
        die(f"Polypus unreachable at {url}: {e.reason}")


def chat_thinking(
    base_url: str,
    system: str,
    user: str,
    *,
    max_tokens: int = MAX_TOKENS,
) -> str:
    url = base_url.rstrip("/") + "/chat/completions"
    payload = {
        "model": MODEL,
        "messages": [
            {"role": "system", "content": system},
            {"role": "user", "content": user},
        ],
        "temperature": 0.25,
        "max_tokens": max_tokens,
        "chat_template_kwargs": {"enable_thinking": True},
    }
    data = json.dumps(payload).encode("utf-8")
    req = urllib.request.Request(
        url,
        data=data,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    try:
        with urllib.request.urlopen(req, timeout=300) as resp:
            raw = resp.read().decode("utf-8")
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        die(f"Gateway HTTP {e.code}:\n{body}")
    except urllib.error.URLError as e:
        die(f"Gateway error: {e.reason}")

    obj = json.loads(raw)
    content = (obj.get("choices") or [{}])[0].get("message", {}).get("content")
    if not content:
        die(f"Empty model response:\n{raw[:800]}")
    return str(content).strip()


def strip_thinking(text: str) -> str:
    """Drop common thinking wrappers; keep YAML-heavy tail."""
    text = re.sub(r"^```(?:yaml)?\s*", "", text.strip())
    text = re.sub(r"\s*```\s*$", "", text)
    if "PRIMARY_PROMPT:" in text or "verdict:" in text or "reviews:" in text:
        return text.strip()
    parts = re.split(r"\n(?=(?:reviews:|candidate_id:|approve:|---\n))", text)
    for p in reversed(parts):
        if re.search(r"verdict:\s*(approve|reject)", p, re.I):
            return p.strip()
    return text.strip()


def load_text(path: Path) -> str:
    if not path.is_file():
        die(f"Missing file: {path}")
    return path.read_text(encoding="utf-8")


def load_research(candidate: dict[str, Any]) -> str:
    export = candidate.get("research_export")
    if export:
        p = Path(export)
        if not p.is_absolute():
            p = REPO_ROOT / p
        if p.is_file() and p.stat().st_size >= 80:
            return load_text(p)
    summary = candidate.get("research_summary")
    if summary and str(summary).strip():
        return str(summary).strip()
    if export:
        die(
            f"Candidate {candidate.get('id')!r}: missing research_export file {export} "
            "and no research_summary"
        )
    die(f"Candidate {candidate.get('id')!r} needs research_export or research_summary")


def parse_explore_excludes(essay_md: str) -> list[str]:
    lines: list[str] = []
    in_explore = False
    for line in essay_md.splitlines():
        if re.match(r"^\s*explore:\s*$", line):
            in_explore = True
            continue
        if in_explore:
            if line and not line[0].isspace() and not line.startswith("-"):
                break
            if "label:" in line:
                m = re.search(r'label:\s*"(.*)"', line)
                if m:
                    lines.append(f"label: {m.group(1)}")
            if "url:" in line:
                m = re.search(r'url:\s*"(.*)"', line)
                if m:
                    lines.append(f"url: {m.group(1)}")
            if "query:" in line:
                m = re.search(r'query:\s*"(.*)"', line)
                if m:
                    lines.append(f"query: {m.group(1)}")
    return lines


def build_review_user(
    essay_md: str,
    packet: dict[str, Any],
    candidates: list[dict[str, Any]],
) -> str:
    excludes = list(parse_explore_excludes(essay_md))
    extra = packet.get("excludes")
    if isinstance(extra, list) and extra:
        excludes = list(dict.fromkeys(excludes + extra))

    blocks = []
    for c in candidates:
        blocks.append(
            f"### {c['id']}\n"
            f"url: {c['url']}\n"
            f"hinge: {c.get('hinge', '')}\n\n"
            f"{load_research(c)}\n"
        )

    prompts = packet.get("prompts") or {}
    prompt_block = ""
    if prompts:
        prompt_block = f"GEMMA PROMPTS USED:\n{json.dumps(prompts, indent=2)}\n\n"

    calibration = packet.get("pipeline_calibration")
    cal_block = ""
    if calibration:
        cal_block = (
            "CALIBRATION: Score research quality and fit even if the URL already appears in EXCLUDE. "
            "Use verdict approve when quality meets bar (for hook refresh or operator validation); "
            "use reject only for weak research or clichés, not merely because the URL is on the page.\n\n"
        )

    return f"""Review Perplexity research candidates for Explore further on this essay.

EDITORIAL: Structural diagnosis; second-order mechanics; ban GPS/pomodoro/generic productivity/ADHD life-hack clichés.
TARGET SHIP COUNT: {packet.get('target_count', 3)} approved threads maximum.

{cal_block}EXCLUDE (already on page or reserved; do not ship duplicates unless CALIBRATION is set):
{chr(10).join(excludes) or '(none)'}

{prompt_block}ESSAY (markdown):
{essay_md[:12000]}

CANDIDATES:
{chr(10).join(blocks)}

OUTPUT (YAML only, no thinking leak):
reviews:
  - candidate_id: ...
    url: ...
    scores:
      evidence: 1-5
      novelty: 1-5
      source_fit: 1-5
      counter_cost_value: 1-5
    duplication_risk: low|medium|high
    verdict: approve|reject
    ship_priority: 1|2|3|null
    reason: one sentence
batch_summary:
  approved_count: N
  counter_cost_present: yes|no
  notes: optional
"""


def build_hooks_user(
    essay_md: str,
    approved: list[dict[str, Any]],
    reviews_yaml: str,
) -> str:
    return f"""Create Explore further link copy for approved Perplexity threads only.

RULES:
- label: 4-8 words, curiosity gap, not a topic tag
- hook: one sentence, max ~20 words, names why this link exists
- No em dash U+2014
- If Spanish sibling exists for this post, add es_label and es_hook (native Spanish)

ESSAY excerpt:
{essay_md[:8000]}

APPROVED REVIEWS:
{reviews_yaml}

OUTPUT YAML only:
hooks:
  - candidate_id: ...
    label: ...
    hook: ...
    es_label: ... (optional)
    es_hook: ... (optional)
review:
  approve: yes|no
  revisions: ...
"""


def main() -> int:
    parser = argparse.ArgumentParser(description="Explore further proposal (Gemma, no Hugo write)")
    parser.add_argument(
        "--essay",
        type=Path,
        required=True,
        help="Path to English index.md",
    )
    parser.add_argument(
        "--candidates",
        type=Path,
        required=True,
        help="JSON packet: target_count, excludes, prompts, candidates[]",
    )
    parser.add_argument(
        "--out",
        type=Path,
        default=None,
        help="Proposal output path (default: tmp/explore-proposals/<slug>.proposal.yaml)",
    )
    parser.add_argument(
        "--base-url",
        default=DEFAULT_BASE,
        help="Polypus OpenAI base URL",
    )
    parser.add_argument(
        "--no-prepare-exports",
        action="store_true",
        help="Do not sync export paths or require export files",
    )
    parser.add_argument(
        "--allow-summary-only",
        action="store_true",
        help="Allow research_summary without export files",
    )
    args = parser.parse_args()

    essay_path = args.essay if args.essay.is_absolute() else REPO_ROOT / args.essay
    cand_path = (
        args.candidates if args.candidates.is_absolute() else REPO_ROOT / args.candidates
    )
    slug = bundle_slug_from_essay(essay_path)
    packet = json.loads(load_text(cand_path))
    if not args.no_prepare_exports:
        manifest = sync_export_paths(packet, slug)
        cand_path.write_text(json.dumps(packet, indent=2) + "\n", encoding="utf-8")
        manifest_path = (
            REPO_ROOT / "tmp" / "explore-proposals" / slug / "export-manifest.json"
        )
        write_json(
            manifest_path,
            {
                "slug": slug,
                "candidates_file": str(cand_path.relative_to(REPO_ROOT)),
                "rows": manifest,
                "mcp_tool": "user-perplexity-browser/perplexity_export",
            },
        )

        if not args.allow_summary_only:
            missing: list[str] = []
            for c in packet.get("candidates") or []:
                rel = c.get("research_export", "")
                if export_file_ready(REPO_ROOT, str(rel)):
                    continue
                if c.get("research_summary"):
                    continue
                missing.append(str(c.get("id")))
            if missing:
                die(
                    "Missing research_export files for: "
                    + ", ".join(missing)
                    + f". Run: make explore-fetch-exports CANDIDATES="
                    f"{cand_path.relative_to(REPO_ROOT)} ESSAY="
                    f"{essay_path.relative_to(REPO_ROOT)} SYNC=1, then agent MCP "
                    f"perplexity_export per {manifest_path.relative_to(REPO_ROOT)}"
                )

    essay_md = load_text(essay_path)
    candidates = packet.get("candidates") or []
    if not candidates:
        die("candidates[] is empty")

    health(args.base_url)

    review_raw = chat_thinking(
        args.base_url,
        "You are a research editor for Behavior Engineering. Output YAML only.",
        build_review_user(essay_md, packet, candidates),
    )
    review_yaml = strip_thinking(review_raw)

    approved_ids = []
    for block in re.split(r"\n\s*-\s*candidate_id:", review_yaml):
        if not block.strip():
            continue
        chunk = block if block.strip().startswith("candidate_id:") else "candidate_id:" + block
        cid_m = re.search(r"candidate_id:\s*(\S+)", chunk)
        if not cid_m:
            continue
        if re.search(r"verdict:\s*approve\b", chunk, re.I):
            approved_ids.append(cid_m.group(1).rstrip(","))
    approved = [c for c in candidates if c["id"] in approved_ids]
    if not approved:
        out_path = args.out
        if out_path is None:
            slug = essay_path.parent.name
            out_path = REPO_ROOT / "tmp" / "explore-proposals" / f"{slug}.proposal.yaml"
        out_path.parent.mkdir(parents=True, exist_ok=True)
        out_path.write_text(
            f"# No approved candidates\n\n## reviews\n{review_yaml}\n",
            encoding="utf-8",
        )
        print(f"Wrote {out_path} (0 approved)", file=sys.stderr)
        print(str(out_path))
        return 0

    hooks_raw = chat_thinking(
        args.base_url,
        "You write click-worthy Explore further labels and hooks. Output YAML only.",
        build_hooks_user(essay_md, approved, review_yaml),
    )
    hooks_yaml = strip_thinking(hooks_raw)

    explore_rows = []
    for c in approved:
        explore_rows.append(
            {
                "type": "perplexity_thread",
                "candidate_id": c["id"],
                "url": c["url"],
                "note": "Apply label/hook from hooks section after operator review",
            }
        )

    slug = essay_path.parent.name
    out_path = args.out
    if out_path is None:
        out_path = REPO_ROOT / "tmp" / "explore-proposals" / f"{slug}.proposal.yaml"
    if not out_path.is_absolute():
        out_path = REPO_ROOT / out_path
    out_path.parent.mkdir(parents=True, exist_ok=True)

    body = f"""# Explore further proposal (do not apply without operator OK)
# Essay: {essay_path.relative_to(REPO_ROOT) if essay_path.is_relative_to(REPO_ROOT) else essay_path}
# Generated by scripts/explore_proposal.py

## reviews
{review_yaml}

## hooks
{hooks_yaml}

## suggested_explore_rows
{json.dumps(explore_rows, indent=2)}
"""
    out_path.write_text(body, encoding="utf-8")
    print(f"Wrote {out_path}", file=sys.stderr)
    print(str(out_path))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
