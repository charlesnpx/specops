package artifacts

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/specops/specops/internal/runstate"
)

type RunContext struct {
	Schema           int                          `json:"schema"`
	RunID            string                       `json:"run_id"`
	Name             string                       `json:"name,omitempty"`
	Status           runstate.Status              `json:"status"`
	SourceSummary    string                       `json:"source_summary,omitempty"`
	OperatorGuidance OperatorGuidance             `json:"operator_guidance"`
	Artifacts        []runstate.ArtifactRef       `json:"artifacts"`
	Decisions        map[string]runstate.Decision `json:"decisions"`
	PatchPlan        *PatchPlan                   `json:"patch_plan,omitempty"`
	NextGate         *runstate.NextAction         `json:"next_gate,omitempty"`
}

type OperatorGuidance struct {
	ContextCommand        string   `json:"context_command,omitempty"`
	NoteCommand           string   `json:"note_command,omitempty"`
	GateKind              string   `json:"gate_kind,omitempty"`
	Stage                 string   `json:"stage,omitempty"`
	HumanInputRecommended bool     `json:"human_input_recommended"`
	SuggestedQuestions    []string `json:"suggested_question_prompts,omitempty"`
	ControlQuestion       string   `json:"control_question"`
}

type NoteResult struct {
	RunID    string               `json:"run_id"`
	Status   runstate.Status      `json:"status"`
	Artifact runstate.ArtifactRef `json:"artifact"`
}

func Context(repo, runID string) (RunContext, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return RunContext{}, err
	}
	var plan *PatchPlan
	if loaded, err := loadPlanForResolvedRun(repo, state.RunID); err == nil {
		plan = &loaded
	}
	guidance := OperatorGuidance{ControlQuestion: "Should the operator continue, pause, or change direction?"}
	if state.Next != nil {
		guidance.ContextCommand = state.Next.ContextCommand
		guidance.NoteCommand = state.Next.NoteCommand
		guidance.GateKind = state.Next.GateKind
		guidance.Stage = state.Next.Stage
		guidance.HumanInputRecommended = state.Next.HumanInputRecommended
		guidance.SuggestedQuestions = state.Next.SuggestedQuestions
	}
	return RunContext{
		Schema:           1,
		RunID:            state.RunID,
		Name:             state.Name,
		Status:           state.Status,
		SourceSummary:    readRunFile(repo, state.RunID, "inputs", "input_summary.md"),
		OperatorGuidance: guidance,
		Artifacts:        state.Artifacts,
		Decisions:        state.Decisions,
		PatchPlan:        plan,
		NextGate:         state.Next,
	}, nil
}

func Note(repo, runID, stage, text string) (NoteResult, error) {
	if stage == "" {
		return NoteResult{}, fmt.Errorf("--stage is required")
	}
	if text == "" {
		return NoteResult{}, fmt.Errorf("--text is required")
	}
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return NoteResult{}, err
	}
	content := text
	if raw, err := os.ReadFile(text); err == nil {
		content = string(raw)
	}
	now := time.Now().UTC()
	slug := safeName(stage)
	if slug == "" {
		slug = "note"
	}
	name := fmt.Sprintf("%s-%s.md", now.Format("20060102-150405"), slug)
	body := fmt.Sprintf("# Operator Note\n\nRun: `%s`\nStage: `%s`\nRecorded: `%s`\n\n%s\n", state.RunID, stage, now.Format(time.RFC3339), strings.TrimSpace(content))
	ref, err := writePrompt(repo, state, name, slug, body)
	if err != nil {
		return NoteResult{}, err
	}
	if err := runstate.Save(repo, state); err != nil {
		return NoteResult{}, err
	}
	return NoteResult{RunID: state.RunID, Status: state.Status, Artifact: ref}, nil
}

func hasStageNote(state *runstate.RunState, stage string) bool {
	want := safeName(stage)
	if want == "" {
		return false
	}
	for _, artifact := range state.Artifacts {
		if artifact.Type != "prompt" {
			continue
		}
		if safeName(artifact.Stage) == want {
			return true
		}
		if artifact.Stage == "" && legacyPromptPathMatchesStage(artifact.Path, want) {
			return true
		}
	}
	return false
}

func legacyPromptPathMatchesStage(path, stage string) bool {
	path = filepath.ToSlash(path)
	return strings.HasPrefix(path, "prompts/") && strings.HasSuffix(path, "-"+stage+".md")
}

func writePrompt(repo string, state *runstate.RunState, filename, stage, content string) (runstate.ArtifactRef, error) {
	path := filepath.Join(runstate.RunDir(repo, state.RunID), "prompts", filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return runstate.ArtifactRef{}, err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return runstate.ArtifactRef{}, err
	}
	ref := runstate.ArtifactRef{ID: fmt.Sprintf("prompt-%03d", len(state.Artifacts)+1), Type: "prompt", Path: "prompts/" + filename, Stage: stage, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	state.Artifacts = append(state.Artifacts, ref)
	return ref, nil
}

func loadPlanForResolvedRun(repo, runID string) (PatchPlan, error) {
	path := filepath.Join(runstate.RunDir(repo, runID), "patches", "patch_plan.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return PatchPlan{}, err
	}
	var plan PatchPlan
	if err := json.Unmarshal(raw, &plan); err != nil {
		return PatchPlan{}, err
	}
	return plan, nil
}
