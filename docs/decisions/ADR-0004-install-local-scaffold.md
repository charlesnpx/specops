---
id: ADR-0004
title: Install a local .specops scaffold into target repositories
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0004: Install a local .specops scaffold into target repositories

## Status

Accepted.

## Context

If the CLI only prints instructions, the process lives in terminal output or chat history. Target repos need durable local operating rules.

## Decision

`specops init` installs a thin `.specops/` control-plane scaffold into target repositories, plus `AGENTS.md`, `CLAUDE.md`, `.specops.yaml`, and `.specops.lock`.

## Consequences

Positive:

- Target repo becomes self-describing.
- Claude/Codex can discover local process rules.
- Git diffs capture process changes.
- Humans can read the process without the CLI.

Negative / tradeoffs:

- Adds files to target repos.
- Requires upgrade/migration strategy for scaffold versions.

## Affected docs

- docs/cross_cutting/target_repo_control_plane.md
- docs/architecture/target_repo_layout.md
- docs/subsystems/specops_scaffold.md
