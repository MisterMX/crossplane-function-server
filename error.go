package server

import "fmt"

type errNotFound string

func (e errNotFound) Error() string {
	return fmt.Sprintf("not found: %s", string(e))
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*errNotFound)
	return ok
}
