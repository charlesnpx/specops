---
id: implementation-strategy
title: Implementation Strategy
doc_type: implementation_plan
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Implementation Strategy

Build the tool in layers:

```text
1. Specification repository and fixtures
2. Go CLI skeleton
3. mise-en-place delegated installer contract
4. target repo scaffold init
5. run state and intake/refine artifacts
6. compile/plan/apply/audit
7. reproduction/eval harness
8. agent and relay backends
9. release and registry integration
```

## Language choice

Use Go for v0 implementation. It gives self-contained release binaries, straightforward filesystem operations, a good fit for installer/staging work, and alignment with `mise-en-place`.

## Testing posture

- Unit tests for path mapping, hashing, JSON reports, state transitions, and schema validation.
- Golden tests for installer output.
- Integration tests with temporary install roots.
- Fixture tests using Context Sink source material.
- Contract tests for `--json` stdout purity.
