---
name: site-spanish-translation-content
description: >-
  Adds or edits Spanish Hugo pages as siblings of English posts (`*.es.md`), sets
  `translationKey`, mirrors front matter per content `type`, keeps taxonomies aligned,
  and applies type skills (claims, video, cognitive-memetics) in Spanish. After the first
  `index.es.md` draft, MUST run local Gemma 4 (`evaluate_spanish.py`) for native-Spanish /
  calque check. Documents Spanish permalinks, internal links, `aliases` for
  `social-protocols`, and `facebook-es.txt` / `linkedin.es.txt` alignment with
  `hugo list all`. Prefers idiomatic Spanish over cognate calques; five-directive native
  opinion-editor audit; hybrid cadence; hard native-Spanish constraints (conceptual
  anglicisms, English skeleton under Spanish words, **clause rebuild** not word-swap,
  reverse-translation test). Use when translating to Spanish, adding `.es.md`, gemma4 /
  Gemma 4 Spanish gate, bilingual posts, Spanish routes, or idiomatic Spanish (not
  English-shaped prose or EN clause order with ES words).
---

# Spanish translation content (Hugo `es`)

## When this applies

- The site is multilingual (`en` default, `es` under `/es/`). English posts live as `content/<section>/<slug>/index.md` or `*.md` without a language suffix.
- Spanish versions use the **same bundle path** with a **language suffix**: `index.es.md` next to `index.md`, or `postname.es.md` next to `postname.md` (see [Hugo multilingual content](https://gohugo.io/content-management/multilingual/#translate-your-content)).

## Workflow

1. **Locate the English source** in `content/<section>/...`. Do not move the section; translation is a **sibling file** in the same folder.
2. **Create** `index.es.md` (or `*.es.md`) with the same `type`, **`date`** (identical string to English; site default **`date: 'YYYY-MM-DDT01:00:00+11:00'`** per **`.cursor/rules/site-content-markdown-writing.mdc`** → **Publish `date`**), `slug` (if set on EN), images paths, and structural front matter as the English file.
3. **Set `translationKey`** to the same stable string on **both** files (for example `translationKey: "2026-04-07-yt-history-of-intelligence"`). Use one identifier per logical article, not per language.
4. **Adapt, do not gloss:** From the English **meaning** (who does what, when, stakes), write Spanish a native editor would say. MUST rebuild clause order per **Clause rebuild** below. MUST NOT substitute Spanish words into the English sentence (same relatives, *thus/making*, *You hear it when*). Keep **proper nouns**, paper titles, and product names as in the source unless a standard Spanish name exists.
5. **Apply the right type skill** for structure and field roles (Claim vs Grounding, video `description`, sayings `tldr`/`fluff`, etc.): **`.cursor/skills/site-claims-content/SKILL.md`**, **video-content**, **cognitive-memetics-content**—and for Markdown **`**bold**`** density in any translated field, **`.cursor/skills/site-revise-emphasis/SKILL.md`**—but **Spanish prose** inside those roles.
6. **Gemma 4 native-Spanish gate (MUST):** After the first **`index.es.md`** draft, run local **Gemma 4** (see **Gemma 4 native-Spanish gate** below). Apply calque / English-skeleton / meaning-load fixes from the report, then re-check **Hard constraints**. Skip this call for **line-by-line author edits** (step 8).
7. **Index first, sidecars later:** Finish and validate **`index.es.md`** (including the Gemma 4 gate) before touching `facebook-es.txt` / `substack.es.md` / `linkedin.es.txt`. Sidecars are **reflections** of the index, not independent sources.
8. **Line-by-line author edits:** If the user pastes a replacement line, apply **only that patch**. Do not rewrite the whole paragraph “for coherence.” Author wording binds even when it looks like a calque. MUST NOT call Gemma 4 to “fix” an intentional author calque.

### Gemma 4 native-Spanish gate

MUST run after the first **`index.es.md`** draft, **before** ES sidecars. Same local gateway and script as **`.cursor/skills/site-revise-spanish/SKILL.md`**.

```bash
python3 .cursor/skills/site-revise-spanish/scripts/evaluate_spanish.py \
  content/<section>/<slug>/index.es.md \
  --scope site
```

Default model: `LOCAL_LLM_MODEL` or `gemini_free/gemma-4-26b-a4b-it` at `LOCAL_LLM_BASE_URL` or `${POLYPUS_BASE_URL}/v1`.

**While drafting with this skill:** apply report fixes that fail native-Spanish, reverse-translation, or meaning-load. Do **not** wait for a separate apply confirmation (that wait is for **revise-spanish** on already-shipped copy). MUST NOT apply fluency-only rewrites that drop contrasts or verb objects (same **Apply / "mejora" guardrails** as **revise-spanish**).

**After Gemma, still own reverse-translation:** a **Ready** score does **not** waive an English clause-map. If `tldr` / `fluff` / `description` / body still remounts as the EN sentence, rewrite syntax and re-check **Clause rebuild**.

**Skip Gemma 4:** user pasted a one-line author patch (workflow step 8).

**Gateway down (no fallback):** MUST report that Gemma 4 did not run and MUST stop the Spanish review there. Ship the draft as written. MUST NOT substitute an agent-run audit, "Forced Self-Correction Loop", or **Native opinion-editor audit** for the gate. MUST NOT apply idiom, calque, or fluency edits the author did not request. The user decides whether to retry the gateway or accept the draft unaudited.

### Index is source; sidecars are reflections

| MUST | MUST NOT |
|------|----------|
| Treat **`index.es.md`** as the only Spanish source of truth for thesis, wording, and section ladder | Edit the index to match a sidecar |
| Sync sidecars **only after** the index reads native and keeps author terms | Let Facebook/Substack invent contrasts or soften clinical precision missing from the index |
| Keep `title` / `subtitle` as **one short line** each (no body bleed) | Dump body paragraphs into scalar front-matter fields |
| Preserve author lexicon from the index (*persona*, *patrón*, *sobre-diagnostican*) | Swap to softer generics (*individuo*, *hábito*, *consultorio*) when syncing |

### Author micro-edits bind

When the user (or `index.es.md` after their edits) chooses a phrase that “sounds like a calque” but carries tension, **author wins**.

| Scenario | Decision |
|----------|----------|
| Phrase looks like EN calque but is author/index wording (*se les cae el piso*) | **Keep** |
| Phrase is raw / dry / café, not “professional” Spanish | **Keep** (do not consultorio-smooth) |
| Soft abstract feeling where the beat needs mechanism | Prefer **active** Spanish (*usan los sentimientos…*) |
| Agent “improvement” softens clinical contrast | **Reject** |

**MUST** keep author micro-edits exact (including intentional calques).  
**MUST NOT** rewrite a whole paragraph when a one-line patch already lands.  
**MUST NOT** professionalize raw manifesto voice.  
**Priority of calque:** if source/author uses a non-standard but loaded phrase, keep it over textbook fluency.  
**Preserve contrast:** do not soften conflict into polite essay Spanish.

#### Spoken manifesto vs soft essay (from live edits)

| Avoid (soft essay) | Prefer (spoken / manifesto) |
|--------------------|-----------------------------|
| *la conversación se vuelve confusa* | *la conversación se va al carajo* / *se convierte en caos* |
| *sienten que su imagen se ve afectada* | *se les cae el piso* (if author chose it) |
| *un hecho difícil de negar* | *un hecho irrefutable* (when force is the point) |
| *intentan cambiar la historia* / *necesitaba borrar* | *continuidad inconveniente… necesita tergiversación* |
| *el sentimiento vale como prueba* | *usan los sentimientos como prueba y acomodan la realidad…* |
| *colgarte sus propios golpes* | *atribuirte sus propios defectos* |
| *poco contacto y límites fríos… que otra conversación heroica* | *evitar el contacto y establecer límites, en vez de perder el tiempo en conversaciones heroicas* |

### Hard constraints (native Spanish, not translation-shaped)

**Goal:** Same **thesis, ladder, and facts** as the English sibling; **not** a sentence-for-sentence English gloss. Spanish may use its own `###` hooks, bridges, and punch lines when they match the paragraph.

| Rule | MUST | MUST NOT |
|------|------|----------|
| **Adaptation** | Write as a **native editor** would: meaning, rhythm, and collocations in Spanish. | Mirror English syntax, section titles word-for-word, or "translated Medium" voice. |
| **English skeleton** | Rebuild clause order in Spanish (dative *le*, short remates, native verb–object). | Keep EN skeletons with ES words (*Para el X, ese Y es peligroso, porque… allí donde…*; *Ahí es donde entra…*). |
| **Conceptual anglicisms** | Replace English-shaped phrases with Spanish idioms (see table below). | Leave calques that only make sense after reading the English (e.g. *viajar gratis* for *free ride*). |
| **Syntax** | Prefer **active voice** and natural **se** constructions. | Stack passive voice that reveals translation (*La decisión fue tomada por…*). |
| **Tone** | Stay **direct, incisive, analytical**; match the source's edge. | Smooth into corporate Spanish (*es importante destacar*, *cabe señalar*, excessive *en conclusión*). |
| **Technical loanwords** | Keep hyper-specific terms in **English** when Spanish has no stable standard (*alignment faking*, *benchmark*, *benchmarks*). | Force awkward Spanish neologisms for established EN literature terms. |
| **Paragraph flow** | Transitions must follow **Spanish logic** (cause → consequence in prose order readers expect). | Translate connective words only (*Now stretch that pattern…* → literal *Ahora estira ese patrón…*). |
| **`###` hooks (claims / long posts)** | MAY use **colloquial Spanish** hooks that differ from English `###` titles; each hook **MUST** fit the paragraph below. | Empty labels, hooks that repeat the Claim verbatim, or hooks that describe a different beat than the section body. |

**Before shipping `index.es.md`:** The **Gemma 4 native-Spanish gate** MUST have run, or its failure MUST be reported per **Gateway down (no fallback)** above. Read the draft aloud. Run **Clause rebuild** and the **reverse-translation test** on every user-facing sentence: if you can put the Spanish back into English almost word-for-word and recover the EN sentence shape, **fail and rewrite syntax** (not only vocabulary). A Gemma **Ready** does not skip this. If it still sounds like dubbed English, apply Gemma 4 findings again and run the **five-directive audit**.

#### English skeleton under Spanish words (MUST check)

| English-shaped (Avoid) | Prefer (idiomatic ES) |
|------------------------|------------------------|
| *Para el [X], ese [Y] es peligroso, porque… allí donde…* | *Al [X] ese [Y] le resulta peligroso: [verbo]…* |
| *Ahí es donde entra [concepto]* | *Ahí entra [concepto]* |
| *los hechos reducen lo desconocido* | *un hecho reduce lo que no se sabe* / *despeja dudas* |
| *su versión se quede arriba* | *su versión mande* / *prevalezca* |
| *se les cae el piso* (*floor drops out*) | *se les viene abajo el suelo* / *se les desmorona el suelo* |
| Long EN chain: *Feeling counts as proof and reality bends…; then your tone…, ignoring…* | Short ES clauses with remate; do not mirror the semicolon stack |

#### Clause rebuild (MUST, before typing)

Do **not** draft by swapping Spanish words into the English sentence. Close the English clause. Write the Spanish thought. Then reverse-translate.

**Spanish default is hypotaxis, not telegram.** Native Spanish often finishes the thought in **one** sentence with *y*, *que*, *cuando*, *mientras*, *porque*. MUST NOT split a native clause into English-style punches (*El otro solo ejecuta.*) just to “sound punchy.”

**MUST NOT** keep these frames (EN clause-map). Prefer **elaborated Spanish**, not a remate stack:

| EN frame | Fail (ES calco) | Prefer |
|----------|-----------------|--------|
| *The X who Y carries Z* (long relative doing two EN jobs) | *El cómplice callado que te allana el terreno carga más culpa que…* (same nest as EN) | Native *y*: *El cómplice que calla te allana el terreno y carga más culpa que el que da el golpe.* Split into two sentences **only** if the nest still copies English. |
| *thus* / *making X possible* / *and so* | *y así el delito se hace posible*; *haciendo que* | *quita las trabas y deja el delito servido, mientras el otro solo ejecuta* |
| *You hear it when* (scene open) | *La oyes cuando…*; *La escuchas cuando…* as the first clause | *Se ve en…*; *Pasa cuando…*; scene first, then the saying |
| *in a vacuum* | *nace en el vacío*; *ocurre en el vacío* | *El delito no viene solo*; *Nadie delinque solo* |
| *You reach for it when* (orphan *it*) | *La sueltas cuando…*; *La alcanzas cuando…* with no noun | *Te sale el dicho cuando…*; *Ahí sueltas la frase…* |
| Stacked *The one who…* | Three *El que…* in a row because EN stacked them | Vary: *el que…* / *al que cubres…* / *el que escala…* |

**MUST**

- Rebuild EN relatives in Spanish. Default: **one elaborated sentence** (author gold: *El cómplice que calla te allana el terreno y carga más culpa que el que da el golpe.*). Split **only** when reverse-translation still recovers the English nest.
- Put cause in a verb (*permite*, *abre paso*, *deja servido*), then keep the contrast in the **same** sentence when Spanish would (*mientras el otro solo ejecuta*).
- Count whether the **nest** matches English, not whether the sentence is “long enough” or “short enough.”
- Sayings use-line **`La usas cuando`** is allowed (series convention, gold: `saying-19`). Name the dicho if *la* has no antecedent.

#### Conceptual anglicisms (MUST check)

| English idea | Avoid (calque) | Prefer (idiomatic ES) |
|--------------|----------------|------------------------|
| *free ride* (no cost to lying) | deja de **viajar** gratis | deja de **salir** gratis; deja de salir impune |
| *trigger* limits / restrictions | **dispara** límites | **activa** límites; activa restricciones |
| *wiring defect* (metaphor) | defecto de **cableado** (often sounds translated) | **defecto de diseño** (when the fix is systemic habit, not neurosurgery) |
| *truth bias* (psych term) | sesgo **hacia** la verdad (softer but less standard) | **sesgo a la verdad**; sesgo de verdad |
| *What does not add up* | Lo que **no cierra** (regional; OK in some voices) | Lo que **no cuadra** (wider default) |
| *on purpose* / deliberate | a propósito (next to manifesto gravity) | a conciencia; deliberadamente |
| *commitment escalation* | escalada de compromiso; compromisos que se encadenan (listy) | **espiral de compromiso** (manifesto / trap tone); *cómo ir atándote al compromiso* if you need a verb phrase |

Repo example (validated): `content/human-condition/2026-06-11-deception-truth-bias-incentive-gaps/index.es.md` (colloquial `###`, *salir gratis*, *activa límites*, Claim *defecto de diseño*).

### Hybrid cadence (punchy-natural body)

English manifesto copy often uses **telegram beats** (bare stats, `Label: value` lines). In Spanish those usually read like PowerPoint or dubbed English unless they are **intentional rhetorical parallels** (see **Anti-staccato guardrail** below).

**Default for `type: claims` / `type: video` body** (not a license to telegram **sayings** `tldr`/`fluff`): **hybrid cadence** for essays. For **sayings** `tldr`/`fluff`, default to **elaborated Spanish** (see **Clause rebuild**). A short remate is allowed after a built-up clause, not as the whole block.

| # | Rule | MUST | MUST NOT |
|---|------|------|----------|
| 1 | **No telegram stats** | Fold bare numbers into context: *En un estudio con 14.779 personas y 17.950 conversaciones…* | Mirror EN *14.779 personas. 17.950 conversaciones.* unless the author wants staccato in ES |
| 2 | **Verbs over label lists** | One flowing clause for scale/cost: *Cuesta apenas céntimos por hora, se replica sin límite y no se cansa.* | Colon-led chips: *Coste: X. Copias: Y. Fatiga: Z.* |
| 3 | **Remate pattern** | Medium sentence builds tension; **short** sentence lands the hook. | Four equal micro-sentences with the same grammar and no parallel job |

**Connectors:** Use *y*, *pero*, *;*, and causal links so the reader can breathe. **Colons:** MUST NOT abuse `:` to fake EN slide bullets; integrate facts in narrative prose.

**Tension with "Punchy and direct" above:** Short lines stay valid when the **previous** sentence prepares them. Do not telegram the whole essay.

**Before / after (body smell → fix):**

| Smell | Prefer |
|-------|--------|
| *La IA era mejor haciéndote actuar… Puedes irte… Apunta a la que hace clic.* (chips) | *La IA resultó ser mucho más eficaz para obligarte a actuar que para convencerte de algo. Al final… solo busca activar tu **clic**.* |

**Repo example:** `content/mind-infrastructure/2026-06-18-ai-persuasion-infrastructure/index.es.md`.

### Spanish idiom, flow, and character (MUST follow)

- **COT para traducción (Chain of Thought):** Antes de escribir la versión final, el modelo debe realizar un proceso interno de razonamiento:
  1. **Desmontar el inglés:** Identificar el **mecanismo o intención** real de la frase (¿qué está pasando realmente? ¿quién hace qué a quién?).
  2. **Transmutación de conceptos (No calcar metáforas):** Las metáforas abstractas del inglés a menudo fallan en español. Aterrizarlas a imágenes visuales y reales. (Ej: "narrative debt" → **montaña de justificaciones**, **factura emocional**; "is a slow way to lose" → **una derrota en cámara lenta**; "it backfires" → **te explota en la cara**).
  3. **Vocabulario de acción:** Listar **verbos con "sabor"** que describan ese mecanismo en español (por ejemplo: en lugar de "avoid", usar **esquivar**, **rehuir**, **capear**; en lugar de "move", usar **desplazar**, **empujar**, **arrastrar**).
  4. **Variación léxica cross-field:** Planear sinónimos diferentes para cada campo (`description`, `tldr`, `fluff`, body) para evitar repeticiones. Si `tldr` usa "trampa", `fluff` usa "aprieto" o "lío"; si `tldr` dice "correr/volar", body usa "velocidad sobrehumana".
  5. **Uso del "Nos" (Autoridad y Conexión):** Cuando el inglés usa formas impersonales o pasivas para describir comportamientos humanos, el español gana fuerza usando la primera persona del plural (**nosotros**). (Ej: "Politics trains people..." → **"La política nos obliga..."**).
  6. **Fuerza y forma directa:** Priorizar la **autoridad y la brevedad**. Sustituir verbos débiles ("significa que", "se trata de") por el **"es"** directo o verbos de identidad fuerte. Si una frase puede decirse con tres palabras en lugar de seis, elegir las tres.
  7. **Filtro de naturalidad (Forced Loop):** Probar la frase en voz alta. Si suena a "robot traduciendo" o a "ensayo corporativo", **es obligatorio volver al paso 1**. Iterar hasta que la estructura suene **100% natural, directa y con autoridad** en español. Si el tema es jerarquía, empresa o diseño organizativo, recorrer también la tabla **Organizational, hierarchy, and workplace English** (más abajo).
- **Goal:** Readers should not feel a **literal gloss** or machine translation. Prefer **meaning and natural collocations** in Spanish over **cognate-for-cognate** wording.
- **No alien literal translations (Cero calcos literales):** If an English phrase relies on specific phrasing that sounds bizarre or robotic when translated word-for-word (e.g., "biological brains" -> "cerebros biológicos"), DO NOT translate it literally. Find the natural Spanish way to express the underlying human experience (e.g., "nuestra biología").
- **Consistent direct address (El "tú"):** MUST maintain a consistent, informal singular "tú" cuando te dirijas al lector. MUST NOT mix "tú", "ustedes" (les), and impersonal forms in the same text.
- **Watch for person switches:** When describing roles, modes, or states (e.g., "El Capitán", "La Nave", "El Ejecutor"), ensure you don't accidentally slip from addressing the reader ("tú") to describing them in third person ("él/ella/ellos"). Either consistently use the "tú" form throughout or clearly mark the switch as intentional character voice. Common trap: "Si no tienes postura, terminan convertidos en..." should be "...acabas como...".
- **Bone-deep simplicity (Al hueso):** Remove filler phrases that explain the process instead of the impact. Avoid "Llega un punto en que...", "Decimos que...", "Se trata de...", or "¿Sabías que...?". Go straight to the feeling or the fact.
  - *Example:* Instead of "Decimos que amamos elegir", use **"Nos gusta elegir"**.
- **Persona común vs Profesor:** El texto debe sonar como una persona inteligente hablando con alguien en el día a día, no como un profesor de sociología o un ensayo académico. Avoid overly clinical or formal translations ("neutralidad acrítica", "te desentiendes del resultado", "amenazas a la identidad"). Use everyday, street-smart equivalents that carry the same punch ("neutralidad por inercia", "te lavas las manos", "miedo a quedar fuera").
- **Grounding as Narrative:** For `grounding` fields, avoid technical lists or academic summaries. Translate the *finding* into a clear, human narrative.
  - *Example:* Instead of "vinculaban los sentimientos al resultado", use **"si alguien rompe una regla y gana, asumen que se siente bien"**.
- **Universal vs Local:** Avoid idioms or slang specific to one region (e.g., "cañas" for forced socials) unless the post's persona explicitly requires it. Use universal Spanish that any speaker understands (e.g., **"reuniones sociales"**).
- **Concrete over abstract (Lenguaje visual):** Translate the *intent* and the *visual*, not just the industry jargon. If an English tech term feels forced or abstract in Spanish (e.g., "feeds change"), rewrite it as a concrete action (e.g., "el algoritmo decide qué mostrar").
  - *Example:* Use **"el sistema por defecto"** or **"lo que ya viene decidido"** instead of "valores por defecto".
- **Management / HR English in Spanish prose:** In longer Spanish copy, **knowledge workers** usually reads more natural as **"trabajadores del conocimiento"** or **"profesionales del conocimiento"** than as a calque like **"obreros del conocimiento"**. MAY keep the English phrase in *italics* or quotes when the audience expects the loanword or you need a sharp aside; default to Spanish when the sentence is expository, not punchy.
- **Punchy and direct:** Avoid wordy, passive, or overly formal constructions. Use short, punchy sentences with character. (e.g., instead of "Habrá desde notas rápidas hasta piezas más pulidas", prefer "Verás desde notas crudas hasta ensayos pulidos").
- **Metaphor:** Cuando el inglés use una imagen (e.g. "aim the beam"), **keep the image** in Spanish if it still works (**dirigir el haz**, **orientar el foco**). Do **not** glue the same English verb onto the wrong object.
- **Check:** If the Spanish phrase would sound odd in a **standalone** sentence (title or Claim line), **rewrite** for idiomatic Spanish while preserving the claim and tone.
- **Quotes:** Primary-source English quotations may stay in English; surrounding commentary should still be natural Spanish.

### Native opinion-editor audit (five directives)

After the first draft, and again after the **Gemma 4 native-Spanish gate** returns a report, act as a **native Spanish opinion editor**: punchy, direct, manifesto-grade prose. Goal: remove **syntactic anglicisms** and **calques** without softening the source.

| # | Directive | MUST do | MUST NOT |
|---|-----------|---------|----------|
| 1 | **Calques and false friends** | Replace literal English locutions with Spanish collocations. Use the calque tables below (science/AI + org/hierarchy). | Leave *en papel*, *pagas* [abstract cost], *en acción*, *a propósito* (grave tone), *alguien a quien culpar*, etc. |
| 2 | **Fluid syntax (anti-staccato)** | Fix **true redundancy** and **robotic** chains (same noun/number repeated, beats with no parallel job). Use a pointer (*de ese tamaño*, *ahí*) instead of repeating the subject. **Vary length** across the paragraph, not by flattening every short line. | **Abrupt cuts:** merge intentional manifesto hooks; replace short parallels with long semicolon sentences; strip phrases that still carry meaning or echo (*La culpa no es solo…* before *echarle la culpa*). |
| 3 | **Natural collocations** | Verb–noun pairs must sound organic (*Agiliza la coordinación*, *El precio que pagas son…*, *tragarse el coste* / *asumir las consecuencias*). If a metaphor sounds **dubbed** (*alguien a quien culpar*), swap for the native equivalent with the **same impact**. | Keep a calque because it is "close enough" to the English metaphor. |
| 4 | **Preserve tone** | Keep raw, sharp, satirical, or challenging energy from the English. | Smooth into corporate, moralizing, or generic AI voice (*es importante destacar*, *en conclusión*, *cabe señalar*). |
| 5 | **Output when auditing in chat** | If the user asks for an idiom audit (not a silent file edit), deliver: **`[Texto final]`** (full corrected text) then **`[Changelog]`** (3–4 bullets only: biggest idiom fixes and **why** the original sounded translated). | Dump every micro-edit; rewrite the argument; soften punch lines "for clarity." |

**Repo edits:** When applying changes to `index.es.md` or sidecars, write the file directly; use **`[Texto final]`** / **`[Changelog]`** only when the user wants a review pass in chat first.

**Cross-check tables:** **Calcos frecuentes en contenido de ciencia e IA** and **Organizational, hierarchy, and workplace English** below.

#### Anti-staccato guardrail (no abrupt cuts)

**Telegram dubbing ≠ intentional staccato.** Isolated EN stats and `Label: value` lines → **Hybrid cadence** above. Keep short parallels when they **land** a hook (*En familia… En clase…*).

Short sentences are **not** automatically wrong in Spanish. Opinion and manifesto prose often **needs** staccato: parallel openings (*En familia… En clase… En la empresa…*), binary contrasts (*Abajo el detalle; arriba, el resumen*), rule-of-three beats.

**Apply directive 2 only when:**

1. The **same key noun or number** repeats in back-to-back sentences without rhetorical gain (fix: *diez personas… Diez personas debaten* → *en una reunión de ese tamaño*).
2. Several beats share **identical syntax** with **no parallel purpose** (generic S+V+O five times in a row).
3. A sentence only **restates** the previous one (echo), not when it **lands** the hook.

**MUST NOT** merge parallel hooks into semicolon chains to "sound more Spanish." **MUST NOT** delete closers (*para poder entenderlo*) unless they are empty filler or a calque. **When in doubt on intentional manifesto parallels, keep short beats.** For sayings `tldr`/`fluff` and most ES prose, default to **elaborated hypotaxis** (**Clause rebuild**); do not telegram-split a native *y* / *porque* clause.

**Changelog discipline:** In idiom audits, list **only** changes that fix calques or objective redundancy. Do **not** report staccato "fixes" that removed intentional rhythm unless the author asked for a compression pass.

### Directrices para Cognitive-Memetics (Sayings/Cows)

**Evitar la redundancia (No al eco):**
- **Description (El Gancho):** MUST NOT repetir palabras clave que ya están en el `title` o el `tldr`. Si el título es "Como ratón al queso", la descripción no debe usar "ratón" ni "queso". Debe ser un "pitch" evocador, no un resumen.
- **Variación entre campos:** Si `tldr` usa "trampa" para el concepto, `fluff` debe usar un sinónimo ("aprieto", "lío", "corral") o una expresión diferente que capture el mismo mecanismo.
- **No repitas estructuras:** Si `tldr` dice "o X, o Y" (o corres, o vuelas), el body debe reformular la idea sin repetir la estructura binaria exacta.

**Eliminar muletillas de traducción:**
- NOT "Literalmente significa...", "Esta frase se refiere a...", "La imagen del X marca...".
- NOT "Refleja una visión de..." (calco del inglés académico); USE "Es el suspiro de...", "Aquí ves cómo...".
- NOT fluff that opens like a subtitle: *La oyes cuando…*, *La escuchas cuando…* (*You hear it when*). USE *Se ve en…*, *Pasa cuando…*, or a human beat first (*Es el suspiro de…*).
- NOT *y así X se hace posible* / *nace en el vacío* (EN glue). Rebuild per **Clause rebuild**.
- USE entradas directas: "Cuando alguien se pierde en...", "El oportunista no pide, acecha", "Un ruego para sobrevivir a...".

**Narrativa con voz humana (Sayings):**
- Empezar con el **humano y su emoción**, no con abstracciones. "Es el suspiro de quien ya entendió..." > "Refleja una visión resignada de..."
- Usar la primera persona implícita ("te", "tu") para conectar con el lector.
- Incluir **escenarios de uso** específicos: el mentor aconsejando, el amigo alertando, la derrota personal.

**Capturar matices de tono:**
- El mismo saying puede ser **advertencia** (a un amigo que duda), **consejo** (de mentor a protegido), **crítica** (del sistema), o **resignación** (tras una derrota).
- El `fluff` debe reflejar estos matices en los ejemplos de uso, no solo el tono negativo.
- Alternar entre carga **positiva** (ingenio, agilidad, avivado) y **negativa** (oportunismo, sin escrúpulos) según el contexto del saying.

**Verbos con "Sabor" y Acción:**
- Preferir verbos que describan la **actitud** o el **mecanismo** detrás del comportamiento (por ejemplo: **esquiva**, **desplaza**, **cubre**, **paga**, **erosiona**).
- "engañar" → "farolear", "mentir", "ocultar".
- "aprovechar" → "sacar tajada", "roer el beneficio", "estar al acecho".
- "pedir paciencia" → "un ruego definitivo", "un suspiro compartido".
- "superar" (a alguien) → "ganarle la mano", "dejar botado", "quedar tieso".

**Equilibrio de negritas (Bold):**
- MUST follow **`.cursor/skills/site-revise-emphasis/SKILL.md`** for all `content/` types (same restrained default as English: a few hooks per block; not wall-to-wall bold).
- MUST NOT saturar el texto con negritas. El objetivo es guiar el ojo, no gritar cada palabra.

**Body sin redundancia:**
- El body debe **complementar**, no repetir lo que `tldr` y `fluff` ya dijeron.
- Si `tldr` explicó la estructura binaria "o corres, o vuelas", el body profundiza en las **consecuencias culturales**, no vuelve a explicar la hipérbole.
- Eliminar frases que "explican el proceso": "Esto significa que...", "La exageración subraya que...". Ir directo al impacto.

### Calcos frecuentes en contenido de ciencia e IA (MUST check)

**Verbos biológicos:**
- "neurons fire" → NOT "sus neuronas disparan" (intransitivo sin reflexivo suena a calco); USE "sus neuronas se disparan" or rewrite with an evocative verb: "sus neuronas recorren por anticipado las rutas posibles".
- "the brain replays the option" → NOT "el cerebro reproduce la opción"; USE "el cerebro rebobina y reproduce" or a similarly active construction.

**Verbos debilitados:**
- "can feel regret" → NOT "pueden sentir arrepentimiento"; USE direct affirmation: "sienten arrepentimiento".
- "had to simulate" → NOT "tuvieron que simular"; PREFER "les tocó simular" or "la evolución los obligó a simular".

**Metáforas biológicas que funcionan bien en español:**
- "atrophied muscle" for cognitive offloading → "músculo atrofiado" works; use it.
- "the brain hallucinates and corrects" → "el cerebro alucina y usa los sentidos para confirmar" works; prefer over "adivina" when the neuroscience framing is explicit.
- "rewind" for memory replay → "rebobina" works and is vivid.

### Organizational, hierarchy, and workplace English (MUST check)

Editorial passes on **human-condition**, **claims**, and org-design posts surfaced **English skeletons** that pass grammar checks but fail the **read-aloud** test. Run this table while drafting and when the **Gemma 4 native-Spanish gate** returns a report, whenever the source discusses startups, meetings, rank, flat teams, or bureaucracy.

| English pattern | Avoid (calque) | Prefer (idiomatic ES) | Notes |
|-----------------|----------------|-------------------------|-------|
| *On paper* | En papel, una startup… | **Sobre el papel**, una startup… | Standard Spanish locution for "in theory." |
| *proved* (study / lab) | Zink lo **probó** en el laboratorio | Zink lo **demostró** / **comprobó** | *Probar* reads as "try" or "test," not "established by evidence." |
| *X in action* | la pirámide **en acción** | **en estado puro**, **la esencia de** la pirámide | *En acción* mirrors *in action*; Spanish favors essence/state framing. |
| *on purpose* (deliberate design) | diseñarlo **a propósito** | **a conciencia**, **deliberadamente** | *A propósito* can sound colloquial or childish next to manifesto gravity; use for casual tone only. |
| *you pay* [abstract cost] | A cambio, **pagas** cuellos de botella | **El precio que pagas son**…, **sufres**…, **a cambio te tragas**… | Spanish rarely *paga* bottlenecks, ego fights, or gaps; price/suffer/swallow collocations work. |
| *someone to blame* | te da **alguien a quien culpar** | te da **a quién echarle la culpa** | Default for direct, street-smart voice. |
| *scapegoat* (optional punch) | (overused default) | **chivo expiatorio** | Stronger and more organizational; can sound essay-like or caricatured. **Default:** *echarle la culpa* unless the beat needs sarcastic HR tone. |
| *counter-design* / design against drift | diseño **a la contra explícito** | diseño **explícito que vaya a contracorriente**, **contrapeso diseñado a conciencia** | Fix awkward word order; keep the "against the default" meaning. |
| *coordinates fast* (funnel benefit) | **Coordina** rápido | **Agiliza la coordinación** | Verb–object pairing from English (*coordinate coordination*) often needs a natural Spanish verb. |
| *gap* between decision and work | un **hueco** enorme entre quien decide… | un **abismo** entre **el que toma la decisión** y **el que** lava los platos | *El que… el que* often beats bare *quien… quien* in contrast pairs. |
| *when you reach your floor* (career ladder) | cuando llegas a tu piso | **en el momento en que** llegas a tu piso | Use the longer frame only when rhythm needs it; do not stack it in every paragraph. |
| *lets X rise* / *boosts X over builders* (prosperity, feeds, algorithms) | **deja subir**, **suben**, **impulsa** (permisivo o vago) | **promueve**, **promueven** | Active causation: the system **promotes** the wrong profile; it is not an open door. Pair with **escalan jerarquías** when naming rank climb. |

**Example rewrite (embudo / org funnel):** After the idiom pass, a block should read like natural manifesto Spanish, not imported *Medium* English:

> Ese embudo es la pirámide en estado puro. Agiliza la coordinación y te da a quién echarle la culpa. El precio que pagas son cuellos de botella, peleas de ego y un abismo entre el que toma la decisión y el que lava los platos rotos. La escalera también te enseña a dejar de pensar en el momento en que llegas a tu piso.

**Cross-field sync:** When Claim or body wording changes (*promueve/promueven*, *a conciencia*, *demostró*), update **`facebook-es.txt`** and **`substack.es.md`** in the same bundle so social copy does not keep stale calques or slip back to **suben** / **deja subir** / **impulsa** for the same hierarchy beat.

6. **Taxonomies (`categories`, `tags`)**: Keep values **identical** to the English post unless the project explicitly adopts translated taxonomy terms. Matching terms avoid duplicate tag/category hubs and keep cross-language analytics consistent. If you must change a tag, coordinate EN and ES and accept a new taxonomy term. When choosing tags on either side, read **`data/tag-register.txt`** and **`.cursor/skills/site-tag-register/SKILL.md`** first.
7. **`related` (claims / video):** When the English page sets **`related`**, mirror the **same path strings** on **`index.es.md`**. `Site.GetPage` in the Spanish language context resolves translated siblings when they exist. Do not invent Spanish-only related paths that the English page lacks.
8. **URLs, permalinks, and `aliases`**: Follow **[Spanish URLs, permalinks, and aliases](#spanish-urls-permalinks-and-aliases)** below (canonical paths, internal links, redirects, `facebook-es.txt` / `linkedin.es.txt`).
9. **Punctuation**: Same as English site rule—**do not** use the em dash (U+2014); use comma, semicolon, colon, or parentheses.

## Spanish URLs, permalinks, and aliases

### Canonical URLs (MUST)

- Read **`hugo.toml`** → **`[languages.es.permalinks]`**. The **folder** under `content/` (for example `social-protocols/`) **does not** always equal the Spanish URL segment. Example: **`social-protocols`** maps to **`/es/protocolos-sociales/:slug/`**, not `/es/social-protocols/...`.
- For each language, the **post slug** in the URL usually comes from the **translated `title`** (Hugo's `:slug` in that permalink pattern), **not** from the English bundle folder name (for example `2026-04-28-candor-zero-sum-politics/` is the **filesystem** slug; the Spanish permalink slug is title-derived).
- MUST **not** guess final URLs. Run **`hugo list all`** (or **`make list`** at the repo root) and copy the **Permalink** column for **`index.es.md`** and **`index.md`** for the same bundle.

### Internal links in Spanish body copy (MUST)

- When a **Spanish page exists** for the target, use its **canonical** path from `hugo list all` (often under `/es/...` with the translated section segment and title slug).
- When the target has **no** `index.es.md`, link to the **English** path **without** `/es/` (for example `/social-protocols/2026-04-26-prisoner-dilemma-tit-for-tat/`). MUST **not** invent `/es/social-protocols/<english-folder>/` for content that only exists in English.

### `aliases` on `index.es.md` (SHOULD, `content/social-protocols/`)

- SHOULD keep **legacy redirects** so old shared links still resolve. For each Spanish post in this section, mirror the pattern used in sibling bundles:
  - `/es/social-protocols/<english-folder-name>/`
  - `/es/protocolos-sociales/<english-folder-name>/`
  - `/es/reality-protocols/<english-folder-name>/`
- If the English bundle lists extra **typo** or legacy paths (for example alternate dates in `aliases`), mirror the **`/es/...`** form on the Spanish page when applicable.
- MUST **not** duplicate a **language-default** alias path on the Spanish page that the **English** page already owns (for example bare **`/reality-protocols/<slug>/`** on `index.md` targets the English URL; the Spanish page should use **`/es/reality-protocols/<slug>/`** only).

### Social files: `facebook-es.txt`, `substack.es.md`, and `linkedin.es.txt` (MUST)

- **`facebook-es.txt`** and **`linkedin.es.txt`** MUST use the **same** ES and EN **Permalink** values as `hugo list all`, in the formats defined in **`.cursor/skills/site-facebook-post/SKILL.md`** and **`.cursor/skills/site-linkedin-post/SKILL.md`** (dual-language site link blocks).
- **`substack.es.md`** (required for **`make sb-es`**) MUST follow **`.cursor/skills/site-substack-post/SKILL.md`**: structured Markdown, no hashtags or read-more block in the sidecar.
- MUST **not** construct Facebook or LinkedIn URLs from the English folder name plus `/es/social-protocols/`; that combination is often **wrong** for this site.

## Field cheat sheet by type

| Type | Translate |
|------|------------|
| **`claims`** | `title`, `description` (Claim), body (Thoughts), `grounding`; optional `image_credit` |
| **`video`** | `title`, `description` (lead), `sowhat` if present; body if any |
| **`sayings` / `panel`** | `title`, `description`, `tldr`, `fluff`, body; follow **`.cursor/skills/site-revise-emphasis/SKILL.md`** for **bold** |
| **Section `_index`** | `title`, `description` in `_index.es.md` only |

## Do not

- MUST NOT replace the English file or duplicate content under a different section path for "Spanish only."
- MUST NOT drop `translationKey` when both languages exist; it powers correct language switching and translation links.
- MUST not translate **identifier-like** tags into Spanish if that would fork the taxonomy (e.g. `StreetWisdom` → keep as-is unless the whole site migrates terms).
- MUST NOT clause-map: Spanish words in English sentence order (same relatives, *thus/making*, *You hear it when*). Rewrite syntax per **Clause rebuild**.

## After editing

- Run `hugo` and confirm the Spanish URL renders and the language switcher lists both languages for that `translationKey`.

## Related

- **Placement and section**: `.cursor/rules/site-content-placement.mdc`
- **English default prose rules**: `.cursor/rules/site-content-markdown-writing.mdc` (English); Spanish files follow this skill for language choice.
- **Hooks / titles**: `.cursor/skills/site-revise-hooks/SKILL.md` after the type skill, adapted for Spanish.
- **Facebook (ES/EN) and LinkedIn (ES) permalinks**: `.cursor/skills/site-facebook-post/SKILL.md`, `.cursor/skills/site-linkedin-post/SKILL.md` (site link blocks; always align with **`hugo list all`**).
- **Gemma 4 evaluate script / standalone audit**: `.cursor/skills/site-revise-spanish/SKILL.md` (this skill runs the same script while drafting; **revise-spanish** is evaluate-only until the user confirms).
- **Step-by-step Spanish revision workflow**: `.cursor/skills/site-revise-post-es/SKILL.md` (use this when doing a full translation or fixing a poor translation).
