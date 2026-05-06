package scaffold

import "fmt"

const Version = "0.1.1"

func files(mode, agent string) map[string]string {
	return map[string]string{
		".specops.yaml":                                    specopsYAML(mode, agent),
		".specops.lock":                                    lockJSON(mode, agent),
		"AGENTS.md":                                        agentsMD,
		"CLAUDE.md":                                        claudeMD,
		".specops/process.md":                              processMD,
		".specops/status_model.md":                         statusModelMD,
		".specops/templates/adr.md":                        adrTemplate,
		".specops/templates/spec_delta.yaml":               specDeltaTemplate,
		".specops/evals/structure_eval.md":                 structureEval,
		".specops/evals/decision_coverage_eval.md":         decisionEval,
		".specops/evals/lineage_eval.md":                   lineageEval,
		".specops/checklists/canonical_patch_checklist.md": patchChecklist,
		".specops/schemas/run_state.schema.json":           runStateSchema,
		".specops/schemas/spec_delta.schema.json":          specDeltaSchema,
		".specops/schemas/decision.schema.json":            decisionSchema,
	}
}

func specopsYAML(mode, agent string) string {
	return fmt.Sprintf("schema: 1\ntool: specops\nscaffold_version: %s\nmode: %s\nagent: %s\n", Version, mode, agent)
}

func lockJSON(mode, agent string) string {
	return fmt.Sprintf("{\n  \"schema\": 1,\n  \"tool\": \"specops\",\n  \"scaffold_version\": \"%s\",\n  \"mode\": \"%s\",\n  \"agent\": \"%s\"\n}\n", Version, mode, agent)
}

const agentsMD = `# Agent Instructions

This repository uses SpecOps.

Before editing canonical specification files, read:

1. docs/CANON.md when present
2. docs/versions/v0_scope.md when present
3. .specops/process.md

Rules:

- Put draft run artifacts under .specops/runs/.
- Put reviewed provenance under docs/research/refinery/.
- Treat ADRs as append-only decision history.
- Do not patch canonical docs before decisions are accepted.
`

const claudeMD = `# Claude Instructions

This repository uses SpecOps. Follow AGENTS.md and .specops/process.md before changing canonical specification files.
`

const processMD = `# SpecOps Process

Canonical loop:

raw input -> intake -> refine -> harden -> synthesize -> decide -> compile -> plan -> apply -> audit -> eval

Rules:

- Do not patch canonical docs before decisions are accepted.
- Keep run artifacts under .specops/runs/.
- Promote reviewed provenance to docs/research/refinery/.
- Use ADRs for consequential accepted decisions.
- Update interfaces when behavior changes.
`

const statusModelMD = `# Status Model

Run statuses:

created
ingested
intake_complete
refined
hardened
synthesized
awaiting_decisions
decisions_accepted
compiled
planned
applied
audited
evaluated
`

const adrTemplate = `---
id: ADR-0000
title: Decision title
doc_type: adr
status: proposed
normative: true
version_scope: v0_required
date: YYYY-MM-DD
---

# ADR-0000: Decision title

## Status

Proposed.

## Context

## Decision

## Consequences

## Affected docs
`

const specDeltaTemplate = `schema: 1
run_id: run-YYYYMMDD-HHMMSS-slug
source_summary: ""
decisions: []
affected_docs: []
patch_plan: []
`

const structureEval = `# Structure Eval

Check that required SpecOps directories, canonical docs, ADRs, interfaces, version scope, and generated indexes exist.
`

const decisionEval = `# Decision Coverage Eval

Check that consequential accepted decisions are represented by ADRs and reflected in affected docs.
`

const lineageEval = `# Lineage Eval

Check that canonical changes can be traced back to run artifacts, decisions, and reviewed provenance.
`

const patchChecklist = `# Canonical Patch Checklist

- Run artifacts exist under .specops/runs/.
- Decisions are accepted before canonical docs change.
- ADRs are append-only.
- Interface behavior changes update interface docs.
- Audit passes after apply.
`

const runStateSchema = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "RunState",
  "type": "object",
  "required": ["schema", "run_id", "status", "created_at", "inputs", "artifacts", "decisions"],
  "properties": {
    "schema": {"const": 1},
    "run_id": {"type": "string"},
    "status": {"type": "string"}
  }
}
`

const specDeltaSchema = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "SpecDelta",
  "type": "object",
  "required": ["schema", "run_id", "source_summary", "decisions", "affected_docs", "patch_plan"],
  "properties": {
    "schema": {"const": 1},
    "run_id": {"type": "string"}
  }
}
`

const decisionSchema = `{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Decision",
  "type": "object",
  "required": ["id", "title", "status", "options", "recommendation"],
  "properties": {
    "id": {"type": "string"},
    "title": {"type": "string"},
    "status": {"type": "string"}
  }
}
`
