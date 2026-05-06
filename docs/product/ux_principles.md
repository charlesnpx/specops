---
id: ux-principles
title: UX Principles
doc_type: product_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# UX Principles

## CLI UX

- Commands should be composable, like `git`, not a brittle wizard.
- Dynamic commands should print a short human summary, artifact paths, and recommended next command.
- Mutating commands must support dry-run or plan modes.
- Canonical repo mutation must be explicit through `apply`.
- JSON mode must produce JSON-only stdout.
- Human-readable logs go to stderr when machine output is requested.

## Agent UX

- The target repo should be legible to agents through `AGENTS.md`, `CLAUDE.md`, and `.specops/` command docs.
- Agent tasks should be small, typed, and stateful.
- The agent should not need the entire repo dumped into context; it should use indexes and targeted files.

## Human review UX

- Every run should explain what it found, what is uncertain, what it recommends, and which files would change.
- Ambiguities should appear as option tables.
- Decisions should be accepted/rejected/deferred explicitly.
- Diffs should be readable and reviewable before merge.
