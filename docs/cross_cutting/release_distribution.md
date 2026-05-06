---
id: release-distribution
title: Release and Distribution
doc_type: cross_cutting_spec
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Release and Distribution

## Release artifacts

Each CLI release should publish:

```text
specops_Darwin_arm64.tar.gz
specops_Darwin_x86_64.tar.gz
specops_Linux_arm64.tar.gz
specops_Linux_x86_64.tar.gz
specops_checksums.txt
```

Each archive should include:

```text
specops binary
install-skill.sh
skills/claude/specops/SKILL.md
skills/codex/specops/SKILL.md
assets/scaffold/...
README.md
LICENSE
```

## Release process

```text
git tag v0.1.1
git push origin v0.1.1
GitHub Actions runs tests
GoReleaser builds archives and checksums
GitHub Release publishes assets
mise-en-place latest-release can resolve v0.1.1
```

## Checksums

Checksums are part of the release acceptance criteria. The install contract reports file-level hashes after staging; the release process reports archive-level checksums.
