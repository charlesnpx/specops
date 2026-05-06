---
id: specops-git
title: SpecOps Git Adapter
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_git

## Purpose

Interacts with the target repo's Git working tree.

## Responsibilities

- Detect dirty working tree.
- Generate diffs for patch plans.
- Optionally create branches.
- Optionally commit applied patches.
- Provide PR summary text.

## Rules

- Git operations are opt-in.
- `apply` may write files without committing.
- Commit/branch operations require explicit flags.

## Acceptance criteria

- Dirty tree warnings are clear.
- Applied files can be reviewed with normal `git diff`.
