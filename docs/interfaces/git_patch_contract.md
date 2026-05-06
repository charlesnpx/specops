---
id: git-patch-contract
title: Git Patch Contract
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Git Patch Contract

Patch plans are structured before application.

```json
{
  "schema": 1,
  "run_id": "run-...",
  "files": [
    {
      "path": "docs/decisions/ADR-0001-example.md",
      "operation": "create",
      "source_decision_ids": ["decision-001"],
      "content_ref": ".specops/runs/.../patches/..."
    }
  ],
  "preconditions": [
    "working_tree_clean_or_allow_dirty",
    "accepted_decisions_present"
  ]
}
```
