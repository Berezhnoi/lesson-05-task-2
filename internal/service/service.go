package service

import (
	"errors"
	"fmt"

	"error-handling-lab/internal/apperrors"
)

type UserRepository interface {
	FetchUser(id int) error
}

type Service struct {
	repository UserRepository
}

func New(repository UserRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) LookupUser(id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: %d", apperrors.ErrInvalidUserID, id)
	}
	return s.repository.FetchUser(id)
}

func (s *Service) ValidateUsername(name string) error {
	if name == "" {
		return errors.New("username must not be empty")
	}
	if len(name) > 32 {
		return fmt.Errorf("username too long: %d characters (maximum 32)", len(name))
	}
	return nil
}
