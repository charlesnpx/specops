---
id: agent-backend-model
title: Agent Backend Model
doc_type: cross_cutting_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Agent Backend Model

## Backend interface

```text
run(prompt_packet, cwd, output_schema, mode) -> agent_result
```

## Modes

```text
intake
refine
harden
synthesize
compile_patch
audit
eval
```

## Backend requirements

- Capture stdout/stderr/transcript paths.
- Validate expected output artifacts.
- Never assume backend success until artifacts validate.
- Allow manual fallback.
