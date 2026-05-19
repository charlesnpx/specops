---
id: ADR-0019
title: Pre-apply synthesis supersession
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-07
---

# ADR-0019: Pre-apply synthesis supersession

## Status

Accepted.

## Context

ADR-0018 exposed `patch_items` as the high-fidelity content path for synthesized spec deltas. During apply review, an operator may discover that a compiled plan is mechanically healthy but semantically too thin because the previous synthesis did not include rich authored `patch_items`.

The normal `synthesize` command should remain a forward semantic transition from `refined` or `hardened`. Loosening it to run from `planned` would blur the state machine and make it unclear whether accepted decisions are being reopened.

## Decision

SpecOps will add `specops supersede-synthesis <run-id> --from <spec_delta.json> [--reopen-decisions]`.

The command is allowed only from `planned` and requires an `apply` stage note. It archives the current `outputs/spec_delta.json` and `patches/patch_plan.json`, reclassifies old current artifact refs as superseded refs, removes the current patch plan, and writes the replacement spec delta as the current `outputs/spec_delta.json`.

By default, settled decisions remain settled. Replacement deltas may omit decisions or repeat existing decision IDs with unchanged substance, but they may not introduce new decision IDs or change settled decision substance. The run returns to `decisions_accepted`.

When `--reopen-decisions` is passed, new or changed decisions from the replacement delta are merged into run state and the run returns to `awaiting_decisions`.

## Consequences

Positive:

- Apply review can reject thin patch content without losing run lineage.
- Agents can refine content repeatedly before apply while preserving settled decisions.
- The normal synthesis transition remains strict.

Negative / tradeoffs:

- Runs may contain superseded artifact archives that need to be understood by context readers.
- Decision-change detection is conservative and requires an explicit reopen path.

## Affected docs

- docs/CANON.md
- docs/versions/v0_scope.md
- docs/interfaces/cli_commands.md
- docs/interfaces/run_directory_contract.md
- docs/interfaces/skill_payload_contract.md
- docs/cross_cutting/run_state_model.md
- docs/cross_cutting/artifact_model.md
