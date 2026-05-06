package runstate

type Status string

const (
	StatusCreated           Status = "created"
	StatusIngested          Status = "ingested"
	StatusIntakeComplete    Status = "intake_complete"
	StatusRefined           Status = "refined"
	StatusHardened          Status = "hardened"
	StatusSynthesized       Status = "synthesized"
	StatusAwaitingDecisions Status = "awaiting_decisions"
	StatusDecisionsAccepted Status = "decisions_accepted"
	StatusCompiled          Status = "compiled"
	StatusPlanned           Status = "planned"
	StatusApplied           Status = "applied"
	StatusAudited           Status = "audited"
	StatusEvaluated         Status = "evaluated"
)

type RunState struct {
	Schema    int                 `json:"schema"`
	RunID     string              `json:"run_id"`
	Name      string              `json:"name,omitempty"`
	Status    Status              `json:"status"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at,omitempty"`
	Inputs    []InputRef          `json:"inputs"`
	Artifacts []ArtifactRef       `json:"artifacts"`
	Decisions map[string]Decision `json:"decisions"`
	Next      *NextAction         `json:"next,omitempty"`
}

type InputRef struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Path     string `json:"path"`
	SHA256   string `json:"sha256,omitempty"`
	AddedAt  string `json:"added_at"`
	Original string `json:"original,omitempty"`
}

type ArtifactRef struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Path      string `json:"path"`
	CreatedAt string `json:"created_at"`
}

type Decision struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Status         string   `json:"status"`
	Options        []string `json:"options,omitempty"`
	Recommendation string   `json:"recommendation,omitempty"`
	Rationale      string   `json:"rationale,omitempty"`
	Text           string   `json:"text,omitempty"`
	ADRRequired    bool     `json:"adr_required,omitempty"`
	AffectedDocs   []string `json:"affected_docs,omitempty"`
}

type NextAction struct {
	Command string `json:"command"`
	Reason  string `json:"reason"`
}
