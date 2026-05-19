package artifacts

import (
	"crypto/sha256"
	"encoding/hex"
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

const compilerContractVersion = "compile-spec-delta-v2"

func Compile(repo, runID string, acceptedOnly bool) (PatchPlan, error) {
	state, err := runstate.Load(repo, runID)
	if err != nil {
		return PatchPlan{}, err
	}
	if acceptedOnly && !canCompileAcceptedOnly(state.Status) {
		return PatchPlan{}, fmt.Errorf("compile --accepted-only requires status %s, %s, or %s, got %s", runstate.StatusDecisionsAccepted, runstate.StatusCompiled, runstate.StatusPlanned, state.Status)
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
	delta, _ := loadSpecDelta(repo, state.RunID)
	items := compileDeltaPatchItems(delta, acceptedOnly, accepted)
	items = append(items, PatchItem{
		Action:      "create",
		Path:        filepath.ToSlash(filepath.Join("docs", "research", "refinery", state.RunID+".md")),
		Title:       "Promote reviewed run provenance",
		Content:     refineryNote(repo, state, accepted),
		DecisionIDs: decisionIDs(accepted),
	})
	items = normalizePatchItems(items)
	inputs := patchPlanInputs(repo, state, accepted)
	plan := PatchPlan{
		Schema:           1,
		RunID:            state.RunID,
		CreatedAt:        time.Now().UTC().Format(time.RFC3339),
		CompilerContract: compilerContractVersion,
		AcceptedOnly:     acceptedOnly,
		Inputs:           inputs,
		Items:            items,
	}
	plan.Health = evaluatePatchPlanHealth(repo, state, plan, accepted)
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

func canCompileAcceptedOnly(status runstate.Status) bool {
	return status == runstate.StatusDecisionsAccepted || status == runstate.StatusCompiled || status == runstate.StatusPlanned
}

func acceptedDecisions(decisions map[string]runstate.Decision) []runstate.Decision {
	accepted := make([]runstate.Decision, 0, len(decisions))
	for _, decision := range decisions {
		if decision.Status == "accepted" {
			accepted = append(accepted, decision)
		}
	}
	sort.Slice(accepted, func(i, j int) bool { return accepted[i].ID < accepted[j].ID })
	return accepted
}

func loadSpecDelta(repo, runID string) (SpecDelta, error) {
	path := filepath.Join(runstate.RunDir(repo, runID), "outputs", "spec_delta.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return SpecDelta{}, err
	}
	var delta SpecDelta
	if err := json.Unmarshal(raw, &delta); err != nil {
		return SpecDelta{}, err
	}
	return delta, nil
}

func patchPlanInputs(repo string, state *runstate.RunState, accepted []runstate.Decision) PatchPlanInputs {
	return PatchPlanInputs{
		SpecDeltaSHA256:         hashRunFile(repo, state.RunID, "outputs", "spec_delta.json"),
		AcceptedDecisionsSHA256: hashJSON(accepted),
		ProvenanceSHA256:        hashStrings([]string{readRunFile(repo, state.RunID, "outputs", "intake.md"), readRunFile(repo, state.RunID, "outputs", "refined.md"), readRunFile(repo, state.RunID, "outputs", "hardened.md")}),
	}
}

func evaluatePatchPlanHealth(repo string, state *runstate.RunState, plan PatchPlan, accepted []runstate.Decision) PatchPlanHealth {
	currentInputs := patchPlanInputs(repo, state, accepted)
	health := PatchPlanHealth{}
	if plan.CompilerContract == "" {
		health.Stale = true
		health.StaleReasons = append(health.StaleReasons, "missing compiler contract")
	} else if plan.CompilerContract != compilerContractVersion {
		health.Stale = true
		health.StaleReasons = append(health.StaleReasons, "compiled with older compiler contract")
	}
	if plan.Inputs.SpecDeltaSHA256 != currentInputs.SpecDeltaSHA256 {
		health.Stale = true
		health.StaleReasons = append(health.StaleReasons, "spec_delta hash changed")
	}
	if plan.Inputs.AcceptedDecisionsSHA256 != currentInputs.AcceptedDecisionsSHA256 {
		health.Stale = true
		health.StaleReasons = append(health.StaleReasons, "accepted decisions changed")
	}
	if plan.Inputs.ProvenanceSHA256 != currentInputs.ProvenanceSHA256 {
		health.Stale = true
		health.StaleReasons = append(health.StaleReasons, "provenance inputs changed")
	}
	for _, missing := range missingAcceptedDeltaPaths(repo, state.RunID, plan, accepted) {
		health.Incomplete = true
		health.IncompleteReasons = append(health.IncompleteReasons, "missing accepted patch path: "+missing)
	}
	return health
}

func missingAcceptedDeltaPaths(repo, runID string, plan PatchPlan, accepted []runstate.Decision) []string {
	delta, err := loadSpecDelta(repo, runID)
	if err != nil || delta.Schema == 0 {
		return nil
	}
	expected := map[string]bool{}
	for _, item := range compileDeltaPatchItems(delta, true, accepted) {
		if item.Path != "" {
			expected[filepath.ToSlash(item.Path)] = true
		}
	}
	present := map[string]bool{}
	for _, item := range plan.Items {
		present[filepath.ToSlash(item.Path)] = true
	}
	var missing []string
	for path := range expected {
		if !present[path] {
			missing = append(missing, path)
		}
	}
	sort.Strings(missing)
	return missing
}

func hashRunFile(repo, runID string, parts ...string) string {
	path := filepath.Join(append([]string{runstate.RunDir(repo, runID)}, parts...)...)
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return hashBytes(raw)
}

func hashJSON(value any) string {
	raw, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return hashBytes(raw)
}

func hashStrings(values []string) string {
	return hashBytes([]byte(strings.Join(values, "\x00")))
}

func hashBytes(raw []byte) string {
	sum := sha256.Sum256(raw)
	return hex.EncodeToString(sum[:])
}

func compileDeltaPatchItems(delta SpecDelta, acceptedOnly bool, accepted []runstate.Decision) []PatchItem {
	if delta.Schema == 0 {
		return nil
	}
	acceptedIDs := map[string]bool{}
	for _, decision := range accepted {
		acceptedIDs[decision.ID] = true
	}
	var items []PatchItem
	for _, item := range delta.PatchItems {
		if !acceptedOnly || patchItemAccepted(item, acceptedIDs) {
			items = append(items, item)
		}
	}
	seenPath := map[string]bool{}
	for _, item := range items {
		seenPath[filepath.ToSlash(item.Path)] = true
	}
	for _, path := range affectedDocPaths(delta, accepted) {
		if seenPath[path] || !isCanonicalDocPath(path) {
			continue
		}
		items = append(items, generatedDocPatchItem(delta, accepted, path))
		seenPath[path] = true
	}
	return items
}

func patchItemAccepted(item PatchItem, acceptedIDs map[string]bool) bool {
	if len(item.DecisionIDs) == 0 {
		return len(acceptedIDs) > 0
	}
	for _, id := range item.DecisionIDs {
		if acceptedIDs[id] {
			return true
		}
	}
	return false
}

func affectedDocPaths(delta SpecDelta, accepted []runstate.Decision) []string {
	var paths []string
	paths = append(paths, delta.AffectedDocs...)
	for _, decision := range accepted {
		paths = append(paths, decision.AffectedDocs...)
	}
	sort.Strings(paths)
	var out []string
	seen := map[string]bool{}
	for _, path := range paths {
		path = filepath.ToSlash(strings.TrimSpace(path))
		if path == "" || seen[path] {
			continue
		}
		seen[path] = true
		out = append(out, path)
	}
	return out
}

func isCanonicalDocPath(path string) bool {
	return strings.HasPrefix(path, "docs/") && strings.HasSuffix(path, ".md") && !strings.HasPrefix(path, "docs/research/refinery/")
}

func generatedDocPatchItem(delta SpecDelta, accepted []runstate.Decision, path string) PatchItem {
	title, docType, id := docMetadata(path)
	return PatchItem{
		Action:      "create",
		Path:        path,
		Title:       "Create " + title,
		Content:     generatedDocContent(delta, accepted, path, title, docType, id),
		DecisionIDs: decisionIDs(accepted),
	}
}

func docMetadata(path string) (title, docType, id string) {
	base := strings.TrimSuffix(filepath.Base(path), ".md")
	id = safeName(base)
	title = strings.TrimSpace(strings.ReplaceAll(base, "_", " "))
	if title == "" {
		title = "Specification Document"
	} else {
		title = strings.Title(title)
	}
	docType = "spec"
	switch {
	case path == "docs/CANON.md":
		return "SpecOps Canon", "canon", "canon"
	case strings.HasPrefix(path, "docs/versions/"):
		return title, "version_scope", id
	case strings.HasPrefix(path, "docs/decisions/"):
		return title, "adr", id
	case strings.HasPrefix(path, "docs/interfaces/"):
		return title, "interface_spec", id
	case strings.HasPrefix(path, "docs/cross_cutting/"):
		return title, "cross_cutting_spec", id
	case strings.HasPrefix(path, "docs/subsystems/"):
		return title, "subsystem_spec", id
	}
	return title, docType, id
}

func generatedDocContent(delta SpecDelta, accepted []runstate.Decision, path, title, docType, id string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "---\nid: %s\ntitle: %s\ndoc_type: %s\nstatus: accepted\nnormative: true\nversion_scope: v0_required\nlast_reviewed: %s\n---\n\n", id, title, docType, time.Now().UTC().Format("2006-01-02"))
	fmt.Fprintf(&b, "# %s\n\n", title)
	fmt.Fprintf(&b, "Generated from accepted SpecOps run `%s` for `%s`.\n", delta.RunID, path)
	writeStringSection(&b, "Source Summary", delta.SourceSummary)
	writeListSection(&b, "Accepted Decisions", decisionLines(accepted))
	writeListSection(&b, "New Concepts", delta.NewConcepts)
	writeListSection(&b, "Requirements", delta.Requirements)
	writeListSection(&b, "Constraints", delta.Constraints)
	writeListSection(&b, "Recommendations", delta.Recommendations)
	writeListSection(&b, "Acceptance Criteria", delta.AcceptanceCriteria)
	writeListSection(&b, "Risks", delta.Risks)
	writeListSection(&b, "Open Questions", delta.OpenQuestions)
	writeListSection(&b, "Patch Plan Notes", delta.PatchPlan)
	return b.String()
}

func writeStringSection(b *strings.Builder, title, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}
	fmt.Fprintf(b, "\n## %s\n\n%s\n", title, value)
}

func writeListSection(b *strings.Builder, title string, values []string) {
	if len(values) == 0 {
		return
	}
	fmt.Fprintf(b, "\n## %s\n\n", title)
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			fmt.Fprintf(b, "- %s\n", value)
		}
	}
}

