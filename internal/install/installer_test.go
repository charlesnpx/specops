package install

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestPlanUsesAbsoluteStagedPaths(t *testing.T) {
	root := t.TempDir()
	report, err := Execute(Options{Operation: OperationPlan, Target: TargetAll, InstallRoot: root, Version: "0.1.0"})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := json.Marshal(report)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(raw) {
		t.Fatal("report did not marshal as valid JSON")
	}
	for name, target := range report.Targets {
		if len(target.Files) == 0 {
			t.Fatalf("target %s has no files", name)
		}
		for _, file := range target.Files {
			if !filepath.IsAbs(file.Path) {
				t.Fatalf("path is not absolute: %s", file.Path)
			}
			if file.SHA256 != "" {
				t.Fatalf("plan should not include sha256 for %s", file.Path)
			}
		}
	}
}

func TestInstallWritesInsideRootAndHashes(t *testing.T) {
	root := t.TempDir()
	report, err := Execute(Options{Operation: OperationInstall, Target: TargetAll, InstallRoot: root, Version: "0.1.0"})
	if err != nil {
		t.Fatal(err)
	}
	for name, target := range report.Targets {
		for _, file := range target.Files {
			if !insideRoot(root, file.Path) {
				t.Fatalf("target %s escaped root: %s", name, file.Path)
			}
			if file.SHA256 == "" {
				t.Fatalf("missing sha256 for %s", file.Path)
			}
			if _, err := os.Stat(file.Path); err != nil {
				t.Fatalf("installed file missing: %s: %v", file.Path, err)
			}
		}
	}
}
