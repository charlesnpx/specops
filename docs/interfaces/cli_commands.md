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
specops refine <run-id> [--from <file>]
specops harden <run-id> [--backend convo-relay] [--from <file>]
specops synthesize <run-id> [--from <spec_delta.json>]
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

`specops note <run-id> --stage <stage> --text <file-or-inline>` records operator or agent guidance under `.specops/runs/<run-id>/prompts/` and adds a `prompt` artifact to run state. It must not advance run status.

## `next --json`

`specops next <run-id> --json` keeps the existing `command` and `reason` fields under `next` and adds gate metadata:

```json
{
  "run_id": "run-...",
  "status": "intake_complete",
  "next": {
    "command": "specops refine run-...",
    "reason": "intake artifact is ready to refine",
    "stage": "refine",
    "gate_kind": "semantic",
    "context_command": "specops context run-...",
    "suggested_question_prompts": [],
    "human_input_recommended": true
  }
}
```

`gate_kind` is `mechanical` when the next step has one safe deterministic path, and `semantic` when the operator loop should pause for bounded human guidance.

## Authored semantic artifacts

`refine --from`, `harden --from`, and `synthesize --from` let an agent or human produce the semantic artifact outside the deterministic fallback. The CLI copies the supplied artifact into the run output without rewriting the content, then applies the same legal state transition as the fallback command. `synthesize --from` must receive a parseable `spec_delta.json`; its decisions are loaded into run state.

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
specops apply <run-id> [--interactive] [--dry-run] [--commit]
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
