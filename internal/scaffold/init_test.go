package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitIsIdempotent(t *testing.T) {
	root := t.TempDir()
	first, err := Init(Options{Path: root, Mode: "minimal", Agent: "both"})
	if err != nil {
		t.Fatal(err)
	}
	if len(first.Files) == 0 {
		t.Fatal("expected scaffold files")
	}
	second, err := Init(Options{Path: root, Mode: "minimal", Agent: "both"})
	if err != nil {
		t.Fatal(err)
	}
	for _, file := range second.Files {
		if file.Status != "unchanged" {
			t.Fatalf("expected idempotent unchanged status for %s, got %s", file.Path, file.Status)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".specops.lock")); err != nil {
		t.Fatal(err)
	}
}

func TestSpecDeltaTemplateExposesPatchItems(t *testing.T) {
	root := t.TempDir()
	if _, err := Init(Options{Path: root, Mode: "minimal", Agent: "both"}); err != nil {
		t.Fatal(err)
	}

	raw, err := os.ReadFile(filepath.Join(root, ".specops", "templates", "spec_delta.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	template := string(raw)
	for _, want := range []string{
		"patch_plan: []",
		"Human-readable compile notes only",
		"patch_items: []",
		"content: |",
	} {
		if !strings.Contains(template, want) {
			t.Fatalf("expected spec delta template to contain %q", want)
		}
	}
}
