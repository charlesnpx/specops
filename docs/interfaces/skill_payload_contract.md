---
id: skill-payload-contract
title: Skill Payload Contract
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Skill Payload Contract

SpecOps installs host-specific skill payloads that teach Claude Code and Codex how to use the `specops` CLI.

## Claude payload

Destination:

```text
~/.claude/skills/specops/SKILL.md
```

Required content:

```yaml
---
name: specops
description: Use SpecOps to turn ideation, traces, and decisions into specification repositories.
---
```

The body should explain:

- when to call the skill
- how to invoke `specops` commands
- how to respect the decision gate
- how to avoid canonical mutation before `apply`
- how to summarize run artifacts to the user

## Codex payload

Destination is configurable; v0 should support the expected Codex user skill directory and keep the exact path behind target resolution.

Required behavior:

- Instruct Codex to read `AGENTS.md` and `.specops/process.md` in target repos.
- Prefer CLI commands over ad hoc editing.
- Write draft artifacts before canonical patches.
- Run audits after applying patches.

## Payload versioning

Payloads should include:

```text
specops_version
payload_version
compatible_cli_range
```

The installer report should hash payload files after staging.
