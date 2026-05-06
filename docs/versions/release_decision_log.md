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
