#!/usr/bin/env python3
"""Prepare and verify Perplexity public-share manifest before explore-apply.

See .cursor/skills/site-extension-pipeline/PERPLEXITY-SHARE.md and PERPLEXITY-SHARE-BROWSER.md.
"""

from __future__ import annotations

import argparse
import json
import re
import sys
from pathlib import Path

SCRIPT_DIR = Path(__file__).resolve().parent
if str(SCRIPT_DIR) not in sys.path:
    sys.path.insert(0, str(SCRIPT_DIR))

from explore_pipeline_common import (  # noqa: E402
    bundle_slug_from_essay,
    load_json,
    repo_root_from_here,
)
from explore_share_manifest import (  # noqa: E402
    canonical_perplexity_url,
    cold_probe_url,
    confirm_share_row,
    empty_manifest,
    load_manifest,
    manifest_path_for_slug,
    merge_manifest_rows,
    prepare_rows_from_packet,
    record_cold_verification,
    save_manifest,
    validate_manifest_for_urls,
)

REPO_ROOT = repo_root_from_here(Path(__file__))


def die(msg: str, code: int = 2) -> None:
    print(msg, file=sys.stderr)
    raise SystemExit(code)


def resolve_slug(candidates_path: Path, essay_path: Path | None, packet: dict) -> str:
    if essay_path is not None:
        return bundle_slug_from_essay(essay_path)
    ep = packet.get("essay_path")
    if ep:
        p = Path(ep)
        if not p.is_absolute():
            p = REPO_ROOT / p
        return bundle_slug_from_essay(p)
    die("pass --essay or essay_path in candidates JSON")


def urls_from_proposal(proposal_text: str) -> list[str]:
    urls: list[str] = []
    m = re.search(r"## reviews\s*\n(.*)(?:\n## |\Z)", proposal_text, re.S)
    if m:
        block = m.group(1)
        for chunk in re.split(r"\n\s*-\s*candidate_id:", block):
            piece = chunk if chunk.strip().startswith("candidate_id:") else "candidate_id:" + chunk
            if not re.search(r"verdict:\s*approve\b", piece, re.I):
                continue
            url_m = re.search(r"url:\s*(https://\S+)", piece)
            if url_m:
                urls.append(url_m.group(1).rstrip(","))
    m2 = re.search(r"## suggested_explore_rows\s*\n(\[.*?\])", proposal_text, re.S)
    if m2:
        try:
            rows = json.loads(m2.group(1))
            for row in rows:
                u = row.get("url")
                if u:
                    urls.append(str(u))
        except json.JSONDecodeError:
            pass
    deduped: list[str] = []
    seen: set[str] = set()
    for u in urls:
        c = canonical_perplexity_url(u)
        if c not in seen:
            seen.add(c)
            deduped.append(c)
    return deduped


def cmd_prepare(args: argparse.Namespace) -> int:
    cand_path = args.candidates if args.candidates.is_absolute() else REPO_ROOT / args.candidates
    packet = load_json(cand_path)
    essay_path = args.essay
    if essay_path:
        essay_path = essay_path if essay_path.is_absolute() else REPO_ROOT / essay_path
    slug = resolve_slug(cand_path, essay_path, packet)
    manifest_path = manifest_path_for_slug(REPO_ROOT, slug)
    new_rows = prepare_rows_from_packet(packet)
    if not new_rows:
        die("no candidates with url in packet")

    if manifest_path.is_file():
        manifest = load_manifest(manifest_path)
    else:
        rel = str(cand_path.relative_to(REPO_ROOT)) if cand_path.is_relative_to(REPO_ROOT) else str(cand_path)
        manifest = empty_manifest(slug, candidates_file=rel)

    manifest = merge_manifest_rows(manifest, new_rows)
    manifest["candidates_file"] = manifest.get("candidates_file") or str(
        cand_path.relative_to(REPO_ROOT) if cand_path.is_relative_to(REPO_ROOT) else cand_path
    )
    save_manifest(manifest_path, manifest)
    print(manifest_path)
    print(f"Prepared {len(manifest['rows'])} row(s)", file=sys.stderr)
    print(
        "\nNext: Perplexity Share → Anyone with the link for each thread, then:\n"
        "  make explore-share-confirm CANDIDATES=... CANDIDATE_ID=<id>\n"
        "  make explore-share-verify-cold CANDIDATES=... CANDIDATE_ID=<id> PROBE=1 OPERATOR=1\n"
        "See .cursor/skills/site-extension-pipeline/PERPLEXITY-SHARE-BROWSER.md",
        file=sys.stderr,
    )
    return 0


def cmd_confirm(args: argparse.Namespace) -> int:
    cand_path = args.candidates if args.candidates.is_absolute() else REPO_ROOT / args.candidates
    packet = load_json(cand_path)
    slug = resolve_slug(cand_path, None, packet)
    manifest_path = manifest_path_for_slug(REPO_ROOT, slug)
    if not manifest_path.is_file():
        die(f"missing manifest {manifest_path}; run explore-share-prepare first")
    manifest = load_manifest(manifest_path)
    cid = args.candidate_id
    confirm_share_row(
        manifest,
        cid,
        url=args.url,
        confirmed_by=args.by or "operator",
    )
    save_manifest(manifest_path, manifest)
    print(f"Confirmed share for {cid}", file=sys.stderr)
    return 0


