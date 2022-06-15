package database

import "github.com/danielcosme/curious-ape/sdk/errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrUniqueCheckFailed = errors.New("unique check failed")
)
