---
id: specops-release
title: SpecOps Release Subsystem
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_release

## Purpose

Owns release packaging, archive contents, checksums, and registry-readiness.

## Responsibilities

- GoReleaser config.
- GitHub Actions release workflow.
- Cross-platform binary naming.
- Archive contents.
- Checksums.
- Version injection into `specops version` and installer report.
- Smoke test archives by running installer contract in a temp root.

## Acceptance criteria

- Release archives contain the binary, installer wrapper, skill payloads, and scaffold assets.
- `specops version --json` reports version, commit, date, and scaffold version.
- `./install-skill.sh --plan --target all --json` works from unpacked release archive.
