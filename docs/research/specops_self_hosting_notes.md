---
id: specops-self-hosting-notes
title: SpecOps Self-hosting Notes
doc_type: research_note
status: accepted
normative: false
version_scope: v0_required
last_reviewed: 2026-05-06
---


# SpecOps Self-hosting Notes

This repository is intentionally self-referential: it is a spec repository for the tool that produces spec repositories.

Self-hosting loop:

```text
this spec repo
  -> implement specops CLI
  -> use specops CLI to evolve this spec repo
  -> use specops to reproduce context-sink-spec
  -> improve specops based on eval failures
```

The first milestone is not perfect automation. It is making the process explicit, testable, and installable.
