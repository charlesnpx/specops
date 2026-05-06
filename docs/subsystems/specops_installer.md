---
id: specops-installer
title: SpecOps Installer
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_installer

## Purpose

Implements the `mise-en-place` delegated repo installer contract.

## Owned responsibilities

- Parse `--plan`, `--install`, `--uninstall`.
- Parse `--target claude|codex|tools|all`.
- Parse `--json` and enforce JSON-only stdout.
- Parse `--install-root <absolute-dir>`.
- Stage Claude payloads, Codex payloads, and tool binary.
- Compute SHA-256 after install.
- Return required JSON shape.

## Does not own

- Collision prompting in real user destinations.
- Final copy from staging root to home.
- mise-en-place state records.

Those are owned by `mise-en-place` after staging.

## Install targets

```text
claude -> ~/.claude/skills/specops/SKILL.md and supporting files
codex  -> ~/.codex/skills/specops/SKILL.md and supporting files
tools  -> ~/.local/bin/specops
all    -> claude + codex + tools
```

## Acceptance criteria

- `--plan --json` emits required keys without writing files.
- `--install --json --install-root /tmp/stage` writes files only under `/tmp/stage`.
- Paths in JSON are absolute.
- `sha256` is present after install.
- Private/human logs go to stderr.
