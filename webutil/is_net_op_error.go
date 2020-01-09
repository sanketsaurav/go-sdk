package webutil

import (
	"net"

	"github.com/blend/go-sdk/ex"
)

// IsNetOpError is an exception filter useful for filtering out
// network errors.
func IsNetOpError() ex.Predicate {
	return func(err error) bool {
		_, ok := err.(*net.OpError)
		return ok
	}
}
