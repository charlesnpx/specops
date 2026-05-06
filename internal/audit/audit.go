package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Report struct {
	Schema   int      `json:"schema"`
	Repo     string   `json:"repo"`
	Passed   bool     `json:"passed"`
	Checks   []Check  `json:"checks"`
	Warnings []string `json:"warnings"`
}

type Check struct {
	Name    string   `json:"name"`
	Passed  bool     `json:"passed"`
	Details []string `json:"details,omitempty"`
}

func Run(repo string) (Report, error) {
	root, err := filepath.Abs(repo)
	if err != nil {
		return Report{}, err
	}
	report := Report{Schema: 1, Repo: root, Passed: true, Checks: []Check{}, Warnings: []string{}}
	add := func(check Check) {
		report.Checks = append(report.Checks, check)
		if !check.Passed {
			report.Passed = false
		}
	}
	add(requiredDirs(root))
	add(versionScope(root))
	add(acceptedADRs(root))
	add(schemaFiles(root))
	add(runStates(root))
	front := frontMatter(root)
	report.Warnings = append(report.Warnings, front.Details...)
	links := localLinks(root)
	if !links.Passed {
		report.Warnings = append(report.Warnings, links.Details...)
	}
	return report, nil
}

func schemaFiles(root string) Check {
	var failures []string
	for _, dir := range []string{filepath.Join(root, "schemas"), filepath.Join(root, ".specops", "schemas")} {
		entries, _ := filepath.Glob(filepath.Join(dir, "*.schema.json"))
		for _, path := range entries {
			raw, err := os.ReadFile(path)
			if err != nil {
				failures = append(failures, filepath.ToSlash(path)+": unreadable")
				continue
			}
			var value any
			if err := json.Unmarshal(raw, &value); err != nil {
				rel, _ := filepath.Rel(root, path)
				failures = append(failures, filepath.ToSlash(rel)+": invalid JSON")
			}
		}
	}
	return Check{Name: "schema files parse", Passed: len(failures) == 0, Details: failures}
}

func runStates(root string) Check {
	var failures []string
	runs := filepath.Join(root, ".specops", "runs")
	_ = filepath.WalkDir(runs, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Base(path) != "run.yaml" {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			failures = append(failures, filepath.ToSlash(path)+": unreadable")
			return nil
		}
		var state struct {
			Schema    int    `json:"schema"`
			RunID     string `json:"run_id"`
			Status    string `json:"status"`
			CreatedAt string `json:"created_at"`
		}
		if err := json.Unmarshal(raw, &state); err != nil {
			rel, _ := filepath.Rel(root, path)
			failures = append(failures, filepath.ToSlash(rel)+": invalid run state JSON")
			return nil
		}
		if state.Schema != 1 || state.RunID == "" || state.Status == "" || state.CreatedAt == "" {
			rel, _ := filepath.Rel(root, path)
			failures = append(failures, filepath.ToSlash(rel)+": missing required run state fields")
		}
		return nil
	})
	return Check{Name: "run state artifacts parse", Passed: len(failures) == 0, Details: failures}
}

func requiredDirs(root string) Check {
	required := []string{
		".specops",
		"docs",
		"docs/decisions",
		"docs/interfaces",
		"docs/versions",
		"docs/research/refinery",
	}
	var missing []string
	for _, rel := range required {
		if info, err := os.Stat(filepath.Join(root, rel)); err != nil || !info.IsDir() {
			missing = append(missing, rel)
		}
	}
	return Check{Name: "required directories", Passed: len(missing) == 0, Details: missing}
}

func versionScope(root string) Check {
	path := filepath.Join(root, "docs", "versions", "v0_scope.md")
	raw, err := os.ReadFile(path)
	if err != nil {
		return Check{Name: "v0 scope", Passed: false, Details: []string{"docs/versions/v0_scope.md missing"}}
	}
	if !strings.Contains(string(raw), "v0 required") {
		return Check{Name: "v0 scope", Passed: false, Details: []string{"v0 required section missing"}}
	}
	return Check{Name: "v0 scope", Passed: true}
}

func acceptedADRs(root string) Check {
	dir := filepath.Join(root, "docs", "decisions")
	entries, err := filepath.Glob(filepath.Join(dir, "ADR-*.md"))
	if err != nil {
		return Check{Name: "accepted ADR affected docs", Passed: false, Details: []string{err.Error()}}
	}
	var failures []string
	for _, path := range entries {
		raw, err := os.ReadFile(path)
		if err != nil {
			failures = append(failures, filepath.Base(path)+": unreadable")
			continue
		}
		text := string(raw)
		if strings.Contains(text, "status: accepted") && !strings.Contains(text, "## Affected docs") {
			failures = append(failures, filepath.Base(path)+": missing affected docs section")
		}
	}
	return Check{Name: "accepted ADR affected docs", Passed: len(failures) == 0, Details: failures}
}

func frontMatter(root string) Check {
	var missing []string
	_ = filepath.WalkDir(filepath.Join(root, "docs"), func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		if filepath.Base(path) == "README.md" {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		if !strings.HasPrefix(string(raw), "---\n") {
			rel, _ := filepath.Rel(root, path)
			missing = append(missing, filepath.ToSlash(rel)+": missing front matter")
		}
		return nil
	})
	return Check{Name: "front matter warnings", Passed: len(missing) == 0, Details: missing}
}

var mdLinkRE = regexp.MustCompile(`\[[^\]]+\]\(([^)]+\.md(?:#[^)]+)?)\)`)

func localLinks(root string) Check {
	var broken []string
	_ = filepath.WalkDir(filepath.Join(root, "docs"), func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".md" {
			return nil
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		for _, match := range mdLinkRE.FindAllStringSubmatch(string(raw), -1) {
			target := strings.Split(match[1], "#")[0]
			if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
				continue
			}
			candidate := filepath.Clean(filepath.Join(filepath.Dir(path), target))
			if _, err := os.Stat(candidate); err != nil {
				rel, _ := filepath.Rel(root, path)
				broken = append(broken, filepath.ToSlash(rel)+" -> "+target)
			}
		}
		return nil
	})
	return Check{Name: "local markdown links", Passed: len(broken) == 0, Details: broken}
}
