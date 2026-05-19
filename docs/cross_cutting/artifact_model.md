---
id: artifact-model
title: Artifact Model
doc_type: cross_cutting_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Artifact Model

## Core artifact types

```text
RawInput
InputSummary
RefineryNote
AmbiguityRegister
DecisionQueue
Decision
SpecDelta
PatchPlan
EvalReport
AuditReport
```

## Status values

```text
draft
proposed
accepted
rejected
deferred
compiled
applied
superseded
```

## Lineage

Every artifact should record:

```text
created_by_command
created_at
input_artifact_ids
source_files
schema_version
backend
model/prompt/config when applicable
```

Prompt artifacts recorded by `specops note` also record the semantic `stage` they apply to. The stage metadata is used to enforce semantic production gates before `refine`, `harden`, or `synthesize` can copy authored `--from` artifacts into the run.

Pre-apply synthesis supersession preserves lineage by reclassifying current `SpecDelta` and `PatchPlan` artifact refs to superseded archive paths before writing a replacement current spec delta.
