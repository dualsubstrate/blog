#!/usr/bin/env python3
"""Wrapper: Explore further proposal runner (see site-extension-pipeline skill)."""
import runpy
from pathlib import Path

TARGET = (
    Path(__file__).resolve().parent.parent
    / ".cursor/skills/site-extension-pipeline/scripts/explore_proposal.py"
)
runpy.run_path(str(TARGET), run_name="__main__")
