#!/usr/bin/env python3
"""Review applicability (Gemma) and optionally merge Explore proposal into Hugo.

Default: Gemma thinking applicability review + operator_question (no file writes).
With --apply --yes: patch index.md and index.es.md explore rows.

Agent workflow: run without --apply, AskQuestion using operator_question, then
go tool task explore-apply APPLY=1 CONFIRM=1 if the human confirms.
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
from typing import Any

SCRIPT_DIR = Path(__file__).resolve().parent
if str(SCRIPT_DIR) not in sys.path:
    sys.path.insert(0, str(SCRIPT_DIR))

from explore_pipeline_common import bundle_slug_from_essay, repo_root_from_here  # noqa: E402
from explore_share import urls_from_proposal  # noqa: E402
from explore_share_manifest import (  # noqa: E402
    load_manifest,
    manifest_path_for_slug,
    validate_manifest_for_urls,
)

REPO_ROOT = repo_root_from_here(Path(__file__))
_SKILLS = REPO_ROOT / ".cursor" / "skills"
if str(_SKILLS) not in sys.path:
    sys.path.insert(0, str(_SKILLS))

from site_local_eval_common.common import DEFAULT_BASE, die  # noqa: E402

MODEL = "gemini_free/gemma-4-26b-a4b-it"


def load_text(path: Path) -> str:
    if not path.is_file():
        die(f"Missing file: {path}")
    return path.read_text(encoding="utf-8")


def health(base: str) -> None:
    url = base.replace("/v1", "").rstrip("/") + "/health"
    try:
        with urllib.request.urlopen(url, timeout=10) as resp:
            if resp.status != 200:
                die(f"Health check failed: HTTP {resp.status}")
    except urllib.error.URLError as e:
        die(f"Polypus unreachable at {url}: {e.reason}")


def chat_thinking(base_url: str, system: str, user: str) -> str:
    url = base_url.rstrip("/") + "/chat/completions"
    payload = {
        "model": MODEL,
        "messages": [
            {"role": "system", "content": system},
            {"role": "user", "content": user},
        ],
        "temperature": 0.2,
        "max_tokens": 2048,
        "chat_template_kwargs": {"enable_thinking": True},
    }
    data = json.dumps(payload).encode("utf-8")
    req = urllib.request.Request(
        url,
        data=data,
        headers={"Content-Type": "application/json"},
        method="POST",
    )
    with urllib.request.urlopen(req, timeout=300) as resp:
        obj = json.loads(resp.read().decode("utf-8"))
    content = (obj.get("choices") or [{}])[0].get("message", {}).get("content")
    if not content:
        die("Empty Gemma response")
    text = str(content).strip()
    text = re.sub(r"^```(?:yaml)?\s*", "", text)
    text = re.sub(r"\s*```\s*$", "", text)
    return text


def parse_proposal_hooks(proposal_text: str) -> list[dict[str, str]]:
    m = re.search(r"## hooks\s*\n(.*)(?:\n## |\Z)", proposal_text, re.S)
    if not m:
        return []
    block = m.group(1)
    rows: list[dict[str, str]] = []
    for chunk in re.split(r"\n\s*-\s*candidate_id:", block):
        stripped = chunk.strip()
        if not stripped or stripped.startswith("hooks:"):
            continue
        piece = chunk if stripped.startswith("candidate_id:") else "candidate_id:" + chunk
        cid = re.search(r"candidate_id:\s*(\S+)", piece)
        label = re.search(r'^\s+label:\s*"([^"]*)"', piece, re.M)
        hook = re.search(r'^\s+hook:\s*"([^"]*)"', piece, re.M)
        es_label = re.search(r'^\s+es_label:\s*"([^"]*)"', piece, re.M)
        es_hook = re.search(r'^\s+es_hook:\s*"([^"]*)"', piece, re.M)
        if not cid or not label or not hook:
            continue
        row = {
            "candidate_id": cid.group(1).rstrip(","),
            "label": label.group(1),
            "hook": hook.group(1),
        }
        if es_label:
            row["es_label"] = es_label.group(1)
        if es_hook:
            row["es_hook"] = es_hook.group(1)
        rows.append(row)
    return rows


def parse_proposal_urls(proposal_text: str) -> dict[str, str]:
    """Map candidate_id -> url from reviews section."""
    urls: dict[str, str] = {}
    m = re.search(r"## reviews\s*\n(.*)(?:\n## |\Z)", proposal_text, re.S)
    if not m:
        return urls
    block = m.group(1)
    for chunk in re.split(r"\n\s*-\s*candidate_id:", block):
        piece = chunk if chunk.strip().startswith("candidate_id:") else "candidate_id:" + chunk
        cid = re.search(r"candidate_id:\s*(\S+)", piece)
        url = re.search(r"url:\s*(https://\S+)", piece)
        verdict = re.search(r"verdict:\s*(\w+)", piece, re.I)
        if not cid or not url or not verdict:
            continue
        if verdict.group(1).lower() != "approve":
            continue
        urls[cid.group(1).rstrip(",")] = url.group(1).rstrip(",")
    return urls


def parse_suggested_rows(proposal_text: str) -> list[dict[str, Any]]:
    m = re.search(r"## suggested_explore_rows\s*\n(\[.*?\])", proposal_text, re.S)
    if not m:
        return []
    try:
        return json.loads(m.group(1))
    except json.JSONDecodeError:
        return []


def build_ship_rows(
    proposal_text: str,
    candidates_path: Path | None,
) -> list[dict[str, str]]:
    hooks = parse_proposal_hooks(proposal_text)
    id_to_url = parse_proposal_urls(proposal_text)
    if candidates_path and candidates_path.is_file():
        packet = json.loads(load_text(candidates_path))
        for c in packet.get("candidates") or []:
            cid = c.get("id")
            if cid and c.get("url"):
                id_to_url.setdefault(str(cid), str(c["url"]))
    for row in parse_suggested_rows(proposal_text):
        cid = row.get("candidate_id")
        url = row.get("url")
        if cid and url:
            id_to_url.setdefault(str(cid), str(url))

    ship: list[dict[str, str]] = []
    for h in hooks:
        cid = h["candidate_id"]
        url = id_to_url.get(cid)
        if not url:
            continue
        ship.append(
            {
                "url": url,
                "label": h["label"],
                "hook": h["hook"],
                "es_label": h.get("es_label", h["label"]),
                "es_hook": h.get("es_hook", h["hook"]),
            }
        )
    return ship


def build_applicability_user(
    essay_en: str,
    essay_es: str | None,
    proposal_excerpt: str,
    ship_rows: list[dict[str, str]],
) -> str:
    es_block = essay_es[:6000] if essay_es else "(no Spanish sibling)"
    return f"""Review whether this Explore further proposal should be merged into Hugo front matter.

