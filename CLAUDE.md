# Claude Code Instructions

Use this repo as a SpecOps target repository.

Useful local command docs live in `.specops/commands/`.

When asked to refine or patch the specification:

1. Create or update a refinery note first.
2. Extract ambiguities, options, recommendations, and decisions.
3. Ask for or infer explicit decision acceptance only when the user has clearly accepted recommendations.
4. Patch canonical docs only for accepted decisions.
5. Preserve traceability from source material to ADRs, canonical docs, interfaces, and implementation phases.

Do not treat `.specops/` as the product specification. `.specops/` is the local process scaffold; the actual specification lives under `docs/`.
