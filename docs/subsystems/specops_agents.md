---
id: specops-agents
title: SpecOps Agent Backend Adapter
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_agents

## Purpose

Provides a common interface for execution backends.

## Backend types

```text
manual
codex_cli
claude_code
openai_api
anthropic_api
convo_relay
```

## Prompt packet

A backend receives:

```text
instructions
repo excerpts
input summaries
schema
expected output paths
mode
```

## Acceptance criteria

- Backends are replaceable.
- Backend output is captured to run traces.
- Backend failure does not corrupt run state.
- Manual backend allows a human/agent to paste output into expected files.
