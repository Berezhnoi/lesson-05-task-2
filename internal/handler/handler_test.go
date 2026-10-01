package handler

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"error-handling-lab/internal/apperrors"
	"error-handling-lab/internal/models"
)

type stubUserService struct {
	lookupErr     error
	validationErr error
	lookedUpID    int
	lookupCalls   int
}

func (s *stubUserService) LookupUser(id int) error {
	s.lookupCalls++
	s.lookedUpID = id
	return s.lookupErr
}

func (s *stubUserService) ValidateUsername(string) error {
	return s.validationErr
}

type stubAuditWriter struct {
	err     error
	message string
	calls   int
}

func (w *stubAuditWriter) WriteAudit(message string) error {
	w.calls++
	w.message = message
	return w.err
}

func TestHandleUserNotFoundReportsQueryAndSentinel(t *testing.T) {
	service := &stubUserService{
		lookupErr: fmt.Errorf("lookup failed: %w", &models.DatabaseError{
			Query:       "SELECT * FROM users WHERE id = 42",
			OriginalErr: apperrors.ErrUserNotFound,
		}),
	}
	audit := &stubAuditWriter{}
	var output bytes.Buffer
	h := New(service, audit, &output)

	if err := h.Handle([]string{"42"}); err != nil {
		t.Fatalf("Handle() error = %v, want nil", err)
	}
	if got, want := output.String(), "failed query: SELECT * FROM users WHERE id = 42\nsentinel matched: user not found\n"; got != want {
		t.Fatalf("output = %q, want %q", got, want)
	}
	if service.lookedUpID != 42 {
		t.Errorf("looked up ID = %d, want 42", service.lookedUpID)
	}
	if audit.message != service.lookupErr.Error() || audit.calls != 1 {
		t.Errorf("audit = (%q, %d calls), want lookup error once", audit.message, audit.calls)
	}
}

func TestHandleSuccessfulLookupIsSilentAndAuditsSuccess(t *testing.T) {
	service := &stubUserService{}
	audit := &stubAuditWriter{}
	var output bytes.Buffer
	h := New(service, audit, &output)

	if err := h.Handle([]string{"1"}); err != nil {
		t.Fatalf("Handle() error = %v, want nil", err)
	}
	if output.Len() != 0 {
		t.Errorf("output = %q, want empty", output.String())
	}
	if audit.message != "fetch user 1: ok" || audit.calls != 1 {
		t.Errorf("audit = (%q, %d calls), want success message once", audit.message, audit.calls)
	}
}

func TestHandleInvalidIDWrapsParseErrorAndStops(t *testing.T) {
	service := &stubUserService{}
	audit := &stubAuditWriter{}
	h := New(service, audit, &bytes.Buffer{})

	err := h.Handle([]string{"abc"})
	if err == nil || !strings.Contains(err.Error(), `parse user id "abc"`) {
		t.Fatalf("Handle() error = %v, want parse context", err)
	}
	var numberErr *strconv.NumError
	if !errors.As(err, &numberErr) {
		t.Errorf("errors.As(%v, *strconv.NumError) = false", err)
	}
	if service.lookupCalls != 0 || audit.calls != 0 {
		t.Errorf("invalid input continued: lookups=%d audits=%d", service.lookupCalls, audit.calls)
	}
}

func TestHandlePropagatesAuditError(t *testing.T) {
	writeErr := errors.New("disk full")
	audit := &stubAuditWriter{err: writeErr}
	h := New(&stubUserService{}, audit, &bytes.Buffer{})

	err := h.Handle([]string{"1"})
	if !errors.Is(err, writeErr) {
		t.Errorf("Handle() error = %v, want wrapped audit error", err)
	}
}

func TestHandlePropagatesOutputError(t *testing.T) {
	writeErr := errors.New("closed output")
	service := &stubUserService{
		lookupErr: &models.DatabaseError{Query: "query", OriginalErr: apperrors.ErrUserNotFound},
	}
	h := New(service, &stubAuditWriter{}, failingWriter{err: writeErr})

	err := h.Handle([]string{"42"})
	if !errors.Is(err, writeErr) {
		t.Errorf("Handle() error = %v, want wrapped output error", err)
	}
}

type failingWriter struct {
	err error
}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, w.err
}
