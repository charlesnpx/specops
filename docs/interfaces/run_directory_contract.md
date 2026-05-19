---
id: run-directory-contract
title: Run Directory Contract
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Run Directory Contract

```text
.specops/runs/<run-id>/
  run.yaml
  inputs/
    source_manifest.yaml
    raw/
    normalized/
  prompts/
    prompt_packet.json
    <timestamp>-<stage>.md
  outputs/
    input_summary.md
    refinery_note.md
    ambiguity_register.yaml
    decision_queue.yaml
    spec_delta.json
    superseded/
      <timestamp>-spec_delta.json
  traces/
    backend_result.json
  patches/
    patch_plan.json
    superseded/
      <timestamp>-patch_plan.json
    files/
  evals/
    eval_report.json
    eval_report.md
```

No file in `outputs/` is canonical until applied or promoted.

Prompt Markdown files created by `specops note` are durable gate guidance. The run state artifact entry records the prompt `stage`; legacy prompt paths with `<timestamp>-<stage>.md` remain stage-checkable when metadata is absent.

When pre-apply synthesis is superseded, old current output refs are reclassified to the archived superseded paths instead of being silently overwritten.
