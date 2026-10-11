"""Perplexity public-share manifest and cold-reader verification for Explore further."""

from __future__ import annotations

import json
import re
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from datetime import datetime, timedelta, timezone
from pathlib import Path
from typing import Any

from explore_pipeline_common import (
    RE_PERPLEXITY_SEARCH,
    thread_id_from_url,
    write_json,
)

MANIFEST_VERSION = 1
SHARE_SETTING_ANYONE = "anyone_with_link"
DEFAULT_COLD_MAX_AGE_DAYS = 30

LOGIN_WALL_PATTERNS = re.compile(
    r"(sign\s*in|log\s*in|create\s+an\s+account|only\s+people\s+with\s+access)",
    re.I,
)


def utc_now() -> datetime:
    return datetime.now(timezone.utc)


def isoformat_z(dt: datetime) -> str:
    if dt.tzinfo is None:
        dt = dt.replace(tzinfo=timezone.utc)
    return dt.astimezone(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def parse_iso_z(value: str) -> datetime | None:
    if not value or not str(value).strip():
        return None
    text = str(value).strip()
    if text.endswith("Z"):
        text = text[:-1] + "+00:00"
    try:
        return datetime.fromisoformat(text)
    except ValueError:
        return None


def canonical_perplexity_url(url: str) -> str:
    tid = thread_id_from_url(url)
    if not tid:
        return url.strip()
    return f"https://www.perplexity.ai/search/{tid}"


def manifest_path_for_slug(repo_root: Path, slug: str) -> Path:
    return repo_root / "tmp" / "explore-proposals" / slug / "share-manifest.json"


def empty_manifest(slug: str, *, candidates_file: str | None = None) -> dict[str, Any]:
    return {
        "manifest_version": MANIFEST_VERSION,
        "slug": slug,
        "candidates_file": candidates_file,
        "updated_at": isoformat_z(utc_now()),
        "rows": [],
    }


def prepare_rows_from_packet(packet: dict[str, Any]) -> list[dict[str, Any]]:
    rows: list[dict[str, Any]] = []
    for c in packet.get("candidates") or []:
        cid = c.get("id")
        url = c.get("url")
        if not cid or not url:
            continue
        tid = thread_id_from_url(str(url))
        if not tid:
            continue
        rows.append(
            {
                "candidate_id": str(cid),
                "url": canonical_perplexity_url(str(url)),
                "thread_id": tid,
                "share_setting": None,
                "share_confirmed_at": None,
                "share_confirmed_by": None,
                "cold_verified": False,
                "cold_verified_at": None,
                "cold_verification_method": None,
                "cold_http_status": None,
                "cold_probe_notes": None,
            }
        )
    return rows


def merge_manifest_rows(
    existing: dict[str, Any],
    new_rows: list[dict[str, Any]],
) -> dict[str, Any]:
    by_id = {str(r["candidate_id"]): r for r in existing.get("rows") or []}
    merged: list[dict[str, Any]] = []
    for nrow in new_rows:
        cid = str(nrow["candidate_id"])
        old = by_id.get(cid)
        if not old:
            merged.append(nrow)
            continue
        url = nrow["url"]
        if old.get("url") and canonical_perplexity_url(str(old["url"])) != url:
            old = dict(nrow)
            old["share_setting"] = None
            old["share_confirmed_at"] = None
            old["cold_verified"] = False
            old["cold_verified_at"] = None
        else:
            old = dict(old)
            old["url"] = url
            old["thread_id"] = nrow["thread_id"]
        merged.append(old)
    for cid, row in by_id.items():
        if not any(str(m["candidate_id"]) == cid for m in merged):
            merged.append(row)
    existing["rows"] = merged
    existing["updated_at"] = isoformat_z(utc_now())
    return existing


def load_manifest(path: Path) -> dict[str, Any]:
    return json.loads(path.read_text(encoding="utf-8"))


def save_manifest(path: Path, manifest: dict[str, Any]) -> None:
    manifest["updated_at"] = isoformat_z(utc_now())
    write_json(path, manifest)


def row_for_thread_id(manifest: dict[str, Any], thread_id: str) -> dict[str, Any] | None:
    for row in manifest.get("rows") or []:
        if str(row.get("thread_id")) == thread_id:
            return row
    return None


def row_for_url(manifest: dict[str, Any], url: str) -> dict[str, Any] | None:
    tid = thread_id_from_url(url)
    if not tid:
        return None
    return row_for_thread_id(manifest, tid)


@dataclass
class SharePreflightResult:
    ok: bool
    errors: list[str] = field(default_factory=list)
    warnings: list[str] = field(default_factory=list)


def validate_row(
    row: dict[str, Any],
    *,
    expected_url: str | None = None,
    now: datetime | None = None,
    max_cold_age_days: int = DEFAULT_COLD_MAX_AGE_DAYS,
) -> list[str]:
    errors: list[str] = []
    cid = row.get("candidate_id", "?")
    url = row.get("url")
    if not url:
        errors.append(f"{cid}: missing url in manifest row")
        return errors
    if expected_url:
        if thread_id_from_url(str(url)) != thread_id_from_url(expected_url):
            errors.append(
                f"{cid}: manifest thread does not match required URL "
                f"(manifest={url}, required={expected_url})"
            )
    if row.get("share_setting") != SHARE_SETTING_ANYONE:
        errors.append(
            f"{cid}: share not confirmed as anyone_with_link "
            f"(run: make explore-share-confirm CANDIDATES=... CANDIDATE_ID={cid})"
        )
    if not row.get("share_confirmed_at"):
        errors.append(f"{cid}: missing share_confirmed_at")
    if not row.get("cold_verified"):
        errors.append(
            f"{cid}: cold_verified is false "
            f"(run: make explore-share-verify-cold CANDIDATES=... CANDIDATE_ID={cid} PROBE=1 "
            f"or OPERATOR=1 after incognito check)"
        )
        return errors
    verified_at = parse_iso_z(str(row.get("cold_verified_at") or ""))
    if not verified_at:
        errors.append(f"{cid}: cold_verified true but cold_verified_at missing")
        return errors
    now = now or utc_now()
    if verified_at > now + timedelta(minutes=5):
        errors.append(f"{cid}: cold_verified_at is in the future")
    age = now - verified_at
    if age > timedelta(days=max_cold_age_days):
        errors.append(
            f"{cid}: cold verification stale ({age.days} days old; re-run verify-cold)"
        )
    return errors


def validate_manifest_for_urls(
    manifest: dict[str, Any],
    required_urls: list[str],
    *,
    max_cold_age_days: int = DEFAULT_COLD_MAX_AGE_DAYS,
    now: datetime | None = None,
) -> SharePreflightResult:
    result = SharePreflightResult(ok=True)
    if manifest.get("manifest_version") != MANIFEST_VERSION:
        result.warnings.append(
            f"manifest_version {manifest.get('manifest_version')} != {MANIFEST_VERSION}"
        )
    seen: set[str] = set()
    for req in required_urls:
        tid = thread_id_from_url(req)
        if not tid:
            result.errors.append(f"invalid Perplexity URL: {req}")
            result.ok = False
            continue
        if tid in seen:
            continue
        seen.add(tid)
        row = row_for_thread_id(manifest, tid)
        if not row:
            result.errors.append(
                f"no manifest row for thread {tid} ({req}); "
                f"run: make explore-share-prepare CANDIDATES=..."
            )
            result.ok = False
            continue
        row_errors = validate_row(
            row,
            expected_url=req,
            now=now,
            max_cold_age_days=max_cold_age_days,
        )
        for err in row_errors:
            result.errors.append(err)
            result.ok = False
    return result


def cold_probe_url(
    url: str,
    *,
    timeout: float = 25.0,
    user_agent: str = "BehaviorEngineering-ExploreShareCheck/1.0",
) -> dict[str, Any]:
    """Best-effort unauthenticated fetch; may false-fail on bot blocking."""
    canonical = canonical_perplexity_url(url)
    req = urllib.request.Request(
        canonical,
        headers={
            "User-Agent": user_agent,
            "Accept": "text/html,application/xhtml+xml",
        },
        method="GET",
    )
    out: dict[str, Any] = {
        "url": canonical,
        "ok": False,
        "http_status": None,
        "notes": "",
    }
    try:
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            status = resp.status
            body = resp.read(500_000).decode("utf-8", errors="replace")
    except urllib.error.HTTPError as e:
        status = e.code
        body = e.read(200_000).decode("utf-8", errors="replace")
    except urllib.error.URLError as e:
        out["notes"] = f"network error: {e.reason}"
        return out

    out["http_status"] = status
    text_sample = re.sub(r"\s+", " ", body[:8000])
    if status >= 400:
        out["notes"] = f"HTTP {status}"
        return out
    if len(body) < 500:
        out["notes"] = "response body too small (possible login wall or block)"
        return out
    if LOGIN_WALL_PATTERNS.search(text_sample) and "perplexity" in text_sample.lower():
        if len(body) < 15_000:
            out["notes"] = "login or access wall heuristics matched"
            return out
    out["ok"] = True
    out["notes"] = "heuristic pass (still confirm in incognito when probe is ambiguous)"
    return out


def confirm_share_row(
    manifest: dict[str, Any],
    candidate_id: str,
    *,
    url: str | None = None,
    confirmed_by: str = "operator",
) -> dict[str, Any]:
    row = None
    for r in manifest.get("rows") or []:
        if str(r.get("candidate_id")) == candidate_id:
            row = r
            break
    if not row:
        raise KeyError(f"candidate_id not in manifest: {candidate_id}")
    if url:
        row["url"] = canonical_perplexity_url(url)
        row["thread_id"] = thread_id_from_url(url)
    row["share_setting"] = SHARE_SETTING_ANYONE
    row["share_confirmed_at"] = isoformat_z(utc_now())
    row["share_confirmed_by"] = confirmed_by
    return manifest


def record_cold_verification(
    manifest: dict[str, Any],
    candidate_id: str,
    *,
    method: str,
    probe: dict[str, Any] | None = None,
    operator_confirmed: bool = False,
) -> dict[str, Any]:
    row = None
    for r in manifest.get("rows") or []:
        if str(r.get("candidate_id")) == candidate_id:
            row = r
            break
    if not row:
        raise KeyError(f"candidate_id not in manifest: {candidate_id}")

    if method == "operator_incognito":
        if not operator_confirmed:
            raise ValueError("operator_incognito requires operator_confirmed=True")
        row["cold_verified"] = True
        row["cold_verified_at"] = isoformat_z(utc_now())
        row["cold_verification_method"] = method
        row["cold_probe_notes"] = "operator incognito / logged-out browser"
        return manifest

    if method == "http_probe":
        if not probe:
            raise ValueError("http_probe requires probe result")
        row["cold_http_status"] = probe.get("http_status")
        row["cold_probe_notes"] = probe.get("notes")
        row["cold_verification_method"] = method
        if probe.get("ok"):
            row["cold_verified"] = True
            row["cold_verified_at"] = isoformat_z(utc_now())
        else:
            row["cold_verified"] = False
            row["cold_verified_at"] = None
        return manifest

    raise ValueError(f"unknown verification method: {method}")
