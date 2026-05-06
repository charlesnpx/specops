---
id: agent-backend-contract
title: Agent Backend Contract
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Agent Backend Contract

## PromptPacket

```json
{
  "schema": 1,
  "mode": "refine",
  "instructions": "...",
  "inputs": [],
  "repo_context": [],
  "output_schema": {},
  "expected_outputs": []
}
```

## AgentResult

```json
{
  "schema": 1,
  "backend": "codex_cli",
  "status": "success",
  "stdout_ref": "...",
  "stderr_ref": "...",
  "outputs": [],
  "warnings": []
}
```

Backends must not be trusted until expected output artifacts validate.
