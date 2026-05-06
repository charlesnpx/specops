---
id: specops-scaffold
title: SpecOps Scaffold Manager
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_scaffold

## Purpose

Installs and upgrades `.specops/` control-plane files inside target repositories.

## Owned responsibilities

- Create `.specops.yaml`.
- Create `.specops.lock`.
- Write `AGENTS.md` and `CLAUDE.md`.
- Vendor or link templates/schemas/evals.
- Maintain scaffold version and migration records.
- Update `.gitignore` for `.specops/runs/` and `.specops/cache/`.

## Modes

```text
minimal
vendor
linked
agent=claude|codex|both
```

## Acceptance criteria

- Re-running `specops init` is idempotent.
- Existing divergent files are not overwritten without `--force`, `--backup`, or explicit confirmation.
- Scaffold version is recorded in `.specops.lock`.
