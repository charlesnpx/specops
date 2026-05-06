---
id: repo-taxonomy
title: Repository Taxonomy
doc_type: cross_cutting_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Repository Taxonomy

SpecOps-produced repos should use a predictable taxonomy:

```text
docs/CANON.md
docs/decisions/
docs/architecture/
docs/subsystems/
docs/cross_cutting/
docs/interfaces/
docs/versions/
docs/implementation/
docs/research/
docs/templates/
docs/_generated/
```

The CLI should support project-specific variants but should default to this taxonomy.
