# SpecOps Process

Canonical loop:

```text
raw input
  -> intake
  -> refine
  -> harden
  -> synthesize
  -> decide
  -> compile
  -> plan
  -> apply
  -> audit
  -> eval
```

Rules:

- Do not patch canonical docs before decisions are accepted.
- Keep run artifacts under `.specops/runs/`.
- Promote reviewed provenance to `docs/research/refinery/`.
- Record a stage note before semantic commands: `refine`, `harden`, and `synthesize`.
- Use ADRs for consequential accepted decisions.
- Update interfaces when behavior changes.