RULES:
- recommend_apply: yes only if hooks match approved research and improve on-page copy without duplicating weak rows
- recommend_apply: partial if only some rows should land
- recommend_apply: no if blockers exist (cliché hooks, wrong essay, stale proposal)
- operator_question: one plain English question the human must answer before any YAML apply (mention EN and ES if both exist)

PROPOSED SHIP ROWS:
{json.dumps(ship_rows, indent=2)}

PROPOSAL EXCERPT:
{proposal_excerpt[:8000]}

EN essay front matter:
{essay_en[:6000]}

ES essay front matter (if any):
{es_block}

OUTPUT YAML only:
applicability:
  recommend_apply: yes|no|partial
  operator_question: "..."
  safe_actions:
    - action: add|update|skip
      url: "..."
      reason: "..."
  blockers: []
"""


def upsert_perplexity_threads(
    md: str,
    rows: list[dict[str, str]],
    *,
    use_spanish: bool,
    replace_all: bool = False,
) -> str:
    if not rows:
        return md

    def block_for(r: dict[str, str]) -> str:
        label = r["es_label"] if use_spanish else r["label"]
        hook = r["es_hook"] if use_spanish else r["hook"]
        return (
            f'    - type: perplexity_thread\n'
            f'      label: "{label}"\n'
            f'      hook: "{hook}"\n'
            f'      url: "{r["url"]}"\n'
        )

    fm_end = md.find("\n---\n", 3)
    if fm_end == -1:
        die("Could not parse front matter")
    fm = md[: fm_end + 1]
    body = md[fm_end + 1 :]

    explore_m = re.search(r"(\n  explore:\n)(.*?)(?=\n  [a-z_]+:|\n---|\Z)", fm, re.S)
    if not explore_m:
        die("No reader_landing.explore block found")

    prefix, explore_body = explore_m.group(1), explore_m.group(2)
    new_blocks = "".join(block_for(r) for r in rows)
    if replace_all:
        explore_body = "\n" + new_blocks
    else:
        for r in rows:
            url = re.escape(r["url"])
            new_block = block_for(r)
            item_pat = re.compile(
                rf"    - type: perplexity_thread\n(?:      .+\n)*?      url: \"{url}\"\n",
                re.M,
            )
            if item_pat.search(explore_body):
                explore_body = item_pat.sub(new_block, explore_body)
            else:
                explore_body = explore_body.rstrip("\n") + "\n" + new_block

    new_fm = fm[: explore_m.start()] + prefix + explore_body + fm[explore_m.end() :]
    return new_fm + body


def main() -> int:
    parser = argparse.ArgumentParser(description="Explore proposal applicability + optional apply")
    parser.add_argument("--essay", type=Path, required=True)
    parser.add_argument("--proposal", type=Path, required=True)
    parser.add_argument("--candidates", type=Path, default=None)
    parser.add_argument("--base-url", default=DEFAULT_BASE)
    parser.add_argument(
        "--apply",
        action="store_true",
        help="Merge into index.md / index.es.md (requires --yes)",
    )
    parser.add_argument(
        "--yes",
        action="store_true",
        help="Confirm apply after operator answered yes",
    )
    parser.add_argument(
        "--share-manifest",
        type=Path,
        default=None,
        help="share-manifest.json (default: tmp/explore-proposals/<slug>/share-manifest.json)",
    )
    parser.add_argument(
        "--skip-share-preflight",
        action="store_true",
        help="Emergency only; prefer SKIP_SHARE_PREFLIGHT=1 in tests",
    )
    args = parser.parse_args()

    essay_path = args.essay if args.essay.is_absolute() else REPO_ROOT / args.essay
    prop_path = args.proposal if args.proposal.is_absolute() else REPO_ROOT / args.proposal
    cand_path = args.candidates
    if cand_path and not cand_path.is_absolute():
        cand_path = REPO_ROOT / cand_path

    proposal_text = load_text(prop_path)
    essay_en = load_text(essay_path)
    es_path = essay_path.parent / "index.es.md"
    essay_es = load_text(es_path) if es_path.is_file() else None

    skip_share = args.skip_share_preflight or os.environ.get("SKIP_SHARE_PREFLIGHT") == "1"
    if not skip_share:
        manifest_path = args.share_manifest
        if manifest_path is None:
            slug = bundle_slug_from_essay(essay_path)
            manifest_path = manifest_path_for_slug(REPO_ROOT, slug)
        elif not manifest_path.is_absolute():
            manifest_path = REPO_ROOT / manifest_path
        if not manifest_path.is_file():
            die(
                f"Missing share manifest {manifest_path}. "
                "Run: go tool task explore-share-prepare CANDIDATES=... then confirm + verify-cold "
                "(see PERPLEXITY-SHARE.md)."
            )
        manifest = load_manifest(manifest_path)
        required_urls = urls_from_proposal(proposal_text)
        if not required_urls:
            die("No approved Perplexity URLs in proposal for share preflight")
        preflight = validate_manifest_for_urls(manifest, required_urls)
        for w in preflight.warnings:
            print(f"share preflight warning: {w}", file=sys.stderr)
        if not preflight.ok:
            for err in preflight.errors:
                print(f"share preflight error: {err}", file=sys.stderr)
            die(
                "Share preflight failed. Fix manifest before explore-apply "
                "(go tool task explore-share-check PROPOSAL=... CANDIDATES=...)."
            )
        print("Share preflight OK.", file=sys.stderr)

    ship_rows = build_ship_rows(proposal_text, cand_path)
    if not ship_rows:
        die("No shippable rows in proposal (hooks + approved URLs)")

    health(args.base_url)
    applicability_yaml = chat_thinking(
        args.base_url,
        "You gate Hugo content changes. Output YAML only.",
        build_applicability_user(essay_en, essay_es, proposal_text, ship_rows),
    )

    sidecar = prop_path.with_suffix(".applicability.yaml")
    sidecar.write_text(applicability_yaml + "\n", encoding="utf-8")
    print(applicability_yaml)
    print(f"\nWrote {sidecar}", file=sys.stderr)

    if not args.apply:
        print(
            "\nNext: ask the operator the operator_question above. "
            "If they confirm, run with --apply --yes\n"
            "Before apply: share manifest must pass (go tool task explore-share-check). "
            "See .cursor/skills/site-extension-pipeline/PERPLEXITY-SHARE.md.",
            file=sys.stderr,
        )
        return 0

    if not args.yes:
        die("Refusing --apply without --yes (operator must confirm after applicability review)")

    rec = re.search(r"recommend_apply:\s*(\w+)", applicability_yaml, re.I)
    if rec and rec.group(1).lower() == "no":
        die("Gemma recommend_apply: no; fix proposal or override manually")

    new_en = upsert_perplexity_threads(
        essay_en, ship_rows, use_spanish=False, replace_all=True
    )
    essay_path.write_text(new_en, encoding="utf-8")
    if essay_es:
        new_es = upsert_perplexity_threads(
            essay_es, ship_rows, use_spanish=True, replace_all=True
        )
        es_path.write_text(new_es, encoding="utf-8")
    print(f"Applied {len(ship_rows)} row(s) to {essay_path}", file=sys.stderr)
    if essay_es:
        print(f"Applied Spanish hooks to {es_path}", file=sys.stderr)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
