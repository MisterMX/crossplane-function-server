package server

import (
	"fmt"

	"github.com/pkg/errors"
)

type errNotFound struct {
	name string
}

func (e errNotFound) Error() string {
	return fmt.Sprintf("not found: %s", e.name)
}

func IsNotFoundError(err error) bool {
	return errors.As(err, &errNotFound{})
}
