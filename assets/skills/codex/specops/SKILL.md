---
name: specops
description: Work with SpecOps specification repositories and the local specops CLI.
specops_version: 0.1.3-dev
payload_version: 0.1.3-dev
compatible_cli_range: ">=0.1.3-dev <0.2.0"
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
8. For refine, harden, or synthesize, author the semantic artifact yourself as a draft run artifact, then run at most the one semantic command for the current stage with --from <file>. Never run these commands without --from.
9. When authoring a synthesized spec_delta.json, put full canonical document bodies in patch_items[].content whenever deterministic generated docs would be too thin. Use patch_plan only for human-readable compile notes, and use affected_docs only for coverage.
10. After every semantic command, refresh with specops context <run-id> before running anything else.
11. If refreshed context reports another semantic gate, you MUST stop and ask that gate's questions before running any command.
12. Never treat "continue" as permission to cross multiple semantic gates.
13. A stage note records guidance and provenance; it is not semantic source material that the CLI will transform.
14. At the apply gate, inspect the patch plan health. Treat stale (compile inputs changed) and incomplete (plan does not cover the accepted delta) as separate unsafe states. If either is true, rerun specops compile <run-id> --accepted-only before asking the operator whether to apply; compile is intentionally rerunnable from planned in 0.1.3-dev.
15. If the apply-gate patch plan is mechanically healthy but semantically too thin, record an apply-stage note, author a replacement spec_delta.json with patch_items, run specops supersede-synthesis <run-id> --from <spec_delta.json>, refresh context, and recompile. Do not pass --reopen-decisions unless the operator explicitly wants settled decisions reopened.
