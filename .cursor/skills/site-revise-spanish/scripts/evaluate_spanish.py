#!/usr/bin/env python3
"""Evaluate-only Spanish naturalness via OpenAI-compatible local gateway.

Default: progressive unit analysis. Default model: Gemma 4.
Does not edit Hugo content. Diagnose only.
"""

from __future__ import annotations

import argparse
import os
import re
import sys
from datetime import date
from pathlib import Path
from typing import Any

_SKILLS = Path(__file__).resolve().parents[2]
if str(_SKILLS) not in sys.path:
    sys.path.insert(0, str(_SKILLS))
from site_local_eval_common.common import (  # noqa: E402
    DEFAULT_BASE,
    DEFAULT_MODEL,
    MAX_TOKENS_WHOLE,
    Unit,
    chat_complete,
    die,
    format_prior_block,
    run_progressive,
    split_paragraphs,
)

LIST_KEYS = ("title", "description", "grounding", "sowhat", "subtitle", "tldr", "fluff")

ES_SIBLINGS = (
    "index.es.md",
    "facebook-es.txt",
    "linkedin.es.txt",
    "substack.es.md",
)

SCRIPT_DIR = Path(__file__).resolve().parent
SKILL_DIR = SCRIPT_DIR.parent
PACK_PATH = SKILL_DIR / "packs" / "evaluate.md"
REFERENCE_PATH = SKILL_DIR / "reference.md"


def split_front_matter(text: str) -> tuple[dict[str, Any], str]:
    if not text.startswith("---"):
        return {}, text
    parts = text.split("---", 2)
    if len(parts) < 3:
        return {}, text
    raw_fm, body = parts[1], parts[2]
    fm: dict[str, Any] = {}
    lines = raw_fm.strip("\n").split("\n")
    i = 0
    while i < len(lines):
        line = lines[i]
        if not line.strip() or line.strip().startswith("#"):
            i += 1
            continue
        m = re.match(r"^([A-Za-z0-9_]+):\s*(.*)$", line)
        if not m:
            i += 1
            continue
        key, val = m.group(1), m.group(2)
        if val in ("|", ">"):
            block: list[str] = []
            i += 1
            while i < len(lines):
                nxt = lines[i]
                if nxt and not nxt[0].isspace() and re.match(
                    r"^[A-Za-z0-9_]+:", nxt
                ):
                    break
                block.append(nxt)
                i += 1
            cleaned = "\n".join(block).strip("\n")
            if cleaned:
                indents = [
                    len(x) - len(x.lstrip(" "))
                    for x in cleaned.split("\n")
                    if x.strip()
                ]
                if indents:
                    n = min(indents)
                    cleaned = "\n".join(
                        x[n:] if len(x) >= n else x for x in cleaned.split("\n")
                    )
            fm[key] = cleaned.strip()
            continue
        val = val.strip()
        if (val.startswith('"') and val.endswith('"')) or (
            val.startswith("'") and val.endswith("'")
        ):
            val = val[1:-1]
        fm[key] = val
        i += 1
    return fm, body.lstrip("\n")


def resolve_bundle(path: Path) -> Path:
    path = path.resolve()
    if path.is_dir():
        return path
    if path.is_file():
        return path.parent
    die(f"Path not found: {path}")
    return path


def collect_es_files(bundle: Path, scope: str) -> list[Path]:
    files: list[Path] = []
    if scope in ("full", "site", "body", "list"):
        site = bundle / "index.es.md"
        if site.is_file():
            files.append(site)
        elif scope != "full":
            die(f"Missing index.es.md in {bundle}")
    if scope == "full":
        for name in ES_SIBLINGS[1:]:
            p = bundle / name
            if p.is_file():
                files.append(p)
    elif scope == "facebook":
        p = bundle / "facebook-es.txt"
        if not p.is_file():
            die(f"Missing facebook-es.txt in {bundle}")
        files = [p]
    elif scope == "linkedin":
        p = bundle / "linkedin.es.txt"
        if not p.is_file():
            die(f"Missing linkedin.es.txt in {bundle}")
        files = [p]
    elif scope == "substack":
        p = bundle / "substack.es.md"
        if not p.is_file():
            die(f"Missing substack.es.md in {bundle}")
        files = [p]

    if not files:
        die(
            f"No Spanish files found in {bundle}. "
            f"Expected at least index.es.md (scope={scope})."
        )
    return files


