# Revise Spanish: reference

Deep rules live in **`.cursor/skills/site-spanish-translation-content/SKILL.md`**. This file is a **fast sniff list** for read-aloud audits and the **local Gemma evaluate** pack.

## Scored criteria (local eval, 0–2 each)

| ID | Pass (2) | Fail (0) |
|----|----------|----------|
| `native_naturalness` | Café / manifesto Spanish for the surface | Dubbing or academic translation |
| `anti_calque` | None or rare intentional loans | EN syntax/verbs as Spanish, **or fluent ES words on an English sentence skeleton** |
| `read_aloud` | Clean aloud | Stumbles, only clear after EN |
| `collocation_grammar` | Natural collocations | Wrong object, EN verb calques |
| `cross_file_coherence` | Same thesis/refrains/close across ES files | Sidecar contradicts site ES |
| `conceptual_fidelity` | Keeps key contrasts, subject/object load, and thesis tension | Softens terms or drops objects so it “sounds better” but the idea thins |

**Pass overall:** Ready / Almost / Re-adapt section (see skill).

---

## English skeleton under Spanish words (hard fail)

A paragraph can sound smooth and still fail **anti_calque** / **native_naturalness** if clause order mirrors English.

| English-shaped | Prefer (ES) |
|----------------|-------------|
| *En una pelea ordinaria, los hechos reducen lo desconocido. Aquí dejan un rastro…* | *En una pelea corriente, un hecho reduce lo que no se sabe. Aquí no: deja un rastro…* |
| *Para el falso yo (…), ese rastro es peligroso, porque fuerza continuidad allí donde…* | *Al falso yo (…) ese rastro le resulta peligroso: obliga a una continuidad que… necesitaba borrar* |
| *Necesitan que su versión se quede arriba* | *Necesitan que su versión mande* |
| *Ahí es donde entra el razonamiento emocional* | *Ahí entra el razonamiento emocional* |
| *El sentimiento cuenta como prueba y la realidad se tuerce…; entonces tu tono…, ignorando…* | Short clauses: *El sentimiento vale como prueba; la realidad se acomoda… Tu tono… y tu intención desaparece.* |
| *se les cae el piso* (*floor drops out*) | *se les viene abajo el suelo* / *se les desmorona el piso* |
| *The X who Y carries Z* | Native hypotaxis: *El cómplice que calla te allana el terreno y carga más culpa que el que da el golpe.* Split **only** if the nest still copies English. |
| *thus* / *making X possible* | *quita las trabas y deja el delito servido, mientras el otro solo ejecuta* (not *y así … se hace posible*, and not two English punches) |
| *You hear it when…* as opener | *Se ve en…* / *Pasa cuando…* (not *La oyes cuando…*) |
| *in a vacuum* | *El delito no viene solo* (not *nace en el vacío*) |
| *You reach for it when* (orphan *it*) | *Te sale el dicho cuando…* |

**Test:** if you can put the Spanish clauses back into English almost word-for-word and recover the original EN sentence, rewrite the Spanish syntax. Do not only swap synonyms.

**Exception:** author-locked lines in `index.es.md` (including intentional calques) are not “fix targets.”

## Cognitive false friends, pop calques, and gerunds (MUST check)

| English pattern | Avoid (calque) | Prefer (idiomatic ES) | Notes |
|---|---|---|---|
| *interest-based* / *interest-driven* | *se mueve por intereses*, *basado en el interés* | *guiado por la curiosidad*, *estimulado por el interés intrínseco*, *funciona por fascinación genuina* | En español, *moverse por intereses* connota egoísmo o oportunismo. |
| *side quests* | *misiones secundarias* | *proyectos paralelos*, *iniciativas personales*, *inquietudes fuera de ruta* | *Misiones secundarias* suena a videojuego y desentona en prosa de ensayo. |
| Frases hechas inglesas literales | *saltos de gigante* (giant leaps), *bala de plata* (silver bullet), *elefante en la habitación*, *cambiador de juego* | *grandes avances*, *solución mágica*, *problema evidente*, *punto de inflexión* | Las frases hechas (idioms) inglesas traducidas palabra por palabra suenan artificiales. |
| Gerundio de posterioridad (*-ing* consecutivo) | *filtrando a los cerebros...*, *reprimiendo las ganas...*, *viendo patrones...* | *lo que filtraba a los cerebros...*, *al tiempo que reprimían...*, *y detectan patrones...* | En español el gerundio no expresa consecuencia o resultado posterior; exige verbo conjugado. |
| Metáforas literales no asimiladas | *tentáculos lineales*, *curiosidad tipo pulpo*, *llevar tentáculos* | *tareas rutinarias y lineales*, *curiosidad expansiva*, *delegar la carga mecánica* | Aterrizar a la imagen o función real; no calcar adjetivos ingleses sobre metáforas biológicas. |
| *career tax* | *impuesto de carrera* | *lastre profesional*, *peaje en la trayectoria* | *Tax* abstracto se traduce por lastre, peaje o coste, no impuesto directo. |
| *routine data entry* | *entrada rutinaria de datos* | *registro mecánico de datos*, *papeleo rutinario* | Evitar ristras de sustantivos abstractos sin verbo. |

## Video / list front matter (description bullets)

Novelty bullets under **What you probably do not know yet** must read as Spanish claims, not translated slide titles.

| English-shaped | Prefer (ES) |
|----------------|-------------|
| *La IA como palanca ante X: de A a B* | *La IA no solo compensa X: puede convertirla en ventaja…* (verbo + afirmación) |
| *Cómo hacer rendir… al delegar…* | *Los proyectos paralelos despegan cuando dejas que…* (sin viñeta *Cómo…*) |
| *El ocaso del especialista de carril único* | *El viejo culto a la hiperespecialización penalizaba…* (evitar *carril único* = traffic calque; evitar *mentes diseñadas para*) |
| *brains built/designed to connect* | *cerebros capaces de conectar* / *volcados en enlazar saberes* |

