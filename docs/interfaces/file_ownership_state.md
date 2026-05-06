---
id: file-ownership-state
title: File Ownership State
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# File Ownership State

SpecOps itself should avoid duplicating `mise-en-place` ownership state for installed files. In target repos, it should track scaffold ownership separately.

## `.specops.lock`

```yaml
schema: 1
scaffold:
  version: 0.1.0
  files:
    - path: .specops/process.md
      sha256: ...
      managed: true
      local_modified: false
```

## Upgrade behavior

- Identical managed files can be updated automatically.
- Locally modified files require plan output and confirmation.
- `--backup` creates `.backup`, `.backup.1`, etc.
