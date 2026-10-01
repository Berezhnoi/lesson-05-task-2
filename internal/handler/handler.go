package handler

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"error-handling-lab/internal/apperrors"
	"error-handling-lab/internal/models"
)

type UserService interface {
	LookupUser(id int) error
	ValidateUsername(name string) error
}

type AuditWriter interface {
	WriteAudit(msg string) error
}

type Handler struct {
	users  UserService
	audit  AuditWriter
	output io.Writer
}

func New(users UserService, audit AuditWriter, output io.Writer) *Handler {
	return &Handler{users: users, audit: audit, output: output}
}

func (h *Handler) Handle(args []string) error {
	rawID := "42"
	if len(args) > 0 {
		rawID = args[0]
	}

	id, err := parseID(rawID)
	if err != nil {
		return err
	}

	fetchErr := h.users.LookupUser(id)
	var dbErr *models.DatabaseError
	if errors.As(fetchErr, &dbErr) {
		if err := h.writeLine("failed query:", dbErr.Query); err != nil {
			return err
		}
	}

	if errors.Is(fetchErr, apperrors.ErrUserNotFound) {
		if err := h.writeLine("sentinel matched: user not found"); err != nil {
			return err
		}
	} else if fetchErr != nil {
		return fetchErr
	}

	if err := h.users.ValidateUsername("alice"); err != nil {
		return fmt.Errorf("validate username: %w", err)
	}

	auditMessage := fmt.Sprintf("fetch user %d: ok", id)
	if fetchErr != nil {
		auditMessage = fetchErr.Error()
	}
	if err := h.audit.WriteAudit(auditMessage); err != nil {
		return err
	}
	return nil
}

func (h *Handler) writeLine(values ...any) error {
	if _, err := fmt.Fprintln(h.output, values...); err != nil {
		return fmt.Errorf("write handler output: %w", err)
	}
	return nil
}

func parseID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("parse user id %q: %w", raw, err)
	}
	return id, nil
}
