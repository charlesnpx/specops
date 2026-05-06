---
id: run-state-model
title: Run State Model
doc_type: cross_cutting_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Run State Model

## Run identity

A run ID should be stable and filesystem-safe:

```text
run-YYYYMMDD-HHMMSS-slug
```

## Run state fields

```yaml
schema: 1
run_id: run-20260506-143000-permission-model
name: permission-model
status: refined
target_repo: .
created_at: ...
inputs: []
artifacts:
  - id: prompt-001
    type: prompt
    path: prompts/20260506-143100-refine.md
    stage: refine
    created_at: ...
decisions:
  accepted: []
  rejected: []
  deferred: []
next:
  command: specops harden <run-id>
  reason: refined artifact can be challenged or synthesized
  stage: harden
  gate_kind: semantic
  context_command: specops context <run-id>
  note_command: specops note <run-id> --stage harden --text <file-or-inline>
```

## State transitions

Commands can be repeated if idempotent. Non-idempotent commands must create new artifact versions.

Semantic production commands (`refine`, `harden`, and `synthesize`) require a matching prompt artifact with `stage` metadata before they can write output artifacts or advance status.
