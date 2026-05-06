---
id: phase-9-self-hosting
title: Phase 9: Self-hosting and Spec Evolution
doc_type: implementation_phase
status: proposed
normative: true
version_scope: post_v0_candidate
last_reviewed: 2026-05-06
---


# Phase 9: Self-hosting and Spec Evolution

## Goal

Use `specops` to evolve its own specification repository.

## Required work

- Import this spec repo as a target repo.
- Run `specops audit` against it.
- Create a new run from a future design conversation.
- Generate a spec delta.
- Accept decisions.
- Apply canonical patches.
- Compare resulting diffs to manual expectations.

## Acceptance criteria

SpecOps can update its own spec docs through the same run/decision/apply/audit loop it provides to downstream users.
