---
id: ADR-0012
title: No autonomous canonical mutation without apply
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0012: No autonomous canonical mutation without apply

## Status

Accepted.

## Context

SpecOps should help agents work, but not surprise users by changing canonical files.

## Decision

Commands may produce drafts, plans, and patch files. Canonical repo mutation requires `specops apply` or an explicit agent action following the same apply semantics.

## Consequences

Positive:

- Safer.
- Reviewable.
- Better alignment with Git workflows.

Negative / tradeoffs:

- Slightly more commands to run.
- Advanced users may want `--auto-apply`, which should remain opt-in.

## Affected docs

- docs/interfaces/cli_commands.md
- docs/cross_cutting/decision_gate.md
