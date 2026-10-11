#!/usr/bin/env python3
"""Generate psych-fitness-28 Hugo bundles for days 4-28 from tmp/articles/petition/campaign/*-first.txt."""

from __future__ import annotations

import re
import textwrap
from datetime import date, timedelta
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
CAMP = ROOT / "tmp/articles/petition/campaign"
OUT = ROOT / "content/cognitive-memetics/psych-fitness-28"

FLUFF_BODY_EN = """We all run on the **same human hardware**. Let's build the support and checks to match the pressure of the job.

Sign the petition:
👉 Petition (EN9806): Stronger Checks and Balances: Psychological Fitness for Australia's Top Leaders
https://www.aph.gov.au/e-petitions/petition/EN9806"""

FLUFF_BODY_ES = """Todos funcionamos con el **mismo hardware humano**. Construyamos el apoyo y los controles para igualar la presion del cargo.

Firma la peticion:
👉 Peticion (EN9806): Verificaciones y Equilibrios Mas Fuertes: Aptitud Psicologica para los Lideres de Australia
https://www.aph.gov.au/e-petitions/petition/EN9806"""

# Spanish titles (episode line) and one-line descriptions for cards; TLDR left as draft stub.
ES_META: dict[int, tuple[str, str]] = {
    4: (
        "La aptitud mental fluye hacia abajo",
        "**Dia 4 de 28**: el tono en la cima moldea instituciones; la **aptitud mental** sostiene la **estabilidad** que baja.",
    ),
    5: (
        "Las guerras empiezan en un sistema nervioso",
        "**Dia 5 de 28**: decisiones de alto riesgo empiezan en **un cuerpo**; el estandar importa antes del choque.",
    ),
    6: (
        "Pilotos y cirujanos aceptan controles",
        "**Dia 6 de 28**: roles de alto riesgo ya aceptan **controles**; el cargo publico puede alinearse con esa logica.",
    ),
    7: (
        "La presion empuja a culpar",
        "**Dia 7 de 28**: bajo estres, la culpa **externa** gana; la aptitud mental reduce ese derrotero.",
    ),
    8: (
        "Los ciudadanos tambien cuentan",
        "**Dia 8 de 28**: el publico y el personal ayudan a **sostener** a los lideres; la estructura importa.",
    ),
    9: (
        "Esperamos la crisis para actuar",
        "**Dia 9 de 28**: el sistema reacciona tarde; pedimos **chequeos** antes del dano.",
    ),
    10: (
        "Estandar quita el estigma",
        "**Dia 10 de 28**: tratar la aptitud como **norma** quita verguenza y mejora la conversacion.",
    ),
    11: (
        "El alto cargo trae presion constante",
        "**Dia 11 de 28**: la presion es **implacable**; sin soporte, la decision se estrecha.",
    ),
    12: (
        "El ejemplo de la rendicion de cuentas",
        "**Dia 12 de 28**: cuando los de arriba se someten a **controles**, modelan **responsabilidad**.",
    ),
    13: (
        "Despues de la crisis, poco se mira al lider",
        "**Dia 13 de 28**: tras el shock, rara vez auditamos al **liderazgo** que dirigio la respuesta.",
    ),
    14: (
        "Mitad del camino",
        "**Dia 14 de 28**: mitad de **28 dias**; seguir firmando y compartiendo.",
    ),
    15: (
        "El estres del lider llega al equipo",
        "**Dia 15 de 28**: el estres **irradia**; el equipo paga el precio del desorden nervioso arriba.",
    ),
    16: (
        "El estre estrecha la vision",
        "**Dia 16 de 28**: el estres **acorta** el horizonte; menos margen para tradeoffs justos.",
    ),
    17: (
        "La confianza cae con lo erratico",
        "**Dia 17 de 28**: actuar de forma **erratica** erosiona la confianza en las instituciones.",
    ),
    18: (
        "El burnout ejecutivo cuesta",
        "**Dia 18 de 28**: el agotamiento en la cima cuesta **dinero** y talento tambien fuera del sector publico.",
    ),
    19: (
        "La empatia cuesta energia mental",
        "**Dia 19 de 28**: la empatia seria requiere **capacidad** cognitiva; no es gratis.",
    ),
    20: (
        "La gobernanza moderna es enorme",
        "**Dia 20 de 28**: la **complejidad** del Estado moderno exige mentes en forma, no solo CV.",
    ),
    21: (
        "Queda una semana",
        "**Dia 21 de 28**: **ultima semana** para firmar EN9806 y dejar constancia.",
    ),
    22: (
        "Medimos la economia con numeros duros",
        "**Dia 22 de 28**: si medimos el PIB, podemos pedir **metricas** honestas de aptitud en el poder.",
    ),
    23: (
        "La politica reactiva cansa",
        "**Dia 23 de 28**: el vaiven **reactivo** agota al publico; los estandares reducen drama.",
    ),
    24: (
        "Las politicas duran mas que los lideres",
        "**Dia 24 de 28**: las leyes de hoy **perduran**; quien decide debe estar en forma.",
    ),
    25: (
        "Redefinir la fuerza politica",
        "**Dia 25 de 28**: la fuerza puede ser **claridad** y calma bajo fuego, no solo volumen.",
    ),
    26: (
        "Ultimo fin de semana para firmar",
        "**Dia 26 de 28**: **ultimo fin de semana** para sumar firmas a EN9806.",
    ),
    27: (
        "Manana cierra la peticion",
        "**Dia 27 de 28**: **ultimo dia** manana; pasa la voz hoy.",
    ),
    28: (
        "Hoy es el ultimo dia para firmar",
        "**Dia 28 de 28**: **hoy** cierra la ventana; firma si aun no lo hiciste.",
    ),
}


