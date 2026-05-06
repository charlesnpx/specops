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
```

## Decision linkage

Every accepted decision should include:

```yaml
id: decision-001
title: ...
status: accepted
adr_required: true
affected_docs: []
```
