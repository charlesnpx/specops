---
id: system-overview
title: System Overview
doc_type: architecture
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# System Overview

## Context

SpecOps is both:

1. A standalone CLI/toolchain repository.
2. A delegated skill/tool installable through `mise-en-place`.
3. A process scaffold installed into target spec repositories.

## Major runtime pieces

```text
+-----------------------+
| mise-en-place         |
| skill manager         |
+-----------+-----------+
            |
            | delegated installer contract
            v
+-----------------------+       +--------------------------+
| specops repo          |       | target spec repo         |
| Go CLI + installer    |-----> | .specops control plane   |
| skills/templates      | init  | canonical docs           |
+-----------+-----------+       +-------------+------------+
            |                                 ^
            | agent backend calls             |
            v                                 |
+-----------------------+       +-------------+------------+
| Claude/Codex/relay    |-----> | run artifacts + patches  |
| reasoning/drafting    |       | human decision gate      |
+-----------------------+       +--------------------------+
```

## Responsibilities

### `mise-en-place`

- Resolves delegated repo.
- Executes installer plan/install/uninstall contract.
- Performs collision checks and state recording.
- Manages dual-target skill deployment.

### `specops` CLI

- Provides delegated installer payloads.
- Initializes target repos.
- Manages run state and artifacts.
- Prepares prompt packets for agents.
- Parses and validates typed outputs.
- Compiles accepted decisions into patch plans.
- Applies patches and runs audits/evals.

### Agent host

- Performs reasoning and drafting.
- Reads `.specops/` local process docs.
- Produces structured artifacts or repository diffs.

### Human

- Accepts/rejects/defer decisions.
- Reviews diffs.
- Owns canonical truth.
