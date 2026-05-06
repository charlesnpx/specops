package runstate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

func RunsDir(repo string) string {
	return filepath.Join(repo, ".specops", "runs")
}

func RunDir(repo, runID string) string {
	return filepath.Join(RunsDir(repo), runID)
}

func StatePath(repo, runID string) string {
	return filepath.Join(RunDir(repo, runID), "run.yaml")
}

func NewRun(repo, name string) (*RunState, error) {
	now := time.Now()
	runID := fmt.Sprintf("run-%s-%s", now.Format("20060102-150405"), slugify(name))
	if strings.HasSuffix(runID, "-") {
		runID += "untitled"
	}
	dir := RunDir(repo, runID)
	for _, sub := range []string{"inputs", "prompts", "outputs", "traces", "evals", "patches"} {
		if err := os.MkdirAll(filepath.Join(dir, sub), 0o755); err != nil {
			return nil, err
		}
	}
	state := &RunState{
		Schema:    1,
		RunID:     runID,
		Name:      name,
		Status:    StatusCreated,
		CreatedAt: now.UTC().Format(time.RFC3339),
		UpdatedAt: now.UTC().Format(time.RFC3339),
		Inputs:    []InputRef{},
		Artifacts: []ArtifactRef{},
		Decisions: map[string]Decision{},
	}
	state.Next = NextForStatus(state.Status, state.RunID)
	if err := Save(repo, state); err != nil {
		return nil, err
	}
	return state, nil
}

func Load(repo, runID string) (*RunState, error) {
	resolved, err := ResolveRunID(repo, runID)
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(StatePath(repo, resolved))
	if err != nil {
		return nil, err
	}
	var state RunState
	if err := json.Unmarshal(raw, &state); err != nil {
		return nil, err
	}
	if state.Decisions == nil {
		state.Decisions = map[string]Decision{}
	}
	state.Next = NextForStatus(state.Status, state.RunID)
	return &state, nil
}

func Save(repo string, state *RunState) error {
	state.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	state.Next = NextForStatus(state.Status, state.RunID)
	raw, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	if err := os.MkdirAll(RunDir(repo, state.RunID), 0o755); err != nil {
		return err
	}
	return os.WriteFile(StatePath(repo, state.RunID), raw, 0o644)
}

func List(repo string) ([]*RunState, error) {
	entries, err := os.ReadDir(RunsDir(repo))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []*RunState{}, nil
		}
		return nil, err
	}
	var states []*RunState
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		state, err := Load(repo, entry.Name())
		if err != nil {
			continue
		}
		states = append(states, state)
	}
	sort.Slice(states, func(i, j int) bool {
		return states[i].CreatedAt < states[j].CreatedAt
	})
	return states, nil
}

func ResolveRunID(repo, query string) (string, error) {
	if query == "" {
		return "", fmt.Errorf("run id is required")
	}
	if _, err := os.Stat(StatePath(repo, query)); err == nil {
		return query, nil
	}
	entries, err := os.ReadDir(RunsDir(repo))
	if err != nil {
		return "", err
	}
	var matches []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), query) {
			matches = append(matches, entry.Name())
		}
	}
	switch len(matches) {
	case 0:
		return "", fmt.Errorf("run %q not found", query)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("run id %q is ambiguous", query)
	}
}

func AddArtifact(state *RunState, typ, path string) {
	state.Artifacts = append(state.Artifacts, ArtifactRef{
		ID:        fmt.Sprintf("%s-%03d", typ, len(state.Artifacts)+1),
		Type:      typ,
		Path:      path,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

func SHA256File(path string) (string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(f)
	return hex.EncodeToString(sum[:]), nil
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = slugRe.ReplaceAllString(value, "-")
	return strings.Trim(value, "-")
}
