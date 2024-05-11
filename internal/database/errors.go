package database

import "errors"

var (
	ErrNotFound           = errors.New("database: not found")
	ErrInvalidCredentials = errors.New("database: invalid credentials")
)

// IfNotFoundErr returns true if the errors exists and is NOT of type ErrNotFound, used
// when we want to ignore that specific error and continue execution.
func IfNotFoundErr(err error) bool {
	return err != nil && !errors.Is(err, ErrNotFound)
}
