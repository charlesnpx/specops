---
id: decision-gate
title: Decision Gate
doc_type: cross_cutting_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Decision Gate

The decision gate prevents unaccepted agent recommendations from becoming canon.

## Decision statuses

```text
proposed
accepted
rejected
deferred
amended
superseded
```

## Rules

- `compile --accepted-only` ignores proposed/deferred decisions.
- `apply` refuses to patch canonical docs from unaccepted decisions unless `--include-proposed` is explicitly provided for draft branches.
- Accepted consequential decisions should become ADRs.
- Rejected options may appear in ADR context or research notes.
