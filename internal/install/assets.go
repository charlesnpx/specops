package install

const claudeSkill = `---
name: specops
description: Produce and maintain SpecOps specification repositories through the local specops CLI.
---

# SpecOps

Use this skill when working in a repository governed by SpecOps, or when asked to turn rough product/process material into canonical specification docs.

Follow the local repository instructions first:

1. Read AGENTS.md or CLAUDE.md.
2. Read .specops/process.md when present.
3. Put run artifacts under .specops/runs/.
4. Promote reviewed provenance to docs/research/refinery/.
5. Do not mutate canonical docs before accepted decisions are available.

Prefer the specops CLI for scaffold, run state, decision, compile, apply, audit, and eval operations.
`

const codexSkill = `---
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
`