def format_md_file(path: Path, scope: str) -> str:
    text = path.read_text(encoding="utf-8")
    name = path.name
    if name.endswith(".txt") or not text.startswith("---"):
        return f"### File: {name}\n\n{text.strip()}\n"

    fm, body = split_front_matter(text)
    title = str(fm.get("title") or path.stem)
    parts = [f"### File: {name}", f"title: {title}"]
    if scope in ("full", "site", "list"):
        for k in LIST_KEYS:
            if k == "title":
                continue
            if k in fm and str(fm[k]).strip():
                parts.append(f"{k}:\n{fm[k]}")
    if scope in ("full", "site", "body"):
        parts.append("## Body\n")
        parts.append(body.strip())
    return "\n\n".join(parts) + "\n"


def build_units_from_files(files: list[Path], scope: str) -> list[Unit]:
    units: list[Unit] = []
    for path in files:
        text = path.read_text(encoding="utf-8")
        name = path.name
        if name.endswith(".txt") or not text.startswith("---"):
            for i, para in enumerate(split_paragraphs(text), start=1):
                units.append(Unit(f"{name} paragraph {i}", para))
            continue

        fm, body = split_front_matter(text)
        md_scope = (
            scope if scope in ("full", "site", "list", "body") else "site"
        )
        title = str(fm.get("title") or "").strip()
        if md_scope in ("full", "site", "list") and title:
            units.append(Unit(f"{name}:title", title))
        if md_scope in ("full", "site", "list"):
            for k in LIST_KEYS:
                if k == "title":
                    continue
                if k in fm and str(fm[k]).strip():
                    units.append(Unit(f"{name}:{k}", str(fm[k]).strip()))
        if md_scope in ("full", "site", "body"):
            for i, para in enumerate(split_paragraphs(body), start=1):
                units.append(Unit(f"{name} body paragraph {i}", para))
    return units


def build_en_context(bundle: Path) -> str:
    en = bundle / "index.md"
    if not en.is_file():
        return "(no index.md present)\n"
    fm, body = split_front_matter(en.read_text(encoding="utf-8"))
    lines = [
        "Use only for meaning gaps. Do not mirror English syntax or ### titles.",
        f"title: {fm.get('title', '')}",
    ]
    for k in ("description", "sowhat", "subtitle"):
        if k in fm and str(fm[k]).strip():
            lines.append(f"{k}:\n{fm[k]}")
    lead = body.strip()[:400]
    if lead:
        lines.append(f"body lead (truncated):\n{lead}…")
    return "\n".join(lines) + "\n"


def load_pack_system() -> str:
    if not PACK_PATH.is_file():
        die(f"Missing pack: {PACK_PATH}")
    text = PACK_PATH.read_text(encoding="utf-8")
    cut = text.find("## User message template")
    system_from_pack = text[:cut].strip() if cut != -1 else text.strip()
    ref = ""
    if REFERENCE_PATH.is_file():
        ref = "\n\n---\n\n# Rubric detail\n\n" + REFERENCE_PATH.read_text(
            encoding="utf-8"
        )
    return system_from_pack + ref


