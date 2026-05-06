---
name: specops
description: Work with SpecOps specification repositories and the local specops CLI.
---

# SpecOps

Use this skill when the task involves a SpecOps specification repository, specification compilation, run artifacts, decision gates, or the mise-en-place delegated installer contract.

Operational rules:

1. Read AGENTS.md and .specops/process.md before canonical edits.
2. Keep draft run artifacts in .specops/runs/.
3. Keep reviewed provenance in docs/research/refinery/.
4. Treat ADRs as append-only decision history.
5. Use the specops CLI for run state, compile, apply, audit, and eval flows where available.
