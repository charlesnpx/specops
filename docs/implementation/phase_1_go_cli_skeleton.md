---
id: phase-1-go-cli-skeleton
title: Phase 1: Go CLI Skeleton
doc_type: implementation_phase
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Phase 1: Go CLI Skeleton

## Goal

Create Go module, root command, command groups, version command, JSON output helpers, and error handling.

## Required work

- Implement the smallest coherent slice for this phase.
- Add unit or integration tests for every public behavior.
- Update relevant interface docs if behavior changes.
- Update `docs/versions/release_decision_log.md` with major scope decisions.

## Acceptance criteria

`specops --help`, `specops version --json`, and core command stubs work on macOS/Linux.

## Out of scope

Work from later phases unless it is necessary to complete this phase safely.
