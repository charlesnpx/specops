---
id: ADR-0013
title: Skill-led interactive operator loop
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0013: Skill-led interactive operator loop

## Status

Accepted.

## Context

SpecOps needs an interactive operator experience without turning the CLI into a full wizard or terminal UI. The operator should be able to inspect run state, review compiled context, answer a bounded set of questions, and then let the agent continue only after the answers have been recorded as durable provenance.

The existing decision history already requires run artifacts before patches, human decisions before canonical mutation, and no autonomous canonical mutation without `apply`.

## Decision

SpecOps will implement the interactive operator loop as a CLI + skill hybrid.

The CLI owns durable state and artifacts:

- `specops context <run-id> [--json]` compiles the current run context without mutating state.
- `specops next <run-id> --json` includes gate metadata while preserving `command` and `reason`.
- `specops note <run-id> --stage <stage> --text <file-or-inline>` records operator guidance under `.specops/runs/<run-id>/prompts/` without advancing status.
- `specops refine`, `specops harden`, and `specops synthesize` accept `--from <file>` to consume human- or agent-authored semantic artifacts while preserving their content and enforcing legal transitions.

The skill owns conversational UX:

- Run mechanical steps when there is only one safe path.
- Pause at semantic gates: refine, harden, synthesize, decisions, and apply.
- Ask at most three open-ended questions per batch.
- Ask a fourth control question: continue, pause, or change direction.
- After each answer batch, record guidance with `specops note`, refresh `specops context`, and then proceed or pause.

## Consequences

Positive:

- Durable provenance no longer depends on chat history alone.
- The CLI remains composable and scriptable.
- Agents get enough context and gate metadata to behave consistently across Claude and Codex.
- Authored semantic artifacts can be preserved exactly while still flowing through the state machine.

Negative / tradeoffs:

- The loop requires agents to perform an extra `context` and `note` step around semantic work.
- Skill payloads and release assets must be kept synchronized.
- The CLI now exposes more structured context, which must remain backward compatible.

## Affected docs

- docs/interfaces/cli_commands.md
- docs/interfaces/skill_payload_contract.md
- docs/product/user_flows.md
- docs/versions/release_decision_log.md
