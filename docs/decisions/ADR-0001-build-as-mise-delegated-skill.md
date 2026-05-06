---
id: ADR-0001
title: Build SpecOps as a mise-en-place delegated skill/tool
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0001: Build SpecOps as a mise-en-place delegated skill/tool

## Status

Accepted.

## Context

The user wants the tool installable through `mise-en-place`, whose delegated repo contract can orchestrate external skill/tool repositories across Claude and Codex targets.

## Decision

SpecOps v0 will be packaged as a delegated repo compatible with the `mise-en-place` installer contract. It will install Claude payloads, Codex payloads, and a `specops` CLI binary through the `tools` target.

## Consequences

Positive:

- Integrates with the user's existing skill manager.
- Reuses `mise-en-place` collision, backup, state, and private-delegated-skill behavior.
- Keeps SpecOps independent while installable from a registry entry.

Negative / tradeoffs:

- Requires a stable `install-skill.sh` wrapper.
- Forces the installer JSON contract to be treated as a public API.

## Affected docs

- docs/cross_cutting/mise_en_place_contract.md
- docs/interfaces/delegated_installer_json.md
- docs/subsystems/specops_installer.md
