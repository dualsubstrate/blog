---
name: site-psych-fitness-28-linkedin
description: Writes LinkedIn posts for the Psych-Fitness-28 petition campaign. Use when creating social copy for content/cognitive-memetics/psych-fitness-28/ posts, when the user mentions Psych-Fitness-28 LinkedIn, petition campaign posts, or EN9806 social copy.
disable-model-invocation: true
---

# Psych-Fitness-28 LinkedIn Post Generator

## Goal

Generate LinkedIn posts for the Psych-Fitness-28 petition campaign (EN9806: Stronger Checks and Balances for Australia's top leaders). Posts deliver the complete argument on-platform; the petition link is the priority CTA.

## Source Fields

Read from `index.md`:

| Field | Maps to |
|-------|---------|
| `heading_code` | Episode label (e.g., `D6:`) |
| `project` | Series name on same line (`Psych-Fitness-28 🙏`) |
| `title` | Quoted saying on line two |
| `tldr` | **THE FACTS** section (strip `**bold**` for plain text) |
| `fluff` | **THE ASK** section (strip `**bold**`) |
| `tags` | Hashtag line |

## Structure

```
{heading_code}: {project}
"{title}"

✍️ Sign the petition → https://www.aph.gov.au/e-petitions/petition/EN9806

✔️ THE FACTS:
[Plain-text tldr]

➕ THE ASK:
[Plain-text fluff with petition link repeated]

❓ THE STAKES:
[Series framing: Day X of 28, petition context, high-stakes argument]

#{tag1} #{tag2} ...

🧷 Full post (site) →
https://blog.dualsubstrate.com/cognitive-memetics/psych-fitness-28/{slug}/

🔗 28-day series → https://blog.dualsubstrate.com/categories/psych-fitness-28/
```

## Rules

1. **Lead with petition link**: Place `✍️ Sign the petition →` in the first 4 lines, before LinkedIn's "see more" fold.

2. **Plain text only**: Strip all Markdown (`**bold**`, backticks) from `tldr` and `fluff`. LinkedIn does not render Markdown.

3. **Fact-checked phrasing**: Use the tightened, defensible language from the `tldr` (e.g., "at least seven ice warnings," "near full speed," "subtle pressure").

4. **No Spanish link**: This campaign targets an Australian petition. Omit the `- ES:` line. Use single English URL format.

5. **Labels**: Use campaign-appropriate section labels:
   - `✔️ THE FACTS:` (historical story and evidence)
   - `➕ THE ASK:` (petition call to action, link repeated)
   - `❓ THE STAKES:` (series context and political argument)

6. **Closing**: End with `🧷 Full post (site) →` and `🔗 28-day series →` (not "Psych-Fitness-28 (English)").

7. **Hashtags**: Include all tags from front matter in original order, formatted as `#TagName`.

## THE STAKES Template

Use this framing for every post:

> This is Day X of 28 in a series counting down the human factors that shape high-stakes decisions. Cabinet ministers make choices affecting millions, yet they rarely face the fitness-for-duty scrutiny that pilots and surgeons accept as normal. The petition asks for stronger checks and balances on psychological fitness for Australia's top leaders because the stakes deserve the same rigor we demand in cockpits and operating rooms.

Adjust the Day number to match `heading_code`.

## Output

Save as `linkedin.txt` in the same bundle folder as `index.md`.

Example: `content/cognitive-memetics/psych-fitness-28/2026-05-07-day-06-pilots-and-surgeons-accept-fitness-checks/linkedin.txt`
