+++
date = '{{ .Date }}'
draft = false
title = '{{ replace .File.ContentBaseName "-" " " | title }}'
type = 'opinion'

# Optional label before the title; styled via .heading-code in _custom.scss
heading_code = ''

# List card and single "Opinion" band (type: opinion).
description = ''

# Optional shorter copy for Open Graph, Twitter, and meta Description. Omit to use `description`.
og_description = ''

# Optional citations or constructs (Markdown). Legacy key `paper` still works. Not required for opinion.
grounding = ''

# Optional freeform dig-deeper markdown (legacy). Prefer structured reader_landing below.
research = ''

# Reader landing: why it matters, takeaways, curated explore links (see layouts/partials/reader-landing.html).
# reader_landing = { why_it_matters = '', takeaways = [], explore = [] }

# Optional keep-reading paths (Hugo GetPage). Prefer sayings/panel; MAY add claims/video/opinion.
# Example: related = ['/cognitive-memetics/sayings/2026-10-05-saying-39/']
# Empty or omitted: layout fills from shared tags. Do not paste related links into the body.
related = []

# Optional Markdown attribution: detail meta above hero; section list under thumbnail. Omit if not needed.
image_credit = ''

# Social preview image (Open Graph / Twitter). Use a page-bundle file name.
# Keep this aligned with `featuredImage` when you add one in the post front matter.
images = []

# Optional page resources. If you set a featured image, name it `featured-image`
# so theme JSON-LD and other partials can find it consistently.
[[resources]]
src = ''
name = 'featured-image'
+++
