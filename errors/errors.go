package errors

import "errors"

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")
