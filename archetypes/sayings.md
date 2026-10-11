+++
date = '{{ .Date }}'
draft = false
title = '{{ replace .File.ContentBaseName "-" " " | title }}'
type = 'sayings'

# Umbrella + Por-Estas-Calles project hub. Remove Por-Estas-Calles if this post is not that series.
categories = ['Cognitive-Memetics', 'Por-Estas-Calles']

# Optional label before the title (e.g. week id); styled via .heading-code in _custom.scss
heading_code = ''

# Series line on the detail hero (canonical: Street-Wisdom 💬🇻🇪 for Por-Estas-Calles). `title` = unique episode name (lists, prev/next). Leave empty for one-line titles.
project = 'Street-Wisdom 💬🇻🇪'

# Used by layouts/sayings/single.html (section list shows TLDR + Context)
tldr = ''
fluff = ''
+++

