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
artifacts: []
decisions:
  accepted: []
  rejected: []
  deferred: []
next:
  recommended:
    - specops decisions <run-id>
```

## State transitions

Commands can be repeated if idempotent. Non-idempotent commands must create new artifact versions.
