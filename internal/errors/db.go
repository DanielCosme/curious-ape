package errors

import "errors"

var (
	ErrRecordNotFound = errors.New("db record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

func Is(err, target error) bool {
	return errors.Is(err, target)
}
