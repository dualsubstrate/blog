#!/usr/bin/env python3
"""Shared helpers for local OpenAI-compatible evaluate-only scripts."""

from __future__ import annotations

import json
import re
import sys
import urllib.error
import urllib.request
from dataclasses import dataclass
from typing import Callable

DEFAULT_BASE = "http://127.0.0.1:1320/v1"
DEFAULT_MODEL = "gemini_free/gemma-4-26b-a4b-it"
TEMPERATURE = 0.25
MAX_TOKENS_UNIT = 1536
MAX_TOKENS_WHOLE = 4096
MAX_TOKENS_SYNTHESIS = 3072


def die(msg: str, code: int = 1) -> None:
    print(msg, file=sys.stderr)
    sys.exit(code)


def chat_complete(
    base_url: str,
    model: str,
    api_key: str | None,
    system: str,
    user: str,
    *,
    max_tokens: int = MAX_TOKENS_WHOLE,
) -> str:
    url = base_url.rstrip("/") + "/chat/completions"
    payload = {
        "model": model,
        "messages": [
            {"role": "system", "content": system},
            {"role": "user", "content": user},
        ],
        "temperature": TEMPERATURE,
        "max_tokens": max_tokens,
    }
    data = json.dumps(payload).encode("utf-8")
    headers = {
        "Content-Type": "application/json",
        "Accept": "application/json",
    }
    if api_key:
        headers["Authorization"] = f"Bearer {api_key}"
    req = urllib.request.Request(url, data=data, headers=headers, method="POST")
    try:
        with urllib.request.urlopen(req, timeout=180) as resp:
            raw = resp.read().decode("utf-8")
    except urllib.error.HTTPError as e:
        body = e.read().decode("utf-8", errors="replace")
        die(f"Gateway HTTP {e.code} for model {model!r}:\n{body}")
    except urllib.error.URLError as e:
        die(f"Gateway unreachable at {url}: {e.reason}")

    try:
        obj = json.loads(raw)
    except json.JSONDecodeError:
        die(f"Non-JSON response from gateway:\n{raw[:800]}")

    if "error" in obj:
        err = obj["error"]
        if isinstance(err, dict):
            die(f"Gateway error: {err.get('message', err)}")
        die(f"Gateway error: {err}")

    choices = obj.get("choices")
    if not choices:
        die(f"No choices in response:\n{raw[:800]}")

    msg = choices[0].get("message") or {}
    content = msg.get("content")
    if not content:
        die(f"Empty message content in response:\n{raw[:800]}")
    return str(content).strip()


@dataclass
class Unit:
    """One analysis unit (list field, paragraph, or heading+paragraph block)."""

    label: str
    text: str


def split_paragraphs(text: str) -> list[str]:
    """Split on blank lines; merge a lone markdown heading into the next block."""
    raw = re.split(r"\n\s*\n+", text.strip())
    chunks = [c.strip() for c in raw if c.strip()]
    units: list[str] = []
    i = 0
    while i < len(chunks):
        c = chunks[i]
        if re.fullmatch(r"[-*_`]{3,}", c):
            i += 1
            continue
        # Pure heading line(s): fold into following paragraph when present
        if re.fullmatch(r"(?:#{1,6}\s+.+\n?)+", c) and i + 1 < len(chunks):
            units.append(c + "\n\n" + chunks[i + 1])
            i += 2
            continue
        if re.fullmatch(r"#{1,6}\s+.+", c) and i + 1 < len(chunks):
            units.append(c + "\n\n" + chunks[i + 1])
            i += 2
            continue
        units.append(c)
        i += 1
    return units


def progressive_system_addon_en() -> str:
    return """
### Progressive analysis mode

You receive **PRIOR units** (already published earlier in the piece) and one **CURRENT unit**.

Rules for each step:
- Judge **CURRENT only** for scores, quotes, and fixes.
- Use PRIOR only as continuity context: setup/payoff, voice drift, echo redundancy, bridges that depend on earlier beats.
- If CURRENT is fine alone but harms flow against PRIOR, say so under **Continuity**.
- Exception: if CURRENT only relabels PRIOR, restates it as a tautology, or restates its implication with a new metaphor pair (**Metaphor-shell restack**), that is a **CURRENT Failures-table row**, not a Continuity footnote. Do not Keep that quote.
- Exception: if CURRENT's only payload is a sales, factory, driving, or hardening verb, a fake scene (`the same room`), or "name the box… and the box" (**Industry-verb shells**), that is a **CURRENT Failures-table row**. Do not Keep `"identity feels solid because the brain keeps assembling it."`
- Do **not** re-audit PRIOR in full; at most one short continuity note per step unless the exception above applies.
- Diagnos only; no full rewrite of the unit.
"""


def progressive_system_addon_es() -> str:
    return """
### Modo de análisis progresivo

Recibes **UNIDADES PREVIAS** (ya leídas antes en la pieza) y una **UNIDAD ACTUAL**.

Reglas por paso:
- Juzga **solo la ACTUAL** para scores, citas y fixes.
- Usa las PREVIAS solo como contexto de continuidad: pago de setups, deriva de voz, eco redundante, puentes.
- Si la ACTUAL funciona sola pero rompe el flujo frente a PREVIAS, dilo en **Continuidad**.
- **No** reevalúes las PREVIAS por completo; como mucho una nota breve de continuidad.
- Solo diagnostica; no reescribas la unidad entera.
"""


def format_prior_block(prior: list[Unit]) -> str:
    if not prior:
        return "(none; this is the first unit)\n"
    parts = []
    for u in prior:
        parts.append(f"#### {u.label}\n\n{u.text}\n")
    return "\n".join(parts)


def run_progressive(
    *,
    units: list[Unit],
    system: str,
    build_unit_user: Callable[[int, int, list[Unit], Unit], str],
    build_synthesis_user: Callable[[list[tuple[Unit, str]]], str],
    base_url: str,
    model: str,
    api_key: str | None,
    language: str = "en",
) -> str:
    """Run one API call per unit, then a synthesis call. Returns full report text."""
    if not units:
        die("No units to evaluate (empty text after split).")

    if language == "es":
        system = system + "\n" + progressive_system_addon_es()
    else:
        system = system + "\n" + progressive_system_addon_en()

    n = len(units)
    sections: list[str] = []
    unit_notes: list[tuple[Unit, str]] = []
    prior: list[Unit] = []

    for i, unit in enumerate(units):
        step = i + 1
        print(f"Evaluating unit {step}/{n}: {unit.label}…", file=sys.stderr)
        user = build_unit_user(step, n, prior, unit)
        note = chat_complete(
            base_url,
            model,
            api_key,
            system,
            user,
            max_tokens=MAX_TOKENS_UNIT,
        )
        header = f"## Unit {step}/{n}: {unit.label}\n\n"
        sections.append(header + note)
        unit_notes.append((unit, note))
        prior.append(unit)

    print("Synthesizing overall report…", file=sys.stderr)
    synth_user = build_synthesis_user(unit_notes)
    synth = chat_complete(
        base_url,
        model,
        api_key,
        system,
        synth_user,
        max_tokens=MAX_TOKENS_SYNTHESIS,
    )
    sections.append("## Overall synthesis\n\n" + synth)
    return "\n\n".join(sections) + "\n"
