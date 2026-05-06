---
id: ADR-0008
title: Use typed artifacts and JSON Schemas
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0008: Use typed artifacts and JSON Schemas

## Status

Accepted.

## Context

A one-prompt system is brittle. The CLI needs parseable, validatable intermediate representations.

## Decision

Use JSON/YAML artifacts validated by schemas for run state, spec deltas, decisions, eval reports, and delegated installer reports.

## Consequences

Positive:

- Enables deterministic validation.
- Supports CLI and CI automation.
- Makes agent outputs checkable.

Negative / tradeoffs:

- Requires schema evolution and migration.
- Some semantic quality still needs human/LLM eval.

## Affected docs

- docs/cross_cutting/artifact_model.md
- schemas/*.schema.json
