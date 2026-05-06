---
id: ADR-0005
title: Separate .specops control plane from repository proper
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0005: Separate .specops control plane from repository proper

## Status

Accepted.

## Context

The process scaffold and generated specification could be confused if both are Markdown files in the repo.

## Decision

`.specops/` is the local control plane. The actual produced specification lives in `docs/` and repository-proper files.

## Consequences

Positive:

- Clear mental model.
- Prevents process logs from polluting canonical truth.
- Allows `.specops/runs/` to be ignored or selectively promoted.

Negative / tradeoffs:

- Users must learn a two-layer repo model.
- Some artifacts require explicit promotion from run scratch to docs/research.

## Affected docs

- docs/cross_cutting/target_repo_control_plane.md
- docs/interfaces/run_directory_contract.md
