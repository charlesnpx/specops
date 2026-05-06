---
name: specops
description: Produce and maintain SpecOps specification repositories through the local specops CLI.
specops_version: 0.1.1
payload_version: 0.1.1
compatible_cli_range: ">=0.1.1 <0.2.0"
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

## Operator loop

Use the CLI as the durable state machine and the skill as the conversational layer.

1. Inspect state with specops context <run-id> before semantic work.
2. Run mechanical steps when there is only one safe path.
3. Pause at semantic gates: refine, harden, synthesize, decisions, and apply.
4. Ask at most three open-ended questions in one batch, then ask whether to continue, pause, or change direction.
5. After an answer batch, record the guidance with specops note <run-id> --stage <stage> --text <file-or-inline>.
6. Refresh with specops context <run-id> before proceeding.
7. Use --from <file> with refine, harden, or synthesize when the semantic artifact was authored by the agent or operator.
