---
name: specops
description: Work with SpecOps specification repositories and the local specops CLI.
specops_version: 0.1.1
payload_version: 0.1.1
compatible_cli_range: ">=0.1.1 <0.2.0"
---

# SpecOps

Use this skill when the task involves a SpecOps specification repository, specification compilation, run artifacts, decision gates, or the mise-en-place delegated installer contract.

Operational rules:

1. Read AGENTS.md and .specops/process.md before canonical edits.
2. Keep draft run artifacts in .specops/runs/.
3. Keep reviewed provenance in docs/research/refinery/.
4. Treat ADRs as append-only decision history.
5. Use the specops CLI for run state, compile, apply, audit, and eval flows where available.

## Operator loop

Use the CLI as the durable state machine and the skill as the conversational layer.

1. Inspect state with specops context <run-id> before semantic work.
2. Run mechanical steps when there is only one safe path.
3. Pause at semantic gates: refine, harden, synthesize, decisions, and apply.
4. Ask at most three open-ended questions in one batch, then ask whether to continue, pause, or change direction.
5. After an answer batch, record the guidance with specops note <run-id> --stage <stage> --text <file-or-inline>.
6. Refresh with specops context <run-id> before proceeding.
7. Use --from <file> with refine, harden, or synthesize when the semantic artifact was authored by the agent or operator.
