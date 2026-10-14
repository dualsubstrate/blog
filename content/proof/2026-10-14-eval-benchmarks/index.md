---
title: "DSS-EVAL: Verified Recall and Abstention Precision"
date: '2026-10-14T08:00:00+11:00'
draft: false
type: claims
featuredImage: 'featured-image.jpg'
resources:
  - name: featured-image
    src: featured-image.jpg
categories:
  - Proof
tags:
  - Benchmarks
  - Precision
  - Abstention
description: >-
  Public benchmark results demonstrating perfect abstention precision and high multi-hop recall on small ledgers, along with honest documentation of current gate limits.
---

Governed memory is only as trustworthy as the verifiable tests behind it. Rather than making sweeping claims about universal intelligence, DSS publishes exact micro-benchmarks on the open `DSS-EVAL` suite.

Key metrics from the latest published run:

- **Perfect Abstention Precision (1.000)**: Across 300 test queries per seed, the retrieval router abstained from answering whenever source truth was absent, avoiding hallucinations.
- **Sub-millisecond Routing**: Retrieval-router lookups consistently measure under 1 ms across small ledgers with hundreds of versioned entries.
- **Honest Limitations**: The evaluation also highlights what is not yet finished, including our two open integrity gates and current 0.14 recall@1 ranking track on dense ablated benchmarks.
