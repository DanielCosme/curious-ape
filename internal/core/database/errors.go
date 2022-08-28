package database

import "github.com/danielcosme/go-sdk/errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrUniqueCheckFailed = errors.New("unique check failed")
)
