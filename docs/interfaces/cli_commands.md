---
id: cli-commands
title: CLI Commands
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# CLI Commands

## Setup

```sh
specops init [path] [--minimal|--vendor|--linked] [--agent claude|codex|both] [--force|--backup]
specops doctor
specops upgrade
specops config get|set|list
specops completion bash|zsh|fish|powershell
```

## Run management

```sh
specops run new --name <name>
specops run list
specops run show <run-id>
specops run status <run-id>
specops context <run-id> [--json]
specops note <run-id> --stage <stage> --text <file-or-inline>
specops next <run-id> [--json]
```

## Input

```sh
specops ingest-file <path> --run <run-id>
specops ingest-chat <path> --slice --run <run-id>
specops ingest-relay <path> --run <run-id>
specops mine-trace --input <path> [--gold <repo>]
specops fixture-build --from <path> --out <dir>
```

## Spec production

```sh
specops intake <run-id>
specops refine <run-id> --from <file>
specops harden <run-id> --from <file> [--backend convo-relay]
specops synthesize <run-id> --from <spec_delta.json>
specops supersede-synthesis <run-id> --from <spec_delta.json> [--reopen-decisions]
specops deepen <run-id> --target <concept-or-doc>
```

## Run context and operator notes

`specops context <run-id>` is non-mutating. It returns the compiled run context so far:

- run id, name, and status
- source summary when present
- operator guidance for the next gate
- current artifacts
- decisions
- patch plan when present
- next gate metadata

`specops note <run-id> --stage <stage> --text <file-or-inline>` records operator or agent guidance under `.specops/runs/<run-id>/prompts/` and adds a `prompt` artifact to run state. The prompt artifact records `stage` metadata so semantic gates can verify that guidance was captured for the current stage. It must not advance run status.

For compatibility with v0.1.1 run state, a prompt artifact without `stage` metadata also satisfies the matching stage when its path has the legacy shape `prompts/<timestamp>-<stage>.md`.

## `next --json`

`specops next <run-id> --json` keeps the existing `command` and `reason` fields under `next` and adds gate metadata:

```json
{
  "run_id": "run-...",
  "status": "intake_complete",
  "next": {
    "command": "specops refine run-... --from <file>",
    "reason": "intake artifact is ready to refine",
    "stage": "refine",
    "gate_kind": "semantic",
    "context_command": "specops context run-...",
    "note_command": "specops note run-... --stage refine --text <file-or-inline>",
    "suggested_question_prompts": [],
    "human_input_recommended": true
  }
}
```

`gate_kind` is `mechanical` when the next step has one safe deterministic path, and `semantic` when the operator loop should pause for bounded human guidance.

For semantic gates, `note_command` gives the command shape for recording the current gate's operator or agent guidance. Human `context` and `next` output must state that semantic production commands require a stage note before execution.

## Authored semantic artifacts

`refine --from`, `harden --from`, and `synthesize --from` let an agent or human produce the semantic artifact outside the CLI. The CLI copies the supplied artifact into the run output without rewriting the content, then applies the legal state transition. `synthesize --from` must receive a parseable `spec_delta.json`; its decisions are loaded into run state.

For `synthesize --from`, `patch_plan` is notes only and `affected_docs` is coverage only. Full canonical document bodies belong in `patch_items[].content` when deterministic generated docs would be too thin.

`specops refine`, `specops harden`, and `specops synthesize` are semantic production commands. Before running, each command must find a recorded stage note for its own stage:

```sh
specops note <run-id> --stage refine --text <file-or-inline>
specops note <run-id> --stage harden --text <file-or-inline>
specops note <run-id> --stage synthesize --text <file-or-inline>
```

If the matching stage note is missing, the command must fail without writing an output artifact or advancing run status. The error must show both:

```sh
specops context <run-id>
specops note <run-id> --stage <stage> --text <file-or-inline>
```

After the matching stage note exists, these semantic production commands still require `--from`. The note records guidance and provenance; it is not semantic source material that the CLI transforms. If `--from` is omitted, the command must fail without writing an output artifact or advancing run status and explain that the CLI is not AI-enabled and requires an authored artifact.

## Pre-apply synthesis supersession

`specops supersede-synthesis <run-id> --from <spec_delta.json>` lets an operator replace a too-thin synthesized delta after `plan` review and before `apply`.

The command requires:

- run status `planned`
- an `apply` stage note
- a parseable replacement `spec_delta.json`

By default, existing settled decisions remain settled. Replacement deltas may omit decisions or repeat existing decision IDs with the same substance, but they must not introduce new decision IDs or change settled decision substance. On success, the run returns to `decisions_accepted` so `compile --accepted-only` can regenerate the patch plan.

With `--reopen-decisions`, the command may merge new or changed decisions from the replacement delta and returns the run to `awaiting_decisions`.

Supersession is append-only. The previous current `outputs/spec_delta.json` and `patches/patch_plan.json` are copied under `outputs/superseded/` and `patches/superseded/`, old current artifact refs are reclassified as superseded refs, the current patch plan is removed, and the replacement delta becomes the current `outputs/spec_delta.json`.

## Compile behavior

`specops compile <run-id> --accepted-only` loads the accepted decisions and the run's `outputs/spec_delta.json`. It may be rerun from `decisions_accepted`, `compiled`, or `planned` to regenerate an unsafe patch plan before apply. The patch plan must include reviewed provenance and must also represent accepted canonical doc work from the spec delta:

- exact authored `patch_items` are preserved when present
- otherwise accepted canonical `affected_docs` receive deterministic create items generated from the structured delta fields and accepted decisions
- `docs/research/refinery/` provenance remains a separate reviewed provenance item

Patch plans expose separate health flags:

- `stale`: compile inputs changed after plan creation, such as spec delta, accepted decisions, provenance inputs, or compiler contract version
- `incomplete`: the plan does not cover the current accepted delta, even when the inputs have not changed

`specops apply` refuses a stale or incomplete patch plan unless `--allow-unsafe-plan` is explicitly passed. `--dry-run` remains available for inspecting unsafe plans without mutating files.

## Decisions

```sh
specops decisions <run-id>
specops accept <run-id> <decision-id>|--all-recommended
specops reject <run-id> <decision-id>
specops defer <run-id> <decision-id> --reason <text>
specops amend <run-id> <decision-id> --text <file-or-inline>
```

## Compile and mutate

```sh
specops compile <run-id> --accepted-only
specops plan <run-id>
specops apply <run-id> [--interactive] [--dry-run] [--commit] [--allow-unsafe-plan]
specops audit
```

## Reproduction/eval

```sh
specops reproduce --fixture <dir> --out <dir>
specops eval --gold <repo> --candidate <repo>
specops diff --gold <repo> --candidate <repo>
specops score <eval-report>
```

## Installer compatibility

```sh
specops install-skill --plan --target all --json
specops install-skill --install --target all --json --install-root /tmp/stage
specops install-skill --uninstall --target all --json
```

`install-skill.sh` may delegate to these commands.
When run from a source checkout, the wrapper injects the exact resolved git tag
version into `go run`. It must not call a `specops` executable from `PATH`,
because delegated installs must report the version and payloads from the
checkout that `mise-en-place` resolved.
