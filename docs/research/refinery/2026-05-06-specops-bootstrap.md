---
id: refinery-specops-bootstrap
title: SpecOps Bootstrap Refinery Note
doc_type: refinery_note
status: accepted
normative: false
version_scope: v0_required
last_reviewed: 2026-05-06
---


# SpecOps Bootstrap Refinery Note

## Raw input summary

The user asked for the specification-production strategy to be used self-referentially to produce a specification repository for implementing the SpecOps tool itself. The tool must implement the `mise-en-place` delegated installer contract so it can be installed through `mise-en-place`. The user suspected Go may be appropriate.

## Extracted decisions

- Build a standalone CLI/toolchain repository.
- Make it installable as a `mise-en-place` delegated skill/tool.
- Implement the core CLI in Go.
- Install a thin `.specops/` local operating scaffold into target repos.
- Keep produced specifications in repository proper under `docs/`.
- Use dynamic run state and typed artifacts, not static prompt dumping.

## Ambiguities

- Final product name.
- Exact Codex skill install path.
- Whether direct model API backends are v0 or post-v0.
- Whether run artifacts are committed by default.

## Patch plan

Create a full spec repository with:

- CANON
- ADRs
- subsystem specs
- mise-en-place contract docs
- CLI command docs
- implementation phases
- schemas
- `.specops/` process scaffold
