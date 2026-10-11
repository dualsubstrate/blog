#!/usr/bin/env python3
"""Prepare Perplexity export paths and MCP manifest for explore_proposal runner.

The explore_proposal runner (make explore-proposal) reads research from
research_export files. This script syncs those paths into candidates.json and
emits a manifest for the Cursor agent to call user-perplexity-browser
perplexity_export (thread_id + save_dir) before re-running the proposal.

Agent loop (one attempt per candidate, no export retry loops):
  For each row in export-manifest.json:
    perplexity_export thread_id=<thread_id> save_dir=<repo>/<save_dir> format=markdown
  Then copy or note the saved file at research_export path.
  Optional: --record export-results.json to patch candidates after MCP.
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path

SCRIPT_DIR = Path(__file__).resolve().parent
if str(SCRIPT_DIR) not in sys.path:
    sys.path.insert(0, str(SCRIPT_DIR))

from explore_pipeline_common import (  # noqa: E402
    bundle_slug_from_essay,
    export_file_ready,
    exports_dir,
    load_json,
    repo_root_from_here,
    sync_export_paths,
    write_json,
)

REPO_ROOT = repo_root_from_here(Path(__file__))


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Sync research_export paths and build Perplexity MCP export manifest"
    )
    parser.add_argument("--candidates", type=Path, required=True)
    parser.add_argument(
        "--essay",
        type=Path,
        default=None,
        help="English index.md (for bundle slug; default from candidates essay_path)",
    )
    parser.add_argument(
        "--sync",
        action="store_true",
        help="Write updated candidates.json with research_export paths",
    )
    parser.add_argument(
        "--record",
        type=Path,
        default=None,
        help="JSON list of {candidate_id, export_path} after MCP exports",
    )
    args = parser.parse_args()

    cand_path = args.candidates if args.candidates.is_absolute() else REPO_ROOT / args.candidates
    packet = load_json(cand_path)
    essay_path = args.essay
    if essay_path is None and packet.get("essay_path"):
        essay_path = Path(packet["essay_path"])
    if essay_path is None:
        print("error: pass --essay or essay_path in candidates JSON", file=sys.stderr)
        return 2
    if not essay_path.is_absolute():
        essay_path = REPO_ROOT / essay_path

    slug = bundle_slug_from_essay(essay_path)
    manifest = sync_export_paths(packet, slug)
    exports_dir(REPO_ROOT, slug).mkdir(parents=True, exist_ok=True)

    if args.record:
        rec_path = args.record if args.record.is_absolute() else REPO_ROOT / args.record
        records = load_json(rec_path)
        by_id = {r["candidate_id"]: r for r in records if r.get("candidate_id")}
        for c in packet.get("candidates") or []:
            cid = c.get("id")
            if cid in by_id:
                rel = by_id[cid].get("export_path") or by_id[cid].get("research_export")
                if rel:
                    c["research_export"] = rel
        if args.sync:
            cand_path.write_text(json.dumps(packet, indent=2) + "\n", encoding="utf-8")

    if args.sync and not args.record:
        cand_path.write_text(json.dumps(packet, indent=2) + "\n", encoding="utf-8")

    manifest_path = REPO_ROOT / "tmp" / "explore-proposals" / slug / "export-manifest.json"
    pending = []
    for row in manifest:
        rel = row["research_export"]
        row["export_ready"] = export_file_ready(REPO_ROOT, rel)
        if not row["export_ready"]:
            pending.append(row)
    write_json(manifest_path, {"slug": slug, "candidates_file": str(cand_path), "rows": manifest})

    print(json.dumps({"manifest": str(manifest_path), "pending": len(pending)}, indent=2))
    if pending:
        print(
            "\nMCP: For each pending row, call perplexity_export with thread_id and "
            f"save_dir under {REPO_ROOT}. Then ensure file exists at research_export.\n",
            file=sys.stderr,
        )
        for row in pending:
            print(
                f"  - {row['candidate_id']}: thread_id={row['thread_id']} -> {row['research_export']}",
                file=sys.stderr,
            )
        return 2
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
