package repository

import (
	"errors"
	"testing"

	"error-handling-lab/internal/apperrors"
	"error-handling-lab/internal/models"
)

func TestFetchUser(t *testing.T) {
	repository := New()

	if err := repository.FetchUser(1); err != nil {
		t.Fatalf("FetchUser(1) error = %v, want nil", err)
	}

	err := repository.FetchUser(42)
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		t.Errorf("FetchUser(42) error = %v, want ErrUserNotFound", err)
	}
	var databaseErr *models.DatabaseError
	if !errors.As(err, &databaseErr) {
		t.Fatalf("errors.As(%v, *DatabaseError) = false", err)
	}
	if databaseErr.Query != "SELECT * FROM users WHERE id = 42" {
		t.Errorf("DatabaseError.Query = %q, want query for ID 42", databaseErr.Query)
	}
}
