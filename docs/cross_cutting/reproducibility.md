---
id: reproducibility
title: Reproducibility
doc_type: cross_cutting_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Reproducibility

## Lockfiles

`.specops.lock` records:

```yaml
schema: 1
specops_version: 0.1.2
scaffold_version: 0.1.2
templates_hash: sha256:...
schemas_hash: sha256:...
evals_hash: sha256:...
```

## Gold/candidate eval

Reproduction evals must separate:

```text
process mining set
generation fixture set
gold evaluation set
```

Do not let a generation run read the gold repo unless the task is explicitly process mining rather than reproduction.
