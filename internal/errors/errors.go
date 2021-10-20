package errors

import "errors"

func New(s string) error {
	return errors.New(s)
}
