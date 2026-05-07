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

type ProductionResult struct {
	RunID    string               `json:"run_id"`
	Status   runstate.Status      `json:"status"`
	Artifact runstate.ArtifactRef `json:"artifact"`
}

func Intake(repo, runID string) (ProductionResult, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return ProductionResult{}, err
	}
	if state.Status != runstate.StatusIngested {
		return ProductionResult{}, fmt.Errorf("intake requires status %s, got %s", runstate.StatusIngested, state.Status)
	}
	summary := readRunFile(repo, state.RunID, "inputs", "input_summary.md")
	content := fmt.Sprintf("# Intake\n\nRun: `%s`\n\n## Source summary\n\n%s\n\n## Observations\n\n- Source material has been normalized into this run.\n- Human review is required before canonical specification changes.\n", state.RunID, summary)
	ref, err := writeOutput(repo, state, "intake.md", "intake", content)
	if err != nil {
		return ProductionResult{}, err
	}
	state.Status = runstate.StatusIntakeComplete
	if err := runstate.Save(repo, state); err != nil {
		return ProductionResult{}, err
	}
	return ProductionResult{RunID: state.RunID, Status: state.Status, Artifact: ref}, nil
}

func Refine(repo, runID string) (ProductionResult, error) {
	return RefineFrom(repo, runID, "")
}

func RefineFrom(repo, runID, from string) (ProductionResult, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return ProductionResult{}, err
	}
	if state.Status != runstate.StatusIntakeComplete {
		return ProductionResult{}, fmt.Errorf("refine requires status %s, got %s", runstate.StatusIntakeComplete, state.Status)
	}
	if err := requireStageNote(state, "refine"); err != nil {
		return ProductionResult{}, err
	}
	if err := requireAuthoredArtifact(state, "refine", from, "specops refine "+state.RunID+" --from <file>"); err != nil {
		return ProductionResult{}, err
	}
	raw, err := os.ReadFile(from)
	if err != nil {
		return ProductionResult{}, err
	}
	content := string(raw)
	ref, err := writeOutput(repo, state, "refined.md", "refined", content)
	if err != nil {
		return ProductionResult{}, err
	}
	state.Status = runstate.StatusRefined
	if err := runstate.Save(repo, state); err != nil {
		return ProductionResult{}, err
	}
	return ProductionResult{RunID: state.RunID, Status: state.Status, Artifact: ref}, nil
}

func Harden(repo, runID, backend string) (ProductionResult, error) {
	return HardenFrom(repo, runID, backend, "")
}

func HardenFrom(repo, runID, backend, from string) (ProductionResult, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return ProductionResult{}, err
	}
	if state.Status != runstate.StatusRefined {
		return ProductionResult{}, fmt.Errorf("harden requires status %s, got %s", runstate.StatusRefined, state.Status)
	}
	if err := requireStageNote(state, "harden"); err != nil {
		return ProductionResult{}, err
	}
	if backend == "" {
		backend = "manual"
	}
	if err := requireAuthoredArtifact(state, "harden", from, "specops harden "+state.RunID+" --from <file>"); err != nil {
		return ProductionResult{}, err
	}
	raw, err := os.ReadFile(from)
	if err != nil {
		return ProductionResult{}, err
	}
	content := string(raw)
	ref, err := writeOutput(repo, state, "hardened.md", "hardened", content)
	if err != nil {
		return ProductionResult{}, err
	}
	state.Status = runstate.StatusHardened
	if err := runstate.Save(repo, state); err != nil {
		return ProductionResult{}, err
	}
	return ProductionResult{RunID: state.RunID, Status: state.Status, Artifact: ref}, nil
}

func Synthesize(repo, runID string) (ProductionResult, error) {
	return SynthesizeFrom(repo, runID, "")
}

