---
id: registry-entry-contract
title: mise-en-place Registry Entry Contract
doc_type: interface_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# mise-en-place Registry Entry Contract

The intended registry entry is delegated:

```yaml
delegated:
  specops:
    repo: github.com/charlesnpx/specops
    channel: latest-release
```

During early private development:

```yaml
delegated:
  specops:
    repo: github.com/charlesnpx/specops
    ref: main
    visibility: private
    optional: true
```

Release-quality entries should prefer stable semver release tags.

## Version channels

- `ref`: exact tag/commit/branch pin.
- `channel: latest-release`: resolve highest stable semver release.
- `fallback_ref`: allowed only for untagged development repos.

## Required delegated compatibility

The repo must expose either:

```text
./install-skill.sh
```

or a configured installer command that supports the same flags and JSON contract.