def build_user_message(
    paths: list[Path],
    scope: str,
    excerpt: str,
    en_context: str,
) -> str:
    path_list = ", ".join(str(p) for p in paths)
    return f"""## Target

Bundle / paths: {path_list}
Scope: {scope}

## Spanish text to evaluate

{excerpt}

## English meaning context only (optional; do not mirror syntax)

{en_context}

## Output format (required)

Responde en español (quotes en español). Estructura:

1. **Veredicto** (1 línea: nativo / a medias / claramente traducido) + 1 párrafo corto
2. **Scores** 0–2 con una línea de por qué: native_naturalness | anti_calque | read_aloud | collocation_grammar | cross_file_coherence | conceptual_fidelity
3. **Suena nativo** (2–5 citas que funcionan)
4. **Huele a traducción / diluye idea** tabla: Cita exacta | Archivo | Patrón | Por qué falla | Dirección de fix (ES, no reescritura completa)
5. **Top 5 fixes** ordenados (si hay dilución conceptual, priorízala sobre estilo)
6. **Sincronía entre archivos** (solo si hay 2+ archivos ES; si no: "n/a")
7. **Ready / Almost / Re-adapt section** con una línea
"""


def build_unit_user(paths: list[Path], scope: str, en_context: str) -> Any:
    path_list = ", ".join(str(p) for p in paths)

    def _fn(step: int, n: int, prior: list[Unit], unit: Unit) -> str:
        return f"""## Target

Paths: {path_list}
Scope: {scope}
Paso progresivo: {step} de {n}
Etiqueta de unidad: {unit.label}

## Contexto de significado EN (opcional; no copies sintaxis)

{en_context}

## UNIDADES PREVIAS (solo contexto; no re-scores a fondo)

{format_prior_block(prior)}

## UNIDAD ACTUAL (evalúa esto)

{unit.text}

## Formato de salida (solo esta unidad)

1. **Veredicto de unidad** (1–2 frases)
2. **Scores** 0–2 solo de la ACTUAL: native_naturalness | anti_calque | read_aloud | collocation_grammar | cross_file_coherence | conceptual_fidelity (si no aplica cross-file con PREVIO, 2)
3. **Suena nativo** (0–3 citas de la ACTUAL)
4. **Huele a traducción / diluye idea** de la ACTUAL: Cita | Patrón | Por qué | Fix (ES)
5. **Continuidad** vs PREVIAS (eco, setup sin pago, deriva) o "ninguna"
6. **Top fixes de esta unidad** (máx. 3; prioriza pérdida de contraste o de objeto)
"""

    return _fn


def build_synthesis_user(paths: list[Path], scope: str) -> Any:
    path_list = ", ".join(str(p) for p in paths)

    def _fn(unit_notes: list[tuple[Unit, str]]) -> str:
        digest_parts = []
        for unit, note in unit_notes:
            digest_parts.append(f"### {unit.label}\n\n{note}\n")
        digest = "\n".join(digest_parts)
        return f"""## Target

Paths: {path_list}
Scope: {scope}

Ya juzgaste cada unidad en orden. Abajo están las notas por unidad.
Sintetiza un informe **de toda la pieza**. Prioriza patrones que se repiten.
No inventes citas que no estén en las notas o en las etiquetas de unidad.

## Notas por unidad

{digest}

## Formato de salida (global)

1. **Veredicto** (nativo / a medias / claramente traducido) + párrafo corto
2. **Scores** 0–2 globales: native_naturalness | anti_calque | read_aloud | collocation_grammar | cross_file_coherence | conceptual_fidelity
3. **Suena nativo** (2–5 mejores citas de las notas)
4. **Huele a traducción / diluye idea** consolidado (cita | archivo/unidad | patrón | fix)
5. **Top 5 fixes** ordenados (dilución conceptual antes que estilo)
6. **Sincronía entre archivos** o "n/a"
7. **Ready / Almost / Re-adapt section**
8. **Mapa de progreso**: etiqueta de unidad → arrastre / impulso / nota
"""

    return _fn


