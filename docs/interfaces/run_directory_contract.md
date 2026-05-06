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
  outputs/
    input_summary.md
    refinery_note.md
    ambiguity_register.yaml
    decision_queue.yaml
    spec_delta.json
    patch_plan.json
  traces/
    backend_result.json
  patches/
    files/
  evals/
    eval_report.json
    eval_report.md
```

No file in `outputs/` is canonical until applied or promoted.
