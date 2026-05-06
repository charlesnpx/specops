package scaffold

import (
	"os"
	"path/filepath"
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
