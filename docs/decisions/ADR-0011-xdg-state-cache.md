---
id: ADR-0011
title: Use XDG-style state/cache locations for user-level state
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0011: Use XDG-style state/cache locations for user-level state

## Status

Accepted.

## Context

The CLI has user-level state, cache, and lock files, and should avoid dotfile sprawl while aligning with Unix-like conventions.

## Decision

Use XDG defaults for config/cache/state where possible, with project overrides and `SPECOPS_HOME` for testing.

## Consequences

Positive:

- Predictable user directories.
- Easier testing.
- Cleaner separation of cache/state/config.

Negative / tradeoffs:

- macOS users may expect `~/Library` conventions.
- Requires documented env vars.

## Affected docs

- docs/subsystems/specops_state.md
- docs/cross_cutting/security_and_filesystem.md
