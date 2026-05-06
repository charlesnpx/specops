package install

import (
	"fmt"
	"os"
	"path/filepath"
)

func BuildReport(opts Options) (Report, error) {
	if opts.Version == "" {
		opts.Version = "0.1.0"
	}
	files, err := plannedFiles(opts)
	if err != nil {
		return Report{}, err
	}
	report := Report{
		Schema:    1,
		Name:      "specops",
		Version:   opts.Version,
		Operation: opts.Operation,
		Kind:      "delegated",
		Targets:   map[string]TargetReport{},
		Warnings:  []string{},
	}
	for _, target := range expandTarget(opts.Target) {
		report.Targets[string(target)] = TargetReport{Files: []FileReport{}}
	}
	for _, file := range files {
		targetReport := report.Targets[string(file.Target)]
		targetReport.Files = append(targetReport.Files, FileReport{Path: file.Path})
		report.Targets[string(file.Target)] = targetReport
	}
	return report, nil
}

func plannedFiles(opts Options) ([]filePlan, error) {
	root, err := installBase(opts.InstallRoot)
	if err != nil {
		return nil, err
	}
	var files []filePlan
	for _, target := range expandTarget(opts.Target) {
		switch target {
		case TargetClaude:
			files = append(files, filePlan{
				Target:  TargetClaude,
				Path:    filepath.Join(root, ".claude", "skills", "specops", "SKILL.md"),
				Content: []byte(claudeSkill),
			})
		case TargetCodex:
			files = append(files, filePlan{
				Target:  TargetCodex,
				Path:    filepath.Join(root, ".codex", "skills", "specops", "SKILL.md"),
				Content: []byte(codexSkill),
			})
		case TargetTools:
			files = append(files, filePlan{
				Target:     TargetTools,
				Path:       filepath.Join(root, ".local", "bin", "specops"),
				Executable: true,
				ToolBinary: true,
			})
		default:
			return nil, fmt.Errorf("unsupported target %q", target)
		}
	}
	return files, nil
}

func installBase(root string) (string, error) {
	if root != "" {
		if !filepath.IsAbs(root) {
			return "", fmt.Errorf("--install-root must be absolute")
		}
		return filepath.Clean(root), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

func expandTarget(target Target) []Target {
	switch target {
	case TargetAll, "":
		return []Target{TargetClaude, TargetCodex, TargetTools}
	case TargetClaude:
		return []Target{TargetClaude}
	case TargetCodex:
		return []Target{TargetCodex}
	case TargetTools:
		return []Target{TargetTools}
	default:
		return []Target{target}
	}
}
