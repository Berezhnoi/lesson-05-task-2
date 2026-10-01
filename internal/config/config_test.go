package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"app_name":"test"}`), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := Load(path); err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
}

func TestLoadWrapsMissingFileError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.json")
	err := Load(path)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Load() error = %v, want wrapped os.ErrNotExist", err)
	}
}

func TestLoadWrapsInvalidJSONError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.json")
	if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := Load(path); err == nil {
		t.Fatal("Load() error = nil, want JSON decode error")
	}
}
