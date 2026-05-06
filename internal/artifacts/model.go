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
}

type PatchPlan struct {
	Schema       int         `json:"schema"`
	RunID        string      `json:"run_id"`
	CreatedAt    string      `json:"created_at"`
	AcceptedOnly bool        `json:"accepted_only"`
	Items        []PatchItem `json:"items"`
}

type PatchItem struct {
	ID          string   `json:"id"`
	Action      string   `json:"action"`
	Path        string   `json:"path"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	DecisionIDs []string `json:"decision_ids"`
}
