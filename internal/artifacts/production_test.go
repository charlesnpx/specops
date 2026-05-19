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
	refinedPath := writeTextArtifact(t, repo, "refined.md", "# Authored refined notes\n\nReviewable semantic content.\n")
	if _, err := RefineFrom(repo, state.RunID, refinedPath); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "synthesize", "prepare explicit decisions"); err != nil {
		t.Fatal(err)
	}
	deltaPath := writeDeltaArtifact(t, repo, state.RunID)
	if _, err := SynthesizeFrom(repo, state.RunID, deltaPath); err != nil {
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
	if _, err := Refine(repo, state.RunID); err == nil || !strings.Contains(err.Error(), "requires an authored artifact via --from") {
		t.Fatalf("expected refine authored-artifact error, got %v", err)
	}
	refinedPath := writeTextArtifact(t, repo, "refined.md", "# Authored refined notes\n")
	if _, err := RefineFrom(repo, state.RunID, refinedPath); err != nil {
		t.Fatal(err)
	}

	if _, err := Harden(repo, state.RunID, "manual"); err == nil || !strings.Contains(err.Error(), "specops note "+state.RunID+" --stage harden") {
		t.Fatalf("expected harden stage-note error, got %v", err)
	}
	if _, err := Note(repo, state.RunID, "harden", "operator guidance for harden"); err != nil {
		t.Fatal(err)
	}
	if _, err := Harden(repo, state.RunID, "manual"); err == nil || !strings.Contains(err.Error(), "requires an authored artifact via --from") {
		t.Fatalf("expected harden authored-artifact error, got %v", err)
	}
	hardenedPath := writeTextArtifact(t, repo, "hardened.md", "# Authored hardened notes\n")
	if _, err := HardenFrom(repo, state.RunID, "manual", hardenedPath); err != nil {
		t.Fatal(err)
	}

	if _, err := Synthesize(repo, state.RunID); err == nil || !strings.Contains(err.Error(), "specops note "+state.RunID+" --stage synthesize") {
		t.Fatalf("expected synthesize stage-note error, got %v", err)
	}
	if _, err := Note(repo, state.RunID, "synthesize", "operator guidance for synthesize"); err != nil {
		t.Fatal(err)
	}
	if _, err := Synthesize(repo, state.RunID); err == nil || !strings.Contains(err.Error(), "requires an authored artifact via --from") {
		t.Fatalf("expected synthesize authored-artifact error, got %v", err)
	}
	deltaPath := writeDeltaArtifact(t, repo, state.RunID)
	if _, err := SynthesizeFrom(repo, state.RunID, deltaPath); err != nil {
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
	refinedPath := writeTextArtifact(t, repo, "refined.md", "# Authored refined notes\n")
	if _, err := RefineFrom(repo, state.RunID, refinedPath); err != nil {
		t.Fatal(err)
	}

	addLegacyPromptArtifact(t, repo, state.RunID, "harden")
	hardenedPath := writeTextArtifact(t, repo, "hardened.md", "# Authored hardened notes\n")
	if _, err := HardenFrom(repo, state.RunID, "manual", hardenedPath); err != nil {
		t.Fatal(err)
	}

	addLegacyPromptArtifact(t, repo, state.RunID, "synthesize")
	deltaPath := writeDeltaArtifact(t, repo, state.RunID)
	if _, err := SynthesizeFrom(repo, state.RunID, deltaPath); err != nil {
		t.Fatal(err)
	}
}

func TestCompileIncludesAcceptedSpecDeltaCanonicalDocs(t *testing.T) {
	repo, state := newIngestedRun(t)
	if _, err := Intake(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "refine", "operator guidance for refine"); err != nil {
		t.Fatal(err)
	}
	refinedPath := writeTextArtifact(t, repo, "refined.md", "# Authored refined notes\n")
	if _, err := RefineFrom(repo, state.RunID, refinedPath); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "synthesize", "operator guidance for synthesize"); err != nil {
		t.Fatal(err)
	}

	exactCanon := "# Exact Canon\n\nAuthor supplied canonical text.\n"
	delta := SpecDelta{
		Schema:        1,
		RunID:         state.RunID,
		SourceSummary: "operator-authored delta",
		NewConcepts:   []string{"governed context graph kernel"},
		Decisions: []runstate.Decision{{
			ID:             "D001",
			Title:          "Create canonical kernel docs",
			Status:         "proposed",
			Recommendation: "accept",
			AffectedDocs:   []string{"docs/CANON.md", "docs/versions/v0_scope.md"},
		}},
		AffectedDocs: []string{"docs/CANON.md", "docs/versions/v0_scope.md"},
		PatchPlan:    []string{"Create canonical frame and v0 scope."},
		PatchItems: []PatchItem{{
			Path:        "docs/CANON.md",
			Action:      "create",
			Title:       "Create exact canon",
			Content:     exactCanon,
			DecisionIDs: []string{"D001"},
		}},
	}
	deltaPath := writeSpecDeltaArtifact(t, repo, "spec_delta.json", delta)
	if _, err := SynthesizeFrom(repo, state.RunID, deltaPath); err != nil {
		t.Fatal(err)
	}
	if _, err := SetDecision(repo, state.RunID, "D001", "accepted", ""); err != nil {
		t.Fatal(err)
	}
	plan, err := Compile(repo, state.RunID, true)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Health.Stale || plan.Health.Incomplete {
		t.Fatalf("fresh plan should be healthy: %+v", plan.Health)
	}
	if item, ok := findPatchItem(plan, "docs/CANON.md"); !ok || item.Content != exactCanon {
		t.Fatalf("expected exact authored CANON patch, got found=%v item=%+v", ok, item)
	}
	if _, ok := findPatchItem(plan, "docs/versions/v0_scope.md"); !ok {
		t.Fatalf("expected generated v0 scope patch in %+v", plan.Items)
	}
	if _, ok := findPatchItem(plan, "docs/research/refinery/"+state.RunID+".md"); !ok {
		t.Fatalf("expected provenance patch in %+v", plan.Items)
	}
	if _, err := MarkPlanned(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	recompiled, err := Compile(repo, state.RunID, true)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := findPatchItem(recompiled, "docs/versions/v0_scope.md"); !ok {
		t.Fatalf("expected replanned v0 scope patch in %+v", recompiled.Items)
	}
}

func TestPatchPlanHealthSeparatesStaleAndIncomplete(t *testing.T) {
	repo, state, plan := newCompiledPlan(t)

	plan.Items = filterPatchItems(plan.Items, "docs/versions/v0_scope.md")
	if err := writePatchPlan(repo, state.RunID, plan); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadPlan(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Health.Stale {
		t.Fatalf("missing patch item should not make unchanged inputs stale: %+v", loaded.Health)
	}
	if !loaded.Health.Incomplete || !strings.Contains(strings.Join(loaded.Health.IncompleteReasons, "\n"), "docs/versions/v0_scope.md") {
		t.Fatalf("expected incomplete plan for missing v0 scope: %+v", loaded.Health)
	}
	if _, err := Apply(repo, state.RunID, false, false, false); err == nil || !strings.Contains(err.Error(), "refusing to apply unsafe patch plan") {
		t.Fatalf("expected unsafe apply refusal, got %v", err)
	}
	if _, err := Apply(repo, state.RunID, true, false, false); err != nil {
		t.Fatalf("dry-run should be allowed for unsafe plan inspection: %v", err)
	}

	recompiled, err := Compile(repo, state.RunID, true)
	if err != nil {
		t.Fatal(err)
	}
	delta, err := loadSpecDelta(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	delta.AffectedDocs = append(delta.AffectedDocs, "docs/interfaces/cli_commands.md")
	if err := os.WriteFile(filepath.Join(runstate.RunDir(repo, state.RunID), "outputs", "spec_delta.json"), mustJSON(t, delta), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writePatchPlan(repo, state.RunID, recompiled); err != nil {
		t.Fatal(err)
	}
	stale, err := LoadPlan(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if !stale.Health.Stale || !strings.Contains(strings.Join(stale.Health.StaleReasons, "\n"), "spec_delta hash changed") {
		t.Fatalf("expected stale plan after spec delta changed: %+v", stale.Health)
	}
}

func TestSynthesizeRejectsPlannedAndSupersedeRequiresApplyNote(t *testing.T) {
	repo, state, _ := newCompiledPlan(t)
	if _, err := MarkPlanned(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	replacementPath := writeSpecDeltaArtifact(t, repo, "replacement_spec_delta.json", replacementSpecDelta(state.RunID, "# Exact v0 Scope\n"))

	if _, err := SynthesizeFrom(repo, state.RunID, replacementPath); err == nil || !strings.Contains(err.Error(), "synthesize requires status") {
		t.Fatalf("expected synthesize to reject planned status, got %v", err)
	}
	if _, err := SupersedeSynthesisFrom(repo, state.RunID, replacementPath, false); err == nil || !strings.Contains(err.Error(), "specops note "+state.RunID+" --stage apply") {
		t.Fatalf("expected apply note requirement, got %v", err)
	}
}

func TestSupersedeSynthesisArchivesCurrentArtifactsAndPreservesSettledDecisions(t *testing.T) {
	repo, state, _ := newCompiledPlan(t)
	if _, err := MarkPlanned(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "apply", "compiled docs are too thin; supersede synthesis before apply"); err != nil {
		t.Fatal(err)
	}
	exactScope := "# Exact v0 Scope\n\nFull authored model.\n"
	replacementPath := writeSpecDeltaArtifact(t, repo, "replacement_spec_delta.json", replacementSpecDelta(state.RunID, exactScope))

	result, err := SupersedeSynthesisFrom(repo, state.RunID, replacementPath, false)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != runstate.StatusDecisionsAccepted || result.ReopenedDecisions {
		t.Fatalf("unexpected result: %+v", result)
	}
	if len(result.ArchivedArtifacts) < 2 {
		t.Fatalf("expected spec delta and patch plan archives, got %+v", result.ArchivedArtifacts)
	}
	if result.Artifact.Type != "spec_delta" || result.Artifact.Path != "outputs/spec_delta.json" {
		t.Fatalf("unexpected current artifact: %+v", result.Artifact)
	}
	if _, err := os.Stat(filepath.Join(runstate.RunDir(repo, state.RunID), "patches", "patch_plan.json")); !os.IsNotExist(err) {
		t.Fatalf("current patch plan should be removed after supersession, got %v", err)
	}
	for _, ref := range result.ArchivedArtifacts {
		if _, err := os.Stat(filepath.Join(runstate.RunDir(repo, state.RunID), filepath.FromSlash(ref.Path))); err != nil {
			t.Fatalf("archived artifact missing: %+v: %v", ref, err)
		}
	}
	loaded, err := runstate.Load(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Decisions["D001"].Status != "accepted" {
		t.Fatalf("settled decision status changed: %+v", loaded.Decisions["D001"])
	}
	if !hasArtifactType(loaded, "superseded_spec_delta") || !hasArtifactType(loaded, "superseded_patch_plan") {
		t.Fatalf("expected superseded artifact refs: %+v", loaded.Artifacts)
	}

	plan, err := Compile(repo, state.RunID, true)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Health.Stale || plan.Health.Incomplete {
		t.Fatalf("replacement plan should be healthy: %+v", plan.Health)
	}
	if item, ok := findPatchItem(plan, "docs/versions/v0_scope.md"); !ok || item.Content != exactScope {
		t.Fatalf("expected authored replacement scope patch, found=%v item=%+v", ok, item)
	}
}

func TestSupersedeSynthesisRejectsDecisionChangesUnlessReopened(t *testing.T) {
	repo, state, _ := newCompiledPlan(t)
	if _, err := MarkPlanned(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "apply", "review rejected patch quality"); err != nil {
		t.Fatal(err)
	}
	changed := replacementSpecDelta(state.RunID, "# Exact v0 Scope\n")
	changed.Decisions[0].Title = "Changed settled decision"
	changedPath := writeSpecDeltaArtifact(t, repo, "changed_spec_delta.json", changed)
	if _, err := SupersedeSynthesisFrom(repo, state.RunID, changedPath, false); err == nil || !strings.Contains(err.Error(), "changes settled decision") {
		t.Fatalf("expected changed decision rejection, got %v", err)
	}

	newDecision := replacementSpecDelta(state.RunID, "# Exact v0 Scope\n")
	newDecision.Decisions = append(newDecision.Decisions, runstate.Decision{ID: "D999", Title: "New decision", Status: "proposed"})
	newDecisionPath := writeSpecDeltaArtifact(t, repo, "new_decision_spec_delta.json", newDecision)
	if _, err := SupersedeSynthesisFrom(repo, state.RunID, newDecisionPath, false); err == nil || !strings.Contains(err.Error(), "introduces decision") {
		t.Fatalf("expected new decision rejection, got %v", err)
	}
}

func TestSupersedeSynthesisCanReopenDecisionGate(t *testing.T) {
	repo, state, _ := newCompiledPlan(t)
	if _, err := MarkPlanned(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "apply", "replacement delta needs a new decision"); err != nil {
		t.Fatal(err)
	}
	reopened := replacementSpecDelta(state.RunID, "# Exact v0 Scope\n")
	reopened.Decisions = append(reopened.Decisions, runstate.Decision{ID: "D999", Title: "Add second canonical doc", Status: "proposed"})
	reopenedPath := writeSpecDeltaArtifact(t, repo, "reopened_spec_delta.json", reopened)

	result, err := SupersedeSynthesisFrom(repo, state.RunID, reopenedPath, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.Status != runstate.StatusAwaitingDecisions || !result.ReopenedDecisions {
		t.Fatalf("expected reopened decision gate, got %+v", result)
	}
	loaded, err := runstate.Load(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := loaded.Decisions["D999"]; !ok {
		t.Fatalf("expected reopened new decision in state: %+v", loaded.Decisions)
	}
}

func newCompiledPlan(t *testing.T) (string, *runstate.RunState, PatchPlan) {
	t.Helper()
	repo, state := newIngestedRun(t)
	if _, err := Intake(repo, state.RunID); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "refine", "operator guidance for refine"); err != nil {
		t.Fatal(err)
	}
	refinedPath := writeTextArtifact(t, repo, "refined.md", "# Authored refined notes\n")
	if _, err := RefineFrom(repo, state.RunID, refinedPath); err != nil {
		t.Fatal(err)
	}
	if _, err := Note(repo, state.RunID, "synthesize", "operator guidance for synthesize"); err != nil {
		t.Fatal(err)
	}
	delta := SpecDelta{
		Schema:        1,
		RunID:         state.RunID,
		SourceSummary: "operator-authored delta",
		Decisions: []runstate.Decision{{
			ID:             "D001",
			Title:          "Create canonical kernel docs",
			Status:         "proposed",
			Recommendation: "accept",
			AffectedDocs:   []string{"docs/CANON.md", "docs/versions/v0_scope.md"},
		}},
		AffectedDocs: []string{"docs/CANON.md", "docs/versions/v0_scope.md"},
		PatchPlan:    []string{"Create canonical frame and v0 scope."},
	}
	deltaPath := writeSpecDeltaArtifact(t, repo, "spec_delta.json", delta)
	if _, err := SynthesizeFrom(repo, state.RunID, deltaPath); err != nil {
		t.Fatal(err)
	}
	if _, err := SetDecision(repo, state.RunID, "D001", "accepted", ""); err != nil {
		t.Fatal(err)
	}
	plan, err := Compile(repo, state.RunID, true)
	if err != nil {
		t.Fatal(err)
	}
	loaded, err := runstate.Load(repo, state.RunID)
	if err != nil {
		t.Fatal(err)
	}
	return repo, loaded, plan
}

func replacementSpecDelta(runID, exactScope string) SpecDelta {
	return SpecDelta{
		Schema:        1,
		RunID:         runID,
		SourceSummary: "replacement authored delta",
		Decisions: []runstate.Decision{{
			ID:             "D001",
			Title:          "Create canonical kernel docs",
			Status:         "proposed",
			Recommendation: "accept",
			AffectedDocs:   []string{"docs/CANON.md", "docs/versions/v0_scope.md"},
		}},
		AffectedDocs: []string{"docs/CANON.md", "docs/versions/v0_scope.md"},
		PatchPlan:    []string{"Supersede thin generated docs with authored canonical content."},
		PatchItems: []PatchItem{{
			Path:        "docs/versions/v0_scope.md",
			Action:      "create",
			Title:       "Create exact v0 scope",
			Content:     exactScope,
			DecisionIDs: []string{"D001"},
		}},
	}
}

func hasArtifactType(state *runstate.RunState, typ string) bool {
	for _, artifact := range state.Artifacts {
		if artifact.Type == typ {
			return true
		}
	}
	return false
}

func filterPatchItems(items []PatchItem, removePath string) []PatchItem {
	var out []PatchItem
	for _, item := range items {
		if item.Path != removePath {
			out = append(out, item)
		}
	}
	return out
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	raw, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	return append(raw, '\n')
}

func writeTextArtifact(t *testing.T, repo, name, content string) string {
	t.Helper()
	path := filepath.Join(repo, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeDeltaArtifact(t *testing.T, repo, runID string) string {
	t.Helper()
	delta := SpecDelta{
		Schema:        1,
		RunID:         runID,
		SourceSummary: "operator-authored delta",
		Decisions: []runstate.Decision{{
			ID:             "DEC-0001",
			Title:          "Accept authored delta",
			Status:         "proposed",
			Recommendation: "accept",
		}},
		AffectedDocs: []string{"docs/research/refinery/"},
		PatchPlan:    []string{"Promote reviewed provenance."},
	}
	raw, err := json.MarshalIndent(delta, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(repo, "spec_delta.json")
	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func writeSpecDeltaArtifact(t *testing.T, repo, name string, delta SpecDelta) string {
	t.Helper()
	raw, err := json.MarshalIndent(delta, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(repo, name)
	if err := os.WriteFile(path, append(raw, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func findPatchItem(plan PatchPlan, path string) (PatchItem, bool) {
	for _, item := range plan.Items {
		if item.Path == path {
			return item, true
		}
	}
	return PatchItem{}, false
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
