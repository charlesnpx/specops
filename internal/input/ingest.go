package input

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/specops/specops/internal/runstate"
)

type IngestOptions struct {
	Repo  string
	RunID string
	Path  string
	Type  string
	Slice bool
}

type IngestResult struct {
	RunID     string                 `json:"run_id"`
	Status    runstate.Status        `json:"status"`
	Input     runstate.InputRef      `json:"input"`
	Artifacts []runstate.ArtifactRef `json:"artifacts"`
}

func IngestFile(opts IngestOptions) (IngestResult, error) {
	if opts.Type == "" {
		opts.Type = "raw_markdown"
	}
	state, err := runstate.Load(opts.Repo, opts.RunID)
	if err != nil {
		return IngestResult{}, err
	}
	source, err := filepath.Abs(opts.Path)
	if err != nil {
		return IngestResult{}, err
	}
	raw, err := os.ReadFile(source)
	if err != nil {
		return IngestResult{}, err
	}
	inputName := uniqueInputName(state, filepath.Base(source))
	dest := filepath.Join(runstate.RunDir(opts.Repo, state.RunID), "inputs", inputName)
	if err := os.WriteFile(dest, raw, 0o644); err != nil {
		return IngestResult{}, err
	}
	sum, err := runstate.SHA256File(dest)
	if err != nil {
		return IngestResult{}, err
	}
	ref := runstate.InputRef{
		ID:       fmt.Sprintf("input-%03d", len(state.Inputs)+1),
		Type:     opts.Type,
		Path:     relToRun(opts.Repo, state.RunID, dest),
		SHA256:   sum,
		AddedAt:  time.Now().UTC().Format(time.RFC3339),
		Original: source,
	}
	state.Inputs = append(state.Inputs, ref)
	if err := writeManifest(opts.Repo, state); err != nil {
		return IngestResult{}, err
	}
	if err := writeSummary(opts.Repo, state, raw); err != nil {
		return IngestResult{}, err
	}
	if opts.Slice {
		if err := writeSegments(opts.Repo, state, raw); err != nil {
			return IngestResult{}, err
		}
	}
	if opts.Type == "relay_transcript" {
		if err := writeRelayPoints(opts.Repo, state, raw); err != nil {
			return IngestResult{}, err
		}
	}
	state.Status = runstate.StatusIngested
	if err := runstate.Save(opts.Repo, state); err != nil {
		return IngestResult{}, err
	}
	return IngestResult{RunID: state.RunID, Status: state.Status, Input: ref, Artifacts: state.Artifacts}, nil
}

func writeManifest(repo string, state *runstate.RunState) error {
	path := filepath.Join(runstate.RunDir(repo, state.RunID), "inputs", "source_manifest.json")
	raw, err := json.MarshalIndent(state.Inputs, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return err
	}
	runstate.AddArtifact(state, "source_manifest", relToRun(repo, state.RunID, path))
	return nil
}

func writeSummary(repo string, state *runstate.RunState, raw []byte) error {
	text := string(raw)
	if len(text) > 1200 {
		text = text[:1200] + "\n\n[truncated]"
	}
	path := filepath.Join(runstate.RunDir(repo, state.RunID), "inputs", "input_summary.md")
	content := fmt.Sprintf("# Input Summary\n\nRun: `%s`\n\nInputs: %d\n\n## Latest excerpt\n\n%s\n", state.RunID, len(state.Inputs), text)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return err
	}
	runstate.AddArtifact(state, "input_summary", relToRun(repo, state.RunID, path))
	return nil
}

func writeSegments(repo string, state *runstate.RunState, raw []byte) error {
	parts := splitSegments(string(raw))
	path := filepath.Join(runstate.RunDir(repo, state.RunID), "inputs", "conversation_segments.json")
	rawJSON, err := json.MarshalIndent(parts, "", "  ")
	if err != nil {
		return err
	}
	rawJSON = append(rawJSON, '\n')
	if err := os.WriteFile(path, rawJSON, 0o644); err != nil {
		return err
	}
	runstate.AddArtifact(state, "conversation_segments", relToRun(repo, state.RunID, path))
	return nil
}

func writeRelayPoints(repo string, state *runstate.RunState, raw []byte) error {
	points := map[string][]string{"settled": {}, "contested": {}, "withdrawn": {}}
	for _, line := range strings.Split(string(raw), "\n") {
		lower := strings.ToLower(line)
		for key := range points {
			if strings.Contains(lower, key) {
				points[key] = append(points[key], strings.TrimSpace(line))
			}
		}
	}
	path := filepath.Join(runstate.RunDir(repo, state.RunID), "inputs", "relay_points.json")
	rawJSON, err := json.MarshalIndent(points, "", "  ")
	if err != nil {
		return err
	}
	rawJSON = append(rawJSON, '\n')
	if err := os.WriteFile(path, rawJSON, 0o644); err != nil {
		return err
	}
	runstate.AddArtifact(state, "relay_points", relToRun(repo, state.RunID, path))
	return nil
}

func splitSegments(text string) []map[string]string {
	chunks := strings.Split(text, "\n\n")
	var segments []map[string]string
	for i, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}
		segments = append(segments, map[string]string{
			"id":   fmt.Sprintf("segment-%03d", i+1),
			"text": chunk,
		})
	}
	if len(segments) == 0 {
		segments = append(segments, map[string]string{"id": "segment-001", "text": strings.TrimSpace(text)})
	}
	return segments
}

func uniqueInputName(state *runstate.RunState, base string) string {
	base = strings.ReplaceAll(base, string(filepath.Separator), "-")
	return fmt.Sprintf("%03d-%s", len(state.Inputs)+1, base)
}

func relToRun(repo, runID, path string) string {
	rel, err := filepath.Rel(runstate.RunDir(repo, runID), path)
	if err != nil {
		return path
	}
	return filepath.ToSlash(rel)
}
