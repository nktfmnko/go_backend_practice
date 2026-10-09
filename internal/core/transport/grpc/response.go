package core_transport_grpc

import (
	"errors"
	core_error "practice/internal/core/errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorResponse(err error, msg string) error {
	var statusCode codes.Code

	switch {
	case errors.Is(err, core_error.ErrInvalidArgument):
		statusCode = codes.InvalidArgument
	case errors.Is(err, core_error.ErrNotFound):
		statusCode = codes.NotFound
	case errors.Is(err, core_error.ErrConflict):
		statusCode = codes.AlreadyExists
	default:
		statusCode = codes.Internal
	}

	return errorResponse(statusCode, err, msg)
}

func errorResponse(statusCode codes.Code, err error, msg string) error {
	return status.Errorf(statusCode, "%s :%v", msg, err)
}
