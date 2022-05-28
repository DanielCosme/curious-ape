package errors

import "runtime/debug"

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
