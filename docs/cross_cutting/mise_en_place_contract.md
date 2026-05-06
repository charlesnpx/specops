---
id: mise-en-place-contract
title: Mise-en-place Contract
doc_type: cross_cutting_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Mise-en-place Contract

## Normative requirement

SpecOps must be installable as a delegated repo managed by `mise-en-place`.

## Required operations

```sh
./install-skill.sh --plan --target all --json
./install-skill.sh --install --target all --json
./install-skill.sh --install --target all --json --install-root /tmp/stage
./install-skill.sh --uninstall --target all --json
```

## Required flags

```text
--target claude|codex|tools|all
--plan
--install
--uninstall
--json
--install-root <absolute-dir>
```

## Required JSON shape

```json
{
  "schema": 1,
  "name": "specops",
  "version": "0.1.2",
  "operation": "install",
  "kind": "delegated",
  "targets": {
    "claude": {"files": [{"path": "/tmp/stage/.claude/skills/specops/SKILL.md", "sha256": "..."}]},
    "codex": {"files": [{"path": "/tmp/stage/.codex/skills/specops/SKILL.md", "sha256": "..."}]},
    "tools": {"files": [{"path": "/tmp/stage/.local/bin/specops", "sha256": "..."}]}
  },
  "warnings": []
}
```

## Interpretation

- `plan` may omit `sha256`.
- `install` must include `sha256`.
- `uninstall` reports files that would be removed or were removed.
- Paths must be absolute.
- Under `--install-root`, all installed paths must be inside the root.
- Stdout is JSON-only when `--json` is set.

## Collision handling

SpecOps delegated installer should not directly prompt about real user destination collisions during staged install. It stages content; `mise-en-place` compares, prompts, backs up, writes, and records ownership.
