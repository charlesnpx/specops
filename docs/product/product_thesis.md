---
id: product-thesis
title: Product Thesis
doc_type: product_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Product Thesis

SpecOps helps a human and an agent produce high-quality specification repositories from unstructured design work.

It exists because useful specification production is not just summarization. The observed Context Sink process followed a stronger pattern:

```text
raw idea
  -> product reframing
  -> design-axis separation
  -> central-risk discovery
  -> adversarial hardening
  -> derived spec
  -> recursive subsystem deepening
  -> repository emission
```

SpecOps should make this repeatable.

## User promise

```text
Give SpecOps a messy conversation, process trace, product brief, or relay transcript.
It helps convert it into accepted decisions and a durable spec repository.
```

## Primary users

- Technical founders and solo builders designing complex products.
- Engineers using Claude Code or Codex to turn strategy into buildable specs.
- Teams that want docs-as-code specification repositories before implementation.
- People already using `mise-en-place` to manage Claude/Codex skills.

## Core value

SpecOps lowers the friction of producing rigorous specs while preserving the things that make specs trustworthy:

- decision history
- ambiguity handling
- explicit options and recommendations
- accepted/rejected/deferred status
- traceability
- implementation phases
- audit/eval
- repo-native diffs

## Product rule

```text
The CLI may draft and patch, but accepted decisions belong to the human.
```