def cmd_verify_cold(args: argparse.Namespace) -> int:
    cand_path = args.candidates if args.candidates.is_absolute() else REPO_ROOT / args.candidates
    packet = load_json(cand_path)
    slug = resolve_slug(cand_path, None, packet)
    manifest_path = manifest_path_for_slug(REPO_ROOT, slug)
    if not manifest_path.is_file():
        die(f"missing manifest {manifest_path}")
    manifest = load_manifest(manifest_path)
    cid = args.candidate_id
    row = next((r for r in manifest.get("rows") or [] if str(r.get("candidate_id")) == cid), None)
    if not row:
        die(f"candidate_id {cid} not in manifest")

    url = str(row.get("url") or "")
    if args.probe:
        probe = cold_probe_url(url)
        print(json.dumps(probe, indent=2))
        record_cold_verification(manifest, cid, method="http_probe", probe=probe)
    if args.operator:
        record_cold_verification(
            manifest,
            cid,
            method="operator_incognito",
            operator_confirmed=True,
        )
    if not args.probe and not args.operator:
        die("pass PROBE=1 and/or OPERATOR=1 (operator = incognito check recorded)")

    save_manifest(manifest_path, manifest)
    row = next(r for r in manifest["rows"] if str(r.get("candidate_id")) == cid)
    if row.get("cold_verified"):
        print(f"Cold verification recorded for {cid}", file=sys.stderr)
        return 0
    die(
        f"Cold verification failed for {cid}: {row.get('cold_probe_notes')}. "
        "Use OPERATOR=1 after a real incognito pass.",
        1,
    )


def cmd_check(args: argparse.Namespace) -> int:
    manifest_path = args.manifest
    if manifest_path:
        manifest_path = manifest_path if manifest_path.is_absolute() else REPO_ROOT / manifest_path
    else:
        cand_path = args.candidates
        if not cand_path:
            die("pass --manifest or --candidates")
        cand_path = cand_path if cand_path.is_absolute() else REPO_ROOT / cand_path
        packet = load_json(cand_path)
        slug = resolve_slug(cand_path, None, packet)
        manifest_path = manifest_path_for_slug(REPO_ROOT, slug)

    if not manifest_path.is_file():
        die(f"missing manifest {manifest_path}")

    prop_path = args.proposal
    if not prop_path:
        die("--proposal required for check")
    prop_path = prop_path if prop_path.is_absolute() else REPO_ROOT / prop_path
    proposal_text = prop_path.read_text(encoding="utf-8")
    urls = urls_from_proposal(proposal_text)
    if not urls:
        die("no approved URLs found in proposal")

    manifest = load_manifest(manifest_path)
    result = validate_manifest_for_urls(manifest, urls, max_cold_age_days=args.max_age_days)
    for w in result.warnings:
        print(f"warning: {w}", file=sys.stderr)
    if result.errors:
        for e in result.errors:
            print(f"error: {e}", file=sys.stderr)
        die("share preflight failed", 1)
    print("share preflight OK for:")
    for u in urls:
        print(f"  - {u}")
    return 0


def main() -> int:
    parser = argparse.ArgumentParser(description="Perplexity share manifest for Explore further")
    sub = parser.add_subparsers(dest="command", required=True)

    p_prepare = sub.add_parser("prepare", help="Build or refresh share-manifest.json from candidates")
    p_prepare.add_argument("--candidates", type=Path, required=True)
    p_prepare.add_argument("--essay", type=Path, default=None)
    p_prepare.set_defaults(func=cmd_prepare)

    p_confirm = sub.add_parser("confirm", help="Record Anyone-with-the-link after operator shares in UI")
    p_confirm.add_argument("--candidates", type=Path, required=True)
    p_confirm.add_argument("--candidate-id", required=True, dest="candidate_id")
    p_confirm.add_argument("--url", default=None, help="Updated URL if Perplexity changed it on share")
    p_confirm.add_argument("--by", default="operator")
    p_confirm.set_defaults(func=cmd_confirm)

    p_cold = sub.add_parser("verify-cold", help="HTTP probe and/or operator incognito confirmation")
    p_cold.add_argument("--candidates", type=Path, required=True)
    p_cold.add_argument("--candidate-id", required=True, dest="candidate_id")
    p_cold.add_argument("--probe", action="store_true")
    p_cold.add_argument("--operator", action="store_true")
    p_cold.set_defaults(func=cmd_verify_cold)

    p_check = sub.add_parser("check", help="Validate manifest against approved proposal URLs")
    p_check.add_argument("--proposal", type=Path, required=True)
    p_check.add_argument("--manifest", type=Path, default=None)
    p_check.add_argument("--candidates", type=Path, default=None)
    p_check.add_argument("--max-age-days", type=int, default=30)
    p_check.set_defaults(func=cmd_check)

    args = parser.parse_args()
    return int(args.func(args))


if __name__ == "__main__":
    raise SystemExit(main())
