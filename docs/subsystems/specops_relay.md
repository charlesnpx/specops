---
id: specops-relay
title: SpecOps Relay Integration
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_relay

## Purpose

Integrates `convo-relay` as an optional adversarial hardening backend.

## Responsibilities

- Build task briefs from run state.
- Select context files.
- Run relay in adversarial/steelman/cooperative mode.
- Import transcript.
- Extract settled, contested, withdrawn, and key insights.
- Produce hardening deltas.

## Acceptance criteria

- `specops harden --backend convo-relay` creates a hardening report.
- Relay failure leaves the run recoverable.
- Relay outputs are provenance, not canonical truth until accepted.
