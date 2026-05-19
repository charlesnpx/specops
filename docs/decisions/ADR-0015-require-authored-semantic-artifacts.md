---
id: ADR-0015
title: Require authored semantic artifacts
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-07
---

# ADR-0015: Require authored semantic artifacts

## Status

Accepted.

Supersedes the deterministic fallback portion of ADR-0014 for `refine`, `harden`, and `synthesize`.

## Context

ADR-0014 made semantic gates operational by requiring a recorded stage note before `refine`, `harden`, or `synthesize` could advance run state.

That prevented a single `continue` answer from crossing multiple gates, but it left a bad path open: after a stage note existed, the CLI could still run its manual deterministic fallback. Because the CLI is not AI-enabled, that fallback only produced thin placeholder content and did not incorporate the operator's content-aware guidance from the note.

The stage note is durable provenance for the gate. It should not be treated as semantic source material that the CLI can transform.

## Decision

`specops refine`, `specops harden`, and `specops synthesize` require both:

- a recorded stage note for the current semantic gate
- an authored artifact passed with `--from`

The CLI will not produce fallback semantic content for these commands. It only copies and validates authored artifacts, then applies legal state transitions.

The next-command metadata for these stages must show the required `--from` form:

```sh
specops refine <run-id> --from <file>
specops harden <run-id> --from <file>
specops synthesize <run-id> --from <spec_delta.json>
```

If `--from` is omitted after the stage note exists, the command fails without writing output artifacts or advancing run status. The error explains that the CLI is not AI-enabled and requires an authored artifact.

## Consequences

Positive:

- Recorded operator guidance cannot be mistaken for semantic output.
- Agents must author the real refined, hardened, or synthesized artifact before advancing the state machine.
- The CLI remains deterministic, non-AI, and honest about what it can produce.

Negative / tradeoffs:

- Fully scripted fallback runs are no longer available for `refine`, `harden`, or `synthesize`.
- Agents and operators need one extra file-authoring step at each semantic production gate.

## Affected docs

- docs/interfaces/cli_commands.md
- docs/interfaces/skill_payload_contract.md
- docs/versions/release_decision_log.md
