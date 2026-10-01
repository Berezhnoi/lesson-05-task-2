package models

import "fmt"

type User struct {
	ID       int
	Username string
}

type DatabaseError struct {
	Query       string
	OriginalErr error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error: query %q failed: %v", e.Query, e.OriginalErr)
}

func (e *DatabaseError) Unwrap() error {
	return e.OriginalErr
}
