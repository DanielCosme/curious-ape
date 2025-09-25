package core

import "errors"

var (
	ErrRepositoryNotFound = errors.New("repository: not found")
)

func IfErrNNotFound(err error) bool {
	return err != nil && !errors.Is(err, ErrRepositoryNotFound)
}
