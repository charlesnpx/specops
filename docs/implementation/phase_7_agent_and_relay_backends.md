---
id: phase-7-agent-and-relay-backends
title: Phase 7: Agent and Relay Backends
doc_type: implementation_phase
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Phase 7: Agent and Relay Backends

## Goal

Integrate Codex CLI, Claude Code command packets, and convo-relay hardening backend.

## Required work

- Implement the smallest coherent slice for this phase.
- Add unit or integration tests for every public behavior.
- Update relevant interface docs if behavior changes.
- Update `docs/versions/release_decision_log.md` with major scope decisions.

## Acceptance criteria

`specops harden --backend convo-relay` imports settled/contested/withdrawn points into a run.

## Out of scope

Work from later phases unless it is necessary to complete this phase safely.
