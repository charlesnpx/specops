---
id: phase-2-mise-installer-contract
title: Phase 2: Mise Installer Contract
doc_type: implementation_phase
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Phase 2: Mise Installer Contract

## Goal

Implement `install-skill.sh` and `specops install-skill` plan/install/uninstall with staged install root.

## Required work

- Implement the smallest coherent slice for this phase.
- Add unit or integration tests for every public behavior.
- Update relevant interface docs if behavior changes.
- Update `docs/versions/release_decision_log.md` with major scope decisions.

## Acceptance criteria

All four normative installer invocations produce valid JSON; install paths are absolute and rooted correctly.

## Out of scope

Work from later phases unless it is necessary to complete this phase safely.
