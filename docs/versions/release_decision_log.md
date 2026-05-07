---
id: release-decision-log
title: Release Decision Log
doc_type: release_log
status: accepted
normative: false
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Release Decision Log

## 2026-05-06

Accepted v0 direction:

```text
Go CLI + mise-en-place delegated installer + target repo .specops scaffold.
```

Deferred:

```text
MCP server, UI, direct model API default, fully autonomous loop.
```

Accepted v0.1.1 patch direction:

```text
CLI + skill hybrid operator loop: CLI owns durable context, notes, authored artifacts, and transitions; skills own bounded conversational UX.
```

Accepted v0.1.2 patch direction:

```text
Semantic production gates are enforced by recorded stage notes; "continue" is scoped to one current gate, and agents must refresh context after every semantic command.
```

Accepted v0.1.3-dev local patch direction:

```text
Semantic production commands require authored artifacts via --from; stage notes are provenance, not source material for CLI-generated semantic output.
```

Accepted v0.1.3-dev compile patch direction:

```text
Compile reads the accepted spec delta and includes canonical affected docs, preserving authored patch_items when supplied and otherwise generating deterministic doc patches from structured delta fields.
```

Accepted v0.1.3-dev patch plan health direction:

```text
Patch plan health separates stale input hashes from incomplete accepted-delta coverage; direct apply blocks either unsafe state unless explicitly overridden.
```
