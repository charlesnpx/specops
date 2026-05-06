package artifacts

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/specops/specops/internal/input"
	"github.com/specops/specops/internal/runstate"
)

func TestContextForRunStages(t *testing.T) {
	repo, state := newIngestedRun(t)

	context, err := Context(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if context.Status != runstate.StatusIngested {
		t.Fatalf("status = %s, want %s", context.Status, runstate.StatusIngested)
	}
	if context.SourceSummary == "" {
		t.Fatal("expected source summary")
	}
	if context.NextGate == nil || context.NextGate.Command == "" || context.NextGate.Reason == "" {
		t.Fatal("next gate should preserve command and reason")
	}

	if _, err := Intake(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	context, err = Context(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if !context.OperatorGuidance.HumanInputRecommended || context.OperatorGuidance.Stage != "refine" {
		t.Fatalf("unexpected guidance: %+v", context.OperatorGuidance)
	}
	if context.NextGate.NoteCommand == "" || context.OperatorGuidance.NoteCommand == "" {
		t.Fatalf("semantic gate should expose note command: %+v / %+v", context.NextGate, context.OperatorGuidance)
	}

	if _, err := Note(repo, state.RunID, "refine", "refine this into reviewable notes"); err != nil {
		t.Fatal(err)
	}
	if _, err := Refine(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "synthesize", "prepare explicit decisions"); err != nil {
		t.Fatal(err)
	}
	if _, err := Synthesize(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	context, err = Context(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if context.Status != runstate.StatusAwaitingDecisions || len(context.Decisions) == 0 {
		t.Fatalf("expected awaiting decisions with decisions, got %s / %d", context.Status, len(context.Decisions))
	}

	if _, err := SetDecision(repo, state.RunID, "DEC-0001", "accepted", ""); err != nil {
		t.Fatal(err)
	}
	if _, err := Compile(repo, state.RunID, true); err != nil {
		t.Fatal(err)
	}
	if _, err := MarkPlanned(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	context, err = Context(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if context.PatchPlan == nil || len(context.PatchPlan.Items) == 0 {
		t.Fatal("expected patch plan in context")
	}
}

func TestNoteCreatesPromptArtifactWithoutChangingStatus(t *testing.T) {
	repo, state := newIngestedRun(t)
	before := state.Status

	result, err := Note(repo, state.RunID, "refine", "please preserve the operator intent")
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != before {
		t.Fatalf("status changed to %s, want %s", result.Status, before)
	}
	if result.Artifact.Type != "prompt" || !strings.HasPrefix(result.Artifact.Path, "prompts/") {
		t.Fatalf("unexpected artifact: %+v", result.Artifact)
	}
	if result.Artifact.Stage != "refine" {
		t.Fatalf("artifact stage = %q, want refine", result.Artifact.Stage)
	}
	if _, err := os.Stat(filepath.Join(runstate.RunDir(repo, state.RunID), result.Artifact.Path)); err != nil {
		t.Fatal(err)
	}
	loaded, err := runstate.Load(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Status != before {
		t.Fatalf("persisted status changed to %s, want %s", loaded.Status, before)
	}
}

func TestFromArtifactsArePreservedAndTransitionsEnforced(t *testing.T) {
	repo, state := newIngestedRun(t)
	if _, err := Intake(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	refinedPath := filepath.Join(repo, "refined.md")
	refined := "# Authored refined notes\n\nKeep this exact body.\n"
	if err := os.WriteFile(refinedPath, []byte(refined), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "refine", "use the authored refined artifact"); err != nil {
		t.Fatal(err)
	}
	if _, err := RefineFrom(repo, state.RunID, refinedPath); err != nil {
		t.Fatal(err)
	}
	if got := readRunFile(repo, state.RunID, "outputs", "refined.md"); got != refined {
		t.Fatalf("refined artifact was not preserved:\n%s", got)
	}
	if _, err := RefineFrom(repo, state.RunID, refinedPath); err == nil {
		t.Fatal("expected illegal transition for second refine")
	}

	hardenedPath := filepath.Join(repo, "hardened.md")
	hardened := "# Authored hardened notes\n\nExact challenge pass.\n"
	if err := os.WriteFile(hardenedPath, []byte(hardened), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "harden", "use the authored hardening artifact"); err != nil {
		t.Fatal(err)
	}
	if _, err := HardenFrom(repo, state.RunID, "manual", hardenedPath); err != nil {
		t.Fatal(err)
	}
	if got := readRunFile(repo, state.RunID, "outputs", "hardened.md"); got != hardened {
		t.Fatalf("hardened artifact was not preserved:\n%s", got)
	}

	delta := SpecDelta{
		Schema:        1,
		RunID:         state.RunID,
		SourceSummary: "operator-authored delta",
		Decisions: []runstate.Decision{{
			ID:             "DEC-9001",
			Title:          "Use authored delta",
			Status:         "proposed",
			Recommendation: "accept",
		}},
		AffectedDocs: []string{"docs/interfaces/cli_commands.md"},
		PatchPlan:    []string{"Update CLI docs."},
	}
	raw, err := json.MarshalIndent(delta, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	deltaPath := filepath.Join(repo, "spec_delta.json")
	if err := os.WriteFile(deltaPath, append(raw, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "synthesize", "use the authored spec delta"); err != nil {
		t.Fatal(err)
	}
	if _, err := SynthesizeFrom(repo, state.RunID, deltaPath); err != nil {
		t.Fatal(err)
	}
	if got := readRunFile(repo, state.RunID, "outputs", "spec_delta.json"); got != string(append(raw, '\n')) {
		t.Fatalf("spec delta was not preserved:\n%s", got)
	}
	loaded, err := runstate.Load(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := loaded.Decisions["DEC-9001"]; !ok {
		t.Fatal("expected decision from authored delta")
	}
}

func TestSemanticProductionRequiresStageNotes(t *testing.T) {
	repo, state := newIngestedRun(t)
	if _, err := Intake(repo, state.RunID); err != nil {
		t.Fatal(err)
	}

	if _, err := Refine(repo, state.RunID); err == nil || !strings.Contains(err.Error(), "specops note "+state.RunID+" --stage refine") {
		t.Fatalf("expected refine stage-note error, got %v", err)
	}
	if _, err := Note(repo, state.RunID, "refine", "operator guidance for refine"); err != nil {
		t.Fatal(err)
	}
	if _, err := Refine(repo, state.RunID); err != nil {
		t.Fatal(err)
	}

	if _, err := Harden(repo, state.RunID, "manual"); err == nil || !strings.Contains(err.Error(), "specops note "+state.RunID+" --stage harden") {
		t.Fatalf("expected harden stage-note error, got %v", err)
	}
	if _, err := Note(repo, state.RunID, "harden", "operator guidance for harden"); err != nil {
		t.Fatal(err)
	}
	if _, err := Harden(repo, state.RunID, "manual"); err != nil {
		t.Fatal(err)
	}

	if _, err := Synthesize(repo, state.RunID); err == nil || !strings.Contains(err.Error(), "specops note "+state.RunID+" --stage synthesize") {
		t.Fatalf("expected synthesize stage-note error, got %v", err)
	}
	if _, err := Note(repo, state.RunID, "synthesize", "operator guidance for synthesize"); err != nil {
		t.Fatal(err)
	}
	if _, err := Synthesize(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
}

func TestSemanticProductionRequiresStageNotesWithFromArtifacts(t *testing.T) {
	repo, state := newIngestedRun(t)
	if _, err := Intake(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	refinedPath := filepath.Join(repo, "refined.md")
	if err := os.WriteFile(refinedPath, []byte("# Authored refined notes\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := RefineFrom(repo, state.RunID, refinedPath); err == nil || !strings.Contains(err.Error(), "specops note "+state.RunID+" --stage refine") {
		t.Fatalf("expected refine --from stage-note error, got %v", err)
	}
}

func TestLegacyPromptPathSatisfiesStageNote(t *testing.T) {
	repo, state := newIngestedRun(t)
	if _, err := Intake(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	addLegacyPromptArtifact(t, repo, state.RunID, "refine")
	if _, err := Refine(repo, state.RunID); err != nil {
		t.Fatal(err)
	}

	addLegacyPromptArtifact(t, repo, state.RunID, "harden")
	if _, err := Harden(repo, state.RunID, "manual"); err != nil {
		t.Fatal(err)
	}

	addLegacyPromptArtifact(t, repo, state.RunID, "synthesize")
	if _, err := Synthesize(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
}

func addLegacyPromptArtifact(t *testing.T, repo, runID, stage string) {
	t.Helper()
	state, err := runstate.Load(repo, runID)
	if err != nil {
		t.Fatal(err)
	}
	filename := "20260506-120000-" + stage + ".md"
	path := filepath.Join(runstate.RunDir(repo, state.RunID), "prompts", filename)
	if err := os.WriteFile(path, []byte("# Legacy note\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	state.Artifacts = append(state.Artifacts, runstate.ArtifactRef{
		ID:        "prompt-legacy-" + stage,
		Type:      "prompt",
		Path:      "prompts/" + filename,
		CreatedAt: "2026-05-06T12:00:00Z",
	})
	if err := runstate.Save(repo, state); err != nil {
		t.Fatal(err)
	}
}

func newIngestedRun(t *testing.T) (string, *runstate.RunState) {
	t.Helper()
	repo := t.TempDir()
	state, err := runstate.NewRun(repo, "operator-loop")
	if err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(repo, "source.md")
	if err := os.WriteFile(source, []byte("# Source\n\nRaw operator material."), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := input.IngestFile(input.IngestOptions{Repo: repo, RunID: state.RunID, Path: source}); err != nil {
		t.Fatal(err)
	}
	loaded, err := runstate.Load(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	return repo, loaded
}
