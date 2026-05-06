---
id: ADR-0003
title: Use a Cobra-style command graph for the Go CLI
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0003: Use a Cobra-style command graph for the Go CLI

## Status

Accepted.

## Context

The CLI has many subcommands with shared flags, local flags, completions, and documentation needs.

## Decision

Use Cobra or an equivalent command-tree framework for v0 unless implementation reveals it is too heavy.

## Consequences

Positive:

- Natural subcommand model.
- Mature flag behavior.
- Supports generated docs and shell completions.
- Familiar to Go CLI users.

Negative / tradeoffs:

- Adds dependency and command framework conventions.
- Small commands can become overly nested if not disciplined.

## Affected docs

- docs/interfaces/cli_commands.md
- docs/subsystems/specops_cli.md
