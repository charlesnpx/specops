---
id: mise-integration
title: Mise-en-place Integration Architecture
doc_type: architecture
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Mise-en-place Integration Architecture

## Installation path

SpecOps should be published as a delegated repo entry in `mise-en-place` registry:

```yaml
delegated:
  specops:
    repo: github.com/charlesnpx/specops
    channel: latest-release
    optional: false
```

`mise-en-place` clones the delegated repo and executes the installer contract.

## Delegated repo contents

```text
specops/
  install-skill.sh
  cmd/specops/
  internal/...
  skills/claude/specops/SKILL.md
  skills/codex/specops/SKILL.md
  payloads/scaffold/...
  registry/manifest.yaml
```

## Installer staging

When called with:

```sh
./install-skill.sh --install --target all --json --install-root /tmp/stage
```

The installer writes staged files under:

```text
/tmp/stage/.claude/skills/specops/SKILL.md
/tmp/stage/.codex/skills/specops/SKILL.md
/tmp/stage/.local/bin/specops
```

and returns JSON paths inside `/tmp/stage`.

`mise-en-place` then maps those staged paths to the user's home, applies collision handling, backups, state recording, and ownership.

## Self-bootstrap issue

Because the installer may need to install the `specops` binary before it is on PATH, `install-skill.sh` should be a small shell wrapper that uses either:

1. A prebuilt release asset vendored or downloaded by `mise-en-place` process, if available.
2. A repo-local compiled binary next to the wrapper, if the release package includes it.
3. `go run -ldflags "-X main.version=<resolved-version>" ./cmd/specops install-skill ...` as a developer fallback when Go is available.

For public release packages, prefer option 2: GitHub Release archive contains the binary plus skill payloads and install wrapper.

The wrapper must not delegate to a `specops` binary found on `PATH`; `mise-en-place`
runs the wrapper from a resolved checkout, and PATH delegation can mix an older
checkout with a newer local binary.
