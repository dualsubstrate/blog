#!/usr/bin/env python3
"""List perplexity_thread URLs from a Hugo bundle for public-share verification."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path

RE_PERPLEXITY = re.compile(
    r'url:\s*"(https://www\.perplexity\.ai/search/[0-9a-f-]{36})"',
    re.I,
)
RE_LABEL = re.compile(r'label:\s*"([^"]*)"')


def main() -> int:
    parser = argparse.ArgumentParser(
        description="Print Explore further Perplexity thread URLs to verify in incognito"
    )
    parser.add_argument(
        "--post",
        required=True,
        help="Bundle path like x-minds/slug or path to index.md",
    )
    args = parser.parse_args()

    repo = Path(__file__).resolve().parents[1]
    post = args.post
    if post.endswith(".md"):
        index = repo / post if Path(post).is_absolute() else repo / post
    else:
        index = repo / "content" / post / "index.md"

    if not index.is_file():
        print(f"Missing {index}", file=sys.stderr)
        return 2

    text = index.read_text(encoding="utf-8")
    in_explore = False
    current_label = ""
    urls: list[tuple[str, str]] = []

    for line in text.splitlines():
        if re.match(r"^\s*explore:\s*$", line):
            in_explore = True
            continue
        if in_explore and line and not line[0].isspace() and not line.startswith("-"):
            break
        if not in_explore:
            continue
        m_label = RE_LABEL.search(line)
        if m_label:
            current_label = m_label.group(1)
        m_url = RE_PERPLEXITY.search(line)
        if m_url:
            urls.append((current_label or m_url.group(1), m_url.group(1)))

    if not urls:
        print("No perplexity_thread URLs found in reader_landing.explore")
        return 0

    print("Verify each URL in an incognito window (logged out of Perplexity):")
    print("Share → Anyone with the link before shipping to readers.\n")
    for label, url in urls:
        print(f"- {label}")
        print(f"  {url}\n")
    print("See .cursor/skills/site-extension-pipeline/PERPLEXITY-SHARE.md")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
