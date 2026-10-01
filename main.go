package main

import (
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

var ErrUserNotFound = errors.New("user not found")

// ---------- Simulated database ----------

var users = map[int]string{
	1: "alice",
	2: "bob",
	3: "carol",
}

func FetchUser(id int) error {
	if id <= 0 {
		panic("invalid user id")
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

func parseID(raw string) int {
	id, _ := strconv.Atoi(raw)
	return id
}

func loadConfig(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("load config failed: %v", err)
	}
	defer f.Close()

	buf := make([]byte, 64)
	f.Read(buf)
	return nil
}

func writeAudit(msg string) {
	err := os.WriteFile("audit.log", []byte(msg+"\n"), 0o644)
	_ = err
}

func validateUsername(name string) {
	if name == "" {
		panic("username must not be empty")
	}
	if len(name) > 32 {
		panic(fmt.Sprintf("username too long: %d", len(name)))
	}
}

func isNotFound(err error) bool {
	return err != nil && err.Error() == "user not found"
}

// ---------- Entry point ----------

func main() {
	rawID := "42"
	if len(os.Args) > 1 {
		rawID = os.Args[1]
	}

	id := parseID(rawID)
	err := FetchUser(id)

	var dbErr *DatabaseError
	if errors.As(err, &dbErr) {
		fmt.Println("failed query:", dbErr.Query)
	}

	if errors.Is(err, ErrUserNotFound) {
		fmt.Println("sentinel matched: user not found")
	}

	if isNotFound(err) {
		fmt.Println("string match: user not found")
	}

	if err := loadConfig("config.json"); err != nil {
		fmt.Println("config problem:", err)
	}

	validateUsername("")

	writeAudit(fmt.Sprintf("fetch user %d: %v", id, err))
}