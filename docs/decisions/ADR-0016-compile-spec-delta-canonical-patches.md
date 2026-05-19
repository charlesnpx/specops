---
id: ADR-0016
title: Compile spec delta canonical patches
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-07
---

# ADR-0016: Compile spec delta canonical patches

## Status

Accepted.

## Context

After `synthesize --from`, the run stores an authored `outputs/spec_delta.json` and proposed decisions. Once those decisions are accepted, `compile --accepted-only` should turn the accepted semantic delta into a file-level patch plan.

The existing compiler only created reviewed provenance under `docs/research/refinery/`. That preserved lineage, but it ignored canonical work described by accepted decisions and `affected_docs`, such as creating `docs/CANON.md` or `docs/versions/v0_scope.md`.

## Decision

`specops compile` will load `outputs/spec_delta.json` and include accepted canonical doc work in the patch plan.

The compiler must:

- preserve exact authored `spec_delta.patch_items` when they are present
- otherwise generate deterministic create items for accepted canonical `affected_docs`
- continue to include the reviewed provenance item under `docs/research/refinery/`
- keep `patch_plan` as human-readable notes, not exact file content
- allow `compile --accepted-only` to be rerun from `compiled` or `planned` before apply so incomplete patch plans can be regenerated

The `SpecDelta` schema must make string-list fields explicit and add optional `patch_items` for exact file-level content.

## Consequences

Positive:

- Accepted canonical doc work is no longer lost during compile.
- Agents can author exact canonical content in `patch_items` when deterministic generated content is not sufficient.
- Provenance promotion remains separate from canonical spec mutation.

Negative / tradeoffs:

- Generated affected-doc content is only a deterministic projection of the authored delta.
- High-fidelity canonical docs should be supplied as authored `patch_items`.

## Affected docs

- docs/interfaces/cli_commands.md
- docs/interfaces/spec_delta_schema.md
- docs/interfaces/git_patch_contract.md
- docs/versions/release_decision_log.md
