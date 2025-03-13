package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrNotFound           = sql.ErrNoRows
	ErrInvalidCredentials = errors.New("database: invalid credentials")
)

// IgnoreIfErrNotFound returns true if the errors exists and is NOT of type ErrNotFound, used
// when we want to ignore that specific error and continue execution.
func IgnoreIfErrNotFound(err error) bool {
	return err != nil && !errors.Is(err, ErrNotFound)
}

func catchDBErr(op string, err error) error {
	if err == nil {
		return nil
	}
	switch e := err.Error(); e {
	case sql.ErrNoRows.Error(): // Don't log this type of error.
	default:
		slog.Error(op + ": " + e)
	}
	return fmt.Errorf("%w %s", ErrNotFound, op)
}
