package core_error

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidArgument    = errors.New("invalid arg")
	ErrConflict           = errors.New("conflict")
	ErrTooManyRequests    = errors.New("too many requests")
	ErrMissingBearerToken = errors.New("missing bearer token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrWrongApp           = errors.New("wrong app")
)
