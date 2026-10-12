---
title: "DSS hops on one cube"
date: '2026-10-11T01:00:00+11:00'
draft: true
type: claims
categories:
  - Lineage
tags:
  - Dual-Substrate
  - Topology
  - Coordinates
description: >-
  Eight stops on one cube, twenty hop kinds, and a separate roll plus face-flip
  bit. The interactive view below is the same hops canvas we use in Cursor,
  published for the blog.
grounding: >-
  Stops are corners 0–7. Roll is 6 faces down × 4 turns (24). Face flip (FF)
  is a binary on the current face. Resolution on the exhibit:
  `8stp * 6fce * 4trn * 2ff = 384`. Walk kinds are 4 Ring + 4 M + 12 cube
  edges; the store is the hop product along a path, not stops × kinds.
---

This page is the public copy of the DSS hops canvas: Ring and M on the two tetrahedra, twelve cube rim edges, Roll, and FF. Drag the cube, pick a hop pill, and open the Resolution tab for the count equations.

{{< dss-hops >}}

### Stops, rolls, and FF

The **Stops × rolls** column is eight landing points times twenty-four cube rolls. **With FF** multiplies by two for face flip on the face the hop already sits on. That is separate from the twenty walk kinds in the hop table.
