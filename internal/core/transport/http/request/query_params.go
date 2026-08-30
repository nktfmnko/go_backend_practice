package core_http_request

import (
	"fmt"
	"net/http"
	core_error "practice/internal/core/errors"
	"strconv"
)

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	v, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param=%s by key=%s not a valid int: %v: %w",
			param,
			key,
			err,
			core_error.ErrInvalidArgument,
		)
	}

	return &v, nil
}
