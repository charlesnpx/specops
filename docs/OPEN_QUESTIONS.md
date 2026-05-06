---
id: open-questions
title: Open Questions
doc_type: open_questions
status: proposed
normative: false
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Open Questions

## OQ-001: Final product name

Current working name is `specops`. Alternatives: `specforge`, `specrepo`, `context-forge`, `specwright`.

Recommendation: keep `specops` for v0 unless package-name conflict appears.

## OQ-002: Where should Codex skills be installed?

The contract assumes a Codex skill target, but Codex skill conventions may evolve. v0 should isolate target paths behind an installer target resolver and keep the actual path configurable.

Recommendation: default to the currently expected user-level Codex skill path but allow `.specops.yaml` and installer flags to override.

## OQ-003: Should the CLI directly call model APIs in v0?

Options:

1. CLI only prepares prompt packets and expects Claude/Codex to run them.
2. CLI can call OpenAI/Anthropic APIs for structured extraction.
3. CLI supports both, behind backend adapters.

Recommendation: support both conceptually, but implement manual/file backend plus Codex/Claude subprocess backends first. Add direct API backends after schemas stabilize.

## OQ-004: Should run artifacts be committed?

Recommendation: `.specops/runs/` is ignored by default. Selected refinery notes and eval reports can be promoted into `docs/research/refinery/` and `docs/_generated/evals/`.

## OQ-005: Should `specops` be installed as an external tool in `mise-en-place` or as a delegated repo tool target?

Recommendation: delegated repo installs all three targets: Claude payload, Codex payload, and `tools` binary. `mise-en-place` owns collision checks and state after staging.

## OQ-006: How much of the process scaffold should be vendored into each target repo?

Recommendation: v0 vendors a thin scaffold for readability and agent independence. Later versions can support linked mode with lockfile pinning.
