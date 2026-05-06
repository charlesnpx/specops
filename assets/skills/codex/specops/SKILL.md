---
name: specops
description: Work with SpecOps specification repositories and the local specops CLI.
specops_version: 0.1.2
payload_version: 0.1.2
compatible_cli_range: ">=0.1.2 <0.2.0"
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
3. You MUST stop at semantic gates: refine, harden, synthesize, decisions, and apply.
4. At a semantic gate, read the run artifacts and ask at most three content-aware open-ended questions, then ask whether to continue, pause, or change direction.
5. Treat the CLI's suggested_question_prompts as generic scaffolding. The CLI is not AI-enabled and does not generate the content-aware questions for you.
6. A user answer of "continue" applies only to exactly the current semantic gate. It is not permission to cross any later semantic gate.
7. After an answer batch, record the guidance with specops note <run-id> --stage <stage> --text <file-or-inline>.
8. Run at most the one semantic command for the current stage. Use --from <file> with refine, harden, or synthesize when the semantic artifact was authored by the agent or operator.
9. After every semantic command, refresh with specops context <run-id> before running anything else.
10. If refreshed context reports another semantic gate, you MUST stop and ask that gate's questions before running any command.
11. Never treat "continue" as permission to cross multiple semantic gates.
