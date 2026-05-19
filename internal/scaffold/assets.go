package scaffold

import "fmt"

const Version = "0.1.3-dev"

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
- Record a stage note and pass an authored artifact with --from before semantic commands: refine, harden, and synthesize.
- During synthesize, put full canonical document bodies in spec_delta.patch_items[].content when generated docs would be too thin. patch_plan is notes only; affected_docs is coverage only.
- At the apply gate, use specops supersede-synthesis <run-id> --from <spec_delta.json> after an apply-stage note when the patch plan is structurally valid but semantically too thin.
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
new_concepts: []
requirements: []
constraints: []
assumptions: []
ambiguities: []
options: []
recommendations: []
decisions: []
affected_docs: []
version_scope_changes: []
implementation_phase_changes: []
acceptance_criteria: []
open_questions: []
risks: []

# Human-readable compile notes only. Do not put full document bodies here.
patch_plan: []

# Optional exact file-level patches. Use this when canonical docs need rich
# authored content; otherwise compile generates skeletal docs from affected_docs.
patch_items: []
# patch_items:
#   - action: create
#     path: docs/CANON.md
#     title: Create canonical frame
#     decision_ids: [D001]
#     content: |
#       ---
#       id: canon
#       title: ...
#       doc_type: canon
#       status: accepted
#       normative: true
#       version_scope: v0_required
#       last_reviewed: 2026-05-06
#       ---
#
#       # Canonical Title
#
#       Full authored document body.
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
    "status": {"type": "string"},
    "artifacts": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "id": {"type": "string"},
          "type": {
            "description": "Artifact kind. Superseded current artifacts may use types such as superseded_spec_delta or superseded_patch_plan.",
            "type": "string"
          },
          "path": {"type": "string"},
          "stage": {"type": "string"},
          "created_at": {"type": "string"}
        }
      }
    },
    "next": {
      "type": "object",
      "properties": {
        "command": {"type": "string"},
        "reason": {"type": "string"},
        "stage": {"type": "string"},
        "gate_kind": {"type": "string"},
        "context_command": {"type": "string"},
        "note_command": {"type": "string"},
        "suggested_question_prompts": {"type": "array", "items": {"type": "string"}},
        "human_input_recommended": {"type": "boolean"}
      }
    }
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
    "run_id": {"type": "string"},
    "source_summary": {"type": "string"},
    "new_concepts": {"type": "array", "items": {"type": "string"}},
    "requirements": {"type": "array", "items": {"type": "string"}},
    "constraints": {"type": "array", "items": {"type": "string"}},
    "assumptions": {"type": "array", "items": {"type": "string"}},
    "ambiguities": {"type": "array", "items": {"type": "string"}},
    "options": {"type": "array", "items": {"type": "string"}},
    "recommendations": {"type": "array", "items": {"type": "string"}},
    "affected_docs": {"type": "array", "items": {"type": "string"}},
    "version_scope_changes": {"type": "array", "items": {"type": "string"}},
    "implementation_phase_changes": {"type": "array", "items": {"type": "string"}},
    "acceptance_criteria": {"type": "array", "items": {"type": "string"}},
    "open_questions": {"type": "array", "items": {"type": "string"}},
    "risks": {"type": "array", "items": {"type": "string"}},
    "patch_plan": {
      "description": "Human-readable compile notes only. Full canonical document bodies belong in patch_items[].content.",
      "type": "array",
      "items": {"type": "string"}
    },
    "patch_items": {
      "description": "Optional exact file-level patches. Use when compile must preserve rich authored canonical document content.",
      "type": "array",
      "items": {
        "type": "object",
        "required": ["path", "content"],
        "properties": {
          "id": {"type": "string"},
          "action": {"type": "string"},
          "path": {"type": "string"},
          "title": {"type": "string"},
          "content": {
            "description": "Exact file content to place in the patch plan item.",
            "type": "string"
          },
          "decision_ids": {"type": "array", "items": {"type": "string"}}
        }
      }
    }
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
