---
id: ADR-0002
title: Implement the core CLI in Go
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0002: Implement the core CLI in Go

## Status

Accepted.

## Context

The tool is a local CLI that must stage installs, hash files, parse registries, run subprocesses, package release binaries, and work on macOS/Linux. `mise-en-place` is also Go-based.

## Decision

The v0 CLI should be implemented in Go, with release binaries for macOS and Linux, arm64 and x86_64.

## Consequences

Positive:

- Self-contained binaries.
- Easy GitHub Release distribution.
- Good fit for filesystem/installer operations.
- Aligns with the existing `mise-en-place` ecosystem.

Negative / tradeoffs:

- Less flexible than Python for rapid LLM/data experimentation.
- Some agent/eval logic may be easier to prototype separately.
- Requires careful design to avoid overengineering early.

## Affected docs

- docs/architecture/system_overview.md
- docs/research/implementation_language_analysis.md
- docs/implementation/phase_1_go_cli_skeleton.md
