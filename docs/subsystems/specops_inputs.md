---
id: specops-inputs
title: SpecOps Input Normalizer
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_inputs

## Purpose

Normalizes raw input files, chat transcripts, relay transcripts, existing specs, and gold repos into run inputs.

## Input types

```text
raw_markdown
conversation_transcript
relay_transcript
existing_spec_repo
gold_repo
product_brief
process_trace
```

## Outputs

```text
source_manifest.yaml
input_summary.md
conversation_segments.json
artifact_inventory.json
```

## Acceptance criteria

- Large conversations can be sliced into topic/time segments.
- Relay transcript import extracts settled/contested/withdrawn points where present.
- Gold repo mining never contaminates clean reproduction fixtures unless explicitly requested.