func SynthesizeFrom(repo, runID, from string) (ProductionResult, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return ProductionResult{}, err
	}
	if state.Status != runstate.StatusHardened && state.Status != runstate.StatusRefined {
		return ProductionResult{}, fmt.Errorf("synthesize requires status %s or %s, got %s", runstate.StatusHardened, runstate.StatusRefined, state.Status)
	}
	if err := requireStageNote(state, "synthesize"); err != nil {
		return ProductionResult{}, err
	}
	if err := requireAuthoredArtifact(state, "synthesize", from, "specops synthesize "+state.RunID+" --from <spec_delta.json>"); err != nil {
		return ProductionResult{}, err
	}
	raw, err := os.ReadFile(from)
	if err != nil {
		return ProductionResult{}, err
	}
	var delta SpecDelta
	if err := json.Unmarshal(raw, &delta); err != nil {
		return ProductionResult{}, err
	}
	if delta.RunID != "" && delta.RunID != state.RunID {
		return ProductionResult{}, fmt.Errorf("spec delta run_id %q does not match %q", delta.RunID, state.RunID)
	}
	path := filepath.Join(runstate.RunDir(repo, state.RunID), "outputs", "spec_delta.json")
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return ProductionResult{}, err
	}
	ref := runstate.ArtifactRef{ID: fmt.Sprintf("spec_delta-%03d", len(state.Artifacts)+1), Type: "spec_delta", Path: "outputs/spec_delta.json", CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	state.Artifacts = append(state.Artifacts, ref)
	for _, decision := range delta.Decisions {
		state.Decisions[decision.ID] = decision
	}
	state.Status = runstate.StatusAwaitingDecisions
	if err := runstate.Save(repo, state); err != nil {
		return ProductionResult{}, err
	}
	return ProductionResult{RunID: state.RunID, Status: state.Status, Artifact: ref}, nil
}

func requireStageNote(state *runstate.RunState, stage string) error {
	if hasStageNote(state, stage) {
		return nil
	}
	return fmt.Errorf("semantic gate %q requires a recorded stage note before execution\nrefresh context: specops context %s\nrecord note: specops note %s --stage %s --text <file-or-inline>", stage, state.RunID, state.RunID, stage)
}

func requireAuthoredArtifact(state *runstate.RunState, stage, from, command string) error {
	if strings.TrimSpace(from) != "" {
		return nil
	}
	return fmt.Errorf("semantic gate %q requires an authored artifact via --from because the CLI is not AI-enabled and cannot generate content-aware %s output\nrefresh context: specops context %s\nrecord note: specops note %s --stage %s --text <file-or-inline>\nrun command: %s", stage, stage, state.RunID, state.RunID, stage, command)
}

func Deepen(repo, runID, target string) (ProductionResult, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return ProductionResult{}, err
	}
	if target == "" {
		return ProductionResult{}, fmt.Errorf("--target is required")
	}
	name := safeName(target)
	content := fmt.Sprintf("# Deepen: %s\n\nRun: `%s`\n\n## Focus\n\n%s\n\n## Notes\n\nAdd source-grounded details here before synthesizing or amending decisions.\n", target, state.RunID, target)
	ref, err := writeOutput(repo, state, "deepen-"+name+".md", "deepen", content)
	if err != nil {
		return ProductionResult{}, err
	}
	if err := runstate.Save(repo, state); err != nil {
		return ProductionResult{}, err
	}
	return ProductionResult{RunID: state.RunID, Status: state.Status, Artifact: ref}, nil
}

func writeOutput(repo string, state *runstate.RunState, filename, typ, content string) (runstate.ArtifactRef, error) {
	path := filepath.Join(runstate.RunDir(repo, state.RunID), "outputs", filename)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return runstate.ArtifactRef{}, err
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return runstate.ArtifactRef{}, err
	}
	ref := runstate.ArtifactRef{ID: fmt.Sprintf("%s-%03d", typ, len(state.Artifacts)+1), Type: typ, Path: "outputs/" + filename, CreatedAt: time.Now().UTC().Format(time.RFC3339)}
	state.Artifacts = append(state.Artifacts, ref)
	return ref, nil
}

func readRunFile(repo, runID string, parts ...string) string {
	path := filepath.Join(append([]string{runstate.RunDir(repo, runID)}, parts...)...)
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(raw)
}

func trimForSummary(text string) string {
	text = strings.TrimSpace(text)
	if len(text) > 2000 {
		return text[:2000] + "\n\n[truncated]"
	}
	return text
}

func safeName(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else if b.Len() > 0 {
			b.WriteByte('-')
		}
	}
	return strings.Trim(b.String(), "-")
}
