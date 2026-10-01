package repository

import (
	"fmt"

	"error-handling-lab/internal/apperrors"
	"error-handling-lab/internal/models"
)

type Repository struct {
	users map[int]models.User
}

func New() *Repository {
	return &Repository{
		users: map[int]models.User{
			1: {ID: 1, Username: "alice"},
			2: {ID: 2, Username: "bob"},
			3: {ID: 3, Username: "carol"},
		},
	}
}

func (r *Repository) FetchUser(id int) error {
	if _, ok := r.users[id]; !ok {
		query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", id)
		return fmt.Errorf("fetch user %d: %w", id, &models.DatabaseError{
			Query:       query,
			OriginalErr: apperrors.ErrUserNotFound,
		})
	}
	return nil
}
