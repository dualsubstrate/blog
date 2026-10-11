---
name: site-cognitive-memetics-content
description: >-
  Authors and edits Hugo posts under the cognitive-memetics section: type panel
  (cube-cows, Raymond, and other cartoon hubs) or type sayings
  (TLDR / Context). **`description`** teasers MUST be drafted via local **Gemma 4**
  (punch the scene; not a caption) for every hub **except Por-Estas-Calles**, which
  keeps **Por-Estas-Calles card teaser** from **`title`** + **`tldr`** + **`fluff`**.
  Also: Street Wisdom tags and LinkedIn hashtags, T-Shirt Art, sayings **`**bold**`**
  per **revise-emphasis**. Footer "But why" cards share one gradient in
  **`assets/css/_custom.scss`**. Plain **title**. Ask before changing existing author
  prose (exceptions: sayings emphasis; Por-Estas-Calles card teaser; Gemma teaser when
  drafting or when the user asks to rewrite the teaser).
---

# Cognitive-Memetics section content

## What this section is

**Cognitive-Memetics** is a **Hugo section**: paths under `content/cognitive-memetics/...` (URL prefix `/cognitive-memetics/`). It is listed in **`hugo.toml`** `params.home.contentSections` so posts can appear on the home feed with other sections.

This section does **not** define a separate `type`. Posts use existing content types:

| `type` | Role | Typical use in this folder |
|--------|------|----------------------------|
| **`panel`** | **Cube-cows** / weekly-style pieces (Cows **storyboards**; each image a **panel**) | **`description`** carries the strip (Teaser on the single); optional markdown body below `<!--more-->` only when you need a real essay. Optional `heading_code`; optional **`project`** + unique **`title`** |
| **`sayings`** | Short entries with TLDR + Context | `tldr` and `fluff`; card **`description`**: **Por-Estas-Calles** uses derived **Por-Estas-Calles card teaser**; other hubs (T-Shirt Art, Reptilocracy, Pawtropolis, …) use **Gemma teaser**. Optional `heading_code`; optional **`project`** + unique **`title`** |

