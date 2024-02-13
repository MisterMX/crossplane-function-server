package server

import (
	"fmt"

	"github.com/pkg/errors"
)

type errNotFound struct {
	name string
}

func NewErrorNotFound(name string) error {
	return &errNotFound{name: name}
}

func (e errNotFound) Error() string {
	return fmt.Sprintf("not found: %s", e.name)
}

func IsErrorNotFound(err error) bool {
	return errors.As(err, &errNotFound{})
}
