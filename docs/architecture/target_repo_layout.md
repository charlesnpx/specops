---
id: target-repo-layout
title: Target Repository Layout
doc_type: architecture
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Target Repository Layout

`specops init` installs a local control plane and ensures a canonical docs layout.

## Minimal target repo

```text
my-product-spec/
  .specops.yaml
  .specops.lock
  AGENTS.md
  CLAUDE.md
  .gitignore

  .specops/
    README.md
    process.md
    status_model.md
    commands/
    templates/
    schemas/
    evals/
    checklists/
    runs/

  docs/
    CANON.md
    GLOSSARY.md
    OPEN_QUESTIONS.md
    decisions/
    architecture/
    subsystems/
    cross_cutting/
    interfaces/
    implementation/
    versions/
    research/
    _generated/
```

## Install modes

### `--minimal`

Installs config, lock, agent entry docs, and command shims.

### `--vendor`

Copies templates, schemas, evals, and command docs into target repo.

### `--linked`

Stores a lockfile reference to the installed CLI/scaffold version and does not vendor all templates.

## Recommendation

v0 should default to `--vendor` for transparency. Later releases may default to `--minimal` once the process stabilizes.