**T-Shirt Art** hub posts also use **`panel`** or **`sayings`**; see [below](#t-shirt-art) for **`categories`**, **`project`**, and tags. **Raymond** (`content/cognitive-memetics/raymond/`) uses **`type: panel`** with category **`Raymond`**; Cube-Cows spinoff starring Raymond (junior dog from *Just smaller*). **Reptilocracy** (`content/cognitive-memetics/reptilocracy/`) uses **`type: sayings`** or **`type: panel`** with category **`Reptilocracy`**; the footer explainer uses the same gradient card as the Cube-Cows hub and Street Wisdom, plus Reptilocracy-only CTA styling when applicable (see **Theme and style** → **Project "But why" explainer cards** in this skill). **Pawtropolis (Under Fire)** (`content/cognitive-memetics/pawtropolis/`) uses **`type: panel`** or **`type: sayings`** with category **`Pawtropolis-Under-Fire`**; copy lives in **`i18n/*`** and **`layouts/partials/pawtropolis-project-about.html`** (no petition CTA).

For **`type: claims`** (Claim / Thoughts / Grounding), use **`.cursor/skills/site-claims-content/SKILL.md`** and a section such as `social-protocols`, not this skill.

## Authoring rules (site-wide)

- When editing **`tags`**, MUST consult **`data/tag-register.txt`** and **`.cursor/skills/site-tag-register/SKILL.md`** (prefer existing tags when they fit; new tags allowed when none do).
- MUST NOT use the em dash character (U+2014). Use comma, semicolon, colon, or parentheses.
- MUST use Markdown **`**bold**`** for scan-style emphasis on key verbs and nouns in the body (restrained density per **`.cursor/skills/site-revise-emphasis/SKILL.md`**).
- MUST keep code, identifiers, and user-facing strings in **English** (see workspace AI protocol).

### Ask before changing existing copy

- Phrases like **apply this skill**, **check the post**, or **make it compliant** do **not** grant permission to rewrite the author’s prose.
- When a file already has author text in **`title`**, **`description`**, **`tldr`**, **`fluff`**, or the markdown body, MUST **ask first** before changing that text (including “fixing” tone, clarity, mismatches you think you see, or aligning fields). Wait for an explicit yes to prose edits (for example: rewrite, revise, tighten, change the title, fix the TLDR). **Exceptions:** (1) **`type: sayings` Sayings emphasis** (adding, trimming, or adjusting **`**bold**`** on **existing** words in **`tldr`** and **`fluff`** per **`.cursor/skills/site-revise-emphasis/SKILL.md`**) is **not** a prose edit and MUST follow that skill when you apply this one. (2) **Por-Estas-Calles card teaser** (replacing **`description`** using only **`title`**, **`tldr`**, and **`fluff`**, only when **`categories`** includes **`Por-Estas-Calles`**) is **not** a prose edit of TLDR/context and MUST run when you apply this skill to that hub. (3) **Gemma teaser** (see **Gemma teaser (`description`)**): for every Cognitive-Memetics hub **except Por-Estas-Calles**, when **creating** a new episode or when the user asks to rewrite / punch up / Gemma the teaser, MUST invoke local Gemma 4 and MUST NOT ship an agent-invented **`description`** as final.
- If the user does not confirm prose edits, restrict changes to **mechanical** front matter, **`type: sayings` emphasis** (per **revise-emphasis**), **Por-Estas-Calles card teaser** when that hub applies, and structure only (see **Editing existing posts** below). Do **not** silently replace an existing non–Por-Estas-Calles **`description`** with a Gemma draft unless the user asked for a teaser rewrite.

### Editing existing posts (preserve author copy)

- When the author already supplied prose in **`tldr`**, **`fluff`**, **`title`**, or the markdown body, MUST **keep their wording** unless they **explicitly** ask for a rewrite, revise, tighten, or similar. For hubs other than **Por-Estas-Calles**, the same applies to an **already shipped** **`description`** until they ask to rewrite the teaser.
- **Exception (Por-Estas-Calles):** **`description`** is **not** protected author copy when you apply this skill to a **`Por-Estas-Calles`** post. MUST replace it per **Por-Estas-Calles card teaser**.
- **Exception (Gemma hubs, new draft):** When the episode is **new** (first **`description`**) or the user asked for a teaser rewrite on any hub **except Por-Estas-Calles**, MUST follow **Gemma teaser (`description`)** instead of agent-authored caption copy.
- MAY apply **mechanical** fixes without asking: YAML safety (quoting, colons), required front matter fields, **`categories`** / **`tags`** / **`date`** / **`heading_code`**, featured image paths, and obvious **wrong-folder** mistakes (for example mismatched **`heading_code`** vs bundle name) when fixing structure.
- For **`type: sayings`**, when applying this skill, MUST apply **Sayings emphasis** on **`tldr`** and **`fluff`** first. Then: if **`Por-Estas-Calles`**, run **Por-Estas-Calles card teaser**; otherwise run **Gemma teaser** for new/rewritten **`description`**. Do **not** wait for a separate request for bold or a new teaser on create.
- For **`type: panel`** and other types (and for markdown **body** copy on any type), when the author **already wrote** the words, MAY add Markdown **`**bold**`** around existing words or short phrases **for emphasis only** if the user asked for emphasis, revise, tighten, or similar; otherwise MUST ask before adding bold. MUST NOT change the underlying words, add sentences, or reorder paragraphs under this rule.

### Tone and nuance (`type: sayings`)

- A saying can carry **multiple tones** depending on context: **warning** (to a friend who hesitates), **advice** (from mentor to protégé), **criticism** (of the system), or **resignation** (after a loss).
- **`fluff`** should reflect these nuances in usage examples, not just the negative/cynical reading.
- Alternate between **positive/resourceful** framing (agility, ingenuity, being sharp) and **negative/opportunistic** framing when the saying supports both.

### Human-first narrative (`type: sayings`)

- Start with the **human and their emotion**, not abstract analysis.
- **Weak:** "This reflects a resigned view of broken systems..."
- **Strong:** "It is the sigh of someone who finally understands how broken systems work..."
- Use direct address ("you", "your") to connect with the reader.
- After the emotional hook, the **next clause** MUST name a situation, use, or who-does-what (see **Hybrid prose**). Emotion alone is not enough.

### Hybrid prose (emotion + claim)

Cognitive-Memetics is **not** essay prose. Punch and street feeling are allowed. Empty polish is not.

Essays / video / Substack use **Explanatory prose** in **`.cursor/rules/site-content-markdown-writing.mdc`**. This section uses a **hybrid**: one emotional or cartoon beat, then a clear claim or scene. Full essay paragraphs are optional only when a post deliberately adds body Article copy.

**Applies when drafting or when the user asks to rewrite, revise, or tighten** `description`, `tldr`, `fluff`, or body under `content/cognitive-memetics/`.

#### By field

| Field / format | Punchy by design? | MUST | MUST NOT |
|----------------|-------------------|------|----------|
| Cube-cows / Raymond / T-Shirt Art / Pawtropolis `description` | Yes (often the whole piece) | **Gemma 4** drafts it (see **Gemma teaser**); name the move or joke claim; punch the scene past a caption | Agent-only caption inventory; oracular closers with no mechanism |
| Reptilocracy `description` / `tldr` / `fluff` | Yes | **`.cursor/skills/site-reptilocracy-gemma/SKILL.md`** (satire-first + psych-fitness knife + thinking) | Generic Gemma teaser-only; agent stamp poetry sold as Gemma |
| Por-Estas-Calles `description` | Short card pitch | **Por-Estas-Calles card teaser** from **`title`** + **`tldr`** + **`fluff`** only | Gemma flow; inventing beats not in those fields |
| Sayings `tldr` | Yes (short) | Name meaning, use, or scene | Trait-dictionary stacks only ("clever, astute, street-smart") |
| Sayings `fluff` | Situational context | Lead with emotion or **one** clarifying metaphor, then when/who uses it; concrete scenes | Metaphor-only fluff with no uses; metaphor stacks; mystical/abstract culture praise |
| Psych-Fitness `tldr` | Campaign vignette + lesson | End on the fitness/mechanism line when needed | Caption pivots: "The long view matters.", "The connection is direct.", "The eve of decision." |
| Reptilocracy clinical edit (non-Gemma polish) | Clinical mechanism OK | Name who does what; continuity vs early correction; late cost when story and reality diverge | Dashboard caption poetry; rebel/activist heat; empty MBA fog; aphorism remix of the thesis |
| T-Shirt Art teaser / `tldr` | Merch line OK | Tie the slogan to one mechanism or threat | Slogan-only copy with no so-what |
| Project "But why" (`i18n`) | Manifesto-tinged OK | Name series mechanism once | Do not copy that lyric tone into every episode field |

#### Rules

- MUST keep **one** street/cartoon emotional hook when it helps; the **next** clause MUST name who does what, when someone says it, or what changes if true. **Exception (Reptilocracy):** prefer clinical mechanism first; see **Reptilocracy voice**.
- MUST prefer one concrete scene or use-case over a prettier restatement of the same thesis.
- MUST cut or merge a line that fails: *If this stood alone, would it add a fact, scene, or use?*
- MUST NOT ship **punch closers** that only remix the claim ("The X matters.", "The connection is…", "The eve of…").
- MUST NOT write **revelation stacks** in episode copy (several one-line poetic beats with no mechanism). That ban is for essays in **revise-post** Step 2; here the hybrid still forbids empty stacks.
- MAY keep a short staccato teaser on cube-cows or shirt lines when the claim is still legible ("**Dream big**, commute **bigger**.").
- Project footers MAY stay manifesto-tinged; **episode** `tldr` / `fluff` / `description` MUST stay hybrid, not essay-oracle.

#### Street Wisdom `fluff` (looser metaphor)

For **`Por-Estas-Calles`** / Street Wisdom **`fluff` only**:

- MAY open with **one** literary or street metaphor when it clarifies feeling or stance (for example "warm, **polite fence paint**", "the **sigh** of someone who…").
- MUST follow that metaphor with **concrete uses** (when you say it, to whom, in what setting). If uses are missing, cut or rewrite the metaphor.
- MUST NOT stack two or more literary metaphors in the same `fluff` block.
- MUST NOT let the metaphor remix the thesis without adding a use ("harmony with names still legible" alone fails; pair it with "You reach for it when…").
- **`tldr`** stays tighter: meaning + scene first; save the optional literary beat for `fluff`.

#### Reptilocracy voice (clinical institutional analysis)

For posts under **`content/cognitive-memetics/reptilocracy/`** (category **`Reptilocracy`**), episode **`description`**, **`tldr`**, and **`fluff`** use a **clinical** register. Satire stays dry: name the institutional move as it works in practice. This **overrides** **Human-first narrative** emotion-first openings for Reptilocracy episode fields (the cartoon may still carry the gag; the copy names the mechanism).

**Thesis shape (MUST make legible early, usually in `tldr`):**

1. Institutions treat a hard signal (red gauges, failed policy, bad metric) as a **messaging / board-deck** problem first.
2. **Continuity** is cheaper than early structural change (comms, PR, reframed talking points).
3. When a real correction arrives, it is often **late**: the gap between the **public story** and day-to-day reality has already widened.

**MUST**

- Call things out plainly: who acts (comms, board, PR), what they do (update **public readout**, revise **board deck**, delay audit), why it is rational for them (cost, tenure, predictability).
- Prefer concrete institutional nouns (**press release**, **board deck**, **audit**, **executive suite**, **quiet restructure**) over abstract poetry.
- Keep each `fluff` paragraph earning its place: rhythm → incentive → who it protects → delayed cost (or a tighter subset of those beats).
- Keep episode constructs when the panel names them (for example **Grandiosity Gauges**, **Leadership Confidence: Stable**).
- Still obey site-wide bans: no em dash; no **Industry-verb shells** / fake scene **`the room`** (see **`.cursor/skills/site-revise-post/reference.md`** and **`.cursor/rules/site-content-markdown-writing.mdc`**).

**MUST NOT**

- Caption the panel beat-by-beat (green/red lever poetry that restates the image without a institutional claim).
- Ship rebel or activist heat: hide / lie / deception / fake / eat the cost / martyr frontline closings, unless the author explicitly asks for that register.
- Ship empty MBA fog: *variance*, *operational reality*, *situational awareness*, *data anomaly*, *display error*, *appearance of control*, *mechanics of survival*, *curate the evidence*, *clean slide*, *incentive structure* as filler without naming the move.
- Soften the thesis into plausible-deniability mush that never states the late-cost beat.

**Do / Don't (Reptilocracy)**

| Do (clinical) | Don't |
|---------------|--------|
| Red gauges handled as a **comms** problem; **board deck** still shows **Leadership Confidence: Stable** (`2026-08-30-stability-is-a-setting`) | "One finger toggles the **public readout**… treats a **systemic collapse** as a **display error**" |
| Continuity cheaper than correction; late **quiet restructure** after story and reality diverge | "**Frontline staff** eat the cost… layoffs land on people who never touched the lever" |
| Wage freeze → executive bonus; name the **bait-and-switch** | Floating aphorism after the mechanism is already clear |

**Gold reference:** `content/cognitive-memetics/reptilocracy/2026-08-30-stability-is-a-setting/` (`tldr` states the thesis; `fluff` is clinical sequence). Photo-op rhythm still OK when concrete: `2026-08-02-photo-op-readiness/`.

#### Do / Don't (from shipped posts)

| Do (hybrid) | Don't (empty polish) |
|-------------|----------------------|
| "It is the **sigh** of someone who finally understands how broken systems work… You use it when someone **beat you to it**…" (`saying-19` fluff) | "extremely **clever**, **astute**, and **street-smart**… that kind of **street smarts**, that **practical intelligence**…" (`saying-01` fluff) |
| One metaphor + uses: "warm, **polite fence paint**… You reach for it when proximity is real but **fusion** would cost trust." (`saying-21` shape) | Metaphor stack with no uses: "fence paint… harmony with names still legible…" and stop |
| Sully vignette → "**Psychological fitness checks** help maintain that capacity." (`psych-fitness` day-03) | Mandela vignette → "The long view matters." (`psych-fitness` day-24) |
| "This week, management refreshes **Our Values**: **think big**… but don’t go overboard." (cube-cows) | "The connection is direct." / "The eve of decision." as standalone closers |
| Wage freeze → executive bonus; levy → consultant fees; then name the **bait-and-switch** (reptilocracy) | "We inherit the future from no one, so we spend it freely." as a floating aphorism after the mechanism is clear |

**Gold reference for Street Wisdom fluff:** `content/cognitive-memetics/sayings/2026-05-18-saying-19/` (`fluff`: emotion → uses → settings).  
**Gold reference for Psych-Fitness:** `content/cognitive-memetics/psych-fitness-28/2026-05-04-day-03-clear-thinking/` (scene → mechanism; no oracle coda).  
**Gold reference for Reptilocracy:** `content/cognitive-memetics/reptilocracy/2026-08-30-stability-is-a-setting/` (clinical thesis → delayed cost; see **Reptilocracy voice**).

### Sayings emphasis (`type: sayings`)

- MUST follow **`.cursor/skills/site-revise-emphasis/SKILL.md`** for how much and which words get **`**bold**`** in **`tldr`**, **`fluff`**, and (lightly) **`description`**. Default is **restrained** emphasis (a few strong hooks per block), not wall-to-wall bold.
- When wrapping author text, MUST wrap **only** substrings that already appear in the author’s text; MUST NOT change, add, or remove words, or reorder sentences.
- MUST NOT “fill” **`tldr`** or **`fluff`** with extra bold to match an older dense-scan style **unless** the user asks for that legacy style. If the user asks to **fix** or **reduce** bold, MUST trim per **revise-emphasis** (markup-only is allowed without a separate prose pass).

### Por-Estas-Calles card teaser (`description`, Street Wisdom only)

**Scope:** Posts whose **`categories`** include **`Por-Estas-Calles`** (Venezuelan Street Wisdom). This is the only Cognitive-Memetics hub that skips Gemma for **`description`**.

- When applying this skill to a **Por-Estas-Calles** post, MUST set or replace **`description`** with a **short** card teaser **derived only** from **`title`**, **`tldr`**, and **`fluff`** (if **`fluff`** is absent or empty, use **`title`** and **`tldr`** only). Read **`title`** as the episode anchor (often the Spanish saying); read **`tldr`** and **`fluff`** after **Sayings emphasis** so the teaser matches the list/detail copy.
- MUST NOT invent facts, examples, or tone that are not supported by those three fields together. MAY use **`**bold**`** on key phrases in the teaser you write. Follow **`.cursor/rules/site-content-markdown-writing.mdc`** (English, no em dash, short sentences).
- MUST NOT change **`title`**, **`tldr`**, or **`fluff`** while drafting the teaser unless the user asked for a prose edit to those fields.
- **No echo rule:** `description` MUST NOT repeat key words already in `title` or `tldr`. If `title` is "The rat and the cheese", `description` should not use "rat" or "cheese". Find an evocative pitch that complements without echoing.
- **MUST NOT** use the **Gemma teaser** flow below for **Por-Estas-Calles**.

### Gemma teaser (`description`)

**Scope:** Every Cognitive-Memetics hub **except Por-Estas-Calles** and **except Reptilocracy**. Includes **Cube-Cows**, **Raymond**, **Pawtropolis**, **T-Shirt Art** (panel or sayings), and any future non–Por-Estas-Calles / non–Reptilocracy hub. Hub membership is the second **`categories`** term (or folder), **not** `type: panel` vs `type: sayings`.

**Reptilocracy:** MUST use **`.cursor/skills/site-reptilocracy-gemma/SKILL.md`** for new or rewritten **`description`**, **`tldr`**, and **`fluff`** (not this generic teaser-only flow).

**MUST** invoke local **Gemma 4** to draft **`description`** when:

1. Creating a **new** episode in a Gemma hub (first teaser), or
2. The user asks to rewrite, punch up, Gemma, or replace the teaser.

**MUST NOT** ship a final Gemma-hub **`description`** that the agent wrote alone as a panel caption or prop inventory. **MUST NOT** pretend a hand-written line came from Gemma.

**Gateway (same stack as revise-prose and carousel captions):**

| Setting | Default |
|---------|---------|
| Base URL | `LOCAL_LLM_BASE_URL` or `http://127.0.0.1:1320/v1` |
| Model | `LOCAL_LLM_MODEL` or `gemini_free/gemma-4-26b-a4b-it` |

MAY call via **`.cursor/skills/site_local_eval_common/common.py`** → **`chat_complete`**, or an equivalent `POST …/chat/completions`.

**Prompt goals (EN):**

- Punch the **scene** past a caption: escalate the absurdity; do **not** inventory every prop or quote every bubble.
- Match Cube-Cows gold density when the hub is cartoon satire, for example: *This week, the bull's entire **leadership philosophy** fits inside a **matryoshka** doll: '**be exactly like me**,' he tells the dog, '**just smaller.**'*
- Obey **Hybrid prose** (and **Reptilocracy voice** when that hub applies).
- No em dash (U+2014). Restrained **`**bold**`** on the joke spine only.
- Prefer opening with **This week,** when it fits the series voice.
- Ban empty MBA fog (*synergy*, *collaborative spirit*, *thriving*) unless the joke is mocking that fog on purpose.
- **Raymond:** each candidate MUST use **Dawg Raymond** (or another clear dog signal once). Spanish adapt MUST include **perro** or **perros** (typically **Raymond el perro**).

**Procedure:**

1. Brief Gemma with: series/hub, episode **`title`**, strip or piece setup (mechanism of the gag, not a full visual inventory), and any gold-tone sample for that hub.
2. Ask for **three** numbered teaser candidates.
3. Present them as **🧭 Options** and wait for a numbered pick (unless the user already supplied the final teaser text).
4. Apply the pick to EN **`description`**. Light mechanical cleanup only (em dash, YAML safety, bold density per **revise-emphasis**). MUST NOT rewrite the joke into a flatter caption.
5. If **`index.es.md`** (or sibling ES) exists or is being created: invoke Gemma again to **adapt** the chosen EN teaser into Spanish (native punch, not a calque). Present three ES candidates or apply the best mirror of the EN pick when the user already chose the EN joke spine; then set ES **`description`**.
6. Keep LinkedIn line one in sync when `linkedin.txt` quotes the teaser.

**Gateway down:** MUST say the gateway is unreachable. MUST NOT invent teasers and claim Gemma wrote them. MAY leave a clearly labeled temporary stub only if the user accepts shipping later.

**Skip Gemma only when:** the user pastes the final teaser themselves, or the post is **Por-Estas-Calles** (use **Por-Estas-Calles card teaser**).

## Front matter conventions

### Shared

- **`date`**, **`title`**, **`draft`** — **`date`** MUST follow **`.cursor/rules/site-content-markdown-writing.mdc`** → **Publish `date`** (default **`date: 'YYYY-MM-DDT01:00:00+11:00'`**).
- **`title`**: MUST **not** use **leading emoji** in **`title`**. **Cognitive-Memetics** opts out of **Optional leading emoji in `title`** in **`.cursor/rules/site-content-markdown-writing.mdc`** (which applies to **`social-protocols`**, **`human-condition`**, **`mind-infrastructure`**, and **`x-minds`** only). Keep the episode or saying line in plain words; use optional **`heading_code`** when you want a compact prefix in the UI.
- **`heading_code`** (optional): short label before the title (e.g. `W6`, `W13`). Rendered via `layouts/partials/heading-title-markup.html` with class `heading-code--tldr`.
- **`categories`**: Use **two** terms so each post belongs to the section **and** to a **project hub** you can link to (Hugo taxonomy list pages under `/categories/<slug>/`).
  - **Umbrella:** always **`Cognitive-Memetics`** (this site area).
  - **Project (pick one):** **`Cube-Cows`** for **`type: panel`** (the **Tales from the Cube Farm** series; Hugo taxonomy term for the shareable hub at `/categories/cube-cows/`). **`Raymond`** for the Raymond junior-dog spinoff (`/categories/raymond/`). **`Por-Estas-Calles`** for **Street Wisdom** **`type: sayings`** posts (Venezuelan sayings series). **`T-Shirt Art`** for visual / merch-style pieces (often **`type: sayings`** with art as featured image; **`type: panel`** if you want a longer essay under the same hub). **`Reptilocracy`** for the reptile-institutions satire line. **`Pawtropolis-Under-Fire`** for **Pawtropolis (Under Fire)** (pets in a cartoon war zone; hub slug is typically `/categories/pawtropolis-under-fire/`).
  Example YAML:

  ```yaml
  categories: ["Cognitive-Memetics", "Cube-Cows"]   # type: panel (cube-cow / Tales from the Cube Farm)
  ```

  ```yaml
  categories: ["Cognitive-Memetics", "Raymond"]   # Raymond spinoff
  ```

  ```yaml
  categories: ["Cognitive-Memetics", "Por-Estas-Calles"]       # Street Wisdom sayings
  ```

  ```yaml
  categories: ["Cognitive-Memetics", "Pawtropolis-Under-Fire"]   # Pawtropolis (Under Fire)
  ```

  Hugo flattens categories (no true parent/child in core), but two terms give a **shareable URL** for “everything in this project” while the section path still scopes content by folder.
- **`tags`**: Punchy **PascalCase** hooks (no placeholder terms). Prefer three to five tags; see **`content/cognitive-memetics/`** examples. Align with **`.cursor/skills/site-claims-content/SKILL.md`** *Tag voice* for shape and reuse, but topics here are cube-cows, sayings, and culture notes—not Claim/Grounding jargon unless you intend it.
- **Featured image**: Prefer **`featuredImage: "file.ext"`** and optional **`featuredImagePreview: "file.ext"`** (page-bundle local resource). Use this for the card / detail hero image (put the file in the page bundle). If you need advanced resource metadata, you MAY instead use `resources` with `name: "featured-image"` and `src: "file.ext"`.

### `type: panel`

- **`description`**: Teaser for cards / list; on the **detail** page it becomes the **Teaser** block (markdownified). For this site, **that is usually the whole piece**; do **not** add body copy unless you deliberately want a long follow-up. **MUST** draft new or rewritten teasers via **Gemma teaser (`description`)** above (not agent-only captions).
- **`project`** (optional but recommended for **Tales from the Cube Farm**): The recurring **series line** on the detail hero (e.g. `Cube-Cows 🐮📈`). **`title`** must be a **unique episode name** (lists, prev/next links, browser tab). Without `project`, the layout falls back to `heading_code` + `title` everywhere (legacy one-line titles).
- **Optional** body markdown below `<!--more-->`: only when you need an **Article** section (`layouts/panel/single.html`). If the body is empty, the single shows **Teaser** only.

Archetype: **`archetypes/panel.md`**.

#### Tales from the Cube Farm (why cube-cows exist)

The **Hugo category** (second `categories` term) for these posts is **`Cube-Cows`**. The **series name** in copy stays **Tales from the Cube Farm**.

**Canonical series explainer** (site "But why" footer, Substack paste, LinkedIn **`🟣`** series block): **`cowsProjectAboutTitle`** / **`cowsProjectAboutBody`** in **`i18n/en.toml`** and **`i18n/es.toml`**. **MUST NOT** duplicate or paraphrase that body in this skill. LinkedIn generation: **`.cursor/skills/site-linkedin-post/SKILL.md`**.

**Publish day (Cube-Cows only):** New **`type: panel`** bundles under **`content/cognitive-memetics/cows/`** MUST use a **Wednesday** calendar day in the bundle folder name and in **`date`** (from **W13 / 2026-05-20** onward). Older episodes may stay on their original Thursday dates; do not retro-shift published weeks unless the author asks.

**Tags for `panel`:** Always include **`CubeCows`**. Add **recurring theme** tags when the joke fits (for example **`AGIHype`** when the strip is about AI hype). Add **one or two episode-specific** tags for that week’s punchline; do not force-reuse narrow joke tags across unrelated weeks.

**LinkedIn post format:** Use **`.cursor/skills/site-linkedin-post/SKILL.md`** (*Cube cows / Tales from the Cube Farm* **fold-first layout**): quoted `description` teaser on **line one** only, then **`🟣`** + **`{heading_code}: Cube-Cows 🐮📈`** and series paragraphs from **`cowsProjectAboutBody`** (no `❓ BUT WHY:` label on image-only strips), one hashtag line with **all** `tags` from front matter, **`🧷 Full post (site) →`** (EN then ES when bilingual), then **English** **`Cube-Cows`** category URL (`/categories/cube-cows/`). Save as `linkedin.txt` in the page bundle.

#### Raymond (Cube-Cows spinoff)

Posts under **`content/cognitive-memetics/raymond/`** use **`categories`**: **`Cognitive-Memetics`** and **`Raymond`**. Set **`project: Everyone ❤️ Raymond 🐕📓`** (heart stands in for “loves” in display; taxonomy term stays **`Raymond`**).

**Premise:** Raymond is the junior dog from Cube-Cows *Just smaller* (`2026-08-19-cow-w26`): fresh out of uni, walks into corporate shenanigans with a notebook open. He does not fix the system; he takes notes.

**Publish day:** New bundles MUST use a **Friday** calendar day in the folder name and in **`date`**.

**Tags:** Always include **`EveryoneLovesRaymond`**, **`PetLife`**, and **`TheCutestDog`** (series flourish + pet discovery; these match the dog character). Always include **`OfficeCulture`** and **`OfficeSatire`** (subject matter: every strip is office satire). Do **not** add bare **`Raymond`**, **`Dog`**, **`Dogs`**, or **`DogLife`** when those series tags are present (category **`Raymond`** already hubs the series). MAY add other recurring office tags when they fit (**`RealityCheck`**, **`KnowledgeWork`**). Add one or two episode-specific hooks. Prefer front-matter order: series + pet discovery, then office anchors, then episode hooks. LinkedIn hashtag line MUST mirror all front matter `tags`. Pet tags describe character discovery; office tags describe subject matter. Both are accurate. Do **not** frame tag choice as “tricking” any platform algorithm.

**Dog in copy:** Episode **`description`** (EN) MUST include a dog signal at least once; prefer the character label **Dawg Raymond** (not “Raymond the dog”). Spanish **`description`** MUST include **perro** or **perros** at least once (typically **Raymond el perro**; do not force “Dawg” into ES). Series LinkedIn body already says junior dog; still keep the dog signal in the quoted teaser line when it mirrors **`description`**.

**Footer explainer:** **`layouts/partials/raymond-project-about.html`** when the category is present; copy from **`raymondProjectAbout*`** in **`i18n/en.toml`** / **`i18n/es.toml`**.

**LinkedIn:** Same fold-first image-only pattern as Cube-Cows (quoted **`description`** on line one; **`🟣`** + **`{heading_code}: Everyone ❤️ Raymond 🐕📓`**; paste **`raymondProjectAboutTitle`** then **`raymondProjectAboutBody`**; hub **`/categories/raymond/`**). Line-one teaser MUST carry the workplace tension of that episode (office joke first; dog signal still required in the same teaser). See **`.cursor/skills/site-linkedin-post/SKILL.md`**.

### `type: sayings`

- **`description`**: Short teaser for cards. **Por-Estas-Calles:** MUST set per **Por-Estas-Calles card teaser** from **`title`**, **`tldr`**, and **`fluff`**. **All other hubs** using **`type: sayings`** (T-Shirt Art, Reptilocracy, Pawtropolis, …): MUST draft new or rewritten teasers via **Gemma teaser**.
- **`tldr`**: Main “TLDR” block (shown in list and on the single layout). When applying this skill, MUST include restrained **`**bold**`** per **`.cursor/skills/site-revise-emphasis/SKILL.md`**.
- **`fluff`**: “Context” block (optional second column on list; shown on single). When present and when applying this skill, MUST include restrained **`**bold**`** per **revise-emphasis**.
- **`project`** (optional): Recurring **series line** on the detail hero. For **Por-Estas-Calles** / Street Wisdom, use the canonical **`Street-Wisdom 💬🇻🇪`** (speech balloon + Venezuelan flag); keep **`ArepaContigo`** in **`tags`** only. **`title`** should be the **unique** episode name (Spanish saying, etc.) for lists, prev/next, and the tab. Without `project`, the layout uses one-line `heading_code` + `title` everywhere.

Archetype: **`archetypes/sayings.md`**.

### Por-Estas-Calles (Venezuelan sayings / Street Wisdom)

**Por-Estas-Calles** is the Street Wisdom hub (often called “sayings” in conversation). Posts usually use **`type: sayings`** with **`categories`**: **`Cognitive-Memetics`** and **`Por-Estas-Calles`** (shareable hub at `/categories/por-estas-calles/`). For the detail hero, set **`project: Street-Wisdom 💬🇻🇪`** on every bundle; **`linkedin.txt`** uses the **fold-first** Por-Estas-Calles layout (quoted saying, then **`❓`** / **`AND:`** / **`🔤 IN ENGLISH`**, then **`🟣`** + **`{heading_code}: {project}`**; see **`.cursor/skills/site-linkedin-post/SKILL.md`** → **Sayings / Street Wisdom**).

**Teaser:** MUST use **Por-Estas-Calles card teaser** (derived from **`title`** / **`tldr`** / **`fluff`**). MUST NOT use **Gemma teaser** for this hub.

**Canonical series explainer** (site "But why" footer, Substack paste, LinkedIn **`🟣`** series block): **`sayingsProjectAboutTitle`**, **`sayingsProjectAboutP1`**, **`sayingsProjectAboutP2`** in **`i18n/en.toml`** and **`i18n/es.toml`**. **MUST NOT** duplicate or paraphrase that body in this skill.

**LinkedIn hashtags** (use on LinkedIn with `#`; in Hugo front matter use **PascalCase** and **no** `#` character):

| LinkedIn | Hugo `tags` value |
|----------|-------------------|
| `#StreetWisdom` | `StreetWisdom` |
| `#CulturalStopwatch` | `CulturalStopwatch` |
| `#TakeBackYourMcDonaldsCulture` | `TakeBackYourMcDonaldsCulture` |
| `#ArepaContigo` | `ArepaContigo` |

For Venezuelan saying posts, **include those four tags** when promoting the series, plus **`VenezuelanSayings`** and optional **post-specific** tags (one or two hooks for that saying). Keep total tags roughly **five to seven** if you add both series tags and a hook.

**LinkedIn post format:** Use **`.cursor/skills/site-linkedin-post/SKILL.md`** (sayings / Street Wisdom **fold-first layout**): quoted **`title`** only on line one; then **`❓`** + `tldr`, **`AND:`** + `fluff`, **`🔤 IN ENGLISH`**, **`🟣`** + **`{heading_code}: Street-Wisdom 💬🇻🇪`** + series paragraphs from **`sayingsProjectAboutP1`** / **`P2`**, hashtags, **`🖋️ Full post (site) →`** (ES then EN when bilingual), **`🔗 Por-Estas-Calles (English) →`** hub URL. Reference: `content/cognitive-memetics/sayings/2026-06-01-saying-21/linkedin.txt`.

### Reptilocracy

Posts under **`content/cognitive-memetics/reptilocracy/`** use **`categories`**: **`Cognitive-Memetics`** and **`Reptilocracy`**. Set **`project: Reptilocracy 🦎🏛️`**.

**Gemma episode copy (`description` / `tldr` / `fluff`):** MUST use **`.cursor/skills/site-reptilocracy-gemma/SKILL.md`** (satire-first brief, thinking mode, psych-fitness knife, three teaser picks). MUST NOT use the generic **Gemma teaser** flow below for this hub when drafting or rewriting those fields.

**Voice:** New or rewritten Gemma fields follow that skill (satire first). Agent clinical edits and older clinical gold (for example `2026-08-30-stability-is-a-setting`) still follow **Reptilocracy voice** under **Hybrid prose** (mechanism, not rebel heat or panel caption). LinkedIn: **`.cursor/skills/site-linkedin-post/SKILL.md`** → Reptilocracy. Footer explainer + petition CTA: **Theme and style** → **Project "But why" explainer cards**.

### T-Shirt Art

Posts in this line use **`categories`**: **`Cognitive-Memetics`** and **`T-Shirt Art`**. The shareable hub lists everything with that category (URL slug is generated by Hugo from the label).

- **`type`:** Prefer **`sayings`** when you want **Teaser** / **TLDR** / **Context** plus a featured image of the graphic. Use **`panel`** when the piece is a longer essay with the art as hero.
- **`project`:** Set a recurring **series line** on the detail hero (for example **`T-Shirt Art`** or a short branded label). Match the voice you want on the card; **`title`** stays the unique episode name.
- **`description` (Teaser):** MUST draft new or rewritten teasers via **Gemma teaser** (not **Por-Estas-Calles card teaser**). Keep it **short**; **`project`** and **`categories`** already show the series; MUST NOT open with redundant meta like “T-Shirt Art piece,” “this post,” or “in this entry.” Jump straight into substance.
- **`tldr`:** Same **Sayings emphasis** rules as other **`type: sayings`** posts (**`.cursor/skills/site-revise-emphasis/SKILL.md`**). MUST NOT change the author’s words; only add, trim, or adjust **`**bold**`** on existing text.
- **`title`:** Unique **episode** name; **two-beat** parallels (for example *cheap X, expensive Y* or *prepare, then advance*) read well as a line without extra labels.
- **`tags`:** Include **`TShirtArt`** plus **two to four** post-specific hooks (themes, mood, format). Do **not** add the **Street Wisdom** LinkedIn set unless the post is also part of that project.
- **Footer explainer:** The Venezuelan **“But why”** block on **`type: sayings`** singles appears only when the post includes the **`Por-Estas-Calles`** category (`layouts/sayings/single.html`). **T-Shirt Art** sayings do not show that block.
- **Facebook friends copy:** **`.cursor/skills/site-facebook-post/SKILL.md`** (not LinkedIn shape); see **Sm(art) / T-Shirt Art (Facebook)** there.

## Layouts and list UI

- Singles: **`layouts/panel/single.html`**, **`layouts/sayings/single.html`**
- Section / home rows: **`layouts/partials/seven-style-row.html`** (claims vs sayings vs default teaser columns)

### Detail singles (parity with social-protocols / `type: claims`)

Single pages for **`panel`** and **`sayings`** follow the same UX pattern as **`layouts/claims/single.html`**:

- **`tags`**: Shown **under the featured image** (not duplicated in the article footer). Footer tag list is suppressed for these types via **`layouts/partials/single/footer.html`**.
- **`type: panel`**: When **`description`** is set, the detail page renders **Teaser** (`#section-teaser`). **Article** (`#section-essay`) appears only when the markdown body is non-empty; merged **Contents** use **`layouts/partials/cows-table-of-contents.html`** when there is a body with headings. With **`project`** set, the detail **h1** uses **`heading-series-single.html`**; lists use **`heading-series-list.html`** (via **`heading-title-markup.html`**).
- **`type: sayings`**: Same **`project`** / **`title`** split as **`panel`** when **`project`** is set (shared **`heading-series-*.html`** partials). Additionally:
  - **Teaser** (`description`), **TLDR** (`tldr`), and **Context** (`fluff`) are rendered as sections with stable anchors **`#section-teaser`**, **`#section-tldr`**, **`#section-context`** when present. Front-matter strings are **markdownified** like claims `description`.
  - If any of those blocks exist **and** the post has a non-empty body, the main markdown is introduced by an **Article** section with anchor **`#section-essay`** (i18n key **`sayingsArticle`**).
  - **Contents** is a **merged** nav: **`layouts/partials/sayings-table-of-contents.html`** combines Teaser / TLDR / Context links with body heading links from the built TOC; when Article wraps the body, in-body heading links nest under **Article** (same idea as **Thoughts** nesting on claims).

### Body extensions (`type: sayings`)

- Use `###` headings for body extensions (e.g. `### In English`, `### Why the devil, why the underwear`).
- Keep extensions concise: focus on cross-language mapping (English equivalents) and specific nuances (e.g. "forsaken" vs "far") rather than general cultural rhetoric.
- Use blockquotes (`>`) for lists of equivalent expressions, including the introductory line (e.g. `> Closest English equivalents include:`).
- MUST NOT use em dashes (—) in body copy; use semicolon or colon instead.

**Body complements, does not repeat:**
- The body should **add depth**, not restate what `tldr` and `fluff` already said.
- If `tldr` explained the structure (e.g., "either X or Y"), the body should explore **consequences** or **cultural context**, not repeat the binary structure.
- Eliminate "process-explaining" phrases: "This means that...", "The metaphor shows that...", "This exaggeration underscores...". Go straight to the impact.
- Plan **vocabulary variation**: if `tldr` uses "trap", `fluff` or body should use a different term ("bind", "dilemma", "squeeze") or a different angle on the same mechanism.

For **`type: claims`** (Claim / Thoughts / Grounding), authoring rules stay in **`.cursor/skills/site-claims-content/SKILL.md`**.

## Bundle layout

Use one folder per post with **`index.md`** and assets beside it, for example:

`content/cognitive-memetics/2006-04-02-cow-w06/index.md` + `agi.png`

## Theme and style

- Site overrides: **`assets/css/_custom.scss`** (prefer not editing `themes/LoveIt/`).
- LoveIt how-tos: **`themes/LoveIt/exampleSite/content/posts/`** (see **`.cursor/rules/site-always-rules-3-hugo.mdc`** index).

### Project "But why" explainer cards (detail footers)

- **Shared look (site-wide):** After the article on **`type: panel`**, **`type: sayings`** (when the Street Wisdom partial shows), **Reptilocracy** singles, **Pawtropolis-Under-Fire** singles, and **Raymond** singles, the **"But why"** explainer is a **gradient card** (warm wash, left accent stripe, soft shadow, circular **❓** mark, slightly roomier body type). Styles live under **`assets/css/_custom.scss`** for **`.cow-project-about`** and **`.sayings-project-about`** together (same shell). Partials: **`layouts/partials/cows-project-about.html`**, **`layouts/partials/sayings-project-about.html`**, **`layouts/partials/reptilocracy-project-about.html`** (Reptilocracy also adds **`reptilocracy-project-about`** for extra rules only), **`layouts/partials/pawtropolis-project-about.html`**, **`layouts/partials/raymond-project-about.html`**.
- **MUST NOT** add duplicate or conflicting card chrome for these explainers in other SCSS files or inline styles unless the user explicitly asks for an exception; extend the shared block in **`_custom.scss`** so Cube-Cows, Street Wisdom, Reptilocracy, Pawtropolis, and Raymond stay visually aligned.
- **Reptilocracy-only:** After **`reptilocracyProjectAboutBody`** (markdown), **`layouts/partials/reptilocracy-project-about.html`** renders a small **CTA row**: **`reptilocracyProjectAboutCtaTitle`** as a **`span`** (not a paragraph, so it lines up cleanly with the pill) plus **`reptilocracyProjectAboutCtaButton`**; the petition URL lives in that partial. Companion styles use **`.reptilocracy-project-about__cta*`** in **`_custom.scss`**. Do not reuse that CTA pattern on Cube-Cows or Street Wisdom explainers.
- **Copy source (canonical; do not paraphrase in this skill):** Tales from the Cube Farm → **`cowsProjectAbout*`**; Street Wisdom → **`sayingsProjectAbout*`**; Reptilocracy → **`reptilocracyProjectAboutBody`** plus **`reptilocracyProjectAboutCtaTitle`** / **`reptilocracyProjectAboutCtaButton`**; Pawtropolis → **`pawtropolisProjectAboutTitle`** / **`pawtropolisProjectAboutBody`**; Raymond → **`raymondProjectAboutTitle`** / **`raymondProjectAboutBody`**. All live in **`i18n/en.toml`** / **`i18n/es.toml`**. LinkedIn series blocks: **`.cursor/skills/site-linkedin-post/SKILL.md`**.

## References in this repo

- **`hugo.toml`**: menu entry and `contentSections` for the home feed.
- **Emphasis (Markdown bold):** **`.cursor/skills/site-revise-emphasis/SKILL.md`**
- Examples: **`content/cognitive-memetics/2006-04-02-cow-w06/`** (`type: panel`), **`content/cognitive-memetics/2006-04-06-saying-13/`** (`type: sayings`).
