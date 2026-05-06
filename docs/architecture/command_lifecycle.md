---
id: command-lifecycle
title: Command Lifecycle
doc_type: architecture
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Command Lifecycle

Dynamic commands should follow this internal sequence:

```text
1. Load project config and scaffold version.
2. Load or create run state.
3. Read relevant target repo indexes and canonical docs.
4. Select command definition and templates.
5. Build a prompt packet or deterministic operation plan.
6. Invoke backend if needed.
7. Parse backend output.
8. Validate against schemas.
9. Write run artifacts.
10. Update run state.
11. Print summary and recommended next command.
```

## Static vs dynamic commands

Static introspection commands:

```sh
specops prompt show refine
specops template show adr
specops schema show spec_delta
```

Dynamic execution commands:

```sh
specops refine run-001
specops synthesize run-001
specops compile run-001
specops eval --gold gold --candidate candidate
```

Static commands may print the same content repeatedly. Dynamic commands must inspect state and produce stateful artifacts.
