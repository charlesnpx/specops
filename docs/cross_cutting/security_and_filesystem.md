---
id: security-and-filesystem
title: Security and Filesystem
doc_type: cross_cutting_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Security and Filesystem

## Filesystem rules

- Never write outside declared target paths.
- Refuse relative `--install-root`.
- In installer JSON, file paths must be absolute.
- Use atomic writes for state and lock-protected shared state.
- Do not store secrets in target spec repos.

## Agent safety

- Agent tasks should default to read-only or workspace-write behavior.
- Canonical mutation requires explicit `apply`.
- External network access is not required for local scaffold/audit operations.
