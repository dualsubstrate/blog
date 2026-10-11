"""Shared helpers for explore_proposal / explore_fetch_exports / explore_apply."""

from __future__ import annotations

import json
import re
from pathlib import Path
from typing import Any

RE_PERPLEXITY_SEARCH = re.compile(
    r"https?://(?:www\.)?perplexity\.ai/search/([0-9a-f-]{36})",
    re.I,
)


def repo_root_from_here(here: Path) -> Path:
    """Site repo root from a script file path under .cursor/skills/.../scripts/."""
    return here.resolve().parents[4]


def thread_id_from_url(url: str) -> str | None:
    m = RE_PERPLEXITY_SEARCH.search(url.strip())
    return m.group(1) if m else None


def bundle_slug_from_essay(essay_path: Path) -> str:
    return essay_path.parent.name


def exports_dir(repo_root: Path, slug: str) -> Path:
    return repo_root / "tmp" / "explore-proposals" / slug / "exports"


def default_export_rel(slug: str, candidate_id: str) -> str:
    return f"tmp/explore-proposals/{slug}/exports/{candidate_id}.md"


def sync_export_paths(
    packet: dict[str, Any],
    slug: str,
) -> list[dict[str, Any]]:
    """Ensure each candidate has research_export path; return export manifest rows."""
    manifest: list[dict[str, Any]] = []
    for c in packet.get("candidates") or []:
        cid = c.get("id")
        url = c.get("url", "")
        if not cid or not url:
            continue
        rel = c.get("research_export") or default_export_rel(slug, str(cid))
        c["research_export"] = rel
        tid = thread_id_from_url(str(url))
        manifest.append(
            {
                "candidate_id": cid,
                "url": url,
                "thread_id": tid,
                "research_export": rel,
                "save_dir": str(Path(rel).parent),
            }
        )
    return manifest


def export_file_ready(repo_root: Path, rel_path: str, *, min_bytes: int = 80) -> bool:
    p = repo_root / rel_path
    return p.is_file() and p.stat().st_size >= min_bytes


def write_json(path: Path, obj: Any) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(obj, indent=2) + "\n", encoding="utf-8")


def load_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))
