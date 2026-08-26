package core_error

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid arg")
	ErrConflict        = errors.New("conflict")
)
