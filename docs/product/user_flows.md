---
id: user-flows
title: User Flows
doc_type: product_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# User Flows

## Flow 1: Install through mise-en-place

```text
user runs mise-en-place install specops
  -> mise-en-place clones delegated repo
  -> runs ./install-skill.sh --plan --target all --json
  -> validates contract and collisions
  -> runs ./install-skill.sh --install --target all --json --install-root <stage>
  -> copies staged files to Claude/Codex/tool destinations
  -> records ownership in mise-en-place state
```

Expected result:

```text
~/.claude/skills/specops/SKILL.md
~/.codex/skills/specops/SKILL.md
~/.local/bin/specops
```

## Flow 2: Initialize a target spec repo

```sh
specops init ./my-product-spec --agent claude --agent codex --vendor
```

Result:

```text
.specops.yaml
.specops.lock
AGENTS.md
CLAUDE.md
.specops/process.md
.specops/commands/*.md
.specops/templates/*.md
.specops/schemas/*.json
.specops/evals/*.md
```

## Flow 3: Turn raw ideation into a run

```sh
specops run new --name permission-model
specops ingest-file notes/raw-discussion.md --run run-001
specops intake run-001
specops refine run-001
specops decisions run-001
```

Result:

```text
.specops/runs/run-001/outputs/refinery_note.md
.specops/runs/run-001/outputs/ambiguity_register.yaml
.specops/runs/run-001/outputs/decision_queue.md
.specops/runs/run-001/outputs/spec_delta.json
```

## Flow 4: Accept recommendations and patch canonical docs

```sh
specops accept run-001 --all-recommended
specops compile run-001 --accepted-only
specops plan run-001
specops apply run-001 --interactive
specops audit
```

Result:

```text
docs/research/refinery/...
docs/decisions/ADR-....md
docs/subsystems/...
docs/cross_cutting/...
docs/versions/...
docs/implementation/...
```

## Flow 5: Reproduce a gold spec

```sh
specops reproduce --fixture fixtures/context-sink --out runs/context-sink-repro-001
specops eval --gold gold/context-sink-spec --candidate runs/context-sink-repro-001
```

Result:

```text
structure eval
concept coverage eval
decision coverage eval
lineage eval
contradiction eval
implementation usefulness eval
```
