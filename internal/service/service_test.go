package service

import (
	"errors"
	"testing"

	"error-handling-lab/internal/apperrors"
)

type stubRepository struct {
	err   error
	id    int
	calls int
}

func (r *stubRepository) FetchUser(id int) error {
	r.calls++
	r.id = id
	return r.err
}

func TestLookupUserRejectsInvalidID(t *testing.T) {
	repository := &stubRepository{}
	s := New(repository)

	err := s.LookupUser(0)
	if !errors.Is(err, apperrors.ErrInvalidUserID) {
		t.Errorf("LookupUser() error = %v, want ErrInvalidUserID", err)
	}
	if repository.calls != 0 {
		t.Errorf("FetchUser calls = %d, want 0", repository.calls)
	}
}

func TestLookupUserDelegatesValidID(t *testing.T) {
	repository := &stubRepository{}
	s := New(repository)

	if err := s.LookupUser(12); err != nil {
		t.Fatalf("LookupUser() error = %v, want nil", err)
	}
	if repository.calls != 1 || repository.id != 12 {
		t.Errorf("FetchUser called with id %d (%d calls), want id 12 once", repository.id, repository.calls)
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "alice"},
		{name: "", wantErr: true},
		{name: "123456789012345678901234567890123", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(&stubRepository{}).ValidateUsername(tt.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUsername() error = %v, wantErr %t", err, tt.wantErr)
			}
		})
	}
}
