---
id: implementation-repository-layout
title: Implementation Repository Layout
doc_type: architecture
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Implementation Repository Layout

The implementation repo for SpecOps should be separate from this specification repo, but can initially vendor this spec under `docs/spec/` if useful.

Recommended Go project layout:

```text
specops/
  go.mod
  go.sum
  README.md
  LICENSE
  install-skill.sh
  .goreleaser.yaml
  .github/workflows/ci.yml
  .github/workflows/release.yml

  cmd/specops/main.go

  internal/cli/
    root.go
    setup_commands.go
    run_commands.go
    production_commands.go
    decision_commands.go
    compile_commands.go
    eval_commands.go
    installer_commands.go

  internal/install/
    contract.go
    planner.go
    stager.go
    targets.go
    report.go
    hash.go
    paths.go

  internal/scaffold/
    init.go
    upgrade.go
    lockfile.go
    assets.go
    merge.go

  internal/runstate/
    model.go
    store.go
    transitions.go
    next.go

  internal/input/
    ingest_file.go
    ingest_chat.go
    ingest_relay.go
    slicer.go

  internal/artifacts/
    spec_delta.go
    decision.go
    patch_plan.go
    eval_report.go
    validation.go

  internal/agents/
    backend.go
    manual.go
    codex.go
    claude.go
    relay.go

  internal/compiler/
    adr.go
    canon.go
    patch.go
    apply.go

  internal/audit/
    frontmatter.go
    links.go
    taxonomy.go
    adrs.go
    version_scope.go

  internal/eval/
    structure.go
    decision_coverage.go
    lineage.go
    report.go

  internal/git/
    status.go
    diff.go
    commit.go

  internal/state/
    xdg.go
    lock.go
    history.go

  internal/output/
    json.go
    human.go
    errors.go

  assets/
    embed.go
    scaffold/
      .specops/
      AGENTS.md
      CLAUDE.md
      schemas/
      templates/
      evals/
    skills/
      claude/specops/SKILL.md
      codex/specops/SKILL.md
```

## Package rules

- `cmd/specops` should contain only `main` wiring.
- `internal/install` must not import agent or compiler packages.
- `internal/scaffold` may use embedded assets but should not call model backends.
- `internal/compiler` should operate on validated artifacts, not raw agent output.
- `internal/output` owns stdout/stderr discipline.

## Build tags / platform notes

- File-locking may need Unix-specific and fallback implementations.
- Windows is not v0 release target unless explicitly added later.
- macOS and Linux arm64/x86_64 are v0 release targets.
