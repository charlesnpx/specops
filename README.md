# SpecOps CLI Specification Repository

Status: **v0 implementation specification**  
Generated: **2026-05-06**

This repository specifies a self-referential tool tentatively named **`specops`**: a CLI-backed specification-production toolchain that helps agents and humans turn messy ideation into a durable specification repository.

The first implementation target is a **Go CLI** that can be installed as a **delegated skill/tool through `mise-en-place`**. It should install dual Claude Code and Codex payloads, expose the delegated installer contract expected by `mise-en-place`, and then operate inside target repositories to scaffold `.specops/`, manage run artifacts, compile accepted decisions into canonical specification docs, and evaluate reproduction attempts.

The specification follows the same document strategy it asks downstream users to produce:

```text
CANON.md
  current accepted truth

decisions/
  ADRs that explain why choices were made

architecture/
  system overview, lifecycles, deployment/runtime shape

subsystems/
  owned boundaries and responsibilities

cross_cutting/
  concepts that cut across subsystems

interfaces/
  CLI, installer, JSON, run, and agent contracts

implementation/
  build phases and acceptance criteria

versions/
  v0 scope and post-v0 candidates

research/
  process reconstruction and alternatives

.specops/
  local process scaffold for how this spec was produced and how it should evolve
```

## Immediate thesis

`specops` should be:

```text
A repo-native specification compiler and review harness.
```

It should not merely print static prompts. It should install local process rails, ingest messy inputs, maintain run state, call an agent or relay backend when useful, produce typed run artifacts, require human decision gates, apply accepted patches, and audit/evaluate the resulting spec repository.

## Contract pressure

This repo treats the `mise-en-place` delegated installer contract as a normative external interface. The tool must provide an installer command compatible with:

```sh
./install-skill.sh --plan --target all --json
./install-skill.sh --install --target all --json
./install-skill.sh --install --target all --json --install-root /tmp/stage
./install-skill.sh --uninstall --target all --json
```

The staged install must report absolute files under the requested install root and must produce JSON-only stdout when `--json` is set.

## Key docs

- [CANON](docs/CANON.md)
- [v0 Scope](docs/versions/v0_scope.md)
- [Mise-en-place Contract](docs/cross_cutting/mise_en_place_contract.md)
- [CLI Commands](docs/interfaces/cli_commands.md)
- [Implementation Strategy](docs/implementation/implementation_strategy.md)
- [Open Questions](docs/OPEN_QUESTIONS.md)
