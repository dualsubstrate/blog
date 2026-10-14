---
title: "Principals and Grants in Governed Memory"
date: '2026-10-12T10:00:00+11:00'
draft: false
type: claims
featuredImage: 'featured-image.jpg'
resources:
  - name: featured-image
    src: featured-image.jpg
categories:
  - Platform
tags:
  - Principals
  - Grants
  - Control-Plane
description: >-
  Every human operator and every runtime agent acts as a named Principal with passkey identity. Grants define exact access windows to versioned ledgers rather than global API keys.
---

In traditional application stacks, agents share broad database credentials and global model keys. When an agent produces an unexpected answer or performs an unintended write, identifying the authorizing context requires digging through disjointed server logs.

In the Dual Substrate System, identity and authorization are structural:

1. **Principals**: Humans, organizations, and autonomous agents are registered as first-class Principals authenticated via passkeys and decentralized identifiers.
2. **Grants**: Permissions to read, append, or bind models to ledgers are issued as discrete, versioned Grants managed in the Control Plane.
3. **Boundaries**: An agent cannot interact with surfaces or memory outside its active grant window.

Because grants are immutable records stored directly in the governed state, accountability is verifiable at the storage layer rather than reconstructed after the fact.
