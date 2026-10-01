package audit

import (
	"fmt"
	"os"
)

type Writer struct {
	path string
}

func New(path string) *Writer {
	return &Writer{path: path}
}

func (w *Writer) WriteAudit(msg string) error {
	return writeAudit(w.path, msg)
}

func writeAudit(path, msg string) error {
	if err := os.WriteFile(path, []byte(msg+"\n"), 0o644); err != nil {
		return fmt.Errorf("write audit log: %w", err)
	}
	return nil
}
