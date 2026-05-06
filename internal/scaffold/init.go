package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Options struct {
	Path   string
	Mode   string
	Agent  string
	Force  bool
	Backup bool
}

type Result struct {
	Path     string       `json:"path"`
	Version  string       `json:"version"`
	Mode     string       `json:"mode"`
	Agent    string       `json:"agent"`
	Files    []FileResult `json:"files"`
	Warnings []string     `json:"warnings"`
}

type FileResult struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

func Init(opts Options) (Result, error) {
	if opts.Path == "" {
		opts.Path = "."
	}
	if opts.Mode == "" {
		opts.Mode = "minimal"
	}
	if opts.Agent == "" {
		opts.Agent = "both"
	}
	root, err := filepath.Abs(opts.Path)
	if err != nil {
		return Result{}, err
	}
	result := Result{Path: root, Version: Version, Mode: opts.Mode, Agent: opts.Agent, Files: []FileResult{}, Warnings: []string{}}
	for _, dir := range []string{
		".specops",
		".specops/runs",
		".specops/cache",
		".specops/commands",
		".specops/templates",
		".specops/schemas",
		".specops/evals",
		".specops/checklists",
		"docs",
		"docs/research/refinery",
	} {
		if err := os.MkdirAll(filepath.Join(root, dir), 0o755); err != nil {
			return Result{}, err
		}
	}
	for rel, content := range files(opts.Mode, opts.Agent) {
		if rel == "CLAUDE.md" && opts.Agent == "codex" {
			continue
		}
		status, warning, err := writeManaged(filepath.Join(root, rel), []byte(content), opts)
		if err != nil {
			return Result{}, err
		}
		result.Files = append(result.Files, FileResult{Path: filepath.Join(root, rel), Status: status})
		if warning != "" {
			result.Warnings = append(result.Warnings, warning)
		}
	}
	status, warning, err := ensureGitignore(root, opts)
	if err != nil {
		return Result{}, err
	}
	result.Files = append(result.Files, FileResult{Path: filepath.Join(root, ".gitignore"), Status: status})
	if warning != "" {
		result.Warnings = append(result.Warnings, warning)
	}
	return result, nil
}

func writeManaged(path string, content []byte, opts Options) (string, string, error) {
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) == string(content) {
			return "unchanged", "", nil
		}
		if opts.Backup {
			backup := fmt.Sprintf("%s.%s.bak", path, time.Now().UTC().Format("20060102T150405Z"))
			if err := os.WriteFile(backup, existing, 0o644); err != nil {
				return "", "", err
			}
		}
		if !opts.Force && !opts.Backup {
			return "skipped", fmt.Sprintf("skipped divergent file %s", path), nil
		}
		if err := os.WriteFile(path, content, 0o644); err != nil {
			return "", "", err
		}
		return "updated", "", nil
	} else if !os.IsNotExist(err) {
		return "", "", err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", "", err
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return "", "", err
	}
	return "created", "", nil
}

func ensureGitignore(root string, opts Options) (string, string, error) {
	path := filepath.Join(root, ".gitignore")
	raw, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return "", "", err
	}
	existed := err == nil
	content := string(raw)
	lines := []string{".specops/runs/", ".specops/cache/"}
	changed := false
	for _, line := range lines {
		if !containsLine(content, line) {
			if content != "" && !strings.HasSuffix(content, "\n") {
				content += "\n"
			}
			content += line + "\n"
			changed = true
		}
	}
	if !changed {
		return "unchanged", "", nil
	}
	_ = opts
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", "", err
	}
	if existed {
		return "updated", "", nil
	}
	return "created", "", nil
}

func containsLine(content, needle string) bool {
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == needle {
			return true
		}
	}
	return false
}
