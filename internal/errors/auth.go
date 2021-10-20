package errors

import "errors"

var (
	ErrTokenExpired = errors.New("token expired")
	ErrUnauthorized = errors.New("server needs to authorize again")
	ErrNoRecord     = errors.New("provider has no record")
)
