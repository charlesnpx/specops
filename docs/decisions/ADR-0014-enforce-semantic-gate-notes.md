---
id: ADR-0014
title: Enforce semantic gate notes
doc_type: adr
status: accepted
normative: true
version_scope: v0_required
date: 2026-05-06
---

# ADR-0014: Enforce semantic gate notes

## Status

Accepted.

## Context

ADR-0013 established a skill-led operator loop where agents pause at semantic gates, ask bounded questions, record answers with `specops note`, refresh context, and then continue.

That contract was advisory in the installed skill payloads. A single operator answer such as `continue` could be misread by an agent as permission to cross several semantic gates in sequence, especially when `refine`, `harden`, and `synthesize` all had deterministic fallback behavior.

The CLI must remain non-AI. It should expose state, generic question scaffolding, note commands, and legal transitions. The agent or human operator remains responsible for reading artifacts and asking content-aware questions.

## Decision

SpecOps will enforce semantic production gate notes for:

- `specops refine`
- `specops harden`
- `specops synthesize`

Before either deterministic fallback or `--from` mode runs, each command must find a recorded `prompt` artifact for its own stage. If the stage note is missing, the command fails without writing output artifacts or advancing run status, and the error shows both:

```sh
specops context <run-id>
specops note <run-id> --stage <stage> --text <file-or-inline>
```

`specops note` records `stage` metadata on prompt artifacts. For compatibility with v0.1.1 run state, a prompt artifact without `stage` metadata also satisfies a stage when its path matches the legacy shape `prompts/<timestamp>-<stage>.md`.

`specops context` and `specops next` continue to preserve existing `command` and `reason` fields. For semantic gates they also expose `note_command` and generic `suggested_question_prompts`.

Skill payloads must define `continue` as scoped to exactly one current semantic gate. After every semantic command, the agent must refresh `specops context <run-id>`; if the refreshed context reports another semantic gate, the agent must stop and ask that gate's content-aware questions before running anything else.

## Consequences

Positive:

- A single `continue` answer cannot carry an agent across multiple semantic production gates.
- The CLI remains deterministic and non-AI while enforcing the durable provenance checkpoint.
- Existing v0.1.1 prompt artifacts remain usable as stage notes.
- Human and JSON outputs give agents the exact note command needed for the current gate.

Negative / tradeoffs:

- Scripted runs must add `specops note` calls before `refine`, `harden`, and `synthesize`.
- Deterministic fallback mode is no longer enough by itself to cross semantic production gates.
- Decisions and apply gates remain enforced by skill wording and existing decision/apply mechanics; CLI hard blocking for those gates can be added later if needed.

## Affected docs

- docs/interfaces/cli_commands.md
- docs/interfaces/skill_payload_contract.md
- docs/versions/release_decision_log.md
