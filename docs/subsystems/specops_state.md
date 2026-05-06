---
id: specops-state
title: SpecOps State Store
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_state

## Purpose

Owns local user-level state, cache, history, and locks.

## User-level locations

```text
state:  $XDG_STATE_HOME/specops/state.json
history:$XDG_STATE_HOME/specops/history.jsonl
lock:   $XDG_STATE_HOME/specops/state.lock
cache:  $XDG_CACHE_HOME/specops/
config: $XDG_CONFIG_HOME/specops/config.yaml
```

Environment override:

```text
SPECOPS_HOME
```

## Acceptance criteria

- Concurrent runs that modify shared state use an advisory lock.
- History is append-only JSONL.
- State is recoverable after partial command failure.
