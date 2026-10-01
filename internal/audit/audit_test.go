package audit

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAudit(t *testing.T) {
	path := filepath.Join(t.TempDir(), "audit.log")
	if err := New(path).WriteAudit("audit entry"); err != nil {
		t.Fatalf("WriteAudit() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != "audit entry\n" {
		t.Errorf("audit file = %q, want %q", got, "audit entry\n")
	}
}

func TestWriteAuditWrapsWriteError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "audit.log")
	err := New(path).WriteAudit("audit entry")
	if err == nil {
		t.Fatal("WriteAudit() error = nil, want write error")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("WriteAudit() error = %v, want wrapped os.ErrNotExist", err)
	}
}
