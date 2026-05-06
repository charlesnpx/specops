---
id: implementation-language-analysis
title: Implementation Language Analysis
doc_type: research_note
status: accepted
normative: false
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Implementation Language Analysis

## Go

Pros:

- Self-contained binaries.
- Strong fit for local CLI and installer workflows.
- Cross-compilation and GoReleaser support.
- Aligns with `mise-en-place` being Go-based.
- Good filesystem, hashing, JSON, subprocess, and release tooling.

Cons:

- Less ergonomic than Python for rapid data/LLM experimentation.
- Requires explicit modeling of state and schemas.

## Python

Pros:

- Fastest for LLM, eval, and Markdown processing experiments.
- Rich data tooling.

Cons:

- Packaging standalone binaries is harder.
- Users need Python/pipx unless separately bundled.
- Weaker fit for the `mise-en-place` self-contained binary pattern.

## Rust

Pros:

- Excellent single-binary CLI.
- Strong type and filesystem safety.

Cons:

- Slower iteration for this research-heavy phase.

## Recommendation

Use **Go** for v0 CLI and installer. Keep eval/agent prompts data-driven so future Python helpers can be added without changing the core contract.
