---
id: delegated-installer-json
title: Delegated Installer JSON
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Delegated Installer JSON

## Report schema

```json
{
  "schema": 1,
  "name": "specops",
  "version": "0.1.0",
  "operation": "install",
  "kind": "delegated",
  "targets": {
    "claude": {"files": []},
    "codex": {"files": []},
    "tools": {"files": []}
  },
  "warnings": []
}
```

## File entry

```json
{
  "path": "/absolute/path",
  "sha256": "hex"
}
```

## Operation semantics

### `plan`

Reports files that would be written. `sha256` optional.

### `install`

Writes files. `sha256` required.

### `uninstall`

Reports files owned by the delegated installer. May not remove real files in staged mode unless invoked directly for uninstall.

## stdout/stderr

- `--json`: stdout contains only report JSON.
- Human logs go to stderr.
