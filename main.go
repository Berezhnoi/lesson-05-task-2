package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

// ---------- Custom error type ----------

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

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidUserID = errors.New("invalid user id")
)

// ---------- Simulated database ----------

var users = map[int]string{
	1: "alice",
	2: "bob",
	3: "carol",
}

func FetchUser(id int) error {
	if id <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidUserID, id)
	}

	if _, ok := users[id]; !ok {
		query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", id)
		return fmt.Errorf("fetch user %d: %w", id, &DatabaseError{
			Query:       query,
			OriginalErr: ErrUserNotFound,
		})
	}
	return nil
}

// ---------- Helpers ----------

func parseID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("parse user id %q: %w", raw, err)
	}
	return id, nil
}

func loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config %q: %w", path, err)
	}

	var config map[string]any
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("decode config %q: %w", path, err)
	}
	return nil
}

func writeAudit(msg string) error {
	if err := os.WriteFile("audit.log", []byte(msg+"\n"), 0o644); err != nil {
		return fmt.Errorf("write audit log: %w", err)
	}
	return nil
}

func validateUsername(name string) error {
	if name == "" {
		return errors.New("username must not be empty")
	}
	if len(name) > 32 {
		return fmt.Errorf("username too long: %d characters (maximum 32)", len(name))
	}
	return nil
}

// ---------- Entry point ----------

func main() {
	rawID := "42"
	if len(os.Args) > 1 {
		rawID = os.Args[1]
	}

	id, err := parseID(rawID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	fetchErr := FetchUser(id)

	var dbErr *DatabaseError
	if errors.As(fetchErr, &dbErr) {
		fmt.Println("failed query:", dbErr.Query)
	}

	if errors.Is(fetchErr, ErrUserNotFound) {
		fmt.Println("sentinel matched: user not found")
	}

	if fetchErr != nil && !errors.Is(fetchErr, ErrUserNotFound) {
		fmt.Fprintln(os.Stderr, "error:", fmt.Errorf("fetch user %d: %w", id, fetchErr))
		os.Exit(1)
	}

	if err := loadConfig("config.json"); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if err := validateUsername("alice"); err != nil {
		fmt.Fprintln(os.Stderr, "error:", fmt.Errorf("validate username: %w", err))
		os.Exit(1)
	}

	auditMessage := fmt.Sprintf("fetch user %d: ok", id)
	if fetchErr != nil {
		auditMessage = fetchErr.Error()
	}
	if err := writeAudit(auditMessage); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
