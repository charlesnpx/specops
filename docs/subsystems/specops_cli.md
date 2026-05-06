---
id: specops-cli
title: SpecOps CLI
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_cli

## Purpose

Owns the command-line entrypoint, command graph, global flags, output modes, and command dispatch.

## Owned responsibilities

- Root command and subcommands.
- Global flags: `--repo`, `--config`, `--json`, `--quiet`, `--verbose`, `--backend`.
- Error formatting.
- Exit codes.
- Shell completions and generated command docs.
- JSON-only stdout discipline.

## Does not own

- Installer staging semantics.
- Run state persistence details.
- Artifact schemas.
- Agent-specific behavior.

## Core commands

```text
init, doctor, upgrade, config
run new/list/show/status
next
ingest-file, ingest-chat, ingest-relay
mine-trace, mine-gold, fixture-build
intake, refine, harden, synthesize, deepen
decisions, accept, reject, defer, amend
compile, plan, apply, audit
reproduce, eval, diff, score
install-skill plan/install/uninstall
```

## Acceptance criteria

- Every command supports `--help`.
- Machine-readable commands support `--json`.
- When `--json` is set, stdout is parseable JSON and logs go to stderr.
- Exit code `0` means success, `1` operational error, `2` usage/validation error, `3` contract violation.
