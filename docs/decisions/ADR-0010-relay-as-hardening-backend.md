---
id: ADR-0010
title: Use convo-relay as an adversarial hardening backend
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0010: Use convo-relay as an adversarial hardening backend

## Status

Accepted.

## Context

The Context Sink spec improved significantly after relay-based adversarial stress testing.

## Decision

SpecOps should support `convo-relay` as an optional backend for `specops harden`, while keeping relay integration replaceable.

## Consequences

Positive:

- Captures an observed high-value process primitive.
- Works with the user's existing toolchain.
- Helps discover seams before canonicalization.

Negative / tradeoffs:

- Requires subprocess integration and transcript parsing.
- Relay availability may vary by user environment.

## Affected docs

- docs/subsystems/specops_relay.md
- docs/research/context_sink_process_reconstruction.md
