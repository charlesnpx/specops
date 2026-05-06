# Agent Instructions for This Specification Repository

You are working in a **specification repository**, not an implementation repository.

Before editing canonical docs, read:

1. `docs/CANON.md`
2. `docs/versions/v0_scope.md`
3. `docs/cross_cutting/mise_en_place_contract.md`
4. `docs/interfaces/cli_commands.md`
5. `.specops/process.md`

Rules:

- Treat ADRs as append-only decision history.
- Do not silently rewrite accepted decisions; create a superseding ADR if needed.
- Separate process artifacts from canonical spec artifacts.
- Put draft run artifacts under `.specops/runs/`.
- Put reviewed provenance under `docs/research/refinery/`.
- Put accepted specification truth under `docs/` proper.
- Every CLI or installer behavior change must update at least one interface doc and usually one ADR.
- The mise-en-place delegated installer contract is normative unless explicitly superseded by a new ADR.