**Test:** if a bullet is a nominal phrase with a colon and no finite verb, or starts with *Cómo*, treat as **anti_calque** unless author-locked.

## Re-adaptación de párrafo vs parches cosméticos

Cuando la prueba de reversión detecte un esqueleto sintáctico inglés:
- **PROHIBIDO:** sustituir solo dos o tres palabras sueltas manteniendo la misma estructura de oraciones calcada del inglés. Eso deja el texto sonando a doblaje.
- **OBLIGATORIO:** Re-adaptar el párrafo completo. Extraer la afirmación central y el mecanismo, cerrar la frase inglesa y redactar en español fluido con hipotaxis natural (conectores orgánicos y subordinación).

## Soft essay vs spoken manifesto (false “improvements”)

| Soft / diluted | Prefer when force is the point |
|----------------|--------------------------------|
| *el sentimiento vale como prueba* | *usan los sentimientos como prueba y acomodan la realidad…* |
| *colgarte sus propios golpes* | *atribuirte sus propios defectos* |
| *poco contacto… que otra conversación heroica* | *evitar el contacto y establecer límites, en vez de perder el tiempo en conversaciones heroicas* |
| *Esa paz, negando la realidad, no dura* | *Esa "paz" que niega la realidad no dura* |

Do not invent these; use when the index/author already pushes that direction.

---

## Meaning load and contrast (do not “fix” away)

Naturalness must not sacrifice argument precision. When editing or proposing fixes after Gemma:

| Edit smell | Rule |
|------------|------|
| Soften technical / author terms | If the text uses *clínico*, *patología*, *diagnóstico*, *patrón*, do not swap for softer generics (*consultorio*, *hábito*, *costumbre*) unless the author asks |
| Drop the object | If the verb loads a person (*sobre-diagnosticar a la persona*, *cargar a la persona*), keep the object |
| Stack three ideas in one opaque colon | Fix direction: clear **cause → effect → what the fight really is**, not prettier compression |
| Lose A vs B tension | If the line contrasts clinic labels vs pattern-first naming, the fix must keep both poles |

**Example (fix direction, not full rewrite):**

- Opaque: *Cuando el patrón es no mirar de frente lo ocurrido, los hechos no cierran la pelea: lo que importa es la coartada.*
- Bad “smooth” fix: *Si no se mira lo ocurrido, los hechos no importan: solo cuenta la coartada.* (changes the claim)
- Better direction: *Si el patrón es negarse a mirar lo ocurrido, traer hechos no cierra la pelea: el tema en disputa solo es coartada.*

Author wording after a micro-edit (*sobre-diagnostican*, *se enfoca*) **wins** over synonym polish.

---

## Perceptual / body / mind posts

| English-shaped | Prefer (ES) |
|----------------|-------------|
| *maquinaria perceptiva* | *circuito perceptivo* / *la misma percepción* |
| *asignaciones sentidas* (every sentence) | once in Claim; then *donde te ubica…* / *dos trucos* |
| *lo que tu cerebro dice que es tu mano/cabeza* (parallel ×2) | *Ahí ubica la mano* / *Ahí te sitúa a ti* |
| *La piel reporta* | *Solo sientes presión* |
| *las palabras aterrizan* | *siguen llegando* |
| *clean story* pasted | *historia limpia* (OK once) |
| *pick up more of the work* | *asumen más peso* |
| *sighted template* | *patrón de quien ve y oye* |
| *on that evidence* | *Por eso* |
| *head-anchored senses* before a list | drop; list *ojos, oídos, equilibrio* |
| *credit your shoulder with a plan* | *le achacas el plan* / *no vas a decir que el hombro…* |

## Social sidecars only

| Smell | Fix |
|-------|-----|
| ✜ blocks, TLDR labels | Facebook: plain paragraphs only |
| *Misma máquina* vs *Mismo sistema* in same bundle | pick one |
| *Si quieren* + *tú* in body | *Si quieres leer más* |
| LinkedIn EN structure copied | rewrite for café (facebook-post) or ladder (linkedin-post), not both |

## Hybrid cadence (telegram dubbing)

Flag when ES mirrors EN **bare stats** or **colon labels** without rhetorical job:

| Smell | Prefer |
|-------|--------|
| *14.779 personas. 17.950 conversaciones.* | *En un estudio con 14.779 personas y 17.950 conversaciones…* |
| *Coste: X. Copias: Y. Fatiga: Z.* | *Cuesta…, se replica… y no se cansa.* (verbs, one idea) |
| *escalada de compromiso* | *espiral de compromiso* |
| Four equal chip sentences before the punch | Setup + **short remate** (*solo busca activar tu clic*) |

Full rules: **spanish-translation-content** → **Hybrid cadence**. Example bundle: `content/mind-infrastructure/2026-06-18-ai-persuasion-infrastructure/`.

## Anti-staccato (do not "fix")

Keep short parallels when they **land** a hook:

- *En familia… En clase…*
- *Mismo sistema, dos trucos*
- *Mismo sistema. Dos sitios.*

Merge only when the **same noun or number** repeats with no rhetorical job.

## Read-aloud fails (grammar)

| Fail | Fix pattern |
|------|-------------|
| *Esa sensación es donde…* | *Ahí es donde…* |
| *te sitúa la cabeza* | *te sitúa a ti* |
| *Esto significa que…* | cut; state the fact |
