---
id: existing-tool-contract-analysis
title: Existing Tool Contract Analysis
doc_type: research_note
status: accepted
normative: false
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Existing Tool Contract Analysis

## mise-en-place

The user-provided contract defines:

- Managed skills in the `mise-en-place` repo.
- Delegated skills in external repos.
- External tools.
- Dual Claude/Codex targets.
- A delegated installer compatibility contract.
- State under `~/.local/state/mise-en-place/`.

SpecOps should be a delegated skill/tool repo, not a managed skill inside `mise-en-place` initially.

## convo-relay

`convo-relay` is an adversarial dialogue backend used to stress-test the Context Sink source specs. SpecOps should integrate it as an optional hardening backend rather than duplicating it.
