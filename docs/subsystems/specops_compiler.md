---
id: specops-compiler
title: SpecOps Compiler
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_compiler

## Purpose

Compiles accepted decisions and spec deltas into patch plans for canonical docs.

## Inputs

```text
accepted decisions
spec_delta.json
current doc index
ADR index
repo taxonomy
```

## Outputs

```text
patch_plan.json
files_to_create
files_to_update
files_to_supersede
validation checklist
```

## Rules

- ADRs are append-only.
- Canon reflects accepted decisions only.
- Research notes may contain speculative material.
- v0 scope only includes accepted v0 work.
- Post-v0 candidates are grouped by theme.

## Acceptance criteria

- `compile` does not write canonical files.
- `plan` shows exact file-level intent.
- plan review may supersede synthesis before `apply` when accepted decisions are correct but generated content is too thin.
- `apply` writes only accepted/approved patch plan items.
