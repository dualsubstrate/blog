#!/usr/bin/env python3
import runpy
from pathlib import Path

runpy.run_path(
    str(
        Path(__file__).resolve().parent.parent
        / ".cursor/skills/site-extension-pipeline/scripts/explore_fetch_exports.py"
    ),
    run_name="__main__",
)
