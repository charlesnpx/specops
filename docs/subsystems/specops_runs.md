---
id: specops-runs
title: SpecOps Run Manager
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_runs

## Purpose

Manages stateful runs that convert input material into draft artifacts and eventually canonical patch plans.

## Owned responsibilities

- Create run IDs.
- Maintain run state file.
- Track inputs, outputs, decisions, backend calls, and next commands.
- Store prompt packets and traces when configured.
- Promote selected run artifacts into `docs/research/refinery/`.

## Run directory

```text
.specops/runs/run-YYYYMMDD-HHMMSS-slug/
  run.yaml
  inputs/
  prompts/
  outputs/
  traces/
  evals/
  patches/
```

## Statuses

```text
created
ingested
intake_complete
refined
hardened
synthesized
awaiting_decisions
decisions_accepted
compiled
planned
applied
audited
evaluated
```

## Acceptance criteria

- `specops next <run>` can recommend the next legal step.
- Illegal state transitions fail clearly.
- Run state is machine-readable and schema-valid.
