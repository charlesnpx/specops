---
id: phase-8-release-and-registry
title: Phase 8: Release and Registry
doc_type: implementation_phase
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Phase 8: Release and Registry

## Goal

Add GoReleaser config, GitHub Actions, release assets, checksums, and mise-en-place registry docs.

## Required work

- Implement the smallest coherent slice for this phase.
- Add unit or integration tests for every public behavior.
- Update relevant interface docs if behavior changes.
- Update `docs/versions/release_decision_log.md` with major scope decisions.

## Acceptance criteria

Pushing a semver tag publishes binaries and skill payload archives usable by mise-en-place.

## Out of scope

Work from later phases unless it is necessary to complete this phase safely.
