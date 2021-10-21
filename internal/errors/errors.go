package errors

import "errors"

func New(s string) error {
	return errors.New(s)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}
