package artifacts

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/specops/specops/internal/runstate"
)

func Compile(repo, runID string, acceptedOnly bool) (PatchPlan, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return PatchPlan{}, err
	}
	if acceptedOnly && state.Status != runstate.StatusDecisionsAccepted {
		return PatchPlan{}, fmt.Errorf("compile --accepted-only requires status %s, got %s", runstate.StatusDecisionsAccepted, state.Status)
	}
	var accepted []runstate.Decision
	for _, decision := range state.Decisions {
		if !acceptedOnly || decision.Status == "accepted" {
			accepted = append(accepted, decision)
		}
	}
	sort.Slice(accepted, func(i, j int) bool { return accepted[i].ID < accepted[j].ID })
	if acceptedOnly && len(accepted) == 0 {
		return PatchPlan{}, fmt.Errorf("no accepted decisions to compile")
	}
	content := refineryNote(repo, state, accepted)
	plan := PatchPlan{
		Schema:       1,
		RunID:        state.RunID,
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
		AcceptedOnly: acceptedOnly,
		Items: []PatchItem{{
			ID:          "patch-001",
			Action:      "create",
			Path:        filepath.ToSlash(filepath.Join("docs", "research", "refinery", state.RunID+".md")),
			Title:       "Promote reviewed run provenance",
			Content:     content,
			DecisionIDs: decisionIDs(accepted),
		}},
	}
	if err := writePatchPlan(repo, state.RunID, plan); err != nil {
		return PatchPlan{}, err
	}
	state.Status = runstate.StatusCompiled
	runstate.AddArtifact(state, "patch_plan", "patches/patch_plan.json")
	if err := runstate.Save(repo, state); err != nil {
		return PatchPlan{}, err
	}
	return plan, nil
}

func LoadPlan(repo, runID string) (PatchPlan, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return PatchPlan{}, err
	}
	path := filepath.Join(runstate.RunDir(repo, state.RunID), "patches", "patch_plan.json")
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

func MarkPlanned(repo, runID string) (*runstate.RunState, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return nil, err
	}
	if state.Status != runstate.StatusCompiled && state.Status != runstate.StatusPlanned {
		return nil, fmt.Errorf("plan requires status %s, got %s", runstate.StatusCompiled, state.Status)
	}
	state.Status = runstate.StatusPlanned
	if err := runstate.Save(repo, state); err != nil {
		return nil, err
	}
	return state, nil
}

type ApplyResult struct {
	RunID    string      `json:"run_id"`
	DryRun   bool        `json:"dry_run"`
	Files    []ApplyFile `json:"files"`
	Warnings []string    `json:"warnings"`
}

type ApplyFile struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

func Apply(repo, runID string, dryRun, commit bool) (ApplyResult, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return ApplyResult{}, err
	}
	if state.Status != runstate.StatusPlanned && state.Status != runstate.StatusCompiled {
		return ApplyResult{}, fmt.Errorf("apply requires status %s, got %s", runstate.StatusPlanned, state.Status)
	}
	plan, err := LoadPlan(repo, state.RunID)
	if err != nil {
		return ApplyResult{}, err
	}
	result := ApplyResult{RunID: state.RunID, DryRun: dryRun, Files: []ApplyFile{}, Warnings: []string{}}
	for _, item := range plan.Items {
		dest, err := safeRepoPath(repo, item.Path)
		if err != nil {
			return ApplyResult{}, err
		}
		status := "would-create"
		if !dryRun {
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				return ApplyResult{}, err
			}
			if existing, err := os.ReadFile(dest); err == nil {
				if string(existing) == item.Content {
					status = "unchanged"
				} else {
					return ApplyResult{}, fmt.Errorf("refusing to overwrite divergent file %s", dest)
				}
			} else if os.IsNotExist(err) {
				if err := os.WriteFile(dest, []byte(item.Content), 0o644); err != nil {
					return ApplyResult{}, err
				}
				status = "created"
			} else {
				return ApplyResult{}, err
			}
		}
		result.Files = append(result.Files, ApplyFile{Path: dest, Status: status})
	}
	if !dryRun {
		state.Status = runstate.StatusApplied
		if err := runstate.Save(repo, state); err != nil {
			return ApplyResult{}, err
		}
	}
	if commit && !dryRun {
		if err := gitCommit(repo, state.RunID); err != nil {
			result.Warnings = append(result.Warnings, err.Error())
		}
	}
	return result, nil
}

func writePatchPlan(repo, runID string, plan PatchPlan) error {
	raw, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(runstate.RunDir(repo, runID), "patches", "patch_plan.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(raw, '\n'), 0o644)
}

func refineryNote(repo string, state *runstate.RunState, decisions []runstate.Decision) string {
	intake := readRunFile(repo, state.RunID, "outputs", "intake.md")
	refined := readRunFile(repo, state.RunID, "outputs", "refined.md")
	hardened := readRunFile(repo, state.RunID, "outputs", "hardened.md")
	var b strings.Builder
	fmt.Fprintf(&b, "---\nid: refinery-%s\ntitle: Reviewed provenance for %s\ndoc_type: research_note\nstatus: accepted\nnormative: false\nversion_scope: v0_required\nlast_reviewed: %s\n---\n\n", state.RunID, state.RunID, time.Now().UTC().Format("2006-01-02"))
	fmt.Fprintf(&b, "# Reviewed Provenance: %s\n\n", state.RunID)
	fmt.Fprintf(&b, "Run `%s` promoted reviewed provenance after accepted decisions.\n\n", state.RunID)
	b.WriteString("## Accepted decisions\n\n")
	for _, decision := range decisions {
		fmt.Fprintf(&b, "- `%s`: %s\n", decision.ID, decision.Title)
	}
	if len(decisions) == 0 {
		b.WriteString("- None recorded.\n")
	}
	b.WriteString("\n## Intake\n\n")
	b.WriteString(trimForSummary(intake))
	b.WriteString("\n\n## Refined notes\n\n")
	b.WriteString(trimForSummary(refined))
	b.WriteString("\n\n## Hardened notes\n\n")
	b.WriteString(trimForSummary(hardened))
	b.WriteByte('\n')
	return b.String()
}

func decisionIDs(decisions []runstate.Decision) []string {
	ids := make([]string, 0, len(decisions))
	for _, decision := range decisions {
		ids = append(ids, decision.ID)
	}
	return ids
}

func safeRepoPath(repo, rel string) (string, error) {
	if filepath.IsAbs(rel) {
		return "", fmt.Errorf("patch path must be relative: %s", rel)
	}
	dest := filepath.Clean(filepath.Join(repo, rel))
	root := filepath.Clean(repo)
	relative, err := filepath.Rel(root, dest)
	if err != nil {
		return "", err
	}
	if relative == ".." || strings.HasPrefix(relative, "../") {
		return "", fmt.Errorf("patch path escapes repo: %s", rel)
	}
	return dest, nil
}

func gitCommit(repo, runID string) error {
	if _, err := os.Stat(filepath.Join(repo, ".git")); err != nil {
		return fmt.Errorf("--commit skipped: target is not a git repository")
	}
	add := exec.Command("git", "add", "docs/research/refinery")
	add.Dir = repo
	if out, err := add.CombinedOutput(); err != nil {
		return fmt.Errorf("git add failed: %s", strings.TrimSpace(string(out)))
	}
	commit := exec.Command("git", "commit", "-m", "Apply SpecOps run "+runID)
	commit.Dir = repo
	if out, err := commit.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit failed: %s", strings.TrimSpace(string(out)))
	}
	return nil
}
