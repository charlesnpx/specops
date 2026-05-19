---
id: canon
title: SpecOps CLI Canon
doc_type: canon
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# SpecOps CLI Canon

## Product thesis

**SpecOps** is a CLI-backed specification-production toolchain for turning ideation, traces, conversations, adversarial relay output, and human decisions into a versioned specification repository.

It should produce the kind of repository this file lives inside: canonical docs, ADRs, subsystem specs, interface contracts, version scope, implementation phases, provenance notes, generated indexes, and evaluation reports.

## Current accepted shape

```text
specops CLI repo
  reusable engine/toolchain
  Go implementation
  delegated installer for mise-en-place
  embedded templates/schemas/prompts/evals
  optional agent backends and relay integration

target spec repository
  produced specification
  .specops/ local operating contract
  CANON.md, ADRs, subsystems, interfaces, versions, implementation docs

agent hosts
  Claude Code
  Codex CLI
  future ChatGPT/MCP/UI clients

human operator
  accepts/rejects/defer decisions
  reviews diffs
```

## Accepted architectural principles

1. **The CLI is the compiler; the target repo is the artifact.**
2. **`.specops/` is the local control plane, not the produced specification.**
3. **The repo proper is the compiled specification.**
4. **The agent proposes; the human accepts.**
5. **Run artifacts precede canonical patches.**
6. **Every accepted consequential decision should have an ADR.**
7. **The delegated installer contract from `mise-en-place` is a first-class interface.**
8. **The first implementation should be a Go CLI.**
9. **The CLI should be dynamic and stateful, not a static prompt printer.**
10. **Generated specs must remain useful without the CLI or any given agent host.**

## Why Go for v0

The tool is a local CLI that must install files, compute hashes, create directories under `$HOME`, stage installs under an arbitrary root, parse YAML/JSON, run subprocesses, package as release binaries, and interoperate with `mise-en-place`, which itself is Go-based. A Go implementation gives the simplest path to self-contained binaries and GitHub Release distribution.

## Normative external contract

The v0 implementation must expose a delegated installer command compatible with `mise-en-place`:

```sh
./install-skill.sh --plan --target claude|codex|tools|all --json
./install-skill.sh --install --target claude|codex|tools|all --json [--install-root /abs/path]
./install-skill.sh --uninstall --target claude|codex|tools|all --json
```

When `--json` is set, stdout must contain only JSON. Human logs go to stderr.

## v0 command families

```text
setup:
  init, doctor, upgrade, config

run management:
  run new, run list, run show, run status, next

input/process:
  ingest, mine-trace, intake, refine, harden, synthesize, supersede-synthesis, deepen

decision gate:
  decisions, accept, reject, defer, amend

compile/mutate:
  compile, plan, apply, audit

reproduction/eval:
  reproduce, eval, diff, score

installer:
  install-skill plan/install/uninstall or compatible wrapper
```

## v0 target repository layout

```text
.target-spec-repo/
  .specops.yaml
  .specops.lock
  AGENTS.md
  CLAUDE.md
  .specops/
    process.md
    commands/
    templates/
    schemas/
    evals/
    checklists/
    runs/        # ignored by default unless promoted

  docs/
    CANON.md
    decisions/
    architecture/
    subsystems/
    cross_cutting/
    interfaces/
    versions/
    implementation/
    research/
    _generated/
```

## Canonical workflow

```text
raw input
  -> ingest
  -> intake/refine
  -> harden if needed
  -> synthesize spec delta
  -> human decision gate
  -> compile patch plan
  -> optionally supersede synthesis before apply if plan review finds semantic content gaps
  -> apply accepted canonical patches
  -> audit
  -> eval/reproduce when relevant
```

## Non-goals for v0

- No hosted UI.
- No full custom autonomous agent loop.
- No background/asynchronous hosted execution.
- No requirement to support every agent host equally on day one.
- No automatic canonical mutation before explicit `apply`.
- No attempt to train or fine-tune a model.
- No secret storage beyond local config/state needed for CLI operation.
