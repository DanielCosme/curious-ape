package database

import (
	"fmt"

	"github.com/danielcosme/go-sdk/errors"
)

var (
	ErrDatabase           = errors.New("database")
	ErrNotFound           = fmt.Errorf("%w: not found", ErrDatabase)
	ErrUniqueCheckFailed  = fmt.Errorf("%w: duplicate field", ErrDatabase)
	ErrInvalidCredentials = fmt.Errorf("%w: invalid credentials", ErrDatabase)
)
