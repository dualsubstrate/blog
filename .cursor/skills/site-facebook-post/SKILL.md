---
name: site-facebook-post
description: >-
  Writes Facebook posts for friends-and-family from Hugo page bundles. Spanish
  (facebook-es.txt from index.es.md) and English (facebook-en.txt from index.md).
  Substack uses substack.md (see substack-post skill). Adaptation over fidelity:
  warm café tone, not LinkedIn copy. Plain text, no em dash. For type video with a
  site Chapter Guide, includes a friends-style topic outline (Label · topic, no clock
  times). Triggers: facebook post, facebook-es, facebook-en, friends Facebook copy.
---

# Facebook post from site content

## Philosophy: adaptation over fidelity

| Translation-first (avoid) | Adaptation-first (this skill) |
|---------------------------|-------------------------------|
| Match the English LinkedIn post | Work for friends on Facebook |
| Copy `linkedin.txt` structure | Write fresh for the platform |
| Professional cadence | Café conversation with friends |
| Word-for-word from the site post | Spoken language that carries the idea |

**Golden rule:** Each post should feel written *for* friends on Facebook, not *exported from* the site or LinkedIn.

## Shared: audience and platform

- **Audience:** friends and family on Facebook, not professional colleagues.
- **Length:** target **150–250 words**; opening lines carry the hook (fold / “see more”).
- **Plain text only.** No `**bold**`, no Markdown headings.
- **Emoji are welcome** where they fit the bundle (e.g. project line on cognitive-memetics).
- **Emoji alignment (MUST):** if the bundle has a leading emoji in **`title`** (or a conventional lead emoji in `linkedin.txt` for that bundle), reuse that same emoji at the start of the Facebook hook line. Do not add extra decorative emoji beyond what the bundle already uses.
- **No em dash (U+2014).** Use commas, semicolons, colons, or spaced hyphens.
- **Short paragraphs:** 1–3 lines; blank lines between blocks.

## Shared: source fields

Read the locale page as **inspiration**, not as a template:

