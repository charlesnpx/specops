---
id: phase-3-init-scaffold
title: Phase 3: Target Repo Init Scaffold
doc_type: implementation_phase
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Phase 3: Target Repo Init Scaffold

## Goal

Implement `specops init`, scaffold installer, lockfile, AGENTS/CLAUDE generation, and idempotent updates.

## Required work

- Implement the smallest coherent slice for this phase.
- Add unit or integration tests for every public behavior.
- Update relevant interface docs if behavior changes.
- Update `docs/versions/release_decision_log.md` with major scope decisions.

## Acceptance criteria

A fresh repo can be initialized and re-initialized without destructive changes.

## Out of scope

Work from later phases unless it is necessary to complete this phase safely.
