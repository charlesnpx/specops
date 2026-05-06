---
id: data-lifecycle
title: Data Lifecycle
doc_type: architecture
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Data Lifecycle

## Source material lifecycle

```text
raw input file/conversation/trace
  -> ingested input artifact
  -> normalized source manifest
  -> intake summary
  -> ambiguity register
  -> decision queue
  -> spec delta
  -> patch plan
  -> canonical docs after apply
```

## Artifact classes

### Ephemeral

```text
.specops/runs/<run-id>/tmp/
.specops/cache/
```

Can be deleted.

### Draft run artifacts

```text
.specops/runs/<run-id>/outputs/refinery_note.md
.specops/runs/<run-id>/outputs/spec_delta.json
.specops/runs/<run-id>/outputs/decision_queue.yaml
.specops/runs/<run-id>/outputs/patch_plan.json
```

Not canonical until applied or promoted.

### Reviewed provenance

```text
docs/research/refinery/YYYY-MM-DD-topic.md
```

Human-readable trace of how canonical decisions emerged.

### Canonical spec

```text
docs/CANON.md
docs/decisions/*.md
docs/subsystems/*.md
docs/cross_cutting/*.md
docs/interfaces/*.md
docs/versions/*.md
docs/implementation/*.md
```

Accepted truth.

## State persistence

User-level CLI state should use XDG-style defaults:

```text
$XDG_STATE_HOME/specops/state.json
$XDG_STATE_HOME/specops/history.jsonl
$XDG_STATE_HOME/specops/state.lock
$XDG_CACHE_HOME/specops/
$XDG_CONFIG_HOME/specops/config.yaml
```

Target repo state should live under:

```text
.specops.yaml
.specops.lock
.specops/runs/
```
