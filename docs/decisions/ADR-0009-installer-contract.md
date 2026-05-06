---
id: ADR-0009
title: Treat mise-en-place installer contract as normative
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0009: Treat mise-en-place installer contract as normative

## Status

Accepted.

## Context

The tool must be installable from `mise-en-place`, whose contract requires stable flags and JSON output.

## Decision

The delegated installer interface is normative v0 API. `install-skill.sh` may be a wrapper around the Go binary, but it must satisfy the exact flag and JSON semantics.

## Consequences

Positive:

- Makes installation reliable.
- Supports plan/stage/install/uninstall.
- Lets `mise-en-place` own collision/backups/state.

Negative / tradeoffs:

- Installer implementation must be tested heavily.
- Breaking changes require semver and migration.

## Affected docs

- docs/interfaces/delegated_installer_json.md
- docs/subsystems/specops_installer.md
