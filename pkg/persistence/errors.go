package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"git.danicos.dev/daniel/curious-ape/pkg/core"
)

var (
	ErrNotFound           = sql.ErrNoRows
	ErrInvalidCredentials = errors.New("database: invalid credentials")
)

func catchDBErr(op string, err error) error {
	if err == nil {
		return nil
	}
	switch e := err.Error(); e {
	case sql.ErrNoRows.Error(): // Don't log this type of error.
		return fmt.Errorf("%w %s", core.ErrRepositoryNotFound, op)
	default:
		slog.Error(op + ": " + e)
	}
	return fmt.Errorf("%w %s", err, op)
}
