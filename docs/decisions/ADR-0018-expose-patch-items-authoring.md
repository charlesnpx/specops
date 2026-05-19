---
id: ADR-0018
title: Expose patch_items authoring
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-07
---

# ADR-0018: Expose patch_items authoring

## Status

Accepted.

## Context

ADR-0016 added `spec_delta.patch_items` as the high-fidelity path for exact canonical document content. The compiler preserves those items when present and otherwise generates skeletal deterministic documents from accepted `affected_docs`.

The CLI behavior was correct, but the scaffolded `spec_delta.yaml` template and installed skill payloads did not make the distinction visible enough. Agents could reasonably put substantive document intent in `patch_plan`, then receive thin compiled docs because `patch_plan` is only human-readable notes.

## Decision

SpecOps will expose `patch_items` authoring in the operational affordances agents see during synthesis.

The scaffolded spec delta template must include:

- all common structured delta list fields
- a warning that `patch_plan` is human-readable notes only
- an explicit `patch_items` placeholder and commented example with `content: |`

Claude and Codex skill payloads must state that, during `synthesize`, full canonical document bodies belong in `patch_items[].content` whenever deterministic generated docs would be too thin. They must also state that `affected_docs` is coverage only.

The JSON schema should describe `patch_plan`, `patch_items`, and patch item `content` so editor/schema surfaces communicate the same contract.

## Consequences

Positive:

- Agents have the exact canonical-content path available before compile.
- Synthesis artifacts are less likely to produce structurally correct but thin canonical docs.
- The CLI remains non-AI and does not invent document bodies from notes.

Negative / tradeoffs:

- The default spec delta template is longer.
- Agents still have to author the real canonical text; the compiler only preserves or deterministically projects what it is given.

## Affected docs

- .specops/process.md
- .specops/templates/spec_delta.yaml
- docs/interfaces/cli_commands.md
- docs/interfaces/skill_payload_contract.md
- docs/interfaces/spec_delta_schema.md
- docs/versions/release_decision_log.md
