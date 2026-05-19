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
- Record a stage note and pass an authored artifact with `--from` before semantic commands: `refine`, `harden`, and `synthesize`.
- During `synthesize`, put full canonical document bodies in `spec_delta.patch_items[].content` when generated docs would be too thin. `patch_plan` is notes only; `affected_docs` is coverage only.
- At the apply gate, use `specops supersede-synthesis <run-id> --from <spec_delta.json>` after an apply-stage note when the patch plan is structurally valid but semantically too thin.
- Use ADRs for consequential accepted decisions.
- Update interfaces when behavior changes.
