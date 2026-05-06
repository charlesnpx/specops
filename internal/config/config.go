package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Store struct {
	Path string
	Data map[string]string
}

func DefaultPath(explicit string) (string, error) {
	if explicit != "" {
		return filepath.Abs(explicit)
	}
	if home := os.Getenv("SPECOPS_HOME"); home != "" {
		return filepath.Join(home, "config.yaml"), nil
	}
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "specops", "config.yaml"), nil
}

func Load(explicit string) (*Store, error) {
	path, err := DefaultPath(explicit)
	if err != nil {
		return nil, err
	}
	store := &Store{Path: path, Data: map[string]string{}}
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return store, nil
		}
		return nil, err
	}
	if len(raw) == 0 {
		return store, nil
	}
	if err := json.Unmarshal(raw, &store.Data); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Store) Save() error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(s.Data, "", "  ")
	if err != nil {
		return err
	}
	raw = append(raw, '\n')
	return os.WriteFile(s.Path, raw, 0o644)
}
