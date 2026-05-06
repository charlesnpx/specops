package input

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type FixtureResult struct {
	From  string   `json:"from"`
	Out   string   `json:"out"`
	Files []string `json:"files"`
}

func BuildFixture(from, out string) (FixtureResult, error) {
	fromAbs, err := filepath.Abs(from)
	if err != nil {
		return FixtureResult{}, err
	}
	outAbs, err := filepath.Abs(out)
	if err != nil {
		return FixtureResult{}, err
	}
	result := FixtureResult{From: fromAbs, Out: outAbs, Files: []string{}}
	if err := filepath.WalkDir(fromAbs, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == ".specops" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, err := filepath.Rel(fromAbs, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(outAbs, rel)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return err
		}
		if err := copyFile(path, dest); err != nil {
			return err
		}
		result.Files = append(result.Files, filepath.ToSlash(rel))
		return nil
	}); err != nil {
		return FixtureResult{}, err
	}
	raw, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return FixtureResult{}, err
	}
	raw = append(raw, '\n')
	if err := os.WriteFile(filepath.Join(outAbs, "artifact_inventory.json"), raw, 0o644); err != nil {
		return FixtureResult{}, err
	}
	return result, nil
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
