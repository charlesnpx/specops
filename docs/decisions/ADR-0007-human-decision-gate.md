---
id: ADR-0007
title: Keep human acceptance as a hard gate
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0007: Keep human acceptance as a hard gate

## Status

Accepted.

## Context

The observed Context Sink process improved because recommendations were accepted deliberately before being encoded into canonical docs.

## Decision

The CLI must model decisions as proposed/accepted/rejected/deferred/amended. Only accepted decisions can become ADRs or canonical patches by default.

## Consequences

Positive:

- Preserves human authority.
- Avoids plausible but unaccepted agent inventions becoming canon.
- Makes traceability explicit.

Negative / tradeoffs:

- Slower than fully automatic doc writing.
- Requires clear CLI UX for decisions.

## Affected docs

- docs/cross_cutting/decision_gate.md
- docs/subsystems/specops_compiler.md
