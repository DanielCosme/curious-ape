package errors

import (
	"errors"
	"runtime/debug"
)

type Error struct {
	err   string
	Stack []byte
}

func New(err string) *Error {
	return &Error{err: err}
}

func NewFatal(err string) *Error {
	return &Error{err: err, Stack: debug.Stack()}
}

func (e *Error) Error() string {
	return e.err
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, &target)
}