func decisionLines(decisions []runstate.Decision) []string {
	lines := make([]string, 0, len(decisions))
	for _, decision := range decisions {
		lines = append(lines, fmt.Sprintf("`%s`: %s", decision.ID, decision.Title))
	}
	return lines
}

func normalizePatchItems(items []PatchItem) []PatchItem {
	out := make([]PatchItem, 0, len(items))
	seenPath := map[string]bool{}
	for _, item := range items {
		item.Path = filepath.ToSlash(strings.TrimSpace(item.Path))
		if item.Path == "" || seenPath[item.Path] {
			continue
		}
		if item.Action == "" {
			item.Action = "create"
		}
		if item.Title == "" {
			item.Title = "Patch " + item.Path
		}
		seenPath[item.Path] = true
		out = append(out, item)
	}
	for i := range out {
		if out[i].ID == "" {
			out[i].ID = fmt.Sprintf("patch-%03d", i+1)
		}
	}
	return out
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
	accepted := acceptedDecisions(state.Decisions)
	plan.Health = evaluatePatchPlanHealth(repo, state, plan, accepted)
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

func Apply(repo, runID string, dryRun, commit, allowUnsafePlan bool) (ApplyResult, error) {
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
	if !dryRun && (plan.Health.Stale || plan.Health.Incomplete) && !allowUnsafePlan {
		return ApplyResult{}, fmt.Errorf("refusing to apply unsafe patch plan: stale=%t incomplete=%t (%s); rerun specops compile %s --accepted-only or pass --allow-unsafe-plan", plan.Health.Stale, plan.Health.Incomplete, strings.Join(append(plan.Health.StaleReasons, plan.Health.IncompleteReasons...), "; "), state.RunID)
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
