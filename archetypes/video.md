+++
date = '{{ .Date }}'
draft = false
title = '{{ replace .File.ContentBaseName "-" " " | title }}'
type = 'video'

# Short pitch for list cards and the lead on the page (Markdown ok).
description = ''

# From the watch URL: …youtube.com/watch?v=THIS_PART
youtube_id = ''

subtitle = ''

# Taxonomy hubs (Mind-Infrastructure, Human-Condition, Social-Protocols, etc.); pick folder under content/ by main job (see .cursor/rules/site-content-placement.mdc).
categories = ['Mind-Infrastructure']

tags = []

# Optional keep-reading paths (Hugo GetPage). Prefer sayings/panel; MAY add claims/video.
# Example: related = ['/cognitive-memetics/sayings/2026-10-05-saying-39/']
# Empty or omitted: layout fills from shared tags. Do not paste related links into the body.
related = []

# Reader landing after the TL;DW (why it matters, takeaways, explore links).
# reader_landing = { why_it_matters = '', takeaways = [], explore = [] }
+++

TLDR / so-what article below the embed (summary for readers who skip the watch). Alternatively, leave `youtube_id` empty and add the built-in YouTube shortcode in Markdown.
