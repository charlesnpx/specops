---
id: spec-delta-interface
title: Spec Delta Interface
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Spec Delta Interface

A `SpecDelta` is the primary intermediate representation between refinement and patch planning.

## Required fields

```yaml
schema: 1
run_id: run-...
source_summary: ...
new_concepts: []
requirements: []
constraints: []
assumptions: []
ambiguities: []
options: []
recommendations: []
decisions: []
affected_docs: []
version_scope_changes: []
implementation_phase_changes: []
acceptance_criteria: []
open_questions: []
risks: []
patch_plan: []
patch_items: []
```

List fields in the Go `SpecDelta` struct are string lists unless otherwise noted. `affected_docs` is a list of relative doc paths, not objects.

`patch_plan` is a human-readable string list. `patch_items` is the optional exact file-level patch form used when the author wants compile to preserve canonical document content exactly:

```yaml
patch_items:
  - action: create
    path: docs/CANON.md
    title: Create canonical frame
    content: |
      ---
      id: canon
      ...
    decision_ids: [D001]
```

When `patch_items` is absent, compile may generate deterministic canonical doc patches for accepted `affected_docs` from the structured delta fields and accepted decisions.

## Decision linkage

Every accepted decision should include:

```yaml
id: decision-001
title: ...
status: accepted
adr_required: true
affected_docs: []
```
