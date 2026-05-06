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
specops next <run-id>
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
specops refine <run-id>
specops harden <run-id> [--backend convo-relay]
specops synthesize <run-id>
specops deepen <run-id> --target <concept-or-doc>
```

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
