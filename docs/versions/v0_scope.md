---
id: v0-scope
title: v0 Scope
doc_type: version_scope
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# v0 Scope

## v0 required

- Go CLI named `specops`.
- `install-skill.sh` delegated installer contract.
- Claude and Codex skill payload installation.
- Tool binary installation under `tools` target.
- `specops init` scaffold installation.
- `.specops.yaml` and `.specops.lock`.
- Run state management.
- Input ingestion from Markdown/transcript/relay files.
- Static/manual backend for intake/refine/synthesize outputs.
- Decision queue and human gate.
- Compile/plan/apply canonical patch flow.
- Pre-apply synthesis supersession when plan review finds incomplete semantic content.
- Audit checks for repo taxonomy, front matter, ADR links, and v0 acceptance criteria.
- Reproduction/eval framework with deterministic structural checks.
- Release artifacts for macOS/Linux arm64/x86_64.

## v0 optional

- Direct Codex CLI backend.
- Direct Claude Code backend.
- `convo-relay` backend.
- Semantic LLM eval graders.
- Git branch/commit automation.

## v0 non-goals

- Hosted UI.
- MCP server.
- Long-running cloud agent.
- Direct model API integration as required default.
- Multi-user approval workflows.
