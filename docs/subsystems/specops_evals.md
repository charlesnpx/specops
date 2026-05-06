---
id: specops-evals
title: SpecOps Evaluation Engine
doc_type: subsystem_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# specops_evals

## Purpose

Evaluates spec repositories and reproduction attempts.

## Eval classes

```text
structure_eval
concept_coverage_eval
decision_coverage_eval
lineage_eval
contradiction_eval
implementation_usefulness_eval
```

## Deterministic checks

- Required directories exist.
- Required front matter exists.
- Accepted ADRs link affected docs.
- v0 items have acceptance criteria.
- No broken links.
- Schema validation passes.

## Semantic checks

- Core concepts recovered.
- Decision logic matches gold.
- Traceability chains are present.
- Candidate does not contradict derived spec.

## Acceptance criteria

- Eval reports are JSON plus Markdown summary.
- Scores explain missing/weak/contradicted items.
- Eval can compare `gold` vs `candidate` repos without allowing generation to read gold.