def write_out(
    out_path: Path,
    paths: list[Path],
    model: str,
    scope: str,
    mode: str,
    report: str,
) -> None:
    out_path.parent.mkdir(parents=True, exist_ok=True)
    path_list = ", ".join(f"`{p}`" for p in paths)
    header = f"""# Local Spanish prose review

- **date:** {date.today().isoformat()}
- **paths:** {path_list}
- **scope:** {scope}
- **mode:** {mode}
- **model:** `{model}`

---

"""
    out_path.write_text(header + report + "\n", encoding="utf-8")
    print(f"Wrote {out_path}", file=sys.stderr)


def main() -> None:
    p = argparse.ArgumentParser(
        description=(
            "Evaluate-only Spanish naturalness via local LLM (no content edits)."
        )
    )
    p.add_argument(
        "path",
        type=Path,
        help="Path to index.es.md, a sidecar, or the page bundle directory",
    )
    p.add_argument(
        "--scope",
        choices=(
            "full",
            "site",
            "list",
            "body",
            "facebook",
            "linkedin",
            "substack",
        ),
        default="full",
        help="full=all ES files; site/list/body=index.es.md only; or one sidecar",
    )
    p.add_argument(
        "--mode",
        choices=("progressive", "whole"),
        default="progressive",
        help="progressive=one unit at a time with prior context (default); whole=single pass",
    )
    p.add_argument(
        "--model",
        default=os.environ.get("LOCAL_LLM_MODEL", DEFAULT_MODEL),
        help=f"Model id (default: {DEFAULT_MODEL})",
    )
    p.add_argument(
        "--base-url",
        default=os.environ.get("LOCAL_LLM_BASE_URL", DEFAULT_BASE),
        help=f"OpenAI-compatible base URL (default: {DEFAULT_BASE})",
    )
    p.add_argument(
        "--out",
        type=Path,
        default=None,
        help="Optional path to write report markdown",
    )
    p.add_argument(
        "--no-en-context",
        action="store_true",
        help="Omit English index.md meaning context",
    )
    args = p.parse_args()

    path: Path = args.path
    if not path.exists():
        die(f"Path not found: {path}")

    if path.is_file() and path.name in (
        "facebook-es.txt",
        "linkedin.es.txt",
        "substack.es.md",
    ) and args.scope == "full":
        scope = {
            "facebook-es.txt": "facebook",
            "linkedin.es.txt": "linkedin",
            "substack.es.md": "substack",
        }[path.name]
    else:
        scope = args.scope

    bundle = resolve_bundle(path)
    files = collect_es_files(bundle, scope)
    en_context = (
        "(omitted)\n" if args.no_en_context else build_en_context(bundle)
    )
    system = load_pack_system()
    api_key = os.environ.get("LOCAL_LLM_API_KEY") or None

    if args.mode == "whole":
        excerpts: list[str] = []
        for f in files:
            if f.name == "index.es.md":
                excerpts.append(
                    format_md_file(
                        f,
                        scope
                        if scope in ("full", "site", "list", "body")
                        else "site",
                    )
                )
            else:
                excerpts.append(format_md_file(f, "full"))
        excerpt = "\n".join(excerpts)
        if not excerpt.strip():
            die("Empty excerpt; nothing to evaluate.")
        user = build_user_message(files, scope, excerpt, en_context)
        report = chat_complete(
            args.base_url,
            args.model,
            api_key,
            system,
            user,
            max_tokens=MAX_TOKENS_WHOLE,
        )
    else:
        units = build_units_from_files(files, scope)
        if not units:
            die("No units to evaluate after split.")
        report = run_progressive(
            units=units,
            system=system,
            build_unit_user=build_unit_user(files, scope, en_context),
            build_synthesis_user=build_synthesis_user(files, scope),
            base_url=args.base_url,
            model=args.model,
            api_key=api_key,
            language="es",
        )

    print(report)

    if args.out is not None:
        write_out(args.out, files, args.model, scope, args.mode, report)


if __name__ == "__main__":
    main()
