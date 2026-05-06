---
id: ADR-0006
title: Use run artifacts before canonical patches
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0006: Use run artifacts before canonical patches

## Status

Accepted.

## Context

Spec production has drafts, ambiguity registers, decisions, and patch plans before accepted canonical changes.

## Decision

Dynamic commands write run artifacts under `.specops/runs/<run-id>/` first. Canonical patches happen only through `specops apply`.

## Consequences

Positive:

- Safer review.
- Enables eval and re-run.
- Makes output inspectable before mutation.

Negative / tradeoffs:

- Adds state machinery.
- Requires cleanup policy for run dirs.

## Affected docs

- docs/cross_cutting/run_state_model.md
- docs/subsystems/specops_runs.md
- docs/interfaces/run_directory_contract.md