| Site field | Use as |
|------------|--------|
| `title` | Hook starting point; adapt for Facebook tone |
| `description` | Core idea; rephrase conversationally |
| `tldr` / `fluff` (sayings) or body | Material to compress into prose (see [Hugo `fluff` vs LinkedIn Context](#hugo-fluff-vs-linkedin-context)) |
| Body **`## Chapter Guide`** / **`## Guía de capítulos`** | **`type: video` only:** topic outline (see [Video chapter outline](#video-chapter-outline)); copy label + topic per row from the **locale** index; **omit clock times** |
| `grounding` | Optional casual citation at the end (claims, etc.) |
| `tags` | Hashtag line (from the locale file you used) |

**MUST ignore `linkedin.txt` structure** unless the user explicitly asks to align with LinkedIn.

## Hugo `fluff` vs LinkedIn Context

On Hugo **`type: sayings`** bundles (and other cognitive-memetics series that carry **`tldr`** / **`fluff`** in front matter), the site field is still named **`fluff`**. LinkedIn uses a labeled block **`➕ Context:`** for that material (never **`➕ FLUFF:`**); see **`.cursor/skills/site-linkedin-post/SKILL.md`** → **Section label: Context**.

| Surface | How to use site `fluff` |
|---------|-------------------------|
| **`linkedin.txt`** | **`➕ Context:`** section after **`✔️ TLDR:`** (Por-Estas-Calles, Sm(art), Reptilocracy, Psych-Fitness-28, and same shape wherever that skill applies labeled blocks) |
| **`facebook-en.txt`** / **`facebook-es.txt`** | **No** section labels; weave **`tldr`** then **`fluff`** ideas into short plain paragraphs in the insight block |
| **`substack.md`** / **`substack.es.md`** | Adapted prose with **`##`** hooks; no **`➕ Context:`** label ( **`.cursor/skills/site-substack-post/SKILL.md`** ) |

**MUST NOT** paste **`✔️ TLDR:`** / **`➕ Context:`** / **`❓ BUT WHY:`** labels into Facebook copy.

## Compression (`type: sayings`, Sm(art), Reptilocracy)

- MUST compress **`tldr`** + **`fluff`** into spoken paragraphs; MUST NOT paste site blocks or restate the thesis in generic form three times.
- SHOULD run **`.cursor/skills/site-revise-post/SKILL.md`** → **Step 2** (standalone worth): cut warmups with no claim, merge redundant beats, one sharp image per metaphor block.
- **Length:** default **150–250 words** for claims and video; **`type: sayings`** / **Sm(art)** **MAY** run **~100–150** when every line earns its place.

## Sm(art) / T-Shirt Art (Facebook)

- **MUST NOT** use LinkedIn product bridges mid-argument (*La camiseta Sm(art) lo plantea en serio:*, *Así lo dice la camiseta…* after setup). The graphic and **`#TShirtArt`** carry the series; deliver the idea as café talk.
- **MAY** name the shirt once at the open only when the whole post is clearly “sharing a design with friends” (see **`2026-05-23-tshirt-predictability-automates-first/facebook-es.txt`**); default for argumentative **`tldr`**: no shirt meta.
- **Metaphor gloss:** when a coinage stays in the post (*feromonas*, feed-as-*reina*), add a **short plain parenthetical** once if friends need the mapping (for example *notificaciones y tendencias que te marcan el ritmo*).

## Shared: post skeleton

```text
[Line 1: hook from title or description, adapted for friends — MUST carry substance in the first line or two; standalone *¿Te suena?* / *Ever notice?* alone is warmup]

[Blank line]

[Setup: why this matters to us]

[Insight: 2–4 short paragraphs, plain prose]

[Blank line]

[Close: one or two punchy lines]

[Blank line]

[Video only, when locale index has Chapter Guide:]
[📑 heading per locale]
[Blank line]
[• Label · chapter title  (one line per site table row; no MM:SS)]

[Blank line]

[#Tag1 #Tag2 ... from front matter tags]

[Blank line]

[Read-more block: both ES and EN URLs when the bundle has index.md + index.es.md]

[Optional: casual research citation]
```

## Shared: site links (ES + EN)

- When the bundle has **`index.md`** and **`index.es.md`**, output **both** URLs.
- Use **final permalinks** from `hugo list all`, not guessed paths.
- **Spanish post (`facebook-es.txt`):** ES first, then EN.
- **English post (`facebook-en.txt`):** EN first, then ES.

Spanish block (**MUST** match body address; see [Read-more line (Spanish)](#read-more-line-spanish)):

```text
Si quieres leer más:

- ES: <full Spanish permalink>

- EN: <full English permalink>
```

Use **`Si quieren leer más:`** only when the **whole** post body stays in **nosotros** / **ustedes** (no **tú** in the argument). Default for direct **tú** posts (most claims-style friends copy): **`Si quieres leer más:`**.

English block:

```text
If you want to read more:

- EN: <full English permalink>

- ES: <full Spanish permalink>
```

### Read-more line (Spanish)

| Body address in `facebook-es.txt` | Read-more heading (MUST) |
|-----------------------------------|---------------------------|
| **tú** (default: hooks like *¿Te ha pasado…?*, *dejas*, *te quedas*) | **`Si quieres leer más:`** |
| **nosotros** / **ustedes** throughout the argument | **`Si quieren leer más:`** |

**MUST NOT** mix **tú** in the body with **`Si quieren leer más:`** (sounds like the author switched mid-post). **MUST NOT** use **`Si quieres`** when the body never uses **tú**.

## Shared: execution

1. Confirm locale(s): user asks for **Spanish**, **English**, or **both**.
2. Read **`index.es.md`** and/or **`index.md`** for the core idea.
3. Write fresh for friends (hook → setup → idea → close).
4. **`type: video`:** when the locale index has a **Chapter Guide** table, add the [Video chapter outline](#video-chapter-outline) after the close (from that locale's rows, friends wording).
5. Run **`hugo list all`**; filter permalinks for the bundle slug.
6. Add **all** `tags` from the locale source file as hashtags.
7. **`type: sayings`** / **Sm(art):** compress **`tldr`** + **`fluff`**; run [Compression](#compression-type-sayings-smart-reptilocracy) and standalone-worth pass.
8. Run the locale **voice pass** below; adjust only real inconsistencies.

## Video chapter outline

When **`type: video`** and the locale **`index.md`** or **`index.es.md`** body has a chapter table (**`## Chapter Guide`** / **`## Guía de capítulos`**):

### Goal

Friends get a **plain topic list** (who / what beat), not a jump map. **Do not** copy `linkedin.txt` (no ✜ ladder, no professional blocks). Clock times and timestamp links stay on the **site** and in **`substack.md`** / **`substack.es.md`**.

### Rules (MUST)

- **Placement:** after the close, **before** hashtags.
- **Heading (locale):**
  - English **`facebook-en.txt`:** `📑 What's in the talk`
  - Spanish **`facebook-es.txt`:** `📑 De qué va el video`
- **Rows:** one line per site table row, same order as that locale's chapter table. Format: `• Label · chapter title` (Unicode bullet `•`, middle dot ` · `). Strip `**` from labels.
- **Omit `MM:SS` and `t=` URLs** on every Facebook row (same rationale as **`.cursor/skills/site-linkedin-post/SKILL.md`** → **Video chapter guide (TOC)**).
- **Wording:** **MAY** shorten titles for spoken friends copy; **MUST** use the **locale** chapter text for **`facebook-es.txt`** (not English bullets translated line by line unless no Spanish table exists).
- **Not counted** toward the **150–250 word** body target.
- **Optional:** one casual watch line after hashtags or after read-more, for example `Video: https://www.youtube.com/watch?v=VIDEO_ID` (plain URL; no Markdown).

### Forbidden

- Pipe tables, box-drawing separators, space-padded columns (Facebook uses a proportional font; columns will not align).
- Clock times on each bullet.
- Copying **`linkedin.txt`** structure or tone for the outline or body.

### Where timestamps stay

| Surface | Chapter data |
|---------|----------------|
| Site `index.md` / `index.es.md` | Full chapter table with YouTube links |
| `substack.md` / `substack.es.md` | Full table ( **`.cursor/skills/site-substack-post/SKILL.md`** ) |
| `linkedin.txt` | `• Label · topic`, no times ( **`.cursor/skills/site-linkedin-post/SKILL.md`** ) |
| `facebook-en.txt` / `facebook-es.txt` | Same outline shape as LinkedIn; **friends** body and headings per locale |

### Reference bundle

**`content/mind-infrastructure/2026-05-21-brain-is-not-a-computer/`** (`facebook-en.txt`, `facebook-es.txt`). See **`.cursor/skills/site-video-content/SKILL.md`** for site table authoring.

```text
📑 What's in the talk

• Friston · free energy and the spherical cow
• Chirimuuta · why scientists abstract
…
```

## Spanish (`facebook-es.txt`)

### Tone

- Conversational Spanish: talk over coffee.
- Avoid corporate or academic cadence: no *potenciar*, *impulsar resultados*, *leverage*.
- Natural openers: *Mira*, *Fíjate:*, *Resulta que*, *¿Te has fijado?*

### Voice pass

- **One address per stretch:** do not mix **tú**, **ustedes**, and **uno**. Singular hook → **tú** body; group hook → **ustedes** or rewrite. **Read-more line:** **`Si quieres leer más:`** when the body uses **tú**; **`Si quieren leer más:`** only when the body stays plural (see [Read-more line (Spanish)](#read-more-line-spanish)).
- **Reflexives:** align with chosen person (*taparte* with **tú**, not *taparse* next to *Buscas*).
- **Concrete, spoken Spanish** over neutral filler; light quoting (*'nos pasa a todos'*) when it helps.
- **Hierarchy verbs:** when **`index.es.md`** uses **promueve/promueven** or **escalan jerarquías** for prosperity, feeds, or algorithms elevating the wrong people, **MUST** carry the same verbs into **`facebook-es.txt`**. **MUST NOT** revert to **suben**, **deja subir**, or **impulsa** for that beat (see **`.cursor/skills/site-spanish-translation-content/SKILL.md`** → **Organizational, hierarchy, and workplace English**).
- **Stop when it sounds like café;** do not over-rewrite.

### Spanish anti-patterns

| Don't | Do instead |
|-------|------------|
| *Tu feed no es un flujo neutral* | *¿Se han fijado en que el celular nos mete en una rutina?* |
| *Es una aldea reconstruida* | *Es como si nos metieran en el pueblo, pero digital* |
| Taxonomy blocks (*✜ Identidad*) | Natural prose |
| *Optimiza estabilidad del grupo* | *Lo que les importa es que no nos salgamos del coro* |
| **Piensa** then **Piensen** / **uno** | One person per block |
| **taparse** mixed with **tu** / **Buscas** | **taparte** or consistent infinitive list |
| *Escuchamos a diario* / *En contraste* | *¿Te suena? Cada semana alguien dice…* / *Y ahora:* |
| *La camiseta Sm(art) lo plantea en serio:* | Thesis in plain prose (see [Sm(art) / T-Shirt Art (Facebook)](#smart--t-shirt-art-facebook)) |
| **tú** mixed with *escuchamos* / *funcionamos* | Pick **tú** or plural; rewrite openers |
| *de forma natural*, *sociedades donde* | Concrete scene (*la calle se los comió sin drama*) |
| **suben** / **deja subir** / **impulsa** (wrong people rise in feeds or hierarchies) | **promueven** / **promueve** (match **`index.es.md`**) |
| *El patrón se siente real* / *The pattern still feels real* (after a quote) | Cut the bridge; go straight to the mechanism or example |

### Output

Save as **`facebook-es.txt`** next to **`index.es.md`**. Plain text only.

## English (`facebook-en.txt`)

### Tone

- Conversational English: friends on Facebook, not LinkedIn.
- Avoid corporate or essay cadence: no *leverage*, *optimize*, *utilize*, *folks should consider*.
- Natural openers: *Look:*, *You know that feeling*, *Picture this*, *Ever notice*.
- Prefer **you** / **we** over passive piles; plain verbs over abstract nouns.

### Voice pass

- **Pick you or we** and stay consistent in the body (mixed *you* / *one* / passive is a common slip).
- **Concrete scenes** (group chat, family thread, news cycle) over thesis statements.
- **Not LinkedIn:** no ✔️ TLDR / ➕ Context / ❓ BUT WHY labeled blocks, or hub-link lines copied from `linkedin.txt`.
- **Stop when it sounds spoken;** do not over-polish into op-ed tone.

### English anti-patterns

| Don't | Do instead |
|-------|------------|
| *Your feed is not a neutral stream* | *Ever notice the phone keeps you in the same loop?* |
| *This is a reconstructed village* | *It is like the old village, but run by algorithms* |
| LinkedIn template blocks | Short plain paragraphs |
| *The system optimizes group stability* | *What they want is nobody breaking from the chorus* |
| *One must consider* / heavy *individuals* | *You* / *we* and plain verbs |
| Copying `linkedin.txt` verbatim | Fresh friends post from `index.md` ideas |

### Output

Save as **`facebook-en.txt`** next to **`index.md`**. Plain text only.

## Substack

Newsletter bodies live in **`substack.md`** / **`substack.es.md`**, not in facebook sidecars. See **`.cursor/skills/site-substack-post/SKILL.md`**.

## Example (Spanish)

**LinkedIn (formal):** taxonomy blocks, *Tu feed no es un flujo neutral*.

**Facebook (friends):** question hook with claim in line 1–2, *Fíjate que…*, warm close, hashtags, *Si quieres leer más* (with **tú** body) or *Si quieren leer más* (plural body only), with ES + EN URLs. **Video + outline:** `content/mind-infrastructure/2026-05-21-brain-is-not-a-computer/facebook-es.txt`. **Sm(art) café tú:** `content/cognitive-memetics/t-shirt-art/2026-05-29-tshirt-not-at-the-border-at-the-table/facebook-es.txt`.

## Example (English)

**Video + chapter outline:** `content/mind-infrastructure/2026-05-21-brain-is-not-a-computer/facebook-en.txt`.

**Tone and link order:** `content/cognitive-memetics/t-shirt-art/2026-05-16-smart-cover-your-ears/facebook-en.txt`, `content/cognitive-memetics/reptilocracy/2026-05-24-missiles-get-the-money/facebook-en.txt`.

## Related skills

- **Substack:** `.cursor/skills/site-substack-post/SKILL.md`
- **Standalone worth (compression pass):** `.cursor/skills/site-revise-post/SKILL.md` → **Step 2**
- **Video picks (site chapter table):** `.cursor/skills/site-video-content/SKILL.md`
- **Spanish naturalness audit:** `.cursor/skills/site-revise-spanish/SKILL.md`
- **Spanish revision:** `.cursor/skills/site-revise-post-es/SKILL.md`
- **LinkedIn (English, professional):** `.cursor/skills/site-linkedin-post/SKILL.md`
- **Spanish site pages:** `.cursor/skills/site-spanish-translation-content/SKILL.md`
