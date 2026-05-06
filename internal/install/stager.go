package install

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Execute(opts Options) (Report, error) {
	switch opts.Operation {
	case OperationPlan:
		return BuildReport(opts)
	case OperationInstall:
		return installFiles(opts)
	case OperationUninstall:
		return uninstallFiles(opts)
	default:
		return Report{}, fmt.Errorf("unsupported operation %q", opts.Operation)
	}
}

func installFiles(opts Options) (Report, error) {
	report, err := BuildReport(opts)
	if err != nil {
		return Report{}, err
	}
	files, err := plannedFiles(opts)
	if err != nil {
		return Report{}, err
	}
	for _, file := range files {
		if opts.InstallRoot != "" && !insideRoot(opts.InstallRoot, file.Path) {
			return Report{}, fmt.Errorf("planned path escapes install root: %s", file.Path)
		}
		if err := os.MkdirAll(filepath.Dir(file.Path), 0o755); err != nil {
			return Report{}, err
		}
		if file.ToolBinary {
			if err := copyCurrentExecutable(file.Path); err != nil {
				return Report{}, err
			}
		} else {
			mode := os.FileMode(0o644)
			if file.Executable {
				mode = 0o755
			}
			if err := os.WriteFile(file.Path, file.Content, mode); err != nil {
				return Report{}, err
			}
		}
		sum, err := hashFile(file.Path)
		if err != nil {
			return Report{}, err
		}
		targetReport := report.Targets[string(file.Target)]
		for i := range targetReport.Files {
			if targetReport.Files[i].Path == file.Path {
				targetReport.Files[i].SHA256 = sum
			}
		}
		report.Targets[string(file.Target)] = targetReport
	}
	return report, nil
}

func uninstallFiles(opts Options) (Report, error) {
	report, err := BuildReport(opts)
	if err != nil {
		return Report{}, err
	}
	files, err := plannedFiles(opts)
	if err != nil {
		return Report{}, err
	}
	for _, file := range files {
		if opts.InstallRoot != "" && !insideRoot(opts.InstallRoot, file.Path) {
			return Report{}, fmt.Errorf("planned path escapes install root: %s", file.Path)
		}
		if err := os.Remove(file.Path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return Report{}, err
		}
	}
	return report, nil
}

func copyCurrentExecutable(dest string) error {
	src, err := os.Executable()
	if err != nil {
		return err
	}
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if destInfo, err := os.Stat(dest); err == nil && os.SameFile(srcInfo, destInfo) {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}
	return os.Chmod(dest, 0o755)
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func insideRoot(root, path string) bool {
	root = filepath.Clean(root)
	path = filepath.Clean(path)
	rel, err := filepath.Rel(root, path)
	return err == nil && rel != ".." && !startsWithDotDot(rel)
}

func startsWithDotDot(rel string) bool {
	return len(rel) >= 3 && rel[:3] == "../"
}
