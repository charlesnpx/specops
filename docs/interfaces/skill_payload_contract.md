---
id: skill-payload-contract
title: Skill Payload Contract
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Skill Payload Contract

SpecOps installs host-specific skill payloads that teach Claude Code and Codex how to use the `specops` CLI.

## Claude payload

Destination:

```text
~/.claude/skills/specops/SKILL.md
```

Required content:

```yaml
---
name: specops
description: Use SpecOps to turn ideation, traces, and decisions into specification repositories.
---
```

The body should explain:

- when to call the skill
- how to invoke `specops` commands
- how to respect the decision gate
- how to avoid canonical mutation before `apply`
- how to summarize run artifacts to the user

## Codex payload

Destination is configurable; v0 should support the expected Codex user skill directory and keep the exact path behind target resolution.

Required behavior:

- Instruct Codex to read `AGENTS.md` and `.specops/process.md` in target repos.
- Prefer CLI commands over ad hoc editing.
- Write draft artifacts before canonical patches.
- Run audits after applying patches.

## Operator loop requirements

Claude and Codex payloads must define the skill-led interactive operator loop:

- The CLI owns durable run state, compiled context, prompt artifacts, authored semantic artifacts, and safe command transitions.
- The skill owns conversational UX.
- Agents inspect `specops context <run-id>` before semantic work.
- Agents run mechanical steps when `specops next <run-id> --json` reports one safe path.
- Agents MUST stop at semantic gates: `refine`, `harden`, `synthesize`, `decisions`, and `apply`.
- Agents treat `suggested_question_prompts` as generic scaffolding. The CLI remains non-AI and does not generate content-aware questions.
- Agents read the run artifacts and ask at most three content-aware open-ended questions in one batch, plus one control question asking whether to continue, pause, or change direction.
- A user answer of `continue` applies only to the current semantic gate and is never permission to cross multiple semantic gates.
- After an answer batch, agents record guidance with `specops note <run-id> --stage <stage> --text <file-or-inline>`.
- For `refine`, `harden`, and `synthesize`, agents author the semantic artifact themselves as a draft run artifact before invoking the CLI transition.
- When authoring `spec_delta.json` for `synthesize`, agents put full canonical document bodies in `patch_items[].content` whenever deterministic generated docs would be too thin. Agents use `patch_plan` only for human-readable compile notes and `affected_docs` only for coverage.
- Agents run at most the one semantic command for the current stage with `--from <file>`, then refresh `specops context <run-id>` before running anything else.
- If refreshed context reports another semantic gate, agents MUST stop and ask that gate's questions before running any command.
- At the apply gate, agents inspect patch plan health and treat both `stale` and `incomplete` as unsafe states. They rerun `specops compile <run-id> --accepted-only` when either flag is true before asking the operator whether to apply.
- At the apply gate, when the patch plan is mechanically healthy but semantically too thin, agents record an apply-stage note, author a replacement `spec_delta.json` with `patch_items`, run `specops supersede-synthesis <run-id> --from <spec_delta.json>`, refresh context, and recompile. They must not pass `--reopen-decisions` unless the operator explicitly wants settled decisions reopened.

The payload must explain that stage notes record guidance and provenance; they are not semantic source material that the CLI transforms. It must also explain that `specops refine`, `specops harden`, and `specops synthesize` require both a matching recorded stage note and an authored artifact passed with `--from`.

The payload source used by the installer and release assets must remain equivalent; installer tests should fail if payloads drift.

## Payload versioning

Payloads should include:

```text
specops_version
payload_version
compatible_cli_range
```

The installer report should hash payload files after staging.