def slugify(s: str) -> str:
    s = s.lower().strip()
    s = re.sub(r"[^a-z0-9]+", "-", s)
    return s.strip("-")


def parse_first(path: Path) -> tuple[str, str]:
    raw = path.read_text(encoding="utf-8")
    lines = [ln.rstrip() for ln in raw.splitlines()]
    m = re.match(r"^\[Day (\d+) of 28\]\s*(.+)\.?\s*$", lines[0])
    if not m:
        raise ValueError(f"bad first line in {path}: {lines[0]!r}")
    day = int(m.group(1))
    title = m.group(2).strip().rstrip(".")
    # body until petition line
    body_lines: list[str] = []
    for ln in lines[1:]:
        if ln.strip().startswith("👉"):
            break
        body_lines.append(ln)
    while body_lines and body_lines[0].strip() == "":
        body_lines.pop(0)
    while body_lines and body_lines[-1].strip() == "":
        body_lines.pop()
    body = "\n".join(body_lines).strip()
    return title, body


def day_date(day: int) -> date:
    # Day 1 published 2026-05-02 in site bundles
    return date(2026, 5, 2) + timedelta(days=day - 1)


def en_description(day: int, title: str, tldr_first_sentence: str) -> str:
    # Short teaser: day hook + bold bits from title words (restrained)
    t = tldr_first_sentence.strip().replace("\n", " ")
    if len(t) > 160:
        t = t[:157] + "..."
    return f"**Day {day} of 28**, **{title}**: {t}"


def main() -> None:
    for day in range(4, 29):
        src = CAMP / f"{day:02d}-first.txt"
        title, body = parse_first(src)
        slug = slugify(title)
        d = day_date(day)
        date_s = d.isoformat()
        folder = f"{date_s}-day-{day:02d}-{slug}"
        bundle = OUT / folder
        bundle.mkdir(parents=True, exist_ok=True)
        feat = f"{day:02d}.webp"

        # First sentence for description (plain, we add bold in template)
        first_para = body.split("\n\n")[0].replace("\n", " ").strip()
        desc = en_description(day, title, first_para)

        tldr_yaml = textwrap.indent(body + "\n", "  ")

        en_fm = f"""---
translationKey: "{date_s}-psych-fitness-day-{day:02d}"
date: '{date_s}T01:00:00+11:00'
heading_code: D{day}
project: Psych-Fitness-28 🙏
title: {title}
type: sayings
description: {desc!r}
tldr: |
{tldr_yaml}
fluff: |
{textwrap.indent(FLUFF_BODY_EN, "  ") + chr(10)}
draft: false

featuredImage: "{feat}"
featuredImagePreview: "{feat}"

images:
  - {feat}

resources:
  - src: {feat}
    name: featured-image

tags:
  - PsychFitness
  - Leadership
  - EN9806
  - MentalFitness
  - Petition

categories: ["Cognitive-Memetics", "Psych-Fitness-28"]

aliases:
  - "/cognitive-memetics/psych-fitness-28/{folder}/"
---

<!--more-->
"""
        (bundle / "index.md").write_text(en_fm, encoding="utf-8")

        es_title, es_desc = ES_META[day]
        es_fm = f"""---
translationKey: "{date_s}-psych-fitness-day-{day:02d}"
date: '{date_s}T01:00:00+11:00'
heading_code: D{day}
project: Psique-Fitness-28 🙏
title: {es_title}
type: sayings
description: {es_desc!r}
tldr: |
  Version en espanol de este dia pendiente de edicion. Mientras tanto, lee la entrada en ingles en la misma ruta con idioma EN.
fluff: |
{textwrap.indent(FLUFF_BODY_ES, "  ") + chr(10)}
draft: false

featuredImage: "{feat}"
featuredImagePreview: "{feat}"

images:
  - {feat}

resources:
  - src: {feat}
    name: featured-image

tags:
  - PsiqueFitness
  - Liderazgo
  - EN9806
  - SaludMental
  - Peticion

categories: ["Cognitive-Memetics", "Psych-Fitness-28"]

aliases:
  - "/es/cognitive-memetics/psych-fitness-28/{folder}/"
---

<!--more-->
"""
        (bundle / "index.es.md").write_text(es_fm, encoding="utf-8")

        # linkedin stub
        en_url = f"https://blog.dualsubstrate.com/cognitive-memetics/psych-fitness-28/{folder}/"
        es_url = f"https://blog.dualsubstrate.com/es/cognitive-memetics/psych-fitness-28/{folder}/"
        li = f"""D{day}: Psych-Fitness-28 🙏
\"{title}\"

✔️ TLDR:
{body.strip()[:1200]}{"..." if len(body.strip()) > 1200 else ""}

➕ FLUFF:
We all run on the same human hardware. Let's build the support and checks to match the pressure of the job.

❓ BUT WHY:
This is part of the Psych-Fitness-28 campaign for petition EN9806: stronger checks and psychological fitness standards for Australia's top leaders.

👉 Petition (EN9806): https://www.aph.gov.au/e-petitions/petition/EN9806

#PsychFitness #Leadership #EN9806 #MentalFitness #Petition

🧷 Full post (site) →
- EN: {en_url}
- ES: {es_url}

🔗 Psych-Fitness-28 (hub) → https://blog.dualsubstrate.com/categories/psych-fitness-28/
"""
        (bundle / "linkedin.txt").write_text(li, encoding="utf-8")

    print("Wrote days 4-28 under", OUT)


if __name__ == "__main__":
    main()
