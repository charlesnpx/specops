---
id: ADR-0017
title: Separate stale and incomplete patch plans
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-07
---

# ADR-0017: Separate stale and incomplete patch plans

## Status

Accepted.

## Context

ADR-0016 made compile include accepted canonical work from the spec delta. The apply gate still needs deterministic safety checks so agents do not have to infer whether a plan is safe to apply.

There are two related but distinct unsafe states:

- A plan is stale when compile inputs changed after plan creation.
- A plan is incomplete when it does not cover the current accepted delta, even if the inputs did not change.

Both should block direct apply unless the operator explicitly overrides the block.

## Decision

Patch plans will expose health metadata with separate `stale` and `incomplete` flags and reason lists.

Stale checks compare the plan's recorded compile inputs against current run inputs, including:

- spec delta hash
- accepted decision hash
- provenance input hash
- compiler contract version

Incomplete checks compare patch item paths against the current accepted delta, including accepted canonical `affected_docs` and authored `patch_items`.

`specops apply` refuses direct mutation when either flag is true unless `--allow-unsafe-plan` is passed. `--dry-run` remains allowed for inspection.

## Consequences

Positive:

- Apply safety no longer depends on agent judgment.
- Old plans produced before compile canonical-doc support are reported as stale.
- Plans missing accepted canonical docs are reported as incomplete even if their input hashes match.

Negative / tradeoffs:

- Patch plan JSON carries additional metadata.
- Operators need an explicit override when applying known-unsafe plans.

## Affected docs

- docs/interfaces/cli_commands.md
- docs/interfaces/git_patch_contract.md
- docs/versions/release_decision_log.md
