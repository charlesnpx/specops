package eval

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Report struct {
	Schema    int                `json:"schema"`
	EvalID    string             `json:"eval_id"`
	Gold      string             `json:"gold"`
	Candidate string             `json:"candidate"`
	Scores    map[string]float64 `json:"scores"`
	Findings  []string           `json:"findings"`
}

type DiffReport struct {
	GoldOnly      []string `json:"gold_only"`
	CandidateOnly []string `json:"candidate_only"`
	Common        []string `json:"common"`
}

func Run(gold, candidate string) (Report, error) {
	goldAbs, err := filepath.Abs(gold)
	if err != nil {
		return Report{}, err
	}
	candidateAbs, err := filepath.Abs(candidate)
	if err != nil {
		return Report{}, err
	}
	diff, err := Diff(goldAbs, candidateAbs)
	if err != nil {
		return Report{}, err
	}
	required := []string{"docs/CANON.md", "docs/versions/v0_scope.md", "docs/interfaces/cli_commands.md", ".specops/process.md"}
	found := 0
	var findings []string
	for _, rel := range required {
		if _, err := os.Stat(filepath.Join(candidateAbs, rel)); err == nil {
			found++
		} else {
			findings = append(findings, "candidate missing "+rel)
		}
	}
	structure := float64(found) / float64(len(required))
	coverage := 1.0
	if len(diff.GoldOnly) > 0 {
		coverage = float64(len(diff.Common)) / float64(len(diff.Common)+len(diff.GoldOnly))
	}
	return Report{
		Schema:    1,
		EvalID:    "eval-" + time.Now().UTC().Format("20060102-150405"),
		Gold:      goldAbs,
		Candidate: candidateAbs,
		Scores: map[string]float64{
			"structure": structure,
			"coverage":  coverage,
		},
		Findings: findings,
	}, nil
}

func Diff(gold, candidate string) (DiffReport, error) {
	goldFiles, err := fileSet(gold)
	if err != nil {
		return DiffReport{}, err
	}
	candidateFiles, err := fileSet(candidate)
	if err != nil {
		return DiffReport{}, err
	}
	report := DiffReport{GoldOnly: []string{}, CandidateOnly: []string{}, Common: []string{}}
	for file := range goldFiles {
		if candidateFiles[file] {
			report.Common = append(report.Common, file)
		} else {
			report.GoldOnly = append(report.GoldOnly, file)
		}
	}
	for file := range candidateFiles {
		if !goldFiles[file] {
			report.CandidateOnly = append(report.CandidateOnly, file)
		}
	}
	sort.Strings(report.GoldOnly)
	sort.Strings(report.CandidateOnly)
	sort.Strings(report.Common)
	return report, nil
}

func Score(path string) (map[string]float64, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var report Report
	if err := json.Unmarshal(raw, &report); err != nil {
		return nil, err
	}
	total := 0.0
	for _, value := range report.Scores {
		total += value
	}
	average := 0.0
	if len(report.Scores) > 0 {
		average = total / float64(len(report.Scores))
	}
	report.Scores["overall"] = average
	return report.Scores, nil
}

func fileSet(root string) (map[string]bool, error) {
	files := map[string]bool{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == ".specops" && path != filepath.Join(root, ".specops") {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = true
		return nil
	})
	return files, err
}
