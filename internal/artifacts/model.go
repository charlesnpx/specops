package artifacts

import "github.com/specops/specops/internal/runstate"

type SpecDelta struct {
	Schema                     int                 `json:"schema"`
	RunID                      string              `json:"run_id"`
	SourceSummary              string              `json:"source_summary"`
	NewConcepts                []string            `json:"new_concepts,omitempty"`
	Requirements               []string            `json:"requirements,omitempty"`
	Constraints                []string            `json:"constraints,omitempty"`
	Assumptions                []string            `json:"assumptions,omitempty"`
	Ambiguities                []string            `json:"ambiguities,omitempty"`
	Options                    []string            `json:"options,omitempty"`
	Recommendations            []string            `json:"recommendations,omitempty"`
	Decisions                  []runstate.Decision `json:"decisions"`
	AffectedDocs               []string            `json:"affected_docs"`
	VersionScopeChanges        []string            `json:"version_scope_changes,omitempty"`
	ImplementationPhaseChanges []string            `json:"implementation_phase_changes,omitempty"`
	AcceptanceCriteria         []string            `json:"acceptance_criteria,omitempty"`
	OpenQuestions              []string            `json:"open_questions,omitempty"`
	Risks                      []string            `json:"risks,omitempty"`
	PatchPlan                  []string            `json:"patch_plan"`
	PatchItems                 []PatchItem         `json:"patch_items,omitempty"`
}

type PatchPlan struct {
	Schema           int             `json:"schema"`
	RunID            string          `json:"run_id"`
	CreatedAt        string          `json:"created_at"`
	CompilerContract string          `json:"compiler_contract,omitempty"`
	AcceptedOnly     bool            `json:"accepted_only"`
	Inputs           PatchPlanInputs `json:"inputs,omitempty"`
	Health           PatchPlanHealth `json:"health"`
	Items            []PatchItem     `json:"items"`
}

type PatchPlanInputs struct {
	SpecDeltaSHA256         string `json:"spec_delta_sha256,omitempty"`
	AcceptedDecisionsSHA256 string `json:"accepted_decisions_sha256,omitempty"`
	ProvenanceSHA256        string `json:"provenance_sha256,omitempty"`
}

type PatchPlanHealth struct {
	Stale             bool     `json:"stale"`
	Incomplete        bool     `json:"incomplete"`
	StaleReasons      []string `json:"stale_reasons,omitempty"`
	IncompleteReasons []string `json:"incomplete_reasons,omitempty"`
}

type PatchItem struct {
	ID          string   `json:"id"`
	Action      string   `json:"action"`
	Path        string   `json:"path"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	DecisionIDs []string `json:"decision_ids"`
}
